package value

import "reflect"

// Elem dereferences a pointer or returns the underlying value of interface.
func (v Value) Elem() Value {
	// Fast path: *int64 pointer (from native int slice)
	if ptr, ok := v.obj.(*int64); ok {
		return MakeInt(*ptr)
	}
	// Fast path: *Value pointer (from OpGlobal / OpAddr on value.Value locals / OpFree)
	if ptr, ok := v.obj.(*Value); ok {
		return *ptr
	}
	if rv, ok := v.ReflectValue(); ok {
		// If the reflect.Value points to a value.Value struct, unwrap it directly.
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			if vp, ok2 := rv.Interface().(*Value); ok2 {
				return *vp
			}
		}
		return MakeFromReflect(rv.Elem())
	}
	panic("invalid reflect.Value in Elem()")
}

// SetElem sets the value pointed to by a pointer.
func (v Value) SetElem(val Value) {
	if setDirectElemPointer(v.obj, val) {
		return
	}
	if rv, ok := v.ReflectValue(); ok {
		if setReflectElem(rv, val) {
			return
		}
	}
	panic("invalid reflect.Value in SetElem()")
}

func setDirectElemPointer(obj any, val Value) bool {
	if ptr, ok := obj.(*int64); ok {
		*ptr = val.num
		return true
	}
	if ptr, ok := obj.(*Value); ok {
		*ptr = val
		return true
	}
	return false
}

func setReflectElem(rv reflect.Value, val Value) bool {
	switch rv.Kind() {
	case reflect.Ptr:
		setReflectPointerElem(rv, val)
	case reflect.Interface:
		rv.Set(ReflectValueForSet(val, rv.Type()))
	case reflect.Struct:
		// Struct values are not element-addresses; retain the previous graceful no-op.
	default:
		return false
	}
	return true
}

func setReflectPointerElem(rv reflect.Value, val Value) {
	elemType := rv.Type().Elem()
	if elemType.Kind() == reflect.Func {
		rv.Elem().Set(ReflectValueForSet(val, elemType))
		return
	}
	if elemType.Name() == "Value" && elemType.PkgPath() == "github.com/t04dJ14n9/gig/model/value" {
		ptr := rv.Interface().(*Value)
		*ptr = val
		return
	}
	if setReflectPointerPayloadElem(rv, val, elemType) {
		return
	}

	targetRV := rv.Elem()
	if !targetRV.CanSet() {
		return
	}
	if setNativeIntSliceTarget(targetRV, val, elemType) {
		return
	}

	valRV := ReflectValueForSet(val, elemType)
	if !valRV.Type().AssignableTo(elemType) {
		if setConcreteSliceAsInterfaceSlice(targetRV, valRV, elemType) {
			return
		}
		if setPointerValueElem(targetRV, valRV, elemType) {
			return
		}
	}
	targetRV.Set(valRV)
}

func setReflectPointerPayloadElem(rv reflect.Value, val Value, elemType reflect.Type) bool {
	if val.Kind() != KindReflect {
		return false
	}
	valRV, ok := val.obj.(reflect.Value)
	if !ok || valRV.Kind() != reflect.Ptr {
		return false
	}
	if elemType.Kind() == reflect.Interface {
		rv.Elem().Set(valRV)
		return true
	}
	if valRV.Type().Elem() == elemType {
		rv.Elem().Set(valRV.Elem())
		return true
	}
	return false
}

func setNativeIntSliceTarget(targetRV reflect.Value, val Value, elemType reflect.Type) bool {
	if val.kind != KindSlice {
		return false
	}
	s, isInt := val.obj.([]int64)
	if !isInt || elemType.Kind() != reflect.Slice {
		return false
	}

	target := reflect.MakeSlice(elemType, len(s), cap(s))
	for i, n := range s {
		target.Index(i).SetInt(n)
	}
	targetRV.Set(target)
	return true
}

func setConcreteSliceAsInterfaceSlice(targetRV, valRV reflect.Value, elemType reflect.Type) bool {
	if elemType.Kind() != reflect.Slice || valRV.Kind() != reflect.Slice {
		return false
	}
	if elemType.Elem().Kind() != reflect.Interface || valRV.Type().Elem().Kind() == reflect.Interface {
		return false
	}
	convertedSlice := convertSliceToInterface(valRV, elemType)
	if !convertedSlice.Type().AssignableTo(elemType) {
		return false
	}
	targetRV.Set(convertedSlice)
	return true
}

func setPointerValueElem(targetRV, valRV reflect.Value, elemType reflect.Type) bool {
	if valRV.Kind() != reflect.Ptr || valRV.IsNil() || !valRV.Type().Elem().AssignableTo(elemType) {
		return false
	}
	targetRV.Set(valRV.Elem())
	return true
}
