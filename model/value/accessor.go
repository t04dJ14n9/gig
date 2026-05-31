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
	switch v.kind {
	case KindNil:
		return nil
	case KindBool:
		return v.Bool()
	case KindInt:
		switch v.size {
		case Size8:
			return int8(v.num)
		case Size16:
			return int16(v.num)
		case Size32:
			return int32(v.num)
		case Size64:
			return v.num // int64
		default:
			return int(v.num) // SizePtr / Size0 → int
		}
	case KindUint:
		switch v.size {
		case Size8:
			return uint8(v.num)
		case Size16:
			return uint16(v.num)
		case Size32:
			return uint32(v.num)
		case Size64:
			return uint64(v.num)
		default:
			return uint(v.num) // SizePtr / Size0 → uint
		}
	case KindFloat:
		f := math.Float64frombits(uint64(v.num))
		if v.size == Size32 {
			return float32(f)
		}
		return f // float64
	case KindString:
		return v.obj.(string)
	case KindComplex:
		c := v.obj.(complex128)
		if v.size == Size32 {
			return complex64(c)
		}
		return c
	case KindInterface:
		if dyn, ok := v.InterpretedInterface(); ok {
			return dyn.Value.Interface()
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Interface()
		}
		return v.obj
	case KindFunc:
		return v.obj
	case KindBytes:
		return v.obj.([]byte)
	case KindSlice:
		// Native int slice: convert []int64 to []int for Go-compatible return
		if s, ok := v.obj.([]int64); ok {
			result := make([]int, len(s))
			for i, n := range s {
				result[i] = int(n)
			}
			return result
		}
		return v.obj
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Interface()
		}
		return v.obj
	default:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Interface()
		}
		return v.obj
	}
}
