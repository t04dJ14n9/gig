# Gig CLI 指南

`gig` 命令行工具提供三个主要功能：依赖包生成和交互式 REPL（Read-Eval-Print-Loop）Go 代码解释器。

## 安装

```bash
# 安装 CLI 工具
go install git.woa.com/youngjin/gig/cmd/gig@latest

# 或者直接运行 (Go 1.21+)
go run git.woa.com/youngjin/gig/cmd/gig@latest --help
```

## 命令

### `gig init` - 初始化依赖包

创建一个新的依赖包目录，包含 `pkgs.go` 模板文件。

```bash
gig init -package <名称>
```

**选项：**
- `-package <名称>` - 要创建的依赖包名称（必需）

**示例：**
```bash
gig init -package mydep
```

创建的目录结构：
```
mydep/
└── pkgs.go    # 用于添加包导入的模板文件
```

生成的 `pkgs.go` 内容：
```go
package mydep

// 编辑此文件以添加/删除包导入。
// 然后运行: gig gen ./mydep

import (
	// 标准库示例
	// _ "fmt"
	// _ "strings"
	// _ "time"

	// 第三方库示例
	// _ "github.com/spf13/cast"
)
```

### `gig gen` - 生成注册代码

为依赖包中导入的包生成注册代码。

```bash
gig gen <目录>
```

**参数：**
- `<目录>` - 依赖包目录路径

**示例：**
```bash
gig gen ./mydep
```

编辑 `mydep/pkgs.go` 后：
```go
package mydep

import (
	_ "fmt"
	_ "strings"
	_ "github.com/spf13/cast"
)
```

运行 `gig gen ./mydep` 生成：
```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    └── github_com_spf13_cast.go
```

每个生成的文件包含该包导出符号（函数、常量、变量和类型）的注册代码。

### `gig repl` - 交互式 REPL

启动支持外部包热加载的交互式 Go 解释器。

```bash
gig repl
```

## REPL 命令

REPL 支持以 `:` 开头的特殊命令：

| 命令 | 别名 | 描述 |
|------|------|------|
| `:help` | `:h`, `:?` | 显示帮助信息 |
| `:quit` | `:q`, `:exit` | 退出 REPL |
| `:clear` | - | 清除会话状态（变量、导入、声明） |
| `:imports` | - | 列出已导入的包 |
| `:vars` | - | 列出已捕获的变量 |
| `:env` | - | 列出所有变量和包 |
| `:source` | - | 显示累积的声明 |
| `:timeout <时长>` | - | 设置执行超时（如 `:timeout 5s`） |
| `:plugins` | - | 列出已加载的插件（热加载的包） |

### `:env` - 环境概览

`:env` 命令提供当前会话状态的全面视图：

```
>>> import "fmt"
>>> x := 5
>>> y := x + 3
>>> :env
Packages:
  fmt (fmt)

Variables:
  x int = 5
  y int = 8
```

## REPL 特性

### 输入类型

REPL 自动分类和处理不同类型的输入：

| 类型 | 示例 | 描述 |
|------|------|------|
| Import | `import "fmt"` | 导入包 |
| Expression | `1+1`, `x + y` | 求值并打印结果 |
| Statement | `x := 1`, `for i:=0; i<5; i++ {}` | 执行语句 |
| Declaration | `func foo() {}`, `type Point struct{}` | 存储以供后续使用 |

### 多行输入

REPL 支持复杂代码的多行输入：

```
>>> func fibonacci(n int) int {
...     if n <= 1 {
...         return n
...     }
...     return fibonacci(n-1) + fibonacci(n-2)
... }
>>> fibonacci(10)
55
```

当输入包含未闭合的括号（`{}`, `()`, `[]`）或反引号（`` ` ``）时，多行模式自动激活。

### Tab 自动补全

REPL 提供智能的 Tab 补全功能：

- **命令**：输入 `:` 后按 `Tab` 循环浏览命令
- **变量**：从当前会话补全变量名
- **包**：从导入和已知包补全包名
- **包符号**：输入 `pkg.` 后，补全导出符号（适用于热加载的包）

**已知包包括**：`fmt`、`strings`、`strconv`、`math`、`time`、`bytes`、`errors`、`sort`、`regexp`、`json`、`xml`、`http`、`url` 等 40 多个。

### 变量持久化

使用短声明（`:=`）定义的变量在语句间持久存在：

```
>>> x := 10
>>> y := 20
>>> x + y
30
>>> :vars
Variables:
  x int = 10
  y int = 20
