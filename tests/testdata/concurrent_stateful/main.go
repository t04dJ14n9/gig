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
// Both value-type (sync.Mutex) and pointer-type (*sync.Mutex) globals work
// correctly. The interpreter heap-allocates value-type struct globals via
// reflect.New(T), so all method calls operate on the same underlying object.
// See docs/concurrent-globals.md for details.
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

// ============================================================================
// Value-type sync.Mutex (not pointer) — tests heap-allocation fix
//
// These use 'var mu sync.Mutex' (value type) not '*sync.Mutex'.
// The interpreter heap-allocates these via reflect.New(T), so they
// behave identically to pointer-form in concurrent scenarios.
// ============================================================================

var (
	valueMu               sync.Mutex
	valueProtectedCounter int
)

func init() {
	valueProtectedCounter = 0
}

// ValueTypeIncrement increments using value-type sync.Mutex.
func ValueTypeIncrement() int {
	valueMu.Lock()
	valueProtectedCounter = valueProtectedCounter + 1
	val := valueProtectedCounter
	valueMu.Unlock()
	return val
}

// ValueTypeGet returns the counter protected by value-type mutex.
func ValueTypeGet() int {
	valueMu.Lock()
	val := valueProtectedCounter
	valueMu.Unlock()
	return val
}

// ============================================================================
// sync.RWMutex — multiple readers, single writer
// ============================================================================

var (
	rwMu      sync.RWMutex
	rwCounter int
)

func init() {
	rwCounter = 0
}

// RWMutexWrite increments with write lock.
func RWMutexWrite() int {
	rwMu.Lock()
	rwCounter++
	val := rwCounter
	rwMu.Unlock()
	return val
}

// RWMutexRead reads with read lock.
func RWMutexRead() int {
	rwMu.RLock()
	val := rwCounter
	rwMu.RUnlock()
	return val
}

// ============================================================================
// sync.Once — exactly-once initialization (anonymous closure)
// ============================================================================

var onceForTest sync.Once
var onceValue int

// OnceInit returns the once-initialized value.
// Tests that sync.Once.Do with an anonymous closure correctly writes
// to a global variable — the closure must operate on the shared globals.
func OnceInit() int {
	onceForTest.Do(func() {
		onceValue = 42
	})
	return onceValue
}

// ============================================================================
// sync.WaitGroup — goroutine synchronization inside guest code
// ============================================================================

// WaitGroupSum uses a local WaitGroup to synchronize goroutines spawned
// inside the guest program using a for-loop with go func(n int) pattern.
func WaitGroupSum() int {
	const N = 50
	var wg sync.WaitGroup
	ch := make(chan int, N)
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(n int) {
			ch <- n
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// ============================================================================
// sync.Map — concurrent-safe map
// ============================================================================

var concurrentMap sync.Map

// MapStore stores a key-value pair.
func MapStore(key string, value int) {
	concurrentMap.Store(key, value)
}

// MapLoad loads a value by key.
func MapLoad(key string) (int, bool) {
	v, ok := concurrentMap.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int), true
}

// MapLoadOrStore loads or stores a value.
func MapLoadOrStore(key string, value int) (int, bool) {
	v, loaded := concurrentMap.LoadOrStore(key, value)
	return v.(int), loaded
}

// MapDelete deletes a key.
func MapDelete(key string) {
	concurrentMap.Delete(key)
}

// ============================================================================
// Nested locks — lock ordering to prevent deadlock
// ============================================================================

var (
	nestedMuA sync.Mutex
	nestedMuB sync.Mutex
	nestedSum int
)

func init() {
	nestedSum = 0
}

// NestedLockAB locks A then B (consistent ordering).
func NestedLockAB() int {
	nestedMuA.Lock()
	nestedMuB.Lock()
	nestedSum = nestedSum + 1
	val := nestedSum
	nestedMuB.Unlock()
	nestedMuA.Unlock()
	return val
}

// NestedLockBA locks B then A (reverse ordering — potential deadlock).
// This tests that the interpreter handles lock ordering correctly.
func NestedLockBA() int {
	nestedMuB.Lock()
	nestedMuA.Lock()
	nestedSum = nestedSum + 1
	val := nestedSum
	nestedMuA.Unlock()
	nestedMuB.Unlock()
	return val
}

// ============================================================================
// Complex mixed synchronization — channel-based signal
// ============================================================================

var (
	complexMu      sync.Mutex
	complexReady   bool
	complexResult  int
	complexSignal  chan bool
)

func init() {
	complexSignal = make(chan bool, 1)
}

// ComplexProducer sets the result and signals via channel.
func ComplexProducer(value int) {
	complexMu.Lock()
	complexResult = value
	complexReady = true
	complexMu.Unlock()
	complexSignal <- true
}

// ComplexConsumer waits for the signal and returns the result.
func ComplexConsumer() int {
	<-complexSignal
	complexMu.Lock()
	val := complexResult
	complexMu.Unlock()
	return val
}

// ResetComplexState resets for testing.
func ResetComplexState() {
	complexMu.Lock()
	complexReady = false
	complexResult = 0
	select {
	case <-complexSignal:
	default:
	}
	complexMu.Unlock()
}
