package divergence_hunt19

import "fmt"

// ============================================================================
// Round 19: Edge cases - empty collections, zero values, nil handling,
// boundary conditions, error handling patterns
// ============================================================================

// EmptySliceOperations tests operations on empty slice
func EmptySliceOperations() int {
	var s []int
	return len(s) + cap(s)
}

// EmptyMapOperations tests operations on empty map
func EmptyMapOperations() int {
	var m map[string]int
	return len(m)
}

// EmptyStringOperations tests operations on empty string
func EmptyStringOperations() int {
	return len("") + len(" ")
}

// ZeroValueInt tests zero value int
func ZeroValueInt() int {
	var x int
	return x
}

// ZeroValueFloat tests zero value float
func ZeroValueFloat() float64 {
	var x float64
	return x
}

// ZeroValueBool tests zero value bool
func ZeroValueBool() bool {
	var x bool
	return x
}

// ZeroValueString tests zero value string
func ZeroValueString() string {
	var x string
	return x
}

// ZeroValueSlice tests zero value slice
func ZeroValueSlice() bool {
	var s []int
	return s == nil
}

// ZeroValueMap tests zero value map
func ZeroValueMap() bool {
	var m map[string]int
	return m == nil
}

// ZeroValuePointer tests zero value pointer
func ZeroValuePointer() bool {
	var p *int
	return p == nil
}

// NilSliceAppend tests appending to nil slice
func NilSliceAppend() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// NilMapRead tests reading from nil map
func NilMapRead() int {
	var m map[string]int
	return m["key"]
}

// NilSliceRange tests ranging over nil slice
func NilSliceRange() int {
	var s []int
	count := 0
	for range s { count++ }
	return count
}

// NilMapRange tests ranging over nil map
func NilMapRange() int {
	var m map[string]int
	count := 0
	for range m { count++ }
	return count
}

// NilChannelRead tests reading from nil channel (non-blocking via select)
func NilChannelRead() int {
	ch := make(chan int, 1)
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

// SliceBoundary tests slice boundary access
func SliceBoundary() int {
	s := []int{1, 2, 3}
	return s[0] + s[len(s)-1]
}

// MapBoundary tests map with boundary keys
func MapBoundary() int {
	m := map[int]string{0: "zero", -1: "neg"}
	return len(m[0]) + len(m[-1])
}

// ErrorHandlingPattern tests error handling pattern
func ErrorHandlingPattern() int {
	mightFail := func(ok bool) (int, error) {
		if ok { return 42, nil }
		return 0, fmt.Errorf("failed")
	}
	if v, err := mightFail(true); err == nil {
		return v
	}
	return -1
}

// MultipleErrorCheck tests multiple error checks
func MultipleErrorCheck() int {
	step1 := func() error { return nil }
	step2 := func() error { return nil }
	step3 := func() (int, error) { return 42, nil }
	if err := step1(); err != nil { return -1 }
	if err := step2(); err != nil { return -2 }
	if v, err := step3(); err == nil { return v }
	return -3
}

// NilFuncVariable tests nil function variable check
func NilFuncVariable() int {
	var f func() int
	if f == nil { return -1 }
	return f()
}

// EmptyInterfaceContains tests empty interface contains
func EmptyInterfaceContains() bool {
	var x any = nil
	return x == nil
}

// StructZeroValueFields tests struct zero value fields
func StructZeroValueFields() int {
	type S struct {
		X int
		Y float64
		Z bool
		W string
	}
	var s S
	if s.X == 0 && s.Y == 0.0 && !s.Z && s.W == "" {
		return 1
	}
	return 0
}
