# Gig Interpreter — Known Issues Bug Fix Report

**Date**: 2026-03-16  
**Author**: AI-assisted fix  
**Branch**: `feature/dev_youngjin`  
**Status**: All 6 bugs fixed, tests migrated, full regression suite passes

---

## Summary

This report documents the investigation and resolution of 6 known bugs in the Gig Go interpreter where interpreted execution results differed from native Go execution. All bugs were tracked in `tests/known_issues_test.go` and have been fixed and migrated to `tests/resolved_issue_test.go`.

| # | Bug | Root Cause | Files Modified | Status |
|---|-----|-----------|----------------|--------|
| 1 | Map with function value type | Closure not wrapped as real Go function | `value/accessor.go`, `vm/vm.go`, `vm/closure.go`, `vm/ops_dispatch.go` | ✅ Fixed |
| 2 | Type switch on interface slice | `int64` stored instead of `int` in interface | `vm/ops_dispatch.go` | ✅ Fixed |
| 3 | Struct with function field | Same root cause as Bug 1 | `value/container.go` | ✅ Fixed |
| 4 | Slice append spread operator | Native `[]int64` / reflect `[]int` type mismatch | `vm/ops_dispatch.go` | ✅ Fixed |
| 5 | Map update during range | Test expected deterministic result for non-deterministic behavior | `tests/known_issues_test.go` | ✅ Fixed |
| 6 | Self-referencing struct type | Infinite recursion in `typeToReflect` | `vm/typeconv.go`, `vm/ops_dispatch.go` | ✅ Fixed |

---

## Bug 1: Map with Function Value Type

### Symptom
```
panic: reflect.Value.SetMapIndex: value of type *vm.Closure is not assignable to type func() int
```
Storing a closure in `map[int]func() int` caused a panic because the VM's internal `*Closure` type is not assignable to a concrete function type.

### Test Case
```go
func MapWithFuncValue() int {
    m := make(map[int]func() int)
    m[1] = func() int { return 10 }
    m[2] = func() int { return 20 }
    return m[1]() + m[2]()   // Expected: 30
}
```

### Root Cause
In Gig's value system, closures are stored as `*vm.Closure` objects with `KindFunc`. When `SetMapIndex` is called with a reflect-based `map[int]func() int`, the value needs to be a real Go `func() int`, not a `*Closure`. The `ToReflectValue` method had no logic to convert closures to real Go functions.

### Fix

**`value/accessor.go`** — Added `ClosureCaller` callback type and `reflect.MakeFunc` wrapping:

```go
type ClosureCaller func(closure any, args []reflect.Value) []reflect.Value

// In ToReflectValue, KindFunc branch:
case KindFunc:
    if typ.Kind() == reflect.Func && closureCaller != nil {
        closure := v.obj
        fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
            results := closureCaller(closure, args)
            // Convert results to match expected return types (e.g., int64 → int)
            ...
        })
        return fn
    }
```

**`vm/vm.go`** — Registered the `ClosureCaller` callback in `init()` to break the `value` → `vm` circular dependency:

```go
func init() {
    value.RegisterClosureCaller(func(closure any, args []reflect.Value) []reflect.Value {
        c := closure.(*Closure)
        // Create temporary VM, execute closure bytecode, return results
        ...
    })
}
```

**`vm/closure.go`** — Added `Program *bytecode.Program` field to `Closure` so the callback can create a child VM.

**`vm/ops_dispatch.go`** — In `OpClosure`, set `closure.Program = vm.program`. In `OpCallIndirect`, added handling for reflect-based function values (closures read back from typed containers become `reflect.Value` of function type).

**`vm/run.go`** — Updated `OpCallIndirect` hot path to handle reflect-based function call via `rv.Call()`.

### Performance Impact
`reflect.MakeFunc` is only invoked when assigning a closure to a typed container (map value, struct field). Normal closure calls through `OpCallIndirect` still use the fast `*Closure` path with zero overhead.

---

## Bug 2: Type Switch on Interface Slice Elements

### Symptom
```
interface slice type switch: got 1110, want 1111
```
Type switch on values extracted from `[]interface{}` failed for `int` type (all other types matched correctly).

### Test Case
```go
func InterfaceSliceTypeSwitch() int {
    var items []interface{}
    items = append(items, 1, "hello", true, 3.14)
    count := 0
    for _, item := range items {
        switch item.(type) {
        case int:    count += 1
        case string: count += 10
        case bool:   count += 100
        case float64: count += 1000
        }
    }
    return count   // Expected: 1111
}
```

### Root Cause
Two issues combined:

1. **Value storage mismatch**: Gig internally stores all integers as `int64`. When `1` (an `int`) is appended to `[]interface{}`, it's stored in the reflect slice as `int64`, not `int`. In native Go, `1` in `interface{}` is stored as `int`.

