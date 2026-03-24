# Panic, Recover, and Dynamic Call Stacks

This document describes the implementation of full `panic`/`recover`/`defer` support in the
Gig interpreter, along with the dynamic call frame stack that was introduced alongside it.
The work spans the compiler, the VM runtime, and the test infrastructure.

**Files changed:**

| File | Description |
|------|-------------|
| `gig.go` | Added `WithAllowPanic()` build option |
| `compiler/build.go` | Added `WithAllowPanic()` compiler option, panic ban check |
| `compiler/parser/parse.go` | Static analysis to reject `panic()` in sandboxed code |
| `compiler/compile_func.go` | `detectResultAllocSlots()` for named return recovery |
| `bytecode/bytecode.go` | `ResultAllocSlots` field on `CompiledFunction` |
| `vm/vm.go` | `panicState`, `panicStack`, `deferDepth`, `growFrames()` |
| `vm/run.go` | Panic-path defer execution, `derefAllocLocal()` |
| `vm/ops_control.go` | `OpRecover` with `panicStack` awareness |
| `vm/call.go` | Dynamic frame growth at call sites |
| `vm/closure.go` | Updated to use `initialFrameDepth` constant |
| `vm/goroutine.go` | Updated to use `initialFrameDepth` constant |

---

## 1. Design Goals

In Go, `panic`/`recover` is a fundamental control flow mechanism. A sandboxed interpreter
that cannot support it will silently break any user code that relies on error recovery
patterns — a common idiom in production Go. The goals were:

1. **Safety by default.** Panic is banned at compile time unless explicitly opted in with
   `gig.WithAllowPanic()`. This preserves the sandbox guarantee for users who don't need it.

2. **Semantic fidelity.** When enabled, `panic`, `recover`, and `defer` must behave exactly
   like native Go — including edge cases like nested panics, panics inside deferred
   functions, named return modification, and LIFO defer ordering after recovery.

3. **Containment.** An unrecovered panic in interpreted code must never crash the host
   process. It is converted to an `error` return from `Program.Run()`.

4. **Zero cost when unused.** The panic machinery adds no overhead to programs that don't
   use it. The `ResultAllocSlots` field is nil, `deferDepth` stays at 0, and `panicStack`
   is never allocated.

---

## 2. Three-Layer Safety Model

Panic handling in Gig operates at three distinct layers:

```
Layer 1: Compile-time ban (default)
    ↓  (if WithAllowPanic)
Layer 2: VM-level panic/recover (guest code)
    ↓  (if unrecovered / host-level panic)
Layer 3: Go-level safety net (defer/recover in Execute)
```

### Layer 1 — Compile-Time Ban

By default, the compiler rejects any source code that calls `panic()`:

```
compile error: panic() is not allowed in sandboxed code (at main.go:10:5);
use gig.WithAllowPanic() to enable
```

This is a static analysis pass in `compiler/parser/parse.go` that walks the AST looking
for calls to the `panic` builtin. It runs before SSA construction, so rejected programs
never reach the compiler backend.

### Layer 2 — VM-Level Panic/Recover

When `WithAllowPanic()` is set, the compiler emits `OpPanic` and `OpRecover` opcodes.
The VM tracks panic state on the `vm` struct:

```go
type vm struct {
    panicking  bool            // a panic is in progress
    panicVal   value.Value     // the panic argument
    panicStack []panicState    // saved states for nested panics
    deferDepth int             // nesting level of deferred execution
    // ...
}
```

When `OpPanic` executes, it sets `panicking = true` and stores the value. The main
execution loop checks `v.panicking` at the top of every iteration and enters the
panic-handling path instead of continuing normal execution.

### Layer 3 — Go-Level Safety Net

Both `Execute()` and `ExecuteWithValues()` wrap `v.run()` in a Go-level `defer/recover`:

```go
func (v *vm) Execute(...) (result value.Value, err error) {
    defer func() {
        if r := recover(); r != nil {
            result = value.MakeNil()
            err = fmt.Errorf("runtime panic: %v", r)
        }
    }()
    result, err = v.run()
    return result, err
}
```

This catches any Go-level panics that escape the VM — nil pointer dereferences, slice
bounds violations, type assertion failures, and the deliberate `panic("gig: call stack
overflow")` from the dynamic frame growth system.

---

## 3. Panic-Path Defer Execution

The most intricate part of the implementation is executing deferred functions during panic
unwinding. This is handled in `vm/run.go` at the top of the main loop:

```go
if v.panicking {
    // Run ALL defers in LIFO order, with panic/recover awareness
}
```

### 3.1 The Core Algorithm

