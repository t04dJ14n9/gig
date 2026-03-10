// Package compiler provides SSA-to-bytecode compilation for the Gig interpreter.
//
// The compiler translates Go SSA (Static Single Assignment) intermediate representation
// into a custom bytecode format defined in the bytecode package.
package compiler

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
)

// Compiler compiles SSA programs into bytecode.
type Compiler interface {
	Compile(mainPkg *ssa.Package) (*bytecode.Program, error)
}

// NewCompiler creates a new compiler with the given package lookup for resolving external functions.
// The PackageLookup dependency is injected to decouple the compiler from the importer package.
func NewCompiler(lookup bytecode.PackageLookup) Compiler {
	return &compiler{
		lookup:    lookup,
		constants: make([]any, 0),
		types:     make([]types.Type, 0),
		globals:   make(map[string]int),
		funcs:     make(map[string]*bytecode.CompiledFunction),
		funcIndex: make(map[*ssa.Function]int),
	}
}

// compiler is the concrete implementation of the Compiler interface.
// It maintains state during compilation including the current function,
// symbol table, and jump targets that need patching.
type compiler struct {
	// lookup resolves external package functions (injected dependency).
	lookup bytecode.PackageLookup

	// program is the output program being compiled.
	program *bytecode.Program

	// constants is the constant pool being built.
	constants []any

	// types is the type pool being built.
	types []types.Type

	// globals maps global names to indices.
	globals map[string]int

	// funcs maps function names to compiled functions.
	funcs map[string]*bytecode.CompiledFunction

	// funcIndex maps SSA functions to call indices.
	funcIndex map[*ssa.Function]int

	// currentFunc is the function being compiled.
	currentFunc *bytecode.CompiledFunction

	// symbolTable tracks SSA values to local slots.
	symbolTable *SymbolTable

	// jumps tracks jump instructions needing target patching.
	jumps []jumpInfo

	// phiSlots maps Phi nodes to their allocated local slots.
	phiSlots map[*ssa.Phi]int
}

// jumpInfo tracks a jump instruction that needs its target patched.
type jumpInfo struct {
	offset      int
	targetBlock *ssa.BasicBlock
}

// phiMove represents a move instruction for Phi elimination.
type phiMove struct {
	sourceValue ssa.Value
	targetSlot  int
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

// Compile is the main entry point that compiles an SSA package to a bytecode Program.
func (c *compiler) Compile(mainPkg *ssa.Package) (*bytecode.Program, error) {
	c.program = &bytecode.Program{
		Functions: make(map[string]*bytecode.CompiledFunction),
		Globals:   make(map[string]int),
		MainPkg:   mainPkg,
		Types:     make([]types.Type, 0),
		FuncIndex: make(map[*ssa.Function]int),
	}

	// Collect all functions (including anonymous/nested and methods)
	var allFuncs []*ssa.Function
	seen := make(map[*ssa.Function]bool)
	var collectFuncs func(fn *ssa.Function)
	collectFuncs = func(fn *ssa.Function) {
		if seen[fn] {
			return
		}
		seen[fn] = true
		allFuncs = append(allFuncs, fn)
		for _, anon := range fn.AnonFuncs {
			collectFuncs(anon)
		}
	}
	for _, member := range mainPkg.Members {
		if fn, ok := member.(*ssa.Function); ok {
			collectFuncs(fn)
		}
	}
	// Also collect methods defined on types in the package.
	// Methods are not package members in SSA — they hang off the type's method set.
	for _, member := range mainPkg.Members {
		t, ok := member.(*ssa.Type)
		if !ok {
			continue
		}
		// Collect methods on both value and pointer receiver types.
		for _, recv := range []types.Type{t.Type(), types.NewPointer(t.Type())} {
			mset := mainPkg.Prog.MethodSets.MethodSet(recv)
			for i := 0; i < mset.Len(); i++ {
				if fn := mainPkg.Prog.MethodValue(mset.At(i)); fn != nil && fn.Package() == mainPkg {
					collectFuncs(fn)
				}
			}
		}
	}

	// First pass: assign indices to all functions
	for idx, fn := range allFuncs {
		c.funcIndex[fn] = idx
		c.program.FuncIndex[fn] = idx
	}

	// Second pass: compile each function
	for _, fn := range allFuncs {
		compiled, err := c.compileFunction(fn)
		if err != nil {
			return nil, fmt.Errorf("compile function %s: %w", fn.Name(), err)
		}
		c.funcs[fn.Name()] = compiled
		c.program.Functions[fn.Name()] = compiled
	}

	// Build direct-index lookup table for O(1) function calls
	c.program.FuncByIndex = make([]*bytecode.CompiledFunction, len(allFuncs))
	for _, fn := range allFuncs {
		idx := c.funcIndex[fn]
		c.program.FuncByIndex[idx] = c.funcs[fn.Name()]
	}

	c.program.Constants = c.constants
	c.program.Types = c.types
	c.program.Globals = c.globals

	// Pre-bake constants for O(1) OpConst (avoids FromInterface per instruction)
	c.program.PrebakedConstants = make([]value.Value, len(c.constants))
	for i, k := range c.constants {
		c.program.PrebakedConstants[i] = value.FromInterface(k)
	}

	// Build int-specialized constant pool for OpInt* superinstructions
	c.program.IntConstants = make([]int64, len(c.constants))
	for i, k := range c.constants {
		if v, ok := k.(int64); ok {
			c.program.IntConstants[i] = v
		}
	}

	return c.program, nil
}

// Compile is a convenience package-level function that compiles an SSA package.
// It creates a compiler with the given PackageLookup and invokes compilation.
func Compile(lookup bytecode.PackageLookup, mainPkg *ssa.Package) (*bytecode.Program, error) {
	return NewCompiler(lookup).Compile(mainPkg)
}
