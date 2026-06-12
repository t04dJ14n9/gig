package vm

import "go/types"

func isErrorInterface(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		return obj != nil && obj.Pkg() == nil && obj.Name() == "error"
	}
	if iface, ok := t.Underlying().(*types.Interface); ok && iface.NumMethods() == 1 {
		method := iface.Method(0)
		sig, ok := method.Type().(*types.Signature)
		return ok &&
			method.Name() == "Error" &&
			sig.Params().Len() == 0 &&
			sig.Results().Len() == 1 &&
			types.Identical(sig.Results().At(0).Type(), types.Typ[types.String])
	}
	return false
}

func isStructLikeType(t types.Type) bool {
	for {
		if t == nil {
			return false
		}
		switch tt := t.(type) {
		case *types.Named:
			_, ok := tt.Underlying().(*types.Struct)
			return ok
		case *types.Pointer:
			t = tt.Elem()
		default:
			_, ok := tt.Underlying().(*types.Struct)
			return ok
		}
	}
}

func isHostCallbackInterface(t types.Type) bool {
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj == nil || obj.Pkg() == nil {
			return false
		}
		pkgPath := obj.Pkg().Path()
		return obj.Name() == "Interface" && (pkgPath == "sort" || pkgPath == "container/heap")
	}
	if iface, ok := t.(*types.Interface); ok {
		return hasInterfaceMethods(iface, "Len", "Less", "Swap")
	}
	return false
}

func hasInterfaceMethods(iface *types.Interface, names ...string) bool {
	if iface.NumMethods() < len(names) {
		return false
	}
	for _, name := range names {
		found := false
		for i := 0; i < iface.NumMethods(); i++ {
			if iface.Method(i).Name() == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
