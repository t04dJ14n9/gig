# Gig Internals

A technical deep-dive into the Gig Go interpreter for engineers who need to understand
how the system works, extend it, or debug production issues.

---

## Why Gig? The Rule Engine Problem

Every business has logic that changes faster than deployment cycles allow. Pricing
rules, eligibility checks, fraud scoring, promotion matching — these shift weekly
or even daily. The standard solutions each have a deal-breaking tradeoff:

| Approach | Problem |
|---|---|
| Hardcoded Go | Requires recompile + redeploy for every rule change |
| Expression languages (CEL, Rego) | Limited power: no loops, no stdlib, new syntax to learn |
| Lua/JS embedding | Different language: Go developers must context-switch, can't reuse Go libraries |
| gRPC microservices | Operational overhead: deploy, version, monitor a separate service per rule |

**What we actually want**: write rules in Go (zero learning cost), call any Go stdlib
or third-party library (full power), but load and execute them dynamically without
recompiling the host application.

### Gig's Approach

Gig is a **full Go interpreter** that compiles Go source code to bytecode and executes
it on a stack-based virtual machine. It is not a subset or a DSL — it handles the
complete Go language including goroutines, closures, defer/panic/recover, interfaces,
methods, and type assertions.

```
┌─────────────┐     ┌──────────┐     ┌────────────┐     ┌──────────┐     ┌──────────┐
│  Go Source   │────▶│  Parser  │────▶│ Type Check │────▶│ SSA Build│────▶│ Compiler │
│  (string)    │     │ go/parser│     │ go/types   │     │ go/ssa   │     │ SSA→BC   │
└─────────────┘     └──────────┘     └────────────┘     └──────────┘     └─────┬────┘
                                                                               │
                    ┌──────────┐     ┌────────────┐     ┌──────────┐           │
                    │  Result  │◀────│    VM       │◀────│ Bytecode │◀──────────┘
                    │  (any)   │     │ stack-based │     │ Program  │
                    └──────────┘     └────────────┘     └──────────┘
```

### Basic Usage

```go
package main

import (
    "fmt"
    "git.woa.com/youngjin/gig"
)

func main() {
    // Compile Go source code to bytecode
    prog, err := gig.Build(`
        import "strings"

        func ProcessName(name string) string {
            return strings.ToUpper(strings.TrimSpace(name))
        }
    `)
    if err != nil {
        panic(err)
    }
    defer prog.Close()

    // Execute the compiled function
    result, err := prog.Run("ProcessName", "  hello world  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: HELLO WORLD
}
```

The `Build` call takes ~1-5ms (parse + type-check + SSA + compile). Each `Run` call
takes microseconds — the bytecode is already compiled and the VM is pooled.

---

## Architecture Overview

### Package Structure

```
gig/
├── gig.go                    # Public API: Build(), Program.Run()
├── compiler/
│   ├── build.go              # Full pipeline: source → parse → SSA → bytecode
│   ├── compiler.go           # SSA → bytecode translation
│   ├── compile_func.go       # Per-function compilation
│   ├── compile_instr.go      # Per-instruction compilation
│   ├── symbol.go             # Symbol table (SSA value → local slot)
│   ├── parser/               # go/parser + security validation
│   ├── ssa/                  # go/ssa builder wrapper
│   ├── peephole/             # Pattern-based superinstruction fusion
│   └── optimize/             # 4-pass bytecode optimization pipeline
├── vm/
│   ├── vm.go                 # VM struct, Execute() entry point
│   ├── run.go                # Main fetch-decode-execute loop (hot path)
│   ├── frame.go              # Call frame + frame pool
│   ├── stack.go              # Operand stack with bounded growth
│   ├── call.go               # External function calls (DirectCall + reflect)
│   ├── closure.go            # Closure type + ClosureExecutor
│   ├── goroutine.go          # GoroutineTracker, child VM construction
│   ├── ops_dispatch.go       # Opcode routing to category handlers
│   ├── ops_arithmetic.go     # Arithmetic, bitwise, comparison ops
│   ├── ops_memory.go         # Stack, locals, globals, fields, addresses
│   ├── ops_container.go      # Slice, map, chan operations
│   ├── ops_control.go        # Control flow, defer, panic/recover
│   ├── ops_convert.go        # Type assertions, conversions
│   └── ops_call.go           # Function/closure calls, goroutines
├── model/
│   ├── value/                # 32-byte tagged-union Value type
│   ├── bytecode/             # CompiledProgram, CompiledFunction, OpCode
│   └── external/             # ExternalFuncInfo, ExternalMethodInfo
├── importer/                 # Package registration, type resolution
├── runner/                   # VM pool, init execution, stateful globals
└── stdlib/packages/          # ~69 pre-generated stdlib wrappers
```

### Build Pipeline in Detail

```go
// gig.go: Build()
func Build(sourceCode string, opts ...BuildOption) (*Program, error) {
    // 1. compiler.Build: source → parse → SSA → bytecode
    result, err := compiler.Build(sourceCode, cfg.registry, compilerOpts...)

    // 2. Execute init() and snapshot globals
    initialGlobals, err := runner.ExecuteInit(result.Program)

    // 3. Create runner (owns VM pool)
    r := runner.New(result.Program, initialGlobals, runnerOpts...)

    return &Program{runner: r, ssaPkg: result.SSAPkg}, nil
}
```

The `compiler.Build` function orchestrates three phases:

