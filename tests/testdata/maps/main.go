package maps

// BasicOps tests basic map operations
func BasicOps() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	return m["a"] + m["b"]
}

// Iteration tests map iteration
func Iteration() int {
	m := make(map[string]int)
	m["x"] = 10
	m["y"] = 20
	m["z"] = 30
	sum := 0
	for _, v := range m {
		sum = sum + v
	}
	return sum
}

// Delete tests map delete
func Delete() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	delete(m, "b")
	return len(m)
}

// Len tests map length
func Len() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	m["d"] = 4
	return len(m)
}

// Overwrite tests map overwrite
func Overwrite() int {
	m := make(map[string]int)
	m["key"] = 10
	m["key"] = 42
	return m["key"]
}

// IntKeys tests map with int keys
func IntKeys() int {
	m := make(map[int]int)
	m[1] = 10
	m[2] = 20
	m[3] = 30
	return m[1] + m[2] + m[3]
}

// PassToFunction tests passing map to function
func PassToFunction() int {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200
	return sumValues(m)
}

func sumValues(m map[string]int) int {
	total := 0
	for _, v := range m {
		total = total + v
	}
	return total
}
