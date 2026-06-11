package vm

import (
	"go/types"

	"github.com/t04dJ14n9/gig/model/bytecode"
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
	if !v.canResolveInterpretedInterface(dyn, iface) {
		return false
	}
	for i := 0; i < iface.NumMethods(); i++ {
		if !v.interpretedMethodSatisfiesInterface(dyn, iface.Method(i).Name()) {
			return false
		}
	}
	return true
}

func (v *vm) canResolveInterpretedInterface(dyn *value.InterpretedInterfaceValue, iface *types.Interface) bool {
	return dyn != nil && iface != nil && v != nil && v.program != nil
}

func (v *vm) interpretedMethodSatisfiesInterface(dyn *value.InterpretedInterfaceValue, methodName string) bool {
	for _, fn := range v.program.MethodsByName[methodName] {
		if compiledMethodMatchesInterpretedType(dyn, fn) {
			return true
		}
		if v.compiledMethodMatchesDynamicReceiver(dyn, methodName, fn) {
			return true
		}
	}
	return false
}

func compiledMethodMatchesInterpretedType(
	dyn *value.InterpretedInterfaceValue,
	fn *bytecode.CompiledFunction,
) bool {
	return fn.ReceiverTypeName == dyn.TypeName && receiverEligibleForInterpretedValue(dyn, fn)
}

func (v *vm) compiledMethodMatchesDynamicReceiver(
	dyn *value.InterpretedInterfaceValue,
	methodName string,
	fn *bytecode.CompiledFunction,
) bool {
	if !receiverEligibleForInterpretedValue(dyn, fn) {
		return false
	}
	_, ok := receiverForCompiledMethodTarget(methodName, dyn.Value, fn, v.program)
	return ok
}

func receiverEligibleForInterpretedValue(
	dyn *value.InterpretedInterfaceValue,
	fn *bytecode.CompiledFunction,
) bool {
	return !fn.ReceiverIsPointer || dyn.IsPointer
}
