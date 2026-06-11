package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

type makeInterfaceOperands struct {
	targetType   types.Type
	concreteType types.Type
	val          value.Value
}

func (v *vm) executeMakeInterface(frame *Frame) {
	// Wraps a concrete value in a Go interface, preserving type information.
	// This is critical for typed nil: var p *T = nil; var e error = p → e is non-nil.
	op := v.readMakeInterfaceOperands(frame)

	if v.pushMakeInterfaceAdapter(op) {
		return
	}
	if v.pushPreservedNonEmptyInterpretedInterface(op) {
		return
	}

	ifaceRT := typeToReflect(op.targetType, v.program)
	if ifaceRT == nil {
		v.push(op.val)
		return
	}
	if v.pushMakeInterfacePassThrough(op, ifaceRT) {
		return
	}

	concreteRV, ok := v.makeInterfaceConcreteReflectValue(op, ifaceRT)
	if !ok {
		v.push(op.val)
		return
	}
	v.pushReflectInterfaceValue(ifaceRT, concreteRV)
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

func (v *vm) pushMakeInterfaceAdapter(op makeInterfaceOperands) bool {
	adapter, ok := v.makeInterpretedInterfaceAdapter(op.targetType, op.concreteType, op.val)
	if !ok {
		return false
	}
	v.push(value.MakeFromReflect(reflect.ValueOf(adapter)))
	return true
}

func (v *vm) pushPreservedNonEmptyInterpretedInterface(op makeInterfaceOperands) bool {
	iface, ok := op.targetType.Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return false
	}
	dyn, ok := v.interpretedInterfacePayload(op, iface)
	if !ok {
		return false
	}
	v.push(value.MakeInterpretedInterface(dyn.Value, dyn.TypeName, dyn.IsPointer))
	return true
}

func (v *vm) interpretedInterfacePayload(
	op makeInterfaceOperands,
	iface *types.Interface,
) (*value.InterpretedInterfaceValue, bool) {
	typeName := namedTypeName(op.concreteType)
	if typeName == "" || !shouldPreserveInterpretedNamedType(op.concreteType, v.program) {
		return nil, false
	}

	dyn := &value.InterpretedInterfaceValue{
		Value:     v.interfaceValueForDynamicType(op),
		TypeName:  typeName,
		IsPointer: isPointerType(op.concreteType),
	}
	if !v.interpretedTypeSatisfiesInterface(dyn, iface) {
		return nil, false
	}
	return dyn, true
}

func (v *vm) interfaceValueForDynamicType(op makeInterfaceOperands) value.Value {
	if !needsTypedNilInterfaceValue(op.val, op.concreteType) {
		return op.val
	}
	concreteRT := typeToReflect(op.concreteType, v.program)
	if concreteRT == nil {
		return op.val
	}
	return value.MakeFromReflect(reflect.Zero(concreteRT))
}

func needsTypedNilInterfaceValue(val value.Value, concreteType types.Type) bool {
	return (val.IsNil() || !val.IsValid()) && isPointerType(concreteType)
}

func (v *vm) pushMakeInterfacePassThrough(op makeInterfaceOperands, ifaceRT reflect.Type) bool {
	if unresolvedDirectNonEmptyInterface(op.targetType, ifaceRT) {
		v.push(op.val)
		return true
	}
	if !emptyInterfaceNeedsNoWrap(op.val, ifaceRT) {
		return false
	}
	if v.pushPreservedEmptyInterface(op) {
		return true
	}
	v.push(op.val)
	return true
}

func unresolvedDirectNonEmptyInterface(targetType types.Type, ifaceRT reflect.Type) bool {
	if ifaceRT.NumMethod() != 0 {
		return false
	}
	iface, ok := targetType.(*types.Interface)
	return ok && iface.NumMethods() > 0
}

func emptyInterfaceNeedsNoWrap(val value.Value, ifaceRT reflect.Type) bool {
	return ifaceRT.NumMethod() == 0 && !val.IsNil() && val.IsValid()
}

func (v *vm) pushPreservedEmptyInterface(op makeInterfaceOperands) bool {
	typeName := namedTypeName(op.concreteType)
	if typeName == "" || !shouldPreserveInterpretedNamedType(op.concreteType, v.program) {
		return false
	}
	v.push(value.MakeInterpretedInterface(op.val, typeName, isPointerType(op.concreteType)))
	return true
}

func (v *vm) makeInterfaceConcreteReflectValue(
	op makeInterfaceOperands,
	ifaceRT reflect.Type,
) (reflect.Value, bool) {
	if op.val.IsNil() || !op.val.IsValid() {
		return v.typedNilInterfaceReflectValue(op.concreteType)
	}
	return op.val.ToReflectValue(ifaceRT), true
}

func (v *vm) typedNilInterfaceReflectValue(concreteType types.Type) (reflect.Value, bool) {
	// A nil concrete value needs the instruction's static type to produce a
	// non-nil interface pair, for example var p *T = nil; var e error = p.
	concreteRT := typeToReflect(concreteType, v.program)
	if concreteRT == nil {
		return reflect.Value{}, false
	}
	return reflect.Zero(concreteRT), true
}

func (v *vm) pushReflectInterfaceValue(ifaceRT reflect.Type, concreteRV reflect.Value) {
	ifaceVal := reflect.New(ifaceRT).Elem()
	ifaceVal.Set(concreteRV)
	v.push(value.MakeFromReflect(ifaceVal))
}
