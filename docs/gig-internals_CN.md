# Gig 内部原理

面向需要理解系统运作方式、进行扩展或排查生产问题的工程师的 Gig Go 解释器技术深度解析。

---

## 为什么选择 Gig？规则引擎的困境

每家企业都有变化速度超出部署周期的业务逻辑。定价规则、资格判定、欺诈评分、促销匹配——这些每周甚至每天都在变化。现有方案各有致命缺陷：

| 方案 | 问题 |
|---|---|
| 硬编码 Go | 每次规则变更都需要重新编译和部署 |
| 表达式语言 (CEL, Rego) | 能力有限：没有循环、没有标准库、还得学新语法 |
| 嵌入 Lua/JS | 语言不同：Go 开发者需要频繁切换上下文，无法复用 Go 库 |
| gRPC 微服务 | 运维开销大：每条规则都要单独部署、版本管理和监控 |

**我们真正想要的是**：用 Go 编写规则（零学习成本），可以调用任何 Go 标准库或第三方库（功能齐全），但无需重新编译宿主应用即可动态加载和执行。

### Gig 的方案

Gig 是一个**完整的 Go 解释器**，将 Go 源代码编译为字节码并在基于栈的虚拟机上执行。它不是子集也不是 DSL——它支持完整的 Go 语言特性，包括 goroutine、闭包、defer/panic/recover、接口、方法和类型断言。

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

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/t04dJ14n9/gig"
)

