package concurrent_stateful

import (
	"sync"
)

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

// ============================================================================
// Atomic-style counter with Mutex — sequential consistency
// ============================================================================

var (
	atomicMu    sync.Mutex
	atomicCount int64
)

func init() {
	atomicCount = 0
}

// AtomicAdd adds delta to the atomic-style counter.
func AtomicAdd(delta int64) int64 {
	atomicMu.Lock()
	atomicCount += delta
	val := atomicCount
	atomicMu.Unlock()
	return val
}

// AtomicGet returns the atomic-style counter value.
func AtomicGet() int64 {
	atomicMu.Lock()
	val := atomicCount
	atomicMu.Unlock()
	return val
}

// AtomicSet sets the atomic-style counter value.
func AtomicSet(val int64) {
	atomicMu.Lock()
	atomicCount = val
	atomicMu.Unlock()
}

// ============================================================================
// Goroutine inside goroutine (nested goroutine spawning)
// ============================================================================

// NestedGoroutineSum spawns goroutines that spawn more goroutines.
// Outer: 5 goroutines, each spawns 10 inner goroutines sending 1.
func NestedGoroutineSum() int {
	ch := make(chan int, 50)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				go func() {
					ch <- 1
				}()
			}
		}()
	}
	wg.Wait()
	// Give time for inner goroutines to complete
	sum := 0
	for i := 0; i < 50; i++ {
		sum += <-ch
	}
	return sum
}

// ============================================================================
// Global slice with mutex protection
// ============================================================================

var (
	sliceMu   sync.Mutex
	globalBuf []int
)

func init() {
	globalBuf = make([]int, 0)
}

// AppendProtected appends a value to the global slice with mutex.
func AppendProtected(v int) int {
	sliceMu.Lock()
	globalBuf = append(globalBuf, v)
	len_ := len(globalBuf)
	sliceMu.Unlock()
	return len_
}

// GetBufLen returns the length of the global slice.
func GetBufLen() int {
	sliceMu.Lock()
	len_ := len(globalBuf)
	sliceMu.Unlock()
	return len_
}

// ResetBuf resets the global slice.
func ResetBuf() {
	sliceMu.Lock()
	globalBuf = make([]int, 0)
	sliceMu.Unlock()
}

// ============================================================================
// Global map with mutex protection
// ============================================================================

var (
	mapMu     sync.Mutex
	globalMap map[string]int
)

func init() {
	globalMap = make(map[string]int)
}

// MapPutProtected puts a key-value pair with mutex.
func MapPutProtected(key string, val int) {
	mapMu.Lock()
	globalMap[key] = val
	mapMu.Unlock()
}

// MapGetProtected gets a value by key with mutex.
func MapGetProtected(key string) (int, bool) {
	mapMu.Lock()
	v, ok := globalMap[key]
	mapMu.Unlock()
	return v, ok
}

// MapLenProtected returns the map size with mutex.
func MapLenProtected() int {
	mapMu.Lock()
	len_ := len(globalMap)
	mapMu.Unlock()
	return len_
}

// ResetProtectedMap resets the global map.
func ResetProtectedMap() {
	mapMu.Lock()
	globalMap = make(map[string]int)
	mapMu.Unlock()
}

// ============================================================================
// Channel-based barrier pattern
// ============================================================================

// BarrierSum uses a channel as a barrier to wait for all goroutines.
func BarrierSum() int {
	const N = 20
	ch := make(chan int, N)
	for i := 0; i < N; i++ {
		go func(n int) {
			ch <- n * n
		}(i)
	}
	sum := 0
	for i := 0; i < N; i++ {
		sum += <-ch
	}
	return sum
}

// ============================================================================
// Multiple sync.Once instances
// ============================================================================

var (
	onceA    sync.Once
	onceB    sync.Once
	onceValA int
	onceValB int
)

func init() {
	onceValA = 0
	onceValB = 0
}

// OnceInitA initializes value A exactly once.
func OnceInitA() int {
	onceA.Do(func() {
		onceValA = 100
	})
	return onceValA
}

// OnceInitB initializes value B exactly once.
func OnceInitB() int {
	onceB.Do(func() {
		onceValB = 200
	})
	return onceValB
}

// ============================================================================
// Defer in goroutine with global state
// ============================================================================

var (
	deferMu    sync.Mutex
	deferCount int
)

func init() {
	deferCount = 0
}

// DeferInGoroutine runs N goroutines that use defer to increment a global.
func DeferInGoroutine() int {
	const N = 30
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				deferMu.Lock()
				deferCount++
				deferMu.Unlock()
			}()
		}()
	}
	wg.Wait()
	return deferCount
}

// ============================================================================
// Context propagation — goroutine with timeout
// ============================================================================

// GoroutineWithResult uses a result channel pattern.
func GoroutineWithResult() int {
	result := make(chan int, 1)
	go func() {
		result <- 42
	}()
	return <-result
}

// ============================================================================
// Bidirectional channel with goroutines
// ============================================================================

// BidirectionalChannel uses both send and receive in goroutines.
func BidirectionalChannel() int {
	ch := make(chan int)
	go func() {
		v := <-ch
		ch <- v * 2
	}()
	ch <- 21
	return <-ch
}

// ============================================================================
// Select-like pattern using multiple channels
// ============================================================================

