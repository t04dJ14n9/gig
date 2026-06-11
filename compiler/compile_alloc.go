package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileAlloc compiles an Alloc instruction (variable allocation).
func (c *compiler) compileAlloc(i *ssa.Alloc) {
	if _, ok := syntheticMakeSliceArrayAlloc(i); ok {
		return
	}
	addrIdx := c.symbolTable.AllocLocal(i)
	elemType := i.Type().(*types.Pointer).Elem()
	typeIdx := c.addType(elemType)
	c.emit(bytecode.OpNew, uint16(typeIdx))
	c.emit(bytecode.OpSetLocal, uint16(addrIdx))
}
