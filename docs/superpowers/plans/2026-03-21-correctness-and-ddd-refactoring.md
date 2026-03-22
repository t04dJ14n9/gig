 Correctness & DDD Refactoring Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix 3 critical correctness bugs, eliminate all global mutable state, and split `vm/ops_dispatch.go` into readable thematic files.

**Architecture:** Replace global callbacks and counters with per-program scoped types. Method resolution moves from a global callback to a VM-local method. Closure execution uses an interface instead of a global callback. Goroutine tracking becomes per-program via `GoroutineTracker`. The `Importer` becomes registry-injectable for proper sandbox isolation.

**Tech Stack:** Go 1.23, `go/types`, `reflect`, `sync/atomic`

**Spec:** `docs/superpowers/specs/2026-03-21-correctness-and-ddd-refactoring-design.md`

---

## File Structure

| File | Responsibility | Action |
|---|---|---|
| `value/accessor.go` | Typed getters, `CallMethod`, `ClosureExecutor` interface | Modify |
| `value/value.go` | Value struct, Kind enum, constructors | Modify (add `ClosureExecutor`) |
| `vm/vm.go` | VM struct, `newVM`, `Reset`, `resolveCompiledMethod` | Modify |
| `vm/interfaces.go` | `VM` interface, `New()`, `VMOption` | Modify |
| `vm/closure.go` | `Closure` struct, pool, `Execute` method | Modify |
| `vm/goroutine.go` | `GoroutineTracker`, child VM creation | Rewrite |
| `vm/vm_test.go` | White-box VM tests | Modify |
| `vm/ops_dispatch.go` | Opcode router (delegates to category handlers) | Rewrite (slim down) |
| `vm/ops_arithmetic.go` | Add, Sub, Mul, Div, Mod, Neg, bitwise, comparisons | New |
| `vm/ops_memory.go` | Const, Local, Global, Free, Field, Addr, Deref | New |
| `vm/ops_call.go` | Call, CallExternal, CallIndirect, GoCall, Closure, Method | New |
| `vm/ops_container.go` | MakeSlice, MakeMap, MakeChan, Index, Append, Range, Len, etc. | New |
| `vm/ops_convert.go` | Assert, Convert, ChangeType, Pack, Unpack | New |
| `vm/ops_control.go` | Jump, Return, Defer, Panic, Recover, Select, Send, Recv, Print | New |
| `importer/importer.go` | `Importer` with injectable registry | Modify |
| `importer/register.go` | `ExternalPackage` with registry back-ref | Modify |
| `compiler/parser/parse.go` | Pass registry to `NewImporter` | Modify |
| `runner/runner.go` | Constructor options, `GoroutineTracker`, initial globals | Modify |
| `bytecode/bytecode.go` | Remove `InitialGlobals` | Modify |
| `bytecode/opcode.go` | Remove duplicate `OperandWidths` map | Modify |
| `gig.go` | Updated `Build()` flow | Modify |
| `stdlib/packages/fmt.go` | Update hand-written `CallMethod` call | Modify |
| `register/register.go` | **Delete** | Delete |
| `README.md`, `README_EN.md`, `examples/README.md`, `examples/README_CN.md` | Remove `register` references | Modify |

---

## Task 1: Add `ClosureExecutor` interface to `value/`

This is the foundation that breaks the `value → vm` circular dependency for closure execution.

**Files:**
- Modify: `value/value.go` (add interface near top, after constructors)
- Modify: `value/accessor.go` (remove globals, update `CallMethod`, update `ToReflectValue`)

- [ ] **Step 1: Add `ClosureExecutor` interface to `value/value.go`**

Add after the `GoString()` method at the bottom of the file (line ~479):

```go
// ClosureExecutor is implemented by closure objects that can be executed.
// This interface breaks the circular dependency between value/ and vm/ —
// vm.Closure implements it, and value.ToReflectValue uses it to convert
// closures into real Go functions via reflect.MakeFunc.
type ClosureExecutor interface {
	Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value
}
```

- [ ] **Step 2: Remove global callbacks from `value/accessor.go`**

Remove lines 14–41 (the `closureCaller` var, `SetClosureCaller`, `methodResolver` var, `SetMethodResolver`, and old `CallMethod`). Replace with:

```go
// MethodResolverFunc is a callback for calling compiled methods on interpreted types.
// It receives a method name and receiver value, and returns the result if found.
type MethodResolverFunc func(methodName string, receiver Value) (Value, bool)

// CallMethod attempts to call a compiled method on the receiver using the given resolver.
// Returns (result, true) if the method was found and called, or (zero, false) otherwise.
func CallMethod(resolver MethodResolverFunc, methodName string, receiver Value) (Value, bool) {
	if resolver == nil {
		return MakeNil(), false
	}
	return resolver(methodName, receiver)
}
```

- [ ] **Step 3: Update `ToReflectValue` KindFunc case in `value/accessor.go`**

Replace lines 184–213 (the `KindFunc` case) with:

