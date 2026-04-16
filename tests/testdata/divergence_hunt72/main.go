package divergence_hunt72

import "math"

// ============================================================================
// Round 72: Floating point edge cases - NaN, Inf, -0, precision
// ============================================================================

func Float64NaN() bool {
	return math.IsNaN(math.NaN())
}

func Float64Inf() bool {
	return math.IsInf(math.Inf(1), 1)
}

func Float64NegInf() bool {
	return math.IsInf(math.Inf(-1), -1)
}

func Float64Zero() bool {
	var z float64
	return z == 0
}

func Float64NegZero() bool {
	negZero := -0.0
	return math.Signbit(negZero)
}

func Float64NaNNotEqual() bool {
	nan := math.NaN()
	return nan != nan
}

func Float64InfArith() float64 {
	inf := math.Inf(1)
	return inf + 1
}

func Float64InfSubInf() bool {
	result := math.Inf(1) - math.Inf(1)
	return math.IsNaN(result)
}

func Float64ZeroDiv() float64 {
	x := 1.0
	y := 0.0
	return x / y
}

func Float32Precision() float32 {
	var x float32 = 0.1
	var y float32 = 0.2
	return x + y
}

func Float64Truncation() int {
	x := 3.7
	return int(x)
}

func Float64NegativeTruncation() int {
	x := -3.7
	return int(x)
}

func Float64Mod() float64 {
	return math.Mod(7.5, 2.5)
}

func Float64Pow() float64 {
	return math.Pow(2, 10)
}

func Float64Sqrt() float64 {
	return math.Sqrt(144)
}

func Float64Abs() float64 {
	return math.Abs(-42.5)
}

func Float64Max() float64 {
	return math.Max(3.14, 2.71)
}

func Float64Min() float64 {
	return math.Min(3.14, 2.71)
}
