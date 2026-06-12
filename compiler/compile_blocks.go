package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func (c *compiler) compileBlocks(fn *ssa.Function) {
	// Reverse postorder keeps emitted blocks close to source control flow while
	// still allowing every branch target to be patched after all offsets exist.
	blocks := reversePostorder(fn)
	blockOffsets := make(map[*ssa.BasicBlock]int)
	for _, block := range blocks {
		blockOffsets[block] = len(c.currentFunc.Instructions)
		c.compileBlock(block)
	}
	c.patchJumps(blockOffsets)
}

// compileBlock compiles a single basic block to bytecode.
func (c *compiler) compileBlock(block *ssa.BasicBlock) {
	for _, instr := range block.Instrs {
		c.compileInstruction(instr)
	}
	c.compileBlockTerminator(block)
}

func (c *compiler) compileBlockTerminator(block *ssa.BasicBlock) {
	// compileInstruction handles the semantic work for return and panic values.
	// Terminator handling here only adds CFG movement: jumps, conditional
	// branches, and the phi moves required at each control-flow edge.
	if block.Instrs == nil {
		return
	}
	last := block.Instrs[len(block.Instrs)-1]
	switch term := last.(type) {
	case *ssa.Return:
		return
	case *ssa.Jump:
		c.emitPhiMoves(block, block.Succs[0])
		c.emitJump(block.Succs[0])
	case *ssa.If:
		c.compileIfTerminator(block, term)
	case *ssa.Panic:
		c.compileValue(term.X)
		c.emit(bytecode.OpPanic)
	}
}

func (c *compiler) compileIfTerminator(block *ssa.BasicBlock, term *ssa.If) {
	// Emit the false edge first. OpJumpTrue is patched to skip over false-edge
	// phi moves, which prevents those moves from running when the condition is
	// true. Reversing this order changes phi semantics.
	c.compileValue(term.Cond)
	jumpTrueOffset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions, byte(bytecode.OpJumpTrue), 0, 0)

	c.emitPhiMoves(block, block.Succs[1])
	c.emitJump(block.Succs[1])

	trueLandingOffset := len(c.currentFunc.Instructions)
	c.currentFunc.Instructions[jumpTrueOffset+1] = byte(trueLandingOffset >> 8)
	c.currentFunc.Instructions[jumpTrueOffset+2] = byte(trueLandingOffset)

	c.emitPhiMoves(block, block.Succs[0])
	c.emitJump(block.Succs[0])
}

// emitPhiMoves emits move instructions for Phi nodes before jumping to a block.
func (c *compiler) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
	// SSA Phi edges are ordered to match targetBlock.Preds. We use that index
	// to select the value coming from predBlock and write it into the phi slot
	// before control reaches targetBlock.
	predIndex := phiPredecessorIndex(predBlock, targetBlock)
	if predIndex < 0 {
		return
	}

	for _, instr := range targetBlock.Instrs {
		phi, ok := instr.(*ssa.Phi)
		if !ok {
			break
		}
		if predIndex >= len(phi.Edges) {
			continue
		}
		sourceValue := phi.Edges[predIndex]
		targetSlot := c.phiSlots[phi]
		c.compileValue(sourceValue)
		c.emit(bytecode.OpSetLocal, uint16(targetSlot))
	}
}

func phiPredecessorIndex(predBlock, targetBlock *ssa.BasicBlock) int {
	for i, pred := range targetBlock.Preds {
		if pred == predBlock {
			return i
		}
	}
	return -1
}
