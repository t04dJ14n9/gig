package divergence_hunt172

import (
	"fmt"
	"sort"
)

// ============================================================================
// Round 172: Sort interface implementations
// ============================================================================

// ByLength implements sort.Interface for sorting strings by length
type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

func SortByLength() string {
	words := []string{"apple", "pie", "banana", "kiwi"}
	sort.Sort(ByLength(words))
	return fmt.Sprintf("%v", words)
}

// ByAge implements sort.Interface for sorting by age
type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }

func SortStructByField() string {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Sort(ByAge(people))
	result := ""
	for _, p := range people {
		result += fmt.Sprintf("%s:%d ", p.Name, p.Age)
	}
	return result
}

// Reverse sort
type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func SortReverse() string {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Sort(Reverse{sort.IntSlice(nums)})
	return fmt.Sprintf("%v", nums)
}

func SortInts() string {
	nums := []int{5, 2, 6, 3, 1, 4}
	sort.Ints(nums)
	return fmt.Sprintf("%v", nums)
}

func SortStrings() string {
	words := []string{"cherry", "apple", "banana"}
	sort.Strings(words)
	return fmt.Sprintf("%v", words)
}

func SortFloats() string {
	nums := []float64{3.14, 1.41, 2.71, 1.73}
	sort.Float64s(nums)
	return fmt.Sprintf("%.2f %.2f %.2f %.2f", nums[0], nums[1], nums[2], nums[3])
}

func SortSearchInts() int {
	nums := []int{1, 3, 5, 7, 9}
	return sort.SearchInts(nums, 5)
}

func SortIsSorted() string {
	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{3, 1, 4, 1, 5}
	return fmt.Sprintf("%v:%v", sort.IntsAreSorted(sorted), sort.IntsAreSorted(unsorted))
}

func SortSlice() string {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	result := ""
	for _, p := range people {
		result += fmt.Sprintf("%s ", p.Name)
	}
	return result
}

func SortStable() string {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 30},
	}
	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	result := ""
	for _, p := range people {
		result += fmt.Sprintf("%s:%d ", p.Name, p.Age)
	}
	return result
}