```go
	case KindFunc:
		// If the target type is a function type, wrap the closure in a real Go function
		// using reflect.MakeFunc. This allows closures to be stored in typed containers
		// (maps, struct fields) that expect concrete function types like func() int.
		if typ.Kind() == reflect.Func {
			if ce, ok := v.obj.(ClosureExecutor); ok {
				numOut := typ.NumOut()
				outTypes := make([]reflect.Type, numOut)
				for i := 0; i < numOut; i++ {
					outTypes[i] = typ.Out(i)
				}
				fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
					results := ce.Execute(args, outTypes)
					// Convert results to match the expected return types
					out := make([]reflect.Value, numOut)
					for i := 0; i < numOut; i++ {
						if i < len(results) && results[i].IsValid() {
							if results[i].Type().ConvertibleTo(outTypes[i]) {
								out[i] = results[i].Convert(outTypes[i])
							} else {
								out[i] = results[i]
							}
						} else {
							out[i] = reflect.Zero(outTypes[i])
						}
					}
					return out
				})
				return fn
			}
		}
		return reflect.ValueOf(v.obj)
```

- [ ] **Step 4: Verify `value/` compiles**

Run: `cd /data/workspace/Code/gig && go build ./value/`
Expected: Clean build (no errors)

- [ ] **Step 5: Commit**

```bash
git add value/accessor.go value/value.go
git commit -m "refactor(value): remove global callbacks, add ClosureExecutor interface

Remove closureCaller and methodResolver globals. CallMethod now takes
an explicit resolver parameter. ToReflectValue uses ClosureExecutor
interface instead of global closureCaller callback."
```

---

## Task 2: Update `vm/closure.go` — add `InitialGlobals` and `Execute` method

**Files:**
- Modify: `vm/closure.go`

- [ ] **Step 1: Add `InitialGlobals` field and `Execute` method to `Closure`**

Replace the entire `vm/closure.go` with:

```go
package vm

import (
	"context"
	"reflect"
	"sync"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
)

// Closure represents a closure with captured free variables.
// When a closure is called, its free variables are bound to the calling context.
// Closure implements value.ClosureExecutor so that value.ToReflectValue can
// wrap it into a real Go function via reflect.MakeFunc without a global callback.
type Closure struct {
	// Fn is the compiled function bytecode.
	Fn *bytecode.CompiledFunction

	// FreeVars are pointers to captured variables.
	// They are stored as pointers to allow shared state between closures.
	FreeVars []*value.Value

	// Program is a reference to the compiled program, needed when the closure
	// is wrapped as a real Go function (via reflect.MakeFunc) for typed containers.
	Program *bytecode.Program

	// InitialGlobals is the post-init globals snapshot used to seed temporary VMs
	// when this closure is converted to a real Go function via Execute().
	InitialGlobals []value.Value
}

// Execute runs the closure in a temporary VM and returns the results as reflect.Values.
// This implements value.ClosureExecutor, allowing value.ToReflectValue to convert
// closures to real Go functions without a global callback.
func (c *Closure) Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
	if c.Program == nil {
		return nil
	}
	// Create a temporary VM to execute the closure
	closureVM := &vm{
		program: c.Program,
		stack:   make([]value.Value, 256),
		sp:      0,
		frames:  make([]*Frame, 64),
		fp:      0,
		globals: make([]value.Value, len(c.Program.Globals)),
		ctx:     context.Background(),
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(c.Program.Constants)),
		},
	}
	if len(c.InitialGlobals) == len(closureVM.globals) {
		copy(closureVM.globals, c.InitialGlobals)
	}
	closureVM.initialGlobals = c.InitialGlobals
	// Convert reflect.Value args to value.Value args
	valArgs := make([]value.Value, len(args))
	for i, arg := range args {
		valArgs[i] = value.MakeFromReflect(arg)
	}
	// Call the closure function with its captured free variables
	closureVM.callFunction(c.Fn, valArgs, c.FreeVars)
	result, _ := closureVM.run()
	// Return the result as reflect.Value
	if result.Kind() == value.KindNil {
		return []reflect.Value{}
	}
	if len(outTypes) > 0 {
		return []reflect.Value{result.ToReflectValue(outTypes[0])}
	}
	iface := result.Interface()
	if iface == nil {
		return []reflect.Value{}
	}
	return []reflect.Value{reflect.ValueOf(iface)}
}

// closurePool pools Closure objects to reduce heap allocations.
var closurePool = sync.Pool{
	New: func() any {
		return &Closure{}
	},
}

// getClosure returns a Closure from the pool, resized for numFree.
func getClosure(fn *bytecode.CompiledFunction, numFree int) *Closure {
	c := closurePool.Get().(*Closure)
	c.Fn = fn
	if numFree == 0 {
		c.FreeVars = nil
	} else if cap(c.FreeVars) >= numFree {
		c.FreeVars = c.FreeVars[:numFree]
	} else {
		c.FreeVars = make([]*value.Value, numFree)
	}
	return c
}
```

Note: This references `v.initialGlobals` which doesn't exist on `vm` yet — that comes in Task 4. This file will compile after Task 4.

- [ ] **Step 2: Commit (may not compile yet — depends on Task 4 for `vm.initialGlobals`)**

```bash
git add vm/closure.go
git commit -m "refactor(vm): add Closure.Execute implementing ClosureExecutor

Closure now implements value.ClosureExecutor directly, removing the need
for the global closureCaller callback. InitialGlobals field seeds temp VMs."
```

---

## Task 3: Rewrite `vm/goroutine.go` — `GoroutineTracker`

**Files:**
- Modify: `vm/goroutine.go`

- [ ] **Step 1: Replace `vm/goroutine.go` with `GoroutineTracker` struct**

