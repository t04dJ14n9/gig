package value

import "reflect"

// ErrorValue extracts a Go error from a value.Value.
// If the value is already a native Go error, returns it directly.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a gigStructWrapper that implements the error interface.
// Otherwise returns nil.
//
// This is the boundary function for generated DirectCall wrappers —
// use it whenever extracting an error-typed parameter from args.
func ErrorValue(v Value) error {
	iface := v.Interface()
	if iface == nil {
		return nil
	}
	// If it's already a Go error (e.g., from fmt.Errorf, errors.New), return as-is
	if e, ok := iface.(error); ok {
		return e
	}
	typeName := isGigStruct(iface)
	if typeName == "" {
		return nil
	}
	// Check if the interpreted type has an Error() method
	errorerFunc, hasError := resolveErrorer(v)
	if !hasError {
		return nil
	}
	// Also capture Stringer/GoStringer lazily so the wrapper is complete.
	stringerFunc, _ := resolveStringer(v)
	gostringerFunc, _ := resolveGoStringer(v)
	constFn := func(s string) func() string { return func() string { return s } }
	errFn := func() string { return errorerFunc() }

	var lazyStringer, lazyErrorer, lazyGoStringer func() (func() string, bool)
	if stringerFunc != nil {
		sf := constFn(stringerFunc())
		lazyStringer = func() (func() string, bool) { return sf, true }
	}
	lazyErrorer = func() (func() string, bool) { return errFn, true }
	if gostringerFunc != nil {
		gf := constFn(gostringerFunc())
		lazyGoStringer = func() (func() string, bool) { return gf, true }
	}
	return &gigStructWrapper{
		iface:          iface,
		typeName:       typeName,
		lazyStringer:   lazyStringer,
		lazyErrorer:    lazyErrorer,
		lazyGoStringer: lazyGoStringer,
	}
}

// ErrorWrap prepares a value.Value for use as a Go error.
// If the value is an interpreter-synthesized struct with an Error() method,
// returns a wrapper that implements the error interface. Otherwise returns
// the raw interface{} value.
//
// Deprecated: Use ErrorValue instead for typed error extraction.
func ErrorWrap(v Value) any {
	if e := ErrorValue(v); e != nil {
		return e
	}
	return v.Interface()
}

// GigErrorsAs implements errors.As semantics for interpreter-defined types.
// It mirrors the standard library's errors.As but uses the interpreter's type
// name registry for matching, since reflect.StructOf types can't implement
// interfaces and have different reflect.Type identities than named Go types.
//
// err is the error value (may be a gigStructWrapper or a native Go error).
// target is a pointer to the target type (e.g., **CustomError as interface{}).
//   - If target is a *Value (frame slot pointer from OpAddr), GigErrorsAs
//     unwraps it and sets the frame slot directly on match.
//   - Otherwise, target is a normal Go double pointer.
//
// Returns true if the error (or any error in its Unwrap chain) matches target.
func GigErrorsAs(err error, target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	if err == nil {
		return false
	}

	// Handle *Value frame slot pointers from OpAddr.
	// The frame slot holds the interpreted value; we need to match against
	// its type and set it directly.
	if vp, ok := target.(*Value); ok {
		slotVal := vp.Interface()
		if slotVal == nil {
			return false
		}
		slotRV := reflect.ValueOf(slotVal)
		if !slotRV.IsValid() || slotRV.Kind() != reflect.Ptr {
			return false
		}
		elemType := slotRV.Type() // The target type (e.g., *myError)
		return gigAsWalkFrameSlot(err, elemType, vp)
	}

	targetVal := reflect.ValueOf(target)
	targetType := targetVal.Type()
	if targetType.Kind() != reflect.Ptr || targetVal.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	elemType := targetType.Elem()

	return gigAsWalkValue(err, elemType, targetVal)
}

// gigAsWalkValue walks the error chain (including multi-unwrap from errors.Join)
// and returns true if any error matches the target type.
func gigAsWalkValue(err error, elemType reflect.Type, targetVal reflect.Value) bool {
	for {
		if gigAsMatchValue(err, elemType, targetVal) {
			return true
		}
		// Multi-unwrap (errors.Join)
		if x, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range x.Unwrap() {
				if e != nil && gigAsWalkValue(e, elemType, targetVal) {
					return true
				}
			}
			return false
		}
		// Single unwrap
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = unwrapper.Unwrap()
		if err == nil {
			return false
		}
	}
}

