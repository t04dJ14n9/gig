# Gig v2 Clean SSA Interpreter Implementation Plan

**Goal:** Rebuild Gig as a clean, lightweight Go interpreter using direct SSA execution, while keeping `RunWithContext`/`Run` as compatibility wrappers.

**Architecture:** Final pipeline is `Go source -> frontend parse/typecheck -> SSA -> direct SSA interpreter`. There is one backend only: no bytecode compiler, opcodes, stack VM, local slots, function indices, VM pools, or int-specialized bytecode optimizations.

**Tech Stack:** Go, `golang.org/x/tools/go/ssa`, `go/parser`, `go/types`, reflection at host boundaries, internal tagged-union `value.Value`.

## Status

| Phase | Steps                 | State        | Notes                                              |
|-------|-----------------------|--------------|----------------------------------------------------|
| 1     | 1, 2, 3 (skeleton)    | **complete** | Branch + skeleton packages, all interfaces.        |
| 2     | 5 (value system)      | **complete** | Tagged-union Value + Converter + 26 unit tests.    |
| 2     | 4 (frontend)          | **complete** | Builder: parse / validate / ssautil + 13 tests.    |
| 3     | 6 (interp slice 1)    | **complete** | Scalar arith + control flow + calls + 26 tests.    |
| 3     | 6 (interp slice 2)    | not started  | Composites, builtins, host calls.                  |
| 3     | 6 (interp slice 3)    | not started  | Closures, interfaces, defer/panic/recover, conc.   |
| 4     | 7 (compat wrappers)   | not started  | gig.Program facade in terms of interp.Program.     |
| 5     | 8 (delete legacy)     | not started  | vm/, runner/, model/bytecode/, optimizer.          |

### Phase 1 deliverables (session 1)

- Branch `feature/clean-ssa-gig` created on top of current main.
- Skeleton packages added at `host/`, `value/`, `diag/`, `internal/frontend/`, `internal/interp/`.
  - `host` declares `Environment`, `Function`, `Variable`, `Constant`, `Type`, `Method`, `InterfaceProxy`, `Import`.
  - `value` declares `Value`, `Kind`, `Size`, `Converter`, `TypeResolver`.
  - `diag` declares `Diagnostic`, `Severity`.
  - `internal/frontend` declares `Builder`, `Unit`, `Source`, `Config`, `PanicPolicy`.
  - `internal/interp` declares `Engine`, `Program`, `Cell`, `Config`.
- Plan step "stdlib" deferred: `stdlib/` already exists at top level for legacy DirectCall wrappers; no new package was created to avoid collision. Decide in Phase 2 whether to nest the new stdlib under `host/stdlib/` or repurpose the existing dir.
- `go build ./...` and `go vet ./...` pass with the skeleton in place.
- All existing tests still pass; legacy code is untouched.

The skeleton is intentionally interface-only at this phase: there are no
constructors, no instance state, and no behaviour. Concrete bridges from
the existing `importer/` to `host.Environment` and from `model/value/`
to `value.Value` land in Phase 2.

### Phase 2 deliverables (session 2 — Step 5 only)

- `value/value.go` lifted from `model/value/value.go` and scrubbed of
  VM-isms. The 32-byte tagged-union design is preserved (kind:1 +
  size:1 + num:8 + obj:16); zero-allocation primitives still hold.
- Removed legacy concepts that pulled the VM into the value package:
  `KindBytes`, `MakeBytes`, `Bytes()`, `MakeIntSlice`/`IntSlice`,
  `MakeIntPtr`/`IntPtr`, `MakeValueSlice`/`ValueSlice`, `MakeFunc`,
  `RawObj`/`RawInt`/`RawBool`, `MethodResolver*`, `ClosureExecutor`,
  `InterpretedInterfaceValue`/`MakeInterpretedInterface`.
- Added `value.Converter` (interface) and `DefaultConverter()`
  (stateless implementation):
  - `FromAny(any) (Value, error)` — type-switch fast path for scalars,
    reflect fallback for everything else.
  - `FromReflect(reflect.Value) (Value, error)` — symmetric, no extra
    `reflect.ValueOf` hop.
  - `ToAny(Value) (any, error)` — preserves original Go type.
  - `ToReflect(Value, reflect.Type) (reflect.Value, error)` — converts
    to caller-specified target type.
  - `Zero(types.Type, TypeResolver) (Value, error)` — basic types
    handled inline; composite types route through `TypeResolver`.
  - `Convert(Value, types.Type, TypeResolver) (Value, error)` — Go
    `T(x)` semantics, including narrowing and float→int truncation.
