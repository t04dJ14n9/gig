package divergence_hunt96

import "fmt"

// ============================================================================
// Round 96: Slice tricks - delete, insert, filter in-place
// ============================================================================

func SliceDelete() string {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return fmt.Sprintf("%v", s)
}

func SliceInsert() string {
	s := []int{1, 2, 5, 6}
	i := 2
	s = append(s[:i], append([]int{3, 4}, s[i:]...)...)
	return fmt.Sprintf("%v", s)
}

func SliceFilter() string {
	s := []int{1, 2, 3, 4, 5, 6}
	filtered := s[:0]
	for _, v := range s {
		if v%2 == 0 {
			filtered = append(filtered, v)
		}
	}
	return fmt.Sprintf("%v", filtered)
}

func SliceReverse() string {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return fmt.Sprintf("%v", s)
}

func SliceUnique() string {
	s := []string{"a", "b", "a", "c", "b"}
	seen := map[string]bool{}
	result := []string{}
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return fmt.Sprintf("%v", result)
}

func SliceFlatten() string {
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	flat := []int{}
	for _, inner := range nested {
		flat = append(flat, inner...)
	}
	return fmt.Sprintf("%v", flat)
}

func SliceBatch() string {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	batchSize := 3
	var batches [][]int
	for batchSize < len(s) {
		s, batches = s[batchSize:], append(batches, s[0:batchSize:batchSize])
	}
	batches = append(batches, s)
	return fmt.Sprintf("%d batches", len(batches))
}

func SliceClone() string {
	original := []int{1, 2, 3}
	clone := make([]int, len(original))
	copy(clone, original)
	clone[0] = 99
	return fmt.Sprintf("%v:%v", original, clone)
}

func SliceAppendGrow() string {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	capBefore := cap(s)
	s = append(s, 3)
	return fmt.Sprintf("cap grew: %v", cap(s) > capBefore)
}

func SliceCut() string {
	s := []int{1, 2, 3, 4, 5}
	// Remove elements at index 1..3 (exclusive)
	i, j := 1, 3
	s = append(s[:i], s[j:]...)
	return fmt.Sprintf("%v", s)
}
