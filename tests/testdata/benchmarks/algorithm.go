package benchmarks

// ============================================================================
// Algorithms
// ============================================================================

func NestedLoops() int {
	sum := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 10; k++ {
				sum = sum + 1
			}
		}
	}
	return sum
}

func BubbleSort() int {
	s := make([]int, 50)
	for i := 0; i < 50; i++ {
		s[i] = 50 - i
	}
	n := len(s)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if s[j] > s[j+1] {
				tmp := s[j]
				s[j] = s[j+1]
				s[j+1] = tmp
			}
		}
	}
	return s[0] + s[49]
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func GCD() int {
	sum := 0
	for i := 1; i <= 100; i++ {
		sum = sum + gcd(i*7, i*13)
	}
	return sum
}

func Sieve() int {
	n := 1000
	sieve := make([]int, n+1)
	for i := 2; i <= n; i++ {
		sieve[i] = 1
	}
	for i := 2; i*i <= n; i++ {
		if sieve[i] == 1 {
			for j := i * i; j <= n; j = j + i {
				sieve[j] = 0
			}
		}
	}
	count := 0
	for i := 2; i <= n; i++ {
		if sieve[i] == 1 {
			count = count + 1
		}
	}
	return count
}
