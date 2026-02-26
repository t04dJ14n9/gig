package tests

import "testing"

// Tests for named return values.

func TestNamedReturnBasic(t *testing.T) {
	runInt(t, `package main
func double(x int) (result int) {
	result = x * 2
	return result
}
func Compute() int { return double(21) }`, 42)
}

func TestNamedReturnMultiple(t *testing.T) {
	runInt(t, `package main
func divmod(a, b int) (quotient int, remainder int) {
	quotient = a / b
	remainder = a % b
	return quotient, remainder
}
func Compute() int {
	q, r := divmod(17, 5)
	return q*10 + r
}`, 32)
}

func TestNamedReturnZeroValue(t *testing.T) {
	runInt(t, `package main
func maybeDouble(x int, doIt int) (result int) {
	if doIt > 0 {
		result = x * 2
	}
	return result
}
func Compute() int {
	return maybeDouble(10, 1) + maybeDouble(10, 0)
}`, 20) // 20 + 0
}
