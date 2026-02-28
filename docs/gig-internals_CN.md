# Gig 内部实现

Gig 是一个用 Go 编写的 Go 语言解释器。它被设计为一个可嵌入的沙箱执行引擎——一种在宿主应用程序内安全运行用户提供的 Go 代码的方式，支持基于 context 的取消机制，且不向被解释的程序暴露 `unsafe`、`reflect` 或 `panic`。

与大多数在运行时遍历抽象语法树的 Go 解释器不同，Gig 通过 SSA 中间表示将 Go 源代码编译为紧凑的字节码，然后由基于栈的虚拟机执行。这种设计——借鉴了传统编译器和字节码虚拟机（如 JVM 或 Lua）——赋予 Gig 独特的性能特征和清晰的编译/执行分离。

本文旨在深入探讨其内部实现。在接下来的内容中，我们将概览全局、探索内部细节并讨论设计决策。目标是提供关键的洞察，阐明架构和代码组织。首先，从概述开始。

## 架构概览

让我们看看当执行以下代码时，Gig 内部发生了什么：

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

下图展示了主要阶段：

```
源代码 ──► 解析器 ──► 类型检查器 ──► SSA 构建器 ──► 编译器 ──► 虚拟机
          go/parser   go/types      go/ssa        bytecode    execute
```

**解析器**（由 `go/parser` 提供）将 Go 源代码转换为抽象语法树。

**类型检查器**（由 `go/types` 提供）解析所有类型、常量和标识符。它使用自定义的 `types.Importer`，针对 Gig 注册的外部包（而非文件系统）来解析导入。

**SSA 构建器**（由 `golang.org/x/tools/go/ssa` 提供）将经过类型检查的 AST 转换为静态单赋值形式——一个由基本块组成的图，其中包含类型化指令，每个值只被赋值一次。

**编译器**（在 `compiler/` 中实现）将 SSA 指令翻译为扁平的字节码流，执行 phi 消除，修补跳转目标，并运行四个优化遍。

**虚拟机**（在 `vm/` 中实现）在取指-译码-执行循环中执行字节码，管理值栈、调用帧和外部函数分发。

该解释器被设计为一个真正的编译器，只不过代码生成到内存中而非目标文件，目标是 Go 运行时本身而非硬件架构。我们不会花时间在解析器、类型检查器或 SSA 构建器上——这些都由标准库及其扩展提供——而是着重考察 Gig 在其之上构建的部分。

## 值系统

在深入编译和执行之前，我们必须理解 Gig 中最基础的数据结构：`Value`。解释器中的每个局部变量、栈槽位、常量、函数参数和返回值都是一个 `Value`。

```go
// value/value.go
type Value struct {
    kind Kind    // 1 字节：类型标签
    num  int64   // 8 字节：bool、int、uint 位模式、float64 位模式
    obj  any     // 16 字节：string、reflect.Value、*Closure、[]int64 等
}
```

在 64 位系统上总大小为 **32 字节**。这是一种标记联合体设计，灵感来自 Lua 和其他动态语言表示值的方式，但适配了 Go 的类型系统。

关键洞察是基本类型和复合类型之间的**两层拆分**：

**基本类型**（bool、int、uint、float、nil）完全存储在 `kind` + `num` 中，`obj` 保持为 nil。创建整数值的方式为：

```go
func MakeInt(i int64) Value { return Value{kind: KindInt, num: i} }
```

无堆分配。无反射。无 GC 压力。栈上两个 64 位字。

**复合类型**（切片、映射、结构体、通道、接口）回落到 `obj`，它持有 `reflect.Value` 或原生 Go 对象。例如，整数切片得到特殊处理：

```go
func MakeIntSlice(s []int64) Value { return Value{kind: KindSlice, obj: s} }
```

`[]int64` 直接存储——不包装在 `reflect.Value` 中——这意味着 VM 可以对其进行索引、设置元素和取地址，完全无需反射开销。

这种设计与 Yaegi 等解释器形成对比，后者将所有值表示为 `reflect.Value`。虽然 `reflect.Value` 提供了通用类型处理，但对基本类型会在堆上分配，且每次操作都需要动态分发。Gig 的标记联合体避免了这一点：一次整数加法字面上就是 `result.num = a.num + b.num`——三次内存访问，无分配，无函数调用。

