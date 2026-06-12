// ops_convert.go handles type assertion, conversion, and change-type operations.
package vm

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// executeConvert handles type assertion, conversion, and change-type opcodes.
func (v *vm) executeConvert(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	case bytecode.OpMakeInterface:
		v.executeMakeInterface(frame)

	case bytecode.OpAssert:
		v.executeAssert(frame)

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
	}

	return nil
}

type makeInterfaceOperands struct {
	targetType   types.Type
	concreteType types.Type
	val          value.Value
}

func (v *vm) executeMakeInterface(frame *Frame) {
	op := v.readMakeInterfaceOperands(frame)

	if v.pushMakeInterfaceAdapter(op) {
		return
	}
	if dyn, ok := v.makeInterpretedInterfaceValue(op); ok {
		v.push(dyn)
		return
	}
	if isEmptyInterfaceType(op.targetType) && !needsTypedNilInterfaceValue(op) {
		v.push(op.val)
		return
	}
	if boxed, ok := v.makeReflectInterfaceValue(op); ok {
		v.push(boxed)
		return
	}
	v.push(op.val)
}

func (v *vm) pushMakeInterfaceAdapter(op makeInterfaceOperands) bool {
	adapter, ok := v.makeInterpretedInterfaceAdapter(op.targetType, op.concreteType, op.val)
	if !ok {
		return false
	}
	v.push(value.MakeFromReflect(reflect.ValueOf(adapter)))
	return true
}

func (v *vm) readMakeInterfaceOperands(frame *Frame) makeInterfaceOperands {
	ifaceTypeIdx := frame.readUint16()
	concreteTypeIdx := frame.readUint16()
	return makeInterfaceOperands{
		targetType:   v.program.Types[ifaceTypeIdx],
		concreteType: v.program.Types[concreteTypeIdx],
		val:          v.pop(),
	}
}

func (v *vm) makeInterpretedInterfaceValue(op makeInterfaceOperands) (value.Value, bool) {
	iface, ok := op.targetType.Underlying().(*types.Interface)
	if !ok {
		return value.Value{}, false
	}
	typeName := namedTypeName(op.concreteType)
	if typeName == "" || !shouldPreserveInterpretedNamedType(op.concreteType, v.program) {
		return value.Value{}, false
	}

	dyn := &value.InterpretedInterfaceValue{
		Value:     v.interfaceValueForDynamicType(op),
		TypeName:  typeName,
		IsPointer: isPointerType(op.concreteType),
	}
	if iface.NumMethods() > 0 && !v.interpretedTypeSatisfiesInterface(dyn, iface) {
		return value.Value{}, false
	}
	return value.MakeInterpretedInterface(dyn.Value, dyn.TypeName, dyn.IsPointer), true
}

func (v *vm) interfaceValueForDynamicType(op makeInterfaceOperands) value.Value {
	if !needsTypedNilInterfaceValue(op) {
		return op.val
	}
	concreteRT := typeToReflect(op.concreteType, v.program)
	if concreteRT == nil {
		return op.val
	}
	return value.MakeFromReflect(reflect.Zero(concreteRT))
}

func isEmptyInterfaceType(t types.Type) bool {
	iface, ok := t.Underlying().(*types.Interface)
	return ok && iface.NumMethods() == 0
}

func needsTypedNilInterfaceValue(op makeInterfaceOperands) bool {
	return (op.val.IsNil() || !op.val.IsValid()) && isPointerType(op.concreteType)
}

func (v *vm) makeReflectInterfaceValue(op makeInterfaceOperands) (value.Value, bool) {
	if op.val.IsNil() || !op.val.IsValid() {
		return v.makeTypedNilInterfaceValue(op)
	}
	concreteRT := typeToReflect(op.concreteType, v.program)
	if concreteRT == nil {
		return value.Value{}, false
	}
	rv := op.val.ToReflectValue(concreteRT)
	if !rv.IsValid() {
		return value.Value{}, false
	}
	return value.MakeFromReflect(rv), true
}

func (v *vm) makeTypedNilInterfaceValue(op makeInterfaceOperands) (value.Value, bool) {
	concreteRT := typeToReflect(op.concreteType, v.program)
	if concreteRT == nil {
		return value.Value{}, false
	}
	return value.MakeFromReflect(reflect.Zero(concreteRT)), true
}

