package complex_tests

// ClosureCaptureLoop tests capturing loop variables correctly.
func ClosureCaptureLoop() int {
	fns := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		v := i
		fns[i] = func() int { return v }
	}
	return fns[0]() + fns[1]()*10 + fns[2]()*100
}

// ClosureMutualRecursion tests closures that call each other.
func ClosureMutualRecursion() int {
	var isEven func(int) bool
	var isOdd func(int) bool
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
		return 1
	}
	return 0
}

// ClosureCurryAdd tests currying with closures.
func ClosureCurryAdd() int {
	add := func(a int) func(int) int {
		return func(b int) int {
			return a + b
		}
	}
	return add(10)(20)
}

// ClosureCounterChain tests chained counter closures.
func ClosureCounterChain() int {
	makeCounter := func(start int) func() int {
		count := start
		return func() int {
			count++
			return count
		}
	}
	c1 := makeCounter(0)
	c2 := makeCounter(100)
	return c1() + c1() + c2() + c2()
}

// ClosureMemoize tests memoization pattern.
func ClosureMemoize() int {
	memo := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		result := fib(n-1) + fib(n-2)
		memo[n] = result
		return result
	}
	return fib(20)
}

// ClosureYieldGenerator tests generator pattern.
func ClosureYieldGenerator() int {
	generate := func(start, end int) func() (int, bool) {
		current := start
		return func() (int, bool) {
			if current > end {
				return 0, false
			}
			v := current
			current++
			return v, true
		}
	}
	gen := generate(1, 5)
	sum := 0
	for {
		v, ok := gen()
		if !ok {
			break
		}
		sum += v
	}
	return sum
}

// ClosureFilter tests filter pattern.
func ClosureFilter() int {
	makeFilter := func(predicate func(int) bool) func([]int) []int {
		return func(nums []int) []int {
			result := []int{}
			for _, n := range nums {
				if predicate(n) {
					result = append(result, n)
				}
			}
			return result
		}
	}
	isEven := func(n int) bool { return n%2 == 0 }
	filter := makeFilter(isEven)
	filtered := filter([]int{1, 2, 3, 4, 5, 6})
	return len(filtered)
}

// ClosureReduce tests reduce pattern.
func ClosureReduce() int {
	makeReducer := func(accumulate func(int, int) int, initial int) func([]int) int {
		return func(nums []int) int {
			result := initial
			for _, n := range nums {
				result = accumulate(result, n)
			}
			return result
		}
	}
	sum := func(a, b int) int { return a + b }
	reducer := makeReducer(sum, 0)
	return reducer([]int{1, 2, 3, 4, 5})
}

// ClosureCompose tests function composition.
func ClosureCompose() int {
	compose := func(f, g func(int) int) func(int) int {
		return func(x int) int {
			return f(g(x))
		}
	}
	double := func(x int) int { return x * 2 }
	increment := func(x int) int { return x + 1 }
	composed := compose(double, increment)
	return composed(5)
}

// ClosurePartial tests partial application.
func ClosurePartial() int {
	partial := func(fn func(int, int) int, a int) func(int) int {
		return func(b int) int {
			return fn(a, b)
		}
	}
	add := func(a, b int) int { return a + b }
	add5 := partial(add, 5)
	return add5(10)
}

// ClosureOnce tests single execution pattern.
func ClosureOnce() int {
	once := func(fn func() int) func() int {
		called := false
		var result int
		return func() int {
			if !called {
				result = fn()
				called = true
			}
			return result
		}
	}
	expensive := func() int { return 42 }
	onceFn := once(expensive)
	return onceFn() + onceFn() + onceFn()
}

// ClosureState tests stateful closure.
func ClosureState() int {
	makeAccumulator := func() func(int) int {
		sum := 0
		return func(n int) int {
			sum += n
			return sum
		}
	}
	acc := makeAccumulator()
	return acc(1) + acc(2) + acc(3)
}

// ClosureDefer tests closure with defer.
func ClosureDefer() int {
	result := 0
	func() {
		defer func() {
			result += 10
		}()
		result += 1
	}()
	return result
}

// ClosureCaptureSlice tests capturing slice.
func ClosureCaptureSlice() int {
	s := []int{1, 2, 3}
	push := func(v int) {
		s = append(s, v)
	}
	push(4)
	push(5)
	return s[len(s)-1]
}

// ClosureCaptureMap tests capturing map.
func ClosureCaptureMap() int {
	m := make(map[string]int)
	set := func(k string, v int) {
		m[k] = v
	}
	set("a", 1)
	set("b", 2)
	return m["a"] + m["b"]
}

// ClosureRecursive tests recursive closure.
func ClosureRecursive() int {
	var factorial func(int) int
	factorial = func(n int) int {
		if n <= 1 {
			return 1
		}
		return n * factorial(n-1)
	}
	return factorial(6)
}

// ClosureInStruct tests closure in struct.
func ClosureInStruct() int {
	type Processor struct {
		Process func(int) int
	}
	p := Processor{
		Process: func(x int) int { return x * 2 },
	}
	return p.Process(21)
}

// ClosureSliceOfFuncs tests slice of closures.
func ClosureSliceOfFuncs() int {
	fns := make([]func(int) int, 5)
	for i := 0; i < 5; i++ {
		factor := i + 1
		fns[i] = func(x int) int { return x * factor }
	}
	sum := 0
	for _, fn := range fns {
		sum += fn(10)
	}
	return sum
}

// ClosureMapOfFuncs tests map of closures.
func ClosureMapOfFuncs() int {
	ops := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	}
	return ops["add"](1, 2) + ops["mul"](3, 4)
}
