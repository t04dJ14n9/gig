package closures_advanced

// Generator tests closure generator pattern
func Generator() int {
	gen := makeRange(0, 3)
	a := gen()
	b := gen()
	c := gen()
	d := gen()
	return a + b + c + d
}

func makeRange(start, step int) func() int {
	current := start
	return func() int {
		val := current
		current = current + step
		return val
	}
}

// Predicate tests closure predicate pattern
func Predicate() int {
	above50 := makeThreshold(50)
	sum := 0
	for i := 0; i <= 100; i = i + 10 {
		sum = sum + above50(i)
	}
	return sum
}

func makeThreshold(limit int) func(int) int {
	return func(x int) int {
		if x > limit {
			return 1
		}
		return 0
	}
}

// StateMachine tests closure state machine
func StateMachine() int {
	add, getSum, getCount := makeStats()
	_ = add(10)
	_ = add(20)
	_ = add(30)
	return getSum()*10 + getCount()
}

func makeStats() (func(int) int, func() int, func() int) {
	sum := 0
	count := 0
	add := func(x int) int {
		sum = sum + x
		count = count + 1
		return sum
	}
	getSum := func() int { return sum }
	getCount := func() int { return count }
	return add, getSum, getCount
}

// RecursiveHelper tests recursive helper with closure
func RecursiveHelper() int { return sumTree(4) }

func sumTree(depth int) int {
	if depth <= 0 {
		return 1
	}
	left := sumTree(depth - 1)
	right := sumTree(depth - 1)
	return left + right + 1
}

// ApplyN tests applying function N times
func ApplyN() int {
	double := func(x int) int { return x * 2 }
	return applyN(double, 1, 10)
}

func applyN(f func(int) int, x int, n int) int {
	result := x
	for i := 0; i < n; i++ {
		result = f(result)
	}
	return result
}

// Compose tests function composition
func Compose() int {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	square := func(x int) int { return x * x }
	return square(double(addOne(5)))
}
