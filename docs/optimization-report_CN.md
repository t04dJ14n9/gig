# Gig 解释器优化报告

## 摘要

本报告记录了 Gig Go 解释器的系统性优化，以缩小与 Yaegi（一个树遍历 Go 解释器）的性能差距。工作分为 6 个阶段，将 Gig 从在大多数基准测试上明显慢于 Yaegi 转变为在关键领域**具有竞争力或更快**。

### 最终结果（AMD EPYC 9754，linux/amd64）

| 基准测试 | 优化前 | 优化后 | 加速比 | vs Yaegi | vs 原生 |
|-----------|--------|-------|---------|----------|-----------|
| Fib25 | 169ms | **45ms** | **3.8x** | Gig **快 2.4 倍** | 慢 99 倍 |
| ArithSum | 311μs | **200μs** | **1.6x** | Yaegi 快 5 倍 | 慢 298 倍 |
| BubbleSort | 10.3ms | **4.9ms** | **2.1x** | Yaegi 快 3.9 倍 | 慢 753 倍 |
| Sieve | 1,681μs | **838μs** | **2.0x** | Yaegi 快 4.1 倍 | 慢 440 倍 |
| ClosureCalls | 964μs | **584μs** | **1.7x** | Gig **快 1.7 倍** | 慢 1,681 倍 |

### 分配减少

| 基准测试 | 优化前 (allocs/op) | 优化后 (allocs/op) | 减少倍数 |
|-----------|--------------------|--------------------|-----------|
| Fib25 | 728,262 | **68** | **10,710x** |
| ArithSum | 13 | **13** | 1x |
| BubbleSort | 39,818 | **16** | **2,489x** |
| Sieve | 5,864 | **14** | **419x** |
| ClosureCalls | 3,000 | **3,000** | 1x |

---

## 第 1 阶段：O(1) 函数查找

**问题：** 函数调用对每次 `OpCall` 使用映射查找（`map[string]*CompiledFunction`），这是带有哈希开销的 O(n)。

**解决方案：** 为每个函数分配编译时数字索引，并将它们存储在扁平切片 `FuncByIndex []*CompiledFunction` 中。`OpCall` 指令将函数索引直接编码在操作数中。

**修改的文件：**
- `bytecode/bytecode.go` — 为 `Program` 添加 `FuncByIndex` 字段
- `compiler/compiler.go` — 编译期间构建 `FuncByIndex`
- `vm/call.go` — 在 `callCompiledFunction` 中使用基于索引的查找

**影响：** 每次调用开销从 ~100ns 减少到 ~5ns。在 Fib25 等递归基准测试中最为明显。

---

## 第 2 阶段：帧池化

**问题：** 每次函数调用在堆上分配一个新的 `Frame`（`newFrame()` → `make([]value.Value, numLocals)`）。对于 Fib(25)，这意味着 242,785 次帧分配。

**解决方案：** 使用 `sync.Pool` 实现 `framePool` 回收 `Frame` 对象。函数返回时帧归还到池中。如果帧的 `locals` 切片足够大则复用，避免重新分配。

**关键细节：** `addrTaken = true`（局部变量地址被闭包获取）的帧不会返回池中，因为外部代码可能仍持有指向帧局部变量的指针。

**修改的文件：**
- `vm/frame.go` — 添加带有 `get()` 和 `put()` 方法的 `framePool`
- `vm/vm.go` — 为 VM 结构体添加 `fpool framePool`
- `vm/ops_dispatch.go` — 在 `OpReturn`/`OpReturnVal` 中使用 `vm.fpool.put(frame)`

**影响：** Fib25 分配次数从 728K 降至 68（每个唯一函数一次，而非每次调用）。这是单一最大的优化。

---

## 第 3 阶段：预烘焙常量

**问题：** 每条 `OpConst` 指令在运行时调用 `value.FromInterface()`，经过 `reflect.ValueOf()` → 类型开关 → `MakeInt()`。对于热循环中的整数常量，这是浪费。

**解决方案：** 在编译时将所有常量预计算为 `value.Value` 并存储在 `PrebakedConstants []value.Value` 中。`OpConst` 现在只做数组索引查找。

