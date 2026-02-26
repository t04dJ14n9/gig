package tests

import "testing"

// Tests for scoping, short var declarations in if/for, and lifetime patterns.

func TestIfInitShortVar(t *testing.T) {
	runInt(t, `package main
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func Compute() int {
	if v := abs(-42); v > 0 {
		return v
	}
	return 0
}`, 42)
}

func TestIfInitMultiCondition(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	result := 0
	for i := 0; i < 10; i++ {
		if rem := i % 3; rem == 0 {
			result = result + i
		}
	}
	return result
}`, 18) // 0+3+6+9
}

func TestNestedScopes(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 1
	y := 0
	if x > 0 {
		x := 10
		y = x
	}
	return x + y
}`, 11) // outer x=1, y=10 from inner scope
}

func TestForScopeIsolation(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	for i := 0; i < 3; i++ {
		x := i * 10
		sum = sum + x
	}
	return sum
}`, 30) // 0+10+20
}

func TestMultipleBlockScopes(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	result := 0
	x := 1
	if x > 0 {
		a := 10
		result = result + a
	}
	if x > 0 {
		b := 20
		result = result + b
	}
	return result
}`, 30)
}

func TestClosureCapturesOuterScope(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 100
	y := 200
	add := func() int { return x + y }
	x = 150
	return add()
}`, 350) // closure captures by reference, so x=150
}
