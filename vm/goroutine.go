// goroutine.go provides GoroutineTracker and child VM construction (newChildVM, newDeferVM).
package vm

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"git.woa.com/youngjin/gig/model/value"
)

// GoroutineTracker tracks active interpreter goroutines for a single program.
// It replaces the old process-wide activeGoroutines counter, making concurrent
// multi-program usage safe.
type GoroutineTracker struct {
	active        int64
	maxGoroutines int64
}

// NewGoroutineTracker creates a new goroutine tracker with the default limit.
func NewGoroutineTracker() *GoroutineTracker {
	return &GoroutineTracker{
		maxGoroutines: defaultMaxGoroutines,
	}
}

// SetMaxGoroutines sets the maximum number of concurrent goroutines.
func (t *GoroutineTracker) SetMaxGoroutines(n int) {
	atomic.StoreInt64(&t.maxGoroutines, int64(n))
}

// Start launches a goroutine and tracks it.
// Returns an error if the goroutine limit would be exceeded.
func (t *GoroutineTracker) Start(fn func()) error {
	max := atomic.LoadInt64(&t.maxGoroutines)
	if max > 0 && atomic.LoadInt64(&t.active) >= max {
		return fmt.Errorf("gig: goroutine limit (%d) exceeded", max)
	}
	atomic.AddInt64(&t.active, 1)
	go func() {
		defer atomic.AddInt64(&t.active, -1)
		fn()
	}()
	return nil
}

// Wait blocks until all tracked goroutines have completed.
// Uses exponential backoff to avoid busy waiting.
func (t *GoroutineTracker) Wait() {
	backoff := time.Microsecond
	for atomic.LoadInt64(&t.active) > 0 {
		time.Sleep(backoff)
		if backoff < 10*time.Millisecond {
			backoff *= 2
		}
	}
}

// WaitContext blocks until all tracked goroutines complete or the context is cancelled.
func (t *GoroutineTracker) WaitContext(ctx context.Context) error {
	backoff := time.Microsecond
	for atomic.LoadInt64(&t.active) > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		time.Sleep(backoff)
		if backoff < 10*time.Millisecond {
			backoff *= 2
		}
	}
	return nil
}

// newChildVM creates a child VM for goroutine execution.
// The child VM shares the globals pointer and external call cache with the parent.
func (v *vm) newChildVM() *vm {
	child := &vm{
		program:        v.program,
		stack:          make([]value.Value, initialStackSize),
		sp:             0,
		frames:         make([]*Frame, initialFrameDepth),
		fp:             0,
		globals:        nil, // Not used when globalsPtr is set
		globalsPtr:     v.globalsPtr,
		ctx:            v.ctx,
		extCallCache:   v.extCallCache,
		initialGlobals: v.initialGlobals,
		goroutines:     v.goroutines,
	}
	if child.globalsPtr == nil {
		child.globalsPtr = &v.globals
	}
	return child
}

// newDeferVM creates a lightweight child VM for deferred function execution.
// Uses getGlobals() to correctly handle both stateful and goroutine modes.
// This consolidates the 3 inline child VM construction sites for defers.
func (v *vm) newDeferVM() *vm {
	return &vm{
		program:      v.program,
		stack:        make([]value.Value, deferVMStackSize),
		sp:           0,
		frames:       make([]*Frame, initialFrameDepth),
		fp:           0,
		globals:      v.getGlobals(),
		globalsPtr:   v.globalsPtr,
		ctx:          v.ctx,
		extCallCache: v.extCallCache,
	}
}
