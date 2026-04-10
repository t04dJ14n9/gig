// shared_globals.go provides a thread-safe wrapper around the globals slice
// for stateful mode concurrent execution.
package vm

import (
	"reflect"
	"sync"

	"git.woa.com/youngjin/gig/model/value"
)

// SharedGlobals wraps a globals slice with a sync.RWMutex for concurrent access.
// Used in stateful mode where multiple Run() calls and goroutines share the same
// global variable state.
//
// Locking strategy:
//   - Get / GetPtr: RLock (multiple concurrent readers allowed)
//   - Set: Lock (exclusive writer)
//   - Slice: returns the raw slice for bulk operations (caller must hold lock)
type SharedGlobals struct {
	mu      sync.RWMutex
	globals []value.Value
}

// NewSharedGlobals creates a SharedGlobals with the given initial values.
func NewSharedGlobals(initial []value.Value, size int) *SharedGlobals {
	globals := make([]value.Value, size)
	if len(initial) == size {
		copy(globals, initial)
	}
	return &SharedGlobals{
		globals: globals,
	}
}

// InitExternalVars applies external variable values to the shared globals.
// Should be called after construction to initialize external package variables
// (e.g., &time.UTC) that aren't captured in the init() snapshot.
func (sg *SharedGlobals) InitExternalVars(extVarValues map[int]any) {
	for idx, ptr := range extVarValues {
		if idx < len(sg.globals) {
			sg.globals[idx] = value.FromInterface(ptr)
		}
	}
}

// InitZeroValues applies zero-valued struct globals from the compiler.
// SSA may store nil constants for zero-valued structs (e.g., sync.Mutex{}),
// so we replace nil/invalid globals with their proper zero reflect.Value.
func (sg *SharedGlobals) InitZeroValues(zeroValues map[int]reflect.Value) {
	for idx, zeroRV := range zeroValues {
		if idx < len(sg.globals) {
			g := sg.globals[idx]
			if !g.IsValid() || g.IsNil() {
				sg.globals[idx] = value.MakeFromReflect(zeroRV)
			}
		}
	}
}

// Get reads the global at idx under a read lock.
func (sg *SharedGlobals) Get(idx int) value.Value {
	sg.mu.RLock()
	v := sg.globals[idx]
	sg.mu.RUnlock()
	return v
}

// GlobalRef is a reference to a shared global variable slot.
// Instead of exposing a raw pointer to globals[idx] (which would bypass the lock),
// GlobalRef defers the actual read/write to locked methods on SharedGlobals.
// Used as the value pushed by OpGlobal in shared mode.
type GlobalRef struct {
	sg  *SharedGlobals
	idx int
}

// Load reads the global value under a read lock.
func (r *GlobalRef) Load() value.Value {
	r.sg.mu.RLock()
	v := r.sg.globals[r.idx]
	r.sg.mu.RUnlock()
	return v
}

// Store writes a value to the global under a write lock.
func (r *GlobalRef) Store(val value.Value) {
	r.sg.mu.Lock()
	r.sg.globals[r.idx] = val
	r.sg.mu.Unlock()
}

// Set writes a value to the global at idx under an exclusive write lock.
func (sg *SharedGlobals) Set(idx int, val value.Value) {
	sg.mu.Lock()
	sg.globals[idx] = val
	sg.mu.Unlock()
}

// Len returns the number of globals.
func (sg *SharedGlobals) Len() int {
	return len(sg.globals)
}

// Globals returns the raw globals slice.
// Used for bulk operations like initialization and Reset.
// Caller is responsible for synchronization.
func (sg *SharedGlobals) Globals() []value.Value {
	return sg.globals
}

// RLock acquires a read lock on the globals.
func (sg *SharedGlobals) RLock() {
	sg.mu.RLock()
}

// RUnlock releases the read lock.
func (sg *SharedGlobals) RUnlock() {
	sg.mu.RUnlock()
}

// Lock acquires an exclusive write lock on the globals.
func (sg *SharedGlobals) Lock() {
	sg.mu.Lock()
}

// Unlock releases the write lock.
func (sg *SharedGlobals) Unlock() {
	sg.mu.Unlock()
}
