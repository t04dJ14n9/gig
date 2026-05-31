// accessor.go provides primitive and interface accessors for Value.
package value

import (
	"fmt"
	"math"
	"reflect"
)

// Bool returns the bool value. Panics if not KindBool.
func (v Value) Bool() bool {
	if v.kind != KindBool {
		panic(fmt.Sprintf("not a bool: %v", v.kind))
	}
	return v.num != 0
}

// Int returns the int value. Panics if not KindInt.
func (v Value) Int() int64 {
	if v.kind != KindInt {
		panic(fmt.Sprintf("not an int: %v", v.kind))
	}
	return v.num
}

// Uint returns the uint value. Panics if not KindUint.
func (v Value) Uint() uint64 {
	if v.kind != KindUint {
		panic(fmt.Sprintf("not a uint: %v", v.kind))
	}
	return uint64(v.num)
}

// Float returns the float value. Panics if not KindFloat.
func (v Value) Float() float64 {
	if v.kind != KindFloat {
		panic(fmt.Sprintf("not a float: %v", v.kind))
	}
	return math.Float64frombits(uint64(v.num))
}

// String returns the string value. Panics if not KindString.
func (v Value) String() string {
	if v.kind != KindString {
		panic(fmt.Sprintf("not a string: %v", v.kind))
	}
	return v.obj.(string)
}

// Complex returns the complex value. Panics if not KindComplex.
func (v Value) Complex() complex128 {
	if v.kind != KindComplex {
		panic(fmt.Sprintf("not a complex: %v", v.kind))
	}
	return v.obj.(complex128)
}

// Interface returns the value as an interface{}.
// For numeric kinds, the returned type matches the original Go type recorded
// by the size field (e.g. int8, int32, int64, float32, etc.).
func (v Value) Interface() any {
	// Keep the public router grouped by representation domain; width-sensitive
	// scalar conversion and native payload extraction stay in focused helpers.
	switch v.kind {
	case KindNil:
		return nil
	case KindBool, KindInt, KindUint, KindFloat, KindString, KindComplex:
		return v.interfaceScalar()
	case KindInterface:
		return v.interfaceInterface()
	case KindFunc, KindBytes, KindSlice:
		return v.interfaceNativePayload()
	case KindReflect:
		return interfaceReflectOrObject(v.obj)
	default:
		return interfaceReflectOrObject(v.obj)
	}
}

func (v Value) interfaceScalar() any {
	switch v.kind {
	case KindBool:
		return v.Bool()
	case KindInt:
		return interfaceSignedInt(v.num, v.size)
	case KindUint:
		return interfaceUnsignedInt(uint64(v.num), v.size)
	case KindFloat:
		return interfaceFloat(v.num, v.size)
	case KindString:
		return v.obj.(string)
	case KindComplex:
		return interfaceComplex(v.obj.(complex128), v.size)
	default:
		return nil
	}
}

func (v Value) interfaceNativePayload() any {
	switch v.kind {
	case KindFunc:
		return v.obj
	case KindBytes:
		return v.obj.([]byte)
	case KindSlice:
		return interfaceSlice(v.obj)
	default:
		return interfaceReflectOrObject(v.obj)
	}
}

func interfaceSignedInt(num int64, size Size) any {
	switch size {
	case Size8:
		return int8(num)
	case Size16:
		return int16(num)
	case Size32:
		return int32(num)
	case Size64:
		return num
	default:
		// SizePtr and Size0 represent the interpreter's default int width.
		return int(num)
	}
}

func interfaceUnsignedInt(num uint64, size Size) any {
	switch size {
	case Size8:
		return uint8(num)
	case Size16:
		return uint16(num)
	case Size32:
		return uint32(num)
	case Size64:
		return num
	default:
		// SizePtr and Size0 represent the interpreter's default uint width.
		return uint(num)
	}
}

func interfaceFloat(bits int64, size Size) any {
	f := math.Float64frombits(uint64(bits))
	if size == Size32 {
		return float32(f)
	}
	return f
}

func interfaceComplex(c complex128, size Size) any {
	if size == Size32 {
		return complex64(c)
	}
	return c
}

func (v Value) interfaceInterface() any {
	if dyn, ok := v.InterpretedInterface(); ok {
		return dyn.Value.Interface()
	}
	return interfaceReflectOrObject(v.obj)
}

func interfaceSlice(obj any) any {
	s, ok := obj.([]int64)
	if !ok {
		return obj
	}

	// Native int slices are stored internally as []int64 but should return as
	// []int so API callers observe the same type native Go would produce.
	result := make([]int, len(s))
	for i, n := range s {
		result[i] = int(n)
	}
	return result
}

func interfaceReflectOrObject(obj any) any {
	if rv, ok := obj.(reflect.Value); ok {
		return rv.Interface()
	}
	return obj
}
