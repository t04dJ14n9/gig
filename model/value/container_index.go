package value

import (
	"fmt"
	"reflect"
)

// Index returns element at index i for slice, array, or string.
func (v Value) Index(i int) Value {
	switch v.kind {
	case KindString:
		// s[i] returns a byte (uint8), not a string
		return MakeUint8(v.obj.(string)[i])
	case KindSlice:
		// Native int slice fast path
		if s, ok := v.obj.([]int64); ok {
			return MakeInt(s[i])
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid obj in Index()")
	case KindArray:
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	case KindReflect:
		// Handle reflect.Value containing a slice
		if rv, ok := v.obj.(reflect.Value); ok {
			return indexReflectSlice(rv, i)
		}
		// Handle native []value.Value slice
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	default:
		panic(fmt.Sprintf("cannot index %v", v.kind))
	}
}

// indexReflectSlice handles indexing into a reflect.Value slice/array.
func indexReflectSlice(rv reflect.Value, i int) Value {
	elem := rv.Index(i)
	if rv.Type().Elem().Kind() == reflect.Func {
		if val, ok := elem.Interface().(Value); ok {
			return val
		}
	}
	if rv.Type().Elem() == reflect.TypeOf(Value{}) {
		return elem.Interface().(Value)
	}
	return MakeFromReflect(elem)
}

// SetIndex sets element at index i for slice or array.
func (v Value) SetIndex(i int, val Value) {
	// Native int slice fast path
	if v.kind == KindSlice {
		if s, ok := v.obj.([]int64); ok {
			s[i] = val.RawInt()
			return
		}
	}
	if rv, ok := v.obj.(reflect.Value); ok {
		elemType := rv.Type().Elem()
		if elemType.Kind() == reflect.Func {
			rv.Index(i).Set(ReflectValueForSet(val, elemType))
			return
		}
		if elemType == reflect.TypeOf(Value{}) {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		rv.Index(i).Set(ReflectValueForSet(val, elemType))
		return
	}
	if slice, ok := v.obj.([]Value); ok {
		slice[i] = val
		return
	}
	panic("invalid reflect.Value in SetIndex()")
}
