package vm

import (
	"go/types"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) assertInterpretedInterfaceValue(dyn *value.InterpretedInterfaceValue, targetType types.Type, original value.Value) (value.Value, bool) {
	if iface, ok := targetType.Underlying().(*types.Interface); ok {
		if iface.NumMethods() == 0 {
			return original, true
		}
		if v.interpretedTypeSatisfiesInterface(dyn, iface) {
			return original, true
		}
		return zeroValueForType(targetType, v.program), false
	}

	if targetName := namedTypeName(targetType); targetName == dyn.TypeName && isPointerType(targetType) == dyn.IsPointer {
		if dyn.IsPointer {
			return dyn.Value, true
		}
		if kindMatchesType(dyn.Value.Kind(), dyn.Value.RawSize(), targetType) {
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
		methodName := iface.Method(i).Name()
		found := false
		for _, fn := range v.program.MethodsByName[methodName] {
			if fn.ReceiverTypeName == dyn.TypeName && (!fn.ReceiverIsPointer || dyn.IsPointer) {
				found = true
				break
			}
			if !fn.ReceiverIsPointer || dyn.IsPointer {
				if _, ok := receiverForCompiledMethodTarget(methodName, dyn.Value, fn, v.program); ok {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}