## 编译

### 从 SSA 到字节码

编译流水线在 `compiler/` 的多个文件中实现。入口是 `compiler.go` 中的 `Compile()`，它接受一个 SSA 程序并生成 `bytecode.Program`。

第一遍为每个函数分配索引，包括匿名函数和闭包。函数既存储在映射中（按名称，用于 `Run("funcName")` 分发），也存储在扁平数组中（`FuncByIndex`，用于运行时 O(1) 调用分发）：

```go
// compiler.go
for idx, fn := range allFuncs {
    c.funcIndex[fn] = idx
}
```

第二遍编译每个函数。`compile_func.go` 中的按函数编译首先构建符号表——将每个 SSA 值映射到局部变量槽位：

```go
// compile_func.go — 槽位分配
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

参数占据前面的槽位，然后是 phi 节点，然后是其他所有 SSA 值。这种扁平编号方案意味着每条 `OpLocal` 和 `OpSetLocal` 指令通过简单的 16 位索引寻址局部变量。

### 基本块遍历

基本块按**逆后序**遍历——这是一种标准的编译器排序，保证每个块的支配者在该块之前被访问。对于每个块，编译器：

1. 为 phi 节点发出 `OpSetLocal` 指令（phi 消除）
2. 将每条 SSA 指令编译为一条或多条字节码指令
3. 发出到后继块的跳转

Phi 消除值得一提。在 SSA 形式中，块入口处的 phi 节点合并来自不同前驱的值。由于我们的字节码没有 phi 概念，我们将其降低为显式的移动：在跳转到目标块之前，我们为每个 phi 节点使用当前前驱的边值发出 `OpSetLocal`：

```go
// compile_func.go — phi 消除
func (c *compiler) emitPhiMoves(predBlock, targetBlock *ssa.BasicBlock) {
    for _, instr := range targetBlock.Instrs {
        phi, ok := instr.(*ssa.Phi)
        if !ok { break }
        sourceValue := phi.Edges[predIndex]
        c.compileValue(sourceValue)          // 将源值压入栈
        c.emit(bytecode.OpSetLocal, slot)    // 弹出到 phi 的局部槽位
    }
}
```

### 指令集

Gig 的字节码是变长编码：1 字节操作码后跟 0–6 字节大端序操作数。指令集约有 100 个操作码，按类别组织：

- **栈操作**：`CONST`、`LOCAL`、`SETLOCAL`、`GLOBAL`、`FREE`、`POP`、`DUP`
- **算术**：`ADD`、`SUB`、`MUL`、`DIV`、`MOD`、`NEG`
- **比较**：`EQUAL`、`LESS`、`GREATER`、`LESSEQ`、`GREATEREQ`
- **控制流**：`JUMP`、`JUMPTRUE`、`JUMPFALSE`、`CALL`、`RETURN`
- **容器**：`MAKESLICE`、`MAKEMAP`、`INDEX`、`SETINDEX`、`FIELD`、`FIELDADDR`
- **指针**：`ADDR`、`DEREF`、`SETDEREF`、`INDEXADDR`
- **外部调用**：`CALLEXTERNAL`、`CALLINDIRECT`
- **并发**：`GOCALL`、`SEND`、`RECV`、`SELECT`、`CLOSE`
- **超级指令**：约 30 个融合操作码（在优化部分讨论）

编译后的函数为：

```go
// bytecode/bytecode.go
type CompiledFunction struct {
    Name         string
    Instructions []byte     // 扁平字节码
    NumLocals    int        // 参数 + phi 节点 + 临时变量
    NumParams    int
    NumFreeVars  int        // 闭包捕获
    HasIntLocals bool       // 是否需要 intLocals 影子数组
}
```

指令是纯 `[]byte`。没有指令结构体，没有指针追踪——只是一个 VM 顺序读取的扁平字节流。这对 CPU 缓存局部性至关重要。

### 常量

常量存储在 `Program` 中的三个并行数组中：

```go
type Program struct {
    Constants         []any          // 原始值：int64、string、ExternalFuncInfo 等
    PrebakedConstants []value.Value  // 编译时预转换
    IntConstants      []int64        // 用于整数特化操作码
}
```

`PrebakedConstants` 数组是关键优化。不是在每次 `OpConst` 执行时将常量从 `any` 转换为 `Value`（涉及类型开关和潜在分配），而是在编译时一次完成。在运行时，`OpConst` 只是一次数组查找：

```go
stack[sp] = prebaked[idx]
sp++
```

## 虚拟机

### 结构

VM 是一个带有独立调用帧栈的栈机器：

```go
// vm/vm.go
type VM struct {
    program      *bytecode.Program
    stack        []value.Value     // 操作数栈（初始 1024）
    sp           int               // 栈指针
    frames       []*Frame          // 调用帧栈（初始 64）
    fp           int               // 帧指针
    globals      []value.Value
    ctx          context.Context
    extCallCache sync.Map          // 外部调用内联缓存
    fpool        framePool         // 帧回收池
}
```

每个调用帧存储一次函数调用的执行状态：

```go
// vm/frame.go
type Frame struct {
    fn        *bytecode.CompiledFunction
    ip        int                // 指令指针，指向 fn.Instructions
    basePtr   int                // 此帧的栈基址
    locals    []value.Value      // 局部变量
    intLocals []int64            // 整数特化影子数组
    freeVars  []*value.Value     // 闭包捕获（共享指针）
    defers    []DeferInfo
}
```

### 分发循环

VM 的核心是 `vm/run.go` 中的单个 `run()` 函数。其结构遵循大多数高性能字节码解释器中的模式——带有 `switch` 语句的紧密循环：

```go
// vm/run.go（简化版）
func (vm *VM) run() (value.Value, error) {
    // 将帧状态提升到局部变量以便寄存器分配
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
                stack[sp] = value.Add(a, b) // 通用路径
            }
            sp++
            continue

        // ... 60 多个热路径操作码在此内联 ...

        default:
            // 冷路径：同步状态，调用 executeOp()
            vm.sp = sp
            vm.executeOp(op) // 处理约 40 个不太常用的操作码
            sp = vm.sp
        }
    }
}
```

这里有几个重要细节：

**寄存器提升**：帧的 `locals`、`intLocals`、`ins`（指令流）以及 VM 的 `stack`、`sp` 在 `run()` 顶部被复制到局部变量中。Go 编译器随后可以将这些放入机器寄存器，避免每条指令都通过 `vm.stack[vm.sp]` 进行重复的指针解引用。当调用或返回改变活动帧时，这些局部变量会重新同步。

**热/冷分离**：`switch` 直接处理约 60 个频繁执行的操作码。剩余约 40 个（类型转换、通道操作、select、defer、panic 恢复）通过单独函数中的 `executeOp()` 处理。这使热循环的机器码更小，改善了指令缓存行为。

**整数快速路径**：像 `OpAdd` 这样的算术操作码首先检查 `a.IsInt()`。由于 `IsInt()` 只是 `v.kind == KindInt`（单字节比较），而整数运算是 `v.num + v.num`（无分配），整数算术的常见情况只需少量机器指令。

### Context 检查

为支持取消和超时，VM 定期检查 context：

```go
instructionCount++
if instructionCount & 0x3FF == 0 {   // 每 1024 条指令
    select {
    case <-vm.ctx.Done():
        return value.MakeNil(), vm.ctx.Err()
    default:
    }
}
```

位与技巧避免了昂贵的取模运算。检查每 1024 条指令进行一次——对于响应式取消足够频繁（大多数工作负载下亚毫秒级），又不至于在性能分析中产生明显影响。

### 帧池化

函数调用是递归程序中的热路径。如果不进行优化，每次调用 `fib(n-1)` 都会分配一个新的 `Frame` 结构体和一个新的 `[]value.Value` 局部变量切片——这正是早期 Gig 在 Fibonacci 上变慢的原因（Fib25 有 728,000 次分配）。

解决方案是帧池：

```go
// vm/frame.go
type framePool struct {
    frames []*Frame
}

