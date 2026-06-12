package divergence_hunt278

import (
	"fmt"
)

// ============================================================================
// Round 278: Integer overflow, shift edge cases, unsigned arithmetic

// Int8Overflow tests int8 overflow wraps around
func Int8Overflow() string {
	var x int8 = 127
	x++
	return fmt.Sprintf("x=%d", x)
}

// Uint8Overflow tests uint8 overflow wraps around
func Uint8Overflow() string {
	var x uint8 = 255
	x++
	return fmt.Sprintf("x=%d", x)
}

// IntAddOverflow tests int addition that would overflow 32-bit
func IntAddOverflow() string {
	var x int32 = 2147483647
	x++
	return fmt.Sprintf("x=%d", x)
}

// ShiftByZero tests shifting by 0
func ShiftByZero() string {
	x := 42
	return fmt.Sprintf("x<<0=%d,x>>0=%d", x<<0, x>>0)
}

// ShiftByLargeAmount tests shifting by large amount
func ShiftByLargeAmount() string {
	x := 1
	return fmt.Sprintf("1<<31=%d", x<<31)
}

// UnsignedShiftRight tests unsigned right shift
func UnsignedShiftRight() string {
	var x uint8 = 128
	return fmt.Sprintf("128>>1=%d", x>>1)
}

// SignedShiftRight tests signed right shift (sign extends)
func SignedShiftRight() string {
	var x int8 = -64
	return fmt.Sprintf("neg>>1=%d", x>>1)
}

// NegateMinInt tests negating minimum int (stays negative due to overflow)
func NegateMinInt() string {
	var x int8 = -128
	y := -x
	return fmt.Sprintf("negate=%d", y)
}

// MultiplyOverflow tests multiplication overflow
func MultiplyOverflow() string {
	var x int8 = 50
	var y int8 = 3
	z := x * y
	return fmt.Sprintf("z=%d", z)
}

// UnsignedSubtraction tests unsigned subtraction underflow
func UnsignedSubtraction() string {
	var x uint8 = 5
	var y uint8 = 10
	z := x - y
	return fmt.Sprintf("z=%d", z)
}

// BitwiseNot tests bitwise NOT
func BitwiseNot() string {
	var x uint8 = 0x0F
	return fmt.Sprintf("^0x0F=0x%02X", ^x)
}

// ComplexBitwise tests complex bitwise expression
func ComplexBitwise() string {
	a := 0x12
	b := 0x0F
	return fmt.Sprintf("0x%02X", a&b|(^a&^b))
}

// ShiftAssignment tests shift assignment operators
func ShiftAssignment() string {
	x := 1
	x <<= 4
	x >>= 2
	return fmt.Sprintf("x=%d", x)
}

// FloatToIntOverflow tests float to int conversion overflow
func FloatToIntOverflow() string {
	f := 1e20
	i := int(f)
	return fmt.Sprintf("i=%d", i)
}
