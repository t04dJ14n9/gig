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
		// Cross-package: use the sanitized import alias
		return fmt.Sprintf("%s.%s", sanitizePkgName(pkg.Path()), obj.Name())
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
		return fmt.Sprintf("%s.%s", sanitizePkgName(pkg.Path()), obj.Name())
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

// collectCrossPkgImports scans a function's parameters for cross-package types
// and adds their import paths to the imports map.
// Only parameters are scanned since return values use value.FromInterface() and don't need imports.
func collectCrossPkgImports(sig *types.Signature, selfPkgPath string, imports map[string]string) {
	params := sig.Params()

	for i := 0; i < params.Len(); i++ {
		collectTypeImports(params.At(i).Type(), selfPkgPath, imports)
	}
}

// collectMethodImports collects cross-package imports for a method's parameters.
func collectMethodImports(named *types.Named, methodName string, selfPkgPath string, imports map[string]string) {
	methodSets := []*types.MethodSet{
		types.NewMethodSet(named),
		types.NewMethodSet(types.NewPointer(named)),
	}
	for _, mset := range methodSets {
		for i := 0; i < mset.Len(); i++ {
			sel := mset.At(i)
			fn := sel.Obj().(*types.Func)
			if fn.Name() != methodName {
				continue
			}
			sig := fn.Type().(*types.Signature)
			params := sig.Params()
			for j := 0; j < params.Len(); j++ {
				collectTypeImports(params.At(j).Type(), selfPkgPath, imports)
			}
			return
		}
	}
}

func collectTypeImports(t types.Type, selfPkgPath string, imports map[string]string) {
	switch tt := t.(type) {
	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg != nil && pkg.Path() != selfPkgPath {
			alias := sanitizePkgName(pkg.Path())
			imports[pkg.Path()] = alias
		}
	case *types.Pointer:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Slice:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Map:
		collectTypeImports(tt.Key(), selfPkgPath, imports)
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	}
}