func (p *framePool) get(fn *bytecode.CompiledFunction, basePtr int, freeVars []*value.Value) *Frame {
    if len(p.frames) > 0 {
        f := p.frames[len(p.frames)-1]
        p.frames = p.frames[:len(p.frames)-1]
        // 如果局部变量容量足够则复用
        if cap(f.locals) >= fn.NumLocals {
            f.locals = f.locals[:fn.NumLocals]
            for i := range f.locals {
                f.locals[i] = value.Value{} // 清零
            }
            // ... 设置 fn、ip、basePtr、freeVars
            return f
        }
    }
    // 分配新帧
    return &Frame{...}
}
```

当函数返回时，其帧回到池中。如果局部变量切片足够大则复用，只需清零。这使 Fib25 的分配次数从 728,000 降至 7——只剩初始的 VM、栈和帧分配。

一个细微之处：如果某个局部变量的地址被获取（对局部变量执行 `OpAddr`），该帧**不会**返回池中。闭包可能持有指向该帧局部变量切片的 `*value.Value`，复用它会破坏闭包的捕获状态。

### 调用分发

Gig 处理三种调用：

**编译函数调用**（`OpCall`）：函数索引嵌入在指令中。VM 查找 `program.FuncByIndex[idx]`（O(1) 数组访问），压入新帧，将参数从栈复制到帧的局部变量中，然后继续执行。

**闭包调用**（`OpCallIndirect`）：栈顶持有一个 `*Closure` 结构体，包含函数索引和一个 `*value.Value` 指针数组（捕获的自由变量）。VM 解包闭包，压入带有自由变量的帧，然后按上述方式继续。

**外部函数调用**（`OpCallExternal`）：这是最有趣的部分——也是大部分优化工作的焦点。

## 外部包集成

一个只能运行纯算法的 Go 解释器并不太有用。真正的价值在于调用 Go 标准库——`fmt.Sprintf`、`strings.Contains`、`json.Marshal`、`http.Get`。但这些是编译后的 Go 函数；解释器不能直接调用它们。存在一个类型边界需要跨越。

### 注册

外部包在 init 时通过生成的代码注册：

```go
// stdlib/packages/strings.go（生成的）
func init() {
    pkg := importer.RegisterPackage("strings", "strings")
    pkg.AddFunction("Contains", strings.Contains, "", directcall_Contains)
    pkg.AddFunction("HasPrefix", strings.HasPrefix, "", directcall_HasPrefix)
    pkg.AddType("Builder", reflect.TypeOf(strings.Builder{}), "")
    pkg.AddMethodDirectCall("Builder", "WriteString", directcall_method_Builder_WriteString)
    // ... 40 多个函数、类型、方法
}
```

每个注册的包提供：
- **函数**：函数值、其 `go/types` 签名（从 `reflect.Type` 解析）以及可选的 DirectCall 包装器
- **类型**：`reflect.Type`，转换为带有所有导出方法的 `types.Named`
- **变量和常量**：以类似方式注册

导入器实现了 `types.Importer`，因此当 Go 类型检查器遇到 `import "strings"` 时，它获得一个具有所有正确类型签名的 `types.Package`，就像从编译后的 `.a` 文件读取一样。这意味着类型检查是精确的——如果解释代码误用了标准库函数，它会得到正确的编译错误，而不是运行时 panic。

### 反射问题

调用外部函数的朴素方法很直接：

```go
// 1. 将 []value.Value 参数转换为 []reflect.Value
reflectArgs := make([]reflect.Value, len(args))
for i, arg := range args {
    reflectArgs[i] = reflect.ValueOf(arg.Interface())
}
// 2. 通过反射调用
results := reflect.ValueOf(fn).Call(reflectArgs)
// 3. 将 []reflect.Value 结果转换为 []value.Value
```

这能工作，但非常慢。步骤 1 分配 `[]reflect.Value` 切片并装箱每个参数。步骤 2 通过 `reflect.Value.Call` 执行安全检查、类型验证，并最终通过间接函数指针调用 `runtime.call`。步骤 3 拆箱结果。

对于 `strings.Contains("hello", "ell")`——一个原生只需 30 纳秒的函数——反射开销增加了约 500 纳秒和 5 次堆分配。

### DirectCall：消除反射

解决方案是代码生成。对于每个具有兼容参数类型的函数，Gig 在构建时生成一个**类型化包装器**：

```go
// stdlib/packages/strings.go（生成的）
func directcall_Contains(args []value.Value) value.Value {
    a0 := args[0].String()         // 直接从 Value 提取字符串
    a1 := args[1].String()
    r0 := strings.Contains(a0, a1) // 原生 Go 函数调用
    return value.FromBool(r0)      // 直接包装结果
}
```

无 `reflect.Value`。无 `Call()`。无分配。参数提取使用 `Value.String()`、`Value.Int()` 等，这些只是标记联合体上的字段访问。结果包装使用 `value.FromBool()`，即 `Value{kind: KindBool, num: ...}`。实际的函数调用编译为机器码中的直接 `CALL` 指令——Go 编译器甚至可以将 `strings.Contains` 内联到包装器中。

这也扩展到**方法**：

```go
func directcall_method_Builder_WriteString(args []value.Value) value.Value {
    recv := args[0].Interface().(*strings.Builder)  // 类型断言
    a1 := args[1].String()
    r0, r1 := recv.WriteString(a1)
    // ... 包装结果
}
```

接收者通过 `Value.Interface()` 上的类型断言提取——调用本身仍然是零反射。

### 覆盖率和类型支持

代码生成器（`gentool/directcall.go`）处理广泛的参数类型：

| 类型 | 提取方式 |
|---|---|
| `string`、`int`、`bool`、`float64` | `.String()`、`.Int()`、`.Bool()`、`.Float()` |
| `[]byte` | `.Bytes()` |
| `io.Reader`、`error` | `.Interface().(io.Reader)` |
| `*bytes.Buffer` | `.Interface().(*bytes.Buffer)` |
| `*int32`、`*int64` | `.Interface().(*int32)` |
| `map[string]bool` | `.Interface().(map[string]bool)` |
| `any` / `interface{}` | `.Interface()` |

带有 `unsafe.Pointer` 参数或某些复杂可变参数签名的函数留在反射路径上——约占标准库函数的 8%。

总计：**1,162 个包装器**（619 个函数 + 543 个方法），覆盖 20 个标准库包，涵盖 92% 的标准库接口。

### 内联缓存

即使有了 DirectCall，VM 仍需解析要调用哪个函数。常量池存储 `ExternalFuncInfo` 和 `ExternalMethodInfo` 结构体，但每次调用都查找它们会很浪费。因此 VM 维护一个**内联缓存**——以常量池索引为键的 `sync.Map`：

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
        // 解析一次，永久缓存
        entry = &extCallCacheEntry{...}
        vm.extCallCache.Store(constIdx, entry)
    }
    if entry.directCall != nil {
        result := entry.directCall(args)  // 快速路径
    } else {
        vm.callExternalReflect(entry, args)  // 回退
    }
}
```

