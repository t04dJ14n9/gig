package divergence_hunt272

import (
	"fmt"
)

// ============================================================================
// Round 272: Map edge cases — nil map, delete during range, zero value

// NilMapRead tests reading from nil map returns zero value
func NilMapRead() string {
	var m map[string]int
	v := m["key"]
	return fmt.Sprintf("v=%d,ok=%t", v, m != nil)
}

// NilMapWritePanics tests that writing to nil map panics
func NilMapWritePanics() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("panic:%v", r)
		}
	}()
	var m map[string]int
	m["key"] = 1
	return "no_panic"
}

// MapZeroValue tests zero value of map value type
func MapZeroValue() string {
	m := map[string]int{}
	return fmt.Sprintf("missing=%d", m["nonexistent"])
}

// MapLenAndDelete tests len after delete
func MapLenAndDelete() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return fmt.Sprintf("len=%d,b_exists=%t,a=%d", len(m), m["b"] != 0 || false, m["a"])
}

// MapDeleteNonExistent tests deleting non-existent key is no-op
func MapDeleteNonExistent() string {
	m := map[string]int{"a": 1}
	delete(m, "nonexistent")
	return fmt.Sprintf("len=%d,a=%d", len(m), m["a"])
}

// MapCommaOk tests map access with comma-ok
func MapCommaOk() string {
	m := map[string]int{"x": 42}
	v, ok := m["x"]
	v2, ok2 := m["y"]
	return fmt.Sprintf("v=%d,ok=%t,v2=%d,ok2=%t", v, ok, v2, ok2)
}

// MapAsSet tests using map as set
func MapAsSet() string {
	set := map[int]bool{}
	for i := 0; i < 3; i++ {
		set[i] = true
	}
	return fmt.Sprintf("len=%d,has1=%t,has5=%t", len(set), set[1], set[5])
}

// MapStructValue tests map with struct values
func MapStructValue() string {
	type Point struct{ X, Y int }
	m := map[string]Point{
		"origin": {0, 0},
		"unit":   {1, 1},
	}
	p := m["unit"]
	return fmt.Sprintf("x=%d,y=%d", p.X, p.Y)
}

// MapUpdateValue tests updating map values
func MapUpdateValue() string {
	m := map[string]int{"a": 1}
	m["a"] = m["a"] + 10
	return fmt.Sprintf("a=%d", m["a"])
}

// NestedMap tests map of maps
func NestedMap() string {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 42}
	return fmt.Sprintf("val=%d", m["outer"]["inner"])
}

// MapKeyWithStruct tests struct as map key
func MapKeyWithStruct() string {
	type Key struct{ A, B int }
	m := map[Key]string{{1, 2}: "found"}
	return fmt.Sprintf("val=%s", m[Key{1, 2}])
}

// MapWithSliceValue tests map with slice values
func MapWithSliceValue() string {
	m := map[string][]int{}
	m["nums"] = append(m["nums"], 1, 2, 3)
	return fmt.Sprintf("len=%d,val=%v", len(m["nums"]), m["nums"])
}
