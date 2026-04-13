// Package compiler provides SSA-to-bytecode compilation for the Gig interpreter.
//
// The compiler translates Go SSA (Static Single Assignment) intermediate representation
// into a custom bytecode format defined in the bytecode package.
package compiler

import (
	"fmt"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// Compiler compiles SSA programs into bytecode.
type Compiler interface {
	Compile(mainPkg *ssa.Package) (*bytecode.CompiledProgram, error)
}

// NewCompiler creates a new compiler with the given package lookup for resolving external functions.
// The PackageLookup dependency is injected to decouple the compiler from the importer package.
func NewCompiler(lookup PackageLookup) Compiler {
	return &compiler{
		lookup:            lookup,
		constants:         make([]any, 0),
		types:             make([]types.Type, 0),
		globals:           make(map[string]int),
		globalZeroValues:  make(map[int]reflect.Value),
		externalVarValues: make(map[int]any),
		funcs:             make(map[string]*bytecode.CompiledFunction),
		funcIndex:         make(map[*ssa.Function]int),
	}
}

// compiler is the concrete implementation of the Compiler interface.
// It maintains state during compilation including the current function,
// symbol table, and jump targets that need patching.
type compiler struct {
	// lookup resolves external package functions (injected dependency).
	lookup PackageLookup

	// program is the output program being compiled.
	program *bytecode.CompiledProgram

	// constants is the constant pool being built.
	constants []any

	// types is the type pool being built.
	types []types.Type

	// globals maps global names to indices.
	globals map[string]int

	// globalZeroValues maps global index to its zero reflect.Value.
	// Resolved at compile time for external named struct types (e.g., sync.Mutex).
	globalZeroValues map[int]reflect.Value

	// externalVarValues stores external variable values indexed by global index.
	// These are resolved at compile time and used to initialize globals in the VM.
	externalVarValues map[int]any

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

// Compile is the main entry point that compiles an SSA package to a bytecode Program.
func (c *compiler) Compile(mainPkg *ssa.Package) (*bytecode.CompiledProgram, error) {
	c.program = &bytecode.CompiledProgram{
		Functions: make(map[string]*bytecode.CompiledFunction),
		Globals:   make(map[string]int),
		Types:     make([]types.Type, 0),
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

	// Collect synthetic wrapper functions ($bound, $thunk) referenced by compiled code.
	// SSA generates these for method values (e.g., obj.Method → $bound) and method
	// expressions (e.g., (*Type).Method → $thunk). They have Pkg==nil but reference
	// methods from the main package. We scan all collected functions' instructions
	// to discover them.
	changed := true
	for changed {
		changed = false
		for _, fn := range allFuncs {
			if fn.Blocks == nil {
				continue
			}
			for _, block := range fn.Blocks {
				for _, instr := range block.Instrs {
					// MakeClosure references the wrapper function directly
					if mc, ok := instr.(*ssa.MakeClosure); ok {
						if wrapperFn, ok := mc.Fn.(*ssa.Function); ok && !seen[wrapperFn] {
							if wrapperFn.Blocks != nil {
								collectFuncs(wrapperFn)
								changed = true
							}
						}
					}
					// Direct calls to $thunk functions
					if call, ok := instr.(*ssa.Call); ok {
						if calledFn, ok := call.Call.Value.(*ssa.Function); ok && !seen[calledFn] {
							if calledFn.Blocks != nil {
								collectFuncs(calledFn)
								changed = true
							}
						}
					}
				}
			}
		}
	}

	// First pass: assign indices to all functions
	for idx, fn := range allFuncs {
		c.funcIndex[fn] = idx
	}

	// Second pass: compile each function and build direct-index lookup table
	c.program.FuncByIndex = make([]*bytecode.CompiledFunction, len(allFuncs))
	for _, fn := range allFuncs {
		compiled, err := c.compileFunction(fn)
		if err != nil {
			return nil, fmt.Errorf("compile function %s: %w", fn.Name(), err)
		}
		c.funcs[fn.Name()] = compiled
		c.program.Functions[fn.Name()] = compiled
		// Use the SSA-pointer-based index to avoid name collisions
		// (e.g., two methods named "Get" on different types).
		idx := c.funcIndex[fn]
		c.program.FuncByIndex[idx] = compiled
	}

	c.program.Constants = c.constants
	c.program.Types = c.types
	c.program.Globals = c.globals
	c.program.GlobalZeroValues = c.globalZeroValues
	c.program.ExternalVarValues = c.externalVarValues
	c.program.TypeResolver = c.lookup

	// Build method lookup index for O(k) dispatch instead of O(n) linear scan.
	methodsByName := make(map[string][]*bytecode.CompiledFunction)
	for _, fn := range c.program.FuncByIndex {
		if fn != nil && fn.HasReceiver {
			methodsByName[fn.Name] = append(methodsByName[fn.Name], fn)
		}
	}
	c.program.MethodsByName = methodsByName

	// Pre-bake constants for O(1) OpConst (avoids FromInterface per instruction)
	c.program.PrebakedConstants = make([]value.Value, len(c.constants))
	for i, k := range c.constants {
		c.program.PrebakedConstants[i] = value.FromInterface(k)
	}

	// Build int-specialized constant pool for OpInt* superinstructions
	c.program.IntConstants = make([]int64, len(c.constants))
	for i, k := range c.constants {
		switch v := k.(type) {
		case int:
			c.program.IntConstants[i] = int64(v)
		case int8:
			c.program.IntConstants[i] = int64(v)
		case int16:
			c.program.IntConstants[i] = int64(v)
		case int32:
			c.program.IntConstants[i] = int64(v)
		case int64:
			c.program.IntConstants[i] = v
		}
	}

	return c.program, nil
}
