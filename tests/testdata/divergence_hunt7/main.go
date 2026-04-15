package divergence_hunt7

import "sort"

// ============================================================================
// Round 7: Sorting, slice manipulation, map operations, struct methods,
// pointer receivers, method values, type assertions
// ============================================================================

// SortInts tests sort.Ints
func SortInts() int {
	s := []int{5, 3, 1, 4, 2}
	sort.Ints(s)
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

// SortStrings tests sort.Strings
func SortStrings() string {
	s := []string{"banana", "apple", "cherry"}
	sort.Strings(s)
	return s[0] + s[1] + s[2]
}

// SortFloat64s tests sort.Float64s
func SortFloat64s() float64 {
	s := []float64{3.14, 1.41, 2.71}
	sort.Float64s(s)
	return s[0] + s[1] + s[2]
}

// SliceDelete tests removing element from slice
func SliceDelete() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return len(s)*10 + s[2]
}

// SliceInsert tests inserting element into slice
func SliceInsert() int {
	s := []int{1, 2, 4, 5}
	s = append(s[:2], append([]int{3}, s[2:]...)...)
	return s[2]
}

// SliceContains tests if slice contains element
func SliceContains() bool {
	s := []int{1, 2, 3, 4, 5}
	for _, v := range s {
		if v == 3 { return true }
	}
	return false
}

// MapKeys tests getting map keys
func MapKeys() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return len(keys)
}

// MapValues tests getting map values
func MapValues() int {
	m := map[string]int{"a": 1, "b": 2}
	sum := 0
	for _, v := range m { sum += v }
	return sum
}

// StructWithMethods tests struct with methods
func StructWithMethods() int {
	type Rect struct{ W, H int }
	area := func(r Rect) int { return r.W * r.H }
	r := Rect{W: 3, H: 4}
	return area(r)
}

// PointerReceiverMethod tests pointer receiver pattern
func PointerReceiverMethod() int {
	type Counter struct{ n int }
	inc := func(c *Counter) { c.n++ }
	c := Counter{n: 0}
	inc(&c)
	inc(&c)
	inc(&c)
	return c.n
}

// TypeAssertion tests type assertion
func TypeAssertion() int {
	var x any = 42
	if v, ok := x.(int); ok {
		return v
	}
	return -1
}

// TypeAssertionString tests type assertion to string
func TypeAssertionString() int {
	var x any = "hello"
	if v, ok := x.(string); ok {
		return len(v)
	}
	return -1
}

// TypeAssertionFail tests failed type assertion
func TypeAssertionFail() int {
	var x any = 42
	if _, ok := x.(string); ok {
		return 1
	}
	return 0
}

// InterfaceTypeSwitch tests interface type switch
func InterfaceTypeSwitch() int {
	describe := func(x any) int {
		switch v := x.(type) {
		case int: return v
		case string: return len(v)
		case bool:
			if v { return 1 }
			return 0
		default: return -1
		}
	}
	return describe(42) + describe("hi") + describe(true)
}

// SliceDedupe tests deduplicating a slice
func SliceDedupe() int {
	s := []int{1, 2, 2, 3, 3, 3, 4}
	seen := map[int]bool{}
	result := []int{}
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return len(result)
}

// MapMerge tests merging two maps
func MapMerge() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	for k, v := range m2 {
		m1[k] = v
	}
	return len(m1)
}

// StructSliceSort tests sorting slice of structs
func StructSliceSort() int {
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

// MapInvert tests inverting a map
func MapInvert() int {
	m := map[string]int{"a": 1, "b": 2}
	inverted := map[int]string{}
	for k, v := range m {
		inverted[v] = k
	}
	return len(inverted)
}

// NestedInterface tests nested interface
func NestedInterface() int {
	var x any = any(42)
	return x.(int)
}

// SliceFlatten tests flattening nested slices conceptually
func SliceFlatten() int {
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	count := 0
	for _, inner := range nested {
		count += len(inner)
	}
	return count
}

// IntSliceSortCustom tests custom sort on int slice
func IntSliceSortCustom() int {
	s := []int{5, 3, 1, 4, 2}
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

// MapCountValues tests counting values in map
func MapCountValues() int {
	m := map[string]int{"a": 1, "b": 2, "c": 1}
	count := map[int]int{}
	for _, v := range m {
		count[v]++
	}
	return count[1]
}
