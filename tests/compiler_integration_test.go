package tests

import (
	"testing"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/compiler"
	"git.woa.com/youngjin/gig/importer"
)

var _ = bytecode.OpCode(0) // silence unused import warning

// TestNewCompiler tests the compiler constructor.
func TestNewCompiler(t *testing.T) {
	lookup := importer.NewPackageLookup(importer.NewRegistry())
	c := compiler.NewCompiler(lookup)
	if c == nil {
		t.Fatal("NewCompiler returned nil")
	}

	// Verify it implements the interface
	var _ compiler.Compiler = c
}

// TestCompileEmptyPackage tests compiling an empty package.
func TestCompileEmptyPackage(t *testing.T) {
	source := `
package main

func main() {
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if result.Program == nil {
		t.Fatal("Program is nil")
	}
	if len(result.Program.Functions) == 0 {
		t.Error("Expected at least main function")
	}
}

// TestCompileSimpleFunction tests compiling a simple function.
func TestCompileSimpleFunction(t *testing.T) {
	source := `
package main

func add(a, b int) int {
	return a + b
}

func main() {
	_ = add(1, 2)
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if _, ok := result.Program.Functions["add"]; !ok {
		t.Error("add function not found")
	}
	if _, ok := result.Program.Functions["main"]; !ok {
		t.Error("main function not found")
	}
}

// TestCompileWithConstants tests compilation with various constant types.
func TestCompileWithConstants(t *testing.T) {
	source := `
package main

func main() {
	i := 42
	s := "hello"
	b := true
	f := 3.14
	_, _, _, _ = i, s, b, f
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify at least some constants were collected
	// (Note: compiler may deduplicate or inline some constants)
	if len(result.Program.Constants) < 1 {
		t.Errorf("Expected at least 1 constant, got %d", len(result.Program.Constants))
	}
}

// TestCompileWithControlFlow tests compilation with control flow.
func TestCompileWithControlFlow(t *testing.T) {
	source := `
package main

func main() {
	x := 10
	if x > 5 {
		x = x + 1
	} else {
		x = x - 1
	}
	for i := 0; i < 10; i++ {
		x += i
	}
	_ = x
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	main := result.Program.Functions["main"]
	if main == nil {
		t.Fatal("main function not found")
	}

	hasJump := false
	for _, b := range main.Instructions {
		if b == byte(bytecode.OpJump) || b == byte(bytecode.OpJumpTrue) || b == byte(bytecode.OpJumpFalse) {
			hasJump = true
			break
		}
	}
	if !hasJump {
		t.Error("Expected jump instructions in control flow")
	}
}

// TestCompileWithStructs tests compilation with struct operations.
func TestCompileWithStructs(t *testing.T) {
	source := `
package main

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) Sum() int {
	return p.X + p.Y
}

func (p *Point) Add(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func main() {
	p := NewPoint(1, 2)
	_ = p.Sum()
	p.Add(10, 20)
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expected := []string{"NewPoint", "Sum", "Add", "main"}
	for _, name := range expected {
		if _, ok := result.Program.Functions[name]; !ok {
			t.Errorf("Function %s not found", name)
		}
	}
}

// TestCompileWithClosures tests compilation with closures.
func TestCompileWithClosures(t *testing.T) {
	source := `
package main

func main() {
	x := 10
	f := func() int {
		return x + 1
	}
	_ = f()
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(result.Program.Functions) < 2 {
		t.Errorf("Expected at least 2 functions (main + closures), got %d", len(result.Program.Functions))
	}
}

// TestCompileRecursiveFunction tests compilation of recursive functions.
func TestCompileRecursiveFunction(t *testing.T) {
	source := `
package main

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	_ = factorial(5)
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if _, ok := result.Program.Functions["factorial"]; !ok {
		t.Error("factorial function not found")
	}
}

// TestCompileMultipleReturn tests compilation with multiple return values.
func TestCompileMultipleReturn(t *testing.T) {
	source := `
package main

func divmod(a, b int) (int, int) {
	return a / b, a % b
}

func main() {
	q, r := divmod(10, 3)
	_, _ = q, r
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if _, ok := result.Program.Functions["divmod"]; !ok {
		t.Error("divmod function not found")
	}
}

// TestCompileArithmetic tests compilation with arithmetic operators.
func TestCompileArithmetic(t *testing.T) {
	source := `
package main

func main() {
	a := 10
	b := 3
	_ = a + b
	_ = a - b
	_ = a * b
	_ = a / b
	_ = a % b
	_ = -a
	a++
	a--
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileBitwise tests compilation with bitwise operators.
func TestCompileBitwise(t *testing.T) {
	source := `
package main

func main() {
	a := 0xFF
	b := 0x0F
	_ = a & b
	_ = a | b
	_ = a ^ b
	_ = a &^ b
	_ = a << 2
	_ = a >> 2
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileLogical tests compilation with logical operators.
func TestCompileLogical(t *testing.T) {
	source := `
package main

func main() {
	a := true
	b := false
	_ = a && b
	_ = a || b
	_ = !a
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompilePointer tests compilation with pointer operations.
func TestCompilePointer(t *testing.T) {
	source := `
package main

func main() {
	x := 10
	p := &x
	_ = *p
	*p = 20
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileChannels tests compilation with channel operations.
func TestCompileChannels(t *testing.T) {
	source := `
package main

func main() {
	ch := make(chan int, 10)
	ch <- 1
	v := <-ch
	_ = v
	close(ch)
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileWithDefer tests compilation with defer.
func TestCompileWithDefer(t *testing.T) {
	source := `
package main

func cleanup() {}

func main() {
	defer cleanup()
	_ = 42
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileWithGoroutine tests compilation with goroutines.
func TestCompileWithGoroutine(t *testing.T) {
	source := `
package main

func worker(n int) {
	_ = n
}

func main() {
	go worker(42)
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileWithSlices tests compilation with slice operations.
func TestCompileWithSlices(t *testing.T) {
	source := `
package main

func main() {
	arr := []int{1, 2, 3, 4, 5}
	_ = arr[0]
	arr[1] = 10
	arr = append(arr, 6)
	s := make([]int, 10)
	_ = s[0]
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}

// TestCompileWithMaps tests compilation with map operations.
func TestCompileWithMaps(t *testing.T) {
	source := `
package main

func main() {
	m := make(map[string]int)
	m["a"] = 1
	v := m["a"]
	_ = v
	m2 := map[string]int{"x": 1, "y": 2}
	_ = m2["x"]
}
`
	result, err := compiler.Build(source, importer.NewRegistry())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if result.Program == nil {
		t.Error("Program is nil")
	}
}
