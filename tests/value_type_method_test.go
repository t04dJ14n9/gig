package tests

import (
	"sync"
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
)

// TestValueTypeMutexNonShared verifies value-type mutex works in non-stateful mode.
func TestValueTypeMutexNonShared(t *testing.T) {
	src := `
package main

import "sync"

var mu sync.Mutex
var counter int

func IncrementAndGet() int {
	mu.Lock()
	counter = counter + 1
	val := counter
	mu.Unlock()
	return val
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("IncrementAndGet")
	if err != nil {
		t.Fatalf("IncrementAndGet error: %v", err)
	}
	got := toInt64(result)
	if got != 1 {
		t.Fatalf("got %d, want 1", got)
	}
	t.Logf("Value-type Mutex (non-shared): Lock/Unlock works correctly")
}

// TestValueTypeMutexStatefulSequential verifies value-type mutex accumulates
// state across sequential calls in stateful mode.
func TestValueTypeMutexStatefulSequential(t *testing.T) {
	src := `
package main

import "sync"

var mu sync.Mutex
var counter int

func IncrementAndGet() int {
	mu.Lock()
	counter = counter + 1
	val := counter
	mu.Unlock()
	return val
}
`
	prog, err := gig.Build(src, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	for i := 1; i <= 10; i++ {
		result, err := prog.Run("IncrementAndGet")
		if err != nil {
			t.Fatalf("call %d error: %v", i, err)
		}
		got := toInt64(result)
		if got != int64(i) {
			t.Fatalf("call %d: got %d, want %d", i, got, i)
		}
	}
	t.Logf("Value-type Mutex (stateful, sequential): 10 increments correct")
}

// TestValueTypeMutexConcurrentExact verifies value-type sync.Mutex provides
// real mutual exclusion under concurrent access — the EXACT same behavior
// as `var mu *sync.Mutex = &sync.Mutex{}`. This is the KEY test proving
// the fix stores a heap-allocated object, not per-call copies.
func TestValueTypeMutexConcurrentExact(t *testing.T) {
	src := `
package main

import "sync"

var mu sync.Mutex
var counter int

func Increment() int {
	mu.Lock()
	counter = counter + 1
	val := counter
	mu.Unlock()
	return val
}

func GetCounter() int {
	mu.Lock()
	val := counter
	mu.Unlock()
	return val
}
`
	prog, err := gig.Build(src, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	const numGoroutines = 50
	const callsPerGoroutine = 20
	totalCalls := numGoroutines * callsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				_, err := prog.Run("Increment")
				if err != nil {
					t.Errorf("Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("GetCounter")
	if err != nil {
		t.Fatalf("GetCounter error: %v", err)
	}
	got := toInt64(result)
	if got != int64(totalCalls) {
		t.Fatalf("counter = %d, want exactly %d (mutex must prevent lost updates)", got, totalCalls)
	}
	t.Logf("Value-type Mutex (concurrent): counter = %d (exact, mutex works!)", got)
}

// TestValueTypeMapGlobal verifies that value-type sync.Map global works.
func TestValueTypeMapGlobal(t *testing.T) {
	src := `
package main

import "sync"

var m sync.Map

func StoreAndLoad() int {
	m.Store("key", 42)
	v, ok := m.Load("key")
	if ok && v == 42 {
		return 1
	}
	return 0
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	defer prog.Close()

	result, err := prog.Run("StoreAndLoad")
	if err != nil {
		t.Fatalf("StoreAndLoad error: %v", err)
	}
	got := toInt64(result)
	if got != 1 {
		t.Fatalf("got %d, want 1", got)
	}
	t.Logf("Value-type sync.Map: Store/Load works correctly")
}
