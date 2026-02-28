# Integer Specialization & Value Shrinking — Optimization Report

## Summary

This document covers two major optimization phases that together delivered a **38–62% speedup** across all benchmarks: **Value struct shrinking** (56B → 32B) and **integer-specialized locals** (`intLocals []int64`). These phases address the two largest remaining performance gaps identified in previous reports — the oversized `Value` struct and the cost of manipulating 32-byte values in the inner loop.

### Final Results (AMD EPYC 9754, linux/amd64)

| Benchmark | Before | After | Speedup | vs Yaegi |
|---|---|---|---|---|
| Fib(25) | 32.2 ms | **19.8 ms** | **1.63x** | **Gig 5.6x faster** |
| ArithSum(1K) | 137 μs | **52.6 μs** | **2.60x** | Yaegi 1.28x faster |
| BubbleSort(100) | 3.93 ms | **2.37 ms** | **1.66x** | Yaegi 1.9x faster |
| Sieve(1000) | 691 μs | **388 μs** | **1.78x** | Yaegi 1.9x faster |
| ClosureCalls(1K) | 535 μs | **374 μs** | **1.43x** | **Gig 2.7x faster** |

### Cumulative Improvement (from project start)

| Benchmark | Original | Current | Total Speedup |
|---|---|---|---|
| Fib(25) | 169 ms | 19.8 ms | **8.5x** |
| ArithSum(1K) | 311 μs | 52.6 μs | **5.9x** |
| BubbleSort(100) | 10.3 ms | 2.37 ms | **4.3x** |
| Sieve(1000) | 1,681 μs | 388 μs | **4.3x** |
| ClosureCalls(1K) | 964 μs | 374 μs | **2.6x** |

---

## Part 1: Value Struct Shrinking (56B → 32B)

### Problem

The `Value` struct — the universal representation for all runtime values in the VM — was 56 bytes:

```go
// Before: 56 bytes
type Value struct {
    kind Kind              // 1 byte + 7 bytes padding
    obj  any               // 16 bytes (interface: type pointer + data pointer)
    num  int64             // 8 bytes
    str  string            // 16 bytes (string header: data pointer + length)
}
```

Every push, pop, local read/write, and function argument pass copies 56 bytes. In a tight loop executing 3 superinstructions per iteration, each touching 2-3 locals, this means moving **336–504 bytes/iteration** of Value data. The 56-byte size also causes poor cache utilization — only 1 Value fits in a cache line (64 bytes), leaving 8 bytes wasted.

### Solution

Merged the `str string` field into the `obj any` field, since strings and other heap-allocated types are never used simultaneously:

```go
// After: 32 bytes
type Value struct {
    kind Kind    // 1 byte + 7 bytes padding
    num  int64   // 8 bytes: bool (0/1), int, uint bits, float64 bits
    obj  any     // 16 bytes: string, complex128, reflect.Value, or nil
}
```

**Key design decisions:**

1. **`num int64` stores all scalar types** — booleans as 0/1, int as-is, uint via bit cast, float64 via `math.Float64bits`/`math.Float64frombits`. This enables `RawInt() int64` to be a simple field access: `return v.num`.

2. **`obj any` is the universal heap field** — Strings are stored as `string` values boxed into the `any` interface. Complex numbers, reflect.Values, maps, slices, channels, etc., all go here. The `kind` field disambiguates at runtime.

3. **`RawInt()` and `RawBool()` are unchecked** — These accessors skip kind-checking, relying on SSA type guarantees. `RawInt()` returns `v.num` directly; `RawBool()` returns `v.num != 0`. Both are small enough for the Go compiler to inline.

### Impact

| Metric | Before | After |
|---|---|---|
| Value size | 56 bytes | 32 bytes |
| Values per cache line | 1 (8B wasted) | 2 (exact fit) |
| Bytes moved per loop iter (ArithSum) | ~504 bytes | ~288 bytes |
| Memory per 1000 locals | 56 KB | 32 KB |

