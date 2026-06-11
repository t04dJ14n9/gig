package divergence_hunt142

import "fmt"

// ============================================================================
// Round 142: Slice three-index slicing and append patterns
// ============================================================================

func ThreeIndexBasic() string {
	s := []int{0, 1, 2, 3, 4, 5}
	sub := s[1:3:5]
	return fmt.Sprintf("val=%v-len=%d-cap=%d", sub, len(sub), cap(sub))
}

func ThreeIndexAppendNoGrow() string {
	s := []int{0, 1, 2, 3, 4, 5}
	sub := s[1:3:4] // len=2, cap=3
	sub = append(sub, 99)
	return fmt.Sprintf("sub=%v", sub)
}

func ThreeIndexFullSlice() string {
	s := []int{0, 1, 2, 3, 4}
	sub := s[1:3:3] // len=2, cap=2 — no room to append without growing
	sub = append(sub, 77)
	return fmt.Sprintf("sub=%v", sub)
}

func AppendCopyPattern() string {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	copy(dst, src)
	dst[0] = 99
	return fmt.Sprintf("src=%v-dst=%v", src, dst)
}

func SliceInsertMiddle() string {
	s := []int{1, 2, 5, 6}
	i := 2
	s = append(s[:i], append([]int{3, 4}, s[i:]...)...)
	return fmt.Sprintf("%v", s)
}

func SliceFilter() string {
	s := []int{1, 2, 3, 4, 5, 6}
	var result []int
	for _, v := range s {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return fmt.Sprintf("%v", result)
}

func SliceReverse() string {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return fmt.Sprintf("%v", s)
}

func SliceClone() string {
	s := []int{10, 20, 30}
	clone := append([]int{}, s...)
	clone[0] = 99
	return fmt.Sprintf("orig=%v-clone=%v", s, clone)
}

func SliceStackPattern() string {
	var stack []int
	stack = append(stack, 1)
	stack = append(stack, 2)
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return fmt.Sprintf("top=%d-remaining=%v", top, stack)
}
