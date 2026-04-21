package divergence_hunt194

import (
	"fmt"
)

// ============================================================================
// Round 194: Shift operations with variable shifts
// ============================================================================

// VariableLeftShift tests left shift with variable count
func VariableLeftShift() string {
	x := uint8(1)
	result := ""
	for i := 0; i < 8; i++ {
		result += fmt.Sprintf("%d:", x<<i)
	}
	return result
}

// VariableRightShift tests right shift with variable count
func VariableRightShift() string {
	x := uint8(128)
	result := ""
	for i := 0; i < 8; i++ {
		result += fmt.Sprintf("%d:", x>>i)
	}
	return result
}

// SignedRightShift tests arithmetic right shift
func SignedRightShift() string {
	x := int8(-128)
	result := ""
	for i := 0; i < 4; i++ {
		result += fmt.Sprintf("%d:", x>>i)
	}
	return result
}

// ShiftByZero tests shift by zero
func ShiftByZero() string {
	x := uint16(0xABCD)
	left := x << 0
	right := x >> 0
	return fmt.Sprintf("%d:%d", left, right)
}

// LargeShiftCount tests shift counts larger than width
func LargeShiftCount() string {
	x := uint8(255)
	shifted := x << 8
	return fmt.Sprintf("%d", shifted)
}

// ShiftAssignment tests shift assignment operators
func ShiftAssignment() string {
	x := uint16(1)
	x <<= 4
	x >>= 2
	return fmt.Sprintf("%d", x)
}

// ShiftInExpression tests shift within complex expression
func ShiftInExpression() string {
	x := uint8(1)
	y := uint8(2)
	result := (x << 3) | (y << 1)
	return fmt.Sprintf("%d", result)
}

// ShiftThenMask tests shift followed by mask
func ShiftThenMask() string {
	x := uint16(0xABCD)
	lowByte := (x >> 0) & 0xFF
	highByte := (x >> 8) & 0xFF
	return fmt.Sprintf("%d:%d", highByte, lowByte)
}

// ShiftBounds tests shifting near type bounds
func ShiftBounds() string {
	a := uint8(255) >> 1
	b := int8(-1) >> 1
	return fmt.Sprintf("%d:%d", a, b)
}

// CircularShiftPattern tests pattern for circular shift
func CircularShiftPattern() string {
	x := uint8(0b10110011)
	left := 3
	right := 8 - left
	result := (x << left) | (x >> right)
	return fmt.Sprintf("%d", result)
}