// MultiChannelMerge merges values from multiple channels.
func MultiChannelMerge() int {
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	ch3 := make(chan int, 5)

	for i := 0; i < 5; i++ {
		ch1 <- i + 1
		ch2 <- (i + 1) * 10
		ch3 <- (i + 1) * 100
	}

	sum := 0
	for i := 0; i < 5; i++ {
		sum += <-ch1
		sum += <-ch2
		sum += <-ch3
	}
	return sum // 15 + 150 + 1500 = 1665
}

// ============================================================================
// Value-type RWMutex (not pointer)
// ============================================================================

var (
	valueRWMu  sync.RWMutex
	valueRWVal int
)

func init() {
	valueRWVal = 0
}

// ValueTypeRWWrite writes with value-type RWMutex.
func ValueTypeRWWrite(v int) {
	valueRWMu.Lock()
	valueRWVal = v
	valueRWMu.Unlock()
}

// ValueTypeRWRead reads with value-type RWMutex.
func ValueTypeRWRead() int {
	valueRWMu.RLock()
	v := valueRWVal
	valueRWMu.RUnlock()
	return v
}

// ============================================================================
// Global string with mutex protection
// ============================================================================

var (
	strMu    sync.RWMutex
	globalStr string
)

func init() {
	globalStr = "initial"
}

// SetGlobalString sets the global string.
func SetGlobalString(s string) {
	strMu.Lock()
	globalStr = s
	strMu.Unlock()
}

// GetGlobalString gets the global string.
func GetGlobalString() string {
	strMu.RLock()
	s := globalStr
	strMu.RUnlock()
	return s
}

// ============================================================================
// Counter with swap (compare-and-swap pattern)
// ============================================================================

var (
	casMu    sync.Mutex
	casValue int
)

func init() {
	casValue = 0
}

// CASIncrement increments using compare-and-swap pattern.
func CASIncrement() int {
	casMu.Lock()
	casValue++
	v := casValue
	casMu.Unlock()
	return v
}

// CASSwap sets value if current equals expected (simulated CAS).
func CASSwap(expected, newValue int) bool {
	casMu.Lock()
	defer casMu.Unlock()
	if casValue == expected {
		casValue = newValue
		return true
	}
	return false
}

// CASGet returns current CAS value.
func CASGet() int {
	casMu.Lock()
	v := casValue
	casMu.Unlock()
	return v
}

// ============================================================================
// Concurrent global bool flag
// ============================================================================

var (
	flagMu  sync.Mutex
	flagVal bool
)

func init() {
	flagVal = false
}

// SetFlag sets the global boolean flag.
func SetFlag(v bool) {
	flagMu.Lock()
	flagVal = v
	flagMu.Unlock()
}

// GetFlag gets the global boolean flag.
func GetFlag() bool {
	flagMu.Lock()
	v := flagVal
	flagMu.Unlock()
	return v
}

// ============================================================================
// Closure capturing globals — closures that read/write package-level vars
// ============================================================================

var closureGlobalVal int

func init() {
	closureGlobalVal = 0
}

// ClosureReadGlobal returns a closure that reads a global variable.
func ClosureReadGlobal() int {
	fn := func() int {
		return closureGlobalVal
	}
	closureGlobalVal = 42
	return fn()
}

// ClosureWriteGlobal uses a closure to write to a global variable.
func ClosureWriteGlobal() int {
	fn := func(v int) {
		closureGlobalVal = v
	}
	fn(99)
	return closureGlobalVal
}

// ClosureAccumulateGlobal uses a closure to accumulate into a global.
func ClosureAccumulateGlobal() int {
	closureGlobalVal = 0
	adder := func(n int) {
		closureGlobalVal = closureGlobalVal + n
	}
	adder(10)
	adder(20)
	adder(30)
	return closureGlobalVal
}

// MultipleClosuresSharedGlobal returns multiple closures sharing a global.
func MultipleClosuresSharedGlobal() int {
	closureGlobalVal = 0
	set := func(v int) { closureGlobalVal = v }
	get := func() int { return closureGlobalVal }
	add := func(n int) { closureGlobalVal = closureGlobalVal + n }

	set(10)
	add(5)
	return get()
}

// ============================================================================
// Recursive closures — var f func(); f = func() { ... f() ... }
// ============================================================================

// RecursiveClosureFib computes fibonacci via recursive closure.
func RecursiveClosureFib() int {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fib(10)
}

// RecursiveClosureFactorial computes factorial via recursive closure.
func RecursiveClosureFactorial() int {
	var fact func(n int) int
	fact = func(n int) int {
		if n <= 1 {
			return 1
		}
		return n * fact(n-1)
	}
	return fact(8)
}

// RecursiveClosureWithGlobal uses a recursive closure that reads a global.
func RecursiveClosureWithGlobal() int {
	closureGlobalVal = 0
	var visit func(n int)
	visit = func(n int) {
		if n <= 0 {
			return
		}
		closureGlobalVal = closureGlobalVal + n
		visit(n - 1)
	}
	visit(10)
	return closureGlobalVal // 10+9+...+1 = 55
}

// ============================================================================
// Goroutines inside guest code accessing shared globals
// ============================================================================

var (
	goroutineMu    sync.Mutex
	goroutineCount int
)

func init() {
	goroutineCount = 0
}

