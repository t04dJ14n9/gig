# Gig Internals

Gig is an interpreter of the Go language, written in Go. It was designed as an embeddable, sandboxed execution engine — a way to run user-provided Go code safely inside a host application, with context-based cancellation and without exposing `unsafe`, `reflect`, or `panic` to the interpreted program.

Unlike most Go interpreters that walk an abstract syntax tree at runtime, Gig compiles Go source code through SSA intermediate representation down to a compact bytecode, which is then executed by a stack-based virtual machine. This design — borrowing from both traditional compilers and bytecode VMs like the JVM or Lua's — gives Gig a distinctive performance profile and a clear separation between compilation and execution.

This document is here to look under the hood. In the following, we get an overview, explore the internals and discuss the design. Our aim is to provide the essential insights, clarify the architecture and the code organization. But first, the overview.

## Overview of Architecture

Let's see what happens inside Gig when one executes the following lines:

```go
prog, _ := gig.Build(`
    package main
    import "fmt"
    func Hello(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    }
`)
result, _ := prog.Run("Hello", "world")
```

The following figure displays the main stages:

```
Source Code ──► Parser ──► Type Checker ──► SSA Builder ──► Compiler ──► VM
               go/parser   go/types        go/ssa          bytecode     execute
```

**The parser** (provided by `go/parser`) transforms Go source into an abstract syntax tree.

**The type checker** (provided by `go/types`) resolves all types, constants, and identifiers. It uses a custom `types.Importer` that resolves imports against Gig's registered external packages rather than the filesystem.

**The SSA builder** (provided by `golang.org/x/tools/go/ssa`) transforms the type-checked AST into Static Single Assignment form — a graph of basic blocks containing typed instructions where every value is assigned exactly once.

**The compiler** (implemented in `compiler/`) translates SSA instructions into a flat bytecode stream, performs phi elimination, patches jump targets, and runs four optimization passes.

**The virtual machine** (implemented in `vm/`) executes the bytecode in a fetch-decode-execute loop, managing a value stack, call frames, and external function dispatch.

The interpreter is designed as a proper compiler, except that the code is generated into memory instead of object files, targeting the Go runtime itself rather than a hardware architecture. We won't spend time on the parser, type checker, or SSA builder — all provided by the standard library and its extensions — and instead examine what Gig builds on top of them.

## The Value System

Before we dive into compilation and execution, we must understand the most fundamental data structure in Gig: the `Value`. Every local variable, stack slot, constant, function argument, and return value in the interpreter is a `Value`.

```go
// value/value.go
type Value struct {
    kind Kind    // 1 byte: type tag
    num  int64   // 8 bytes: bool, int, uint bits, float64 bits
    obj  any     // 16 bytes: string, reflect.Value, *Closure, []int64, etc.
}
```

The total size is **32 bytes** on 64-bit systems. This is a tagged-union design, inspired by how Lua and other dynamic languages represent values, but adapted for Go's type system.

The key insight is the **two-tier split** between primitive and composite types:

**Primitive types** (bool, int, uint, float, nil) are stored entirely in `kind` + `num`, with `obj` remaining nil. Creating an integer value is:

```go
func MakeInt(i int64) Value { return Value{kind: KindInt, num: i} }
```

No heap allocation. No reflection. No GC pressure. Two 64-bit words on the stack.

**Composite types** (slices, maps, structs, channels, interfaces) fall through to `obj`, which holds either a `reflect.Value` or a native Go object. For example, integer slices get special treatment:

```go
func MakeIntSlice(s []int64) Value { return Value{kind: KindSlice, obj: s} }
```

The `[]int64` is stored directly — not wrapped in `reflect.Value` — which means the VM can index it, set elements, and take addresses without any reflection overhead.

This design stands in contrast to interpreters like Yaegi, which represent all values as `reflect.Value`. While `reflect.Value` provides universal type handling, it allocates on the heap for primitives and requires dynamic dispatch for every operation. Gig's tagged-union avoids this: an integer addition is literally `result.num = a.num + b.num` — three memory accesses, no allocations, no function calls.

## Compilation

### From SSA to Bytecode

The compilation pipeline is implemented across several files in `compiler/`. The entry point is `Compile()` in `compiler.go`, which takes an SSA program and produces a `bytecode.Program`.