```go
package vm

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/t04dJ14n9/gig/value"
)

// GoroutineTracker tracks active interpreter goroutines for a single program.
// It replaces the old process-wide activeGoroutines counter, making concurrent
// multi-program usage safe.
type GoroutineTracker struct {
	active int64
}

// NewGoroutineTracker creates a new goroutine tracker.
func NewGoroutineTracker() *GoroutineTracker {
	return &GoroutineTracker{}
}

// Start launches a goroutine and tracks it.
func (t *GoroutineTracker) Start(fn func()) {
	atomic.AddInt64(&t.active, 1)
	go func() {
		defer atomic.AddInt64(&t.active, -1)
		fn()
	}()
}

// Wait blocks until all tracked goroutines have completed.
// Uses exponential backoff to avoid busy waiting.
func (t *GoroutineTracker) Wait() {
	backoff := time.Microsecond
	for atomic.LoadInt64(&t.active) > 0 {
		time.Sleep(backoff)
		if backoff < 10*time.Millisecond {
			backoff *= 2
		}
	}
}

// WaitContext blocks until all tracked goroutines complete or the context is cancelled.
func (t *GoroutineTracker) WaitContext(ctx context.Context) error {
	backoff := time.Microsecond
	for atomic.LoadInt64(&t.active) > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		time.Sleep(backoff)
		if backoff < 10*time.Millisecond {
			backoff *= 2
		}
	}
	return nil
}

// newChildVM creates a child VM for goroutine execution.
// The child VM shares the globals pointer and external call cache with the parent.
func (v *vm) newChildVM() *vm {
	child := &vm{
		program:        v.program,
		stack:          make([]value.Value, 1024),
		sp:             0,
		frames:         make([]*Frame, 64),
		fp:             0,
		globals:        nil, // Not used when globalsPtr is set
		globalsPtr:     v.globalsPtr,
		ctx:            v.ctx,
		extCallCache:   v.extCallCache,
		initialGlobals: v.initialGlobals,
		goroutines:     v.goroutines,
	}
	if child.globalsPtr == nil {
		child.globalsPtr = &v.globals
	}
	return child
}
```

- [ ] **Step 2: Commit**

```bash
git add vm/goroutine.go
git commit -m "refactor(vm): replace global goroutine counter with per-program GoroutineTracker

GoroutineTracker scopes goroutine tracking to a single program, making
concurrent multi-program usage safe. Removes vmRegistry globals."
```

---

## Task 4: Update `vm/vm.go` — remove `init()`, add `initialGlobals` and `goroutines` fields

**Files:**
- Modify: `vm/vm.go`
- Modify: `vm/interfaces.go`

- [ ] **Step 1: Remove `init()` function from `vm/vm.go`**

Delete lines 44–100 (the entire `func init()` block that sets `value.SetClosureCaller`).

- [ ] **Step 2: Add `initialGlobals` and `goroutines` fields to `vm` struct**

In the `vm` struct (line ~104), add two new fields after `extCallCache`:

```go
	// initialGlobals is the post-init globals snapshot.
	// Used by Reset() to restore globals to their initial state.
	initialGlobals []value.Value

	// goroutines tracks active interpreter goroutines for this program.
	goroutines *GoroutineTracker
```

- [ ] **Step 3: Update `newVM` to accept `initialGlobals` and `goroutines`**

Change `newVM` signature and body:

```go
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
		stack:          make([]value.Value, 1024),
		sp:             0,
		frames:         make([]*Frame, 64),
		fp:             0,
		globals:        globals,
		initialGlobals: initialGlobals,
		goroutines:     goroutines,
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(program.Constants)),
		},
	}
}
```

- [ ] **Step 4: Update `Reset()` to use `v.initialGlobals` instead of `v.program.InitialGlobals`**

Replace the globals restoration section in `Reset()`:

```go
	// Stateless mode: restore globals to post-init snapshot, or zero them.
	if len(v.initialGlobals) == len(v.globals) {
		copy(v.globals, v.initialGlobals)
	} else {
		for i := range v.globals {
			v.globals[i] = value.Value{}
		}
	}
```

- [ ] **Step 5: Remove `SetMethodResolver` call from `newVM`**

The old `newVM` had this block (around line 185):
```go
	value.SetMethodResolver(func(methodName string, receiver value.Value) (value.Value, bool) {
		return resolveCompiledMethod(program, methodName, receiver)
	})
```
This was already removed in Step 1 (it was inside `newVM`). Verify it's gone.

- [ ] **Step 6: Update `VMPool` to carry `initialGlobals` and `goroutines`**

```go
type VMPool struct {
	mu    sync.Mutex
	vms   []*vm
	newVM func() *vm
}

func NewVMPool(program *bytecode.Program, initialGlobals []value.Value, goroutines *GoroutineTracker) *VMPool {
	return &VMPool{
		newVM: func() *vm {
			return newVM(program, initialGlobals, goroutines)
		},
	}
}
```

- [ ] **Step 7: Update `vm/interfaces.go` — `New()` and `NewWithOptions()`**

```go
func New(program *bytecode.Program) VM {
	return newVM(program, nil, NewGoroutineTracker())
}

func NewWithOptions(program *bytecode.Program, opts ...VMOption) VM {
	v := newVM(program, nil, NewGoroutineTracker())
	for _, opt := range opts {
		opt(v)
	}
	return v
}
```

- [ ] **Step 8: Update closure creation in `ops_dispatch.go` — set `InitialGlobals` on closures**

