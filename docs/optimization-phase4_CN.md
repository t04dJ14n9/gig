# 第四阶段优化：数据流分析与帧清零优化

## 概述

本文档介绍 Gig VM 第四轮优化，目标是通过**编译期数据流分析**消除冗余运行时开销。新增两个编译器分析 pass：

1. **IntOnlyLocals 分析**——识别仅被 `OpInt*` 指令访问的局部变量（为未来双写消除提供基础设施）
2. **部分帧清零**——使用 `ZeroFrom` 阈值 + `clear()` 跳过参数和直线代码中已写入的槽位清零

这些是**编译期**分析。帧清零优化减少每次调用的开销，对不受益的函数零额外开销。IntOnlyLocals 分析作为编译器基础设施保留，供未来操作码级双写消除使用。

### 设计考量

Phase 4 初始实现在 `vm/run.go` 的 7 个 `OpInt*` 处理器中添加了逐指令 `if intOnlyLocals == nil || !intOnlyLocals[idxC]` 分支，并使用 `NeedsZero []bool` 位图进行选择性帧清零。基准测试发现**三处性能回退**：

1. **热路径分支开销**：7 个 `OpInt*` 处理器中的逐指令 nil 检查 + 位图查找在 Fib(25) 的 242K 次递归调用中造成约 1–3% 的开销，即使优化从未生效（`fib` 函数的 IntOnlyLocals 为 nil）。
2. **memclr 内建函数失效**：将 `for i := range s { s[i] = zero }`（被 Go 编译器识别为 memclr）替换为逐元素条件 `if need { s[i] = zero }` 循环，破坏了批量清零优化。
3. **结构体大小增加**：向 `CompiledFunction` 添加 `IntOnlyLocals []bool`（24 字节切片头）在递归工作负载中造成可测量的缓存效应。

修订后的方案：
- **移除所有 VM 热路径变更**——`vm/run.go` 与 Phase 3 完全相同
- **用 `ZeroFrom int` 替代 `NeedsZero []bool`**——单一整数阈值保留连续 `clear()`（memclr 内建函数）
- **从 `CompiledFunction` 移除 `IntOnlyLocals`**——分析结果计算后丢弃，避免结构体膨胀
- **使用 `clear()` 内建函数**（Go 1.21+）——直接映射到 memclr，比依赖 range 循环模式匹配更可靠

### 结果（AMD EPYC 9754 128-Core Processor，linux/amd64，A/B 背靠背对比）

#### 核心工作负载

| 工作负载 | 原生 Go | Gig（第4阶段） | Yaegi | GopherLua | Gig vs 原生 | Gig vs Yaegi | Gig vs Lua |
|---|---:|---:|---:|---:|---:|---:|---:|
| Fibonacci(25) | 448 μs | 19.7 ms | 107 ms | 21.1 ms | 慢 44x | **快 5.4x** | **快 1.07x** |
| ArithSum(1K) | 659 ns | 34.3 μs | 39.8 μs | 38.8 μs | 慢 52x | **快 1.2x** | **快 1.1x** |
| BubbleSort(100) | 6.4 μs | 935 μs | 1,206 μs | 768 μs | 慢 146x | **快 1.3x** | 慢 1.2x |
| Sieve(1000) | 1.86 μs | 188 μs | 206 μs | 206 μs | 慢 101x | **快 1.1x** | **快 1.1x** |
| ClosureCalls(1K) | 345 ns | 315 μs | 956 μs | 122 μs | 慢 913x | **快 3.0x** | 慢 2.6x |

#### 外部函数调用（Gig vs Yaegi，无 Lua/原生对等测试）

| 工作负载 | 原生 Go | Gig | Yaegi | Gig vs 原生 | Gig vs Yaegi |
|---|---:|---:|---:|---:|---:|
| DirectCall | 28.3 μs | 495 μs | 1,520 μs | 慢 17x | **快 3.1x** |
| Reflect | 24.1 μs | 328 μs | 998 μs | 慢 14x | **快 3.0x** |
| Method | 18.2 μs | 418 μs | 1,227 μs | 慢 23x | **快 2.9x** |
| Mixed | 11.4 μs | 294 μs | 874 μs | 慢 26x | **快 3.0x** |

#### 内存效率（allocs/op）

