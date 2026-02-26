package closures

// Counter tests closure counter
func Counter() int {
	counter := makeCounter()
	return counter() + counter() + counter()
}

func makeCounter() func() int {
	count := 0
	return func() int {
		count = count + 1
		return count
	}
}

// CaptureMutation tests closure capture mutation
func CaptureMutation() int {
	x := 10
	add := func(n int) int {
		x = x + n
		return x
	}
	a := add(5)
	b := add(3)
	return a + b
}

// Factory tests closure factory pattern
func Factory() int {
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	return add5(3) + add10(7)
}

func makeAdder(base int) func(int) int {
	return func(x int) int { return base + x }
}

// MultipleInstances tests multiple closure instances
func MultipleInstances() int {
	c1 := makeCounter()
	c2 := makeCounter()
	a := c1()
	b := c1()
	c := c2()
	return a + b + c
}

// OverLoop tests closure over loop
func OverLoop() int {
	sum := 0
	adder := func(x int) int {
		sum = sum + x
		return sum
	}
	for i := 1; i <= 5; i++ {
		adder(i)
	}
	return sum
}

// Chain tests closure chain
func Chain() int {
	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }
	return addOne(double(5))
}

// Accumulator tests accumulator pattern
func Accumulator() int {
	acc := makeAccumulator(100)
	acc(10)
	acc(20)
	return acc(30)
}

func makeAccumulator(init int) func(int) int {
	total := init
	return func(n int) int {
		total = total + n
		return total
	}
}
