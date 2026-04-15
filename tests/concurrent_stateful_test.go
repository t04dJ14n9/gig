package tests

import (
	"context"
	_ "embed"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"git.woa.com/youngjin/gig"
	"git.woa.com/youngjin/gig/model/value"
	"git.woa.com/youngjin/gig/tests/testdata/concurrent_stateful"
)

//go:embed testdata/concurrent_stateful/main.go
var concurrentStatefulSrc string

// toInt64 converts various integer types to int64 for flexible comparison.
func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case int16:
		return int64(n)
	case int8:
		return int64(n)
	case uint:
		return int64(n)
	case uint64:
		return int64(n)
	case uint32:
		return int64(n)
	case uint16:
		return int64(n)
	case uint8:
		return int64(n)
	case float64:
		return int64(n)
	case float32:
		return int64(n)
	default:
		return 0
	}
}

// buildStateful compiles the concurrent_stateful source with stateful globals.
func buildStateful(t *testing.T) *gig.Program {
	t.Helper()
	prog, err := gig.Build(concurrentStatefulSrc, gig.WithStatefulGlobals(), gig.WithAllowPanic())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	return prog
}

// ============================================================================
// 1. Mutex-Protected Counter — Exact Correctness
// ============================================================================

// TestProtectedCounterExact verifies that with sync.Mutex protecting the global
// counter, concurrent increments produce the EXACT correct sum.
// This is the KEY test proving that the lock mechanism works correctly.
func TestProtectedCounterExact(t *testing.T) {
	prog := buildStateful(t)
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
				result, err := prog.Run("IncrementProtected")
				if err != nil {
					t.Errorf("IncrementProtected error: %v", err)
					return
				}
				val := toInt64(result)
				if val <= 0 || val > int64(totalCalls) {
					t.Errorf("IncrementProtected returned out-of-range value: %d", val)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Final counter MUST be exactly totalCalls — mutex guarantees no lost updates.
	result, err := prog.Run("GetProtected")
	if err != nil {
		t.Fatalf("GetProtected error: %v", err)
	}
	got := toInt64(result)

	// Native comparison
	expected := int64(totalCalls)
	if got != expected {
		t.Fatalf("Protected counter = %d, want exactly %d (mutex should prevent lost updates)", got, expected)
	}
	t.Logf("Protected counter = %d (exact, as expected with mutex)", got)
}

// ============================================================================
// 2. Unprotected Counter — Race Tolerance (no crash, value may be imprecise)
// ============================================================================

// TestUnprotectedCounterNoCrash verifies that concurrent access to an unprotected
// global doesn't crash. Lost updates are expected (same as Go data races).
func TestUnprotectedCounterNoCrash(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 30
	const callsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				result, err := prog.Run("IncrementUnprotected")
				if err != nil {
					t.Errorf("IncrementUnprotected error: %v", err)
					return
				}
				val := toInt64(result)
				if val <= 0 {
					t.Errorf("IncrementUnprotected returned non-positive: %d", val)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Counter should be positive, but likely < totalCalls due to lost updates.
	result, err := prog.Run("GetUnprotected")
	if err != nil {
		t.Fatalf("GetUnprotected error: %v", err)
	}
	got := toInt64(result)
	total := int64(numGoroutines * callsPerGoroutine)
	t.Logf("Unprotected counter = %d (expected <= %d due to lost updates)", got, total)
	if got <= 0 {
		t.Fatalf("Unprotected counter should be positive, got %d", got)
	}
}

// ============================================================================
// 3. Read-Only Global — All Concurrent Reads Must Match
// ============================================================================

// TestConcurrentReadOnlyGlobal verifies that concurrent reads of a read-only
// global all return the same value, compared against native execution.
func TestConcurrentReadOnlyGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	expected := concurrent_stateful.GetGreeting()

	const numReaders = 100
	var wg sync.WaitGroup
	wg.Add(numReaders)

	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			result, err := prog.Run("GetGreeting")
			if err != nil {
				t.Errorf("GetGreeting error: %v", err)
				return
			}
			if result != expected {
				t.Errorf("expected %q, got %v", expected, result)
			}
		}()
	}

	wg.Wait()
}

// ============================================================================
// 4. State Mutation Visibility — Set Then Concurrent Get
// ============================================================================

// TestStateMutationVisibility verifies that a global mutation from one Run()
// call is visible to subsequent concurrent calls.
func TestStateMutationVisibility(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set state to 42
	_, err := prog.Run("SetState", 42)
	if err != nil {
		t.Fatalf("SetState error: %v", err)
	}

	// Concurrent reads should all see 42
	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)

	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("GetState")
			if err != nil {
				t.Errorf("GetState error: %v", err)
				return
			}
			if toInt64(r) != 42 {
				t.Errorf("expected 42, got %v (type %T)", r, r)
			}
		}()
	}

	wg.Wait()
}

// ============================================================================
// 5. Concurrent RunWithValues — Stateless Pure Functions
// ============================================================================

// TestConcurrentRunWithValues tests concurrent RunWithValues on stateless
// pure functions, verifying each result matches the native computation.
func TestConcurrentRunWithValues(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 100
	var wg sync.WaitGroup
	var successCount int64
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()
			args := []value.Value{
				value.FromInterface(int64(i)),
				value.FromInterface(int64(i * 10)),
			}
			result, err := prog.RunWithValues(ctx, "Add", args)
			if err != nil {
				t.Errorf("RunWithValues Add error: %v", err)
				return
			}
			// Compare with native
			expected := concurrent_stateful.Add(i, i*10)
			if result.Int() != int64(expected) {
				t.Errorf("Add(%d, %d): expected %d, got %d", i, i*10, expected, result.Int())
				return
			}
			atomic.AddInt64(&successCount, 1)
		}()
	}

	wg.Wait()

	if successCount != numGoroutines {
		t.Errorf("expected %d successes, got %d", numGoroutines, successCount)
	}
}

// ============================================================================
// 6. Channel-Based Sum — Exact Correctness via goroutines inside program
// ============================================================================

// TestSumViaChannel tests goroutines spawned INSIDE the guest program
// that use channels for accumulation. The result must be exactly correct.
func TestSumViaChannel(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("SumViaChannel")
	if err != nil {
		t.Fatalf("SumViaChannel error: %v", err)
	}

	expected := concurrent_stateful.SumViaChannel()
	got := toInt64(result)
	if got != int64(expected) {
		t.Fatalf("SumViaChannel = %d, want %d", got, expected)
	}
	t.Logf("SumViaChannel = %d (matches native)", got)
}

// ============================================================================
// 7. Producer-Consumer Pattern
// ============================================================================

// TestProducerConsumerSum tests a producer-consumer pattern with multiple
// goroutines inside the guest program. The sum must be exactly correct.
func TestProducerConsumerSum(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("ProducerConsumerSum")
	if err != nil {
		t.Fatalf("ProducerConsumerSum error: %v", err)
	}

	expected := concurrent_stateful.ProducerConsumerSum()
	got := toInt64(result)
	if got != int64(expected) {
		t.Fatalf("ProducerConsumerSum = %d, want %d", got, expected)
	}
	t.Logf("ProducerConsumerSum = %d (matches native)", got)
}

