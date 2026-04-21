package divergence_hunt218

import (
	"fmt"
	"sync"
)

func SyncMapBasic() string {
	var m sync.Map
	m.Store("key1", "value1")
	m.Store("key2", 42)

	v1, _ := m.Load("key1")
	v2, _ := m.Load("key2")

	return fmt.Sprintf("%v, %v", v1, v2)
}

func SyncMapLoadOrStore() string {
	var m sync.Map

	v1, loaded1 := m.LoadOrStore("key", "first")
	v2, loaded2 := m.LoadOrStore("key", "second")

	return fmt.Sprintf("v1=%v,loaded=%v; v2=%v,loaded=%v", v1, loaded1, v2, loaded2)
}

func SyncMapDelete() string {
	var m sync.Map
	m.Store("key", "value")

	m.Delete("key")
	_, ok := m.Load("key")

	return fmt.Sprintf("found after delete: %v", ok)
}

func SyncMapRange() string {
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	result := ""
	m.Range(func(key, value interface{}) bool {
		result += fmt.Sprintf("%v=%v;", key, value)
		return true
	})

	return result
}

func SyncMapRangeEarlyExit() string {
	var m sync.Map
	for i := 0; i < 10; i++ {
		m.Store(i, i*i)
	}

	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return count < 3
	})

	return fmt.Sprintf("iterated: %d", count)
}

func SyncMapLoadAndDelete() string {
	var m sync.Map
	m.Store("key", "value")

	v, loaded := m.LoadAndDelete("key")
	_, found := m.Load("key")

	return fmt.Sprintf("v=%v,loaded=%v,found=%v", v, loaded, found)
}

func SyncMapSwap() string {
	var m sync.Map
	m.Store("key", "old")

	prev, loaded := m.Swap("key", "new")

	return fmt.Sprintf("prev=%v,loaded=%v", prev, loaded)
}

func SyncMapCompareAndSwap() string {
	var m sync.Map
	m.Store("key", "old")

	swapped := m.CompareAndSwap("key", "old", "new")
	swapped2 := m.CompareAndSwap("key", "old", "another")

	return fmt.Sprintf("first=%v,second=%v", swapped, swapped2)
}

func SyncMapCompareAndDelete() string {
	var m sync.Map
	m.Store("key", "value")

	deleted := m.CompareAndDelete("key", "wrong")
	deleted2 := m.CompareAndDelete("key", "value")

	return fmt.Sprintf("first=%v,second=%v", deleted, deleted2)
}

func SyncMapWithGoroutines() string {
	var m sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Store(n, n*n)
		}(i)
	}

	wg.Wait()

	count := 0
	m.Range(func(_, _ interface{}) bool {
		count++
		return true
	})

	return fmt.Sprintf("stored: %d", count)
}

func SyncMapTypeSafety() string {
	var m sync.Map
	m.Store("int", 42)
	m.Store("string", "hello")

	v1, _ := m.Load("int")
	v2, _ := m.Load("string")

	i, ok1 := v1.(int)
	s, ok2 := v2.(string)

	return fmt.Sprintf("int=%d(%v),string=%s(%v)", i, ok1, s, ok2)
}
