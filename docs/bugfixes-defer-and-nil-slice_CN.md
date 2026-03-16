# Bug 修复报告 — Defer 执行与 Nil Slice Append

本文档描述了本次发现并修复的七个已确认 bug。
所有 bug 涉及 **defer 语义**（命名返回值修改、闭包捕获、栈序执行）和
**nil slice append** 处理。

**修改文件：**

| 文件 | 说明 |
|------|------|
| `bytecode/opcode.go` | 新增 `OpRunDefers` 和 `OpDeferIndirect` 操作码 |
| `compiler/compile_instr.go` | 修复 `*ssa.RunDefers` 编译（原来错误地生成 `OpRecover`） |
| `compiler/compile_value.go` | 新增 `compileDefer` 支持直接调用和闭包 defer |
| `value/container.go` | 清理 `Elem()` / `SetElem()` 中无用的 `**Value` 处理 |
| `vm/frame.go` | `DeferInfo` 新增 `closure` 字段支持间接 defer |
| `vm/ops_dispatch.go` | 新增 `OpRunDefers` / `OpDeferIndirect` / `OpFree` / `OpClosure` / `OpAppend` 修复 |
| `vm/vm.go` | 新增 `runDefers` 辅助方法 |

---

## Bug 1：`*ssa.RunDefers` 被错误编译为 `OpRecover`（根本原因）

**现象：** 所有 defer 相关测试用例均失败 — `DeferNamedReturn`、`DeferModifyNamed`、
`DeferStackOrder` 和 `ClosureWithDefer` 均返回错误结果或 panic。

**根因：** 在 `compiler/compile_instr.go` 中，SSA 指令 `*ssa.RunDefers` 被错误地
编译为 `OpRecover`：

```go
// 修复前（错误）
case *ssa.RunDefers:
    c.emit(bytecode.OpRecover)
```

在 Go 的 SSA 中间表示中，`RunDefers` 是一条关键指令，它出现在函数**读取命名返回值并返回之前**。
SSA 序列如下：

```
*result = 5          // 设置命名返回值
rundefers            // 执行所有 defer 函数（可能修改 *result）
t = *result          // 读取（可能已修改的）命名返回值
return t
```

由于生成了 `OpRecover` 而非正确的操作码，defer 函数在正确的时间点从未被执行，
命名返回值也从未被 defer 修改。

**修复：** 生成正确的 `OpRunDefers` 操作码：

```go
// 修复后（正确）
case *ssa.RunDefers:
    c.emit(bytecode.OpRunDefers)
```

**测试用例：** `DeferNamedReturn` — 预期返回 10（5 × 2），实际返回 5。

---

## Bug 2：`OpRunDefers` 操作码不存在

**现象：** 即使修复了编译器，VM 中也没有同步执行 defer 的处理器。

**根因：** 字节码指令集中没有"立即执行所有待处理 defer"的操作码。
Defer 仅在函数返回时（`OpReturn`/`OpReturnVal`）执行，而此时命名返回值**已经被读取** —
对 defer 来说为时已晚。

**修复：** 在 `bytecode/opcode.go` 中新增 `OpRunDefers`，并在 `vm/ops_dispatch.go`
中实现处理器：

1. 按 **LIFO** 顺序（后进先出）遍历待处理 defer
2. 为每个 defer 创建**子 VM**（避免干扰父帧栈）
3. 与父 VM 共享 globals、context、program 和 extCallCache
4. **同步**执行每个 defer，然后继续

```go
case bytecode.OpRunDefers:
    for len(frame.defers) > 0 {
        d := frame.defers[len(frame.defers)-1]
        frame.defers = frame.defers[:len(frame.defers)-1]
        var freeVars []*value.Value
        if d.closure != nil {
            freeVars = d.closure.FreeVars
        }
        childVM := &VM{
            program: vm.program, stack: make([]value.Value, 256),
            globals: vm.globals, globalsPtr: vm.globalsPtr,
            ctx: vm.ctx, extCallCache: vm.extCallCache,
        }
        deferFrame := newFrame(d.fn, 0, d.args, freeVars)
        childVM.frames[0] = deferFrame
        childVM.fp = 1
        _, _ = childVM.run()
    }
```

