package namedreturn

// Basic tests named return basic
func Basic() int { return double(21) }

func double(x int) (result int) {
	result = x * 2
	return result
}

// Multiple tests multiple named returns
func Multiple() int {
	q, r := divmod(17, 5)
	return q*10 + r
}

func divmod(a, b int) (quotient int, remainder int) {
	quotient = a / b
	remainder = a % b
	return quotient, remainder
}

// ZeroValue tests named return zero value
func ZeroValue() int {
	return maybeDouble(10, 1) + maybeDouble(10, 0)
}

func maybeDouble(x int, doIt int) (result int) {
	if doIt > 0 {
		result = x * 2
	}
	return result
}
