package gentool

import "go/types"

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

		// Cross-package named types: allow if we can extract via .Interface().(Type).
		return canWrapType(t.Underlying(), true)
	}

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
// If crossPkg is true, all representable types are allowed and extracted via Interface().
func canWrapType(t types.Type, crossPkg bool) bool {
	switch ut := t.(type) {
	case *types.Basic:
		return ut.Kind() != types.UnsafePointer && ut.Kind() != types.Invalid
	case *types.Slice:
		if crossPkg {
			return true
		}
		_, ok := ut.Elem().Underlying().(*types.Basic)
		return ok
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
	case *types.Struct, *types.Map, *types.Chan, *types.Signature, *types.Array:
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

// isNamedOrAlias returns true if the type is a package-defined named type or alias.
func isNamedOrAlias(t types.Type) bool {
	switch tt := t.(type) {
	case *types.Named:
		return tt.Obj().Pkg() != nil
	case *types.Alias:
		return tt.Obj().Pkg() != nil
	default:
		return false
	}
}
