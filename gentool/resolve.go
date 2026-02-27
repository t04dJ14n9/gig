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
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg == nil {
			return obj.Name()
		}
		if pkg.Path() == currentPkgPath {
			return fmt.Sprintf("%s.%s", pkgRef, obj.Name())
		}
		return ""
	case *types.Alias:
		// Handle type aliases (Go 1.23+)
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg == nil {
			return obj.Name()
		}
		if pkg.Path() == currentPkgPath {
			return fmt.Sprintf("%s.%s", pkgRef, obj.Name())
		}
		return ""
	case *types.Basic:
		return tt.Name()
	case *types.Pointer:
		elem := resolveTypeName(tt.Elem(), pkgRef)
		if elem == "" {
			return ""
		}
		return "*" + elem
	case *types.Slice:
		elem := resolveTypeName(tt.Elem(), pkgRef)
		if elem == "" {
			return ""
		}
		return "[]" + elem
	case *types.Interface:
		if tt.NumMethods() == 0 {
			return "interface{}"
		}
		return ""
	default:
		return ""
	}
}

// --- Utility ---

func isEmptyInterface(t types.Type) bool {
	iface, ok := t.Underlying().(*types.Interface)
	return ok && iface.NumMethods() == 0
}

func typeToReflectExpr(t types.Type, pkgRef string) string {
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	name := named.Obj().Name()

	switch t.Underlying().(type) {
	case *types.Interface:
		return fmt.Sprintf("reflect.TypeOf((*%s.%s)(nil)).Elem()", pkgRef, name)
	case *types.Struct:
		return fmt.Sprintf("reflect.TypeOf(%s.%s{})", pkgRef, name)
	default:
		return fmt.Sprintf("reflect.TypeOf((*%s.%s)(nil)).Elem()", pkgRef, name)
	}
}

func sanitizePkgName(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"-", "_",
		".", "_",
	).Replace(path)
}