第一次调用后，对同一外部函数的后续调用命中缓存——一次 `sync.Map` 查找（在无竞争情况下本质上是一次指针读取）和一次直接函数调用。

## 优化遍

初始编译后，四个优化遍变换字节码：

### 第一遍：窥孔超级指令

优化器扫描常见的多指令序列，并将其替换为单个**超级指令**。这个思想来自 Forth 传统，已被 CPython 到 Lua 等解释器使用。

示例：Go 语句 `sum += a` 编译为 4 条指令共 11 字节：

```
LOCAL(sum)  LOCAL(a)  ADD  SETLOCAL(sum)
```

窥孔遍将其融合为单条 7 字节超级指令：

```
OpLocalLocalAddSetLocal(sum, a, sum)
```

这消除了 3 次分发周期、2 次栈压入和 2 次栈弹出。操作数直接编码在指令中——完全没有栈流量。

优化器识别 **17 种模式**，包括：

| 模式 | 融合操作码 |
|---|---|
| `LOCAL(A) LOCAL(B) ADD SETLOCAL(C)` | `OpLocalLocalAddSetLocal(A,B,C)` |
| `LOCAL(A) CONST(B) ADD SETLOCAL(C)` | `OpLocalConstAddSetLocal(A,B,C)` |
| `LOCAL(A) CONST(B) LESS JUMPTRUE(off)` | `OpLessLocalConstJumpTrue(A,B,off)` |
| `LOCAL(A) CONST(B) LESSEQ JUMPFALSE(off)` | `OpLessEqLocalConstJumpFalse(A,B,off)` |
| `ADD SETLOCAL(A)` | `OpAddSetLocal(A)` |

