package divergence_hunt131

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Round 131: Goroutine edge cases with sync primitives
// ============================================================================

func GoroutineWaitGroup() string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sum := 0
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()
			sum += val
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return fmt.Sprintf("sum=%d", sum)
}

func GoroutineOnce() string {
	var once sync.Once
	count := 0
	for i := 0; i < 5; i++ {
		once.Do(func() {
			count++
		})
	}
	return fmt.Sprintf("count=%d", count)
}

func GoroutineChannelSum() string {
	ch := make(chan int, 3)
	go func() {
		ch <- 10
		ch <- 20
		ch <- 30
		close(ch)
	}()
	sum := 0
	for v := range ch {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

func GoroutineMutex() string {
	var mu sync.Mutex
	x := 0
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			x++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return fmt.Sprintf("x=%d", x)
}

func GoroutineSelectTimeout() string {
	ch := make(chan int)
	select {
	case v := <-ch:
		return fmt.Sprintf("got-%d", v)
	case <-time.After(10 * time.Millisecond):
		return "timeout"
	}
}

func GoroutineSendReceive() string {
	ch := make(chan string, 1)
	go func() {
		ch <- "hello"
	}()
	v := <-ch
	return v
}

func GoroutineCloseSignal() string {
	done := make(chan struct{})
	go func() {
		close(done)
	}()
	<-done
	return "done"
}

func GoroutinePanicRecover() string {
	ch := make(chan string, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- fmt.Sprintf("recovered-%v", r)
			}
		}()
		panic("goroutine-panic")
	}()
	v := <-ch
	return v
}
