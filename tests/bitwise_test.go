package tests

import "testing"

func TestBitwiseAnd(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 0xFF & 0x0F }`, 15)
}

func TestBitwiseOr(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 0xFF | 0x100 }`, 0x1FF)
}

func TestBitwiseXor(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 0xAA ^ 0x55 }`, 0xFF)
}

func TestBitwiseLeftShift(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 1 << 10 }`, 1024)
}

func TestBitwiseRightShift(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 1024 >> 5 }`, 32)
}

func TestBitwiseCombined(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := 0xFF
	b := 0x0F
	andResult := a & b
	orResult := a | 0x100
	xorResult := 0xAA ^ 0x55
	shifted := 1 << 10
	return andResult + orResult + xorResult + shifted
}`, 1805)
}

func TestBitwiseAndNot(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 0xFF &^ 0x0F }`, 0xF0)
}

func TestBitwisePowerOfTwo(t *testing.T) {
	runInt(t, `package main
func isPowerOfTwo(n int) int {
	if n > 0 && (n & (n - 1)) == 0 {
		return 1
	}
	return 0
}
func Compute() int {
	return isPowerOfTwo(16)*10 + isPowerOfTwo(15)
}`, 10)
}
