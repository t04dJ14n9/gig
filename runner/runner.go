// Package runner provides program execution orchestration for the Gig interpreter.
//
// It encapsulates VM pool management, global state handling (stateless snapshot
// vs stateful shared globals), and result unpacking, providing a clean execution
// API on top of the low-level VM.
package runner

import (
	"context"
	"time"
	"unsafe"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
	"github.com/t04dJ14n9/gig/vm"
)

// runnerConfig holds internal options parsed from RunnerOption values.
type runnerConfig struct {
	stateful      bool
	maxGoroutines int
}

// RunnerOption configures Runner construction.
type RunnerOption func(*runnerConfig)

// WithStatefulGlobals enables persistent package-level globals across Run calls.
func WithStatefulGlobals() RunnerOption {
	return func(c *runnerConfig) {
		c.stateful = true
	}
}

// WithMaxGoroutines sets the maximum number of concurrent interpreter goroutines.
// The default is 10,000. Set to 0 to use the default.
func WithMaxGoroutines(n int) RunnerOption {
	return func(c *runnerConfig) {
		c.maxGoroutines = n
	}
}

// Runner orchestrates program execution with VM pool management and global state handling.
type Runner struct {
	program        *bytecode.CompiledProgram
	initialGlobals []value.Value
	vmPool         *vm.VMPool
	goroutines     *vm.GoroutineTracker

	// progKey is the key used for RegisterMethodResolver, stored for cleanup.
	progKey uintptr

	// stateful mode fields
	stateful bool
	shared   *vm.SharedGlobals // thread-safe shared globals for concurrent stateful execution
}

// New creates a new Runner for the given compiled bytecode program.
func New(program *bytecode.CompiledProgram, initialGlobals []value.Value, opts ...RunnerOption) *Runner {
	cfg := runnerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	gt := vm.NewGoroutineTracker()
	if cfg.maxGoroutines > 0 {
		gt.SetMaxGoroutines(cfg.maxGoroutines)
	}
	r := &Runner{
		program:        program,
		initialGlobals: initialGlobals,
		vmPool:         vm.NewVMPool(program, initialGlobals, gt),
		goroutines:     gt,
		stateful:       cfg.stateful,
		progKey:        uintptr(unsafe.Pointer(program)),
	}

	if cfg.stateful {
		r.shared = vm.NewSharedGlobals(initialGlobals, len(program.Globals))
		r.shared.InitExternalVars(program.ExternalVarValues)
		r.shared.InitZeroValues(program.GlobalZeroValues)
	}

	// Register per-program method resolver for fmt.Stringer support.
	// Uses program pointer as unique key. Thread-safe via sync.Map.
	value.RegisterMethodResolver(r.progKey, func(methodName string, receiver value.Value) (value.Value, bool) {
		return vm.ResolveCompiledMethod(program, methodName, receiver)
	})

	return r
}

// Run executes a function with the default timeout.
func (r *Runner) Run(funcName string, params ...any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return r.RunWithContext(ctx, funcName, params...)
}

// RunWithContext executes a function with context for timeout/cancellation.
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
// In stateful mode, multiple concurrent calls are allowed. Each call gets its
// own VM from the pool but shares the same locked SharedGlobals.
func (r *Runner) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
	if r.stateful {
		v := r.vmPool.Get()
		v.BindSharedGlobals(r.shared)
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

// Wait blocks until all interpreter goroutines for this program have completed.
func (r *Runner) Wait() {
	r.goroutines.Wait()
}

// WaitContext blocks until all goroutines complete or the context is cancelled.
func (r *Runner) WaitContext(ctx context.Context) error {
	return r.goroutines.WaitContext(ctx)
}

// InternalProgram exposes the compiled bytecode program for testing/debugging.
func (r *Runner) InternalProgram() *bytecode.CompiledProgram { return r.program }

// Close releases resources associated with the Runner.
// It unregisters the per-program method resolver to prevent memory leaks.
// Callers should defer Close() after creating a Runner.
func (r *Runner) Close() {
	value.UnregisterMethodResolver(r.progKey)
}

// UnwrapResult converts internal multi-return value.Value slices to []any.
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

// ExecuteInit runs the program's init() function if present and returns the globals snapshot.
// The snapshot should be passed to runner.New as initialGlobals.
func ExecuteInit(program *bytecode.CompiledProgram) ([]value.Value, error) {
	if _, hasInit := program.Functions["init#1"]; hasInit {
		initVM := vm.New(program)
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		if _, err := initVM.Execute("init", ctx); err != nil {
			return nil, err
		}
		snap := make([]value.Value, len(initVM.Globals()))
		copy(snap, initVM.Globals())
		return snap, nil
	}
	return nil, nil
}
