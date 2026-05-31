package value

import "reflect"

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

	// Direct type match: err's type is assignable to target element type.
	if errType.AssignableTo(elemType) {
		targetVal.Elem().Set(errVal)
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name.
	if wrapper, ok := err.(*gigStructWrapper); ok {
		// Case 1: target is **StructType (errors.As(&ce) where ce is *CustomError)
		if elemType.Kind() == reflect.Ptr {
			ptrElemType := elemType.Elem()

			// Check if the wrapper's underlying value type is assignable.
			ifaceType := reflect.TypeOf(wrapper.iface)
			if ifaceType.AssignableTo(elemType) {
				targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
				return true
			}

			// Match by gig type name: compare wrapper's typeName with target's gig tag.
			wrapperTypeName := extractBareTypeName(wrapper.typeName)
			targetTypeName := extractGigTagFromType(ptrElemType)
			if targetTypeName == "" {
				targetTypeName = ptrElemType.Name()
			}
			targetTypeName = extractBareTypeName(targetTypeName)

			if wrapperTypeName != "" && wrapperTypeName == targetTypeName {
				// Type names match: set the target to the underlying value.
				if ifaceType.AssignableTo(elemType) {
					targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
					return true
				}
				// Try converting through interface{} if the pointer element is also a gig struct.
				if ifaceType.Kind() == reflect.Ptr && ptrElemType.Kind() == reflect.Struct {
					// Both are pointers to structs: try setting the value directly.
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

	// For non-gig errors, try standard interface check.
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

	// Direct type match: err's type equals target type.
	if errType == targetType || errType.AssignableTo(targetType) {
		*slot = MakeFromReflect(errVal)
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name.
	if wrapper, ok := err.(*gigStructWrapper); ok {
		ifaceType := reflect.TypeOf(wrapper.iface)
		if ifaceType == targetType || ifaceType.AssignableTo(targetType) {
			*slot = MakeFromReflect(reflect.ValueOf(wrapper.iface))
			return true
		}
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