```go
// compiler/build.go
func Build(source string, reg importer.PackageRegistry, opts ...BuildOption) (*BuildResult, error) {
    // Phase 1: Parse + type-check + validate
    parseResult, err := parser.Parse(source, reg, parseOpts...)

    // Phase 2: Build SSA
    ssaResult, err := ssabuilder.Build(parseResult.FSet, parseResult.Pkg, ...)

    // Phase 3: Compile SSA → bytecode
    lookup := importer.NewPackageLookup(reg)
    compiled, err := NewCompiler(lookup).Compile(ssaResult.Pkg)

    return &BuildResult{Program: compiled, SSAPkg: ssaResult.Pkg}, nil
}
```

---

## The Value System

### The Problem

An interpreter needs a universal type to represent any Go value at runtime: `int`,
`string`, `[]byte`, `*http.Request`, closures, etc. The naive approach is `interface{}`,
but in an interpreter's operand stack, intermediate values constantly flow through
push/pop — escape analysis can't optimize this, so most arithmetic results heap-allocate.
For numeric-heavy rule engines, this means massive GC pressure.

### The Solution: 32-Byte Tagged Union

```
┌──────────┬──────────┬───────────┬────────────────────────────────────┐
│ kind (1B)│ size (1B)│ pad (6B)  │          num (8B)                  │
├──────────┴──────────┴───────────┼────────────────────────────────────┤
│                                 │          obj (16B)                  │
│         (interface{})           │    string / reflect.Value / etc.    │
└─────────────────────────────────┴────────────────────────────────────┘
                            Total: 32 bytes
```

```go
// model/value/value.go
type Value struct {
    kind Kind    // 1 byte: type tag (KindInt, KindString, KindFloat, ...)
    size Size    // 1 byte: original Go bit-width (8, 16, 32, 64)
    num  int64   // 8 bytes: stores bool (0/1), int, uint bits, float64 bits
    obj  any     // 16 bytes: string, complex128, reflect.Value, or nil
}
```

**The key insight**: for primitive types (`bool`, `int`, `uint`, `float64`, `nil`),
everything lives in `kind` + `num`. The `obj` field is `nil`. **Zero heap allocation,
zero GC pressure.**

### Kind Types

```go
const (
    KindInvalid Kind = iota  // 0 — zero value of Value (uninitialized globals)
    KindNil                  // 1 — explicit nil
    KindBool                 // 2 — stored in num: 0=false, 1=true
    KindInt                  // 3 — stored in num as int64
    KindUint                 // 4 — stored in num as uint64 bits
    KindFloat                // 5 — stored in num as float64 bits (math.Float64bits)
    KindString               // 6 — stored in obj as Go string
    KindComplex              // 7 — stored in obj as complex128
    KindPointer              // 8
    KindSlice                // 9
    KindArray                // 10
    KindMap                  // 11
    KindChan                 // 12
    KindFunc                 // 13
    KindStruct               // 14
    KindInterface            // 15
    KindReflect              // 16 — fallback: reflect.Value in obj
    KindBytes                // 17 — native []byte (avoids reflect overhead)
)
```

### Constructors

```go
// Primitives: zero allocation
value.MakeInt(42)         // kind=KindInt, num=42, obj=nil
value.MakeFloat(3.14)     // kind=KindFloat, num=float64bits(3.14), obj=nil
value.MakeBool(true)      // kind=KindBool, num=1, obj=nil
value.MakeNil()           // kind=KindNil, num=0, obj=nil

// Strings: obj holds the Go string
value.MakeString("hello") // kind=KindString, num=0, obj="hello"

// Composites: obj holds reflect.Value or native Go type
value.MakeBytes([]byte{1,2,3})  // kind=KindBytes, obj=[]byte{1,2,3}
value.MakeIntSlice([]int64{...}) // kind=KindSlice, obj=[]int64{...}
value.FromInterface(anyValue)    // auto-detect: fast-path type switch
```

### Size Tag: Preserving Original Types

Go has `int8`, `int16`, `int32`, `int64`, `int` — all stored as `int64` internally.
The `size` field remembers which one it was, so `Interface()` returns the correct type:

```go
value.MakeInt8(42).Interface()  // returns int8(42), not int64(42)
value.MakeInt32(42).Interface() // returns int32(42)
value.MakeInt(42).Interface()   // returns int(42)
```

This matters when passing values to external Go functions via reflection — if you
call `strings.Repeat(s, n)` and `n` is `int64` instead of `int`, `reflect.Call` panics.

### Why Not `interface{}`?

Consider a Fibonacci benchmark computing Fib(25) = 75,025:

| Value representation | Allocations | Reason |
|---|---|---|
| `interface{}` for all values | ~2.1M | Intermediate values escape to heap in interpreter stack |
| `value.Value` tagged union | ~7 | Only the initial frame + stack allocation |

The 32-byte `Value` fits in two cache lines and never escapes to the heap for
primitive operations.

---

## Compilation: From Go Source to Bytecode

### A Concrete Example

Let's trace this function through the entire compilation pipeline:

```go
func ProcessOrder(price float64, quantity int, coupon string) (float64, bool) {
    total := price * float64(quantity)
    if coupon == "HALF" {
        total *= 0.5
    }
    valid := total > 0 && total < 10000
    return total, valid
}
```

### Phase 1: Parsing and Type Checking

The parser does three things:

1. **`go/parser.ParseFile`** produces an AST
2. **Security validation**: checks for banned imports (`unsafe`, `reflect`) and banned
   `panic()` usage (configurable via `WithAllowPanic()`)
3. **Auto-import**: if the source references `strings.Contains(...)` but doesn't
   have `import "strings"`, the parser adds it automatically (since registered
   packages are known)
4. **`go/types.Config.Check`**: type-checks with a custom `types.Importer` that
   resolves packages against the registry

### Phase 2: SSA Construction

`golang.org/x/tools/go/ssa` converts the typed AST into Static Single Assignment form.
For our example, the SSA looks approximately like this:

```
func ProcessOrder(price float64, quantity int, coupon string) (float64, bool):
  entry:                                              ; block 0
    t0 = Convert quantity int → float64               ; float64(quantity)
    t1 = Mul price t0                                 ; price * float64(quantity)
    t2 = BinOp coupon == "HALF"                       ; coupon == "HALF"
    If t2 → if.then, if.done

  if.then:                                            ; block 1
    t3 = Mul t1 0.5:float64                           ; total *= 0.5
    Jump → if.done

  if.done:                                            ; block 2
    t4 = Phi [entry: t1, if.then: t3]                 ; total (merge)
    t5 = BinOp t4 > 0:float64                         ; total > 0
    If t5 → and.rhs, and.done(false)

  and.rhs:                                            ; block 3
    t6 = BinOp t4 < 10000:float64                     ; total < 10000
    Jump → and.done

  and.done:                                           ; block 4
    t7 = Phi [if.done: false, and.rhs: t6]            ; valid (short-circuit &&)
    t8 = MakeResult t4 t7                             ; (total, valid)
    Return t8
```

**Key SSA concepts**:
- Every value is assigned exactly once (SSA property)
- **Phi nodes** merge values from different control paths (e.g., `t4` picks `t1` or `t3`)
- Short-circuit `&&` becomes explicit control flow with a Phi

### Phase 3: Bytecode Generation

The compiler translates SSA instructions to stack-based bytecode. Here's the
compilation process:

#### Symbol Table Construction

First, the compiler allocates local slots:

```
Slot 0: price     (parameter)
Slot 1: quantity  (parameter)
Slot 2: coupon    (parameter)
Slot 3: t0        (float64(quantity))
Slot 4: t1        (price * float64(quantity))
Slot 5: t4        (phi: total after merge)
Slot 6: t3        (total * 0.5)
Slot 7: t5        (total > 0)
Slot 8: t6        (total < 10000)
Slot 9: t7        (phi: valid)
```

Parameters get the first slots. Phi nodes get dedicated slots (the compiler emits
explicit `SetLocal` moves to resolve phis). Temporaries get the remaining slots.

#### Generated Bytecode

```
; entry block
0000: LOCAL    0         ; push price
0003: LOCAL    1         ; push quantity
0006: CONVERT  [float64] ; convert int → float64, result on stack
0009: SETLOCAL 3         ; store t0 = float64(quantity)
000C: LOCAL    0         ; push price
000F: LOCAL    3         ; push t0
0012: MUL                ; price * float64(quantity)
0013: SETLOCAL 4         ; store t1
0016: SETLOCAL 5         ; phi pre-copy: total = t1 (entry path)
0019: LOCAL    2         ; push coupon
001C: CONST    0         ; push "HALF" from constant pool
001F: EQUAL              ; coupon == "HALF"
0020: JUMPFALSE 002E     ; skip to if.done if false

; if.then block
0023: LOCAL    4         ; push t1
0026: CONST    1         ; push 0.5
0029: MUL                ; t1 * 0.5
002A: SETLOCAL 6         ; store t3
002D: SETLOCAL 5         ; phi pre-copy: total = t3 (if.then path)

; if.done block (total is in slot 5 via phi resolution)
002E: LOCAL    5         ; push total
0031: CONST    2         ; push 0.0
0034: GREATER            ; total > 0
0035: JUMPFALSE 0042     ; short-circuit: if false, skip to and.done

; and.rhs block
0038: LOCAL    5         ; push total
003B: CONST    3         ; push 10000.0
003E: LESS               ; total < 10000
003F: SETLOCAL 9         ; store valid
0040: JUMP     0045      ; jump to return

; and.done (false path)
0042: FALSE              ; push false
0043: SETLOCAL 9         ; valid = false

; return
0045: LOCAL    5         ; push total
0048: LOCAL    9         ; push valid
004B: PACK     2         ; pack into multi-return tuple
004D: RETURNVAL          ; return
```

#### Phi Elimination

SSA Phi nodes don't map to hardware instructions. The compiler eliminates them by:
1. Allocating a dedicated local slot for each Phi (slot 5 for `total`, slot 9 for `valid`)
2. At the end of each predecessor block, emitting `SETLOCAL` to write the correct
   value into the Phi slot

This is why you see `SETLOCAL 5` in both the entry block and the if.then block —
each writes its version of `total` into the shared phi slot before control reaches
the merge point.

### Phase 4: Optimization

After initial bytecode generation, four optimization passes run:

```go
// compiler/optimize/optimize.go
func Optimize(code []byte, localIsInt, constIsInt, localIsIntSlice []bool) ([]byte, bool) {
    code = Peephole(code)                                    // Pass 1: superinstruction fusion
    code = FuseSliceOps(code, localIsInt, localIsIntSlice)   // Pass 2: slice op fusion
    code, hasInt := IntSpecialize(code, localIsInt, constIsInt) // Pass 3: int specialization
    code = FuseIntMoves(code)                                // Pass 4: move fusion
    return code, hasInt
}
```

#### Pass 1: Peephole — Superinstruction Fusion

Peephole patterns detect common multi-instruction sequences and replace them with
a single fused opcode. For example:

```
Before (3 instructions, 3 dispatch cycles):
    LOCAL    0        ; push a
    LOCAL    1        ; push b
    ADD               ; a + b

After (1 instruction, 1 dispatch cycle):
    ADDLOCALLOCAL 0 1  ; push locals[0] + locals[1]
```

The peephole optimizer has 17+ pattern rules covering:

| Pattern | Fused Opcode | Saves |
|---|---|---|
| `LOCAL a` + `LOCAL b` + `ADD` | `OpAddLocalLocal a b` | 2 dispatches |
| `LOCAL a` + `CONST c` + `ADD` | `OpAddLocalConst a c` | 2 dispatches |
| `ADD` + `SETLOCAL x` | `OpAddSetLocal x` | 1 dispatch |
| `LOCAL a` + `LOCAL b` + `ADD` + `SETLOCAL c` | `OpLocalLocalAddSetLocal a b c` | 3 dispatches |
| `LOCAL a` + `LOCAL b` + `LESS` + `JUMPTRUE off` | `OpLessLocalLocalJumpTrue a b off` | 3 dispatches |

Each pattern is registered in the global pattern registry:

```go
// compiler/peephole/pattern.go
type Pattern interface {
    Match(code []byte, i int) (consumed int, newBytes []byte, ok bool)
}
```

#### Pass 2: Slice Operation Fusion

Detects the common pattern of `LOCAL(slice)` + `LOCAL(index)` + `INDEXADDR` +
`SETLOCAL(ptr)` + `LOCAL(ptr)` + `DEREF` + `SETLOCAL(val)` and fuses it into
a single `OpIntSliceGet slice index val` when all types are `int`.

#### Pass 3: Integer Specialization

When the compiler can prove all locals in a superinstruction are `int`-typed, it
upgrades the generic superinstruction to an `OpInt*` variant:

```
OpLocalConstAddSetLocal → OpIntLocalConstAddSetLocal
OpLessLocalLocalJumpFalse → OpIntLessLocalLocalJumpFalse
```

The `OpInt*` variants operate on a shadow `intLocals []int64` array — pure 8-byte
int64 arithmetic instead of 32-byte Value operations:

```go
// vm/run.go — OpIntLocalConstAddSetLocal handler
case bytecode.OpIntLocalConstAddSetLocal:
    idxA := readU16()
    idxB := readU16()
    idxC := readU16()
    r := intLocals[idxA] + intConsts[idxB]   // raw int64 add — no kind check
    intLocals[idxC] = r                       // write to int shadow array
    locals[idxC] = value.MakeInt(r)           // keep Value array in sync
    continue
```

This gives 4x better cache utilization: `int64` is 8 bytes vs `Value`'s 32 bytes,
so 4x more operands fit in the same cache line.

#### Pass 4: Move Fusion

Replaces `OpIntLocal(src)` + `OpIntSetLocal(dst)` with `OpIntMoveLocal(src, dst)`,
eliminating the stack round-trip.

---

## The Virtual Machine

### VM Structure

```go
// vm/vm.go
type vm struct {
    program        *bytecode.CompiledProgram  // compiled bytecode
    stack          []value.Value              // operand stack
    sp             int                        // stack pointer
    frames         []*Frame                   // call frame stack
    fp             int                        // frame pointer
    globals        []value.Value              // package-level variables
    globalsPtr     *[]value.Value             // shared globals (goroutines)
    ctx            context.Context            // cancellation/timeout
    panicking      bool                       // panic in progress?
    panicVal       value.Value                // current panic value
    panicStack     []panicState               // saved panics (nested)
    deferDepth     int                        // defer nesting level
    extCallCache   *externalCallCache         // inline cache for ext calls
    initialGlobals []value.Value              // post-init globals snapshot
    goroutines     *GoroutineTracker          // goroutine limiter
    fpool          framePool                  // frame recycler
}
```

### Key Constants

```go
const (
    initialStackSize     = 1024     // starting operand stack slots
    maxStackSize         = 1 << 20  // 1M slots = 32 MB per VM
    initialFrameDepth    = 64       // starting call frame slots
    maxFrameDepth        = 1024     // max call depth (8 KB)
    contextCheckInterval = 1024     // check ctx.Done() every N instructions
    defaultMaxGoroutines = 10000    // goroutine limit per program
)
```

### Stack Management

The operand stack is a `[]value.Value` slice that doubles in size when full:

```go
// vm/stack.go
func (v *vm) push(val value.Value) {
    if v.sp >= len(v.stack) {
        if len(v.stack) >= maxStackSize {
            panic("gig: operand stack overflow")
        }
        newCap := len(v.stack) * 2
        if newCap > maxStackSize {
            newCap = maxStackSize
        }
        newStack := make([]value.Value, newCap)
        copy(newStack, v.stack)
        v.stack = newStack
    }
    v.stack[v.sp] = val
    v.sp++
}

func (v *vm) pop() value.Value {
    v.sp--
    return v.stack[v.sp]
}
```

Note: the panic from stack overflow is caught by the safety net (a `defer/recover`
wrapping `vm.run()`), so it becomes an error return — never crashing the host process.

### Frame Management and Pooling

Each function call creates a `Frame`:

```go
// vm/frame.go
type Frame struct {
    fn        *bytecode.CompiledFunction  // which function
    ip        int                         // instruction pointer
    basePtr   int                         // operand stack base for this call
    locals    []value.Value               // local variable array
    intLocals []int64                     // int-specialized shadow array
    freeVars  []*value.Value              // closure captures (pointers!)
    defers    []DeferInfo                 // deferred calls (LIFO)
    addrTaken bool                        // true if OpAddr pointed into locals
}
```

#### Frame Pooling

Without pooling, every function call allocates a `Frame` + `[]value.Value` on the
heap. For recursive functions like Fibonacci, this means millions of allocations.

The `framePool` is a simple LIFO stack that recycles frames:

