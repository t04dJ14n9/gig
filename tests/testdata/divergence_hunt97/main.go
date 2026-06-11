package divergence_hunt97

import "fmt"

// ============================================================================
// Round 97: Map iteration and deletion patterns
// ============================================================================

func MapDeleteKey() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return fmt.Sprintf("%d", len(m))
}

func MapDeleteNonExistent() string {
	m := map[string]int{"a": 1}
	delete(m, "z")
	return fmt.Sprintf("%d", len(m))
}

func MapDoubleDelete() string {
	m := map[string]int{"a": 1}
	delete(m, "a")
	delete(m, "a")
	return fmt.Sprintf("%d", len(m))
}

func MapClear() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return fmt.Sprintf("%d", len(m))
}

func MapAccessMissing() string {
	m := map[string]int{"a": 1}
	v, ok := m["z"]
	return fmt.Sprintf("%d:%v", v, ok)
}

func MapSetDefault() string {
	m := map[string]int{}
	v, ok := m["key"]
	if !ok {
		m["key"] = 42
		v = 42
	}
	return fmt.Sprintf("%d:%v", v, ok)
}

func MapCountValues() string {
	m := map[string]int{"a": 1, "b": 2, "c": 1, "d": 2}
	counts := map[int]int{}
	for _, v := range m {
		counts[v]++
	}
	return fmt.Sprintf("%d:%d", counts[1], counts[2])
}

func MapInvert() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	inv := map[int]string{}
	for k, v := range m {
		inv[v] = k
	}
	return fmt.Sprintf("%s", inv[2])
}

func MapMerge() string {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	for k, v := range m2 {
		m1[k] = v
	}
	return fmt.Sprintf("%d:%d:%d", m1["a"], m1["b"], m1["c"])
}

func MapKeys() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return fmt.Sprintf("%d", len(keys))
}

func MapNestedAccess() string {
	type Inner struct{ Val int }
	m := map[string]Inner{"x": {Val: 42}}
	return fmt.Sprintf("%d", m["x"].Val)
}
