# 整数特化与 Value 缩减 — 优化报告

## 摘要

本文档涵盖了两个主要优化阶段，它们共同在所有基准测试中实现了 **38–62% 的加速**：**Value 结构体缩减**（56B → 32B）和**整数特化局部变量**（`intLocals []int64`）。这些阶段解决了之前报告中确定的两个最大剩余性能差距——过大的 `Value` 结构体和在内层循环中操作 32 字节值的开销。

### 最终结果（AMD EPYC 9754，linux/amd64）

| 基准测试 | 优化前 | 优化后 | 加速比 | vs Yaegi |
|---|---|---|---|---|
| Fib(25) | 32.2 ms | **19.8 ms** | **1.63x** | **Gig 快 5.6 倍** |
| ArithSum(1K) | 137 μs | **52.6 μs** | **2.60x** | Yaegi 快 1.28 倍 |
| BubbleSort(100) | 3.93 ms | **2.37 ms** | **1.66x** | Yaegi 快 1.9 倍 |
| Sieve(1000) | 691 μs | **388 μs** | **1.78x** | Yaegi 快 1.9 倍 |
| ClosureCalls(1K) | 535 μs | **374 μs** | **1.43x** | **Gig 快 2.7 倍** |

### 累积改进（从项目开始）

| 基准测试 | 原始 | 当前 | 总加速比 |
|---|---|---|---|
| Fib(25) | 169 ms | 19.8 ms | **8.5x** |
| ArithSum(1K) | 311 μs | 52.6 μs | **5.9x** |
| BubbleSort(100) | 10.3 ms | 2.37 ms | **4.3x** |
| Sieve(1000) | 1,681 μs | 388 μs | **4.3x** |
| ClosureCalls(1K) | 964 μs | 374 μs | **2.6x** |

---

## 第一部分：Value 结构体缩减（56B → 32B）

### 问题

`Value` 结构体——VM 中所有运行时值的通用表示——为 56 字节：

```go
// 优化前：56 字节
type Value struct {
    kind Kind              // 1 字节 + 7 字节填充
    obj  any               // 16 字节（接口：类型指针 + 数据指针）
    num  int64             // 8 字节
    str  string            // 16 字节（字符串头：数据指针 + 长度）
}
```

每次压栈、弹栈、局部变量读写和函数参数传递都复制 56 字节。在每次迭代执行 3 条超级指令、每条触及 2-3 个局部变量的紧密循环中，这意味着每次迭代移动 **336–504 字节**的 Value 数据。56 字节的大小还导致较差的缓存利用率——一个缓存行（64 字节）只能容纳 1 个 Value，浪费 8 字节。

### 解决方案

将 `str string` 字段合并到 `obj any` 字段中，因为字符串和其他堆分配类型永远不会同时使用：

```go
// 优化后：32 字节
type Value struct {
    kind Kind    // 1 字节 + 7 字节填充
    num  int64   // 8 字节：bool (0/1)、int、uint 位模式、float64 位模式
    obj  any     // 16 字节：string、complex128、reflect.Value 或 nil
}
```

**关键设计决策：**

1. **`num int64` 存储所有标量类型** — 布尔值为 0/1，int 原样存储，uint 通过位转换，float64 通过 `math.Float64bits`/`math.Float64frombits`。这使 `RawInt() int64` 成为简单的字段访问：`return v.num`。

2. **`obj any` 是通用堆字段** — 字符串作为 `string` 值装箱到 `any` 接口中存储。复数、reflect.Value、映射、切片、通道等都放在这里。`kind` 字段在运行时进行区分。

3. **`RawInt()` 和 `RawBool()` 是无检查的** — 这些访问器跳过 kind 检查，依赖 SSA 类型保证。`RawInt()` 直接返回 `v.num`；`RawBool()` 返回 `v.num != 0`。两者都足够小，Go 编译器可以内联。

### 影响

| 指标 | 优化前 | 优化后 |
|---|---|---|
| Value 大小 | 56 字节 | 32 字节 |
| 每缓存行 Value 数 | 1（浪费 8B） | 2（精确匹配） |
| 每循环迭代移动字节数（ArithSum） | ~504 字节 | ~288 字节 |
| 每 1000 个局部变量的内存 | 56 KB | 32 KB |

