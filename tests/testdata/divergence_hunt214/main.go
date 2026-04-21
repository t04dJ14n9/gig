package divergence_hunt214

import (
	"fmt"
	"time"
)

func NilChannelIsNil() string {
	var ch chan int
	return fmt.Sprintf("nil: %v", ch == nil)
}

func NilChannelSendBlock() string {
	var ch chan int

	done := make(chan string)
	go func() {
		ch <- 1
		done <- "sent"
	}()

	select {
	case <-done:
		return "unexpected"
	case <-time.After(50 * time.Millisecond):
		return "blocked on nil send"
	}
}

func NilChannelReceiveBlock() string {
	var ch chan int

	done := make(chan string)
	go func() {
		<-ch
		done <- "received"
	}()

	select {
	case <-done:
		return "unexpected"
	case <-time.After(50 * time.Millisecond):
		return "blocked on nil receive"
	}
}

func NilChannelInSelect() string {
	var nilCh chan int
	ch := make(chan int, 1)
	ch <- 42

	select {
	case <-nilCh:
		return "nil case"
	case val := <-ch:
		return fmt.Sprintf("got %d", val)
	}
}

func NilChannelClosePanic() string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from nil close")
		}
	}()

	var ch chan int
	close(ch)
	return "no panic"
}

func NilChannelLenCap() string {
	var ch chan int
	return fmt.Sprintf("len: %d, cap: %d", len(ch), cap(ch))
}

func NilChannelTypeAssertion() string {
	var ch interface{} = (chan int)(nil)

	if intCh, ok := ch.(chan int); ok {
		return fmt.Sprintf("is chan int, nil: %v", intCh == nil)
	}
	return "type assertion failed"
}

func NilChannelAssignAndUse() string {
	var ch1 chan string
	ch2 := ch1
	_ = ch2 // ch2 is nil too

	ch3 := make(chan string, 1)
	ch3 <- "test"

	ch1 = ch3
	return fmt.Sprintf("was nil, now: %v", <-ch1)
}

func NilChannelComparison() string {
	var ch1 chan int
	var ch2 chan int
	ch3 := make(chan int)

	return fmt.Sprintf("nil==nil: %v, nil!=real: %v", ch1 == ch2, ch1 != ch3)
}

func NilChannelInSlice() string {
	channels := []chan int{nil, make(chan int, 1), nil}
	channels[1] <- 42

	count := 0
	for _, ch := range channels {
		if ch == nil {
			count++
		}
	}
	return fmt.Sprintf("nil count: %d", count)
}

func NilChannelSelectMultiple() string {
	var nilCh1, nilCh2 chan int
	ch := make(chan int, 1)
	ch <- 1

	select {
	case <-nilCh1:
		return "nil1"
	case <-nilCh2:
		return "nil2"
	case val := <-ch:
		return fmt.Sprintf("got %d", val)
	default:
		return "default"
	}
}
