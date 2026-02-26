package bitwise

// And tests bitwise AND
func And() int { return 0xFF & 0x0F }

// Or tests bitwise OR
func Or() int { return 0xFF | 0x100 }

// Xor tests bitwise XOR
func Xor() int { return 0xAA ^ 0x55 }

// LeftShift tests left shift
func LeftShift() int { return 1 << 10 }

// RightShift tests right shift
func RightShift() int { return 1024 >> 5 }

// Combined tests combined bitwise operations
func Combined() int {
	a := 0xFF
	b := 0x0F
	andResult := a & b
	orResult := a | 0x100
	xorResult := 0xAA ^ 0x55
	shifted := 1 << 10
	return andResult + orResult + xorResult + shifted
}

// AndNot tests bitwise AND NOT
func AndNot() int { return 0xFF &^ 0x0F }

// PowerOfTwo tests power of two check using bitwise
func PowerOfTwo() int {
	return isPowerOfTwo(16)*10 + isPowerOfTwo(15)
}

func isPowerOfTwo(n int) int {
	if n > 0 && (n&(n-1)) == 0 {
		return 1
	}
	return 0
}
