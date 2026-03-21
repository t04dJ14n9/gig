# Correctness Fixes & DDD Refactoring Design

**Date:** 2026-03-21
**Status:** Approved
**Scope:** Fix critical correctness bugs, eliminate global mutable state, restructure vm/ for readability

## Problem Statement

The Gig interpreter has three critical correctness bugs and pervasive global mutable state that makes concurrent multi-program usage unsafe and the code hard to reason about.

### Critical Bugs

1. **Global method resolver race (C1):** `value.SetMethodResolver()` is called in `vm.newVM()`, overwriting a package-level var. Concurrent programs sharing different function tables clobber each other's method dispatch.

2. **Sandbox registry bypass (C2):** `importer.NewImporter()` hardcodes `GetPackageByPath()` (global registry). `WithRegistry(sandbox)` only affects `autoImport()`, not type-checking — sandboxed programs can still resolve globally-registered types.

3. **`register.NewType()` drops the type (C3):** `p.inner.AddType(name, nil, doc)` always passes `nil` for the reflect.Type, making type registration via `register` package a silent no-op.

### Global Mutable State

| Global | Location | Risk |
|---|---|---|
| `value.methodResolver` | `value/accessor.go` | Race condition (C1) |
| `value.closureCaller` | `value/accessor.go` | Set once by `vm.init()`, captures no state — low risk but poor design |
| `vm.activeGoroutines` | `vm/goroutine.go` | Process-wide counter; `WaitGoroutines()` waits for ALL programs |
| `vm.vmRegistry` / `vmIDCounter` | `vm/goroutine.go` | Process-wide VM registry |
| `importer.globalRegistry` | `importer/register.go` | Read-only after init() — acceptable for stdlib, but sandbox bypass (C2) |
| `importer.typeCache` / `typePkgCache` | `importer/typeconv.go` | Pure caches, immutable mappings — safe, keep as-is |
| `importer.typeOf` | `importer/register.go` | Init-order dependent function var |

### Structural Issues

- `vm/ops_dispatch.go` is 1761 lines (single `executeOp()` function)
- `register/` package is a thin broken wrapper over `importer`
- `bytecode.Program` is semi-mutable (`InitialGlobals` set post-construction)
- `runner.Runner.Stateful` is a public mutable field

## Design

### Approach: Scoped Runtime Context

Eliminate all global mutable state by moving per-program concerns into scoped types. The VM resolves methods locally, goroutine tracking is per-program, and the Importer accepts an injectable registry.

### 1. Fix Method Resolver Race (C1)

**Before:** `value.SetMethodResolver(fn)` sets a global; `value.CallMethod()` reads it.

**After:**
- Remove `value.methodResolver`, `value.SetMethodResolver()`, `value.CallMethod()` from `value/accessor.go`
- Method resolution becomes a method on `*vm` in `vm/ops_call.go`
- Where `ops_dispatch.go` currently calls `value.CallMethod()`, it calls `v.resolveCompiledMethod()` instead
- `resolveCompiledMethod()` already exists on `vm/vm.go` and reads from `v.program` — no global needed

### 2. Fix Sandbox Registry Bypass (C2)

**Before:** `importer.NewImporter()` creates an Importer that always reads from `globalRegistry`.

**After:**
- `NewImporter(reg PackageRegistry)` — takes a registry parameter
- `Importer.Import(path)` calls `reg.GetPackageByPath(path)` instead of the global `GetPackageByPath(path)`
- `compiler/parser/parse.go` passes the registry it receives from Build: `importer.NewImporter(reg)`
- `Importer.buildPackage()` calls `reg.SetExternalType()` instead of global `SetExternalType()`

### 3. Fix register.NewType() (C3) & Remove register/ Package

**Before:** `register.NewType(name, typ, doc)` calls `p.inner.AddType(name, nil, doc)` — always nil.

**After:**
- Delete the `register/` package entirely
- Users use `importer.RegisterPackage()` or `gig.RegisterPackage()` directly
- Update any examples/docs that reference `register`

### 4. Eliminate Closure Caller Global

**Before:** `vm.init()` sets `value.SetClosureCaller(fn)` — a global callback.

**After:**
- Remove `value.closureCaller`, `value.SetClosureCaller()` from `value/accessor.go`
- Remove `func init()` from `vm/vm.go`
- Add `Closure.Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value` method on `vm/closure.go`
- `value.ToReflectValue()` for `KindFunc` needs to call closure execution. Since `value` can't import `vm`, the closure object stored in `Value.obj` will implement an interface:
  ```go
  // value/value.go
  type ClosureExecutor interface {
      Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value
  }
  ```
  The `vm.Closure` type implements this interface. `ToReflectValue` type-asserts `v.obj.(ClosureExecutor)` instead of calling a global callback.

### 5. Per-Program Goroutine Tracking

**Before:** `vm.activeGoroutines` is a process-wide `int64`; `WaitGoroutines()` waits for all programs.

**After:**
- New exported type `vm.GoroutineTracker`:
  ```go
  type GoroutineTracker struct {
      active int64
  }
  func (t *GoroutineTracker) Start(fn func())
  func (t *GoroutineTracker) Wait()
  func (t *GoroutineTracker) WaitContext(ctx context.Context) error
  ```
