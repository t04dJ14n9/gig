package divergence_hunt211

import (
	"fmt"
	"time"
)

func SelectWithTimeout() string {
	ch := make(chan string)
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "value"
	}()

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(100 * time.Millisecond):
		return "timeout"
	}
}

func SelectTimeoutNotTriggered() string {
	ch := make(chan string, 1)
	ch <- "quick"

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(200 * time.Millisecond):
		return "timeout"
	}
}

func SelectMultipleWithTimeout() string {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch1 <- "first"
	}()
	go func() {
		time.Sleep(60 * time.Millisecond)
		ch2 <- "second"
	}()

	select {
	case val := <-ch1:
		return fmt.Sprintf("ch1: %s", val)
	case val := <-ch2:
		return fmt.Sprintf("ch2: %s", val)
	case <-time.After(100 * time.Millisecond):
		return "timeout"
	}
}

func SelectDefaultWithTimeout() string {
	ch := make(chan string)

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(10 * time.Millisecond):
		return "timeout"
	default:
		return "default case"
	}
}

func NestedSelectTimeout() string {
	outerCh := make(chan string)
	innerCh := make(chan string)

	go func() {
		time.Sleep(20 * time.Millisecond)
		outerCh <- "outer"
	}()

	select {
	case <-outerCh:
		go func() {
			time.Sleep(20 * time.Millisecond)
			innerCh <- "inner"
		}()
		select {
		case val := <-innerCh:
			return fmt.Sprintf("nested: %s", val)
		case <-time.After(50 * time.Millisecond):
			return "inner timeout"
		}
	case <-time.After(100 * time.Millisecond):
		return "outer timeout"
	}
}

func SelectTimeoutChannelClosed() string {
	ch := make(chan string)
	close(ch)

	select {
	case val, ok := <-ch:
		return fmt.Sprintf("received: %s, ok: %v", val, ok)
	case <-time.After(10 * time.Millisecond):
		return "timeout"
	}
}

func SelectTimeoutZeroDuration() string {
	ch := make(chan string)

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(0):
		return "zero timeout"
	}
}

func BufferedChannelSelectTimeout() string {
	ch := make(chan string, 2)
	ch <- "one"
	ch <- "two"

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(50 * time.Millisecond):
		return "timeout"
	}
}

func SelectTimeoutReuse() string {
	timeout := time.After(100 * time.Millisecond)
	ch := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- "value"
	}()

	select {
	case val := <-ch:
		return fmt.Sprintf("first: %s", val)
	case <-timeout:
		return "first timeout"
	}
}

func LongTimeoutNotTriggered() string {
	ch := make(chan string)

	go func() {
		ch <- "quick"
	}()

	select {
	case val := <-ch:
		return fmt.Sprintf("received: %s", val)
	case <-time.After(5 * time.Second):
		return "long timeout"
	}
}