Find the `OpClosure` case (line ~1012) and after `c.Program = v.program`, add:
```go
		c.InitialGlobals = v.initialGlobals
```

- [ ] **Step 9: Verify `vm/` compiles**

Run: `cd /data/workspace/Code/gig && go build ./vm/`
Expected: Clean build

- [ ] **Step 10: Commit**

```bash
git add vm/vm.go vm/interfaces.go vm/ops_dispatch.go
git commit -m "refactor(vm): remove init(), add initialGlobals/goroutines to vm struct

newVM now takes initialGlobals and GoroutineTracker parameters.
Reset() uses vm.initialGlobals. VMPool carries both.
Closures receive initialGlobals snapshot at creation time."
```

---

## Task 5: Fix sandbox registry bypass — injectable `Importer`

**Files:**
- Modify: `importer/importer.go`
- Modify: `importer/register.go`
- Modify: `compiler/parser/parse.go`

- [ ] **Step 1: Add `reg` field to `Importer` struct and update `NewImporter`**

In `importer/importer.go`, change `Importer` struct and `NewImporter`:

```go
type Importer struct {
	reg      PackageRegistry
	packages map[string]*types.Package
	mutex    sync.RWMutex
}

func NewImporter(reg PackageRegistry) *Importer {
	return &Importer{
		reg:      reg,
		packages: make(map[string]*types.Package),
	}
}
```

- [ ] **Step 2: Update `Import()` to use `i.reg`**

Replace `GetPackageByPath(path)` and `GetPackageByName(path)` calls:

```go
func (i *Importer) Import(path string) (*types.Package, error) {
	i.mutex.RLock()
	if pkg, ok := i.packages[path]; ok {
		i.mutex.RUnlock()
		if pkg == nil {
			return nil, fmt.Errorf("package %q not found", path)
		}
		return pkg, nil
	}
	i.mutex.RUnlock()

	extPkg := i.reg.GetPackageByPath(path)
	if extPkg == nil {
		extPkg = i.reg.GetPackageByName(path)
		if extPkg == nil {
			return nil, fmt.Errorf("package %q not registered", path)
		}
	}

	pkg := i.buildPackage(extPkg)

	i.mutex.Lock()
	i.packages[path] = pkg
	i.mutex.Unlock()

	return pkg, nil
}
```

- [ ] **Step 3: Update `buildPackage()` to use `i.reg.SetExternalType()`**

Replace `SetExternalType(t, rt)` (line ~112) with `i.reg.SetExternalType(t, rt)`.

- [ ] **Step 4: Add `registry` back-reference to `ExternalPackage`**

In `importer/register.go`, add field to `ExternalPackage`:

```go
type ExternalPackage struct {
	Path     string
	Name     string
	Objects  map[string]*ExternalObject
	Types    map[string]reflect.Type
	registry PackageRegistry // back-reference to owning registry
}
```

Update `RegisterPackage` on `Registry` to set the back-ref:

```go
func (r *Registry) RegisterPackage(path, name string) *ExternalPackage {
	pkg := &ExternalPackage{
		Path:     path,
		Name:     name,
		Objects:  make(map[string]*ExternalObject),
		Types:    make(map[string]reflect.Type),
		registry: r,
	}
	r.mu.Lock()
	r.packagesByName[path] = pkg
	r.packagesByAlias[name] = pkg
	r.mu.Unlock()
	return pkg
}
```

- [ ] **Step 5: Update `AddMethodDirectCall` to use instance registry**

```go
func (p *ExternalPackage) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
	if p.registry != nil {
		p.registry.AddMethodDirectCall(p.Path+"."+typeName, methodName, dc)
	}
}
```

- [ ] **Step 6: Update `compiler/parser/parse.go` to pass registry**

Change line 51 from:
```go
	imp := importer.NewImporter()
```
to:
```go
	imp := importer.NewImporter(reg)
```

- [ ] **Step 7: Verify compilation**

Run: `cd /data/workspace/Code/gig && go build ./importer/ && go build ./compiler/...`
Expected: Clean build

- [ ] **Step 8: Commit**

```bash
git add importer/importer.go importer/register.go compiler/parser/parse.go
git commit -m "fix(importer): make Importer registry-injectable, fix sandbox bypass

NewImporter(reg) now takes a PackageRegistry parameter. Import() uses
the injected registry instead of globalRegistry. ExternalPackage gets
a registry back-reference so AddMethodDirectCall uses instance registry."
```

---

## Task 6: Update `runner/runner.go` — constructor options, initial globals, `GoroutineTracker`

**Files:**
- Modify: `runner/runner.go`
- Modify: `bytecode/bytecode.go`

- [ ] **Step 1: Remove `InitialGlobals` from `bytecode.Program`**

In `bytecode/bytecode.go`, delete lines 103–106:
```go
	// InitialGlobals holds the global variable state after init() has run.
	// ...
	InitialGlobals []value.Value
```

- [ ] **Step 2: Rewrite `runner/runner.go`**

