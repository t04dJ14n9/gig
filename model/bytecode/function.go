package bytecode

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

	// HasIntLocals indicates that this function uses OpInt* superinstructions
	// and needs intLocals []int64 allocated in its Frame.
	HasIntLocals bool

	// FuncIdx is the index of this function in FuncByIndex.
	// Used by method dispatch to call compiled functions by index.
	FuncIdx int

	// ReceiverTypeName is the unqualified receiver type name (e.g., "Reader").
	// Empty for non-method functions. Used by method dispatch.
	ReceiverTypeName string

	// HasReceiver indicates this function is a method (has a receiver parameter).
	HasReceiver bool

	// ResultAllocSlots holds the local slot indices of Alloc instructions that
	// correspond to named return values or variables captured by defer closures
	// and used in the return path. During panic recovery, the VM dereferences
	// these slots to reconstruct the correct return value instead of pushing nil.
	// Populated at compile time by examining SSA Return instructions.
	ResultAllocSlots []int
}
