package divergence_hunt187

import (
	"fmt"
)

// ============================================================================
// Round 187: Map iteration order independence
// ============================================================================

func MapLength() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	return fmt.Sprintf("%d", len(m))
}

func MapContainsKey() string {
	m := map[string]int{"a": 1, "b": 2}
	_, ok1 := m["a"]
	_, ok2 := m["z"]
	return fmt.Sprintf("%v:%v", ok1, ok2)
}

func MapDelete() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return fmt.Sprintf("%d", len(m))
}

func MapDeleteNonExistent() string {
	m := map[string]int{"a": 1}
	delete(m, "z")
	return fmt.Sprintf("%d", len(m))
}

func MapGetZeroValue() string {
	m := map[string]int{"a": 1}
	v := m["z"]
	return fmt.Sprintf("%d", v)
}

func MapOverwrite() string {
	m := map[string]int{"a": 1}
	m["a"] = 10
	return fmt.Sprintf("%d", m["a"])
}

func MapAddKey() string {
	m := map[string]int{"a": 1}
	m["b"] = 2
	return fmt.Sprintf("%d", len(m))
}

func MapNil() string {
	var m map[string]int
	return fmt.Sprintf("%v", m == nil)
}

func MapEmpty() string {
	m := map[string]int{}
	return fmt.Sprintf("%d", len(m))
}

func MapKeyTypes() string {
	type MyString string
	m := map[MyString]int{"hello": 42}
	return fmt.Sprintf("%d", m["hello"])
}

func MapValueSlice() string {
	m := map[string][]int{
		"a": {1, 2},
		"b": {3, 4},
	}
	return fmt.Sprintf("%d:%d", len(m["a"]), len(m["b"]))
}

func MapClearAll() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return fmt.Sprintf("%d", len(m))
}
