package divergence_hunt217

import (
	"fmt"
	"sync"
)

func OnceBasic() string {
	var once sync.Once
	count := 0

	once.Do(func() {
		count++
	})
	once.Do(func() {
		count++
	})

	return fmt.Sprintf("count: %d", count)
}

func OnceWithGoroutines() string {
	var once sync.Once
	var wg sync.WaitGroup
	count := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				count++
			})
		}()
	}

	wg.Wait()
	return fmt.Sprintf("count: %d", count)
}

func OnceWithDataInitialization() string {
	var once sync.Once
	var data string

	getData := func() string {
		once.Do(func() {
			data = "initialized"
		})
		return data
	}

	_ = getData()
	_ = getData()
	return getData()
}

func OnceWithPanic() string {
	var once sync.Once
	count := 0

	defer func() {
		recover()
	}()

	once.Do(func() {
		count++
		panic("oops")
	})

	once.Do(func() {
		count++
	})

	return fmt.Sprintf("count: %d", count)
}

func OncePerInstance() string {
	var once1, once2 sync.Once
	count := 0

	once1.Do(func() { count++ })
	once2.Do(func() { count++ })

	return fmt.Sprintf("count: %d", count)
}

func OnceWithComplexInitialization() string {
	var once sync.Once
	config := make(map[string]string)

	once.Do(func() {
		config["host"] = "localhost"
		config["port"] = "8080"
	})

	return fmt.Sprintf("host=%s, port=%s", config["host"], config["port"])
}

func MultipleOnceVariables() string {
	var initA, initB sync.Once
	result := ""

	initA.Do(func() { result += "A" })
	initB.Do(func() { result += "B" })
	initA.Do(func() { result += "A2" })
	initB.Do(func() { result += "B2" })

	return result
}

func OnceWithMutexCombo() string {
	var once sync.Once
	var mu sync.Mutex
	value := 0

	init := func() {
		mu.Lock()
		value = 42
		mu.Unlock()
	}

	once.Do(init)
	mu.Lock()
	v := value
	mu.Unlock()

	return fmt.Sprintf("value: %d", v)
}

func OnceInStructSlice() string {
	type Item struct {
		Once  sync.Once
		Value int
	}

	items := make([]Item, 3)
	for i := range items {
		idx := i
		items[idx].Once.Do(func() {
			items[idx].Value = idx * 10
		})
	}

	return fmt.Sprintf("values: %d,%d,%d", items[0].Value, items[1].Value, items[2].Value)
}

func OnceWithLazyLoading() string {
	type Service struct {
		once   sync.Once
		client string
	}

	s := &Service{}

	getClient := func() string {
		s.once.Do(func() {
			s.client = "connected"
		})
		return s.client
	}

	_ = getClient()
	return getClient()
}

func OnceWithChannelClose() string {
	var once sync.Once
	ch := make(chan struct{})

	ready := func() {
		once.Do(func() {
			close(ch)
		})
	}

	ready()
	ready()

	select {
	case <-ch:
		return "closed"
	default:
		return "open"
	}
}
