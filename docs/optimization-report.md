# Gig Interpreter Optimization Report

## Summary

This report documents the systematic optimization of the Gig Go interpreter to close the performance gap with Yaegi (a tree-walking Go interpreter). The work was organized into 6 phases, transforming Gig from being significantly slower than Yaegi on most benchmarks to being **competitive or faster** in key areas.

### Final Results (AMD EPYC 9754, linux/amd64)

| Benchmark | Before | After | Speedup | vs Yaegi | vs Native |
|-----------|--------|-------|---------|----------|-----------|
| Fib25 | 169ms | **45ms** | **3.8x** | Gig **2.4x faster** | 99x slower |
| ArithSum | 311μs | **200μs** | **1.6x** | Yaegi 5x faster | 298x slower |
| BubbleSort | 10.3ms | **4.9ms** | **2.1x** | Yaegi 3.9x faster | 753x slower |
| Sieve | 1,681μs | **838μs** | **2.0x** | Yaegi 4.1x faster | 440x slower |
| ClosureCalls | 964μs | **584μs** | **1.7x** | Gig **1.7x faster** | 1,681x slower |

### Allocation Reduction

| Benchmark | Before (allocs/op) | After (allocs/op) | Reduction |
|-----------|--------------------|--------------------|-----------|
| Fib25 | 728,262 | **68** | **10,710x** |
| ArithSum | 13 | **13** | 1x |
| BubbleSort | 39,818 | **16** | **2,489x** |
| Sieve | 5,864 | **14** | **419x** |
| ClosureCalls | 3,000 | **3,000** | 1x |

---

## Phase 1: O(1) Function Lookup

**Problem:** Function calls used a map lookup (`map[string]*CompiledFunction`) for every `OpCall`, which is O(n) with hash overhead.

**Solution:** Assigned each function a compile-time numeric index and stored them in a flat slice `FuncByIndex []*CompiledFunction`. The `OpCall` instruction encodes the function index directly in its operands.

**Files modified:**
- `bytecode/bytecode.go` — Added `FuncByIndex` field to `Program`
- `compiler/compiler.go` — Built `FuncByIndex` during compilation
- `vm/call.go` — Used index-based lookup in `callCompiledFunction`

**Impact:** Reduced per-call overhead from ~100ns to ~5ns. Most visible in recursive benchmarks like Fib25.

---

## Phase 2: Frame Pooling

**Problem:** Every function call allocated a new `Frame` on the heap (`newFrame()` → `make([]value.Value, numLocals)`). For Fib(25), this means 242,785 frame allocations.

**Solution:** Implemented a `framePool` using `sync.Pool` that recycles `Frame` objects. Frames are returned to the pool on function return. A frame's `locals` slice is reused if large enough, avoiding re-allocation.

