package divergence_hunt186

import (
	"fmt"
)

// ============================================================================
// Round 186: Slice tricks (copy, clear, clip)
// ============================================================================

func SliceCopyBasic() string {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src)
	return fmt.Sprintf("%d:%v", n, dst)
}

func SliceCopyPartial() string {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return fmt.Sprintf("%d:%v", n, dst)
}

func SliceCopyOverlap() string {
	nums := []int{1, 2, 3, 4, 5}
	copy(nums[1:], nums)
	return fmt.Sprintf("%v", nums)
}

func SliceCopyString() string {
	s := "hello"
	b := make([]byte, len(s))
	copy(b, s)
	return fmt.Sprintf("%s", string(b))
}

func SliceClear() string {
	nums := []int{1, 2, 3, 4, 5}
	for i := range nums {
		nums[i] = 0
	}
	return fmt.Sprintf("%v", nums)
}

func SliceAppendGrow() string {
	s := make([]int, 0, 3)
	s = append(s, 1)
	s = append(s, 2, 3)
	s = append(s, 4)
	return fmt.Sprintf("%v:%d:%d", s, len(s), cap(s))
}

func SliceClip() string {
	s := make([]int, 5, 10)
	s = s[:3]
	return fmt.Sprintf("%d:%d", len(s), cap(s))
}

func SliceDeleteElement() string {
	nums := []int{1, 2, 3, 4, 5}
	i := 2
	nums = append(nums[:i], nums[i+1:]...)
	return fmt.Sprintf("%v", nums)
}

func SliceInsertElement() string {
	nums := []int{1, 2, 4, 5}
	nums = append(nums[:2], append([]int{3}, nums[2:]...)...)
	return fmt.Sprintf("%v", nums)
}

func SliceReverse() string {
	nums := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return fmt.Sprintf("%v", nums)
}

func SliceFilter() string {
	nums := []int{1, 2, 3, 4, 5, 6}
	filtered := nums[:0]
	for _, n := range nums {
		if n%2 == 0 {
			filtered = append(filtered, n)
		}
	}
	return fmt.Sprintf("%v", filtered)
}