重写必须感知偏移：当指令被缩短时，所有跳转目标必须重映射。优化器在重写后构建偏移映射，并调整每条跳转指令。

### 第二遍：切片操作融合

整数切片访问模式被识别并融合。Go 语句 `v = arr[j]`，当两者都是 `int` 类型时，编译为一个 7 条指令（17 字节）的序列，涉及 `LOCAL`、`INDEXADDR`、`SETLOCAL`、`DEREF`。优化器将其融合为：

```
OpIntSliceGet(arr, j, v)    // 7 字节，直接 []int64 索引访问
```

写入同理：`arr[j] = v` 变为 `OpIntSliceSet(arr, j, v)`。

### 第三遍：整数特化

这是最激进的优化。它在常规 `[]value.Value` 局部变量旁引入了一个原生 `int64` 值的**影子数组**：

```go
// vm/frame.go
type Frame struct {
    locals    []value.Value   // 每个槽位 32 字节
    intLocals []int64         // 每个槽位 8 字节（影子）
}
```

优化器对字节码执行两遍：

**第一遍（分析）**：识别哪些局部变量索引专门参与整数运算——作为 `OpLocalLocalAddSetLocal`、`OpLessLocalConstJumpTrue` 等的源或目标的局部变量。

**第二遍（升级）**：将符合条件的超级指令替换为 `OpInt*` 变体：

