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
	return &VM{
		program: program,
		stack:   make([]value.Value, 1024), // initial stack size
		sp:      0,
		frames:  make([]*Frame, 64), // max call depth
		fp:      0,
		globals: make([]value.Value, len(program.Globals)),
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(program.Constants)),
		},
	}
}

// Reset prepares the VM for reuse by clearing execution state.
// The stack, frames, and globals slices are retained (zero-alloc reuse).
func (vm *VM) Reset() {
	vm.sp = 0
	vm.fp = 0
	vm.panicking = false
	vm.panicVal = value.MakeNil()
	vm.ctx = nil
	vm.globalsPtr = nil
	// Clear globals to zero values (prevent stale state between runs)
	for i := range vm.globals {
		vm.globals[i] = value.Value{}
	}
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

// checkContext returns ctx.Err() if the context is cancelled, nil otherwise.
// Use this for fast context checks in hot paths.
func (vm *VM) checkContext() error {
	select {
	case <-vm.ctx.Done():
		return vm.ctx.Err()
	default:
		return nil
	}
}
