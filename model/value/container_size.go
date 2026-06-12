package value

import (
	"fmt"
	"reflect"
)

// Len returns the length of string, slice, array, map, or chan.
func (v Value) Len() int {
	switch v.kind {
	case KindString:
		return len(v.obj.(string))
	case KindSlice:
		// Native int slice fast path
		if s, ok := v.obj.([]int64); ok {
			return len(s)
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Len()
		}
		panic("invalid obj in Len()")
	case KindArray, KindMap, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Len()
		}
		panic("invalid reflect.Value in Len()")
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			switch rv.Kind() {
			case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
				return rv.Len()
			}
		}
		panic(fmt.Sprintf("cannot take len of reflect kind %v", v.obj))
	default:
		panic(fmt.Sprintf("cannot take len of %v", v.kind))
	}
}

// Cap returns the capacity of slice, array, or chan.
func (v Value) Cap() int {
	switch v.kind {
	case KindSlice:
		if s, ok := v.obj.([]int64); ok {
			return cap(s)
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Cap()
		}
		panic("invalid obj in Cap()")
	case KindArray, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Cap()
		}
		panic("invalid reflect.Value in Cap()")
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.Chan:
				return rv.Cap()
			}
		}
		panic(fmt.Sprintf("cannot take cap of reflect kind %v", v.obj))
	default:
		panic(fmt.Sprintf("cannot take cap of %v", v.kind))
	}
}
