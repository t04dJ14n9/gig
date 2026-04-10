package tests

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	"git.woa.com/youngjin/gig/model/value"
)

// buildStatefulSrc is a helper to build source code with stateful globals.
func buildStatefulSrc(t *testing.T, src string) *gig.Program {
	t.Helper()
	prog, err := gig.Build(src, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	return prog
}

// ============================================================================
// 1. Basic Concurrent Counter Patterns
// ============================================================================

// TestConcurrentCounterAtomic tests atomic-like counter with mutex.
func TestConcurrentCounterAtomic(t *testing.T) {
	src := `
package main

import "sync"

var (
	mu      sync.Mutex
	counter int
)

func Increment() int {
	mu.Lock()
	counter++
	v := counter
	mu.Unlock()
	return v
}

func GetCounter() int {
	mu.Lock()
	v := counter
	mu.Unlock()
	return v
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numGoroutines = 100
	const incrementsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
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

	expected := int64(numGoroutines * incrementsPerGoroutine)
	got := toInt64(result)
	if got != expected {
		t.Fatalf("counter = %d, want %d", got, expected)
	}
	t.Logf("Atomic counter = %d (exact)", got)
}

// TestConcurrentCounterRWMutex tests RWMutex with many readers and few writers.
func TestConcurrentCounterRWMutex(t *testing.T) {
	src := `
package main

import "sync"

var (
	rwmu    sync.RWMutex
	counter int
)

func WriteCounter() int {
	rwmu.Lock()
	counter++
	v := counter
	rwmu.Unlock()
	return v
}

func ReadCounter() int {
	rwmu.RLock()
	v := counter
	rwmu.RUnlock()
	return v
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numWriters = 20
	const numReaders = 100
	const writesPerWriter = 5

	var wg sync.WaitGroup
	wg.Add(numWriters + numReaders)

	// Writers
	for i := 0; i < numWriters; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				_, err := prog.Run("WriteCounter")
				if err != nil {
					t.Errorf("WriteCounter error: %v", err)
					return
				}
			}
		}()
	}

	// Readers
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, err := prog.Run("ReadCounter")
				if err != nil {
					t.Errorf("ReadCounter error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	result, err := prog.Run("ReadCounter")
	if err != nil {
		t.Fatalf("ReadCounter error: %v", err)
	}

	expected := int64(numWriters * writesPerWriter)
	got := toInt64(result)
	if got != expected {
		t.Fatalf("counter = %d, want %d", got, expected)
	}
	t.Logf("RWMutex counter = %d after %d writers + %d readers", got, numWriters, numReaders)
}

// ============================================================================
// 2. sync.Map Comprehensive Tests
// ============================================================================

// TestSyncMapConcurrentOps tests concurrent store/load/delete operations.
func TestSyncMapConcurrentOps(t *testing.T) {
	src := `
package main

import "sync"

var m sync.Map

func Store(key string, value int) {
	m.Store(key, value)
}

func Load(key string) (int, bool) {
	v, ok := m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

func Delete(key string) {
	m.Delete(key)
}

func LoadOrStore(key string, value int) (int, bool) {
	v, loaded := m.LoadOrStore(key, value)
	return v.(int), loaded
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numGoroutines = 50
	const opsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent stores
	for i := 0; i < numGoroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", i, j)
				_, err := prog.Run("Store", key, i*1000+j)
				if err != nil {
					t.Errorf("Store error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Verify all keys exist
	missing := 0
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < opsPerGoroutine; j++ {
			key := fmt.Sprintf("key-%d-%d", i, j)
			result, err := prog.Run("Load", key)
			if err != nil {
				t.Errorf("Load error: %v", err)
				continue
			}
			if result == nil {
				missing++
			}
		}
	}

	if missing > 0 {
		t.Errorf("sync.Map: %d keys missing after concurrent stores", missing)
	} else {
		t.Logf("sync.Map: all %d keys present", numGoroutines*opsPerGoroutine)
	}
}

// ============================================================================
// 3. Channel-Based Patterns
// ============================================================================

// TestChannelPipeline tests pipeline pattern with channels.
func TestChannelPipeline(t *testing.T) {
	src := `
package main

func Pipeline() int {
	// Stage 1: generate numbers
	gen := make(chan int, 10)
	go func() {
		for i := 1; i <= 10; i++ {
			gen <- i
		}
		close(gen)
	}()

	// Stage 2: square numbers
	squares := make(chan int, 10)
	go func() {
		for n := range gen {
			squares <- n * n
		}
		close(squares)
	}()

	// Stage 3: sum
	sum := 0
	for n := range squares {
		sum += n
	}
	return sum
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	result, err := prog.Run("Pipeline")
	if err != nil {
		t.Fatalf("Pipeline error: %v", err)
	}

	// Sum of squares 1..10 = 385
	got := toInt64(result)
	expected := int64(385)
	if got != expected {
		t.Fatalf("Pipeline = %d, want %d", got, expected)
	}
	t.Logf("Pipeline = %d (correct)", got)
}

// ============================================================================
// 4. WaitGroup Patterns
// ============================================================================

// TestWaitGroupMultipleBatches tests multiple sequential WaitGroup uses.
func TestWaitGroupMultipleBatches(t *testing.T) {
	src := `
package main

import "sync"

func BatchWork(batchSize int) int {
	sum := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	wg.Add(batchSize)
	for i := 0; i < batchSize; i++ {
		go func(n int) {
			mu.Lock()
			sum += n
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	return sum
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	// Test multiple batch sizes
	for _, size := range []int{5, 10, 20} {
		result, err := prog.Run("BatchWork", size)
		if err != nil {
			t.Fatalf("BatchWork(%d) error: %v", size, err)
		}
		got := toInt64(result)
		expected := int64((size - 1) * size / 2) // sum(0..size-1)
		if got != expected {
			t.Fatalf("BatchWork(%d) = %d, want %d", size, got, expected)
		}
		t.Logf("BatchWork(%d) = %d (correct)", size, got)
	}
}

// ============================================================================
// 5. Complex Synchronization
// ============================================================================

// ============================================================================
// 6. Global State Persistence
// ============================================================================

// TestGlobalStatePersistence tests that global state persists across multiple Run calls.
func TestGlobalStatePersistence(t *testing.T) {
	src := `
package main

import "sync"

var (
	mu    sync.Mutex
	state int
)

func SetState(v int) {
	mu.Lock()
	state = v
	mu.Unlock()
}

func GetState() int {
	mu.Lock()
	v := state
	mu.Unlock()
	return v
}

func IncState() int {
	mu.Lock()
	state++
	v := state
	mu.Unlock()
	return v
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	// Set initial state
	_, err := prog.Run("SetState", 100)
	if err != nil {
		t.Fatalf("SetState error: %v", err)
	}

	// Verify state persists
	result, err := prog.Run("GetState")
	if err != nil {
		t.Fatalf("GetState error: %v", err)
	}
	if toInt64(result) != 100 {
		t.Fatalf("state = %d, want 100", toInt64(result))
	}

	// Increment multiple times
	for i := 0; i < 10; i++ {
		_, err := prog.Run("IncState")
		if err != nil {
			t.Fatalf("IncState error: %v", err)
		}
	}

	// Verify final state
	result, err = prog.Run("GetState")
	if err != nil {
		t.Fatalf("GetState error: %v", err)
	}
	if toInt64(result) != 110 {
		t.Fatalf("final state = %d, want 110", toInt64(result))
	}

	t.Log("Global state persistence works correctly")
}

// ============================================================================
// 7. Concurrent RunWithValues
// ============================================================================

// TestConcurrentRunWithValuesComprehensive tests concurrent RunWithValues calls.
func TestConcurrentRunWithValuesComprehensive(t *testing.T) {
	src := `
package main

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numGoroutines = 100
	var wg sync.WaitGroup
	var successCount int64
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			ctx := context.Background()
			args := []value.Value{
				value.FromInterface(int64(i)),
				value.FromInterface(int64(i * 2)),
				}
			result, err := prog.RunWithValues(ctx, "Add", args)
			if err != nil {
				t.Errorf("Add error: %v", err)
				return
			}
			expected := int64(i + i*2)
			if result.Int() != expected {
				t.Errorf("Add(%d, %d) = %d, want %d", i, i*2, result.Int(), expected)
				return
			}
			atomic.AddInt64(&successCount, 1)
		}()
	}

	wg.Wait()

	if successCount != numGoroutines {
		t.Fatalf("expected %d successes, got %d", numGoroutines, successCount)
	}
	t.Logf("Concurrent RunWithValues: %d successes", successCount)
}

// ============================================================================
// 8. Error Handling and Timeouts
// ============================================================================

// TestConcurrentTimeoutComprehensive tests that context timeout works under concurrent load.
func TestConcurrentTimeoutComprehensive(t *testing.T) {
	src := `
package main

func InfiniteLoop() int {
	i := 0
	for {
		i = i + 1
	}
	return i
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			_, err := prog.RunWithContext(ctx, "InfiniteLoop")
			if err == nil {
				t.Error("expected timeout error, got nil")
			}
		}()
	}

	wg.Wait()
	t.Log("Concurrent timeout handling works")
}

// ============================================================================
// 9. Nested Mutex Locks
// ============================================================================

// TestNestedMutexLocks tests that nested lock acquisition works correctly.
func TestNestedMutexLocks(t *testing.T) {
	src := `
package main

import "sync"

var (
	muA sync.Mutex
	muB sync.Mutex
	sum int
)

func NestedLockAB() int {
	muA.Lock()
	muB.Lock()
	sum = sum + 1
	v := sum
	muB.Unlock()
	muA.Unlock()
	return v
}

func NestedLockBA() int {
	muB.Lock()
	muA.Lock()
	sum = sum + 1
	v := sum
	muA.Unlock()
	muB.Unlock()
	return v
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	_, err := prog.Run("NestedLockAB")
	if err != nil {
		t.Fatalf("NestedLockAB error: %v", err)
	}

	_, err = prog.Run("NestedLockBA")
	if err != nil {
		t.Fatalf("NestedLockBA error: %v", err)
	}

	t.Log("Nested mutex locks with both orderings work correctly")
}

// ============================================================================
// 10. sync.Map Store and LoadAndDelete
// ============================================================================

// TestSyncMapStoreAndDelete tests sequential store + LoadAndDelete.
func TestSyncMapStoreAndDelete(t *testing.T) {
	src := `
package main

import "sync"

var m sync.Map

func StoreAndDelete(key string, value int) bool {
	m.Store(key, value)
	v, loaded := m.LoadAndDelete(key)
	if !loaded {
		return false
	}
	return v.(int) == value
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	const numOps = 100
	for i := 0; i < numOps; i++ {
		result, err := prog.Run("StoreAndDelete", fmt.Sprintf("key-%d", i), i)
		if err != nil {
			t.Fatalf("StoreAndDelete error at iteration %d: %v", i, err)
		}
		if result != true {
			t.Fatalf("StoreAndDelete(%d) = false, want true", i)
		}
	}
	t.Logf("All %d store/delete operations succeeded", numOps)
}

// ============================================================================
// 11. Channel Fan-Out Pattern
// ============================================================================

// TestChannelFanOut tests fan-out pattern where multiple workers read from a shared channel.
// Results are non-deterministic because job-to-worker assignment is random.
func TestChannelFanOut(t *testing.T) {
	src := `
package main

func FanOut() int {
	const workers = 10
	jobs := make(chan int, workers)
	results := make(chan int, workers)

	// Start workers — each reads from the shared jobs channel
	for w := 0; w < workers; w++ {
		go func(id int) {
			for job := range jobs {
				results <- job * id
			}
		}(w)
	}

	// Send 10 jobs
	for j := 1; j <= workers; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect all results
	sum := 0
	for i := 0; i < workers; i++ {
		sum += <-results
	}
	return sum
}
`
	prog := buildStatefulSrc(t, src)
	defer prog.Close()

	result, err := prog.Run("FanOut")
	if err != nil {
		t.Fatalf("FanOut error: %v", err)
	}

	got := toInt64(result)
	// Result is non-deterministic: depends on which worker picks which job.
	// Just verify it's non-negative (all jobs processed, no crash).
	if got < 0 {
		t.Fatalf("FanOut = %d (negative), expected non-negative", got)
	}
	t.Logf("FanOut = %d (non-deterministic, verified non-negative)", got)
}

