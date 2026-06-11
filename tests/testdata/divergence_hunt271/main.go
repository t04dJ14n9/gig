package divergence_hunt271

import (
	"fmt"
)

// ============================================================================
// Round 271: Slice tricks — 3-index slicing, append, nil vs empty, capacity

// ThreeIndexSlice tests 3-index slice expression [low:high:max]
func ThreeIndexSlice() string {
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3:3]
	return fmt.Sprintf("len=%d,cap=%d,val=%v", len(b), cap(b), b)
}

// ThreeIndexSliceFull tests a[::cap] 3-index
func ThreeIndexSliceFull() string {
	a := []int{1, 2, 3, 4, 5}
	b := a[0:3:5]
	return fmt.Sprintf("len=%d,cap=%d,val=%v", len(b), cap(b), b)
}

// AppendToNil tests appending to nil slice
func AppendToNil() string {
	var s []int
	s = append(s, 1, 2, 3)
	return fmt.Sprintf("len=%d,val=%v", len(s), s)
}

// NilSliceVsEmpty tests nil vs empty slice behavior
func NilSliceVsEmpty() string {
	var nilSlice []int
	emptySlice := []int{}
	return fmt.Sprintf("nil_len=%d,nil_cap=%d,nil_nil=%t,empty_len=%d,empty_cap=%d,empty_nil=%t",
		len(nilSlice), cap(nilSlice), nilSlice == nil,
		len(emptySlice), cap(emptySlice), emptySlice == nil)
}

// AppendToSliceWithCap tests append to slice with spare capacity
func AppendToSliceWithCap() string {
	a := make([]int, 2, 5)
	a[0] = 10
	a[1] = 20
	b := append(a, 30)
	b[0] = 99 // shares underlying array with a
	return fmt.Sprintf("a=%v,b=%v", a, b)
}

// SliceSharedBacking tests that slices share backing arrays
func SliceSharedBacking() string {
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3]
	b[0] = 99
	return fmt.Sprintf("a=%v,b=%v", a, b)
}

// SliceCopyIndependent tests that copy makes independent slice
func SliceCopyIndependent() string {
	a := []int{1, 2, 3}
	b := make([]int, len(a))
	copy(b, a)
	b[0] = 99
	return fmt.Sprintf("a=%v,b=%v", a, b)
}

// SliceGrowWithAppend tests slice growth
func SliceGrowWithAppend() string {
	s := make([]int, 0, 2)
	for i := 0; i < 5; i++ {
		s = append(s, i)
	}
	return fmt.Sprintf("len=%d,val=%v", len(s), s)
}

// SliceDeleteElement tests removing an element from slice
func SliceDeleteElement() string {
	s := []int{1, 2, 3, 4, 5}
	i := 2 // remove index 2
	s = append(s[:i], s[i+1:]...)
	return fmt.Sprintf("len=%d,val=%v", len(s), s)
}

// SliceInsertElement tests inserting into slice
func SliceInsertElement() string {
	s := []int{1, 2, 4, 5}
	s = append(s[:2], append([]int{3}, s[2:]...)...)
	return fmt.Sprintf("len=%d,val=%v", len(s), s)
}

// SliceFromArray tests slicing an array
func SliceFromArray() string {
	a := [5]int{10, 20, 30, 40, 50}
	s := a[1:4]
	return fmt.Sprintf("len=%d,val=%v", len(s), s)
}

// SliceOfString tests slice of string operations
func SliceOfString() string {
	s := []string{"a", "b", "c"}
	return fmt.Sprintf("len=%d,join=%s", len(s), s[0]+s[1]+s[2])
}
