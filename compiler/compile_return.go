package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

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