**Files modified:**
- `value/value.go` — Restructured `Value`, updated all constructors (`MakeInt`, `MakeString`, `MakeFloat`, etc.)
- `value/accessor.go` — Updated `String()`, `Int()`, `Float64()`, `ToReflectValue()` to read from `obj`
- `value/arithmetic.go` — Updated `Add`, `Sub`, `Mul`, `Div`, `Cmp` to use new field layout
- `value/convert.go` — Updated `Convert()` and `FromInterface()` for new field layout
- `value/container.go` — Updated container operations for new string/slice storage
- `value/value_test.go` — Updated tests for new struct layout

---

## Part 2: Integer-Specialized Locals (`intLocals []int64`)

### Problem

Even after shrinking `Value` to 32 bytes, integer-heavy inner loops still move 32-byte structs through `locals[]` for every read and write. For a simple loop like `sum += i; i++`, each iteration reads and writes multiple `value.Value` structs (32 bytes each) when the actual payload is just an `int64` (8 bytes). This is a 4x overhead in data movement.

Additionally, every integer operation must:
1. Check `Kind() == KindInt` (branch)
2. Extract via `RawInt()` (field access)
3. Compute the result
4. Box via `MakeInt()` (construct a 32-byte struct)
5. Store the 32-byte struct back

Steps 1, 2, 4, and 5 are pure overhead when the compiler can statically prove a variable is always `int`.

### Solution: Shadow `int64` Array

Added a parallel `intLocals []int64` array to each frame for functions that contain integer-typed locals. Integer-specialized opcodes (`OpInt*`) operate directly on this 8-byte array, bypassing the 32-byte `Value` entirely in the hot path.

```
┌─────────────────────────────────────────────┐
│ Frame                                        │
│                                              │
│  locals []value.Value   (32 bytes per slot)  │  ← used by generic opcodes
│  ┌──────┬──────┬──────┬──────┐              │
│  │ V[0] │ V[1] │ V[2] │ V[3] │              │
│  └──────┴──────┴──────┴──────┘              │
│                                              │
│  intLocals []int64       (8 bytes per slot)  │  ← used by OpInt* opcodes
│  ┌──────┬──────┬──────┬──────┐              │
│  │ i[0] │ i[1] │ i[2] │ i[3] │              │
│  └──────┴──────┴──────┴──────┘              │
│                                              │
│  Both arrays are kept in sync via dual-write │
└─────────────────────────────────────────────┘
```

### The Dual-Write Invariant

**Every write to an int-specialized local must update BOTH `intLocals[idx]` and `locals[idx]`.**

This invariant exists because non-specialized code (generic `OpLocal`, closures, function returns, the debugger) reads from `locals[]`. If we only wrote to `intLocals[]`, those code paths would see stale data.

```go
// OpIntLocalConstAddSetLocal — the hot path
r := intLocals[idxA] + intConsts[idxB]   // pure int64 arithmetic
intLocals[idxC] = r                        // fast 8-byte write
locals[idxC] = value.MakeInt(r)            // sync 32-byte shadow
```

The dual-write costs one extra `MakeInt` + 32-byte store per operation, but the critical benefit is that **reads** (`OpIntLocal`) only touch the 8-byte `intLocals[]` array, and the compare+jump fusions (`OpIntLess*`) skip both `Kind()` checks and `RawInt()` extraction.

### Integer-Specialized Opcodes

13 new `OpInt*` opcodes operate on `intLocals []int64` and `intConsts []int64`:

**Fused arithmetic (zero stack traffic, 8-byte operands):**

| Opcode | Semantics | Width |
|---|---|---|
| `OpIntLocalConstAddSetLocal` | `intLocals[C] = intLocals[A] + intConsts[B]` | 7B |
| `OpIntLocalConstSubSetLocal` | `intLocals[C] = intLocals[A] - intConsts[B]` | 7B |
| `OpIntLocalLocalAddSetLocal` | `intLocals[C] = intLocals[A] + intLocals[B]` | 7B |

**Fused compare+branch (zero stack traffic, no kind check):**

