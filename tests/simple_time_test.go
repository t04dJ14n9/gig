package tests

import (
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

func TestSimpleTimeSub(t *testing.T) {
	code := `package main
import "time"
func Test() int {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	return int(t2.Sub(t1).Hours())
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Test")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result != 24 {
		t.Errorf("expected 24, got %v", result)
	}
}

func TestSimpleTimeAdd(t *testing.T) {
	code := `package main
import "time"
func Test() int {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t.Add(24 * time.Hour)
	return t2.Day()
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Test")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result != 2 {
		t.Errorf("expected 2, got %v", result)
	}
}

func TestSimpleContextValue(t *testing.T) {
	code := `package main
import "context"
func Test() int {
	ctx := context.Background()
	ctx2 := context.WithValue(ctx, "key", "value")
	v := ctx2.Value("key")
	if v == nil {
		return -1
	}
	s, ok := v.(string)
	if !ok {
		return -2
	}
	if s == "value" {
		return 1
	}
	return 0
}
`
	prog, err := gig.Build(code)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	result, err := prog.Run("Test")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}