| 工作负载 | Gig | Yaegi | GopherLua | Gig vs Yaegi |
|---|---:|---:|---:|---:|
| Fibonacci(25) | 6 | 2,138,703 | 41 | **少 356,450x** |
| ArithSum(1K) | 6 | 8 | 93 | 少 1.3x |
| BubbleSort(100) | 9 | 5,085 | 12 | **少 565x** |
| Sieve(1000) | 7 | 43 | 207 | **少 6x** |
| ClosureCalls(1K) | 1,995 | 13,018 | 96 | **少 6.5x** |

---

## 编译器分析 1：IntOnlyLocals（基础设施）

### 问题

在 int 特化 VM 路径中，每条 `OpInt*SetLocal` 指令执行**双写**：

```go
case bytecode.OpIntLocalConstAddSetLocal:
    r := intLocals[idxA] + intConsts[idxB]
    intLocals[idxC] = r              // 8 字节 int64 写入（快）
    locals[idxC] = value.MakeInt(r)  // 32 字节 Value 构造 + 存储（慢）
```

`locals[idxC]` 写入的目的是让通用（非 int 特化）代码能从 `locals[]` 读取正确值。但如果一个局部变量**仅被** `OpInt*` 指令访问，`locals[]` 副本永远不会被读取——该写入纯属浪费。

### 解决方案

新增编译期两阶段分析：

**阶段 1** —— `intSpecialize()` 现在返回 `intUsed` 位图（哪些局部变量参与 `OpInt*` 运算）。

**阶段 2** —— `buildIntOnlyLocals()` 扫描特化后的字节码，找出仅被 `OpInt*` 指令访问的局部变量：

分析以所有 `intUsed` 局部变量为候选，然后**撤销**被通用指令访问的变量的 int-only 状态：

| 撤销指令 | 原因 |
|---|---|
| `OpLocal(idx)` | 通用读取 locals[] |
| `OpSetLocal(idx)` | 通用写入 locals[] |
| `OpAddSetLocal(idx)` | 融合通用操作写入 locals[] |
| `OpIntSetLocal(idx)` | 桥接指令——局部变量跨 int/通用边界 |
| `OpIntMoveLocal(src,dst)` | 保守处理：双方撤销 |

### 当前状态：仅作为基础设施

`buildIntOnlyLocals()` 的结果**计算后即丢弃**（`_ = buildIntOnlyLocals(...)`），因为在 VM 热路径中通过逐指令分支应用会导致性能回退。分析代码保留作为未来方案的基础设施：

- **专用无双写操作码**（如 `OpIntLocalConstAddSetLocal_NoDual`）——无条件跳过 `locals[]` 写入，无需逐指令分支
- 编译器为 int-only 局部变量发出这些操作码，为其他变量发出标准双写操作码

### 修改文件

| 文件 | 变更 |
|---|---|
| `compiler/optimize.go` | `intSpecialize` 返回 `intUsed`；新增 `buildIntOnlyLocals()` |
| `compiler/compile_func.go` | 调用 `buildIntOnlyLocals`（结果用 `_` 丢弃） |

---

## 优化 2：基于 `ZeroFrom` 的部分帧清零

### 问题

函数调用时，VM 分配帧并为 `NumLocals` 个值槽位分配空间。对复用帧（来自帧池），**所有**槽位都需清零：

```go
// 修改前：每次函数调用清零所有 locals
f.locals = f.locals[:fn.NumLocals]
for i := range f.locals {
    f.locals[i] = value.Value{}  // 32 字节 × NumLocals
}
```

对 `fib(25)` 的 ~242K 次递归调用，每次清零 8 个 locals × 32 字节 = 256 字节，总清零工作量约 **62 MB** 内存填充——仅仅为了确保即将被覆盖的槽位的正确性。

在 SSA 形式中，大多数局部变量遵循严格的"先定义后使用"模式：参数由调用者写入，SSA 临时变量在使用前定义。只有 Phi 节点目标和可能未初始化的变量需要清零。

### 解决方案

新增 `computeZeroFrom()` 编译期分析，找到需要清零的最低局部变量索引：

```go
func computeZeroFrom(code []byte, numLocals, numParams int) int
```

分析执行字节码的正向扫描：

1. **参数**（索引 0..NumParams-1）：永远不需要清零——调用者总是写入它们
2. **直线代码分析**：对每个局部变量，如果首次访问是写入（SetLocal），则不需要清零
3. **在任何分支处**（Jump/JumpTrue/JumpFalse）：保守地标记所有未解析的局部变量为需要清零

结果是存储在 `CompiledFunction` 中的单个 `ZeroFrom int`。`needsZero` 位图的模式（参数和直线写入的局部变量形成连续的 `false` 前缀）自然映射为单一截断索引。这保留了 `clear()`（memclr）内建函数：

