// typecheck.go provides compile-time type flow analysis for custom type safety.
package compiler

import (
	"fmt"
	"go/types"
	"strings"
)

// isStdlibPath returns true if the import path belongs to the Go standard library.
// Stdlib paths have no dots in the first path segment (e.g., "fmt", "encoding/json").
// Third-party paths contain dots (e.g., "github.com/foo/bar", "golang.org/x/tools").
//
// Note: golang.org/x/* packages are treated as third-party. While they are
// maintained by the Go team, they are not part of the standard library and
// may use reflection on argument types. This conservative classification
// avoids false negatives; users who need to pass custom types to x/ packages
// can use WithAllowUnsafeTypePass().
func isStdlibPath(path string) bool {
	if path == "" || path == "command-line-arguments" || path == "main" {
		return false
	}
	firstSlash := strings.IndexByte(path, '/')
	firstSegment := path
	if firstSlash >= 0 {
		firstSegment = path[:firstSlash]
	}
	return !strings.ContainsRune(firstSegment, '.')
}

// isUserDefinedNamedType returns true if the type is a named type defined in the
// user's script (the "main" / "command-line-arguments" package), not from an
// external registered package.
// Unwraps pointers: *MyStruct → MyStruct → check package.
func isUserDefinedNamedType(t types.Type) bool {
	if t == nil {
		return false
	}
	for {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		} else {
			break
		}
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	pkgPath := pkg.Path()
	return pkgPath == "command-line-arguments" || pkgPath == "main"
}

// containsUserDefinedType recursively checks if a type contains a user-defined
// named type. This catches nested cases like []MyStruct, map[string]MyStruct,
// *MyStruct, [][]MyStruct, etc.
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
	if containsUserDefinedTypeParamConstraint(t, seen) {
		return true
	}
	return containsUserDefinedUnderlying(t.Underlying(), seen)
}

func containsUserDefinedTypeParamConstraint(t types.Type, seen map[types.Type]bool) bool {
	typeParam, ok := t.(*types.TypeParam)
	return ok && containsUserDefinedTypeSeen(typeParam.Constraint(), seen)
}

func containsUserDefinedUnderlying(t types.Type, seen map[types.Type]bool) bool {
	if containsUserDefinedElementContainer(t, seen) {
		return true
	}
	if containsUserDefinedAggregateContainer(t, seen) {
		return true
	}
	return containsUserDefinedCallableOrInterface(t, seen)
}

func containsUserDefinedElementContainer(t types.Type, seen map[types.Type]bool) bool {
	switch tt := t.(type) {
	case *types.Slice:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Pointer:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Array:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	case *types.Chan:
		return containsUserDefinedTypeSeen(tt.Elem(), seen)
	default:
		return false
	}
}

func containsUserDefinedAggregateContainer(t types.Type, seen map[types.Type]bool) bool {
	switch tt := t.(type) {
	case *types.Map:
		return containsUserDefinedMap(tt, seen)
	case *types.Struct:
		return containsUserDefinedStruct(tt, seen)
	default:
		return false
	}
}

func containsUserDefinedCallableOrInterface(t types.Type, seen map[types.Type]bool) bool {
	switch tt := t.(type) {
	case *types.Signature:
		return containsUserDefinedSignature(tt, seen)
	case *types.Interface:
		return containsUserDefinedInterface(tt, seen)
	default:
		return false
	}
}

func containsUserDefinedMap(t *types.Map, seen map[types.Type]bool) bool {
	return containsUserDefinedTypeSeen(t.Key(), seen) || containsUserDefinedTypeSeen(t.Elem(), seen)
}

func containsUserDefinedStruct(t *types.Struct, seen map[types.Type]bool) bool {
	for i := 0; i < t.NumFields(); i++ {
		if containsUserDefinedTypeSeen(t.Field(i).Type(), seen) {
			return true
		}
	}
	return false
}

func containsUserDefinedSignature(t *types.Signature, seen map[types.Type]bool) bool {
	if recv := t.Recv(); recv != nil && containsUserDefinedTypeSeen(recv.Type(), seen) {
		return true
	}
	return containsUserDefinedTuple(t.Params(), seen) || containsUserDefinedTuple(t.Results(), seen)
}

func containsUserDefinedInterface(t *types.Interface, seen map[types.Type]bool) bool {
	return containsUserDefinedInterfaceMethods(t, seen) || containsUserDefinedEmbeddedInterfaces(t, seen)
}

func containsUserDefinedInterfaceMethods(t *types.Interface, seen map[types.Type]bool) bool {
	for i := 0; i < t.NumMethods(); i++ {
		if containsUserDefinedTypeSeen(t.Method(i).Type(), seen) {
			return true
		}
	}
	return false
}

func containsUserDefinedEmbeddedInterfaces(t *types.Interface, seen map[types.Type]bool) bool {
	for i := 0; i < t.NumEmbeddeds(); i++ {
		if containsUserDefinedTypeSeen(t.EmbeddedType(i), seen) {
			return true
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

// validateExternalCallArgs checks whether any argument to an external function
// call is a user-defined type being passed to a third-party (non-stdlib) package.
func validateExternalCallArgs(pkgPath, funcName string, argTypes []types.Type) error {
	args := make([]externalCallArg, len(argTypes))
	for i, argType := range argTypes {
		args[i] = externalCallArg{SourceType: argType}
	}
	return validateExternalCallBoundary(pkgPath, funcName, args)
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

// validateExternalCallArgBoundary applies the per-argument policy after the
// caller has established that the destination package is third-party code.
func validateExternalCallArgBoundary(pkgPath, funcName string, argIndex int, arg externalCallArg) error {
	if arg.AllowInterfaceProxy || !containsUserDefinedType(arg.SourceType) {
		return nil
	}
	return externalCallBoundaryError(pkgPath, funcName, argIndex, describeType(arg.SourceType))
}

// externalCallBoundaryError keeps the diagnostic text centralized so proxy,
// type-scan, and package-trust decisions do not duplicate user-facing wording.
func externalCallBoundaryError(pkgPath, funcName string, argIndex int, typeName string) error {
	return fmt.Errorf(
		"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
			"custom types are not compatible with external libraries that use reflection. "+
			"Use primitive types, slices, maps, types from registered packages, or a registered interface proxy instead",
		typeName, pkgPath, funcName, argIndex+1,
	)
}

// describeType returns a human-readable name for a type for error messages.
func describeType(t types.Type) string {
	if t == nil {
		return "<nil>"
	}
	return t.String()
}