func main() {
    // 将 Go 源码编译为字节码
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

    // 执行编译后的函数
    result, err := prog.Run("ProcessName", "  hello world  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: HELLO WORLD
}
```

`Build` 调用耗时约 1-5ms（解析 + 类型检查 + SSA + 编译）。每次 `Run` 调用仅需微秒级——字节码已经编译完毕，VM 通过池化复用。

---

## 架构概览

### 包结构

```
gig/
├── gig.go                    # 公共 API: Build(), Program.Run()
├── compiler/
│   ├── build.go              # 完整流水线: source → parse → SSA → bytecode
│   ├── compiler.go           # SSA → 字节码翻译
│   ├── compile_func.go       # 逐函数编译
│   ├── compile_instr.go      # 逐指令编译
│   ├── symbol.go             # 符号表 (SSA value → local slot)
│   ├── parser/               # go/parser + 安全性校验
│   ├── ssa/                  # go/ssa builder 封装
│   ├── peephole/             # 基于模式的超级指令融合
│   └── optimize/             # 4 遍字节码优化流水线
├── vm/
│   ├── vm.go                 # VM 结构体，Execute() 入口
│   ├── run.go                # 主取指-解码-执行循环（热路径）
│   ├── frame.go              # 调用帧 + 帧池
│   ├── stack.go              # 操作数栈，支持有界增长
│   ├── call.go               # 外部函数调用 (DirectCall + reflect)
│   ├── closure.go            # 闭包类型 + ClosureExecutor
│   ├── goroutine.go          # GoroutineTracker，子 VM 构造
│   ├── ops_dispatch.go       # 操作码路由到分类处理器
│   ├── ops_arithmetic.go     # 算术、位运算、比较操作
│   ├── ops_memory.go         # 栈、局部变量、全局变量、字段、地址操作
│   ├── ops_container.go      # Slice、map、channel 操作
│   ├── ops_control.go        # 控制流、defer、panic/recover
│   ├── ops_convert.go        # 类型断言、类型转换
│   └── ops_call.go           # 函数/闭包调用、goroutine
├── model/
│   ├── value/                # 32 字节标记联合体 Value 类型
│   ├── bytecode/             # CompiledProgram, CompiledFunction, OpCode
│   └── external/             # ExternalFuncInfo, ExternalMethodInfo
├── importer/                 # 包注册、类型解析
├── runner/                   # VM 池、init 执行、有状态全局变量
└── stdlib/packages/          # ~69 个预生成的标准库封装
```

### 构建流水线详解

```go
// gig.go: Build()
func Build(sourceCode string, opts ...BuildOption) (*Program, error) {
    // 1. compiler.Build: source → parse → SSA → bytecode
    result, err := compiler.Build(sourceCode, cfg.registry, compilerOpts...)

    // 2. 执行 init() 并快照全局变量
    initialGlobals, err := runner.ExecuteInit(result.Program)

    // 3. 创建 runner（拥有 VM 池）
    r := runner.New(result.Program, initialGlobals, runnerOpts...)

    return &Program{runner: r, ssaPkg: result.SSAPkg}, nil
}
```

`compiler.Build` 函数协调三个阶段：

```go
// compiler/build.go
func Build(source string, reg importer.PackageRegistry, opts ...BuildOption) (*BuildResult, error) {
    // 阶段 1：解析 + 类型检查 + 校验
    parseResult, err := parser.Parse(source, reg, parseOpts...)

    // 阶段 2：构建 SSA
    ssaResult, err := ssabuilder.Build(parseResult.FSet, parseResult.Pkg, ...)

    // 阶段 3：将 SSA 编译为字节码
    lookup := importer.NewPackageLookup(reg)
    compiled, err := NewCompiler(lookup).Compile(ssaResult.Pkg)

    return &BuildResult{Program: compiled, SSAPkg: ssaResult.Pkg}, nil
}
```

---

## 值系统

### 问题

解释器需要一个通用类型来在运行时表示任意 Go 值：`int`、`string`、`[]byte`、`*http.Request`、闭包等。最直接的做法是 `interface{}`，但在解释器的操作数栈中，中间值不断地 push/pop，逃逸分析无法优化——大量整数运算的结果会在堆上分配。对于数值密集型的规则引擎来说，这意味着巨大的 GC 压力。

### 解决方案：32 字节标记联合体

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
    kind Kind    // 1 byte: 类型标记 (KindInt, KindString, KindFloat, ...)
    size Size    // 1 byte: 原始 Go 位宽 (8, 16, 32, 64)
    num  int64   // 8 bytes: 存储 bool (0/1), int, uint 位, float64 位
    obj  any     // 16 bytes: string, complex128, reflect.Value, 或 nil
}
```

**关键洞察**：对于基本类型（`bool`、`int`、`uint`、`float64`、`nil`），所有数据都存在 `kind` + `num` 中。`obj` 字段为 `nil`。**零堆分配，零 GC 压力。**

### Kind 类型

```go
const (
    KindInvalid Kind = iota  // 0 — Value 的零值（未初始化的全局变量）
    KindNil                  // 1 — 显式 nil
    KindBool                 // 2 — 存储在 num 中: 0=false, 1=true
    KindInt                  // 3 — 以 int64 形式存储在 num 中
    KindUint                 // 4 — 以 uint64 位模式存储在 num 中
    KindFloat                // 5 — 以 float64 位模式存储在 num 中 (math.Float64bits)
    KindString               // 6 — 以 Go string 形式存储在 obj 中
    KindComplex              // 7 — 以 complex128 形式存储在 obj 中
    KindPointer              // 8
    KindSlice                // 9
    KindArray                // 10
    KindMap                  // 11
    KindChan                 // 12
    KindFunc                 // 13
    KindStruct               // 14
    KindInterface            // 15
    KindReflect              // 16 — 兜底: obj 中存储 reflect.Value
    KindBytes                // 17 — 原生 []byte（避免 reflect 开销）
)
```

### 构造函数

```go
// 基本类型：零分配
value.MakeInt(42)         // kind=KindInt, num=42, obj=nil
value.MakeFloat(3.14)     // kind=KindFloat, num=float64bits(3.14), obj=nil
value.MakeBool(true)      // kind=KindBool, num=1, obj=nil
value.MakeNil()           // kind=KindNil, num=0, obj=nil

// 字符串：obj 持有 Go string
value.MakeString("hello") // kind=KindString, num=0, obj="hello"

// 复合类型：obj 持有 reflect.Value 或原生 Go 类型
value.MakeBytes([]byte{1,2,3})  // kind=KindBytes, obj=[]byte{1,2,3}
value.MakeIntSlice([]int64{...}) // kind=KindSlice, obj=[]int64{...}
value.FromInterface(anyValue)    // 自动检测：快速路径类型 switch
```

### Size 标记：保留原始类型

Go 有 `int8`、`int16`、`int32`、`int64`、`int`——在内部都以 `int64` 存储。`size` 字段记住了原始类型，这样 `Interface()` 就能返回正确的类型：

```go
value.MakeInt8(42).Interface()  // 返回 int8(42)，而非 int64(42)
value.MakeInt32(42).Interface() // 返回 int32(42)
value.MakeInt(42).Interface()   // 返回 int(42)
```

这在通过反射将值传递给外部 Go 函数时非常重要——如果你调用 `strings.Repeat(s, n)` 而 `n` 是 `int64` 而不是 `int`，`reflect.Call` 会 panic。

### 为什么不用 `interface{}`？

以计算 Fib(25) = 75,025 的斐波那契基准测试为例：

| 值表示方式 | 分配次数 | 原因 |
|---|---|---|
| 所有值用 `interface{}` | ~2.1M | 中间值在解释器栈中逃逸到堆上 |
| `value.Value` 标记联合体 | ~7 | 仅初始帧 + 栈分配 |

32 字节的 `Value` 恰好占两条缓存行，基本类型操作永远不会逃逸到堆上。

---

## 编译：从 Go 源码到字节码

### 一个具体示例

让我们跟踪这个函数经过完整编译流水线的过程：

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

### 阶段 1：解析与类型检查

解析器做三件事：

1. **`go/parser.ParseFile`** 生成 AST
2. **安全校验**：检查被禁止的导入（`unsafe`、`reflect`）以及被禁止的 `panic()` 用法（可通过 `WithAllowPanic()` 配置）
3. **自动导入**：如果源码引用了 `strings.Contains(...)` 但没有 `import "strings"`，解析器会自动添加（因为已注册的包是已知的）
4. **`go/types.Config.Check`**：使用自定义的 `types.Importer` 进行类型检查，它根据注册表解析包

### 阶段 2：SSA 构建

`golang.org/x/tools/go/ssa` 将带类型信息的 AST 转换为静态单赋值形式。对于我们的例子，SSA 大致如下：

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

**SSA 核心概念**：
- 每个值只被赋值一次（SSA 性质）
- **Phi 节点**合并来自不同控制路径的值（例如 `t4` 在 `t1` 和 `t3` 之间选择）
- 短路求值 `&&` 变成了带 Phi 的显式控制流

### 阶段 3：字节码生成

编译器将 SSA 指令翻译为基于栈的字节码。以下是编译过程：

#### 符号表构建

首先，编译器分配局部变量槽位：

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

参数占据前面的槽位。Phi 节点获得专用槽位（编译器发出显式的 `SetLocal` 移动指令来解析 Phi）。临时变量占据剩余的槽位。

#### 生成的字节码

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

#### Phi 消除

SSA 的 Phi 节点无法直接映射到硬件指令。编译器通过以下方式消除它们：
1. 为每个 Phi 分配一个专用局部变量槽位（`total` 对应 slot 5，`valid` 对应 slot 9）
2. 在每个前驱基本块末尾，发出 `SETLOCAL` 将正确的值写入 Phi 槽位

这就是为什么你在 entry 块和 if.then 块中都能看到 `SETLOCAL 5`——每个块在控制流到达合并点之前，将自己版本的 `total` 写入共享的 Phi 槽位。

### 阶段 4：优化

初始字节码生成之后，运行四遍优化：

```go
// compiler/optimize/optimize.go
func Optimize(code []byte, localIsInt, constIsInt, localIsIntSlice []bool) ([]byte, bool) {
    code = Peephole(code)                                    // 第 1 遍：超级指令融合
    code = FuseSliceOps(code, localIsInt, localIsIntSlice)   // 第 2 遍：slice 操作融合
    code, hasInt := IntSpecialize(code, localIsInt, constIsInt) // 第 3 遍：整数特化
    code = FuseIntMoves(code)                                // 第 4 遍：移动融合
    return code, hasInt
}
```

#### 第 1 遍：窥孔优化——超级指令融合

窥孔模式检测常见的多指令序列，并将其替换为单条融合操作码。例如：

```
优化前（3 条指令，3 次分发）：
    LOCAL    0        ; push a
    LOCAL    1        ; push b
    ADD               ; a + b

优化后（1 条指令，1 次分发）：
    ADDLOCALLOCAL 0 1  ; push locals[0] + locals[1]
```

窥孔优化器有 17+ 条模式规则，覆盖：

| 模式 | 融合操作码 | 节省 |
|---|---|---|
| `LOCAL a` + `LOCAL b` + `ADD` | `OpAddLocalLocal a b` | 2 次分发 |
| `LOCAL a` + `CONST c` + `ADD` | `OpAddLocalConst a c` | 2 次分发 |
| `ADD` + `SETLOCAL x` | `OpAddSetLocal x` | 1 次分发 |
| `LOCAL a` + `LOCAL b` + `ADD` + `SETLOCAL c` | `OpLocalLocalAddSetLocal a b c` | 3 次分发 |
| `LOCAL a` + `LOCAL b` + `LESS` + `JUMPTRUE off` | `OpLessLocalLocalJumpTrue a b off` | 3 次分发 |

每个模式注册在全局模式注册表中：

```go
// compiler/peephole/pattern.go
type Pattern interface {
    Match(code []byte, i int) (consumed int, newBytes []byte, ok bool)
}
```

#### 第 2 遍：Slice 操作融合

检测常见模式 `LOCAL(slice)` + `LOCAL(index)` + `INDEXADDR` + `SETLOCAL(ptr)` + `LOCAL(ptr)` + `DEREF` + `SETLOCAL(val)`，当所有类型都是 `int` 时，将其融合为单条 `OpIntSliceGet slice index val`。

#### 第 3 遍：整数特化

当编译器能证明超级指令中所有局部变量都是 `int` 类型时，将泛型超级指令升级为 `OpInt*` 变体：

```
OpLocalConstAddSetLocal → OpIntLocalConstAddSetLocal
OpLessLocalLocalJumpFalse → OpIntLessLocalLocalJumpFalse
```

`OpInt*` 变体操作影子数组 `intLocals []int64`——纯 8 字节 int64 运算，而非 32 字节的 Value 操作：

```go
// vm/run.go — OpIntLocalConstAddSetLocal handler
case bytecode.OpIntLocalConstAddSetLocal:
    idxA := readU16()
    idxB := readU16()
    idxC := readU16()
    r := intLocals[idxA] + intConsts[idxB]   // 原始 int64 加法——无 kind 检查
    intLocals[idxC] = r                       // 写入 int 影子数组
    locals[idxC] = value.MakeInt(r)           // 保持 Value 数组同步
    continue
```

缓存利用率提升 4 倍：`int64` 是 8 字节，而 `Value` 是 32 字节，同一条缓存行能容纳 4 倍的操作数。

#### 第 4 遍：移动融合

将 `OpIntLocal(src)` + `OpIntSetLocal(dst)` 替换为 `OpIntMoveLocal(src, dst)`，消除栈的往返开销。

---

## 虚拟机

### VM 结构

```go
// vm/vm.go
type vm struct {
    program        *bytecode.CompiledProgram  // 编译后的字节码
    stack          []value.Value              // 操作数栈
    sp             int                        // 栈指针
    frames         []*Frame                   // 调用帧栈
    fp             int                        // 帧指针
    globals        []value.Value              // 包级变量
    globalsPtr     *[]value.Value             // 共享全局变量（goroutine 间）
    ctx            context.Context            // 取消/超时控制
    panicking      bool                       // 是否正在 panic？
    panicVal       value.Value                // 当前 panic 值
    panicStack     []panicState               // 保存的 panic（嵌套场景）
    deferDepth     int                        // defer 嵌套层级
    extCallCache   *externalCallCache         // 外部调用内联缓存
    initialGlobals []value.Value              // init 后的全局变量快照
    goroutines     *GoroutineTracker          // goroutine 限制器
    fpool          framePool                  // 帧回收器
}
```

### 关键常量

```go
const (
    initialStackSize     = 1024     // 操作数栈初始槽数
    maxStackSize         = 1 << 20  // 1M 槽位 = 每个 VM 32 MB
    initialFrameDepth    = 64       // 调用帧初始槽数
    maxFrameDepth        = 1024     // 最大调用深度（8 KB）
    contextCheckInterval = 1024     // 每 N 条指令检查一次 ctx.Done()
    defaultMaxGoroutines = 10000    // 每个 program 的 goroutine 上限
)
```

### 栈管理

操作数栈是一个 `[]value.Value` 切片，满时容量翻倍：

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

注意：栈溢出导致的 panic 会被安全网（包裹 `vm.run()` 的 `defer/recover`）捕获，因此会转化为错误返回——永远不会崩溃宿主进程。

### 帧管理与帧池化

每次函数调用都会创建一个 `Frame`：

```go
// vm/frame.go
type Frame struct {
    fn        *bytecode.CompiledFunction  // 哪个函数
    ip        int                         // 指令指针
    basePtr   int                         // 本次调用的操作数栈基址
    locals    []value.Value               // 局部变量数组
    intLocals []int64                     // 整数特化影子数组
    freeVars  []*value.Value              // 闭包捕获（指针！）
    defers    []DeferInfo                 // 延迟调用（LIFO）
    addrTaken bool                        // OpAddr 指向过 locals 则为 true
}
```

#### 帧池化

没有池化时，每次函数调用都会在堆上分配一个 `Frame` + `[]value.Value`。对于像斐波那契这样的递归函数，这意味着数百万次分配。

`framePool` 是一个简单的 LIFO 栈，用于回收帧：

```go
// vm/frame.go
func (p *framePool) get(fn *bytecode.CompiledFunction, basePtr int, freeVars []*value.Value) *Frame {
    n := len(p.frames)
    if n > 0 {
        f = p.frames[n-1]
        p.frames = p.frames[:n-1]
        // 如果 locals 切片容量足够则复用
        if cap(f.locals) >= fn.NumLocals {
            f.locals = f.locals[:fn.NumLocals]
            for i := range f.locals {
                f.locals[i] = value.Value{}  // 清零以保证正确性
            }
        } else {
            f.locals = make([]value.Value, fn.NumLocals)
        }
        // ... 重置所有其他字段 ...
    } else {
        f = &Frame{locals: make([]value.Value, fn.NumLocals)}
    }
    return f
}

func (p *framePool) put(f *Frame) {
    if f.addrTaken {
        return  // 闭包可能持有活跃引用——不回收
    }
    f.fn = nil
    f.freeVars = nil
    p.frames = append(p.frames, f)
}
```

**`addrTaken` 守卫**：当 `OpAddr` 创建了指向帧局部变量的指针时，闭包或延迟函数可能持有这些槽位的引用。回收帧会损坏这些指针。因此 `addrTaken` 帧会交给 GC 处理。

**效果**：Fib(25) 从 ~728K 次分配降至 **7 次分配**。帧池吸收了所有递归开销。

### 分发循环

`vm/run.go` 中的分发循环是性能关键的热路径。它使用多种技术来最大化吞吐量：

#### 寄存器提升

```go
func (v *vm) run() (value.Value, error) {
    // 将热字段提升为局部变量以优化寄存器分配。
    // Go 编译器会在迭代间将这些保持在 CPU 寄存器中。
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

通过将 `v.stack`、`v.sp`、`frame.locals` 等复制到局部变量中，Go 编译器可以将它们保持在 CPU 寄存器中。否则，每条指令都需要通过两次指针间接寻址来访问 `v.stack[v.sp]`。

#### 热/冷路径分离

`run()` 中的主 `switch` 直接内联了**约 60 个热操作码**。较少使用的操作码会落入 `executeOp()`：

```go
    for v.fp > 0 {
        op := bytecode.OpCode(ins[frame.ip])
        frame.ip++

        switch op {
        case bytecode.OpLocal:       // 内联——热路径
            idx := readU16()
            stack[sp] = locals[idx]
            sp++
            continue

        case bytecode.OpAdd:         // 内联——热路径
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

        // ... 约 58 个更多热操作码 ...

        default:
            // 冷路径：同步 sp，调用 executeOp
            v.sp = sp
            if err := v.executeOp(op, frame); err != nil {
                return value.MakeNil(), err
            }
            sp = v.sp
            stack = v.stack
        }
    }
```

#### 整数快速路径

每个算术和比较操作都会先检查两个操作数是否都是 `KindInt`。如果是（在规则引擎中是常见情况），则进行原始 `int64` 运算，零开销：

```go
case bytecode.OpLess:
    sp--
    b := stack[sp]
    sp--
    a := stack[sp]
    if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
        stack[sp] = value.MakeBool(a.RawInt() < b.RawInt())  // 快速路径
    } else {
        stack[sp] = value.MakeBool(a.Cmp(b) < 0)             // 泛型路径
    }
    sp++
    continue
