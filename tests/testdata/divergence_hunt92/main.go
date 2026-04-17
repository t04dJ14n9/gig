package divergence_hunt92

import "fmt"

// ============================================================================
// Round 92: Channel operations - select, buffered, closing
// ============================================================================

func BufferedSendRecv() int {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	a := <-ch
	b := <-ch
	c := <-ch
	return a + b + c
}

func BufferedLenCap() string {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return fmt.Sprintf("%d:%d", len(ch), cap(ch))
}

func CloseChannel() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func CloseAndRecv() string {
	ch := make(chan int, 2)
	ch <- 42
	close(ch)
	a, ok1 := <-ch
	b, ok2 := <-ch
	return fmt.Sprintf("%d:%v:%v:%v", a, b, ok1, ok2)
}

func SelectBasic() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 100
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	}
}

func SelectDefault() int {
	ch := make(chan int)
	select {
	case v := <-ch:
		return v
	default:
		return -1
	}
}

func ChannelNilBlock() string {
	var ch chan int
	return fmt.Sprintf("%v", ch == nil)
}

func NilChannelSelect() int {
	var ch chan int
	select {
	case <-ch:
		return 1
	default:
		return 2
	}
}

func BufferedStringChan() string {
	ch := make(chan string, 2)
	ch <- "hello"
	ch <- "world"
	a := <-ch
	b := <-ch
	return a + " " + b
}

func ChannelOfStruct() string {
	type Msg struct {
		ID   int
		Text string
	}
	ch := make(chan Msg, 1)
	ch <- Msg{ID: 1, Text: "hi"}
	m := <-ch
	return fmt.Sprintf("%d:%s", m.ID, m.Text)
}

func ChannelDirection() int {
	producer := func(ch chan<- int) {
		ch <- 42
	}
	consumer := func(ch <-chan int) int {
		return <-ch
	}
	ch := make(chan int, 1)
	producer(ch)
	return consumer(ch)
}

func SelectMultipleReady() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	ch2 <- 20
	total := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			total += v
		case v := <-ch2:
			total += v
		}
	}
	return total
}

func CloseRangeSum() int {
	ch := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}