**修改的文件：**
- `value/value.go` — 重构 `Value`，更新所有构造函数（`MakeInt`、`MakeString`、`MakeFloat` 等）
- `value/accessor.go` — 更新 `String()`、`Int()`、`Float64()`、`ToReflectValue()` 从 `obj` 读取
- `value/arithmetic.go` — 更新 `Add`、`Sub`、`Mul`、`Div`、`Cmp` 使用新字段布局
- `value/convert.go` — 更新 `Convert()` 和 `FromInterface()` 适配新字段布局
- `value/container.go` — 更新容器操作适配新的字符串/切片存储
- `value/value_test.go` — 更新测试适配新结构体布局

---

## 第二部分：整数特化局部变量（`intLocals []int64`）

### 问题

即使将 `Value` 缩减到 32 字节后，整数密集的内层循环仍然通过 `locals[]` 为每次读写移动 32 字节的结构体。对于像 `sum += i; i++` 这样的简单循环，每次迭代读写多个 `value.Value` 结构体（每个 32 字节），而实际有效载荷只是一个 `int64`（8 字节）。数据移动有 4 倍的开销。

此外，每次整数运算还必须：
1. 检查 `Kind() == KindInt`（分支）
2. 通过 `RawInt()` 提取（字段访问）
3. 计算结果
4. 通过 `MakeInt()` 装箱（构造 32 字节结构体）
5. 将 32 字节结构体存回

当编译器可以静态证明变量始终是 `int` 时，步骤 1、2、4 和 5 都是纯开销。

### 解决方案：影子 `int64` 数组

为包含整数类型局部变量的函数，在每个帧中添加并行的 `intLocals []int64` 数组。整数特化操作码（`OpInt*`）直接在这个 8 字节数组上操作，在热路径中完全绕过 32 字节的 `Value`。

```
┌─────────────────────────────────────────────┐
│ Frame                                        │
│                                              │
│  locals []value.Value   (每槽位 32 字节)      │  ← 通用操作码使用
│  ┌──────┬──────┬──────┬──────┐              │
│  │ V[0] │ V[1] │ V[2] │ V[3] │              │
│  └──────┴──────┴──────┴──────┘              │
│                                              │
│  intLocals []int64       (每槽位 8 字节)      │  ← OpInt* 操作码使用
│  ┌──────┬──────┬──────┬──────┐              │
│  │ i[0] │ i[1] │ i[2] │ i[3] │              │
│  └──────┴──────┴──────┴──────┘              │
│                                              │
│  两个数组通过双写保持同步                      │
└─────────────────────────────────────────────┘
```

### 双写不变量

**每次写入整数特化局部变量都必须同时更新 `intLocals[idx]` 和 `locals[idx]`。**

此不变量存在是因为非特化代码（通用 `OpLocal`、闭包、函数返回、调试器）从 `locals[]` 读取。如果我们只写入 `intLocals[]`，这些代码路径会看到过期数据。

```go
// OpIntLocalConstAddSetLocal — 热路径
r := intLocals[idxA] + intConsts[idxB]   // 纯 int64 算术
intLocals[idxC] = r                        // 快速 8 字节写入
locals[idxC] = value.MakeInt(r)            // 同步 32 字节影子
```

双写每次操作多花费一次 `MakeInt` + 32 字节存储，但关键好处是**读取**（`OpIntLocal`）只触及 8 字节的 `intLocals[]` 数组，且比较+跳转融合（`OpIntLess*`）跳过了 `Kind()` 检查和 `RawInt()` 提取。

### 整数特化操作码

13 个新的 `OpInt*` 操作码在 `intLocals []int64` 和 `intConsts []int64` 上操作：

**融合算术（零栈流量，8 字节操作数）：**

| 操作码 | 语义 | 宽度 |
|---|---|---|
| `OpIntLocalConstAddSetLocal` | `intLocals[C] = intLocals[A] + intConsts[B]` | 7B |
| `OpIntLocalConstSubSetLocal` | `intLocals[C] = intLocals[A] - intConsts[B]` | 7B |
| `OpIntLocalLocalAddSetLocal` | `intLocals[C] = intLocals[A] + intLocals[B]` | 7B |

**融合比较+分支（零栈流量，无 kind 检查）：**

