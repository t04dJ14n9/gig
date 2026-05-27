// ops_convert.go handles type assertion, conversion, and change-type operations.
package vm

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// interfaceMethodNames returns the method names of an interface type.
func interfaceMethodNames(iface *types.Interface) []string {
	names := make([]string, iface.NumMethods())
	for i := 0; i < iface.NumMethods(); i++ {
		names[i] = iface.Method(i).Name()
	}
	return names
}

// typeMatchesKind checks if a compiled receiver type name corresponds to a
// named type whose underlying kind matches the given value kind. This prevents
// false positives when checking if a primitive value implements an interface.
func typeMatchesKind(receiverTypeName string, k value.Kind, prog *bytecode.CompiledProgram) bool {
	for _, typ := range prog.Types {
		if named, ok := typ.(*types.Named); ok {
			if named.Obj().Name() == receiverTypeName {
				underlying := named.Underlying()
				if basic, ok := underlying.(*types.Basic); ok {
					switch basic.Kind() {
					case types.String:
						return k == value.KindString
					case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
						return k == value.KindInt
					case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
						return k == value.KindUint
					case types.Float32, types.Float64:
						return k == value.KindFloat
					case types.Bool:
						return k == value.KindBool
					}
				}
				return false
			}
		}
	}
	return false
}

// concreteImplementsInterface checks if a concrete reflect.Value implements an
// interface by checking if the program's compiled methods contain all required
// interface methods for that type name. Handles both named types and
// reflect.StructOf types (which have empty names but may have gig tags).
// For primitive types (string, int, etc.), checks if any compiled type with
// the same underlying kind implements the interface (handles named primitive
// types like MyString string).
func concreteImplementsInterface(rv reflect.Value, iface *types.Interface, prog *bytecode.CompiledProgram) bool {
	required := interfaceMethodNames(iface)
	if len(required) == 0 {
		return true
	}

	// Check if this is a gigStructWrapper — if so, use the original type name
	// to look up methods, since the wrapper only implements a few stdlib interfaces.
	if wrapperTypeName := gigStructWrapperTypeName(rv); wrapperTypeName != "" {
		shortName := wrapperTypeName
		if dotIdx := strings.LastIndex(wrapperTypeName, "."); dotIdx >= 0 {
			shortName = wrapperTypeName[dotIdx+1:]
		}
		for _, fn := range prog.MethodsByName[required[0]] {
			if fn.ReceiverTypeName != "" {
				if fn.ReceiverTypeName == wrapperTypeName || fn.ReceiverTypeName == shortName ||
					strings.HasSuffix(fn.ReceiverTypeName, "."+wrapperTypeName) {
					if implementsInterface(fn.ReceiverTypeName, required, prog) {
						return true
					}
				}
			}
		}
		return false
	}

	rt := rv.Type()
	typeName := resolveTypeName(rt, prog)
	if typeName != "" {
		if implementsInterface(typeName, required, prog) {
			return true
		}
		// For basic types (string, int, etc.), check if any compiled type with
		// the same underlying kind implements the interface. This handles named
		// primitive types like "type MyString string" which lose their name
		// when stored as KindString in the value system.
		if rt.Kind() == reflect.String || rt.Kind() == reflect.Int || rt.Kind() == reflect.Float64 ||
			rt.Kind() == reflect.Bool || rt.Kind() == reflect.Int32 || rt.Kind() == reflect.Uint8 {
			for _, fn := range prog.MethodsByName[required[0]] {
				if fn.ReceiverTypeName != "" && implementsInterface(fn.ReceiverTypeName, required, prog) {
					return true
				}
			}
		}
	}
	// For interpreter-defined structs (created via reflect.StructOf), methods aren't
	// on the reflect type. But the gig struct tag encodes the original type name.
	// Extract it and check compiled method tables.
	if gigName := gigStructTagName(rt); gigName != "" {
		// gigName is "pkg.Type" (short package name). ReceiverTypeName may be
		// "full/import/path.Type" or just "Type" (for main/test packages).
		// Try suffix match with ".TypeName" and bare "TypeName".
		shortName := gigName
		if dotIdx := strings.LastIndex(gigName, "."); dotIdx >= 0 {
			shortName = gigName[dotIdx+1:]
		}
		for _, fn := range prog.MethodsByName[required[0]] {
			if fn.ReceiverTypeName != "" {
				if fn.ReceiverTypeName == gigName || fn.ReceiverTypeName == shortName ||
					strings.HasSuffix(fn.ReceiverTypeName, "."+gigName) {
					if implementsInterface(fn.ReceiverTypeName, required, prog) {
						return true
					}
				}
			}
		}
	}

	// Fallback: use reflect to check if the concrete type implements the interface.
	// This handles internal types (e.g., *fmt.wrapError) whose type names aren't
	// registered in the program's method tables.
	// Only applies to struct/pointer types where reflect carries method info.
	// Named primitive types (type MyString string) lose their methods in reflect,
	// so they must go through the compiled method tables above.
	if rt.Kind() == reflect.Struct || rt.Kind() == reflect.Ptr {
		for _, methodName := range required {
			if _, ok := rt.MethodByName(methodName); !ok {
				if rt.Kind() != reflect.Ptr {
					if _, ok := reflect.PointerTo(rt).MethodByName(methodName); !ok {
						return false
					}
				} else {
					return false
				}
			}
		}
		return true
	}
	return false
}

