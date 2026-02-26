package tests

import (
	"testing"

	"gig"
	_ "gig/packages"
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
	if result.(int64) != expected {
		t.Errorf("expected %d, got %v", expected, result)
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
