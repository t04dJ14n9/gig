// Package vm provides a stack-based bytecode virtual machine for executing compiled Gig programs.
//
// The VM executes bytecode instructions produced by the compiler. It uses a stack-based
// architecture for operand handling and a frame-based call stack for function calls.
//
// # Architecture
//
// The VM maintains:
//   - An operand stack for intermediate values
//   - A call frame stack for function calls
//   - A global variable array for package-level variables
//   - An inline cache for external function calls
//
// # Execution Model
//
// The VM fetches, decodes, and executes bytecode instructions in a loop.
// Each instruction may push/pop values from the operand stack and modify the call stack.
// Execution continues until all frames return or an error occurs.
//
// # Context Support
//
// The VM supports context-based cancellation and timeout. It checks the context
// every 1024 instructions to avoid blocking on long-running operations.
//
// # Closures
//
// Closures are represented as Closure structs containing a function reference
// and captured free variables. Free variables are stored as pointers to allow
// shared state between closures.
package vm

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

func init() { //nolint:gochecknoinits // registers cross-package callback, must run before any VM usage
	// Register the closure caller so that value.ToReflectValue can wrap
	// *vm.Closure objects into real Go functions via reflect.MakeFunc.
	value.RegisterClosureCaller(func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
		c, ok := closure.(*Closure)
		if !ok {
			return nil
		}
		prog := c.Program
		if prog == nil {
			return nil
		}
		// Create a temporary VM to execute the closure
		closureVM := &VM{
			program: prog,
			stack:   make([]value.Value, 256),
			sp:      0,
			frames:  make([]*Frame, 64),
			fp:      0,
			globals: make([]value.Value, len(prog.Globals)),
			ctx:     context.Background(),
			extCallCache: &externalCallCache{
				cache: make([]*extCallCacheEntry, len(prog.Constants)),
			},
		}
		if len(prog.InitialGlobals) == len(closureVM.globals) {
			copy(closureVM.globals, prog.InitialGlobals)
		}
		// Convert reflect.Value args to value.Value args
		valArgs := make([]value.Value, len(args))
		for i, arg := range args {
			valArgs[i] = value.MakeFromReflect(arg)
		}
		// Call the closure function with its captured free variables
		closureVM.callFunction(c.Fn, valArgs, c.FreeVars)
		result, _ := closureVM.run()
		// Return the result as an interface{} wrapped in reflect.Value.
		// reflect.MakeFunc will handle the type conversion from the returned
		// reflect.Value to the expected function signature return type.
		if result.Kind() == value.KindNil {
			return []reflect.Value{}
		}
		// If we have expected output types and the result is a closure (KindFunc),
		// use ToReflectValue to properly wrap it as a real Go function.
		// This handles nested closures like func() → func() int.
		if len(outTypes) > 0 {
			return []reflect.Value{result.ToReflectValue(outTypes[0])}
		}
		// Fallback: use reflect.ValueOf + Convert to match the expected return type.
		// This handles int64 → int conversion etc.
		iface := result.Interface()
		if iface == nil {
			return []reflect.Value{}
		}
		return []reflect.Value{reflect.ValueOf(iface)}
	})
}

// Executor executes compiled bytecode programs.
type Executor interface {
	Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error)
	ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error)
}

// VM is the bytecode virtual machine.
// It executes compiled programs using a stack-based architecture.
type VM struct {
	// program is the compiled program to execute.
	program *bytecode.Program

	// stack is the operand stack for intermediate values.
	stack []value.Value

	// sp is the stack pointer (index of next free slot).
	sp int

	// frames is the call frame stack.
	frames []*Frame

	// fp is the frame pointer (number of active frames).
	fp int

	// globals stores global variables.
	globals []value.Value

	// globalsPtr is a pointer to shared globals (for goroutine communication).
	// If set, globals operations use this pointer instead of the local slice.
	globalsPtr *[]value.Value

	// ctx is the execution context for cancellation/timeout.
	ctx context.Context

	// panicking indicates a panic is in progress.
	panicking bool

	// panicVal is the current panic value.
	panicVal value.Value

	// extCallCache caches resolved external function info for fast dispatch.
	// Uses a shared cache pointer for concurrent access from goroutines.
	extCallCache *externalCallCache

	// fpool recycles Frame objects to reduce heap allocations.
	fpool framePool
}

// externalCallCache is a shared cache for external function lookups.
// It is shared between a parent VM and all its child goroutine VMs.
type externalCallCache struct {
	mu    sync.RWMutex
	cache []*extCallCacheEntry
}

// extCallCacheEntry caches resolved external function info for fast dispatch.
// This avoids repeated reflection lookups for external function calls.
type extCallCacheEntry struct {
	// fn is the reflect.Value of the function.
	fn reflect.Value

	// fnType is the function's type.
	fnType reflect.Type

	// directCall is a typed wrapper that bypasses reflect.Call.
	directCall func([]value.Value) value.Value

	// isVariadic indicates if the function takes variadic arguments.
	isVariadic bool

	// numIn is the number of declared parameters.
	numIn int
}

