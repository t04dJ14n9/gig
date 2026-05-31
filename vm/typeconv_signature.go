package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func signatureToReflect(tt *types.Signature, cache map[types.Type]reflect.Type, prog *bytecode.CompiledProgram, depth int) reflect.Type {
	paramTypes := tupleToReflect(tt.Params(), cache, prog, depth+1)
	if paramTypes == nil {
		return nil
	}
	resultTypes := tupleToReflect(tt.Results(), cache, prog, depth+1)
	if resultTypes == nil {
		return nil
	}
	return reflect.FuncOf(paramTypes, resultTypes, tt.Variadic())
}

func tupleToReflect(tuple *types.Tuple, cache map[types.Type]reflect.Type, prog *bytecode.CompiledProgram, depth int) []reflect.Type {
	if tuple == nil {
		return []reflect.Type{}
	}
	out := make([]reflect.Type, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		rt := typeToReflectWithCache(tuple.At(i).Type(), cache, "", prog, depth)
		if rt == nil {
			return nil
		}
		out[i] = rt
	}
	return out
}
