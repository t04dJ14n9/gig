package divergence_hunt46

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 46: Channel patterns - buffered, unbuffered, select, close
// ============================================================================

func BufferedChannelSendRecv() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	return <-ch + <-ch + <-ch
}

func BufferedChannelLenCap() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return len(ch)*10 + cap(ch)
}

func ChannelCloseAndRange() int {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func ChannelRecvAfterClose() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	v3, ok3 := <-ch
	result := v1 + v2 + v3
	if ok1 { result += 10 }
	if ok2 { result += 10 }
	if ok3 { result += 10 } // ok3 is false
	return result
}

func SelectTwoChannels() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 42
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
	case <-ch:
		return 1
	default:
		return 2
	}
}

func SelectNilChannel() int {
	ch := make(chan int, 1)
	ch <- 42
	var nilCh chan int
	select {
	case v := <-ch:
		return v
	case <-nilCh:
		return -1
	default:
		return 99
	}
}

func ChannelDirection() int {
	// Test that channels work as function arguments
	send := func(ch chan<- int, v int) { ch <- v }
	recv := func(ch <-chan int) int { return <-ch }
	ch := make(chan int, 1)
	send(ch, 42)
	return recv(ch)
}

func ChannelAsSignal() int {
	done := make(chan struct{}, 1)
	done <- struct{}{}
	<-done
	return 1
}

func ChannelStruct() int {
	type Msg struct{ Type string; Value int }
	ch := make(chan Msg, 1)
	ch <- Msg{Type: "test", Value: 42}
	msg := <-ch
	return msg.Value
}

func ChannelSlice() int {
	ch := make(chan []int, 1)
	ch <- []int{1, 2, 3}
	s := <-ch
	return s[0] + s[1] + s[2]
}

func ChannelMap() int {
	ch := make(chan map[string]int, 1)
	ch <- map[string]int{"a": 1}
	m := <-ch
	return m["a"]
}

func JSONThroughChannel() int {
	type Data struct{ Value int `json:"value"` }
	ch := make(chan Data, 1)
	d := Data{Value: 42}
	b, _ := json.Marshal(d)
	var decoded Data
	json.Unmarshal(b, &decoded)
	ch <- decoded
	return (<-ch).Value
}

func MultipleSelects() int {
	ch := make(chan int, 1)
	result := 0
	ch <- 10
	select {
	case v := <-ch: result += v
	default:
	}
	ch <- 20
	select {
	case v := <-ch: result += v
	default:
	}
	return result
}

func FmtChannel() string {
	ch := make(chan int, 1)
	return fmt.Sprintf("%d", cap(ch))
}

func SortThroughChannel() int {
	ch := make(chan []int, 1)
	data := []int{3, 1, 2}
	sort.Ints(data)
	ch <- data
	s := <-ch
	return s[0]*100 + s[1]*10 + s[2]
}

func StringsThroughChannel() int {
	ch := make(chan string, 1)
	ch <- strings.ToUpper("hello")
	return len(<-ch)
}