相应地，从 `OpReturn` 和 `OpReturnVal` 中**移除**了 defer 执行逻辑，
因为 defer 现在通过 `OpRunDefers` 在 SSA 定义的正确时间点运行。

**测试用例：** `DeferStackOrder` — 预期返回 1111（1000 + 100 + 10 + 1），实际 panic。

---

## Bug 3：`OpFree` 双重包装指针值

**现象：** `ClosureWithDefer` 和 `DeferModifyNamed` 失败，
因为闭包无法通过自由变量捕获正确读写共享变量。

**根因：** `OpFree` 使用 `value.FromInterface(frame.freeVars[idx])` 将 `*value.Value`
指针包装进新的 `Value`，产生双重间接引用（`Value` → `*Value` → `Value`）。
期望读取简单 `int` 值（通过 reflect 指针）的闭包会得到一个 `*value.Value` 包装器 — 类型不匹配。

**修复：** 将 `OpFree` 改为直接解引用 slot：

```go
// 修复前
vm.push(value.FromInterface(frame.freeVars[idx]))

// 修复后
vm.push(*frame.freeVars[idx])
```

这确保闭包看到的是**实际捕获的值**（例如类型为 `*int` 的 `reflect.Value`），
而不是 slot 指针的包装器。

**测试用例：** `DeferModifyNamed` — 预期返回 999，实际返回 42。

---

## Bug 4：`OpClosure` 自由变量 slot 创建错误

**现象：** 多个共享同一捕获变量的闭包无法看到彼此的修改。

**根因：** `OpClosure` 有复杂的逻辑尝试从栈上检测 `*value.Value` 和 `**value.Value`，
但行为不一致。slot 共享机制（多个闭包引用同一 `*value.Value` slot 来共享状态）已损坏。

**修复：** 简化 `OpClosure`，始终为每个捕获变量创建新的 `*value.Value` slot：

```go
slot := new(value.Value)
*slot = v
closure.FreeVars[i] = slot
```

如果捕获的值是 reflect 指针（例如 `Alloc` 产生的 `*int`），
所有共享该指针的闭包都能看到彼此的修改 — 因为它们共享同一个底层堆分配的 int，
即使每个闭包有自己的 `*value.Value` 包装 slot。

**测试用例：** `ClosureWithDefer` — 预期返回 30，实际返回 nil/panic。

---

## Bug 5：`OpDeferIndirect` — 基于闭包的 defer

**现象：** `defer func() { result *= 2 }()` — 使用匿名闭包的 defer 未被正确编译或执行。

**根因：** 编译器不支持 `OpDeferIndirect`（defer 一个闭包调用），
它与 `OpDefer`（defer 一个具名函数调用）不同。
SSA 的 `Defer` 指令可以指向 `*ssa.Function` 或 `*ssa.MakeClosure`，
后者需要在 defer 时捕获自由变量。

**修复：** 在字节码、编译器和 VM 中新增 `OpDeferIndirect`：

- **字节码：** 新操作码，操作数 `[num_args:2]`
- **编译器：** `compileDefer` 现在处理 `*ssa.Function`、`*ssa.MakeClosure` 和回退情况
- **VM：** `OpDeferIndirect` 从栈上弹出参数和闭包，
  存储包含闭包 `FreeVars` 的 `DeferInfo` 以供后续执行

**测试用例：** `DeferNamedReturn` — `defer func() { result *= 2 }()` 是闭包 defer。

---

## Bug 6：`Slice_AppendToNil` — nil slice + 原生 `[]int64` 元素

**现象：** `var s []int; s = append(s, 1)` panic：
`reflect.Set: value of type int is not assignable to type []int64`

