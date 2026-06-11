package divergence_hunt279

import (
	"fmt"
)

// ============================================================================
// Round 279: Deeply nested composite types

// SliceOfMaps tests []map[string]int
func SliceOfMaps() string {
	s := []map[string]int{
		{"a": 1},
		{"b": 2},
	}
	return fmt.Sprintf("s[0][a]=%d,s[1][b]=%d", s[0]["a"], s[1]["b"])
}

// MapOfSlices tests map[string][]int
func MapOfSlices() string {
	m := map[string][]int{
		"evens": {2, 4, 6},
		"odds":  {1, 3, 5},
	}
	return fmt.Sprintf("evens=%v,odds=%v", m["evens"], m["odds"])
}

// MapOfMaps tests map[string]map[string]int
func MapOfMaps() string {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 99}
	return fmt.Sprintf("val=%d", m["outer"]["inner"])
}

// NestedStructInSlice tests struct in slice
func NestedStructInSlice() string {
	type Item struct {
		Name  string
		Price float64
	}
	items := []Item{
		{"apple", 1.5},
		{"banana", 0.75},
	}
	return fmt.Sprintf("%s=%.2f", items[0].Name, items[0].Price)
}

// SliceOfSlices tests [][]int
func SliceOfSlices() string {
	s := [][]int{
		{1, 2, 3},
		{4, 5},
		{6},
	}
	return fmt.Sprintf("s[0]=%v,s[1]=%v,s[2]=%v", s[0], s[1], s[2])
}

// MapWithStructKey tests map with struct key containing nested data
func MapWithStructKey() string {
	type Coord struct{ X, Y int }
	type Cell struct{ Value string }
	grid := map[Coord]Cell{
		{0, 0}: {"origin"},
		{1, 1}: {"center"},
	}
	return fmt.Sprintf("origin=%s,center=%s", grid[Coord{0, 0}].Value, grid[Coord{1, 1}].Value)
}

// SliceOfPointers tests []*int
func SliceOfPointers() string {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}
	return fmt.Sprintf("vals=%d,%d,%d", *s[0], *s[1], *s[2])
}

// StructWithSliceOfMaps tests struct containing []map
func StructWithSliceOfMaps() string {
	type Config struct {
		Name    string
		Entries []map[string]string
	}
	c := Config{
		Name: "test",
		Entries: []map[string]string{
			{"key": "val1"},
			{"key": "val2"},
		},
	}
	return fmt.Sprintf("name=%s,e0=%s,e1=%s", c.Name, c.Entries[0]["key"], c.Entries[1]["key"])
}

// ThreeLevelNesting tests map[string]map[int][]string
func ThreeLevelNesting() string {
	m := map[string]map[int][]string{}
	m["level1"] = map[int][]string{}
	m["level1"][0] = []string{"a", "b"}
	m["level1"][1] = []string{"c"}
	return fmt.Sprintf("l1_0=%v,l1_1=%v", m["level1"][0], m["level1"][1])
}

// AppendToNestedSlice tests appending to nested slice
func AppendToNestedSlice() string {
	m := map[string][]int{}
	m["nums"] = []int{1, 2}
	m["nums"] = append(m["nums"], 3)
	return fmt.Sprintf("nums=%v", m["nums"])
}

// ArrayOfStructs tests [3]struct{...}
func ArrayOfStructs() string {
	type Point struct{ X, Y int }
	arr := [3]Point{{1, 2}, {3, 4}, {5, 6}}
	return fmt.Sprintf("p0=%v,p2=%v", arr[0], arr[2])
}

// ModifyNestedValue tests modifying value deep in nested structure
func ModifyNestedValue() string {
	s := [][]int{{1, 2, 3}, {4, 5, 6}}
	s[0][1] = 99
	return fmt.Sprintf("s=%v", s)
}