The first pass assigns an index to every function, including anonymous functions and closures. Functions are stored both in a map (by name, for `Run("funcName")` dispatch) and in a flat array (`FuncByIndex`, for O(1) call dispatch at runtime):

```go
// compiler.go
for idx, fn := range allFuncs {
    c.funcIndex[fn] = idx
}
```

The second pass compiles each function. Per-function compilation in `compile_func.go` begins by building a symbol table — mapping each SSA value to a local variable slot:

```go
// compile_func.go — slot allocation
slot := 0
for _, param := range fn.Params {
    c.symbols[param] = slot
    slot++
}
for _, block := range fn.Blocks {
    for _, instr := range block.Instrs {
        if phi, ok := instr.(*ssa.Phi); ok {
            c.phiSlots[phi] = slot
            slot++
        }
    }
}
```

Parameters occupy the first slots, then phi nodes, then all other SSA values. This flat numbering scheme means every `OpLocal` and `OpSetLocal` instruction addresses locals by a simple 16-bit index.

### Basic Block Traversal

Basic blocks are visited in **reverse postorder** — a standard compiler ordering that guarantees every block's dominators are visited before the block itself. For each block, the compiler:

1. Emits `OpSetLocal` instructions for phi nodes (phi elimination)
2. Compiles each SSA instruction into one or more bytecode instructions
3. Emits jumps to successor blocks

Phi elimination deserves a word. In SSA form, phi nodes at block entries merge values from different predecessors. Since our bytecode has no phi concept, we lower them to explicit moves: before jumping to a target block, we emit `OpSetLocal` for each phi node using the edge value from the current predecessor:

```go
// compile_func.go — phi elimination
func (c *compiler) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
    for _, instr := range targetBlock.Instrs {
        phi, ok := instr.(*ssa.Phi)
        if !ok { break }
        sourceValue := phi.Edges[predIndex]
        c.compileValue(sourceValue)          // push source onto stack
        c.emit(bytecode.OpSetLocal, slot)    // pop into phi's local slot
    }
}
```

### The Instruction Set

Gig's bytecode is a variable-length encoding: 1 byte opcode followed by 0–6 bytes of operands in big-endian format. The instruction set has about 100 opcodes, organized into categories:

- **Stack**: `CONST`, `LOCAL`, `SETLOCAL`, `GLOBAL`, `FREE`, `POP`, `DUP`
- **Arithmetic**: `ADD`, `SUB`, `MUL`, `DIV`, `MOD`, `NEG`
- **Comparison**: `EQUAL`, `LESS`, `GREATER`, `LESSEQ`, `GREATEREQ`
- **Control flow**: `JUMP`, `JUMPTRUE`, `JUMPFALSE`, `CALL`, `RETURN`
- **Containers**: `MAKESLICE`, `MAKEMAP`, `INDEX`, `SETINDEX`, `FIELD`, `FIELDADDR`
- **Pointers**: `ADDR`, `DEREF`, `SETDEREF`, `INDEXADDR`
- **External**: `CALLEXTERNAL`, `CALLINDIRECT`
- **Concurrency**: `GOCALL`, `SEND`, `RECV`, `SELECT`, `CLOSE`
- **Superinstructions**: ~30 fused opcodes (discussed in Optimization)

A compiled function is:

```go
// bytecode/bytecode.go
type CompiledFunction struct {
    Name         string
    Instructions []byte     // flat bytecode
    NumLocals    int        // params + phis + temporaries
    NumParams    int
    NumFreeVars  int        // closure captures
    HasIntLocals bool       // needs intLocals shadow array
}
```

Instructions are a plain `[]byte`. There is no instruction struct, no pointer chasing — just a flat byte stream that the VM reads sequentially. This is crucial for CPU cache locality.

### Constants

Constants are stored in three parallel arrays in the `Program`:

```go
type Program struct {
    Constants         []any          // raw: int64, string, ExternalFuncInfo, ...
    PrebakedConstants []value.Value  // pre-converted at compile time
    IntConstants      []int64        // for int-specialized opcodes
}
```

The `PrebakedConstants` array is the key optimization. Instead of converting constants from `any` to `Value` at every `OpConst` execution (which involves a type switch and potential allocation), we do it once at compile time. At runtime, `OpConst` is a single array lookup:

```go
stack[sp] = prebaked[idx]
sp++
```

## The Virtual Machine

### Structure

The VM is a stack machine with a separate call frame stack:

```go
// vm/vm.go
type VM struct {
    program      *bytecode.Program
    stack        []value.Value     // operand stack (initial 1024)
    sp           int               // stack pointer
    frames       []*Frame          // call frame stack (initial 64)
    fp           int               // frame pointer
    globals      []value.Value
    ctx          context.Context
    extCallCache sync.Map          // inline cache for external calls
    fpool        framePool         // frame recycling
}
```

Each call frame stores the execution state for a function invocation:

```go
// vm/frame.go
type Frame struct {
    fn        *bytecode.CompiledFunction
    ip        int                // instruction pointer into fn.Instructions
    basePtr   int                // stack base for this frame
    locals    []value.Value      // local variables
    intLocals []int64            // integer-specialized shadow array
    freeVars  []*value.Value     // closure captures (shared pointers)
    defers    []DeferInfo
}
```

### The Dispatch Loop

The core of the VM is a single `run()` function in `vm/run.go`. Its structure follows a pattern found in most high-performance bytecode interpreters — a tight loop with a `switch` statement:

```go
// vm/run.go (simplified)
func (vm *VM) run() (value.Value, error) {
    // Hoist frame state into locals for register allocation
    stack := vm.stack
    sp := vm.sp
    prebaked := vm.program.PrebakedConstants
    intConsts := vm.program.IntConstants
    var frame *Frame
    var ins []byte
    var locals []value.Value
    var intLocals []int64

    for vm.fp > 0 {
        op := bytecode.OpCode(ins[frame.ip])
        frame.ip++

        switch op {
        case bytecode.OpLocal:
            idx := uint16(ins[frame.ip])<<8 | uint16(ins[frame.ip+1])
            frame.ip += 2
            stack[sp] = locals[idx]
            sp++
            continue

        case bytecode.OpConst:
            idx := uint16(ins[frame.ip])<<8 | uint16(ins[frame.ip+1])
            frame.ip += 2
            stack[sp] = prebaked[idx]
            sp++
            continue

        case bytecode.OpAdd:
            sp--
            b := stack[sp]
            sp--
            a := stack[sp]
            if a.IsInt() {
                stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())
            } else {
                stack[sp] = value.Add(a, b) // generic path
            }
            sp++
            continue

        // ... 60+ hot-path opcodes inlined here ...

        default:
            // Cold path: sync state, call executeOp()
            vm.sp = sp
            vm.executeOp(op) // handles ~40 less common opcodes
            sp = vm.sp
        }
    }
}
```

There are several important details here:

**Register hoisting**: The frame's `locals`, `intLocals`, `ins` (instruction stream), and the VM's `stack`, `sp` are copied into local variables at the top of `run()`. The Go compiler can then place these in machine registers, avoiding repeated pointer dereferencing through `vm.stack[vm.sp]` on every instruction. When a call or return changes the active frame, these locals are re-synced.

**Hot/cold split**: The `switch` directly handles ~60 frequently-executed opcodes. The remaining ~40 (type conversions, channel operations, select, defer, panic recovery) go through `executeOp()` in a separate function. This keeps the hot loop's machine code smaller, improving instruction cache behavior.

**Integer fast-paths**: Arithmetic opcodes like `OpAdd` check `a.IsInt()` first. Since `IsInt()` is just `v.kind == KindInt` (a single byte comparison), and integer operations are `v.num + v.num` (no allocation), the common case of integer arithmetic is a handful of machine instructions.

### Context Checking

To support cancellation and timeouts, the VM checks the context periodically:

```go
instructionCount++
if instructionCount & 0x3FF == 0 {   // every 1024 instructions
    select {
    case <-vm.ctx.Done():
        return value.MakeNil(), vm.ctx.Err()
    default:
    }
}
```

The bitwise AND trick avoids an expensive modulo operation. The check happens every 1024 instructions — frequent enough for responsive cancellation (sub-millisecond on most workloads), infrequent enough to be negligible in profiling.

### Frame Pooling