| Opcode | Semantics | Width |
|---|---|---|
| `OpIntLessLocalConstJumpTrue` | `if intLocals[A] < intConsts[B] { goto off }` | 7B |
| `OpIntLessLocalConstJumpFalse` | `if intLocals[A] >= intConsts[B] { goto off }` | 7B |
| `OpIntLessEqLocalConstJumpTrue` | `if intLocals[A] <= intConsts[B] { goto off }` | 7B |
| `OpIntLessEqLocalConstJumpFalse` | `if intLocals[A] > intConsts[B] { goto off }` | 7B |
| `OpIntLessLocalLocalJumpTrue` | `if intLocals[A] < intLocals[B] { goto off }` | 7B |
| `OpIntLessLocalLocalJumpFalse` | `if intLocals[A] >= intLocals[B] { goto off }` | 7B |
| `OpIntGreaterLocalLocalJumpTrue` | `if intLocals[A] > intLocals[B] { goto off }` | 7B |

**Bridge opcodes (sync between intLocals and the stack):**

| Opcode | Semantics | Width |
|---|---|---|
| `OpIntSetLocal` | `intLocals[idx] = pop().RawInt(); locals[idx] = pop()` | 3B |
| `OpIntLocal` | `push(MakeInt(intLocals[idx]))` | 3B |

**Register move (fused phi-move):**

| Opcode | Semantics | Width |
|---|---|---|
| `OpIntMoveLocal` | `intLocals[dst] = intLocals[src]; locals[dst] = locals[src]` | 5B |

### The `intConsts []int64` Constant Pool

A parallel constant pool stores pre-extracted `int64` values alongside the existing `PrebakedConstants []value.Value` pool. This allows `OpInt*` compare and arithmetic opcodes to read constants as raw `int64` without any `RawInt()` call or kind check.

```go
// In Program (bytecode/bytecode.go)
type Program struct {
    PrebakedConstants []value.Value   // for generic opcodes
    IntConstants      []int64         // for OpInt* opcodes (parallel array)
    // ...
}
```

---

## The Compilation Pipeline

The optimization pipeline runs in three sequential passes after bytecode generation:

```
Raw bytecode
     │
     ▼
┌─────────────────────────────────┐
│ Pass 1: Peephole Fusion         │   optimizeBytecode()
│  • 6-instruction fusions        │   Fuses common patterns into
│  • 4-instruction fusions        │   superinstructions (Value-typed)
│  • 3-instruction fusions        │
│  • 2-instruction fusions        │
│  • No-op jump elimination       │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│ Pass 2: Int Specialization      │   intSpecialize()
│  • Two-pass algorithm           │   Upgrades superinstructions
│  • Candidate identification     │   to OpInt* variants + inserts
│  • Opcode upgrade + bridging    │   bridge instructions
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│ Pass 3: Int Move Fusion         │   fuseIntMoves()
│  • OpIntLocal + OpIntSetLocal   │   Fuses phi-move pairs into
│    → OpIntMoveLocal             │   single register moves
└────────────┬────────────────────┘
             │
             ▼
  Optimized bytecode
```

### Pass 1: Peephole Superinstruction Fusion

`optimizeBytecode()` in `compiler/optimize.go` scans the instruction stream for known patterns and fuses them into single superinstructions. Patterns are checked longest-first to avoid partial matches.

**6-instruction fusions (16B → 7B, 6 dispatches → 1):**

Handles the SSA compare pattern where a comparison result is stored to a temp, then loaded and branched on:

```
LOCAL(A) CONST(B) LESS SETLOCAL(X) LOCAL(X) JUMPFALSE(off)
└──────────────────────────────────────────────────────────┘
  → OpLessLocalConstJumpFalse(A, B, off)
```

The optimizer validates that `SETLOCAL(X)` and `LOCAL(X)` reference the same slot (the SSA temporary), confirming this is a pure compare-and-branch with no side effects.

**4-instruction fusions (10B → 7B, 4 dispatches → 1):**

Arithmetic: `LOCAL(A) LOCAL(B) ADD SETLOCAL(C)` → `OpLocalLocalAddSetLocal(A,B,C)`
Compare: `LOCAL(A) CONST(B) LESS JUMPTRUE(off)` → `OpLessLocalConstJumpTrue(A,B,off)`

**3-instruction fusions (7B → 5B, 3 dispatches → 1):**

`LOCAL(A) LOCAL(B) ADD` → `OpAddLocalLocal(A,B)` (result stays on stack)

**2-instruction fusions (4B → 3B, 2 dispatches → 1):**