func (v *vm) executeAssert(frame *Frame) {
	typeIdx := frame.readUint16()
	targetType := v.program.Types[typeIdx]
	obj := v.pop()

	if iface, ok := targetType.Underlying().(*types.Interface); ok && iface.NumMethods() == 0 {
		v.pushCommaOk(obj, true)
		return
	}
	if dyn, ok := obj.InterpretedInterface(); ok {
		result, assertionOk := v.assertInterpretedInterfaceValue(dyn, targetType, obj)
		v.pushCommaOk(result, assertionOk)
		return
	}
	if obj.Kind() == value.KindInterface {
		v.assertReflectInterfaceValue(obj, targetType)
		return
	}

	var result value.Value
	var assertionOk bool
	if obj.Kind() == value.KindReflect {
		result, assertionOk = v.assertReflectValue(obj, targetType)
	} else {
		result, assertionOk = v.assertPrimitiveValue(obj, targetType)
	}
	v.pushCommaOk(result, assertionOk)
}

func (v *vm) assertReflectInterfaceValue(obj value.Value, targetType types.Type) {
	rv, isReflect := obj.ReflectValue()
	if !isReflect || rv.Kind() != reflect.Interface {
		v.pushCommaOk(obj, true)
		return
	}
	if rv.IsNil() {
		v.pushCommaOk(value.MakeNil(), false)
		return
	}

	underlying := rv.Elem()
	targetReflectType := typeToReflect(targetType, v.program)
	if targetReflectType == nil {
		v.pushCommaOk(value.MakeNil(), false)
		return
	}
	if underlying.Type().AssignableTo(targetReflectType) {
		v.pushCommaOk(value.MakeFromReflect(underlying), true)
		return
	}
	v.pushCommaOk(value.MakeFromReflect(reflect.Zero(targetReflectType)), false)
}

func (v *vm) assertReflectValue(obj value.Value, targetType types.Type) (value.Value, bool) {
	rv, ok := obj.ReflectValue()
	if !ok {
		return zeroValueForType(targetType, v.program), false
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	targetReflectType := typeToReflect(targetType, v.program)
	if targetReflectType == nil {
		return value.MakeNil(), false
	}
	if rv.Type().AssignableTo(targetReflectType) {
		return value.MakeFromReflect(rv), true
	}
	if sameReflectKindFamily(rv.Type(), targetReflectType) {
		return value.MakeFromReflect(rv.Convert(targetReflectType)), true
	}
	return value.MakeNil(), false
}

func (v *vm) assertPrimitiveValue(obj value.Value, targetType types.Type) (value.Value, bool) {
	if kindMatchesType(obj.Kind(), obj.RawSize(), targetType) {
		return obj, true
	}
	return zeroValueForType(targetType, v.program), false
}

func (v *vm) assertInterpretedInterfaceValue(dyn *value.InterpretedInterfaceValue, targetType types.Type, original value.Value) (value.Value, bool) {
	if iface, ok := targetType.Underlying().(*types.Interface); ok {
		if iface.NumMethods() == 0 || v.interpretedTypeSatisfiesInterface(dyn, iface) {
			return original, true
		}
		return zeroValueForType(targetType, v.program), false
	}
	if targetName := namedTypeName(targetType); targetName == dyn.TypeName && isPointerType(targetType) == dyn.IsPointer {
		if dyn.IsPointer || kindMatchesType(dyn.Value.Kind(), dyn.Value.RawSize(), targetType) {
			return dyn.Value, true
		}
	}
	return zeroValueForType(targetType, v.program), false
}

func (v *vm) interpretedTypeSatisfiesInterface(dyn *value.InterpretedInterfaceValue, iface *types.Interface) bool {
	if dyn == nil || iface == nil || v == nil || v.program == nil {
		return false
	}
	for i := 0; i < iface.NumMethods(); i++ {
		if !v.interpretedMethodSatisfiesInterface(dyn, iface.Method(i).Name()) {
			return false
		}
	}
	return true
}

func (v *vm) interpretedMethodSatisfiesInterface(dyn *value.InterpretedInterfaceValue, methodName string) bool {
	for _, fn := range v.program.MethodsByName[methodName] {
		if fn.ReceiverTypeName == dyn.TypeName && (!fn.ReceiverIsPointer || dyn.IsPointer) {
			return true
		}
	}
	return false
}

func zeroValueForType(t types.Type, program *bytecode.CompiledProgram) value.Value {
	if rt := typeToReflect(t, program); rt != nil {
		return value.MakeFromReflect(reflect.Zero(rt))
	}
	return value.MakeNil()
}

func shouldPreserveInterpretedNamedType(t types.Type, program *bytecode.CompiledProgram) bool {
	typeName := namedTypeName(t)
	if typeName == "" {
		return false
	}
	rt := typeToReflect(t, program)
	return rt == nil || rt.Name() != typeName
}

func isPointerType(t types.Type) bool {
	_, ok := t.(*types.Pointer)
	return ok
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