Function calls are the hot path in recursive programs. Without optimization, each call to `fib(n-1)` would allocate a new `Frame` struct with a fresh `[]value.Value` locals slice — exactly what made early Gig slow on Fibonacci (728,000 allocations for Fib25).

The solution is a frame pool:

```go
// vm/frame.go
type framePool struct {
    frames []*Frame
}

func (p *framePool) get(fn *bytecode.CompiledFunction, basePtr int, freeVars []*value.Value) *Frame {
    if len(p.frames) > 0 {
        f := p.frames[len(p.frames)-1]
        p.frames = p.frames[:len(p.frames)-1]
        // Reuse if locals capacity is sufficient
        if cap(f.locals) >= fn.NumLocals {
            f.locals = f.locals[:fn.NumLocals]
            for i := range f.locals {
                f.locals[i] = value.Value{} // zero
            }
            // ... set fn, ip, basePtr, freeVars
            return f
        }
    }
    // Allocate new
    return &Frame{...}
}
```

When a function returns, its frame goes back to the pool. The locals slice is reused if it's large enough, just zeroed. This brought Fib25 allocations from 728,000 down to 7 — only the initial VM, stack, and frame allocations remain.

One subtlety: frames where a local's address has been taken (`OpAddr` on a local variable) are **not** returned to the pool. A closure might hold a `*value.Value` pointing into that frame's locals slice, and reusing it would corrupt the closure's captured state.

### Call Dispatch

Gig handles three kinds of calls:

**Compiled function calls** (`OpCall`): The function index is embedded in the instruction. The VM looks up `program.FuncByIndex[idx]` (O(1) array access), pushes a new frame, copies arguments from the stack into the frame's locals, and continues execution.

**Closure calls** (`OpCallIndirect`): The top of the stack holds a `*Closure` struct containing the function index and an array of `*value.Value` pointers (the captured free variables). The VM unwraps the closure, pushes a frame with the free vars attached, and proceeds as above.

**External function calls** (`OpCallExternal`): This is where it gets interesting — and where a large part of the optimization work was focused.

## External Package Integration

A Go interpreter that can only run pure algorithms isn't very useful. The real value comes from calling the Go standard library — `fmt.Sprintf`, `strings.Contains`, `json.Marshal`, `http.Get`. But these are compiled Go functions; the interpreter can't just call them directly. There's a type boundary to cross.

### Registration

External packages are registered at init time via generated code:

```go
// stdlib/packages/strings.go (generated)
func init() {
    pkg := importer.RegisterPackage("strings", "strings")
    pkg.AddFunction("Contains", strings.Contains, "", directcall_Contains)
    pkg.AddFunction("HasPrefix", strings.HasPrefix, "", directcall_HasPrefix)
    pkg.AddType("Builder", reflect.TypeOf(strings.Builder{}), "")
    pkg.AddMethodDirectCall("Builder", "WriteString", directcall_method_Builder_WriteString)
    // ... 40+ functions, types, methods
}
```

Each registered package provides:
- **Functions**: the function value, its `go/types` signature (resolved from `reflect.Type`), and optionally a DirectCall wrapper
- **Types**: the `reflect.Type`, converted to `types.Named` with all exported methods added
- **Variables and constants**: registered similarly

The importer implements `types.Importer`, so when the Go type checker encounters `import "strings"`, it gets a `types.Package` with all the correct type signatures, as if it were reading from compiled `.a` files. This means type checking is exact — if the interpreted code misuses a standard library function, it gets a proper compile error, not a runtime panic.

### The Reflection Problem

The naive approach to calling external functions is straightforward:

```go
// 1. Convert []value.Value → []reflect.Value
reflectArgs := make([]reflect.Value, len(args))
for i, arg := range args {
    reflectArgs[i] = reflect.ValueOf(arg.Interface())
}
// 2. Call via reflection
results := reflect.ValueOf(fn).Call(reflectArgs)
// 3. Convert []reflect.Value → []value.Value
```

This works, but is devastatingly slow. Step 1 allocates a `[]reflect.Value` slice and boxes every argument. Step 2 goes through `reflect.Value.Call`, which performs safety checks, type validation, and ultimately calls `runtime.call` through an indirect function pointer. Step 3 unboxes results.