`ADD SETLOCAL(A)` → `OpAddSetLocal(A)` (pop two, add, store)

**No-op jump elimination:**

`JUMP(target)` where target is the next instruction is removed entirely. These arise from SSA compilation of unconditional branches between adjacent basic blocks.

### Jump Target Fixup

All fusions change instruction sizes, which invalidates jump targets. The optimizer uses a `rewrite` list that records `(start, end, newBytes)` for each fusion. After collecting all rewrites, `applyRewrites()` rebuilds the bytecode and `fixJumpTargets()` adjusts every jump offset using an offset mapping table.

### Pass 2: Integer Specialization (`intSpecialize`)

This is a **two-pass** algorithm that upgrades qualifying superinstructions to their `OpInt*` variants.

**Why two passes?** A single pass can't correctly handle bridge instructions. Consider:

```
OpSetLocal(x)              ← appears before any OpInt* instruction
OpLocalConstAddSetLocal(x, 1, x)   ← upgradeable to OpInt*
```

In a single pass, when we see `OpSetLocal(x)`, we don't yet know that `x` will participate in int-specialized operations. We'd miss the bridge upgrade. The two-pass approach solves this:

**Pass 1 — Candidate identification:**

Scans the entire bytecode for superinstructions whose ALL local and constant operands are statically confirmed as `int` type (via `localIsInt[]` and `constIsInt[]` built from SSA type information). For each qualifying instruction, marks the involved local indices in an `intUsed []bool` set.

```go
// Example: OpLocalConstAddSetLocal(A, B, C) qualifies if
// localIsInt[A] && constIsInt[B] && localIsInt[C]
// → marks intUsed[A] = true, intUsed[C] = true
```

If no qualifying instructions exist, returns `(code, false)` — no `intLocals` array is allocated.

**Pass 2 — Opcode upgrade + bridge insertion:**

Two types of transformations:

1. **Superinstruction upgrade** — Replaces the opcode byte in-place (operand layout is identical):
   ```
   OpLocalConstAddSetLocal → OpIntLocalConstAddSetLocal
   OpLessLocalConstJumpFalse → OpIntLessLocalConstJumpFalse
   // ... (all qualifying superinstructions)
   ```

2. **Bridge insertion** — For any `OpSetLocal`/`OpLocal` that references a local in the `intUsed` set:
   ```
   OpSetLocal(idx) → OpIntSetLocal(idx)   // dual-writes to both arrays
   OpLocal(idx)    → OpIntLocal(idx)      // reads from intLocals
   ```

Bridges ensure that non-specialized code paths (function arguments, generic arithmetic on the same variable) keep `intLocals` in sync with `locals`.

### Pass 3: Int Move Fusion (`fuseIntMoves`)

After `intSpecialize` creates `OpIntLocal` + `OpIntSetLocal` pairs (typically from SSA phi-move patterns), this pass fuses them:

```
OpIntLocal(src) OpIntSetLocal(dst)
→ OpIntMoveLocal(src, dst)
```

This eliminates a push+pop cycle (load to stack, then store from stack) and replaces it with a direct register-to-register move:

```go
intLocals[dst] = intLocals[src]   // 8-byte copy
locals[dst] = locals[src]          // 32-byte copy
```

Savings: 6 bytes → 5 bytes, 2 dispatches → 1 dispatch, and eliminates 2 stack operations.

---

## Runtime Support

### Frame Allocation (`vm/frame.go`)

```go
type Frame struct {
    fn        *bytecode.CompiledFunction
    ip        int
    basePtr   int
    locals    []value.Value    // 32 bytes per slot — universal
    intLocals []int64          // 8 bytes per slot — int shadow
    freeVars  []*value.Value
    defers    []DeferInfo
    addrTaken bool
}
```

`intLocals` is only allocated when `fn.HasIntLocals == true` (set by `intSpecialize`). The `framePool` reuses both slices when recycling frames, avoiding allocation on every call.

**Argument mirroring in `newFrame()`:**

```go
if fn.HasIntLocals {
    f.intLocals = make([]int64, fn.NumLocals)
    for i, arg := range args {
        if i < fn.NumLocals {
            f.intLocals[i] = arg.RawInt()  // mirror int params
        }
    }
}
```

