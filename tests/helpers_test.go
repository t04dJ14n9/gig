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
