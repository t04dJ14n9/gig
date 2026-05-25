# External Call Architecture: A General Approach

## Problem Statement

When interpreted Go code calls external Go functions (stdlib, third-party), values must cross the interpreter→Go boundary. Currently this boundary is handled by a patchwork of per-function wrappers (`ErrorValue`, `FmtWrap`, `GigErrorsAs`, `SprintfExtern`, etc.). Each new function that inspects types or calls methods on interpreted values requires a new wrapper. This document proposes a general approach that eliminates per-function custom logic.

## Current Architecture

### How Values Cross the Boundary

```
Interpreted Code          │  VM                     │  External Go Function
                          │                         │
var err *myError          │  Value{KindReflect,     │
= &myError{"test"}        │    obj: reflect.Value(  │
                          │      *myError{...})}    │
                          │                         │
errors.As(err, &target)   │  args[0] = Value        │  GigErrorsAs(err, target)
                          │  args[1] = Value        │
                          │                         │
                          │  → Interface()          │  receives interface{} values
                          │  → *value.Value         │  type identity lost
```

### The Three Type Identity Problems

**Problem 1: `reflect.StructOf` types have no named identity**

Interpreted structs are created via `reflect.StructOf(fields)`. This produces anonymous struct types like `struct{msg string}`. They cannot implement Go interfaces (`error`, `fmt.Stringer`, etc.) because `reflect.StructOf` types cannot have methods.

**Problem 2: `OpAddr` creates `*value.Value` pointers**

When interpreted code takes `&localVar`, `OpAddr` creates `&frame.locals[idx]` which is `*value.Value`. When this crosses the boundary via `Interface()`, external functions receive `*value.Value` instead of the actual Go pointer.

**Problem 3: Frame slot pointers can't receive assignments**

`errors.As(err, &target)` needs to modify `target` via `reflect.Value.Set()`. But the target is a `*value.Value` frame slot, not a `**myError`. The `Value` type doesn't implement `error`, so even if we pass the frame slot, type matching fails.

### Current Per-Function Wrappers

| Wrapper | File | Purpose | Affected Functions |
|---------|------|---------|-------------------|
| `ErrorValue` | `extern.go:444` | Wraps interpreted structs as `error` via `gigStructWrapper` | `errors.Is`, `errors.Unwrap`, `context.WithCancel`, `io.Copy`, `os.IsNotExist`, etc. |
| `GigErrorsAs` | `extern.go:509` | Custom `errors.As` with type-name matching | `errors.As` only |
| `FmtWrap` | `extern.go:295` | Wraps interpreted structs for `fmt.*` variadic args | `fmt.Sprintf`, `fmt.Fprintf`, etc. |
| `SprintfExtern` | `extern.go:648` | Custom `%T` handling for format strings | `fmt.Sprintf` with `%T` |
| `interpretedInterfaceAdapter` | `interface_adapter.go` | Wraps interpreted types for `sort.Interface`, `heap.Interface` | `sort.Sort`, `heap.Init`, etc. |

**Problems with this approach:**
1. Every new function that inspects types needs a new wrapper
2. Wrappers are scattered across `extern.go`, `stdlib/packages/`, and `vm/`
3. The `gigStructWrapper` only implements 4 interfaces (`error`, `fmt.Stringer`, `fmt.Formatter`, `fmt.GoStringer`)
4. Functions that do `v.(error)` on `Interface()` results (like `errors.Join`) break silently
5. Third-party packages have no way to handle interpreted types

## Proposed General Approach

### Core Idea: Transparent Type Reconstruction at the Boundary

Instead of wrapping values per-function, reconstruct proper Go types at the boundary automatically. The key insight: **the interpreter already knows the correct Go type for every value** (stored in the SSA type system). We just need to use it when crossing the boundary.

### Architecture

```
Interpreted Code          │  VM                     │  External Go Function
                          │                         │
var err *myError          │  Value{                 │
= &myError{"test"}        │    kind: KindReflect,   │
                          │    obj: reflect.Value(  │
                          │      *myError{...}),    │
                          │    typeInfo: *TypeInfo{  │  ← NEW: type metadata
                          │      name: "myError",   │
                          │      pkg: "main",       │
                          │      prog: *Program,    │
                          │    }                    │
                          │  }                      │
                          │                         │
errors.As(err, &target)   │  → AsExternalValue()   │  receives proper Go types
                          │  → reconstructs         │  type identity preserved
                          │    *gigStructWrapper    │
                          │    or direct assignment │
```

### Component 1: Type Metadata on Value

Add a `typeInfo` field to `Value` that carries the original SSA type information:

```go
type Value struct {
    kind     Kind
    size     Size
    num      int64
    obj      any
    typeInfo *TypeInfo  // NEW: nil for primitives, set for structs/pointers
}

type TypeInfo struct {
    Name     string                // "myError"
    PkgPath  string                // "main"
    Prog     *CompiledProgram      // back-reference for method lookup
    GigTag   string                // "#main.myError" (from PkgPath)
}
```

This metadata is set when `MakeFromReflect` is called for `reflect.StructOf` types during compilation. It flows through the VM without any per-operation overhead (just an extra pointer field).

### Component 2: Unified Boundary Conversion

Replace all per-function wrappers with a single `AsExternalValue(v Value) any` function:

```go
// AsExternalValue converts an interpreter Value to a Go interface{}
// that preserves type identity for external consumption.
func AsExternalValue(v Value) any {
    // Primitives: return directly
    if v.typeInfo == nil {
        return v.Interface()
    }

    // Struct/pointer with type metadata: wrap in gigStructWrapper
    iface := v.Interface()
    if iface == nil {
        return nil
    }

    // Already a native Go type (e.g., from external package)
    if _, ok := iface.(error); ok {
        return iface
    }

    // Interpreter-defined type: create wrapper with full interface support
    return newGigStructWrapper(v, iface, v.typeInfo)
}
```

### Component 3: Enhanced gigStructWrapper

Expand `gigStructWrapper` to support arbitrary interfaces, not just 4 hardcoded ones:

```go
type gigStructWrapper struct {
    iface     any
    typeName  string
    typeInfo  *TypeInfo
    prog      *CompiledProgram

    // Lazy method resolution cache
    methodCache map[string]reflect.Value
}

// Implements checks if the wrapped type satisfies the given interface.
func (w *gigStructWrapper) Implements(ifaceType reflect.Type) bool {
    // Check compiled methods against interface methods
    for i := 0; i < ifaceType.NumMethod(); i++ {
        methodName := ifaceType.Method(i).Name
        if !w.hasMethod(methodName) {
            return false
        }
    }
    return true
}

// MethodByName returns a reflect.Value for the named method,
// dispatching to compiled interpreter methods.
func (w *gigStructWrapper) MethodByName(name string) (reflect.Value, bool) {
    if cached, ok := w.methodCache[name]; ok {
        return cached, true
    }
    // Look up in compiled program
    fn := w.prog.LookupMethod(w.typeInfo.Name, name)
    if fn == nil {
        return reflect.Value{}, false
    }
    // Create closure that dispatches to interpreter
    method := createMethodClosure(w.prog, w.iface, fn)
    w.methodCache[name] = method
    return method, true
}
```

### Component 4: Dynamic Interface Implementation via reflect.MakeFunc

For functions that need the wrapper to satisfy specific interfaces (not just `error`), use `reflect.MakeFunc` to generate method implementations dynamically:

```go
// satisfyInterface creates a new type that implements the given interface
// by dispatching method calls to the interpreter.
func satisfyInterface(w *gigStructWrapper, ifaceType reflect.Type) reflect.Value {
    methods := make([]reflect.Method, ifaceType.NumMethod())
    for i := 0; i < ifaceType.NumMethod(); i++ {
        m := ifaceType.Method(i)
        methods[i] = reflect.Method{
            Name: m.Name,
            Type: m.Type,
            Func: createMethodClosure(w.prog, w.iface, w.typeInfo.Name, m.Name),
        }
    }
    // Create a new struct type that embeds the wrapper and implements the interface
    // ...
}
```

### Component 5: Automatic Frame Slot Assignment

For functions like `errors.As` that need to write back to frame slots:

```go
// GigErrorsAs with automatic frame slot handling
func GigErrorsAs(err error, target any) bool {
    // Detect *Value frame slot pointer
    if vp, ok := target.(*Value); ok {
        return gigErrorsAsToSlot(err, vp)
    }
    // Normal Go double pointer
    return gigErrorsAsNormal(err, target)
}

func gigErrorsAsToSlot(err error, slot *Value) bool {
    slotVal := slot.Interface()
    if slotVal == nil {
        return false
    }
    slotRV := reflect.ValueOf(slotVal)
    if !slotRV.IsValid() || slotRV.Kind() != reflect.Ptr {
        return false
    }
    targetType := slotRV.Type()

    for {
        if matchAndSetSlot(err, targetType, slot) {
            return true
        }
        unwrapper, ok := err.(interface{ Unwrap() error })
        if !ok {
            return false
        }
        err = unwrapper.Unwrap()
        if err == nil {
            return false
        }
    }
}
```

## Implementation Plan

### Phase 1: Type Metadata on Value (Low Risk)