### Function Calls (`vm/call.go`)

Both `callCompiledFunction` and `callFunction` (for closures) mirror arguments:

```go
// In callCompiledFunction:
v := vm.pop()
frame.locals[i] = v
if intL != nil {
    intL[i] = v.RawInt()    // mirror every argument to intLocals
}
```

This ensures `intLocals` is initialized correctly when entering a function, even if the caller used generic (non-int-specialized) code.

### VM Dispatch (`vm/run.go`)

The `run()` loop hoists `intLocals` and `intConsts` into Go local variables alongside `stack`, `sp`, `locals`, and `prebaked`:

```go
intLocals := frame.intLocals     // may be nil
intConsts := vm.program.IntConstants
```

All `OpInt*` handlers operate directly on these locals, avoiding repeated `frame.intLocals` dereferences. The loop re-loads these after any frame change (function call, return).

**Example hot-path handler:**

```go
case bytecode.OpIntLocalConstAddSetLocal:
    idxA := readU16(code, ip+1)
    idxB := readU16(code, ip+3)
    idxC := readU16(code, ip+5)
    ip += 7
    r := intLocals[idxA] + intConsts[idxB]   // pure int64 add
    intLocals[idxC] = r                        // 8-byte write
    locals[idxC] = value.MakeInt(r)            // 32-byte sync
    continue
```

---

## Concrete Example: ArithSum Inner Loop

**Source:**
```go
sum := 0
for i := 0; i < 1000; i++ {
    sum += i
}
```

### Before (Value-typed superinstructions, 3 dispatches):
```
OpLocalLocalAddSetLocal(sum, i, sum)   — reads 2×32B, writes 1×32B, kind checks
OpLocalConstAddSetLocal(i, 1, i)       — reads 1×32B + 1×32B const, writes 1×32B
OpLessLocalConstJumpTrue(i, 1000, -)   — reads 1×32B + 1×32B const, 2 kind checks
```
**Per iteration:** 3 dispatches, ~288 bytes moved, 4 kind checks, 3 `RawInt()` calls, 3 `MakeInt()` calls.

### After (int-specialized, 3 dispatches):
```
OpIntLocalLocalAddSetLocal(sum, i, sum) — reads 2×8B, writes 1×8B + 1×32B sync
OpIntLocalConstAddSetLocal(i, 1, i)     — reads 1×8B + 1×8B, writes 1×8B + 1×32B sync
OpIntLessLocalConstJumpTrue(i, 1000, -) — reads 1×8B + 1×8B, no write
```
**Per iteration:** 3 dispatches, ~112 bytes moved, 0 kind checks, 0 `RawInt()` calls, 2 `MakeInt()` calls (dual-write only).

**Net savings:** 176 fewer bytes moved per iteration, 4 fewer branches (kind checks), 1 fewer `MakeInt()` call. On 1000 iterations, this eliminates 176KB of data movement and 4000 branch instructions.

---

## Concrete Example: Fibonacci Recursion

**Source:**
```go
func fib(n int) int {
    if n <= 1 { return n }
    return fib(n-1) + fib(n-2)
}
```

### Key optimization:

The `n <= 1` comparison compiles through SSA as:
```
LOCAL(n) CONST(1) LESSEQ SETLOCAL(t) LOCAL(t) JUMPFALSE(else)
```

The 6-instruction fusion reduces this to:
```
OpLessEqLocalConstJumpFalse(n, 1, else)
```

Then `intSpecialize` upgrades it to:
```
OpIntLessEqLocalConstJumpFalse(n, 1, else)
```

This single instruction replaces 6 dispatches and eliminates all intermediate stack operations. Combined with the recursive calls (which mirror `n` into `intLocals` at each frame), the entire hot path operates on 8-byte `int64` values.

---

## Correctness Guarantees

### Why dual-write is necessary

Consider a function that has both int-specialized and generic code paths for the same variable:

```go
x := computeInt()    // OpIntSetLocal — writes intLocals[x] and locals[x]
if condition {
    useAsInt(x)      // OpIntLocal — reads intLocals[x] ✓
} else {
    useAsAny(x)      // OpLocal — reads locals[x] ✓ (dual-write kept it in sync)
}
```

