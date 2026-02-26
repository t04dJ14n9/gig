package controlflow

// IfTrue tests if with true condition
func IfTrue() int {
	x := 10
	if x > 5 {
		return 1
	}
	return 0
}

// IfFalse tests if with false condition
func IfFalse() int {
	x := 3
	if x > 5 {
		return 1
	}
	return 0
}

// IfElse tests if-else
func IfElse() int {
	x := 3
	if x > 5 {
		return 1
	} else {
		return -1
	}
}

// IfElseChainNegative tests if-else chain with negative
func IfElseChainNegative() int { return classify(-5) }

// IfElseChainZero tests if-else chain with zero
func IfElseChainZero() int { return classify(0) }

// IfElseChainPositive tests if-else chain with positive
func IfElseChainPositive() int { return classify(42) }

func classify(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	} else {
		return 1
	}
}

// ForLoop tests basic for loop
func ForLoop() int {
	sum := 0
	for i := 1; i <= 10; i++ {
		sum = sum + i
	}
	return sum
}

// ForConditionOnly tests for with condition only
func ForConditionOnly() int {
	i := 0
	sum := 0
	for i < 5 {
		sum = sum + i
		i = i + 1
	}
	return sum
}

// NestedFor tests nested for loops
func NestedFor() int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			sum = sum + 1
		}
	}
	return sum
}

// ForBreak tests for loop with break
func ForBreak() int {
	sum := 0
	for i := 0; i < 100; i++ {
		if i >= 5 {
			break
		}
		sum = sum + i
	}
	return sum
}

// ForContinue tests for loop with continue
func ForContinue() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		sum = sum + i
	}
	return sum
}

// BooleanAndOr tests boolean operators
func BooleanAndOr() int {
	a := true
	b := false
	result := 0
	if a && !b {
		result = result + 1
	}
	if a || b {
		result = result + 10
	}
	if !b {
		result = result + 100
	}
	return result
}
