package divergence_hunt221

import "fmt"

// ============================================================================
// Round 221: Map iteration patterns
// ============================================================================

// MapIterationSum tests iterating over map to sum values
func MapIterationSum() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapIterationCount tests counting map entries
func MapIterationCount() string {
	m := map[int]string{1: "one", 2: "two", 3: "three", 4: "four"}
	count := 0
	for range m {
		count++
	}
	return fmt.Sprintf("count=%d", count)
}

// MapIterationKeys collects keys during iteration
func MapIterationKeys() string {
	m := map[int]bool{10: true, 20: true, 30: true}
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	return fmt.Sprintf("len=%d", len(keys))
}

// MapIterationValues collects values during iteration
func MapIterationValues() string {
	m := map[string]int{"x": 100, "y": 200, "z": 300}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapIterationOrder tests that iteration order is not guaranteed
func MapIterationOrder() string {
	m := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4}
	result := 0
	for i := 0; i < 5; i++ {
		for k := range m {
			result += k
		}
	}
	return fmt.Sprintf("result=%d", result)
}

// MapIterationBreak tests breaking from map iteration
func MapIterationBreak() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	count := 0
	for _, v := range m {
		count += v
		if count >= 6 {
			break
		}
	}
	return fmt.Sprintf("count=%d", count)
}

// MapIterationContinue tests continue in map iteration
func MapIterationContinue() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	sum := 0
	for k, v := range m {
		if k == "c" {
			continue
		}
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapIterationNested tests nested map iteration
func MapIterationNested() string {
	m1 := map[int]int{1: 10, 2: 20}
	m2 := map[int]int{1: 100, 2: 200}
	sum := 0
	for _, v1 := range m1 {
		for _, v2 := range m2 {
			sum += v1 + v2
		}
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapIterationWithModify tests reading from map during iteration
func MapIterationWithModify() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := 0
	for k, v := range m {
		m[k] = v * 2
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapIterationEmpty tests iteration over empty map
func MapIterationEmpty() string {
	m := map[string]int{}
	count := 0
	for range m {
		count++
	}
	return fmt.Sprintf("count=%d", count)
}

// MapIterationNil tests iteration over nil map
func MapIterationNil() string {
	var m map[string]int
	count := 0
	for range m {
		count++
	}
	return fmt.Sprintf("count=%d", count)
}
