package divergence_hunt156

import "fmt"

// ============================================================================
// Round 156: Map advanced patterns
// ============================================================================

// MapZeroValue tests map zero value behavior
func MapZeroValue() string {
	var m map[string]int
	v := m["missing"] // Should return 0, not panic
	return fmt.Sprintf("zero=%d-nil=%t", v, m == nil)
}

// MapEmptyVsNil tests empty vs nil maps
func MapEmptyVsNil() string {
	var nilMap map[string]int
	emptyMap := map[string]int{}
	return fmt.Sprintf("nil-nil=%t-empty-nil=%t-nil-len=%d-empty-len=%d",
		nilMap == nil, emptyMap == nil, len(nilMap), len(emptyMap))
}

// MapMakeWithCapacity tests make with capacity hint
func MapMakeWithCapacity() string {
	m := make(map[string]int, 100)
	m["key"] = 42
	return fmt.Sprintf("len=%d", len(m))
}

// MapStringKey tests string keys
func MapStringKey() string {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	return fmt.Sprintf("sum=%d", m["one"]+m["two"]+m["three"])
}

// MapIntKey tests int keys
func MapIntKey() string {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	return fmt.Sprintf("%s%s%s", m[1], m[2], m[3])
}

// MapStructKey tests struct keys
func MapStructKey() string {
	type Point struct{ X, Y int }
	m := map[Point]string{
		{1, 2}: "A",
		{3, 4}: "B",
	}
	return fmt.Sprintf("A=%s-B=%s", m[Point{1, 2}], m[Point{3, 4}])
}

// MapArrayKey tests array keys
func MapArrayKey() string {
	m := map[[3]int]string{
		{1, 2, 3}: "first",
		{4, 5, 6}: "second",
	}
	return fmt.Sprintf("first=%s", m[[3]int{1, 2, 3}])
}

// MapPointerKey tests pointer keys
func MapPointerKey() string {
	k1 := &struct{ Val int }{1}
	k2 := &struct{ Val int }{1}
	m := map[*struct{ Val int }]string{
		k1: "one",
	}
	// k1 and k2 have same value but different pointers
	v1 := m[k1]
	v2 := m[k2]
	return fmt.Sprintf("v1=%s-v2=%s", v1, v2)
}

// MapInterfaceKey tests interface keys
func MapInterfaceKey() string {
	m := make(map[interface{}]string)
	m[1] = "int"
	m["hello"] = "string"
	m[true] = "bool"
	return fmt.Sprintf("int=%s-str=%s", m[1], m["hello"])
}

// MapSliceValue tests map with slice values
func MapSliceValue() string {
	m := map[string][]int{
		"even": {2, 4, 6},
		"odd":  {1, 3, 5},
	}
	return fmt.Sprintf("even-len=%d-odd-len=%d", len(m["even"]), len(m["odd"]))
}

// MapMapValue tests map with map values
func MapMapValue() string {
	m := map[string]map[string]int{
		"math":    {"alice": 90, "bob": 85},
		"science": {"alice": 88, "bob": 92},
	}
	return fmt.Sprintf("alice-math=%d", m["math"]["alice"])
}
