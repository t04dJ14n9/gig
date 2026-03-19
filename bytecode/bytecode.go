package bytecode

import (
	"go/types"
	"reflect"
	"sync"

	"golang.org/x/tools/go/ssa"

	"git.woa.com/youngjin/gig/value"
)

// PackageLookup resolves external package functions for the compiler.
// This interface enables dependency injection: the compiler depends on this
// abstraction rather than importing the importer package directly.
type PackageLookup interface {
	LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool)
	LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool)
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

	// HasIntLocals indicates that this function uses OpInt* superinstructions
	// and needs intLocals []int64 allocated in its Frame.
	HasIntLocals bool

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

	// ReceiverTypeName is the fully qualified name of the receiver type
	// (e.g., "GetterImpl", "AdderStruct"). Used by callCompiledMethod
	// to disambiguate when multiple compiled methods share the same name.
	// Empty string means "match any receiver" (backward compatible).
	ReceiverTypeName string

	// DirectCall is an optional typed wrapper that avoids reflect.Call for this method.
	// args[0] is the receiver, args[1:] are method arguments.
	// If nil, the VM will use reflection for the call.
	DirectCall func([]value.Value) value.Value
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

	// IntConstants is an int-specialized constant pool.
	// For each constant that is an int64, IntConstants[i] holds the value.
	// Used by OpInt* superinstructions for zero-overhead constant access.
	IntConstants []int64

	// Globals maps global variable names to their indices.
	Globals map[string]int

	// MainPkg is the SSA package (for debugging/inspection).
	MainPkg *ssa.Package

	// Types is the type pool for runtime type operations.
	Types []types.Type

	// FuncIndex maps SSA functions to their indices for call instructions.
	FuncIndex map[*ssa.Function]int

	// InitialGlobals holds the global variable state after init() has run.
	// New VMs copy this slice as their starting globals so each call sees a
	// fully-initialised package state.  Nil when there is no init() function.
	InitialGlobals []value.Value

	// ReflectTypeCache caches types.Type → reflect.Type conversions at the
	// program level. This prevents reflect.StructOf from returning different
	// reflect.Type objects for the same types.Type across multiple VM executions,
	// which would cause "reflect.Set: value not assignable" panics.
	// Key: types.Type, Value: reflect.Type.
	ReflectTypeCache sync.Map
}

// CachedReflectType looks up a cached reflect.Type for the given types.Type.
// Returns the cached type and true, or nil and false if not cached.
func (p *Program) CachedReflectType(t types.Type) (reflect.Type, bool) {
	if v, ok := p.ReflectTypeCache.Load(t); ok {
		return v.(reflect.Type), true
	}
	return nil, false
}

// CacheReflectType stores a types.Type → reflect.Type mapping.
// Uses LoadOrStore to handle concurrent writes safely.
func (p *Program) CacheReflectType(t types.Type, rt reflect.Type) reflect.Type {
	actual, _ := p.ReflectTypeCache.LoadOrStore(t, rt)
	return actual.(reflect.Type)
}