```go
// vm/frame.go
func (p *framePool) get(fn *bytecode.CompiledFunction, basePtr int, freeVars []*value.Value) *Frame {
    n := len(p.frames)
    if n > 0 {
        f = p.frames[n-1]
        p.frames = p.frames[:n-1]
        // Reuse the locals slice if it has enough capacity
        if cap(f.locals) >= fn.NumLocals {
            f.locals = f.locals[:fn.NumLocals]
            for i := range f.locals {
                f.locals[i] = value.Value{}  // zero out for correctness
            }
        } else {
            f.locals = make([]value.Value, fn.NumLocals)
        }
        // ... reset all other fields ...
    } else {
        f = &Frame{locals: make([]value.Value, fn.NumLocals)}
    }
    return f
}

func (p *framePool) put(f *Frame) {
    if f.addrTaken {
        return  // closures may hold live references — don't recycle
    }
    f.fn = nil
    f.freeVars = nil
    p.frames = append(p.frames, f)
}
```

**The `addrTaken` guard**: when `OpAddr` creates a pointer into a frame's locals,
closures or deferred functions may hold references to those slots. Recycling the
frame would corrupt those pointers. So `addrTaken` frames are left for the GC.

**Result**: Fib(25) drops from ~728K allocations to **7 allocations**. The frame
pool absorbs all recursion overhead.

### The Dispatch Loop

The dispatch loop in `vm/run.go` is the performance-critical hot path. It uses
several techniques to maximize throughput:

#### Register Hoisting

```go
func (v *vm) run() (value.Value, error) {
    // Hoist hot fields into local variables for register allocation.
    // The Go compiler keeps these in CPU registers across iterations.
    stack := v.stack
    sp := v.sp
    prebaked := v.program.PrebakedConstants

    var frame *Frame
    var ins []byte
    var locals []value.Value
    var intLocals []int64
    intConsts := v.program.IntConstants

    loadFrame := func() {
        frame = v.frames[v.fp-1]
        ins = frame.fn.Instructions
        locals = frame.locals
        intLocals = frame.intLocals
    }

    readU16 := func() uint16 {
        v := uint16(ins[frame.ip])<<8 | uint16(ins[frame.ip+1])
        frame.ip += 2
        return v
    }
    // ...
```

By copying `v.stack`, `v.sp`, `frame.locals`, etc. into local variables, the Go
compiler can keep them in CPU registers. Without this, every instruction would
dereference `v.stack[v.sp]` through two pointer indirections.

#### Hot/Cold Split

The main `switch` in `run()` inlines **~60 hot opcodes** directly. Less frequent
opcodes fall through to `executeOp()`:

```go
    for v.fp > 0 {
        op := bytecode.OpCode(ins[frame.ip])
        frame.ip++

        switch op {
        case bytecode.OpLocal:       // INLINED — hot
            idx := readU16()
            stack[sp] = locals[idx]
            sp++
            continue

        case bytecode.OpAdd:         // INLINED — hot
            sp--
            b := stack[sp]
            sp--
            a := stack[sp]
            if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
                stack[sp] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
            } else {
                stack[sp] = a.Add(b)
            }
            sp++
            continue

        // ... ~58 more hot opcodes ...

        default:
            // Cold path: sync sp, call executeOp
            v.sp = sp
            if err := v.executeOp(op, frame); err != nil {
                return value.MakeNil(), err
            }
            sp = v.sp
            stack = v.stack
        }
    }
```

#### Integer Fast Paths

Every arithmetic and comparison operation checks if both operands are `KindInt`
first. When they are (the common case in rule engines), it does raw `int64` math
with zero overhead:

```go
case bytecode.OpLess:
    sp--
    b := stack[sp]
    sp--
    a := stack[sp]
    if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
        stack[sp] = value.MakeBool(a.RawInt() < b.RawInt())  // fast path
    } else {
        stack[sp] = value.MakeBool(a.Cmp(b) < 0)             // generic path
    }
    sp++
    continue
```

#### Context Checking

The VM checks for cancellation every 1024 instructions using a bitmask:

```go
    instructionCount++
    if instructionCount & 0x3FF == 0 {  // contextCheckMask = 1023
        select {
        case <-v.ctx.Done():
            return value.MakeNil(), v.ctx.Err()
        default:
        }
    }
```

The bitmask avoids a modulo operation, and `select` with `default` is a non-blocking
channel check.

### Execution Walkthrough

Let's trace the execution of `ProcessOrder(100.0, 3, "HALF")` step by step.
Only showing the key operations (omitting some SetLocals for brevity):

```
Frame: ProcessOrder
  locals[0] = Float(100.0)    ; price
  locals[1] = Int(3)          ; quantity
  locals[2] = String("HALF")  ; coupon

Step  IP     Instruction          Stack (top→)              locals[5]
───── ────── ──────────────────── ──────────────────────── ─────────
  1   0000   LOCAL 0              [Float(100.0)]
  2   0003   LOCAL 1              [Float(100.0), Int(3)]
  3   0006   CONVERT float64      [Float(100.0), Float(3.0)]
  4   0009   SETLOCAL 3           [Float(100.0)]              ; t0=3.0
  5   000C   LOCAL 0              [Float(100.0)]
  6   000F   LOCAL 3              [Float(100.0), Float(3.0)]
  7   0012   MUL                  [Float(300.0)]
  8   0013   SETLOCAL 4           []                          ; t1=300.0
  9   0016   SETLOCAL 5           []                          300.0
 10   0019   LOCAL 2              [String("HALF")]
 11   001C   CONST 0              [String("HALF"), String("HALF")]
 12   001F   EQUAL                [Bool(true)]
 13   0020   JUMPFALSE 002E       []                          ; taken? NO → fall through
 14   0023   LOCAL 4              [Float(300.0)]
 15   0026   CONST 1              [Float(300.0), Float(0.5)]
 16   0029   MUL                  [Float(150.0)]
 17   002A   SETLOCAL 6           []                          ; t3=150.0
 18   002D   SETLOCAL 5           []                          150.0
 19   002E   LOCAL 5              [Float(150.0)]              ; total=150.0
 20   0031   CONST 2              [Float(150.0), Float(0.0)]
 21   0034   GREATER              [Bool(true)]
 22   0035   JUMPFALSE 0042       []                          ; not taken
 23   0038   LOCAL 5              [Float(150.0)]
 24   003B   CONST 3              [Float(150.0), Float(10000.0)]
 25   003E   LESS                 [Bool(true)]
 26   003F   SETLOCAL 9           []                          ; valid=true
 27   0045   LOCAL 5              [Float(150.0)]
 28   0048   LOCAL 9              [Float(150.0), Bool(true)]
 29   004B   PACK 2               [Tuple(150.0, true)]
 30   004D   RETURNVAL                                        ; → (150.0, true)
```

