package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

// intSliceToReflect converts a native []int64 to a reflect []int slice.
func intSliceToReflect(s []int64) reflect.Value {
	rs := reflect.MakeSlice(reflect.TypeOf([]int{}), len(s), cap(s))
	for i, v := range s {
		rs.Index(i).SetInt(v)
	}
	return rs
}

// appendValue implements the append builtin for the VM.
// It handles native int slices, byte slices, reflect slices, and nil slices.
func appendValue(slice, elem value.Value) value.Value {
	// Fast path: native []int64 slice
	if s, ok := slice.IntSlice(); ok {
		return appendToIntSlice(s, elem)
	}

	// Byte slice ([]byte / KindBytes)
	if slice.Kind() == value.KindBytes {
		if b, ok := slice.Bytes(); ok {
			if elem.Kind() == value.KindUint || elem.Kind() == value.KindInt {
				return value.MakeBytes(append(b, byte(elem.Uint())))
			}
			// If elem is a byte slice (spread append)
			if elem.Kind() == value.KindBytes {
				if eb, ok := elem.Bytes(); ok {
					return value.MakeBytes(append(b, eb...))
				}
			}
			// If elem is a string (append(b, "str"...))
			if elem.Kind() == value.KindString {
				return value.MakeBytes(append(b, elem.String()...))
			}
			// Fallback: convert via interface
			if v := elem.Interface(); v != nil {
				if bv, ok := v.(byte); ok {
					return value.MakeBytes(append(b, bv))
				}
				if bv, ok := v.(uint8); ok {
					return value.MakeBytes(append(b, bv))
				}
			}
		}
		return slice
	}

	// Native []int64 that needs reflect conversion (e.g., stored in [][]int)
	if slice.Kind() == value.KindSlice {
		if intSlice, ok := slice.IntSlice(); ok {
			return appendIntSliceViaReflect(intSlice, elem)
		}
	}

	// Reflect-based slice
	if rv, ok := slice.ReflectValue(); ok {
		return appendToReflectSlice(rv, elem)
	}

	// Nil slice: create a new slice
	if slice.IsNil() || slice.Kind() == value.KindInvalid {
		return appendToNilSlice(elem)
	}

	return slice
}

// appendToIntSlice appends to a native []int64.
func appendToIntSlice(s []int64, elem value.Value) value.Value {
	if es, ok := elem.IntSlice(); ok {
		return value.MakeIntSlice(append(s, es...))
	}
	if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice {
		// elem is a reflect-based integer slice (e.g. []int from a [][]int range)
		for i := 0; i < elemRV.Len(); i++ {
			s = append(s, elemRV.Index(i).Int())
		}
		return value.MakeIntSlice(s)
	}
	return value.MakeIntSlice(append(s, elem.RawInt()))
}

// appendIntSliceViaReflect converts []int64 to reflect []int and appends.
func appendIntSliceViaReflect(intSlice []int64, elem value.Value) value.Value {
	rv := intSliceToReflect(intSlice)
	if elem.Kind() == value.KindInt {
		return value.MakeFromReflect(reflect.Append(rv, reflect.ValueOf(int(elem.RawInt()))))
	}
	if elem.Kind() == value.KindSlice {
		if elemIntSlice, ok := elem.IntSlice(); ok {
			return value.MakeFromReflect(reflect.AppendSlice(rv, intSliceToReflect(elemIntSlice)))
		}
	}
	return value.MakeFromReflect(reflect.Append(rv, elem.ToReflectValue(reflect.TypeOf(int(0)))))
}

var valueValueType = reflect.TypeOf(value.Value{})

// appendToReflectSlice appends to a reflect.Value slice.
func appendToReflectSlice(rv reflect.Value, elem value.Value) value.Value {
	sliceElemType := rv.Type().Elem()

	// []value.Value slices (function slices)
	if sliceElemType == valueValueType {
		if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice && elemRV.Type().Elem() == sliceElemType {
			return value.MakeFromReflect(reflect.AppendSlice(rv, elemRV))
		}
		return value.MakeFromReflect(reflect.Append(rv, reflect.ValueOf(elem)))
	}

	// Check if elem is native []int64 needing spread-append
	if elem.Kind() == value.KindSlice {
		if elemIntSlice, ok := elem.IntSlice(); ok {
			for _, v := range elemIntSlice {
				rv = reflect.Append(rv, reflect.ValueOf(int(v)))
			}
			return value.MakeFromReflect(rv)
		}
	}

	// SSA-packed variadic slice spread
	if elemRV, ok := elem.ReflectValue(); ok && elemRV.Kind() == reflect.Slice {
		if elemRV.Type().AssignableTo(rv.Type()) {
			return value.MakeFromReflect(reflect.AppendSlice(rv, elemRV))
		}
	}

	return value.MakeFromReflect(reflect.Append(rv, elem.ToReflectValue(sliceElemType)))
}

// appendToNilSlice creates a new slice from a nil/zero slice and appends.
func appendToNilSlice(elem value.Value) value.Value {
	// Native []int64 spread
	if es, ok := elem.IntSlice(); ok {
		return value.MakeIntSlice(append([]int64(nil), es...))
	}

	elemRV, ok := elem.ReflectValue()
	if ok && elemRV.Kind() == reflect.Slice {
		sliceType := reflect.SliceOf(elemRV.Type().Elem())
		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		return value.MakeFromReflect(reflect.AppendSlice(newSlice, elemRV))
	}

	// Single-element append: infer type from value
	elemIface := elem.Interface()
	if elemIface != nil {
		elemRV2 := reflect.ValueOf(elemIface)
		sliceType := reflect.SliceOf(elemRV2.Type())
		newSlice := reflect.MakeSlice(sliceType, 0, 0)
		return value.MakeFromReflect(reflect.Append(newSlice, elemRV2))
	}
	return value.MakeNil()
}
