package compiler

import (
	"go/constant"
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func typedNilConstValue(t types.Type, val constant.Value) any {
	// Typed nil constants must be represented as reflect.Zero(type), not nil.
	// The VM relies on that reflect value to preserve map/slice/chan/function,
	// pointer, interface, and struct type information after OpConst.
	if val != nil {
		return nil
	}
	if rt := constTypeToReflect(t); rt != nil {
		return reflect.Zero(rt)
	}
	return nil
}

func constTypeToReflect(t types.Type) reflect.Type {
	if isEmptyStruct(t) {
		return emptyStructReflectType(t)
	}

	switch typ := t.Underlying().(type) {
	case *types.Basic:
		return bytecode.BasicKindToReflectType[typ.Kind()]
	case *types.Map:
		return mapConstReflectType(typ)
	case *types.Slice:
		return elemConstReflectType(typ.Elem(), reflect.SliceOf)
	case *types.Pointer:
		return elemConstReflectType(typ.Elem(), reflect.PointerTo)
	case *types.Chan:
		return chanConstReflectType(typ)
	case *types.Interface:
		if typ.NumMethods() == 0 {
			return reflect.TypeFor[any]()
		}
	case *types.Signature:
		return buildFuncType(typ)
	}
	return nil
}

func mapConstReflectType(typ *types.Map) reflect.Type {
	keyRT := constTypeToReflect(typ.Key())
	elemRT := constTypeToReflect(typ.Elem())
	if keyRT != nil && elemRT != nil {
		return reflect.MapOf(keyRT, elemRT)
	}
	return nil
}

func elemConstReflectType(elem types.Type, build func(reflect.Type) reflect.Type) reflect.Type {
	elemRT := constTypeToReflect(elem)
	if elemRT == nil {
		return nil
	}
	return build(elemRT)
}

func chanConstReflectType(typ *types.Chan) reflect.Type {
	elemRT := constTypeToReflect(typ.Elem())
	if elemRT == nil {
		return nil
	}
	return reflect.ChanOf(chanDirection(typ), elemRT)
}

func chanDirection(typ *types.Chan) reflect.ChanDir {
	switch typ.Dir() {
	case types.SendOnly:
		return reflect.SendDir
	case types.RecvOnly:
		return reflect.RecvDir
	default:
		return reflect.BothDir
	}
}

func buildFuncType(sig *types.Signature) reflect.Type {
	params := make([]reflect.Type, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		pt := constTypeToReflect(sig.Params().At(i).Type())
		if pt == nil {
			return nil
		}
		params[i] = pt
	}
	results := make([]reflect.Type, sig.Results().Len())
	for i := 0; i < sig.Results().Len(); i++ {
		rt := constTypeToReflect(sig.Results().At(i).Type())
		if rt == nil {
			return nil
		}
		results[i] = rt
	}
	return reflect.FuncOf(params, results, sig.Variadic())
}
