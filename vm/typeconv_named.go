package vm

import (
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func namedToReflect(tt *types.Named, cache map[types.Type]reflect.Type, prog *bytecode.CompiledProgram, depth int) reflect.Type {
	if rt := externalNamedTypeToReflect(tt, cache, prog); rt != nil {
		return rt
	}

	cache[tt] = nil
	typeName := tt.Obj().Name()
	result := typeToReflectWithCache(tt.Underlying(), cache, namedTypeSuffix(tt), prog, depth+1)
	if result == nil {
		return nil
	}

	cache[tt] = result
	if prog != nil && result.Kind() == reflect.Struct {
		prog.RegisterTypeName(result, typeName)
	}
	return result
}

func externalNamedTypeToReflect(tt *types.Named, cache map[types.Type]reflect.Type, prog *bytecode.CompiledProgram) reflect.Type {
	if prog == nil || prog.TypeResolver == nil {
		return nil
	}
	if rt, ok := prog.TypeResolver.LookupExternalType(tt); ok {
		cache[tt] = rt
		return rt
	}

	obj := tt.Obj()
	if pkg := obj.Pkg(); pkg != nil {
		if rt, ok := prog.TypeResolver.LookupExternalTypeByName(pkg.Path(), obj.Name()); ok {
			cache[tt] = rt
			return rt
		}
	}
	return nil
}

func namedTypeSuffix(tt *types.Named) string {
	typeName := tt.Obj().Name()
	if pkg := tt.Obj().Pkg(); pkg != nil {
		return "#" + pkg.Name() + "." + typeName
	}
	return "#" + typeName
}
