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
	"go/types"
	"reflect"
	"strings"
	"sync"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
)

// vm is the bytecode virtual machine struct.
// It executes compiled programs using a stack-based architecture.
type vm struct {
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

	// initialGlobals is the post-init globals snapshot.
	// Used by Reset() to restore globals to their initial state.
	initialGlobals []value.Value

	// goroutines tracks active interpreter goroutines for this program.
	goroutines *GoroutineTracker

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

	directCall func(args []value.Value) value.Value

	// isVariadic indicates if the function takes variadic arguments.
	isVariadic bool

	// numIn is the number of declared parameters.
	numIn int
}

// newVM creates a new VM for executing the given program.
func newVM(program *bytecode.Program, initialGlobals []value.Value, goroutines *GoroutineTracker) *vm {
	globals := make([]value.Value, len(program.Globals))
	if len(initialGlobals) == len(globals) {
		copy(globals, initialGlobals)
	}
	// Initialize external variable values
	for idx, ptr := range program.ExternalVarValues {
		if idx < len(globals) {
			globals[idx] = value.FromInterface(ptr)
		}
	}

	return &vm{
		program:        program,
		stack:          make([]value.Value, 1024), // initial stack size
		sp:             0,
		frames:         make([]*Frame, 64), // max call depth
		fp:             0,
		globals:        globals,
		initialGlobals: initialGlobals,
		goroutines:     goroutines,
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(program.Constants)),
		},
	}
}

// Reset prepares the VM for reuse by clearing execution state.
func (v *vm) Reset() {
	v.sp = 0
	v.fp = 0
	v.panicking = false
	v.panicVal = value.MakeNil()
	v.ctx = nil
	// Clear all frames (prevents stale frame references from previous execution).
	for i := range v.frames {
		v.frames[i] = nil
	}
	// If globalsPtr is set (shared globals from stateful mode or goroutine),
	// do not restore the local globals copy — the caller manages the shared state.
	if v.globalsPtr != nil {
		v.globalsPtr = nil
		return
	}
	// Stateless mode: restore globals to post-init snapshot, or zero them.
	if len(v.initialGlobals) == len(v.globals) {
		copy(v.globals, v.initialGlobals)
	} else {
		for i := range v.globals {
			v.globals[i] = value.Value{}
		}
	}
	// Restore external variable values (they should always be the same)
	for idx, ptr := range v.program.ExternalVarValues {
		if idx < len(v.globals) {
			v.globals[idx] = value.FromInterface(ptr)
		}
	}
}

// BindSharedGlobals makes this VM execute against the provided shared globals slice.
func (v *vm) BindSharedGlobals(globals *[]value.Value) {
	v.globalsPtr = globals
}

// UnbindSharedGlobals detaches the VM from shared globals.
func (v *vm) UnbindSharedGlobals() {
	v.globalsPtr = nil
}

// Globals returns the VM's global variable slice.
func (v *vm) Globals() []value.Value {
	return v.globals
}

// resolveCompiledMethod finds a compiled method in the program's function table
// ResolveCompiledMethod searches for a compiled method and executes it with the given receiver.
func ResolveCompiledMethod(program *bytecode.Program, methodName string, receiver value.Value) (value.Value, bool) {
	rv, ok := receiver.ReflectValue()
	if !ok {
		return value.MakeNil(), false
	}
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}

	// Extract the type name from the _gig_id field for matching
	receiverTypeName := ""
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)
			if sf.Name == "_gig_id" {
				if idx := strings.LastIndex(sf.PkgPath, "#"); idx >= 0 {
					qualName := sf.PkgPath[idx+1:]
					if dotIdx := strings.LastIndex(qualName, "."); dotIdx >= 0 {
						receiverTypeName = qualName[dotIdx+1:]
					} else {
						receiverTypeName = qualName
					}
				}
				break
			}
		}
	}

	if receiverTypeName == "" {
		return value.MakeNil(), false
	}

	// Search the compiled function table for a method with matching name and receiver type
	for _, fn := range program.FuncByIndex {
		if fn == nil || fn.Source == nil {
			continue
		}
		if fn.Source.Name() != methodName {
			continue
		}
		sig := fn.Source.Signature
		recv := sig.Recv()
		if recv == nil {
			continue
		}
		recvType := recv.Type()
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}
		named, ok := recvType.(*types.Named)
		if !ok {
			continue
		}
		if named.Obj().Name() != receiverTypeName {
			continue
		}
		// Found the method! Execute it with a temporary VM.
		tempVM := &vm{
			program: program,
			stack:   make([]value.Value, 256),
			sp:      0,
			frames:  make([]*Frame, 64),
			fp:      0,
			globals: make([]value.Value, len(program.Globals)),
			ctx:     context.Background(),
			extCallCache: &externalCallCache{
				cache: make([]*extCallCacheEntry, len(program.Constants)),
			},
		}
		// Note: tempVM does not have initialGlobals since resolveCompiledMethod
		// is called without a VM context. This is acceptable because method resolution
		// only needs to execute the method, not full program init.
		tempVM.callFunction(fn, []value.Value{receiver}, nil)
		result, err := tempVM.run()
		if err != nil {
			return value.MakeNil(), false
		}
		return result, true
	}
	return value.MakeNil(), false
}

// VMPool is a thread-safe pool of VMs for a given program.
type VMPool struct {
	mu    sync.Mutex
	vms   []*vm // available VMs
	newVM func() *vm
}

// NewVMPool creates a VM pool for the given program.
func NewVMPool(program *bytecode.Program, initialGlobals []value.Value, goroutines *GoroutineTracker) *VMPool {
	return &VMPool{
		newVM: func() *vm {
			return newVM(program, initialGlobals, goroutines)
		},
	}
}

// Get returns an idle VM from the pool.
func (p *VMPool) Get() VM {
	p.mu.Lock()
	if len(p.vms) > 0 {
		v := p.vms[len(p.vms)-1]
		p.vms = p.vms[:len(p.vms)-1]
		p.mu.Unlock()
		return v
	}
	p.mu.Unlock()
	return p.newVM()
}

// Put returns a VM to the pool for reuse.
func (p *VMPool) Put(x VM) {
	x.Reset()
	p.mu.Lock()
	p.vms = append(p.vms, x.(*vm))
	p.mu.Unlock()
}

// Execute runs the specified function with the given arguments.
func (v *vm) Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error) {
	v.ctx = ctx

	fn, ok := v.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	valArgs := make([]value.Value, len(args))
	copy(valArgs, args)

	frame := v.fpool.get(fn, 0, nil)
	for i, arg := range valArgs {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	v.frames[0] = frame
	v.fp = 1

	result, err := v.run()
	return result, err
}

// ExecuteWithValues runs the specified function with pre-converted Value arguments.
func (v *vm) ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error) {
	v.ctx = ctx

	fn, ok := v.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	frame := v.fpool.get(fn, 0, nil)
	for i, arg := range args {
		if i < fn.NumLocals {
			frame.locals[i] = arg
			if frame.intLocals != nil {
				frame.intLocals[i] = arg.RawInt()
			}
		}
	}
	v.frames[0] = frame
	v.fp = 1

	return v.run()
}

// getGlobals returns the globals slice, using the shared pointer if available.
func (v *vm) getGlobals() []value.Value {
	if v.globalsPtr != nil {
		return *v.globalsPtr
	}
	return v.globals
}