```

#### Context 检查

VM 每 1024 条指令通过位掩码检查一次取消：

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

位掩码避免了取模运算，带 `default` 的 `select` 是非阻塞的 channel 检查。

### 执行过程详解

让我们逐步跟踪 `ProcessOrder(100.0, 3, "HALF")` 的执行过程。仅展示关键操作（为简洁省略部分 SetLocal）：

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

结果：`(150.0, true)`——订单有效，使用 50% 优惠券后总价 = $150。

---

### Panic、Defer 和 Recover

Gig 忠实地实现了 Go 的 panic/defer/recover 语义，包括嵌套 panic（在延迟函数中 panic）。

#### 数据结构

```go
// vm/frame.go
type DeferInfo struct {
    fn      *bytecode.CompiledFunction  // 要调用的函数
    args    []value.Value               // 捕获的参数
    closure *Closure                    // 闭包（间接 defer 场景）
}

// vm/vm.go — panic 状态
type panicState struct {
    panicking bool
    panicVal  value.Value
}

type vm struct {
    // ...
    panicking  bool              // panic 进行中
    panicVal   value.Value       // 当前 panic 值
    panicStack []panicState      // 保存的 panic（嵌套 panic 场景）
    deferDepth int               // defer 执行的嵌套层级
}
```

#### 工作原理

1. **`OpDefer`**：将函数 + 参数捕获到 `frame.defers`
2. **`OpRunDefers`**：正常返回路径——按 LIFO 顺序执行 defer，每个在子 VM 中运行
3. **`OpPanic`**：设置 `v.panicking = true` 和 `v.panicVal`
4. **Panic 处理器**（分发循环顶部）：当 `v.panicking` 为 true 时：
   - 按 LIFO 顺序遍历 `frame.defers`
   - 在每个延迟调用之前，将当前 panic 状态压入 `panicStack`
   - 通过递归 `v.run()` 执行 defer
   - defer 执行后，检查 `recover()` 是否清除了保存的状态
   - 如果已恢复：以正常模式继续执行剩余 defer
   - 如果未恢复：将 panic 向调用者帧传播
5. **`OpRecover`**：弹出 `panicStack` 栈顶，清除 `panicking`，返回值

#### 具体示例

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

执行 `SafeDivide(10, 0)`：

```
1. 为 SafeDivide 创建帧
   locals[0] = Int(10), locals[1] = Int(0)

