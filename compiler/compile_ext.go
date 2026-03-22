package compiler

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
)

// compileExternalStaticCall compiles a call to an external package function.
func (ctx *funcContext) compileExternalStaticCall(i *ssa.Call, fn *ssa.Function, resultIdx int) {
	for _, arg := range i.Call.Args {
		ctx.compileValue(arg)
	}

	sig := fn.Signature
	if sig.Recv() != nil {
		methodName := fn.Name()
		if idx := strings.LastIndex(methodName, "."); idx >= 0 {
			methodName = methodName[idx+1:]
		}
		if idx := strings.LastIndex(methodName, ")"); idx >= 0 {
			rest := methodName[idx+1:]
			if len(rest) > 0 && rest[0] == '.' {
				methodName = rest[1:]
			}
		}

		methodInfo := &bytecode.ExternalMethodInfo{
			MethodName: methodName,
		}

		// Try to resolve method DirectCall at compile time
		if ctx.c.registry != nil {
			typeName := extractReceiverTypeName(sig.Recv().Type())
			if typeName != "" {
				if dc, ok := ctx.c.registry.LookupMethodDirectCall(typeName, methodName); ok {
					methodInfo.DirectCall = dc
				}
			}
		}

		funcIdx := ctx.c.addConstant(methodInfo)
		numArgs := len(i.Call.Args)
		ctx.cf.Instructions = append(ctx.cf.Instructions,
			byte(bytecode.OpCallExternal),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	// Resolve external function
	var extFuncInfo *bytecode.ExternalFuncInfo
	if fn.Pkg != nil && ctx.c.registry != nil {
		pkgPath := fn.Pkg.Pkg.Path()
		if fnVal, directCall, ok := ctx.c.registry.LookupExternalFunc(pkgPath, fn.Name()); ok {
			extFuncInfo = &bytecode.ExternalFuncInfo{
				Func:       fnVal,
				DirectCall: directCall,
			}
		}
	}

	if extFuncInfo == nil {
		extFuncInfo = &bytecode.ExternalFuncInfo{
			Func:       fn,
			DirectCall: nil,
		}
	}

	funcIdx := ctx.c.addConstant(extFuncInfo)
	numArgs := len(i.Call.Args)
	ctx.cf.Instructions = append(ctx.cf.Instructions,
		byte(bytecode.OpCallExternal),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndirectCall compiles an indirect call (closure or function value).
func (ctx *funcContext) compileIndirectCall(i *ssa.Call) {
	resultIdx := ctx.symbolTable.AllocLocal(i)

	ctx.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		ctx.compileValue(arg)
	}

	numArgs := len(i.Call.Args)
	ctx.emit(bytecode.OpCallIndirect, uint16(numArgs))
	ctx.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// extractReceiverTypeName extracts the package-path-qualified type name from a receiver type.
// For pointer receivers like *Reader, it unwraps the pointer.
// Returns "pkgPath.TypeName" (e.g., "encoding/json.Encoder") for use as a DirectCall lookup key.
func extractReceiverTypeName(recvType types.Type) string {
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}
	named, ok := recvType.(*types.Named)
	if !ok {
		return ""
	}
	obj := named.Obj()
	pkg := obj.Pkg()
	if pkg != nil {
		return pkg.Path() + "." + obj.Name()
	}
	return obj.Name()
}

// extractNamedType unwraps pointer types to find the underlying *types.Named type.
// Returns nil if the type is not named (e.g., interface types without a name).
func extractNamedType(t types.Type) *types.Named {
	for {
		switch tt := t.(type) {
		case *types.Named:
			return tt
		case *types.Pointer:
			t = tt.Elem()
		default:
			return nil
		}
	}
}

// compileReturn compiles a Return instruction.
func (ctx *funcContext) compileReturn(i *ssa.Return) {
	if len(i.Results) == 0 {
		ctx.emit(bytecode.OpReturn)
		return
	}

	if len(i.Results) == 1 {
		ctx.compileValue(i.Results[0])
		ctx.emit(bytecode.OpReturnVal)
		return
	}

	for _, result := range i.Results {
		ctx.compileValue(result)
	}

	ctx.emit(bytecode.OpPack, uint16(len(i.Results)))
	ctx.emit(bytecode.OpReturnVal)
}
