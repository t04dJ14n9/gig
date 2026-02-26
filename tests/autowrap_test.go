package tests

import "testing"

// Tests for auto-wrap feature (package main is prepended if missing).

func TestAutoWrapNoPackage(t *testing.T) {
	// No package declaration - should auto-wrap
	runInt(t, `
func Compute() int {
	return 42
}`, 42)
}

func TestAutoWrapWithPackage(t *testing.T) {
	// Explicit package main - should work as before
	runInt(t, `package main
func Compute() int {
	return 99
}`, 99)
}

func TestAutoWrapWithImport(t *testing.T) {
	// No package but has import - should auto-wrap correctly
	runStr(t, `
import "fmt"
func Compute() string {
	return fmt.Sprintf("hello %d", 42)
}`, "hello 42")
}

func TestAutoWrapMultipleFunctions(t *testing.T) {
	// No package, multiple functions
	runInt(t, `
func add(a, b int) int { return a + b }
func Compute() int {
	return add(10, 20)
}`, 30)
}
