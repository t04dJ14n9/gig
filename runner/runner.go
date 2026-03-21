// Package runner provides program execution orchestration for the Gig interpreter.
//
// It encapsulates VM pool management, global state handling (stateless snapshot
// vs stateful shared globals), and result unpacking, providing a clean execution
// API on top of the low-level VM.
package runner

import (
	"context"
	"sync"
	"time"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
	"git.woa.com/youngjin/gig/vm"
)

// Executor is the main execution interface for running compiled programs.
type Executor interface {
	Run(funcName string, params ...any) (any, error)
	RunWithContext(ctx context.Context, funcName string, params ...any) (any, error)
	RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error)
	InternalProgram() *bytecode.Program
}

// Runner orchestrates program execution with VM pool management and global state handling.
type Runner struct {
	program *bytecode.Program
	vmPool  *vm.VMPool

	// stateful mode fields (only used when Stateful is true)
	Stateful      bool          // whether stateful globals mode is enabled
	sharedGlobals []value.Value // program-owned globals shared across runs
	runMu         sync.Mutex    // serializes top-level Run calls in stateful mode
}

// New creates a new Runner for the given compiled bytecode program.
// It sets up a VM pool for efficient reuse across executions.
func New(program *bytecode.Program) *Runner {
	return &Runner{
		program: program,
		vmPool:  vm.NewVMPool(program),
	}
}

// Run executes a function in the program with the given arguments.
// Parameters are automatically converted to value.Value using FromInterface.
func (r *Runner) Run(funcName string, params ...any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return r.RunWithContext(ctx, funcName, params...)
}

// RunWithContext executes a function with context for timeout/cancellation control.
func (r *Runner) RunWithContext(ctx context.Context, funcName string, params ...any) (any, error) {
	args := make([]value.Value, len(params))
	for i, param := range params {
		args[i] = value.FromInterface(param)
	}
	result, err := r.RunWithValues(ctx, funcName, args)
	if err != nil {
		return nil, err
	}
	return UnwrapResult(result), nil
}

// RunWithValues executes a function with pre-converted Value arguments.
// This is more efficient when calling the same function repeatedly with the same types.
func (r *Runner) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
	if r.Stateful {
		r.runMu.Lock()
		defer r.runMu.Unlock()

		v := r.vmPool.Get()
		v.BindSharedGlobals(&r.sharedGlobals)
		result, err := v.ExecuteWithValues(funcName, ctx, args)
		v.UnbindSharedGlobals()
		r.vmPool.Put(v)
		return result, err
	}

	v := r.vmPool.Get()
	result, err := v.ExecuteWithValues(funcName, ctx, args)
	r.vmPool.Put(v)
	return result, err
}

// InternalProgram exposes the compiled bytecode program for testing/debugging.
func (r *Runner) InternalProgram() *bytecode.Program { return r.program }

// InitSharedGlobals initializes the shared globals slice for stateful mode.
// This should be called after Build() with stateful mode enabled.
func (r *Runner) InitSharedGlobals() {
	if r.Stateful {
		r.sharedGlobals = make([]value.Value, len(r.program.Globals))
		if len(r.program.InitialGlobals) == len(r.sharedGlobals) {
			copy(r.sharedGlobals, r.program.InitialGlobals)
		}
	}
}

// UnwrapResult converts internal multi-return value.Value slices to []any.
// A single return value is unwrapped directly; multiple return values become []any.
func UnwrapResult(result value.Value) any {
	iface := result.Interface()
	if vals, ok := iface.([]value.Value); ok {
		out := make([]any, len(vals))
		for i, v := range vals {
			out[i] = v.Interface()
		}
		return out
	}
	return iface
}

// DefaultTimeout is the default execution timeout.
const DefaultTimeout = 10 * time.Second

// ExecuteInit runs the program's init() function if present and snapshots the globals.
// This must be called after compilation and before any user-facing Run calls.
func ExecuteInit(program *bytecode.Program) error {
	if _, hasInit := program.Functions["init#1"]; hasInit {
		initVM := vm.New(program)
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		if _, err := initVM.Execute("init", ctx); err != nil {
			return err
		}
		snap := make([]value.Value, len(initVM.Globals()))
		copy(snap, initVM.Globals())
		program.InitialGlobals = snap
	}
	return nil
}
