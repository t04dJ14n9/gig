package compiler

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"
)

// extractMethodName strips SSA receiver qualification from a method name.
// SSA names look like "(*Type).Method" or "pkgpath.Method"; this extracts just "Method".
func extractMethodName(ssaName string) string {
	name := ssaName
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ")"); idx >= 0 {
		rest := name[idx+1:]
		if len(rest) > 0 && rest[0] == '.' {
			name = rest[1:]
		}
	}
	return name
}

func methodOwnerPkgPath(fn *ssa.Function) string {
	if fn == nil || fn.Signature == nil {
		return ""
	}
	if fn.Pkg != nil && fn.Pkg.Pkg != nil {
		return fn.Pkg.Pkg.Path()
	}
	recv := fn.Signature.Recv()
	if recv == nil {
		return ""
	}
	recvType := recv.Type()
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return ""
	}
	return named.Obj().Pkg().Path()
}

// extractReceiverTypeName extracts the package-path-qualified type name from a receiver type.
// For pointer receivers like *Reader, it unwraps the pointer.
// Returns "pkgPath.TypeName" (e.g., "encoding/json.Encoder") for use as a DirectCall lookup key.
func extractReceiverTypeName(recvType types.Type) string {
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok {
		return ""
	}
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg != nil {
		return pkg.Path() + "." + obj.Name()
	}
	return obj.Name()
}

// extractNamedType unwraps pointer types to find the underlying *types.Named type.
// Returns nil if the type is not named (e.g., interface types without a name).
func extractNamedType(t types.Type) *types.Named {
	for {
		switch tt := t.(type) {
		case *types.Named:
			return tt
		case *types.Pointer:
			t = tt.Elem()
		default:
			return nil
		}
	}
}

// extractReceiverShortName extracts the unqualified type name from a receiver type.
// For pointer receivers like *Reader, it unwraps the pointer.
// Returns just the type name (e.g., "Reader"), without package path.
func extractReceiverShortName(recvType types.Type) string {
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	if named, ok := recvType.(*types.Named); ok {
		return named.Obj().Name()
	}
	return ""
}
