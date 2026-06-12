// typecheck.go provides compile-time type flow analysis for custom type safety.
package compiler

import (
	"fmt"
	"go/types"
	"strings"
)

const mainPackagePath = "main"

func isStdlibPath(path string) bool {
	if path == "" || path == "command-line-arguments" || path == mainPackagePath {
		return false
	}
	firstSlash := strings.IndexByte(path, '/')
	firstSegment := path
	if firstSlash >= 0 {
		firstSegment = path[:firstSlash]
	}
	return !strings.ContainsRune(firstSegment, '.')
}

func isUserDefinedNamedType(t types.Type) bool {
	named, ok := namedTypeAfterPointerUnwrap(t)
	if !ok {
		return false
	}
	return isScriptPackagePath(namedTypePackagePath(named))
}

func namedTypeAfterPointerUnwrap(t types.Type) (*types.Named, bool) {
	if t == nil {
		return nil, false
	}
	named, ok := unwrapPointerType(t).(*types.Named)
	return named, ok
}

func unwrapPointerType(t types.Type) types.Type {
	for ptr, ok := t.(*types.Pointer); ok; ptr, ok = t.(*types.Pointer) {
		t = ptr.Elem()
	}
	return t
}

func namedTypePackagePath(named *types.Named) string {
	obj := named.Obj()
	if obj == nil {
		return ""
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return ""
	}
	return pkg.Path()
}

func isScriptPackagePath(pkgPath string) bool {
	return pkgPath == "command-line-arguments" || pkgPath == mainPackagePath
}

func containsUserDefinedType(t types.Type) bool {
	return containsUserDefinedTypeSeen(t, make(map[types.Type]bool))
}

func containsUserDefinedTypeSeen(t types.Type, seen map[types.Type]bool) bool {
	if t == nil {
		return false
	}
	if seen[t] {
		return false
	}
	seen[t] = true

	if isUserDefinedNamedType(t) {
		return true
	}
	if typeParam, ok := t.(*types.TypeParam); ok {
		return containsUserDefinedTypeSeen(typeParam.Constraint(), seen)
	}
	return containsUserDefinedUnderlying(t.Underlying(), seen)
}

func containsUserDefinedUnderlying(t types.Type, seen map[types.Type]bool) bool {
	switch tt := t.(type) {
	case *types.Slice:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Pointer:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Array:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Chan:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Map:
		return containsUserDefinedTypeSeen(tt.Key(), seen) || containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Struct:
		for i := 0; i < tt.NumFields(); i++ {
			if containsUserDefinedTypeSeen(tt.Field(i).Type(), seen) {
				return true
			}
		}
	case *types.Signature:
		if recv := tt.Recv(); recv != nil && containsUserDefinedTypeSeen(recv.Type(), seen) {
			return true
		}
		return containsUserDefinedTuple(tt.Params(), seen) || containsUserDefinedTuple(tt.Results(), seen)
	case *types.Interface:
		for i := 0; i < tt.NumMethods(); i++ {
			if containsUserDefinedTypeSeen(tt.Method(i).Type(), seen) {
				return true
			}
		}
		for i := 0; i < tt.NumEmbeddeds(); i++ {
			if containsUserDefinedTypeSeen(tt.EmbeddedType(i), seen) {
				return true
			}
		}
	}
	return false
}

func containsUserDefinedTuple(tuple *types.Tuple, seen map[types.Type]bool) bool {
	if tuple == nil {
		return false
	}
	for i := 0; i < tuple.Len(); i++ {
		if containsUserDefinedTypeSeen(tuple.At(i).Type(), seen) {
			return true
		}
	}
	return false
}

type externalCallArg struct {
	SourceType          types.Type
	AllowInterfaceProxy bool
}

func validateExternalCallBoundary(pkgPath, funcName string, args []externalCallArg) error {
	if isStdlibPath(pkgPath) {
		return nil
	}
	for i, arg := range args {
		if err := validateExternalCallArgBoundary(pkgPath, funcName, i, arg); err != nil {
			return err
		}
	}
	return nil
}

func validateExternalCallArgBoundary(pkgPath, funcName string, argIndex int, arg externalCallArg) error {
	if arg.AllowInterfaceProxy || !containsUserDefinedType(arg.SourceType) {
		return nil
	}
	return fmt.Errorf(
		"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
			"custom types are not compatible with external libraries that use reflection. "+
			"Use primitive types, slices, maps, types from registered packages, or a registered interface proxy instead",
		describeType(arg.SourceType), pkgPath, funcName, argIndex+1,
	)
}

func describeType(t types.Type) string {
	if t == nil {
		return "<nil>"
	}
	return t.String()
}
