package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/compiler/optimize"
)

// isIntType returns true if the type is a signed integer (int, int8..int64).
func isIntType(t types.Type) bool {
	if t == nil {
		return false
	}
	basic, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return true
	}
	return false
}

// isIntSliceType returns true if the type is []int or []int64 (matches native int slice fast path).
func isIntSliceType(t types.Type) bool {
	if t == nil {
		return false
	}
	sl, ok := t.Underlying().(*types.Slice)
	if !ok {
		return false
	}
	elem := sl.Elem()
	if elem == nil {
		return false
	}
	basic, ok := elem.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Int, types.Int64:
		return true
	}
	return false
}

// compileFunction compiles a single SSA function to bytecode.
func (c *compiler) compileFunction(fn *ssa.Function) (*bytecode.CompiledFunction, error) {
	ctx := &funcContext{
		cf: &bytecode.CompiledFunction{
			Name:         fn.Name(),
			Instructions: make([]byte, 0),
			Source:       fn,
			NumParams:    len(fn.Params),
		},
		symbolTable: NewSymbolTable(),
		jumps:       nil,
		phiSlots:    make(map[*ssa.Phi]int),
		c:           c,
	}

	// Allocate locals for parameters
	for _, param := range fn.Params {
		ctx.symbolTable.AllocLocal(param)
	}

	// Allocate locals for free variables (for closures)
	for i, freeVar := range fn.FreeVars {
		ctx.symbolTable.freeVars[freeVar] = i
	}

	// First pass: collect Phi nodes and allocate slots for them
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if phi, ok := instr.(*ssa.Phi); ok {
				slot := ctx.symbolTable.AllocLocal(phi)
				ctx.phiSlots[phi] = slot
			}
		}
	}

	// Allocate locals for all other values in the function
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if val, ok := instr.(ssa.Value); ok {
				if _, isPhi := instr.(*ssa.Phi); !isPhi {
					if _, isAlloc := instr.(*ssa.Alloc); !isAlloc {
						ctx.symbolTable.AllocLocal(val)
					}
				}
			}
		}
	}

	// Pre-allocate slots for Alloc instructions too
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if _, isAlloc := instr.(*ssa.Alloc); isAlloc {
				ctx.symbolTable.AllocLocal(instr.(ssa.Value))
			}
		}
	}

	ctx.cf.NumLocals = ctx.symbolTable.NumLocals()
	ctx.cf.NumFreeVars = len(fn.FreeVars)

	// Build local-is-int map for int-specialization
	localIsInt := make([]bool, ctx.symbolTable.NumLocals())
	localIsIntSlice := make([]bool, ctx.symbolTable.NumLocals())
	for v, idx := range ctx.symbolTable.locals {
		if isIntType(v.Type()) {
			localIsInt[idx] = true
		}
		if isIntSliceType(v.Type()) {
			localIsIntSlice[idx] = true
		}
	}

	// Compile basic blocks in reverse postorder
	blocks := reversePostorder(fn)
	blockOffsets := make(map[*ssa.BasicBlock]int)

	for _, block := range blocks {
		blockOffsets[block] = len(ctx.cf.Instructions)
		ctx.compileBlock(block)
	}

	// Patch jump targets
	ctx.patchJumps(blockOffsets)

	// Build const-is-int map (must recognize all integer types stored by compileConst)
	constIsInt := make([]bool, len(c.constants))
	for i, k := range c.constants {
		switch k.(type) {
		case int, int8, int16, int32, int64:
			constIsInt[i] = true
		}
	}

	// Apply all optimization passes (peephole, slice fusion, int-specialization, int move fusion)
	ctx.cf.Instructions, ctx.cf.HasIntLocals = optimize.Optimize(ctx.cf.Instructions, localIsInt, constIsInt, localIsIntSlice)

	return ctx.cf, nil
}

// compileBlock compiles a single basic block to bytecode.
func (ctx *funcContext) compileBlock(block *ssa.BasicBlock) {
	for _, instr := range block.Instrs {
		ctx.compileInstruction(instr)
	}

	// Handle block terminator
	if block.Instrs != nil {
		last := block.Instrs[len(block.Instrs)-1]
		switch term := last.(type) {
		case *ssa.Return:
			// Already handled in compileInstruction
		case *ssa.Jump:
			ctx.emitPhiMoves(block, block.Succs[0])
			ctx.emitJump(block.Succs[0])
		case *ssa.If:
			ctx.compileValue(term.Cond)
			// Emit OpJumpTrue with a placeholder offset — we'll patch it after
			// emitting the false-branch phi moves so they only execute when
			// the condition is false.
			jumpTrueOffset := len(ctx.cf.Instructions)
			ctx.cf.Instructions = append(ctx.cf.Instructions,
				byte(bytecode.OpJumpTrue), 0, 0)

			// False branch: phi moves + jump (only reached when condition is false)
			ctx.emitPhiMoves(block, block.Succs[1])
			ctx.emitJump(block.Succs[1])

			// Patch the OpJumpTrue to land here (true branch)
			trueLandingOffset := len(ctx.cf.Instructions)
			ctx.cf.Instructions[jumpTrueOffset+1] = byte(trueLandingOffset >> 8)
			ctx.cf.Instructions[jumpTrueOffset+2] = byte(trueLandingOffset)

			// True branch: phi moves + jump (only reached when condition is true)
			ctx.emitPhiMoves(block, block.Succs[0])
			ctx.emitJump(block.Succs[0])
		case *ssa.Panic:
			ctx.compileValue(term.X)
			ctx.emit(bytecode.OpPanic)
		}
	}
}

// emitPhiMoves emits move instructions for Phi nodes before jumping to a block.
func (ctx *funcContext) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
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
			targetSlot := ctx.phiSlots[phi]

			ctx.compileValue(sourceValue)
			ctx.emit(bytecode.OpSetLocal, uint16(targetSlot))
		}
	}
}
