package divergence_hunt87

// ============================================================================
// Round 87: Bit manipulation edge cases - shifts, masks, XOR
// ============================================================================

func BitwiseAnd() int {
	return 0xFF & 0x0F
}

func BitwiseOr() int {
	return 0xF0 | 0x0F
}

func BitwiseXor() int {
	return 0xFF ^ 0x0F
}

func BitwiseNot() int {
	return ^0x0F & 0xFF
}

func BitwiseShiftLeft() int {
	return 1 << 10
}

func BitwiseShiftRight() int {
	return 1024 >> 5
}

func BitwiseShiftLeftOverflow() uint8 {
	var x uint8 = 1
	return x << 8 // shifts all bits out
}

func BitwiseAndNot() int {
	return 0xFF &^ 0x0F
}

func BitMask() int {
	return (1<<4 - 1) // 0b1111 = 15
}

func BitSet() int {
	x := 0
	x |= 1 << 3 // set bit 3
	return x
}

func BitClear() int {
	x := 0xFF
	x &^= 1 << 3 // clear bit 3
	return x
}

func BitToggle() int {
	x := 0
	x ^= 1 << 3 // toggle bit 3
	x ^= 1 << 3 // toggle back
	return x
}

func BitCheck() bool {
	x := 8 // 0b1000
	return x&(1<<3) != 0
}

func ShiftByVariable() int {
	s := 3
	return 1 << s
}

func Uint8BitOps() uint8 {
	var x uint8 = 0xAA
	var y uint8 = 0x55
	return x & y
}

func IntBitSign() int {
	x := -1
	return x >> 1 // arithmetic shift right (sign extension)
}

func BitCount() int {
	n := 0xFF
	count := 0
	for n != 0 {
		count++
		n &= n - 1
	}
	return count
}

func ReverseBits() uint8 {
	var b uint8 = 0b11010110
	var result uint8
	for i := 0; i < 8; i++ {
		result = (result << 1) | (b & 1)
		b >>= 1
	}
	return result
}