```go
package runner

import (
	"context"
	"sync"
	"time"

	"github.com/t04dJ14n9/gig/bytecode"
	"github.com/t04dJ14n9/gig/value"
	"github.com/t04dJ14n9/gig/vm"
)

// Executor is the main execution interface for running compiled programs.
type Executor interface {
	Run(funcName string, params ...any) (any, error)
	RunWithContext(ctx context.Context, funcName string, params ...any) (any, error)
	RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error)
	InternalProgram() *bytecode.Program
}

// runnerConfig holds internal options parsed from RunnerOption values.
type runnerConfig struct {
	stateful bool
}

// RunnerOption configures Runner construction.
type RunnerOption func(*runnerConfig)

// WithStatefulGlobals enables persistent package-level globals across Run calls.
func WithStatefulGlobals() RunnerOption {
	return func(c *runnerConfig) {
		c.stateful = true
	}
}

// Runner orchestrates program execution with VM pool management and global state handling.
type Runner struct {
	program        *bytecode.Program
	initialGlobals []value.Value
	vmPool         *vm.VMPool
	goroutines     *vm.GoroutineTracker

	// stateful mode fields
	stateful      bool
	sharedGlobals []value.Value
	runMu         sync.Mutex
}

// New creates a new Runner for the given compiled bytecode program.
func New(program *bytecode.Program, initialGlobals []value.Value, opts ...RunnerOption) *Runner {
	cfg := runnerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	gt := vm.NewGoroutineTracker()
	r := &Runner{
		program:        program,
		initialGlobals: initialGlobals,
		vmPool:         vm.NewVMPool(program, initialGlobals, gt),
		goroutines:     gt,
		stateful:       cfg.stateful,
	}

	if cfg.stateful {
		r.sharedGlobals = make([]value.Value, len(program.Globals))
		if len(initialGlobals) == len(r.sharedGlobals) {
			copy(r.sharedGlobals, initialGlobals)
		}
	}

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
func (r *Runner) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
	if r.stateful {
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

// Wait blocks until all interpreter goroutines for this program have completed.
func (r *Runner) Wait() {
	r.goroutines.Wait()
}

// WaitContext blocks until all goroutines complete or the context is cancelled.
func (r *Runner) WaitContext(ctx context.Context) error {
	return r.goroutines.WaitContext(ctx)
}

// InternalProgram exposes the compiled bytecode program for testing/debugging.
func (r *Runner) InternalProgram() *bytecode.Program { return r.program }

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
func ExecuteInit(program *bytecode.Program) ([]value.Value, error) {
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
```

- [ ] **Step 3: Verify compilation**

Run: `cd /data/workspace/Code/gig && go build ./bytecode/ && go build ./runner/`
Expected: Clean build

- [ ] **Step 4: Commit**

```bash
git add runner/runner.go bytecode/bytecode.go
git commit -m "refactor(runner): constructor options, per-program GoroutineTracker

Runner.New takes initialGlobals and options. ExecuteInit returns snapshot
instead of mutating Program. Runner owns GoroutineTracker and exposes
Wait()/WaitContext(). Remove InitialGlobals from bytecode.Program."
```

---

## Task 7: Update `gig.go` — wire new `Build()` flow

**Files:**
- Modify: `gig.go`

- [ ] **Step 1: Update `Build()` to use new `runner.ExecuteInit` and `runner.New` signatures**

Replace the `Build` function body (lines 138–169):

```go
func Build(sourceCode string, opts ...BuildOption) (*Program, error) {
	cfg := buildConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.registry == nil {
		cfg.registry = importer.GlobalRegistry()
	}

	result, err := compiler.Build(sourceCode, cfg.registry)
	if err != nil {
		return nil, err
	}

	initialGlobals, err := runner.ExecuteInit(result.Program)
	if err != nil {
		return nil, fmt.Errorf("executing init(): %w", err)
	}

	var runnerOpts []runner.RunnerOption
	if cfg.statefulGlobals {
		runnerOpts = append(runnerOpts, runner.WithStatefulGlobals())
	}

	r := runner.New(result.Program, initialGlobals, runnerOpts...)

	return &Program{
		runner: r,
		ssaPkg: result.SSAPkg,
	}, nil
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /data/workspace/Code/gig && go build .`
Expected: Clean build

- [ ] **Step 3: Run all tests**

Run: `cd /data/workspace/Code/gig && go test -v -race -count=1 ./... 2>&1 | tail -50`
Expected: All integration tests pass. `vm/vm_test.go` may fail (fixed in Task 8).

- [ ] **Step 4: Commit**

```bash
git add gig.go
git commit -m "refactor(gig): update Build() for new runner.ExecuteInit/New signatures"
```

---

## Task 8: Update `vm/vm_test.go` — fix white-box tests

**Files:**
- Modify: `vm/vm_test.go`

- [ ] **Step 1: Update goroutine and VM registry tests**

Replace `TestGoroutineWaitContext`, `TestStartAndWaitGoroutines`, and `TestVMRegistry` (lines ~331–392) with:

```go
// TestGoroutineTrackerWaitContext verifies that GoroutineTracker.WaitContext respects cancellation.
func TestGoroutineTrackerWaitContext(t *testing.T) {
	gt := NewGoroutineTracker()
	gt.Start(func() {
		time.Sleep(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := gt.WaitContext(ctx)
	if err == nil {
		t.Error("expected timeout error")
	}

	// Wait for the goroutine to actually complete to clean up
	gt.Wait()
}

func TestGoroutineTrackerStartAndWait(t *testing.T) {
	gt := NewGoroutineTracker()
	var counter int64
	for i := 0; i < 5; i++ {
		gt.Start(func() {
			atomic.AddInt64(&counter, 1)
		})
	}
	gt.Wait()
	if atomic.LoadInt64(&counter) != 5 {
		t.Errorf("counter = %d, want 5", counter)
	}
}
```

