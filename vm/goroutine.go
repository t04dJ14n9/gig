package vm

import (
	"context"
	"sync"

	"github.com/t04dJ14n9/gig/value"
)

// newChildVM creates a child VM for goroutine execution.
// The child VM shares the globals pointer and external call cache with the parent.
func (vm *VM) newChildVM() *VM {
	child := &VM{
		program:      vm.program,
		stack:        make([]value.Value, 1024),
		sp:           0,
		frames:       make([]*Frame, 64),
		fp:           0,
		globals:      nil, // Not used when globalsPtr is set
		globalsPtr:   vm.globalsPtr,
		ctx:          vm.ctx,
		extCallCache: vm.extCallCache, // Share cache (thread-safe via shared mutex)
	}
	// If parent doesn't have a globalsPtr yet, create one for sharing
	if child.globalsPtr == nil {
		child.globalsPtr = &vm.globals
	}
	return child
}

// goroutineTracker provides efficient goroutine lifecycle tracking using sync.WaitGroup.
type goroutineTracker struct {
	wg sync.WaitGroup
}

// Global tracker instance
var globalGoroutineTracker = &goroutineTracker{}

// StartGoroutine starts a new goroutine and tracks it.
// Used for the "go" statement implementation.
func StartGoroutine(fn func()) {
	globalGoroutineTracker.wg.Add(1)
	go func() {
		defer globalGoroutineTracker.wg.Done()
		fn()
	}()
}

// WaitGoroutines waits for all tracked goroutines to complete.
// This blocks until all goroutines finish or the context is cancelled.
func WaitGoroutines() {
	globalGoroutineTracker.wg.Wait()
}

// WaitGoroutinesContext waits for all tracked goroutines to complete with context cancellation.
// Returns ctx.Err() if the context is cancelled before all goroutines complete.
func WaitGoroutinesContext(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		globalGoroutineTracker.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Global VM registry for concurrent execution.
var (
	vmRegistryMutex sync.Mutex
	vmRegistry      = make(map[int64]*VM)
	vmIDCounter     int64
)

// RegisterVM registers a VM for later use in concurrent execution.
// Returns a unique ID for the VM.
func RegisterVM(vm *VM) int64 {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	vmIDCounter++
	vmRegistry[vmIDCounter] = vm
	return vmIDCounter
}

// UnregisterVM removes a VM from the registry.
func UnregisterVM(id int64) {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	delete(vmRegistry, id)
}
