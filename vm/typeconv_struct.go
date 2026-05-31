package vm

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func structToReflect(tt *types.Struct, cache map[types.Type]reflect.Type, uniqueSuffix string, prog *bytecode.CompiledProgram, depth int) reflect.Type {
	fields := make([]reflect.StructField, 0, tt.NumFields())
	hasUnexported := false

	for i := 0; i < tt.NumFields(); i++ {
		sf, ok, unexported := structFieldToReflect(tt, i, cache, uniqueSuffix, prog, depth)
		if !ok {
			continue
		}
		hasUnexported = hasUnexported || unexported
		fields = append(fields, sf)
	}

	if len(fields) == 0 {
		return emptyStructToReflect(uniqueSuffix)
	}
	addUniqueStructTags(fields, uniqueSuffix, hasUnexported)
	return reflect.StructOf(fields)
}

func structFieldToReflect(
	tt *types.Struct,
	fieldIndex int,
	cache map[types.Type]reflect.Type,
	uniqueSuffix string,
	prog *bytecode.CompiledProgram,
	depth int,
) (reflect.StructField, bool, bool) {
	f := tt.Field(fieldIndex)
	ft := typeToReflectWithCache(f.Type(), cache, fieldTypeSuffix(f.Type()), prog, depth+1)
	if ft == nil {
		return reflect.StructField{}, false, false
	}

	sf := reflect.StructField{
		Name:      f.Name(),
		Type:      ft,
		Anonymous: f.Anonymous(),
	}
	if sf.Anonymous && sf.Type.Kind() == reflect.Interface && sf.Type.NumMethod() > 0 {
		sf.Anonymous = false
	}

	unexported := !f.Exported()
	if unexported {
		applyUnexportedFieldMetadata(&sf, f, bareTypeSuffix(uniqueSuffix))
	}
	if tag := tt.Tag(fieldIndex); tag != "" {
		appendStructTag(&sf, reflect.StructTag(tag))
	}
	return sf, true, unexported
}

func fieldTypeSuffix(t types.Type) string {
	if named, ok := t.(*types.Named); ok {
		return "#" + named.Obj().Name()
	}
	return ""
}

func applyUnexportedFieldMetadata(sf *reflect.StructField, f *types.Var, suffix string) {
	if sf.Anonymous {
		sf.Anonymous = false
		sf.Tag = reflect.StructTag(`gig_embed:"1"`)
	}
	if pkg := f.Pkg(); pkg != nil {
		sf.PkgPath = pkg.Path() + suffix
		return
	}
	if suffix != "" {
		sf.PkgPath = "gig/internal" + suffix
	}
}

func bareTypeSuffix(uniqueSuffix string) string {
	if idx := strings.LastIndex(uniqueSuffix, "."); idx > 0 && uniqueSuffix[0] == '#' {
		return "#" + uniqueSuffix[idx+1:]
	}
	return uniqueSuffix
}

func addUniqueStructTags(fields []reflect.StructField, uniqueSuffix string, hasUnexported bool) {
	if hasUnexported || uniqueSuffix == "" {
		return
	}
	gigTag := reflect.StructTag(`gig:"` + uniqueSuffix + `"`)
	for i := range fields {
		appendStructTag(&fields[i], gigTag)
	}
}

func appendStructTag(sf *reflect.StructField, tag reflect.StructTag) {
	if sf.Tag == "" {
		sf.Tag = tag
		return
	}
	sf.Tag += " " + tag
}

func emptyStructToReflect(uniqueSuffix string) reflect.Type {
	if uniqueSuffix == "" {
		return reflect.TypeOf(struct{}{})
	}
	return reflect.StructOf([]reflect.StructField{{
		Name:    "gigType",
		Type:    reflect.TypeOf(struct{}{}),
		PkgPath: "gig/internal",
		Tag:     reflect.StructTag(`gig:"` + uniqueSuffix + `"`),
	}})
}