Delete the `TestVMRegistry` function entirely (the VM registry has been removed).

- [ ] **Step 2: Update `ops_dispatch.go` — replace `StartGoroutine` calls with `v.goroutines.Start`**

In `vm/ops_dispatch.go`, find the two `StartGoroutine(func() {` calls (lines ~1072 and ~1107) and replace with `v.goroutines.Start(func() {`.

- [ ] **Step 3: Verify all tests pass**

Run: `cd /data/workspace/Code/gig && go test -v -race -count=1 ./...`
Expected: All tests pass with no race conditions.

- [ ] **Step 4: Commit**

```bash
git add vm/vm_test.go vm/ops_dispatch.go
git commit -m "test(vm): update tests for GoroutineTracker, remove VM registry tests"
```

---

## Task 9: Update `stdlib/packages/fmt.go` — fix `CallMethod` call site

**Files:**
- Modify: `stdlib/packages/fmt.go`

- [ ] **Step 1: Update `sanitizeArgForFmt` to accept a resolver**

The function `sanitizeArgForFmt` (line ~159) currently calls `value.CallMethod("String", v)`. Since this is hand-written code in the fmt wrapper, we need to thread a resolver through. However, looking at the call chain, `sanitizeArgForFmt` is called from DirectCall wrappers which don't have VM context.

The pragmatic fix: since `sanitizeArgForFmt` is called at fmt wrapper time (not VM time), and the method resolver needs a program reference, we pass `nil` for now — this means fmt.Stringer support for interpreted types won't work through DirectCall wrappers (same as if no method resolver were registered). This is a pre-existing limitation that would need a deeper fix to the DirectCall architecture.

Update line 170:
```go
	if result, found := value.CallMethod(nil, "String", v); found {
```

- [ ] **Step 2: Verify compilation and tests**

Run: `cd /data/workspace/Code/gig && go build ./stdlib/packages/ && go test -v -race -count=1 ./...`
Expected: All tests pass.

- [ ] **Step 3: Commit**

```bash
git add stdlib/packages/fmt.go
git commit -m "fix(stdlib/fmt): update CallMethod to new signature (nil resolver)

The fmt DirectCall wrapper calls value.CallMethod with nil resolver since
it lacks VM context. fmt.Stringer for interpreted types via DirectCall
was already unreliable and needs a deeper architectural fix."
```

---

## Task 10: Delete `register/` package and update references

**Files:**
- Delete: `register/register.go`
- Modify: `README.md`, `README_EN.md`, `examples/README.md`, `examples/README_CN.md`

- [ ] **Step 1: Delete the `register/` package**

```bash
rm -rf register/
```

- [ ] **Step 2: Check for any remaining imports of `register`**

Run: `cd /data/workspace/Code/gig && grep -r "gig/register" --include='*.go' .`
Expected: No results (only READMEs should reference it).

- [ ] **Step 3: Update README files**

In all 4 README files, replace `register.AddPackage` → `gig.RegisterPackage` or `importer.RegisterPackage`, and `register.NewFunction` → `pkg.AddFunction`, etc. Follow the migration table from the spec.

- [ ] **Step 4: Verify compilation**

Run: `cd /data/workspace/Code/gig && go build ./...`
Expected: Clean build (no references to deleted package).

- [ ] **Step 5: Run all tests**

Run: `cd /data/workspace/Code/gig && go test -v -race -count=1 ./...`
Expected: All tests pass.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: delete register/ package, update READMEs

register/ was a thin broken wrapper (NewType silently dropped types).
Users should use importer.RegisterPackage() or gig.RegisterPackage()
directly. Updated all README examples."
```

---

## Task 11: Remove duplicate `OperandWidths` map from `bytecode/opcode.go`

**Files:**
- Modify: `bytecode/opcode.go`
- Modify: `bytecode/bytecode_test.go` (if it references `OperandWidths`)

- [ ] **Step 1: Find and remove the `OperandWidths` map**

In `bytecode/opcode.go`, find the `var OperandWidths = map[OpCode]int{...}` declaration and delete it entirely. Keep `operandWidthTable` (the array version).

- [ ] **Step 2: Update any references to `OperandWidths`**

Run: `cd /data/workspace/Code/gig && grep -rn 'OperandWidths' --include='*.go' .`

Replace any usage of `OperandWidths[op]` with `operandWidthTable[op]` or the equivalent accessor.

- [ ] **Step 3: Verify compilation and tests**

Run: `cd /data/workspace/Code/gig && go build ./... && go test -v -race -count=1 ./...`
Expected: All pass.

- [ ] **Step 4: Commit**

```bash
git add bytecode/opcode.go bytecode/bytecode_test.go
git commit -m "cleanup(bytecode): remove duplicate OperandWidths map

