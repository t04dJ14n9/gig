package tests

import (
	"context"
	_ "embed"
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
