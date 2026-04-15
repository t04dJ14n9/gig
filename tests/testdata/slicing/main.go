package slicing

// SubSliceBasic tests basic sub-slice
func SubSliceBasic() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = (i + 1) * 10
	}
	sub := s[1:4]
	return sub[0] + sub[1] + sub[2]
}

// SubSliceLen tests sub-slice length
func SubSliceLen() int {
	s := make([]int, 10)
	return len(s[2:7])
}

// SubSliceFromStart tests sub-slice from start
func SubSliceFromStart() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[:3]
	sum := 0
	for _, v := range sub {
		sum = sum + v
	}
	return sum
}

// SubSliceToEnd tests sub-slice to end
func SubSliceToEnd() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[3:]
	sum := 0
	for _, v := range sub {
		sum = sum + v
	}
	return sum
}

// SubSliceCopy tests full slice copy
func SubSliceCopy() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[:]
	return len(sub)
}

// SubSliceChained tests chained sub-slice
func SubSliceChained() int {
	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		s[i] = i
	}
	sub := s[2:8]
	sub2 := sub[1:4]
	return sub2[0] + sub2[1] + sub2[2]
}

// SubSliceModifiesOriginal tests sub-slice modifies original
func SubSliceModifiesOriginal() int {
	s := make([]int, 5)
	for i := 0; i < 5; i++ {
		s[i] = i
	}
	sub := s[1:4]
	sub[0] = 99
	return s[1]
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// SliceLen returns len(s[from:to])
func SliceLen(s []int, from, to int) int { return len(s[from:to]) }

// SliceSumRange returns the sum of elements in s[from:to]
func SliceSumRange(s []int, from, to int) int {
	sum := 0
	for _, v := range s[from:to] {
		sum += v
	}
	return sum
}

// SliceFirstElement returns s[idx]
func SliceFirstElement(s []int, idx int) int { return s[idx] }

// ThreeIndexSlice tests 3-index slice expression s[low:high:max]
func ThreeIndexSlice() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3]
	return len(sub)*10 + cap(sub)
}

// ThreeIndexSliceFull tests 3-index slice with full bounds
func ThreeIndexSliceFull() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:4:5] // len=3, cap=4
	return len(sub)*10 + cap(sub)
}

// ThreeIndexSliceCapControl tests that 3-index slice limits capacity
func ThreeIndexSliceCapControl() int {
	s := make([]int, 5, 10)
	sub := s[1:3:4]
	return cap(sub)
}
