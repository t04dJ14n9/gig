package compiler

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"golang.org/x/tools/go/ssa"
)

// compileTypeAssert compiles a TypeAssert instruction.
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
	typeIdx := c.addType(i.AssertedType)
	c.compileValue(i.X)
	c.emit(bytecode.OpAssert, uint16(typeIdx))

	if !i.CommaOk {
		// Non-comma-ok assertion: SSA's `typeassert t.(T)` (without comma-ok)
		// panics on failure by branching to the recover block. We must check
		// the ok value and emit OpPanic if the assertion fails.
		// Stack has: [result, ok] tuple
		// Duplicate the tuple, extract ok (index 1), check if false → panic.
		c.emit(bytecode.OpDup)                             // [tuple, tuple]
		c.emit(bytecode.OpConst, uint16(c.addConstant(1))) // [tuple, tuple, 1]
		c.emit(bytecode.OpIndex)                           // [tuple, ok]
		// Emit OpJumpTrue with placeholder offset (3 bytes: opcode + u16 offset)
		jumpTrueOffset := len(c.currentFunc.Instructions)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpJumpTrue), 0, 0)
		// ok was false — panic with a type assertion error
		c.emit(bytecode.OpConst, uint16(c.addConstant("interface conversion: type assertion failed")))
		c.emit(bytecode.OpPanic)
		// Patch the JumpTrue to land here (ok case)
		skipOffset := len(c.currentFunc.Instructions)
		c.currentFunc.Instructions[jumpTrueOffset+1] = byte(skipOffset >> 8)
		c.currentFunc.Instructions[jumpTrueOffset+2] = byte(skipOffset)
		// Stack still has: [tuple] — extract the value
		c.emit(bytecode.OpConst, uint16(c.addConstant(0)))
		c.emit(bytecode.OpIndex)
	}

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeInterface compiles a ChangeInterface instruction.
func (c *compiler) compileChangeInterface(i *ssa.ChangeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeType compiles a ChangeType instruction.
// ChangeType converts between types with identical underlying types (e.g., []int -> sort.IntSlice).
// We emit OpChangeType which carries both the target type and the source local index,
// so the VM can update the source variable to share the same backing array after conversion.
func (c *compiler) compileChangeType(i *ssa.ChangeType) {
	resultIdx := c.symbolTable.AllocLocal(i)
	typeIdx := c.addType(i.Type())

	// Try to find the source local index. If the source is a local variable,
	// we pass its index so the VM can update it for slice aliasing.
	srcLocalIdx := uint16(bytecode.NoSourceLocal)
	if srcIdx, ok := c.symbolTable.GetLocal(i.X); ok {
		srcLocalIdx = uint16(srcIdx)
	}

	c.compileValue(i.X)
	// Emit OpChangeType with 4 bytes of operands: type_idx(2) + src_local(2)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpChangeType),
		byte(typeIdx>>8), byte(typeIdx),
		byte(srcLocalIdx>>8), byte(srcLocalIdx))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileConvert compiles a Convert instruction.
func (c *compiler) compileConvert(i *ssa.Convert) {
	resultIdx := c.symbolTable.AllocLocal(i)
	typeIdx := c.addType(i.Type())
	c.compileValue(i.X)
	c.emit(bytecode.OpConvert, uint16(typeIdx))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileExtract compiles an Extract instruction.
func (c *compiler) compileExtract(i *ssa.Extract) {
	c.compileValue(i.Tuple)
	c.emit(bytecode.OpConst, uint16(c.addConstant(i.Index)))
	c.emit(bytecode.OpIndex)
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
