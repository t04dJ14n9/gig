package tests

import (
	"testing"

	"github.com/t04dJ14n9/gig"
	_ "github.com/t04dJ14n9/gig/stdlib/packages"
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

// TestMapByteKey
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
