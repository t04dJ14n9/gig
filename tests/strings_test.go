package tests

import "testing"

func TestStringConcat(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	s := "hello"
	return s + " world"
}`, "hello world")
}

func TestStringConcatLoop(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	s := ""
	for i := 0; i < 3; i++ {
		s = s + "ab"
	}
	return s
}`, "ababab")
}

func TestStringLen(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	return len("hello")
}`, 5)
}

func TestStringIndex(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	s := "abcde"
	return string(s[0]) + string(s[4])
}`, "ae")
}

func TestStringComparison(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := "abc"
	b := "abd"
	if a < b {
		return 1
	}
	return 0
}`, 1)
}

func TestStringEquality(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := "hello"
	b := "hello"
	c := "world"
	result := 0
	if a == b { result = result + 1 }
	if a != c { result = result + 10 }
	return result
}`, 11)
}

func TestStringEmptyCheck(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := ""
	if len(s) == 0 {
		return 1
	}
	return 0
}`, 1)
}
