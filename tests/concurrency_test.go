package tests

import (
	"sync"
	"testing"

	"github.com/t04dJ14n9/gig"
)

// TestConcurrentProgramMethodResolution verifies that two programs compiled
// and run concurrently don't interfere with each other's method dispatch.
// This is a regression test for the global method resolver race (C1).
func TestConcurrentProgramMethodResolution(t *testing.T) {
	srcA := `
package main

type Greeter struct{ Name string }
func (g Greeter) Greet() string { return "Hello, " + g.Name }
func Run() string {
	g := Greeter{Name: "Alice"}
	return g.Greet()
}
`
	srcB := `
package main

type Greeter struct{ Name string }
func (g Greeter) Greet() string { return "Hi, " + g.Name }
func Run() string {
	g := Greeter{Name: "Bob"}
	return g.Greet()
}
`
	var wg sync.WaitGroup
	errCh := make(chan error, 20)

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			prog, err := gig.Build(srcA)
			if err != nil {
				errCh <- err
				return
			}
			result, err := prog.Run("Run")
			if err != nil {
				errCh <- err
				return
			}
			if result != "Hello, Alice" {
				t.Errorf("srcA: got %q, want %q", result, "Hello, Alice")
			}
		}()
		go func() {
			defer wg.Done()
			prog, err := gig.Build(srcB)
			if err != nil {
				errCh <- err
				return
			}
			result, err := prog.Run("Run")
			if err != nil {
				errCh <- err
				return
			}
			if result != "Hi, Bob" {
				t.Errorf("srcB: got %q, want %q", result, "Hi, Bob")
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Errorf("unexpected error: %v", err)
	}
}
