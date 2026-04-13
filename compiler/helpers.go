// helpers.go provides helper functions for common instruction compilation patterns.
package compiler

import (
	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/model/bytecode"
)

// compileSimpleInstruction handles the common pattern for simple instructions:
// 1. Allocate a local slot for the result
// 2. Compile operand value(s)
// 3. Emit the opcode (with optional operands)
// 4. Emit OpSetLocal to store the result
//
// This consolidates the pattern used by Field, FieldAddr, Index, IndexAddr,
// Range, Next, and similar instructions.
func (c *compiler) compileSimpleInstruction(
	resultValue ssa.Value,
	operands []ssa.Value,
	opcode bytecode.OpCode,
	emitOperands ...uint16,
) {
	resultIdx := c.symbolTable.AllocLocal(resultValue)

	// Compile all operand values in order
	for _, operand := range operands {
		c.compileValue(operand)
	}

	// Emit the operation with its operands
	args := make([]uint16, len(emitOperands))
	copy(args, emitOperands)
	c.emit(opcode, args...)

	// Store the result
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSimpleUnaryOp is a specialized helper for unary operations (single operand).
// This includes: Field, FieldAddr, Range, Next, MakeInterface
func (c *compiler) compileSimpleUnaryOp(
	resultValue ssa.Value,
	operand ssa.Value,
	opcode bytecode.OpCode,
	emitOperands ...uint16,
) {
	c.compileSimpleInstruction(resultValue, []ssa.Value{operand}, opcode, emitOperands...)
}

// compileSimpleBinaryOp is a specialized helper for binary operations (two operands).
// This includes: Index, IndexAddr, Lookup
func (c *compiler) compileSimpleBinaryOp(
	resultValue ssa.Value,
	operand1, operand2 ssa.Value,
	opcode bytecode.OpCode,
	emitOperands ...uint16,
) {
	c.compileSimpleInstruction(resultValue, []ssa.Value{operand1, operand2}, opcode, emitOperands...)
}