2. **Strict `AssignableTo` check**: The `OpAssert` handler used `reflect.Type.AssignableTo()` to match types. Since `int64` is not assignable to `int` (they're different types), the `case int:` branch never matched.

### Fix

**`vm/ops_dispatch.go`** — Added `sameReflectKindFamily()` fallback in the `KindReflect` branch of `OpAssert`:

```go
if targetReflectType != nil && concreteRV.Type().AssignableTo(targetReflectType) {
    result = value.MakeFromReflect(concreteRV)
    assertionOk = true
} else if targetReflectType != nil && sameReflectKindFamily(concreteRV.Type(), targetReflectType) {
    // Gig stores int as int64 internally; for type switch, match by kind family
    result = value.MakeFromReflect(concreteRV.Convert(targetReflectType))
    assertionOk = true
}
```

Also added the `kindMatchesType()` helper for the non-reflect path (primitive `KindInt`/`KindString`/etc. values), replacing the previous "assume success" logic.

**`sameReflectKindFamily`** matches types within the same numeric family:
- Signed integers: `int`, `int8`, `int16`, `int32`, `int64`
- Unsigned integers: `uint`, `uint8`, ... `uintptr`
- Floats: `float32`, `float64`
- Complex: `complex64`, `complex128`

### Side Effect Fix
The previous code had an `else { assertionOk = true }` branch that always assumed type assertions succeed. This also caused `TestCompiler_TypeAssertionCommaOk` to expect wrong behavior (`ok=true` for `"hello".(int)`). The test was updated to expect the correct result (`ok=false`).

---

## Bug 3: Struct with Function Field

### Symptom
```
panic: reflect.Set: value of type value.Value is not assignable to type func() int
```
Assigning a closure to a struct's function field caused a panic.

### Test Case
```go
type structWithFunc struct {
    f func() int
}

func StructWithFuncField() int {
    s := structWithFunc{f: func() int { return 42 }}
    return s.f()   // Expected: 42
}
```

### Root Cause
Same underlying issue as Bug 1. Additionally, `value/container.go`'s `SetElem` method for pointer-to-struct used `reflect.ValueOf(val)` instead of `val.ToReflectValue(elemType)`, bypassing the closure-to-function wrapping.

### Fix

**`value/container.go`** — Changed `SetElem` to use `ToReflectValue`:

```go
// Before (broken):
rv.Elem().Set(reflect.ValueOf(val))

// After (fixed):
rv.Elem().Set(val.ToReflectValue(elemType))
```

This ensures the closure is properly wrapped via `reflect.MakeFunc` before being assigned to the struct field.

---

## Bug 4: Slice Append with Spread Operator

### Symptom
```
slice flatten: got 2, want 4
```
`append(result, inner...)` only appended the first element instead of all elements.

### Test Case
```go
func SliceFlatten() int {
    s := [][]int{{1, 2}, {3, 4}}
    result := []int{}
    for _, inner := range s {
        result = append(result, inner...)
    }
    return len(result)   // Expected: 4
}
```

### Root Cause
In the `OpAppend` handler's native `[]int64` fast path, when `elem.IntSlice()` failed (because `inner` was a reflect-based `[]int` from ranging over `[][]int`), the code fell through to `elem.RawInt()` which treated the slice as a single integer, appending only one element.

### Fix

**`vm/ops_dispatch.go`** — Added a reflect slice detection branch in the native `[]int64` fast path:

```go
if s, ok := slice.IntSlice(); ok {
    if es, ok2 := elem.IntSlice(); ok2 {
        vm.push(value.MakeIntSlice(append(s, es...)))
    } else if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice {
        // elem is a reflect-based integer slice (e.g. []int from a [][]int range).
        // Convert each element to int64 and spread-append.
        for i := 0; i < elemRV.Len(); i++ {
            s = append(s, elemRV.Index(i).Int())
        }
        vm.push(value.MakeIntSlice(s))
    } else {
        vm.push(value.MakeIntSlice(append(s, elem.RawInt())))
    }
}
```

---

## Bug 5: Map Update During Range

### Symptom
```
map update during range: got 6, want 7
```
Adding keys during `range` iteration produced fewer entries than expected.

### Test Case
```go
func MapUpdateDuringRange() int {
    m := map[int]int{1: 10, 2: 20}
    for k := range m {
        m[k+10] = k
    }
    return len(m)   // Go spec: non-deterministic, but >= 4
}
```

### Root Cause
The Go specification explicitly states that adding keys during `range` iteration is allowed, but **whether** newly-added keys are visited is non-deterministic. The original test expected exactly 7, which assumed all new keys are always visited.

Investigation showed that `reflect.MapRange()` (used by the VM for map iteration) already correctly observes some newly-added keys, consistent with Go's native behavior. The issue was the test's overly strict expected value.

### Fix

**`tests/known_issues_test.go`** — Changed the assertion from exact match to range check:

```go
// Before: if n != 7 { ... }
// After:
if n < 4 {
    t.Errorf("map update during range: got %d, want >= 4", n)
}
```

The minimum valid result is 4 (2 original keys + 2 keys added from visiting only the original keys).

---

## Bug 6: Self-Referencing Struct Type

### Symptom
```
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp: ... stack: [...
fatal error: stack overflow
```
Creating a self-referencing struct like `type node struct { next *node }` caused infinite recursion in `typeToReflect`.

### Test Case
```go
type node struct {
    value int
    next  *node
}

func StructSelfRef() int {
    n1 := &node{value: 1}
    n2 := &node{value: 2, next: n1}
    return n2.value + n2.next.value   // Expected: 3
}
```

### Root Cause
`typeToReflect` recursively converts `go/types.Type` to `reflect.Type`. For `type node struct { next *node }`, this creates an infinite loop: `node → struct{int, *node} → *node → node → ...`

### Fix (Two Parts)

**Part 1: `vm/typeconv.go`** — Added cycle detection cache:

```go
func typeToReflectWithCache(t types.Type, cache map[types.Type]reflect.Type) reflect.Type {
    if cached, ok := cache[t]; ok {
        return cached  // Returns nil for in-progress types (cycle detected)
    }
    // For *types.Named: mark as in-progress before recursing
    case *types.Named:
        cache[tt] = nil  // sentinel for cycle detection
        result := typeToReflectWithCache(tt.Underlying(), cache)
        cache[tt] = result
        return result
    // For *types.Pointer: when elem is nil (cycle), use interface{} as placeholder
    case *types.Pointer:
        elem := typeToReflectWithCache(tt.Elem(), cache)
        if elem != nil {
            return reflect.PointerTo(elem)
        }
        return reflect.TypeOf(&emptyIface).Elem()  // interface{} placeholder
}
```

**Part 2: `vm/ops_dispatch.go`** — Updated `OpFieldAddr` to handle the `interface{}` placeholder:

When accessing a field through the self-referencing pointer (stored as `interface{}` in the reflect struct), the reflect.Value will be of kind `reflect.Interface` wrapping the actual struct pointer. Added unwrapping logic:

```go
case bytecode.OpFieldAddr:
    // ... existing pointer dereference ...
    // For self-referencing struct types, the recursive pointer field is
    // stored as interface{} by typeToReflect. Unwrap it here.
    if s.Kind() == reflect.Interface && !s.IsNil() {
        s = s.Elem()
        if s.Kind() == reflect.Ptr {
            s = s.Elem()
        }
    }
    if s.Kind() == reflect.Struct { ... }
```

### Design Decision: `interface{}` vs `unsafe.Pointer`
Initially considered `unsafe.Pointer` as the placeholder type, but this caused `reflect.Set` errors because `*struct{...}` is not assignable to `unsafe.Pointer`. `interface{}` works because any Go value can be stored in an interface, and the VM already handles interface unwrapping in other contexts.

---

## Verification

### Test Results
```
$ go test -race ./...
ok   github.com/t04dJ14n9/gig              1.055s
ok   github.com/t04dJ14n9/gig/bytecode     1.014s
ok   github.com/t04dJ14n9/gig/compiler     1.013s
ok   github.com/t04dJ14n9/gig/importer     1.014s
ok   github.com/t04dJ14n9/gig/tests       49.273s
ok   github.com/t04dJ14n9/gig/value        1.014s
ok   github.com/t04dJ14n9/gig/vm           1.124s
```

### Lint Results
```
$ golangci-lint-v2 run
(no issues)
```

### Test Migration
All 6 tests migrated from `tests/known_issues_test.go` → `tests/resolved_issue_test.go`:
- `TestResolved_MapWithFuncValue`
- `TestResolved_InterfaceSliceTypeSwitch`
- `TestResolved_StructWithFuncField`
- `TestResolved_SliceFlatten`
- `TestResolved_MapUpdateDuringRange`
- `TestResolved_StructSelfRef`

---

## Files Modified

| File | Changes |
|------|---------|
| `value/accessor.go` | Added `ClosureCaller` type, `reflect.MakeFunc` wrapping in `ToReflectValue` |
| `value/container.go` | Fixed `SetElem` to use `ToReflectValue` for func-typed pointer fields |
| `vm/vm.go` | Registered `ClosureCaller` callback in `init()` |
| `vm/closure.go` | Added `Program` field to `Closure` struct |
| `vm/ops_dispatch.go` | Fixed `OpAssert` (type switch), `OpAppend` (spread), `OpFieldAddr` (self-ref), `OpClosure` (set Program), `OpCallIndirect` (reflect func) |
| `vm/run.go` | Updated `OpCallIndirect` hot path for reflect functions |
| `vm/typeconv.go` | Added cycle detection cache to `typeToReflect` |
| `tests/known_issues_test.go` | Cleared (all issues resolved) |
| `tests/testdata/known_issues/main.go` | Cleared (all issues resolved) |
| `tests/resolved_issue_test.go` | Added 6 migrated test functions |
| `tests/testdata/resolved_issue/main.go` | Added 6 migrated test data functions |
| `tests/compiler_vm_test.go` | Updated `TypeAssertionCommaOk` to expect correct behavior |