- 26 unit tests in `value/value_test.go` covering primitive
  round-trips, kind-mismatch panics, `IsValid`/`IsNil`, size
  preservation, reflect handling, basic-type `Zero`, narrowing
  `Convert`, and float bit-pattern preservation.
- Legacy `model/value/` is **untouched**. The new `value/` package
  ships alongside it; the cutover happens in Phase 5 (Step 8) when the
  legacy backend is deleted.

### Phase 2 deliverables (session 2 — Step 4 frontend)

- `internal/frontend/builder.go` implements `frontend.Builder` and
  `frontend.Unit` using the same pipeline gofun and legacy gig use:
  parse → banned-import check → panic-policy check → auto-import
  splice → `ssautil.BuildPackage`. No bytecode emission.
- The implementation is ~250 lines (vs. the legacy `compiler/parser/` +
  `compiler/ssa/` which together are ~270 lines of similar logic but
  spread across two packages and tied to the legacy registry).
- Banned-import default list (`unsafe`, `reflect`) preserved, but now
  configurable via `Config.BannedImports` (legacy hardcoded it).
- Panic policy translates to the same compile-time AST scan.
- Type-check errors are now collected into `[]diag.Diagnostic` rather
  than collapsed into a single `error` string. An aggregate `*Errors`
  type renders them all when the build fails.
- `Config.AutoImport` toggles the gofun-style identifier scan that
  splices imports in for unresolved selector expressions; off by default
  to keep the path easy to reason about.
- 13 unit tests covering: simple program builds, package-main wrap,
  explicit package decl honoured, unsafe/reflect rejection, custom
  banned list, panic policy on/off, type-check error reporting, nil
  environment rejection, context cancellation, AutoImport on/off.
- Tests use a stub `host.Environment` that imports nothing — enough to
  exercise everything that touches only the universe block. The
  concrete `host.Environment` constructor lands in Phase 3 alongside
  the interpreter.

### Phase 3 deliverables (session 3 — Step 6 vertical slice 1)

The first slice of the SSA interpreter is now running interpreted Go
programs end-to-end. Source flows through `frontend.Builder` and the
SSA IR is walked directly by `interp.Engine`, no bytecode in sight.

**Files added under `internal/interp/`:**

- `engine.go` (~170 lines) — `defaultEngine`, `program`, `typeResolver`.
  Allocates per-package globals at construction, runs `init()` once,
  resolves basic types via reflect.
- `frame.go` (~120 lines) — `callSSA`, `runFrame`, `runBlockPhis`. The
  Phi-at-block-entry semantics is preserved: all Phis evaluate from a
  pre-snapshot of the cell map, then commit together.
- `ops.go` (~190 lines) — instruction dispatcher (`visitInstr`) and
  per-instruction runners. Pattern-matches on `ssa.Instruction` types
  and routes to small handlers, mirroring gofun's `visitInstr`.
- `arith.go` (~180 lines) — `evalBinOp`, `evalUnOp`, `evalEquality`.
  All scalar operations (int / uint / float / bool / string), including
  shifts, bitwise ops, and comparisons. Returns are routed through
  `value.Converter.Convert` so result types stay correct (e.g. `int8`
  arithmetic preserves the `int8` size tag).

**SSA opcodes covered in this slice:**

| Op                | Status                                  |
|-------------------|-----------------------------------------|
| `*ssa.Const`      | done                                    |
| `*ssa.Parameter`  | done                                    |
| `*ssa.Return`     | done                                    |
| `*ssa.BinOp`      | done                                    |
| `*ssa.UnOp`       | done (incl. MUL deref)                  |
| `*ssa.Convert`    | done                                    |
| `*ssa.ChangeType` | done                                    |
| `*ssa.If`         | done                                    |
| `*ssa.Jump`       | done                                    |
| `*ssa.Phi`        | done (block-entry resolve)              |
| `*ssa.Call`       | done (interpreted-only, single return)  |
| `*ssa.Alloc`      | done (heap + local)                     |
| `*ssa.Store`      | done (locals + globals)                 |
| `*ssa.DebugRef`   | no-op skip                              |

**Programs that now run on the new interpreter:**