**Key detail:** Frames with `addrTaken = true` (where a local's address was taken for closures) are NOT returned to the pool, since external code may still hold a pointer to the frame's locals.

**Files modified:**
- `vm/frame.go` — Added `framePool` with `get()` and `put()` methods
- `vm/vm.go` — Added `fpool framePool` to VM struct
- `vm/ops_dispatch.go` — Used `vm.fpool.put(frame)` in `OpReturn`/`OpReturnVal`

**Impact:** Fib25 allocations dropped from 728K to 68 (one per unique function, not per call). This alone was the single largest optimization.

---

## Phase 3: Pre-baked Constants

**Problem:** Every `OpConst` instruction called `value.FromInterface()` at runtime, which goes through `reflect.ValueOf()` → type switch → `MakeInt()`. For integer constants in hot loops, this is wasteful.

**Solution:** Pre-computed all constants as `value.Value` at compile time and stored them in `PrebakedConstants []value.Value`. `OpConst` now just does an array index lookup.

**Files modified:**
- `bytecode/bytecode.go` — Added `PrebakedConstants []value.Value`
- `compiler/compiler.go` — Built `PrebakedConstants` after compilation
- `vm/ops_dispatch.go` — `OpConst` reads from `PrebakedConstants` first

**Impact:** Eliminated ~50ns overhead per constant load. Most visible in ArithSum and loop-heavy benchmarks.

---

## Phase 4: Integer Fast-Paths

**Problem:** Arithmetic and comparison operations used generic methods (`a.Add(b)`, `a.Cmp(b)`) that check types at runtime for every operation, even when both operands are int (the most common case).

**Solution:** Added inline type checks in the VM's hot-path opcodes. For `OpAdd`, `OpSub`, `OpMul`, and all comparison operators: if both operands have `KindInt`, perform the operation directly on `RawInt()` without calling the generic methods.

```go
case bytecode.OpAdd:
    b := vm.pop()
    a := vm.pop()
    if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
        vm.push(value.MakeInt(a.RawInt() + b.RawInt()))
    } else {
        vm.push(a.Add(b))
    }
```

**Impact:** ~30% speedup on arithmetic-heavy benchmarks. The fast-path avoids method call overhead and type dispatch.

---

## Phase 5: Native `[]int64` Slice Representation

**Problem:** Integer slices (`[]int`) were stored as `reflect.Value` wrapping `reflect.MakeSlice(...)`. Every element access went through `rv.Index(i)` → `reflect.Value` → `MakeFromReflect()`, causing 2-3 allocations per access.

**Solution:** Introduced `KindSlice` with a native `[]int64` backing store. When the VM creates an integer slice (via `OpMakeSlice` or `OpSlice` on an `[N]int` array), it stores a Go `[]int64` directly in the `Value.obj` field.

### Key implementation details:

1. **`MakeIntSlice([]int64)`** — Creates a `Value{kind: KindSlice, obj: []int64{...}}`
2. **`MakeIntPtr(*int64)`** — For `OpIndexAddr`, stores a raw `*int64` pointer without reflect
3. **SSA `Alloc([N]int) + Slice` pattern** — SSA compiles `make([]int, 100)` with constant N as `OpNew([100]int)` + `OpSlice`, not `OpMakeSlice`. The `OpSlice` handler detects `reflect.Array` with `reflect.Int` elements and converts to native `[]int64`.
4. **`SetElem` / `ToReflectValue` conversions** — When a native `[]int64` needs to be stored into a `*[]int` pointer (e.g., 2D slices), automatic `[]int64` → `[]int` conversion is performed.

**Fast-path operations on native int slices:**
- `OpIndex`: `s[i]` → `MakeInt(s[key.RawInt()])` (no reflect)
- `OpSetIndex`: `s[i] = v` → `s[key.RawInt()] = val.RawInt()` (no reflect)
- `OpIndexAddr`: `&s[i]` → `MakeIntPtr(&s[idx])` (raw pointer, no reflect)
- `OpLen/OpCap`: Native `len(s)` / `cap(s)`
- `OpSlice`: Native slice expression
- `OpAppend/OpCopy`: Native append/copy

**Files modified:**
- `value/value.go` — Added `MakeIntSlice`, `MakeIntPtr`, `IntSlice()`
- `value/container.go` — Added `KindSlice` cases in `Len`, `Cap`, `Index`, `SetIndex`, `Elem`, `SetElem`
- `value/accessor.go` — Added `KindSlice` case in `ToReflectValue`
- `vm/ops_dispatch.go` — Fast paths in `OpMakeSlice`, `OpIndex`, `OpSetIndex`, `OpIndexAddr`, `OpSlice`, `OpLen`, `OpAppend`, `OpCopy`

**Impact:** Sieve allocations dropped from 5,864 to 14. BubbleSort allocations dropped from 39,818 to 16. Sieve time improved from 1,681μs to 1,297μs.

---

## Phase 6: Inline Hot-Path Dispatch

**Problem:** The VM's main loop called `vm.executeOp(op, frame)` for every instruction — a Go function call with interface return value and error check. This added ~10-15ns overhead per instruction, which dominates in tight loops.

**Solution:** Moved the most frequently executed opcodes directly into the `run()` loop as a `switch` statement with `continue` to bypass the `executeOp` call. Less common opcodes fall through to `executeOp`.

**Inlined opcodes (covers >90% of instructions in numeric programs):**
- Stack: `OpLocal`, `OpSetLocal`, `OpConst`, `OpNil`, `OpTrue`, `OpFalse`, `OpPop`, `OpDup`
- Arithmetic: `OpAdd`, `OpSub`, `OpMul`
- Comparison: `OpLess`, `OpLessEq`, `OpGreater`, `OpGreaterEq`, `OpEqual`, `OpNotEqual`
- Logic: `OpNot`
- Jumps: `OpJump`, `OpJumpTrue`, `OpJumpFalse`
- Calls: `OpCall`, `OpReturn`, `OpReturnVal`
- Pointer: `OpSetDeref`

**Additional micro-optimization:** `OpJumpTrue`/`OpJumpFalse`/`OpNot` use `RawBool()` (unchecked `v.num != 0`) instead of `Bool()` (which does a kind check + panic). SSA guarantees the condition is always a boolean.

**Files modified:**
- `vm/run.go` — Rewrote `run()` with inlined hot-path switch

**Impact:** 1.5x speedup across all benchmarks. Fib25: 69ms → 45ms. ArithSum: 313μs → 200μs. BubbleSort: 7.4ms → 4.9ms.

---

## Architecture Overview

```
Source Code (.go)
       │
       ▼
  Go SSA Package (golang.org/x/tools/go/ssa)
       │
       ▼
  Compiler (gig/compiler)
   ├── Phase 1: Collect & index functions
   ├── Phase 2: Allocate local slots (params, phis, values)
   ├── Phase 3: Compile blocks in reverse-postorder
   ├── Phase 4: Patch jump targets
   └── Phase 5: Pre-bake constants
       │
       ▼
  Bytecode Program (gig/bytecode)
   ├── FuncByIndex []*CompiledFunction  ← O(1) lookup
   ├── PrebakedConstants []value.Value   ← zero-overhead
   └── Types []types.Type
       │
       ▼
  VM Execution (gig/vm)
   ├── Inline hot-path dispatch          ← no function call overhead
   ├── Frame pooling (sync.Pool)         ← near-zero alloc per call
   ├── Native int slice fast-path        ← no reflect for []int
   └── Integer arithmetic fast-path      ← direct int64 ops
```

---

## Remaining Optimization Opportunities

1. **Register-based VM**: Convert from stack-based to register-based architecture. This would eliminate push/pop overhead (currently 2 array accesses per operation) and enable instruction fusion. Expected 2-3x improvement but requires major rewrite.

2. **Super-instructions**: Combine common instruction sequences (e.g., `OpLocal+OpLocal+OpAdd+OpSetLocal` → `OpAddLocals`) into single opcodes. 1.2-1.5x expected improvement for arithmetic loops.

3. **Native bool slice**: Similar to `[]int64`, implement `[]bool` fast-path for Sieve-like workloads that use `[]bool`.

4. **Closure allocation reduction**: ClosureCalls still has 3,000 allocs/op from closure creation. Could pool closures or use a different representation.

5. **Global variable optimization**: `OpGlobal` currently creates a pointer via `FromInterface(&globals[idx])` which allocates. Could use a direct index approach similar to locals.

6. **Instruction encoding optimization**: Current encoding uses 3 bytes per instruction (1 opcode + 2 operand). Could use variable-length encoding to reduce instruction cache pressure.

---

## Testing

All optimizations maintain full backward compatibility. The test suite (`go test ./...`) passes with zero failures:

```
ok  gig/bytecode    0.002s
ok  gig/compiler    0.002s
ok  gig/importer    0.003s
ok  gig/tests       0.852s
ok  gig/value       0.003s
ok  gig/vm          0.003s
```

The benchmark test in `benchmarks/` verifies correct results for all interpreters (Gig, Yaegi, GopherLua, Native Go).
