package divergence_hunt220

import (
	"fmt"
	"sync"
)

func RWMutexBasic() string {
	var mu sync.RWMutex
	data := "initial"

	mu.Lock()
	data = "modified"
	mu.Unlock()

	mu.RLock()
	result := data
	mu.RUnlock()

	return fmt.Sprintf("data: %s", result)
}

func RWMutexMultipleReaders() string {
	var mu sync.RWMutex
	count := 0

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.RLock()
			count++
			mu.RUnlock()
		}()
	}

	wg.Wait()
	return fmt.Sprintf("readers: %d", count)
}

func RWMutexWriterExclusion() string {
	// Two concurrent writers produce different final values non-deterministically;
	// assert only that the final value is one of the expected writes.
	var mu sync.RWMutex
	data := 0

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mu.Lock()
		data = 42
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu.Lock()
		data = 100
		mu.Unlock()
	}()

	wg.Wait()
	ok := data == 42 || data == 100
	return fmt.Sprintf("data-ok=%v", ok)
}

func RWMutexDeferredUnlock() string {
	var mu sync.RWMutex
	result := ""

	func() {
		mu.Lock()
		defer mu.Unlock()
		result = "locked"
	}()

	return result
}

func RWMutexDeferredRUnlock() string {
	var mu sync.RWMutex
	data := "shared"

	func() string {
		mu.RLock()
		defer mu.RUnlock()
		return data
	}()

	return data
}

func RWMutexPromote() string {
	var mu sync.RWMutex
	data := 0

	mu.RLock()
	_ = data
	mu.RUnlock()

	mu.Lock()
	data = 1
	mu.Unlock()

	return fmt.Sprintf("promoted: %d", data)
}

func RWMutexNestedReadLock() string {
	var mu sync.RWMutex
	count := 0

	mu.RLock()
	count++
	mu.RLock()
	count++
	mu.RUnlock()
	mu.RUnlock()

	return fmt.Sprintf("nested reads: %d", count)
}

func RWMutexWithWaitGroup() string {
	var mu sync.RWMutex
	data := 0
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if n%2 == 0 {
				mu.Lock()
				data += n
				mu.Unlock()
			} else {
				mu.RLock()
				_ = data
				mu.RUnlock()
			}
		}(i)
	}

	wg.Wait()
	return fmt.Sprintf("data: %d", data)
}

func RWMutexTryLock() string {
	var mu sync.RWMutex

	locked := mu.TryLock()
	if locked {
		mu.Unlock()
	}

	return fmt.Sprintf("trylock: %v", locked)
}

func RWMutexTryRLock() string {
	var mu sync.RWMutex

	locked := mu.TryRLock()
	if locked {
		mu.RUnlock()
	}

	return fmt.Sprintf("tryrlock: %v", locked)
}

func RWMutexInStruct() string {
	type Cache struct {
		mu   sync.RWMutex
		data map[string]string
	}

	c := Cache{data: make(map[string]string)}

	c.mu.Lock()
	c.data["key"] = "value"
	c.mu.Unlock()

	c.mu.RLock()
	val := c.data["key"]
	c.mu.RUnlock()

	return fmt.Sprintf("value: %s", val)
}

func RWMutexReadPrefer() string {
	var mu sync.RWMutex
	var cmu sync.Mutex
	reads := 0
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.RLock()
			cmu.Lock()
			reads++
			cmu.Unlock()
			mu.RUnlock()
		}()
	}

	wg.Wait()
	return fmt.Sprintf("reads: %d", reads)
}
