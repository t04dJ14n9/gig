# Known Issues Bugfix Report – Batch 2

**Date**: 2026-03-16
**Branch**: feature/dev_youngjin

## Summary

Resolved all 6 remaining known issues tracked in `tests/known_issues_test.go`:

| # | Bug | Root Cause | Fix Location |
|---|-----|-----------|-------------|
| 11 | DeferInClosureWithArg | Compiler stack ordering error | `compiler/compile_value.go` |
| 12 | PointerSwapInStruct | OpDeref alias on pointer fields | `vm/ops_dispatch.go` |
| 13 | StructWithFuncSlice | OpNew created `[]value.Value` for func arrays | `vm/ops_dispatch.go` |
| 14 | StructAnonymousField | Missing PkgPath for anonymous unexported fields | `vm/typeconv.go` |
| 15 | StructEmbeddedInterface | Already passing (tracked as known issue) | — |
| 16 | MapRangeWithBreak | Non-deterministic by Go spec | — |

## Bug 11: DeferInClosureWithArg

**Symptom**: `defer func(v int){ result += v }(10)` inside a closure returned `result = 1` instead of `11`.

**Root Cause**: In `compileDefer` for `*ssa.Function` with free variables, arguments were pushed onto the stack **before** the closure was created. `OpDeferIndirect` pops arguments first (from top of stack), then pops the closure. With the wrong ordering, it popped the closure as an "argument" and the actual argument (10) as the "closure".

**Fix**: Reordered code generation in `compileDefer` to push free variable bindings → `OpClosure` → push arguments → `OpDeferIndirect`. Applied the same fix to both the `*ssa.Function` with FreeVars branch and the `*ssa.MakeClosure` branch.

## Bug 12: PointerSwapInStruct

**Symptom**: `p.a, p.b = p.b, p.a` on a `PtrPair{a: &x, b: &y}` resulted in `*p.b = 2` (same as `*p.a`) instead of `*p.b = 1`.

**Root Cause**: `OpFieldAddr` uses `reflect.NewAt` to create a `**int` pointer to the struct's `*int` field. `OpDeref` on this `**int` calls `rv.Elem()` which returns a `reflect.Value` that is **addressable and settable** — it directly references the struct field's memory. When SSA loads `old_a = *FieldAddr(p, 0)` and later stores `*FieldAddr(p, 0) = old_b`, the store also mutates `old_a` because it's not an independent copy.

**Fix**: In `OpDeref`, when `rv.Elem()` returns a pointer type (`reflect.Ptr`) and is settable (indicating it's from a FieldAddr), create an independent copy via `reflect.ValueOf(elem.Interface())`. This produces a non-addressable value that won't be affected by subsequent stores to the struct field.

## Bug 13: StructWithFuncSlice

**Symptom**: `FuncSliceHolder{funcs: []func() int{...}}` panicked with `[]value.Value is not assignable to []func() int`.

**Root Cause**: SSA compiles slice literals as `Alloc([N]func() int)` + element stores + `Slice`. `OpNew` for `*types.Array` with function elements created a `[]value.Value` instead of a proper `[N]func() int` array. The same issue existed for `*types.Slice` with function elements. When this `[]value.Value` was assigned to the struct's `[]func() int` field via `reflect.Set`, it panicked due to type mismatch.

**Fix**:
1. Removed the `[]value.Value` special case in `OpNew` for both `*types.Slice` and `*types.Array` with function elements. Now uses `typeToReflect` + `reflect.New` to create properly typed arrays/slices (e.g., `[2]func() int`).
2. Updated `SetIndex` in `value/container.go` to use `ToReflectValue(elemType)` when storing closures to function-typed slice elements, enabling the `ClosureCaller + reflect.MakeFunc` wrapping.
3. Added `[]value.Value` → typed slice conversion in `ToReflectValue` (KindSlice branch) for cases where a `[]value.Value` needs to be converted to a typed function slice.
4. Removed the stale `[]value.Value` function slice special case in `OpMakeSlice`.

## Bug 14: StructAnonymousField

**Symptom**: `AnonField{int: 42, name: "test"}` panicked with `reflect.StructOf: field "int" is unexported but missing PkgPath`.

**Root Cause**: `typeToReflectWithCache` in `vm/typeconv.go` checked for unexported fields using `f.Name()[0] >= 'a' && f.Name()[0] <= 'z'` and skipped PkgPath for anonymous fields (`if !f.Anonymous()`). However, `reflect.StructOf` has two conflicting constraints:
- Unexported fields **must** have PkgPath set
- Anonymous fields **must not** have PkgPath set

This makes anonymous unexported fields (like embedded `int`) impossible to represent directly.

**Fix**:
1. Replaced manual lowercase check with `f.Exported()` from `go/types` (handles edge cases like `_` prefix, Unicode).
2. For anonymous unexported fields, demoted them to regular unexported fields (`Anonymous = false`) with PkgPath set. This is the only valid representation in `reflect.StructOf`.

## Files Modified

- `compiler/compile_value.go` — compileDefer stack ordering
- `vm/ops_dispatch.go` — OpDeref pointer copy, OpNew func slice/array, OpMakeSlice cleanup
- `vm/typeconv.go` — anonymous unexported field handling
- `value/accessor.go` — ToReflectValue []value.Value → typed slice conversion
- `value/container.go` — SetIndex function element wrapping
- `tests/testdata/resolved_issue/main.go` — migrated 6 test functions
- `tests/resolved_issue_test.go` — added 6 test cases
- `tests/testdata/known_issues/main.go` — cleared (all issues resolved)
- `tests/known_issues_test.go` — cleared (all issues resolved)

## Test Results

All 24 resolved issue tests pass with `-race` flag. All compiler, vm, value, bytecode, and importer tests pass.
