package divergence_hunt133

import "fmt"

// ============================================================================
// Round 133: Recursive closures and mutual recursion
// ============================================================================

func RecursiveClosure() string {
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fmt.Sprintf("fib10=%d", fib(10))
}

func ClosureCounter() string {
	makeCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}
	c := makeCounter()
	c()
	c()
	return fmt.Sprintf("count=%d", c())
}

func ClosureCapture() string {
	x := 10
	f := func() int {
		return x
	}
	x = 20
	return fmt.Sprintf("captured=%d", f())
}

func ClosureParamCapture() string {
	makeAdder := func(base int) func(int) int {
		return func(x int) int {
			return base + x
		}
	}
	add5 := makeAdder(5)
	return fmt.Sprintf("5+3=%d", add5(3))
}

func ClosureSliceCapture() string {
	funcs := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		i := i // capture loop variable
		funcs[i] = func() int { return i }
	}
	return fmt.Sprintf("%d-%d-%d", funcs[0](), funcs[1](), funcs[2]())
}

func ClosureSliceCaptureNoCopy() string {
	funcs := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		funcs[i] = func() int { return i }
	}
	// All closures share the same i, which ends at 3
	return fmt.Sprintf("%d-%d-%d", funcs[0](), funcs[1](), funcs[2]())
}

func MutualRecursion() string {
	isEven := func(n int) bool { return true }
	isOdd := func(n int) bool { return false }

	isEven = func(n int) bool {
		if n == 0 {
			return true
		}
		return isOdd(n - 1)
	}
	isOdd = func(n int) bool {
		if n == 0 {
			return false
		}
		return isEven(n - 1)
	}

	if isEven(10) && isOdd(7) {
		return "correct"
	}
	return "wrong"
}

func ClosureAsParam() string {
	apply := func(f func(int) int, x int) int {
		return f(x)
	}
	double := func(x int) int { return x * 2 }
	return fmt.Sprintf("double(5)=%d", apply(double, 5))
}

func ClosureReturnClosure() string {
	makeMultiplier := func(factor int) func(int) int {
		return func(x int) int {
			return x * factor
		}
	}
	triple := makeMultiplier(3)
	return fmt.Sprintf("3*7=%d", triple(7))
}
