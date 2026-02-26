package mapadvanced

// LookupExistingKey tests map lookup of existing key
func LookupExistingKey() int {
	m := make(map[string]int)
	m["key"] = 42
	return m["key"]
}

// LookupWithDefault tests map lookup with default
func LookupWithDefault() int {
	m := make(map[string]int)
	m["a"] = 10
	m["b"] = 20
	return m["a"] + m["b"]
}

// AsCounter tests map as counter
func AsCounter() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 2)
	s = append(s, 1)
	s = append(s, 2)

	counts := make(map[int]int)
	counts[1] = 0
	counts[2] = 0
	counts[3] = 0
	for _, v := range s {
		counts[v] = counts[v] + 1
	}
	return counts[1]*100 + counts[2]*10 + counts[3]
}

// WithStringValues tests map with string values
func WithStringValues() string {
	m := make(map[int]string)
	m[1] = "one"
	m[2] = "two"
	m[3] = "three"
	return m[1] + "-" + m[2] + "-" + m[3]
}

// BuildFromLoop tests building map from loop
func BuildFromLoop() int {
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i * i
	}
	return m[10] + m[50]
}

// DeleteAndReinsert tests delete and reinsert
func DeleteAndReinsert() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	delete(m, "a")
	m["a"] = 99
	return m["a"] + m["b"]
}
