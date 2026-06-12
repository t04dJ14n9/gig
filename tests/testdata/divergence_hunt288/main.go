package divergence_hunt288

import (
	"fmt"
)

// ============================================================================
// Round 288: Channel edge cases — buffered, closed, direction, nil channel

// BufferedChannelSendRecv tests buffered channel send/recv
func BufferedChannelSendRecv() string {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	a := <-ch
	b := <-ch
	return fmt.Sprintf("a=%d,b=%d", a, b)
}

// CloseChannelThenRecv tests receiving from closed channel returns zero value
func CloseChannelThenRecv() string {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v, ok := <-ch
	v2, ok2 := <-ch
	return fmt.Sprintf("v=%d,ok=%t,v2=%d,ok2=%t", v, ok, v2, ok2)
}

// ChannelLenCap tests len and cap on buffered channel
func ChannelLenCap() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	return fmt.Sprintf("len=%d,cap=%d", len(ch), cap(ch))
}

// RangeOverChannel tests ranging over channel until closed
func RangeOverChannel() string {
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

// NilChannelBlocksForever tests nil channel select (default path taken)
func NilChannelBlocksForever() string {
	var ch chan int
	select {
	case <-ch:
		return "received"
	default:
		return "default"
	}
}

// SendOnClosedChannelPanics tests sending on closed channel panics
func SendOnClosedChannelPanics() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	ch := make(chan int, 1)
	close(ch)
	ch <- 1
	return "no_panic"
}

// ChannelOfString tests channel of string type
func ChannelOfString() string {
	ch := make(chan string, 1)
	ch <- "hello"
	v := <-ch
	return v
}

// ChannelOfStruct tests channel of struct type
func ChannelOfStruct() string {
	type Msg struct {
		Code int
		Text string
	}
	ch := make(chan Msg, 1)
	ch <- Msg{Code: 200, Text: "OK"}
	v := <-ch
	return fmt.Sprintf("code=%d,text=%s", v.Code, v.Text)
}

// SelectWithMultipleReady tests select when multiple channels ready
func SelectWithMultipleReady() string {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	// Both ready, select picks one (non-deterministic, but one will be picked)
	select {
	case v := <-ch1:
		return fmt.Sprintf("ready:%t", v == 1)
	case v := <-ch2:
		return fmt.Sprintf("ready:%t", v == 2)
	}
}

// BidirectionalChannelAsSend tests using bidirectional channel
func BidirectionalChannelAsSend() string {
	ch := make(chan int, 1)
	ch <- 42
	v := <-ch
	return fmt.Sprintf("v=%d", v)
}
