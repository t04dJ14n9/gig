# Known Issues Bugfix Report – Batch 3 (Final)

**Date**: 2026-03-17  
**Author**: AI-assisted fix  
**Branch**: `feature/dev_youngjin`  
**Status**: All 4 bugs fixed, tests migrated, full regression suite passes  
**Milestone**: 🎉 **Zero known issues remaining** — all 20 issues across 3 batches resolved

---

## Summary

Resolved the final 4 known issues tracked in `tests/known_issues_test.go`, bringing the total known issue count to **zero**. These bugs all involved type assertions and closure type wrapping — two fundamental operations in the interpreter's type system.

| # | Bug | Root Cause | Fix Location | Status |
|---|-----|-----------|-------------|--------|
| 17 | PointerToInterface | `compileTypeAssert` didn't extract value from tuple for non-comma-ok assertions | `compiler/compile_value.go` | ✅ Fixed |
| 18 | StructWithPointerToInterface | Same root cause as #17 | `compiler/compile_value.go` | ✅ Fixed |
| 19 | StructWithNestedFunc | `closureCaller` couldn't wrap nested closure return values | `value/accessor.go`, `vm/vm.go` | ✅ Fixed |
| 20 | StructWithInterfaceMap | Same root cause as #17 | `compiler/compile_value.go` | ✅ Fixed |

### Cross-Batch Overview

| Batch | Date | Issues | Key Areas |
|-------|------|--------|-----------|
| Batch 1 | 2026-03-16 | #1–#6 (6 bugs) | Closure→function wrapping, type switch, slice append, self-referencing structs |
| Batch 2 | 2026-03-16 | #7–#16 (10 bugs) | Defer ordering, pointer aliasing, func slices, anonymous fields, map range |
| Batch 3 | 2026-03-17 | #17–#20 (4 bugs) | Type assertions, nested closures |
| **Total** | | **20 bugs** | **All resolved** ✅ |

---

## Bug 17: PointerToInterface — Non-Comma-Ok Type Assertion

### Symptom
```
pointer to interface: got []value.Value, want 42
```
Dereferencing a pointer to interface and type-asserting (`(*p).(int)`) returned the raw `[result, ok]` tuple instead of the extracted `int` value.

### Test Case
```go
func PointerToInterface() int {
    var i interface{} = 42
    p := &i
    return (*p).(int)  // Expected: 42
}
```

### Root Cause

The SSA IR for `(*p).(int)` generates a `TypeAssert` instruction with `CommaOk = false` (non-comma-ok variant). In Go semantics:
- **Comma-ok**: `val, ok := x.(T)` — returns `(T, bool)` tuple
- **Non-comma-ok**: `val := x.(T)` — returns just `T` (panics on failure)

The `OpAssert` opcode in the VM **always** pushes a `[result, ok]` tuple onto the stack (a `[]value.Value` with 2 elements). The compiler's `compileTypeAssert` function was supposed to handle this difference, but it unconditionally stored the raw tuple:

```go
// Before (broken):
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
    typeIdx := c.addType(i.AssertedType)
    c.compileValue(i.X)
    c.emit(bytecode.OpAssert, uint16(typeIdx))
    // ❌ Stores the [result, ok] tuple as-is, even for non-comma-ok
    resultIdx := c.symbolTable.AllocLocal(i)
    c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
```

When the stored `[]value.Value` was later used as an `int`, it caused a type mismatch.

### Fix

Added a `CommaOk` check in `compileTypeAssert`. For non-comma-ok assertions, emit `OpConst(0) + OpIndex` after `OpAssert` to extract element #0 (the value) from the tuple:

```go
// After (fixed):
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
    typeIdx := c.addType(i.AssertedType)
    c.compileValue(i.X)
    c.emit(bytecode.OpAssert, uint16(typeIdx))

    if !i.CommaOk {
        // Non-comma-ok: extract just the value from the [result, ok] tuple
        c.emit(bytecode.OpConst, uint16(c.addConstant(0)))
        c.emit(bytecode.OpIndex)
    }

    resultIdx := c.symbolTable.AllocLocal(i)
    c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
```

