package vm

import (
	"fmt"
	"reflect"
	"strings"

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

func externalBoundaryReflectArgType(fnType reflect.Type, argIndex int) reflect.Type {
	if fnType == nil || fnType.Kind() != reflect.Func || argIndex < 0 {
		return nil
	}
	numIn := fnType.NumIn()
	if argIndex < numIn {
		if fnType.IsVariadic() && argIndex == numIn-1 {
			return fnType.In(argIndex).Elem()
		}
		return fnType.In(argIndex)
	}
	if fnType.IsVariadic() && numIn > 0 {
		return fnType.In(numIn - 1).Elem()
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

func canPassInterpretedFuncToThirdParty(targetType reflect.Type) bool {
	if targetType == nil || targetType.Kind() != reflect.Func {
		return false
	}
	for i := 0; i < targetType.NumOut(); i++ {
		if reflectTypeContainsInterface(targetType.Out(i), make(map[reflect.Type]bool)) {
			return false
		}
	}
	return true
}

func reflectTypeContainsInterface(rt reflect.Type, seen map[reflect.Type]bool) bool {
	if rt == nil || seen[rt] {
		return false
	}
	seen[rt] = true

	switch rt.Kind() {
	case reflect.Interface:
		return true
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Chan:
		return reflectTypeContainsInterface(rt.Elem(), seen)
	case reflect.Map:
		return reflectTypeContainsInterface(rt.Key(), seen) || reflectTypeContainsInterface(rt.Elem(), seen)
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			if reflectTypeContainsInterface(rt.Field(i).Type, seen) {
				return true
			}
		}
	case reflect.Func:
		for i := 0; i < rt.NumOut(); i++ {
			if reflectTypeContainsInterface(rt.Out(i), seen) {
				return true
			}
		}
	}
	return false
}

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

func isStdlibExternalPath(path string) bool {
	if path == "" || path == "command-line-arguments" || path == "main" {
		return true
	}
	firstSlash := strings.IndexByte(path, '/')
	firstSegment := path
	if firstSlash >= 0 {
		firstSegment = path[:firstSlash]
	}
	return !strings.ContainsRune(firstSegment, '.')
}
