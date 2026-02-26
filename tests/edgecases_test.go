package tests

import "testing"

// Edge case tests: boundary conditions, special values, unusual patterns.

func TestMaxInt64(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	return 9223372036854775807
}`, 9223372036854775807)
}

func TestMinInt64(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	return -9223372036854775807 - 1
}`, -9223372036854775808)
}

func TestDivisionByMinusOne(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 42
	return x / (-1)
}`, -42)
}

func TestModuloNegative(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	return (-7) % 3
}`, -1)
}

func TestEmptyString(t *testing.T) {
	runStr(t, `package main
func Compute() string { return "" }`, "")
}

func TestLargeSlice(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		s[i] = i
	}
	sum := 0
	for _, v := range s {
		sum = sum + v
	}
	return sum
}`, 49995000)
}

func TestNestedMapLookup(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	keys := make([]string, 0)
	keys = append(keys, "a")
	keys = append(keys, "b")
	keys = append(keys, "c")
	for i, k := range keys {
		m[k] = (i + 1) * 10
	}
	sum := 0
	for _, k := range keys {
		sum = sum + m[k]
	}
	return sum
}`, 60) // 10+20+30
}

func TestZeroDivisionGuard(t *testing.T) {
	runInt(t, `package main
func safeDivide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
func Compute() int {
	return safeDivide(10, 2) + safeDivide(10, 0)
}`, 5)
}

func TestBooleanComplexExpr(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := 10
	b := 20
	c := 30
	result := 0
	if a < b && b < c { result = result + 1 }
	if a > b || c > b { result = result + 10 }
	if !(a > c) { result = result + 100 }
	if a < b && (c > 20 || b < 10) { result = result + 1000 }
	return result
}`, 1111)
}

func TestSingleElementSlice(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, 42)
	return s[0] + len(s)
}`, 43)
}

func TestEmptyMap(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	return len(m)
}`, 0)
}

func TestTightLoop(t *testing.T) {
	// Test that tight loops with complex operations work
	runInt(t, `package main
func Compute() int {
	result := 0
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if (i+j)%2 == 0 {
				result = result + 1
			}
		}
	}
	return result
}`, 5000)
}
