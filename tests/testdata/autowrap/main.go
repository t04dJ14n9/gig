package autowrap

import "fmt"

// WithPackage tests explicit package main
func WithPackage() int {
	return 99
}

// WithImport tests auto-wrap with import
func WithImport() string {
	return fmt.Sprintf("hello %d", 42)
}

// add is a helper function
func add(a, b int) int { return a + b }

// Compute tests multiple functions
func Compute() int {
	return add(10, 20)
}