- Scalar arithmetic (`Add(a,b int) int`)
- Comparisons returning bool
- Unary negation, bitwise complement
- `if/else` with multiple return paths
- `for` loops with index variable mutation
- Nested `for` loops (proves Phi correctness across multiple back-edges)
- Recursive `Fib`, `Fact` (proves call-stack discipline + arg binding)
- Mutual function call chains
- Type narrowing via `int8(x)` truncation (200 → -56)
- Divide-by-zero produces a structured error
- Max call depth backstop

**26 unit tests in `internal/interp/interp_test.go`**, all green.
Total runtime under 200 ms.

**Out of slice 1 (deferred to slice 2 / 3):** composites
(slice/map/struct/array), built-ins (`len`, `cap`, `append`, `print`),
closures (`MakeClosure`/free-vars), interfaces (`MakeInterface`,
`TypeAssert`), `defer`/`panic`/`recover`, goroutines, channels,
`select`, multi-return calls + `Extract`, host function calls. Each is
a self-contained increment on the working slice.

---

## Summary

- Start fresh in a new worktree/branch; do not implement in the current dirty `main`.
- Write a new SSA interpreter core from scratch, using gofun as a design reference, not as code to fork.
- Keep domain boundaries: packages interact only through interfaces.
- Keep current Gig user ergonomics where useful: `RunWithContext` remains compatible.
- Remove public `Wait`, `Package`, and `Close` from the core `Program` API unless a future feature proves they are necessary.

## Public API

### `gig`

```go
package gig

type Compiler interface {
    Compile(ctx context.Context, source string) (Program, error)
}

type Program interface {
    // Clean core API.
    Call(ctx context.Context, name string, args ...any) ([]any, error)

    // Compatibility API.
    RunWithContext(ctx context.Context, name string, args ...any) (any, error)
    Run(name string, args ...any) (any, error)
}

func NewCompiler(opts ...Option) Compiler
func Compile(ctx context.Context, source string, opts ...Option) (Program, error)

func WithEnvironment(env host.Environment) Option
func WithPanic(policy PanicPolicy) Option
```

Note: there is no `WithGlobals` option. The interpreter has one global
state model — globals are allocated once at compile time, `init()` runs
once, and every `Call` observes/mutates the same globals. This matches
real Go semantics and gofun's behaviour. Callers that want isolation
between requests compile a fresh `Program` per request.

Compatibility behavior:

```go
func (p *program) RunWithContext(ctx context.Context, name string, args ...any) (any, error) {
    results, err := p.Call(ctx, name, args...)
    if err != nil {
        return nil, err
    }
    switch len(results) {
    case 0:
        return nil, nil
    case 1:
        return results[0], nil
    default:
        return results, nil
    }
}

func (p *program) Run(name string, args ...any) (any, error) {
    ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
    defer cancel()
    return p.RunWithContext(ctx, name, args...)
}
```

Do not expose these on core `Program`:

```go
Wait(ctx context.Context) error
Package() *ssa.Package
Close() error
```

SSA inspection can be added later under a separate debug-only API if needed.

## Domain Packages

### `host`

Explicit external Go environment. Replaces global importer/registry.

```go
type Environment interface {
    types.Importer

    AutoImport(name string) (Import, bool)

    LookupFunc(pkgPath, name string) (Function, bool)
    LookupVar(pkgPath, name string) (Variable, bool)
    LookupConst(pkgPath, name string) (Constant, bool)
    LookupType(pkgPath, name string) (Type, bool)
    LookupReflectType(t types.Type) (reflect.Type, bool)
    LookupMethod(typeName, methodName string) (Method, bool)
    LookupInterfaceProxy(iface *types.Interface) (InterfaceProxy, bool)
}

func NewEnvironment() Environment
func StandardEnvironment() Environment
```

Host function/method calls use `[]value.Value` and return `[]value.Value`; DirectCall remains hidden behind the interface.

### `value`

Keep concrete tagged-union runtime values.

```go
type Value struct {
    kind Kind
    size Size
    num  int64
    obj  any
}

type Converter interface {
    FromAny(any) (Value, error)
    FromReflect(reflect.Value) (Value, error)
    ToAny(Value) (any, error)
    ToReflect(Value, reflect.Type) (reflect.Value, error)
    Zero(types.Type, TypeResolver) (Value, error)
    Convert(Value, types.Type, TypeResolver) (Value, error)
}
```

Mutability belongs to interpreter `Cell`, not to `value.Value`.

### `internal/frontend`

Source to SSA only.

