// Package compiler provides SSA-to-bytecode compilation for the Gig interpreter.
//
// The compiler translates Go SSA (Static Single Assignment) intermediate representation
// into a custom bytecode format defined in the bytecode package.
package compiler

import (
	"go/types"
	"reflect"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// Compiler compiles SSA programs into bytecode.
type Compiler interface {
	Compile(mainPkg *ssa.Package) (*bytecode.CompiledProgram, error)
}

// NewCompiler creates a new compiler with the given package lookup for resolving external functions.
// The PackageLookup dependency is injected to decouple the compiler from the importer package.
func NewCompiler(lookup PackageLookup, allowUnsafeTypePass bool) Compiler {
	return &compiler{
		lookup:              lookup,
		allowUnsafeTypePass: allowUnsafeTypePass,
		constants:           make([]any, 0),
		types:               make([]types.Type, 0),
		globals:             make(map[string]int),
		globalZeroValues:    make(map[int]reflect.Value),
		globalElemTypes:     make(map[int]types.Type),
		externalVarValues:   make(map[int]any),
		funcs:               make(map[string]*bytecode.CompiledFunction),
		funcIndex:           make(map[*ssa.Function]int),
	}
}

// addError records a compilation error for later reporting.
func (c *compiler) addError(err error) {
	c.errors = append(c.errors, err)
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

	// globalElemTypes maps global index to the element type (T from *T).
	// Used by the VM to compute zero values for anonymous structs, arrays, etc.
	globalElemTypes map[int]types.Type

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

	// errors collects non-fatal compilation errors (e.g., type safety violations).
	errors []error

	// allowUnsafeTypePass disables type safety validation for external calls.
	allowUnsafeTypePass bool
}

// Compile is the main entry point that compiles an SSA package to a bytecode Program.
func (c *compiler) Compile(mainPkg *ssa.Package) (*bytecode.CompiledProgram, error) {
	c.program = c.newProgram()

	functions := collectPackageFunctions(mainPkg)
	c.indexFunctions(functions)
	if err := c.compilePackageFunctions(functions); err != nil {
		return nil, err
	}
	c.finalizeProgram()
	if err := c.compilationError(); err != nil {
		return nil, err
	}
	c.program.ResolveExternCalls()

	return c.program, nil
}
