# Concurrent Global Variables in Gig

This document explains how Gig handles concurrent access to package-level (global) variables, where it matches native Go semantics, and where it differs.

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Two Execution Modes](#two-execution-modes)
- [How Global Variables Work Internally](#how-global-variables-work-internally)
- [Value-Type Struct Globals (sync.Mutex, sync.WaitGroup, etc.)](#value-type-struct-globals)
- [Semantic Comparison with Native Go](#semantic-comparison-with-native-go)
- [Known Differences from Native Go](#known-differences-from-native-go)
- [Best Practices](#best-practices)
- [Implementation Details](#implementation-details)

---

## Quick Start

```go
source := `
package main

import "sync"

var mu sync.Mutex     // value-type — works correctly ✅
var counter int

func Increment() int {
    mu.Lock()
    counter++
    val := counter
    mu.Unlock()
    return val
}
`

// Enable stateful globals for concurrent access
prog, _ := gig.Build(source, gig.WithStatefulGlobals())
defer prog.Close()

// Safe for concurrent use from multiple goroutines
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        prog.Run("Increment")
    }()
}
wg.Wait()

result, _ := prog.Run("GetCounter")
// result == 100 (exact — mutex prevents lost updates)
```

Both `var mu sync.Mutex` and `var mu *sync.Mutex = &sync.Mutex{}` work identically in Gig.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│                    Host Go Program               │
│                                                   │
│   prog.Run("Increment")  prog.Run("Increment")  │
│         │                        │                │
│         ▼                        ▼                │
│   ┌──────────┐            ┌──────────┐           │
│   │   VM 1   │            │   VM 2   │           │
│   │ (pooled) │            │ (pooled) │           │
│   └────┬─────┘            └────┬─────┘           │
│        │                       │                  │
│        ▼                       ▼                  │
│   ┌────────────────────────────────────┐         │
│   │         SharedGlobals              │         │
│   │  ┌──────────────────────────────┐  │         │
│   │  │ [0] int counter              │  │         │
│   │  │ [1] *sync.Mutex (heap ptr)   │  │         │
│   │  └──────────────────────────────┘  │         │
│   │         sync.RWMutex               │         │
│   └────────────────────────────────────┘         │
└─────────────────────────────────────────────────┘
```

In stateful mode, all VMs share a single `SharedGlobals` instance. Global variable reads and writes go through `sync.RWMutex` locks. Value-type struct globals (like `sync.Mutex`) are heap-allocated at compile time so all VMs operate on the same underlying object.

---

## Two Execution Modes

### Default Mode (Stateless)

```go
prog, _ := gig.Build(source)
```

- Each `Run()` starts from the post-`init()` globals snapshot
- Mutations are discarded after each call
- No concurrent state sharing — safe but isolated
- Like calling a pure function each time

### Stateful Mode

```go
prog, _ := gig.Build(source, gig.WithStatefulGlobals())
```

- Global variable mutations persist across `Run()` calls
- Multiple concurrent `Run()` calls are supported
- Global access protected by `sync.RWMutex` (via `SharedGlobals`)
- Like a long-running Go program with shared package-level state

---

## How Global Variables Work Internally

### Storage

All global variables are stored in a `[]value.Value` slice. Each slot holds one global.

| Mode | Storage | Concurrency |
|------|---------|-------------|
| Default | `vm.globals []value.Value` (per-VM copy) | N/A (single VM per call) |
| Stateful | `SharedGlobals.globals []value.Value` (single shared instance) | `sync.RWMutex` protected |

### Access Opcodes

| Opcode | Default Mode | Stateful Mode |
|--------|-------------|---------------|
| `OpGlobal` (read address) | pushes `*value.Value` (raw pointer to slot) | pushes `GlobalRef` (locked proxy) |
| `OpSetGlobal` (write) | direct `globals[idx] = val` | `SharedGlobals.Set()` with write lock |
| `OpDeref` (load via address) | dereference `*value.Value` | `GlobalRef.Load()` with read lock |
| `OpSetDeref` (store via address) | direct write through pointer | `GlobalRef.Store()` with write lock |

### GlobalRef — The Locked Proxy

In stateful mode, `OpGlobal` does NOT expose a raw pointer. Instead it pushes a `GlobalRef`:

```go
type GlobalRef struct {
    sg  *SharedGlobals
    idx int
}

func (r *GlobalRef) Load() value.Value {
    r.sg.mu.RLock()
    v := r.sg.globals[r.idx]
    r.sg.mu.RUnlock()
    return v
}

func (r *GlobalRef) Store(val value.Value) {
    r.sg.mu.Lock()
    r.sg.globals[r.idx] = val
    r.sg.mu.Unlock()
}
```

This ensures every global read/write is atomically protected.

---

## Value-Type Struct Globals

### The Problem (Before Fix)

In native Go, `var mu sync.Mutex` allocates the struct at a fixed memory address. `mu.Lock()` takes `&mu` implicitly and calls the pointer-receiver method on the actual struct.

In Gig, globals are stored in `value.Value` slots. A zero-valued `sync.Mutex{}` would be stored as a non-addressable `reflect.Value` copy. Each method call would see a different copy — `Lock()` would lock copy A, `Unlock()` would try to unlock copy B.

### The Fix: Heap-Allocated Pointer

At compile time, the compiler detects external named struct types and allocates a heap object:

```go
// compiler/compile_value.go
c.globalZeroValues[globalIdx] = reflect.New(rt)  // *sync.Mutex, NOT sync.Mutex
```

This stores a **pointer** (`*sync.Mutex`) in the global slot, not a struct value. All method calls on all goroutines operate on the **same heap-allocated object**.

### What This Means

| User Code | What's Actually Stored in Global Slot |
|-----------|--------------------------------------|
| `var mu sync.Mutex` | `*sync.Mutex` (heap pointer via `reflect.New`) |
| `var mu *sync.Mutex = &sync.Mutex{}` | `*sync.Mutex` (pointer from user code) |
| `var counter int` | `int(0)` (primitive, stored directly) |
| `var m sync.Map` | `*sync.Map` (heap pointer via `reflect.New`) |

Both `var mu sync.Mutex` and `var mu *sync.Mutex` end up with the same runtime representation. Method calls like `mu.Lock()` are dispatched identically.

---

## Semantic Comparison with Native Go

### ✅ Identical Behavior (Verified by Tests)

| Feature | Native Go | Gig (Stateful Mode) | Test Coverage |
|---------|-----------|---------------------|---------------|
| `var mu sync.Mutex; mu.Lock()` | Works | Works | `TestValueTypeMutexExact` — 1000 concurrent increments exact |
| `var rwmu sync.RWMutex` | Works | Works | `TestRWMutex` — 20 writers + 50 readers, 200 exact writes |
| `var m sync.Map` concurrent ops | Works | Works | `TestSyncMap` — 30 goroutines × 20 stores each, all 600 keys present |
| `sync.Once.Do(func(){...})` | Works | Works | `TestOnce` — 100 concurrent calls, all return 42 |
| `sync.WaitGroup` local | Works | Works | `TestWaitGroupGlobal` — 50 goroutines in for-loop, sum=1225 exact |
| Channel-based synchronization | Works | Works | `TestSumViaChannel`, `TestProducerConsumerSum`, `TestComplexSync` |
| Nested lock ordering | Works | Works | `TestNestedLocks` — 50 goroutines × 10 calls, no deadlock |
| Unprotected counter lost updates | Expected | Expected | `TestUnprotectedCounterNoCrash` — same race semantics |
| `init()` function execution | Runs once at startup | Runs once at Build time | Verified by all stateful tests |
| Global zero-value initialization | `int` = 0, `string` = "" | Same | Via `GlobalZeroValues` |

### ⚠️ Behavioral Differences

| Feature | Native Go | Gig (Stateful Mode) | Severity |
|---------|-----------|---------------------|----------|
| **Global slot atomicity** | Each global has independent memory | All globals share one `sync.RWMutex` | Low — coarser granularity, still correct |
| **Read-modify-write on primitives** | Data race (undefined behavior) | No torn reads (RWMutex), but lost updates | Low — Gig is actually safer |
| **`-race` detector** | Available | Not applicable | Info — Gig's RWMutex prevents data races at the slot level |
| **`sync/atomic` operations** | Hardware atomics | Interpreted — slower but correct | Low — semantically correct |
| **Struct field mutation visibility** | Per-field granularity | Per-slot granularity | Low — entire value.Value is swapped |

### Detailed Analysis: Global Slot Locking Granularity

**Native Go**: Each global variable has its own memory address. Two goroutines can write to `counterA` and `counterB` simultaneously with no interference.

**Gig (Stateful Mode)**: All global variables share a single `sync.RWMutex`. A write to `counterA` blocks reads from `counterB`. This is **coarser** than native Go but **never incorrect** — it may cause slight contention under extreme concurrent load.

```
Native Go:     [counterA (own lock)] [counterB (own lock)]  → no contention
Gig Stateful:  [counterA] [counterB]  ← shared RWMutex     → possible contention
```

This is a performance trade-off, not a correctness issue.

### Detailed Analysis: Read-Modify-Write

```go
var counter int
// Goroutine 1:          // Goroutine 2:
counter = counter + 1    counter = counter + 1
```

| Step | Native Go | Gig (Stateful Mode) |
|------|-----------|---------------------|
| Read | No synchronization | `GlobalRef.Load()` (RLock) |
| Compute | No synchronization | Local computation |
| Write | No synchronization | `GlobalRef.Store()` (WLock) |
| **Result** | Data race — undefined behavior | No torn reads, but lost updates possible |

In Gig, each individual read or write is atomic (protected by RWMutex), but the compound read-modify-write is NOT atomic. This matches Go's semantic model where unprotected concurrent read-modify-write causes lost updates — but Gig is actually **safer** because it prevents torn reads/writes.

---

## Known Differences from Native Go

### 1. User-Defined Struct Type Globals

The heap-pointer fix currently applies only to **external named struct types** (types from registered packages like `sync`, `bytes`, etc.). User-defined struct types declared in the interpreted source are compiled by Gig and have different method dispatch paths — they don't hit the `callExternalMethod` code path.

```go
// ✅ Works — external type
var mu sync.Mutex
mu.Lock()

// ✅ Works — user-defined types use compiled method dispatch
type Counter struct { n int }
func (c *Counter) Inc() { c.n++ }
var c Counter
c.Inc()  // dispatched via compiled function table, not external method
```

### 2. No `go vet` copylocks Check

In native Go, `go vet` warns when a `sync.Mutex` is copied. In Gig, the heap-pointer approach means the global slot holds a pointer, so copies of the slot value are pointer copies (safe). However, if user code explicitly copies the struct via assignment:

```go
var mu sync.Mutex
mu2 := mu  // In Gig: copies the pointer — mu2 and mu share the same mutex!
           // In Go:  copies the struct — mu2 is independent
```

This is a **semantic difference**: in native Go, `mu2 := mu` copies the Mutex struct (which `go vet` warns about). In Gig, it copies the pointer, so `mu` and `mu2` alias the same mutex. This is unlikely to be a problem in practice (copying a Mutex is almost always a bug anyway), but it's worth noting.

### 3. Struct Field Access on Value-Type Globals

Since `var mu sync.Mutex` is stored as `*sync.Mutex` internally, accessing the struct's fields (if they were exported) would go through the pointer. For most use cases (method calls), this is transparent. For direct field access patterns, the SSA-compiled `OpField`/`OpFieldAddr` will automatically dereference the pointer.

### 4. ~~sync.Once with Anonymous Functions~~ (FIXED)

`sync.Once.Do()` with anonymous closures now works correctly in concurrent scenarios.
The fix ensures the closure's temporary VM shares the parent's `SharedGlobals`,
so writes to globals inside `Once.Do(func() { globalVar = 42 })` are visible.

### 5. ~~Goroutine Closure Capture in For Loops~~ (FIXED)

Spawning goroutines with closures inside `for` loops now works correctly.
The fix ensures `compileGo` pushes the callee (closure) before arguments
for `OpGoCallIndirect`, matching the VM's stack pop order.

```go
// ✅ All patterns work correctly now
for i := 0; i < N; i++ {
    go func(n int) {
        ch <- n
        wg.Done()
    }(i)
}
```

---

## Best Practices

### For Concurrent Access

```go
// ✅ Both work identically now
var mu sync.Mutex          // value-type — heap-allocated automatically
var mu *sync.Mutex         // pointer-type — explicit allocation

// ✅ Always protect shared mutable state
var mu sync.Mutex
var counter int

func Increment() int {
    mu.Lock()
    counter++
    val := counter
    mu.Unlock()
    return val
}
```

### For Maximum Compatibility with Native Go

```go
// ✅ Use sync primitives the same way you would in native Go
var mu sync.Mutex           // works
var wg sync.WaitGroup       // works
var m sync.Map              // works
var once sync.Once          // works

// ✅ Channel-based patterns work identically
ch := make(chan int, 100)
go func() { ch <- 42 }()
val := <-ch
```

### What to Avoid

```go
// ⚠️ Don't rely on unprotected concurrent writes for exact results
var counter int
// Multiple goroutines doing counter++ without a mutex
// → Lost updates expected (same as native Go)

// ⚠️ Don't copy sync types by assignment
var mu sync.Mutex
mu2 := mu  // Shares the same mutex in Gig (pointer copy)
           // Would be a struct copy in native Go (go vet warns)
```

---

## Implementation Details

### File Map

| File | Responsibility |
|------|---------------|
| `compiler/compile_value.go` | Detects external struct globals, emits `reflect.New(T)` into `GlobalZeroValues`. Fixed `compileGo` stack order for `OpGoCallIndirect`. |
| `compiler/interfaces.go` | `LookupExternalTypeByName()` interface for type resolution |
| `model/bytecode/program.go` | `GlobalZeroValues map[int]reflect.Value` field |
| `vm/vm.go` | `newVM()` and `Reset()` apply `GlobalZeroValues` |
| `vm/shared_globals.go` | `SharedGlobals`, `GlobalRef`, `InitZeroValues()` |
| `vm/call.go` | `callExternalMethod()` resolves GlobalRef → loads stored pointer |
| `vm/closure.go` | `Closure` struct with `Shared`/`Goroutines`/`ExtCallCache`/`Ctx` for SharedGlobals propagation to `reflect.MakeFunc` closures |
| `vm/ops_memory.go` | `OpGlobal` pushes GlobalRef (shared) or `*value.Value` (default) |
| `vm/ops_call.go` | `OpClosure` populates `Shared`/`Goroutines` on closures; `OpGoCall`/`OpGoCallIndirect` goroutine spawning |
| `runner/runner.go` | Creates `SharedGlobals`, calls `InitZeroValues()` |
| `importer/interfaces.go` | `LookupExternalTypeByName()` implementation |

### Data Flow: `var mu sync.Mutex; mu.Lock()`

```
Compile Time:
  1. Compiler sees `var mu sync.Mutex` (SSA Global, type *sync.Mutex)
  2. Extracts elem type: sync.Mutex (Named, pkg="sync", name="Mutex")
  3. LookupExternalTypeByName("sync", "Mutex") → reflect.TypeOf(sync.Mutex{})
  4. reflect.New(sync.Mutex) → *sync.Mutex (heap pointer)
  5. Stores in GlobalZeroValues[idx] = reflect.New(sync.Mutex)

Init Time:
  6. newVM() copies GlobalZeroValues[idx] into globals[idx]
  7. SSA init may store nil over it (zero struct const)
  8. GlobalZeroValues loop replaces nil with heap pointer

Runtime (mu.Lock()):
  9.  OpGlobal idx → pushes GlobalRef{sg, idx}
  10. OpCallExternal (Mutex.Lock) → pops GlobalRef as args[0]
  11. callExternalMethod detects GlobalRef → args[0] = ref.Load()
  12. ref.Load() returns value.Value wrapping *sync.Mutex (heap pointer)
  13. DirectCall: args[0].Interface().(*sync.Mutex) → the actual heap mutex
  14. recv.Lock() → locks the real mutex object
  15. No write-back needed — the heap object is mutated in place
```

### Concurrency Safety Proof

For `var mu sync.Mutex` under concurrent access:

1. **Single heap object**: `reflect.New(sync.Mutex)` allocates exactly one `sync.Mutex` on the heap at compile time. This pointer is stored in the global slot.

2. **Pointer identity**: `GlobalRef.Load()` returns the `value.Value` wrapping this pointer. `value.Value.Interface()` returns `*sync.Mutex` — always the **same pointer** to the **same heap object**.

3. **No copy**: `callExternalMethod` does NOT call `reflect.New` or `Set` at runtime. It loads the existing pointer and passes it to the DirectCall wrapper.

4. **Method operates in-place**: `recv.Lock()` / `recv.Unlock()` operates directly on the heap-allocated `sync.Mutex`. The `sync.Mutex` internal state (`state int32`, `sema uint32`) is modified atomically by the Go runtime.

5. **SharedGlobals RWMutex**: Protects the `value.Value` slot read/write, NOT the mutex itself. The `sync.Mutex` has its own internal synchronization.

**Result**: Multiple goroutines calling `mu.Lock()` / `mu.Unlock()` all operate on the same `sync.Mutex` object, achieving true mutual exclusion — identical to native Go.
