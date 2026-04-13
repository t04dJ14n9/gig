package vm

import (
	"context"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// VM is the interface for the bytecode virtual machine.
// It executes compiled programs using a stack-based architecture.
type VM interface {
	Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error)
	ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error)
	Globals() []value.Value
	Reset()
	// BindSharedGlobals makes this VM execute against the provided SharedGlobals.
	// All global loads/stores will go through the shared (locked) globals,
	// enabling concurrent execution in stateful mode. This must be called before
	// Execute and paired with UnbindSharedGlobals after execution finishes.
	BindSharedGlobals(sg *SharedGlobals)

	// UnbindSharedGlobals detaches the VM from shared globals so that Reset (called
	// when the VM is returned to the pool) does not clobber the shared state.
	UnbindSharedGlobals()
}

// New creates a new VM for executing the given program.
func New(program *bytecode.CompiledProgram) VM {
	return newVM(program, nil, NewGoroutineTracker())
}

// NewWithOptions creates a VM with custom options.
func NewWithOptions(program *bytecode.CompiledProgram, opts ...VMOption) VM {
	v := newVM(program, nil, NewGoroutineTracker())
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// VMOption configures a VM.
type VMOption func(*vm)

// WithContext sets the execution context for a VM.
func WithContext(ctx context.Context) VMOption {
	return func(v *vm) {
		v.ctx = ctx
	}
}
