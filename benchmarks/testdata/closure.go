package main

func ClosureCalls() int {
	sum := 0
	adder := func(x int) int {
		sum = sum + x
		return sum
	}
	for i := 0; i < 1000; i++ {
		adder(i)
	}
	return sum
}