Without dual-write, the `OpLocal` read in the else branch would see a stale `locals[x]` value.

### Non-int operations on intUsed locals

When `OpAddSetLocal` operates on a local that is also in the `intUsed` set, it must update `intLocals` too:

```go
case bytecode.OpAddSetLocal:
    // ... compute result ...
    locals[idx] = result
    if intLocals != nil {
        intLocals[idx] = result.RawInt()  // keep shadow in sync
    }
```

### Frame transitions

Function calls (`callCompiledFunction`, `callFunction`) and frame creation (`newFrame`) all mirror int arguments into `intLocals`. This ensures `intLocals` is correctly initialized regardless of whether the caller used int-specialized opcodes.

---

## Files Modified

| File | Changes |
|---|---|
| `value/value.go` | Shrunk Value from 56B to 32B (merged `str` into `obj`) |
| `value/accessor.go` | Updated string/type accessors for new field layout |
| `value/arithmetic.go` | Updated arithmetic ops for new field layout |
| `value/convert.go` | Updated type conversion for new field layout |
| `value/container.go` | Updated container ops for new string storage |
| `value/value_test.go` | Updated unit tests for 32B Value |
| `bytecode/opcode.go` | Added 13 `OpInt*` opcodes + `OpLessEqLocalConstJumpFalse/True` |
| `bytecode/bytecode.go` | Added `HasIntLocals` flag, `IntConstants []int64` pool |
| `compiler/compile_func.go` | Integrated `intSpecialize()` and `fuseIntMoves()` into pipeline |
| `compiler/optimize.go` | Added 6-instr fusions, `intSpecialize()`, `fuseIntMoves()`, no-op jump elimination |
| `compiler/compiler.go` | Built `IntConstants` pool during compilation |
| `vm/frame.go` | Added `intLocals []int64` to Frame, updated pool |
| `vm/run.go` | Added all `OpInt*` handlers with dual-write, hoisted `intLocals`/`intConsts` |
| `vm/call.go` | Added intLocals mirroring in `callCompiledFunction` and `callFunction` |

---

## Testing

All optimizations maintain full backward compatibility. The complete test suite passes:

```
ok  gig/bytecode    0.002s
ok  gig/compiler    0.002s
ok  gig/importer    0.003s
ok  gig/tests       0.852s
ok  gig/value       0.003s
ok  gig/vm          0.003s
```

---

## Updated Architecture

```
Source Code (.go)
       │
       ▼
  Go SSA Package (golang.org/x/tools/go/ssa)
       │
       ▼
  Compiler (gig/compiler)
   ├── Allocate local slots + build localIsInt/constIsInt maps
   ├── Compile basic blocks → raw bytecode
   ├── Patch jump targets
   ├── Pass 1: Peephole fusion (superinstructions)
   ├── Pass 2: intSpecialize (Value → int64 upgrade)
   └── Pass 3: fuseIntMoves (phi-move fusion)
       │
       ▼
  Bytecode (gig/bytecode)
   ├── PrebakedConstants []value.Value   ← generic constant pool
   ├── IntConstants []int64              ← int-specialized constant pool
   ├── HasIntLocals bool                 ← per-function flag
   └── 13 OpInt* opcodes + bridges
       │
       ▼
  VM Execution (gig/vm)
   ├── Frame: locals []Value + intLocals []int64 (shadow)
   ├── Dual-write invariant on every int local mutation
   ├── Register-hoisted dispatch (stack, sp, locals, intLocals, intConsts)
   ├── Frame pooling reuses both locals and intLocals slices
   └── Argument mirroring on function entry
```

---

## Remaining Optimization Opportunities

| Priority | Optimization | Expected Impact |
|---|---|---|
| P0 | Direct-threaded dispatch (computed goto via assembly) | ~1.2–1.5x on all workloads |
| P1 | Eliminate dual-write via escape analysis (prove locals[] never read for int vars) | ~1.1–1.2x on arithmetic loops |
| P2 | Native `[]bool` slices for Sieve-like workloads | ~1.3x on Sieve |
| P3 | Closure allocation pooling | Reduce 3,000 allocs/op in ClosureCalls |
| P4 | Float64 specialization (analogous to int specialization) | ~1.3x on float-heavy workloads |
