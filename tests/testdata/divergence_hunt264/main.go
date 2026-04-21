package divergence_hunt264

import (
	"fmt"
)

// ============================================================================
// Round 264: Map iteration order, deletion during range, nil map ops
// ============================================================================

// MapDeleteDuringRange tests deleting from map during range
func MapDeleteDuringRange() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		if k == "b" {
			delete(m, k)
		}
	}
	_, hasB := m["b"]
	_, hasA := m["a"]
	return fmt.Sprintf("hasA=%t,hasB=%t,len=%d", hasA, hasB, len(m))
}

// MapNilRead tests reading from nil map
func MapNilRead() string {
	var m map[string]int
	v, ok := m["key"]
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// MapNilWritePanic tests writing to nil map (should panic, recovered)
func MapNilWritePanic() string {
	var m map[string]int
	defer func() {
		recover()
	}()
	m["key"] = 1
	return "should_not_reach"
}

// MapLenCap tests map length
func MapLenCap() string {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	return fmt.Sprintf("len=%d", len(m))
}

// MapOverwriteKey tests overwriting a key
func MapOverwriteKey() string {
	m := map[string]int{"x": 1}
	m["x"] = 2
	return fmt.Sprintf("x=%d,len=%d", m["x"], len(m))
}

// MapTwoValLookup tests comma-ok map lookup
func MapTwoValLookup() string {
	m := map[string]int{"a": 1}
	v1, ok1 := m["a"]
	v2, ok2 := m["b"]
	return fmt.Sprintf("v1=%d,ok1=%t,v2=%d,ok2=%t", v1, ok1, v2, ok2)
}

// MapWithStructKey tests struct as map key
func MapWithStructKey() string {
	type Point struct{ X, Y int }
	m := map[Point]string{
		{1, 2}: "origin+",
		{0, 0}: "origin",
	}
	return fmt.Sprintf("p00=%s,p12=%s", m[Point{0, 0}], m[Point{1, 2}])
}

// MapOfSlices tests map with slice values (not keys)
func MapOfSlices() string {
	m := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5},
	}
	return fmt.Sprintf("a_len=%d,b_len=%d", len(m["a"]), len(m["b"]))
}

// MapGrowDynamic tests map growing dynamically
func MapGrowDynamic() string {
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i * 2
	}
	return fmt.Sprintf("len=%d,v50=%d", len(m), m[50])
}

// MapDeleteAndReinsert tests delete then re-insert
func MapDeleteAndReinsert() string {
	m := map[string]int{"x": 10}
	delete(m, "x")
	_, ok1 := m["x"]
	m["x"] = 20
	v, ok2 := m["x"]
	return fmt.Sprintf("ok1=%t,ok2=%t,v=%d", ok1, ok2, v)
}
