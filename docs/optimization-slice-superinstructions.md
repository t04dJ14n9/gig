# Fused Int Slice Superinstructions — Optimization Report

## Summary

This optimization adds **3 fused superinstructions** for `[]int` element access, replacing 7-instruction (17-byte) sequences with single-dispatch 7-byte opcodes. This is the culmination of the optimization series, flipping BubbleSort and Sieve from **losing** to **winning** vs Yaegi.

### Results (AMD EPYC 9754, linux/amd64)

| Benchmark | Before Slice Fusion | After Slice Fusion | Change | vs Yaegi |
|---|---|---|---|---|
| Fib25 | 20.7 ms | 20.2 ms | — | **Gig 5.6x faster** |
| ArithSum | 37.2 μs | 37.5 μs | — | **Gig 12% faster** |
| BubbleSort | 1,770 μs (1.40x slower) | **963 μs** | **1.84x faster** | **Gig 32% faster** |
| Sieve | 301 μs (1.44x slower) | **203 μs** | **1.48x faster** | **Gig 3% faster** |
| ClosureCalls | 392 μs | 381 μs | — | **Gig 2.7x faster** |

**Gig now wins all 5 benchmarks against Yaegi.**

### Cumulative Improvement (Full Optimization Series)

| Benchmark | Original | Current | Total Speedup |
|---|---|---|---|
| Fib25 | 169 ms | 20.2 ms | **8.4x** |
| ArithSum | 311 μs | 37.5 μs | **8.3x** |
| BubbleSort | 10.3 ms | 963 μs | **10.7x** |
| Sieve | 1,681 μs | 203 μs | **8.3x** |
| ClosureCalls | 964 μs | 381 μs | **2.5x** |

---

## Problem

Integer slice access (`s[i]` read and `s[i] = v` write) compiled to 7-instruction sequences:

**Read pattern (`v = s[j]`):**
```
LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) DEREF SETLOCAL(v)
```
= 7 dispatches, 17 bytes, multiple stack push/pops, kind checks, and pointer indirection.

**Write pattern (`s[j] = val`):**
```
LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) LOCAL(val) SETDEREF
```
= 7 dispatches, 17 bytes, same overhead.

These patterns dominate BubbleSort (swaps = 2 reads + 2 writes per comparison) and Sieve (1 read + 1 write per inner loop iteration). Each pattern requires 7 instruction dispatches, 7 bytecode decodes, and multiple stack operations — all for what is logically a single array access.

---

## Solution: 3 Fused Opcodes

### OpIntSliceGet(s, j, v) — 7 bytes, 1 dispatch

Fuses the read pattern. At runtime:
```go
slice := locals[s].([]int64)    // type-assert once
result := slice[intLocals[j]]   // direct indexed read
intLocals[v] = result            // 8-byte write
locals[v] = value.MakeInt(result) // 32-byte sync
```

### OpIntSliceSet(s, j, val) — 7 bytes, 1 dispatch

Fuses the write-from-local pattern. At runtime:
```go
slice := locals[s].([]int64)
slice[intLocals[j]] = intLocals[val]  // direct indexed write
```

### OpIntSliceSetConst(s, j, c) — 7 bytes, 1 dispatch

Fuses the write-from-constant pattern (`s[j] = 0` or `s[j] = 1`). At runtime:
```go
slice := locals[s].([]int64)
slice[intLocals[j]] = intConsts[c]    // direct indexed write from const pool
```

This pattern is critical for the Sieve benchmark which does `sieve[i] = 1` and `sieve[j] = 0`.

### Safety: Fallback Path

All three opcodes include a fallback path for non-native slices (e.g., `[]int` backed by `reflect.Value` instead of `[]int64`). The fast path checks `IntSlice()` and falls back to executing the equivalent generic operations:

```go
if s, ok := locals[sIdx].IntSlice(); ok {
    // fast path: direct []int64 access
} else {
    // fallback: execute IndexAddr + Deref/SetDeref generically
}
```

---

## Compiler Changes

### New Pass: `fuseSliceOps()`

A new peephole pass inserted between `optimizeBytecode()` and `intSpecialize()`:

