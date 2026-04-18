package divergence_hunt124

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 124: Map operations — iteration, delete, len, nil map
// ============================================================================

func MapLiteral() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	return fmt.Sprintf("len=%d", len(m))
}

func MapDeleteLen() string {
	m := map[string]int{"x": 10, "y": 20}
	delete(m, "x")
	return fmt.Sprintf("len=%d-val=%d", len(m), m["y"])
}

func MapNilWrite() string {
	var m map[string]int
	// Writing to nil map panics — we test reading instead
	v := m["key"]
	return fmt.Sprintf("val=%d", v)
}

func MapZeroValue() string {
	m := map[string]int{}
	v := m["nonexistent"]
	return fmt.Sprintf("val=%d", v)
}

func MapOkCheck() string {
	m := map[string]int{"a": 1}
	v, ok := m["a"]
	v2, ok2 := m["b"]
	return fmt.Sprintf("a=%d-%t-b=%d-%t", v, ok, v2, ok2)
}

func MapSortedKeys() string {
	m := map[string]int{"banana": 2, "apple": 1, "cherry": 3}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func MapIntKey() string {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	return fmt.Sprintf("2=%s", m[2])
}

func MapNestedMap() string {
	outer := map[string]map[string]int{}
	outer["inner"] = map[string]int{"x": 42}
	return fmt.Sprintf("val=%d", outer["inner"]["x"])
}

func MapUpdateValue() string {
	m := map[string]int{"a": 1}
	m["a"] = 99
	return fmt.Sprintf("val=%d", m["a"])
}

func MapBoolValue() string {
	m := map[string]bool{"active": true, "closed": false}
	return fmt.Sprintf("active=%t-closed=%t", m["active"], m["closed"])
}
