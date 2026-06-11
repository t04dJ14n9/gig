package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileIndirectCall compiles an indirect call (closure or function value).
func (c *compiler) compileIndirectCall(i *ssa.Call) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.validateExternalFuncValueBoundary(i.Call.Value, i.Call.Args)

	c.compileValue(i.Call.Value)

	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	numArgs := len(i.Call.Args)
	c.emit(bytecode.OpCallIndirect, uint16(numArgs))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
