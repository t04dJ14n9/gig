package compiler

import (
	"go/types"
	"reflect"
)

func isEmptyStruct(t types.Type) bool {
	switch u := t.(type) {
	case *types.Named:
		t = u.Underlying()
	case *types.Alias:
		t = u.Underlying()
	}
	if st, ok := t.(*types.Struct); ok {
		return st.NumFields() == 0
	}
	return false
}

func emptyStructReflectType(t types.Type) reflect.Type {
	// Named empty structs need a synthetic field carrying a gig tag. Plain
	// reflect.StructOf cannot attach methods, but the tag preserves the
	// interpreter type name so later formatting and type checks can distinguish
	// package-level named empty structs from anonymous struct{} values.
	named, ok := t.(*types.Named)
	if !ok {
		return reflect.TypeFor[struct{}]()
	}
	obj := named.Obj()
	if obj == nil {
		return reflect.TypeFor[struct{}]()
	}
	typeName := obj.Name()
	qualName := "#" + typeName
	if pkg := obj.Pkg(); pkg != nil {
		qualName = "#" + pkg.Name() + "." + typeName
	}
	return reflect.StructOf([]reflect.StructField{{
		Name:    "gigType",
		Type:    reflect.TypeFor[struct{}](),
		PkgPath: "gig/internal",
		Tag:     reflect.StructTag(`gig:"` + qualName + `"`),
	}})
}
