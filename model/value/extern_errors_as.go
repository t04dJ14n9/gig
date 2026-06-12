package value

import "reflect"

// GigErrorsAs implements errors.As semantics for interpreter-defined types.
func GigErrorsAs(err error, target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	if err == nil {
		return false
	}

	if vp, ok := target.(*Value); ok {
		slotVal := vp.Interface()
		if slotVal == nil {
			return false
		}
		slotRV := reflect.ValueOf(slotVal)
		if !slotRV.IsValid() || slotRV.Kind() != reflect.Ptr {
			return false
		}
		return gigAsWalkFrameSlot(err, slotRV.Type(), vp)
	}

	targetVal := reflect.ValueOf(target)
	targetType := targetVal.Type()
	if targetType.Kind() != reflect.Ptr || targetVal.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	return gigAsWalkValue(err, targetType.Elem(), targetVal)
}

func gigAsWalkValue(err error, elemType reflect.Type, targetVal reflect.Value) bool {
	for {
		if gigAsMatchValue(err, elemType, targetVal) {
			return true
		}
		if x, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range x.Unwrap() {
				if e != nil && gigAsWalkValue(e, elemType, targetVal) {
					return true
				}
			}
			return false
		}
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

func gigAsWalkFrameSlot(err error, targetType reflect.Type, slot *Value) bool {
	for {
		if gigAsMatchFrameSlot(err, targetType, slot) {
			return true
		}
		if x, ok := err.(interface{ Unwrap() []error }); ok {
			for _, e := range x.Unwrap() {
				if e != nil && gigAsWalkFrameSlot(e, targetType, slot) {
					return true
				}
			}
			return false
		}
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

func gigAsMatchValue(err error, elemType reflect.Type, targetVal reflect.Value) bool {
	errVal := reflect.ValueOf(err)
	errType := errVal.Type()

	wrapper, isWrapper := asGigStructError(err)
	if isWrapper && gigAsMatchEmptyInterface(wrapper, elemType, targetVal) {
		return true
	}
	if gigAsSetAssignable(errVal, elemType, targetVal) {
		return true
	}
	if isWrapper && gigAsMatchWrapper(wrapper, errVal, errType, elemType, targetVal) {
		return true
	}
	return gigAsSetImplementedInterface(errVal, errType, elemType, targetVal)
}

func gigAsMatchEmptyInterface(wrapper *gigStructError, elemType reflect.Type, targetVal reflect.Value) bool {
	if elemType.Kind() != reflect.Interface || elemType.NumMethod() != 0 {
		return false
	}
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
	wrapper *gigStructError,
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
		return gigAsSetAssignable(reflect.ValueOf(wrapper.iface), elemType, targetVal)
	default:
		return false
	}
}

func gigAsMatchWrapperPointer(wrapper *gigStructError, elemType reflect.Type, targetVal reflect.Value) bool {
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

func gigAsWrapperTypeNameMatches(wrapper *gigStructError, targetType reflect.Type) bool {
	wrapperTypeName := extractBareTypeName(wrapper.typeName)
	targetTypeName := extractGigTagFromType(targetType)
	if targetTypeName == "" {
		targetTypeName = targetType.Name()
	}
	return wrapperTypeName != "" && wrapperTypeName == extractBareTypeName(targetTypeName)
}

func gigAsSetConvertedWrapperPointer(wrapper *gigStructError, targetType reflect.Type, targetVal reflect.Value) bool {
	ifaceVal := reflect.ValueOf(wrapper.iface)
	if ifaceVal.Kind() != reflect.Ptr || targetType.Kind() != reflect.Struct || !ifaceVal.Type().Elem().ConvertibleTo(targetType) {
		return false
	}
	converted := ifaceVal.Elem().Convert(targetType)
	ptr := reflect.New(targetType)
	ptr.Elem().Set(converted)
	targetVal.Elem().Set(ptr)
	return true
}

func gigAsSetImplementedInterface(errVal reflect.Value, errType, elemType reflect.Type, targetVal reflect.Value) bool {
	if elemType.Kind() != reflect.Interface || !errType.Implements(elemType) {
		return false
	}
	targetVal.Elem().Set(errVal)
	return true
}

func gigAsMatchFrameSlot(err error, targetType reflect.Type, slot *Value) bool {
	errVal := reflect.ValueOf(err)
	errType := errVal.Type()

	if errType == targetType || errType.AssignableTo(targetType) {
		*slot = MakeFromReflect(errVal)
		return true
	}

	if wrapper, ok := asGigStructError(err); ok {
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