2. OpDefer: 将匿名闭包捕获到 frame.defers
   defers = [{fn: anon$1, args: [], closure: {FreeVars: [&result, &err]}}]

3. OpDiv: Int(10) / Int(0)
   → Go 层面的 panic: "runtime error: integer divide by zero"
   → 安全网捕获: v.panicking = true, v.panicVal = String("integer divide by zero")

4. Panic 处理器激活（分发循环顶部）：
   - 保存 panic 状态: panicStack = [{panicking:true, val:"integer divide by zero"}]
   - 清除 v.panicking
   - 通过递归 v.run() 调用 anon$1

5. 在 anon$1 内部：
   - OpRecover: 弹出 panicStack，发现 panicking=true
     → 清除保存的状态（panicking=false）
     → 返回 String("integer divide by zero")
   - fmt.Errorf 包装它 → 通过自由变量写入 &err
   - anon$1 返回

6. 回到 panic 处理器：
   - 检查保存的状态：panicking 为 false → 已恢复！
   - 从 ResultAllocSlots 读取命名返回值
   - 返回 (0, error("caught: integer divide by zero"))
```

`ResultAllocSlots` 机制至关重要：在 Go 中，`defer` 可以修改命名返回值。编译器记录哪些局部变量槽位对应命名返回值，恢复路径通过解引用它们来获取最终值。

#### 安全网

所有这些都包裹在 Go 层面的 `defer/recover` 中：

```go
// vm/vm.go
func (v *vm) Execute(funcName string, ctx context.Context, args ...value.Value) (result value.Value, err error) {
    // 安全网：捕获 VM 执行中的 Go 层面 panic
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

这确保即使被解释的代码触发了宿主层面的 panic（nil map 写入、切片越界、类型断言失败），宿主进程也永远不会崩溃。Panic 会被捕获并作为 error 返回。

---

### Goroutine 支持

```go
// 被解释的代码中：
go processItem(item)
```

VM 通过 `OpGoCall` 处理 `go` 语句：

```go
// vm/goroutine.go
func (v *vm) newChildVM() *vm {
    child := &vm{
        program:      v.program,
        stack:        make([]value.Value, initialStackSize),  // 全新的栈
        frames:       make([]*Frame, initialFrameDepth),
        globalsPtr:   v.globalsPtr,       // 通过指针共享全局变量！
        ctx:          v.ctx,              // 共享 context
        extCallCache: v.extCallCache,     // 共享缓存
        goroutines:   v.goroutines,       // 共享追踪器
    }
    if child.globalsPtr == nil {
        child.globalsPtr = &v.globals     // 父级全局变量变为共享
    }
    return child
}
```

**关键设计决策**：
- 每个 goroutine 获得**全新的栈**（无竞争）
- 全局变量通过**指针共享**（正确的 Go 语义）
- 外部调用缓存是**共享的**（通过 RWMutex 保证线程安全）
- Context 是**共享的**（取消操作传播到所有 goroutine）

`GoroutineTracker` 防止 goroutine 失控创建：

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

## 外部包集成

### 注册机制

外部 Go 包必须在编译之前注册。注册通过 `stdlib/packages/` 中代码生成的文件在 `init()` 时完成：

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

每个注册的包提供：
- **函数值**用于 `reflect.Call`（慢路径）
- **DirectCall 封装**用于零反射调用（快速路径）
- **类型信息**用于类型检查和运行时类型断言
- **方法 DirectCall**用于零反射方法分发

### DirectCall：零反射函数调用

Gig 中最大的性能提升来自 **DirectCall 封装**——代码生成的类型化封装函数，完全绕过 `reflect.Call`。

#### reflect.Call 的问题

```go
// 慢路径：reflect.Call
fn := reflect.ValueOf(strings.Contains)
args := []reflect.Value{
    reflect.ValueOf("hello world"),
    reflect.ValueOf("world"),
}
result := fn.Call(args)  // ~400ns：反射开销，内存分配
```

#### DirectCall：解决方案

```go
// 由 gig gen 生成——零反射
func direct_strings_Contains(args []value.Value) value.Value {
    a0 := args[0].String()   // 直接字段访问，无 reflect
    a1 := args[1].String()
    return value.MakeBool(strings.Contains(a0, a1))  // 直接 Go 调用
}
```

这个封装函数：
1. 直接从 `value.Value` 中提取带类型的参数（通过 `String()`、`Int()` 等）
2. 直接调用真正的 Go 函数（无 `reflect.ValueOf`，无 `reflect.Call`）
3. 将结果包装为 `value.Value`（通过 `MakeBool`、`MakeInt` 等）

**效果**：对于典型的标准库函数，速度约为 `reflect.Call` 的 5 倍。

#### 分发流程

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

内联缓存（`externalCallCache`）确保函数解析在每个 program 的生命周期内只发生一次。首次调用后，后续调用只需一次指针解引用 + 函数调用。

```go
// vm/call.go — callExternal 快速路径
func (v *vm) callExternal(funcIdx, numArgs int) error {
    // 弹出参数
    args := make([]value.Value, numArgs)
    for i := numArgs - 1; i >= 0; i-- {
        args[i] = v.pop()
    }

    // 内联缓存查找（读路径用 RLock）
    v.extCallCache.mu.RLock()
    cacheEntry := v.extCallCache.cache[funcIdx]
    v.extCallCache.mu.RUnlock()

    if cacheEntry == nil {
        // 首次调用：解析并缓存（写锁）
        v.extCallCache.mu.Lock()
        cacheEntry = v.resolveExternalFunc(funcIdx)
        v.extCallCache.cache[funcIdx] = cacheEntry
        v.extCallCache.mu.Unlock()
    }

    // 快速路径：DirectCall
    if cacheEntry.directCall != nil {
        result := cacheEntry.directCall(args)
        v.push(result)
        return nil
    }

    // 慢路径：reflect.Call
    return v.callExternalReflect(cacheEntry, args)
}
```

### 外部调用的闭包转换

当被解释的代码将闭包传递给 Go 标准库函数（例如 `sort.Slice` 的比较函数）时，Gig 必须将闭包转换为真正的 Go 函数：

```go
// 被解释的代码：
sort.Slice(items, func(i, j int) bool {
    return items[i].Price < items[j].Price
})
```

`sort.Slice` 函数期望一个 `func(int, int) bool`——它无法接受 `*vm.Closure`。Gig 使用 `reflect.MakeFunc` 创建一个真正的 Go 函数，当被调用时，创建一个临时 VM 并执行闭包的字节码：

```go
// vm/closure.go — Closure implements value.ClosureExecutor
func (c *Closure) Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value {
    // 创建临时 VM 来执行闭包
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

## Init 与执行流程

### 完整生命周期

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

### 无状态模式 vs 有状态模式

**无状态（默认）**：每次 `Run()` 从 `init()` 后的全局变量快照开始。调用结束后对全局变量的修改会被丢弃。这对并发调用是安全的。

```go
prog, _ := gig.Build(`
    var counter int
    func Increment() int {
        counter++
        return counter
    }
`)
prog.Run("Increment") // 返回 1
prog.Run("Increment") // 返回 1（全局变量被重置！）
```

**有状态**（`WithStatefulGlobals()`）：全局变量在调用间持久化。调用通过互斥锁串行化。

```go
prog, _ := gig.Build(`
    var counter int
    func Increment() int {
        counter++
        return counter
    }
`, gig.WithStatefulGlobals())
prog.Run("Increment") // 返回 1
prog.Run("Increment") // 返回 2（全局变量持久化！）
```

```go
// runner/runner.go — 有状态执行
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
    // 无状态：无需加锁
    v := r.vmPool.Get()
    result, err := v.ExecuteWithValues(funcName, ctx, args)
    r.vmPool.Put(v)
    return result, err
}
```

---

## 安全模型

Gig 专为**沙箱执行**不受信任的代码而设计：

### 编译时检查

| 检查项 | 目的 |
|---|---|
| 禁止 `import "unsafe"` | 防止内存损坏 |
| 禁止 `import "reflect"` | 防止类型系统被绕过 |
| 禁止 `panic()`（默认） | 防止通过未恢复的 panic 进行 DoS 攻击。可通过 `WithAllowPanic()` 启用 |
| 自动导入仅限已注册的包 | 代码无法导入任意包 |

### 运行时检查

| 检查项 | 目的 |
|---|---|
| 每 1024 条指令检查 Context 取消 | 防止无限循环 |
| 栈溢出检测（最大 1M 槽位） | 防止内存耗尽 |
| 调用栈深度限制（1024 帧） | 防止栈溢出 |
| Goroutine 限制（默认 10K） | 防止 goroutine 炸弹 |
| 安全网 `defer/recover` | 宿主层面 panic → error 返回 |

### 沙箱注册表

为获得最大隔离性，使用一个初始为空的沙箱注册表：

```go
reg := gig.NewSandboxRegistry()
// 只暴露你想要的：
// reg.RegisterPackage("strings", "strings")
// （或者什么都不注册——纯计算模式）

