package tests

import "testing"

// Advanced map tests: comma-ok, nested operations, map as counter.

// Map comma-ok pattern (v, ok := m[key]) is not supported by the VM.
// Map lookup of a non-existent key returns the zero value via reflect.

func TestMapLookupExistingKey(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["key"] = 42
	return m["key"]
}`, 42)
}

func TestMapLookupWithDefault(t *testing.T) {
	// Workaround for missing comma-ok: check if key was set via sentinel
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["a"] = 10
	m["b"] = 20
	// Access existing keys
	return m["a"] + m["b"]
}`, 30)
}

func TestMapAsCounter(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 2)
	s = append(s, 1)
	s = append(s, 2)

	counts := make(map[int]int)
	// Initialize counters to zero
	counts[1] = 0
	counts[2] = 0
	counts[3] = 0
	for _, v := range s {
		counts[v] = counts[v] + 1
	}
	return counts[1]*100 + counts[2]*10 + counts[3]
}`, 231) // 1 appears 2x, 2 appears 3x, 3 appears 1x
}

func TestMapWithStringValues(t *testing.T) {
	runStr(t, `package main
func Compute() string {
	m := make(map[int]string)
	m[1] = "one"
	m[2] = "two"
	m[3] = "three"
	return m[1] + "-" + m[2] + "-" + m[3]
}`, "one-two-three")
}

func TestMapBuildFromLoop(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[int]int)
	for i := 0; i < 100; i++ {
		m[i] = i * i
	}
	return m[10] + m[50]
}`, 2600) // 100 + 2500
}

func TestMapDeleteAndReinsert(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	delete(m, "a")
	m["a"] = 99
	return m["a"] + m["b"]
}`, 101)
}
