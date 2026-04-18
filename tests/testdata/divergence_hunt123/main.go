package divergence_hunt123

import "fmt"

// ============================================================================
// Round 123: Slice tricks — append, copy, delete, three-index
// ============================================================================

func SliceAppendNil() string {
	var s []int
	s = append(s, 1, 2, 3)
	return fmt.Sprintf("%v", s)
}

func SliceAppendExpand() string {
	s := []int{1, 2}
	s = append(s, 3, 4, 5)
	return fmt.Sprintf("%v", s)
}

func SliceCopyCount() string {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return fmt.Sprintf("n=%d-dst=%v", n, dst)
}

func SliceCopyOverlap() string {
	s := []int{1, 2, 3, 4, 5}
	copy(s[0:], s[2:])
	return fmt.Sprintf("%v", s)
}

func SliceDeleteElement() string {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return fmt.Sprintf("%v", s)
}

func SliceThreeIndex() string {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:4] // len=2, cap=3
	return fmt.Sprintf("%v-len=%d-cap=%d", sub, len(sub), cap(sub))
}

func SliceNilAppend() string {
	var s []int
	result := append(s, 42)
	return fmt.Sprintf("%v", result)
}

func SliceNilCopy() string {
	var src []int
	dst := make([]int, 3)
	n := copy(dst, src)
	return fmt.Sprintf("n=%d", n)
}

func SliceAppendSlice() string {
	a := []int{1, 2}
	b := []int{3, 4, 5}
	c := append(a, b...)
	return fmt.Sprintf("%v", c)
}

func SliceCapAfterAppend() string {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	// Next append should double capacity
	before := cap(s)
	s = append(s, 3)
	after := cap(s)
	return fmt.Sprintf("before=%d-after=%d", before, after)
}