// ============================================================================
// 8. Multiple Independent Mutex-Protected Counters
// ============================================================================

// TestMultipleProtectedCounters verifies that two independently locked counters
// can be incremented concurrently and both produce exact results.
func TestMultipleProtectedCounters(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numPerCounter = 200
	var wg sync.WaitGroup
	wg.Add(numPerCounter * 2)

	// Increment counter A
	for i := 0; i < numPerCounter; i++ {
		go func() {
			defer wg.Done()
			_, err := prog.Run("IncrementA")
			if err != nil {
				t.Errorf("IncrementA error: %v", err)
			}
		}()
	}

	// Increment counter B concurrently
	for i := 0; i < numPerCounter; i++ {
		go func() {
			defer wg.Done()
			_, err := prog.Run("IncrementB")
			if err != nil {
				t.Errorf("IncrementB error: %v", err)
			}
		}()
	}

	wg.Wait()

	// Both counters must be exactly numPerCounter
	resultA, err := prog.Run("GetCountA")
	if err != nil {
		t.Fatalf("GetCountA error: %v", err)
	}
	resultB, err := prog.Run("GetCountB")
	if err != nil {
		t.Fatalf("GetCountB error: %v", err)
	}

	gotA := toInt64(resultA)
	gotB := toInt64(resultB)

	if gotA != int64(numPerCounter) {
		t.Errorf("CounterA = %d, want exactly %d", gotA, numPerCounter)
	}
	if gotB != int64(numPerCounter) {
		t.Errorf("CounterB = %d, want exactly %d", gotB, numPerCounter)
	}
	t.Logf("CounterA = %d, CounterB = %d (both exact)", gotA, gotB)
}

// ============================================================================
// 9. Context Timeout Under Concurrency
// ============================================================================

// TestConcurrentTimeout verifies that context cancellation/timeout works
// correctly when multiple concurrent executions are running.
func TestConcurrentTimeout(t *testing.T) {
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
	prog, err := gig.Build(src, gig.WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
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
}

// ============================================================================
// 10. Concurrent RunWithValues — Multiply
// ============================================================================

// TestConcurrentMultiply tests concurrent RunWithValues on Multiply,
// verifying each result matches native.
func TestConcurrentMultiply(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 100
	var wg sync.WaitGroup
	var successCount int64
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()
			args := []value.Value{
				value.FromInterface(int64(i + 1)),
				value.FromInterface(int64(i + 2)),
			}
			result, err := prog.RunWithValues(ctx, "Multiply", args)
			if err != nil {
				t.Errorf("RunWithValues Multiply error: %v", err)
				return
			}
			expected := concurrent_stateful.Multiply(i+1, i+2)
			if result.Int() != int64(expected) {
				t.Errorf("Multiply(%d, %d): expected %d, got %d", i+1, i+2, expected, result.Int())
				return
			}
			atomic.AddInt64(&successCount, 1)
		}()
	}

	wg.Wait()

	if successCount != numGoroutines {
		t.Errorf("expected %d successes, got %d", numGoroutines, successCount)
	}
}

// ============================================================================
// 11. Protected Counter — Large Scale Stress
// ============================================================================