- `Runner` creates a `GoroutineTracker` per program
- `vm` receives a `*GoroutineTracker` field, uses it for `go` statement dispatch
- Remove global `activeGoroutines`, `vmRegistry`, `vmIDCounter`, `RegisterVM`, `UnregisterVM`

### 6. Immutable Program

**Before:** `runner.ExecuteInit()` mutates `program.InitialGlobals` after compilation.

**After:**
- `runner.ExecuteInit()` returns `([]value.Value, error)` — the globals snapshot
- `Runner` stores the snapshot internally: `runner.initialGlobals []value.Value`
- `bytecode.Program.InitialGlobals` field removed
- VM receives initial globals from `Runner` (via VMPool constructor or per-Get)

### 7. Runner Constructor Options

**Before:** `runner.New(program)` then `r.Stateful = true; r.InitSharedGlobals()`

**After:**
```go
type RunnerOption func(*runnerConfig)
func WithStatefulGlobals() RunnerOption
func New(program *bytecode.Program, initialGlobals []value.Value, opts ...RunnerOption) *Runner
```
`Stateful` becomes a private field set at construction. `InitSharedGlobals()` is called internally by the constructor when the option is set.

### 8. ExternalPackage.AddMethodDirectCall Uses Instance Registry

**Before:** `AddMethodDirectCall` delegates to `globalRegistry.AddMethodDirectCall()`.

**After:** Each `ExternalPackage` holds a reference to the registry it was created from. `AddMethodDirectCall` delegates to that instance's registry.

```go
type ExternalPackage struct {
    Path     string
    Name     string
    Objects  map[string]*ExternalObject
    Types    map[string]reflect.Type
    registry PackageRegistry // back-reference to owning registry
}
```

### 9. Split ops_dispatch.go into Thematic Files

Split the 1761-line `executeOp()` switch into category handler methods:

| New file | Opcodes | ~Lines |
|---|---|---|
| `ops_dispatch.go` | Router switch (delegates to handlers) | ~100 |
| `ops_arithmetic.go` | Add, Sub, Mul, Div, Mod, Neg, Cmp, bitwise | ~300 |
| `ops_control.go` | Jump, If, Return, Panic, Defer, Select | ~250 |
| `ops_memory.go` | Local, Global, Alloc, Field, FieldAddr, Store, Load | ~300 |
| `ops_call.go` | Call, CallExternal, CallDirect, Go, method resolution | ~350 |
| `ops_container.go` | MakeSlice, MakeMap, Index, SetIndex, Append, Range, Len, Cap | ~300 |
| `ops_convert.go` | Convert, TypeAssert, ChangeType, ChangeInterface | ~200 |

Each handler is a method on `*vm`: `v.executeArithmetic(op, frame)`, etc.

`run.go` (1048 lines) stays monolithic — it's the inlined hot path and splitting it would harm performance.

### 10. Remove Duplicate Operand Width Table

`bytecode/opcode.go` has both `operandWidthTable` (array) and `OperandWidths` (map). Remove the map, keep the array.

## Files Changed

| File | Action |
|---|---|
| `value/accessor.go` | Remove `closureCaller`, `methodResolver`, `SetClosureCaller`, `SetMethodResolver`, `CallMethod` |
| `value/value.go` | Add `ClosureExecutor` interface |
| `vm/vm.go` | Remove `init()`, remove `SetMethodResolver` call from `newVM()` |
| `vm/goroutine.go` | Replace globals with `GoroutineTracker` struct; remove `vmRegistry` |
| `vm/closure.go` | Add `Execute()` method implementing `ClosureExecutor` |
| `vm/ops_dispatch.go` | Refactor into router; split cases into `ops_*.go` files |
| `vm/ops_arithmetic.go` | New — arithmetic opcode handlers |
| `vm/ops_control.go` | New — control flow opcode handlers |
| `vm/ops_memory.go` | New — memory opcode handlers |
| `vm/ops_call.go` | New — call + method resolution handlers |
| `vm/ops_container.go` | New — container opcode handlers |
| `vm/ops_convert.go` | New — conversion opcode handlers |
| `importer/importer.go` | `NewImporter(reg)` — injectable registry |
| `importer/register.go` | `ExternalPackage` gets `registry` back-ref; `AddMethodDirectCall` uses instance |
| `compiler/parser/parse.go` | Pass registry to `NewImporter(reg)` |
| `runner/runner.go` | Constructor options; store initial globals; own `GoroutineTracker` |
| `bytecode/bytecode.go` | Remove `InitialGlobals` field |
| `bytecode/opcode.go` | Remove duplicate `OperandWidths` map |
| `gig.go` | Update Build() flow; remove `register` package references |
| `register/register.go` | **Delete** |

## Testing Strategy

1. **All existing tests must pass unchanged** — integration tests exercise the public API
2. **New: Concurrency test** — two programs compiled concurrently, each with methods, both called in goroutines. Validates race condition fix.
3. **New: Sandbox isolation test** — `WithRegistry(sandbox)` prevents access to globally-registered packages during type-checking
4. **Run with `-race` flag** — validates no data races in concurrent usage

## Scope Boundary (NOT Changing)

- `compiler/` internal structure — already clean with good interfaces
- `stdlib/packages/` generated wrappers — generated code stays as-is
- `run.go` hot path — stays monolithic for performance
- `importer/typeconv.go` global caches — pure caches, safe for concurrent use
- `importer.globalRegistry` existence — still needed for stdlib init() registration
