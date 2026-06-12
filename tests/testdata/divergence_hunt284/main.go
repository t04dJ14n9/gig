package divergence_hunt284

import (
	"fmt"
)

// ============================================================================
// Round 284: Slice and map len/cap edge cases, make vs literal

// MakeSliceLenCap tests make([]int, len, cap)
func MakeSliceLenCap() string {
	s := make([]int, 3, 5)
	return fmt.Sprintf("len=%d,cap=%d,val=%v", len(s), cap(s), s)
}

// MakeSliceLenOnly tests make([]int, len)
func MakeSliceLenOnly() string {
	s := make([]int, 3)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// MakeMap tests make(map[string]int)
func MakeMap() string {
	m := make(map[string]int)
	m["key"] = 42
	return fmt.Sprintf("len=%d,val=%d", len(m), m["key"])
}

// MakeMapWithCapacity tests make(map[string]int, cap)
func MakeMapWithCapacity() string {
	m := make(map[string]int, 10)
	m["a"] = 1
	return fmt.Sprintf("len=%d", len(m))
}

// SliceLiteralCreatesNewArray tests that slice literal creates new array each time
func SliceLiteralCreatesNewArray() string {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s1[0] = 99
	return fmt.Sprintf("s1=%v,s2=%v", s1, s2)
}

// ArrayValueCopy tests array assignment copies values
func ArrayValueCopy() string {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 99
	return fmt.Sprintf("a=%v,b=%v", a, b)
}

// SliceAppendGrowsCap tests append growth strategy
func SliceAppendGrowsCap() string {
	s := make([]int, 0, 1)
	s = append(s, 1)
	cap1 := cap(s)
	s = append(s, 2)
	cap2 := cap(s)
	return fmt.Sprintf("cap1=%d,cap2=%d", cap1, cap2)
}

// NilSliceAppendReturnsNew tests append to nil slice
func NilSliceAppendReturnsNew() string {
	var s []int
	s2 := append(s, 1, 2, 3)
	return fmt.Sprintf("s_nil=%t,s2=%v", s == nil, s2)
}

// MapDeleteAndLen tests delete doesn't leave phantom entries
func MapDeleteAndLen() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return fmt.Sprintf("len=%d", len(m))
}

// SliceOfInterface tests []interface{} holding different types
func SliceOfInterface() string {
	s := []interface{}{1, "two", 3.0, true}
	return fmt.Sprintf("%v,%v,%v,%v", s[0], s[1], s[2], s[3])
}

// MapIterationCount tests counting map iterations
func MapIterationCount() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for range m {
		count++
	}
	return fmt.Sprintf("count=%d", count)
}