// gigAsWalkFrameSlot walks the error chain for the frame-slot path.
func gigAsWalkFrameSlot(err error, targetType reflect.Type, slot *Value) bool {
	for {
		if gigAsMatchFrameSlot(err, targetType, slot) {
			return true
		}
		// Multi-unwrap (errors.Join)
		if x, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range x.Unwrap() {
				if e != nil && gigAsWalkFrameSlot(e, targetType, slot) {
					return true
				}
			}
			return false
		}
		// Single unwrap
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = unwrapper.Unwrap()
		if err == nil {
			return false
		}
	}
}

// gigAsMatchValue checks if an error value matches the target type.
// It handles both native Go types (via reflect.AssignableTo) and
// interpreter-defined types (via gig type name matching).
func gigAsMatchValue(err error, elemType reflect.Type, targetVal reflect.Value) bool {
	errVal := reflect.ValueOf(err)
	errType := errVal.Type()

	if wrapper, ok := err.(*gigStructWrapper); ok && elemType.Kind() == reflect.Interface && elemType.NumMethod() == 0 {
		targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
		return true
	}

	// Direct type match: err's type is assignable to target element type
	if errType.AssignableTo(elemType) {
		targetVal.Elem().Set(errVal)
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name
	if wrapper, ok := err.(*gigStructWrapper); ok {
		// Case 1: target is **StructType (errors.As(&ce) where ce is *CustomError)
		if elemType.Kind() == reflect.Ptr {
			ptrElemType := elemType.Elem()

			// Check if the wrapper's underlying value type is assignable
			ifaceType := reflect.TypeOf(wrapper.iface)
			if ifaceType.AssignableTo(elemType) {
				targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
				return true
			}

			// Match by gig type name: compare wrapper's typeName with target's gig tag
			wrapperTypeName := extractBareTypeName(wrapper.typeName)
			targetTypeName := extractGigTagFromType(ptrElemType)
			if targetTypeName == "" {
				targetTypeName = ptrElemType.Name()
			}
			targetTypeName = extractBareTypeName(targetTypeName)

			if wrapperTypeName != "" && wrapperTypeName == targetTypeName {
				// Type names match — set the target to the underlying value
				if ifaceType.AssignableTo(elemType) {
					targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
					return true
				}
				// Try converting through interface{} if the pointer element is also a gig struct
				if ifaceType.Kind() == reflect.Ptr && ptrElemType.Kind() == reflect.Struct {
					// Both are pointers to structs — try setting the value directly
					if ifaceType.Elem().ConvertibleTo(ptrElemType) {
						converted := reflect.ValueOf(wrapper.iface).Elem().Convert(ptrElemType)
						ptr := reflect.New(ptrElemType)
						ptr.Elem().Set(converted)
						targetVal.Elem().Set(ptr)
						return true
					}
				}
			}
		}

		// Case 2: target is an interface type (e.g., error)
		if elemType.Kind() == reflect.Interface {
			if errType.Implements(elemType) {
				targetVal.Elem().Set(errVal)
				return true
			}
		}

		// Case 3: target is a struct type (value receiver, unlikely for errors)
		if elemType.Kind() == reflect.Struct {
			ifaceType := reflect.TypeOf(wrapper.iface)
			if ifaceType.AssignableTo(elemType) {
				targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
				return true
			}
		}
	}

	// For non-gig errors, try standard interface check
	if elemType.Kind() == reflect.Interface && errType.Implements(elemType) {
		targetVal.Elem().Set(errVal)
		return true
	}

	return false
}

// gigAsMatchFrameSlot checks if an error matches a target type and sets
// the frame slot directly. Used when the target is a *Value from OpAddr.
func gigAsMatchFrameSlot(err error, targetType reflect.Type, slot *Value) bool {
	errVal := reflect.ValueOf(err)
	errType := errVal.Type()

	// Direct type match: err's type equals target type
	if errType == targetType || errType.AssignableTo(targetType) {
		*slot = MakeFromReflect(errVal)
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name
	if wrapper, ok := err.(*gigStructWrapper); ok {
		ifaceType := reflect.TypeOf(wrapper.iface)
		if ifaceType == targetType || ifaceType.AssignableTo(targetType) {
			*slot = MakeFromReflect(reflect.ValueOf(wrapper.iface))
			return true
		}
		// Match by gig type name
		wrapperTypeName := extractBareTypeName(wrapper.typeName)
		targetTypeName := extractGigTagFromType(targetType)
		if targetTypeName == "" {
			targetTypeName = targetType.Name()
		}
		targetTypeName = extractBareTypeName(targetTypeName)
		if wrapperTypeName != "" && wrapperTypeName == targetTypeName {
			*slot = MakeFromReflect(reflect.ValueOf(wrapper.iface))
			return true
		}
	}

	return false
}

// extractBareTypeName extracts the short type name from a qualified name.
// "known_issues.CustomError" → "CustomError"

// GigErrorsIs implements errors.Is semantics for interpreter-defined types.
// It replicates the standard library algorithm but uses gig's method resolution
// to invoke custom Is(error) bool and Unwrap() error methods on gig types.
func GigErrorsIs(errVal Value, targetVal Value) bool {
	err := ErrorValue(errVal)
	target := ErrorValue(targetVal)
	if err == nil && target == nil {
		return true
	}
	if err == nil || target == nil {
		return err == target
	}

	for {
		// Direct comparison — also handles gigStructWrapper by comparing underlying values
		if err == target {
			return true
		}
		if gigErrorsEqual(err, target) {
			return true
		}

		// Check custom Is() method on gig types
		if _, ok := err.(*gigStructWrapper); ok {
			// Try to call Is(target) via compiled method
			if result, found := callMethodWithArgs("Is", errVal, []Value{targetVal}); found {
				if result.Kind() == KindBool && result.Bool() {
					return true
				}
			}
		} else if x, ok := err.(interface{ Is(error) bool }); ok {
			if x.Is(target) {
				return true
			}
		}

		// Unwrap
		if _, ok := err.(*gigStructWrapper); ok {
			unwrapResult, found := callMethod(nil, "Unwrap", errVal)
			if !found {
				return false
			}
			unwrapped := ErrorValue(unwrapResult)
			if unwrapped == nil {
				return false
			}
			err = unwrapped
			errVal = unwrapResult
		} else if x, ok := err.(interface{ Unwrap() []error }); ok {
			// Multi-unwrap (errors.Join): recursively check each wrapped error
			for _, e := range x.Unwrap() {
				if e != nil && GigErrorsIs(FromInterface(e), targetVal) {
					return true
				}
			}
			return false
		} else if x, ok := err.(interface{ Unwrap() error }); ok {
			err = x.Unwrap()
			if err == nil {
				return false
			}
			errVal = FromInterface(err)
		} else {
			return false
		}
	}
}

// gigErrorsEqual compares two errors for equality, handling gigStructWrapper.
// Two gigStructWrappers are equal if they wrap the same type and underlying value.
func gigErrorsEqual(a, b error) bool {
	wa, aIsGig := a.(*gigStructWrapper)
	wb, bIsGig := b.(*gigStructWrapper)
	if aIsGig && bIsGig {
		return wa.typeName == wb.typeName && reflect.DeepEqual(wa.iface, wb.iface)
	}
	return false
}

// GigErrorsUnwrap implements errors.Unwrap for interpreter-defined types.
// If the error is a gig type with an Unwrap() method, invokes it via the
// compiled method resolver. Otherwise delegates to standard errors.Unwrap.
func GigErrorsUnwrap(errVal Value) Value {
	err := ErrorValue(errVal)
	if err == nil {
		return MakeNil()
	}

	// For gig types, use compiled method resolution
	if _, ok := err.(*gigStructWrapper); ok {
		result, found := callMethod(nil, "Unwrap", errVal)
		if found {
			return result
		}
		return MakeNil()
	}

	// For native Go errors, use standard unwrap
	if x, ok := err.(interface{ Unwrap() error }); ok {
		unwrapped := x.Unwrap()
		if unwrapped == nil {
			return MakeNil()
		}
		return FromInterface(unwrapped)
	}
	return MakeNil()
}
