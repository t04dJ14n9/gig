package divergence_hunt63

import "math"

// ============================================================================
// Round 63: Map edge cases - NaN keys, struct keys, delete during range
// ============================================================================

func MapNaNKey() int {
	m := map[float64]int{}
	m[math.NaN()] = 1
	m[math.NaN()] = 2
	return len(m) // NaN != NaN, so both keys exist
}

func MapNaNKeyLookup() int {
	m := map[float64]int{math.NaN(): 42}
	v, ok := m[math.NaN()]
	if ok {
		return v
	}
	return -1 // NaN lookup should fail since NaN != NaN
}

func MapStructKey() int {
	type Point struct{ X, Y int }
	m := map[Point]int{}
	m[Point{1, 2}] = 10
	m[Point{3, 4}] = 20
	return m[Point{1, 2}]
}

func MapArrayKey() int {
	m := map[[3]int]int{}
	m[[3]int{1, 2, 3}] = 100
	return m[[3]int{1, 2, 3}]
}

func MapDeleteDuringRange() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	for k := range m {
		if k%2 == 0 {
			delete(m, k)
		}
	}
	return len(m)
}

func MapDeleteAndReadd() int {
	m := map[int]int{1: 10, 2: 20}
	delete(m, 1)
	m[1] = 30
	return m[1]
}

func MapNilDelete() int {
	var m map[int]int
	delete(m, 1) // should not panic
	return 0
}

func MapLenAfterDelete() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	delete(m, 2)
	return len(m)
}

func MapZeroValueAccess() int {
	m := map[string]int{}
	return m["missing"] // returns 0 for missing key
}

func MapBoolKey() int {
	m := map[bool]int{true: 1, false: 0}
	return m[true] + m[false]
}

func MapStringKeyEmpty() int {
	m := map[string]int{"": 1, "a": 2}
	return m[""]
}

func MapIntKeyZero() int {
	m := map[int]int{0: 100, 1: 200}
	return m[0]
}

func MapNestedMap() int {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 42}
	return m["outer"]["inner"]
}

func MapOverwritePreservesType() string {
	m := map[string]any{}
	m["key"] = 42
	m["key"] = "hello"
	return m["key"].(string)
}

func MapCommaOkDelete() int {
	m := map[int]int{1: 10}
	delete(m, 1)
	_, ok := m[1]
	if ok {
		return 1
	}
	return 0
}
