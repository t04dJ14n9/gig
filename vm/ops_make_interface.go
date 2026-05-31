package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeMakeInterface(frame *Frame) {
	// Wraps a concrete value in a Go interface, preserving type information.
	// This is critical for typed nil: var p *T = nil; var e error = p → e is non-nil.
	ifaceTypeIdx := frame.readUint16()
	concreteTypeIdx := frame.readUint16()
	targetType := v.program.Types[ifaceTypeIdx]
	concreteType := v.program.Types[concreteTypeIdx]
	val := v.pop()

	if adapter, ok := v.makeInterpretedInterfaceAdapter(targetType, concreteType, val); ok {
		v.push(value.MakeFromReflect(reflect.ValueOf(adapter)))
		return
	}

	if iface, ok := targetType.Underlying().(*types.Interface); ok && iface.NumMethods() > 0 {
		if typeName := namedTypeName(concreteType); typeName != "" && shouldPreserveInterpretedNamedType(concreteType, v.program) {
			interfaceValue := val
			if (val.IsNil() || !val.IsValid()) && isPointerType(concreteType) {
				if concreteRT := typeToReflect(concreteType, v.program); concreteRT != nil {
					interfaceValue = value.MakeFromReflect(reflect.Zero(concreteRT))
				}
			}
			dyn := &value.InterpretedInterfaceValue{
				Value:     interfaceValue,
				TypeName:  typeName,
				IsPointer: isPointerType(concreteType),
			}
			if v.interpretedTypeSatisfiesInterface(dyn, iface) {
				v.push(value.MakeInterpretedInterface(interfaceValue, dyn.TypeName, dyn.IsPointer))
				return
			}
		}
	}

	// Get the interface's reflect.Type
	ifaceRT := typeToReflect(targetType, v.program)
	if ifaceRT == nil {
		// Can't determine interface type, pass through
		v.push(val)
		return
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
			return
		}
	}

	// For empty interface (interface{}), wrapping is only needed for typed nil values.
	// For non-nil concrete values, the existing value is already correct.
	if ifaceRT.NumMethod() == 0 && !val.IsNil() && val.IsValid() {
		if typeName := namedTypeName(concreteType); typeName != "" && shouldPreserveInterpretedNamedType(concreteType, v.program) {
			v.push(value.MakeInterpretedInterface(val, typeName, isPointerType(concreteType)))
			return
		}
		v.push(val)
		return
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
			return
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
