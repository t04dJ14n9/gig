package bytecode

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"gig/value"
)

// PackageLookup resolves external package functions for the compiler.
// This interface enables dependency injection: the compiler depends on this
// abstraction rather than importing the importer package directly.
type PackageLookup interface {
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)
}

// CompiledFunction represents a function compiled to bytecode.
// It contains the bytecode instructions, local variable count, and metadata.
type CompiledFunction struct {
	// Name is the function name for debugging.
	Name string

	// Instructions is the compiled bytecode.
	Instructions []byte

	// NumLocals is the number of local variable slots.
	// This includes parameters and intermediate values.
	NumLocals int

	// NumParams is the number of function parameters.
	NumParams int

	// NumFreeVars is the number of free variables (for closures).
	NumFreeVars int

	// MaxStack is the maximum stack depth (for future optimization).
	MaxStack int

	// Source is the original SSA function for debugging.
	Source *ssa.Function
}

// ExternalFuncInfo contains pre-resolved external function info for fast calls.
// This allows the VM to bypass reflection when calling external functions.
type ExternalFuncInfo struct {
	// Func is the actual function value.
	Func any

	// DirectCall is a typed wrapper that avoids reflect.Call.
	// If nil, the VM will use reflection for the call.
	DirectCall func([]value.Value) value.Value
}

// ExternalMethodInfo contains method dispatch information.
// It is stored in the constant pool and used by OpMethodCall.
type ExternalMethodInfo struct {
	// MethodName is the name of the method to call.
	MethodName string
}

// Program represents a compiled program ready for execution.
// It contains all compiled functions, constants, types, and global variables.
type Program struct {
	// Functions maps function names to their compiled bytecode.
	Functions map[string]*CompiledFunction

	// FuncByIndex provides O(1) function lookup by index.
	// Populated at compile time so the VM can skip the FuncIndex map scan.
	FuncByIndex []*CompiledFunction

	// Constants is the constant pool for literal values and external references.
	Constants []any

	// PrebakedConstants is the pre-converted constant pool.
	// Built once at startup to avoid per-OpConst value.FromInterface() calls.
	PrebakedConstants []value.Value

	// Globals maps global variable names to their indices.
	Globals map[string]int

	// MainPkg is the SSA package (for debugging/inspection).
	MainPkg *ssa.Package

	// Types is the type pool for runtime type operations.
	Types []types.Type

	// FuncIndex maps SSA functions to their indices for call instructions.
	FuncIndex map[*ssa.Function]int
}
