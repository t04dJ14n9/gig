package tests

import "testing"

func TestClosureCounter(t *testing.T) {
	runInt(t, `package main
func makeCounter() func() int {
	count := 0
	return func() int {
		count = count + 1
		return count
	}
}
func Compute() int {
	counter := makeCounter()
	return counter() + counter() + counter()
}`, 6)
}

func TestClosureCaptureMutation(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 10
	add := func(n int) int {
		x = x + n
		return x
	}
	a := add(5)
	b := add(3)
	return a + b
}`, 33)
}

func TestClosureFactory(t *testing.T) {
	runInt(t, `package main
func makeAdder(base int) func(int) int {
	return func(x int) int { return base + x }
}
func Compute() int {
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	return add5(3) + add10(7)
}`, 25)
}

func TestClosureMultipleInstances(t *testing.T) {
	runInt(t, `package main
func makeCounter() func() int {
	count := 0
	return func() int {
		count = count + 1
		return count
	}
}
func Compute() int {
	c1 := makeCounter()
	c2 := makeCounter()
	a := c1()
	b := c1()
	c := c2()
	return a + b + c
}`, 4)
}

func TestClosureOverLoop(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	adder := func(x int) int {
		sum = sum + x
		return sum
	}
	for i := 1; i <= 5; i++ {
		adder(i)
	}
	return sum
}`, 15)
}

func TestClosureChain(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }
	return addOne(double(5))
}`, 11)
}

func TestClosureAccumulator(t *testing.T) {
	runInt(t, `package main
func makeAccumulator(init int) func(int) int {
	total := init
	return func(n int) int {
		total = total + n
		return total
	}
}
func Compute() int {
	acc := makeAccumulator(100)
	acc(10)
	acc(20)
	return acc(30)
}`, 160)
}
