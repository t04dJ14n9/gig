package thirdparty

import (
	"context"
	"sync"
	"time"
)

// ============================================================================
// TIME COMPLEX OPERATIONS
// ============================================================================

// TimeDurationParse tests ParseDuration.
func TimeDurationParse() int {
	d, _ := time.ParseDuration("1h30m")
	return int(d.Minutes())
}

// TimeUnixNano tests UnixNano.
func TimeUnixNano() int64 {
	t := time.Date(2024, 1, 1, 0, 0, 0, 123456789, time.UTC)
	return t.UnixNano()
}

// ============================================================================
// CONTEXT WITH CANCELLATION AND VALUES
// ============================================================================

// ContextWithCancelAndValue tests combined cancel and value.
func ContextWithCancelAndValue() int {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "key", "value")

	cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.Canceled {
			return 1
		}
	default:
	}
	return 0
}

// ContextNestedValues tests nested context values.
func ContextNestedValues() int {
	ctx := context.Background()
	ctx1 := context.WithValue(ctx, "level", 1)
	ctx2 := context.WithValue(ctx1, "level", 2)
	ctx3 := context.WithValue(ctx2, "level", 3)

	if ctx.Value("level") != nil {
		return 0
	}
	if ctx1.Value("level") != 1 {
		return 0
	}
	if ctx2.Value("level") != 2 {
		return 0
	}
	if ctx3.Value("level") != 3 {
		return 0
	}
	return 1
}

// ============================================================================
// SYNC WITH POOL
// ============================================================================

// SyncPool tests sync.Pool.
func SyncPool() int {
	var pool sync.Pool
	pool.New = func() interface{} {
		return make([]byte, 1024)
	}

	item := pool.Get()
	if item != nil {
		pool.Put(item)
	}
	return 1
}
