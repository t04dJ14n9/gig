package compiler

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"golang.org/x/tools/go/ssa"
)

// compileField compiles a Field instruction.
func (c *compiler) compileField(i *ssa.Field) {
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpField, uint16(i.Field))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (c *compiler) compileFieldAddr(i *ssa.FieldAddr) {
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpFieldAddr, uint16(i.Field))
}

// compileIndex compiles an Index instruction.
func (c *compiler) compileIndex(i *ssa.Index) {
	c.compileSimpleBinaryOp(i, i.X, i.Index, bytecode.OpIndex)
}

// compileIndexAddr compiles an IndexAddr instruction.
func (c *compiler) compileIndexAddr(i *ssa.IndexAddr) {
	c.compileSimpleBinaryOp(i, i.X, i.Index, bytecode.OpIndexAddr)
}

// compileLookup compiles a Lookup instruction.
func (c *compiler) compileLookup(i *ssa.Lookup) {
	opcode := bytecode.OpIndex
	if i.CommaOk {
		opcode = bytecode.OpIndexOk
	}
	c.compileSimpleBinaryOp(i, i.X, i.Index, opcode)
}

// compileStore compiles a Store instruction.
func (c *compiler) compileStore(i *ssa.Store) {
	c.compileValue(i.Addr)
	c.compileValue(i.Val)
	c.emit(bytecode.OpSetDeref)
}
