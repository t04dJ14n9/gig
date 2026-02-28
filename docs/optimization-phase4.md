# Phase 4 Optimization: Dataflow Analysis & Frame Zeroing

## Summary

This document covers the fourth round of optimizations to the Gig VM, targeting **compile-time dataflow analysis** to eliminate redundant runtime work. Two new compiler analysis passes were added:

1. **IntOnlyLocals Analysis** — identifies locals accessed only by `OpInt*` instructions (infrastructure for future dual-write elimination)
2. **Partial Frame Zeroing** — uses `ZeroFrom` threshold + `clear()` to skip zeroing parameter and straightline-written slots

These are **compile-time** analyses. The frame zeroing optimization reduces per-call overhead with zero runtime cost for functions that don't benefit. The IntOnlyLocals analysis is retained as compiler infrastructure for future opcode-level dual-write elimination.

### Design Rationale

An initial Phase 4 implementation added per-instruction `if intOnlyLocals == nil || !intOnlyLocals[idxC]` branches to 7 `OpInt*` handlers in `vm/run.go` and used a `NeedsZero []bool` bitmap for selective frame zeroing. Benchmarking revealed **three regressions**:

1. **Hot-path branch overhead**: The per-instruction nil-check + bitmap-lookup in 7 `OpInt*` handlers cost ~1–3% on Fib(25) with 242K recursive calls, even when the optimization never fired (IntOnlyLocals was nil for `fib`).
2. **memclr intrinsic loss**: Replacing `for i := range s { s[i] = zero }` (recognized by Go compiler as memclr) with per-element conditional `if need { s[i] = zero }` loops destroyed the bulk-zero optimization.
3. **Struct size increase**: Adding `IntOnlyLocals []bool` (24-byte slice header) to `CompiledFunction` caused measurable cache effects on recursive workloads.

The revised approach:
- **Removes all VM hot-path changes** — `vm/run.go` is identical to Phase 3
- **Replaces `NeedsZero []bool`** with `ZeroFrom int` — a single integer threshold that preserves contiguous `clear()` (memclr intrinsic)
- **Removes `IntOnlyLocals` from `CompiledFunction`** — the analysis is computed but discarded, avoiding struct bloat
- **Uses `clear()` builtin** (Go 1.21+) — directly maps to memclr, more reliable than pattern-matching range loops

### Results (AMD EPYC 9754 128-Core Processor, linux/amd64, A/B back-to-back comparison)

#### Core Workloads

| Workload | Native Go | Gig (Phase 4) | Yaegi | GopherLua | Gig vs Native | Gig vs Yaegi | Gig vs Lua |
|---|---:|---:|---:|---:|---:|---:|---:|
| Fibonacci(25) | 448 μs | 19.7 ms | 107 ms | 21.1 ms | 44x slower | **5.4x faster** | **1.07x faster** |
| ArithSum(1K) | 659 ns | 34.3 μs | 39.8 μs | 38.8 μs | 52x slower | **1.2x faster** | **1.1x faster** |
| BubbleSort(100) | 6.4 μs | 935 μs | 1,206 μs | 768 μs | 146x slower | **1.3x faster** | 1.2x slower |
| Sieve(1000) | 1.86 μs | 188 μs | 206 μs | 206 μs | 101x slower | **1.1x faster** | **1.1x faster** |
| ClosureCalls(1K) | 345 ns | 315 μs | 956 μs | 122 μs | 913x slower | **3.0x faster** | 2.6x slower |

#### External Function Calls (Gig vs Yaegi, no Lua/Native equivalent)

| Workload | Native Go | Gig | Yaegi | Gig vs Native | Gig vs Yaegi |
|---|---:|---:|---:|---:|---:|
| DirectCall | 28.3 μs | 495 μs | 1,520 μs | 17x slower | **3.1x faster** |
| Reflect | 24.1 μs | 328 μs | 998 μs | 14x slower | **3.0x faster** |
| Method | 18.2 μs | 418 μs | 1,227 μs | 23x slower | **2.9x faster** |
| Mixed | 11.4 μs | 294 μs | 874 μs | 26x slower | **3.0x faster** |

#### Memory Efficiency (allocs/op)

| Workload | Gig | Yaegi | GopherLua | Gig vs Yaegi |
|---|---:|---:|---:|---:|
| Fibonacci(25) | 6 | 2,138,703 | 41 | **356,450x fewer** |
| ArithSum(1K) | 6 | 8 | 93 | 1.3x fewer |
| BubbleSort(100) | 9 | 5,085 | 12 | **565x fewer** |
| Sieve(1000) | 7 | 43 | 207 | **6x fewer** |
| ClosureCalls(1K) | 1,995 | 13,018 | 96 | **6.5x fewer** |

