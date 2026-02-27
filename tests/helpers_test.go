package tests

import (
	"testing"

	"gig"

	_ "gig/stdlib/packages"
)

// runInt builds and runs source, expecting an int64 result.
func runInt(t *testing.T, source string, expected int64) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	// Handle different integer types
	var got int64
	switch v := result.(type) {
	case int64:
		got = v
	case int:
		got = int64(v)
	case int32:
		got = int64(v)
	default:
		t.Fatalf("expected int type, got %T", result)
	}
	if got != expected {
		t.Errorf("expected %d, got %v", expected, got)
	}
}

// runStr builds and runs source, expecting a string result.
func runStr(t *testing.T, source string, expected string) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result.(string) != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// runBool builds and runs source, expecting a bool result.
func runBool(t *testing.T, source string, expected bool) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result.(bool) != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// expectBuildError asserts that Build returns an error.
func expectBuildError(t *testing.T, source string) {
	t.Helper()
	_, err := gig.Build(source)
	if err == nil {
		t.Error("expected build error, got nil")
	}
}

// runFloat builds and runs source, expecting a float64 result.
func runFloat(t *testing.T, source string, expected float64) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	var got float64
	switch v := result.(type) {
	case float64:
		got = v
	case float32:
		got = float64(v)
	case int:
		got = float64(v)
	case int64:
		got = float64(v)
	default:
		t.Fatalf("expected float type, got %T", result)
	}
	// Allow small floating point tolerance
	diff := got - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.0001 {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

// runIntCompare builds and runs source with interpreter, comparing with native Go result.
// funcName specifies which function to call.
func runIntCompare(t *testing.T, source string, funcName string, expected int64) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run(funcName)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	var got int64
	switch v := result.(type) {
	case int64:
		got = v
	case int:
		got = int64(v)
	case int32:
		got = int64(v)
	default:
		t.Fatalf("expected int type, got %T", result)
	}
	if got != expected {
		t.Errorf("interpreter result mismatch for %s: expected %d, got %d", funcName, expected, got)
	}
}

// runStrCompare builds and runs source with interpreter, comparing with native Go result.
func runStrCompare(t *testing.T, source string, funcName string, expected string) {
	t.Helper()
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run(funcName)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	got, ok := result.(string)
	if !ok {
		t.Fatalf("expected string type, got %T", result)
	}
	if got != expected {
		t.Errorf("interpreter result mismatch for %s: expected %q, got %q", funcName, expected, got)
	}
}

// TestMapByteKey tests map with byte key operations
func TestMapByteKey(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	m['a'] = 1
	return m['a']
}
`
	runInt(t, code, 1)
}

// TestMapByteKeyMissing tests map access with missing key
func TestMapByteKeyMissing(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	return m['a']
}
`
	runInt(t, code, 0)
}

// TestMapByteKeyInc tests map increment
func TestMapByteKeyInc(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	m['a']++
	return m['a']
}
`
	runInt(t, code, 1)
}

// TestMapByteLen tests map length with byte keys
func TestMapByteLen(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	m['a'] = 1
	m['b'] = 2
	return len(m)
}
`
	runInt(t, code, 2)
}

// TestMapCommaOk tests comma-ok idiom
func TestMapCommaOk(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	m['a'] = 1
	
	v, ok := m['a']
	if ok {
		return v
	}
	return 0
}
`
	runInt(t, code, 1)
}

// TestMapCommaOkMissing tests comma-ok with missing key
func TestMapCommaOkMissing(t *testing.T) {
	code := `
package main

func Compute() int {
	m := make(map[byte]int)
	m['a'] = 1
	
	v, ok := m['b']
	if ok {
		return v
	}
	return -1
}
`
	runInt(t, code, -1)
}
