# Gig 解释器 — 已知问题修复报告

**日期**: 2026-03-16  
**作者**: AI 辅助修复  
**分支**: `feature/dev_youngjin`  
**状态**: 全部 6 个 Bug 已修复，测试已迁移，全量回归测试通过

---

## 概述

本报告记录了 Gig Go 解释器中 6 个已知 Bug 的调查与修复过程。这些 Bug 的共同表现是：解释执行的结果与原生 Go 执行结果不一致。所有 Bug 原先在 `tests/known_issues_test.go` 中跟踪，现已全部修复并迁移至 `tests/resolved_issue_test.go`。

| # | Bug 名称 | 根因 | 修改文件 | 状态 |
|---|---------|------|---------|------|
| 1 | Map 中存储函数值 | Closure 未包装为真实 Go 函数 | `value/accessor.go`, `vm/vm.go`, `vm/closure.go`, `vm/ops_dispatch.go` | ✅ 已修复 |
| 2 | interface 切片的 type switch | interface 中 `int` 被存为 `int64` | `vm/ops_dispatch.go` | ✅ 已修复 |
| 3 | 结构体中的函数字段 | 与 Bug 1 根因相同 | `value/container.go` | ✅ 已修复 |
| 4 | 切片 append 展开运算符 | 原生 `[]int64` 与 reflect `[]int` 类型不匹配 | `vm/ops_dispatch.go` | ✅ 已修复 |
| 5 | range 遍历中修改 Map | 测试对非确定性行为使用了精确断言 | `tests/known_issues_test.go` | ✅ 已修复 |
| 6 | 自引用结构体类型 | `typeToReflect` 无限递归 | `vm/typeconv.go`, `vm/ops_dispatch.go` | ✅ 已修复 |

---

## Bug 1: Map 中存储函数值

### 症状
```
panic: reflect.Value.SetMapIndex: value of type *vm.Closure is not assignable to type func() int
```
将闭包存入 `map[int]func() int` 时触发 panic，原因是 VM 内部的 `*Closure` 类型无法赋值给具体的函数类型。

### 测试用例
```go
func MapWithFuncValue() int {
    m := make(map[int]func() int)
    m[1] = func() int { return 10 }
    m[2] = func() int { return 20 }
    return m[1]() + m[2]()   // 期望: 30
}
```

### 根因分析
在 Gig 的值系统中，闭包以 `*vm.Closure` 对象（`KindFunc`）的形式存储。当通过 reflect 对 `map[int]func() int` 调用 `SetMapIndex` 时，需要的是一个真正的 Go `func() int`，而不是 `*Closure`。`ToReflectValue` 方法原先缺少将闭包转换为真实 Go 函数的逻辑。

### 修复方案

**`value/accessor.go`** — 新增 `ClosureCaller` 回调类型，通过 `reflect.MakeFunc` 包装闭包：

```go
type ClosureCaller func(closure any, args []reflect.Value) []reflect.Value

// 在 ToReflectValue 的 KindFunc 分支:
case KindFunc:
    if typ.Kind() == reflect.Func && closureCaller != nil {
        closure := v.obj
        fn := reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
            results := closureCaller(closure, args)
            // 将返回值转换为期望类型（如 int64 → int）
            ...
        })
        return fn
    }
```

**`vm/vm.go`** — 在 `init()` 中注册 `ClosureCaller` 回调，打破 `value` → `vm` 的循环依赖：

```go
func init() {
    value.RegisterClosureCaller(func(closure any, args []reflect.Value) []reflect.Value {
        c := closure.(*Closure)
        // 创建临时 VM，执行闭包字节码，返回结果
        ...
    })
}
```

**`vm/closure.go`** — 为 `Closure` 结构体添加 `Program *bytecode.Program` 字段，使回调能够创建子 VM。

**`vm/ops_dispatch.go`** — 在 `OpClosure` 中设置 `closure.Program = vm.program`；在 `OpCallIndirect` 中增加对 reflect 函数值的处理（从类型化容器中读出的闭包会变成 `reflect.Value` 函数类型）。

