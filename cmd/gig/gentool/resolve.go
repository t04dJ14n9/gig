package gentool

import (
	"fmt"
	"go/constant"
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
		if tt.Kind() == types.UnsafePointer {
			return "unsafe.Pointer"
		}
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
		return resolveInterfaceTypeName(tt, pkgRef)
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

// resolveInterfaceTypeName returns the Go type string for an unnamed interface
// with methods, e.g. "interface{ Printf(string, ...interface{}) }".
// If any method signature contains an unresolvable type, it returns "".
func resolveInterfaceTypeName(iface *types.Interface, pkgRef string) string {
	var methodStrs []string
	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		sig := m.Type().(*types.Signature)
		methodStr := resolveMethodSigStr(m.Name(), sig, pkgRef)
		if methodStr == "" {
			return ""
		}
		methodStrs = append(methodStrs, methodStr)
	}
	return "interface{ " + strings.Join(methodStrs, "; ") + " }"
}

// resolveMethodSigStr returns the Go method signature string, e.g. "Printf(string, ...interface{})".
func resolveMethodSigStr(name string, sig *types.Signature, pkgRef string) string {
	params := sig.Params()
	results := sig.Results()

	var paramStrs []string
	for i := 0; i < params.Len(); i++ {
		pType := params.At(i).Type()
		// For variadic parameters, the last param is a slice type.
		// If the function is variadic and this is the last param, show ...Elem instead of []Elem.
		if sig.Variadic() && i == params.Len()-1 {
			sliceType, ok := pType.(*types.Slice)
			if !ok {
				return ""
			}
			elemName := resolveTypeName(sliceType.Elem(), pkgRef)
			if elemName == "" {
				return ""
			}
			paramStrs = append(paramStrs, "..."+elemName)
		} else {
			pName := resolveTypeName(pType, pkgRef)
			if pName == "" {
				return ""
			}
			paramStrs = append(paramStrs, pName)
		}
	}

	var resultPart string
	switch results.Len() {
	case 0:
		resultPart = ""
	case 1:
		rName := resolveTypeName(results.At(0).Type(), pkgRef)
		if rName == "" {
			return ""
		}
		resultPart = " " + rName
	default:
		var retStrs []string
		for i := 0; i < results.Len(); i++ {
			rName := resolveTypeName(results.At(i).Type(), pkgRef)
			if rName == "" {
				return ""
			}
			retStrs = append(retStrs, rName)
		}
		resultPart = " (" + strings.Join(retStrs, ", ") + ")"
	}

	return fmt.Sprintf("%s(%s)%s", name, strings.Join(paramStrs, ", "), resultPart)
}

// --- Type checking ---

// canWrapParam checks whether a parameter type can be wrapped in a DirectCall.
// It handles named types, aliases, and unnamed types with appropriate restrictions.
func canWrapParam(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		pkg := obj.Pkg()

		if pkg == nil {
			return obj.Name() == errorTypeName
		}

		if pkg.Path() == currentPkgPath {
			return canWrapType(t.Underlying(), false)
		}

		// Cross-package named types: allow if we can extract via .Interface().(Type)
		return canWrapType(t.Underlying(), true)
	}

	// Handle type aliases (Go 1.23+), including builtin 'any'
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		pkg := obj.Pkg()

		if pkg == nil {
			return canWrapType(t.Underlying(), false)
		}

		if pkg.Path() == currentPkgPath {
			return canWrapType(t.Underlying(), false)
		}

		return canWrapType(t.Underlying(), true)
	}

	return canWrapType(t.Underlying(), false)
}

// canWrapType checks whether a type can be wrapped in a DirectCall.
// If crossPkg is true, all representable types are allowed (extracted via .Interface()).
// If false, stricter checks apply (e.g., slices only with basic element types).
func canWrapType(t types.Type, crossPkg bool) bool {
	switch ut := t.(type) {
	case *types.Basic:
		return ut.Kind() != types.UnsafePointer && ut.Kind() != types.Invalid
	case *types.Slice:
		if crossPkg {
			return true
		}
		// Same-package: allow slices with basic element types only
		if _, ok := ut.Elem().Underlying().(*types.Basic); ok {
			return true
		}
		return false
	case *types.Interface:
		return true
	case *types.Pointer:
		if crossPkg {
			return true
		}
		if bt, ok := ut.Elem().(*types.Basic); ok {
			return bt.Kind() != types.UnsafePointer && bt.Kind() != types.Invalid
		}
		_, isNamed := ut.Elem().(*types.Named)
		return isNamed
	case *types.Struct:
		return true
	case *types.Map:
		return true
	case *types.Chan:
		return true
	case *types.Signature:
		return true
	case *types.Array:
		return true
	default:
		return false
	}
}

// isEmptyInterface returns true if the type is (or underlies) an empty interface.
func isEmptyInterface(t types.Type) bool {
	iface, ok := t.Underlying().(*types.Interface)
	return ok && iface.NumMethods() == 0
}

// isNamedOrAlias returns true if the type is a named type or a type alias
// (as opposed to a literal interface{} or other unnamed type).
// Builtin aliases like 'any' and 'error' are excluded because they are
// semantically equivalent to their underlying types and should be treated
// as unnamed for variadic slice generation.
func isNamedOrAlias(t types.Type) bool {
	switch tt := t.(type) {
	case *types.Named:
		if tt.Obj().Pkg() == nil {
			// Builtin named types (like 'error') are effectively unnamed
			// for variadic slice generation purposes.
			return false
		}
		return true
	case *types.Alias:
		if tt.Obj().Pkg() == nil {
			// Builtin aliases (like 'any') are effectively unnamed
			// for variadic slice generation purposes.
			return false
		}
		return true
	}
	return false
}

// --- Import collection ---

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

// collectTypeImports recursively collects import paths for cross-package types.
func collectTypeImports(t types.Type, selfPkgPath string, imports map[string]string) {
	switch tt := t.(type) {
	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg != nil && pkg.Path() != selfPkgPath {
			alias := sanitizePkgName(pkg.Path())
			imports[pkg.Path()] = alias
		}
	case *types.Alias:
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg != nil && pkg.Path() != selfPkgPath {
			alias := sanitizePkgName(pkg.Path())
			imports[pkg.Path()] = alias
		}
	case *types.Basic:
		if tt.Kind() == types.UnsafePointer {
			imports["unsafe"] = "unsafe"
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
	case *types.Interface:
		// Recurse into interface method signatures for unnamed interfaces
		for i := 0; i < tt.NumMethods(); i++ {
			collectTypeImports(tt.Method(i).Type(), selfPkgPath, imports)
		}
	case *types.Chan:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	}
}

// --- Utilities ---

// needsUintCast checks whether a constant needs a uint64() cast to prevent
// overflow when passed as an untyped int to AddConstant. This happens when
// the constant's underlying type is unsigned and its value exceeds math.MaxInt64.
func needsUintCast(c *types.Const) bool {
	basic, ok := c.Type().Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint, types.Uint64, types.Uintptr, types.UntypedInt:
		// Check if the value overflows int64
		val := c.Val()
		if val.Kind() == constant.Int {
			if v, ok := constant.Uint64Val(val); ok {
				return v > (1<<63 - 1) // > math.MaxInt64
			}
		}
	}
	return false
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