prog, _ := gig.Build(untrustedCode, gig.WithRegistry(reg))
```

---

## 性能

### 关键优化总结

| 优化 | 技术 | 效果 |
|---|---|---|
| 帧池化 | LIFO 帧回收器 | Fib(25): 728K → 7 次分配 |
| Value 标记联合体 | 32 字节内联基本类型 | int/float/bool 零 GC |
| DirectCall 封装 | 代码生成的类型化封装 | 比 reflect.Call 快约 5 倍 |
| 预烘焙常量 | `[]value.Value` 编译时一次性构建 | 消除逐指令的 `FromInterface` |
| 整数特化 | `intLocals []int64` 影子数组 | 4 倍缓存利用率（8B vs 32B） |
| 超级指令 | 融合操作码（17 种模式） | 热循环中分发次数减少 3-4 倍 |
| 寄存器提升 | 栈/sp/locals 存入 Go 局部变量 | 更优的 CPU 寄存器分配 |
| 内联缓存 | 每 program 的函数解析缓存 | O(1) 外部调用分发 |
| Slice 融合 | `OpIntSliceGet/Set` | `[]int` 访问从 5 条指令降至 1 条 |
| 移动融合 | `OpIntMoveLocal` | 消除拷贝操作的栈往返 |

### 优化效果叠加

以整数密集型循环（如冒泡排序）为例：

```
基线（朴素字节码）：
    LOCAL 0          ; 1 次分发 + 栈写入
    LOCAL 1          ; 1 次分发 + 栈写入
    LESS             ; 1 次分发 + 2 次栈读取 + kind 检查 + 比较 + 栈写入
    JUMPFALSE off    ; 1 次分发 + 栈读取 + 分支
    ─────────────────
    总计：4 次分发，6 次栈操作，1 次 kind 检查

窥孔融合后：
    LessLocalLocalJumpFalse 0 1 off   ; 1 次分发，2 次局部变量读取，比较，分支
    ─────────────────
    总计：1 次分发，0 次栈操作，1 次 kind 检查

整数特化后：
    IntLessLocalLocalJumpFalse 0 1 off ; 1 次分发，2 次 intLocal 读取，比较，分支
    ─────────────────
    总计：1 次分发，0 次栈操作，0 次 kind 检查，8B 操作数
```

这实现了**分发次数减少 4 倍**，且操作数现在存储在 8 字节数组中而非 32 字节的 Value 中。

---

来自：youngjin，2026 年 3 月
