package divergence_hunt28

// ============================================================================
// Round 28: Channel patterns, select patterns, goroutine-free concurrency
// ============================================================================

func ChannelSendRecv() int {
	ch := make(chan int, 1)
	ch <- 42
	return <-ch
}

func ChannelBuffered() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	return <-ch + <-ch + <-ch
}

func ChannelCloseRange() int {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	sum := 0
	for v := range ch { sum += v }
	return sum
}

func ChannelSelectTwo() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	select {
	case v := <-ch1: return v * 10
	case v := <-ch2: return v * 20
	}
}

func ChannelSelectDefault2() int {
	ch := make(chan int)
	select {
	case v := <-ch: return v
	default: return -1
	}
}

func ChannelNilSelect() int {
	var nilCh chan int
	ch := make(chan int, 1)
	ch <- 42
	select {
	case v := <-ch: return v
	case <-nilCh: return -1
	default: return -2
	}
}

func ChannelLen() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return len(ch)
}

func ChannelCap2() int {
	ch := make(chan int, 5)
	return cap(ch)
}

func ChannelRecvAfterClose() int {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	_, ok3 := <-ch
	result := v1 + v2
	if ok1 { result += 100 }
	if ok2 { result += 200 }
	if ok3 { result += 300 }
	return result
}

func ChannelDirection() int {
	// Test sending and receiving through a channel in the same function
	ch := make(chan int, 1)
	send := func(c chan<- int, v int) { c <- v }
	recv := func(c <-chan int) int { return <-c }
	send(ch, 42)
	return recv(ch)
}

func SelectMultipleReady() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	ch2 <- 20
	sum := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1: sum += v
		case v := <-ch2: sum += v
		}
	}
	return sum
}

func ChannelAsSignal() int {
	done := make(chan struct{}, 1)
	done <- struct{}{}
	<-done
	return 1
}
