package divergence_hunt212

import (
	"fmt"
)

func RangeOverClosedChannel() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	result := ""
	for val := range ch {
		result += fmt.Sprintf("%d", val)
	}
	return result
}

func RangeOverEmptyClosedChannel() string {
	ch := make(chan int)
	close(ch)

	result := ""
	for val := range ch {
		result += fmt.Sprintf("%d", val)
	}
	return fmt.Sprintf("empty: %s", result)
}

func CloseChannelReceiveValues() string {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)

	v1, ok1 := <-ch
	v2, ok2 := <-ch
	v3, ok3 := <-ch

	return fmt.Sprintf("v1=%d,ok1=%v; v2=%d,ok2=%v; v3=%d,ok3=%v", v1, ok1, v2, ok2, v3, ok3)
}

func CloseChannelWithGoroutine() string {
	ch := make(chan string)
	result := make(chan string)

	go func() {
		for val := range ch {
			result <- fmt.Sprintf("got: %s", val)
		}
		result <- "done"
	}()

	ch <- "a"
	ch <- "b"
	close(ch)

	r1 := <-result
	r2 := <-result
	r3 := <-result

	return fmt.Sprintf("%s; %s; %s", r1, r2, r3)
}

func DoubleClosePanic() string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic")
		}
	}()

	ch := make(chan int)
	close(ch)
	close(ch)

	return "no panic"
}

func CloseNilChannelPanic() string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from nil close panic")
		}
	}()

	var ch chan int
	close(ch)

	return "no panic"
}

func RangeWithBreak() string {
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)

	result := ""
	for val := range ch {
		if val == 3 {
			break
		}
		result += fmt.Sprintf("%d", val)
	}
	return result
}

func RangeStringChannel() string {
	ch := make(chan string, 3)
	ch <- "hello"
	ch <- "world"
	ch <- "!"
	close(ch)

	result := ""
	for s := range ch {
		result += s + " "
	}
	return result
}

func RangeStructChannel() string {
	type Point struct{ X, Y int }
	ch := make(chan Point, 2)
	ch <- Point{1, 2}
	ch <- Point{3, 4}
	close(ch)

	result := ""
	for p := range ch {
		result += fmt.Sprintf("(%d,%d)", p.X, p.Y)
	}
	return result
}

func PartialReceiveThenRange() string {
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)

	v1 := <-ch
	v2 := <-ch

	result := fmt.Sprintf("manual: %d,%d; range:", v1, v2)
	for val := range ch {
		result += fmt.Sprintf("%d", val)
	}
	return result
}

func CloseSentinelPattern() string {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		ch <- 1
		ch <- 2
		close(done)
	}()

	<-done
	close(ch)

	result := ""
	for val := range ch {
		result += fmt.Sprintf("%d", val)
	}
	return result
}
