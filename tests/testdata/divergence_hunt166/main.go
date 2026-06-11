package divergence_hunt166

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Round 166: Goroutine basic patterns (launch, wait, sync)
// ============================================================================

// BasicGoroutineLaunch tests basic goroutine launch
func BasicGoroutineLaunch() string {
	result := make(chan string, 1)
	go func() {
		result <- "goroutine executed"
	}()
	return <-result
}

// MultipleGoroutines tests launching multiple goroutines
func MultipleGoroutines() string {
	results := make(chan int, 3)
	for i := 1; i <= 3; i++ {
		go func(n int) {
			results <- n * 10
		}(i)
	}
	sum := 0
	for i := 0; i < 3; i++ {
		sum += <-results
	}
	return fmt.Sprintf("sum=%d", sum)
}

// WaitGroupBasic tests WaitGroup basic usage
func WaitGroupBasic() string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	counter := 0
	for i := 0; i < 5; i++ {
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

// WaitGroupWithData tests WaitGroup with data collection
func WaitGroupWithData() string {
	var wg sync.WaitGroup
	results := make([]int, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = idx * idx
		}(i)
	}
	wg.Wait()
	sum := 0
	for _, v := range results {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ChannelSynchronization tests channel-based synchronization
func ChannelSynchronization() string {
	done := make(chan bool)
	result := ""
	go func() {
		result = "completed"
		done <- true
	}()
	<-done
	return fmt.Sprintf("result=%s", result)
}

// BufferedChannelGoroutines tests buffered channels with goroutines
func BufferedChannelGoroutines() string {
	ch := make(chan int, 3)
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
	}()
	time.Sleep(10 * time.Millisecond)
	sum := 0
	for i := 0; i < 3; i++ {
		sum += <-ch
	}
	return fmt.Sprintf("sum=%d", sum)
}

// GoroutineClosureCapture tests closure capture in goroutines
func GoroutineClosureCapture() string {
	results := make(chan int, 5)
	for i := 0; i < 5; i++ {
		go func(n int) {
			results <- n * 2
		}(i)
	}
	sum := 0
	for i := 0; i < 5; i++ {
		sum += <-results
	}
	return fmt.Sprintf("sum=%d", sum)
}

// SelectStatement tests select statement with goroutines.
// Select order between multiple-ready channels is non-deterministic;
// verify only that both messages are received.
func SelectStatement() string {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	go func() { ch1 <- "first" }()
	go func() { ch2 <- "second" }()
	got := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			got[msg] = true
		case msg := <-ch2:
			got[msg] = true
		}
	}
	return fmt.Sprintf("first=%v,second=%v", got["first"], got["second"])
}

// FanOutPattern tests fan-out pattern
func FanOutPattern() string {
	work := make(chan int, 10)
	results := make(chan int, 10)
	// Worker function
	worker := func() {
		for n := range work {
			results <- n * n
		}
	}
	// Start 3 workers
	for i := 0; i < 3; i++ {
		go worker()
	}
	// Send work
	for i := 1; i <= 5; i++ {
		work <- i
	}
	close(work)
	time.Sleep(10 * time.Millisecond)
	// Collect results
	sum := 0
	for i := 0; i < 5; i++ {
		sum += <-results
	}
	return fmt.Sprintf("sum=%d", sum)
}

// FanInPattern tests fan-in pattern
func FanInPattern() string {
	ch1 := make(chan int)
	ch2 := make(chan int)
	merged := make(chan int, 6)
	// Send to ch1
	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- i
		}
		close(ch1)
	}()
	// Send to ch2
	go func() {
		for i := 4; i <= 6; i++ {
			ch2 <- i
		}
		close(ch2)
	}()
	// Fan in
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for n := range ch1 {
			merged <- n
		}
	}()
	go func() {
		defer wg.Done()
		for n := range ch2 {
			merged <- n
		}
	}()
	go func() {
		wg.Wait()
		close(merged)
	}()
	sum := 0
	for n := range merged {
		sum += n
	}
	return fmt.Sprintf("sum=%d", sum)
}

// GoroutineWithError tests error handling in goroutines
func GoroutineWithError() string {
	type Result struct {
		value int
		err   string
	}
	result := make(chan Result, 1)
	go func() {
		if true {
			result <- Result{value: 0, err: "error occurred"}
		} else {
			result <- Result{value: 42, err: ""}
		}
	}()
	r := <-result
	return fmt.Sprintf("err=%s", r.err)
}

// TimeoutPattern tests timeout pattern with select
func TimeoutPattern() string {
	ch := make(chan string)
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "result"
	}()
	select {
	case msg := <-ch:
		return fmt.Sprintf("msg=%s", msg)
	case <-time.After(5 * time.Millisecond):
		return "timeout"
	}
}
