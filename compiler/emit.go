// emit.go provides bytecode emission helpers: emit, emitJump, patch, addConstant, addType.
package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// emit appends an opcode and its operands to the current function's bytecode.
func (c *compiler) emit(op bytecode.OpCode, operands ...uint16) {
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(op))

	width := bytecode.OperandWidth(op)

	for _, operand := range operands {
		switch width {
		case 2:
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(operand>>8), byte(operand))
		case 1:
			c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(operand))
		default:
			if operand > 0xFF {
				c.currentFunc.Instructions = append(c.currentFunc.Instructions,
					byte(operand>>8), byte(operand))
			} else {
				c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(operand))
			}
		}
	}
}

// emitCallOp emits a call-family instruction with a uint16 index and a byte arg count.
// Used for OpCall, OpCallExternal, and OpGoCall which share the [op, idx_hi, idx_lo, numArgs] format.
func (c *compiler) emitCallOp(op bytecode.OpCode, funcIdx uint16, numArgs int) {
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(op),
		byte(funcIdx>>8), byte(funcIdx),
		byte(numArgs))
}

// emitClosure emits an OpClosure instruction: [OpClosure, fnIdx_hi, fnIdx_lo, numFreeVars].
// fnIdx is int because funcIndex map values are int.
func (c *compiler) emitClosure(fnIdx int, numFreeVars int) {
	idx := uint16(fnIdx)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpClosure),
		byte(idx>>8), byte(idx),
		byte(numFreeVars))
}

// emitJump emits an unconditional jump instruction.
func (c *compiler) emitJump(target *ssa.BasicBlock) {
	offset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(bytecode.OpJump), 0, 0)
	c.jumps = append(c.jumps, jumpInfo{offset: offset, targetBlock: target})
}

// patchJumps resolves jump targets with actual bytecode offsets.
func (c *compiler) patchJumps(blockOffsets map[*ssa.BasicBlock]int) {
	for _, jump := range c.jumps {
		targetOffset := blockOffsets[jump.targetBlock]
		c.currentFunc.Instructions[jump.offset+1] = byte(targetOffset >> 8)
		c.currentFunc.Instructions[jump.offset+2] = byte(targetOffset)
	}
}

// addConstant adds a value to the constant pool and returns its index.
func (c *compiler) addConstant(val any) uint16 {
	idx := len(c.constants)
	c.constants = append(c.constants, val)
	return uint16(idx)
}

// addType adds a types.Type to the type pool and returns its index.
func (c *compiler) addType(t types.Type) uint16 {
	idx := len(c.types)
	c.types = append(c.types, t)
	return uint16(idx)
}

// reversePostorder returns basic blocks in reverse postorder.
func reversePostorder(fn *ssa.Function) []*ssa.BasicBlock {
	if len(fn.Blocks) == 0 {
		return nil
	}

	visited := make(map[*ssa.BasicBlock]bool)
	var order []*ssa.BasicBlock

	var visit func(b *ssa.BasicBlock)
	visit = func(b *ssa.BasicBlock) {
		if visited[b] {
			return
		}
		visited[b] = true
		for _, succ := range b.Succs {
			visit(succ)
		}
		order = append(order, b)
	}

	visit(fn.Blocks[0])

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}
