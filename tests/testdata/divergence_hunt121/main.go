package divergence_hunt121

import "fmt"

// ============================================================================
// Round 121: Channel select and send/receive edge cases
// ============================================================================

func ChanSelectDefault() string {
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		return fmt.Sprintf("got %d", v)
	default:
		return "default"
	}
}

func ChanSelectReady() string {
	ch := make(chan int, 1)
	ch <- 42
	select {
	case v := <-ch:
		return fmt.Sprintf("got %d", v)
	default:
		return "default"
	}
}

func ChanBufferedSend() string {
	ch := make(chan string, 2)
	ch <- "hello"
	ch <- "world"
	close(ch)
	var result string
	for v := range ch {
		result += v
	}
	return result
}

func ChanNilBlock() string {
	var ch chan int
	// Reading from nil channel blocks forever, but select with default works
	select {
	case <-ch:
		return "received"
	default:
		return "nil-default"
	}
}

func ChanClosedReceive() string {
	ch := make(chan int, 1)
	ch <- 10
	close(ch)
	v, ok := <-ch
	return fmt.Sprintf("%d-%t", v, ok)
}

func ChanClosedEmpty() string {
	ch := make(chan int)
	close(ch)
	v, ok := <-ch
	return fmt.Sprintf("%d-%t", v, ok)
}

func ChanSelectMultiReady() string {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	// Both ready — Go picks one pseudo-randomly, so just check that one was picked
	select {
	case <-ch1:
		return "one-picked"
	case <-ch2:
		return "one-picked"
	}
}

func ChanLenCap() string {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	return fmt.Sprintf("len=%d-cap=%d", len(ch), cap(ch))
}

func ChanSelectWithAssign() string {
	ch := make(chan int, 2)
	ch <- 100
	ch <- 200
	sum := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch:
			sum += v
		}
	}
	return fmt.Sprintf("sum=%d", sum)
}
