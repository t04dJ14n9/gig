package divergence_hunt40

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 40: Map edge cases - various key/value types, nil, delete, iteration
// ============================================================================

func MapIntKey() int {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	return len(m)
}

func MapFloatKey() int {
	m := map[float64]string{1.0: "a", 2.5: "b"}
	return len(m)
}

func MapBoolKey() int {
	m := map[bool]int{true: 1, false: 0}
	return m[true] + m[false]
}

func MapStructKey() int {
	type Key struct{ X, Y int }
	m := map[Key]string{{1, 2}: "a", {3, 4}: "b"}
	return len(m)
}

func MapStringKey() int {
	m := map[string]int{"hello": 5, "world": 5}
	return m["hello"] + m["world"]
}

func MapWithSliceValue() int {
	m := map[string][]int{}
	m["a"] = []int{1, 2}
	m["a"] = append(m["a"], 3)
	return len(m["a"])
}

func MapWithMapValue() int {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 42}
	return m["outer"]["inner"]
}

func MapDeleteAndLen() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return len(m)
}

func MapDeleteNonExistent() int {
	m := map[string]int{"a": 1}
	delete(m, "nonexistent")
	return len(m)
}

func MapNilDelete() int {
	var m map[string]int
	delete(m, "key") // no-op on nil map
	return len(m)
}

func MapNilLookup() int {
	var m map[string]int
	return m["key"] // returns zero value
}

func MapCommaOkPresent() int {
	m := map[string]int{"a": 1}
	v, ok := m["a"]
	if ok { return v }
	return -1
}

func MapCommaOkMissing() int {
	m := map[string]int{"a": 1}
	v, ok := m["b"]
	if ok { return v }
	return -1
}

func MapOverwrite() int {
	m := map[string]int{"a": 1}
	m["a"] = 2
	return m["a"]
}

func MapIterationSum() int {
	m := map[string]int{"a": 10, "b": 20, "c": 30}
	sum := 0
	for _, v := range m { sum += v }
	return sum
}

func MapMakeWithSize() int {
	m := make(map[string]int, 10)
	m["x"] = 42
	return m["x"]
}

func MapEmptyLiteral() int {
	m := map[string]int{}
	m["key"] = 1
	return len(m)
}

func JSONRoundTripMap() int {
	m := map[string]int{"x": 10, "y": 20}
	data, _ := json.Marshal(m)
	var decoded map[string]int
	json.Unmarshal(data, &decoded)
	return decoded["x"] + decoded["y"]
}

func FmtMap() string {
	return fmt.Sprintf("%d", len(map[string]int{"a": 1, "b": 2}))
}

func SortMapKeys() int {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	sort.Strings(keys)
	return len(keys)
}

func MapStringJoin() string {
	m := map[string]int{"x": 1, "y": 2}
	parts := []string{}
	for k := range m { parts = append(parts, k) }
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
