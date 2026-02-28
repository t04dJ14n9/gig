# 超级指令优化 — 性能报告

## 摘要

此优化为 Gig VM 引入了**窥孔超级指令**和**寄存器提升分发**，针对 Yaegi 此前保持 4-5 倍优势的紧密循环性能。

### 结果（跨解释器基准测试）

| 工作负载 | 优化前 (μs) | 优化后 (μs) | 加速比 | vs Yaegi 优化前 | vs Yaegi 优化后 |
|---|---|---|---|---|---|
| Fibonacci(25) | 46,550 | 32,045 | **1.45x** | Gig 快 2.4 倍 | Gig 快 3.5 倍 |
| ArithmeticSum(1K) | 200 | 136 | **1.47x** | Yaegi 快 5.0 倍 | Yaegi 快 3.3 倍 |
| BubbleSort(100) | 4,979 | 3,904 | **1.28x** | Yaegi 快 4.0 倍 | Yaegi 快 3.1 倍 |
| Sieve(1000) | 840 | 687 | **1.22x** | Yaegi 快 4.1 倍 | Yaegi 快 3.3 倍 |
| ClosureCalls(1K) | 586 | 528 | **1.11x** | Gig 快 1.7 倍 | Gig 快 1.9 倍 |

### 完整对比表（优化后）

| 工作负载 | 原生 Go | Gig | Yaegi | GopherLua | Gig/原生 | Yaegi/原生 | Gig vs Yaegi |
|---|---|---|---|---|---|---|---|
| Fibonacci(25) | 453 μs | 32.0 ms | 111 ms | 21.0 ms | 71x | 244x | **Gig 快 3.5 倍** |
| ArithmeticSum(1K) | 664 ns | 136 μs | 41 μs | 41 μs | 205x | 62x | Yaegi 快 3.3 倍 |
| BubbleSort(100) | 6.4 μs | 3.90 ms | 1.26 ms | 774 μs | 609x | 197x | Yaegi 快 3.1 倍 |
| Sieve(1000) | 1.88 μs | 687 μs | 209 μs | 212 μs | 366x | 111x | Yaegi 快 3.3 倍 |
| ClosureCalls(1K) | 347 ns | 528 μs | 1,009 μs | 122 μs | 1,522x | 2,908x | **Gig 快 1.9 倍** |

## 所做变更

### 1. 窥孔优化器（`compiler/optimize.go`）

一个后编译遍，扫描常见的字节码模式并将其融合为单个超级指令。在跳转目标修补之后运行。优化器：

- 扫描指令流中的已知模式
- 用更短的融合操作码替换多指令序列
- 重建字节码并修正所有跳转目标以适应大小变化

### 2. 新超级指令（`bytecode/opcode.go`）

17 个新的融合操作码：

**融合算术 + 存储（完全消除栈流量）：**
- `OpLocalLocalAddSetLocal` — `local[A] + local[B] → local[C]`（零栈操作）
- `OpLocalConstAddSetLocal` — `local[A] + const[B] → local[C]`
- `OpLocalConstSubSetLocal` — `local[A] - const[B] → local[C]`

**融合加载 + 算术（消除 2 次压入操作）：**
- `OpAddLocalLocal` — 压入 `local[A] + local[B]`
- `OpSubLocalLocal` — 压入 `local[A] - local[B]`
- `OpMulLocalLocal` — 压入 `local[A] * local[B]`
- `OpAddLocalConst` — 压入 `local[A] + const[B]`
- `OpSubLocalConst` — 压入 `local[A] - const[B]`

**融合比较+分支（消除 bool 压入/弹出 + 独立跳转）：**
- `OpLessLocalLocalJumpTrue` — 若 `local[A] < local[B]` 则跳转
- `OpLessLocalConstJumpTrue` — 若 `local[A] < const[B]` 则跳转
- `OpLessEqLocalConstJumpTrue` — 若 `local[A] <= const[B]` 则跳转
- `OpGreaterLocalLocalJumpTrue` — 若 `local[A] > local[B]` 则跳转
- `OpLessLocalLocalJumpFalse` — 若 `local[A] >= local[B]` 则跳转
- `OpLessLocalConstJumpFalse` — 若 `local[A] >= const[B]` 则跳转

