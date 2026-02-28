# Superinstruction Optimization — Performance Report

## Summary

This optimization introduces **peephole superinstructions** and **register-hoisted dispatch** to the Gig VM, targeting tight loop performance where Yaegi previously held a 4-5x advantage.

### Results (Cross-Interpreter Benchmarks)

| Workload | Before (μs) | After (μs) | Speedup | vs Yaegi Before | vs Yaegi After |
|---|---|---|---|---|---|
| Fibonacci(25) | 46,550 | 32,045 | **1.45x** | Gig 2.4x faster | Gig 3.5x faster |
| ArithmeticSum(1K) | 200 | 136 | **1.47x** | Yaegi 5.0x faster | Yaegi 3.3x faster |
| BubbleSort(100) | 4,979 | 3,904 | **1.28x** | Yaegi 4.0x faster | Yaegi 3.1x faster |
| Sieve(1000) | 840 | 687 | **1.22x** | Yaegi 4.1x faster | Yaegi 3.3x faster |
| ClosureCalls(1K) | 586 | 528 | **1.11x** | Gig 1.7x faster | Gig 1.9x faster |

### Full Comparison Table (After Optimization)

| Workload | Native Go | Gig | Yaegi | GopherLua | Gig/Native | Yaegi/Native | Gig vs Yaegi |
|---|---|---|---|---|---|---|---|
| Fibonacci(25) | 453 μs | 32.0 ms | 111 ms | 21.0 ms | 71x | 244x | **Gig 3.5x faster** |
| ArithmeticSum(1K) | 664 ns | 136 μs | 41 μs | 41 μs | 205x | 62x | Yaegi 3.3x faster |
| BubbleSort(100) | 6.4 μs | 3.90 ms | 1.26 ms | 774 μs | 609x | 197x | Yaegi 3.1x faster |
| Sieve(1000) | 1.88 μs | 687 μs | 209 μs | 212 μs | 366x | 111x | Yaegi 3.3x faster |
| ClosureCalls(1K) | 347 ns | 528 μs | 1,009 μs | 122 μs | 1,522x | 2,908x | **Gig 1.9x faster** |

## Changes Made

### 1. Peephole Optimizer (`compiler/optimize.go`)

A post-compilation pass that scans for common bytecode patterns and fuses them into single superinstructions. Runs after jump target patching. The optimizer:

- Scans the instruction stream for known patterns
- Replaces multi-instruction sequences with shorter fused opcodes
- Rebuilds the bytecode and fixes all jump targets to account for size changes

### 2. New Superinstructions (`bytecode/opcode.go`)

17 new fused opcodes:

**Fused arithmetic + store (eliminate stack traffic entirely):**
- `OpLocalLocalAddSetLocal` — `local[A] + local[B] → local[C]` (zero stack ops)
- `OpLocalConstAddSetLocal` — `local[A] + const[B] → local[C]`
- `OpLocalConstSubSetLocal` — `local[A] - const[B] → local[C]`

**Fused load + arithmetic (eliminate 2 push ops):**
- `OpAddLocalLocal` — push `local[A] + local[B]`
- `OpSubLocalLocal` — push `local[A] - local[B]`
- `OpMulLocalLocal` — push `local[A] * local[B]`
- `OpAddLocalConst` — push `local[A] + const[B]`
- `OpSubLocalConst` — push `local[A] - const[B]`

**Fused compare-and-branch (eliminate bool push/pop + separate jump):**
- `OpLessLocalLocalJumpTrue` — jump if `local[A] < local[B]`
- `OpLessLocalConstJumpTrue` — jump if `local[A] < const[B]`
- `OpLessEqLocalConstJumpTrue` — jump if `local[A] <= const[B]`
- `OpGreaterLocalLocalJumpTrue` — jump if `local[A] > local[B]`
- `OpLessLocalLocalJumpFalse` — jump if `local[A] >= local[B]`
- `OpLessLocalConstJumpFalse` — jump if `local[A] >= const[B]`

