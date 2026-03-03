# Context 取消支持

Gig 支持 Go 的 `context.Context` 机制，用于在脚本执行期间进行超时和取消控制。这允许您安全地中断长时间运行的脚本、阻塞操作和无限循环。

## 概述

Context 取消系统提供以下功能：

- **协作式取消**: VM 在关键点检查取消信号
- **阻塞操作支持**: Channel 发送/接收/select 可以被中断
- **外部调用边界检查**: 长时间运行的 Go 函数调用可以被取消
- **低开销**: 对正常执行的性能影响最小

## API 使用

### 基本超时

```go
package main

import (
    "context"
    "time"
    "git.woa.com/youngjin/gig"
    _ "git.woa.com/youngjin/gig/stdlib/packages"
)

func main() {
    source := `
package main

func InfiniteLoop() int {
    sum := 0
    for {
        sum++
    }
    return sum
}
`
    prog, _ := gig.Build(source)
    
    // 使用 1 秒超时运行
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    _, err := prog.RunWithContext(ctx, "InfiniteLoop")
    if err != nil {
        // 将收到 context.DeadlineExceeded 错误
        println("执行超时:", err.Error())
    }
}
```

### 手动取消

```go
package main

import (
    "context"
    "git.woa.com/youngjin/gig"
    _ "git.woa.com/youngjin/gig/stdlib/packages"
)

func main() {
    source := `
package main

func LongRunning() int {
    sum := 0
    for i := 0; i < 100000000; i++ {
        sum += i
    }
    return sum
}
`
    prog, _ := gig.Build(source)
    
    ctx, cancel := context.WithCancel(context.Background())
    
    // 一段时间后取消
    go func() {
        time.Sleep(100 * time.Millisecond)
        cancel()
    }()
    
    _, err := prog.RunWithContext(ctx, "LongRunning")
    if err == context.Canceled {
        println("执行被取消")
    }
}
```

### 阻塞 Channel 操作

```go
package main

import (
    "context"
    "time"
    "git.woa.com/youngjin/gig"
    _ "git.woa.com/youngjin/gig/stdlib/packages"
)

func main() {
    source := `
package main

func BlockOnChannel() int {
    ch := make(chan int)
    // 除非被取消，否则会永远阻塞
    val := <-ch
    return val
}
`
    prog, _ := gig.Build(source)
    
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    _, err := prog.RunWithContext(ctx, "BlockOnChannel")
    // 500ms 后执行将被取消
}
```

## 取消检查点

VM 在以下位置检查取消信号：

### 1. 循环回边
每 128 次向后跳转（循环），VM 检查 context 是否被取消。这确保无限循环可以被中断。

### 2. 阻塞 Channel 操作
- **发送**: `ch <- value` 操作在等待时可以被取消
- **接收**: `<-ch` 操作在等待时可以被取消
- **Select**: 阻塞的 select 语句可以被取消

### 3. 外部函数调用
调用外部 Go 函数（通过反射或直接调用）后，VM 检查调用期间 context 是否被取消。

### 4. 顺序代码
每 1024 条指令，VM 执行周期性检查，确保即使在没有循环的长时间顺序代码中也能响应取消。

## 性能影响

Context 取消系统的性能开销很小：

| 工作负载类型 | 开销 |
|-------------|------|
| 纯计算 (Fibonacci) | ~1% |
| 算术循环 | ~3% |
| 外部调用 | ~6-7% |

开销来自：
- 主执行循环中的周期性计数器检查
- `select` 语句检查 context.Done() 的开销
- 外部函数的后调用检查

## 错误类型

当执行被取消时，Gig 返回标准的 Go context 错误：

- `context.Canceled` - 执行被 `cancel()` 函数取消
- `context.DeadlineExceeded` - 执行超过超时时间

您可以使用标准 Go 错误处理检查这些错误：

```go
result, err := prog.RunWithContext(ctx, "Function")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // 处理超时
    } else if errors.Is(err, context.Canceled) {
        // 处理取消
    }
}
```

## 实现细节

### Value 包扩展

`value` 包提供 context 感知的 channel 操作：

- `SendContext(ctx, val)` - 带取消支持的发送
- `RecvContext(ctx)` - 带取消支持的接收

这些使用 Go 的 `reflect.Select` 同时等待 channel 操作和 context 的 Done channel。

### VM 集成

VM 在多个层面集成取消检查：

1. **主循环** (`vm/run.go`): 周期性指令计数器检查
2. **操作码分发** (`vm/ops_dispatch.go`): Channel 操作取消
3. **外部调用** (`vm/call.go`): 调用后取消检查
4. **Goroutine 跟踪** (`vm/goroutine.go`): 基于 WaitGroup 的跟踪，支持 context

### 默认超时

`Run()` 方法使用默认 10 秒超时，防止无限循环挂起您的应用程序：

```go
const DefaultTimeout = 10 * time.Second
```

使用 `RunWithContext()` 获取自定义超时值或取消控制。

## 最佳实践

1. **始终使用超时** 处理不受信任或可能无限的脚本
2. **优雅地处理取消错误** 在您的应用程序中
3. **使用 `RunWithContext()`** 替代 `Run()` 用于生产代码
4. **设置适当的超时** 基于预期执行时间

## 示例：Web 服务器集成

```go
package main

import (
    "context"
    "net/http"
    "time"
    "git.woa.com/youngjin/gig"
    _ "git.woa.com/youngjin/gig/stdlib/packages"
)

var prog *gig.Program

func handler(w http.ResponseWriter, r *http.Request) {
    // 使用请求 context 进行取消
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    result, err := prog.RunWithContext(ctx, "ProcessRequest", r.URL.Path)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Write([]byte(result.(string)))
}

func main() {
    source := `
package main

func ProcessRequest(path string) string {
    // 处理请求
    return "Processed: " + path
}
`
    prog, _ = gig.Build(source)
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```