```
OpLocalConstAddSetLocal(A, B, C)  →  OpIntLocalConstAddSetLocal(A, B, C)
```

`OpInt*` 变体直接在 `intLocals` 上操作：

```go
// vm/run.go
case bytecode.OpIntLocalConstAddSetLocal:
    r := intLocals[idxA] + intConsts[idxB]   // 原始 int64 加法
    intLocals[idxC] = r                       // 8 字节写入
    locals[idxC] = value.MakeInt(r)           // 同步到 Value 局部变量
```

`ArithmeticSum` 的内层循环——`sum += i; i++; i < n`——编译为每次迭代仅 3 次分发，全部在 8 字节 `int64` 槽位上操作而非 32 字节 `Value` 槽位。这是 4 倍的缓存利用率提升。

一个关键不变量是**双写**：每条 `OpInt*` 指令同时写入 `intLocals[idx]`（用于快速整数运算）和 `locals[idx]`（供可能读取同一局部变量的非特化代码使用）。这在无需数据流分析来确定何时需要 Value 副本的情况下保持了正确性。

### 第四遍：移动融合

最后一遍消除整数局部变量的 phi 移动开销：

```
OpIntLocal(A)  OpIntSetLocal(B)  →  OpIntMoveLocal(A, B)
```

这将压入-弹出对替换为直接的寄存器到寄存器复制。

### 累积效果

四个遍协同工作。考虑一个简单循环：

```go
for i := 0; i < n; i++ {
    sum += arr[i]
}
```

**第一遍后**：`LOCAL(i) CONST(1) ADD SETLOCAL(i)` → `OpLocalConstAddSetLocal(i, 1, i)`

**第二遍后**：`LOCAL(arr) LOCAL(i) INDEXADDR... DEREF... SETLOCAL(v)` → `OpIntSliceGet(arr, i, v)`

**第三遍后**：`OpLocalConstAddSetLocal(i, 1, i)` → `OpIntLocalConstAddSetLocal(i, 1, i)`（原生 int64）

**第四遍后**：循环入口处的 phi 移动 → `OpIntMoveLocal`

结果：内层循环是 3–4 条融合指令，在 8 字节整数上操作并直接进行切片访问。没有栈流量，没有 32 字节值复制，热路径中没有类型检查。

## 协程与并发

