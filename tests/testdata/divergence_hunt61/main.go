package divergence_hunt61

// ============================================================================
// Round 61: Integer overflow and wrapping edge cases
// ============================================================================

func Uint8Overflow() uint8 {
	var x uint8 = 255
	x += 1
	return x
}

func Uint8Underflow() uint8 {
	var x uint8 = 0
	x -= 1
	return x
}

func Int8Overflow() int8 {
	var x int8 = 127
	x += 1
	return x
}

func Int8Underflow() int8 {
	var x int8 = -128
	x -= 1
	return x
}

func Uint16Overflow() uint16 {
	var x uint16 = 65535
	x += 1
	return x
}

func Uint32Overflow() uint32 {
	var x uint32 = 4294967295
	x += 1
	return x
}

func Int16Overflow() int16 {
	var x int16 = 32767
	x += 1
	return x
}

func Int16Underflow() int16 {
	var x int16 = -32768
	x -= 1
	return x
}

func IntNegateMin() int8 {
	var x int8 = -128
	return -x // remains -128 due to overflow
}

func UintMulOverflow() uint8 {
	var x uint8 = 200
	var y uint8 = 2
	return x * y
}

func IntDivTruncation() int {
	return -7 / 2
}

func IntModNegative() int {
	return -7 % 2
}

func ShiftLeftLarge() uint8 {
	var x uint8 = 1
	return x << 7
}

func ShiftRightSigned() int8 {
	var x int8 = -128
	return x >> 3
}

func UintConvertNegative() uint8 {
	var x int8 = -1
	return uint8(x)
}

func IntConvertLargeUint() int8 {
	var x uint8 = 200
	return int8(x)
}

func FloatTruncateToInt() int {
	x := 3.9
	return int(x)
}

func FloatTruncateNegToInt() int {
	x := -3.9
	return int(x)
}

func ComplexRealImag() (float64, float64) {
	z := complex(3.0, 4.0)
	return real(z), imag(z)
}