**`vm/run.go`** — 更新 `OpCallIndirect` 热路径，支持通过 `rv.Call()` 调用 reflect 函数。

### 性能影响
`reflect.MakeFunc` 仅在闭包被赋值给类型化容器（map 值、结构体字段）时才会触发。普通闭包通过 `OpCallIndirect` 调用时仍走快速的 `*Closure` 路径，零额外开销。

---

## Bug 2: interface 切片元素的 type switch

### 症状
```
interface slice type switch: got 1110, want 1111
```
对 `[]interface{}` 中提取的值做 type switch 时，`int` 类型匹配失败（其他类型均正常匹配）。

### 测试用例
```go
func InterfaceSliceTypeSwitch() int {
    var items []interface{}
    items = append(items, 1, "hello", true, 3.14)
    count := 0
    for _, item := range items {
        switch item.(type) {
        case int:    count += 1
        case string: count += 10
        case bool:   count += 100
        case float64: count += 1000
        }
    }
    return count   // 期望: 1111
}
```

### 根因分析
两个问题叠加导致：

1. **值存储不匹配**: Gig 内部将所有整数存储为 `int64`。当 `1`（一个 `int`）被 append 到 `[]interface{}` 时，在 reflect 切片中实际存储为 `int64` 而非 `int`。而原生 Go 中 `interface{}` 里的 `1` 存储为 `int`。

2. **严格的 `AssignableTo` 检查**: `OpAssert` 处理器使用 `reflect.Type.AssignableTo()` 进行类型匹配。由于 `int64` 不能赋值给 `int`（它们是不同类型），`case int:` 分支永远不会匹配。

### 修复方案

**`vm/ops_dispatch.go`** — 在 `OpAssert` 的 `KindReflect` 分支增加 `sameReflectKindFamily()` 回退逻辑：

```go
if targetReflectType != nil && concreteRV.Type().AssignableTo(targetReflectType) {
    result = value.MakeFromReflect(concreteRV)
    assertionOk = true
} else if targetReflectType != nil && sameReflectKindFamily(concreteRV.Type(), targetReflectType) {
    // Gig 内部将 int 存储为 int64；type switch 时按 kind 族匹配
    result = value.MakeFromReflect(concreteRV.Convert(targetReflectType))
    assertionOk = true
}
```

同时增加了 `kindMatchesType()` 辅助函数，用于处理非 reflect 路径（原始 `KindInt`/`KindString` 等值），替换了之前的"默认成功"逻辑。

**`sameReflectKindFamily`** 在同一数值族内匹配类型：
- 有符号整数: `int`, `int8`, `int16`, `int32`, `int64`
- 无符号整数: `uint`, `uint8`, ... `uintptr`
- 浮点数: `float32`, `float64`
- 复数: `complex64`, `complex128`

### 附带修复
之前代码中有一个 `else { assertionOk = true }` 分支，会导致所有类型断言都判定为成功。这也导致 `TestCompiler_TypeAssertionCommaOk` 期望了错误的行为（`"hello".(int)` 的 `ok=true`）。测试已更新为期望正确结果（`ok=false`）。

---

## Bug 3: 结构体中的函数字段

### 症状
```
panic: reflect.Set: value of type value.Value is not assignable to type func() int
```
将闭包赋值给结构体的函数字段时触发 panic。

### 测试用例
```go
type structWithFunc struct {
    f func() int
}

func StructWithFuncField() int {
    s := structWithFunc{f: func() int { return 42 }}
    return s.f()   // 期望: 42
}
```

### 根因分析
与 Bug 1 是同一底层问题。另外，`value/container.go` 的 `SetElem` 方法对指向结构体的指针使用了 `reflect.ValueOf(val)` 而非 `val.ToReflectValue(elemType)`，跳过了闭包到函数的包装过程。

### 修复方案

**`value/container.go`** — 将 `SetElem` 改为使用 `ToReflectValue`：

```go
// 修复前（有问题）：
rv.Elem().Set(reflect.ValueOf(val))

// 修复后：
rv.Elem().Set(val.ToReflectValue(elemType))
```

这确保了闭包在赋值给结构体字段前，通过 `reflect.MakeFunc` 正确包装。

