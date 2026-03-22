package thirdparty

import "sort"

// SortStrings tests sort.Strings.
func SortStrings() int {
	s := []string{"banana", "apple", "cherry"}
	sort.Strings(s)
	if s[0] == "apple" && s[1] == "banana" && s[2] == "cherry" {
		return 1
	}
	return 0
}

// SortInts tests sort.Ints.
func SortInts() int {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Ints(s)
	if s[0] == 1 && s[7] == 9 {
		return 1
	}
	return 0
}

// SortFloat64s tests sort.Float64s.
func SortFloat64s() int {
	s := []float64{3.14, 1.41, 2.71}
	sort.Float64s(s)
	return int(s[0] * 100)
}

// SortSearchInts tests sort.SearchInts.
func SortSearchInts() int {
	s := []int{1, 3, 5, 7, 9}
	return sort.SearchInts(s, 5)
}

// SortSearchStrings tests sort.SearchStrings.
func SortSearchStrings() int {
	s := []string{"apple", "banana", "cherry"}
	return sort.SearchStrings(s, "banana")
}

// SortSlice tests sort.Slice with custom comparator.
func SortSlice() int {
	s := []int{3, 1, 4, 1, 5}
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	if s[0] == 1 && s[4] == 5 {
		return 1
	}
	return 0
}

// SortSliceStable tests sort.SliceStable.
func SortSliceStable() int {
	s := []int{3, 1, 4, 1, 5}
	sort.SliceStable(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	if s[0] == 1 && s[4] == 5 {
		return 1
	}
	return 0
}

// SortIsSorted tests sort.IntsAreSorted.
func SortIsSorted() int {
	s := []int{1, 2, 3}
	if sort.IntsAreSorted(s) {
		return 1
	}
	return 0
}