**根因：** SSA 将 `append(s, 1)` 编译为：
1. 创建 `[1]int{1}` 数组
2. 切片为 `[]int{1}`（内部存储为 `[]int64{1}`）
3. 调用 `append(s, sliced_result)`

在 `OpAppend` 处理器中，当 `s` 为 nil 时，nil slice 分支检查 `elem.ReflectValue()`。
但对于原生 `[]int64`，`ReflectValue()` 返回 `false`（原生切片不以 `reflect.Value` 存储）。
代码落入单元素 append 路径，将 `[]int64{1}` 当作单个元素，
尝试创建 `[][]int64` — 类型不匹配。

**修复：** 在 nil slice 分支开头新增原生 `[]int64` 的快速路径：

```go
if es, ok2 := elem.IntSlice(); ok2 {
    vm.push(value.MakeIntSlice(append([]int64(nil), es...)))
    break
}
```

**测试用例：** `Slice_AppendToNil` — 预期返回 3，实际 panic。

---

## Bug 7：`NilSliceAppend` — `append(nil, 1, 2, 3)`（原有 bug）

**现象：** `var s []int; s = append(s, 1, 2, 3); return len(s)` 返回 1 而非 3。

**根因：** 与 Bug 6 相同。SSA 将可变参数 `1, 2, 3` 打包为 `[]int{1, 2, 3}`
（存储为 `[]int64{1, 2, 3}`），然后调用 `append(nil, packed_slice)`。
nil slice 分支无法识别原生 `[]int64` 元素，仅追加了第一个值。

**修复：** 与 Bug 6 相同 — `IntSlice()` 快速路径正确地展开追加所有元素。

**测试用例：** `NilSliceAppend` — 预期返回 3，实际返回 1。

---

## 代码清理：移除无用的 `**Value` 处理

修复 `OpFree`（Bug 3）后，`value/container.go` 中 `Elem()` 和 `SetElem()` 的
`**Value` 检测逻辑以及 `vm/ops_dispatch.go` 中 `OpDeref` 的 `**value.Value` 解包逻辑
成为死代码。这些路径是为 `OpFree` 以前创建的双重包装而添加的变通方案。
由于 `OpFree` 现在直接解引用 slot，`**Value` 不再出现在栈上。

移除的死代码：
- `Elem()`：`**Value` 解引用检查
- `SetElem()`：`**Value` 写穿检查
- `OpDeref`：`**value.Value` 解包检查

---

## Lint 修复：`gci` 导入格式化

**现象：** `golangci-lint-v2` 报告 "File is not properly formatted (gci)"，
涉及 `bytecode/opcode.go` 和 `vm/ops_dispatch.go`。

**根因：** `gci` 格式化器（在 `.golangci.yml` 的 `formatters` 部分配置）
要求导入分组顺序为：`standard → default → prefix(gig) → localmodule`。
修改后的文件不符合此分组规则。

**修复：** 执行 `golangci-lint-v2 run --fix` 自动格式化两个文件。

---

## 测试结果

所有 7 个 bug 均有专用测试用例验证：

| 测试用例 | 预期值 | 修复前 | 修复后 |
|----------|--------|--------|--------|
| `DeferNamedReturn` | 10 | 5 | ✅ 10 |
| `DeferModifyNamed` | 999 | 42 | ✅ 999 |
| `DeferStackOrder` | 1111 | panic | ✅ 1111 |
| `ClosureWithDefer` | 30 | panic | ✅ 30 |
| `MultipleNamedReturn` | 1021 | （未测试） | ✅ 1021 |
| `Slice_AppendToNil` | 3 | panic | ✅ 3 |
| `NilSliceAppend` | 3 | 1 | ✅ 3 |

完整测试套件：`go test ./...` — 所有包 **0 失败**。

Lint 检查：`golangci-lint-v2 run` — **0 问题**。
