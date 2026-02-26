package tests

import (
	"testing"

	"gig"

	_ "gig/packages"
)

func TestFunctionCall(t *testing.T) {
	runInt(t, `package main
func add(a int, b int) int { return a + b }
func Compute() int { return add(5, 7) }`, 12)
}

func TestMultipleReturn(t *testing.T) {
	runInt(t, `package main
func swap(a, b int) (int, int) { return b, a }
func Compute() int {
	x, y := swap(3, 7)
	return x + y
}`, 10)
}

func TestMultipleReturnDivmod(t *testing.T) {
	runInt(t, `package main
func divmod(a, b int) (int, int) { return a / b, a % b }
func Compute() int {
	q, r := divmod(17, 5)
	return q*10 + r
}`, 32)
}

func TestRecursionFactorial(t *testing.T) {
	runInt(t, `package main
func factorial(n int) int {
	if n <= 1 { return 1 }
	return n * factorial(n - 1)
}
func Compute() int { return factorial(5) }`, 120)
}

func TestMutualRecursion(t *testing.T) {
	runInt(t, `package main
func isEven(n int) bool {
	if n == 0 { return true }
	return isOdd(n - 1)
}
func isOdd(n int) bool {
	if n == 0 { return false }
	return isEven(n - 1)
}
func Compute() int {
	if isEven(10) { return 1 }
	return 0
}`, 1)
}

func TestFibonacciIterative(t *testing.T) {
	runInt(t, `package main
func fib(n int) int {
	if n <= 1 { return n }
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}
func Compute() int { return fib(20) }`, 6765)
}

func TestFibonacciRecursive(t *testing.T) {
	runInt(t, `package main
func fib(n int) int {
	if n <= 1 { return n }
	return fib(n-1) + fib(n-2)
}
func Compute() int { return fib(15) }`, 610)
}

func TestParametersPassed(t *testing.T) {
	source := `package main
func Compute(a int, b int) int { return a * b }`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("Compute", 6, 7)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result.(int64) != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestVariadicFunction(t *testing.T) {
	runInt(t, `package main
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total = total + n
	}
	return total
}
func Compute() int { return sum(1, 2, 3, 4, 5) }`, 15)
}

func TestFunctionAsValue(t *testing.T) {
	runInt(t, `package main
func apply(f func(int) int, x int) int { return f(x) }
func double(x int) int { return x * 2 }
func triple(x int) int { return x * 3 }
func Compute() int { return apply(double, 5) + apply(triple, 5) }`, 25)
}

func TestHigherOrderMap(t *testing.T) {
	runInt(t, `package main
func mapSlice(s []int, f func(int) int) []int {
	result := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = f(s[i])
	}
	return result
}
func Compute() int {
	s := make([]int, 3)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	doubled := mapSlice(s, func(x int) int { return x * 2 })
	return doubled[0] + doubled[1] + doubled[2]
}`, 12)
}

func TestHigherOrderFilter(t *testing.T) {
	runInt(t, `package main
func count(s []int, pred func(int) bool) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if pred(s[i]) {
			n = n + 1
		}
	}
	return n
}
func Compute() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	s = append(s, 5)
	return count(s, func(x int) bool { return x > 3 })
}`, 2)
}

func TestHigherOrderReduce(t *testing.T) {
	runInt(t, `package main
func reduce(s []int, init int, f func(int, int) int) int {
	acc := init
	for i := 0; i < len(s); i++ {
		acc = f(acc, s[i])
	}
	return acc
}
func Compute() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	return reduce(s, 0, func(a, b int) int { return a + b })
}`, 10)
}