// GoroutineIncrementGlobal spawns goroutines that increment a mutex-protected global.
func GoroutineIncrementGlobal() int {
	goroutineCount = 0
	const N = 50
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			goroutineMu.Lock()
			goroutineCount++
			goroutineMu.Unlock()
		}()
	}
	wg.Wait()
	return goroutineCount
}

// GoroutineReadGlobal spawns goroutines that read a global concurrently.
func GoroutineReadGlobal() int {
	goroutineCount = 42
	const N = 50
	ok := int64(0)
	ch := make(chan int, N)
	for i := 0; i < N; i++ {
		go func() {
			goroutineMu.Lock()
			v := goroutineCount
			goroutineMu.Unlock()
			ch <- v
		}()
	}
	for i := 0; i < N; i++ {
		if <-ch == 42 {
			ok++
		}
	}
	return int(ok)
}

// GoroutineClosureGlobal spawns goroutines with closures that access a global.
func GoroutineClosureGlobal() int {
	goroutineCount = 0
	const N = 30
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		inc := func() {
			goroutineMu.Lock()
			goroutineCount++
			goroutineMu.Unlock()
		}
		go func() {
			defer wg.Done()
			inc()
		}()
	}
	wg.Wait()
	return goroutineCount
}

// ============================================================================
// Diverse global types under concurrent access
// ============================================================================

var (
	globalFloatMu  sync.Mutex
	globalFloatVal float64
	globalStrMu    sync.RWMutex
	globalStrVal   string
	globalBoolMu   sync.Mutex
	globalBoolVal  bool
)

func init() {
	globalFloatVal = 0.0
	globalStrVal = ""
	globalBoolVal = false
}

// FloatGlobalIncrement increments a float64 global with mutex protection.
func FloatGlobalIncrement() float64 {
	globalFloatMu.Lock()
	globalFloatVal += 1.0
	v := globalFloatVal
	globalFloatMu.Unlock()
	return v
}

// FloatGlobalGet returns the float64 global.
func FloatGlobalGet() float64 {
	globalFloatMu.Lock()
	v := globalFloatVal
	globalFloatMu.Unlock()
	return v
}

// StringGlobalSet sets the string global.
func StringGlobalSet(s string) {
	globalStrMu.Lock()
	globalStrVal = s
	globalStrMu.Unlock()
}

// StringGlobalGet returns the string global.
func StringGlobalGet() string {
	globalStrMu.RLock()
	v := globalStrVal
	globalStrMu.RUnlock()
	return v
}

// BoolGlobalSet sets the bool global.
func BoolGlobalSet(v bool) {
	globalBoolMu.Lock()
	globalBoolVal = v
	globalBoolMu.Unlock()
}

// BoolGlobalGet returns the bool global.
func BoolGlobalGet() bool {
	globalBoolMu.Lock()
	v := globalBoolVal
	globalBoolMu.Unlock()
	return v
}

// ============================================================================
// Closure in goroutine with global: closure captures loop variable, writes to global
// ============================================================================

var (
	loopClosureMu    sync.Mutex
	loopClosureCount int
)

func init() {
	loopClosureCount = 0
}

// LoopClosureGoroutineGlobal spawns goroutines in a loop, each closure
// captures its iteration variable and writes to a global.
// This tests the Go 1.22+ per-iteration variable semantics with global mutation.
func LoopClosureGoroutineGlobal() int {
	loopClosureCount = 0
	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(n int) {
			defer wg.Done()
			loopClosureMu.Lock()
			loopClosureCount += n
			loopClosureMu.Unlock()
		}(i)
	}
	wg.Wait()
	return loopClosureCount // sum(0..19) = 190
}

// ============================================================================
// Global int64 with mutex — large integer concurrent access
// ============================================================================

var (
	int64Mu    sync.Mutex
	globalInt64 int64
)

func init() {
	globalInt64 = 0
}

// Int64Increment increments a global int64 with mutex.
func Int64Increment() int64 {
	int64Mu.Lock()
	globalInt64++
	v := globalInt64
	int64Mu.Unlock()
	return v
}

// Int64Get returns the global int64 value.
func Int64Get() int64 {
	int64Mu.Lock()
	v := globalInt64
	int64Mu.Unlock()
	return v
}

// Int64Set sets the global int64 value.
func Int64Set(v int64) {
	int64Mu.Lock()
	globalInt64 = v
	int64Mu.Unlock()
}

// ============================================================================
// Global map concurrent read-only — multiple readers, no writers
// ============================================================================

var (
	roMapMu sync.RWMutex
	roMap   map[string]int
)

func init() {
	roMap = map[string]int{
		"alpha": 1, "bravo": 2, "charlie": 3,
		"delta": 4, "echo": 5, "foxtrot": 6,
		"golf": 7, "hotel": 8, "india": 9, "juliet": 10,
	}
}

// ROMapGet retrieves a key from the read-only map.
func ROMapGet(key string) (int, bool) {
	roMapMu.RLock()
	v, ok := roMap[key]
	roMapMu.RUnlock()
	return v, ok
}

// ROMapLen returns the number of entries in the read-only map.
func ROMapLen() int {
	roMapMu.RLock()
	n := len(roMap)
	roMapMu.RUnlock()
	return n
}

// ============================================================================
// Defer recover in goroutine with global state
// ============================================================================

var (
	panicRecoverMu    sync.Mutex
	panicRecoverCount int
)

func init() {
	panicRecoverCount = 0
}