| 操作码 | 语义 | 宽度 |
|---|---|---|
| `OpIntLessLocalConstJumpTrue` | `if intLocals[A] < intConsts[B] { goto off }` | 7B |
| `OpIntLessLocalConstJumpFalse` | `if intLocals[A] >= intConsts[B] { goto off }` | 7B |
| `OpIntLessEqLocalConstJumpTrue` | `if intLocals[A] <= intConsts[B] { goto off }` | 7B |
| `OpIntLessEqLocalConstJumpFalse` | `if intLocals[A] > intConsts[B] { goto off }` | 7B |
| `OpIntLessLocalLocalJumpTrue` | `if intLocals[A] < intLocals[B] { goto off }` | 7B |
| `OpIntLessLocalLocalJumpFalse` | `if intLocals[A] >= intLocals[B] { goto off }` | 7B |
| `OpIntGreaterLocalLocalJumpTrue` | `if intLocals[A] > intLocals[B] { goto off }` | 7B |

**桥接操作码（在 intLocals 和栈之间同步）：**

| 操作码 | 语义 | 宽度 |
|---|---|---|
| `OpIntSetLocal` | `intLocals[idx] = pop().RawInt(); locals[idx] = pop()` | 3B |
| `OpIntLocal` | `push(MakeInt(intLocals[idx]))` | 3B |

**寄存器移动（融合 phi 移动）：**

| 操作码 | 语义 | 宽度 |
|---|---|---|
| `OpIntMoveLocal` | `intLocals[dst] = intLocals[src]; locals[dst] = locals[src]` | 5B |

### `intConsts []int64` 常量池

并行常量池在现有的 `PrebakedConstants []value.Value` 池旁存储预提取的 `int64` 值。这允许 `OpInt*` 比较和算术操作码将常量作为原始 `int64` 读取，无需任何 `RawInt()` 调用或 kind 检查。

```go
// 在 Program 中（bytecode/bytecode.go）
type Program struct {
    PrebakedConstants []value.Value   // 用于通用操作码
    IntConstants      []int64         // 用于 OpInt* 操作码（并行数组）
    // ...
}
```

---

## 编译流水线

优化流水线在字节码生成后按三个顺序遍运行：

```
原始字节码
     │
     ▼
┌─────────────────────────────────┐
│ 第 1 遍：窥孔融合               │   optimizeBytecode()
│  • 6 指令融合                    │   将常见模式融合为
│  • 4 指令融合                    │   超级指令（Value 类型）
│  • 3 指令融合                    │
│  • 2 指令融合                    │
│  • 空跳转消除                    │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│ 第 2 遍：整数特化               │   intSpecialize()
│  • 两遍算法                     │   将超级指令升级为
│  • 候选识别                     │   OpInt* 变体 + 插入
│  • 操作码升级 + 桥接            │   桥接指令
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│ 第 3 遍：整数移动融合           │   fuseIntMoves()
│  • OpIntLocal + OpIntSetLocal   │   将 phi 移动对融合为
│    → OpIntMoveLocal             │   单个寄存器移动
└────────────┬────────────────────┘
             │
             ▼
  优化后字节码
```

### 第 1 遍：窥孔超级指令融合

`compiler/optimize.go` 中的 `optimizeBytecode()` 扫描指令流中的已知模式并将其融合为单个超级指令。模式按最长优先检查以避免部分匹配。

**6 指令融合（16B → 7B，6 次分发 → 1 次）：**

处理 SSA 比较模式，其中比较结果存储到临时变量，然后加载并分支：

```
LOCAL(A) CONST(B) LESS SETLOCAL(X) LOCAL(X) JUMPFALSE(off)
└──────────────────────────────────────────────────────────┘
  → OpLessLocalConstJumpFalse(A, B, off)
```

优化器验证 `SETLOCAL(X)` 和 `LOCAL(X)` 引用同一槽位（SSA 临时变量），确认这是纯比较-分支，无副作用。

**4 指令融合（10B → 7B，4 次分发 → 1 次）：**

算术：`LOCAL(A) LOCAL(B) ADD SETLOCAL(C)` → `OpLocalLocalAddSetLocal(A,B,C)`
比较：`LOCAL(A) CONST(B) LESS JUMPTRUE(off)` → `OpLessLocalConstJumpTrue(A,B,off)`

**3 指令融合（7B → 5B，3 次分发 → 1 次）：**

`LOCAL(A) LOCAL(B) ADD` → `OpAddLocalLocal(A,B)`（结果留在栈上）

**2 指令融合（4B → 3B，2 次分发 → 1 次）：**

