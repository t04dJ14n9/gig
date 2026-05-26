# Plan: Interface Adapter Correctness & Custom Type Safety

## Requirements Summary

The `interface_adapter.go` mechanism — introduced to bridge interpreted types into host
`sort.Interface` / `heap.Interface` calls — has three correctness bugs that break the
goal of 100% behavioral parity between interpreted and natively compiled Go code.

Additionally, user-defined custom types (structs, interfaces) in gig scripts are
fundamentally incompatible with host/third-party code that uses reflection or type
assertions, because gig represents them internally as `reflect.StructOf` types rather
than native Go types. Users need a hard opt-out.

## Acceptance Criteria

1. **Swap correctness**: A user-defined type implementing `sort.Interface` with a
   non-trivial `Swap` method (updating auxiliary state) produces correct results when
   passed to `sort.Sort`.
2. **Host callback context**: A `Less`/`Len`/`Swap` callback that reads package-level
   globals or relies on stateful mode observes the correct state.
3. **Host callback error propagation**: A callback that panics propagates the panic
   back to the host caller instead of silently swallowing it.
4. **Adapter scope**: Only exact `sort.Interface` and `heap.Interface` boundary cases
   are adapted. Arbitrary user-defined interfaces with `Len`/`Less`/`Swap` methods
   retain their concrete type identity through `OpMakeInterface`.
5. **Type identity**: After `OpMakeInterface`, type assertions and reflection on the
   resulting interface value see the original concrete type, not
   `*interpretedInterfaceAdapter`.
6. **`WithDisableCustomTypes()`**: When set, compilation fails with a clear error
   for any user-defined `type Foo struct {...}` or `type Bar interface {...}`.
7. **No regressions**: All existing tests pass (`go test -v -race ./...`).
8. **Backward compatibility**: Default behavior (no new options) preserves current
   behavior for `sort.Interface`/`heap.Interface` adapters but fixes the three bugs.

## Implementation Steps

### Phase 1: Fix adapter injection scope ([vm/ops_convert.go:392-409](vm/ops_convert.go#L392-L409))

**Problem**: `isHostCallbackInterface` triggers for ANY interface with methods named
`Len`, `Less`, `Swap` — including user-defined interfaces that merely share the same
shape. The resulting `*interpretedInterfaceAdapter` replaces the concrete type, so
type assertions, reflection, and third-party code see the adapter instead of the
original type.

**Fix**: Remove the shape-based fallback in `isHostCallbackInterface`. Only match the
exact `sort.Interface` and `heap.Interface` named types (already done on lines
393-397). Delete lines 400-408 (the `len(iface.NumMethods())>0` shape check).

**File**: [vm/ops_convert.go](vm/ops_convert.go) — `isHostCallbackInterface` function.

### Phase 2: Fix Swap bypass bug ([vm/interface_adapter.go:41-46](vm/interface_adapter.go#L41-L46))

**Problem**: `interpretedInterfaceAdapter.Swap` calls `callReceiverSliceSwap` first,
and if it succeeds (receiver is slice-backed), returns immediately without ever
calling the interpreted `Swap` method. For user-defined types whose receiver is
a named slice type with a non-trivial `Swap` (updating auxiliary indexes, counters,
parallel state), this silently skips the user's logic.

**Fix**: Invert the dispatch order — call the compiled method first via
`a.call("Swap", ...)`. Only fall back to `callReceiverSliceSwap` when no compiled
method is found.

**File**: [vm/interface_adapter.go](vm/interface_adapter.go) — `Swap` method.

### Phase 3: Fix detached VM in host callbacks ([vm/interface_adapter.go:71-116](vm/interface_adapter.go#L71-L116))

**Problem**: `callCompiledMethodValue` creates a temp VM with zeroed globals and
`context.Background()`. This severs callbacks from:
- Package-level globals (they see zero values, not post-init state)
- Stateful mode (mutations from prior `Run` calls are invisible)
- Goroutine context and cancellation deadline
- The goroutine tracker (goroutine accounting is lost)

Additionally, interpreter panics/errors are recovered and silently swallowed via
`errCompiledMethodPanic{}`, causing the callback to return bogus zero values.

**Fix**: The adapter needs access to the caller VM's execution context. Two approaches:

**Option A (preferred)**: Store the caller VM's context in the adapter at construction
time. Pass `globals`, `initialGlobals`, `goroutines`, and `ctx` from the caller VM to
the temp VM in `callCompiledMethodValue`. Propagate errors instead of swallowing them.

**Option B**: Run callbacks on the live caller VM directly by pushing a new frame
instead of creating a temp VM. More invasive but eliminates the detached-VM problem
entirely.

Given the complexity of Option B (it would require suspending/resuming the caller
frame), go with **Option A**: thread the caller's VM context through the adapter.

**Files**:
- [vm/interface_adapter.go](vm/interface_adapter.go) — `interpretedInterfaceAdapter` struct, `newInterpretedInterfaceAdapter`, `callCompiledMethodValue`
- [vm/ops_convert.go](vm/ops_convert.go) — `makeInterpretedInterfaceAdapter` call site (needs to receive caller VM)

### Phase 4: Add `WithDisableCustomTypes()` build option

#### 4a. Add option to build config ([gig.go](gig.go))

Add `disableCustomTypes bool` to `buildConfig`, add `WithDisableCustomTypes()`
BuildOption function, and thread through to `compiler.BuildOption`.

#### 4b. Add compiler option ([compiler/](compiler/))

