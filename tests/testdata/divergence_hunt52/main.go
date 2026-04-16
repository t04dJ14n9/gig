package divergence_hunt52

import "sort"

// ============================================================================
// Round 52: Sort patterns - ints, strings, floats, structs, custom comparators
// ============================================================================

func SortInts() int {
	s := []int{5, 3, 1, 4, 2}
	sort.Ints(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SortStrings() string {
	s := []string{"cherry", "apple", "banana"}
	sort.Strings(s)
	return s[0] + s[1] + s[2]
}

func SortFloat64s() float64 {
	s := []float64{3.14, 1.41, 2.71}
	sort.Float64s(s)
	return s[0] + s[1] + s[2]
}

func SortIntSlice() int {
	s := []int{5, 3, 1, 4, 2}
	sort.Ints(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SortReverse() int {
	s := []int{1, 2, 3, 4, 5}
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SortStructSlice() int {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Bob", 30},
		{"Alice", 25},
		{"Charlie", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	return people[0].Age*100 + people[1].Age*10 + people[2].Age
}

func SortSliceDesc() int {
	s := []int{1, 5, 3, 2, 4}
	sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SortStable() int {
	type Item struct {
		Name     string
		Priority int
	}
	items := []Item{
		{"a", 2},
		{"b", 1},
		{"c", 2},
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Priority < items[j].Priority
	})
	return items[0].Priority*100 + items[1].Priority*10 + items[2].Priority
}

func SortSearch() int {
	data := []int{1, 3, 5, 7, 9}
	i := sort.SearchInts(data, 5)
	return data[i]
}

func SortSearchString() int {
	data := []string{"apple", "banana", "cherry"}
	i := sort.SearchStrings(data, "banana")
	return i
}

func SortFloat64Search() int {
	data := []float64{1.1, 2.2, 3.3}
	i := sort.SearchFloat64s(data, 2.2)
	return i
}

func SortIsSorted() bool {
	s := []int{1, 2, 3, 4, 5}
	return sort.IntsAreSorted(s)
}

func SortEmptySlice() int {
	s := []int{}
	sort.Ints(s)
	return len(s)
}

func SortSingleElement() int {
	s := []int{42}
	sort.Ints(s)
	return s[0]
}

func SortDuplicate() int {
	s := []int{3, 1, 2, 3, 1}
	sort.Ints(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}
