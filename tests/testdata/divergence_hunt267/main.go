package divergence_hunt267

import (
	"fmt"
)

// ============================================================================
// Round 267: Channel edge cases — nil channel, closed channel, select
// ============================================================================

// ChannelSendRecv tests basic send and receive
func ChannelSendRecv() string {
	ch := make(chan int, 1)
	ch <- 42
	v := <-ch
	return fmt.Sprintf("v=%d", v)
}

// ChannelClosedRecv tests receiving from closed channel returns zero
func ChannelClosedRecv() string {
	ch := make(chan int, 1)
	ch <- 10
	close(ch)
	v1 := <-ch
	v2, ok := <-ch
	return fmt.Sprintf("v1=%d,v2=%d,ok=%t", v1, v2, ok)
}

// ChannelBuffered tests buffered channel
func ChannelBuffered() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	result := ""
	for i := 0; i < 3; i++ {
		v := <-ch
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", v)
	}
	return result
}

// ChannelLenCap tests len and cap on buffered channel
func ChannelLenCap() string {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return fmt.Sprintf("len=%d,cap=%d", len(ch), cap(ch))
}

// ChannelSelectDefault tests select with default case
func ChannelSelectDefault() string {
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		return fmt.Sprintf("got:%d", v)
	default:
		return "nothing"
	}
}

// ChannelSelectReady tests select with ready channel
func ChannelSelectReady() string {
	ch := make(chan int, 1)
	ch <- 99
	select {
	case v := <-ch:
		return fmt.Sprintf("got:%d", v)
	default:
		return "nothing"
	}
}

// ChannelNilBlock tests nil channel blocks (using select to avoid deadlock)
func ChannelNilBlock() string {
	var nilCh chan int
	ch := make(chan int, 1)
	ch <- 42
	select {
	case v := <-ch:
		return fmt.Sprintf("v=%d", v)
	case <-nilCh:
		return "nil"
	}
}

// ChannelCloseRange tests ranging over channel until closed
func ChannelCloseRange() string {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ChannelDirectionSend tests send-only channel direction
func ChannelDirectionSend() string {
	ch := make(chan int, 1)
	sendOnly := func(c chan<- int) { c <- 7 }
	sendOnly(ch)
	v := <-ch
	return fmt.Sprintf("v=%d", v)
}

// ChannelDirectionRecv tests receive-only channel direction
func ChannelDirectionRecv() string {
	ch := make(chan int, 1)
	ch <- 99
	recvOnly := func(c <-chan int) int { return <-c }
	v := recvOnly(ch)
	return fmt.Sprintf("v=%d", v)
}