For `strings.Contains("hello", "ell")` — a function that takes 30 nanoseconds natively — the reflection overhead adds about 500 nanoseconds and 5 heap allocations.

### DirectCall: Eliminating Reflection

The solution is code generation. For each function with compatible parameter types, Gig generates a **typed wrapper** at build time:

```go
// stdlib/packages/strings.go (generated)
func directcall_Contains(args []value.Value) value.Value {
    a0 := args[0].String()         // extract string directly from Value
    a1 := args[1].String()
    r0 := strings.Contains(a0, a1) // native Go function call
    return value.FromBool(r0)      // wrap result directly
}
```

No `reflect.Value`. No `Call()`. No allocation. The argument extraction uses `Value.String()`, `Value.Int()`, etc., which are just field accesses on the tagged union. The result wrapping uses `value.FromBool()`, which is `Value{kind: KindBool, num: ...}`. The actual function call compiles to a direct `CALL` instruction in machine code — the Go compiler can even inline `strings.Contains` into the wrapper.

This extends to **methods** as well:

```go
func directcall_method_Builder_WriteString(args []value.Value) value.Value {
    recv := args[0].Interface().(*strings.Builder)  // type assertion
    a1 := args[1].String()
    r0, r1 := recv.WriteString(a1)
    // ... wrap results
}
```

The receiver is extracted via a type assertion on `Value.Interface()` — still zero-reflection on the call itself.

### Coverage and Type Support

The code generator (`gentool/directcall.go`) handles a wide range of parameter types:

| Type | Extraction |
|---|---|
| `string`, `int`, `bool`, `float64` | `.String()`, `.Int()`, `.Bool()`, `.Float()` |
| `[]byte` | `.Bytes()` |
| `io.Reader`, `error` | `.Interface().(io.Reader)` |
| `*bytes.Buffer` | `.Interface().(*bytes.Buffer)` |
| `*int32`, `*int64` | `.Interface().(*int32)` |
| `map[string]bool` | `.Interface().(map[string]bool)` |
| `any` / `interface{}` | `.Interface()` |

Functions with `unsafe.Pointer` parameters or certain complex variadic signatures are left on the reflection path — about 8% of stdlib functions.

In total: **1,162 wrappers** (619 functions + 543 methods) across 20 standard library packages, covering 92% of the standard library surface.

### Inline Caching

Even with DirectCall, the VM still needs to resolve which function to call. The constant pool stores `ExternalFuncInfo` and `ExternalMethodInfo` structs, but looking them up on every call would be wasteful. So the VM maintains an **inline cache** — a `sync.Map` keyed by constant pool index:

```go
// vm/call.go
type extCallCacheEntry struct {
    funcInfo   *bytecode.ExternalFuncInfo
    directCall func([]value.Value) value.Value
    reflectFn  reflect.Value
}

func (vm *VM) callExternal(constIdx int, numArgs int) {
    entry, cached := vm.extCallCache.Load(constIdx)
    if !cached {
        // Resolve once, cache forever
        entry = &extCallCacheEntry{...}
        vm.extCallCache.Store(constIdx, entry)
    }
    if entry.directCall != nil {
        result := entry.directCall(args)  // fast path
    } else {
        vm.callExternalReflect(entry, args)  // fallback
    }
}
```

After the first call, subsequent calls to the same external function hit the cache — one `sync.Map` lookup (essentially a pointer read in the uncontended case) and a direct function call.

## Optimization Passes

After initial compilation, four optimization passes transform the bytecode:

### Pass 1: Peephole Superinstructions

The optimizer scans for common multi-instruction sequences and replaces them with single **superinstructions**. The idea comes from the Forth tradition and has been used in interpreters from CPython to Lua.

Example: the Go statement `sum += a` compiles to 4 instructions totaling 11 bytes:

```
LOCAL(sum)  LOCAL(a)  ADD  SETLOCAL(sum)
```

The peephole pass fuses this into a single 7-byte superinstruction:

```
OpLocalLocalAddSetLocal(sum, a, sum)
```

This eliminates 3 dispatch cycles, 2 stack pushes, and 2 stack pops. The operands are encoded directly in the instruction — no stack traffic at all.

The optimizer recognizes **17 patterns**, including:

