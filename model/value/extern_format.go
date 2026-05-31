package value

import "reflect"

func FmtWrap(v Value) any {
	iface := v.Interface()
	if iface == nil {
		return nil
	}

	typeName := isGigStruct(iface)
	if typeName == "" {
		rv := reflect.ValueOf(iface)
		if shouldFormatSequence(rv) {
			return &gigSequenceFormatter{rv: rv}
		}
		if shouldFormatMap(rv) {
			return &gigMapFormatter{rv: rv}
		}
		return iface
	}

	return makeGigStructWrapper(v, iface, typeName)
}

func makeGigStructWrapper(v Value, iface any, typeName string) any {
	captured := v

	// Lazy resolvers — each is called at most once by the wrapper, and only
	// when the corresponding fmt verb triggers its interface.
	var (
		stringerFunc       func() string
		stringerResolved   bool
		errorerFunc        func() string
		errorerResolved    bool
		gostringerFunc     func() string
		gostringerResolved bool
	)
	lazyStringer := func() (func() string, bool) {
		if !stringerResolved {
			stringerFunc, _ = resolveStringer(captured)
			stringerResolved = true
		}
		return stringerFunc, stringerFunc != nil
	}
	lazyErrorer := func() (func() string, bool) {
		if !errorerResolved {
			errorerFunc, _ = resolveErrorer(captured)
			errorerResolved = true
		}
		return errorerFunc, errorerFunc != nil
	}
	lazyGoStringer := func() (func() string, bool) {
		if !gostringerResolved {
			gostringerFunc, _ = resolveGoStringer(captured)
			gostringerResolved = true
		}
		return gostringerFunc, gostringerFunc != nil
	}

	// Always return the wrapper for gig structs - it handles all fmt verbs correctly
	return &gigStructWrapper{
		iface:          iface,
		typeName:       typeName,
		lazyStringer:   lazyStringer,
		lazyErrorer:    lazyErrorer,
		lazyGoStringer: lazyGoStringer,
	}
}

func shouldFormatSequence(rv reflect.Value) bool {
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return false
	}
	if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
		return false
	}
	if typeContainsGigStruct(rv.Type().Elem()) {
		return true
	}
	for i := 0; i < rv.Len(); i++ {
		if shouldWrapReflectValue(rv.Index(i)) {
			return true
		}
	}
	return false
}

func shouldFormatMap(rv reflect.Value) bool {
	if !rv.IsValid() || rv.Kind() != reflect.Map {
		return false
	}
	if typeContainsGigStruct(rv.Type().Key()) || typeContainsGigStruct(rv.Type().Elem()) {
		return true
	}
	for _, key := range rv.MapKeys() {
		if shouldWrapReflectValue(key) || shouldWrapReflectValue(rv.MapIndex(key)) {
			return true
		}
	}
	return false
}

func shouldWrapReflectValue(rv reflect.Value) bool {
	if !rv.IsValid() {
		return false
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.CanInterface() && isGigStruct(rv.Interface()) != "" {
		return true
	}
	return typeContainsGigStruct(rv.Type())
}
