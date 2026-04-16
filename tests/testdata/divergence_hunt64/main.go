package divergence_hunt64

// ============================================================================
// Round 64: Slice tricks - 3-index slice, append, copy, overlapping
// ============================================================================

func ThreeIndexSlice() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3] // len=2, cap=2
	return len(sub)*10 + cap(sub)
}

func ThreeIndexSliceFull() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:5] // len=2, cap=4
	return len(sub)*10 + cap(sub)
}

func SliceAppendWithinCap() []int {
	s := make([]int, 2, 5)
	s[0] = 1
	s[1] = 2
	s = append(s, 3)
	return s
}

func SliceAppendNil() []int {
	var s []int
	s = append(s, 1, 2, 3)
	return s
}

func SliceCopyCount() int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n
}

func SliceCopyFromSub() int {
	s := []int{1, 2, 3, 4, 5}
	dst := make([]int, 2)
	copy(dst, s[2:])
	return dst[0]*10 + dst[1]
}

func SliceAppendGrow() int {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	oldCap := cap(s)
	s = append(s, 3) // should grow
	newCap := cap(s)
	if newCap > oldCap {
		return 1
	}
	return 0
}

func SliceNilSubslice() []int {
	var s []int
	return s[0:0]
}

func SliceEmptySubslice() []int {
	s := []int{1, 2, 3}
	return s[2:2]
}

func SliceCapAfterAppend() int {
	s := []int{1, 2, 3}
	s = s[:2]
	return cap(s)
}

func SliceMakeZeroLen() int {
	s := make([]int, 0, 10)
	return len(s)*10 + cap(s)
}

func SliceOfString() int {
	s := []string{"hello", "world"}
	return len(s)
}

func SliceOfBool() int {
	s := []bool{true, false, true}
	count := 0
	for _, v := range s {
		if v {
			count++
		}
	}
	return count
}

func SliceOverlappingCopy() []int {
	s := []int{1, 2, 3, 4, 5}
	copy(s[1:], s[0:]) // overlapping copy
	return s
}

func SliceDoubleAppend() []int {
	s := []int{1, 2, 3}
	s = append(s, 4)
	s = append(s, 5)
	return s
}
