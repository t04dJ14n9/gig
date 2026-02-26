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

// TestSliceAppendAndIndex tests slice append, indexing, and len/cap.
func TestSliceAppendAndIndex(t *testing.T) {
	source := `
package main

func Compute() int {
	s := make([]int, 0)
	s = append(s, 10)
	s = append(s, 20)
	s = append(s, 30)
	return s[0] + s[1] + s[2] + len(s)
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

	if result.(int64) != 63 {
		t.Errorf("expected 63, got %d", result)
	}
}

// TestSliceElementAssignment tests assigning to slice elements via indexing.
func TestSliceElementAssignment(t *testing.T) {
	source := `
package main

func Compute() int {
	s := make([]int, 3)
	s[0] = 100
	s[1] = 200
	s[2] = 300
	return s[0] + s[1] + s[2]
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

	if result.(int64) != 600 {
		t.Errorf("expected 600, got %d", result)
	}
}

// TestMapIteration tests iterating over a map with a for-range loop.
func TestMapIteration(t *testing.T) {
	source := `
package main

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

	if result.(int64) != 60 {
		t.Errorf("expected 60, got %d", result)
	}
}

// TestMapDelete tests deleting keys from a map.
func TestMapDelete(t *testing.T) {
	source := `
package main

func Compute() int {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	delete(m, "b")
	return len(m)
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

	if result.(int64) != 2 {
		t.Errorf("expected 2, got %d", result)
	}
}

// TestClosureCapture tests that closures properly capture and mutate variables.
func TestClosureCapture(t *testing.T) {
	source := `
package main

func Compute() int {
	x := 10
	add := func(n int) int {
		x = x + n
		return x
	}
	a := add(5)
	b := add(3)
	return a + b
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

	// add(5) -> x=15, returns 15; add(3) -> x=18, returns 18; 15+18=33
	if result.(int64) != 33 {
		t.Errorf("expected 33, got %d", result)
	}
}

// TestNestedClosure tests nested closures capturing variables from multiple scopes.
func TestNestedClosure(t *testing.T) {
	source := `
package main

func makeAdder(base int) func(int) int {
	return func(x int) int {
		return base + x
	}
}

func Compute() int {
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	return add5(3) + add10(7)
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

	// add5(3)=8, add10(7)=17, total=25
	if result.(int64) != 25 {
		t.Errorf("expected 25, got %d", result)
	}
}

// TestIfElseChain tests if/else if/else chains.
func TestIfElseChain(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int64
	}{
		{
			name: "first branch",
			source: `
package main

func classify(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	} else {
		return 1
	}
}

func Compute() int {
	return classify(-5)
}
`,
			expected: -1,
		},
		{
			name: "middle branch",
			source: `
package main

func classify(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	} else {
		return 1
	}
}

func Compute() int {
	return classify(0)
}
`,
			expected: 0,
		},
		{
			name: "else branch",
			source: `
package main

func classify(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	} else {
		return 1
	}
}

func Compute() int {
	return classify(42)
}
`,
			expected: 1,
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

// TestBooleanOperations tests boolean logic and comparisons.
func TestBooleanOperations(t *testing.T) {
	source := `
package main

func Compute() int {
	a := true
	b := false
	result := 0
	if a && !b {
		result = result + 1
	}
	if a || b {
		result = result + 10
	}
	if !b {
		result = result + 100
	}
	return result
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

	if result.(int64) != 111 {
		t.Errorf("expected 111, got %d", result)
	}
}

// TestNestedLoops tests nested for loops.
func TestNestedLoops(t *testing.T) {
	source := `
package main

func Compute() int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			sum = sum + 1
		}
	}
	return sum
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

	if result.(int64) != 25 {
		t.Errorf("expected 25, got %d", result)
	}
}

// TestMutualRecursion tests mutually recursive functions.
func TestMutualRecursion(t *testing.T) {
	source := `
package main

func isEven(n int) bool {
	if n == 0 {
		return true
	}
	return isOdd(n - 1)
}

func isOdd(n int) bool {
	if n == 0 {
		return false
	}
	return isEven(n - 1)
}

func Compute() int {
	if isEven(10) {
		return 1
	}
	return 0
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

	if result.(int64) != 1 {
		t.Errorf("expected 1, got %d", result)
	}
}

// TestMultipleReturnValues tests extracting multiple return values.
func TestMultipleReturnValues(t *testing.T) {
	source := `
package main

func divmod(a int, b int) (int, int) {
	return a / b, a % b
}

func Compute() int {
	q, r := divmod(17, 5)
	return q*10 + r
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

	// 17/5=3, 17%5=2 => 3*10+2=32
	if result.(int64) != 32 {
		t.Errorf("expected 32, got %d", result)
	}
}

// TestStringConcat tests string concatenation in a loop.
func TestStringConcat(t *testing.T) {
	source := `
package main

func Compute() string {
	s := ""
	for i := 0; i < 3; i++ {
		s = s + "ab"
	}
	return s
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

	if result.(string) != "ababab" {
		t.Errorf("expected 'ababab', got '%s'", result)
	}
}

// TestBitwiseOperations tests bitwise AND, OR, XOR, shifts.
func TestBitwiseOperations(t *testing.T) {
	source := `
package main

func Compute() int {
	a := 0xFF
	b := 0x0F
	andResult := a & b
	orResult := a | 0x100
	xorResult := 0xAA ^ 0x55
	shifted := 1 << 10
	return andResult + orResult + xorResult + shifted
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

	// 0xFF & 0x0F = 15
	// 0xFF | 0x100 = 0x1FF = 511
	// 0xAA ^ 0x55 = 0xFF = 255
	// 1 << 10 = 1024
	// total = 15 + 511 + 255 + 1024 = 1805
	if result.(int64) != 1805 {
		t.Errorf("expected 1805, got %d", result)
	}
}

// TestFibonacci tests a classic iterative fibonacci.
func TestFibonacci(t *testing.T) {
	source := `
package main

func fib(n int) int {
	if n <= 1 {
		return n
	}
	a := 0
	b := 1
	for i := 2; i <= n; i++ {
		c := a + b
		a = b
		b = c
	}
	return b
}

func Compute() int {
	return fib(20)
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

	if result.(int64) != 6765 {
		t.Errorf("expected 6765, got %d", result)
	}
}

// TestExternalStrings tests calling strings package functions.
func TestExternalStrings(t *testing.T) {
	source := `
package main

import "strings"

func Compute() string {
	return strings.ToUpper("hello world")
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

	if result.(string) != "HELLO WORLD" {
		t.Errorf("expected 'HELLO WORLD', got '%s'", result)
	}
}

// TestExternalStrconv tests calling strconv package functions.
func TestExternalStrconv(t *testing.T) {
	source := `
package main

import "strconv"

func Compute() string {
	return strconv.Itoa(42)
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

	if result.(string) != "42" {
		t.Errorf("expected '42', got '%s'", result)
	}
}

// TestSliceForRange tests iterating over a slice with for-range.
func TestSliceForRange(t *testing.T) {
	source := `
package main

func Compute() int {
	nums := make([]int, 0)
	nums = append(nums, 10)
	nums = append(nums, 20)
	nums = append(nums, 30)
	sum := 0
	for _, v := range nums {
		sum = sum + v
	}
	return sum
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

	if result.(int64) != 60 {
		t.Errorf("expected 60, got %d", result)
	}
}

// TestFunctionAsValue tests passing functions as values.
func TestFunctionAsValue(t *testing.T) {
	source := `
package main

func apply(f func(int) int, x int) int {
	return f(x)
}

func double(x int) int {
	return x * 2
}

func triple(x int) int {
	return x * 3
}

func Compute() int {
	return apply(double, 5) + apply(triple, 5)
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

	// double(5)=10, triple(5)=15, total=25
	if result.(int64) != 25 {
		t.Errorf("expected 25, got %d", result)
	}
}

// TestExternalFmtSprintf tests fmt.Sprintf with multiple format verbs.
func TestExternalFmtSprintf(t *testing.T) {
	source := `
package main

import "fmt"

func Compute() string {
	return fmt.Sprintf("%s is %d years old", "Alice", 30)
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

	if result.(string) != "Alice is 30 years old" {
		t.Errorf("expected 'Alice is 30 years old', got '%s'", result)
	}
}