Follow the existing pattern from `WithAllowPanic()`:
- Add `disableCustomTypes bool` to the compiler's internal config
- Add `WithDisableCustomTypes()` compiler `BuildOption`
- Thread the value through `gig.Build()` → `compiler.Build()`

#### 4c. Add compile-time check

In the compiler, when processing `*ssa.Type` instructions for named types:
- If `disableCustomTypes` is set and the type has an underlying struct or interface
  that is NOT from an external package (i.e., user-defined in the script), emit a
  clear error like: `"custom type 'Foo' is not allowed when WithDisableCustomTypes is set"`

The check should go in the type declaration handling within the compiler, likely in
`compile_instr.go` or `compile_value.go` where `*ssa.Type` is processed.

#### 4d. What counts as "custom":
- `type Foo struct { ... }` — blocked
- `type Bar interface { ... }` — blocked
- `type MyInt int` — could go either way; named primitive types are less problematic
  since they're backed by real Go primitives, but for safety, block them too
- External types (`sort.IntSlice`, `time.Time`) — always allowed

### Phase 5: Testing

#### 5a. Swap regression test
New test case in `tests/testdata/` for a user-defined sort.Interface with a Swap that
updates a counter:
```go
type CountingSlice []int
func (c CountingSlice) Len() int           { return len(c) }
func (c CountingSlice) Less(i, j int) bool { return c[i] < c[j] }
func (c CountingSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i]; swapCount++ }
```

#### 5b. Globals/context regression test
New test case where a sort callback reads a package-level variable that was mutated
before the sort call:
```go
var multiplier int = 2
type ByScaledValue []int
func (b ByScaledValue) Less(i, j int) bool { return b[i]*multiplier < b[j]*multiplier }
```

#### 5c. Stateful mode callback test
Verify that a callback sees globals mutated by a prior `Run()` call when stateful
mode is enabled.

#### 5d. Adapter scope regression test
Define an interface with `Len`/`Less`/`Swap` that is NOT `sort.Interface`, make an
interface value, and verify type assertion to the concrete type works:
```go
type MyIface interface { Len() int; Less(i, j int) bool; Swap(i, j int) }
type MyType struct { data []int }
// ... implement MyIface ...
var i MyIface = MyType{...}
_, ok := i.(MyType) // must be true after fix
```

#### 5e. `WithDisableCustomTypes` tests
- `gig.Build("type Foo struct { X int }", gig.WithDisableCustomTypes())` → error
- `gig.Build("type Bar interface { M() }", gig.WithDisableCustomTypes())` → error
- `gig.Build("func Add(a, b int) int { return a+b }", gig.WithDisableCustomTypes())` → OK

#### 5f. Panic propagation test
Callback method that panics should not return a bogus zero value:
```go
type PanicSlice []int
func (p PanicSlice) Len() int { panic("boom") }
// sort.Sort(PanicSlice{1,2,3}) should propagate the panic
```

### Phase 6: Verification

```bash
go test -v -race ./...
go test -v -run TestKnownIssues ./tests/
go test -v -run TestDivergenceHunt ./tests/
```

## Risks and Mitigations

| Risk | Mitigation |
|------|-----------|
| Removing shape-based adapter matching breaks existing sort/heap usage | The exact-name check for `sort.Interface`/`heap.Interface` already covers the stdlib use cases; the shape check was an over-broadening that introduced the type-identity bug |
| Threading VM context through adapter is invasive | Option A (store context in adapter) is minimally invasive — only the adapter struct and construction call site change |
| `WithDisableCustomTypes` may block legitimate use cases | The option is opt-in; default behavior is unchanged |
| Performance regression from always calling compiled Swap | The compiled method lookup is cheap (map lookup in MethodsByName). For stdlib types like `sort.IntSlice` where the compiled method is the trivial swap, the fast path can be preserved by checking if the method is from an external package |

## ADR

### Decision
Narrow the interface adapter to exact `sort.Interface`/`heap.Interface` boundary cases,
fix the Swap dispatch order and detached-VM bugs, and add an opt-in
`WithDisableCustomTypes()` build option.

### Drivers
1. 100% behavioral correctness for interpreted code vs compiled Go
2. Type identity preservation for user-defined types across interpreter/host boundary
3. Safe defaults — users shouldn't hit reflection bugs by accident

### Alternatives considered
- **Do nothing**: Unacceptable — the three bugs cause silent data corruption and
  hard-to-diagnose failures
- **Remove the adapter entirely**: Would break `sort.Sort`/`heap.Push` for interpreted
  types, a major regression
- **Add `DisableCustomTypes` only**: Doesn't fix the adapter bugs that affect users
  who DON'T disable custom types
- **Build a full Go type system in gig**: Months of work, not practical now

### Why chosen
The combination of fixing adapter bugs + opt-in disable option provides:
- Correctness for existing users (adapter fixes)
- Safety for users who want guarantees (disable option)
- No breaking changes to the default API

### Consequences
- Shape-based adapter matching is removed; only exact `sort.Interface`/`heap.Interface` match
- Users with custom `Len`/`Less`/`Swap` interfaces that are NOT `sort.Interface` will
  no longer get automatic adapter injection (which is correct — they shouldn't have been)
- `WithDisableCustomTypes()` users cannot define custom types; they must use
  primitives, slices, maps, and external types only

### Follow-ups
- Consider a "type registry" that lets users pre-register custom types so host code
  sees native Go types instead of internal representations
- Explore whether `reflect.TypeOf` can be made to return the "right" type for
  interpreted values (long-term architectural question)
