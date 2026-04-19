package divergence_hunt151

import "fmt"

// ============================================================================
// Round 151: Select statements and channel operations
// ============================================================================

// SelectBasic tests basic select with default
func SelectBasic() string {
	ch := make(chan int, 1)
	ch <- 42

	select {
	case v := <-ch:
		return fmt.Sprintf("received-%d", v)
	default:
		return "default"
	}
}

// SelectDefaultOnly tests select with only default case
func SelectDefaultOnly() string {
	select {
	default:
		return "default-case"
	}
}

// SelectNoDefault tests blocking select without default
func SelectNoDefault() string {
	ch := make(chan int, 1)
	ch <- 100

	var result int
	select {
	case result = <-ch:
	}
	return fmt.Sprintf("got-%d", result)
}

// SelectMultipleChannels tests select with multiple channels
func SelectMultipleChannels() string {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	ch1 <- "first"

	select {
	case v := <-ch1:
		return fmt.Sprintf("ch1-%s", v)
	case v := <-ch2:
		return fmt.Sprintf("ch2-%s", v)
	default:
		return "default"
	}
}

// ChannelBufferedCapacity tests buffered channel capacity
func ChannelBufferedCapacity() string {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	return fmt.Sprintf("len=%d-cap=%d", len(ch), cap(ch))
}

// ChannelNilSend tests send on nil channel (should block forever, but we use select)
func ChannelNilSend() string {
	var ch chan int

	select {
	case ch <- 42:
		return "sent"
	default:
		return "cannot-send-nil"
	}
}

// ChannelNilReceive tests receive from nil channel (should block forever, but we use select)
func ChannelNilReceive() string {
	var ch chan int

	select {
	case <-ch:
		return "received"
	default:
		return "cannot-recv-nil"
	}
}

// ChannelCloseCheck tests checking if channel is closed
func ChannelCloseCheck() string {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)

	v, ok := <-ch
	if ok {
		return fmt.Sprintf("open-%d", v)
	}
	return "closed"
}

// ChannelRange tests ranging over channel
func ChannelRange() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	sum := 0
	for v := range ch {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ChannelSendReceiveOrder tests send and receive order
func ChannelSendReceiveOrder() string {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30

	v1 := <-ch
	ch <- 40
	v2 := <-ch
	return fmt.Sprintf("v1=%d-v2=%d", v1, v2)
}
