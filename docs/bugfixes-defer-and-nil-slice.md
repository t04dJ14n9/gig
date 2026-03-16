# Bug Fixes — Defer Execution & Nil Slice Append

This document describes seven confirmed bugs that were discovered and fixed in this release.
All bugs relate to **defer semantics** (named return value modification, closure capture,
stack ordering) and **nil-slice append** handling.

**Files changed:**

| File | Description |
|------|-------------|
| `bytecode/opcode.go` | Added `OpRunDefers` and `OpDeferIndirect` opcodes |
| `compiler/compile_instr.go` | Fixed `*ssa.RunDefers` compilation (was emitting `OpRecover`) |
| `compiler/compile_value.go` | Added `compileDefer` for direct and closure-based defers |
| `value/container.go` | Cleaned up dead `**Value` handling in `Elem()` / `SetElem()` |
| `vm/frame.go` | Added `closure` field to `DeferInfo` for indirect defers |
| `vm/ops_dispatch.go` | New `OpRunDefers` / `OpDeferIndirect` / `OpFree` / `OpClosure` / `OpAppend` fixes |
| `vm/vm.go` | Added `runDefers` helper method |

---

## Bug 1: `*ssa.RunDefers` compiled as `OpRecover` (Root Cause)

**Symptom:** All defer-related test cases failed — `DeferNamedReturn`, `DeferModifyNamed`,
`DeferStackOrder`, and `ClosureWithDefer` all produced incorrect results or panicked.

**Root cause:** In `compiler/compile_instr.go`, the SSA instruction `*ssa.RunDefers` was
incorrectly compiled as `OpRecover`:

```go
// BEFORE (incorrect)
case *ssa.RunDefers:
    c.emit(bytecode.OpRecover)
```

In Go's SSA IR, `RunDefers` is a critical instruction that appears **before** the function
reads its named return values and returns. The sequence is:

```
*result = 5          // set named return
rundefers            // execute all deferred functions (may modify *result)
t = *result          // read (possibly modified) named return
return t
```

By emitting `OpRecover` instead, deferred functions were never executed at the correct
point, and named return values were never modified by defers.

**Fix:** Emit the correct `OpRunDefers` opcode:

```go
// AFTER (correct)
case *ssa.RunDefers:
    c.emit(bytecode.OpRunDefers)
```

**Test case:** `DeferNamedReturn` — expected 10 (5 × 2), was returning 5.

---

## Bug 2: `OpRunDefers` opcode did not exist

**Symptom:** Even after fixing the compiler, there was no VM handler for executing defers
synchronously before return.

**Root cause:** The bytecode instruction set had no opcode for "run all pending defers now."
Defers were only executed at return time (`OpReturn`/`OpReturnVal`), which happens **after**
the named return values have already been read — too late for defers to modify them.

**Fix:** Added `OpRunDefers` to `bytecode/opcode.go` and implemented a handler in
`vm/ops_dispatch.go` that:

1. Iterates pending defers in **LIFO** order (last deferred = first executed)
2. Creates a **child VM** for each deferred call (avoids interfering with parent frame stack)
3. Shares globals, context, program, and external call cache with the parent VM
4. Executes each defer **synchronously** before continuing

```go
case bytecode.OpRunDefers:
    for len(frame.defers) > 0 {
        d := frame.defers[len(frame.defers)-1]
        frame.defers = frame.defers[:len(frame.defers)-1]
        var freeVars []*value.Value
        if d.closure != nil {
            freeVars = d.closure.FreeVars
        }
        childVM := &VM{
            program: vm.program, stack: make([]value.Value, 256),
            globals: vm.globals, globalsPtr: vm.globalsPtr,
            ctx: vm.ctx, extCallCache: vm.extCallCache,
        }
        deferFrame := newFrame(d.fn, 0, d.args, freeVars)
        childVM.frames[0] = deferFrame
        childVM.fp = 1
        _, _ = childVM.run()
    }
```

Correspondingly, defer execution was **removed** from `OpReturn` and `OpReturnVal`, since
defers now run at the correct SSA-defined point via `OpRunDefers`.

**Test case:** `DeferStackOrder` — expected 1111 (1000 + 100 + 10 + 1), was panicking.

---

## Bug 3: `OpFree` double-wrapped pointer values

**Symptom:** `ClosureWithDefer` and `DeferModifyNamed` failed because closures could not
correctly read/write shared variables through free variable capture.

**Root cause:** `OpFree` was using `value.FromInterface(frame.freeVars[idx])` which wrapped
the `*value.Value` pointer in a new `Value`, creating a double indirection
(`Value` → `*Value` → `Value`). Closures expecting to read a simple `int` via a reflect
pointer would instead get a `*value.Value` wrapper — type mismatch.

**Fix:** Changed `OpFree` to directly dereference the slot:

```go
// BEFORE
vm.push(value.FromInterface(frame.freeVars[idx]))

// AFTER
vm.push(*frame.freeVars[idx])
```

This ensures the closure sees the **actual captured value** (e.g., a `reflect.Value` of
type `*int`), not a wrapper around the slot pointer.

**Test case:** `DeferModifyNamed` — expected 999, was returning 42.

---

## Bug 4: `OpClosure` incorrect free variable slot creation

**Symptom:** Multiple closures sharing the same captured variable would not see each
other's modifications.

