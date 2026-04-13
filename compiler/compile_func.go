// compile_func.go handles per-function SSA→bytecode compilation: blocks, phis, locals.
package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/compiler/optimize"
	"git.woa.com/youngjin/gig/model/bytecode"
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

// buildTypeMap constructs a boolean map for a set of SSA values, marking entries
// where the given predicate returns true. This consolidates the pattern used for
// type-specialization maps like localIsInt and localIsIntSlice.
func buildTypeMap(values map[ssa.Value]int, predicate func(types.Type) bool) []bool {
	result := make([]bool, len(values))
	for v, idx := range values {
		if predicate(v.Type()) {
			result[idx] = true
		}
	}
	return result
}

// compileFunction compiles a single SSA function to bytecode.
func (c *compiler) compileFunction(fn *ssa.Function) (*bytecode.CompiledFunction, error) { //nolint:unparam // error return reserved for future compilation errors
	cf := &bytecode.CompiledFunction{
		Name:         fn.Name(),
		Instructions: make([]byte, 0),
		NumParams:    len(fn.Params),
		FuncIdx:      c.funcIndex[fn],
	}

	// Populate method dispatch metadata from SSA signature.
	sig := fn.Signature
	if sig.Recv() != nil {
		cf.HasReceiver = true
		cf.ReceiverTypeName = extractReceiverShortName(sig.Recv().Type())
	}

	c.currentFunc = cf

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

	// Single-pass allocation for Phi, value, and Alloc instructions
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			switch instr := instr.(type) {
			case *ssa.Phi:
				slot := c.symbolTable.AllocLocal(instr)
				c.phiSlots[instr] = slot
			case *ssa.Alloc:
				c.symbolTable.AllocLocal(instr)
			case ssa.Value:
				// Allocate for all other value-producing instructions
				c.symbolTable.AllocLocal(instr)
			}
		}
	}

	c.currentFunc.NumLocals = c.symbolTable.NumLocals()
	c.currentFunc.NumFreeVars = len(fn.FreeVars)

	// Detect result Alloc slots for panic-recovery return value reconstruction.
	// In SSA, named return values and variables captured by defer closures are
	// represented as Alloc instructions. The return path loads them via
	// OpLocal → OpDeref → OpReturnVal. We record their local slot indices so
	// the VM can deref them after panic recovery instead of returning nil.
	c.currentFunc.ResultAllocSlots = detectResultAllocSlots(fn, c.symbolTable)

	// Build local type maps for int-specialization (single pass)
	localIsInt := buildTypeMap(c.symbolTable.locals, isIntType)
	localIsIntSlice := buildTypeMap(c.symbolTable.locals, isIntSliceType)

	// Compile basic blocks in reverse postorder
	blocks := reversePostorder(fn)
	blockOffsets := make(map[*ssa.BasicBlock]int)

	for _, block := range blocks {
		blockOffsets[block] = len(c.currentFunc.Instructions)
		c.compileBlock(block)
	}

	// Patch jump targets
	c.patchJumps(blockOffsets)

	// Build const-is-int map (must recognize all integer types stored by compileConst)
	constIsInt := make([]bool, len(c.constants))
	for i, k := range c.constants {
		switch k.(type) {
		case int, int8, int16, int32, int64:
			constIsInt[i] = true
		}
	}

	// Apply all optimization passes (peephole, slice fusion, int-specialization, int move fusion)
	c.currentFunc.Instructions, c.currentFunc.HasIntLocals = optimize.Optimize(c.currentFunc.Instructions, localIsInt, constIsInt, localIsIntSlice)

	return c.currentFunc, nil
}

// compileBlock compiles a single basic block to bytecode.
func (c *compiler) compileBlock(block *ssa.BasicBlock) {
	for _, instr := range block.Instrs {
		c.compileInstruction(instr)
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
			// Emit OpJumpTrue with a placeholder offset — we'll patch it after
			// emitting the false-branch phi moves so they only execute when
			// the condition is false.
			jumpTrueOffset := len(c.currentFunc.Instructions)
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpJumpTrue), 0, 0)

			// False branch: phi moves + jump (only reached when condition is false)
			c.emitPhiMoves(block, block.Succs[1])
			c.emitJump(block.Succs[1])

			// Patch the OpJumpTrue to land here (true branch)
			trueLandingOffset := len(c.currentFunc.Instructions)
			c.currentFunc.Instructions[jumpTrueOffset+1] = byte(trueLandingOffset >> 8)
			c.currentFunc.Instructions[jumpTrueOffset+2] = byte(trueLandingOffset)

			// True branch: phi moves + jump (only reached when condition is true)
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

// detectResultAllocSlots finds Alloc instructions that correspond to named
// return values. In Go, only named return variables can be modified by deferred
// functions during panic recovery. Unnamed returns use the zero value after recovery.
//
// These are identified as Allocs whose value (or a deref thereof) appears
// directly in a Return instruction's results.
func detectResultAllocSlots(fn *ssa.Function, st *SymbolTable) []int {
	if fn.Blocks == nil {
		return nil
	}

	// Collect all Alloc instructions.
	allocSet := make(map[ssa.Value]bool)
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if _, ok := instr.(*ssa.Alloc); ok {
				allocSet[instr.(ssa.Value)] = true
			}
		}
	}
	if len(allocSet) == 0 {
		return nil
	}

	// Find Allocs referenced in Return instructions (named return variables).
	// SSA represents named returns as Alloc → Store → UnOp(deref) → Return.
	// The Return may reference the Alloc directly or via an UnOp deref.
	slotSet := make(map[int]bool)
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok {
				continue
			}
			for _, result := range ret.Results {
				if allocSet[result] {
					if idx, ok := st.GetLocal(result); ok {
						slotSet[idx] = true
					}
				}
				if unop, ok := result.(*ssa.UnOp); ok {
					if allocSet[unop.X] {
						if idx, ok := st.GetLocal(unop.X); ok {
							slotSet[idx] = true
						}
					}
				}
			}
		}
	}

	if len(slotSet) == 0 {
		return nil
	}

	slots := make([]int, 0, len(slotSet))
	for idx := range slotSet {
		slots = append(slots, idx)
	}
	return slots
}
