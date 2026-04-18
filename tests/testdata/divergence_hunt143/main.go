package divergence_hunt143

import "fmt"

// ============================================================================
// Round 143: Map type assertion and interface maps
// ============================================================================

func MapInterfaceValues() string {
	m := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"admin": true,
	}
	return fmt.Sprintf("name=%s", m["name"])
}

func MapInterfaceTypeSwitch() string {
	m := map[string]interface{}{
		"int":    42,
		"string": "hello",
		"bool":   true,
	}
	counts := map[string]int{}
	for _, v := range m {
		switch v.(type) {
		case int:
			counts["int"]++
		case string:
			counts["string"]++
		case bool:
			counts["bool"]++
		}
	}
	return fmt.Sprintf("int:%d-string:%d-bool:%d", counts["int"], counts["string"], counts["bool"])
}

func MapInterfaceAssertion() string {
	m := map[string]interface{}{"val": 42.0}
	if v, ok := m["val"].(float64); ok {
		return fmt.Sprintf("float=%v", v)
	}
	return "not-float"
}

func MapStringSlice() string {
	m := map[string][]string{
		"fruits": {"apple", "banana"},
		"colors": {"red", "blue"},
	}
	return fmt.Sprintf("len=%d", len(m["fruits"]))
}

func MapStringFunc() string {
	ops := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
	return fmt.Sprintf("5+3=%d", ops["add"](5, 3))
}

func MapDeleteAndRead() string {
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	v, ok := m["a"]
	return fmt.Sprintf("val=%d-ok=%t", v, ok)
}

func MapLengthAfterDelete() string {
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	delete(m, "y")
	return fmt.Sprintf("len=%d", len(m))
}

func MapNilVsEmptyAccess() string {
	var m map[string]int
	v := m["key"]
	return fmt.Sprintf("val=%d", v)
}

func MapCompositeLiteral() string {
	m := map[[2]int]string{
		{1, 2}: "one-two",
		{3, 4}: "three-four",
	}
	return m[[2]int{1, 2}]
}