// DeferRecoverInGoroutine spawns goroutines that panic and recover,
// incrementing a global counter in the deferred recovery function.
func DeferRecoverInGoroutine() int {
	panicRecoverCount = 0
	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicRecoverMu.Lock()
					panicRecoverCount++
					panicRecoverMu.Unlock()
				}
			}()
			panic("test")
		}()
	}
	wg.Wait()
	return panicRecoverCount
}

// ============================================================================
// Goroutine ping-pong via channels with global state tracking
// ============================================================================

var (
	pingPongMu    sync.Mutex
	pingPongCount int
)

func init() {
	pingPongCount = 0
}

// PingPongGlobal uses two goroutines that ping-pong through channels
// and increment a global counter each round.
func PingPongGlobal() int {
	pingPongCount = 0
	const rounds = 10
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Ping goroutine
	go func() {
		for i := 0; i < rounds; i++ {
			<-ch1 // wait for signal
			pingPongMu.Lock()
			pingPongCount++
			pingPongMu.Unlock()
			ch2 <- i // signal back
		}
	}()

	// Pong goroutine (main)
	for i := 0; i < rounds; i++ {
		ch1 <- i // signal
		<-ch2    // wait for response
		pingPongMu.Lock()
		pingPongCount++
		pingPongMu.Unlock()
	}

	return pingPongCount // 20 (10 ping + 10 pong)
}

// ============================================================================
// Concurrent global uint access
// ============================================================================

var (
	uintMu     sync.Mutex
	globalUint uint
)

func init() {
	globalUint = 0
}

// UintIncrement increments a global uint with mutex.
func UintIncrement() uint {
	uintMu.Lock()
	globalUint++
	v := globalUint
	uintMu.Unlock()
	return v
}

// UintGet returns the global uint value.
func UintGet() uint {
	uintMu.Lock()
	v := globalUint
	uintMu.Unlock()
	return v
}

// ============================================================================
// Select with global state — goroutines communicate via select
// ============================================================================

var (
	selectMu    sync.Mutex
	selectCount int
)

func init() {
	selectCount = 0
}

// SelectIncrementGlobal uses select to receive from one of two channels
// and increment a global counter.
func SelectIncrementGlobal() int {
	selectCount = 0
	const N = 20
	ch1 := make(chan int, N)
	ch2 := make(chan int, N)

	// Fill channels
	for i := 0; i < N; i++ {
		ch1 <- i
		ch2 <- i + 100
	}

	// Read from both channels via select
	for i := 0; i < N*2; i++ {
		select {
		case <-ch1:
			selectMu.Lock()
			selectCount++
			selectMu.Unlock()
		case <-ch2:
			selectMu.Lock()
			selectCount++
			selectMu.Unlock()
		}
	}
	return selectCount
}

// ============================================================================
// Global struct pointer — concurrent access to struct fields
// ============================================================================

type SharedConfig struct {
	Enabled  bool
	Count    int
	Name     string
	Priority float64
}

var (
	configMu sync.Mutex
	config   *SharedConfig
)

func init() {
	config = &SharedConfig{Enabled: true, Count: 0, Name: "default", Priority: 1.0}
}

// ConfigIncrementCount increments the config's Count field.
func ConfigIncrementCount() int {
	configMu.Lock()
	config.Count++
	v := config.Count
	configMu.Unlock()
	return v
}

// ConfigGetCount returns the config's Count field.
func ConfigGetCount() int {
	configMu.Lock()
	v := config.Count
	configMu.Unlock()
	return v
}

// ConfigSetName sets the config's Name field.
func ConfigSetName(name string) {
	configMu.Lock()
	config.Name = name
	configMu.Unlock()
}

// ConfigGetName returns the config's Name field.
func ConfigGetName() string {
	configMu.Lock()
	v := config.Name
	configMu.Unlock()
	return v
}

// ConfigReset resets config to default.
func ConfigReset() {
	configMu.Lock()
	config.Enabled = true
	config.Count = 0
	config.Name = "default"
	config.Priority = 1.0
	configMu.Unlock()
}

// ============================================================================
// Nested goroutines with global counter — goroutine spawns goroutine
// ============================================================================

var (
	nestedGoroutineMu    sync.Mutex
	nestedGoroutineCount int
)

func init() {
	nestedGoroutineCount = 0
}

// NestedGoroutineGlobal spawns outer goroutines that spawn inner goroutines,
// each incrementing a global counter.
func NestedGoroutineGlobal() int {
	nestedGoroutineCount = 0
	const outerN = 5
	const innerN = 10
	var outerWg sync.WaitGroup
	outerWg.Add(outerN)
	for i := 0; i < outerN; i++ {
		go func() {
			defer outerWg.Done()
			var innerWg sync.WaitGroup
			innerWg.Add(innerN)
			for j := 0; j < innerN; j++ {
				go func() {
					defer innerWg.Done()
					nestedGoroutineMu.Lock()
					nestedGoroutineCount++
					nestedGoroutineMu.Unlock()
				}()
			}
			innerWg.Wait()
		}()
	}
	outerWg.Wait()
	return nestedGoroutineCount // 5 * 10 = 50
}

// ============================================================================
// Channel close + range with global accumulation
// ============================================================================

var (
	chanRangeMu    sync.Mutex
	chanRangeSum   int
)

func init() {
	chanRangeSum = 0
}