**融合算术 + 存储（弹出 2 个，计算，存储）：**
- `OpAddSetLocal` — 弹出 a,b；存储 `a+b` 到 `local[A]`
- `OpSubSetLocal` — 弹出 a,b；存储 `a-b` 到 `local[A]`

### 3. 寄存器提升分发（`vm/run.go`）

主执行循环现在将关键 VM 字段提升到 Go 局部变量中：

```go
stack := vm.stack     // 栈切片头在寄存器中
sp := vm.sp           // 栈指针在寄存器中
prebaked := vm.program.PrebakedConstants
```

所有热路径操作码直接在这些局部变量上操作：
```go
// 优化前：vm.push(frame.locals[idx]) → 方法调用 + 间接访问
// 优化后：stack[sp] = locals[idx]; sp++ → 直接基于寄存器的访问
```

这消除了每条指令的开销：
- push/pop 的方法调用开销（Go 无法内联所有方法调用）
- 每次指令从内存重复加载 `vm.stack`、`vm.sp`
- 每次 push 的栈增长检查（现在仅在函数边界检查）

局部变量在任何可能读取它们的操作（executeOp、callFunction、context 检查）之前同步回 `vm.*` 字段。

### 4. 集成（`compiler/compile_func.go`）

优化器在跳转修补之后自动调用：
```go
c.patchJumps(blockOffsets)
c.currentFunc.Instructions = optimizeBytecode(c.currentFunc.Instructions)
```

## 改进原因

### 优化前：ArithSum 内层循环（`sum += i; i++; i < 1000`）
```
每次迭代 12 条字节码：
  OpLocal(sum) OpLocal(i) OpAdd OpSetLocal(sum)    — 4 次分发，6 次栈操作
  OpLocal(i) OpConst(1) OpAdd OpSetLocal(i)        — 4 次分发，6 次栈操作
  OpLocal(i) OpConst(1000) OpLess OpJumpTrue       — 4 次分发，6 次栈操作
总计：12 次分发，每次迭代 18 次 48 字节 Value 的压入/弹出
```

### 优化后：ArithSum 内层循环
```
每次迭代 3 条超级指令：
  OpLocalLocalAddSetLocal(sum, i, sum)              — 1 次分发，0 次栈操作
  OpLocalConstAddSetLocal(i, 1, i)                  — 1 次分发，0 次栈操作
  OpLessLocalConstJumpTrue(i, 1000, target)         — 1 次分发，0 次栈操作
总计：3 次分发，每次迭代 0 次压入/弹出（4 倍更少的分发，零栈流量）
```

## 架构对比更新

| 维度 | Gig（优化前） | Gig（优化后） | Yaegi |
|---|---|---|---|
| 每循环迭代分发次数 | 12 | 3 | 3-4 |
| 每循环迭代栈操作 | ~18 次压入/弹出 | 0 | 0 |
| 每迭代移动数据 | ~1,536 字节 | ~0 字节 | ~24 字节 |
| 分发方式 | 在 vm.* 上 switch | 在局部变量上 switch | 闭包链 |
| 剩余差距原因 | 48 字节 Value 结构体 | 48 字节 Value 结构体 | 24 字节 reflect.Value |

## 剩余差距分析

Gig 在纯算术循环上仍比 Yaegi 慢约 3.3 倍。剩余差距在于：

1. **Value 结构体大小（48 字节 vs 24 字节）** — 即使零栈操作，locals[] 仍移动 48 字节结构体。Yaegi 的 reflect.Value 为 24 字节。
2. **`obj any` 字段的 GC 压力** — `Value.obj` 字段（类型 `any`，16 字节）迫使 GC 扫描所有 Value 切片。Yaegi 使用 reflect.Value 有同样问题但内存减半。
3. **闭包线程 vs switch 分发** — Yaegi 的闭包线程执行有更好的分支预测。每个闭包"知道"其下一个闭包，而 switch 分发每次都是间接分支。

### 未来优化优先级

| 优先级 | 优化 | 预期影响 |
|---|---|---|
| P0 | 将 Value 缩减到 24 字节（NaN 装箱或分离联合体） | 算术循环 ~1.5-2 倍 |
| P1 | 直接线程分发（通过汇编的计算跳转） | 所有工作负载 ~1.2-1.5 倍 |
| P2 | 热函数的整数特化局部变量 | int 密集循环 ~1.3 倍 |
