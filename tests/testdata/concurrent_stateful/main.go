package concurrent_stateful

import "sync"

// ============================================================================
// Unprotected globals (race-tolerant tests — verify no crash, not exact value)
// ============================================================================

var unprotectedCounter int

func init() {
	unprotectedCounter = 0
}

// IncrementUnprotected does a non-atomic read-modify-write on a global.
// Under concurrent execution, lost updates are expected (same as Go).
func IncrementUnprotected() int {
	unprotectedCounter = unprotectedCounter + 1
	return unprotectedCounter
}

// GetUnprotected returns the current unprotected counter value.
func GetUnprotected() int {
	return unprotectedCounter
}

// ============================================================================
// Mutex-protected globals (exact correctness under concurrency)
//
// NOTE: Global struct variables with pointer-receiver methods (like sync.Mutex)
// must be stored as pointers (*sync.Mutex), not values (sync.Mutex), because
// the interpreter stores globals in value.Value slots, and value-type globals
// cannot be addressed for pointer-receiver method calls.
// ============================================================================

var (
	mu               *sync.Mutex
	protectedCounter int
)

func init() {
	mu = &sync.Mutex{}
	protectedCounter = 0
}

// IncrementProtected increments a mutex-protected global counter.
// With proper locking, the final sum MUST be exactly N after N calls.
func IncrementProtected() int {
	mu.Lock()
	protectedCounter = protectedCounter + 1
	val := protectedCounter
	mu.Unlock()
	return val
}

// GetProtected returns the current protected counter value.
func GetProtected() int {
	mu.Lock()
	val := protectedCounter
	mu.Unlock()
	return val
}

// ============================================================================
// Channel-based accumulation (exact correctness, no shared mutable state)
// ============================================================================

// SumViaChannel spawns N goroutines, each sending 1 to a channel,
// then sums all received values. Result must be exactly N.
func SumViaChannel() int {
	const N = 100
	ch := make(chan int, N)
	for i := 0; i < N; i++ {
		go func() {
			ch <- 1
		}()
	}
	sum := 0
	for i := 0; i < N; i++ {
		sum += <-ch
	}
	return sum
}

// ============================================================================
// Read-only global (concurrent reads must all see the same value)
// ============================================================================

var greeting string

func init() {
	greeting = "hello"
}

// GetGreeting returns a read-only global. All concurrent reads must match.
func GetGreeting() string {
	return greeting
}

// ============================================================================
// State mutation visibility (set then get across calls)
// ============================================================================

var (
	stateMu *sync.Mutex
	state   int
)

func init() {
	stateMu = &sync.Mutex{}
	state = 0
}

// SetState writes a value to the global state variable.
func SetState(v int) {
	stateMu.Lock()
	state = v
	stateMu.Unlock()
}

// GetState reads the current state.
func GetState() int {
	stateMu.Lock()
	val := state
	stateMu.Unlock()
	return val
}

// ============================================================================
// Pure function (no globals, stateless — verify concurrent RunWithValues)
// ============================================================================

// Add returns a + b. Stateless, safe for concurrent execution.
func Add(a, b int) int {
	return a + b
}

// Multiply returns a * b.
func Multiply(a, b int) int {
	return a * b
}

// ============================================================================
// Producer-consumer pattern with buffered channels
// ============================================================================

// ProducerConsumerSum spawns producers that send values to a channel
// and a consumer that sums them. Tests goroutine coordination correctness.
func ProducerConsumerSum() int {
	ch := make(chan int, 50)
	done := make(chan bool)

	// 10 producers, each sends values 1..10
	for p := 0; p < 10; p++ {
		go func() {
			for i := 1; i <= 10; i++ {
				ch <- i
			}
			done <- true
		}()
	}

	// Wait for all producers via done channel
	go func() {
		for i := 0; i < 10; i++ {
			<-done
		}
		close(ch)
	}()

	sum := 0
	for v := range ch {
		sum += v
	}
	// 10 producers * sum(1..10) = 10 * 55 = 550
	return sum
}

// ============================================================================
// Multiple mutex-protected globals (independent locks)
// ============================================================================

var (
	muA    *sync.Mutex
	countA int
	muB    *sync.Mutex
	countB int
)

func init() {
	muA = &sync.Mutex{}
	countA = 0
	muB = &sync.Mutex{}
	countB = 0
}

// IncrementA increments counter A with its own lock.
func IncrementA() int {
	muA.Lock()
	countA = countA + 1
	val := countA
	muA.Unlock()
	return val
}

// IncrementB increments counter B with its own lock.
func IncrementB() int {
	muB.Lock()
	countB = countB + 1
	val := countB
	muB.Unlock()
	return val
}

// GetCountA returns counter A.
func GetCountA() int {
	muA.Lock()
	val := countA
	muA.Unlock()
	return val
}

// GetCountB returns counter B.
func GetCountB() int {
	muB.Lock()
	val := countB
	muB.Unlock()
	return val
}
