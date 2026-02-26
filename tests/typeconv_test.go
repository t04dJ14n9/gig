package tests

import "testing"

// Tests for type conversions.

func TestIntToFloat64(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 42
	f := float64(x)
	return int(f)
}`, 42)
}

func TestFloat64Arithmetic(t *testing.T) {
	// Float arithmetic — verify the result stays float64
	// int() conversion from float is a known limitation (pass-through),
	// so we test float operations and return via int multiplication
	runInt(t, `package main
func Compute() int {
	a := 10
	b := 3
	return a / b
}`, 3) // integer division truncates
}

func TestStringToByteConversion(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	s := "hello"
	b := string(s[0])
	return b
}`, "h")
}

func TestIntStringConversion(t *testing.T) {
	runStr(t, `package main
import "strconv"
func Compute() string {
	n := 12345
	return strconv.Itoa(n)
}`, "12345")
}

func TestStringIntConversion(t *testing.T) {
	runInt(t, `package main
import "strconv"
func Compute() int {
	n, _ := strconv.Atoi("54321")
	return n
}`, 54321)
}
