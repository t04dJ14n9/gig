package divergence_hunt132

import "fmt"

// ============================================================================
// Round 132: Slice append and capacity growth patterns
// ============================================================================

func SliceGrowFromEmpty() string {
	s := []int{}
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return fmt.Sprintf("%v", s)
}

func SliceGrowWithCap() string {
	s := make([]int, 0, 5)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return fmt.Sprintf("len=%d", len(s))
}

func SliceReslice() string {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	sub := s[2:5]
	return fmt.Sprintf("%v", sub)
}

func SliceResliceCap() string {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	sub := s[2:5:7] // len=3, cap=5
	return fmt.Sprintf("%v-len=%d-cap=%d", sub, len(sub), cap(sub))
}

func SliceAppendBeyondCap() string {
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 1, 2, 3
	s = append(s, 4)
	s = append(s, 5)
	s = append(s, 6) // triggers growth
	return fmt.Sprintf("len=%d-cap=%d", len(s), cap(s))
}

func SliceMakeZeroLen() string {
	s := make([]int, 0)
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("len=%d-nil=%t", len(s), s == nil)
}

func SliceNilVsEmpty() string {
	var s1 []int
	s2 := []int{}
	s3 := make([]int, 0)
	return fmt.Sprintf("nil=%t-empty=%t-make=%t", s1 == nil, s2 == nil, s3 == nil)
}

func SliceOfString() string {
	s := []string{"hello", "world"}
	s = append(s, "!")
	return fmt.Sprintf("%v", s)
}

func SliceBool() string {
	s := []bool{true, false, true}
	count := 0
	for _, v := range s {
		if v {
			count++
		}
	}
	return fmt.Sprintf("count=%d", count)
}

func SliceStructLiteral() string {
	type Point struct{ X, Y int }
	s := []Point{{1, 2}, {3, 4}, {5, 6}}
	return fmt.Sprintf("%v", s)
}
