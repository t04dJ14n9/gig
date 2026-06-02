package compiler

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"golang.org/x/tools/go/ssa"
)

// compileSend compiles a Send instruction.
func (c *compiler) compileSend(i *ssa.Send) {
	c.compileValue(i.Chan)
	c.compileValue(i.X)
	c.emit(bytecode.OpSend)
}

// compileDefer compiles a Defer instruction.
func (c *compiler) compileDefer(i *ssa.Defer) {
	// Interface method invocation (e.g., defer iface.Method())
	if i.Call.IsInvoke() {
		c.compileDeferInvoke(i)
		return
	}

	switch val := i.Call.Value.(type) {
	case *ssa.Function:
		c.compileDeferFunction(i, val)
	case *ssa.MakeClosure:
		c.compileDeferMakeClosure(i, val)
	default:
		// Other cases: compile the callable, then push args
		c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)
		c.compileValue(i.Call.Value)
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
	}
}

// compileDeferInvoke handles defer of an interface method call.
func (c *compiler) compileDeferInvoke(i *ssa.Defer) {
	c.compileValue(i.Call.Value)
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}
	methodInfo := &external.ExternalMethodInfo{
		MethodName: i.Call.Method.Name(),
		IsStdlib:   true,
	}
	if recvType := i.Call.Value.Type(); recvType != nil {
		if named := extractNamedType(recvType); named != nil {
			methodInfo.ReceiverTypeName = named.Obj().Name()
		}
	}
	funcIdx := c.addConstant(methodInfo)
	c.emitCallOp(bytecode.OpDeferExternal, uint16(funcIdx), len(i.Call.Args)+1)
}

// compileDeferFunction handles defer of a static function call.
func (c *compiler) compileDeferFunction(i *ssa.Defer, val *ssa.Function) {
	// Known internal function
	if _, known := c.funcIndex[val]; known {
		if len(val.FreeVars) > 0 {
			// Has free variables — create closure, then push args
			fnIdx := c.funcIndex[val]
			c.compileAndEmitClosureFromFreeVars(val.FreeVars, fnIdx)
			c.compileDeferCallArgs(i.Call.Args)
			c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
			return
		}
		// No free variables — push args, use OpDefer directly
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDefer, uint16(c.funcIndex[val]))
		return
	}

	c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)

	// External method wrapper (not in funcIndex)
	if val.Signature.Recv() != nil {
		c.compileDeferCallArgs(i.Call.Args)
		methodInfo := c.lookupExternalMethodInfo(val)
		funcIdx := c.addConstant(methodInfo)
		c.emitCallOp(bytecode.OpDeferExternal, uint16(funcIdx), len(i.Call.Args))
		return
	}

	// External package function
	if extFuncInfo := c.lookupExternalFuncInfo(val); extFuncInfo != nil {
		c.compileDeferCallArgs(i.Call.Args)
		funcIdx := c.addConstant(extFuncInfo)
		c.emitCallOp(bytecode.OpDeferExternal, uint16(funcIdx), len(i.Call.Args))
		return
	}

	c.compileValue(i.Call.Value)
	c.compileDeferCallArgs(i.Call.Args)
	c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
}

// compileDeferMakeClosure handles defer of a closure expression.
func (c *compiler) compileDeferMakeClosure(i *ssa.Defer, val *ssa.MakeClosure) {
	// Already compiled — load from local
	if idx, ok := c.symbolTable.GetLocal(val); ok {
		c.emit(bytecode.OpLocal, uint16(idx))
		c.compileDeferCallArgs(i.Call.Args)
		c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
		return
	}
	// Not yet compiled — create the closure now
	for _, binding := range val.Bindings {
		if fv, ok := binding.(*ssa.FreeVar); ok {
			if idx, ok := c.symbolTable.freeVars[fv]; ok {
				c.emit(bytecode.OpFree, uint16(idx))
				continue
			}
		}
		if alloc, ok := binding.(*ssa.Alloc); ok {
			if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
				c.emit(bytecode.OpLocal, uint16(slotIdx))
				continue
			}
		}
		c.compileValue(binding)
	}
	c.emitClosure(c.funcIndex[val.Fn.(*ssa.Function)], len(val.Bindings))
	c.compileDeferCallArgs(i.Call.Args)
	c.emit(bytecode.OpDeferIndirect, uint16(len(i.Call.Args)))
}

// compileGo compiles a Go instruction.
func (c *compiler) compileGo(i *ssa.Go) {
	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		// Check if this is a known internal function
		if _, known := c.funcIndex[fn]; known {
			// If the function has free variables, we need to create a closure
			// so the child goroutine can access captured variables (e.g., channels,
			// mutexes from the enclosing scope). Without this, OpGoCall passes
			// nil for freeVars and the child VM cannot access them.
			if len(fn.FreeVars) > 0 {
				// Create closure, then push arguments
				fnIdx := c.funcIndex[fn]
				c.compileAndEmitClosureFromFreeVars(fn.FreeVars, fnIdx)
				// Push arguments AFTER closure
				for _, arg := range i.Call.Args {
					c.compileValue(arg)
				}
				c.emit(bytecode.OpGoCallIndirect, uint16(len(i.Call.Args)))
				return
			}
			// No free variables — use OpGoCall directly
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			c.emitCallOp(bytecode.OpGoCall, uint16(c.funcIndex[fn]), len(i.Call.Args))
			return
		}

		c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)

		if fn.Signature.Recv() != nil {
			methodInfo := c.lookupExternalMethodInfo(fn)
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			funcIdx := c.addConstant(methodInfo)
			c.emitCallOp(bytecode.OpGoCallExternal, uint16(funcIdx), len(i.Call.Args))
			return
		}
		if extFuncInfo := c.lookupExternalFuncInfo(fn); extFuncInfo != nil {
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}
			funcIdx := c.addConstant(extFuncInfo)
			c.emitCallOp(bytecode.OpGoCallExternal, uint16(funcIdx), len(i.Call.Args))
			return
		}
	}

	// Indirect call (closure or MakeClosure result): push callee FIRST,
	// then args. OpGoCallIndirect pops args first, then callee.
	c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)
	c.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	c.emit(bytecode.OpGoCallIndirect, uint16(len(i.Call.Args)))
}
