package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

// appendValue implements the append builtin for the VM.
// It handles native int slices, byte slices, reflect slices, and nil slices.
func appendValue(slice, elem value.Value) value.Value {
	if s, ok := slice.IntSlice(); ok {
		return appendToIntSlice(s, elem)
	}

	if slice.Kind() == value.KindBytes {
		return appendToByteSlice(slice, elem)
	}

	if rv, ok := slice.ReflectValue(); ok {
		return appendToReflectSlice(rv, elem)
	}

	if slice.IsNil() || slice.Kind() == value.KindInvalid {
		return appendToNilSlice(elem)
	}

	return slice
}

func appendToByteSlice(slice, elem value.Value) value.Value {
	b, ok := slice.Bytes()
	if !ok {
		return slice
	}
	appended, ok := appendByteElement(b, elem)
	if !ok {
		return slice
	}
	return value.MakeBytes(appended)
}

func appendByteElement(b []byte, elem value.Value) ([]byte, bool) {
	switch elem.Kind() {
	case value.KindUint:
		return append(b, byte(elem.Uint())), true
	case value.KindInt:
		return append(b, byte(elem.RawInt())), true
	case value.KindBytes:
		return appendByteSlice(b, elem)
	case value.KindString:
		return append(b, elem.String()...), true
	default:
		return appendByteInterface(b, elem)
	}
}

func appendByteSlice(b []byte, elem value.Value) ([]byte, bool) {
	eb, ok := elem.Bytes()
	if !ok {
		return b, false
	}
	return append(b, eb...), true
}

func appendByteInterface(b []byte, elem value.Value) ([]byte, bool) {
	v := elem.Interface()
	if v == nil {
		return b, false
	}
	if bv, ok := v.(byte); ok {
		return append(b, bv), true
	}
	if bv, ok := v.(uint8); ok {
		return append(b, bv), true
	}
	return b, false
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
