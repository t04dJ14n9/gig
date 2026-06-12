package divergence_hunt262

import (
	"fmt"
)

// ============================================================================
// Round 262: Slice tricks — three-index slicing, append behavior, copy
// ============================================================================

// ThreeIndexSlice tests three-index slice expression
func ThreeIndexSlice() string {
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3:4] // len=2, cap=3
	return fmt.Sprintf("len=%d,cap=%d,vals=%v", len(b), cap(b), b)
}

// ThreeIndexSliceFullCap tests three-index with full capacity
func ThreeIndexSliceFullCap() string {
	a := []int{10, 20, 30, 40}
	b := a[1:3:4] // len=2, cap=3
	return fmt.Sprintf("len=%d,cap=%d", len(b), cap(b))
}

// AppendExtendBeyondLen tests append extending past length but within cap
func AppendExtendBeyondLen() string {
	a := make([]int, 2, 5)
	a[0] = 10
	a[1] = 20
	b := append(a, 30)
	return fmt.Sprintf("a_len=%d,b_len=%d,b=%v", len(a), len(b), b)
}

// AppendCausesReallocation tests append triggering reallocation
func AppendCausesReallocation() string {
	a := make([]int, 3, 3)
	a[0], a[1], a[2] = 1, 2, 3
	b := append(a, 4) // new allocation
	b[0] = 99
	return fmt.Sprintf("a[0]=%d,b[0]=%d", a[0], b[0]) // a[0] should still be 1
}

// CopySliceBasic tests basic copy between slices
func CopySliceBasic() string {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	n := copy(dst, src)
	return fmt.Sprintf("n=%d,dst=%v", n, dst)
}

// CopySlicePartial tests copy with different lengths
func CopySlicePartial() string {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return fmt.Sprintf("n=%d,dst=%v", n, dst)
}

// SliceModifyAfterSlice tests modifying through sub-slice
func SliceModifyAfterSlice() string {
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3]
	b[0] = 99
	return fmt.Sprintf("a=%v", a) // a should be [1 99 3 4 5]
}

// SliceNilAppend tests appending to nil slice
func SliceNilAppend() string {
	var s []int
	s = append(s, 1, 2, 3)
	return fmt.Sprintf("s=%v,len=%d,cap=%d", s, len(s), cap(s))
}

// SliceGrowPattern tests repeated append growth
func SliceGrowPattern() string {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return fmt.Sprintf("len=%d,first=%d,last=%d", len(s), s[0], s[9])
}

// SliceDeleteElement tests deleting element from middle
func SliceDeleteElement() string {
	s := []int{1, 2, 3, 4, 5}
	i := 2 // remove element at index 2
	s = append(s[:i], s[i+1:]...)
	return fmt.Sprintf("s=%v,len=%d", s, len(s))
}

// SliceInsertElement tests inserting element at position
func SliceInsertElement() string {
	s := []int{1, 2, 4, 5}
	i := 2
	s = append(s[:i], append([]int{3}, s[i:]...)...)
	return fmt.Sprintf("s=%v", s)
}
