package tests

import "testing"

// Tests for multiple assignment and tuple patterns.

func TestMultiAssignSwap(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := 10
	b := 20
	a, b = b, a
	return a*100 + b
}`, 2010)
}

func TestMultiAssignFromFunction(t *testing.T) {
	runInt(t, `package main
func twoVals() (int, int) { return 42, 58 }
func Compute() int {
	a, b := twoVals()
	return a + b
}`, 100)
}

func TestMultiAssignThreeValues(t *testing.T) {
	runInt(t, `package main
func threeVals(x int) (int, int, int) {
	return x, x*2, x*3
}
func Compute() int {
	a, b, c := threeVals(10)
	return a + b + c
}`, 60)
}

func TestMultiAssignInLoop(t *testing.T) {
	// Fibonacci via multiple assignment
	runInt(t, `package main
func Compute() int {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		a, b = b, a+b
	}
	return a
}`, 55)
}

func TestDiscardWithBlank(t *testing.T) {
	runInt(t, `package main
func divmod(a, b int) (int, int) { return a / b, a % b }
func Compute() int {
	q, _ := divmod(17, 5)
	return q
}`, 3)
}
