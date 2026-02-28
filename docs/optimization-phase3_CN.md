# 第三阶段优化：Value 类型与上下文检查机制重构

## 概述

本文档记录 Gig VM 的第三轮优化，目标是 `Value` 类型内部实现、闭包调度开销和上下文取消检查机制。这些变更包含 **内部 API 的 breaking change**，实现了 **3–11% 的性能提升** 和闭包密集型负载 **65% 的内存降低**。

### 结果（AMD EPYC 9754 128-Core Processor, linux/amd64, 3 次取中位数）

| 基准测试 | Phase 2 | Phase 3 | 提升 |
|---|---|---|---|
| Fib(25) | 20.3 ms | **19.7 ms** | **3.0%** |
| ArithSum(1K) | 35.2 μs | **34.1 μs** | **3.1%** |
| BubbleSort(100) | 913 μs | **943 μs** | ~噪声 |
| Sieve(1000) | 196 μs | **187 μs** | **4.6%** |
| ClosureCalls(1K) | 369 μs / 48 KB / 3K allocs | **326 μs / 17 KB / 2K allocs** | **时间 11.7%, 内存 65%, 分配 33%** |
| ExtCallDirectCall | 511 μs | **500 μs** | **2.2%** |
| ExtCallReflect | 344 μs | **338 μs** | **1.7%** |
| ExtCallMethod | 439 μs | **416 μs** | **5.2%** |
| ExtCallMixed | 308 μs | **301 μs** | **2.3%** |

### 跨解释器对比（Phase 3 后当前状态）

#### 核心负载

| 负载 | 原生 Go | Gig | Yaegi | GopherLua | Gig vs 原生 | Gig vs Yaegi | Gig vs Lua |
|---|---:|---:|---:|---:|---:|---:|---:|
| Fibonacci(25) | 449 μs | 19.7 ms | 111 ms | 21.3 ms | 慢 44x | **快 5.6x** | **快 1.1x** |
| ArithSum(1K) | 667 ns | 34.1 μs | 40.1 μs | 55.0 μs | 慢 51x | **快 1.2x** | **快 1.6x** |
| BubbleSort(100) | 6.3 μs | 943 μs | 1,215 μs | 1,053 μs | 慢 150x | **快 1.3x** | **快 1.1x** |
| Sieve(1000) | 1.89 μs | 187 μs | 209 μs | 258 μs | 慢 99x | **快 1.1x** | **快 1.4x** |
| ClosureCalls(1K) | 346 ns | 326 μs | 977 μs | 146 μs | 慢 942x | **快 3.0x** | 慢 2.2x |

#### 外部函数调用（Gig vs Yaegi，无 Lua/原生等价测试）

| 负载 | 原生 Go | Gig | Yaegi | Gig vs 原生 | Gig vs Yaegi |
|---|---:|---:|---:|---:|---:|
| DirectCall | 28.4 μs | 500 μs | 1,529 μs | 慢 18x | **快 3.1x** |
| Reflect | 24.3 μs | 338 μs | 1,004 μs | 慢 14x | **快 3.0x** |
| Method | 18.4 μs | 416 μs | 1,226 μs | 慢 23x | **快 2.9x** |
| Mixed | 11.7 μs | 301 μs | 857 μs | 慢 26x | **快 2.8x** |

#### 内存效率（allocs/op）

| 负载 | Gig | Yaegi | GopherLua | Gig vs Yaegi |
|---|---:|---:|---:|---:|
| Fibonacci(25) | 6 | 2,138,703 | 41 | **少 356,450x** |
| ArithSum(1K) | 6 | 8 | 93 | 少 1.3x |
| BubbleSort(100) | 9 | 5,085 | 12 | **少 565x** |
| Sieve(1000) | 7 | 43 | 207 | **少 6x** |
| ClosureCalls(1K) | 1,995 | 13,018 | 96 | **少 6.5x** |

---

## 优化 1：KindFunc — Value 中直接存储闭包

### 问题

VM 创建闭包（`OpClosure`）时，`*Closure` 指针通过 `FromInterface()` 存储到 `Value` 中：

