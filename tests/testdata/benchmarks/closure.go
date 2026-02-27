package benchmarks

// ============================================================================
// Closure Operations
// ============================================================================

func ClosureCalls() int {
	sum := 0
	adder := func(x int) int {
		sum = sum + x
		return sum
	}
	for i := 0; i < 1000; i++ {
		adder(i)
	}
	return sum
}

func HigherOrder() int {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}
	double := func(x int) int { return x * 2 }
	sum := 0
	for _, v := range s {
		sum = sum + double(v)
	}
	return sum
}