---

## Compiler Analysis 1: IntOnlyLocals (Infrastructure)

### Problem

In the int-specialized VM path, every `OpInt*SetLocal` instruction performs a **dual write**:

```go
case bytecode.OpIntLocalConstAddSetLocal:
    r := intLocals[idxA] + intConsts[idxB]
    intLocals[idxC] = r              // 8-byte int64 write (fast)
    locals[idxC] = value.MakeInt(r)  // 32-byte Value construction + store (slow)
```

The `locals[idxC]` write exists so that generic (non-int-specialized) code can read the correct value from `locals[]`. However, if a local variable is **only ever accessed** by `OpInt*` instructions, the `locals[]` copy is never read — making the write pure waste.

### Solution

Added a two-phase compile-time analysis:

**Phase 1 — `intSpecialize()` now returns the `intUsed` bitmap** (which locals participate in `OpInt*` ops):

```go
func intSpecialize(code []byte, localIsInt, constIsInt []bool) ([]byte, bool, []bool)
//                                                              ^^^^^^^^^^^^^^^^^^^^^^
//                                                              now returns intUsed
```

**Phase 2 — `buildIntOnlyLocals()` scans the post-specialization bytecode** to find locals accessed ONLY by `OpInt*` instructions:

```go
func buildIntOnlyLocals(code []byte, numLocals int, intUsed []bool) []bool
```

The analysis starts with all `intUsed` locals as candidates, then **revokes** int-only status for any local accessed by a generic instruction:

| Revoking instruction | Reason |
|---|---|
| `OpLocal(idx)` | Generic read of locals[] |
| `OpSetLocal(idx)` | Generic write to locals[] |
| `OpAddSetLocal(idx)` | Fused generic op writes locals[] |
| `OpLocal*Local*SetLocal(A,B,C)` | Fused generic op reads A,B; writes C |
| `OpIntSetLocal(idx)` | Bridge instruction — local crosses int/generic boundary |
| `OpIntMoveLocal(src,dst)` | Conservative: both revoked (reads locals[src]) |

### Current Status: Infrastructure Only

The `buildIntOnlyLocals()` result is **computed but discarded** (`_ = buildIntOnlyLocals(...)`) because applying it in the VM hot path via per-instruction branches caused regressions. The analysis is retained as infrastructure for a future approach:

- **Dedicated no-dual-write opcodes** (e.g., `OpIntLocalConstAddSetLocal_NoDual`) that skip the `locals[]` write unconditionally — no per-instruction branch needed
- The compiler would emit these opcodes for int-only locals, and the standard dual-write opcodes for everything else

### Files Changed

| File | Change |
|---|---|
| `compiler/optimize.go` | `intSpecialize` returns `intUsed`; added `buildIntOnlyLocals()` |
| `compiler/compile_func.go` | Calls `buildIntOnlyLocals` (result discarded with `_`) |

---

## Optimization 2: Partial Frame Zeroing with `ZeroFrom`

### Problem

When a function is called, the VM allocates a frame with `NumLocals` value slots. For recycled frames (from the frame pool), **all** slots are zeroed:

```go
// Before: zero ALL locals on every function call
f.locals = f.locals[:fn.NumLocals]
for i := range f.locals {
    f.locals[i] = value.Value{}  // 32 bytes × NumLocals
}
```

For `fib(25)` with ~242K recursive calls, each zeroing 8 locals × 32 bytes = 256 bytes, the total zeroing work is **~62 MB** of memory fills — just to ensure correctness for slots that are immediately overwritten.

In SSA form, most locals follow a strict "def-before-use" pattern: parameters are always written by the caller, and SSA temporaries are defined before any use. Only Phi-node targets and potentially uninitialized variables need zeroing.

### Solution

Added `computeZeroFrom()` — a compile-time analysis that finds the lowest local index requiring zeroing:

```go
func computeZeroFrom(code []byte, numLocals, numParams int) int
```

The analysis performs a forward scan of the bytecode:

1. **Parameters** (indices 0..NumParams-1): never need zeroing — the caller always writes them
2. **Straightline analysis**: For each local, if the first access is a WRITE (SetLocal/OpInt*SetLocal), the local doesn't need zeroing
3. **At any branch** (Jump/JumpTrue/JumpFalse): conservatively mark all unresolved locals as needs-zero

The result is a single `ZeroFrom int` stored in `CompiledFunction`. The `needsZero` bitmap pattern (where params and straightline-written locals form a contiguous `false` prefix) maps naturally to a single cutoff index. This preserves the `clear()` (memclr) intrinsic:

