package divergence_hunt33

// ============================================================================
// Round 33: Integer arithmetic edge cases - int8/int16/int32/uint16/uint32
// ============================================================================

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

func Int16Overflow() int16 {
	var x int16 = 32767
	x += 1
	return x
}

func Uint16Overflow() uint16 {
	var x uint16 = 65535
	x += 1
	return x
}

func Uint32Arith() uint32 {
	var x uint32 = 4000000000
	var y uint32 = 1000000000
	return x - y
}

func Int32Arith() int32 {
	var x int32 = 1000000000
	var y int32 = 500000000
	return x + y
}

func ShiftLeft8() uint8 {
	var x uint8 = 1
	return x << 7
}

func ShiftRight8() uint8 {
	var x uint8 = 128
	return x >> 1
}

func NegateInt8() int8 {
	var x int8 = 10
	return -x
}

func NegateInt16() int16 {
	var x int16 = -100
	return -x
}

func MixedIntArith() int {
	var a int8 = 10
	var b int16 = 20
	var c int32 = 30
	return int(a) + int(b) + int(c)
}

func IntDivTruncation() int {
	return 7 / 2
}

func IntModNegative() int {
	return -7 % 3
}

func UintDivTruncation() uint {
	var x uint = 7
	var y uint = 2
	return x / y
}

func UintMod() uint {
	var x uint = 7
	var y uint = 2
	return x % y
}

func BitwiseAndNot() int {
	return 0xFF &^ 0x0F // 0xF0 = 240
}

func BitwiseXor() int {
	return 0xFF ^ 0x0F // 0xF0 = 240
}

func BitwiseOr() int {
	return 0xF0 | 0x0F // 0xFF = 255
}

func ComplexShift() int {
	x := 1
	return x << 20
}

func ShiftWithUintAmount() uint {
	var s uint = 4
	var x uint = 1
	return x << s
}