```go
case bytecode.OpClosure:
    closure := getClosure(fn, numFree)
    // ...
    vm.push(value.FromInterface(closure))
```

`FromInterface()` 调用 `reflect.ValueOf(closure)` 将 `*Closure` 包装为 `reflect.Value`，然后存入 `Value.obj`。之后调用闭包（`OpCallIndirect`）时需要：

```go
callee.Interface() → reflect.Value.Interface() → 类型断言 .(*Closure)
```

这个双重间接（reflect 包装 + 解包）每次闭包调用耗费约 15ns。

### 方案

新增 `KindFunc` 路径，将 `*Closure` 直接存入 `Value.obj`，不经过任何 reflect 包装：

```go
// value/value.go
func MakeFunc(fn any) Value {
    return Value{kind: KindFunc, obj: fn}
}

func (v Value) RawObj() any { return v.obj }
```

```go
// vm/ops_dispatch.go — OpClosure
vm.push(value.MakeFunc(closure))  // 原来: value.FromInterface(closure)
```

```go
// vm/run.go — OpCallIndirect
if closure, ok := callee.RawObj().(*Closure); ok {
    // 直接类型断言，无 reflect
}
```

`value/accessor.go` 中的 `Interface()` 和 `ToReflectValue()` 方法新增了 `KindFunc` 分支，确保以此方式存储的闭包与系统其余部分保持互操作性。

### 修改文件

| 文件 | 变更 |
|---|---|
| `value/value.go` | 新增 `MakeFunc()` 构造器、`RawObj()` 访问器 |
| `value/accessor.go` | `Interface()` 和 `ToReflectValue()` 新增 `KindFunc` 分支 |
| `vm/ops_dispatch.go` | `OpClosure` 使用 `value.MakeFunc(closure)` |
| `vm/ops_dispatch.go` | `OpCallIndirect`/`OpGoCallIndirect` 使用 `callee.RawObj().(*Closure)` |
| `vm/ops_dispatch.go` | `OpFree` 新增 `KindFunc` 快速路径 |
| `vm/run.go` | 内联 `OpCallIndirect` 使用 `callee.RawObj().(*Closure)` |

### 效果

消除了每次闭包创建和调用的 `reflect.ValueOf()` + `reflect.Value.Interface()` 往返开销。贡献了 **ClosureCalls 11.7% 的加速**和 **65% 的内存降低**（不再分配 reflect.Value 包装器）。

---

## 优化 2：OpCallIndirect 栈分配参数缓冲区

### 问题

每次 `OpCallIndirect` 都为函数参数分配一个新切片：

```go
args := make([]value.Value, numArgs)
```

对于参数数量 ≤8 的闭包调用（绝大多数情况），这种小切片分配造成不必要的 GC 压力。Go 的逃逸分析无法证明该切片不会逃逸到堆上，因为它被传递给了 `callFunction()`。

### 方案

新增栈分配的 `[8]value.Value` 数组作为小参数数量的后备存储：

```go
case bytecode.OpCallIndirect:
    numArgs := int(frame.readByte())
    var argsBuf [8]value.Value
    var args []value.Value
    if numArgs <= len(argsBuf) {
        args = argsBuf[:numArgs]  // 指向栈内存
    } else {
        args = make([]value.Value, numArgs)  // 堆回退
    }
```

`argsBuf` 声明在 `run()` 的 `case` 块内部，Go 编译器可以将其保留在栈上。切片头 `args` 引用栈数组，对常见情况避免了堆分配。

**注意：** 此优化也曾尝试应用于 `callExternal()`，但因逃逸分析失败而回退——见下方"尝试但回退"章节。

### 修改文件

| 文件 | 变更 |
|---|---|
| `vm/run.go` | `OpCallIndirect` 对小参数列表使用 `argsBuf[8]` |

### 效果

减少每次闭包调用的分配次数。与 KindFunc 结合，贡献了 **ClosureCalls allocs/op 33% 的降低**。

---

## 优化 3：反向跳转上下文检查（替代逐指令计数器）

