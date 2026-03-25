package compiler

import "golang.org/x/tools/go/ssa"

// jumpInfo tracks a jump instruction that needs its target patched.
type jumpInfo struct {
	offset      int
	targetBlock *ssa.BasicBlock
}

// SymbolTable tracks SSA values to local slots.
type SymbolTable struct {
	locals    map[ssa.Value]int
	freeVars  map[ssa.Value]int
	numLocals int
}

// NewSymbolTable creates a new symbol table for tracking SSA values.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		locals:   make(map[ssa.Value]int),
		freeVars: make(map[ssa.Value]int),
	}
}

// AllocLocal allocates a new local slot for an SSA value.
func (s *SymbolTable) AllocLocal(v ssa.Value) int {
	if idx, ok := s.locals[v]; ok {
		return idx
	}
	idx := s.numLocals
	s.locals[v] = idx
	s.numLocals++
	return idx
}

// GetLocal returns the local slot index for an SSA value.
func (s *SymbolTable) GetLocal(v ssa.Value) (int, bool) {
	idx, ok := s.locals[v]
	return idx, ok
}

// NumLocals returns the number of allocated local slots.
func (s *SymbolTable) NumLocals() int {
	return s.numLocals
}