`ADD SETLOCAL(A)` → `OpAddSetLocal(A)`（弹出两个，相加，存储）

**空跳转消除：**

`JUMP(target)` 中 target 是下一条指令时被完全移除。这些来自相邻基本块之间无条件分支的 SSA 编译。

### 跳转目标修正

所有融合改变了指令大小，使跳转目标失效。优化器使用 `rewrite` 列表记录每次融合的 `(start, end, newBytes)`。收集所有重写后，`applyRewrites()` 重建字节码，`fixJumpTargets()` 使用偏移映射表调整每条跳转指令。

### 第 2 遍：整数特化（`intSpecialize`）

这是一个**两遍**算法，将符合条件的超级指令升级为其 `OpInt*` 变体。

**为什么需要两遍？** 单遍无法正确处理桥接指令。考虑：

```
OpSetLocal(x)              ← 出现在任何 OpInt* 指令之前
OpLocalConstAddSetLocal(x, 1, x)   ← 可升级为 OpInt*
```

在单遍中，当我们看到 `OpSetLocal(x)` 时，我们还不知道 `x` 将参与整数特化操作。我们会错过桥接升级。两遍方法解决了这个问题：

**第 1 遍 — 候选识别：**

扫描整个字节码，查找所有局部变量和常量操作数都被静态确认为 `int` 类型（通过 SSA 类型信息构建的 `localIsInt[]` 和 `constIsInt[]`）的超级指令。对于每条符合条件的指令，在 `intUsed []bool` 集合中标记涉及的局部变量索引。

```go
// 示例：OpLocalConstAddSetLocal(A, B, C) 在以下条件下符合资格
// localIsInt[A] && constIsInt[B] && localIsInt[C]
// → 标记 intUsed[A] = true, intUsed[C] = true
```

如果没有符合条件的指令，返回 `(code, false)` — 不分配 `intLocals` 数组。

**第 2 遍 — 操作码升级 + 桥接插入：**

两种类型的变换：

1. **超级指令升级** — 原地替换操作码字节（操作数布局相同）：
   ```
   OpLocalConstAddSetLocal → OpIntLocalConstAddSetLocal
   OpLessLocalConstJumpFalse → OpIntLessLocalConstJumpFalse
   // ...（所有符合条件的超级指令）
   ```

2. **桥接插入** — 对于引用 `intUsed` 集合中局部变量的任何 `OpSetLocal`/`OpLocal`：
   ```
   OpSetLocal(idx) → OpIntSetLocal(idx)   // 双写到两个数组
   OpLocal(idx)    → OpIntLocal(idx)      // 从 intLocals 读取
   ```

桥接确保非特化代码路径（函数参数、同一变量上的通用算术）保持 `intLocals` 与 `locals` 同步。

### 第 3 遍：整数移动融合（`fuseIntMoves`）

在 `intSpecialize` 创建 `OpIntLocal` + `OpIntSetLocal` 对（通常来自 SSA phi 移动模式）后，此遍将其融合：

```
OpIntLocal(src) OpIntSetLocal(dst)
→ OpIntMoveLocal(src, dst)
```

这消除了压入+弹出周期（加载到栈，然后从栈存储），替换为直接的寄存器到寄存器移动：

```go
intLocals[dst] = intLocals[src]   // 8 字节复制
locals[dst] = locals[src]          // 32 字节复制
```

节省：6 字节 → 5 字节，2 次分发 → 1 次分发，并消除 2 次栈操作。

---

## 运行时支持

### 帧分配（`vm/frame.go`）

```go
type Frame struct {
    fn        *bytecode.CompiledFunction
    ip        int
    basePtr   int
    locals    []value.Value    // 每槽位 32 字节 — 通用
    intLocals []int64          // 每槽位 8 字节 — int 影子
    freeVars  []*value.Value
    defers    []DeferInfo
    addrTaken bool
}
```

`intLocals` 仅在 `fn.HasIntLocals == true`（由 `intSpecialize` 设置）时分配。`framePool` 在回收帧时复用两个切片，避免每次调用的分配。

**`newFrame()` 中的参数镜像：**

```go
if fn.HasIntLocals {
    f.intLocals = make([]int64, fn.NumLocals)
    for i, arg := range args {
        if i < fn.NumLocals {
            f.intLocals[i] = arg.RawInt()  // 镜像 int 参数
        }
    }
}
```

### 函数调用（`vm/call.go`）

