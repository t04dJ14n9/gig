package tests

import (
	"context"
	"testing"
	"time"

	"gig"
	_ "gig/packages"
)

func TestIfTrue(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 10
	if x > 5 { return 1 }
	return 0
}`, 1)
}

func TestIfFalse(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 3
	if x > 5 { return 1 }
	return 0
}`, 0)
}

func TestIfElse(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	x := 3
	if x > 5 {
		return 1
	} else {
		return -1
	}
}`, -1)
}

func TestIfElseChain(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int64
	}{
		{"negative", `package main
func classify(x int) int {
	if x < 0 { return -1 } else if x == 0 { return 0 } else { return 1 }
}
func Compute() int { return classify(-5) }`, -1},
		{"zero", `package main
func classify(x int) int {
	if x < 0 { return -1 } else if x == 0 { return 0 } else { return 1 }
}
func Compute() int { return classify(0) }`, 0},
		{"positive", `package main
func classify(x int) int {
	if x < 0 { return -1 } else if x == 0 { return 0 } else { return 1 }
}
func Compute() int { return classify(42) }`, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { runInt(t, tt.source, tt.expected) })
	}
}

func TestForLoop(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	for i := 1; i <= 10; i++ {
		sum = sum + i
	}
	return sum
}`, 55)
}

func TestForConditionOnly(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	i := 0
	sum := 0
	for i < 5 {
		sum = sum + i
		i = i + 1
	}
	return sum
}`, 10)
}

func TestNestedFor(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			sum = sum + 1
		}
	}
	return sum
}`, 25)
}

func TestForBreak(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	for i := 0; i < 100; i++ {
		if i >= 5 {
			break
		}
		sum = sum + i
	}
	return sum
}`, 10)
}

func TestForContinue(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		sum = sum + i
	}
	return sum
}`, 25)
}

func TestBooleanAndOr(t *testing.T) {
	runInt(t, `package main
func Compute() int {
	a := true
	b := false
	result := 0
	if a && !b { result = result + 1 }
	if a || b { result = result + 10 }
	if !b { result = result + 100 }
	return result
}`, 111)
}

func TestContextTimeout(t *testing.T) {
	source := `package main
func Compute() int {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum = sum + i
	}
	return sum
}`
	prog, err := gig.Build(source)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err = prog.RunWithContext("Compute", ctx)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
