package value

import "reflect"

// resolveStringer attempts to resolve the String() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard because side-channel method invocation via
// ResolveCompiledMethod can fail on method bodies that depend on features
// only wired up by a full program VM (e.g., global state, extCallCache). On
// such failure we return (nil, false) so the caller falls back to default
// formatting instead of propagating a panic through fmt.
func resolveStringer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	// Try to call String() method via the global resolver registry
	result, found := callMethod(nil, "String", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "String", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// resolveErrorer attempts to resolve the Error() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard, see resolveStringer for details.
func resolveErrorer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	result, found := callMethod(nil, "Error", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "Error", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}

// resolveGoStringer attempts to resolve the GoString() method for a value.
// Returns a function that can be called later, and a boolean indicating if found.
// Uses a panic recovery guard, see resolveStringer for details.
func resolveGoStringer(v Value) (fn func() string, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fn, ok = nil, false
		}
	}()
	result, found := callMethod(nil, "GoString", v)
	if !found {
		// If not found, try with pointer to the value (for pointer receiver methods)
		if rv, rvOK := v.ReflectValue(); rvOK && rv.Kind() == reflect.Struct {
			ptrRV := reflect.New(rv.Type())
			ptrRV.Elem().Set(rv)
			ptrValue := MakeFromReflect(ptrRV)
			result, found = callMethod(nil, "GoString", ptrValue)
		}
	}

	if !found {
		return nil, false
	}
	str := result.String()
	return func() string { return str }, true
}
