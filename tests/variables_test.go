package tests

import "testing"

func TestVariablesDeclareAndUse(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 10
	y := 20
	z := x + y
	return z
}`, 30)
}

func TestVariablesReassignment(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 1
	x = x + 10
	x = x * 2
	return x
}`, 22)
}

func TestVariablesMultipleDecl(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	return a + b + c + d + e
}`, 15)
}

func TestVariablesZeroValues(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	var x int
	return x
}`, 0)
}

func TestVariablesStringZeroValue(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	var s string
	return s
}`, "")
}

func TestVariablesShadowing(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 10
	y := 1
	if y > 0 {
		x := 20
		_ = x
	}
	return x
}`, 10)
}