---

## Bug 4: 切片 append 展开运算符

### 症状
```
slice flatten: got 2, want 4
```
`append(result, inner...)` 只 append 了第一个元素，而不是全部元素。

### 测试用例
```go
func SliceFlatten() int {
    s := [][]int{{1, 2}, {3, 4}}
    result := []int{}
    for _, inner := range s {
        result = append(result, inner...)
    }
    return len(result)   // 期望: 4
}
```

### 根因分析
在 `OpAppend` 处理器的原生 `[]int64` 快速路径中，当 `elem.IntSlice()` 失败时（因为 `inner` 是从 `[][]int` range 得到的 reflect `[]int`），代码直接降级到 `elem.RawInt()`，将整个切片当作单个整数处理，只 append 了一个元素。

### 修复方案

**`vm/ops_dispatch.go`** — 在原生 `[]int64` 快速路径中增加 reflect 切片检测分支：

```go
if s, ok := slice.IntSlice(); ok {
    if es, ok2 := elem.IntSlice(); ok2 {
        vm.push(value.MakeIntSlice(append(s, es...)))
    } else if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice {
        // elem 是 reflect 整数切片（如从 [][]int range 得到的 []int）
        // 逐个转换为 int64 并展开 append
        for i := 0; i < elemRV.Len(); i++ {
            s = append(s, elemRV.Index(i).Int())
        }
        vm.push(value.MakeIntSlice(s))
    } else {
        vm.push(value.MakeIntSlice(append(s, elem.RawInt())))
    }
}
```

---

## Bug 5: range 遍历中修改 Map

### 症状
```
map update during range: got 6, want 7
```
在 `range` 遍历中添加键，产生的条目数少于预期。

### 测试用例
```go
func MapUpdateDuringRange() int {
    m := map[int]int{1: 10, 2: 20}
    for k := range m {
        m[k+10] = k
    }
    return len(m)   // Go 规范: 非确定性，但 >= 4
}
```

### 根因分析
Go 语言规范明确指出：在 `range` 遍历过程中添加键是允许的，但新添加的键**是否**会被遍历到是非确定性的。原始测试期望精确值 7，这假设了所有新键总是会被遍历到。

调查发现，`reflect.MapRange()`（VM 用于 map 遍历）已经能正确观察到部分新添加的键，符合 Go 的原生行为。问题出在测试的期望值过于严格。

### 修复方案

**`tests/known_issues_test.go`** — 将精确匹配改为范围检查：

```go
// 修复前: if n != 7 { ... }
// 修复后:
if n < 4 {
    t.Errorf("map update during range: got %d, want >= 4", n)
}
```

最小有效结果为 4（2 个原始键 + 仅遍历原始键时添加的 2 个新键）。

---

## Bug 6: 自引用结构体类型

### 症状
```
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp: ... stack: [...
fatal error: stack overflow
```
创建自引用结构体（如 `type node struct { next *node }`）时，`typeToReflect` 陷入无限递归导致栈溢出。

### 测试用例
```go
type node struct {
    value int
    next  *node
}

func StructSelfRef() int {
    n1 := &node{value: 1}
    n2 := &node{value: 2, next: n1}
    return n2.value + n2.next.value   // 期望: 3
}
```

### 根因分析
`typeToReflect` 递归地将 `go/types.Type` 转换为 `reflect.Type`。对于 `type node struct { next *node }`，这会产生无限循环：`node → struct{int, *node} → *node → node → ...`

### 修复方案（两部分）

**第 1 部分: `vm/typeconv.go`** — 添加循环检测缓存：

```go
func typeToReflectWithCache(t types.Type, cache map[types.Type]reflect.Type) reflect.Type {
    if cached, ok := cache[t]; ok {
        return cached  // 对正在处理的类型返回 nil（检测到循环）
    }
    // 对 *types.Named：递归前先标记为处理中
    case *types.Named:
        cache[tt] = nil  // 循环检测哨兵值
        result := typeToReflectWithCache(tt.Underlying(), cache)
        cache[tt] = result
        return result
    // 对 *types.Pointer：当 elem 为 nil（循环）时，使用 interface{} 作为占位符
    case *types.Pointer:
        elem := typeToReflectWithCache(tt.Elem(), cache)
        if elem != nil {
            return reflect.PointerTo(elem)
        }
        return reflect.TypeOf(&emptyIface).Elem()  // interface{} 占位符
}
```

