package divergence_hunt215

import (
	"fmt"
)

func SendOnlyChannel() string {
	ch := make(chan<- int, 2)
	ch <- 1
	ch <- 2
	return "sent to send-only"
}

func ReceiveOnlyChannel() string {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20

	recvCh := (<-chan int)(ch)
	v1 := <-recvCh
	v2 := <-recvCh
	return fmt.Sprintf("received: %d, %d", v1, v2)
}

func BidirectionalAsSendOnly() string {
	ch := make(chan int, 1)

	var sendOnly chan<- int = ch
	sendOnly <- 42

	return fmt.Sprintf("received from bidir: %d", <-ch)
}

func BidirectionalAsReceiveOnly() string {
	ch := make(chan int, 1)
	ch <- 99

	var recvOnly <-chan int = ch
	val := <-recvOnly
	return fmt.Sprintf("received: %d", val)
}

func DirectionConversion() string {
	ch := make(chan int, 1)

	var sendCh chan<- int = ch
	sendCh <- 1

	var recvCh <-chan int = ch
	val := <-recvCh

	return fmt.Sprintf("converted: %d", val)
}

func FunctionWithDirectionalChannels() string {
	sender := func(ch chan<- int, val int) {
		ch <- val
	}

	receiver := func(ch <-chan int) int {
		return <-ch
	}

	ch := make(chan int, 1)
	sender(ch, 42)
	val := receiver(ch)

	return fmt.Sprintf("sent and received: %d", val)
}

func PipelinePattern() string {
	gen := func() <-chan int {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)
		return ch
	}

	square := func(in <-chan int) <-chan int {
		out := make(chan int, 3)
		for n := range in {
			out <- n * n
		}
		close(out)
		return out
	}

	result := ""
	for n := range square(gen()) {
		result += fmt.Sprintf("%d", n)
	}
	return result
}

func CannotReceiveFromSendOnly() string {
	ch := make(chan int, 1)
	ch <- 1

	var sendOnly chan<- int = ch
	_ = sendOnly

	return "send-only channel cannot receive"
}

func CannotSendToReceiveOnly() string {
	ch := make(chan int, 1)
	ch <- 1

	var recvOnly <-chan int = ch
	_ = recvOnly

	return "receive-only channel cannot send"
}

func DirectionalChannelInStruct() string {
	type Pipe struct {
		In  chan<- int
		Out <-chan int
	}

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	p := Pipe{In: ch1, Out: ch2}
	p.In <- 42
	ch1 <- 42
	<-ch1

	ch2 <- 99
	val := <-p.Out

	return fmt.Sprintf("struct directional: %d", val)
}

func CloseSendOnlyChannel() string {
	ch := make(chan int, 1)
	var sendCh chan<- int = ch

	close(sendCh)
	_, ok := <-ch
	return fmt.Sprintf("closed ok: %v", !ok)
}