// ChannelCloseRangeGlobal sends values on a channel, closes it, then
// a goroutine ranges over it and accumulates into a global.
func ChannelCloseRangeGlobal() int {
	chanRangeSum = 0
	ch := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {
			chanRangeMu.Lock()
			chanRangeSum += v
			chanRangeMu.Unlock()
		}
	}()
	wg.Wait()
	return chanRangeSum // 1+2+...+10 = 55
}

// ============================================================================
// Global []byte with mutex — concurrent append to byte slice
// ============================================================================

var (
	bytesMu    sync.Mutex
	globalBytes []byte
)

func init() {
	globalBytes = make([]byte, 0)
}

// BytesAppendProtected appends bytes to the global slice with mutex.
func BytesAppendProtected(b byte) int {
	bytesMu.Lock()
	globalBytes = append(globalBytes, b)
	n := len(globalBytes)
	bytesMu.Unlock()
	return n
}

// BytesLenProtected returns the length of the global byte slice.
func BytesLenProtected() int {
	bytesMu.Lock()
	n := len(globalBytes)
	bytesMu.Unlock()
	return n
}

// BytesReset resets the global byte slice.
func BytesReset() {
	bytesMu.Lock()
	globalBytes = make([]byte, 0)
	bytesMu.Unlock()
}

// ============================================================================
// Concurrent interleaved set+get on same global
// ============================================================================

var (
	interleaveMu    sync.Mutex
	interleaveVal   int
	interleaveReads int64
)

func init() {
	interleaveVal = 0
	interleaveReads = 0
}

// InterleaveSet sets the interleaved global value.
func InterleaveSet(v int) {
	interleaveMu.Lock()
	interleaveVal = v
	interleaveMu.Unlock()
}

// InterleaveGet gets the interleaved global value and increments read counter.
func InterleaveGet() int {
	interleaveMu.Lock()
	v := interleaveVal
	interleaveReads++
	interleaveMu.Unlock()
	return v
}

// InterleaveReadCount returns the number of reads performed.
func InterleaveReadCount() int64 {
	interleaveMu.Lock()
	n := interleaveReads
	interleaveMu.Unlock()
	return n
}

// InterleaveReset resets the interleaved state.
func InterleaveReset() {
	interleaveMu.Lock()
	interleaveVal = 0
	interleaveReads = 0
	interleaveMu.Unlock()
}

// ============================================================================
// Global complex128 with mutex — concurrent access
// ============================================================================

var (
	complexGlobalMu sync.Mutex
	complexGlobal   complex128
)

func init() {
	complexGlobal = 0 + 0i
}

// ComplexGlobalAdd adds a value to the global complex128.
func ComplexGlobalAdd(re, im float64) complex128 {
	complexGlobalMu.Lock()
	complexGlobal = complexGlobal + complex(re, im)
	v := complexGlobal
	complexGlobalMu.Unlock()
	return v
}

// ComplexGlobalGet returns the global complex128 value.
func ComplexGlobalGet() complex128 {
	complexGlobalMu.Lock()
	v := complexGlobal
	complexGlobalMu.Unlock()
	return v
}

// ComplexGlobalReset resets the global complex128.
func ComplexGlobalReset() {
	complexGlobalMu.Lock()
	complexGlobal = 0 + 0i
	complexGlobalMu.Unlock()
}

// ============================================================================
// Global int32 with mutex — concurrent access
// ============================================================================

var (
	int32Mu    sync.Mutex
	globalInt32 int32
)

func init() {
	globalInt32 = 0
}

// Int32Increment increments a global int32 with mutex.
func Int32Increment() int32 {
	int32Mu.Lock()
	globalInt32++
	v := globalInt32
	int32Mu.Unlock()
	return v
}

// Int32Get returns the global int32 value.
func Int32Get() int32 {
	int32Mu.Lock()
	v := globalInt32
	int32Mu.Unlock()
	return v
}

// Int32Set sets the global int32 value.
func Int32Set(v int32) {
	int32Mu.Lock()
	globalInt32 = v
	int32Mu.Unlock()
}

// ============================================================================
// Struct method that accesses a global
// ============================================================================

type GlobalAccessor struct {
	label string
}

// IncrementGlobal increments a global counter (method on struct).
func (g *GlobalAccessor) IncrementGlobal() int {
	goroutineMu.Lock()
	goroutineCount++
	v := goroutineCount
	goroutineMu.Unlock()
	return v
}

// GetGlobalCount returns the global counter value (method on struct).
func (g *GlobalAccessor) GetGlobalCount() int {
	goroutineMu.Lock()
	v := goroutineCount
	goroutineMu.Unlock()
	return v
}

var accessor = &GlobalAccessor{label: "test"}

// StructMethodIncrementGlobal uses a struct method to increment a global.
func StructMethodIncrementGlobal() int {
	goroutineCount = 0
	const N = 30
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			accessor.IncrementGlobal()
		}()
	}
	wg.Wait()
	return accessor.GetGlobalCount()
}

// ============================================================================
// Fan-out pattern — one producer, multiple consumers via channels
// ============================================================================

var (
	fanoutMu    sync.Mutex
	fanoutSum   int
)

func init() {
	fanoutSum = 0
}

