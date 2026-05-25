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
	if t == nil {
		return false
	}
	if isUserDefinedNamedType(t) {
		return true
	}
	switch tt := t.Underlying().(type) {
	case *types.Slice:
		return containsUserDefinedType(tt.Elem())
	case *types.Map:
		return containsUserDefinedType(tt.Key()) || containsUserDefinedType(tt.Elem())
	case *types.Pointer:
		return containsUserDefinedType(tt.Elem())
	case *types.Array:
		return containsUserDefinedType(tt.Elem())
	}
	return false
}

// validateExternalCallArgs checks whether any argument to an external function
// call is a user-defined type being passed to a third-party (non-stdlib) package.
func validateExternalCallArgs(pkgPath, funcName string, argTypes []types.Type) error {
	if isStdlibPath(pkgPath) {
		return nil
	}
	for i, argType := range argTypes {
		if containsUserDefinedType(argType) {
			typeName := describeType(argType)
			return fmt.Errorf(
				"cannot pass interpreter-defined type %q to third-party function %s.%s (argument %d): "+
					"custom types are not compatible with external libraries that use reflection. "+
					"Use primitive types, slices, maps, or types from registered packages instead",
				typeName, pkgPath, funcName, i+1,
			)
		}
	}
	return nil
}

// describeType returns a human-readable name for a type for error messages.
func describeType(t types.Type) string {
	if t == nil {
		return "<nil>"
	}
	return t.String()
}
