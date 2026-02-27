package compiler

import (
	"golang.org/x/tools/go/ssa"

	"gig/bytecode"
)

// compileFunction compiles a single SSA function to bytecode.
func (c *compiler) compileFunction(fn *ssa.Function) (*bytecode.CompiledFunction, error) {
	c.currentFunc = &bytecode.CompiledFunction{
		Name:         fn.Name(),
		Instructions: make([]byte, 0),
		Source:       fn,
		NumParams:    len(fn.Params),
	}

	c.symbolTable = NewSymbolTable()
	c.jumps = nil
	c.phiSlots = make(map[*ssa.Phi]int)

	// Allocate locals for parameters
	for _, param := range fn.Params {
		c.symbolTable.AllocLocal(param)
	}

	// Allocate locals for free variables (for closures)
	for i, freeVar := range fn.FreeVars {
		c.symbolTable.freeVars[freeVar] = i
	}

	// First pass: collect Phi nodes and allocate slots for them
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if phi, ok := instr.(*ssa.Phi); ok {
				slot := c.symbolTable.AllocLocal(phi)
				c.phiSlots[phi] = slot
			}
		}
	}

	// Allocate locals for all other values in the function
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if val, ok := instr.(ssa.Value); ok {
				if _, isPhi := instr.(*ssa.Phi); !isPhi {
					if _, isAlloc := instr.(*ssa.Alloc); !isAlloc {
						c.symbolTable.AllocLocal(val)
					}
				}
			}
		}
	}

	// Pre-allocate slots for Alloc instructions too
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if _, isAlloc := instr.(*ssa.Alloc); isAlloc {
				c.symbolTable.AllocLocal(instr.(ssa.Value))
			}
		}
	}

	c.currentFunc.NumLocals = c.symbolTable.NumLocals()
	c.currentFunc.NumFreeVars = len(fn.FreeVars)

	// Compile basic blocks in reverse postorder
	blocks := reversePostorder(fn)
	blockOffsets := make(map[*ssa.BasicBlock]int)

	for _, block := range blocks {
		blockOffsets[block] = len(c.currentFunc.Instructions)
		c.compileBlock(fn, block)
	}

	// Patch jump targets
	c.patchJumps(blockOffsets)

	return c.currentFunc, nil
}

// compileBlock compiles a single basic block to bytecode.
func (c *compiler) compileBlock(fn *ssa.Function, block *ssa.BasicBlock) {
	for _, instr := range block.Instrs {
		c.compileInstruction(fn, instr)
	}

	// Handle block terminator
	if block.Instrs != nil {
		last := block.Instrs[len(block.Instrs)-1]
		switch term := last.(type) {
		case *ssa.Return:
			// Already handled in compileInstruction
		case *ssa.Jump:
			c.emitPhiMoves(block, block.Succs[0])
			c.emitJump(block.Succs[0])
		case *ssa.If:
			c.compileValue(term.Cond)
			c.emitPhiMoves(block, block.Succs[1])
			c.emitJumpFalse(block.Succs[1])
			c.emitPhiMoves(block, block.Succs[0])
			c.emitJump(block.Succs[0])
		case *ssa.Panic:
			c.compileValue(term.X)
			c.emit(bytecode.OpPanic)
		}
	}
}

// emitPhiMoves emits move instructions for Phi nodes before jumping to a block.
func (c *compiler) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
	predIndex := -1
	for i, pred := range targetBlock.Preds {
		if pred == predBlock {
			predIndex = i
			break
		}
	}
	if predIndex < 0 {
		return
	}

	for _, instr := range targetBlock.Instrs {
		phi, ok := instr.(*ssa.Phi)
		if !ok {
			break
		}

		if predIndex < len(phi.Edges) {
			sourceValue := phi.Edges[predIndex]
			targetSlot := c.phiSlots[phi]

			c.compileValue(sourceValue)
			c.emit(bytecode.OpSetLocal, uint16(targetSlot))
		}
	}
}
