package divergence_hunt167

import (
	"fmt"
	"sync"
)

// ============================================================================
// Round 167: Mutex and sync patterns
// ============================================================================

// BasicMutex tests basic mutex usage
func BasicMutex() string {
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return fmt.Sprintf("counter=%d", counter)
}

// RWMutexBasic tests RWMutex basic usage
func RWMutexBasic() string {
	var rw sync.RWMutex
	value := 42
	var wg sync.WaitGroup
	// Readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rw.RLock()
			_ = value
			rw.RUnlock()
		}()
	}
	// Writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		rw.Lock()
		value = 100
		rw.Unlock()
	}()
	wg.Wait()
	return fmt.Sprintf("value=%d", value)
}

// DeferredUnlock tests deferred unlock
func DeferredUnlock() string {
	var mu sync.Mutex
	counter := 0
	increment := func() {
		mu.Lock()
		defer mu.Unlock()
		counter++
	}
	increment()
	increment()
	increment()
	return fmt.Sprintf("counter=%d", counter)
}

// MutexOperations tests mutex-protected operations (replaces atomic)
func MutexOperations() string {
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return fmt.Sprintf("counter=%d", counter)
}

// MutexLoadStore tests mutex-protected load/store (replaces atomic)
func MutexLoadStore() string {
	var mu sync.Mutex
	value := 10
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		mu.Lock()
		value = 20
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu.Lock()
		_ = value
		mu.Unlock()
	}()
	wg.Wait()
	mu.Lock()
	result := value
	mu.Unlock()
	return fmt.Sprintf("value=%d", result)
}

// OncePattern tests sync.Once pattern
func OncePattern() string {
	var once sync.Once
	counter := 0
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				counter++
			})
		}()
	}
	wg.Wait()
	return fmt.Sprintf("counter=%d", counter)
}

// PoolPattern tests sync.Pool pattern
func PoolPattern() string {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	// Get from pool
	item := pool.Get().([]byte)
	item[0] = 42
	// Put back to pool
	pool.Put(item)
	// Get again (might be same object)
	item2 := pool.Get().([]byte)
	return fmt.Sprintf("item2[0]=%d", item2[0])
}

// CondPattern tests sync.Cond pattern
func CondPattern() string {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		for !ready {
			cond.Wait()
		}
		mu.Unlock()
	}()
	// Signal the goroutine
	mu.Lock()
	ready = true
	cond.Signal()
	mu.Unlock()
	wg.Wait()
	return "signaled"
}

// MapPattern tests sync.Map pattern
func MapPattern() string {
	var m sync.Map
	// Store values
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	// Load values
	a, _ := m.Load("a")
	b, _ := m.Load("b")
	// Load or store
	d, loaded := m.LoadOrStore("d", 4)
	return fmt.Sprintf("a=%v,b=%v,d=%v,loaded=%t", a, b, d, loaded)
}

// ProtectedStruct tests struct with embedded mutex
func ProtectedStruct() string {
	type Counter struct {
		sync.Mutex
		value int
	}
	c := Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Lock()
			c.value++
			c.Unlock()
		}()
	}
	wg.Wait()
	return fmt.Sprintf("value=%d", c.value)
}

// TryLockPattern tests TryLock (Go 1.18+)
func TryLockPattern() string {
	var mu sync.Mutex
	locked := mu.TryLock()
	result := ""
	if locked {
		result = "locked"
		mu.Unlock()
	} else {
		result = "not locked"
	}
	return fmt.Sprintf("result=%s", result)
}