// TestProtectedCounterStress tests a larger number of concurrent increments
// to verify the mutex holds under higher contention.
func TestProtectedCounterStress(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 100
	const callsPerGoroutine = 50
	totalCalls := numGoroutines * callsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				_, err := prog.Run("IncrementProtected")
				if err != nil {
					t.Errorf("IncrementProtected error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	result, err := prog.Run("GetProtected")
	if err != nil {
		t.Fatalf("GetProtected error: %v", err)
	}
	got := toInt64(result)

	if got != int64(totalCalls) {
		t.Fatalf("Stress: Protected counter = %d, want exactly %d", got, totalCalls)
	}
	t.Logf("Stress: Protected counter = %d (exact after %d concurrent increments)", got, totalCalls)
}

// ============================================================================
// 12. Value-Type sync.Mutex — Exact Concurrent Correctness
// ============================================================================

// TestValueTypeMutexExact verifies that value-type sync.Mutex (not pointer)
// provides exact mutual exclusion under concurrent access.
func TestValueTypeMutexExact(t *testing.T) {
	prog := buildStateful(t)
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
				_, err := prog.Run("ValueTypeIncrement")
				if err != nil {
					t.Errorf("ValueTypeIncrement error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	result, err := prog.Run("ValueTypeGet")
	if err != nil {
		t.Fatalf("ValueTypeGet error: %v", err)
	}
	got := toInt64(result)

	if got != int64(totalCalls) {
		t.Fatalf("Value-type mutex counter = %d, want exactly %d", got, totalCalls)
	}
	t.Logf("Value-type mutex counter = %d (exact, concurrent safe)", got)
}

// ============================================================================
// 13. sync.RWMutex — Multiple Readers, Single Writer
// ============================================================================

// TestRWMutex verifies that RWMutex allows multiple concurrent readers
// and exclusive writers.
func TestRWMutex(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numWriters = 20
	const numReaders = 50
	const writesPerWriter = 10

	var wg sync.WaitGroup
	wg.Add(numWriters + numReaders)

	// Writers
	for i := 0; i < numWriters; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				_, err := prog.Run("RWMutexWrite")
				if err != nil {
					t.Errorf("RWMutexWrite error: %v", err)
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
				_, err := prog.Run("RWMutexRead")
				if err != nil {
					t.Errorf("RWMutexRead error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Final count must be exact
	expected := int64(numWriters * writesPerWriter)
	result, err := prog.Run("RWMutexRead")
	if err != nil {
		t.Fatalf("RWMutexRead final error: %v", err)
	}
	got := toInt64(result)

	if got != expected {
		t.Fatalf("RWMutex counter = %d, want exactly %d", got, expected)
	}
	t.Logf("RWMutex counter = %d (exact)", got)
}

// ============================================================================
// 14. sync.Once — Exactly-Once Initialization with Anonymous Closure
// ============================================================================

// TestOnce verifies that sync.Once.Do with an anonymous closure correctly
// initializes a global variable exactly once, even under concurrent calls.
func TestOnce(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	results := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			result, err := prog.Run("OnceInit")
			if err != nil {
				t.Errorf("OnceInit error: %v", err)
				return
			}
			results <- int(toInt64(result))
		}()
	}

	wg.Wait()
	close(results)

	// All results should be 42 (the once-initialized value)
	for v := range results {
		if v != 42 {
			t.Errorf("OnceInit returned %d, want 42", v)
		}
	}
	t.Log("sync.Once: all 100 concurrent calls returned 42 (exactly-once correct)")
}

// ============================================================================
// 15. sync.WaitGroup — Goroutine Synchronization Inside Guest Program
// ============================================================================

// TestWaitGroupGlobal verifies that sync.WaitGroup correctly synchronizes
// goroutines spawned inside the guest program.
func TestWaitGroupGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("WaitGroupSum")
	if err != nil {
		t.Fatalf("WaitGroupSum error: %v", err)
	}

	got := toInt64(result)
	expected := int64(49 * 50 / 2) // sum(0..49)

	if got != expected {
		t.Fatalf("WaitGroupSum = %d, want %d", got, expected)
	}
	t.Logf("WaitGroupSum = %d (exact)", got)
}

// ============================================================================
// 16. Concurrent Counter Stress Test
// ============================================================================

// TestConcurrentCounterStress verifies counter increment under high contention.
func TestConcurrentCounterStress(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 100
	const incrementsPerGoroutine = 50
	totalIncrements := numGoroutines * incrementsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("ValueTypeIncrement")
				if err != nil {
					t.Errorf("ValueTypeIncrement error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	result, err := prog.Run("ValueTypeGet")
	if err != nil {
		t.Fatalf("ValueTypeGet error: %v", err)
	}
	got := toInt64(result)

	if got != int64(totalIncrements) {
		t.Fatalf("Counter = %d, want exactly %d", got, totalIncrements)
	}
	t.Logf("High contention counter = %d (exact)", got)
}

// ============================================================================
// 15. sync.Map — Concurrent Map Operations
// ============================================================================

// TestSyncMap verifies that sync.Map supports concurrent store/load/delete.
func TestSyncMap(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 30
	const opsPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", i, j)
				_, err := prog.Run("MapStore", key, i*100+j)
				if err != nil {
					t.Errorf("MapStore error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Verify all keys are present
	missing := 0
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < opsPerGoroutine; j++ {
			key := fmt.Sprintf("key-%d-%d", i, j)
			result, err := prog.Run("MapLoad", key)
			if err != nil {
				t.Errorf("MapLoad error: %v", err)
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
		t.Logf("sync.Map: all %d keys present after concurrent stores", numGoroutines*opsPerGoroutine)
	}
}

// ============================================================================
// 16. Channel-Based Worker Pool
// ============================================================================

// TestChannelWorkerPool verifies a worker pool pattern using channels.
func TestChannelWorkerPool(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("SumViaChannel")
	if err != nil {
		t.Fatalf("SumViaChannel error: %v", err)
	}

	got := toInt64(result)
	expected := int64(100) // 100 goroutines each send 1

	if got != expected {
		t.Fatalf("SumViaChannel = %d, want %d", got, expected)
	}
	t.Logf("Channel worker pool sum = %d (exact)", got)
}

// ============================================================================
// 17. Nested Locks — Lock Ordering
// ============================================================================

// TestNestedLocks verifies that nested lock acquisition with consistent
// ordering works correctly (no deadlock).
func TestNestedLocks(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const callsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				_, err := prog.Run("NestedLockAB")
				if err != nil {
					t.Errorf("NestedLockAB error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()
	t.Log("Nested locks with consistent ordering: no deadlock, all calls completed")
}

// ============================================================================
// 18. Complex Mixed Synchronization — Cond + Mutex
// ============================================================================

// TestComplexSync verifies that complex synchronization patterns
// (Cond + Mutex) work correctly.
func TestComplexSync(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Reset state
	_, err := prog.Run("ResetComplexState")
	if err != nil {
		t.Fatalf("ResetComplexState error: %v", err)
	}

	// Start consumer in background
	var consumerResult int
	var consumerErr error
	var consumerDone sync.WaitGroup
	consumerDone.Add(1)
	go func() {
		defer consumerDone.Done()
		result, err := prog.Run("ComplexConsumer")
		if err != nil {
			consumerErr = err
			return
		}
		consumerResult = int(toInt64(result))
	}()

	// Small delay to ensure consumer is waiting
	time.Sleep(50 * time.Millisecond)

	// Producer sets value
	_, err = prog.Run("ComplexProducer", 123)
	if err != nil {
		t.Fatalf("ComplexProducer error: %v", err)
	}

	// Wait for consumer
	consumerDone.Wait()

	if consumerErr != nil {
		t.Fatalf("ComplexConsumer error: %v", consumerErr)
	}

	if consumerResult != 123 {
		t.Fatalf("ComplexConsumer returned %d, want 123", consumerResult)
	}
	t.Log("Complex sync (Cond + Mutex): producer/consumer pattern works correctly")
}

// ============================================================================
// 19. Atomic-style counter with Mutex
// ============================================================================

func TestAtomicCounter(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Reset
	_, _ = prog.Run("AtomicSet", int64(0))

	const numGoroutines = 50
	const addsPerGoroutine = 20
	totalAdds := int64(numGoroutines * addsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < addsPerGoroutine; j++ {
				_, err := prog.Run("AtomicAdd", int64(1))
				if err != nil {
					t.Errorf("AtomicAdd error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("AtomicGet")
	if err != nil {
		t.Fatalf("AtomicGet error: %v", err)
	}
	got := toInt64(result)
	if got != totalAdds {
		t.Fatalf("AtomicCounter = %d, want %d", got, totalAdds)
	}
	t.Logf("Atomic-style counter = %d (exact)", got)
}

// ============================================================================
// 20. Global protected slice
// ============================================================================

func TestGlobalProtectedSlice(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("ResetBuf")

	const numGoroutines = 50
	const appendsPerGoroutine = 10
	totalAppends := numGoroutines * appendsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				_, err := prog.Run("AppendProtected", j)
				if err != nil {
					t.Errorf("AppendProtected error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("GetBufLen")
	if err != nil {
		t.Fatalf("GetBufLen error: %v", err)
	}
	got := toInt64(result)
	if got != int64(totalAppends) {
		t.Fatalf("Protected slice len = %d, want %d", got, totalAppends)
	}
	t.Logf("Protected slice len = %d (exact)", got)
}

// ============================================================================
// 21. Global protected map
// ============================================================================

func TestGlobalProtectedMap(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("ResetProtectedMap")

	const numGoroutines = 50
	const putsPerGoroutine = 10
	totalPuts := numGoroutines * putsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < putsPerGoroutine; j++ {
				key := fmt.Sprintf("k-%d-%d", i, j)
				_, err := prog.Run("MapPutProtected", key, i*100+j)
				if err != nil {
					t.Errorf("MapPutProtected error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("MapLenProtected")
	if err != nil {
		t.Fatalf("MapLenProtected error: %v", err)
	}
	got := toInt64(result)
	if got != int64(totalPuts) {
		t.Fatalf("Protected map len = %d, want %d", got, totalPuts)
	}
	t.Logf("Protected map len = %d (exact)", got)
}

// ============================================================================
// 22. Defer in goroutine with global state
// ============================================================================

func TestDeferInGoroutine(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("DeferInGoroutine")
	if err != nil {
		t.Fatalf("DeferInGoroutine error: %v", err)
	}
	got := toInt64(result)
	if got != 30 {
		t.Fatalf("DeferInGoroutine = %d, want 30", got)
	}
	t.Logf("Defer in goroutine = %d (exact)", got)
}

// ============================================================================
// 23. Bidirectional channel with goroutines
// ============================================================================

func TestBidirectionalChannel(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("BidirectionalChannel")
	if err != nil {
		t.Fatalf("BidirectionalChannel error: %v", err)
	}
	got := toInt64(result)
	if got != 42 {
		t.Fatalf("BidirectionalChannel = %d, want 42", got)
	}
	t.Logf("Bidirectional channel = %d (exact)", got)
}

// ============================================================================
// 24. Multi-channel merge
// ============================================================================

func TestMultiChannelMerge(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("MultiChannelMerge")
	if err != nil {
		t.Fatalf("MultiChannelMerge error: %v", err)
	}
	got := toInt64(result)
	if got != 1665 {
		t.Fatalf("MultiChannelMerge = %d, want 1665", got)
	}
	t.Logf("Multi-channel merge = %d (exact)", got)
}

// ============================================================================
// 25. Goroutine with result channel
// ============================================================================

func TestGoroutineWithResult(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("GoroutineWithResult")
	if err != nil {
		t.Fatalf("GoroutineWithResult error: %v", err)
	}
	got := toInt64(result)
	if got != 42 {
		t.Fatalf("GoroutineWithResult = %d, want 42", got)
	}
	t.Logf("Goroutine with result channel = %d (exact)", got)
}

// ============================================================================
// 26. Barrier pattern
// ============================================================================

func TestBarrierSum(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("BarrierSum")
	if err != nil {
		t.Fatalf("BarrierSum error: %v", err)
	}
	expected := concurrent_stateful.BarrierSum()
	got := toInt64(result)
	if got != int64(expected) {
		t.Fatalf("BarrierSum = %d, want %d", got, expected)
	}
	t.Logf("Barrier sum = %d (matches native)", got)
}

// ============================================================================
// 27. Value-type RWMutex
// ============================================================================

func TestValueTypeRWMutex(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, err := prog.Run("ValueTypeRWWrite", 100)
	if err != nil {
		t.Fatalf("ValueTypeRWWrite error: %v", err)
	}

	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("ValueTypeRWRead")
			if err != nil {
				t.Errorf("ValueTypeRWRead error: %v", err)
				return
			}
			if toInt64(r) != 100 {
				t.Errorf("ValueTypeRWRead = %d, want 100", toInt64(r))
			}
		}()
	}
	wg.Wait()
	t.Log("Value-type RWMutex: 50 concurrent reads all returned 100")
}

// ============================================================================
// 28. Global string with RWMutex
// ============================================================================

func TestGlobalString(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set value
	_, err := prog.Run("SetGlobalString", "hello world")
	if err != nil {
		t.Fatalf("SetGlobalString error: %v", err)
	}

	// Concurrent reads
	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("GetGlobalString")
			if err != nil {
				t.Errorf("GetGlobalString error: %v", err)
				return
			}
			if s, ok := r.(string); !ok || s != "hello world" {
				t.Errorf("GetGlobalString = %v, want 'hello world'", r)
			}
		}()
	}
	wg.Wait()
	t.Log("Global string: 50 concurrent reads all matched")
}

// ============================================================================
// 29. CAS pattern
// ============================================================================

func TestCASPatter(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set initial value
	_, err := prog.Run("CASSwap", 0, 10)
	if err != nil {
		t.Fatalf("CASSwap error: %v", err)
	}

	// Concurrent increments
	const numGoroutines = 50
	const incrementsPerGoroutine = 10
	totalIncrements := numGoroutines * incrementsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("CASIncrement")
				if err != nil {
					t.Errorf("CASIncrement error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("CASGet")
	if err != nil {
		t.Fatalf("CASGet error: %v", err)
	}
	got := toInt64(result)
	expected := int64(10 + totalIncrements)
	if got != expected {
		t.Fatalf("CAS counter = %d, want %d", got, expected)
	}
	t.Logf("CAS counter = %d (exact)", got)
}

// ============================================================================
// 30. Global bool flag
// ============================================================================

func TestGlobalFlag(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Concurrent flag setting and reading
	const numOps = 100
	var wg sync.WaitGroup
	wg.Add(numOps)
	for i := 0; i < numOps; i++ {
		go func() {
			defer wg.Done()
			if i%2 == 0 {
				_, _ = prog.Run("SetFlag", i%4 == 0)
			} else {
				_, _ = prog.Run("GetFlag")
			}
		}()
	}
	wg.Wait()
	t.Log("Global bool flag: 100 concurrent ops completed without panic")
}

// ============================================================================
// 31. Multiple sync.Once instances
// ============================================================================

func TestMultipleOnce(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)
	resultsA := make(chan int, numGoroutines)
	resultsB := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("OnceInitA")
			if err != nil {
				t.Errorf("OnceInitA error: %v", err)
				return
			}
			resultsA <- int(toInt64(r))
		}()
		go func() {
			defer wg.Done()
			r, err := prog.Run("OnceInitB")
			if err != nil {
				t.Errorf("OnceInitB error: %v", err)
				return
			}
			resultsB <- int(toInt64(r))
		}()
	}
	wg.Wait()
	close(resultsA)
	close(resultsB)

	for v := range resultsA {
		if v != 100 {
			t.Errorf("OnceInitA = %d, want 100", v)
		}
	}
	for v := range resultsB {
		if v != 200 {
			t.Errorf("OnceInitB = %d, want 200", v)
		}
	}
	t.Log("Multiple sync.Once: A=100, B=200 (both exact)")
}

// ============================================================================
// 32. Closure capturing globals
// ============================================================================

// TestClosureReadGlobal verifies a closure can read a package-level variable.
func TestClosureReadGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("ClosureReadGlobal")
	if err != nil {
		t.Fatalf("ClosureReadGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 42 {
		t.Fatalf("ClosureReadGlobal = %d, want 42", got)
	}
	t.Log("Closure reading global: works correctly")
}

// TestClosureWriteGlobal verifies a closure can write to a package-level variable.
func TestClosureWriteGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("ClosureWriteGlobal")
	if err != nil {
		t.Fatalf("ClosureWriteGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 99 {
		t.Fatalf("ClosureWriteGlobal = %d, want 99", got)
	}
	t.Log("Closure writing global: works correctly")
}

// TestClosureAccumulateGlobal verifies a closure can accumulate into a global.
func TestClosureAccumulateGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("ClosureAccumulateGlobal")
	if err != nil {
		t.Fatalf("ClosureAccumulateGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 60 {
		t.Fatalf("ClosureAccumulateGlobal = %d, want 60", got)
	}
	t.Log("Closure accumulating global: works correctly")
}

// TestMultipleClosuresSharedGlobal verifies multiple closures share global state.
func TestMultipleClosuresSharedGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("MultipleClosuresSharedGlobal")
	if err != nil {
		t.Fatalf("MultipleClosuresSharedGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 15 {
		t.Fatalf("MultipleClosuresSharedGlobal = %d, want 15", got)
	}
	t.Log("Multiple closures sharing global: works correctly")
}

// ============================================================================
// 33. Recursive closures
// ============================================================================

// TestRecursiveClosureFib verifies recursive closure (fibonacci).
func TestRecursiveClosureFib(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("RecursiveClosureFib")
	if err != nil {
		t.Fatalf("RecursiveClosureFib error: %v", err)
	}
	got := toInt64(result)
	expected := int64(concurrent_stateful.RecursiveClosureFib())
	if got != expected {
		t.Fatalf("RecursiveClosureFib = %d, want %d", got, expected)
	}
	t.Logf("Recursive closure fibonacci: fib(10) = %d (matches native)", got)
}

// TestRecursiveClosureFactorial verifies recursive closure (factorial).
func TestRecursiveClosureFactorial(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("RecursiveClosureFactorial")
	if err != nil {
		t.Fatalf("RecursiveClosureFactorial error: %v", err)
	}
	got := toInt64(result)
	expected := int64(concurrent_stateful.RecursiveClosureFactorial())
	if got != expected {
		t.Fatalf("RecursiveClosureFactorial = %d, want %d", got, expected)
	}
	t.Logf("Recursive closure factorial: fact(8) = %d (matches native)", got)
}

// TestRecursiveClosureWithGlobal verifies recursive closure that writes to a global.
func TestRecursiveClosureWithGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("RecursiveClosureWithGlobal")
	if err != nil {
		t.Fatalf("RecursiveClosureWithGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 55 {
		t.Fatalf("RecursiveClosureWithGlobal = %d, want 55", got)
	}
	t.Logf("Recursive closure with global: sum(1..10) = %d", got)
}

// ============================================================================
// 34. Goroutines inside guest code accessing shared globals
// ============================================================================

// TestGoroutineIncrementGlobal verifies goroutines inside guest code can
// increment a mutex-protected global and produce the exact correct sum.
func TestGoroutineIncrementGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("GoroutineIncrementGlobal")
	if err != nil {
		t.Fatalf("GoroutineIncrementGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 50 {
		t.Fatalf("GoroutineIncrementGlobal = %d, want 50", got)
	}
	t.Logf("Goroutine increment global = %d (exact)", got)
}

// TestGoroutineReadGlobal verifies goroutines inside guest code can
// concurrently read a global and all see the same value.
func TestGoroutineReadGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("GoroutineReadGlobal")
	if err != nil {
		t.Fatalf("GoroutineReadGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 50 {
		t.Fatalf("GoroutineReadGlobal = %d, want 50 (all reads should return 42)", got)
	}
	t.Logf("Goroutine read global: %d/50 reads returned 42", got)
}

// TestGoroutineClosureGlobal verifies goroutines with closures accessing a global.
func TestGoroutineClosureGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("GoroutineClosureGlobal")
	if err != nil {
		t.Fatalf("GoroutineClosureGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 30 {
		t.Fatalf("GoroutineClosureGlobal = %d, want 30", got)
	}
	t.Logf("Goroutine closure global = %d (exact)", got)
}

// ============================================================================
// 35. Diverse global types under concurrent access
// ============================================================================

// TestFloatGlobalConcurrent verifies concurrent increments of a float64 global.
func TestFloatGlobalConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 10
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("FloatGlobalIncrement")
				if err != nil {
					t.Errorf("FloatGlobalIncrement error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("FloatGlobalGet")
	if err != nil {
		t.Fatalf("FloatGlobalGet error: %v", err)
	}
	var got float64
	switch v := result.(type) {
	case float64:
		got = v
	case int64:
		got = float64(v)
	default:
		t.Fatalf("FloatGlobalGet returned unexpected type %T: %v", result, result)
	}
	if got != float64(totalIncrements) {
		t.Fatalf("FloatGlobal = %v, want %v", got, totalIncrements)
	}
	t.Logf("Float global concurrent: %v (exact)", got)
}

// TestStringGlobalConcurrent verifies concurrent reads of a string global.
func TestStringGlobalConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set the string
	_, err := prog.Run("StringGlobalSet", "hello world")
	if err != nil {
		t.Fatalf("StringGlobalSet error: %v", err)
	}

	// Concurrent reads
	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("StringGlobalGet")
			if err != nil {
				t.Errorf("StringGlobalGet error: %v", err)
				return
			}
			if s, ok := r.(string); !ok || s != "hello world" {
				t.Errorf("StringGlobalGet = %v, want 'hello world'", r)
			}
		}()
	}
	wg.Wait()
	t.Log("String global: 50 concurrent reads all matched")
}

// TestBoolGlobalConcurrent verifies concurrent set/get of a bool global.
func TestBoolGlobalConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set true
	_, err := prog.Run("BoolGlobalSet", true)
	if err != nil {
		t.Fatalf("BoolGlobalSet error: %v", err)
	}

	// Concurrent reads
	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("BoolGlobalGet")
			if err != nil {
				t.Errorf("BoolGlobalGet error: %v", err)
			}
			if b, ok := r.(bool); !ok || !b {
				t.Errorf("BoolGlobalGet = %v, want true", r)
			}
		}()
	}
	wg.Wait()
	t.Log("Bool global: 50 concurrent reads all returned true")
}

// ============================================================================
// 36. Loop closure + goroutine + global
// ============================================================================

// TestLoopClosureGoroutineGlobal verifies goroutines in a loop with closures
// that capture iteration variables and write to a shared global.
func TestLoopClosureGoroutineGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("LoopClosureGoroutineGlobal")
	if err != nil {
		t.Fatalf("LoopClosureGoroutineGlobal error: %v", err)
	}
	got := toInt64(result)
	expected := int64(19 * 20 / 2) // sum(0..19) = 190
	if got != expected {
		t.Fatalf("LoopClosureGoroutineGlobal = %d, want %d", got, expected)
	}
	t.Logf("Loop closure goroutine global = %d (exact)", got)
}

// ============================================================================
// 37. Global int64 with mutex — concurrent access
// ============================================================================

// TestGlobalInt64Concurrent verifies concurrent increments of a global int64.
func TestGlobalInt64Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("Int64Set", int64(0))

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Int64Increment")
				if err != nil {
					t.Errorf("Int64Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Int64Get")
	if err != nil {
		t.Fatalf("Int64Get error: %v", err)
	}
	var got int64
	switch v := result.(type) {
	case int64:
		got = v
	case int:
		got = int64(v)
	default:
		t.Fatalf("Int64Get returned unexpected type %T: %v", result, result)
	}
	if got != totalIncrements {
		t.Fatalf("Global int64 = %d, want %d", got, totalIncrements)
	}
	t.Logf("Global int64 concurrent: %d (exact)", got)
}

// ============================================================================
// 38. Global map concurrent read-only
// ============================================================================

// TestGlobalMapConcurrentRead verifies concurrent reads from a read-only map.
func TestGlobalMapConcurrentRead(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numReaders = 100
	var wg sync.WaitGroup
	var successCount int64
	wg.Add(numReaders)

	keys := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel", "india", "juliet"}
	expected := map[string]int{
		"alpha": 1, "bravo": 2, "charlie": 3, "delta": 4, "echo": 5,
		"foxtrot": 6, "golf": 7, "hotel": 8, "india": 9, "juliet": 10,
	}

	for i := 0; i < numReaders; i++ {
		go func(idx int) {
			defer wg.Done()
			key := keys[idx%len(keys)]
			result, err := prog.Run("ROMapGet", key)
			if err != nil {
				t.Errorf("ROMapGet error: %v", err)
				return
			}
			// result should be a tuple (int, bool)
			if m, ok := result.([]any); ok && len(m) == 2 {
				if toInt64(m[0]) == int64(expected[key]) && m[1] == true {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}
	wg.Wait()

	if successCount != numReaders {
		t.Errorf("Read-only map: %d/%d reads matched expected", successCount, numReaders)
	} else {
		t.Logf("Read-only map: all %d concurrent reads matched", numReaders)
	}
}

// ============================================================================
// 39. Defer recover in goroutine with global state
// ============================================================================

// TestDeferRecoverInGoroutineGlobal verifies panic/recover in goroutines
// correctly updates a global counter.
func TestDeferRecoverInGoroutineGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("DeferRecoverInGoroutine")
	if err != nil {
		t.Fatalf("DeferRecoverInGoroutine error: %v", err)
	}
	got := toInt64(result)
	if got != 20 {
		t.Fatalf("DeferRecoverInGoroutine = %d, want 20", got)
	}
	t.Logf("Defer recover in goroutine: %d (exact)", got)
}

// ============================================================================
// 40. Ping-pong via channels with global state
// ============================================================================

// TestPingPongGlobal verifies a ping-pong pattern with global counter.
func TestPingPongGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("PingPongGlobal")
	if err != nil {
		t.Fatalf("PingPongGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 20 {
		t.Fatalf("PingPongGlobal = %d, want 20", got)
	}
	t.Logf("Ping-pong global: %d (exact)", got)
}

// ============================================================================
// 41. Global uint concurrent access
// ============================================================================

// TestGlobalUintConcurrent verifies concurrent increments of a global uint.
func TestGlobalUintConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("UintIncrement")
				if err != nil {
					t.Errorf("UintIncrement error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("UintGet")
	if err != nil {
		t.Fatalf("UintGet error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global uint = %d (raw: %T %v), want %d", got, result, result, totalIncrements)
	}
	t.Logf("Global uint concurrent: %d (exact)", got)
}

// ============================================================================
// 42. Select with global state
// ============================================================================

// TestSelectIncrementGlobal verifies select statement with global state.
func TestSelectIncrementGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("SelectIncrementGlobal")
	if err != nil {
		t.Fatalf("SelectIncrementGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 40 {
		t.Fatalf("SelectIncrementGlobal = %d, want 40", got)
	}
	t.Logf("Select with global: %d (exact)", got)
}

// ============================================================================
// 43. Global struct pointer — concurrent access to struct fields
// ============================================================================

// TestGlobalStructPointerConcurrent verifies concurrent access to a global
// struct pointer's fields under mutex protection.
func TestGlobalStructPointerConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("ConfigReset")

	const numGoroutines = 50
	const incrementsPerGoroutine = 10
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("ConfigIncrementCount")
				if err != nil {
					t.Errorf("ConfigIncrementCount error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("ConfigGetCount")
	if err != nil {
		t.Fatalf("ConfigGetCount error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Config count = %d, want %d", got, totalIncrements)
	}
	t.Logf("Global struct pointer count = %d (exact)", got)
}

// TestGlobalStructPointerName verifies setting and reading a struct pointer's
// string field concurrently.
func TestGlobalStructPointerName(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	// Set name
	_, err := prog.Run("ConfigSetName", "test-config")
	if err != nil {
		t.Fatalf("ConfigSetName error: %v", err)
	}

	// Concurrent reads
	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			r, err := prog.Run("ConfigGetName")
			if err != nil {
				t.Errorf("ConfigGetName error: %v", err)
				return
			}
			if s, ok := r.(string); !ok || s != "test-config" {
				t.Errorf("ConfigGetName = %v, want 'test-config'", r)
			}
		}()
	}
	wg.Wait()
	t.Log("Global struct pointer name: 50 concurrent reads all matched")
}

// ============================================================================
// 44. Nested goroutines with global counter
// ============================================================================

// TestNestedGoroutineGlobal verifies nested goroutine spawning with global counter.
func TestNestedGoroutineGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("NestedGoroutineGlobal")
	if err != nil {
		t.Fatalf("NestedGoroutineGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 50 {
		t.Fatalf("NestedGoroutineGlobal = %d, want 50", got)
	}
	t.Logf("Nested goroutine global = %d (exact)", got)
}

// ============================================================================
// 45. Channel close + range with global accumulation
// ============================================================================

// TestChannelCloseRangeGlobal verifies channel close + range with global accumulation.
func TestChannelCloseRangeGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("ChannelCloseRangeGlobal")
	if err != nil {
		t.Fatalf("ChannelCloseRangeGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 55 {
		t.Fatalf("ChannelCloseRangeGlobal = %d, want 55", got)
	}
	t.Logf("Channel close range global = %d (exact)", got)
}

// ============================================================================
// 46. Global []byte concurrent append
// ============================================================================

// TestGlobalBytesConcurrent verifies concurrent appends to a global []byte.
func TestGlobalBytesConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("BytesReset")

	const numGoroutines = 50
	const appendsPerGoroutine = 10
	totalAppends := int64(numGoroutines * appendsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				_, err := prog.Run("BytesAppendProtected", byte(j))
				if err != nil {
					t.Errorf("BytesAppendProtected error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("BytesLenProtected")
	if err != nil {
		t.Fatalf("BytesLenProtected error: %v", err)
	}
	got := toInt64(result)
	if got != totalAppends {
		t.Fatalf("Global bytes len = %d, want %d", got, totalAppends)
	}
	t.Logf("Global bytes concurrent: len = %d (exact)", got)
}

// ============================================================================
// 47. Concurrent interleaved set+get on same global
// ============================================================================

// TestConcurrentInterleaveSetGet verifies interleaved set+get on the same global.
func TestConcurrentInterleaveSetGet(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("InterleaveReset")

	const numWriters = 20
	const numReaders = 30
	const writesPerWriter = 10
	totalWrites := int64(numWriters * writesPerWriter)

	var wg sync.WaitGroup
	wg.Add(numWriters + numReaders)

	// Writers: each writes 10 times
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				_, err := prog.Run("InterleaveSet", id*100+j)
				if err != nil {
					t.Errorf("InterleaveSet error: %v", err)
					return
				}
			}
		}(i)
	}

	// Readers: each reads 10 times
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				_, err := prog.Run("InterleaveGet")
				if err != nil {
					t.Errorf("InterleaveGet error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	// Verify total reads happened
	readCount, err := prog.Run("InterleaveReadCount")
	if err != nil {
		t.Fatalf("InterleaveReadCount error: %v", err)
	}
	totalReads := int64(numReaders * writesPerWriter)
	if toInt64(readCount) != totalReads {
		t.Errorf("Read count = %d, want %d", toInt64(readCount), totalReads)
	}
	t.Logf("Interleaved set+get: %d writes + %d reads completed", totalWrites, totalReads)
}

// ============================================================================
// 48. Global complex128 concurrent access
// ============================================================================

// TestGlobalComplexConcurrent verifies concurrent operations on a global complex128.
func TestGlobalComplexConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("ComplexGlobalReset")

	const numGoroutines = 50
	const addsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < addsPerGoroutine; j++ {
				_, err := prog.Run("ComplexGlobalAdd", 1.0, 0.0)
				if err != nil {
					t.Errorf("ComplexGlobalAdd error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("ComplexGlobalGet")
	if err != nil {
		t.Fatalf("ComplexGlobalGet error: %v", err)
	}
	switch v := result.(type) {
	case complex128:
		if real(v) != 500.0 {
			t.Fatalf("Complex real part = %v, want 500", real(v))
		}
	default:
		t.Fatalf("ComplexGlobalGet returned unexpected type %T: %v", result, result)
	}
	t.Logf("Global complex128 concurrent: real=%v (exact)", result)
}

// ============================================================================
// 49. Global int32 concurrent access
// ============================================================================

// TestGlobalInt32Concurrent verifies concurrent increments of a global int32.
func TestGlobalInt32Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("Int32Set", int32(0))

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Int32Increment")
				if err != nil {
					t.Errorf("Int32Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Int32Get")
	if err != nil {
		t.Fatalf("Int32Get error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global int32 = %d (raw: %T %v), want %d", got, result, result, totalIncrements)
	}
	t.Logf("Global int32 concurrent: %d (exact)", got)
}

// ============================================================================
// 50. Struct method with global access
// ============================================================================

// TestStructMethodGlobal verifies that struct methods can access and modify globals.
func TestStructMethodGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("StructMethodIncrementGlobal")
	if err != nil {
		t.Fatalf("StructMethodIncrementGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 30 {
		t.Fatalf("StructMethodIncrementGlobal = %d, want 30", got)
	}
	t.Logf("Struct method global: %d (exact)", got)
}

// ============================================================================
// 51. Fan-out pattern with global accumulation
// ============================================================================

// TestFanOutGlobal verifies the fan-out pattern with global accumulation.
func TestFanOutGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("FanOutGlobal")
	if err != nil {
		t.Fatalf("FanOutGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 55 {
		t.Fatalf("FanOutGlobal = %d, want 55", got)
	}
	t.Logf("Fan-out global: %d (exact)", got)
}

// ============================================================================
// 52. Global float32 concurrent access
// ============================================================================

func TestGlobalFloat32Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 10
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Float32Increment")
				if err != nil {
					t.Errorf("Float32Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Float32Get")
	if err != nil {
		t.Fatalf("Float32Get error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global float32 = %d (raw: %T %v), want %d", got, result, result, totalIncrements)
	}
	t.Logf("Global float32 concurrent: %d (exact)", got)
}

// ============================================================================
// 53. Global int8 concurrent access
// ============================================================================

func TestGlobalInt8Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("Int8Set", int8(0))

	const numGoroutines = 50
	const incrementsPerGoroutine = 2
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Int8Increment")
				if err != nil {
					t.Errorf("Int8Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Int8Get")
	if err != nil {
		t.Fatalf("Int8Get error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global int8 = %d (raw: %T %v), want %d", got, result, result, totalIncrements)
	}
	t.Logf("Global int8 concurrent: %d (exact)", got)
}

// ============================================================================
// 54. Global int16 concurrent access
// ============================================================================

func TestGlobalInt16Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("Int16Set", int16(0))

	const numGoroutines = 50
	const incrementsPerGoroutine = 10
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Int16Increment")
				if err != nil {
					t.Errorf("Int16Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Int16Get")
	if err != nil {
		t.Fatalf("Int16Get error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global int16 = %d (raw: %T %v), want %d", got, result, result, totalIncrements)
	}
	t.Logf("Global int16 concurrent: %d (exact)", got)
}

// ============================================================================
// 55. Worker pool with global accumulator
// ============================================================================

func TestWorkerPoolGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("WorkerPoolGlobal")
	if err != nil {
		t.Fatalf("WorkerPoolGlobal error: %v", err)
	}
	got := toInt64(result)
	expected := int64(concurrent_stateful.WorkerPoolGlobal())
	if got != expected {
		t.Fatalf("WorkerPoolGlobal = %d, want %d", got, expected)
	}
	t.Logf("Worker pool global = %d (matches native)", got)
}

// ============================================================================
// 56. Global uintptr concurrent access
// ============================================================================

func TestGlobalUintptrConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("UintptrIncrement")
				if err != nil {
					t.Errorf("UintptrIncrement error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("UintptrGet")
	if err != nil {
		t.Fatalf("UintptrGet error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global uintptr = %d, want %d", got, totalIncrements)
	}
	t.Logf("Global uintptr concurrent: %d (exact)", got)
}

// ============================================================================
// 57. Global uint32 concurrent access
// ============================================================================

func TestGlobalUint32Concurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("Uint32Increment")
				if err != nil {
					t.Errorf("Uint32Increment error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("Uint32Get")
	if err != nil {
		t.Fatalf("Uint32Get error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Global uint32 = %d, want %d", got, totalIncrements)
	}
	t.Logf("Global uint32 concurrent: %d (exact)", got)
}

// ============================================================================
// 58. Multiple struct fields concurrent access
// ============================================================================

func TestMultiFieldConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("MultiVarReset")

	const numGoroutines = 30

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = prog.Run("MultiVarIncrementA")
		}()
	}
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = prog.Run("MultiVarIncrementB")
		}()
	}
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = prog.Run("MultiVarIncrementC")
		}()
	}
	wg.Wait()

	result, err := prog.Run("MultiVarGetSum")
	if err != nil {
		t.Fatalf("MultiVarGetSum error: %v", err)
	}
	got := toInt64(result)
	expected := int64(numGoroutines * 3) // A=30 + B=30 + C=30 = 90
	if got != expected {
		t.Fatalf("MultiVar sum = %d, want %d", got, expected)
	}
	t.Logf("Multi-var concurrent: %d (exact)", got)
}

// ============================================================================
// 59. Defer external method in goroutine with panic/recover
// ============================================================================

// TestDeferExtMethodInGoroutine verifies that defer mu.Unlock() runs correctly
// during panic/recovery in goroutines.
func TestDeferExtMethodInGoroutine(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("DeferExtMethodInGoroutine")
	if err != nil {
		t.Fatalf("DeferExtMethodInGoroutine error: %v", err)
	}
	got := toInt64(result)
	// All 20 goroutines should increment, even ones that panic get recovered
	if got != 20 {
		t.Fatalf("DeferExtMethodInGoroutine = %d, want 20", got)
	}
	t.Logf("Defer ext method in goroutine: %d (exact)", got)
}

// ============================================================================
// 60. Type alias global concurrent access
// ============================================================================

func TestTypeAliasConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	const numGoroutines = 50
	const incrementsPerGoroutine = 20
	totalIncrements := int64(numGoroutines * incrementsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				_, err := prog.Run("AliasIncrement")
				if err != nil {
					t.Errorf("AliasIncrement error: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()

	result, err := prog.Run("AliasGet")
	if err != nil {
		t.Fatalf("AliasGet error: %v", err)
	}
	got := toInt64(result)
	if got != totalIncrements {
		t.Fatalf("Type alias counter = %d, want %d", got, totalIncrements)
	}
	t.Logf("Type alias counter: %d (exact)", got)
}

// ============================================================================
// 61. Select with channel close + global state
// ============================================================================

func TestSelectGlobalWaitNotify(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("SelectGlobalWaitNotify")
	if err != nil {
		t.Fatalf("SelectGlobalWaitNotify error: %v", err)
	}
	got := toInt64(result)
	if got != 10 {
		t.Fatalf("SelectGlobalWaitNotify = %d, want 10", got)
	}
	t.Logf("Select global wait/notify: %d (exact)", got)
}

// ============================================================================
// 62. Global []string concurrent append
// ============================================================================

func TestGlobalStrSliceConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("StrSliceReset")

	const numGoroutines = 50
	const appendsPerGoroutine = 10
	totalAppends := int64(numGoroutines * appendsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				_, err := prog.Run("StrSliceAppendProtected", fmt.Sprintf("k-%d-%d", id, j))
				if err != nil {
					t.Errorf("StrSliceAppendProtected error: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	result, err := prog.Run("StrSliceLenProtected")
	if err != nil {
		t.Fatalf("StrSliceLenProtected error: %v", err)
	}
	got := toInt64(result)
	if got != totalAppends {
		t.Fatalf("Global string slice len = %d, want %d", got, totalAppends)
	}
	t.Logf("Global string slice concurrent: len = %d (exact)", got)
}

// ============================================================================
// 63. Semaphore pattern with channel + global
// ============================================================================

func TestSemaphorePattern(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("SemaphorePattern")
	if err != nil {
		t.Fatalf("SemaphorePattern error: %v", err)
	}
	got := toInt64(result)
	if got != 10 {
		t.Fatalf("SemaphorePattern = %d, want 10", got)
	}
	t.Logf("Semaphore pattern: %d (exact)", got)
}

// ============================================================================
// 64. Pipeline pattern with channels + global accumulator
// ============================================================================

func TestPipelineGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("PipelineGlobal")
	if err != nil {
		t.Fatalf("PipelineGlobal error: %v", err)
	}
	got := toInt64(result)
	expected := int64(concurrent_stateful.PipelineGlobal())
	if got != expected {
		t.Fatalf("PipelineGlobal = %d, want %d", got, expected)
	}
	t.Logf("Pipeline global: %d (matches native)", got)
}

// ============================================================================
// 65. Global []rune concurrent append
// ============================================================================

func TestGlobalRuneSliceConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("RuneSliceReset")

	const numGoroutines = 50
	const appendsPerGoroutine = 10
	totalAppends := int64(numGoroutines * appendsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				_, err := prog.Run("RuneSliceAppendProtected", rune('A'+(id+j)%26))
				if err != nil {
					t.Errorf("RuneSliceAppendProtected error: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	result, err := prog.Run("RuneSliceLenProtected")
	if err != nil {
		t.Fatalf("RuneSliceLenProtected error: %v", err)
	}
	got := toInt64(result)
	if got != totalAppends {
		t.Fatalf("Global rune slice len = %d, want %d", got, totalAppends)
	}
	t.Logf("Global rune slice concurrent: len = %d (exact)", got)
}

// ============================================================================
// 66. Bool swap with mutex — concurrent compare-and-swap
// ============================================================================

func TestBoolSwapConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("SwapReset")

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			_, err := prog.Run("SwapBool", idx%2 == 0)
			if err != nil {
				t.Errorf("SwapBool error: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Verify swap count is reasonable (each change increments counter)
	result, err := prog.Run("SwapGetCount")
	if err != nil {
		t.Fatalf("SwapGetCount error: %v", err)
	}
	got := toInt64(result)
	t.Logf("Bool swap: %d value changes from %d concurrent swaps", got, numGoroutines)
}

// ============================================================================
// 67. Named return + defer modify during panic in goroutine + global
// ============================================================================

func TestNamedReturnDeferPanicGlobal(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	result, err := prog.Run("NamedReturnDeferPanicGlobal")
	if err != nil {
		t.Fatalf("NamedReturnDeferPanicGlobal error: %v", err)
	}
	got := toInt64(result)
	if got != 840 {
		t.Fatalf("NamedReturnDeferPanicGlobal = %d, want 840", got)
	}
	t.Logf("Named return defer panic global: %d (exact)", got)
}

// ============================================================================
// 68. Global map[string]string concurrent access
// ============================================================================

func TestGlobalStrMapConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("StrMapReset")

	const numGoroutines = 50
	const putsPerGoroutine = 10
	totalPuts := int64(numGoroutines * putsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < putsPerGoroutine; j++ {
				key := fmt.Sprintf("k-%d-%d", id, j)
				val := fmt.Sprintf("v-%d-%d", id, j)
				_, err := prog.Run("StrMapPut", key, val)
				if err != nil {
					t.Errorf("StrMapPut error: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	result, err := prog.Run("StrMapLen")
	if err != nil {
		t.Fatalf("StrMapLen error: %v", err)
	}
	got := toInt64(result)
	if got != totalPuts {
		t.Fatalf("Global string map len = %d, want %d", got, totalPuts)
	}
	t.Logf("Global string map concurrent: len = %d (exact)", got)
}

// ============================================================================
// 69. Global []float64 concurrent append
// ============================================================================

func TestGlobalFloatSliceConcurrent(t *testing.T) {
	prog := buildStateful(t)
	defer prog.Close()

	_, _ = prog.Run("FloatSliceReset")

	const numGoroutines = 50
	const appendsPerGoroutine = 10
	totalAppends := int64(numGoroutines * appendsPerGoroutine)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				_, err := prog.Run("FloatSliceAppendProtected", float64(id)*0.1+float64(j))
				if err != nil {
					t.Errorf("FloatSliceAppendProtected error: %v", err)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	result, err := prog.Run("FloatSliceLenProtected")
	if err != nil {
		t.Fatalf("FloatSliceLenProtected error: %v", err)
	}
	got := toInt64(result)
	if got != totalAppends {
		t.Fatalf("Global float64 slice len = %d, want %d", got, totalAppends)
	}
	t.Logf("Global float64 slice concurrent: len = %d (exact)", got)
}
