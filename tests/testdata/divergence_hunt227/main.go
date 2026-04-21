package divergence_hunt227

import "fmt"

// ============================================================================
// Round 227: Slice capacity and length
// ============================================================================

// SliceLenCapBasic tests basic len/cap operations
func SliceLenCapBasic() string {
	s := make([]int, 5, 10)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceLenCapAfterAppend tests len/cap after append
func SliceLenCapAfterAppend() string {
	s := make([]int, 3, 5)
	s = append(s, 1, 2)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceLenCapReslice tests len/cap after reslicing
func SliceLenCapReslice() string {
	s := make([]int, 10, 20)
	s = s[2:8]
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceLenCapNil tests len/cap of nil slice
func SliceLenCapNil() string {
	var s []int
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceLenCapEmptyLiteral tests len/cap of empty literal
func SliceLenCapEmptyLiteral() string {
	s := []int{}
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceLenCapSlicingBounds tests len/cap with different slice bounds
func SliceLenCapSlicingBounds() string {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[1:3]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// SliceLenCapFullSlice tests full slice expression
func SliceLenCapFullSlice() string {
	s := make([]int, 5, 10)
	s2 := s[1:3:5]
	return fmt.Sprintf("len=%d,cap=%d", len(s2), cap(s2))
}

// SliceLenCapZeroLength tests zero length but non-zero capacity
func SliceLenCapZeroLength() string {
	s := make([]int, 0, 100)
	return fmt.Sprintf("len=%d,cap=%d,can_append=%t", len(s), cap(s), cap(s) > len(s))
}

// SliceLenCapByteSlice tests len/cap with byte slices
func SliceLenCapByteSlice() string {
	b := make([]byte, 10, 50)
	return fmt.Sprintf("len=%d,cap=%d", len(b), cap(b))
}

// SliceLenCapStringConversion tests len after string conversion
func SliceLenCapStringConversion() string {
	str := "hello world"
	b := []byte(str)
	return fmt.Sprintf("strlen=%d,bytelen=%d,cap=%d", len(str), len(b), cap(b))
}

// SliceLenCapAfterClear tests len/cap after clearing
func SliceLenCapAfterClear() string {
	s := []int{1, 2, 3, 4, 5}
	s = s[:0]
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}