**Bytecode before fix** (for `(*p).(int)` non-comma-ok):
```
OpAssert <typeIdx>     ; pushes [result, ok] tuple
OpSetLocal <idx>       ; stores the WHOLE tuple ❌
```

**Bytecode after fix**:
```
OpAssert <typeIdx>     ; pushes [result, ok] tuple
OpConst 0              ; pushes index 0
OpIndex                ; extracts tuple[0] = result value
OpSetLocal <idx>       ; stores just the value ✅
```

### Impact
This single fix resolves **three** bugs (17, 18, 20), since all three involve non-comma-ok type assertions in different contexts.

---

## Bug 18: StructWithPointerToInterface

### Symptom
Same as Bug 17 — type assertion on a dereferenced `*interface{}` inside a struct returned the raw tuple.

### Test Case
```go
type PtrToInterface struct {
    data *interface{}
}

func StructWithPointerToInterface() int {
    var i interface{} = 42
    s := PtrToInterface{data: &i}
    return (*s.data).(int)  // Expected: 42
}
```

### Root Cause
Identical to Bug 17. The expression `(*s.data).(int)` generates the same non-comma-ok `TypeAssert` SSA instruction. The fix in `compileTypeAssert` resolves this case automatically.

### Fix
Same fix as Bug 17 — no additional changes needed.

---

## Bug 19: StructWithNestedFunc — Nested Closure Return Type

### Symptom
```
panic: reflect: function created by MakeFunc using closure has wrong type:
  have func() *vm.Closure
  want func() func() int
```
Calling a struct's function field that returns another function panicked because the inner closure was not properly wrapped.

### Test Case
```go
type NestedFuncHolder struct {
    get func() func() int
}

func StructWithNestedFunc() int {
    h := NestedFuncHolder{
        get: func() func() int {
            return func() int { return 42 }
        },
    }
    return h.get()()  // Expected: 42
}
```

### Root Cause

The execution flow reveals the issue:

1. `h.get` is stored as a `*vm.Closure` in the struct field of type `func() func() int`
2. When `h.get` is read from the struct, `ToReflectValue` wraps it via `reflect.MakeFunc` ✅
3. `h.get()` is called → the `closureCaller` callback executes the outer closure in a child VM
4. The outer closure returns a `*vm.Closure` (the inner `func() int`)
5. `closureCaller` converts this to `reflect.ValueOf(*vm.Closure)` → type `*vm.Closure` ❌
6. `reflect.MakeFunc` expects the return type to be `func() int`, not `*vm.Closure` → **panic**

The `closureCaller` function had no knowledge of the expected output types, so it couldn't know to wrap the returned `*vm.Closure` as a `func() int`.

### Fix (Two Parts)

**Part 1: `value/accessor.go`** — Extended `ClosureCaller` signature to accept expected output types:

```go
// Before:
type ClosureCaller func(closure any, args []reflect.Value) []reflect.Value

// After:
type ClosureCaller func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value
```

In the `ToReflectValue` MakeFunc callback, the output types are already computed. Pass them through to `closureCaller`:

```go
fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
    results := closureCaller(closure, args, outTypes)  // ← pass outTypes
    // ... result conversion ...
})
```

**Part 2: `vm/vm.go`** — Updated the `closureCaller` implementation to use `outTypes` for recursive wrapping:

```go
value.RegisterClosureCaller(func(closure any, args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
    // ... execute closure in child VM ...
    result, _ := closureVM.run()

    if result.Kind() == value.KindNil {
        return []reflect.Value{}
    }

    // NEW: When output types are available, use ToReflectValue for recursive wrapping.
    // This handles nested closures: the inner *vm.Closure gets wrapped via
    // reflect.MakeFunc into a proper func() int.
    if len(outTypes) > 0 {
        return []reflect.Value{result.ToReflectValue(outTypes[0])}
    }

    // Fallback: direct reflect.ValueOf conversion
    iface := result.Interface()
    if iface == nil {
        return []reflect.Value{}
    }
    return []reflect.Value{reflect.ValueOf(iface)}
})
```

