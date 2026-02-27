package benchmarks

import "strconv"

// ============================================================================
// Slice & Map Operations
// ============================================================================

func SliceAppend() int {
	s := make([]int, 0)
	for i := 0; i < 1000; i++ {
		s = append(s, i)
	}
	return len(s)
}

func SliceSum() int {
	s := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		s[i] = i
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}

func MapOps() int {
	m := make(map[string]int)
	for i := 0; i < 100; i++ {
		m[strconv.Itoa(i)] = i
	}
	sum := 0
	for _, v := range m {
		sum = sum + v
	}
	return sum
}
