package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/external"
)

func externalBoundaryArgType(arg ssa.Value) types.Type {
	switch v := arg.(type) {
	case *ssa.MakeInterface:
		return externalBoundaryArgType(v.X)
	case *ssa.ChangeInterface:
		return externalBoundaryArgType(v.X)
	case *ssa.Convert:
		return externalBoundaryArgType(v.X)
	default:
		return arg.Type()
	}
}

func callTargetType(sig *types.Signature, argIndex int) types.Type {
	if sig == nil || sig.Params() == nil || argIndex < 0 {
		return nil
	}
	if sig.Recv() != nil {
		if argIndex == 0 {
			return sig.Recv().Type()
		}
		return callParamTargetType(sig, argIndex-1)
	}
	return callParamTargetType(sig, argIndex)
}

func callParamTargetType(sig *types.Signature, argIndex int) types.Type {
	if sig == nil || sig.Params() == nil || argIndex < 0 {
		return nil
	}
	params := sig.Params()
	if argIndex < params.Len() {
		return params.At(argIndex).Type()
	}
	if sig.Variadic() && params.Len() > 0 {
		return params.At(params.Len() - 1).Type()
	}
	return nil
}

func (c *compiler) externalTargetAllowsInterfaceProxy(sig *types.Signature, argIndex int) bool {
	return c.externalTargetHasInterfaceProxy(callTargetType(sig, argIndex))
}

func (c *compiler) externalTargetHasInterfaceProxy(t types.Type) bool {
	if t == nil {
		return false
	}
	iface, ok := t.Underlying().(*types.Interface)
	if !ok || iface.NumMethods() == 0 {
		return false
	}
	if c.lookup == nil {
		return false
	}
	nameLookup, hasNameLookup := c.lookup.(interface {
		LookupInterfaceProxy(pkgPath string, typeName string) (*external.InterfaceProxyInfo, bool)
	})
	if !hasNameLookup {
		return false
	}
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil {
			if _, ok := nameLookup.LookupInterfaceProxy(obj.Pkg().Path(), obj.Name()); ok {
				return true
			}
		}
	}
	return false
}

func (c *compiler) validateExternalStaticCallBoundary(pkgPath, funcName string, sig *types.Signature, args []ssa.Value) {
	if c.compileErr != nil {
		return
	}
	callArgs := make([]externalCallArg, len(args))
	for idx, arg := range args {
		callArgs[idx] = externalCallArg{
			SourceType:          externalBoundaryArgType(arg),
			AllowInterfaceProxy: c.externalTargetAllowsInterfaceProxy(sig, idx),
		}
	}
	if err := validateExternalCallBoundary(pkgPath, funcName, callArgs); err != nil {
		c.compileErr = err
	}
}
