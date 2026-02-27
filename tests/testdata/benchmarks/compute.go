package benchmarks

// ============================================================================
// Arithmetic & Fibonacci
// ============================================================================

func ArithmeticSum() int {
	sum := 0
	for i := 1; i <= 1000; i++ {
		sum = sum + i
	}
	return sum
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func FibRecursive() int { return fib(15) }

func FibIterative() int {
	a, b := 0, 1
	for i := 0; i < 50; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func Factorial() int { return factorial(12) }
