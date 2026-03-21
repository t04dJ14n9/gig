package vm

import (
	"fmt"
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// executeConvert handles type assertion, conversion, and change-type opcodes.
func (v *vm) executeConvert(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	case bytecode.OpAssert:
		typeIdx := frame.readUint16()
		targetType := v.program.Types[typeIdx]
		obj := v.pop()

		// Type assertion - check if obj can be asserted to targetType
		// Returns (value, ok) tuple on stack
		var result value.Value
		var assertionOk bool

		if obj.Kind() == value.KindInterface {
			// Get the underlying interface
			if rv, isReflect := obj.ReflectValue(); isReflect && rv.Kind() == reflect.Interface {
				if rv.IsNil() {
					// Interface is nil, assertion fails
					result = value.MakeNil()
					assertionOk = false
				} else {
					// Get the underlying value
					underlying := rv.Elem()
					targetReflectType := typeToReflect(targetType, v.program)
					if targetReflectType == nil {
						result = value.MakeNil()
						assertionOk = false
					} else if underlying.Type().AssignableTo(targetReflectType) {
						// Successful assertion - create a value from the underlying
						result = value.MakeFromReflect(underlying)
						assertionOk = true
					} else {
						result = value.MakeNil()
						assertionOk = false
					}
				}
			} else {
				result = obj
				assertionOk = true
			}
		} else if obj.Kind() == value.KindReflect {
			// Already a reflect value — may be an interface wrapping a concrete type
			if rv, isReflect := obj.ReflectValue(); isReflect {
				// If the reflect value is an interface, unwrap it to the concrete type
				concreteRV := rv
				if rv.Kind() == reflect.Interface && !rv.IsNil() {
					concreteRV = rv.Elem()
				}
				targetReflectType := typeToReflect(targetType, v.program)
				if targetReflectType != nil && concreteRV.Type().AssignableTo(targetReflectType) {
					result = value.MakeFromReflect(concreteRV)
					assertionOk = true
				} else if targetReflectType != nil && sameReflectKindFamily(concreteRV.Type(), targetReflectType) {
					// Gig stores numeric values with internal types (int64 for int, float64 for float32, etc.).
					// When these are stored in interface{}, the concrete type may differ from what Go
					// would use (e.g., int64 instead of int). For type switch correctness, we match
					// by reflect.Kind family and convert to the target type.
					result = value.MakeFromReflect(concreteRV.Convert(targetReflectType))
					assertionOk = true
				} else {
					result = value.MakeNil()
					assertionOk = false
				}
			} else {
				result = obj
				assertionOk = false
			}
		} else {
			// For primitive kinds (KindInt, KindString, KindBool, KindFloat, etc.),
			// check whether the concrete value kind actually matches the target type.
			// This is critical for type switches on interface slice elements.
			assertionOk = kindMatchesType(obj.Kind(), targetType)
			if assertionOk {
				result = obj
			} else {
				result = value.MakeNil()
			}
		}

		// Push result as a tuple [result, ok]
		// Use a slice to represent the tuple
		tuple := []value.Value{result, value.MakeBool(assertionOk)}
		v.push(value.FromInterface(tuple))

	case bytecode.OpConvert:
		typeIdx := frame.readUint16()
		targetType := v.program.Types[typeIdx]
		val := v.pop()

		// Handle type conversion
		switch t := targetType.(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.String:
				// Convert to string
				switch val.Kind() {
				case value.KindInt:
					// int -> string: convert rune to string
					v.push(value.MakeString(string(rune(val.Int()))))
				case value.KindUint:
					// byte/uint8 -> string: convert byte to string
					v.push(value.MakeString(string(byte(val.Uint()))))
				case value.KindString:
					v.push(val)
				case value.KindBytes:
					if b, ok := val.Bytes(); ok {
						v.push(value.MakeString(string(b)))
					} else {
						v.push(value.MakeString(""))
					}
				case value.KindReflect:
					// Handle string([]rune) / string([]int32)
					if rv, ok := val.ReflectValue(); ok && rv.Kind() == reflect.Slice {
						elemKind := rv.Type().Elem().Kind()
						if elemKind == reflect.Int32 {
							runes := make([]rune, rv.Len())
							for i := 0; i < rv.Len(); i++ {
								runes[i] = rune(rv.Index(i).Int())
							}
							v.push(value.MakeString(string(runes)))
						} else {
							v.push(value.MakeString(fmt.Sprintf("%v", val.Interface())))
						}
					} else {
						v.push(value.MakeString(fmt.Sprintf("%v", val.Interface())))
					}
				default:
					// Use reflection for other types
					v.push(value.MakeString(fmt.Sprintf("%v", val.Interface())))
				}
			case types.Int:
				v.push(value.MakeInt(toInt64(val)))
			case types.Int8:
				v.push(value.MakeInt8(int8(toInt64(val))))
			case types.Int16:
				v.push(value.MakeInt16(int16(toInt64(val))))
			case types.Int32:
				v.push(value.MakeInt32(int32(toInt64(val))))
			case types.Int64:
				v.push(value.MakeInt64(toInt64(val)))
			case types.Uint:
				v.push(value.MakeUint(toUint64(val)))
			case types.Uint8:
				v.push(value.MakeUint8(uint8(toUint64(val))))
			case types.Uint16:
				v.push(value.MakeUint16(uint16(toUint64(val))))
			case types.Uint32:
				v.push(value.MakeUint32(uint32(toUint64(val))))
			case types.Uint64:
				v.push(value.MakeUint64(toUint64(val)))
			case types.Uintptr:
				v.push(value.MakeUint64(toUint64(val)))
			case types.Float32:
				v.push(value.MakeFloat32(float32(toFloat64(val))))
			case types.Float64:
				v.push(value.MakeFloat(toFloat64(val)))
			default:
				v.push(val)
			}
		case *types.Slice:
			// Handle string -> []rune or string -> []byte conversions
			elem := t.Elem()
			if basic, ok := elem.(*types.Basic); ok {
				switch basic.Kind() {
				case types.Int32: // []rune(string) or []int32(string)
					if val.Kind() == value.KindString {
						runes := []rune(val.String())
						rs := reflect.MakeSlice(reflect.TypeOf([]int32{}), len(runes), len(runes))
						for i, r := range runes {
							rs.Index(i).SetInt(int64(r))
						}
						v.push(value.MakeFromReflect(rs))
					} else {
						v.push(val)
					}
				case types.Uint8:
					if val.Kind() == value.KindString {
						// Use make+copy to ensure cap==len, matching native Go compiler behavior.
						// Direct []byte(s) uses runtime's stringtoslicebyte which rounds cap up
						// to allocator size classes (e.g. 5→8), but the native compiler optimizes
						// to exact capacity via stack allocation or constant folding.
						s := val.String()
						b := make([]byte, len(s))
						copy(b, s)
						v.push(value.MakeBytes(b))
					} else {
						v.push(val)
					}
				default:
					v.push(val)
				}
			} else {
				v.push(val)
			}
		case *types.Named:
			// Named-type conversion (e.g., []int -> sort.IntSlice, []string -> sort.StringSlice).
			// Resolve the target reflect.Type via external type registry, then convert via reflect.
			targetRT := typeToReflect(t, v.program)
			if targetRT != nil {
				// Get a reflect.Value of the source. Use the target's underlying type
				// so ToReflectValue can do element-type conversion (e.g., []int64 -> []int).
				rv := val.ToReflectValue(targetRT)
				if rv.IsValid() {
					// If ToReflectValue gave us the underlying type (e.g., []float64 instead
					// of sort.Float64Slice), convert to the actual named type.
					if rv.Type() != targetRT && rv.Type().ConvertibleTo(targetRT) {
						rv = rv.Convert(targetRT)
					}
					v.push(value.MakeFromReflect(rv))
				} else {
					v.push(val)
				}
			} else {
				v.push(val)
			}
		default:
			// For non-basic types, just pass through for now
			v.push(val)
		}

	// Function operations
	case bytecode.OpChangeType:
		typeIdx := frame.readUint16()
		srcLocalIdx := frame.readUint16()
		targetType := v.program.Types[typeIdx]
		val := v.pop()

		// Named-type conversion (e.g., []int -> sort.IntSlice).
		if named, ok := targetType.(*types.Named); ok {
			targetRT := typeToReflect(named, v.program)
			if targetRT != nil {
				// Get a reflect.Value, using the target type so ToReflectValue
				// handles element-type conversion (e.g., []int64 -> []int).
				rv := val.ToReflectValue(targetRT)
				if rv.IsValid() {
					// If ToReflectValue returned the underlying type, Convert to the named type.
					if rv.Type() != targetRT && rv.Type().ConvertibleTo(targetRT) {
						rv = rv.Convert(targetRT)
					}
					// For slices: update the source local to share the same backing array.
					// This ensures that sort.IntSlice(s) and s refer to the same data,
					// matching Go's semantics where ChangeType on slices shares memory.
					if srcLocalIdx != 0xFFFF && rv.Kind() == reflect.Slice {
						if int(srcLocalIdx) < len(frame.locals) {
							// Create a view of the same backing array as the underlying slice type.
							// e.g., for sort.IntSlice -> create a []int sharing the same backing.
							underlyingRV := rv.Convert(reflect.SliceOf(rv.Type().Elem()))
							frame.locals[srcLocalIdx] = value.MakeFromReflect(underlyingRV)
						}
					}
					v.push(value.MakeFromReflect(rv))
				} else {
					v.push(val)
				}
			} else {
				// Named type not in external registry (interpreted type) — pass through.
				v.push(val)
			}
		} else {
			// Not a named type, fall back to simple pass-through.
			v.push(val)
		}
	}

	return nil
}