**Root cause:** `OpClosure` had complex logic trying to detect `*value.Value` vs
`**value.Value` from the stack, but it was inconsistent. The slot-sharing mechanism
(where multiple closures reference the same `*value.Value` slot to share state) was broken.

**Fix:** Simplified `OpClosure` to always create a fresh `*value.Value` slot for each
captured variable:

```go
slot := new(value.Value)
*slot = v
closure.FreeVars[i] = slot
```

If the captured value is a reflect pointer (e.g., `*int` from `Alloc`), all closures
sharing that pointer will see each other's modifications through the shared underlying
heap-allocated int, even though each has its own `*value.Value` wrapper slot.

**Test case:** `ClosureWithDefer` — expected 30, was returning nil/panicking.

---

## Bug 5: `OpDeferIndirect` for closure-based defers

**Symptom:** `defer func() { result *= 2 }()` — defers using anonymous closures were
not compiled or executed correctly.

**Root cause:** The compiler had no support for `OpDeferIndirect` which defers a closure
call (as opposed to `OpDefer` which defers a named function call). SSA's `Defer` instruction
can target either a `*ssa.Function` or a `*ssa.MakeClosure`, and the latter requires
capturing free variables at defer time.

**Fix:** Added `OpDeferIndirect` to the bytecode, compiler, and VM:

- **Bytecode:** New opcode with operand `[num_args:2]`
- **Compiler:** `compileDefer` now handles `*ssa.Function`, `*ssa.MakeClosure`, and
  fallback cases
- **VM:** `OpDeferIndirect` pops args and closure from stack, stores `DeferInfo` with
  the closure's `FreeVars` for later execution

**Test case:** `DeferNamedReturn` — `defer func() { result *= 2 }()` is a closure defer.

---

## Bug 6: `Slice_AppendToNil` — nil slice + native `[]int64` element

**Symptom:** `var s []int; s = append(s, 1)` panicked with:
`reflect.Set: value of type int is not assignable to type []int64`

**Root cause:** SSA compiles `append(s, 1)` as:
1. Create `[1]int{1}` array
2. Slice it to `[]int{1}` (internally stored as `[]int64{1}`)
3. Call `append(s, sliced_result)`

In the `OpAppend` handler, when `s` is nil, the nil-slice branch checks
`elem.ReflectValue()`. But for a native `[]int64`, `ReflectValue()` returns `false`
(native slices are not stored as `reflect.Value`). The code fell through to the
single-element append path, which treated `[]int64{1}` as a single element and tried to
create `[][]int64` — type mismatch.

**Fix:** Added a fast path at the start of the nil-slice branch to check for native
`[]int64`:

```go
if es, ok2 := elem.IntSlice(); ok2 {
    vm.push(value.MakeIntSlice(append([]int64(nil), es...)))
    break
}
```

**Test case:** `Slice_AppendToNil` — expected 3, was panicking.

---

## Bug 7: `NilSliceAppend` — `append(nil, 1, 2, 3)` (pre-existing)

**Symptom:** `var s []int; s = append(s, 1, 2, 3); return len(s)` returned 1 instead of 3.

**Root cause:** Same as Bug 6. SSA packs variadic args `1, 2, 3` into `[]int{1, 2, 3}`
(stored as `[]int64{1, 2, 3}`), then calls `append(nil, packed_slice)`. The nil-slice
branch failed to recognize the native `[]int64` element and only appended the first value.

**Fix:** Same as Bug 6 — the `IntSlice()` fast path correctly spread-appends all elements.

**Test case:** `NilSliceAppend` — expected 3, was returning 1.

---

## Cleanup: Removed dead `**Value` handling

After fixing `OpFree` (Bug 3), the `**Value` detection logic in `value/container.go`
(`Elem()`, `SetElem()`) and `vm/ops_dispatch.go` (`OpDeref`) became dead code. These
paths were added as workarounds for the double-wrapping that `OpFree` used to create.
Since `OpFree` now directly dereferences the slot, `**Value` never appears on the stack.

Removed the following dead code:
- `Elem()`: `**Value` dereference check
- `SetElem()`: `**Value` write-through check
- `OpDeref`: `**value.Value` unwrap check

---

## Lint Fix: `gci` import formatting

**Symptom:** `golangci-lint-v2` reported "File is not properly formatted (gci)" for
`bytecode/opcode.go` and `vm/ops_dispatch.go`.

**Root cause:** The `gci` formatter (configured in `.golangci.yml` under `formatters`)
enforces import section ordering: `standard → default → prefix(gig) → localmodule`.
The modified files did not conform to this grouping.

**Fix:** Applied `golangci-lint-v2 run --fix` to auto-format both files.

---

## Test Results

All 7 bugs are verified with dedicated test cases:

| Test Case | Expected | Before Fix | After Fix |
|-----------|----------|------------|-----------|
| `DeferNamedReturn` | 10 | 5 | ✅ 10 |
| `DeferModifyNamed` | 999 | 42 | ✅ 999 |
| `DeferStackOrder` | 1111 | panic | ✅ 1111 |
| `ClosureWithDefer` | 30 | panic | ✅ 30 |
| `MultipleNamedReturn` | 1021 | (untested) | ✅ 1021 |
| `Slice_AppendToNil` | 3 | panic | ✅ 3 |
| `NilSliceAppend` | 3 | 1 | ✅ 3 |

Full test suite: `go test ./...` — **0 failures** across all packages.

Lint check: `golangci-lint-v2 run` — **0 issues**.