// gigStructWrapperTypeName checks if a reflect.Value is a *gigStructWrapper
// and returns its typeName field value. Returns "" if not a wrapper.
func gigStructWrapperTypeName(rv reflect.Value) string {
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ""
	}
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return ""
	}
	// Check for the typeName field which is unique to gigStructWrapper
	field := elem.FieldByName("typeName")
	if !field.IsValid() || field.Kind() != reflect.String {
		return ""
	}
	return field.String()
}

// gigStructTagName extracts the type name from the gig struct tag on
// interpreter-defined struct types. Returns "" if not a gig struct.
func gigStructTagName(rt reflect.Type) string {
	// Unwrap pointer
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct || rt.NumField() == 0 {
		return ""
	}
	gigTag := rt.Field(0).Tag.Get("gig")
	if gigTag == "" {
		return ""
	}
	if strings.HasPrefix(gigTag, "#") {
		return gigTag[1:]
	}
	return gigTag
}

// implementsInterface checks if a type with the given name has all required methods.
func implementsInterface(typeName string, required []string, prog *bytecode.CompiledProgram) bool {
	for _, methodName := range required {
		found := false
		for _, fn := range prog.MethodsByName[methodName] {
			if fn.ReceiverTypeName == typeName {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

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

		// Special case: any value can be asserted to interface{} (empty interface).
		// This handles nested interface assertions like: outer.(interface{})
		// Use Underlying() to also handle *types.Named wrapping interface{} (e.g., "any").
		if iface, isInterface := targetType.Underlying().(*types.Interface); isInterface && iface.NumMethods() == 0 {
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
						// Check: typeToReflect returns interface{} for custom interfaces.
						// AssignableTo(interface{}) is always true, so we must verify the
						// concrete type actually implements the target interface's methods.
						if targetReflectType.NumMethod() == 0 {
							if iface, isIface := targetType.Underlying().(*types.Interface); isIface && iface.NumMethods() > 0 {
								assertionOk = concreteImplementsInterface(underlying, iface, v.program)
							} else {
								assertionOk = true
							}
						} else {
							assertionOk = true
						}
						if assertionOk {
							result = value.MakeFromReflect(underlying)
						} else {
							result = value.MakeFromReflect(reflect.Zero(targetReflectType))
						}
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
					// typeToReflect returns interface{} for custom interfaces.
					// AssignableTo(interface{}) is always true, so verify the concrete
					// type actually implements the target interface's methods.
					if targetReflectType.NumMethod() == 0 {
						if iface, isIface := targetType.Underlying().(*types.Interface); isIface && iface.NumMethods() > 0 {
							assertionOk = concreteImplementsInterface(concreteRV, iface, v.program)
						} else {
							assertionOk = true
						}
					} else {
						assertionOk = true
					}
					if assertionOk {
						result = value.MakeFromReflect(concreteRV)
					} else {
						result = value.MakeNil()
					}
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
				// If target is an interface, check if any compiled type with
				// matching methods implements the interface. This handles named
				// primitive types like "type MyString string" which lose their
				// type name when stored as KindString.
				if iface, isInterface := targetType.Underlying().(*types.Interface); isInterface && iface.NumMethods() > 0 {
					requiredMethods := interfaceMethodNames(iface)
					if len(requiredMethods) > 0 {
						for _, fn := range v.program.MethodsByName[requiredMethods[0]] {
							if fn.ReceiverTypeName != "" && implementsInterface(fn.ReceiverTypeName, requiredMethods, v.program) {
								if typeMatchesKind(fn.ReceiverTypeName, obj.Kind(), v.program) {
									assertionOk = true
									result = obj
									break
								}
							}
						}
					}
				}
				if !assertionOk {
					// Failed assertion — return zero value of target type, not nil
					targetReflectType := typeToReflect(targetType, v.program)
					if targetReflectType != nil {
						result = value.MakeFromReflect(reflect.Zero(targetReflectType))
					} else {
						result = value.MakeNil()
					}
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

		// Handle nil pointer dereference for value-receiver methods.
		// When calling v.Method() on a nil *T where Method has a value receiver,
		// Go creates a zero value of T. The SSA generates ChangeType(*T -> T).
		if val.IsNil() || !val.IsValid() {
			if named, ok := targetType.(*types.Named); ok {
				if targetRT := typeToReflect(named, v.program); targetRT != nil {
					v.push(value.MakeFromReflect(reflect.Zero(targetRT)))
					break
				}
			}
			v.push(val)
			break
		}
		_ = srcLocalIdx // used below for slice aliasing

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

		if adapter, ok := v.makeInterpretedInterfaceAdapter(targetType, concreteType, val); ok {
			v.push(value.MakeFromReflect(reflect.ValueOf(adapter)))
			break
		}

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

func (v *vm) makeInterpretedInterfaceAdapter(targetType, concreteType types.Type, receiver value.Value) (*interpretedInterfaceAdapter, bool) {
	if !isHostCallbackInterface(targetType) {
		return nil, false
	}
	receiverTypeName := namedTypeName(concreteType)
	if receiverTypeName == "" {
		return nil, false
	}
	return newInterpretedInterfaceAdapter(
		v.program, receiver, receiverTypeName,
		v.getGlobals(), v.initialGlobals, v.shared, v.ctx, v.goroutines,
	), true
}

// isHostCallbackInterface checks whether the target type is sort.Interface or
// container/heap.Interface — the two stdlib interfaces that receive callbacks
// and for which gig provides interpreted-to-native adapters.
//
// The target may arrive as *types.Named (when the compiler stores the named
// type) or as *types.Interface (when it stores the underlying interface).
// We handle both.
func isHostCallbackInterface(t types.Type) bool {
	// Fast path: *types.Named wrapping sort.Interface or heap.Interface.
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj == nil || obj.Pkg() == nil {
			return false
		}
		pkgPath := obj.Pkg().Path()
		return obj.Name() == "Interface" && (pkgPath == "sort" || pkgPath == "container/heap")
	}
	// Slow path: *types.Interface directly. Match by method signature.
	if iface, ok := t.(*types.Interface); ok {
		return hasInterfaceMethods(iface, "Len", "Less", "Swap")
	}
	return false
}

// hasInterfaceMethods checks that an interface type has methods with exactly
// the given names (order-independent).
func hasInterfaceMethods(iface *types.Interface, names ...string) bool {
	if iface.NumMethods() < len(names) {
		return false
	}
	for _, name := range names {
		found := false
		for i := 0; i < iface.NumMethods(); i++ {
			if iface.Method(i).Name() == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func namedTypeName(t types.Type) string {
	for {
		switch tt := t.(type) {
		case *types.Named:
			return tt.Obj().Name()
		case *types.Pointer:
			t = tt.Elem()
		default:
			return ""
		}
	}
}
