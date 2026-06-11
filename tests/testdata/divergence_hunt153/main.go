package divergence_hunt153

import (
	"fmt"
	"math"
	"math/bits"
)

// ============================================================================
// Round 153: Math and bits operations
// ============================================================================

// BitsLeadingZeros tests bits.LeadingZeros
func BitsLeadingZeros() string {
	return fmt.Sprintf("16=%d-32=%d", bits.LeadingZeros16(0x00FF), bits.LeadingZeros32(0x0000FFFF))
}

// BitsTrailingZeros tests bits.TrailingZeros
func BitsTrailingZeros() string {
	return fmt.Sprintf("16=%d-32=%d", bits.TrailingZeros16(0xFF00), bits.TrailingZeros32(0xFFFF0000))
}

// BitsOnesCount tests bits.OnesCount
func BitsOnesCount() string {
	return fmt.Sprintf("8=%d-16=%d", bits.OnesCount8(0b10101010), bits.OnesCount16(0xFFFF))
}

// BitsRotate tests bits.RotateLeft
func BitsRotate() string {
	return fmt.Sprintf("16=%d-32=%d", bits.RotateLeft16(0x8001, 1), bits.RotateLeft32(0x80000001, 1))
}

// BitsReverse tests bits.Reverse
func BitsReverse() string {
	return fmt.Sprintf("8=%d-16=%d", bits.Reverse8(0b11001100), bits.Reverse16(0xFF00))
}

// BitsLen tests bits.Len
func BitsLen() string {
	return fmt.Sprintf("8=%d-16=%d-32=%d", bits.Len8(0b1000), bits.Len16(0x0100), bits.Len32(0x00010000))
}

// MathAbs tests math.Abs
func MathAbs() string {
	return fmt.Sprintf("pos=%.0f-neg=%.0f", math.Abs(5.5), math.Abs(-5.5))
}

// MathMinMax tests math.Min and math.Max
func MathMinMax() string {
	return fmt.Sprintf("min=%.0f-max=%.0f", math.Min(3.0, 7.0), math.Max(3.0, 7.0))
}

// MathFloorCeil tests math.Floor and math.Ceil
func MathFloorCeil() string {
	return fmt.Sprintf("floor=%.0f-ceil=%.0f", math.Floor(3.7), math.Ceil(3.2))
}

// MathRound tests math.Round
func MathRound() string {
	return fmt.Sprintf("down=%.0f-up=%.0f", math.Round(3.4), math.Round(3.5))
}

// MathTrunc tests math.Trunc
func MathTrunc() string {
	return fmt.Sprintf("pos=%.0f-neg=%.0f", math.Trunc(3.7), math.Trunc(-3.7))
}

// MathMod tests math.Mod
func MathMod() string {
	return fmt.Sprintf("mod=%.1f", math.Mod(17.5, 4.0))
}

// MathPow tests math.Pow
func MathPow() string {
	return fmt.Sprintf("pow=%.0f", math.Pow(2, 10))
}

// MathSqrt tests math.Sqrt
func MathSqrt() string {
	return fmt.Sprintf("sqrt=%.0f", math.Sqrt(256))
}

// MathNaNCheck tests NaN detection
func MathNaNCheck() string {
	return fmt.Sprintf("isnan=%t", math.IsNaN(math.NaN()))
}

// MathInfCheck tests Inf detection
func MathInfCheck() string {
	return fmt.Sprintf("posinf=%t-neginf=%t", math.IsInf(math.Inf(1), 1), math.IsInf(math.Inf(-1), -1))
}

// MathCopySign tests math.Copysign
func MathCopySign() string {
	return fmt.Sprintf("sign=%.0f", math.Copysign(5.0, -1.0))
}

// MathPiE tests math constants
func MathPiE() string {
	return fmt.Sprintf("pi=%.2f-e=%.2f", math.Pi, math.E)
}
