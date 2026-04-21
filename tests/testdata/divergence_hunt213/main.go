package divergence_hunt213

import (
	"fmt"
	"time"
)

func BufferedChannelCapacity() string {
	ch := make(chan int, 5)
	return fmt.Sprintf("capacity: %d", cap(ch))
}

func UnbufferedChannelCapacity() string {
	ch := make(chan int)
	return fmt.Sprintf("capacity: %d", cap(ch))
}

func BufferedChannelLength() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	return fmt.Sprintf("len: %d, cap: %d", len(ch), cap(ch))
}

func BufferedChannelBlocking() string {
	ch := make(chan int, 1)

	ch <- 1

	done := make(chan string)
	go func() {
		ch <- 2
		done <- "sent"
	}()

	select {
	case msg := <-done:
		return msg
	case <-time.After(50 * time.Millisecond):
		<-ch
		return "blocked then unblocked"
	}
}

func UnbufferedChannelSynchronization() string {
	ch := make(chan string)
	result := ""

	go func() {
		result += "receiver started;"
		val := <-ch
		result += fmt.Sprintf("received %s;", val)
	}()

	result += "sender waiting;"
	ch <- "hello"
	result += "sender done"

	time.Sleep(10 * time.Millisecond)
	return result
}

func BufferedVsUnbufferedSend() string {
	buf := make(chan int, 1)
	unbuf := make(chan int)

	buf <- 1

	done := make(chan string)
	go func() {
		unbuf <- 2
		done <- "unbuffered sent"
	}()

	select {
	case <-done:
		return "unexpected"
	case <-time.After(10 * time.Millisecond):
		<-unbuf
		return fmt.Sprintf("buffered ok, unbuffered blocked")
	}
}

func BufferedChannelDrainOrder() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	result := ""
	result += fmt.Sprintf("%d", <-ch)
	result += fmt.Sprintf("%d", <-ch)
	result += fmt.Sprintf("%d", <-ch)
	return result
}

func ZeroBufferChannel() string {
	ch := make(chan int, 0)
	return fmt.Sprintf("len: %d, cap: %d", len(ch), cap(ch))
}

func BufferedStringChannel() string {
	ch := make(chan string, 2)
	ch <- "first"
	ch <- "second"

	return fmt.Sprintf("%s-%s", <-ch, <-ch)
}

func BufferedChannelCloseSemantics() string {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)

	result := ""
	for v := range ch {
		result += fmt.Sprintf("%d", v)
	}
	return result
}

func MixedBufferedUnbuffered() string {
	buf := make(chan int, 1)
	unbuf := make(chan int)

	buf <- 1

	go func() {
		unbuf <- <-buf
	}()

	val := <-unbuf
	return fmt.Sprintf("transferred: %d", val)
}

func BufferedChannelOverwriteProtection() string {
	ch := make(chan int, 1)
	ch <- 42

	select {
	case ch <- 99:
		return "overwrote"
	default:
		return "blocked"
	}
}
