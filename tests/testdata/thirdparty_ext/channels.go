package thirdparty

import (
	"bytes"
	"io"
)

// ============================================================================
// CHANNEL WITH EXTERNAL TYPES AS ELEMENTS
// ============================================================================

// ChannelWithStructElement tests channel with struct.
func ChannelWithStructElement() int {
	type Point struct{ X, Y int }
	ch := make(chan Point, 1)
	ch <- Point{X: 10, Y: 20}
	p := <-ch
	return p.X + p.Y
}

// ChannelWithInterfaceElement tests channel with interface.
func ChannelWithInterfaceElement() int {
	ch := make(chan io.Writer, 1)
	ch <- &bytes.Buffer{}
	w := <-ch
	if w != nil {
		return 1
	}
	return 0
}

// ChannelWithMapElement tests channel with map.
func ChannelWithMapElement() int {
	ch := make(chan map[string]int, 1)
	ch <- map[string]int{"a": 1}
	m := <-ch
	return m["a"]
}

// ChannelWithSliceElement tests channel with slice.
func ChannelWithSliceElement() int {
	ch := make(chan []int, 1)
	ch <- []int{1, 2, 3}
	s := <-ch
	return len(s)
}

// ChannelWithFuncElement tests channel with function.
func ChannelWithFuncElement() int {
	ch := make(chan func() int, 1)
	ch <- func() int { return 42 }
	fn := <-ch
	return fn()
}

// ============================================================================
// SELECT WITH DEFAULT AND EXTERNAL OPERATIONS
// ============================================================================

// SelectDefaultAlways tests select with default case.
func SelectDefaultAlways() int {
	ch := make(chan int)
	select {
	case <-ch:
		return 0
	default:
		return 1
	}
}

// SelectNonBlockingSend tests non-blocking send.
func SelectNonBlockingSend() int {
	ch := make(chan int, 1)
	sent := false
	select {
	case ch <- 42:
		sent = true
	default:
	}
	if sent {
		return <-ch
	}
	return 0
}

// SelectMultipleChannels tests select with multiple channels.
func SelectMultipleChannels() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	ch2 <- 20

	sum := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			sum += v
		case v := <-ch2:
			sum += v
		}
	}
	return sum
}
