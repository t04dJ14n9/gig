package divergence_hunt195

import (
	"fmt"
)

// ============================================================================
// Round 195: Overflow/underflow edge cases
// ============================================================================

// Uint8Overflow tests uint8 overflow
func Uint8Overflow() string {
	x := uint8(255)
	x = x + 1
	return fmt.Sprintf("%d", x)
}

// Uint16Overflow tests uint16 overflow
func Uint16Overflow() string {
	x := uint16(65535)
	x = x + 1
	return fmt.Sprintf("%d", x)
}

// Int8PositiveOverflow tests int8 positive overflow
func Int8PositiveOverflow() string {
	x := int8(127)
	x = x + 1
	return fmt.Sprintf("%d", x)
}

// Int8NegativeOverflow tests int8 negative overflow
func Int8NegativeOverflow() string {
	x := int8(-128)
	x = x - 1
	return fmt.Sprintf("%d", x)
}

// Int16Overflow tests int16 overflow
func Int16Overflow() string {
	x := int16(32767)
	x = x + 1
	return fmt.Sprintf("%d", x)
}

// MultiplicationOverflow tests multiplication overflow
func MultiplicationOverflow() string {
	x := uint8(100)
	y := uint8(3)
	result := x * y
	return fmt.Sprintf("%d", result)
}

// CompoundOverflow tests compound assignment overflow
func CompoundOverflow() string {
	x := uint8(200)
	x += 100
	return fmt.Sprintf("%d", x)
}

// UnderflowSubtraction tests subtraction underflow
func UnderflowSubtraction() string {
	x := uint8(0)
	y := uint8(1)
	result := x - y
	return fmt.Sprintf("%d", result)
}

// IncrementOverflow tests increment overflow
func IncrementOverflow() string {
	x := uint8(255)
	x++
	return fmt.Sprintf("%d", x)
}

// DecrementUnderflow tests decrement underflow
func DecrementUnderflow() string {
	x := uint8(0)
	x--
	return fmt.Sprintf("%d", x)
}

// AdditiveChainOverflow tests chained addition overflow
func AdditiveChainOverflow() string {
	x := uint8(100)
	x = x + 100 + 100
	return fmt.Sprintf("%d", x)
}
