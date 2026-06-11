package gentool

import (
	"fmt"
	"go/types"
	"strings"
)

// currentPkgPath is set per-package during generation to resolve same-package named types.
var currentPkgPath string

// errorTypeName is the string constant for the "error" type name.
const errorTypeName = "error"

// --- Type name resolution ---

func resolveTypeName(t types.Type, pkgRef string) string {
	switch tt := t.(type) {
	case *types.Named:
		return resolveQualifiedObjectName(tt.Obj(), pkgRef)
	case *types.Alias:
		return resolveQualifiedObjectName(tt.Obj(), pkgRef)
	case *types.Basic:
		if tt.Kind() == types.UnsafePointer {
			return "unsafe.Pointer"
		}
		return tt.Name()
	case *types.Pointer:
		return resolvePrefixedTypeName("*", tt.Elem(), pkgRef)
	case *types.Slice:
		return resolvePrefixedTypeName("[]", tt.Elem(), pkgRef)
	case *types.Interface:
		if tt.NumMethods() == 0 {
			return "interface{}"
		}
		return resolveInterfaceTypeName(tt, pkgRef)
	case *types.Chan:
		return resolveChanTypeName(tt, pkgRef)
	case *types.Signature:
		return resolveFuncTypeName(tt, pkgRef)
	case *types.Array:
		return resolveArrayTypeName(tt, pkgRef)
	default:
		return ""
	}
}

func resolveQualifiedObjectName(obj types.Object, pkgRef string) string {
	pkg := obj.Pkg()
	if pkg == nil {
		return obj.Name()
	}
	if pkg.Path() == currentPkgPath {
		return fmt.Sprintf("%s.%s", pkgRef, obj.Name())
	}
	return fmt.Sprintf("%s.%s", sanitizePkgName(pkg.Path()), obj.Name())
}

func resolvePrefixedTypeName(prefix string, elem types.Type, pkgRef string) string {
	elemName := resolveTypeName(elem, pkgRef)
	if elemName == "" {
		return ""
	}
	return prefix + elemName
}

func resolveChanTypeName(ch *types.Chan, pkgRef string) string {
	elemName := resolveTypeName(ch.Elem(), pkgRef)
	if elemName == "" {
		return ""
	}
	switch ch.Dir() {
	case types.SendRecv:
		return fmt.Sprintf("chan %s", elemName)
	case types.SendOnly:
		return fmt.Sprintf("chan<- %s", elemName)
	case types.RecvOnly:
		return fmt.Sprintf("<-chan %s", elemName)
	default:
		return ""
	}
}

// resolveArrayTypeName returns the Go type string for an array type, e.g. "[32]byte".
func resolveArrayTypeName(arr *types.Array, pkgRef string) string {
	elemName := resolveTypeName(arr.Elem(), pkgRef)
	if elemName == "" {
		return ""
	}
	return fmt.Sprintf("[%d]%s", arr.Len(), elemName)
}

func sanitizePkgName(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"-", "_",
		".", "_",
	).Replace(path)
}
