package divergence_hunt188

import (
	"fmt"
)

// ============================================================================
// Round 188: Channel directionality (send-only, receive-only)
// ============================================================================

func sendOnly(ch chan<- int, val int) {
	ch <- val
}

func receiveOnly(ch <-chan int) int {
	return <-ch
}

func ChannelBidirectional() string {
	ch := make(chan int, 1)
	ch <- 42
	val := <-ch
	return fmt.Sprintf("%d", val)
}

func ChannelSendOnly() string {
	ch := make(chan int, 1)
	sendOnly(ch, 99)
	val := <-ch
	return fmt.Sprintf("%d", val)
}

func ChannelReceiveOnly() string {
	ch := make(chan int, 1)
	ch <- 77
	val := receiveOnly(ch)
	return fmt.Sprintf("%d", val)
}

func ChannelBufferedFull() string {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	return fmt.Sprintf("%d:%d", len(ch), cap(ch))
}

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
	return fmt.Sprintf("%d", sum)
}

func ChannelCloseCheck() string {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v, ok := <-ch
	return fmt.Sprintf("%d:%v", v, ok)
}

func ChannelSelectDefault() string {
	ch := make(chan int)
	select {
	case <-ch:
		return fmt.Sprintf("received")
	default:
		return fmt.Sprintf("default")
	}
}

func ChannelNil() string {
	var ch chan int
	return fmt.Sprintf("%v", ch == nil)
}

func ChannelMakeZero() string {
	ch := make(chan int)
	return fmt.Sprintf("%d", cap(ch))
}

func ChannelSendReceiveOrder() string {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	v1 := <-ch
	v2 := <-ch
	return fmt.Sprintf("%d:%d", v1, v2)
}
