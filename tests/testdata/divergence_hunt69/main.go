package divergence_hunt69

// ============================================================================
// Round 69: Channel/select edge cases - nil channel, closed channel, select
// ============================================================================

func ChannelBasic() int {
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

func ChannelClose() int {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)
	v1 := <-ch
	v2 := <-ch
	return v1 + v2
}

func ChannelClosedReadZero() int {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	<-ch
	v, ok := <-ch
	if ok {
		return v
	}
	return -1
}

func ChannelLen() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return len(ch)
}

func ChannelCap() int {
	ch := make(chan int, 5)
	return cap(ch)
}

func SelectBasic() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
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

func ChannelCloseAndRange() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func NilChannelBlocks() int {
	var ch chan int // nil channel
	// Reading from nil channel blocks forever, so we check it's nil
	if ch == nil {
		return 0
	}
	return -1
}

func ChannelDirection() int {
	producer := func(ch chan<- int) {
		ch <- 42
	}
	ch := make(chan int, 1)
	producer(ch)
	return <-ch
}

func ChannelSelectMultiple() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 100
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v * 2
	}
}