When a panic is active and the current frame has deferred functions, the VM runs them
in last-in-first-out order. Each defer follows this protocol:

```
1. Push current panic state onto panicStack
2. Clear v.panicking (so recursive run() doesn't re-enter the panic handler)
3. Execute the deferred function via v.run()
4. Inspect the result:
   a. If v.panicking is true:  the defer itself panicked (new panic replaces old)
      → Pop saved state, continue to next defer with new panic active
   b. If saved state still has panicking=true:  defer didn't call recover()
      → Restore the original panic, continue to next defer
   c. If saved state has panicking=false:  recover() was called
      → Mark as recovered, continue running remaining defers normally
```

### 3.2 Why panicStack?

Consider this code:

```go
defer func() {
    defer func() {
        recover()  // recovers "inner"
    }()
    panic("inner")
}()
panic("outer")
```

When `panic("outer")` triggers, the outer defer runs. Inside it, `panic("inner")` triggers
a new panic. The inner defer calls `recover()`. At this point, `OpRecover` needs to find
the correct panic value — which is `"inner"`, not `"outer"`.

The solution is `panicStack`: a stack of saved `(panicking, panicVal)` states. Each level
of deferred execution pushes the current panic state before entering the defer's `run()`.
`OpRecover` checks both `v.panicking` (direct context) and the top of `panicStack`:

```go
case bytecode.OpRecover:
    if v.panicking {
        v.push(v.panicVal)
        v.panicking = false
    } else if len(v.panicStack) > 0 && v.panicStack[len(v.panicStack)-1].panicking {
        v.push(v.panicStack[len(v.panicStack)-1].panicVal)
        v.panicStack[len(v.panicStack)-1].panicking = false
    } else {
        v.push(value.MakeNil())
    }
```

### 3.3 All Defers Run After Recovery

A critical Go semantic: **all deferred functions execute regardless of whether a panic
was recovered**. After `recover()` clears the panic, remaining defers in the same frame
still run — just in normal mode (no panic context).

The VM tracks this with a `recovered` flag. Once set, subsequent defers execute via a
child VM (like `OpRunDefers`) rather than through the shared panic-path `v.run()`:

```go
if v.panicking {
    // Panic-path: save/restore panicStack, use shared v.run()
} else {
    // Normal mode: run in isolated child VM
    childVM := &vm{program: v.program, ...}
    childVM.run()
}
```

Using a child VM for post-recovery defers prevents them from interfering with the parent
frame stack. If a post-recovery defer itself panics, the child VM's `panicking` flag is
checked and the parent re-enters panic mode.

### 3.4 deferDepth and Frame Boundaries

The `deferDepth` counter tracks how many levels deep we are in deferred execution. It
serves two purposes:

1. **OpReturn / OpReturnVal**: When `deferDepth > 0`, a function return means the
   deferred function finished — return from the recursive `v.run()` immediately instead
   of continuing with the caller's frame.

2. **Recovery return**: After panic recovery pops a frame, if `deferDepth > 0`, the VM
   must return from the current `v.run()` immediately. Without this check, the recursive
   `v.run()` would continue executing the *outer* function's instructions — the frame
   below the one that was just popped — causing stack corruption and index-out-of-range
   panics.

---

## 4. Named Return Value Recovery

In Go, only **named return values** can be modified by deferred functions during panic
recovery. This is because named returns are heap-allocated variables (SSA `Alloc`
instructions) whose pointers are shared between the function body, deferred closures,
and the return path.

Consider:

```go
func example() (result int) {
    defer func() { recover(); result = 42 }()
    panic("test")
}
// Returns 42 — 'result' is a named return, modified by defer
```

vs:

```go
func example() int {
    var result int
    defer func() { recover(); result = 42 }()
    panic("test")
    return result
}
// Returns 0 — 'result' is a local variable, not a named return
```

### 4.1 Compile-Time Detection

The compiler identifies named return Allocs by scanning SSA `Return` instructions:

```go
func detectResultAllocSlots(fn *ssa.Function, st *SymbolTable) []int {
    // Find all Alloc instructions
    // Check which Allocs appear (directly or via UnOp deref) in Return results
    // Record their local slot indices
}
```

These slot indices are stored on `CompiledFunction.ResultAllocSlots`.

The SSA representation of named returns follows this pattern:

```
result = Alloc *int           ; heap-allocate the named return
Store 0:int result            ; initialize to zero
...
RunDefers                     ; execute defers (may modify *result)
t1 = *result                  ; deref to get final value
Return t1
```

The compiler detects that `result` (an Alloc) appears in the Return instruction via
the `UnOp` deref (`t1 = *result`), and records `result`'s local slot index.

