package gentool

import (
	"go/ast"
	"go/types"
)

type packageSymbols struct {
	Funcs  []*funcInfo
	Consts []*constInfo
	Vars   []*varInfo
	Types  []*typeInfo
}

func (s packageSymbols) empty() bool {
	return len(s.Funcs) == 0 && len(s.Consts) == 0 && len(s.Vars) == 0 && len(s.Types) == 0
}

func collectPackageSymbols(scope *types.Scope, pkgRef string) packageSymbols {
	// Only exported, non-generic, non-alias symbols are emitted. This mirrors
	// what external Go packages make available through normal imports and keeps
	// generated registration files compileable without type parameters.
	var symbols packageSymbols
	for _, name := range scope.Names() {
		if !ast.IsExported(name) {
			continue
		}
		obj := scope.Lookup(name)

		switch o := obj.(type) {
		case *types.Func:
			if fi := newFuncInfo(name, o, pkgRef); fi != nil {
				symbols.Funcs = append(symbols.Funcs, fi)
			}
		case *types.Const:
			symbols.Consts = append(symbols.Consts, &constInfo{Name: name, Obj: o})
		case *types.Var:
			symbols.Vars = append(symbols.Vars, &varInfo{Name: name, Obj: o})
		case *types.TypeName:
			if shouldRegisterTypeName(o) {
				symbols.Types = append(symbols.Types, &typeInfo{Name: name, Obj: o})
			}
		}
	}
	return symbols
}

func newFuncInfo(name string, fn *types.Func, pkgRef string) *funcInfo {
	sig := fn.Type().(*types.Signature)
	if sig.TypeParams().Len() > 0 {
		return nil
	}
	fi := &funcInfo{Name: name, Sig: sig}
	fi.DirectCall = generateDirectCall(fi, pkgRef)
	return fi
}

func shouldRegisterTypeName(typeName *types.TypeName) bool {
	// Generic aliases and constraint-like interfaces do not have a stable
	// runtime reflect.Type registration in the generated package model.
	if typeName.IsAlias() {
		return false
	}
	if named, ok := typeName.Type().(*types.Named); ok {
		if named.TypeParams().Len() > 0 {
			return false
		}
	}
	if iface, ok := typeName.Type().Underlying().(*types.Interface); ok {
		return iface.IsMethodSet()
	}
	return true
}

type funcInfo struct {
	Name       string
	Sig        *types.Signature
	DirectCall string
}

type constInfo struct {
	Name string
	Obj  *types.Const
}

type varInfo struct {
	Name string
	Obj  *types.Var
}

type typeInfo struct {
	Name string
	Obj  *types.TypeName
}
