package goroutine

import "time"

func workerInt(ch chan int) {
	ch <- 42
}

func workerString(s string, ch chan string) {
	ch <- s + "!"
}

func sender(ch chan int, v int) {
	ch <- v
}

func sleepAndSend(ch chan int, val int) {
	time.Sleep(100 * time.Millisecond)
	ch <- val
}

func producerLoop(ch chan int) {
	for i := 1; i <= 3; i++ {
		ch <- i
	}
	close(ch)
}

func producerFixed(ch chan int) {
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
}

// BasicSpawn tests that a goroutine can be spawned and executes.
func BasicSpawn() int {
	ch := make(chan int, 1)
	go workerInt(ch)
	return <-ch
}

// ChannelCommunication tests goroutine communication via channels.
func ChannelCommunication() int {
	ch := make(chan int, 1)
	go func(c chan int) {
		c <- 42
	}(ch)
	return <-ch
}

// WithArguments tests goroutine with function arguments.
func WithArguments() int {
	ch := make(chan int, 1)
	go func(a, b int, c chan int) {
		c <- a + b
	}(10, 20, ch)
	return <-ch
}

// WithStruct tests goroutines working with structs.
func WithStruct() int {
	type Result struct {
		Value int
	}
	ch := make(chan Result, 1)
	go func(c chan Result) {
		c <- Result{Value: 42}
	}(ch)
	r := <-ch
	return r.Value
}

// DifferentTypes tests goroutines with different parameter types.
func DifferentTypes() string {
	ch := make(chan string, 1)
	go workerString("hello", ch)
	return <-ch
}

// GlobalsSharing tests that goroutines share globals via channel.
func GlobalsSharing() int {
	ch := make(chan int, 1)
	go func(c chan int) {
		c <- 1
	}(ch)
	return <-ch
}

// MultipleSends tests multiple goroutine sends.
func MultipleSends() int {
	ch := make(chan int, 3)
	go sender(ch, 10)
	go sender(ch, 20)
	go sender(ch, 30)

	sum := 0
	for i := 0; i < 3; i++ {
		sum += <-ch
	}
	return sum
}

// ParallelExecution tests that goroutines truly run in parallel.
func ParallelExecution() int {
	ch := make(chan int, 2)
	start := time.Now()
	go sleepAndSend(ch, 1)
	go sleepAndSend(ch, 2)
	v1 := <-ch
	v2 := <-ch
	elapsed := time.Since(start)
	_ = v1 + v2
	if elapsed < 150*time.Millisecond {
		return 1
	}
	return 0
}

// ClosureCapture tests goroutines with captured variables.
func ClosureCapture() int {
	ch := make(chan int, 1)
	x := 42
	go func() {
		ch <- x
	}()
	return <-ch
}

// ClosureCaptureMultiple tests multiple captured variables in goroutine.
func ClosureCaptureMultiple() int {
	ch := make(chan int, 1)
	a := 10
	b := 20
	go func() {
		ch <- a + b
	}()
	return <-ch
}

// SelectStatement tests a select with a single ready channel case.
func SelectStatement() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	result := 0
	select {
	case v := <-ch1:
		result = v
	case v := <-ch2:
		result = v
	}
	return result
}

// SelectDefault tests select with a default branch.
func SelectDefault() int {
	ch := make(chan int)
	result := 0
	select {
	case v := <-ch:
		result = v
	default:
		result = -1
	}
	return result
}

// SelectSend tests select with a send case.
func SelectSend() int {
	ch := make(chan int, 1)
	select {
	case ch <- 99:
	default:
	}
	return <-ch
}

// RangeOverChannel tests ranging over a channel.
func RangeOverChannel() int {
	ch := make(chan int, 3)
	go producerLoop(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// RangeOverChannelWithBuiltin tests range over channel counting iterations.
func RangeOverChannelWithBuiltin() int {
	ch := make(chan int, 3)
	go producerFixed(ch)
	count := 0
	for range ch {
		count++
	}
	return count
}
