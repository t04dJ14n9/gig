// compile_value.go routes SSA values to the value-lowering domain helpers.
package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileValue compiles an SSA value to push it onto the stack.
func (c *compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		c.compileConst(val)
	case *ssa.Function:
		if fnIdx, ok := c.funcIndex[val]; ok {
			c.emitClosure(fnIdx, 0)
		} else {
			// External function not in funcIndex — look up the actual Go function
			// and store it as a constant so it can be used as a value (e.g., passed
			// as a callback argument). OpCallIndirect handles reflect.Func values.
			c.compileExternalFuncValue(val)
		}
	case *ssa.Phi:
		if slot, ok := c.phiSlots[val]; ok {
			c.emit(bytecode.OpLocal, uint16(slot))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.FreeVar:
		if idx, ok := c.symbolTable.freeVars[val]; ok {
			c.emit(bytecode.OpFree, uint16(idx))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.Global:
		c.compileGlobalValue(val)
	default:
		if idx, ok := c.symbolTable.GetLocal(v); ok {
			c.emit(bytecode.OpLocal, uint16(idx))
			return
		}
		if idx, ok := c.symbolTable.freeVars[v]; ok {
			c.emit(bytecode.OpFree, uint16(idx))
			return
		}
		c.emit(bytecode.OpNil)
	}
}
