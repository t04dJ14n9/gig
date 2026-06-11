package divergence_hunt109

import (
	"fmt"
	"sort"
)

// ============================================================================
// Round 109: Slice sorting with custom comparators
// ============================================================================

func SortIntSlice() string {
	s := []int{5, 3, 1, 4, 2}
	sort.Ints(s)
	return fmt.Sprintf("%v", s)
}

func SortStringSlice() string {
	s := []string{"cherry", "apple", "banana"}
	sort.Strings(s)
	return fmt.Sprintf("%v", s)
}

func SortByLen() string {
	s := []string{"bb", "aaa", "c"}
	sort.Slice(s, func(i, j int) bool { return len(s[i]) < len(s[j]) })
	return fmt.Sprintf("%v", s)
}

func SortStructByField() string {
	type Item struct{ Name string; Val int }
	items := []Item{
		{"c", 3}, {"a", 1}, {"b", 2},
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Val < items[j].Val })
	return fmt.Sprintf("%s%d", items[0].Name, items[0].Val)
}

func SortReverse() string {
	s := []int{1, 2, 3, 4, 5}
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return fmt.Sprintf("%v", s)
}

func SortFloatSlice() string {
	s := []float64{3.14, 1.41, 2.71}
	sort.Float64s(s)
	return fmt.Sprintf("%.2f", s[0])
}

func SortStable() string {
	type Pair struct{ Key, Val int }
	items := []Pair{
		{2, 1}, {1, 2}, {2, 3}, {1, 4},
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Key < items[j].Key })
	return fmt.Sprintf("%d%d", items[0].Val, items[1].Val)
}

func SortIsSorted() string {
	s := []int{1, 2, 3}
	return fmt.Sprintf("%v", sort.IntsAreSorted(s))
}

func SortEmpty() string {
	s := []int{}
	sort.Ints(s)
	return fmt.Sprintf("%v", s)
}

func SortSingleElement() string {
	s := []int{42}
	sort.Ints(s)
	return fmt.Sprintf("%v", s)
}