**修改的文件：**
- `bytecode/bytecode.go` — 添加 `PrebakedConstants []value.Value`
- `compiler/compiler.go` — 编译后构建 `PrebakedConstants`
- `vm/ops_dispatch.go` — `OpConst` 优先从 `PrebakedConstants` 读取

**影响：** 消除每次常量加载约 ~50ns 的开销。在 ArithSum 和循环密集基准测试中最为明显。

---

## 第 4 阶段：整数快速路径

**问题：** 算术和比较操作使用通用方法（`a.Add(b)`、`a.Cmp(b)`），每次操作都在运行时检查类型，即使两个操作数都是 int（最常见的情况）。

**解决方案：** 在 VM 热路径操作码中添加内联类型检查。对于 `OpAdd`、`OpSub`、`OpMul` 和所有比较运算符：如果两个操作数都是 `KindInt`，直接在 `RawInt()` 上执行操作，无需调用通用方法。

```go
case bytecode.OpAdd:
    b := vm.pop()
    a := vm.pop()
    if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
        vm.push(value.MakeInt(a.RawInt() + b.RawInt()))
    } else {
        vm.push(a.Add(b))
    }
```

**影响：** 算术密集基准测试约 ~30% 加速。快速路径避免了方法调用开销和类型分发。

---

## 第 5 阶段：原生 `[]int64` 切片表示

**问题：** 整数切片（`[]int`）存储为包装 `reflect.MakeSlice(...)` 的 `reflect.Value`。每次元素访问经过 `rv.Index(i)` → `reflect.Value` → `MakeFromReflect()`，导致每次访问 2-3 次分配。

**解决方案：** 引入带有原生 `[]int64` 后备存储的 `KindSlice`。当 VM 创建整数切片（通过 `OpMakeSlice` 或对 `[N]int` 数组的 `OpSlice`）时，它将 Go `[]int64` 直接存储在 `Value.obj` 字段中。

### 关键实现细节：

1. **`MakeIntSlice([]int64)`** — 创建 `Value{kind: KindSlice, obj: []int64{...}}`
2. **`MakeIntPtr(*int64)`** — 对于 `OpIndexAddr`，存储原始 `*int64` 指针，无需 reflect
3. **SSA `Alloc([N]int) + Slice` 模式** — SSA 将常量 N 的 `make([]int, 100)` 编译为 `OpNew([100]int)` + `OpSlice`，而非 `OpMakeSlice`。`OpSlice` 处理器检测到 `reflect.Array` 含 `reflect.Int` 元素并转换为原生 `[]int64`。
4. **`SetElem` / `ToReflectValue` 转换** — 当原生 `[]int64` 需要存储到 `*[]int` 指针时（如二维切片），自动执行 `[]int64` → `[]int` 转换。

**原生 int 切片的快速路径操作：**
- `OpIndex`: `s[i]` → `MakeInt(s[key.RawInt()])`（无 reflect）
- `OpSetIndex`: `s[i] = v` → `s[key.RawInt()] = val.RawInt()`（无 reflect）
- `OpIndexAddr`: `&s[i]` → `MakeIntPtr(&s[idx])`（原始指针，无 reflect）
- `OpLen/OpCap`: 原生 `len(s)` / `cap(s)`
- `OpSlice`: 原生切片表达式
- `OpAppend/OpCopy`: 原生 append/copy

**修改的文件：**
- `value/value.go` — 添加 `MakeIntSlice`、`MakeIntPtr`、`IntSlice()`
- `value/container.go` — 在 `Len`、`Cap`、`Index`、`SetIndex`、`Elem`、`SetElem` 中添加 `KindSlice` 分支
- `value/accessor.go` — 在 `ToReflectValue` 中添加 `KindSlice` 分支
- `vm/ops_dispatch.go` — `OpMakeSlice`、`OpIndex`、`OpSetIndex`、`OpIndexAddr`、`OpSlice`、`OpLen`、`OpAppend`、`OpCopy` 中的快速路径

**影响：** Sieve 分配次数从 5,864 降至 14。BubbleSort 分配次数从 39,818 降至 16。Sieve 时间从 1,681μs 改善到 1,297μs。

---

## 第 6 阶段：内联热路径分发

