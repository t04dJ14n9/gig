// accessor.go provides the method resolver registry and type accessors
// (Bool, Int, Float, String, Interface, ToReflectValue).
package value

import (
	"fmt"
	"math"
	"reflect"
	"sync"
)

// MethodResolverFunc is a callback for calling compiled methods on interpreted types.
// It receives a method name and receiver value, and returns the result if found.
type MethodResolverFunc func(methodName string, receiver Value) (Value, bool)

// methodResolverRegistry is a thread-safe global registry of per-program method resolvers.
// This allows fmt DirectCall wrappers (which lack VM context) to resolve compiled methods
// on interpreted types. Each program registers its resolver on creation and unregisters
// on cleanup. Using sync.Map eliminates the data race that the old single-global approach had.
var methodResolverRegistry sync.Map // map[uintptr]MethodResolverFunc

// RegisterMethodResolver registers a method resolver for a program identified by key.
// The key should be a unique identifier per program (e.g., uintptr of program pointer).
func RegisterMethodResolver(key uintptr, resolver MethodResolverFunc) {
	methodResolverRegistry.Store(key, resolver)
}

// UnregisterMethodResolver removes a method resolver for the given program key.
func UnregisterMethodResolver(key uintptr) {
	methodResolverRegistry.Delete(key)
}

// CallMethod attempts to call a compiled method on the receiver using the given resolver.
// If resolver is nil, it falls back to searching all registered per-program resolvers.
// Returns (result, true) if the method was found and called, or (zero, false) otherwise.
func CallMethod(resolver MethodResolverFunc, methodName string, receiver Value) (Value, bool) {
	if resolver != nil {
		return resolver(methodName, receiver)
	}
	// Fallback: try all registered resolvers (for fmt DirectCall wrappers lacking VM context)
	var result Value
	var found bool
	methodResolverRegistry.Range(func(_, v any) bool {
		if r, ok := v.(MethodResolverFunc); ok {
			result, found = r(methodName, receiver)
			if found {
				return false // stop iteration
			}
		}
		return true // continue
	})
	if found {
		return result, true
	}
	return MakeNil(), false
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
		// Create the correct integer type based on size
		var intRV reflect.Value
		switch v.size {
		case Size8:
			intRV = reflect.ValueOf(int8(v.num))
		case Size16:
			intRV = reflect.ValueOf(int16(v.num))
		case Size32:
			intRV = reflect.ValueOf(int32(v.num))
		case Size64:
			intRV = reflect.ValueOf(v.num) // int64
		default:
			intRV = reflect.ValueOf(int(v.num)) // SizePtr / Size0 → int
		}
		if intRV.Type().ConvertibleTo(typ) {
			return intRV.Convert(typ)
		}
		return intRV
	case KindUint:
		// Create the correct unsigned integer type based on size
		var uintRV reflect.Value
		switch v.size {
		case Size8:
			uintRV = reflect.ValueOf(uint8(v.num))
		case Size16:
			uintRV = reflect.ValueOf(uint16(v.num))
		case Size32:
			uintRV = reflect.ValueOf(uint32(v.num))
		case Size64:
			uintRV = reflect.ValueOf(uint64(v.num))
		default:
			uintRV = reflect.ValueOf(uint(v.num)) // SizePtr / Size0 → uint
		}
		if uintRV.Type().ConvertibleTo(typ) {
			return uintRV.Convert(typ)
		}
		return uintRV
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
		if typ.Kind() == reflect.Func {
			if ce, ok := v.obj.(ClosureExecutor); ok {
				numOut := typ.NumOut()
				outTypes := make([]reflect.Type, numOut)
				for i := 0; i < numOut; i++ {
					outTypes[i] = typ.Out(i)
				}
				fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
					results := ce.Execute(args, outTypes)
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