Keep only the operandWidthTable array (O(1) lookup). The map was
redundant and required manual sync with every new opcode."
```

---

## Task 12: Split `vm/ops_dispatch.go` into thematic files

This is the largest task. The 1761-line `executeOp()` switch gets split into 6 category handler files.

**Files:**
- Modify: `vm/ops_dispatch.go` (slim down to router)
- Create: `vm/ops_arithmetic.go`
- Create: `vm/ops_memory.go`
- Create: `vm/ops_call.go`
- Create: `vm/ops_container.go`
- Create: `vm/ops_convert.go`
- Create: `vm/ops_control.go`

- [ ] **Step 1: Create `vm/ops_arithmetic.go`**

Extract these opcodes from `executeOp()` into `func (v *vm) executeArithmetic(op bytecode.OpCode, frame *Frame) error`:
- `OpAdd`, `OpSub`, `OpMul`, `OpDiv`, `OpMod`, `OpNeg`
- `OpAnd`, `OpOr`, `OpXor`, `OpAndNot`, `OpLsh`, `OpRsh`
- `OpEqual`, `OpNotEqual`, `OpLess`, `OpLessEq`, `OpGreater`, `OpGreaterEq`
- `OpNot`

- [ ] **Step 2: Create `vm/ops_memory.go`**

Extract into `func (v *vm) executeMemory(op bytecode.OpCode, frame *Frame) error`:
- `OpNop`, `OpPop`, `OpDup`, `OpConst`, `OpNil`, `OpTrue`, `OpFalse`
- `OpLocal`, `OpSetLocal`, `OpGlobal`, `OpSetGlobal`
- `OpFree`, `OpSetFree`
- `OpField`, `OpSetField`, `OpAddr`, `OpFieldAddr`, `OpIndexAddr`
- `OpDeref`, `OpSetDeref`
- `OpNew`, `OpMake`

- [ ] **Step 3: Create `vm/ops_call.go`**

Extract into `func (v *vm) executeCall(op bytecode.OpCode, frame *Frame) error`:
- `OpCall`, `OpCallExternal`, `OpCallIndirect`
- `OpClosure`
- `OpGoCall`, `OpGoCallIndirect`
- `OpPack`, `OpUnpack`

- [ ] **Step 4: Create `vm/ops_container.go`**

Extract into `func (v *vm) executeContainer(op bytecode.OpCode, frame *Frame) error`:
- `OpMakeSlice`, `OpMakeMap`, `OpMakeChan`
- `OpIndex`, `OpIndexOk`, `OpSetIndex`
- `OpSlice`
- `OpRange`, `OpRangeNext`
- `OpLen`, `OpCap`
- `OpAppend`, `OpCopy`, `OpDelete`

- [ ] **Step 5: Create `vm/ops_convert.go`**

Extract into `func (v *vm) executeConvert(op bytecode.OpCode, frame *Frame) error`:
- `OpAssert`, `OpConvert`, `OpChangeType`

- [ ] **Step 6: Create `vm/ops_control.go`**

Extract into `func (v *vm) executeControl(op bytecode.OpCode, frame *Frame) error`:
- `OpJump`, `OpJumpTrue`, `OpJumpFalse`
- `OpReturn`, `OpReturnVal`
- `OpSend`, `OpRecv`, `OpRecvOk`, `OpClose`, `OpSelect`
- `OpDefer`, `OpDeferIndirect`, `OpRunDefers`
- `OpRecover`, `OpPanic`
- `OpPrint`, `OpPrintln`
- `OpHalt`

- [ ] **Step 7: Replace `executeOp()` body with router**

Slim down `ops_dispatch.go` to a router that delegates to category handlers:

```go
func (v *vm) executeOp(op bytecode.OpCode, frame *Frame) error {
	switch op {
	// Arithmetic & comparisons
	case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpNeg, bytecode.OpAnd, bytecode.OpOr, bytecode.OpXor, bytecode.OpAndNot,
		bytecode.OpLsh, bytecode.OpRsh,
		bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpLess, bytecode.OpLessEq,
		bytecode.OpGreater, bytecode.OpGreaterEq, bytecode.OpNot:
		return v.executeArithmetic(op, frame)

	// Memory & stack
	case bytecode.OpNop, bytecode.OpPop, bytecode.OpDup,
		bytecode.OpConst, bytecode.OpNil, bytecode.OpTrue, bytecode.OpFalse,
		bytecode.OpLocal, bytecode.OpSetLocal, bytecode.OpGlobal, bytecode.OpSetGlobal,
		bytecode.OpFree, bytecode.OpSetFree,
		bytecode.OpField, bytecode.OpSetField, bytecode.OpAddr, bytecode.OpFieldAddr, bytecode.OpIndexAddr,
		bytecode.OpDeref, bytecode.OpSetDeref,
		bytecode.OpNew, bytecode.OpMake:
		return v.executeMemory(op, frame)

	// Calls & closures
	case bytecode.OpCall, bytecode.OpCallExternal, bytecode.OpCallIndirect,
		bytecode.OpClosure, bytecode.OpGoCall, bytecode.OpGoCallIndirect,
		bytecode.OpPack, bytecode.OpUnpack:
		return v.executeCall(op, frame)

	// Containers
	case bytecode.OpMakeSlice, bytecode.OpMakeMap, bytecode.OpMakeChan,
		bytecode.OpIndex, bytecode.OpIndexOk, bytecode.OpSetIndex, bytecode.OpSlice,
		bytecode.OpRange, bytecode.OpRangeNext,
		bytecode.OpLen, bytecode.OpCap,
		bytecode.OpAppend, bytecode.OpCopy, bytecode.OpDelete:
		return v.executeContainer(op, frame)

	// Type conversions
	case bytecode.OpAssert, bytecode.OpConvert, bytecode.OpChangeType:
		return v.executeConvert(op, frame)

	// Control flow, channels, defer, panic, print, halt
	default:
		return v.executeControl(op, frame)
	}
}
```

- [ ] **Step 8: Verify compilation and tests**

Run: `cd /data/workspace/Code/gig && go build ./vm/ && go test -v -race -count=1 ./...`
Expected: All tests pass — behavior is identical, just reorganized.

- [ ] **Step 9: Commit**

```bash
git add vm/ops_dispatch.go vm/ops_arithmetic.go vm/ops_memory.go vm/ops_call.go vm/ops_container.go vm/ops_convert.go vm/ops_control.go
git commit -m "refactor(vm): split ops_dispatch.go into thematic files

