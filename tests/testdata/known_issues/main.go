package known_issues

import (
	"container/heap"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// ============================================================================
// Resolved issue regression test data.
//
// This file contains test cases for interpreter bugs and limitations that are
// now expected to behave like native Go.
// ============================================================================

// --- Category 1: Interface with nil concrete type ---

type Stringer interface {
	String() string
}

type PointerReceiver struct {
	Name string
}

func (p *PointerReceiver) String() string {
	if p == nil {
		return "<nil>"
	}
	return p.Name
}

// InterfaceWithNilConcrete tests interface holding nil concrete type.
// In Go, an interface holding a nil pointer is NOT nil (it has a type but no value).
func InterfaceWithNilConcrete() any {
	var p *PointerReceiver
	var s Stringer = p
	if s == nil {
		return "interface is nil"
	}
	return s.String()
}

// NilInterfaceCall tests calling method on nil interface.
func NilInterfaceCall() any {
	var s Stringer
	defer func() {
		recover()
	}()
	_ = s.String()
	return "no panic"
}

// NestedNilReceiver tests nested struct with nil embedded pointer.
// Accessing promoted field on nil embedded pointer should return zero value.
func NestedNilReceiver() (result any) {
	type Inner struct {
		Name string
	}
	type Outer struct {
		*Inner
	}
	outer := Outer{}
	defer func() { recover() }()
	return outer.Name
}

// --- Category 2: Three-index slicing ---

// ThreeIndexByteSlice tests three-index slicing on byte slice.
func ThreeIndexByteSlice() any {
	b := []byte("hello world")
	b2 := b[0:5:11]
	return fmt.Sprintf("len=%d,cap=%d", len(b2), cap(b2))
}

// --- Category 3: sort.Interface callback dispatch ---

type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

// SortByLength tests sort.Sort with custom sort.Interface.
func SortByLength() any {
	words := []string{"apple", "pie", "banana", "kiwi"}
	sort.Sort(ByLength(words))
	return fmt.Sprintf("%v", words)
}

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

// SortStructByField tests sort.Sort with struct slice.
func SortStructByField() any {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Sort(ByAge(people))
	result := ""
	for _, p := range people {
		result += fmt.Sprintf("%s:%d ", p.Name, p.Age)
	}
	return result
}

// SortReverse tests sort.Reverse wrapper.
type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func SortReverse() any {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Sort(Reverse{sort.IntSlice(nums)})
	return fmt.Sprintf("%v", nums)
}

// Descending wraps sort.Interface like Reverse but deliberately uses a different
// type name so adapter dispatch cannot depend on a hard-coded receiver name.
type Descending struct {
	sort.Interface
}

func (d Descending) Less(i, j int) bool {
	return d.Interface.Less(j, i)
}

func SortEmbeddedInterfaceDescending() any {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Sort(Descending{sort.IntSlice(nums)})
	return fmt.Sprintf("%v", nums)
}

// --- Category 4: heap.Interface callback dispatch ---

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// HeapInit tests heap.Init with custom heap.Interface.
func HeapInit() any {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	return fmt.Sprintf("%v", *h)
}

// HeapPush tests heap.Push.
func HeapPush() any {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	return fmt.Sprintf("%v", *h)
}

// HeapPop tests heap.Pop.
func HeapPop() any {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	result := heap.Pop(h).(int)
	return fmt.Sprintf("%d:%v", result, *h)
}

// --- Category 5: errors.As with interface target ---

type CustomError struct {
	Code int
	Msg  string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// ErrorsAsInterface tests errors.As with interface target.
func ErrorsAsInterface() any {
	err := &CustomError{Code: 404, Msg: "not found"}
	var ce *CustomError
	if errors.As(err, &ce) {
		return ce.Error()
	}
	return "not matched"
}

// ErrorsAsNotMatching tests errors.As when target doesn't match.
func ErrorsAsNotMatching() any {
	err := fmt.Errorf("plain error")
	var ce *CustomError
	if errors.As(err, &ce) {
		return "matched"
	}
	return "not matched"
}

// ErrorsIsAndAsTogether tests errors.Is and errors.As combined.
func ErrorsIsAndAsTogether() any {
	base := fmt.Errorf("base error")
	wrapped := fmt.Errorf("wrapped: %w", base)
	err := &CustomError{Code: 500, Msg: "internal"}
	final := fmt.Errorf("%w: %v", wrapped, err)

	if errors.Is(final, base) {
		var ce *CustomError
		if errors.As(final, &ce) {
			return ce.Code
		}
	}
	return -1
}

// --- Category 6: io.MultiWriter callback ---

// MultiWriterTest tests io.MultiWriter.
func MultiWriterTest() any {
	buf1 := &strings.Builder{}
	buf2 := &strings.Builder{}
	mw := io.MultiWriter(buf1, buf2)
	mw.Write([]byte("hello"))
	return fmt.Sprintf("buf1=%s,buf2=%s", buf1.String(), buf2.String())
}

// --- Category 7: Sync primitives ---

// SyncMapRange tests sync.Map.Range callback.
func SyncMapRange() any {
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	count := 0
	m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// OnceWithPanic tests sync.Once with panic in initialization.
// If the function panics, subsequent calls should also panic.
func OnceWithPanic() any {
	var once sync.Once
	panicked := false

	func() {
		defer func() { recover() }()
		once.Do(func() { panic("init failed") })
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		once.Do(func() {})
	}()

	return panicked
}

// --- Category 8: Channel semantics ---

// SelectWithMultipleReady tests select with multiple ready channels.
// When multiple cases are ready, select should pick one randomly.
func SelectWithMultipleReady() any {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2
	ch3 <- 3

	counts := map[int]int{}
	for i := 0; i < 30; i++ {
		// Refill channels
		select {
		case v := <-ch1:
			counts[v]++
			ch1 <- 1
		case v := <-ch2:
			counts[v]++
			ch2 <- 2
		case v := <-ch3:
			counts[v]++
			ch3 <- 3
		}
	}

	// In native Go, all three channels should be selected roughly equally
	return len(counts)
}