```
Raw bytecode
     │
     ▼
Pass 1: optimizeBytecode()     — arithmetic/compare fusion
     │
     ▼
Pass 2: fuseSliceOps()         — ★ NEW: slice access fusion
     │
     ▼
Pass 3: intSpecialize()        — Value → int64 upgrade
     │
     ▼
Pass 4: fuseIntMoves()         — phi-move fusion
     │
     ▼
Optimized bytecode
```

### Pattern Matching

The pass scans for three 17-byte patterns, all starting with `LOCAL LOCAL INDEXADDR SETLOCAL LOCAL`:

| Pattern | Suffix | Fused Opcode |
|---|---|---|
| 1 (read) | `DEREF SETLOCAL(v)` | `OpIntSliceGet(s,j,v)` |
| 2 (write from local) | `LOCAL(val) SETDEREF` | `OpIntSliceSet(s,j,val)` |
| 3 (write from const) | `CONST(val) SETDEREF` | `OpIntSliceSetConst(s,j,c)` |

### Type Requirements

- `s` must be in `localIsIntSlice` (SSA type is `[]int` or `[]int64`)
- `j` must be in `localIsInt` (SSA type is `int` or `int64`)
- `v`/`val` must be in `localIsInt` (for patterns 1 and 2)
- `ptr` and `ptrGet` must reference the same local (confirming no aliasing)

### `isIntSliceType()` Detection

Only matches `[]int` and `[]int64` (corresponding to the native `[]int64` fast path in the VM):

```go
func isIntSliceType(t types.Type) bool {
    sl, ok := t.Underlying().(*types.Slice)
    if !ok { return false }
    basic, ok := sl.Elem().Underlying().(*types.Basic)
    switch basic.Kind() {
    case types.Int, types.Int64:
        return true
    }
    return false
}
```

### Integration with intSpecialize

The fused opcodes' index/value operands are registered in `intSpecialize` pass 1 so their `intLocals[]` slots are kept in sync:

```go
case bytecode.OpIntSliceGet:
    intUsed[j] = true; intUsed[v] = true
case bytecode.OpIntSliceSet:
    intUsed[j] = true; intUsed[val] = true
case bytecode.OpIntSliceSetConst:
    intUsed[j] = true
```

---

## Concrete Example: BubbleSort Inner Loop

**Source:**
```go
if arr[j] > arr[j+1] {
    arr[j], arr[j+1] = arr[j+1], arr[j]
}
```

**Before (28+ instructions per swap):**
```
LOCAL(arr) LOCAL(j) INDEXADDR SETLOCAL(ptr1)     — 4 dispatches
LOCAL(ptr1) DEREF SETLOCAL(aj)                   — 3 dispatches
LOCAL(arr) LOCAL(j1) INDEXADDR SETLOCAL(ptr2)    — 4 dispatches
LOCAL(ptr2) DEREF SETLOCAL(aj1)                  — 3 dispatches
... compare, swap writes ...                      — 14+ dispatches
```

**After (4 fused + compare + 2 fused = ~7 dispatches for the swap path):**
```
OpIntSliceGet(arr, j, aj)         — 1 dispatch
OpIntSliceGet(arr, j1, aj1)       — 1 dispatch
OpIntGreaterLocalLocalJumpTrue ... — 1 dispatch
OpIntSliceSet(arr, j, aj1)        — 1 dispatch
OpIntSliceSet(arr, j1, aj)        — 1 dispatch
```

**Savings:** ~21 fewer dispatches, ~210 fewer bytes of bytecode decoded, zero intermediate stack traffic.

---

## Files Modified

| File | Changes |
|---|---|
| `bytecode/opcode.go` | Added `OpIntSliceGet`, `OpIntSliceSet`, `OpIntSliceSetConst` with 6-byte operands |
| `compiler/compile_func.go` | Added `isIntSliceType()`, `localIsIntSlice[]` map, wired `fuseSliceOps()` into pipeline |
| `compiler/optimize.go` | Added `fuseSliceOps()` with 3 pattern matchers + `intSpecialize` pass 1 registration |
| `vm/run.go` | Added 3 hot-path handlers with `IntSlice()` fast path + generic fallback |
| `value/value.go` | (Cleanup: removed unused `UnsafeIntSlice()`) |
