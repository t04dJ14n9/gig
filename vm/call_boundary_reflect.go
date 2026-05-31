package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

const maxBoundaryValidationDepth = 64

func (v *vm) interpreterDefinedReflectValueType(rv reflect.Value, seen map[reflect.Type]bool, depth int) (string, bool) {
	if !rv.IsValid() {
		return "", false
	}
	if depth > maxBoundaryValidationDepth {
		return "<unknown>", true
	}
	if typeName, ok := v.interpreterDefinedReflectType(rv.Type(), seen); ok {
		return typeName, true
	}

	switch rv.Kind() {
	case reflect.Interface:
		if rv.IsNil() {
			return "", false
		}
		return v.interpreterDefinedReflectValueType(rv.Elem(), seen, depth+1)
	case reflect.Ptr:
		if rv.IsNil() {
			return "", false
		}
		return v.interpreterDefinedReflectValueType(rv.Elem(), seen, depth+1)
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if typeName, ok := v.interpreterDefinedReflectValueType(rv.Index(i), seen, depth+1); ok {
				return typeName, true
			}
		}
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			if typeName, ok := v.interpreterDefinedReflectValueType(iter.Key(), seen, depth+1); ok {
				return typeName, true
			}
			if typeName, ok := v.interpreterDefinedReflectValueType(iter.Value(), seen, depth+1); ok {
				return typeName, true
			}
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if typeName, ok := v.interpreterDefinedReflectValueType(rv.Field(i), seen, depth+1); ok {
				return typeName, true
			}
		}
	}

	return "", false
}

func (v *vm) interpreterDefinedReflectType(rt reflect.Type, seen map[reflect.Type]bool) (string, bool) {
	if rt == nil || seen[rt] {
		return "", false
	}
	seen[rt] = true

	if typeName := resolveTypeName(rt, v.program); typeName != "" && isInterpreterSynthesizedReflectType(rt, v.program) {
		return typeName, true
	}

	switch rt.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Chan:
		return v.interpreterDefinedReflectType(rt.Elem(), seen)
	case reflect.Map:
		if typeName, ok := v.interpreterDefinedReflectType(rt.Key(), seen); ok {
			return typeName, true
		}
		return v.interpreterDefinedReflectType(rt.Elem(), seen)
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			if typeName, ok := v.interpreterDefinedReflectType(rt.Field(i).Type, seen); ok {
				return typeName, true
			}
		}
	}

	return "", false
}

func isInterpreterSynthesizedReflectType(rt reflect.Type, prog *bytecode.CompiledProgram) bool {
	if rt == nil {
		return false
	}
	if prog != nil {
		if name := prog.LookupTypeName(rt); name != "" {
			return true
		}
	}
	return pkgPathTypeName(rt) != ""
}