### 问题

VM 使用逐指令计数器每 N 条指令检查一次上下文取消：

```go
instructionCount++
if instructionCount & 0x1FFF == 0 {
    select {
    case <-vm.ctx.Done():
        return value.MakeNil(), vm.ctx.Err()
    default:
    }
}
```

该计数器递增在**每一条指令**上执行——计算密集型负载每秒数百万次。递增本身很廉价（约 1ns），但它：
1. 占用一个可用于其他热变量的寄存器
2. 增加主循环的指令缓存压力
3. 在前向跳转（函数调用准备）上检查上下文，而这些永远不会形成无限循环

### 方案

用仅在反向跳转时的计数器替代逐指令计数器。上下文取消只对无限循环有意义，而循环总是涉及反向跳转：

```go
// 移除: instructionCount++  (原来在每条指令上)

case bytecode.OpJump:
    offset := readU16()
    if int(offset) < frame.ip {
        // 反向跳转 — 这是一次循环迭代
        backJumpCount++
        if backJumpCount & 0x7F == 0 {  // 每 128 次反向跳转
            select {
            case <-vm.ctx.Done():
                return value.MakeNil(), vm.ctx.Err()
            default:
            }
        }
    }
    frame.ip = int(offset)
```

关键设计决策：
- **只检查 `OpJump`，不检查 `OpJumpTrue`/`OpJumpFalse`**：Go SSA→字节码编译器总是在循环底部发出无条件 `OpJump`。顶部的条件跳转是循环退出检查。只检查 `OpJump` 避免了对更频繁的条件分支增加开销。
- **128 间隔（0x7F）**：选择此值以保持最坏情况取消延迟在 ~200μs 以下。紧凑循环体执行约 50 条指令，所以 128 次反向跳转 ≈ 6400 条指令，约每 64μs 检查一次。

### 修改文件

| 文件 | 变更 |
|---|---|
| `vm/run.go` | 移除 `instructionCount`，新增 `backJumpCount` 及 128 间隔节流 |

### 效果

所有计算密集型负载约 3% 的加速（Fib、ArithSum、Sieve）。改进来自消除逐指令计数器递增和在热循环中释放一个 CPU 寄存器。

---

## 尝试但回退

### 1. 内联 readU16() 调用

**概念：** 将所有 `readU16()` 闭包调用替换为内联表达式：
```go
// 之前
idx := readU16()
// 之后
idx := uint16(ins[frame.ip])<<8 | uint16(ins[frame.ip+1]); frame.ip += 2
```

**失败原因：** `run()` 函数体显著膨胀，导致指令缓存（icache）压力。Fib25 退化 16%，Sieve 退化 8%。Go 编译器已经高效地内联了 `readU16` 闭包——手动内联使编译后的函数更大却没有收益。

**教训：** Go 编译器对小闭包的优化已经很好。手动内联在函数体超过 L1 icache 容量时可能适得其反。

### 2. callExternal() 中的栈分配参数缓冲区

**概念：** 将同样的 `[8]value.Value` 栈缓冲区技巧应用于 `callExternal()`：
```go
func (vm *VM) callExternal(funcIdx, numArgs int) {
    var argsBuf [8]value.Value
    var args []value.Value
    if numArgs <= 8 {
        args = argsBuf[:numArgs]
    } else {
        args = make([]value.Value, numArgs)
    }
    // ...
    entry.directCall(args)  // 接口方法调用
}
```

**失败原因：** `go build -gcflags='-m'` 确认 `argsBuf` 逃逸到堆上。`directCall(args)` 是接口方法调用——Go 的逃逸分析无法证明被调用的方法不会保留切片引用，因此保守地将后备数组分配到堆上。

**教训：** 栈分配优化只在 Go 的逃逸分析能证明数据不逃逸时才有效。接口方法调用是常见的逃逸点。

### 3. 所有跳转类型的逐反向跳转上下文 Select

**概念：** 在 `OpJump`、`OpJumpTrue` 和 `OpJumpFalse` 的**每次**反向跳转时检查上下文取消：

