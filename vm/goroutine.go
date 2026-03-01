package vm

import (
	"sync"
	"sync/atomic"

	"github.com/t04dJ14n9/gig/value"
)

// newChildVM creates a child VM for goroutine execution.
// The child VM shares the globals pointer with the parent for communication.
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
		extCallCache: vm.extCallCache, // Share cache (read-mostly, safe for goroutines)
	}
	// If parent doesn't have a globalsPtr yet, create one for sharing
	if child.globalsPtr == nil {
		child.globalsPtr = &vm.globals
	}
	return child
}

// Goroutine tracking for concurrent execution.
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
// Uses busy waiting - could be improved with sync.WaitGroup.
func WaitGoroutines() {
	for atomic.LoadInt64(&activeGoroutines) > 0 {
		// Busy wait - could use a WaitGroup instead
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