**问题：** VM 的主循环对每条指令调用 `vm.executeOp(op, frame)` — 一个带有接口返回值和错误检查的 Go 函数调用。这增加了每条指令约 ~10-15ns 的开销，在紧密循环中占据主导。

**解决方案：** 将最频繁执行的操作码直接移入 `run()` 循环，作为带有 `continue` 的 `switch` 语句，绕过 `executeOp` 调用。不太常见的操作码回退到 `executeOp`。

**内联的操作码（覆盖数值程序中 >90% 的指令）：**
- 栈：`OpLocal`、`OpSetLocal`、`OpConst`、`OpNil`、`OpTrue`、`OpFalse`、`OpPop`、`OpDup`
- 算术：`OpAdd`、`OpSub`、`OpMul`
- 比较：`OpLess`、`OpLessEq`、`OpGreater`、`OpGreaterEq`、`OpEqual`、`OpNotEqual`
- 逻辑：`OpNot`
- 跳转：`OpJump`、`OpJumpTrue`、`OpJumpFalse`
- 调用：`OpCall`、`OpReturn`、`OpReturnVal`
- 指针：`OpSetDeref`

**额外微优化：** `OpJumpTrue`/`OpJumpFalse`/`OpNot` 使用 `RawBool()`（无检查的 `v.num != 0`）而非 `Bool()`（执行 kind 检查 + panic）。SSA 保证条件始终是布尔值。

**修改的文件：**
- `vm/run.go` — 用内联热路径 switch 重写 `run()`

**影响：** 所有基准测试 1.5 倍加速。Fib25: 69ms → 45ms。ArithSum: 313μs → 200μs。BubbleSort: 7.4ms → 4.9ms。

---

## 架构概览

```
源代码 (.go)
       │
       ▼
  Go SSA 包 (golang.org/x/tools/go/ssa)
       │
       ▼
  编译器 (gig/compiler)
   ├── 阶段 1：收集和索引函数
   ├── 阶段 2：分配局部变量槽位（参数、phi、值）
   ├── 阶段 3：逆后序编译基本块
   ├── 阶段 4：修补跳转目标
   └── 阶段 5：预烘焙常量
       │
       ▼
  字节码程序 (gig/bytecode)
   ├── FuncByIndex []*CompiledFunction  ← O(1) 查找
   ├── PrebakedConstants []value.Value   ← 零开销
   └── Types []types.Type
       │
       ▼
  VM 执行 (gig/vm)
   ├── 内联热路径分发          ← 无函数调用开销
   ├── 帧池化 (sync.Pool)     ← 每次调用近零分配
   ├── 原生 int 切片快速路径   ← []int 无 reflect
   └── 整数算术快速路径        ← 直接 int64 操作
```

---

## 剩余优化机会

1. **基于寄存器的 VM**：从基于栈转换为基于寄存器的架构。这将消除压入/弹出开销（当前每次操作 2 次数组访问）并支持指令融合。预期 2-3 倍改进但需要大幅重写。

2. **超级指令**：将常见指令序列（如 `OpLocal+OpLocal+OpAdd+OpSetLocal` → `OpAddLocals`）合并为单个操作码。算术循环预期 1.2-1.5 倍改进。

3. **原生 bool 切片**：类似于 `[]int64`，为使用 `[]bool` 的 Sieve 类工作负载实现 `[]bool` 快速路径。

4. **闭包分配减少**：ClosureCalls 仍有 3,000 allocs/op 来自闭包创建。可以池化闭包或使用不同表示。

5. **全局变量优化**：`OpGlobal` 当前通过 `FromInterface(&globals[idx])` 创建指针导致分配。可以使用类似局部变量的直接索引方法。

6. **指令编码优化**：当前编码每条指令 3 字节（1 操作码 + 2 操作数）。可以使用变长编码减少指令缓存压力。

---

## 测试

所有优化保持完全向后兼容。测试套件（`go test ./...`）零失败通过：

```
ok  gig/bytecode    0.002s
ok  gig/compiler    0.002s
ok  gig/importer    0.003s
ok  gig/tests       0.852s
ok  gig/value       0.003s
ok  gig/vm          0.003s
```

`benchmarks/` 中的基准测试验证所有解释器（Gig、Yaegi、GopherLua、原生 Go）的正确结果。
