package divergence_hunt176

import (
	"fmt"
	"sort"
)

// ============================================================================
// Round 176: Binary search variations
// ============================================================================

func BinarySearchIntsFound() int {
	nums := []int{1, 3, 5, 7, 9, 11, 13}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 7
	})
	return idx
}

func BinarySearchIntsNotFound() int {
	nums := []int{1, 3, 5, 7, 9, 11, 13}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 8
	})
	return idx
}

func BinarySearchIntsBefore() int {
	nums := []int{1, 3, 5, 7, 9}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 0
	})
	return idx
}

func BinarySearchIntsAfter() int {
	nums := []int{1, 3, 5, 7, 9}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 10
	})
	return idx
}

func BinarySearchStrings() int {
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	idx := sort.Search(len(words), func(i int) bool {
		return words[i] >= "cherry"
	})
	return idx
}

func BinarySearchFloats() int {
	nums := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 3.3
	})
	return idx
}

func BinarySearchIntsFunc() int {
	nums := []int{1, 3, 5, 7, 9}
	idx := sort.SearchInts(nums, 5)
	return idx
}

func BinarySearchStringsFunc() int {
	words := []string{"apple", "banana", "cherry", "date"}
	idx := sort.SearchStrings(words, "cherry")
	return idx
}

func BinarySearchEmpty() int {
	nums := []int{}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 5
	})
	return idx
}

func BinarySearchSingleElement() int {
	nums := []int{5}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 5
	})
	return idx
}

func BinarySearchDuplicates() int {
	nums := []int{1, 2, 2, 2, 3, 4, 5}
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 2
	})
	return idx
}

func BinarySearchFindRange() string {
	nums := []int{1, 2, 2, 2, 3, 4, 5}
	// Find first position >= 2
	first := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= 2
	})
	// Find first position > 2
	last := sort.Search(len(nums), func(i int) bool {
		return nums[i] > 2
	})
	return fmt.Sprintf("first=%d last=%d", first, last)
}
