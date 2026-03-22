# Context Cancellation Support

Gig supports Go's `context.Context` for timeout and cancellation control during script execution. This allows you to safely interrupt long-running scripts, blocking operations, and infinite loops.

## Overview

The context cancellation system provides:

- **Cooperative cancellation**: The VM checks for cancellation at strategic points
- **Blocking operation support**: Channel send/recv/select can be interrupted
- **External call boundary checks**: Long-running Go function calls can be cancelled
- **Low overhead**: Minimal performance impact on normal execution

## API Usage

### Basic Timeout

```go
package main

import (
    "context"
    "time"
    "github.com/t04dJ14n9/gig"
    _ "github.com/t04dJ14n9/gig/stdlib/packages"
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
    
    // Run with 1 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    _, err := prog.RunWithContext(ctx, "InfiniteLoop")
    if err != nil {
        // Will receive context.DeadlineExceeded error
        println("Execution timed out:", err.Error())
    }
}
```

### Manual Cancellation

```go
package main

import (
    "context"
    "github.com/t04dJ14n9/gig"
    _ "github.com/t04dJ14n9/gig/stdlib/packages"
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
    
    // Cancel after some time
    go func() {
        time.Sleep(100 * time.Millisecond)
        cancel()
    }()
    
    _, err := prog.RunWithContext(ctx, "LongRunning")
    if err == context.Canceled {
        println("Execution was cancelled")
    }
}
```

### Blocking Channel Operations

```go
package main

import (
    "context"
    "time"
    "github.com/t04dJ14n9/gig"
    _ "github.com/t04dJ14n9/gig/stdlib/packages"
)

func main() {
    source := `
package main

func BlockOnChannel() int {
    ch := make(chan int)
    // This will block forever unless cancelled
    val := <-ch
    return val
}
`
    prog, _ := gig.Build(source)
    
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    _, err := prog.RunWithContext(ctx, "BlockOnChannel")
    // Execution will be cancelled after 500ms
}
```

## Cancellation Points

The VM checks for cancellation at the following points:

### 1. Loop Back-Edges
Every 128 backward jumps (loops), the VM checks if the context is cancelled. This ensures infinite loops can be interrupted.

### 2. Blocking Channel Operations
- **Send**: `ch <- value` operations can be cancelled while waiting
- **Receive**: `<-ch` operations can be cancelled while waiting  
- **Select**: Blocking select statements can be cancelled

### 3. External Function Calls
After calling external Go functions (via reflection or direct call), the VM checks if the context was cancelled during the call.

### 4. Sequential Code
Every 1024 instructions, the VM performs a periodic check to ensure cancellation is responsive even in long-running sequential code without loops.

## Performance Impact

The context cancellation system has minimal performance overhead:

| Workload Type | Overhead |
|--------------|----------|
| Pure computation (Fibonacci) | ~1% |
| Arithmetic loops | ~3% |
| External calls | ~6-7% |

The overhead comes from:
- Periodic counter checks in the main execution loop
- `select` statement overhead for context.Done() checks
- Post-call checks for external functions

## Error Types

When execution is cancelled, Gig returns standard Go context errors:

- `context.Canceled` - Execution was cancelled via `cancel()` function
- `context.DeadlineExceeded` - Execution exceeded the timeout duration

You can check for these errors using standard Go error handling:

```go
result, err := prog.RunWithContext(ctx, "Function")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        // Handle timeout
    } else if errors.Is(err, context.Canceled) {
        // Handle cancellation
    }
}
```

## Implementation Details

### Value Package Extensions

The `value` package provides context-aware channel operations:

- `SendContext(ctx, val)` - Send with cancellation support
- `RecvContext(ctx)` - Receive with cancellation support

These use Go's `reflect.Select` to wait on both the channel operation and the context's Done channel.

### VM Integration

The VM integrates cancellation checks at multiple levels:

1. **Main loop** (`vm/run.go`): Periodic instruction counter checks
2. **Opcode dispatch** (`vm/ops_dispatch.go`): Channel operation cancellation
3. **External calls** (`vm/call.go`): Post-call cancellation checks
4. **Goroutine tracking** (`vm/goroutine.go`): WaitGroup-based tracking with context support

### Default Timeout

The `Run()` method uses a default timeout of 10 seconds to prevent infinite loops from hanging your application:

```go
const DefaultTimeout = 10 * time.Second
```

Use `RunWithContext()` for custom timeout values or cancellation control.

## Best Practices

1. **Always use timeouts** for untrusted or potentially infinite scripts
2. **Handle cancellation errors** gracefully in your application
3. **Use `RunWithContext()`** instead of `Run()` for production code
4. **Set appropriate timeouts** based on expected execution time

## Example: Web Server Integration

```go
package main

import (
    "context"
    "net/http"
    "time"
    "github.com/t04dJ14n9/gig"
    _ "github.com/t04dJ14n9/gig/stdlib/packages"
)

var prog *gig.Program

func handler(w http.ResponseWriter, r *http.Request) {
    // Use request context for cancellation
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
    // Process the request
    return "Processed: " + path
}
`
    prog, _ = gig.Build(source)
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```