**Fused arithmetic + store (pop 2, compute, store):**
- `OpAddSetLocal` — pop a,b; store `a+b` to `local[A]`
- `OpSubSetLocal` — pop a,b; store `a-b` to `local[A]`

### 3. Register-Hoisted Dispatch (`vm/run.go`)

The main execution loop now hoists critical VM fields into Go local variables:

```go
stack := vm.stack     // stack slice header in register
sp := vm.sp           // stack pointer in register
prebaked := vm.program.PrebakedConstants
```

And all hot-path opcodes operate on these locals directly:
```go
// Before: vm.push(frame.locals[idx]) → method call + indirect access
// After: stack[sp] = locals[idx]; sp++ → direct register-based access
```

This eliminates per-instruction overhead from:
- Method call overhead for push/pop (Go can't inline all method calls)
- Repeated loads of `vm.stack`, `vm.sp` from memory
- Stack growth check on every push (now only checked at function boundaries)

The locals are synced back to `vm.*` fields before any operation that might read them (executeOp, callFunction, context check).

### 4. Integration (`compiler/compile_func.go`)

The optimizer is called automatically after jump patching:
```go
c.patchJumps(blockOffsets)
c.currentFunc.Instructions = optimizeBytecode(c.currentFunc.Instructions)
```

## Why the Improvement

### Before: ArithSum inner loop (`sum += i; i++; i < 1000`)
```
12 bytecodes per iteration:
  OpLocal(sum) OpLocal(i) OpAdd OpSetLocal(sum)    — 4 dispatches, 6 stack ops
  OpLocal(i) OpConst(1) OpAdd OpSetLocal(i)        — 4 dispatches, 6 stack ops  
  OpLocal(i) OpConst(1000) OpLess OpJumpTrue       — 4 dispatches, 6 stack ops
Total: 12 dispatches, 18 push/pop of 48-byte Values per iteration
```

### After: ArithSum inner loop
```
3 superinstructions per iteration:
  OpLocalLocalAddSetLocal(sum, i, sum)              — 1 dispatch, 0 stack ops
  OpLocalConstAddSetLocal(i, 1, i)                  — 1 dispatch, 0 stack ops
  OpLessLocalConstJumpTrue(i, 1000, target)         — 1 dispatch, 0 stack ops
Total: 3 dispatches, 0 push/pop per iteration (4x fewer dispatches, zero stack traffic)
```

## Architecture Comparison Update

| Dimension | Gig (Before) | Gig (After) | Yaegi |
|---|---|---|---|
| Dispatches per loop iter | 12 | 3 | 3-4 |
| Stack ops per loop iter | ~18 push/pop | 0 | 0 |
| Data moved per iter | ~1,536 bytes | ~0 bytes | ~24 bytes |
| Dispatch method | switch on vm.* | switch on local vars | closure chain |
| Remaining gap cause | 48-byte Value struct | 48-byte Value struct | 24-byte reflect.Value |

## Remaining Gap Analysis

Gig is still ~3.3x slower than Yaegi on pure arithmetic loops. The remaining gap is:

1. **Value struct size (48 bytes vs 24 bytes)** — Even with zero stack ops, locals[] still moves 48-byte structs. Yaegi's reflect.Value is 24 bytes.
2. **GC pressure from `obj any` field** — The `Value.obj` field (type `any`, 16 bytes) forces the GC to scan all Value slices. Yaegi uses reflect.Value which has the same issue but at half the memory.
3. **Closure-threaded vs switch dispatch** — Yaegi's closure-threaded execution has better branch prediction. Each closure "knows" its next closure, while switch dispatch is an indirect branch every time.

### Future Optimization Priorities

| Priority | Optimization | Expected Impact |
|---|---|---|
| P0 | Shrink Value to 24 bytes (NaN-boxing or split union) | ~1.5-2x on arithmetic loops |
| P1 | Direct-threaded dispatch (computed goto via assembly) | ~1.2-1.5x on all workloads |
| P2 | Integer-specialized locals for hot functions | ~1.3x on int-heavy loops |
