package divergence_hunt114

import (
	"fmt"
	"math"
)

// ============================================================================
// Round 114: Math and numeric edge cases
// ============================================================================

func MathAbs() string {
	return fmt.Sprintf("%.1f", math.Abs(-3.14))
}

func MathMax() string {
	return fmt.Sprintf("%.1f", math.Max(1.5, 2.5))
}

func MathMin() string {
	return fmt.Sprintf("%.1f", math.Min(1.5, 2.5))
}

func MathCeil() string {
	return fmt.Sprintf("%.0f", math.Ceil(3.14))
}

func MathFloor() string {
	return fmt.Sprintf("%.0f", math.Floor(3.14))
}

func MathRound() string {
	return fmt.Sprintf("%.0f", math.Round(3.5))
}

func MathPow() string {
	return fmt.Sprintf("%.0f", math.Pow(2, 10))
}

func MathSqrt() string {
	return fmt.Sprintf("%.1f", math.Sqrt(16))
}

func IntOverflow() string {
	var x int8 = 127
	x += 1
	return fmt.Sprintf("%d", x)
}

func FloatPrecision() string {
	a := 0.1 + 0.2
	return fmt.Sprintf("%.1f", a)
}

func IntegerDivision() string {
	a := 7 / 2
	b := 7 % 2
	return fmt.Sprintf("%d:%d", a, b)
}

func UintRange() string {
	var x uint8 = 255
	return fmt.Sprintf("%d", x)
}

func NegativeModulo() string {
	a := -7 % 2
	return fmt.Sprintf("%d", a)
}

func FloatToIntTruncation() string {
	x := int(3)
	return fmt.Sprintf("%d", x)
}