// FanOutGlobal produces a value and fans out to N consumer goroutines
// that each add to a global sum.
func FanOutGlobal() int {
	fanoutSum = 0
	const N = 10
	ch := make(chan int, N)

	// Producer: send value to all consumers
	producer := func() {
		for i := 0; i < N; i++ {
			ch <- i + 1
		}
		close(ch)
	}
	go producer()

	// Consumers: read from channel and add to global
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			if v, ok := <-ch; ok {
				fanoutMu.Lock()
				fanoutSum += v
				fanoutMu.Unlock()
			}
		}()
	}
	wg.Wait()
	return fanoutSum // 1+2+...+10 = 55
}

// ============================================================================
// Global float32 with mutex — concurrent access
// ============================================================================

var (
	float32Mu     sync.Mutex
	globalFloat32 float32
)

func init() {
	globalFloat32 = 0.0
}

// Float32Increment increments a global float32 with mutex.
func Float32Increment() float32 {
	float32Mu.Lock()
	globalFloat32 += 1.0
	v := globalFloat32
	float32Mu.Unlock()
	return v
}

// Float32Get returns the global float32 value.
func Float32Get() float32 {
	float32Mu.Lock()
	v := globalFloat32
	float32Mu.Unlock()
	return v
}

// ============================================================================
// Global int8/int16 with mutex — concurrent access
// ============================================================================

var (
	int8Mu     sync.Mutex
	globalInt8  int8
	int16Mu    sync.Mutex
	globalInt16 int16
)

func init() {
	globalInt8 = 0
	globalInt16 = 0
}

// Int8Increment increments a global int8 with mutex.
func Int8Increment() int8 {
	int8Mu.Lock()
	globalInt8++
	v := globalInt8
	int8Mu.Unlock()
	return v
}

// Int8Get returns the global int8 value.
func Int8Get() int8 {
	int8Mu.Lock()
	v := globalInt8
	int8Mu.Unlock()
	return v
}

// Int8Set sets the global int8 value.
func Int8Set(v int8) {
	int8Mu.Lock()
	globalInt8 = v
	int8Mu.Unlock()
}

// Int16Increment increments a global int16 with mutex.
func Int16Increment() int16 {
	int16Mu.Lock()
	globalInt16++
	v := globalInt16
	int16Mu.Unlock()
	return v
}

// Int16Get returns the global int16 value.
func Int16Get() int16 {
	int16Mu.Lock()
	v := globalInt16
	int16Mu.Unlock()
	return v
}

// Int16Set sets the global int16 value.
func Int16Set(v int16) {
	int16Mu.Lock()
	globalInt16 = v
	int16Mu.Unlock()
}

// ============================================================================
// Worker pool with result channel and global accumulator
// ============================================================================

var (
	workerPoolMu    sync.Mutex
	workerPoolTotal int
)

func init() {
	workerPoolTotal = 0
}

// WorkerPoolGlobal distributes work to N workers via a channel,
// each worker adds to a global total.
func WorkerPoolGlobal() int {
	workerPoolTotal = 0
	const numWorkers = 5
	const numJobs = 20
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- j * j
			}
		}()
	}

	for j := 0; j < numJobs; j++ {
		jobs <- j + 1
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		workerPoolMu.Lock()
		workerPoolTotal += r
		workerPoolMu.Unlock()
	}

	return workerPoolTotal // sum of 1^2 + ... + 20^2 = 2870
}

// ============================================================================
// Global uintptr with mutex — concurrent access
// ============================================================================

var (
	uintptrMu     sync.Mutex
	globalUintptr uintptr
)

func init() {
	globalUintptr = 0
}

// UintptrIncrement increments a global uintptr with mutex.
func UintptrIncrement() uintptr {
	uintptrMu.Lock()
	globalUintptr++
	v := globalUintptr
	uintptrMu.Unlock()
	return v
}

// UintptrGet returns the global uintptr value.
func UintptrGet() uintptr {
	uintptrMu.Lock()
	v := globalUintptr
	uintptrMu.Unlock()
	return v
}

// ============================================================================
// Global uint32 with mutex — concurrent access
// ============================================================================

var (
	uint32Mu     sync.Mutex
	globalUint32 uint32
)

func init() {
	globalUint32 = 0
}

// Uint32Increment increments a global uint32 with mutex.
func Uint32Increment() uint32 {
	uint32Mu.Lock()
	globalUint32++
	v := globalUint32
	uint32Mu.Unlock()
	return v
}

// Uint32Get returns the global uint32 value.
func Uint32Get() uint32 {
	uint32Mu.Lock()
	v := globalUint32
	uint32Mu.Unlock()
	return v
}

// ============================================================================
// Multiple goroutines writing to different global variables
// ============================================================================

var (
	multiVarMu *sync.Mutex
	multiVarA  int
	multiVarB  int
	multiVarC  int
)

func init() {
	multiVarMu = &sync.Mutex{}
	multiVarA = 0
	multiVarB = 0
	multiVarC = 0
}

// MultiVarIncrementA increments var A.
func MultiVarIncrementA() int {
	multiVarMu.Lock()
	multiVarA++
	v := multiVarA
	multiVarMu.Unlock()
	return v
}

// MultiVarIncrementB increments var B.
func MultiVarIncrementB() int {
	multiVarMu.Lock()
	multiVarB++
	v := multiVarB
	multiVarMu.Unlock()
	return v
}

// MultiVarIncrementC increments var C.
func MultiVarIncrementC() int {
	multiVarMu.Lock()
	multiVarC++
	v := multiVarC
	multiVarMu.Unlock()
	return v
}

