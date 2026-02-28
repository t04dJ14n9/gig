# DirectCall：消除外部函数和方法的 reflect.Value.Call()

## 问题

当 Gig 调用外部 Go 包的函数（如 `strings.Contains`、`fmt.Sprintf`）时，它必须跨越解释器的 `value.Value` 表示和 Go 原生类型之间的边界。在此优化之前，**每次外部调用**都经过：

1. 将 `[]value.Value` 参数转换为 `[]reflect.Value`（分配密集）
2. `reflect.Value.Call()`（慢：安全检查、类型验证、`runtime.call`）
3. 将 `[]reflect.Value` 结果转换为 `[]value.Value`

这使得外部函数调用成为 VM 中最慢的操作——与原生 Go 调用相比大约有 **50 倍的开销**。

## 解决方案：生成类型化包装器（DirectCall）

### 函数 DirectCall

对于每个具有兼容参数类型的外部函数，我们在代码生成时生成一个类型化的 Go 包装器。例如，`strings.Contains(s, substr string) bool` 获得：

```go
func directcall_Contains(args []value.Value) value.Value {
    a0 := args[0].String()
    a1 := args[1].String()
    r0 := strings.Contains(a0, a1)
    return value.FromBool(r0)
}
```

无 `reflect.Value` 分配。无 `reflect.Value.Call()`。只有直接从 `value.Value` 的类型化提取、原生 Go 函数调用和直接的结果包装。

### 方法 DirectCall（新增）

将同样的方法扩展到**外部类型上的方法**。例如，`(*bytes.Buffer).WriteString(s string) (int, error)` 获得：

```go
func directcall_method_Buffer_WriteString(args []value.Value) value.Value {
    recv := args[0].Interface().(*bytes.Buffer)
    a1 := args[1].String()
    r0, r1 := recv.WriteString(a1)
    // ... 返回元组
}
```

接收者通过 `.Interface().(T)` 类型断言提取——运行时无需 `reflect.MethodByName` 查找。

## 架构

### 代码生成流水线（`gentool/`）

```
包类型信息 (go/types)
    │
    ├── directcall.go: generateDirectCall()       → 函数包装器
    ├── directcall.go: generateMethodDirectCalls() → 方法包装器
    ├── resolve.go:    collectCrossPkgImports()    → 导入解析
    │                  collectMethodImports()
    └── generator.go:  编排 + 输出
           │
           ▼
    stdlib/packages/*.go  （生成的，1162 个包装器）
```

### 参数类型支持

| 类型类别 | 示例 | 提取方式 |
|---|---|---|
| 基本类型 | `int`、`string`、`bool`、`float64` | `.Int()`、`.String()`、`.Bool()`、`.Float()` |
| 同包命名类型 | `Regexp`、`Template` | `.Interface().(TypeName)` |
| 跨包命名类型 | `time.Time`、`io.Reader` | `.Interface().(pkg.Type)` |
| 命名类型指针 | `*bytes.Buffer`、`*http.Request` | `.Interface().(*pkg.Type)` |
| 基本类型指针 | `*int32`、`*int64` | `.Interface().(*int32)` |
| 切片类型 | `[]byte`、`[]string` | `.Bytes()`、`.Interface().([]string)` |
| 映射类型 | `map[string]bool` | `.Interface().(map[string]bool)` |
| 空接口 | `any` / `interface{}` | `.Interface()` |
| 错误接口 | `error` | 通过 `value.ErrorFromValue()` 转换 |

### 编译时解析

编译器在编译时解析 DirectCall 包装器（`compiler/compile_ext.go`），将它们存储在 `ExternalFuncInfo.DirectCall` 和 `ExternalMethodInfo.DirectCall` 中。在运行时，VM 检查 `DirectCall != nil` 并直接调用——零映射查找。

### SSA 外部方法 Pkg=nil 问题

对于 SSA 中的外部包方法，`fn.Pkg`、`fn.Object().Pkg()` 甚至 `named.Obj().Pkg()` 都是 `nil`，因为外部类型缺少正确的包绑定。我们通过以 `typeName.methodName`（不含包路径）为键来索引方法 DirectCall 注册表来解决这个问题，因为类型名在标准库包之间已足够唯一。

## 覆盖率

| 类别 | 包装器数量 | 覆盖率 |
|---|---|---|
| 函数 DirectCall | 619 / 671 | 92.2% |
| 方法 DirectCall | 543 | 所有符合条件的方法 |
| **总计** | **1,162** | — |

剩余约 8% 的函数使用无法静态包装的参数类型（如 `unsafe.Pointer`、具有复杂元素类型的可变参数）。

## 基准测试结果

### 外部调用基准测试（5 次运行，`benchstat`）

| 基准测试 | 基线 (ns/op) | 优化后 (ns/op) | 加速比 | 内存变化 | 分配变化 |
|---|---|---|---|---|---|
| ExtCallReflect | 1,319,800 | 359,100 | **3.7x** (−72.8%) | −62.9% | −64.6% |
| ExtCallMethod | 1,216,000 | 460,100 | **2.6x** (−62.2%) | −49.0% | −50.3% |
| ExtCallMixed | 730,300 | 330,500 | **2.2x** (−54.8%) | −39.4% | −45.1% |
| ExtCallDirectCall | 588,000 | 583,500 | ~不变 | ~ | ~ |

ExtCallDirectCall 在基线中已使用函数 DirectCall——改进来自之前的工作。其他三个基准测试的巨大收益来自**方法 DirectCall** 和**扩展的函数 DirectCall 覆盖率**（从约 460 个增加到 619 个包装器）。

### Gig vs Yaegi（优化后）

| 基准测试 | Gig (ns/op) | Yaegi (ns/op) | Gig 优势 |
|---|---|---|---|
| ExtCallDirectCall | 583,500 | 1,551,000 | **快 2.7 倍** |
| ExtCallReflect | 359,100 | 1,001,500 | **快 2.8 倍** |
| ExtCallMethod | 460,100 | 1,214,000 | **快 2.6 倍** |
| ExtCallMixed | 330,500 | 845,900 | **快 2.6 倍** |

### 核心 VM 基准测试（无回退）

所有核心 VM 基准测试（Fib25、ArithSum、BubbleSort、Sieve、ClosureCalls）无统计显著变化——该优化纯粹是增量式的。

## 修改的文件

| 文件 | 角色 |
|---|---|
| `gentool/directcall.go` | 函数和方法的核心包装器生成 |
| `gentool/generator.go` | 编排生成，输出方法包装器 |
| `gentool/resolve.go` | 方法签名的跨包导入收集 |
| `bytecode/bytecode.go` | 添加 `ExternalMethodInfo.DirectCall` 字段 |
| `compiler/compile_ext.go` | 方法的编译时 DirectCall 解析 |
| `vm/call.go` | 运行时快速路径：`DirectCall != nil` → 直接调用 |
| `importer/register.go` | 方法 DirectCall 注册表（`AddMethodDirectCall` / `LookupMethodDirectCall`） |
| `gig.go` | `packageLookupAdapter` 连接方法 DirectCall |
| `stdlib/packages/*.go` | 20 个重新生成的包，共 1,162 个包装器 |
| `benchmarks/bench_test.go` | 12 个新基准测试（4 Gig + 4 Yaegi + 4 原生） |