**The recursive wrapping chain**:
```
h.get stored as *Closure → ToReflectValue wraps as func() func() int
  ↓ h.get() called via reflect.MakeFunc
  closureCaller runs outer closure → returns *Closure (inner func)
  ↓ outTypes[0] = func() int
  result.ToReflectValue(func() int) → wraps inner *Closure via reflect.MakeFunc
  ↓ h.get()() called
  closureCaller runs inner closure → returns 42
  ↓ outTypes[0] = int
  result.ToReflectValue(int) → returns reflect.ValueOf(42)
```

### Design Decision: outTypes vs. Post-hoc Conversion

Two approaches were considered:

1. **Post-hoc conversion** in MakeFunc callback: detect when `results[i]` is a `*vm.Closure` and wrap it. This requires recognizing `*vm.Closure` in the `value` package, which introduces a dependency on `vm` types.

2. **Pass outTypes to closureCaller** (chosen): the closure caller already has access to `value.ToReflectValue`, and knows the result type. This approach is cleaner because it leverages existing infrastructure and handles arbitrary nesting depth through recursion.

---

## Bug 20: StructWithInterfaceMap

### Symptom
Same as Bug 17 — type assertion on a value retrieved from `map[string]interface{}` inside a struct returned the raw tuple.

### Test Case
```go
type InterfaceMapHolder struct {
    data map[string]interface{}
}

func StructWithInterfaceMap() int {
    h := InterfaceMapHolder{
        data: map[string]interface{}{
            "a": 1,
            "b": "hello",
        },
    }
    return h.data["a"].(int)  // Expected: 1
}
```

### Root Cause
Identical to Bug 17. The expression `h.data["a"].(int)` generates a non-comma-ok `TypeAssert`.

### Fix
Same fix as Bug 17 — no additional changes needed.

---

## Files Modified

| File | Changes |
|------|---------|
| `compiler/compile_value.go` | Added `i.CommaOk` check in `compileTypeAssert`: emit `OpConst(0)+OpIndex` for non-comma-ok assertions |
| `value/accessor.go` | Extended `ClosureCaller` signature to accept `outTypes []reflect.Type`; pass `outTypes` in MakeFunc callback |
| `vm/vm.go` | Updated `closureCaller` to accept `outTypes`; use `result.ToReflectValue(outTypes[0])` for recursive closure wrapping |
| `tests/testdata/resolved_issue/main.go` | Added `PointerToInterface` test function + comments for issues 18-20 |
| `tests/resolved_issue_test.go` | Added 4 test cases (issues 17-20); tests 18-20 use isolated inline source |
| `tests/testdata/known_issues/main.go` | Cleared — all issues resolved |
| `tests/known_issues_test.go` | Cleared — `TestKnownIssues` now skips with "No known issues remaining" |

---

## Verification

### Test Results

All 4 new resolved issue tests pass:
```
$ go test ./tests/ -run "TestResolved_PointerToInterface|TestResolved_StructWith" -v
=== RUN   TestResolved_PointerToInterface
--- PASS: TestResolved_PointerToInterface (0.00s)
=== RUN   TestResolved_StructWithPointerToInterface
--- PASS: TestResolved_StructWithPointerToInterface (0.00s)
=== RUN   TestResolved_StructWithNestedFunc
--- PASS: TestResolved_StructWithNestedFunc (0.00s)
=== RUN   TestResolved_StructWithInterfaceMap
--- PASS: TestResolved_StructWithInterfaceMap (0.00s)
PASS
```