```go
// After: only zero locals from ZeroFrom onward
if zf := fn.ZeroFrom; zf > 0 {
    clear(f.locals[zf:])   // memclr intrinsic, only zeros what's needed
} else {
    clear(f.locals)        // ZeroFrom=0 → zero all (backward compatible default)
}
```

Key design choices:
- **`ZeroFrom int` instead of `NeedsZero []bool`**: The bitmap would require per-element conditional zeroing, destroying the memclr intrinsic. A single threshold preserves contiguous `clear()`.
- **`clear()` builtin (Go 1.21+)**: Directly calls memclr, more reliable than the `for i := range s { s[i] = zero }` pattern which depends on Go compiler pattern recognition.
- **Zero value = zero all**: `ZeroFrom = 0` means all locals need zeroing, which is the safe backward-compatible default.

### Files Changed

| File | Change |
|---|---|
| `bytecode/bytecode.go` | Added `ZeroFrom int` to `CompiledFunction` |
| `compiler/optimize.go` | Added `computeZeroFrom()` function |
| `compiler/compile_func.go` | Wired `computeZeroFrom` into pipeline |
| `vm/frame.go` | `framePool.get()` uses `clear()` + `ZeroFrom` for both `locals` and `intLocals` |

### Impact

For functions where parameters dominate the local count (like `fib` with 3 params, 8 locals), saves ~37% of zeroing work per call. For straight-line code with many SSA temporaries, savings can reach 80%+. The `clear()` builtin ensures maximum throughput via memclr intrinsic regardless of the slice size.

---

## Status Update

### What Was Already Implemented (Pre-Phase 4)

During the investigation for this phase, we confirmed that several planned optimizations were **already implemented in Phase 3**:

1. **Backward-Jump Context Check** — `backJumpCount` throttle in `OpJump` handler (implemented Phase 3)
2. **Closure Pooling** — `sync.Pool` for `*Closure` objects in `vm/closure.go` (implemented Phase 3)
3. **DirectCall Expansion** — `canWrapUnderlying` already supports structs, pointers, non-empty interfaces, maps, cross-package named types. Method DirectCall wrappers already generated:
   - **619 / 671 functions** have DirectCall wrappers (92.3% coverage)
   - **543 method DirectCalls** across 24 stdlib packages
   - Remaining 52 `nil` cases are all **function-parameter types** (callbacks like `strings.Map`, `sort.Slice`) — deferred to future work

### Cumulative Progress

| Benchmark | Baseline | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|---|---:|---:|---:|---:|---:|
| Fib(25) | 24.1 ms | 22.4 ms | 20.3 ms | 19.7 ms | 19.7 ms |
| ArithSum(1K) | 99.1 μs | 39.6 μs | 35.2 μs | 34.1 μs | 34.3 μs |
| BubbleSort(100) | 2,124 μs | 1,027 μs | 913 μs | 943 μs | 935 μs |
| Sieve(1000) | — | — | — | 187 μs | 188 μs |
| ClosureCalls(1K) | — | — | 369 μs | 326 μs | 315 μs |

### Interpreter Comparison Summary

| | Gig | Yaegi | GopherLua |
|---|---|---|---|
| **Language** | Go subset | Go subset | Lua |
| **Architecture** | Stack-based bytecode VM | AST walker | Register-based bytecode VM |
| **Fib(25) speed** | **5.4x faster** than Yaegi | baseline | ~same as Gig |
| **Ext call speed** | **3x faster** than Yaegi | baseline | N/A |
| **Memory efficiency** | 6 allocs/op (Fib) | 2.1M allocs/op | 41 allocs/op |
| **DirectCall coverage** | 92% funcs, 543 methods | N/A | N/A |

### Architecture (with Phase 4 additions marked *)

```
Source Code
    │
    ▼
SSA Compiler (golang.org/x/tools/go/ssa)
    │
    ▼
Bytecode Compiler (compiler/compile_func.go)
    │
    ├─── Pass 1: optimizeBytecode()     [fuse generic superinstructions]
    ├─── Pass 2: fuseSliceOps()         [fuse []int access patterns]
    ├─── Pass 3: intSpecialize()        [upgrade to OpInt* variants]
    ├─── Pass 3.5*: buildIntOnlyLocals() [identify dual-write-safe locals — infrastructure]
    ├─── Pass 4: fuseIntMoves()         [fuse OpIntLocal+OpIntSetLocal]
    └─── Pass 5*: computeZeroFrom()     [compute partial-zeroing threshold]
    │
    ▼
Optimized Bytecode + Metadata
    (ZeroFrom threshold*)
    │
    ▼
VM Execution (vm/run.go — unchanged from Phase 3)
    ├── framePool.get(): partial zeroing via ZeroFrom + clear()*
    ├── Backward-jump context checks (Phase 3)
    └── Closure pooling via sync.Pool (Phase 3)
```