Result: `(150.0, true)` — the order is valid, total = $150 after 50% coupon.

---

### Panic, Defer, and Recover

Gig implements Go's panic/defer/recover semantics faithfully, including nested
panics (a panic inside a deferred function).

#### Data Structures

```go
// vm/frame.go
type DeferInfo struct {
    fn      *bytecode.CompiledFunction  // function to call
    args    []value.Value               // captured arguments
    closure *Closure                    // closure (for indirect defers)
}

// vm/vm.go — panic state
type panicState struct {
    panicking bool
    panicVal  value.Value
}

type vm struct {
    // ...
    panicking  bool              // panic in progress
    panicVal   value.Value       // current panic value
    panicStack []panicState      // saved panics (for nested panics)
    deferDepth int               // nesting level of defer execution
}
```

#### How It Works

1. **`OpDefer`**: captures the function + arguments into `frame.defers`
2. **`OpRunDefers`**: normal return path — executes defers in LIFO order, each in a child VM
3. **`OpPanic`**: sets `v.panicking = true` and `v.panicVal`
4. **Panic handler** (top of dispatch loop): when `v.panicking` is true:
   - Iterates `frame.defers` in LIFO order
   - Before each deferred call, pushes current panic state onto `panicStack`
   - Executes the defer via recursive `v.run()`
   - After the defer, checks if `recover()` cleared the saved state
   - If recovered: continues with remaining defers in normal mode
   - If not recovered: propagates panic to caller frame
5. **`OpRecover`**: pops the top of `panicStack`, clears `panicking`, returns the value

#### Concrete Example

```go
func SafeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("caught: %v", r)
        }
    }()
    return a / b, nil
}
```

Execution of `SafeDivide(10, 0)`:

```
1. Frame created for SafeDivide
   locals[0] = Int(10), locals[1] = Int(0)

2. OpDefer: captures anonymous closure into frame.defers
   defers = [{fn: anon$1, args: [], closure: {FreeVars: [&result, &err]}}]

3. OpDiv: Int(10) / Int(0)
   → Go-level panic: "runtime error: integer divide by zero"
   → Safety net catches it: v.panicking = true, v.panicVal = String("integer divide by zero")

4. Panic handler activates (top of dispatch loop):
   - Saves panic state: panicStack = [{panicking:true, val:"integer divide by zero"}]
   - Clears v.panicking
   - Calls anon$1 via recursive v.run()

5. Inside anon$1:
   - OpRecover: pops panicStack, finds panicking=true
     → Clears saved state (panicking=false)
     → Returns String("integer divide by zero")
   - fmt.Errorf wraps it → writes to &err via free variable
   - anon$1 returns

6. Back in panic handler:
   - Checks saved state: panicking is false → recovered!
   - Reads named return values from ResultAllocSlots
   - Returns (0, error("caught: integer divide by zero"))
```

The `ResultAllocSlots` mechanism is crucial: in Go, `defer` can modify named return
values. The compiler records which local slots correspond to named returns, and the
recovery path dereferences them to get the final values.

#### Safety Net

All of this is wrapped in a Go-level `defer/recover`:

