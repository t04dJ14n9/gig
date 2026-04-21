package divergence_hunt216

import (
	"fmt"
	"sync"
)

func WaitGroupBasic() string {
	var wg sync.WaitGroup
	result := ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		result += "goroutine done"
	}()

	wg.Wait()
	return result
}

func WaitGroupWithParam() string {
	// Goroutine execution order is non-deterministic. Verify that exactly three
	// workers ran by counting, not by comparing output order.
	var wg sync.WaitGroup
	var mu sync.Mutex
	seen := map[int]bool{}

	worker := func(id int) {
		defer wg.Done()
		mu.Lock()
		seen[id] = true
		mu.Unlock()
	}

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i)
	}

	wg.Wait()
	return fmt.Sprintf("workers=%d all=%v", len(seen), seen[1] && seen[2] && seen[3])
}

func WaitGroupMultipleAdd() string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}

	wg.Wait()
	return fmt.Sprintf("count: %d", count)
}

func WaitGroupAddInsideGoroutine() string {
	var wg sync.WaitGroup
	result := ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		wg.Add(1)
		go func() {
			defer wg.Done()
			result = "nested"
		}()
	}()

	wg.Wait()
	return result
}

func WaitGroupWithMutex() string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}

	wg.Wait()
	return fmt.Sprintf("count: %d", count)
}

func NestedWaitGroup() string {
	// Concurrent goroutines contribute "a" and "b" in non-deterministic order.
	// Verify only set membership and that the outer goroutine's "c" comes last.
	var outer, inner sync.WaitGroup
	var mu sync.Mutex
	pieces := []string{}

	outer.Add(1)
	go func() {
		defer outer.Done()

		inner.Add(2)
		go func() {
			defer inner.Done()
			mu.Lock()
			pieces = append(pieces, "a")
			mu.Unlock()
		}()
		go func() {
			defer inner.Done()
			mu.Lock()
			pieces = append(pieces, "b")
			mu.Unlock()
		}()
		inner.Wait()
		mu.Lock()
		pieces = append(pieces, "c")
		mu.Unlock()
	}()

	outer.Wait()
	// Sort first two to hide non-determinism; last item is always "c".
	if len(pieces) == 3 && pieces[0] > pieces[1] {
		pieces[0], pieces[1] = pieces[1], pieces[0]
	}
	return fmt.Sprintf("%v", pieces)
}

func WaitGroupWithErrorResult() string {
	var wg sync.WaitGroup
	errors := make([]error, 0)
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			errors = append(errors, fmt.Errorf("error %d", id))
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return fmt.Sprintf("errors: %d", len(errors))
}

func WaitGroupReuse() string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := 0

	for batch := 0; batch < 2; batch++ {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				result++
				mu.Unlock()
			}()
		}
		wg.Wait()
	}

	return fmt.Sprintf("result: %d", result)
}

func WaitGroupWithChannelResult() string {
	var wg sync.WaitGroup
	results := make(chan int, 3)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			results <- n * n
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for n := range results {
		sum += n
	}
	return fmt.Sprintf("sum: %d", sum)
}

func WaitGroupDeferDone() string {
	// Goroutine scheduling is non-deterministic; instead of testing output
	// order, just verify that exactly 3 start/end pairs completed.
	var wg sync.WaitGroup
	var mu sync.Mutex
	starts := 0
	ends := 0

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			mu.Lock()
			starts++
			mu.Unlock()
			defer func() {
				mu.Lock()
				ends++
				mu.Unlock()
				wg.Done()
			}()
		}(i)
	}

	wg.Wait()
	return fmt.Sprintf("starts=%d ends=%d", starts, ends)
}

func WaitGroupZeroValue() string {
	var wg sync.WaitGroup
	wg.Add(1)

	result := ""
	go func() {
		defer wg.Done()
		result = "zero value works"
	}()

	wg.Wait()
	return result
}
