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
	return v.str
}

// Complex returns the complex value. Panics if not KindComplex.
func (v Value) Complex() complex128 {
	if v.kind != KindComplex {
		panic(fmt.Sprintf("not a complex: %v", v.kind))
	}
	return complex(math.Float64frombits(uint64(v.num)), math.Float64frombits(uint64(v.num2)))
}

// Interface returns the value as an interface{}.
func (v Value) Interface() any {
	switch v.kind {
	case KindNil:
		return nil
	case KindBool:
		return v.Bool()
	case KindInt:
		return v.Int()
	case KindUint:
		return v.Uint()
	case KindFloat:
		return v.Float()
	case KindString:
		return v.str
	case KindComplex:
		return v.Complex()
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

// ToReflectValue converts to reflect.Value.
func (v Value) ToReflectValue(typ reflect.Type) reflect.Value {
	switch v.kind {
	case KindNil:
		return reflect.Zero(typ)
	case KindBool:
		return reflect.ValueOf(v.Bool())
	case KindInt:
		return reflect.ValueOf(v.num).Convert(typ)
	case KindUint:
		return reflect.ValueOf(uint64(v.num)).Convert(typ)
	case KindFloat:
		return reflect.ValueOf(v.Float()).Convert(typ)
	case KindString:
		return reflect.ValueOf(v.str)
	case KindComplex:
		return reflect.ValueOf(v.Complex())
	case KindSlice:
		// Native int slice → target type conversion
		if s, ok := v.obj.([]int64); ok && typ.Kind() == reflect.Slice {
			target := reflect.MakeSlice(typ, len(s), cap(s))
			for i, n := range s {
				target.Index(i).SetInt(n)
			}
			return target
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	default:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	}
}

// ReflectValue returns the internal reflect.Value if stored.
func (v Value) ReflectValue() (reflect.Value, bool) {
	rv, ok := v.obj.(reflect.Value)
	return rv, ok
}