### 4.2 VM-Side Reconstruction

After panic recovery, the VM reads the named return value from the frame's locals:

```go
if slots := frame.fn.ResultAllocSlots; len(slots) > 0 {
    ptr := frame.locals[slots[0]]
    retVal = derefAllocLocal(ptr)  // dereference the Alloc pointer
}
```

`derefAllocLocal` mirrors what `OpDeref` does in the normal return path — it reads
the value behind the pointer that the deferred closure wrote to:

```go
func derefAllocLocal(ptr value.Value) value.Value {
    switch ptr.Kind() {
    case value.KindPointer:
        return ptr.Elem()
    case value.KindInterface:
        if rv, ok := ptr.ReflectValue(); ok {
            if rv.Kind() == reflect.Ptr && !rv.IsNil() {
                return value.MakeFromReflect(rv.Elem())
            }
        }
    }
    // ...
}
```

### 4.3 Why Not Mark All Allocs?

An earlier approach marked *all* Allocs as result slots when the function had any defers.
This produced incorrect results: local variables like `order` and `result` in
`MultipleDefersOnPanic` were all returned as a tuple `[0, 3, 123]` instead of the
expected `0` (the zero value, since none are named returns).

The correct approach marks **only** Allocs that appear in Return instructions. Functions
with no named returns have `ResultAllocSlots = nil`, and the VM returns `value.MakeNil()`
(which evaluates to 0 for int, "" for string, etc.) — matching Go's behavior.

---

## 5. Dynamic Call Frame Stack

### 5.1 The Problem

The VM's frame stack was a fixed-size slice of 64 `*Frame` pointers. Any recursion
deeper than 64 levels would write past the end of the slice, causing an index-out-of-range
panic caught by the Layer 3 safety net — an opaque error message for the user.

Several legitimate test cases were disabled because of this:

- `RecursionAckermann(3, 4)` — needs ~500 frames
- `RecursionSum(100)` — needs 101 frames
- `RecursionStaircase(10)` and `RecursionHanoi(10)` — moderate depth but disabled conservatively

### 5.2 The Solution

The frame stack now starts at 64 and grows dynamically by doubling, up to a maximum
of 1024:

```go
const (
    initialFrameDepth = 64
    maxFrameDepth     = 1024
)

func (v *vm) growFrames() bool {
    cur := len(v.frames)
    if cur >= maxFrameDepth {
        return false
    }
    newCap := cur * 2
    if newCap > maxFrameDepth {
        newCap = maxFrameDepth
    }
    grown := make([]*Frame, newCap)
    copy(grown, v.frames)
    v.frames = grown
    return true
}
```

The growth check is placed at the two frame-push sites (`callCompiledFunction` and
`callFunction`):

```go
if v.fp >= len(v.frames) {
    if !v.growFrames() {
        panic("gig: call stack overflow")
    }
}
v.frames[v.fp] = frame
v.fp++
```

### 5.3 Performance Impact

**Zero cost on the hot path.** The bounds check (`v.fp >= len(v.frames)`) is a single
comparison that the branch predictor will always predict as not-taken, since 99%+ of
programs never exceed 64 frames. The check is at function call boundaries, which already
perform frame pool allocation — so the relative cost is negligible.

**Memory:** The baseline allocation is identical (64 × 8 = 512 bytes). Only programs
that exceed 64 frames pay for the doubled allocation (128 × 8 = 1 KB, then 256, 512,
up to 8 KB at the 1024 cap).

**Pooled VMs:** When a VM is returned to the pool via `Reset()`, the grown frame slice
is retained. This means a pooled VM that previously ran a deep-recursion program will
start with the larger slice on reuse — a natural warm-up effect.

### 5.4 Stack Overflow Reporting

When the hard cap is hit, `growFrames()` returns false and the call site issues
`panic("gig: call stack overflow")`. This Go-level panic is caught by the Layer 3
safety net and returned as:

```
error: runtime panic: gig: call stack overflow
```

This is a clear, actionable error message compared to the previous opaque
"runtime error: index out of range [64] with length 64".

---

## 6. Test Coverage

### Panic/Recover Tests — 27 cases

All test functions live in `tests/testdata/panic_recover/main.go` and are verified
against native Go execution via reflection.

