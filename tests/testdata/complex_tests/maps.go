package complex_tests

// ============================================================================
// Complex Map Operations Tests (25 tests)
// ============================================================================

// MapNested tests nested maps.
func MapNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"x": 3, "y": 4},
	}
	return m["a"]["x"] + m["b"]["y"]
}

// MapMerge tests merging maps.
func MapMerge() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	result := make(map[string]int)
	for k, v := range m1 {
		result[k] = v
	}
	for k, v := range m2 {
		result[k] = v
	}
	return result["a"] + result["b"] + result["c"]
}

// MapInvert tests inverting key-value.
func MapInvert() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	result := make(map[string]int)
	for k, v := range m {
		result[v] = k
	}
	return result["a"] + result["b"]*10 + result["c"]*100
}

// MapKeys tests getting keys.
func MapKeys() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	return len(keys)
}

// MapValues tests getting values.
func MapValues() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	values := []string{}
	for _, v := range m {
		values = append(values, v)
	}
	return len(values)
}

// MapFilterKeys tests filtering by key.
func MapFilterKeys() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	result := make(map[int]int)
	for k, v := range m {
		if k%2 == 0 {
			result[k] = v
		}
	}
	return len(result)
}

// MapFilterValues tests filtering by value.
func MapFilterValues() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
	result := make(map[int]int)
	for k, v := range m {
		if v > 25 {
			result[k] = v
		}
	}
	return len(result)
}

// MapMapKeys tests mapping keys.
func MapMapKeys() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	result := make(map[int]int)
	for k, v := range m {
		result[k*10] = v
	}
	return result[10] + result[20] + result[30]
}

// MapMapValues tests mapping values.
func MapMapValues() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	result := make(map[int]int)
	for k, v := range m {
		result[k] = v * 2
	}
	return result[1] + result[2] + result[3]
}

// MapCounter tests using map as counter.
func MapCounter() int {
	str := "hello world"
	counts := make(map[rune]int)
	for _, r := range str {
		counts[r]++
	}
	return counts['l']
}

// MapHistogram tests histogram.
func MapHistogram() int {
	arr := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	hist := make(map[int]int)
	for _, v := range arr {
		hist[v]++
	}
	return hist[4]
}

// MapGroupBy tests grouping.
func MapGroupBy() int {
	arr := []string{"apple", "banana", "apricot", "cherry", "avocado"}
	groups := make(map[byte][]string)
	for _, s := range arr {
		key := s[0]
		groups[key] = append(groups[key], s)
	}
	return len(groups['a'])
}

// MapSet tests set operations with map.
func MapSet() int {
	set := make(map[int]bool)
	add := func(v int) { set[v] = true }
	contains := func(v int) bool { return set[v] }
	remove := func(v int) { delete(set, v) }
	add(1)
	add(2)
	add(3)
	if contains(2) {
		remove(2)
	}
	return len(set)
}

// MapMultiSet tests multiset.
func MapMultiSet() int {
	mset := make(map[int]int)
	add := func(v int) { mset[v]++ }
	remove := func(v int) {
		if mset[v] > 0 {
			mset[v]--
			if mset[v] == 0 {
				delete(mset, v)
			}
		}
	}
	add(1)
	add(1)
	add(2)
	remove(1)
	return mset[1]
}

// MapBiMap tests bidirectional map.
func MapBiMap() int {
	forward := make(map[string]int)
	backward := make(map[int]string)
	insert := func(k string, v int) {
		forward[k] = v
		backward[v] = k
	}
	insert("a", 1)
	insert("b", 2)
	return forward["a"] + len(backward[2])
}

// MapUpdateNested tests updating nested map.
func MapUpdateNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1},
	}
	if inner, ok := m["a"]; ok {
		inner["y"] = 2
	}
	return m["a"]["x"] + m["a"]["y"]
}

// MapDeleteWhileIterating tests safe deletion pattern.
func MapDeleteWhileIterating() int {
	m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	toDelete := []int{}
	for k, v := range m {
		if v%2 == 0 {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(m, k)
	}
	return len(m)
}

// MapComplexKey tests complex key.
func MapComplexKey() int {
	type Key struct {
		X, Y int
	}
	m := make(map[Key]int)
	m[Key{1, 2}] = 10
	m[Key{3, 4}] = 20
	return m[Key{1, 2}] + m[Key{3, 4}]
}

// MapDefaultValue tests default value pattern.
func MapDefaultValue() int {
	m := map[int]int{1: 10, 2: 20}
	get := func(k int) int {
		if v, ok := m[k]; ok {
			return v
		}
		return 0
	}
	return get(1) + get(3)
}

// MapAccumulate tests accumulation.
func MapAccumulate() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// MapFindKey tests finding key by value.
func MapFindKey() int {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	findKey := func(v string) int {
		for k, val := range m {
			if val == v {
				return k
			}
		}
		return -1
	}
	return findKey("b")
}

// MapFrequency tests frequency count.
func MapFrequency() int {
	str := "mississippi"
	freq := make(map[rune]int)
	for _, r := range str {
		freq[r]++
	}
	return freq['s']*10 + freq['i']
}

// MapMemoize tests memoization with map.
func MapMemoize() int {
	cache := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := cache[n]; ok {
			return v
		}
		result := fib(n-1) + fib(n-2)
		cache[n] = result
		return result
	}
	return fib(20)
}

// MapIncrement tests increment operations.
func MapIncrement() int {
	m := make(map[int]int)
	for i := 0; i < 10; i++ {
		m[i%3]++
	}
	return m[0] + m[1] + m[2]
}

// MapCopy tests copying map.
func MapCopy() int {
	m1 := map[int]int{1: 10, 2: 20, 3: 30}
	m2 := make(map[int]int)
	for k, v := range m1 {
		m2[k] = v
	}
	m2[1] = 100
	return m1[1] + m2[1]
}
