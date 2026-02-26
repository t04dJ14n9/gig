package gig_test

import (
	"context"
	"testing"
	"time"

	"gig"
	_ "gig/packages" // register stdlib packages (fmt, strings, etc.)
)

// TestBasicArithmetic tests basic arithmetic operations.
func TestBasicArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int64
	}{
		{
			name: "addition",
			source: `
package main

func Compute() int {
	return 2 + 3
}
`,
			expected: 5,
		},
		{
			name: "subtraction",
			source: `
package main

func Compute() int {
	return 10 - 4
}
`,
			expected: 6,
		},
		{
			name: "multiplication",
			source: `
package main

func Compute() int {
	return 6 * 7
}
`,
			expected: 42,
		},
		{
			name: "division",
			source: `
package main

func Compute() int {
	return 20 / 4
}
`,
			expected: 5,
		},
		{
			name: "modulo",
			source: `
package main

func Compute() int {
	return 17 % 5
}
`,
			expected: 2,
		},
		{
			name: "complex expression",
			source: `
package main

func Compute() int {
	return (2 + 3) * 4 - 10 / 2
}
`,
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := gig.Build(tt.source)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			result, err := prog.Run("Compute")
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			if result.(int64) != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestVariables tests variable declarations and assignments.
func TestVariables(t *testing.T) {
	source := `
package main

func Compute() int {
	x := 10
	y := 20
	z := x + y
	return z
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 30 {
		t.Errorf("expected 30, got %d", result)
	}
}

// TestControlFlow tests if statements and loops.
func TestControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int64
	}{
		{
			name: "if statement true",
			source: `
package main

func Compute() int {
	x := 10
	if x > 5 {
		return 1
	}
	return 0
}
`,
			expected: 1,
		},
		{
			name: "if statement false",
			source: `
package main

func Compute() int {
	x := 3
	if x > 5 {
		return 1
	}
	return 0
}
`,
			expected: 0,
		},
		{
			name: "for loop",
			source: `
package main

func Compute() int {
	sum := 0
	for i := 1; i <= 10; i++ {
		sum = sum + i
	}
	return sum
}
`,
			expected: 55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := gig.Build(tt.source)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			result, err := prog.Run("Compute")
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			if result.(int64) != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestFunctions tests function definitions and calls.
func TestFunctions(t *testing.T) {
	source := `
package main

func add(a int, b int) int {
	return a + b
}

func Compute() int {
	return add(5, 7)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}

// TestMultipleReturn tests functions with multiple return values.
func TestMultipleReturn(t *testing.T) {
	source := `
package main

func swap(a int, b int) (int, int) {
	return b, a
}

func Compute() int {
	x, y := swap(3, 7)
	return x + y
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 10 {
		t.Errorf("expected 10, got %d", result)
	}
}

// TestRecursion tests recursive function calls.
func TestRecursion(t *testing.T) {
	source := `
package main

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n - 1)
}

func Compute() int {
	return factorial(5)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 120 {
		t.Errorf("expected 120, got %d", result)
	}
}

// TestWithContext tests context-based timeout.
func TestWithContext(t *testing.T) {
	source := `
package main

func Compute() int {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum = sum + i
	}
	return sum
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	// Test with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = prog.RunWithContext("Compute", ctx)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

// TestParameters tests passing parameters to functions.
func TestParameters(t *testing.T) {
	source := `
package main

func Compute(a int, b int) int {
	return a * b
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute", 6, 7)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

// TestStringOperations tests string operations.
func TestStringOperations(t *testing.T) {
	source := `
package main

func Compute() string {
	s := "hello"
	return s + " world"
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(string) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

// TestSecurityBannedImports tests that banned imports are rejected.
func TestSecurityBannedImports(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name: "unsafe import",
			source: `
package main

import "unsafe"

func Compute() int {
	return 42
}
`,
		},
		{
			name: "reflect import",
			source: `
package main

import "reflect"

func Compute() int {
	return 42
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gig.Build(tt.source)
			if err == nil {
				t.Error("expected error for banned import, got nil")
			}
		})
	}
}

// BenchmarkCompute runs performance benchmarks.
func BenchmarkCompute(b *testing.B) {
	source := `
package main

func Compute() int {
	sum := 0
	for i := 1; i <= 100; i++ {
		sum = sum + i
	}
	return sum
}
`
	prog, err := gig.Build(source)
	if err != nil {
		b.Fatalf("Build error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = prog.Run("Compute")
	}
}

// TestSliceBasic tests basic slice operations.
func TestSliceBasic(t *testing.T) {
	source := `
package main

func Compute() int {
	nums := make([]int, 3)
	return len(nums)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

// TestMapBasic tests basic map operations.
func TestMapBasic(t *testing.T) {
	source := `
package main

func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	return m["a"] + m["b"]
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

// TestVariadic tests variadic functions.
func TestVariadic(t *testing.T) {
	source := `
package main

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total = total + n
	}
	return total
}

func Compute() int {
	return sum(1, 2, 3, 4, 5)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 15 {
		t.Errorf("expected 15, got %d", result)
	}
}

// TestClosure tests closure support.
func TestClosure(t *testing.T) {
	source := `
package main

func makeCounter() func() int {
	count := 0
	return func() int {
		count = count + 1
		return count
	}
}

func Compute() int {
	counter := makeCounter()
	return counter() + counter() + counter()
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(int64) != 6 {
		t.Errorf("expected 6, got %d", result)
	}
}

// TestExternalCall tests calling external registered functions.
func TestExternalCall(t *testing.T) {
	source := `
package main

import "fmt"

func Compute() string {
	return fmt.Sprintf("hello %d", 42)
}
`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}

	if result.(string) != "hello 42" {
		t.Errorf("expected 'hello 42', got '%s'", result)
	}
}
