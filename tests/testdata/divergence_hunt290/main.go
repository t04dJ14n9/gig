package divergence_hunt290

import (
	"fmt"
)

// ============================================================================
// Round 290: Complex expressions — short-circuit evaluation, ternary-like, compound assignment

// ShortCircuitAnd tests && short-circuit: second operand not evaluated
func ShortCircuitAnd() string {
	called := false
	sideEffect := func() bool {
		called = true
		return true
	}
	_ = false && sideEffect()
	return fmt.Sprintf("called=%t", called)
}

// ShortCircuitOr tests || short-circuit: second operand not evaluated
func ShortCircuitOr() string {
	called := false
	sideEffect := func() bool {
		called = true
		return true
	}
	_ = true || sideEffect()
	return fmt.Sprintf("called=%t", called)
}

// TernaryLike tests ternary-like expression using if
func TernaryLike() string {
	x := 5
	result := ""
	if x > 3 {
		result = "big"
	} else {
		result = "small"
	}
	return result
}

// CompoundAssignmentAdd tests += operator
func CompoundAssignmentAdd() string {
	x := 10
	x += 5
	return fmt.Sprintf("x=%d", x)
}

// CompoundAssignmentMul tests *= operator
func CompoundAssignmentMul() string {
	x := 3
	x *= 7
	return fmt.Sprintf("x=%d", x)
}

// CompoundAssignmentOnSliceIndex tests compound assignment on slice element
func CompoundAssignmentOnSliceIndex() string {
	s := []int{1, 2, 3}
	s[1] += 10
	return fmt.Sprintf("s=%v", s)
}

// CompoundAssignmentOnMapValue tests compound assignment on map value
func CompoundAssignmentOnMapValue() string {
	m := map[string]int{"a": 5}
	m["a"] *= 3
	return fmt.Sprintf("a=%d", m["a"])
}

// IncrementDecrement tests ++ and -- operators
func IncrementDecrement() string {
	x := 10
	x++
	x++
	x--
	return fmt.Sprintf("x=%d", x)
}

// ComplexArithmeticExpression tests operator precedence
func ComplexArithmeticExpression() string {
	a := 2 + 3*4     // 14, not 20
	b := (2 + 3) * 4 // 20
	c := 10 - 3 - 2  // 5, left-to-right
	d := 8 / 4 * 2   // 4, left-to-right
	return fmt.Sprintf("a=%d,b=%d,c=%d,d=%d", a, b, c, d)
}

// BooleanExpression tests boolean expression with precedence
func BooleanExpression() string {
	a := true
	b := false
	c := true
	result := a && b || c // (true && false) || true = true
	return fmt.Sprintf("result=%t", result)
}

// NegationOperator tests ! operator
func NegationOperator() string {
	a := true
	b := !a
	return fmt.Sprintf("a=%t,b=%t", a, b)
}

// UnaryMinus tests unary minus
func UnaryMinus() string {
	x := 5
	y := -x
	return fmt.Sprintf("y=%d", y)
}

// StringComparisonInIf tests string comparison in if
func StringComparisonInIf() string {
	s := "hello"
	if s == "hello" {
		return "match"
	}
	return "no_match"
}
