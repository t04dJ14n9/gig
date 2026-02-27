package compiler

import (
	"strings"

	"golang.org/x/tools/go/ssa"

	"gig/bytecode"
)

// compileExternalCall compiles an external function call (non-static / unknown callee).
func (c *compiler) compileExternalCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	funcIdx := c.addConstant(i.Call.Value)
	numArgs := len(i.Call.Args)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpCallExternal),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileExternalStaticCall compiles a call to an external package function.
// It uses the injected PackageLookup to resolve the function, avoiding direct importer dependency.
func (c *compiler) compileExternalStaticCall(i *ssa.Call, fn *ssa.Function, resultIdx int) {
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
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
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpCallExternal),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	// Use injected PackageLookup instead of direct importer access
	var extFuncInfo *bytecode.ExternalFuncInfo
	if fn.Pkg != nil && c.lookup != nil {
		pkgPath := fn.Pkg.Pkg.Path()
		if fnVal, directCall, ok := c.lookup.LookupExternalFunc(pkgPath, fn.Name()); ok {
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

	funcIdx := c.addConstant(extFuncInfo)
	numArgs := len(i.Call.Args)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpCallExternal),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndirectCall compiles an indirect call (closure or function value).
func (c *compiler) compileIndirectCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	numArgs := len(i.Call.Args)
	c.emit(bytecode.OpCallIndirect, uint16(numArgs))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileReturn compiles a Return instruction.
func (c *compiler) compileReturn(i *ssa.Return) {
	if len(i.Results) == 0 {
		c.emit(bytecode.OpReturn)
		return
	}

	if len(i.Results) == 1 {
		c.compileValue(i.Results[0])
		c.emit(bytecode.OpReturnVal)
		return
	}

	for _, result := range i.Results {
		c.compileValue(result)
	}

	c.emit(bytecode.OpPack, uint16(len(i.Results)))
	c.emit(bytecode.OpReturnVal)
}
