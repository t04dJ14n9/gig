package tests

import "testing"

func TestMapBasicOps(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	return m["a"] + m["b"]
}`, 3)
}

func TestMapIteration(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["x"] = 10
	m["y"] = 20
	m["z"] = 30
	sum := 0
	for _, v := range m {
		sum = sum + v
	}
	return sum
}`, 60)
}

func TestMapDelete(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	delete(m, "b")
	return len(m)
}`, 2)
}

func TestMapLen(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	m["d"] = 4
	return len(m)
}`, 4)
}

func TestMapOverwrite(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["key"] = 10
	m["key"] = 42
	return m["key"]
}`, 42)
}

func TestMapIntKeys(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[int]int)
	m[1] = 10
	m[2] = 20
	m[3] = 30
	return m[1] + m[2] + m[3]
}`, 60)
}

func TestMapPassToFunction(t *testing.T) {
	runInt(t, `package main
func sumValues(m map[string]int) int {
	total := 0
	for _, v := range m {
		total = total + v
	}
	return total
}
func Compute() int {
	m := make(map[string]int)
	m["a"] = 100
	m["b"] = 200
	return sumValues(m)
}`, 300)
}
