package divergence_hunt152

import (
	"fmt"
	"sort"
)

// ============================================================================
// Round 152: Advanced slice operations and sorting
// ============================================================================

// SortIntsReverse tests sorting ints in reverse order
func SortIntsReverse() string {
	s := []int{5, 2, 8, 1, 9, 3}
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return fmt.Sprintf("%v", s)
}

// SortStrings tests sorting strings
func SortStrings() string {
	s := []string{"banana", "apple", "cherry", "date"}
	sort.Strings(s)
	return fmt.Sprintf("%v", s)
}

// SortFloat64s tests sorting float64s
func SortFloat64s() string {
	s := []float64{3.14, 1.41, 2.71, 0.57}
	sort.Float64s(s)
	return fmt.Sprintf("%.2f", s[0])
}

// SortSearch tests binary search in sorted slice
func SortSearch() string {
	s := []int{1, 3, 5, 7, 9, 11, 13}
	idx := sort.Search(len(s), func(i int) bool { return s[i] >= 7 })
	return fmt.Sprintf("idx=%d-val=%d", idx, s[idx])
}

// SortSearchInts tests sort.SearchInts convenience function
func SortSearchInts() string {
	s := []int{10, 20, 30, 40, 50}
	idx := sort.SearchInts(s, 35)
	return fmt.Sprintf("idx=%d", idx)
}

// SliceIsSorted tests checking if slice is sorted
func SliceIsSorted() string {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{5, 3, 1, 4, 2}
	return fmt.Sprintf("s1=%t-s2=%t", sort.IntsAreSorted(s1), sort.IntsAreSorted(s2))
}

// SliceSortStable tests stable sort
func SliceSortStable() string {
	type Item struct {
		Val   int
		Order int
	}
	s := []Item{
		{3, 0}, {1, 1}, {3, 2}, {2, 3}, {1, 4},
	}
	// Sort by Val, keeping original order for equal values
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Val < s[j].Val
	})
	return fmt.Sprintf("order=%d,%d", s[0].Order, s[1].Order)
}

// SliceCut tests slice cutting (removing middle)
func SliceCut() string {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	// Cut out elements 2-4 (indices 2, 3, 4)
	s = append(s[:2], s[5:]...)
	return fmt.Sprintf("%v", s)
}

// SliceInsert tests inserting into slice
func SliceInsert() string {
	s := []int{1, 2, 4, 5}
	// Insert 3 at index 2
	s = append(s[:2], append([]int{3}, s[2:]...)...)
	return fmt.Sprintf("%v", s)
}

// SliceDeleteUnordered tests unordered delete (swap with last)
func SliceDeleteUnordered() string {
	s := []int{1, 2, 3, 4, 5}
	// Delete element at index 1 (value 2) by swapping with last
	s[1] = s[len(s)-1]
	s = s[:len(s)-1]
	return fmt.Sprintf("len=%d", len(s))
}

// SliceCompact tests removing duplicates
func SliceCompact() string {
	s := []int{1, 1, 2, 2, 2, 3, 4, 4, 5}
	if len(s) == 0 {
		return "empty"
	}
	result := s[:1]
	for _, v := range s[1:] {
		if v != result[len(result)-1] {
			result = append(result, v)
		}
	}
	return fmt.Sprintf("%v", result)
}