| Category | Tests |
|----------|-------|
| Basic recovery | `PanicRecoverBasic`, `PanicRecoverWithValue`, `PanicRecoverInt`, `NoPanicNoRecover` |
| Defer + panic | `DeferRunsOnPanic`, `MultipleDefersOnPanic`, `DeferModifyBeforePanic` |
| Nested panics | `NestedPanicRecover`, `NestedRecover`, `PanicInDefer`, `PanicChain` |
| Recover values | `RecoverReturnsNilWhenNotPanicking`, `RecoverReturnsPanicValueCheck` |
| Named returns | `NamedReturnPanicRecover`, `NamedReturnDeferModify` |
| Complex scenarios | `PanicInLoop`, `PanicInClosure`, `MultiplePanicSameDefer`, `PanicInRecursiveFunction`, `DeferClosureCapturePanic`, `PanicInDeferWithRecoverInDefer`, `RecoverOnlyInDefer`, `RecoverInGoroutine` |
| Edge cases | `EmptyDeferPanic`, `DeferOrderWithMultiplePanics`, `DeferPanicRecoverChain`, `PanicNil` |

### Recursion Tests — 4 re-enabled

Previously disabled due to the 64-frame limit, now passing with dynamic growth:

- `RecursionAckermannCheck` — Ackermann(3, 4) = 125
- `RecursionSumCheck` — Sum(100) = 5050
- `RecursionHanoiCheck` — Hanoi(10) = 1023
- `RecursionStaircaseCheck` — Staircase(10) = 89

### Full Suite

1938 tests pass, 0 failures.

---

## 7. Go Semantics — Lessons Learned

Several Go semantics were initially misunderstood, leading to bugs that were caught by
native comparison testing:

### 7.1 Only Named Returns Are Visible After Recovery

The most impactful discovery. Six test cases had incorrect expectations:

```go
func f() int {           // unnamed return
    var result int
    defer func() {
        recover()
        result = 42       // modifies local variable
    }()
    panic("x")
    return result
}
// Native Go returns: 0 (zero value of int), NOT 42
```

Only `(result int)` named returns are modified by defers during panic recovery.
This is because unnamed returns are evaluated *before* `RunDefers`, and after panic
recovery, the function returns the zero value of the declared return type.

### 7.2 All Defers Run After Recovery

An initial implementation broke out of the defer loop after `recover()` cleared the
panic. This is wrong — Go runs **all** deferred functions in the frame regardless:

```go
defer func() { fmt.Println("A") }()    // runs third
defer func() { recover() }()            // runs second, recovers
panic("x")                              // runs first
// Output: B then A — both defers run
```

### 7.3 Nested Panics Don't Replace the Outer Panic

When a deferred function panics during panic unwinding, the new panic is handled by
*that defer's own defers*. If recovered internally, the original panic continues
unwinding the outer function's remaining defers:

```go
defer func() { recover() }()        // recovers "outer"
defer func() {
    defer func() { recover() }()    // recovers "inner"
    panic("inner")
}()
panic("outer")
// Both panics are recovered separately; function returns normally
```

### 7.4 Recovery Must Return Immediately in Recursive run()

After panic recovery pops a frame inside a recursive `v.run()` (deferDepth > 0), the
VM must `return` immediately. Without this, the loop `for v.fp > 0` would continue
executing the *outer* function's frame — the one below the popped frame — because the
recursive `v.run()` has access to the shared frame stack.

---

## 8. Architecture Diagram

```
                    ┌─────────────────────────────────┐
                    │          Source Code             │
                    └───────────────┬─────────────────┘
                                    │
                    ┌───────────────▼─────────────────┐
                    │   Layer 1: Compile-Time Ban     │
                    │   (reject panic unless opted in) │
                    └───────────────┬─────────────────┘
                                    │ WithAllowPanic()
                    ┌───────────────▼─────────────────┐
                    │   Compiler: SSA → Bytecode      │
                    │   OpPanic, OpRecover, OpDefer   │
                    │   detectResultAllocSlots()       │
                    └───────────────┬─────────────────┘
                                    │
                    ┌───────────────▼─────────────────┐
                    │        VM Execution             │
                    │                                 │
                    │  ┌──────────────────────────┐   │
                    │  │  Layer 2: VM Panic Path  │   │
                    │  │  panicStack + deferDepth │   │
                    │  │  ResultAllocSlots deref  │   │
                    │  └──────────────────────────┘   │
                    │                                 │
                    │  ┌──────────────────────────┐   │
                    │  │ Layer 3: Go Safety Net   │   │
                    │  │ defer { recover() } in   │   │
                    │  │ Execute/ExecuteWithValues │   │
                    │  └──────────────────────────┘   │
                    │                                 │
                    │  ┌──────────────────────────┐   │
                    │  │ Dynamic Frame Stack      │   │
                    │  │ 64 → 128 → ... → 1024   │   │
                    │  │ growFrames() at call     │   │
                    │  └──────────────────────────┘   │
                    └─────────────────────────────────┘
```
