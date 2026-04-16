package divergence_hunt65

// ============================================================================
// Round 65: Complex closures - shared state, recursive closures, closure chains
// ============================================================================

func ClosureCounter() int {
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
	return c()
}

func ClosureSharedState() int {
	x := 0
	f1 := func() { x += 10 }
	f2 := func() { x += 1 }
	f1()
	f2()
	return x
}

func ClosureChain() int {
	f := func(x int) int { return x + 1 }
	g := func(x int) int { return f(x) * 2 }
	h := func(x int) int { return g(x) + 3 }
	return h(5)
}

func ClosureOverLoopVar() int {
	var fns []func() int
	for i := 0; i < 5; i++ {
		v := i // capture by copy
		fns = append(fns, func() int { return v })
	}
	sum := 0
	for _, f := range fns {
		sum += f()
	}
	return sum
}

func RecursiveClosure() int {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fib(10)
}

func ClosureReturnClosure() int {
	makeAdder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := makeAdder(5)
	return add5(3)
}

func ClosureCaptureSlice() []int {
	s := []int{1, 2, 3}
	f := func() {
		s[0] = 10
	}
	f()
	return s
}

func ClosureCaptureMap() int {
	m := map[string]int{"a": 1}
	f := func() {
		m["b"] = 2
	}
	f()
	return len(m)
}

func ClosureMultipleReturns() (int, int) {
	split := func(x int) (int, int) {
		return x / 10, x % 10
	}
	return split(42)
}

func ClosureCurry() int {
	add := func(a int) func(int) int {
		return func(b int) int {
			return a + b
		}
	}
	return add(10)(20)
}

func ClosureCaptureModify() int {
	x := 1
	inc := func() int {
		x *= 2
		return x
	}
	inc()
	inc()
	return x
}

func ClosureNoCapture() int {
	f := func(x int) int { return x * x }
	return f(7)
}

func ClosureAsArg() int {
	apply := func(f func(int) int, x int) int {
		return f(x)
	}
	double := func(x int) int { return x * 2 }
	return apply(double, 5)
}

func ClosureSliceMap() int {
	doubler := func(x int) int { return x * 2 }
	data := []int{1, 2, 3, 4, 5}
	result := 0
	for _, v := range data {
		result += doubler(v)
	}
	return result
}
