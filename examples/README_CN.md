# Gig 示例

[![中文](https://img.shields.io/badge/lang-中文-red.svg)](README_CN.md) [![English](https://img.shields.io/badge/lang-English-blue.svg)](README.md)

本目录包含演示如何使用 Gig（高性能 Go 解释器）的示例。

## 示例概览

| 示例 | 描述 | 难度 |
|------|------|------|
| [simple](./simple) | 使用 Gig 内置标准库 | 初学者 |
| [custom](./custom) | 使用 Gig 自定义/第三方依赖 | 中级 |

## 快速开始

### 简单示例（推荐首次使用）

使用 Gig 最简单的方式是使用其内置标准库支持。只需导入 `gig/stdlib/packages`，即可访问 40+ 标准库包。

```bash
cd simple
go run main.go
```

### 自定义示例（用于第三方库）

当需要第三方库或想要最小化依赖时，使用 CLI 工具生成注册代码。

```bash
cd custom
go run main.go
```

---

## 简单示例详情

**位置：** `./simple/`

**使用场景：** 快速原型开发、脚本、仅使用标准库的规则引擎。

**演示的关键特性：**
- 使用循环和变量的基本计算
- 使用标准库包（`fmt`、`strings`、`math`、`time`）
- 基于上下文的超时控制
- 多函数程序

### 代码结构

```
simple/
├── go.mod
├── go.sum
└── main.go
```

### 使用模式

```go
import (
    "github.com/t04dJ14n9/gig"
    _ "github.com/t04dJ14n9/gig/stdlib/packages" // 导入内置标准库（40+ 包）
)

func main() {
    source := `
    package main
    
    import "fmt"
    import "strings"
    
    func Greet(name string) string {
        return fmt.Sprintf("Hello, %s!", strings.ToUpper(name))
    }
    `
    
    prog, err := gig.Build(source)
    if err != nil {
        panic(err)
    }
    
    result, err := prog.Run("Greet", "world")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result) // 输出: Hello, WORLD!
}
```

### 内置包

Gig 的内置标准库包括：

| 类别 | 包 |
|------|-----|
| **I/O** | `fmt`、`io`、`bufio`、`bytes`、`strings` |
| **编码** | `encoding/json`、`encoding/base64`、`encoding/hex`、`encoding/xml` |
| **文本** | `strings`、`strconv`、`text/template`、`regexp` |
| **数学** | `math`、`math/rand` |
| **时间** | `time` |
| **集合** | `sort`、`container/list`、`container/heap` |
| **网络** | `net/url`、`net/http`（部分） |
| **加密** | `crypto/hmac`、`crypto/sha256` |
| **其他** | `errors`、`sync`、`context`、`path`、`path/filepath`、`os`（部分） |

### 上下文超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

result, err := prog.RunWithContext(ctx, "SlowFunction", args...)
if errors.Is(err, context.DeadlineExceeded) {
    // 处理超时
}
```

---

## 自定义示例详情

**位置：** `./custom/`

**使用场景：** 需要第三方库的生产应用、最小依赖占用或自定义包子集。

**演示的关键特性：**
- 使用第三方库（github.com/tidwall/gjson）
- 外部类型上的方法调用（`.String()`、`.Int()`、`.Bool()`）
- 自定义依赖包生成工作流

### 代码结构

```
custom/
├── go.mod
├── go.sum
├── main.go
└── mydep/
    ├── pkgs.go              # 依赖声明
    └── packages/            # 生成的注册代码
        ├── fmt.go
        ├── strings.go
        ├── github_com_tidwall_gjson.go
        └── ...
```

### 设置工作流

#### 步骤 1：安装 CLI

```bash
go install github.com/t04dJ14n9/gig/cmd/gig@latest
```

#### 步骤 2：初始化依赖包

```bash
gig init -package mydep
```

这将创建：
```
mydep/
└── pkgs.go    # 编辑此文件以声明依赖
```

#### 步骤 3：声明依赖

编辑 `mydep/pkgs.go`：

```go
package mydep

import (
    // 标准库（保留需要的）
    _ "fmt"
    _ "strings"
    _ "time"

    // 第三方库
    _ "github.com/tidwall/gjson"
    _ "github.com/spf13/cast"
)
```

#### 步骤 4：生成注册代码

```bash
gig gen ./mydep
```

这将生成：
```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    ├── time.go
    ├── github_com_tidwall_gjson.go
    └── github_com_spf13_cast.go
```

#### 步骤 5：在程序中使用

```go
import (
    "github.com/t04dJ14n9/gig"
    _ "myapp/mydep/packages" // 你的自定义依赖包
)

func main() {
    source := `
    package main
    
    import "github.com/tidwall/gjson"
    
    func GetValue(json string, path string) string {
        return gjson.Get(json, path).String()
    }
    `
    
    prog, _ := gig.Build(source)
    result, _ := prog.Run("GetValue", `{"name":"Alice"}`, "name")
    fmt.Println(result) // 输出: Alice
}
```

### 外部类型方法

Gig 支持在外部类型上调用方法：

```go
source := `
package main

import "github.com/tidwall/gjson"

func GetUserAge(json string) int64 {
    return gjson.Get(json, "age").Int()  // gjson.Result 上的方法调用
}
`
```

支持的方法模式：
- 值接收者方法：`.String()`、`.Int()`、`.Bool()`、`.Float()`
- 指针接收者方法：`.Scan()`、`.ForEach()`
- 链式方法调用：`gjson.Get(json, "arr").Array()[0].String()`

---

## 对比：简单 vs 自定义

| 方面 | 简单 | 自定义 |
|------|------|--------|
| **设置** | 无（只需导入） | CLI 工具 + 代码生成 |
| **依赖** | 40+ 标准库包 | 只包含你需要的 |
| **第三方库** | 不支持 | 完全支持 |
| **二进制大小** | 较大 | 较小（tree-shaking） |
| **构建时间** | 更快（预构建） | 较慢（代码生成） |
| **使用场景** | 原型开发、脚本 | 生产环境、嵌入式 |

---

## API 参考

### 构建

```go
// Build 解析并编译 Go 源代码
prog, err := gig.Build(source string) (*Program, error)
```

### 运行

```go
// Run 按名称执行函数（默认 10 秒超时）
result, err := prog.Run(funcName string, args ...interface{}) (interface{}, error)

// RunWithContext 带上下文执行，支持取消
result, err := prog.RunWithContext(ctx context.Context, funcName string, args ...interface{}) (interface{}, error)
```

### 包注册

```go
import "github.com/t04dJ14n9/gig/register"

pkg := register.AddPackage("mypkg", "mypkg")
pkg.NewFunction("MyFunc", MyFunc, "文档说明")
pkg.NewConst("MyConst", MyConst, "文档说明")
pkg.NewVar("MyVar", &MyVar, "文档说明")
```

---

## 安全考虑

Gig 通过在解释代码中禁止某些导入来强制安全性：

| 禁止 | 原因 |
|------|------|
| `unsafe` | 内存安全 |
| `reflect` | 类型安全 |
| `panic` | 受控执行 |

---

## 支持的语言特性

| 特性 | 状态 |
|------|------|
| 算术运算 | 完全支持 |
| 变量和赋值 | 完全支持 |
| 控制流（if/else、for、switch） | 完全支持 |
| 函数和递归 | 完全支持 |
| 多返回值 | 完全支持 |
| 闭包 | 完全支持 |
| 字符串操作 | 完全支持 |
| 切片和数组 | 完全支持 |
| 映射（Map） | 完全支持 |
| 结构体和方法 | 完全支持 |
| 接口 | 完全支持 |
| Goroutine | 基础支持 |
| Channel | 基础支持 |
| 基于上下文的超时 | 完全支持 |
| 外部 Go 函数调用 | 完全支持 |

---

## 故障排除

### "package not registered"（包未注册）

**问题：** 你正在导入 Gig 不知道的包。

**解决方案：**
- 对于标准库：确保你导入了 `_ "github.com/t04dJ14n9/gig/stdlib/packages"`
- 对于第三方库：使用 CLI 生成注册代码

### "method not found on external type"（外部类型上找不到方法）

**问题：** 在外部类型上调用方法失败。

**解决方案：** 确保方法是导出的，并且类型已通过 `gig gen` 正确注册。

### 超时错误

**问题：** 函数执行超时。

**解决方案：**
- 增加超时时间：`context.WithTimeout(ctx, longerDuration)`
- 检查解释代码中是否存在无限循环

### 全局变量问题

**问题：** 全局变量未按预期工作。

**解决方案：** 这是一个已知限制。使用函数局部变量或将数据作为参数传递。
