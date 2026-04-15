package divergence_hunt12

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 12: Encoding, string manipulation, data structure operations
// ============================================================================

// JSONNestedStruct tests JSON with nested struct
func JSONNestedStruct() string {
	type Address struct{ City string }
	type Person struct {
		Name    string
		Address Address
	}
	p := Person{Name: "Alice", Address: Address{City: "NYC"}}
	b, _ := json.Marshal(p)
	return string(b)
}

// JSONSliceOfStructs tests JSON with slice of structs
func JSONSliceOfStructs() string {
	type Item struct{ Name string; Value int }
	items := []Item{{"a", 1}, {"b", 2}}
	b, _ := json.Marshal(items)
	return string(b)
}

// JSONUnmarshalIntoMap tests JSON unmarshal into map
func JSONUnmarshalIntoMap() int {
	data := `{"x":10,"y":20}`
	m := map[string]int{}
	json.Unmarshal([]byte(data), &m)
	return m["x"] + m["y"]
}

// StringTitle tests strings.Title (deprecated but still works)
func StringTitle() string { return strings.ToTitle("hello world") }

// StringEqualFold tests strings.EqualFold
func StringEqualFold() bool { return strings.EqualFold("Hello", "HELLO") }

// StringIndex tests strings.Index
func StringIndex() int { return strings.Index("hello world", "world") }

// StringLastIndex tests strings.LastIndex
func StringLastIndex() int { return strings.LastIndex("hello hello", "hello") }

// StringIndexAny tests strings.IndexAny
func StringIndexAny() int { return strings.IndexAny("hello", "aeiou") }

// StringNewReplacer tests strings.NewReplacer
func StringNewReplacer() string {
	r := strings.NewReplacer("a", "b", "b", "c")
	return r.Replace("abc")
}

// StringBuilderGrow tests strings.Builder.Grow
func StringBuilderGrow() int {
	var b strings.Builder
	b.Grow(100)
	b.WriteString("hello")
	return b.Len()
}

// SortSliceStable tests sort.SliceStable
func SortSliceStable() int {
	type Item struct{ Name string; Priority int }
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

// SortSearch tests sort.Search
func SortSearch() int {
	data := []int{1, 3, 5, 7, 9}
	i := sort.Search(len(data), func(i int) bool { return data[i] >= 5 })
	return data[i]
}

// FmtSprintfBoolean tests fmt with boolean
func FmtSprintfBoolean() string {
	return fmt.Sprintf("%v %t", true, false)
}

// FmtSprintfFloat tests fmt with various float formats
func FmtSprintfFloat() string {
	return fmt.Sprintf("%.1f %e %g", 3.14, 3.14, 3.14)
}

// FmtSprintfInt tests fmt with various int formats
func FmtSprintfInt() string {
	return fmt.Sprintf("%d %x %o %b", 42, 42, 42, 42)
}

// FmtSprintfString tests fmt with string formats
func FmtSprintfString() string {
	return fmt.Sprintf("%s %q %5s", "hi", "hi", "hi")
}

// JSONMarshalBool tests JSON with bool
func JSONMarshalBool() string {
	b, _ := json.Marshal(true)
	return string(b)
}

// JSONUnmarshalBool tests JSON unmarshal bool
func JSONUnmarshalBool() bool {
	var v bool
	json.Unmarshal([]byte("true"), &v)
	return v
}

// JSONMarshalNil tests JSON with nil
func JSONMarshalNil() string {
	var s []int
	b, _ := json.Marshal(s)
	return string(b)
}

// SliceMinMaxInt tests finding min/max in int slice
func SliceMinMaxInt() int {
	s := []int{5, 3, 8, 1, 9, 2, 7}
	min, max := s[0], s[0]
	for _, v := range s[1:] {
		if v < min { min = v }
		if v > max { max = v }
	}
	return min * 10 + max
}

// StringCountSubstring tests counting substrings
func StringCountSubstring() int {
	return strings.Count("abababab", "ab")
}

// MapHasKey tests map key existence
func MapHasKey() bool {
	m := map[string]int{"a": 1}
	_, ok := m["a"]
	_, ok2 := m["b"]
	return ok && !ok2
}
