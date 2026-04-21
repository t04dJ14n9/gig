package divergence_hunt222

import "fmt"

// ============================================================================
// Round 222: Map delete during iteration
// ============================================================================

// MapDeleteDuringIteration deletes current element during iteration
func MapDeleteDuringIteration() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	for k := range m {
		if k == "b" || k == "d" {
			delete(m, k)
		}
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapDeleteOtherDuringIteration deletes other elements during iteration.
// The Go spec permits implementations to either produce or not produce the
// deleted entry, so sum is non-deterministic. Only the final length is defined.
func MapDeleteOtherDuringIteration() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	for k := range m {
		if k == "a" {
			delete(m, "c")
			break
		}
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapDeleteAllDuringIteration attempts to delete all during iteration
func MapDeleteAllDuringIteration() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapDeleteAndAddDuringIteration deletes and adds during iteration
func MapDeleteAndAddDuringIteration() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for k, v := range m {
		delete(m, k)
		m[k+"_new"] = v * 10
		count++
		if count > 10 {
			break
		}
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapSafeDeleteCollect collects keys then deletes
func MapSafeDeleteCollect() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	toDelete := []string{}
	for k := range m {
		if m[k]%2 == 0 {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(m, k)
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapDeleteNonExistent deletes non-existent key during iteration
func MapDeleteNonExistent() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for k := range m {
		delete(m, k+"_missing")
		count++
	}
	return fmt.Sprintf("count=%d,len=%d", count, len(m))
}

// MapDeleteThenLookup deletes then looks up same key
func MapDeleteThenLookup() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	results := []int{}
	for k := range m {
		delete(m, k)
		v, ok := m[k]
		results = append(results, v)
		if !ok {
			results = append(results, -1)
		}
	}
	return fmt.Sprintf("results=%d", len(results))
}

// MapIterDeleteWithComplicatedKeys deletes with int keys during iteration
func MapIterDeleteWithComplicatedKeys() string {
	m := map[int]string{10: "a", 20: "b", 30: "c", 40: "d", 50: "e"}
	for k := range m {
		if k > 25 {
			delete(m, k)
		}
	}
	sum := 0
	for k := range m {
		sum += k
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapDeletePreservesIteration deletes but continues iteration
func MapDeletePreservesIteration() string {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	sum := 0
	for k, v := range m {
		if k%2 == 0 {
			delete(m, k)
		}
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapConditionalDeleteDuringIteration conditionally deletes elements
func MapConditionalDeleteDuringIteration() string {
	m := map[string]int{"apple": 5, "banana": 6, "cherry": 7, "date": 4}
	for k := range m {
		if len(k) > 5 {
			delete(m, k)
		}
	}
	return fmt.Sprintf("len=%d", len(m))
}