```go
// vm/vm.go
func (v *vm) Execute(funcName string, ctx context.Context, args ...value.Value) (result value.Value, err error) {
    // Safety net: catch Go-level panics from VM execution
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

This ensures that even if the interpreted code triggers a host-level panic (nil map
write, slice out of bounds, type assertion failure), the host process never crashes.
The panic is caught and returned as an error.

---

### Goroutine Support

```go
// In interpreted code:
go processItem(item)
```

The VM handles `go` statements via `OpGoCall`:

```go
// vm/goroutine.go
func (v *vm) newChildVM() *vm {
    child := &vm{
        program:      v.program,
        stack:        make([]value.Value, initialStackSize),  // fresh stack
        frames:       make([]*Frame, initialFrameDepth),
        globalsPtr:   v.globalsPtr,       // shared globals via pointer!
        ctx:          v.ctx,              // shared context
        extCallCache: v.extCallCache,     // shared cache
        goroutines:   v.goroutines,       // shared tracker
    }
    if child.globalsPtr == nil {
        child.globalsPtr = &v.globals     // parent's globals become shared
    }
    return child
}
```

**Key design decisions**:
- Each goroutine gets a **fresh stack** (no contention)
- Globals are **shared via pointer** (correct Go semantics)
- The external call cache is **shared** (thread-safe via RWMutex)
- Context is **shared** (cancellation propagates to all goroutines)

The `GoroutineTracker` prevents runaway goroutine creation:

```go
func (t *GoroutineTracker) Start(fn func()) error {
    max := atomic.LoadInt64(&t.maxGoroutines)
    if max > 0 && atomic.LoadInt64(&t.active) >= max {
        return fmt.Errorf("gig: goroutine limit (%d) exceeded", max)
    }
    atomic.AddInt64(&t.active, 1)
    go func() {
        defer atomic.AddInt64(&t.active, -1)
        fn()
    }()
    return nil
}
```

---

## External Package Integration

### How Registration Works

External Go packages must be registered before compilation. Registration happens at
`init()` time via code-generated files in `stdlib/packages/`:

```go
// stdlib/packages/strings.go (generated by `gig gen`)
func init() {
    pkg := importer.RegisterPackage("strings", "strings")

    // Functions
    pkg.AddFunction("Contains", strings.Contains, "", direct_strings_Contains)
    pkg.AddFunction("HasPrefix", strings.HasPrefix, "", direct_strings_HasPrefix)
    // ... 60+ more functions

    // Types
    pkg.AddType("Builder", reflect.TypeOf(strings.Builder{}), "")
    pkg.AddType("Reader", reflect.TypeOf(strings.Reader{}), "")

    // Method DirectCalls
    pkg.AddMethodDirectCall("Builder", "WriteString", direct_method_strings_Builder_WriteString)
}
```

Each registered package provides:
- **Function values** for `reflect.Call` (slow path)
- **DirectCall wrappers** for zero-reflection calls (fast path)
- **Type information** for type checking and runtime type assertions
- **Method DirectCalls** for zero-reflection method dispatch

### DirectCall: Zero-Reflection Function Calls

The biggest performance win in Gig comes from **DirectCall wrappers** — code-generated
typed wrappers that bypass `reflect.Call` entirely.

#### The Problem with reflect.Call

```go
// Slow path: reflect.Call
fn := reflect.ValueOf(strings.Contains)
args := []reflect.Value{
    reflect.ValueOf("hello world"),
    reflect.ValueOf("world"),
}
result := fn.Call(args)  // ~400ns: reflection overhead, allocations
```

#### DirectCall: The Solution

```go
// Generated by gig gen — zero reflection
func direct_strings_Contains(args []value.Value) value.Value {
    a0 := args[0].String()   // direct field access, no reflect
    a1 := args[1].String()
    return value.MakeBool(strings.Contains(a0, a1))  // direct Go call
}
```

This wrapper:
1. Extracts typed arguments directly from `value.Value` (via `String()`, `Int()`, etc.)
2. Calls the real Go function directly (no `reflect.ValueOf`, no `reflect.Call`)
3. Wraps the result in a `value.Value` (via `MakeBool`, `MakeInt`, etc.)

**Result**: ~5x faster than `reflect.Call` for typical stdlib functions.

#### Dispatch Flow

```
OpCallExternal(funcIdx, numArgs)
        │
        ▼
┌─ Check inline cache ─┐
│  cache[funcIdx] hit?  │
│    YES → use entry    │
│    NO  → resolve once │
└───────┬───────────────┘
        │
        ▼
┌─ DirectCall available? ─┐
│  YES → call wrapper     │──▶ result = direct_strings_Contains(args)
│  NO  → reflect.Call     │──▶ result = fn.Call(reflectArgs)
└─────────────────────────┘
```

The inline cache (`externalCallCache`) ensures that function resolution happens
exactly once per function per program lifetime. After the first call, subsequent
calls are a simple pointer dereference + function call.

```go
// vm/call.go — callExternal fast path
func (v *vm) callExternal(funcIdx, numArgs int) error {
    // Pop arguments
    args := make([]value.Value, numArgs)
    for i := numArgs - 1; i >= 0; i-- {
        args[i] = v.pop()
    }

    // Inline cache lookup (RLock for read path)
    v.extCallCache.mu.RLock()
    cacheEntry := v.extCallCache.cache[funcIdx]
    v.extCallCache.mu.RUnlock()

    if cacheEntry == nil {
        // First call: resolve and cache (write lock)
        v.extCallCache.mu.Lock()
        cacheEntry = v.resolveExternalFunc(funcIdx)
        v.extCallCache.cache[funcIdx] = cacheEntry
        v.extCallCache.mu.Unlock()
    }

    // Fast path: DirectCall
    if cacheEntry.directCall != nil {
        result := cacheEntry.directCall(args)
        v.push(result)
        return nil
    }

    // Slow path: reflect.Call
    return v.callExternalReflect(cacheEntry, args)
}
```

### Closure Conversion for External Calls

When interpreted code passes a closure to a Go stdlib function (e.g., `sort.Slice`
with a comparison function), Gig must convert the closure to a real Go function:

```go
// Interpreted code:
sort.Slice(items, func(i, j int) bool {
    return items[i].Price < items[j].Price
})
```

The `sort.Slice` function expects a `func(int, int) bool` — it can't accept a
`*vm.Closure`. Gig uses `reflect.MakeFunc` to create a real Go function that, when
called, creates a temporary VM and executes the closure's bytecode:

```go
// vm/closure.go — Closure implements value.ClosureExecutor
func (c *Closure) Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
    // Create a temporary VM to execute the closure
    closureVM := &vm{
        program: c.Program,
        stack:   make([]value.Value, 256),
        // ...
    }
    valArgs := make([]value.Value, len(args))
    for i, arg := range args {
        valArgs[i] = value.MakeFromReflect(arg)
    }
    closureVM.callFunction(c.Fn, valArgs, c.FreeVars)
    result, _ := closureVM.run()
    return []reflect.Value{result.ToReflectValue(outTypes[0])}
}
```

---

## Init and Execution Flow

### Full Lifecycle

```
Build(source)
    │
    ├── compiler.Build(source, registry)
    │       ├── parser.Parse(source)           → typed AST
    │       ├── ssabuilder.Build(AST)          → SSA IR
    │       └── compiler.Compile(SSA)          → CompiledProgram
    │
    ├── runner.ExecuteInit(program)
    │       ├── Check for "init#1" function
    │       ├── Create temp VM, run "init"
    │       └── Snapshot globals → initialGlobals
    │
    └── runner.New(program, initialGlobals)
            ├── Create VMPool
            ├── Create GoroutineTracker
            └── Register method resolver (for fmt.Stringer)

