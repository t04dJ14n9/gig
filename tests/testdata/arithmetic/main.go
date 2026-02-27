package arithmetic

// Addition tests basic addition
func Addition() int { return 2 + 3 }

// Subtraction tests basic subtraction
func Subtraction() int { return 10 - 4 }

// Multiplication tests basic multiplication
func Multiplication() int { return 6 * 7 }

// Division tests basic division
func Division() int { return 20 / 4 }

// Modulo tests modulo operation
func Modulo() int { return 17 % 5 }

// ComplexExpr tests complex arithmetic expression
func ComplexExpr() int { return (2+3)*4 - 10/2 }

// Negation tests unary negation
func Negation() int {
	x := 42
	return -x
}

// ChainedOps tests chained arithmetic operations
func ChainedOps() int {
	a := 10
	b := a * 2
	c := b + a
	d := c - 5
	return d / 5
}

// Overflow tests int64 wrapping behavior
func Overflow() int {
	x := 9223372036854775807
	return x + 1
}

// Precedence tests operator precedence
func Precedence() int { return 2 + 3*4 }
