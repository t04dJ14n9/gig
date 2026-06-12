package divergence_hunt207

import (
	"fmt"
)

// ============================================================================
// Round 207: Empty struct and zero-size types
// ============================================================================

type Empty207 struct{}

type ZeroArray207 [0]int

type EmptyWithMethod207 struct{}

func (e EmptyWithMethod207) Method() string {
	return "empty method"
}

type HasEmpty207 struct {
	X     int
	Empty Empty207
	Y     int
}

// EmptyStructSize tests empty struct size
func EmptyStructSize() string {
	var e Empty207
	_ = e
	// Cannot use unsafe.Sizeof, so we test behavior instead
	return fmt.Sprintf("empty-struct-ok")
}

// EmptyStructEquality tests empty struct equality
func EmptyStructEquality() string {
	var a, b Empty207
	return fmt.Sprintf("%v", a == b)
}

// ZeroSizeArray tests zero-size array
func ZeroSizeArray() string {
	var a ZeroArray207
	return fmt.Sprintf("len:%d,cap:%d", len(a), cap(a))
}

// EmptyStructSlice tests empty struct slice
func EmptyStructSlice() string {
	s := make([]Empty207, 100)
	return fmt.Sprintf("len:%d,cap:%d", len(s), cap(s))
}

// EmptyStructMap tests empty struct map
func EmptyStructMap() string {
	m := make(map[int]Empty207)
	m[1] = Empty207{}
	m[2] = Empty207{}
	return fmt.Sprintf("len:%d", len(m))
}

// EmptyStructChannel tests empty struct channel
func EmptyStructChannel() string {
	ch := make(chan Empty207, 10)
	ch <- Empty207{}
	ch <- Empty207{}
	<-ch
	return fmt.Sprintf("len:%d,cap:%d", len(ch), cap(ch))
}

// EmptyStructMethod tests empty struct with method
func EmptyStructMethod() string {
	var e EmptyWithMethod207
	return e.Method()
}

// EmptyStructInStruct tests empty struct embedded in struct
func EmptyStructInStruct() string {
	h := HasEmpty207{X: 1, Y: 2}
	return fmt.Sprintf("X=%d,Y=%d", h.X, h.Y)
}

// EmptyStructPointer tests pointer to empty struct
func EmptyStructPointer() string {
	e := &Empty207{}
	_ = e
	return fmt.Sprintf("ptr-ok")
}

// EmptyStructInterface tests empty struct in interface
func EmptyStructInterface() string {
	var i interface{} = Empty207{}
	_, ok := i.(Empty207)
	return fmt.Sprintf("ok=%v", ok)
}
