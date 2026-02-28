# 融合 Int 切片超级指令 — 优化报告

## 摘要

此优化添加了 **3 个融合超级指令**用于 `[]int` 元素访问，将 7 条指令（17 字节）的序列替换为单次分发的 7 字节操作码。这是优化系列的集大成之作，将 BubbleSort 和 Sieve 从**落后**翻转为**领先** Yaegi。

### 结果（AMD EPYC 9754，linux/amd64）

| 基准测试 | 切片融合前 | 切片融合后 | 变化 | vs Yaegi |
|---|---|---|---|---|
| Fib25 | 20.7 ms | 20.2 ms | — | **Gig 快 5.6 倍** |
| ArithSum | 37.2 μs | 37.5 μs | — | **Gig 快 12%** |
| BubbleSort | 1,770 μs（慢 1.40 倍） | **963 μs** | **快 1.84 倍** | **Gig 快 32%** |
| Sieve | 301 μs（慢 1.44 倍） | **203 μs** | **快 1.48 倍** | **Gig 快 3%** |
| ClosureCalls | 392 μs | 381 μs | — | **Gig 快 2.7 倍** |

**Gig 现在在所有 5 个基准测试中都胜过 Yaegi。**

### 累积改进（完整优化系列）

| 基准测试 | 原始 | 当前 | 总加速比 |
|---|---|---|---|
| Fib25 | 169 ms | 20.2 ms | **8.4x** |
| ArithSum | 311 μs | 37.5 μs | **8.3x** |
| BubbleSort | 10.3 ms | 963 μs | **10.7x** |
| Sieve | 1,681 μs | 203 μs | **8.3x** |
| ClosureCalls | 964 μs | 381 μs | **2.5x** |

---

## 问题

整数切片访问（`s[i]` 读取和 `s[i] = v` 写入）编译为 7 条指令的序列：

**读取模式（`v = s[j]`）：**
```
LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) DEREF SETLOCAL(v)
```
= 7 次分发，17 字节，多次栈压入/弹出，kind 检查和指针间接寻址。

**写入模式（`s[j] = val`）：**
```
LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) LOCAL(val) SETDEREF
```
= 7 次分发，17 字节，相同开销。

这些模式在 BubbleSort（每次比较需要 2 次读取 + 2 次写入的交换）和 Sieve（每次内层循环迭代 1 次读取 + 1 次写入）中占主导。每个模式需要 7 次指令分发、7 次字节码解码和多次栈操作——而逻辑上只是一次数组访问。

---

## 解决方案：3 个融合操作码

### OpIntSliceGet(s, j, v) — 7 字节，1 次分发

融合读取模式。运行时：
```go
slice := locals[s].([]int64)    // 一次类型断言
result := slice[intLocals[j]]   // 直接索引读取
intLocals[v] = result            // 8 字节写入
locals[v] = value.MakeInt(result) // 32 字节同步
```

### OpIntSliceSet(s, j, val) — 7 字节，1 次分发

融合从局部变量写入的模式。运行时：
```go
slice := locals[s].([]int64)
slice[intLocals[j]] = intLocals[val]  // 直接索引写入
```

### OpIntSliceSetConst(s, j, c) — 7 字节，1 次分发

融合从常量写入的模式（`s[j] = 0` 或 `s[j] = 1`）。运行时：
```go
slice := locals[s].([]int64)
slice[intLocals[j]] = intConsts[c]    // 从常量池直接索引写入
```

此模式对 Sieve 基准测试至关重要，它执行 `sieve[i] = 1` 和 `sieve[j] = 0`。

### 安全性：回退路径

所有三个操作码都包含对非原生切片（例如由 `reflect.Value` 而非 `[]int64` 支持的 `[]int`）的回退路径。快速路径检查 `IntSlice()` 并回退到执行等效的通用操作：

```go
if s, ok := locals[sIdx].IntSlice(); ok {
    // 快速路径：直接 []int64 访问
} else {
    // 回退：通用执行 IndexAddr + Deref/SetDeref
}
```

---

## 编译器变更

### 新遍：`fuseSliceOps()`

在 `optimizeBytecode()` 和 `intSpecialize()` 之间插入的新窥孔遍：

