package vm

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/t04dJ14n9/gig/value"
)

// GoroutineTracker tracks active interpreter goroutines for a single program.
// It replaces the old process-wide activeGoroutines counter, making concurrent
// multi-program usage safe.
type GoroutineTracker struct {
	active int64
}

// NewGoroutineTracker creates a new goroutine tracker.
func NewGoroutineTracker() *GoroutineTracker {
	return &GoroutineTracker{}
}

// Start launches a goroutine and tracks it.
func (t *GoroutineTracker) Start(fn func()) {
	atomic.AddInt64(&t.active, 1)
	go func() {
		defer atomic.AddInt64(&t.active, -1)
		fn()
	}()
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
		stack:          make([]value.Value, 1024),
		sp:             0,
		frames:         make([]*Frame, 64),
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

// activeGoroutines tracks the number of active goroutines using atomic operations.
var activeGoroutines int64

// StartGoroutine starts a new goroutine and tracks it.
// Used for the "go" statement implementation.
func StartGoroutine(fn func()) {
	atomic.AddInt64(&activeGoroutines, 1)
	go func() {
		defer atomic.AddInt64(&activeGoroutines, -1)
		fn()
	}()
}

// WaitGoroutines waits for all tracked goroutines to complete.
// Uses exponential backoff to avoid busy waiting.
func WaitGoroutines() {
	backoff := time.Microsecond
	for atomic.LoadInt64(&activeGoroutines) > 0 {
		time.Sleep(backoff)
		// Cap backoff at 10ms to avoid waiting too long
		if backoff < 10*time.Millisecond {
			backoff *= 2
		}
	}
}

// WaitGoroutinesContext waits for all tracked goroutines to complete with context cancellation.
// Returns ctx.Err() if the context is cancelled before all goroutines complete.
func WaitGoroutinesContext(ctx context.Context) error {
	backoff := time.Microsecond
	for atomic.LoadInt64(&activeGoroutines) > 0 {
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

// Global VM registry for concurrent execution.
var (
	vmRegistryMutex sync.Mutex
	vmRegistry      = make(map[int64]*vm)
	vmIDCounter     int64
)

// RegisterVM registers a VM for later use in concurrent execution.
// Returns a unique ID for the VM.
func RegisterVM(v *vm) int64 {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	vmIDCounter++
	vmRegistry[vmIDCounter] = v
	return vmIDCounter
}

// UnregisterVM removes a VM from the registry.
func UnregisterVM(id int64) {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	delete(vmRegistry, id)
}
