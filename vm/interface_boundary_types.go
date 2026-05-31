package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

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
	// A script-defined named type may be backed by a synthesized reflect type
	// with incomplete method metadata. Preserve the interpreted wrapper whenever
	// the reflected name cannot prove it is the same native named type.
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