```go
// 修改后：仅从 ZeroFrom 开始清零
if zf := fn.ZeroFrom; zf > 0 {
    clear(f.locals[zf:])   // memclr 内建函数，只清零需要的部分
} else {
    clear(f.locals)        // ZeroFrom=0 → 全部清零（向后兼容默认值）
}
```

关键设计选择：
- **`ZeroFrom int` 而非 `NeedsZero []bool`**：位图需要逐元素条件清零，破坏 memclr 内建函数。单一阈值保留连续 `clear()`。
- **`clear()` 内建函数（Go 1.21+）**：直接调用 memclr，比依赖 `for i := range s { s[i] = zero }` 的模式匹配更可靠。
- **零值即全部清零**：`ZeroFrom = 0` 表示所有 locals 都需要清零，是安全的向后兼容默认值。

### 修改文件

| 文件 | 变更 |
|---|---|
| `bytecode/bytecode.go` | `CompiledFunction` 新增 `ZeroFrom int` |
| `compiler/optimize.go` | 新增 `computeZeroFrom()` 函数 |
| `compiler/compile_func.go` | 管线接入 `computeZeroFrom` |
| `vm/frame.go` | `framePool.get()` 使用 `clear()` + `ZeroFrom` 对 `locals` 和 `intLocals` 部分清零 |

### 影响

对参数占局部变量比例较大的函数（如 `fib` 有 3 个参数、8 个 locals），每次调用节省约 37% 的清零工作。对有大量 SSA 临时变量的直线代码，节省可达 80% 以上。`clear()` 内建函数确保通过 memclr 内建函数获得最大吞吐量。

---

## 状态更新

### 已在第 4 阶段前实现的优化

调研本阶段时，确认以下计划中的优化**已在第 3 阶段实现**：

1. **反向跳转上下文检查** —— `OpJump` 处理器中的 `backJumpCount` 节流（第 3 阶段实现）
2. **闭包池化** —— `vm/closure.go` 中的 `sync.Pool`（第 3 阶段实现）
3. **DirectCall 扩展** —— `canWrapUnderlying` 已支持 struct、pointer、非空 interface、map、跨包命名类型。方法 DirectCall 包装器已生成：
   - **619 / 671 个函数**有 DirectCall 包装器（92.3% 覆盖率）
   - **543 个方法 DirectCall**跨 24 个标准库包
   - 剩余 52 个 `nil` 均为**函数参数类型**（回调如 `strings.Map`、`sort.Slice`）—— 推迟到未来工作

### 累计进展

| 基准测试 | 基线 | 第 1 阶段 | 第 2 阶段 | 第 3 阶段 | 第 4 阶段 |
|---|---:|---:|---:|---:|---:|
| Fib(25) | 24.1 ms | 22.4 ms | 20.3 ms | 19.7 ms | 19.7 ms |
| ArithSum(1K) | 99.1 μs | 39.6 μs | 35.2 μs | 34.1 μs | 34.3 μs |
| BubbleSort(100) | 2,124 μs | 1,027 μs | 913 μs | 943 μs | 935 μs |
| Sieve(1000) | — | — | — | 187 μs | 188 μs |
| ClosureCalls(1K) | — | — | 369 μs | 326 μs | 315 μs |

### 架构（标 * 为第 4 阶段新增）

```
源代码
    │
    ▼
SSA 编译器 (golang.org/x/tools/go/ssa)
    │
    ▼
字节码编译器 (compiler/compile_func.go)
    │
    ├─── Pass 1: optimizeBytecode()      [融合通用超级指令]
    ├─── Pass 2: fuseSliceOps()          [融合 []int 访问模式]
    ├─── Pass 3: intSpecialize()         [升级为 OpInt* 变体]
    ├─── Pass 3.5*: buildIntOnlyLocals() [识别可跳过双写的局部变量——基础设施]
    ├─── Pass 4: fuseIntMoves()          [融合 OpIntLocal+OpIntSetLocal]
    └─── Pass 5*: computeZeroFrom()      [计算部分清零阈值]
    │
    ▼
优化后字节码 + 元数据
    (ZeroFrom 阈值*)
    │
    ▼
VM 执行 (vm/run.go——与 Phase 3 完全相同)
    ├── framePool.get()：通过 ZeroFrom + clear() 部分清零*
    ├── 反向跳转上下文检查（第 3 阶段）
    └── 闭包池化 via sync.Pool（第 3 阶段）
```