```

### 热加载外部包

REPL 可以使用 Go 的插件系统热加载外部 Go 包（仅限 Linux 和 macOS）：

```
>>> import "github.com/spf13/cast"
Imported: github.com/spf13/cast (hot-loaded)
>>> cast.ToString(123)
"123"
>>> cast.ToInt("456")
456
```

**平台支持：**
- ✅ Linux - 完全支持
- ✅ macOS - 完全支持
- ❌ Windows - 不支持（返回错误）

**工作原理：**
1. 使用 `go get` 下载包
2. 使用 `go doc` 发现符号，生成包装代码
3. 使用 `go build -buildmode=plugin` 编译为 `.so` 插件
4. 加载插件并向 gig 注册符号

**插件缓存：**
插件缓存在 `~/.gig/plugins/` 目录：
```
~/.gig/plugins/
├── github.com/
│   └── spf13/
│       └── cast/
│           ├── cast.so      # 编译的插件
│           └── wrapper.go   # 生成的包装代码
├── go.mod                   # 插件模块文件
└── plugin_registry.json     # 插件元数据
```

## 示例

### 基本 REPL 会话

```
>>> import "fmt"
Imported: fmt
>>> import "strings"
Imported: strings
>>> name := "world"
>>> fmt.Sprintf("Hello, %s!", strings.ToUpper(name))
"Hello, WORLD!"
>>> :env
Packages:
  fmt (fmt)
  strings (strings)

Variables:
  name string = "world"
```

### 使用函数

```
>>> func add(a, b int) int { return a + b }

>>> func multiply(a, b int) int { return a * b }

>>> add(2, 3)
5
>>> multiply(4, 5)
20
>>> :source
Accumulated declarations:
func add(a, b int) int { return a + b }

func multiply(a, b int) int { return a * b }
```

### 自定义类型

```
>>> type Point struct { X, Y int }

>>> p := Point{X: 10, Y: 20}
>>> p.X + p.Y
30
>>> :vars
Variables:
  p main.Point = {X:10 Y:20}
```

### 使用第三方包

```
>>> import "github.com/spf13/cast"
Imported: github.com/spf13/cast (hot-loaded)
>>> cast.ToInt("123")
123
>>> cast.ToBool("true")
true
>>> cast.ToString(456.78)
"456.78"
>>> :plugins
Loaded plugins:
  github.com/spf13/cast
```

### 控制流

```
>>> sum := 0
>>> for i := 1; i <= 10; i++ { sum += i }

>>> sum
55
>>> :vars
Variables:
  sum int = 55
```

### 使用上下文和超时

```
>>> :timeout 5s
Timeout set to: 5s
>>> :timeout
Current timeout: 5s
```

超时应用于表达式求值和语句执行，防止无限循环导致 REPL 无响应。

## 故障排除

### 插件加载错误

**"plugin was built with a different version of package"**

当插件使用的主机不同版本的 gig 时会发生此错误。插件管理器通过在插件的 `go.mod` 中使用 `replace` 指令自动处理此问题。

**"plugin loading requires Linux or macOS"**

Windows 不支持 Go 的插件系统。建议：
- 使用 WSL（Windows Subsystem for Linux）
- 使用 `gig gen` 预生成注册代码
- 使用内置的标准库包

### 导入错误

**"package not found; run 'go get <package>' first"**

包未安装。运行：
```bash
go get <包路径>
```

然后在 REPL 中重试导入。

## 配置

### 超时

默认执行超时为 10 秒。可通过以下方式修改：
```
:timeout 30s
```

### 已知包

REPL 内置了常见包名到导入路径的映射：

| 名称 | 导入路径 |
|------|----------|
| `fmt` | `fmt` |
| `strings` | `strings` |
| `json` | `encoding/json` |
| `http` | `net/http` |
| `url` | `net/url` |
| `rand` | `math/rand` |
| ... | ... |

这允许在未显式导入时自动检测包函数的使用。

## 相关文档

- [README_CN.md](../README_CN.md) - 主要文档
- [gig-internals_CN.md](gig-internals_CN.md) - 内部架构
- [examples/](../examples/) - 示例程序
