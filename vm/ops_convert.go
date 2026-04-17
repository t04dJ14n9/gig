// ops_convert.go handles type assertion, conversion, and change-type operations.
package vm

import (
	"fmt"
	"go/types"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
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

		// Special case: any value can be asserted to interface{}
		// This handles nested interface assertions like: outer.(interface{})
		if _, isInterface := targetType.(*types.Interface); isInterface {
			// Any value can be asserted to interface{}, including nil
			result = obj
			assertionOk = true
			v.pushCommaOk(result, assertionOk)
			break
		}

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
						// Failed assertion — return zero value of target type, not nil
						result = value.MakeFromReflect(reflect.Zero(targetReflectType))
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
				targetReflectType := typeToReflect(targetType, v.program)
				if targetReflectType != nil {
					result = value.MakeFromReflect(reflect.Zero(targetReflectType))
				} else {
					result = value.MakeNil()
				}
				assertionOk = false
			}
		} else {
			// For primitive kinds (KindInt, KindString, KindBool, KindFloat, etc.),
			// check whether the concrete value kind actually matches the target type.
			// This is critical for type switches on interface slice elements.
			assertionOk = kindMatchesType(obj.Kind(), obj.RawSize(), targetType)
			if assertionOk {
				result = obj
			} else {
				// Failed assertion — return zero value of target type, not nil
				targetReflectType := typeToReflect(targetType, v.program)
				if targetReflectType != nil {
					result = value.MakeFromReflect(reflect.Zero(targetReflectType))
				} else {
					result = value.MakeNil()
				}
			}
		}

		// Push result as a tuple [result, ok]
		v.pushCommaOk(result, assertionOk)

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
					if srcLocalIdx != noSourceLocalSentinel && rv.Kind() == reflect.Slice {
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

	case bytecode.OpMakeInterface:
		// Wraps a concrete value in a Go interface, preserving type information.
		// This is critical for typed nil: var p *T = nil; var e error = p → e is non-nil.
		ifaceTypeIdx := frame.readUint16()
		concreteTypeIdx := frame.readUint16()
		targetType := v.program.Types[ifaceTypeIdx]
		concreteType := v.program.Types[concreteTypeIdx]
		val := v.pop()

		// Get the interface's reflect.Type
		ifaceRT := typeToReflect(targetType, v.program)
		if ifaceRT == nil {
			// Can't determine interface type, pass through
			v.push(val)
			break
		}

		// If typeToReflect returned interface{} (any) but the target is a
		// non-empty interface (e.g., io.Reader, error), try to look up the
		// real type via the external type registry. typeToReflect converts
		// all *types.Interface to interface{} which loses method information.
		if ifaceRT.NumMethod() == 0 {
			if iface, ok := targetType.(*types.Interface); ok && iface.NumMethods() > 0 {
				// This is a real Go interface with methods, but typeToReflect
				// returned interface{}. Try to get the real type from external registry.
				if named, ok2 := targetType.(*types.Named); ok2 {
					if rt, ok3 := v.program.TypeResolver.LookupExternalType(named); ok3 {
						ifaceRT = rt
					}
				}
			}
		}

		// If we still have interface{} for a non-empty interface, just pass through.
		// The concrete value is already assignable to the target.
		if ifaceRT.NumMethod() == 0 {
			if iface, ok := targetType.(*types.Interface); ok && iface.NumMethods() > 0 {
				v.push(val)
				break
			}
		}

		// For empty interface (interface{}), wrapping is only needed for typed nil values.
		// For non-nil concrete values, the existing value is already correct.
		if ifaceRT.NumMethod() == 0 && !val.IsNil() && val.IsValid() {
			v.push(val)
			break
		}

		// Get the concrete value as a reflect.Value
		var concreteRV reflect.Value
		if val.IsNil() || !val.IsValid() {
			// The value is nil, but we need a TYPED nil to create a non-nil interface.
			// Use the concrete type from the MakeInterface instruction to create
			// a properly typed nil reflect.Value.
			concreteRT := typeToReflect(concreteType, v.program)
			if concreteRT != nil {
				concreteRV = reflect.Zero(concreteRT) // typed nil for pointers, nil slices, etc.
			} else {
				// Can't determine concrete type — the interface will be nil
				v.push(val)
				break
			}
		} else {
			concreteRV = val.ToReflectValue(ifaceRT)
		}

		// Create a new interface value and set it to the concrete value.
		// This properly creates a (type, value) pair so that even a typed nil
		// pointer results in a non-nil interface.
		ifaceVal := reflect.New(ifaceRT).Elem()
		ifaceVal.Set(concreteRV)
		v.push(value.MakeFromReflect(ifaceVal))
	}

	return nil
}