1. Add `TypeInfo` struct to `model/value/value.go`
2. Add `typeInfo` field to `Value` (zero-value for primitives, no overhead)
3. Set `typeInfo` in `MakeFromReflect` when creating values from `reflect.StructOf` types
4. Thread `TypeInfo` through `typeconv.go` → `MakeFromReflect` → `Value`

**Files**: `model/value/value.go`, `model/value/container.go`, `vm/typeconv.go`

### Phase 2: Unified AsExternalValue (Medium Risk)

1. Create `AsExternalValue(v Value) any` in `model/value/extern.go`
2. Replace `ErrorValue`, `FmtWrap`, and `SprintfExtern` with calls to `AsExternalValue`
3. Enhance `gigStructWrapper` to use `TypeInfo` for method lookup
4. Add `gigAsMatchFrameSlot` for frame slot assignment

**Files**: `model/value/extern.go`, `stdlib/packages/errors.go`, `stdlib/packages/fmt.go`

### Phase 3: Dynamic Interface Satisfaction (High Risk)

1. Implement `satisfyInterface` using `reflect.MakeFunc`
2. Replace `interpretedInterfaceAdapter` with dynamic interface satisfaction
3. Support arbitrary Go interfaces (not just `sort.Interface`, `error`)
4. Add caching for generated method implementations

**Files**: `vm/interface_adapter.go`, `model/value/extern.go`

### Phase 4: DirectCall Wrapper Simplification (Low Risk)

1. Update `cmd/gig/gentool/directcall.go` to use `AsExternalValue` for all error/interface parameters
2. Remove per-function custom overrides from `customCallOverrides`
3. Regenerate all 71 stdlib wrappers

**Files**: `cmd/gig/gentool/directcall.go`, `stdlib/packages/*.go`

### Phase 5: Third-Party Package Support (Low Risk)

1. Document the `AsExternalValue` API for third-party package authors
2. Add `RegisterInterface(path, name, ifaceType)` for custom interface registration
3. Provide `gig.StructWrapper(v Value) any` public API

**Files**: `gig.go`, `docs/cli-guide.md`

## Migration Path

### Backward Compatibility

- `ErrorValue`, `FmtWrap`, `GigErrorsAs` remain as thin wrappers around `AsExternalValue`
- Existing DirectCall wrappers continue to work (just less efficiently)
- No changes to the public API (`Build`, `Run`, `RunWithContext`)

### Incremental Adoption

Each phase is independently deployable:
- Phase 1 alone improves type debugging
- Phase 2 alone eliminates most per-function wrappers
- Phase 3 alone enables arbitrary interface satisfaction
- Phase 4 alone simplifies code generation

## Testing Strategy

### Unit Tests

- `model/value/extern_test.go`: Test `AsExternalValue` with all value kinds
- `model/value/typeinfo_test.go`: Test `TypeInfo` creation and propagation

### Integration Tests

- `tests/testdata/parity/main.go`: Add tests for `errors.As`, `errors.Join`, `sort.Sort` with interpreted types
- `tests/known_issue_test.go`: Verify all known issues are resolved

### Regression Tests

- Run full test suite after each phase
- Benchmark before/after to ensure no performance regression

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| `TypeInfo` adds 8 bytes per Value | Low (most Values are primitives with nil typeInfo) | Profile memory usage; use pointer only for struct/pointer kinds |
| `reflect.MakeFunc` is slow | Medium (called per-method per-wrapper) | Cache generated methods; use DirectCall fast path when available |
| Breaking existing wrappers | High | Keep wrappers as thin delegates; add deprecation warnings |
| Third-party packages can't handle interpreted types | Medium | Provide public `gig.StructWrapper` API |

## Decision Record

**Decision**: Implement general type reconstruction at the boundary via `TypeInfo` + `AsExternalValue`.

**Drivers**:
- Eliminate per-function wrapper whack-a-mole
- Support arbitrary Go interfaces (not just 4 hardcoded ones)
- Enable third-party package compatibility
- Reduce maintenance burden

**Alternatives Considered**:
1. **Keep current approach**: Rejected — doesn't scale, requires new wrapper per function
2. **Fix `Interface()` to unwrap `*value.Value`**: Rejected — breaks `ErrorValue` which needs `Value` metadata
3. **Use code generation for all wrappers**: Rejected — still per-function, just automated

**Consequences**:
- `Value` grows by 8 bytes (pointer to `TypeInfo`)
- `gigStructWrapper` becomes more complex but handles all interfaces
- Per-function wrappers become thin delegates (can be removed later)
- Third-party packages can use `gig.StructWrapper` for compatibility

**Follow-ups**:
- Benchmark `TypeInfo` memory overhead
- Prototype `reflect.MakeFunc` performance
- Design public API for third-party packages