Full test suite passes with zero failures:
```
$ go test ./...
ok   git.woa.com/youngjin/gig              0.014s
ok   git.woa.com/youngjin/gig/bytecode     (cached)
ok   git.woa.com/youngjin/gig/compiler     (cached)
ok   git.woa.com/youngjin/gig/importer     0.003s
ok   git.woa.com/youngjin/gig/tests       40.856s
ok   git.woa.com/youngjin/gig/value        (cached)
ok   git.woa.com/youngjin/gig/vm           (cached)
```

### Test Architecture Note

Tests for issues 18, 19, and 20 use **isolated inline source** rather than the shared `testdata/resolved_issue/main.go` file. This is because these tests define package-level types (`PtrToInterface`, `NestedFuncHolder`, `InterfaceMapHolder`) that would conflict with types from other tests when `reflect.StructOf` is used — the Go reflect package caches struct types globally and may panic on duplicate definitions with different layouts.

---

## Cumulative Statistics

### All 20 Resolved Issues

| # | Bug Name | Category | Batch |
|---|----------|----------|-------|
| 1 | MapWithFuncValue | Closure wrapping | 1 |
| 2 | InterfaceSliceTypeSwitch | Type assertion | 1 |
| 3 | StructWithFuncField | Closure wrapping | 1 |
| 4 | SliceFlatten | Slice operations | 1 |
| 5 | MapUpdateDuringRange | Map semantics | 1 |
| 6 | StructSelfRef | Type conversion | 1 |
| 7 | ClosureCapture | Closure semantics | 2 |
| 8 | NilSliceAppend | Slice operations | 2 |
| 9 | ChannelDirections | Channel semantics | 2 |
| 10 | StringConversion | Type conversion | 2 |
| 11 | DeferInClosureWithArg | Defer semantics | 2 |
| 12 | PointerSwapInStruct | Pointer aliasing | 2 |
| 13 | StructWithFuncSlice | Closure wrapping | 2 |
| 14 | StructAnonymousField | Struct reflection | 2 |
| 15 | StructEmbeddedInterface | Struct semantics | 2 |
| 16 | MapRangeWithBreak | Map semantics | 2 |
| 17 | PointerToInterface | Type assertion | 3 |
| 18 | StructWithPointerToInterface | Type assertion | 3 |
| 19 | StructWithNestedFunc | Closure wrapping | 3 |
| 20 | StructWithInterfaceMap | Type assertion | 3 |

### Bug Categories

| Category | Count | Issues |
|----------|-------|--------|
| **Closure wrapping** | 4 | #1, #3, #13, #19 |
| **Type assertion** | 4 | #2, #17, #18, #20 |
| **Slice operations** | 2 | #4, #8 |
| **Map semantics** | 2 | #5, #16 |
| **Type conversion** | 2 | #6, #10 |
| **Pointer aliasing** | 1 | #12 |
| **Defer semantics** | 1 | #11 |
| **Closure semantics** | 1 | #7 |
| **Channel semantics** | 1 | #9 |
| **Struct reflection** | 1 | #14 |
| **Struct semantics** | 1 | #15 |

### Key Components Modified (Across All 3 Batches)

| Component | Files Modified | Purpose |
|-----------|---------------|---------|
| **Compiler** | `compiler/compile_value.go` | Type assertion codegen, defer ordering |
| **VM Dispatch** | `vm/ops_dispatch.go` | OpAssert, OpAppend, OpDeref, OpFieldAddr, OpNew, OpClosure |
| **VM Core** | `vm/vm.go`, `vm/run.go` | ClosureCaller registration, OpCallIndirect |
| **Type Conv** | `vm/typeconv.go` | Self-referencing structs, anonymous fields |
| **Value System** | `value/accessor.go`, `value/container.go` | ClosureCaller, ToReflectValue, SetElem, SetIndex |
| **Closure** | `vm/closure.go` | Program field for child VM creation |

### Test Coverage

- **20 resolved issue tests** in `tests/resolved_issue_test.go`
- **600+ tricky tests** in `tests/tricky_test.go`
- **0 remaining known issues** in `tests/known_issues_test.go`
- Full suite: **all packages pass** with zero failures
