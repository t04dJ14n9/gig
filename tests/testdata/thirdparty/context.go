package thirdparty

import (
	"context"
	"time"
)

// ContextBackground tests context.Background.
func ContextBackground() int {
	ctx := context.Background()
	if ctx != nil {
		return 1
	}
	return 0
}

// ContextTODO tests context.TODO.
func ContextTODO() int {
	ctx := context.TODO()
	if ctx != nil {
		return 1
	}
	return 0
}

// ContextWithValue tests context.WithValue.
func ContextWithValue() int {
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "key", "value")
	v := ctx2.Value("key")
	if v == "value" {
		return 1
	}
	return 0
}

// ContextWithCancel tests context.WithCancel.
func ContextWithCancel() int {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	select {
	case <-ctx.Done():
		return 1
	default:
		return 0
	}
}

// ContextWithTimeout tests context.WithTimeout.
func ContextWithTimeout() int {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return 1
	}
	return 0
}

// ContextWithCancelParent tests that parent cancellation propagates to child.
func ContextWithCancelParent() int {
	parent, parentCancel := context.WithCancel(context.Background())
	child, _ := context.WithCancel(parent)
	parentCancel()
	select {
	case <-child.Done():
		return 1
	default:
		return 0
	}
}