// MultiVarGetSum returns the sum of all variables.
func MultiVarGetSum() int {
	multiVarMu.Lock()
	v := multiVarA + multiVarB + multiVarC
	multiVarMu.Unlock()
	return v
}

// MultiVarReset resets all variables.
func MultiVarReset() {
	multiVarMu.Lock()
	multiVarA = 0
	multiVarB = 0
	multiVarC = 0
	multiVarMu.Unlock()
}

// ============================================================================
// Defer external method (mu.Unlock) in goroutine with panic/recover
// Tests that external method defers run during panic recovery in goroutines
// ============================================================================

var (
	deferExtMu    *sync.Mutex
	deferExtCount int
)

func init() {
	deferExtMu = &sync.Mutex{}
	deferExtCount = 0
}

// DeferExtMethodInGoroutine uses defer mu.Unlock() pattern with panic/recover.
// This tests that external method defers (OpDeferExternal) run correctly
// during panic recovery inside goroutines.
func DeferExtMethodInGoroutine() int {
	deferExtCount = 0
	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			deferExtMu.Lock()
			defer deferExtMu.Unlock()
			defer func() {
				if r := recover(); r != nil {
					// recovered
				}
			}()
			deferExtCount++
			if deferExtCount%2 == 0 {
				panic("even")
			}
		}()
	}
	wg.Wait()
	return deferExtCount
}

// ============================================================================
// Global type alias with concurrent access
// ============================================================================

type CounterInt = int // type alias

var (
	aliasMu    *sync.Mutex
	aliasCount CounterInt
)

func init() {
	aliasMu = &sync.Mutex{}
	aliasCount = 0
}

// AliasIncrement increments a type-aliased global counter.
func AliasIncrement() int {
	aliasMu.Lock()
	aliasCount++
	v := aliasCount
	aliasMu.Unlock()
	return v
}

// AliasGet returns the type-aliased global counter.
func AliasGet() int {
	aliasMu.Lock()
	v := aliasCount
	aliasMu.Unlock()
	return v
}

// ============================================================================
// Concurrent goroutine + channel + select with global state
// ============================================================================

var (
	selectGlobalMu    *sync.Mutex
	selectGlobalCount int
	selectGlobalReady chan struct{}
)

func init() {
	selectGlobalMu = &sync.Mutex{}
	selectGlobalReady = make(chan struct{})
}

// SelectGlobalWaitNotify uses select to wait for a notification and update a global.
func SelectGlobalWaitNotify() int {
	selectGlobalCount = 0
	const N = 10
	done := make(chan int, N)

	// Start N waiters
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			<-selectGlobalReady
			selectGlobalMu.Lock()
			selectGlobalCount++
			selectGlobalMu.Unlock()
			done <- 1
		}()
	}

	// Signal all waiters
	close(selectGlobalReady)

	// Wait for all to complete
	wg.Wait()
	return selectGlobalCount
}

// ============================================================================
// Global []string with mutex — concurrent append
// ============================================================================

var (
	strSliceMu     *sync.Mutex
	globalStrSlice []string
)

func init() {
	strSliceMu = &sync.Mutex{}
	globalStrSlice = make([]string, 0)
}

// StrSliceAppendProtected appends a string to the global string slice with mutex.
func StrSliceAppendProtected(s string) int {
	strSliceMu.Lock()
	globalStrSlice = append(globalStrSlice, s)
	n := len(globalStrSlice)
	strSliceMu.Unlock()
	return n
}

// StrSliceLenProtected returns the length of the global string slice.
func StrSliceLenProtected() int {
	strSliceMu.Lock()
	n := len(globalStrSlice)
	strSliceMu.Unlock()
	return n
}

// StrSliceReset resets the global string slice.
func StrSliceReset() {
	strSliceMu.Lock()
	globalStrSlice = make([]string, 0)
	strSliceMu.Unlock()
}

// ============================================================================
// Semaphore pattern — limit concurrency with buffered channel + global
// ============================================================================

var (
	semMu       *sync.Mutex
	semTotal    int
)

func init() {
	semMu = &sync.Mutex{}
	semTotal = 0
}

// SemaphorePattern limits concurrency using a buffered channel as semaphore.
func SemaphorePattern() int {
	semTotal = 0
	const maxConcurrent = 3
	const totalTasks = 10
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	wg.Add(totalTasks)

	for i := 0; i < totalTasks; i++ {
		sem <- struct{}{} // acquire
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // release
			semMu.Lock()
			semTotal++
			semMu.Unlock()
		}()
	}
	wg.Wait()
	return semTotal
}

// ============================================================================
// Pipeline pattern — stages connected by channels, global accumulator
// ============================================================================

var (
	pipelineMu    *sync.Mutex
	pipelineSum   int
)

func init() {
	pipelineMu = &sync.Mutex{}
	pipelineSum = 0
}