`callCompiledFunction` 和 `callFunction`（用于闭包）都镜像参数：

```go
// 在 callCompiledFunction 中：
v := vm.pop()
frame.locals[i] = v
if intL != nil {
    intL[i] = v.RawInt()    // 将每个参数镜像到 intLocals
}
```

这确保进入函数时 `intLocals` 正确初始化，即使调用者使用了通用（非整数特化）代码。

### VM 分发（`vm/run.go`）

`run()` 循环将 `intLocals` 和 `intConsts` 与 `stack`、`sp`、`locals` 和 `prebaked` 一起提升到 Go 局部变量中：

```go
intLocals := frame.intLocals     // 可能为 nil
intConsts := vm.program.IntConstants
```

所有 `OpInt*` 处理器直接在这些局部变量上操作，避免重复的 `frame.intLocals` 解引用。循环在任何帧变更（函数调用、返回）后重新加载这些。

**示例热路径处理器：**

```go
case bytecode.OpIntLocalConstAddSetLocal:
    idxA := readU16(code, ip+1)
    idxB := readU16(code, ip+3)
    idxC := readU16(code, ip+5)
    ip += 7
    r := intLocals[idxA] + intConsts[idxB]   // 纯 int64 加法
    intLocals[idxC] = r                        // 8 字节写入
    locals[idxC] = value.MakeInt(r)            // 32 字节同步
    continue
```

---

## 具体示例：ArithSum 内层循环

**源码：**
```go
sum := 0
for i := 0; i < 1000; i++ {
    sum += i
}
```

### 优化前（Value 类型超级指令，3 次分发）：
```
OpLocalLocalAddSetLocal(sum, i, sum)   — 读 2×32B，写 1×32B，kind 检查
OpLocalConstAddSetLocal(i, 1, i)       — 读 1×32B + 1×32B 常量，写 1×32B
OpLessLocalConstJumpTrue(i, 1000, -)   — 读 1×32B + 1×32B 常量，2 次 kind 检查
```
**每次迭代：** 3 次分发，~288 字节移动，4 次 kind 检查，3 次 `RawInt()` 调用，3 次 `MakeInt()` 调用。

### 优化后（整数特化，3 次分发）：
```
OpIntLocalLocalAddSetLocal(sum, i, sum) — 读 2×8B，写 1×8B + 1×32B 同步
OpIntLocalConstAddSetLocal(i, 1, i)     — 读 1×8B + 1×8B，写 1×8B + 1×32B 同步
OpIntLessLocalConstJumpTrue(i, 1000, -) — 读 1×8B + 1×8B，无写入
```
**每次迭代：** 3 次分发，~112 字节移动，0 次 kind 检查，0 次 `RawInt()` 调用，2 次 `MakeInt()` 调用（仅双写）。

**净节省：** 每次迭代少移动 176 字节，少 4 次分支（kind 检查），少 1 次 `MakeInt()` 调用。在 1000 次迭代中，消除了 176KB 的数据移动和 4000 条分支指令。

---

## 具体示例：Fibonacci 递归

**源码：**
```go
func fib(n int) int {
    if n <= 1 { return n }
    return fib(n-1) + fib(n-2)
}
```

### 关键优化：

`n <= 1` 比较通过 SSA 编译为：
```
LOCAL(n) CONST(1) LESSEQ SETLOCAL(t) LOCAL(t) JUMPFALSE(else)
```

6 指令融合将其简化为：
```
OpLessEqLocalConstJumpFalse(n, 1, else)
```

然后 `intSpecialize` 将其升级为：
```
OpIntLessEqLocalConstJumpFalse(n, 1, else)
```

这条单指令替换了 6 次分发并消除了所有中间栈操作。结合递归调用（在每个帧中将 `n` 镜像到 `intLocals`），整个热路径在 8 字节 `int64` 值上操作。

---

## 正确性保证

### 为什么双写是必要的

考虑一个对同一变量同时拥有整数特化和通用代码路径的函数：

```go
x := computeInt()    // OpIntSetLocal — 写 intLocals[x] 和 locals[x]
if condition {
    useAsInt(x)      // OpIntLocal — 读 intLocals[x] ✓
} else {
    useAsAny(x)      // OpLocal — 读 locals[x] ✓（双写保持同步）
}
```

没有双写，else 分支中的 `OpLocal` 读取会看到过期的 `locals[x]` 值。

### intUsed 局部变量上的非 int 操作

