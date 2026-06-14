// host_iface_check.go enforces the interpreter's interpreted-struct ↔
// host-interface boundary rule. The interpreter does not synthesise
// host-interface proxies for user-defined types declared in interpreted
// source. Instead of letting such programs explode at runtime with
// confusing reflect panics, the frontend rejects them at type-check
// time with a clear diagnostic.
//
// The rule, precisely:
//
//   - For every Call expression whose callee is a function declared in
//     a host package (i.e. its types.Func.Pkg() != source pkg), look at
//     each parameter whose type is a non-empty interface.
//   - If the corresponding argument's static type is an interpreted
//     struct or interpreted-struct pointer (declared in source pkg),
//     reject.
//
// Empty interface (any) is allowed because it carries no required
// methods. Method values, type assertions, and conversions remain
// unaffected; this rule only fires on host-call boundaries.
package frontend

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/t04dJ14n9/gig/diag"
)

// checkHostInterfaceBoundary walks the AST and reports the first call
// that violates the boundary rule. Returns nil if everything checks.
func checkHostInterfaceBoundary(fset *token.FileSet, file *ast.File, info *types.Info, srcPkg *types.Package) error {
	var firstErr error
	ast.Inspect(file, func(n ast.Node) bool {
		if firstErr != nil {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		// Resolve the callee to a *types.Func in a host package.
		callee := calleeFunc(info, call)
		if callee == nil {
			return true
		}
		if callee.Pkg() == nil || callee.Pkg() == srcPkg {
			return true // interpreted call — fine
		}
		sig, ok := callee.Type().(*types.Signature)
		if !ok {
			return true
		}
		params := sig.Params()
		variadic := sig.Variadic()
		for i, arg := range call.Args {
			pi := i
			if variadic && i >= params.Len()-1 {
				pi = params.Len() - 1
			}
			if pi >= params.Len() {
				break
			}
			paramT := params.At(pi).Type()
			if variadic && pi == params.Len()-1 {
				if slice, ok := paramT.(*types.Slice); ok {
					paramT = slice.Elem()
				}
			}
			iface, ok := paramT.Underlying().(*types.Interface)
			if !ok || iface.NumMethods() == 0 {
				continue // not a non-empty interface
			}
			argTV, ok := info.Types[arg]
			if !ok {
				continue
			}
			if isInterpretedConcrete(argTV.Type, srcPkg) {
				pos := fset.Position(call.Pos())
				firstErr = fmt.Errorf(
					"frontend: cannot pass interpreted type %s to host parameter %s of type %s (%s); "+
						"interpreted types cannot satisfy host interfaces (G_iface_ban)",
					argTV.Type, params.At(pi).Name(), paramT, pos,
				)
				return false
			}
		}
		return true
	})
	if firstErr != nil {
		return &BuildError{Diags: []diag.Diagnostic{{
			Severity: diag.SeverityError,
			Message:  firstErr.Error(),
		}}}
	}
	return nil
}

// calleeFunc returns the declared function object for a call, or nil if
// the callee is not a static function reference (e.g. method value,
// closure, dynamic dispatch, builtin).
func calleeFunc(info *types.Info, call *ast.CallExpr) *types.Func {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		if obj := info.Uses[fn]; obj != nil {
			if f, ok := obj.(*types.Func); ok {
				return f
			}
		}
	case *ast.SelectorExpr:
		if sel, ok := info.Selections[fn]; ok {
			if f, ok := sel.Obj().(*types.Func); ok {
				return f
			}
		}
		if obj := info.Uses[fn.Sel]; obj != nil {
			if f, ok := obj.(*types.Func); ok {
				return f
			}
		}
	}
	return nil
}

// isInterpretedConcrete reports whether t is a concrete type (struct,
// pointer-to-struct, named, etc.) declared in the source package. Empty
// interfaces, host-package types, and built-in types do not count.
func isInterpretedConcrete(t types.Type, srcPkg *types.Package) bool {
	if t == nil {
		return false
	}
	if iface, ok := t.Underlying().(*types.Interface); ok {
		// An interpreted-defined interface flowing into a host
		// interface is also a problem unless it's empty.
		_ = iface
		// We allow interfaces through — only concrete types matter.
		return false
	}
	switch tt := t.(type) {
	case *types.Named:
		return tt.Obj() != nil && tt.Obj().Pkg() == srcPkg
	case *types.Pointer:
		return isInterpretedConcrete(tt.Elem(), srcPkg)
	case *types.Struct:
		// Anonymous struct literal — treat as interpreted (it has no
		// home package and was defined in source).
		return true
	}
	return false
}
