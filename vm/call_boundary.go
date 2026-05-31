package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) validateExternalBoundary(rc *bytecode.ResolvedCall, args []value.Value) error {
	if rc == nil || v.program.AllowUnsafeTypePass || isStdlibExternalPath(rc.PkgPath) {
		return nil
	}
	for i, arg := range args {
		targetType := externalBoundaryReflectArgType(rc.FnType, i)
		if typeName, ok := v.interpreterDefinedBoundaryType(arg, targetType); ok {
			return fmt.Errorf(
				"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
					"value crossed the boundary through an interface. "+
					"Use primitive types, slices, maps, types from registered packages, or a registered interface proxy instead",
				typeName, rc.PkgPath, rc.FuncName, i+1,
			)
		}
	}
	return nil
}

func (v *vm) interpreterDefinedBoundaryType(arg value.Value, targetType reflect.Type) (string, bool) {
	if dyn, ok := arg.InterpretedInterface(); ok {
		if dyn.TypeName != "" {
			return dyn.TypeName, true
		}
		return "<unknown>", true
	}
	if arg.Kind() == value.KindFunc {
		closure, ok := arg.RawObj().(*Closure)
		if !ok {
			return "", false
		}
		if canPassInterpretedFuncToThirdParty(targetType) {
			return "", false
		}
		if closure.Fn != nil && closure.Fn.Name != "" {
			return "func " + closure.Fn.Name, true
		}
		return "func", true
	}
	if rv, ok := arg.ReflectValue(); ok {
		return v.interpreterDefinedReflectValueType(rv, make(map[reflect.Type]bool), 0)
	}
	return "", false
}