Program.Run("FuncName", args...)
    │
    ├── Convert args to []value.Value
    │
    ├── vmPool.Get() → vm
    │       ├── If pool empty: newVM() with fresh globals from snapshot
    │       └── If pool has idle: return recycled vm
    │
    ├── vm.Execute("FuncName", ctx, args)
    │       ├── Lookup function in program.Functions
    │       ├── Create frame, copy args to locals
    │       ├── defer { recover → error } (safety net)
    │       └── vm.run() → main dispatch loop
    │
    ├── vmPool.Put(vm)
    │       └── vm.Reset() → clear stack, restore globals from snapshot
    │
    └── UnwrapResult(result) → any
```

### Stateless vs Stateful Modes

**Stateless (default)**: each `Run()` starts from the post-`init()` globals snapshot.
Mutations to globals are discarded after the call. This is safe for concurrent calls.

```go
prog, _ := gig.Build(`
    var counter int
    func Increment() int {
        counter++
        return counter
    }
`)
prog.Run("Increment") // returns 1
prog.Run("Increment") // returns 1 (globals reset!)
```

**Stateful** (`WithStatefulGlobals()`): globals persist across calls. Calls are
serialized with a mutex.

```go
prog, _ := gig.Build(`
    var counter int
    func Increment() int {
        counter++
        return counter
    }
`, gig.WithStatefulGlobals())
prog.Run("Increment") // returns 1
prog.Run("Increment") // returns 2 (globals persist!)
```

```go
// runner/runner.go — stateful execution
func (r *Runner) RunWithValues(ctx context.Context, funcName string, args []value.Value) (value.Value, error) {
    if r.stateful {
        r.runMu.Lock()
        defer r.runMu.Unlock()
        v := r.vmPool.Get()
        v.BindSharedGlobals(&r.sharedGlobals)
        result, err := v.ExecuteWithValues(funcName, ctx, args)
        v.UnbindSharedGlobals()
        r.vmPool.Put(v)
        return result, err
    }
    // Stateless: no lock needed
    v := r.vmPool.Get()
    result, err := v.ExecuteWithValues(funcName, ctx, args)
    r.vmPool.Put(v)
    return result, err
}
```

---

## Security Model

Gig is designed for **sandboxed execution** of untrusted code:

### Compile-Time Checks

| Check | Purpose |
|---|---|
| Ban `import "unsafe"` | Prevents memory corruption |
| Ban `import "reflect"` | Prevents type system bypass |
| Ban `panic()` (default) | Prevents DoS via unrecovered panics. Enable with `WithAllowPanic()` |
| Auto-import only registered packages | Code can't import arbitrary packages |

### Runtime Checks

| Check | Purpose |
|---|---|
| Context cancellation every 1024 instructions | Prevents infinite loops |
| Stack overflow detection (1M slots max) | Prevents memory exhaustion |
| Call stack depth limit (1024 frames) | Prevents stack overflow |
| Goroutine limit (10K default) | Prevents goroutine bomb |
| Safety net `defer/recover` | Host-level panics → error returns |

### Sandbox Registry

For maximum isolation, use a sandbox registry that starts empty:

```go
reg := gig.NewSandboxRegistry()
// Only expose what you want:
// reg.RegisterPackage("strings", "strings")
// (or register nothing — pure computation only)

prog, _ := gig.Build(untrustedCode, gig.WithRegistry(reg))
```

---

## Performance

### Key Optimizations Summary

| Optimization | Technique | Impact |
|---|---|---|
| Frame pooling | LIFO frame recycler | Fib(25): 728K → 7 allocations |
| Value tagged union | 32-byte inline primitives | Zero GC for int/float/bool |
| DirectCall wrappers | Code-gen typed wrappers | ~5x faster than reflect.Call |
| Prebaked constants | `[]value.Value` built once at compile | Eliminates per-instruction `FromInterface` |
| Integer specialization | `intLocals []int64` shadow array | 4x cache utilization (8B vs 32B) |
| Superinstructions | Fused opcodes (17 patterns) | 3-4x fewer dispatch cycles in hot loops |
| Register hoisting | Stack/sp/locals in Go locals | Better CPU register allocation |
| Inline caching | Per-program function resolution cache | O(1) external call dispatch |
| Slice fusion | `OpIntSliceGet/Set` | 5 instructions → 1 for `[]int` access |
| Move fusion | `OpIntMoveLocal` | Eliminates stack round-trip for copies |

### How the Optimizations Stack

For a typical integer-heavy loop like bubble sort:

```
Baseline (naive bytecode):
    LOCAL 0          ; 1 dispatch + stack write
    LOCAL 1          ; 1 dispatch + stack write
    LESS             ; 1 dispatch + 2 stack reads + kind check + compare + stack write
    JUMPFALSE off    ; 1 dispatch + stack read + branch
    ─────────────────
    Total: 4 dispatches, 6 stack ops, 1 kind check

After peephole fusion:
    LessLocalLocalJumpFalse 0 1 off   ; 1 dispatch, 2 local reads, compare, branch
    ─────────────────
    Total: 1 dispatch, 0 stack ops, 1 kind check

After int specialization:
    IntLessLocalLocalJumpFalse 0 1 off ; 1 dispatch, 2 intLocal reads, compare, branch
    ─────────────────
    Total: 1 dispatch, 0 stack ops, 0 kind checks, 8B operands
```

That's a **4x reduction in dispatch cycles** and the operands are now in an 8-byte
array instead of 32-byte Values.

---

From: youngjin, March 2026
