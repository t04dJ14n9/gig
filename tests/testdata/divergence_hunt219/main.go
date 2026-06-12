package divergence_hunt219

import (
	"fmt"
	"sync"
)

// ============================================================================
// Round 219: Mutex-based atomic-like operations
// ============================================================================

// MutexAddInt32 tests mutex-protected int32 addition
func MutexAddInt32() string {
	var mu sync.Mutex
	val := int32(10)
	mu.Lock()
	val += 5
	mu.Unlock()
	return fmt.Sprintf("result: %d, val: %d", val, val)
}

// MutexAddInt64 tests mutex-protected int64 addition
func MutexAddInt64() string {
	var mu sync.Mutex
	val := int64(100)
	mu.Lock()
	val += 50
	mu.Unlock()
	return fmt.Sprintf("result: %d, val: %d", val, val)
}

// MutexLoadStoreInt32 tests mutex-protected int32 load/store
func MutexLoadStoreInt32() string {
	var mu sync.Mutex
	val := int32(42)
	mu.Lock()
	loaded := val
	val = 99
	newVal := val
	mu.Unlock()
	return fmt.Sprintf("loaded: %d, new: %d", loaded, newVal)
}

// MutexLoadStoreInt64 tests mutex-protected int64 load/store
func MutexLoadStoreInt64() string {
	var mu sync.Mutex
	val := int64(1000)
	mu.Lock()
	loaded := val
	val = 2000
	newVal := val
	mu.Unlock()
	return fmt.Sprintf("loaded: %d, new: %d", loaded, newVal)
}

// MutexSwapInt32 tests mutex-protected int32 swap
func MutexSwapInt32() string {
	var mu sync.Mutex
	val := int32(10)
	mu.Lock()
	old := val
	val = 20
	mu.Unlock()
	return fmt.Sprintf("old: %d, new: %d", old, val)
}

// MutexCompareAndSwapInt32 tests mutex-protected int32 CAS
func MutexCompareAndSwapInt32() string {
	var mu sync.Mutex
	val := int32(10)
	mu.Lock()
	swapped := val == 10
	if swapped {
		val = 20
	}
	mu.Unlock()
	mu.Lock()
	swapped2 := val == 10
	if swapped2 {
		val = 30
	}
	mu.Unlock()
	return fmt.Sprintf("first=%v,second=%v,val=%d", swapped, swapped2, val)
}

// MutexCounterPattern tests mutex-protected counter
func MutexCounterPattern() string {
	var mu sync.Mutex
	counter := int64(0)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			mu.Lock()
			counter += int64(n)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return fmt.Sprintf("sum: %d", counter)
}

// MutexFlagPattern tests mutex-protected flag
func MutexFlagPattern() string {
	var mu sync.Mutex
	ready := false
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		ready = true
		mu.Unlock()
	}()
	wg.Wait()
	mu.Lock()
	isReady := ready
	mu.Unlock()
	return fmt.Sprintf("ready: %v", isReady)
}

// MutexMaxValue tests finding max with mutex
func MutexMaxValue() string {
	var mu sync.Mutex
	max := int32(0)
	values := []int32{5, 10, 3, 20, 8}
	for _, v := range values {
		mu.Lock()
		if v > max {
			max = v
		}
		mu.Unlock()
	}
	return fmt.Sprintf("max: %d", max)
}

// ProtectedInt32 tests protected int32 operations
func ProtectedInt32() string {
	type ProtectedInt struct {
		sync.Mutex
		value int32
	}
	p := &ProtectedInt{value: 100}
	p.Lock()
	p.value += 50
	result := p.value
	p.Unlock()
	return fmt.Sprintf("value: %d", result)
}

// ProtectedInt64 tests protected int64 operations
func ProtectedInt64() string {
	type ProtectedInt struct {
		sync.Mutex
		value int64
	}
	p := &ProtectedInt{value: 1000}
	p.Lock()
	p.value += 500
	result := p.value
	p.Unlock()
	return fmt.Sprintf("value: %d", result)
}

// MutexUint32 tests mutex-protected uint32
func MutexUint32() string {
	var mu sync.Mutex
	val := uint32(10)
	mu.Lock()
	val += 5
	val--
	mu.Unlock()
	return fmt.Sprintf("value: %d", val)
}

// MutexUint64 tests mutex-protected uint64
func MutexUint64() string {
	var mu sync.Mutex
	val := uint64(100)
	mu.Lock()
	val += 50
	mu.Unlock()
	return fmt.Sprintf("value: %d", val)
}
