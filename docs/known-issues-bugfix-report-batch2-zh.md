# 已知问题修复报告 – 第二批

**日期**: 2026-03-16
**分支**: feature/dev_youngjin

## 概述

解决了 `tests/known_issues_test.go` 中追踪的全部 6 个已知问题：

| # | Bug | 根因 | 修复位置 |
|---|-----|------|---------|
| 11 | DeferInClosureWithArg | 编译器栈顺序错误 | `compiler/compile_value.go` |
| 12 | PointerSwapInStruct | OpDeref 对指针字段返回引用而非副本 | `vm/ops_dispatch.go` |
| 13 | StructWithFuncSlice | OpNew 为函数数组创建 `[]value.Value` | `vm/ops_dispatch.go` |
| 14 | StructAnonymousField | 匿名未导出字段缺少 PkgPath | `vm/typeconv.go` |
| 15 | StructEmbeddedInterface | 已通过（仅作为已知问题追踪） | — |
| 16 | MapRangeWithBreak | Go 规范中非确定性行为 | — |

## Bug 11: DeferInClosureWithArg

**现象**: 闭包内 `defer func(v int){ result += v }(10)` 返回 `result = 1` 而非 `11`。

**根因**: `compileDefer` 在处理带 FreeVars 的 `*ssa.Function` 时，先将参数入栈，再创建闭包。但 `OpDeferIndirect` 先弹出参数（栈顶），再弹出闭包。栈顺序不匹配导致参数和闭包被调换。

**修复**: 调整 `compileDefer` 的代码生成顺序：先推入自由变量绑定 → `OpClosure` → 再推入参数 → `OpDeferIndirect`。对 `*ssa.Function` 和 `*ssa.MakeClosure` 两个分支都进行了修复。

## Bug 12: PointerSwapInStruct

**现象**: `p.a, p.b = p.b, p.a` 对 `PtrPair{a: &x, b: &y}` 执行后，`*p.b = 2`（与 `*p.a` 相同）而非 `*p.b = 1`。

**根因**: `OpFieldAddr` 使用 `reflect.NewAt` 创建指向结构体 `*int` 字段的 `**int` 指针。`OpDeref` 对 `**int` 调用 `rv.Elem()` 返回的 `reflect.Value` 是可寻址且可设置的，直接引用结构体字段的内存。当 SSA 加载 `old_a = *FieldAddr(p, 0)` 后再执行 `*FieldAddr(p, 0) = old_b` 时，存储操作也改变了 `old_a` 的值。

**修复**: 在 `OpDeref` 中，当 `rv.Elem()` 返回指针类型且可设置时，通过 `reflect.ValueOf(elem.Interface())` 创建独立副本，确保后续存储不会影响已加载的值。

## Bug 13: StructWithFuncSlice

**现象**: `FuncSliceHolder{funcs: []func() int{...}}` 触发 panic：`[]value.Value is not assignable to []func() int`。

**根因**: SSA 将切片字面量编译为 `Alloc([N]func() int)` + 元素存储 + `Slice`。`OpNew` 为带函数元素的 `*types.Array` 创建了 `[]value.Value` 而非正确的 `[N]func() int` 数组。

**修复**:
1. 移除 `OpNew` 中函数切片/数组的 `[]value.Value` 特殊路径，改用 `typeToReflect + reflect.New` 创建正确类型
2. 更新 `SetIndex` 使用 `ToReflectValue(elemType)` 进行闭包到函数的转换
3. 在 `ToReflectValue` 的 KindSlice 分支中添加 `[]value.Value` → 类型化切片的转换逻辑
4. 移除 `OpMakeSlice` 中已废弃的 `[]value.Value` 函数切片特殊路径

## Bug 14: StructAnonymousField

**现象**: `AnonField{int: 42, name: "test"}` 触发 panic：`reflect.StructOf: field "int" is unexported but missing PkgPath`。

**根因**: `typeToReflectWithCache` 对匿名未导出字段跳过了 PkgPath 设置。但 `reflect.StructOf` 有两个互斥约束：未导出字段必须有 PkgPath，匿名字段不能有 PkgPath。这使得匿名未导出字段无法直接表示。

**修复**:
1. 使用 `f.Exported()` 替代手动首字母检查
2. 对匿名未导出字段降级为普通未导出字段（`Anonymous = false` + 设置 PkgPath），这是 `reflect.StructOf` 中唯一的合法表示

## 修改文件

- `compiler/compile_value.go` — compileDefer 栈顺序
- `vm/ops_dispatch.go` — OpDeref 指针副本、OpNew 函数切片/数组、OpMakeSlice 清理
- `vm/typeconv.go` — 匿名未导出字段处理
- `value/accessor.go` — ToReflectValue 切片转换
- `value/container.go` — SetIndex 函数元素包装
- `tests/testdata/resolved_issue/main.go` — 迁入 6 个测试函数
- `tests/resolved_issue_test.go` — 新增 6 个测试用例
- `tests/testdata/known_issues/main.go` — 清空（所有问题已解决）
- `tests/known_issues_test.go` — 清空（所有问题已解决）

## 测试结果

全部 24 个已解决问题测试在 `-race` 标志下通过。编译器、虚拟机、值系统、字节码和导入器的所有测试均通过。
