package divergence_hunt57

import "sync"

// ============================================================================
// Round 57: Sync primitives - Mutex, RWMutex, Once, WaitGroup patterns
// ============================================================================

func MutexBasic() int {
	var mu sync.Mutex
	x := 0
	mu.Lock()
	x++
	mu.Unlock()
	return x
}

func MutexInDefer() int {
	var mu sync.Mutex
	x := 0
	func() {
		mu.Lock()
		defer mu.Unlock()
		x++
	}()
	return x
}

func RWMutexBasic() int {
	var mu sync.RWMutex
	x := 0
	mu.Lock()
	x = 42
	mu.Unlock()
	mu.RLock()
	v := x
	mu.RUnlock()
	return v
}

func OnceBasic() int {
	var once sync.Once
	count := 0
	once.Do(func() { count++ })
	once.Do(func() { count++ })
	once.Do(func() { count++ })
	return count
}

func MutexCounter() int {
	var mu sync.Mutex
	count := 0
	for i := 0; i < 10; i++ {
		mu.Lock()
		count++
		mu.Unlock()
	}
	return count
}

func OnceInClosure() int {
	var once sync.Once
	result := 0
	init := func() { result = 42 }
	once.Do(init)
	once.Do(func() { result = 99 })
	return result
}

func RWMutexMultipleReaders() int {
	var mu sync.RWMutex
	data := 100
	mu.RLock()
	v1 := data
	mu.RUnlock()
	mu.RLock()
	v2 := data
	mu.RUnlock()
	return v1 + v2
}

func MutexSwapPattern() int {
	var mu sync.Mutex
	a, b := 1, 2
	mu.Lock()
	a, b = b, a
	mu.Unlock()
	return a*10 + b
}

func MutexMapProtect() int {
	var mu sync.Mutex
	m := map[string]int{}
	mu.Lock()
	m["key"] = 42
	mu.Unlock()
	mu.Lock()
	v := m["key"]
	mu.Unlock()
	return v
}

func OnceLazyInit() int {
	type Config struct {
		Value int
	}
	var config *Config
	var once sync.Once
	getConfig := func() *Config {
		once.Do(func() {
			config = &Config{Value: 99}
		})
		return config
	}
	return getConfig().Value
}

func FmtMutex() string {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()
	return "ok"
}

func MutexNestedLock() int {
	var mu sync.Mutex
	x := 0
	mu.Lock()
	x++
	mu.Unlock()
	// Lock again (not nested - sequential)
	mu.Lock()
	x++
	mu.Unlock()
	return x
}
