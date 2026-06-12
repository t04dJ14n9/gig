package divergence_hunt81

// ============================================================================
// Round 81: Append/copy edge cases - nil append, grow, overlapping
// ============================================================================

func AppendToNil() []int {
	var s []int
	s = append(s, 1, 2, 3)
	return s
}

func AppendMultiple() []int {
	s := []int{1, 2}
	s = append(s, 3, 4, 5)
	return s
}

func AppendSlice() []int {
	s := []int{1, 2}
	extra := []int{3, 4, 5}
	s = append(s, extra...)
	return s
}

func AppendEmptySlice() []int {
	s := []int{1, 2}
	empty := []int{}
	s = append(s, empty...)
	return s
}

func AppendNilSlice() []int {
	s := []int{1, 2}
	var nilSlice []int
	s = append(s, nilSlice...)
	return s
}

func CopyBasic() int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n
}

func CopyLargerDst() int {
	src := []int{1, 2}
	dst := make([]int, 5)
	n := copy(dst, src)
	return n + dst[0] + dst[1]
}

func CopySlice() []int {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func CopyPartial() []int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	copy(dst, src[1:])
	return dst
}

func AppendBool() []bool {
	var s []bool
	s = append(s, true, false, true)
	return s
}

func AppendString() []string {
	var s []string
	s = append(s, "hello", "world")
	return s
}

func AppendFloat() []float64 {
	var s []float64
	s = append(s, 1.1, 2.2, 3.3)
	return s
}

func CopyStringSlice() []string {
	src := []string{"a", "b", "c"}
	dst := make([]string, 3)
	copy(dst, src)
	return dst
}

func AppendGrow() int {
	s := make([]int, 0)
	for i := 0; i < 100; i++ {
		s = append(s, i)
	}
	return len(s)
}
