package divergence_hunt168

import (
	"context"
	"fmt"
	"time"
)

// ============================================================================
// Round 168: Context usage patterns
// ============================================================================

// BackgroundContext tests background context
func BackgroundContext() string {
	ctx := context.Background()
	return fmt.Sprintf("ctx=%v", ctx != nil)
}

// TODOContext tests TODO context
func TODOContext() string {
	ctx := context.TODO()
	return fmt.Sprintf("ctx=%v", ctx != nil)
}

// WithValue tests context with value
func WithValue() string {
	ctx := context.WithValue(context.Background(), "key", "value")
	val := ctx.Value("key").(string)
	return fmt.Sprintf("val=%s", val)
}

// WithCancel tests cancelable context
func WithCancel() string {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan string, 1)
	go func() {
		<-ctx.Done()
		done <- "cancelled"
	}()
	cancel()
	result := <-done
	return fmt.Sprintf("result=%s", result)
}

// WithDeadline tests deadline context
func WithDeadline() string {
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	dl, ok := ctx.Deadline()
	return fmt.Sprintf("has_deadline=%t", ok && !dl.IsZero())
}

// WithTimeout tests timeout context
func WithTimeout() string {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	time.Sleep(100 * time.Millisecond)
	err := ctx.Err()
	return fmt.Sprintf("err=%v", err != nil)
}

// NestedContext tests nested contexts with values
func NestedContext() string {
	ctx1 := context.WithValue(context.Background(), "a", "1")
	ctx2 := context.WithValue(ctx1, "b", "2")
	ctx3 := context.WithValue(ctx2, "c", "3")
	a := ctx3.Value("a").(string)
	b := ctx3.Value("b").(string)
	c := ctx3.Value("c").(string)
	return fmt.Sprintf("a=%s,b=%s,c=%s", a, b, c)
}

// ContextValueOverride tests value override in nested contexts
func ContextValueOverride() string {
	ctx1 := context.WithValue(context.Background(), "key", "original")
	ctx2 := context.WithValue(ctx1, "key", "overridden")
	val1 := ctx1.Value("key").(string)
	val2 := ctx2.Value("key").(string)
	return fmt.Sprintf("ctx1=%s,ctx2=%s", val1, val2)
}

// ContextPropagation tests context cancellation propagation
func ContextPropagation() string {
	parent, cancel := context.WithCancel(context.Background())
	child1, _ := context.WithCancel(parent)
	child2, _ := context.WithCancel(parent)
	cancel()
	result := ""
	if parent.Err() != nil {
		result += "parent-cancelled "
	}
	if child1.Err() != nil {
		result += "child1-cancelled "
	}
	if child2.Err() != nil {
		result += "child2-cancelled"
	}
	return fmt.Sprintf("result=%s", result)
}

// SelectWithContext tests select with context
func SelectWithContext() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ch := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "result"
	}()
	select {
	case <-ch:
		return "received"
	case <-ctx.Done():
		return "timeout"
	}
}

// ContextInStruct tests context stored in struct
func ContextInStruct() string {
	type Worker struct {
		ctx context.Context
	}
	ctx, cancel := context.WithCancel(context.Background())
	w := Worker{ctx: ctx}
	cancel()
	return fmt.Sprintf("cancelled=%v", w.ctx.Err() != nil)
}

// MultipleValues tests multiple values in context
func MultipleValues() string {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "12345")
	ctx = context.WithValue(ctx, "user_id", "user123")
	ctx = context.WithValue(ctx, "trace_id", "trace456")
	reqID := ctx.Value("request_id").(string)
	userID := ctx.Value("user_id").(string)
	return fmt.Sprintf("req=%s,user=%s", reqID, userID)
}