```go
type Builder interface {
    Build(ctx context.Context, src Source, env host.Environment, cfg Config) (Unit, error)
}

type Unit interface {
    Package() *ssa.Package
    FileSet() *token.FileSet
    Diagnostics() []diag.Diagnostic
}
```

### `internal/interp`

Direct SSA interpreter.

```go
type Engine interface {
    NewProgram(ctx context.Context, unit frontend.Unit, env host.Environment, cfg Config) (Program, error)
}

type Program interface {
    Call(ctx context.Context, name string, args []value.Value) ([]value.Value, error)
}

type Cell struct {
    Name  string
    Type  types.Type
    Value value.Value
}
```

Internal frame model:

```go
type frame struct {
    fn        *ssa.Function
    block     *ssa.BasicBlock
    prevBlock *ssa.BasicBlock
    cells     map[ssa.Value]*Cell
    freeVars  []*Cell
}
```

## Implementation Changes

1. **Create clean branch**
   - New worktree: `feature/clean-ssa-gig`.
   - Run `go test ./...` and record current baseline.

2. **Add domain skeleton**
   - Add `host`, `value`, `diag`, `stdlib`, `internal/frontend`, `internal/interp`.
   - Add interface assertions and dependency-direction checks where practical.

3. **Build explicit host environment**
   - Refactor importer/registry concepts into `host.Environment`.
   - Remove process-global registry as the primary API.
   - Preserve funcs, vars, constants, named types, methods, DirectCall wrappers, and interface proxies.

4. **Replace compiler with SSA frontend**
   - Keep parse/typecheck/validation and SSA construction.
   - Remove SSA-to-bytecode compilation.
   - Preserve auto package insertion, auto-import, banned imports, and panic policy.

5. **Refactor value system**
   - Keep tagged union.
   - Remove VM/bytecode assumptions.
   - Add conversion APIs for `any`, `reflect.Value`, typed nil, zero values, and type conversions.

6. **Implement direct SSA interpreter**
   - Use `map[ssa.Value]*Cell`.
   - Process Phi nodes at block entry with temporary values, then assign together.
   - Implement core SSA ops, data structures, calls, closures, methods, host interop, globals/init, defer/panic/recover, goroutines/channels/select.
   - `Call` returns `[]value.Value`; `gig.Program.Call` converts to `[]any`.

7. **Add compatibility wrappers**
   - `RunWithContext` calls `Call` and unwraps results to existing Gig shape.
   - `Run` uses `DefaultTimeout` and delegates to `RunWithContext`.

8. **Delete legacy stack VM**
   - Remove `vm/`, `runner/`, `model/bytecode/`, bytecode compiler files, optimizer files, opcode tests, stack-frame tests, VM-pool tests, and int-specialization docs.
   - Rewrite `cmd/gig dump` to print SSA/interpreter metadata only.
   - Update docs to describe `source -> SSA -> interpreter`.

## Test Plan

Required final commands:

```bash
go test ./...
git diff --check
```

Required behavior coverage:

- `Call` returns `[]any`;
- `RunWithContext` preserves old unwrap behavior;
- `Run` preserves default-timeout behavior;
- explicit host environments and sandbox environments;
- zero, one, and multiple returns;
- arithmetic, branches, loops, Phi;
- pointers, load/store, addressable mutation;
- structs, methods, method values, method expressions;
- arrays, slices, maps, range;
- closures and captured mutation;
- external functions, vars, constants, types, methods;
- external callback into interpreted closure;
- interface proxy and rejected unsafe host boundary;
- init globals run once; globals shared across Call invocations;
- defer, panic, recover;
- goroutines, channels, select, timeout.

Final cleanup search:

```bash
rg -n 'bytecode|opcode|stack VM|stack-based|CompiledProgram|CompiledFunction|FuncByIndex|SymbolTable|Op[A-Z]' --glob '*.go' --glob '*.md'
```

Remaining matches must be deleted or explicitly marked as historical archived docs.

## Assumptions

- Clean design wins over backward compatibility, except `RunWithContext` and `Run` are intentionally kept as migration-friendly wrappers.
- There is only one backend: direct SSA interpreter.
- No public `Wait`, `Package`, or `Close` on core `Program`.
- No global mutable package registry in the final design.
- `value.Value` remains a concrete tagged union.
- `Cell` is the mutable storage unit.
- Domains interact only through declared interfaces.
- gofun is a reference for execution flow, not code to fork directly.