// New creates a new VM for executing the given program.
// The VM is created with an empty stack and call frame array.
func New(program *bytecode.Program) *VM {
	globals := make([]value.Value, len(program.Globals))
	if len(program.InitialGlobals) == len(globals) {
		copy(globals, program.InitialGlobals)
	}
	// Initialize external variable values
	for idx, ptr := range program.ExternalVarValues {
		if idx < len(globals) {
			globals[idx] = value.FromInterface(ptr)
		}
	}
	return &VM{
		program: program,
		stack:   make([]value.Value, 1024), // initial stack size
		sp:      0,
		frames:  make([]*Frame, 64), // max call depth
		fp:      0,
		globals: globals,
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(program.Constants)),
		},
	}
}

// Reset prepares the VM for reuse by clearing execution state.
// The stack, frames, and globals slices are retained (zero-alloc reuse).
// If the VM is currently bound to shared globals (stateful mode), the shared
// globals are left untouched and only execution state is cleared.
func (vm *VM) Reset() {
	vm.sp = 0
	vm.fp = 0
	vm.panicking = false
	vm.panicVal = value.MakeNil()
	vm.ctx = nil
	// If globalsPtr is set (shared globals from stateful mode or goroutine),
	// do not restore the local globals copy — the caller manages the shared
	// state.  Only clear globalsPtr so the VM is detached.
	if vm.globalsPtr != nil {
		// VM was bound to shared globals; detach but don't reset local copy.
		vm.globalsPtr = nil
		return
	}
	// Stateless mode: restore globals to post-init snapshot, or zero them.
	if len(vm.program.InitialGlobals) == len(vm.globals) {
		copy(vm.globals, vm.program.InitialGlobals)
	} else {
		for i := range vm.globals {
			vm.globals[i] = value.Value{}
		}
	}
	// Restore external variable values (they should always be the same)
	for idx, ptr := range vm.program.ExternalVarValues {
		if idx < len(vm.globals) {
			vm.globals[idx] = value.FromInterface(ptr)
		}
	}
}

// BindSharedGlobals makes this VM execute against the provided shared globals
// slice.  All global loads/stores will go through globalsPtr, which points at
// the Program-owned backing store.  This must be called before Execute and
// paired with UnbindSharedGlobals after execution finishes.
func (vm *VM) BindSharedGlobals(globals *[]value.Value) {
	vm.globalsPtr = globals
}

// UnbindSharedGlobals detaches the VM from shared globals so that Reset (called
// when the VM is returned to the pool) does not clobber the shared state.
func (vm *VM) UnbindSharedGlobals() {
	vm.globalsPtr = nil
}

// Globals returns the VM's global variable slice.
// Used to snapshot global state after init() has run.
func (vm *VM) Globals() []value.Value {
	return vm.globals
}

// VMPool is a pool of VMs for a given program, eliminating per-call allocation overhead.
type VMPool struct {
	pool sync.Pool
}

// NewVMPool creates a VM pool for the given program.
func NewVMPool(program *bytecode.Program) *VMPool {
	return &VMPool{
		pool: sync.Pool{
			New: func() any {
				return New(program)
			},
		},
	}
}

// Get returns a VM from the pool (or creates a new one).
func (p *VMPool) Get() *VM {
	return p.pool.Get().(*VM)
}

// Put returns a VM to the pool for reuse.
func (p *VMPool) Put(v *VM) {
	v.Reset()
	p.pool.Put(v)
}

// Execute runs the specified function with the given arguments.
// It creates an initial call frame and starts the execution loop.
// Returns the result value or an error if execution fails.
func (vm *VM) Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error) {
	vm.ctx = ctx

	// Look up the function
	fn, ok := vm.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	// Convert args to []value.Value
	valArgs := make([]value.Value, len(args))
	copy(valArgs, args)

	// Create initial frame using pool
	frame := vm.fpool.get(fn, 0, nil)
	for i, arg := range valArgs {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			// Mirror int parameters into intLocals for OpInt* opcodes
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	vm.frames[0] = frame
	vm.fp = 1

	// Run the VM
	result, err := vm.run()
	return result, err
}

// ExecuteWithValues runs the specified function with pre-converted Value arguments.
// This is more efficient than Execute when the arguments are already Value types.
func (vm *VM) ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error) {
	vm.ctx = ctx

	// Look up the function
	fn, ok := vm.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	// Create initial frame using pool
	frame := vm.fpool.get(fn, 0, nil)
	for i, arg := range args {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			// Mirror int parameters into intLocals for OpInt* opcodes
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	vm.frames[0] = frame
	vm.fp = 1

	// Run the VM
	return vm.run()
}

// getGlobals returns the globals slice, using the shared pointer if available.
// This allows goroutines to share globals for communication.
func (vm *VM) getGlobals() []value.Value {
	if vm.globalsPtr != nil {
		return *vm.globalsPtr
	}
	return vm.globals
}
