package divergence_hunt191

import (
	"fmt"
	"math"
)

// ============================================================================
// Round 191: Floating point special values (NaN, Inf)
// ============================================================================

// NaNEquality tests that NaN != NaN
func NaNEquality() string {
	nan := math.NaN()
	result := nan != nan
	return fmt.Sprintf("%v", result)
}

// NaNIsNaN tests math.IsNaN function
func NaNIsNaN() string {
	nan := math.NaN()
	regular := 3.14
	return fmt.Sprintf("%v:%v", math.IsNaN(nan), math.IsNaN(regular))
}

// PositiveInf tests positive infinity
func PositiveInf() string {
	inf := math.Inf(1)
	return fmt.Sprintf("%v:%v", inf > 1e308, math.IsInf(inf, 1))
}

// NegativeInf tests negative infinity
func NegativeInf() string {
	inf := math.Inf(-1)
	return fmt.Sprintf("%v:%v", inf < -1e308, math.IsInf(inf, -1))
}

// InfArithmetic tests arithmetic with infinity
func InfArithmetic() string {
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	sum := posInf + negInf
	product := posInf * 2
	return fmt.Sprintf("%v:%v", math.IsNaN(sum), product == posInf)
}

// InfComparison tests infinity comparisons
func InfComparison() string {
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	return fmt.Sprintf("%v:%v:%v", negInf < 0, 0 < posInf, negInf < posInf)
}

// ZeroSign tests positive and negative zero
func ZeroSign() string {
	posZero := 0.0
	negZero := math.Copysign(0, -1)
	return fmt.Sprintf("%v:%v", posZero == negZero, 1/posZero == 1/negZero)
}

// NaNPropagation tests NaN propagation in operations
func NaNPropagation() string {
	nan := math.NaN()
	a := nan + 5
	b := nan * 0
	c := nan / 1
	return fmt.Sprintf("%v:%v:%v", math.IsNaN(a), math.IsNaN(b), math.IsNaN(c))
}

// InfDivision tests division by infinity
func InfDivision() string {
	posInf := math.Inf(1)
	a := 1.0 / posInf
	b := -1.0 / posInf
	return fmt.Sprintf("%v:%v", a == 0, b < 0)
}

// FiniteCheck tests math.IsInf with sign 0
func FiniteCheck() string {
	regular := 3.14
	inf := math.Inf(1)
	return fmt.Sprintf("%v:%v", !math.IsInf(regular, 0), math.IsInf(inf, 0))
}
