package divergence_hunt4

import (
	"math"
	"strconv"
)

// ============================================================================
// Round 4: Numeric edge cases, type conversions, math functions,
// string conversions, edge cases with special values
// ============================================================================

// Float64NaN tests NaN propagation
func Float64NaN() bool {
	return math.IsNaN(math.NaN())
}

// Float64Inf tests infinity
func Float64Inf() bool {
	return math.IsInf(math.Inf(1), 1) && math.IsInf(math.Inf(-1), -1)
}

// Float64NegZero tests negative zero
func Float64NegZero() bool {
	return math.Signbit(-0.0) && !math.Signbit(0.0)
}

// Int16Conversion tests int16 conversion
func Int16Conversion() int16 {
	var x int32 = 1000
	return int16(x)
}

// Uint32Conversion tests uint32 conversion
func Uint32Conversion() uint32 {
	var x int64 = 50000
	return uint32(x)
}

// FloatToIntTruncation tests float to int truncation
func FloatToIntTruncation() int {
	x := 9.99
	return int(x)
}

// NegativeFloatToInt tests negative float to int
func NegativeFloatToInt() int {
	x := -3.7
	return int(x)
}

// StrconvAtoi tests strconv.Atoi
func StrconvAtoi() int {
	n, _ := strconv.Atoi("12345")
	return n
}

// StrconvItoa tests strconv.Itoa
func StrconvItoa() string {
	return strconv.Itoa(42)
}

// StrconvFormatInt tests strconv.FormatInt
func StrconvFormatInt() string {
	return strconv.FormatInt(-42, 10)
}

// StrconvParseFloat tests strconv.ParseFloat
func StrconvParseFloat() float64 {
	f, _ := strconv.ParseFloat("3.14", 64)
	return f
}

// MathAbs tests math.Abs
func MathAbs() float64 { return math.Abs(-42.5) }

// MathMax tests math.Max
func MathMax() float64 { return math.Max(3.14, 2.71) }

// MathMin tests math.Min
func MathMin() float64 { return math.Min(3.14, 2.71) }

// MathPow tests math.Pow
func MathPow() float64 { return math.Pow(2, 10) }

// MathSqrt tests math.Sqrt
func MathSqrt() float64 { return math.Sqrt(144) }

// MathCeil tests math.Ceil
func MathCeil() float64 { return math.Ceil(3.14) }

// MathFloor tests math.Floor
func MathFloor() float64 { return math.Floor(3.14) }

// IntMin tests integer minimum
func IntMin() int {
	a, b := 3, 7
	if a < b { return a }
	return b
}

// IntMax tests integer maximum
func IntMax() int {
	a, b := 3, 7
	if a > b { return a }
	return b
}

// UintptrSize tests uintptr
func UintptrSize() int {
	var p uintptr = 42
	return int(p)
}

// ByteArith tests byte arithmetic
func ByteArith() byte {
	var a byte = 200
	var b byte = 55
	return a + b
}

// Int32Overflow tests int32 overflow
func Int32Overflow() int32 {
	var x int32 = 2147483647
	x += 1
	return x
}

// Uint8Wrap tests uint8 wrapping
func Uint8Wrap() uint8 {
	var x uint8 = 0
	x -= 1
	return x
}

// ComplexConj tests complex conjugate
func ComplexConj() complex128 {
	z := complex(3, 4)
	return complex(real(z), -imag(z))
}

// Float32Precision tests float32 precision
func Float32Precision() float32 {
	var a float32 = 0.1
	var b float32 = 0.2
	return a + b
}

// MapLenAfterDelete tests map len after delete
func MapLenAfterDelete() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	delete(m, "nonexistent")
	return len(m)
}

// SliceCapAfterAppend tests slice capacity after append
func SliceCapAfterAppend() int {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3) // triggers grow
	return cap(s)
}

// StringFromRunes tests building string from runes
func StringFromRunes() string {
	runes := []rune{'H', 'e', 'l', 'l', 'o'}
	return string(runes)
}

// RuneToInt tests rune to int
func RuneToInt() int {
	var r rune = 'A'
	return int(r)
}

// BoolToInt tests bool conversion
func BoolToInt() int {
	v := true
	if v { return 1 }
	return 0
}