```
原始字节码
     │
     ▼
第 1 遍：optimizeBytecode()     — 算术/比较融合
     │
     ▼
第 2 遍：fuseSliceOps()         — ★ 新增：切片访问融合
     │
     ▼
第 3 遍：intSpecialize()        — Value → int64 升级
     │
     ▼
第 4 遍：fuseIntMoves()         — phi 移动融合
     │
     ▼
优化后字节码
```

### 模式匹配

该遍扫描三个 17 字节模式，都以 `LOCAL LOCAL INDEXADDR SETLOCAL LOCAL` 开头：

| 模式 | 后缀 | 融合操作码 |
|---|---|---|
| 1（读取） | `DEREF SETLOCAL(v)` | `OpIntSliceGet(s,j,v)` |
| 2（从局部变量写入） | `LOCAL(val) SETDEREF` | `OpIntSliceSet(s,j,val)` |
| 3（从常量写入） | `CONST(val) SETDEREF` | `OpIntSliceSetConst(s,j,c)` |

### 类型要求

- `s` 必须在 `localIsIntSlice` 中（SSA 类型为 `[]int` 或 `[]int64`）
- `j` 必须在 `localIsInt` 中（SSA 类型为 `int` 或 `int64`）
- `v`/`val` 必须在 `localIsInt` 中（模式 1 和 2）
- `ptr` 和 `ptrGet` 必须引用同一局部变量（确认无别名）

### `isIntSliceType()` 检测

仅匹配 `[]int` 和 `[]int64`（对应 VM 中的原生 `[]int64` 快速路径）：

```go
func isIntSliceType(t types.Type) bool {
    sl, ok := t.Underlying().(*types.Slice)
    if !ok { return false }
    basic, ok := sl.Elem().Underlying().(*types.Basic)
    switch basic.Kind() {
    case types.Int, types.Int64:
        return true
    }
    return false
}
```

### 与 intSpecialize 的集成

融合操作码的索引/值操作数在 `intSpecialize` 第 1 遍中注册，以保持其 `intLocals[]` 槽位同步：

```go
case bytecode.OpIntSliceGet:
    intUsed[j] = true; intUsed[v] = true
case bytecode.OpIntSliceSet:
    intUsed[j] = true; intUsed[val] = true
case bytecode.OpIntSliceSetConst:
    intUsed[j] = true
```

---

## 具体示例：BubbleSort 内层循环

**源码：**
```go
if arr[j] > arr[j+1] {
    arr[j], arr[j+1] = arr[j+1], arr[j]
}
```

**优化前（每次交换 28+ 条指令）：**
```
LOCAL(arr) LOCAL(j) INDEXADDR SETLOCAL(ptr1)     — 4 次分发
LOCAL(ptr1) DEREF SETLOCAL(aj)                   — 3 次分发
LOCAL(arr) LOCAL(j1) INDEXADDR SETLOCAL(ptr2)    — 4 次分发
LOCAL(ptr2) DEREF SETLOCAL(aj1)                  — 3 次分发
... 比较、交换写入 ...                             — 14+ 次分发
```

**优化后（4 个融合 + 比较 + 2 个融合 = ~7 次分发完成交换路径）：**
```
OpIntSliceGet(arr, j, aj)         — 1 次分发
OpIntSliceGet(arr, j1, aj1)       — 1 次分发
OpIntGreaterLocalLocalJumpTrue ... — 1 次分发
OpIntSliceSet(arr, j, aj1)        — 1 次分发
OpIntSliceSet(arr, j1, aj)        — 1 次分发
```

**节省：** ~21 次更少的分发，~210 字节更少的字节码解码，零中间栈流量。

---

## 修改的文件

| 文件 | 变更 |
|---|---|
| `bytecode/opcode.go` | 添加 `OpIntSliceGet`、`OpIntSliceSet`、`OpIntSliceSetConst`（6 字节操作数） |
| `compiler/compile_func.go` | 添加 `isIntSliceType()`、`localIsIntSlice[]` 映射，将 `fuseSliceOps()` 接入流水线 |
| `compiler/optimize.go` | 添加 `fuseSliceOps()` 含 3 个模式匹配器 + `intSpecialize` 第 1 遍注册 |
| `vm/run.go` | 添加 3 个热路径处理器（含 `IntSlice()` 快速路径 + 通用回退） |
| `value/value.go` | （清理：移除未使用的 `UnsafeIntSlice()`） |
