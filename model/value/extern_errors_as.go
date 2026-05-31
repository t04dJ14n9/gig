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

	wrapper, isWrapper := err.(*gigStructWrapper)
	if isWrapper && gigAsMatchEmptyInterface(wrapper, elemType, targetVal) {
		return true
	}

	// Direct type match: err's type is assignable to target element type.
	if gigAsSetAssignable(errVal, elemType, targetVal) {
		return true
	}

	// If err is a *gigStructWrapper, try matching by interpreter type name.
	if isWrapper && gigAsMatchWrapper(wrapper, errVal, errType, elemType, targetVal) {
		return true
	}

	// For non-gig errors, try standard interface check.
	return gigAsSetImplementedInterface(errVal, errType, elemType, targetVal)
}

func gigAsMatchEmptyInterface(wrapper *gigStructWrapper, elemType reflect.Type, targetVal reflect.Value) bool {
	if elemType.Kind() != reflect.Interface || elemType.NumMethod() != 0 {
		return false
	}
	// For *any targets, expose the script value itself. Setting the wrapper
	// would leak the fmt/error adapter instead of the interpreter struct.
	targetVal.Elem().Set(reflect.ValueOf(wrapper.iface))
	return true
}

func gigAsSetAssignable(val reflect.Value, elemType reflect.Type, targetVal reflect.Value) bool {
	if !val.Type().AssignableTo(elemType) {
		return false
	}
	targetVal.Elem().Set(val)
	return true
}

func gigAsMatchWrapper(
	wrapper *gigStructWrapper,
	errVal reflect.Value,
	errType reflect.Type,
	elemType reflect.Type,
	targetVal reflect.Value,
) bool {
	switch elemType.Kind() {
	case reflect.Ptr:
		return gigAsMatchWrapperPointer(wrapper, elemType, targetVal)
	case reflect.Interface:
		return gigAsSetImplementedInterface(errVal, errType, elemType, targetVal)
	case reflect.Struct:
		return gigAsMatchWrapperStruct(wrapper, elemType, targetVal)
	default:
		return false
	}
}

func gigAsMatchWrapperPointer(wrapper *gigStructWrapper, elemType reflect.Type, targetVal reflect.Value) bool {
	ifaceVal := reflect.ValueOf(wrapper.iface)
	if gigAsSetAssignable(ifaceVal, elemType, targetVal) {
		return true
	}

	ptrElemType := elemType.Elem()
	if !gigAsWrapperTypeNameMatches(wrapper, ptrElemType) {
		return false
	}
	return gigAsSetConvertedWrapperPointer(wrapper, ptrElemType, targetVal)
}

func gigAsWrapperTypeNameMatches(wrapper *gigStructWrapper, targetType reflect.Type) bool {
	// reflect.StructOf gives script structs anonymous host identities, so
	// errors.As must compare the embedded Gig type name before conversion.
	wrapperTypeName := extractBareTypeName(wrapper.typeName)
	targetTypeName := extractGigTagFromType(targetType)
	if targetTypeName == "" {
		targetTypeName = targetType.Name()
	}
	return wrapperTypeName != "" && wrapperTypeName == extractBareTypeName(targetTypeName)
}

func gigAsSetConvertedWrapperPointer(wrapper *gigStructWrapper, targetType reflect.Type, targetVal reflect.Value) bool {
	ifaceVal := reflect.ValueOf(wrapper.iface)
	if ifaceVal.Kind() != reflect.Ptr || targetType.Kind() != reflect.Struct || !ifaceVal.Type().Elem().ConvertibleTo(targetType) {
		return false
	}
	// Allocate a fresh pointer of the target type so the caller receives the
	// shape requested by errors.As, not the anonymous StructOf pointer.
	converted := ifaceVal.Elem().Convert(targetType)
	ptr := reflect.New(targetType)
	ptr.Elem().Set(converted)
	targetVal.Elem().Set(ptr)
	return true
}

func gigAsMatchWrapperStruct(wrapper *gigStructWrapper, elemType reflect.Type, targetVal reflect.Value) bool {
	return gigAsSetAssignable(reflect.ValueOf(wrapper.iface), elemType, targetVal)
}

func gigAsSetImplementedInterface(errVal reflect.Value, errType, elemType reflect.Type, targetVal reflect.Value) bool {
	if elemType.Kind() != reflect.Interface || !errType.Implements(elemType) {
		return false
	}
	// Non-empty interfaces such as error should receive the wrapper, because
	// that is the host value carrying the method implementation.
	targetVal.Elem().Set(errVal)
	return true
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