| Pattern | Fused Opcode |
|---|---|
| `LOCAL(A) LOCAL(B) ADD SETLOCAL(C)` | `OpLocalLocalAddSetLocal(A,B,C)` |
| `LOCAL(A) CONST(B) ADD SETLOCAL(C)` | `OpLocalConstAddSetLocal(A,B,C)` |
| `LOCAL(A) CONST(B) LESS JUMPTRUE(off)` | `OpLessLocalConstJumpTrue(A,B,off)` |
| `LOCAL(A) CONST(B) LESSEQ JUMPFALSE(off)` | `OpLessEqLocalConstJumpFalse(A,B,off)` |
| `ADD SETLOCAL(A)` | `OpAddSetLocal(A)` |

The rewriting must be offset-aware: when instructions are shortened, all jump targets must be remapped. The optimizer builds an offset map after rewriting and adjusts every jump instruction.

### Pass 2: Slice Operation Fusion

Integer slice access patterns are recognized and fused. The Go statement `v = arr[j]`, when both are `int` typed, compiles to a 7-instruction sequence (17 bytes) involving `LOCAL`, `INDEXADDR`, `SETLOCAL`, `DEREF`. The optimizer fuses this into:

```
OpIntSliceGet(arr, j, v)    // 7 bytes, direct []int64 indexed access
```

Similarly for writes: `arr[j] = v` becomes `OpIntSliceSet(arr, j, v)`.

### Pass 3: Integer Specialization

This is the most aggressive optimization. It introduces a **shadow array** of native `int64` values alongside the regular `[]value.Value` locals:

```go
// vm/frame.go
type Frame struct {
    locals    []value.Value   // 32 bytes per slot
    intLocals []int64         // 8 bytes per slot (shadow)
}
```

The optimizer performs two passes over the bytecode:

**Pass 1 (analysis)**: Identify which local indices participate exclusively in integer operations — locals that are sources or destinations of `OpLocalLocalAddSetLocal`, `OpLessLocalConstJumpTrue`, etc.

**Pass 2 (upgrade)**: Replace eligible superinstructions with `OpInt*` variants:

```
OpLocalConstAddSetLocal(A, B, C)  →  OpIntLocalConstAddSetLocal(A, B, C)
```

The `OpInt*` variant operates directly on `intLocals`:

```go
// vm/run.go
case bytecode.OpIntLocalConstAddSetLocal:
    r := intLocals[idxA] + intConsts[idxB]   // raw int64 add
    intLocals[idxC] = r                       // 8-byte write
    locals[idxC] = value.MakeInt(r)           // sync to Value locals
```

The inner loop of `ArithmeticSum` — `sum += i; i++; i < n` — compiles to just 3 dispatches per iteration, all operating on 8-byte `int64` slots instead of 32-byte `Value` slots. This is 4x better cache utilization.

A critical invariant is the **dual write**: every `OpInt*` instruction writes to both `intLocals[idx]` (for fast int operations) and `locals[idx]` (for non-specialized code that might read the same local). This maintains correctness without requiring dataflow analysis to determine when the Value copy is needed.

### Pass 4: Move Fusion

The final pass eliminates phi-move overhead for integer locals:

```
OpIntLocal(A)  OpIntSetLocal(B)  →  OpIntMoveLocal(A, B)
```

This replaces a push-pop pair with a direct register-to-register copy.

### Cumulative Effect

The four passes work together. Consider a simple loop:

```go
for i := 0; i < n; i++ {
    sum += arr[i]
}
```

**After Pass 1**: `LOCAL(i) CONST(1) ADD SETLOCAL(i)` → `OpLocalConstAddSetLocal(i, 1, i)`

**After Pass 2**: `LOCAL(arr) LOCAL(i) INDEXADDR... DEREF... SETLOCAL(v)` → `OpIntSliceGet(arr, i, v)`

**After Pass 3**: `OpLocalConstAddSetLocal(i, 1, i)` → `OpIntLocalConstAddSetLocal(i, 1, i)` (native int64)

**After Pass 4**: Phi moves at loop entry → `OpIntMoveLocal`

The result: the inner loop is 3–4 fused instructions operating on 8-byte integers with direct slice access. No stack traffic, no 32-byte value copies, no type checks in the hot path.

