package tests

import (
	"context"
	_ "embed"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig"
	"github.com/t04dJ14n9/gig/model/value"
	"github.com/t04dJ14n9/gig/tests/testdata/concurrent_stateful"
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
	default:
		return 0
	}
}

// buildStateful compiles the concurrent_stateful source with stateful globals.
func buildStateful(t *testing.T) *gig.Program {
	t.Helper()
	prog, err := gig.Build(concurrentStatefulSrc, gig.WithStatefulGlobals())
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
		i := i
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
		i := i
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
		i := i
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
		i := i
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
		i := i
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
