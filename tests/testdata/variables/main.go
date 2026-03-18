package variables

// DeclareAndUse tests variable declaration and use
func DeclareAndUse() int {
	x := 10
	y := 20
	z := x + y
	return z
}

// Reassignment tests variable reassignment
func Reassignment() int {
	x := 1
	x = x + 10
	x = x * 2
	return x
}

// MultipleDecl tests multiple variable declarations
func MultipleDecl() int {
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	return a + b + c + d + e
}

// ZeroValues tests zero values
func ZeroValues() int {
	var x int
	return x
}

// StringZeroValue tests string zero value
func StringZeroValue() string {
	var s string
	return s
}

// Shadowing tests variable shadowing
func Shadowing() int {
	x := 10
	y := 1
	if y > 0 {
		x := 20
		_ = x
	}
	return x
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// SumThree returns a + b + c
func SumThree(a, b, c int) int { return a + b + c }

// Multiply returns a * b
func Multiply(a, b int) int { return a * b }

// Max returns the maximum of a and b
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IsPositive returns true if x > 0
func IsPositive(x int) bool { return x > 0 }
