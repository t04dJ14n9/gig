# Gig - Go 语言实现的Go 解释器

[![中文](https://img.shields.io/badge/lang-中文-red.svg)](README.md) [![English](https://img.shields.io/badge/lang-English-blue.svg)](README_EN.md)

Gig 是一个用 Go 语言编写的高性能 Go 解释器，采用 SSA 到字节码的编译方式和基于栈的虚拟机。

> **说明**：本项目大量使用 AI 工具进行开发。它包含了全面的测试（40+ 测试文件）和基准测试，以确保正确性和性能。

## 特性

- **基于 SSA 的编译**：使用 `golang.org/x/tools/go/ssa` 作为中间表示
- **基于栈的虚拟机**：高效字节码执行，开销极小
- **Tagged-Union 值系统**：基本类型零反射开销
- **安全性**：在解释代码中禁止 `unsafe`、`reflect` 和 `panic`
- **可扩展**：支持注册外部 Go 包（内置 40+ 标准库包）
- **Context 取消支持**：完整支持 `context.Context` 超时和取消（[文档](docs/context-cancellation_CN.md)）

## 安装

```bash
go get git.woa.com/youngjin/gig
```

## 快速开始

### 方式一：使用内置标准库（推荐）

Gig 内置了 40+ 标准库包，只需导入 `gig/stdlib/packages`：

```go
package main

import (
    "fmt"
    _ "git.woa.com/youngjin/gig/stdlib/packages" // 导入 gig 的内置标准库
    "git.woa.com/youngjin/gig"
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

**内置包包括**：`fmt`、`strings`、`strconv`、`math`、`time`、`bytes`、`errors`、`sort`、`regexp`、`encoding/json`、`encoding/base64`、`net/url` 等 30 多个。

### 方式二：使用自定义依赖

如果需要第三方库或标准库子集，使用 `gig` CLI 工具：

#### 步骤 1：安装 CLI

```bash
# 安装 CLI 工具
go install git.woa.com/youngjin/gig/cmd/gig@latest

# 或直接运行（Go 1.21+）
go run git.woa.com/youngjin/gig/cmd/gig@latest --help
```

#### 步骤 2：初始化依赖包

```bash
# 创建名为 "mydep" 的依赖包
gig init -package mydep
```

这将创建：

```
mydep/
└── pkgs.go    # 编辑此文件以添加/移除包
```

#### 步骤 3：自定义依赖

编辑 `mydep/pkgs.go` 添加第三方库：

```go
package mydep

import (
    // 标准库（保留需要的）
    _ "fmt"
    _ "strings"
    _ "time"

    // 第三方库
    _ "github.com/spf13/cast"
    _ "github.com/tidwall/gjson"
)
```

#### 步骤 4：生成注册代码

```bash
# 从 pkgs.go 生成注册代码
gig gen ./mydep
```

这将生成：

```
mydep/
├── pkgs.go
└── packages/
    ├── fmt.go
    ├── strings.go
    ├── github_com_spf13_cast.go
    └── github_com_tidwall_gjson.go
```

#### 步骤 5：在程序中使用

```go
package main

import (
    "fmt"
    _ "myapp/mydep/packages" // 你的自定义依赖包
    "git.woa.com/youngjin/gig"
)

func main() {
    source := `
package main

import "github.com/tidwall/gjson"

func GetJsonValue(json string, path string) string {
    return gjson.Get(json, path).String()
}
`

    prog, _ := gig.Build(source)
    result, _ := prog.Run("GetJsonValue", `{"name":"Alice"}`, "name")
    fmt.Println(result) // 输出: Alice
}
```

## API 参考

### 构建和运行

```go
// Build 解析并编译 Go 源代码
prog, err := gig.Build(source string) (*Program, error)

// Run 按名称执行函数
result, err := prog.Run(funcName string, args ...interface{}) (interface{}, error)

// RunWithContext 带上下文执行，支持取消（ctx 是第一个参数）
result, err := prog.RunWithContext(ctx context.Context, funcName string, args ...interface{}) (interface{}, error)
```

### 注册包（高级）

```go
import "git.woa.com/youngjin/gig/importer"

// 手动注册包（通常通过生成的代码完成）
pkg := importer.RegisterPackage("mypkg", "mypkg")
pkg.AddFunction("MyFunc", MyFunc, "", directCall_MyFunc)
pkg.AddConstant("MyConst", MyConst, "")
pkg.AddVariable("MyVar", &MyVar, "")
pkg.AddType("MyType", reflect.TypeOf(MyType{}), "")
```

## 示例

参见 `examples/` 目录：

- **`examples/simple/`** - 使用内置标准库（最简单）
- **`examples/custom/`** - 使用自定义依赖

运行示例：

```bash
# 简单示例（使用内置标准库）
cd gig/examples/simple
go run main.go

# 自定义示例
cd gig/examples/custom
go run main.go
```

## gig CLI 命令

```bash
# 初始化依赖包
gig init -package <名称>

# 生成注册代码
gig gen <目录>

# 示例
gig init -package mydep         # 创建 mydep/pkgs.go
gig gen ./mydep                 # 在 myapp/mydep/packages/ 生成注册代码
```

## 支持的特性

- ✅ 算术运算
- ✅ 变量和赋值
- ✅ 控制流（if/else、for 循环、switch）
- ✅ 函数和递归
- ✅ 多返回值
- ✅ 闭包
- ✅ Defer、Panic 和 Recover
- ✅ 字符串操作
- ✅ 切片和数组
- ✅ 映射（Map）
- ✅ 结构体和方法
- ✅ 接口
- ✅ Goroutine（基础支持）
- ✅ 基于上下文的超时控制
- ✅ 外部 Go 函数调用

## 性能

在同一台机器上使用相同算法，对比 **Gig**、**Yaegi**（Go 解释器）、**GopherLua**（Lua 解释器）和 **原生 Go** 的真实基准测试。

> **测试环境**：AMD EPYC 9754 128 核, 32 线程, Linux amd64, Go 1.23.1
> 使用 `-count=3` 运行。源码：[`benchmarks/bench_test.go`](benchmarks/bench_test.go), [`benchmarks/stress_test.go`](benchmarks/stress_test.go), [`benchmarks/extreme_stress_test.go`](benchmarks/extreme_stress_test.go)

### 核心工作负载 (Gig vs Yaegi vs GopherLua vs 原生 Go)

| 工作负载                      | 原生 Go |         Gig |      Yaegi |  GopherLua |      Gig vs Yaegi | 原生/Gig |
| ----------------------------- | ------: | ----------: | ---------: | ---------: | ----------------: | -------: |
| **Fibonacci(25)** 递归        |  450 μs | **23.7 ms** |   102.4 ms |   24.3 ms | **Gig 快 4.3 倍** |   52.7x |
| **ArithmeticSum(1K)** 循环    |  336 ns | **74.0 μs** |    40.8 μs |    37.4 μs |    Yaegi 快 1.8 倍 |  220x |
| **BubbleSort(100)** 嵌套循环  |  6.2 μs | **953.5 μs** | 1,235.3 μs |   775.9 μs | **Gig 快 1.3 倍** |  154x |
| **Sieve(1000)** 质数筛        |  1.7 μs | **270.3 μs** |   202.1 μs |   203.7 μs |    Yaegi 快 1.3 倍 |  159x |
| **ClosureCalls(1K)** 闭包调用 |  660 ns | **389.4 μs** |   943.4 μs |   118.6 μs | **Gig 快 2.4 倍** |  590x |

### 外部函数调用 (Gig vs Yaegi vs 原生 Go)

从解释代码调用 Go 标准库函数 —— 最常见的实际使用场景：

| 工作负载                           | 原生 Go |        Gig |      Yaegi |      Gig vs Yaegi | 原生/Gig |
| ---------------------------------- | ------: | ---------: | ---------: | ----------------: | -------: |
| **DirectCall** (strings/strconv)   | 26.8 μs | **541.4 μs** | 1,470.5 μs | **Gig 快 2.7 倍** |  20.2x |
| **Reflect** (fmt/encoding)         | 22.8 μs | **397.1 μs** |   962.1 μs | **Gig 快 2.4 倍** |  17.4x |
| **Method** (Builder/Buffer/Regexp) | 17.1 μs | **486.3 μs** | 1,171.4 μs | **Gig 快 2.4 倍** |  28.4x |
| **Mixed** (函数 + 方法)            | 11.2 μs | **338.2 μs** |   823.2 μs | **Gig 快 2.4 倍** |  30.2x |

### 内存效率

| 工作负载        | Gig 分配次数/op | Yaegi 分配次数/op |      Gig 优势 |
| --------------- | --------------: | ----------------: | ------------: |
| Fibonacci(25)   |           **7** |         2,138,705 | 少 305,529 倍 |
| BubbleSort(100) |           **9** |             5,085 |     少 565 倍 |
| Sieve(1000)     |           **7** |                43 |       少 6 倍 |
| ExtCallMethod   |       **6,906** |            13,916 |     少 2.0 倍 |
| ExtCallMixed    |       **4,258** |             9,125 |     少 2.1 倍 |

### 并发压力测试

使用真实规则引擎工作负载（字符串处理 + 数学运算 + 条件逻辑 + stdlib 调用），3 轮取中位数，每轮持续 3 秒：

**Gig（32 核 AMD EPYC 9754，优化后 VMPool）：**

| 并发度     | 吞吐量              | 平均延迟    | 错误数 | 堆分配     | GC 次数 |
| ---------: | ------------------: | ----------: | -----: | ---------: | ------: |
|          1 |     **430K ops/s** |      2.3 μs |      0 |  ~1,130 MB |     163 |
|         10 |   **2,061K ops/s** |      4.9 μs |      0 |  ~5,275 MB |     687 |
|        100 |   **3,081K ops/s** |     32.5 μs |      0 |  ~7,990 MB |     474 |
|        500 |   **3,765K ops/s** |    132.8 μs |      0 |  ~9,818 MB |     308 |
|      1,000 |   **4,077K ops/s** |    245.3 μs |      0 | ~10,624 MB |     240 |
|      2,000 |   **4,426K ops/s** |    451.9 μs |      0 | ~11,550 MB |     156 |
|      5,000 |   **4,686K ops/s** |  1,067.0 μs |      0 | ~12,306 MB |      81 |
|     10,000 |   **4,776K ops/s** |  2,093.7 μs |      0 | ~12,597 MB |      55 |

**Native Go 基线（相同工作负载）：**

| 并发度     | 吞吐量           | 平均延迟    | 堆分配     | GC 次数 |
| ---------: | ---------------: | ----------: | ---------: | ------: |
|          1 |    6,423K ops/s |      0.2 μs |    ~904 MB |      89 |
|        100 |   30,148K ops/s |      3.3 μs |  ~4,123 MB |     308 |
|      1,000 |   31,287K ops/s |     32.0 μs |  ~4,296 MB |     223 |
|     10,000 |   28,375K ops/s |    352.6 μs |  ~4,206 MB |     153 |

**吞吐量比值（Native / Gig）：** 单核 14.9x → 10K 并发 5.9x

**关键发现**：
- **零错误**：10,000 个并发 goroutine，3 轮测试，0 错误
- **并发扩展性良好**：吞吐量从 1G 的 430K 增长到 10000G 的 4,776K（11.1 倍提升）
- **与 Native 差距收窄**：随并发增加，差距从 14.9x 降至 5.9x
- **GC 友好**：高并发时 GC 次数随并发增加反而降低（内存分配压力分散）

### 并发有状态压力测试

使用 `WithStatefulGlobals()` 模式，全局变量带 `*sync.Mutex` 保护的规则引擎工作负载：

**Gig Stateful（32 核 AMD EPYC 9754，WithStatefulGlobals + *sync.Mutex）：**

| 并发度     | 吞吐量            | 平均延迟      | 错误数 | 堆分配     | GC 次数 |
| ---------: | ----------------: | ------------: | -----: | ---------: | ------: |
|          1 |     **296K ops/s** |      3.4 μs |      0 |    ~960 MB |     217 |
|        100 |     **372K ops/s** |    268.7 μs |      0 |  ~1,210 MB |     155 |
|      1,000 |     **402K ops/s** |  2,489.7 μs |      0 |  ~1,332 MB |      36 |
|      5,000 |     **376K ops/s** | 13,309.8 μs |      0 |  ~1,386 MB |      11 |
|     10,000 |     **372K ops/s** | 26,861.2 μs |      0 |  ~1,544 MB |       6 |

**计数器正确性验证**：100G × 100 次操作 = 10,000 次计数器自增，最终值 = 10,000 ✓（零丢失更新）

**关键发现**：
- **有状态模式吞吐量约为无状态模式的 8-13%**：受 SharedGlobals RWMutex + *sync.Mutex 双重锁开销影响
- **零错误、零丢失更新**：mutex 保护下全局变量自增完全精确
- **高并发下吞吐量趋于饱和**：5K+ 并发时 ~375K ops/s，锁竞争成为瓶颈

### 分析

**Gig 在 3/5 项核心基准测试中优于 Yaegi**：

- **递归快 4.3 倍**（Fib25）—— O(1) 函数查找、帧池化，仅 7 次分配 vs 210 万次
- **外部调用快 2.4–2.7 倍** —— 1,162 个生成的 DirectCall 包装器消除了 92% 标准库函数和方法的 `reflect.Value.Call()`
- **闭包快 2.4 倍** —— 高效的闭包表示，通过共享 `*value.Value` 捕获变量
- **排序快 1.3 倍**（BubbleSort）—— 优化的切片操作和循环
- **紧凑循环**（ArithSum、Sieve）—— Yaegi 快 1.3-1.8 倍；Gig 的字节码解释开销在极短循环中更明显

**GopherLua vs Gig**：GopherLua 在纯数值循环上接近 Gig，但是：

- **GopherLua 需要手动注册函数** —— 每个 Go 函数都需要单独包装和注册；无法直接导入包
- **没有 Goroutine/Channel** —— Lua 有协程，但不是 Go 的 CSP 并发模型
- **没有结构体/接口/方法** —— Lua 使用表（table），不是 Go 的类型系统
- **不同的语法** —— 团队需要学习 Lua；Gig 使用熟悉的 Go 语法

关键优化：SSA 到字节码编译、32 字节 tagged-union 值、超级指令融合（17 种模式）、`intLocals []int64` 特化、`[]int64` 切片融合、DirectCall 代码生成、帧池化和内联缓存、**lock-free VMPool（sync.Pool，吞吐量提升 50%）**。

**为什么选择 Gig：**

|                       | Gig                             | Yaegi        | GopherLua  | Expr       |
| --------------------- | ------------------------------- | ------------ | ---------- | ---------- |
| **语言**              | Go                              | Go           | Lua        | 表达式 DSL |
| **完整 Go 语法**      | ✅                              | ✅           | ❌         | ❌         |
| **Goroutine/Channel** | ✅                              | ✅           | ❌         | ❌         |
| **Defer/Panic/Recover** | ✅                            | ✅           | ❌         | ❌         |
| **安全沙箱**          | ✅（禁止 unsafe/reflect/panic） | ❌           | ❌         | ✅         |
| **结构体/接口/方法**  | ✅                              | ✅           | ❌         | 有限       |
| **40+ 标准库包**      | ✅                              | ✅           | 需手动注册 | N/A        |
| **自定义 Go 包导入**  | ✅（代码生成）                  | ✅（符号表） | 需手动包装 | N/A        |
| **Context 取消**      | ✅                              | ❌           | ❌         | ❌         |
| **并发压力测试**      | ✅（10KG, 477万/秒, 0 错误）  | 未测试       | 未测试     | 未测试     |
| **可嵌入**            | ✅                              | ✅           | ✅         | ✅         |

**复现这些基准测试：**

```bash
cd benchmarks
# 单线程基准测试
go test -bench='^Benchmark(Gig|Yaegi|Lua|Native|Expr)' -benchmem -count=3 -timeout=30m -run='^$'
# 并发压力测试
go test -bench='BenchmarkStress' -benchmem -count=3 -timeout=10m -run='^$'
# 持续吞吐量测试
go test -run='TestStress_Gig_Sustained5s' -v -timeout=5m
# 并发有状态压力测试
go test -run='TestStatefulStress' -v -timeout=10m
# Go Native vs Gig 并发全局变量对比
go test -run='TestConcurrentGlobals_GoNative_vs_Gig' -v -timeout=2m
```

## 安全性

Gig 通过禁止某些导入来强制安全性：

- `unsafe` - 内存安全
- `reflect` - 类型安全
- `panic` 使用 - 受控执行

## 架构

Gig 使用多阶段编译流水线将 Go 源代码转换为高效字节码，然后由基于栈的虚拟机执行。

### 高层架构

```mermaid
flowchart TB
    subgraph Input["📥 输入"]
        SRC["Go 源代码"]
    end

    subgraph Frontend["🔍 前端"]
        PARSER["go/parser<br/>AST 生成"]
        TYPECHECK["go/types<br/>类型检查"]
        SSA["golang.org/x/tools/go/ssa<br/>SSA IR 生成"]
    end

    subgraph Compiler["⚙️ 编译器"]
        COMP["SSA → 字节码<br/>~100 操作码"]
        CONST["常量池"]
        TYPES["类型池"]
        FUNCS["函数注册表"]
    end

    subgraph Runtime["🚀 运行时"]
        VM["基于栈的虚拟机"]
        VALUE["Tagged-Union<br/>值系统"]
        EXT["外部包<br/>注册表"]
    end

    subgraph Output["📤 输出"]
        RESULT["结果<br/>(interface{})"]
    end

    SRC --> PARSER
    PARSER --> TYPECHECK
    TYPECHECK --> SSA
    SSA --> COMP
    COMP --> CONST
    COMP --> TYPES
    COMP --> FUNCS
    CONST --> VM
    TYPES --> VM
    FUNCS --> VM
    VM --> VALUE
    VM --> EXT
    EXT --> VALUE
    VALUE --> RESULT
```

### 详细组件架构

```mermaid
flowchart LR
    subgraph "用户代码"
        UC["gig.Build(source)"]
        UR["prog.Run(func, args...)"]
    end

    subgraph "gig.go [入口点]"
        BUILD["Build()"]
        SECURITY["安全检查<br/>(unsafe, reflect, panic)"]
        RUN["Run() / RunWithContext()"]
    end

    subgraph "前端流水线"
        PARSE["解析器<br/>go/parser"]
        TC["类型检查器<br/>go/types"]
        SSABUILD["SSA 构建器<br/>golang.org/x/tools/go/ssa"]
    end

    subgraph "compiler/ [编译器]"
        COMPILER["编译器"]
        SYMTAP["符号表"]
        BYTECODE["字节码生成器"]
        subgraph "输出"
            PROG["Program"]
            CFUNC["CompiledFunction[]"]
            CCONST["Constants[]"]
            CTYPES["Types[]"]
        end
    end

    subgraph "vm/ [虚拟机]"
        VM["VM"]
        STACK["栈<br/>(Value[])"]
        FRAMES["调用帧"]
        OPS["操作码处理器<br/>(~100 ops)"]
    end

    subgraph "value/ [值系统]"
        VAL["Value (tagged-union)"]
        PRIM["基本类型<br/>(int, float, string, bool)"]
        REF["反射回退<br/>(复杂类型)"]
    end

    subgraph "importer/ [包系统]"
        IMP["导入器<br/>(types.Importer)"]
        REG["包注册表"]
        EXTTYPE["ExternalPackage"]
    end

    subgraph "register/ [公开 API]"
        REGAPI["AddPackage()"]
        REGFUNC["NewFunction()"]
        REGVAR["NewVar()"]
        REGCONST["NewConst()"]
    end

    subgraph "packages/ [标准库]"
        STDLIB["40+ 包<br/>(fmt, strings, math, ...)"]
    end

    UC --> BUILD
    BUILD --> SECURITY
    SECURITY --> PARSE
    PARSE --> TC
    TC --> SSABUILD
    SSABUILD --> COMPILER
    COMPILER --> SYMTAP
    SYMTAP --> BYTECODE
    BYTECODE --> PROG
    PROG --> CFUNC
    PROG --> CCONST
    PROG --> CTYPES

    UR --> RUN
    RUN --> VM
    VM --> STACK
    VM --> FRAMES
    VM --> OPS
    OPS --> VAL
    VAL --> PRIM
    VAL --> REF

    TC --> IMP
    IMP --> REG
    REG --> EXTTYPE

    REGAPI --> REG
    REGFUNC --> EXTTYPE
    REGVAR --> EXTTYPE
    REGCONST --> EXTTYPE

    STDLIB --> REG
```

### 执行时数据流

```mermaid
sequenceDiagram
    participant User as 用户代码
    participant Gig as gig.go
    participant Frontend as 前端
    participant Comp as 编译器
    participant VM as 虚拟机
    participant Value as 值系统
    participant Ext as 外部包

    User->>Gig: Build(source)
    Gig->>Frontend: Parse(source)
    Frontend-->>Gig: AST
    Gig->>Frontend: TypeCheck(AST)
    Frontend->>Ext: 导入包
    Ext-->>Frontend: types.Package
    Frontend-->>Gig: 类型检查后的 AST
    Gig->>Frontend: BuildSSA(AST)
    Frontend-->>Gig: SSA IR
    Gig->>Comp: Compile(SSA)
    Comp-->>Gig: Program{Functions, Constants, Types}
    Gig-->>User: *Program

    User->>Gig: Run(funcName, args)
    Gig->>Value: FromInterface(args)
    Value-->>Gig: []Value
    Gig->>VM: Execute(funcName, args)
    VM->>VM: 取指-译码-执行 循环
    VM->>Value: 操作 (Add, Sub, 等)
    VM->>Ext: 外部函数调用
    Ext-->>VM: 结果
    VM-->>Gig: Value
    Gig->>Value: Interface()
    Value-->>Gig: interface{}
    Gig-->>User: result
```

### 编译流水线详情

```mermaid
flowchart TB
    subgraph Input
        SRC["Go 源码"]
    end

    subgraph "阶段 1: 解析"
        P1["词法/语法分析器<br/>go/parser"]
        AST["抽象语法树"]
    end

    subgraph "阶段 2: 类型检查"
        P2["类型检查器<br/>go/types"]
        TCINFO["types.Info<br/>(Types, Defs, Uses, Scopes)"]
        PKG["types.Package"]
    end

    subgraph "阶段 3: SSA 生成"
        P3["SSA 构建器<br/>golang.org/x/tools/go/ssa"]
        SSAFN["ssa.Function"]
        SSABLK["ssa.BasicBlock"]
        SSAINST["ssa.Instruction"]
    end

    subgraph "阶段 4: 字节码编译"
        P4["编译器"]
        SYM["符号表<br/>(Value → Local Slot)"]
        PHI["Phi 消除"]
        JMP["跳转修补"]
    end

    subgraph Output
        PROG["Program"]
        FN["CompiledFunction"]
        CODE["字节码<br/>(~100 opcodes)"]
        CONSTPOOL["常量池"]
        TYPEPOOL["类型池"]
    end

    SRC --> P1 --> AST
    AST --> P2 --> TCINFO --> PKG
    PKG --> P3 --> SSAFN --> SSABLK --> SSAINST
    SSAINST --> P4
    P4 --> SYM --> PHI --> JMP
    JMP --> PROG
    PROG --> FN --> CODE
    PROG --> CONSTPOOL
    PROG --> TYPEPOOL
```

### 虚拟机架构

```mermaid
flowchart TB
    subgraph VM["虚拟机"]
        subgraph State["执行状态"]
            STACK["栈<br/>Value[1024]"]
            SP["栈指针"]
            FRAMES["调用帧[64]"]
            FP["帧指针"]
            GLOBALS["全局变量"]
        end

        subgraph Frame["调用帧"]
            FN["函数"]
            IP["指令指针"]
            LOCALS["局部变量[]"]
            FREE["自由变量[]"]
            DEFER["延迟调用"]
        end

        subgraph Execution["执行循环"]
            FETCH["取操作码"]
            DECODE["解码操作数"]
            EXEC["执行"]
            CHECK["上下文检查<br/>(每 1024 条指令)"]
        end
    end

    subgraph Opcodes["操作码类别"]
        STACK_OP["栈操作<br/>(Push, Pop, Dup)"]
        ARITH["算术<br/>(Add, Sub, Mul, Div)"]
        CMP["比较<br/>(Eq, Lt, Gt)"]
        CTRL["控制流<br/>(Jump, Call, Return)"]
        CONTAINER["容器<br/>(Index, Slice, Map)"]
        FUNC["函数<br/>(Closure, CallExternal)"]
        BUILTIN["内置<br/>(Len, Append, Make)"]
    end

    STACK --> FETCH --> DECODE --> EXEC --> CHECK --> FETCH
    EXEC --> STACK_OP
    EXEC --> ARITH
    EXEC --> CMP
    EXEC --> CTRL
    EXEC --> CONTAINER
    EXEC --> FUNC
    EXEC --> BUILTIN

    FRAMES --> Frame
    Frame --> LOCALS --> STACK
```

### 值系统设计

```mermaid
flowchart LR
    subgraph Value["Value (16 字节 + obj)"]
        KIND["Kind (uint8)"]
        NUM["num (int64)"]
        NUM2["num2 (int64)"]
        STR["str (string)"]
        OBJ["obj (any)"]
    end

    subgraph Kinds["值类型"]
        PRIM["基本类型<br/>(零分配)"]
        COMP["复合类型<br/>(反射回退)"]
    end

    subgraph Primitives["基本类型快速路径"]
        BOOL["KindBool<br/>num: 0|1"]
        INT["KindInt<br/>num: int64"]
        UINT["KindUint<br/>num: uint64 bits"]
        FLOAT["KindFloat<br/>num: float64 bits"]
        STRV["KindString<br/>str: string"]
        CPLX["KindComplex<br/>num+num2: 实部+虚部"]
    end

    subgraph Composite["复合类型慢速路径"]
        SLICE["KindSlice/Array<br/>obj: reflect.Value"]
        MAP["KindMap<br/>obj: reflect.Value"]
        STRUCT["KindStruct<br/>obj: reflect.Value"]
        FUNC["KindFunc<br/>obj: *Closure"]
        IFACE["KindInterface<br/>obj: interface{}"]
        REFLECT["KindReflect<br/>obj: reflect.Value"]
    end

    KIND --> Kinds
    Kinds --> PRIM
    Kinds --> COMP
    PRIM --> Primitives
    COMP --> Composite
```

### 外部包集成

```mermaid
flowchart TB
    subgraph Registration["包注册"]
        CLI["gig CLI"]
        PKGS["pkgs.go<br/>(imports)"]
        GEN["gig gen"]
        GENERATED["packages/*.go<br/>(生成的)"]
    end

    subgraph Runtime["运行时集成"]
        REG["包注册表"]
        IMPORT["导入器<br/>(types.Importer)"]
        TYPECONV["类型转换器<br/>(reflect.Type → types.Type)"]
        METHOD["方法内省<br/>(addMethodsToNamed)"]
    end

    subgraph Execution["VM 执行"]
        CALL["OpCallExternal"]
        CACHE["内联缓存"]
        DIRECT["DirectCall<br/>(快速路径)"]
        REFLECT["reflect.Call<br/>(慢速路径)"]
        METHODCALL["方法分发<br/>(MethodByName)"]
    end

    CLI --> PKGS --> GEN --> GENERATED
    GENERATED --> REG
    REG --> IMPORT --> TYPECONV --> METHOD

    IMPORT --> CALL
    CALL --> CACHE
    CACHE --> DIRECT
    CACHE --> REFLECT
    CALL --> METHODCALL
```

---

### 组件概览

| 组件         | 包                | 用途                                             |
| ------------ | ----------------- | ------------------------------------------------ |
| **入口点**   | `gig.go`          | 公开 API：`Build()`、`Run()`、`RunWithContext()` |
| **编译器**   | `compiler/`       | SSA 到字节码编译（~100 操作码）                  |
| **虚拟机**   | `vm/`             | 基于栈的字节码执行                               |
| **值系统**   | `model/value/`    | Tagged-union 值，基本类型零分配                  |
| **字节码**   | `model/bytecode/` | CompiledProgram、OpCode 定义                     |
| **外部类型** | `model/external/` | ExternalFuncInfo、ExternalObject 等共享类型       |
| **导入器**   | `importer/`       | 外部包注册、类型解析                             |
| **运行器**   | `runner/`         | VM 池、init 执行、有状态/无状态模式              |
| **标准库包** | `stdlib/packages/`| 40+ 预注册标准库包                               |
| **CLI**      | `cmd/gig`         | 代码生成工具                                     |

### 关键设计决策

1. **基于 SSA 的编译**：使用 Go 官方 SSA 库，正确处理复杂控制流、闭包和方法调用。

2. **Tagged-Union 值**：基本类型操作避免反射开销，将值存储在联合体中的原生 Go 类型中。

3. **内联缓存**：外部函数调用缓存已解析的函数信息，实现快速分发。

4. **上下文集成**：虚拟机每 1024 条指令检查一次上下文取消，实现响应式超时处理。

5. **默认安全**：在解释代码中禁止 `unsafe`、`reflect` 和 `panic`，实现受控执行。

## 更新日志

### v0.4.0 - VMPool Performance Optimization

**优化**：VMPool 从 mutex 实现改为 lock-free `sync.Pool`，显著提升并发性能。

- **吞吐量提升 50%**：从 1.15M ops/s → 1.72M ops/s（并发 20 goroutines）
- **内存减少 99.7%**：从 3 GB → 9 MB（高并发场景）
- **无锁竞争**：`sync.Pool` 使用 per-P 本地缓存，完全消除 mutex 竞争
- **无内存泄漏**：长时间压力测试（10 分钟）验证，堆增长 < 5 MB
- **基准测试验证**：
  - VMPool 并发：146.6 ns/op
  - VMPool 串行：1428 ns/op
  - 10 并发 goroutines：355.4 μs/op

**技术细节**：
- 移除 `sync.Mutex` + `[]*vm` 实现
- 采用 `sync.Pool` 的 lock-free 设计
- GC 友好：未使用的 VM 自动回收
- 特别适合高并发场景（20+ goroutines）

### v0.3.0 - Fmt Sanitization Integration

**改进**：`fmt` 包的参数清理逻辑现已完全集成到生成的代码中。

此前，`fmt.Sprintf` 等函数在打印 Gig 结构体时会输出冗长的内部表示。此问题通过在生成的 `fmt.go` 中嵌入参数清理辅助函数（`sanitizeArgForFmt`、`sprintfWithTypeAwareness` 等）解决：

- **`gentool/directcall.go`**：新增 `fmtSanitizeHelperCode()` 函数，返回参数清理辅助代码的 Go 源码字符串
- **`gentool/generator.go`**：在生成 `fmt` 包时，将辅助代码嵌入到生成的文件中
- **删除 `fmt_sanitize.go`**：不再需要单独维护的支持文件，所有代码均由 `gig gen` 生成

现在 `fmt.go` 是一个完全生成的、自包含的文件，包含 DirectCall 包装器和参数清理辅助函数。

### v0.2.0 - 外部类型方法支持

**修复**：现在完全支持外部（已注册）类型上的方法。

此前，调用外部类型的方法如 `gjson.Get(json, path).String()` 会因类型检查错误而失败，因为方法未在 `types.Named` 类型上注册。此问题已在三个层面修复：

- **`importer/importer.go`**：添加了 `addMethodsToNamed()` —— 在将 `reflect.Type` 转换为 `types.Named` 时，枚举所有导出方法（值接收者和指针接收者）并通过 `named.AddMethod()` 添加。这使 Go 类型检查器能够解析外部类型上的方法调用。

- **`compiler/compiler.go`**：添加了 `ExternalMethodInfo` 并更新 `compileExternalStaticCall` 以检测方法调用（`sig.Recv() != nil`）并使用方法分发元数据发出，而不是查找静态函数对象。

- **`vm/vm.go`**：添加了 `callExternalMethod()`，通过 `reflect.Value.MethodByName()` 分发外部类型上的方法调用，处理可变参数、指针接收者和多返回值。

## 许可证

MIT License
