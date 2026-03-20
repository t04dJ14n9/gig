package thirdparty

import "sync"

// SyncMutex tests sync.Mutex.
func SyncMutex() int {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()
	return 1
}

// SyncMutexCounter tests mutex-protected counter.
func SyncMutexCounter() int {
	var mu sync.Mutex
	counter := 0
	for i := 0; i < 100; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
	return counter
}

// SyncRWMutex tests sync.RWMutex.
func SyncRWMutex() int {
	var mu sync.RWMutex
	mu.RLock()
	mu.RUnlock()
	mu.Lock()
	mu.Unlock()
	return 1
}

// SyncWaitGroup tests sync.WaitGroup.
func SyncWaitGroup() int {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
	}()
	wg.Wait()
	return 1
}

// SyncOnce tests sync.Once.
func SyncOnce() int {
	var once sync.Once
	counter := 0
	for i := 0; i < 10; i++ {
		once.Do(func() {
			counter++
		})
	}
	return counter
}

// SyncOnceFunc tests sync.OnceFunc.
func SyncOnceFunc() int {
	counter := 0
	fn := sync.OnceFunc(func() {
		counter++
	})
	for i := 0; i < 10; i++ {
		fn()
	}
	return counter
}

// SyncMap tests sync.Map.
func SyncMap() int {
	var m sync.Map
	m.Store("key", 42)
	v, ok := m.Load("key")
	if ok && v == 42 {
		return 1
	}
	return 0
}

// SyncMapLoadOrStore tests LoadOrStore.
func SyncMapLoadOrStore() int {
	var m sync.Map
	v, loaded := m.LoadOrStore("key", 42)
	if !loaded && v == 42 {
		return 1
	}
	return 0
}
