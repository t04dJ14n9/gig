package tests

import "testing"

// Advanced closure tests: generators, iterators, memoization patterns.

func TestClosureGenerator(t *testing.T) {
	runInt(t, `package main
func makeRange(start, step int) func() int {
	current := start
	return func() int {
		val := current
		current = current + step
		return val
	}
}
func Compute() int {
	gen := makeRange(0, 3)
	a := gen()
	b := gen()
	c := gen()
	d := gen()
	return a + b + c + d
}`, 18) // 0+3+6+9
}

func TestClosurePredicate(t *testing.T) {
	runInt(t, `package main
func makeThreshold(limit int) func(int) int {
	return func(x int) int {
		if x > limit {
			return 1
		}
		return 0
	}
}
func Compute() int {
	above50 := makeThreshold(50)
	sum := 0
	for i := 0; i <= 100; i = i + 10 {
		sum = sum + above50(i)
	}
	return sum
}`, 5) // 60,70,80,90,100 -> 5 values above 50
}

func TestClosureStateMachine(t *testing.T) {
	// Simple state machine: accumulate values, track count
	runInt(t, `package main
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
func Compute() int {
	add, getSum, getCount := makeStats()
	_ = add(10)
	_ = add(20)
	_ = add(30)
	return getSum()*10 + getCount()
}`, 603) // sum=60, count=3 -> 60*10+3
}

func TestClosureRecursiveHelper(t *testing.T) {
	// Use a closure inside a function to help with recursion
	runInt(t, `package main
func sumTree(depth int) int {
	if depth <= 0 {
		return 1
	}
	left := sumTree(depth - 1)
	right := sumTree(depth - 1)
	return left + right + 1
}
func Compute() int { return sumTree(4) }`, 31) // 2^5 - 1
}

func TestClosureApplyN(t *testing.T) {
	// Apply a function N times
	runInt(t, `package main
func applyN(f func(int) int, x int, n int) int {
	result := x
	for i := 0; i < n; i++ {
		result = f(result)
	}
	return result
}
func Compute() int {
	double := func(x int) int { return x * 2 }
	return applyN(double, 1, 10)
}`, 1024)
}

func TestClosureCompose(t *testing.T) {
	// Manual function composition
	runInt(t, `package main
func Compute() int {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }
	square := func(x int) int { return x * x }
	// compose manually: square(double(addOne(5)))
	return square(double(addOne(5)))
}`, 144) // (5+1)*2 = 12, 12^2 = 144
}
