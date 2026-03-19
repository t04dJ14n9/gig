package value

import (
	"fmt"
	"math"
	"reflect"
)

// ClosureCaller is a callback registered by the vm package to execute a closure.
// It receives the raw closure object, reflect.Value arguments, and the expected
// output types. It returns reflect.Value results matching the target function signature.
// The outTypes parameter enables recursive wrapping of nested closures.
// This breaks the circular dependency between value and vm packages.
type ClosureCaller func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value

// closureCaller is the registered callback for executing closures.
// It is set by the vm package during initialization.
var closureCaller ClosureCaller //nolint:gochecknoglobals // cross-package callback, must be global

// RegisterClosureCaller registers the closure execution callback.
// This must be called by the vm package before any closure-to-func conversion occurs.
func RegisterClosureCaller(caller ClosureCaller) {
	closureCaller = caller
}

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
		return v.obj.(complex128)
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
		rv := reflect.ValueOf(v.obj.(string))
		if rv.Type() != typ {
			rv = rv.Convert(typ)
		}
		return rv
	case KindComplex:
		return reflect.ValueOf(v.obj.(complex128))
	case KindFunc:
		// If the target type is a function type, wrap the closure in a real Go function
		// using reflect.MakeFunc. This allows closures to be stored in typed containers
		// (maps, struct fields) that expect concrete function types like func() int.
		if typ.Kind() == reflect.Func && closureCaller != nil {
			closure := v.obj
			numOut := typ.NumOut()
			outTypes := make([]reflect.Type, numOut)
			for i := 0; i < numOut; i++ {
				outTypes[i] = typ.Out(i)
			}
			fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
				results := closureCaller(closure, args, outTypes)
				// Convert results to match the expected return types
				out := make([]reflect.Value, numOut)
				for i := 0; i < numOut; i++ {
					if i < len(results) && results[i].IsValid() {
						if results[i].Type().ConvertibleTo(outTypes[i]) {
							out[i] = results[i].Convert(outTypes[i])
						} else {
							out[i] = results[i]
						}
					} else {
						out[i] = reflect.Zero(outTypes[i])
					}
				}
				return out
			})
			return fn
		}
		return reflect.ValueOf(v.obj)
	case KindBytes:
		return reflect.ValueOf(v.obj.([]byte))
	case KindSlice:
		// Native int slice → target type conversion
		if s, ok := v.obj.([]int64); ok && typ.Kind() == reflect.Slice {
			target := reflect.MakeSlice(typ, len(s), cap(s))
			for i, n := range s {
				target.Index(i).SetInt(n)
			}
			return target
		}
		// []value.Value → typed slice conversion (e.g. []func() int)
		if s, ok := v.obj.([]Value); ok && typ.Kind() == reflect.Slice {
			target := reflect.MakeSlice(typ, len(s), cap(s))
			elemType := typ.Elem()
			for i, elem := range s {
				target.Index(i).Set(elem.ToReflectValue(elemType))
			}
			return target
		}
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			// Handle *func(...) target type: when the value is a *value.Value pointer
			// (created by OpNew for Signature types) containing a closure, convert
			// to a real Go function pointer via reflect.MakeFunc.
			if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Func {
				if rv.Kind() == reflect.Ptr && !rv.IsNil() {
					if vp, ok2 := rv.Interface().(*Value); ok2 {
						// Convert the inner closure to a real function
						funcRV := vp.ToReflectValue(typ.Elem())
						// Allocate a new pointer to hold the function
						ptr := reflect.New(typ.Elem())
						ptr.Elem().Set(funcRV)
						return ptr
					}
				}
			}
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