// PipelineGlobal creates a 3-stage pipeline with channels, accumulating results
// into a global counter.
func PipelineGlobal() int {
	pipelineSum = 0
	const N = 10
	stage1 := make(chan int, N)
	stage2 := make(chan int, N)
	stage3 := make(chan int, N)

	// Stage 1: increment
	go func() {
		for v := range stage1 {
			stage2 <- v + 1
		}
		close(stage2)
	}()

	// Stage 2: double
	go func() {
		for v := range stage2 {
			stage3 <- v * 2
		}
		close(stage3)
	}()

	// Stage 3: accumulate into global
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range stage3 {
			pipelineMu.Lock()
			pipelineSum += v
			pipelineMu.Unlock()
		}
	}()

	// Feed input
	for i := 0; i < N; i++ {
		stage1 <- i
	}
	close(stage1)

	wg.Wait()
	// (0+1)*2 + (1+1)*2 + ... + (9+1)*2 = 2+4+6+...+20 = 110
	return pipelineSum
}

// ============================================================================
// Global rune slice with mutex — concurrent append
// ============================================================================

var (
	runeSliceMu     *sync.Mutex
	globalRuneSlice []rune
)

func init() {
	runeSliceMu = &sync.Mutex{}
	globalRuneSlice = make([]rune, 0)
}

// RuneSliceAppendProtected appends a rune to the global slice with mutex.
func RuneSliceAppendProtected(r rune) int {
	runeSliceMu.Lock()
	globalRuneSlice = append(globalRuneSlice, r)
	n := len(globalRuneSlice)
	runeSliceMu.Unlock()
	return n
}

// RuneSliceLenProtected returns the length of the global rune slice.
func RuneSliceLenProtected() int {
	runeSliceMu.Lock()
	n := len(globalRuneSlice)
	runeSliceMu.Unlock()
	return n
}

// RuneSliceReset resets the global rune slice.
func RuneSliceReset() {
	runeSliceMu.Lock()
	globalRuneSlice = make([]rune, 0)
	runeSliceMu.Unlock()
}

// ============================================================================
// Global bool with atomic-style swap via mutex
// ============================================================================

var (
	swapMu    *sync.Mutex
	swapVal   bool
	swapCount int
)

func init() {
	swapMu = &sync.Mutex{}
	swapVal = false
	swapCount = 0
}

// SwapBool atomically swaps the global bool and returns the old value.
func SwapBool(newVal bool) bool {
	swapMu.Lock()
	old := swapVal
	swapVal = newVal
	if old != newVal {
		swapCount++
	}
	swapMu.Unlock()
	return old
}

// SwapGetCount returns how many times the value actually changed.
func SwapGetCount() int {
	swapMu.Lock()
	v := swapCount
	swapMu.Unlock()
	return v
}

// SwapReset resets the swap state.
func SwapReset() {
	swapMu.Lock()
	swapVal = false
	swapCount = 0
	swapMu.Unlock()
}

// ============================================================================
// Goroutine with named return + defer modify during panic
// ============================================================================

var (
	namedReturnMu    *sync.Mutex
	namedReturnCount int
)

func init() {
	namedReturnMu = &sync.Mutex{}
	namedReturnCount = 0
}

// NamedReturnDeferPanicGlobal tests named returns modified by defers
// during panic recovery in goroutines, updating a global.
func NamedReturnDeferPanicGlobal() int {
	namedReturnCount = 0
	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			result := func() (v int) {
				defer func() {
					if r := recover(); r != nil {
						v = 42
					}
					namedReturnMu.Lock()
					namedReturnCount += v
					namedReturnMu.Unlock()
				}()
				panic("test")
				return 0
			}()
			_ = result
		}()
	}
	wg.Wait()
	return namedReturnCount // 20 * 42 = 840
}

// ============================================================================
// Global map with string values — concurrent read+write
// ============================================================================

var (
	strMapMu     *sync.Mutex
	globalStrMap map[string]string
)

func init() {
	strMapMu = &sync.Mutex{}
	globalStrMap = make(map[string]string)
}

// StrMapPut puts a key-value pair into the global string map.
func StrMapPut(key, val string) {
	strMapMu.Lock()
	globalStrMap[key] = val
	strMapMu.Unlock()
}

// StrMapGet gets a value from the global string map.
func StrMapGet(key string) (string, bool) {
	strMapMu.Lock()
	v, ok := globalStrMap[key]
	strMapMu.Unlock()
	return v, ok
}

// StrMapLen returns the size of the global string map.
func StrMapLen() int {
	strMapMu.Lock()
	n := len(globalStrMap)
	strMapMu.Unlock()
	return n
}

// StrMapReset resets the global string map.
func StrMapReset() {
	strMapMu.Lock()
	globalStrMap = make(map[string]string)
	strMapMu.Unlock()
}

// ============================================================================
// Global float64 slice with mutex — concurrent append
// ============================================================================

var (
	floatSliceMu     *sync.Mutex
	globalFloatSlice []float64
)

func init() {
	floatSliceMu = &sync.Mutex{}
	globalFloatSlice = make([]float64, 0)
}

// FloatSliceAppendProtected appends a float64 to the global slice with mutex.
func FloatSliceAppendProtected(v float64) int {
	floatSliceMu.Lock()
	globalFloatSlice = append(globalFloatSlice, v)
	n := len(globalFloatSlice)
	floatSliceMu.Unlock()
	return n
}

// FloatSliceLenProtected returns the length of the global float64 slice.
func FloatSliceLenProtected() int {
	floatSliceMu.Lock()
	n := len(globalFloatSlice)
	floatSliceMu.Unlock()
	return n
}

// FloatSliceReset resets the global float64 slice.
func FloatSliceReset() {
	floatSliceMu.Lock()
	globalFloatSlice = make([]float64, 0)
	floatSliceMu.Unlock()
}