当解释代码用 `go func()` 启动协程时，VM 创建一个共享相同全局变量的子 VM：

```go
// vm/goroutine.go
func (vm *VM) newChildVM() *VM {
    child := &VM{
        program:    vm.program,
        globalsPtr: &vm.globals,  // 共享用于跨协程通信
        ctx:        vm.ctx,
    }
    child.initStack()
    return child
}
```

子 VM 拥有自己的栈和帧栈，但共享程序、全局变量和 context。通道通过 Go 原生的 `reflect.Value` 通道操作工作——解释器不重新实现通道语义。

`select` 语句通过构建 `reflect.SelectCase` 切片并调用 `reflect.Select()` 来处理，后者委托给 Go 运行时的 select 实现。这是一个无法避免反射的领域，但它发生得不够频繁，不会成为瓶颈。

## 安全模型

Gig 在尽可能早的阶段——编译之前——实施安全沙箱：

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

通过在 AST 层面禁止 `unsafe` 和 `reflect`，解释代码无法：
- 读写任意内存
- 绕过类型安全
- 访问未导出的字段
- 伪造接口值

内建的 `panic` 也被限制——解释代码无法使宿主进程崩溃。`defer` 和 `recover` 在解释器的帧栈内工作，由 VM 包含。

Context 取消确保宿主可以随时终止失控的脚本：

```go
result, err := prog.RunWithContext(ctx, "ProcessData", input)
if err == context.DeadlineExceeded {
    log.Warn("脚本超时")
}
```

## 设计选择：Gig vs Yaegi

将 Gig 的设计与 Yaegi 进行比较很有启发性，因为它们以根本不同的方法解决同一个问题。

**AST 遍历 vs 字节码 VM**：Yaegi 在运行时遍历 AST，为每个节点动态生成闭包。Gig 通过 SSA 编译为字节码。权衡是：Yaegi 编译开销更低（无 SSA 构建，无优化遍），但 Gig 执行开销更低（线性字节码、超级指令、整数特化）。

**`reflect.Value` vs 标记联合体**：Yaegi 将每个值表示为 `reflect.Value`。Gig 使用 32 字节标记联合体，避免基本类型的分配。结果：Fibonacci(25) 在 Yaegi 中执行 210 万次分配；在 Gig 中，7 次。

**控制流表示**：Yaegi 用 `tnext`/`fnext` 指针注释 AST，在树中形成控制流图。Gig 使用带有显式跳转偏移的扁平字节码，支持顺序预取和超级指令融合——这些在树结构上不切实际的优化。

**外部调用策略**：两个解释器都必须通过反射调用外部包。Yaegi 围绕 `reflect.Value` 操作生成闭包包装器。Gig 在构建时生成类型化的 Go 函数（DirectCall），对 92% 的标准库调用完全避免反射。

基准测试说明了一切：Gig 在所有工作负载上比 Yaegi 快 1.1–5.2 倍，分配次数大幅减少。差距在递归上最大（Fib25 上 5.2 倍——帧池化占主导），外部调用上（2.6–2.8 倍——DirectCall 消除了反射），以及闭包上（2.7 倍——共享指针捕获 vs Yaegi 的作用域链）。

## 总结

我们描述了一个走出 AST 遍历不同路径的 Go 解释器的架构：基于 SSA 的字节码编译，由带有激进特化的基于栈的 VM 执行。关键的设计决策——标记联合体值、超级指令融合、整数特化局部变量和生成的 DirectCall 包装器——每一个都针对特定的性能瓶颈，同时保持完整的 Go 语言兼容性。

代码库组织为清晰的层次：`bytecode/` 作为共享内核，`compiler/` 和 `vm/` 作为独立的消费者，`value/` 作为通用数据表示，`importer/` + `gentool/` 弥合与宿主 Go 运行时的差距。整个项目编译为单个二进制文件，无外部依赖。

一些领域仍有待未来工作：
- 基于寄存器的 VM（完全消除栈流量）
- 热函数的 JIT 编译
- 用于更智能帧池化的逃逸分析
- 编译时更激进的常量折叠

来自：youngjin，2026 年 2 月 28 日
