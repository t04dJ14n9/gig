package vm

import (
	"context"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// VM is the interface for the bytecode virtual machine.
// It executes compiled programs using a stack-based architecture.
type VM interface {
	Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error)
	ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error)
	Globals() []value.Value
	Reset()
	// BindSharedGlobals makes this VM execute against the provided shared globals
	// slice. All global loads/stores will go through globalsPtr, which points at
	// the Program-owned backing store. This must be called before Execute and
	// paired with UnbindSharedGlobals after execution finishes.
	BindSharedGlobals(globals *[]value.Value)

	// UnbindSharedGlobals detaches the VM from shared globals so that Reset (called
	// when the VM is returned to the pool) does not clobber the shared state.
	UnbindSharedGlobals()
}

// New creates a new VM for executing the given program.
func New(program *bytecode.Program) VM {
	return newVM(program, nil, NewGoroutineTracker())
}

// NewWithOptions creates a VM with custom options.
func NewWithOptions(program *bytecode.Program, opts ...VMOption) VM {
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