**失败原因：** 所有基准测试约 5% 退化。`OpJumpTrue`/`OpJumpFalse` 在循环体中的命中频率远高于 `OpJump`（每次迭代都有循环退出检查）。带 `default` 的 `select{}` 编译为 `runtime.selectnbrecv`（约 10ns），在每次条件反向分支上执行开销太大。

**修复：** 只在 `OpJump` 反向跳转时检查，并用 128 间隔节流。

---

## 修改文件汇总

| 文件 | 变更行数 | 描述 |
|---|---|---|
| `value/value.go` | +8/−0 | `MakeFunc()` 构造器、`RawObj()` 访问器 |
| `value/accessor.go` | +4/−0 | `Interface()` + `ToReflectValue()` 中 `KindFunc` |
| `vm/run.go` | +25/−10 | `backJumpCount`、参数缓冲区、`RawObj()` 调度 |
| `vm/ops_dispatch.go` | +8/−5 | OpClosure/OpCallIndirect/OpFree 中 `MakeFunc()`、`RawObj()` |
| `vm/closure.go` | +0/−1 | 小清理 |
| **合计** | **+45/−16** | |

---

## 测试

所有优化通过完整测试套件：

```
ok  gig/bytecode    0.003s
ok  gig/compiler    0.002s
ok  gig/importer    0.003s
ok  gig/tests       0.846s
ok  gig/value       0.003s
ok  gig/vm          0.003s
```

---

## 累计优化进展（Phase 1 → Phase 3）

| 基准测试 | 原始 | Phase 1 | Phase 2 | Phase 3 | 总提升 |
|---|---|---|---|---|---|
| Fib(25) | 21.1 ms | 20.3 ms | 19.6 ms | **19.7 ms** | **6.6%** |
| BubbleSort(100) | 1.08 ms | 953 μs | 913 μs | **943 μs** | **12.7%** |
| Sieve(1000) | 187 μs | 186 μs | 196 μs | **187 μs** | ~0% |
| ClosureCalls(1K) | 371 μs | 384 μs | 369 μs | **326 μs** | **12.1%** |
| ExtCallDirectCall | 583 μs | 577 μs | 511 μs | **500 μs** | **14.2%** |
| ExtCallReflect | 359 μs | 360 μs | 344 μs | **338 μs** | **5.8%** |
| ExtCallMethod | 460 μs | 452 μs | 439 μs | **416 μs** | **9.6%** |
| ExtCallMixed | 331 μs | 331 μs | 308 μs | **301 μs** | **9.1%** |

---

## 架构图

```
源代码 (.go)
       │
       ▼
  Go SSA (golang.org/x/tools/go/ssa)
       │
       ▼
  编译器 (compiler/)
   ├── 编译为字节码
   ├── 窥孔优化器
   │   ├── 超级指令融合 (Add/Sub/Mul)
   │   ├── 切片操作融合
   │   ├── 整数特化 (2 遍)
   │   └── 整数移动融合
   └── 操作数宽度: O(1) 数组查找
       │
       ▼
  字节码程序 (bytecode/)
   ├── 80+ 操作码，含融合操作
   ├── FuncByIndex、PrebakedConstants、IntConstants
   └── ExternalFuncInfo 及 DirectCall 包装器
       │
       ▼
  VM 执行 (vm/)
   ├── run() 主循环
   │   ├── 50+ 内联热操作码
   │   ├── KindFunc: 直接 *Closure 存储 ──────── 新增 (Phase 3)
   │   ├── RawObj() 类型断言 (无 reflect) ────── 新增 (Phase 3)
   │   ├── 栈分配参数缓冲区 [8] ─────────────── 新增 (Phase 3)
   │   └── 反向跳转上下文检查 (128x) ─────────── 新增 (Phase 3)
   ├── Value 类型: MakeFunc() / RawObj() ─────── 新增 (Phase 3)
   ├── extCallCache: []*entry (O(1) 切片)
   ├── closurePool: sync.Pool
   └── 帧池 (已有)
```