## Goroutines and Concurrency

When the interpreted code spawns a goroutine with `go func()`, the VM creates a child VM sharing the same globals:

```go
// vm/goroutine.go
func (vm *VM) newChildVM() *VM {
    child := &VM{
        program:    vm.program,
        globalsPtr: &vm.globals,  // shared for cross-goroutine communication
        ctx:        vm.ctx,
    }
    child.initStack()
    return child
}
```

The child gets its own stack and frame stack but shares the program, globals, and context. Channels work through Go's native channel operations on `reflect.Value` — the interpreter doesn't reimplement channel semantics.

`select` statements are handled by building a `reflect.SelectCase` slice and calling `reflect.Select()`, which delegates to the Go runtime's select implementation. This is one area where reflection cannot be avoided, but it occurs infrequently enough to not be a bottleneck.

## Security Model

Gig enforces a security sandbox at the earliest possible stage — before compilation:

```go
// gig.go
func checkBannedImports(file *ast.File) error {
    for _, imp := range file.Imports {
        path := strings.Trim(imp.Path.Value, `"`)
        if path == "unsafe" || path == "reflect" {
            return fmt.Errorf("import %q is not allowed", path)
        }
    }
    return nil
}
```

By banning `unsafe` and `reflect` at the AST level, interpreted code cannot:
- Read or write arbitrary memory
- Circumvent type safety
- Access unexported fields
- Forge interface values

The `panic` built-in is also restricted — interpreted code cannot crash the host process. `defer` and `recover` work within the interpreter's frame stack, contained by the VM.

Context cancellation ensures the host can always terminate a runaway script:

```go
result, err := prog.RunWithContext(ctx, "ProcessData", input)
if err == context.DeadlineExceeded {
    log.Warn("script timed out")
}
```

## Design Choices: Gig vs Yaegi

It's instructive to compare Gig's design with Yaegi's, as they solve the same problem with fundamentally different approaches.

**AST-walking vs Bytecode VM**: Yaegi walks the AST at runtime, generating closures on the fly for each node. Gig compiles through SSA to bytecode. The tradeoff: Yaegi has lower compilation overhead (no SSA construction, no optimization passes), but Gig has lower execution overhead (linear bytecode, superinstructions, integer specialization).

**`reflect.Value` vs Tagged-union**: Yaegi represents every value as a `reflect.Value`. Gig uses a 32-byte tagged-union that avoids allocation for primitives. The result: Fibonacci(25) in Yaegi performs 2.1 million allocations; in Gig, 7.

**Control-flow representation**: Yaegi annotates the AST with `tnext`/`fnext` pointers, forming a control-flow graph embedded in the tree. Gig uses flat bytecode with explicit jump offsets, enabling sequential prefetch and superinstruction fusion — optimizations that are impractical on a tree structure.

**External call strategy**: Both interpreters must call through reflection for external packages. Yaegi generates closure wrappers around `reflect.Value` operations. Gig generates typed Go functions at build time (DirectCall) that avoid reflection entirely for 92% of standard library calls.

The benchmarks tell the story: Gig is 1.1–5.2x faster than Yaegi across all workloads, with dramatically fewer allocations. The gap is largest on recursion (5.2x on Fib25 — frame pooling dominates), external calls (2.6–2.8x — DirectCall eliminates reflection), and closures (2.7x — shared pointer captures vs Yaegi's scope chain).

## Conclusion

We have described the architecture of a Go interpreter that takes a different path from AST-walking: SSA-based compilation to bytecode, executed by a stack-based VM with aggressive specialization. The key design decisions — tagged-union values, superinstruction fusion, integer-specialized locals, and generated DirectCall wrappers — each address a specific performance bottleneck while maintaining full Go language compatibility.

The codebase is organized into clean layers: `bytecode/` as the shared kernel, `compiler/` and `vm/` as independent consumers, `value/` as the universal data representation, and `importer/` + `gentool/` bridging the gap to the host Go runtime. The whole thing compiles to a single binary with no external dependencies.

Some areas remain for future work:
- Register-based VM (eliminating stack traffic entirely)
- JIT compilation for hot functions
- Escape analysis for smarter frame pooling
- More aggressive constant folding at compile time

From: youngjin, 28 Feb 2026
