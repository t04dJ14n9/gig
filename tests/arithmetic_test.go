package tests

import "testing"

func TestArithmeticAddition(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 2 + 3 }`, 5)
}

func TestArithmeticSubtraction(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 10 - 4 }`, 6)
}

func TestArithmeticMultiplication(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 6 * 7 }`, 42)
}

func TestArithmeticDivision(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 20 / 4 }`, 5)
}

func TestArithmeticModulo(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 17 % 5 }`, 2)
}

func TestArithmeticComplexExpr(t *testing.T) {
	runInt(t, `package main; func Compute() int { return (2 + 3) * 4 - 10 / 2 }`, 15)
}

func TestArithmeticNegation(t *testing.T) {
	runInt(t, `package main; func Compute() int { x := 42; return -x }`, -42)
}

func TestArithmeticChainedOps(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := 10
	b := a * 2
	c := b + a
	d := c - 5
	return d / 5
}`, 5)
}

func TestArithmeticOverflow(t *testing.T) {
	// int64 wrapping behavior
	runInt(t, `package main
func Compute() int {
	x := 9223372036854775807
	return x + 1
}`, -9223372036854775808)
}

func TestArithmeticPrecedence(t *testing.T) {
	runInt(t, `package main; func Compute() int { return 2 + 3*4 }`, 14)
}
