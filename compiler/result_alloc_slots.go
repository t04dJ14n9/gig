package compiler

import "golang.org/x/tools/go/ssa"

// detectResultAllocSlots finds Alloc instructions that correspond to named
// return values. In Go, only named return variables can be modified by deferred
// functions during panic recovery. Unnamed returns use the zero value after recovery.
//
// These are identified as Allocs whose value (or a deref thereof) appears
// directly in a Return instruction's results.
func detectResultAllocSlots(fn *ssa.Function, st *SymbolTable) []int {
	// Panic recovery needs named return variables after deferred code has had a
	// chance to mutate them. We record the backing Alloc local slots at compile
	// time so the VM can reconstruct those values without re-walking SSA.
	if fn.Blocks == nil {
		return nil
	}

	allocSet := collectFunctionAllocs(fn)
	if len(allocSet) == 0 {
		return nil
	}

	slotSet := collectReturnAllocSlots(fn, st, allocSet)
	if len(slotSet) == 0 {
		return nil
	}
	return resultSlotsFromSet(slotSet)
}

func collectFunctionAllocs(fn *ssa.Function) map[ssa.Value]bool {
	allocSet := make(map[ssa.Value]bool)
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			if _, ok := instr.(*ssa.Alloc); ok {
				allocSet[instr.(ssa.Value)] = true
			}
		}
	}
	return allocSet
}

func collectReturnAllocSlots(fn *ssa.Function, st *SymbolTable, allocSet map[ssa.Value]bool) map[int]bool {
	// SSA represents named returns as either direct Alloc references or
	// UnOp(deref) of an Alloc. Both shapes must map back to the same local slot.
	slotSet := make(map[int]bool)
	for _, block := range fn.Blocks {
		for _, instr := range block.Instrs {
			ret, ok := instr.(*ssa.Return)
			if !ok {
				continue
			}
			for _, result := range ret.Results {
				recordReturnAllocSlot(st, slotSet, allocSet, result)
			}
		}
	}
	return slotSet
}

func recordReturnAllocSlot(st *SymbolTable, slotSet map[int]bool, allocSet map[ssa.Value]bool, result ssa.Value) {
	if allocSet[result] {
		recordAllocSlot(st, slotSet, result)
	}
	if unop, ok := result.(*ssa.UnOp); ok && allocSet[unop.X] {
		recordAllocSlot(st, slotSet, unop.X)
	}
}

func recordAllocSlot(st *SymbolTable, slotSet map[int]bool, alloc ssa.Value) {
	if idx, ok := st.GetLocal(alloc); ok {
		slotSet[idx] = true
	}
}

func resultSlotsFromSet(slotSet map[int]bool) []int {
	slots := make([]int, 0, len(slotSet))
	for idx := range slotSet {
		slots = append(slots, idx)
	}
	return slots
}
