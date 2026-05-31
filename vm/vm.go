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
//
// # File Organization
//
// The vm package is split across files by responsibility:
//
//   - vm.go          — VM struct, constructor, Execute entry points
//   - vm_execute.go  — Execute entry points and entry-frame setup
//   - vm_entry_args.go — External entry argument validation and variadic packing
//   - vm_lifecycle.go  — Reset, frame growth, shared globals, and global access
//   - pool.go        — VMPool, ResolveCompiledMethod
//   - interfaces.go  — VM interface and constructors
//   - run.go         — Main fetch-decode-execute loop with hot-path inlined instructions
//   - frame.go       — Frame (call stack entry) and DeferInfo (deferred call metadata)
//   - stack.go       — Operand stack push/pop/peek with bounded growth
//   - call.go        — External function calls (DirectCall + reflect), method dispatch
//   - closure.go     — Closure type and ClosureExecutor for reflect.MakeFunc integration
//   - goroutine.go   — GoroutineTracker, child/defer VM construction
//   - iterator.go    — Range iteration over slices, arrays, maps, and strings
//   - constants.go   — Named constants replacing magic numbers across the VM
//   - typeconv.go    — go/types.Type → reflect.Type conversion with cycle detection
//   - ops_dispatch.go — Opcode routing: executeOp dispatches to category handlers
//   - ops_memory.go   — Stack ops, constants, locals/globals/free vars, fields, addresses
//   - ops_arithmetic.go — Arithmetic, bitwise, comparison, and logical operations
//   - ops_container.go  — Slice/map/chan creation, index, append, copy, delete, range
//   - ops_control.go    — Control flow, channels, select, defer, panic/recover, halt
//   - ops_convert.go    — Type assertion, conversion, and change-type operations
//   - ops_call.go       — Function/closure calls, goroutine spawning, tuple pack/unpack
//
// For detailed internals, see docs/gig-internals.md and docs/value-system.md.
package vm

import (
	"context"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// panicState stores a saved panic state for nested panics.
type panicState struct {
	panicking bool
	panicVal  value.Value
}

const (
	// initialFrameDepth is the starting size of the call frame stack.
	// Covers the vast majority of programs without any growth.
	initialFrameDepth = 64

	// maxFrameDepth is the hard ceiling for the call frame stack.
	// Prevents runaway recursion from consuming unbounded memory.
	// Each slot is a single pointer (8 bytes), so 1024 slots = 8 KB.
	maxFrameDepth = 1024
)

// vm is the bytecode virtual machine struct.
// It executes compiled programs using a stack-based architecture.
type vm struct {
	// program is the compiled program to execute.
	program *bytecode.CompiledProgram

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

	// shared is a pointer to the SharedGlobals for stateful/goroutine mode.
	// If set, global variable operations use the shared (locked) globals
	// instead of the local slice. This enables concurrent execution.
	shared *SharedGlobals

	// ctx is the execution context for cancellation/timeout.
	ctx context.Context

	// panicking indicates a panic is in progress.
	panicking bool

	// panicVal is the current panic value.
	panicVal value.Value

	// panicStack stores saved panic states for nested panics (panic inside defer).
	// When a panic occurs during defer execution, the current panic is pushed
	// onto the stack, and the new panic becomes active. When the inner panic
	// is recovered, the outer panic is restored from the stack.
	panicStack []panicState

	// lastPanicVal preserves the original panic value after run() returns an error.
	// When run() reaches the top frame with an active panic, it formats the panic
	// into an error string (losing type info) and clears panicVal. lastPanicVal
	// keeps the original typed value so callers (e.g. OpRunDefers) can recover it.
	lastPanicVal value.Value

	// deferDepth tracks the nesting level of deferred function execution.
	// 0 = normal execution, 1+ = inside deferred function(s).
	// Replaces the boolean runningDefer flag to support nested defer panics.
	deferDepth int

	// initialGlobals is the post-init globals snapshot.
	// Used by Reset() to restore globals to their initial state.
	initialGlobals []value.Value

	// goroutines tracks active interpreter goroutines for this program.
	goroutines *GoroutineTracker

	// fpool recycles Frame objects to reduce heap allocations.
	fpool framePool
}

// newVM creates a new VM for executing the given program.
func newVM(program *bytecode.CompiledProgram, initialGlobals []value.Value, goroutines *GoroutineTracker) *vm {
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

	// Initialize zero-valued and deferred-type globals.
	for idx, zeroRV := range program.GlobalZeroValues {
		if idx < len(globals) {
			g := globals[idx]
			if !g.IsValid() || g.IsNil() {
				globals[idx] = value.MakeFromReflect(zeroRV)
			}
		}
	}
	for idx, typeIdx := range program.GlobalTypes {
		if idx >= len(globals) || typeIdx >= len(program.Types) {
			continue
		}
		g := globals[idx]
		if !g.IsValid() || g.IsNil() {
			t := program.Types[typeIdx]
			if rt := typeToReflect(t, program); rt != nil && rt.Kind() == reflect.Ptr {
				globals[idx] = value.MakeFromReflect(reflect.New(rt.Elem()))
			}
		}
	}

	// For globals with element types but no pre-computed zero values (anonymous
	// structs, arrays, user-defined named types), compute the zero value using
	// typeToReflect. This must run AFTER GlobalZeroValues so it doesn't override
	// pre-computed values for external types like sync.Mutex.
	for idx, elemType := range program.GlobalElemTypes {
		if idx < len(globals) {
			g := globals[idx]
			if !g.IsValid() || g.IsNil() {
				if _, hasZero := program.GlobalZeroValues[idx]; !hasZero {
					if rt := typeToReflect(elemType, program); rt != nil {
						globals[idx] = value.MakeFromReflect(reflect.New(rt))
					}
				}
			}
		}
	}

	return &vm{
		program:        program,
		stack:          make([]value.Value, initialStackSize),
		sp:             0,
		frames:         make([]*Frame, initialFrameDepth),
		fp:             0,
		globals:        globals,
		initialGlobals: initialGlobals,
		goroutines:     goroutines,
	}
}
