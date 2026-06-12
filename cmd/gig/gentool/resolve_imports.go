package gentool

import "go/types"

// collectCrossPkgImports scans a function's parameters for cross-package types.
// Return values use value.FromInterface() and do not need imports here.
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
			params := fn.Type().(*types.Signature).Params()
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
		collectObjectImport(tt.Obj(), selfPkgPath, imports)
	case *types.Alias:
		collectObjectImport(tt.Obj(), selfPkgPath, imports)
	case *types.Basic:
		collectBasicImport(tt, imports)
	case *types.Map:
		collectTypeImports(tt.Key(), selfPkgPath, imports)
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Signature:
		collectTupleImports(tt.Params(), selfPkgPath, imports)
		collectTupleImports(tt.Results(), selfPkgPath, imports)
	case *types.Interface:
		collectInterfaceImports(tt, selfPkgPath, imports)
	default:
		collectSingleElementTypeImports(t, selfPkgPath, imports)
	}
}

func collectBasicImport(t *types.Basic, imports map[string]string) {
	if t.Kind() == types.UnsafePointer {
		imports["unsafe"] = "unsafe"
	}
}

func collectSingleElementTypeImports(t types.Type, selfPkgPath string, imports map[string]string) {
	// Pointers, slices, arrays, and channels all expose exactly one nested type.
	// Keeping that shared recursion here leaves the main walker focused on the
	// shapes that need special handling: objects, maps, signatures, interfaces,
	// and unsafe.Pointer.
	switch tt := t.(type) {
	case *types.Pointer:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Slice:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Array:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	case *types.Chan:
		collectTypeImports(tt.Elem(), selfPkgPath, imports)
	}
}

func collectInterfaceImports(t *types.Interface, selfPkgPath string, imports map[string]string) {
	for i := 0; i < t.NumMethods(); i++ {
		collectTypeImports(t.Method(i).Type(), selfPkgPath, imports)
	}
}

func collectObjectImport(obj types.Object, selfPkgPath string, imports map[string]string) {
	pkg := obj.Pkg()
	if pkg != nil && pkg.Path() != selfPkgPath {
		imports[pkg.Path()] = sanitizePkgName(pkg.Path())
	}
}

func collectTupleImports(tuple *types.Tuple, selfPkgPath string, imports map[string]string) {
	for i := 0; i < tuple.Len(); i++ {
		collectTypeImports(tuple.At(i).Type(), selfPkgPath, imports)
	}
}
