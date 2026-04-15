package divergence_hunt14

import (
	"fmt"
	"math"
)

// ============================================================================
// Round 14: Float edge cases, math edge cases, numeric precision,
// division edge cases, comparison edge cases
// ============================================================================

// FloatAddPrecision tests float addition precision
func FloatAddPrecision() float64 {
	return 0.1 + 0.2
}

// FloatMultiplyPrecision tests float multiplication precision
func FloatMultiplyPrecision() float64 {
	return 0.1 * 0.2
}

// FloatDivPrecision tests float division precision
func FloatDivPrecision() float64 {
	return 1.0 / 3.0
}

// FloatNegative tests negative float
func FloatNegative() float64 {
	return -3.14
}

// FloatZeroDivision tests float division by zero (produces Inf, not panic)
func FloatZeroDivision() bool {
	x := 1.0
	y := 0.0
	return math.IsInf(x/y, 1)
}

// FloatNaNArithmetic tests NaN arithmetic
func FloatNaNArithmetic() bool {
	nan := math.NaN()
	return nan != nan && nan+1 != nan+1
}

// FloatInfArithmetic tests infinity arithmetic
func FloatInfArithmetic() bool {
	inf := math.Inf(1)
	return inf+1 == inf && math.IsInf(-inf, -1)
}

// FloatComparisonPrecision tests float comparison
func FloatComparisonPrecision() bool {
	a := 0.1 + 0.2
	b := 0.3
	return a != b // floating point: 0.1+0.2 != 0.3
}

// IntDivisionTruncation tests integer division truncation
func IntDivisionTruncation() int {
	return 7 / 3
}

// IntModulo tests integer modulo
func IntModulo() int {
	return 7 % 3
}

// NegativeModulo tests negative modulo
func NegativeModulo() int {
	return -7 % 3
}

// Float32NaN tests float32 NaN
func Float32NaN() bool {
	return float32(math.NaN()) != float32(math.NaN())
}

// Float32Inf tests float32 infinity
func Float32Inf() bool {
	return float32(math.Inf(1)) > float32(1e38)
}

// MathSin tests math.Sin
func MathSin() float64 { return math.Sin(math.Pi / 2) }

// MathCos tests math.Cos
func MathCos() float64 { return math.Cos(0) }

// MathTan tests math.Tan
func MathTan() float64 { return math.Tan(0) }

// MathAtan2 tests math.Atan2
func MathAtan2() float64 { return math.Atan2(1, 1) }

// MathLog2 tests math.Log2
func MathLog2() float64 { return math.Log2(1024) }

// MathLog10 tests math.Log10
func MathLog10() float64 { return math.Log10(1000) }

// FmtFloatFormat tests various float formats
func FmtFloatFormat() string {
	return fmt.Sprintf("%.2f|%e|%g", math.Pi, math.Pi, math.Pi)
}

// FmtIntFormat tests various int formats
func FmtIntFormat() string {
	return fmt.Sprintf("%d|%x|%#x|%o|%#o", 255, 255, 255, 255, 255)
}

// FloatMaxMin tests float max/min edge cases
func FloatMaxMin() float64 {
	return math.Max(1.0, 2.0) + math.Min(1.0, 2.0)
}

// Float32Limits tests float32 limits
func Float32Limits() float32 {
	max := math.MaxFloat32
	return float32(max)
}

// ComplexMagnitude tests complex number magnitude
func ComplexMagnitude() float64 {
	z := complex(3, 4)
	return math.Sqrt(real(z)*real(z) + imag(z)*imag(z))
}