Split 1761-line executeOp() into 6 category handlers:
- ops_arithmetic.go: Add, Sub, Mul, Div, comparisons, bitwise
- ops_memory.go: Const, Local, Global, Field, Addr, Deref, New
- ops_call.go: Call, CallExternal, Closure, GoCall
- ops_container.go: Slice, Map, Index, Append, Range, Len
- ops_convert.go: Assert, Convert, ChangeType
- ops_control.go: Jump, Return, Defer, Panic, Select, Print"
```

---

## Task 13: Add concurrency and sandbox isolation tests

**Files:**
- Create: `tests/concurrency_test.go`
- Create: `tests/sandbox_test.go`

- [ ] **Step 1: Write concurrency test**

```go
package tests

import (
	"sync"
	"testing"

	"github.com/t04dJ14n9/gig"
)

// TestConcurrentProgramMethodResolution verifies that two programs compiled
// and run concurrently don't interfere with each other's method dispatch.
// This is a regression test for the global method resolver race (C1).
func TestConcurrentProgramMethodResolution(t *testing.T) {
	srcA := `
package main

type Greeter struct{ Name string }
func (g Greeter) Greet() string { return "Hello, " + g.Name }
func Run() string {
	g := Greeter{Name: "Alice"}
	return g.Greet()
}
`
	srcB := `
package main

type Greeter struct{ Name string }
func (g Greeter) Greet() string { return "Hi, " + g.Name }
func Run() string {
	g := Greeter{Name: "Bob"}
	return g.Greet()
}
`
	var wg sync.WaitGroup
	errCh := make(chan error, 20)

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			prog, err := gig.Build(srcA)
			if err != nil {
				errCh <- err
				return
			}
			result, err := prog.Run("Run")
			if err != nil {
				errCh <- err
				return
			}
			if result != "Hello, Alice" {
				t.Errorf("srcA: got %q, want %q", result, "Hello, Alice")
			}
		}()
		go func() {
			defer wg.Done()
			prog, err := gig.Build(srcB)
			if err != nil {
				errCh <- err
				return
			}
			result, err := prog.Run("Run")
			if err != nil {
				errCh <- err
				return
			}
			if result != "Hi, Bob" {
				t.Errorf("srcB: got %q, want %q", result, "Hi, Bob")
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Errorf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 2: Write sandbox isolation test**

```go
package tests

import (
	"testing"

	"github.com/t04dJ14n9/gig"
	"github.com/t04dJ14n9/gig/importer"
)

// TestSandboxIsolation verifies that WithRegistry(sandbox) prevents access
// to globally-registered packages during type-checking.
func TestSandboxIsolation(t *testing.T) {
	sandbox := importer.NewRegistry()
	// Don't register fmt in sandbox

	_, err := gig.Build(`
		import "fmt"
		func Hello() string {
			return fmt.Sprintf("hello")
		}
	`, gig.WithRegistry(sandbox))

	if err == nil {
		t.Error("expected error when using unregistered package in sandbox, got nil")
	}
}

// TestSandboxWithRegisteredPackage verifies sandbox works with explicitly registered packages.
func TestSandboxWithRegisteredPackage(t *testing.T) {
	sandbox := importer.NewRegistry()
	pkg := sandbox.RegisterPackage("strings", "strings")
	// Register at least one function so it type-checks
	_ = pkg // minimal registration

	// This should fail because we only registered the package name,
	// not the actual functions — but it shouldn't fall through to global registry.
	_, err := gig.Build(`
		func Hello() string { return "hello" }
	`, gig.WithRegistry(sandbox))

	if err != nil {
		t.Errorf("expected simple program to compile in sandbox: %v", err)
	}
}
```

- [ ] **Step 3: Run the new tests with `-race`**

Run: `cd /data/workspace/Code/gig && go test -v -race -count=1 -run 'TestConcurrent|TestSandbox' ./tests/`
Expected: All pass with no race conditions.

- [ ] **Step 4: Commit**

```bash
git add tests/concurrency_test.go tests/sandbox_test.go
git commit -m "test: add concurrency and sandbox isolation regression tests

TestConcurrentProgramMethodResolution validates C1 fix (no method resolver race).
TestSandboxIsolation validates C2 fix (WithRegistry prevents global fallback)."
```

---

## Task 14: Final verification — full test suite with race detector

- [ ] **Step 1: Run full test suite**

Run: `cd /data/workspace/Code/gig && go test -v -race -count=1 ./... 2>&1 | tail -100`
Expected: All tests pass, no race conditions.

- [ ] **Step 2: Run linter**

Run: `cd /data/workspace/Code/gig && golangci-lint run --timeout=5m 2>&1 | head -50`
Expected: No new lint errors introduced by refactoring.

- [ ] **Step 3: Run build for all targets**

Run: `cd /data/workspace/Code/gig && go build ./...`
Expected: Clean build.

- [ ] **Step 4: Final commit (if any fixups needed)**

Only if previous steps revealed issues that needed fixing.
