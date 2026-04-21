package divergence_hunt193

import (
	"fmt"
)

// ============================================================================
// Round 193: Bitwise operations on signed integers
// ============================================================================

// BitwiseAndSigned tests bitwise AND on signed integers
func BitwiseAndSigned() string {
	a := int8(-1)  // 0xFF
	b := int8(0x0F)
	result := a & b
	return fmt.Sprintf("%d", result)
}

// BitwiseOrSigned tests bitwise OR on signed integers
func BitwiseOrSigned() string {
	a := int8(-128) // 0x80
	b := int8(0x7F) // 0x7F
	result := a | b
	return fmt.Sprintf("%d", result)
}

// BitwiseXorSigned tests bitwise XOR on signed integers
func BitwiseXorSigned() string {
	a := int16(-1)
	b := int16(-1)
	result := a ^ b
	return fmt.Sprintf("%d", result)
}

// BitwiseNotSigned tests bitwise NOT on signed integers
func BitwiseNotSigned() string {
	a := int8(0)
	result := ^a
	return fmt.Sprintf("%d", result)
}

// SignBitIsolation tests isolating sign bit
func SignBitIsolation() string {
	a := int32(-5)
	signBit := uint32(a) >> 31
	return fmt.Sprintf("%d", signBit)
}

// SignExtension tests sign extension
func SignExtension() string {
	a := int8(-5)
	extended := int16(a)
	return fmt.Sprintf("%d", extended)
}

// AbsoluteValueBitwise tests absolute value using bitwise trick
func AbsoluteValueBitwise() string {
	n := int8(-42)
	mask := n >> 7
	abs := (n + mask) ^ mask
	return fmt.Sprintf("%d", abs)
}

// MinMaxBitwise tests min/max using bitwise operations
func MinMaxBitwise() string {
	a := int8(10)
	b := int8(20)
	// Min
	min := b + ((a - b) & ((a - b) >> 7))
	return fmt.Sprintf("%d", min)
}

// SwapXor tests swapping using XOR
func SwapXor() string {
	a := int16(100)
	b := int16(200)
	a = a ^ b
	b = a ^ b
	a = a ^ b
	return fmt.Sprintf("%d:%d", a, b)
}

// ClearLowestBit tests clearing lowest set bit
func ClearLowestBit() string {
	n := int16(12) // 1100
	result := n & (n - 1)
	return fmt.Sprintf("%d", result)
}

// IsolateLowestBit tests isolating lowest set bit
func IsolateLowestBit() string {
	n := int16(12) // 1100
	result := n & (-n)
	return fmt.Sprintf("%d", result)
}