当 `OpAddSetLocal` 操作的局部变量也在 `intUsed` 集合中时，它也必须更新 `intLocals`：

```go
case bytecode.OpAddSetLocal:
    // ... 计算结果 ...
    locals[idx] = result
    if intLocals != nil {
        intLocals[idx] = result.RawInt()  // 保持影子同步
    }
```

### 帧转换

函数调用（`callCompiledFunction`、`callFunction`）和帧创建（`newFrame`）都将 int 参数镜像到 `intLocals`。这确保无论调用者是否使用了整数特化操作码，`intLocals` 都正确初始化。

---

## 修改的文件

| 文件 | 变更 |
|---|---|
| `value/value.go` | 将 Value 从 56B 缩减到 32B（将 `str` 合并到 `obj`） |
| `value/accessor.go` | 更新字符串/类型访问器适配新字段布局 |
| `value/arithmetic.go` | 更新算术操作适配新字段布局 |
| `value/convert.go` | 更新类型转换适配新字段布局 |
| `value/container.go` | 更新容器操作适配新字符串存储 |
| `value/value_test.go` | 更新单元测试适配 32B Value |
| `bytecode/opcode.go` | 添加 13 个 `OpInt*` 操作码 + `OpLessEqLocalConstJumpFalse/True` |
| `bytecode/bytecode.go` | 添加 `HasIntLocals` 标志，`IntConstants []int64` 池 |
| `compiler/compile_func.go` | 将 `intSpecialize()` 和 `fuseIntMoves()` 集成到流水线 |
| `compiler/optimize.go` | 添加 6 指令融合、`intSpecialize()`、`fuseIntMoves()`、空跳转消除 |
| `compiler/compiler.go` | 编译期间构建 `IntConstants` 池 |
| `vm/frame.go` | 为 Frame 添加 `intLocals []int64`，更新池 |
| `vm/run.go` | 添加所有 `OpInt*` 处理器（含双写），提升 `intLocals`/`intConsts` |
| `vm/call.go` | 在 `callCompiledFunction` 和 `callFunction` 中添加 intLocals 镜像 |

---

## 测试

所有优化保持完全向后兼容。完整测试套件通过：

```
ok  gig/bytecode    0.002s
ok  gig/compiler    0.002s
ok  gig/importer    0.003s
ok  gig/tests       0.852s
ok  gig/value       0.003s
ok  gig/vm          0.003s
```

---

## 更新后的架构

```
源代码 (.go)
       │
       ▼
  Go SSA 包 (golang.org/x/tools/go/ssa)
       │
       ▼
  编译器 (gig/compiler)
   ├── 分配局部变量槽位 + 构建 localIsInt/constIsInt 映射
   ├── 编译基本块 → 原始字节码
   ├── 修补跳转目标
   ├── 第 1 遍：窥孔融合（超级指令）
   ├── 第 2 遍：intSpecialize（Value → int64 升级）
   └── 第 3 遍：fuseIntMoves（phi 移动融合）
       │
       ▼
  字节码 (gig/bytecode)
   ├── PrebakedConstants []value.Value   ← 通用常量池
   ├── IntConstants []int64              ← 整数特化常量池
   ├── HasIntLocals bool                 ← 按函数标志
   └── 13 个 OpInt* 操作码 + 桥接
       │
       ▼
  VM 执行 (gig/vm)
   ├── Frame: locals []Value + intLocals []int64（影子）
   ├── 每次 int 局部变量变更的双写不变量
   ├── 寄存器提升分发（stack、sp、locals、intLocals、intConsts）
   ├── 帧池化复用 locals 和 intLocals 切片
   └── 函数入口时参数镜像
```

---

## 剩余优化机会

| 优先级 | 优化 | 预期影响 |
|---|---|---|
| P0 | 直接线程分发（通过汇编的计算跳转） | 所有工作负载 ~1.2–1.5 倍 |
| P1 | 通过逃逸分析消除双写（证明 int 变量的 locals[] 永远不被读取） | 算术循环 ~1.1–1.2 倍 |
| P2 | 原生 `[]bool` 切片用于 Sieve 类工作负载 | Sieve ~1.3 倍 |
| P3 | 闭包分配池化 | 减少 ClosureCalls 中 3,000 allocs/op |
| P4 | Float64 特化（类似于整数特化） | 浮点密集工作负载 ~1.3 倍 |
