package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
)

// compileCall compiles a function call instruction.
func (c *compiler) compileCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	// Handle interface method invocation (e.g., iface.Method())
	// SSA represents this with IsInvoke() == true, where Call.Value is the interface value
	// and Call.Method is the method being invoked.
	if i.Call.IsInvoke() {
		// Push the receiver (interface value) as the first argument
		c.compileValue(i.Call.Value)
		// Push remaining arguments
		for _, arg := range i.Call.Args {
			c.compileValue(arg)
		}
		methodInfo := &external.ExternalMethodInfo{
			MethodName: i.Call.Method.Name(),
			IsStdlib:   true,
		}
		// For invoke calls, try to extract the concrete receiver type from the
		// interface value. This helps callCompiledMethod disambiguate methods
		// when multiple types define the same method name (e.g., Get, Add).
		if recvType := i.Call.Value.Type(); recvType != nil {
			// The receiver is an interface type; the concrete type is unknown statically.
			// However, the interface's method set constrains which types are valid.
			// We store the interface type name as a hint for fallback dispatch.
			if iface, ok := recvType.Underlying().(*types.Interface); ok {
				_ = iface // interface type available for future use
			}
			// If the value itself has a known concrete type (rare for invoke), use it.
			if named := extractNamedType(recvType); named != nil {
				methodInfo.ReceiverTypeName = named.Obj().Name()
			}
		}
		funcIdx := c.addConstant(methodInfo)
		numArgs := len(i.Call.Args) + 1 // +1 for receiver
		c.emitCallOp(bytecode.OpCallExternal, funcIdx, numArgs)
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	if builtin, ok := i.Call.Value.(*ssa.Builtin); ok {
		c.compileBuiltinCall(builtin, i.Call.Args, resultIdx)
		return
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		if _, known := c.funcIndex[fn]; known {
			c.validateExternalFuncValueBoundary(fn, i.Call.Args)
			for _, arg := range i.Call.Args {
				c.compileValue(arg)
			}

			funcIdx := c.funcIndex[fn]
			numArgs := len(i.Call.Args)
			c.emitCallOp(bytecode.OpCall, uint16(funcIdx), numArgs)
			c.emit(bytecode.OpSetLocal, uint16(resultIdx))
			return
		}

		c.compileExternalStaticCall(i, fn, resultIdx)
		return
	}

	c.compileIndirectCall(i)
}
