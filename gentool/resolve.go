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
	case *types.Chan:
		elemName := resolveTypeName(tt.Elem(), pkgRef)
		if elemName == "" {
			return ""
		}
		switch tt.Dir() {
		case types.SendRecv:
			return fmt.Sprintf("chan %s", elemName)
		case types.SendOnly:
			return fmt.Sprintf("chan<- %s", elemName)
		case types.RecvOnly:
			return fmt.Sprintf("<-chan %s", elemName)
		}
		return ""
	case *types.Signature:
		return resolveFuncTypeName(tt, pkgRef)
	case *types.Array:
		return resolveArrayTypeName(tt, pkgRef)
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

// resolveFuncTypeName returns the Go type string for a function signature,
// e.g. "func(string) bool" or "func(int, int) int".
func resolveFuncTypeName(sig *types.Signature, pkgRef string) string {
	params := sig.Params()
	results := sig.Results()

	var paramStrs []string
	for i := 0; i < params.Len(); i++ {
		pName := resolveTypeName(params.At(i).Type(), pkgRef)
		if pName == "" {
			return ""
		}
		paramStrs = append(paramStrs, pName)
	}

	switch results.Len() {
	case 0:
		return fmt.Sprintf("func(%s)", strings.Join(paramStrs, ", "))
	case 1:
		rName := resolveTypeName(results.At(0).Type(), pkgRef)
		if rName == "" {
			return ""
		}
		return fmt.Sprintf("func(%s) %s", strings.Join(paramStrs, ", "), rName)
	default:
		var retStrs []string
		for i := 0; i < results.Len(); i++ {
			rName := resolveTypeName(results.At(i).Type(), pkgRef)
			if rName == "" {
				return ""
			}
			retStrs = append(retStrs, rName)
		}
		return fmt.Sprintf("func(%s) (%s)", strings.Join(paramStrs, ", "), strings.Join(retStrs, ", "))
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
	case *types.Array:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Map:
		collectTypeImports(tt.Key(), selfPkgPath, imports)
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Signature:
		// Recurse into function parameter and return types
		params := tt.Params()
		for i := 0; i < params.Len(); i++ {
			collectTypeImports(params.At(i).Type(), selfPkgPath, imports)
		}
		results := tt.Results()
		for i := 0; i < results.Len(); i++ {
			collectTypeImports(results.At(i).Type(), selfPkgPath, imports)
		}
	case *types.Chan:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	}
}