**第 2 部分: `vm/ops_dispatch.go`** — 更新 `OpFieldAddr` 以处理 `interface{}` 占位符：

当通过自引用指针访问字段时（在 reflect 结构体中存储为 `interface{}`），reflect.Value 的 kind 会是 `reflect.Interface`，内部包装着实际的结构体指针。增加了解包装逻辑：

```go
case bytecode.OpFieldAddr:
    // ... 现有的指针解引用 ...
    // 对于自引用结构体类型，typeToReflect 将递归指针字段
    // 存储为 interface{}。在这里进行解包装。
    if s.Kind() == reflect.Interface && !s.IsNil() {
        s = s.Elem()
        if s.Kind() == reflect.Ptr {
            s = s.Elem()
        }
    }
    if s.Kind() == reflect.Struct { ... }
```

### 设计决策: `interface{}` vs `unsafe.Pointer`
最初考虑使用 `unsafe.Pointer` 作为占位符类型，但这会导致 `reflect.Set` 错误，因为 `*struct{...}` 无法赋值给 `unsafe.Pointer`。最终选择 `interface{}`，因为任何 Go 值都可以存入 interface，而且 VM 在其他场景中已经具备了 interface 解包装的能力。

---

## 验证结果

### 测试结果
```
$ go test -race ./...
ok   github.com/t04dJ14n9/gig              1.055s
ok   github.com/t04dJ14n9/gig/bytecode     1.014s
ok   github.com/t04dJ14n9/gig/compiler     1.013s
ok   github.com/t04dJ14n9/gig/importer     1.014s
ok   github.com/t04dJ14n9/gig/tests       49.273s
ok   github.com/t04dJ14n9/gig/value        1.014s
ok   github.com/t04dJ14n9/gig/vm           1.124s
```

### Lint 结果
```
$ golangci-lint-v2 run
（无问题）
```

### 测试迁移
全部 6 个测试已从 `tests/known_issues_test.go` 迁移至 `tests/resolved_issue_test.go`：
- `TestResolved_MapWithFuncValue`
- `TestResolved_InterfaceSliceTypeSwitch`
- `TestResolved_StructWithFuncField`
- `TestResolved_SliceFlatten`
- `TestResolved_MapUpdateDuringRange`
- `TestResolved_StructSelfRef`

---

## 修改文件一览

| 文件 | 修改内容 |
|------|---------|
| `value/accessor.go` | 新增 `ClosureCaller` 类型，在 `ToReflectValue` 中通过 `reflect.MakeFunc` 包装闭包 |
| `value/container.go` | 修复 `SetElem`，对函数类型指针字段使用 `ToReflectValue` |
| `vm/vm.go` | 在 `init()` 中注册 `ClosureCaller` 回调 |
| `vm/closure.go` | 为 `Closure` 结构体添加 `Program` 字段 |
| `vm/ops_dispatch.go` | 修复 `OpAssert`（type switch）、`OpAppend`（展开）、`OpFieldAddr`（自引用）、`OpClosure`（设置 Program）、`OpCallIndirect`（reflect 函数） |
| `vm/run.go` | 更新 `OpCallIndirect` 热路径，支持 reflect 函数 |
| `vm/typeconv.go` | 为 `typeToReflect` 添加循环检测缓存 |
| `tests/known_issues_test.go` | 已清空（所有问题已解决） |
| `tests/testdata/known_issues/main.go` | 已清空（所有问题已解决） |
| `tests/resolved_issue_test.go` | 新增 6 个迁移后的测试函数 |
| `tests/testdata/resolved_issue/main.go` | 新增 6 个迁移后的测试数据函数 |
| `tests/compiler_vm_test.go` | 更新 `TypeAssertionCommaOk` 期望正确行为 |
