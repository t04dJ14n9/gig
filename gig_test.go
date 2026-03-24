package gig

import (
	"strings"
	"sync"
	"testing"

	_ "git.woa.com/youngjin/gig/stdlib/packages" // register stdlib packages
)

// TestAutoImport_SinglePackage verifies that a program referencing fmt without
// an explicit import declaration is compiled and executed successfully.
func TestAutoImport_SinglePackage(t *testing.T) {
	source := `
package main

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed (expected autoImport to inject fmt): %v", err)
	}

	result, err := prog.Run("Greet", "World")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	want := "Hello, World!"
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}
}

// TestAutoImport_MultiplePackages verifies that multiple missing imports are
// all injected automatically in a single Build call.
func TestAutoImport_MultiplePackages(t *testing.T) {
	source := `
package main

func Format(name string) string {
	upper := strings.ToUpper(name)
	return fmt.Sprintf("Hello, %s! Pi=%.2f", upper, math.Pi)
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed (expected autoImport to inject fmt/strings/math): %v", err)
	}

	result, err := prog.Run("Format", "world")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, ok := result.(string)
	if !ok {
		t.Fatalf("result type = %T, want string", result)
	}
	if !strings.Contains(got, "WORLD") {
		t.Errorf("result %q does not contain upper-cased name", got)
	}
	if !strings.Contains(got, "3.14") {
		t.Errorf("result %q does not contain Pi value", got)
	}
}

// TestAutoImport_NoDuplicateImport verifies that when the user already has an
// explicit import, autoImport does not inject a duplicate, and the program
// still compiles and runs correctly.
func TestAutoImport_NoDuplicateImport(t *testing.T) {
	source := `
package main

import "fmt"

func Greet(name string) string {
	return fmt.Sprintf("Hi, %s!", name)
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed with explicit import: %v", err)
	}

	result, err := prog.Run("Greet", "Alice")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	want := "Hi, Alice!"
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}
}

// TestAutoImport_NoPackageUsed verifies that a program with no external package
// references compiles and runs without any auto-imported packages.
func TestAutoImport_NoPackageUsed(t *testing.T) {
	source := `
package main

func Compute() int {
	a, b := 3, 4
	return a + b
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	result, err := prog.Run("Compute")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	var got int64
	switch v := result.(type) {
	case int64:
		got = v
	case int:
		got = int64(v)
	default:
		t.Fatalf("unexpected result type %T", result)
	}
	if got != 7 {
		t.Errorf("result = %v, want 7", got)
	}
}

// TestAutoImport_UnknownPackageFails verifies that referencing a completely
// unknown package (not registered) still produces a compile error.
func TestAutoImport_UnknownPackageFails(t *testing.T) {
	source := `
package main

func Foo() string {
	return unknownpkg.DoSomething()
}
`
	_, err := Build(source)
	if err == nil {
		t.Fatal("expected Build to fail for unregistered package, but it succeeded")
	}
}

// ---------------------------------------------------------------------------
// Panic ban tests
// ---------------------------------------------------------------------------

// TestPanicBan_DefaultRejectsPanic verifies that panic() is rejected at compile
// time by default (without WithAllowPanic).
func TestPanicBan_DefaultRejectsPanic(t *testing.T) {
	source := `
package main

func Fail() int {
	panic("boom")
	return 0
}
`
	_, err := Build(source)
	if err == nil {
		t.Fatal("expected Build to fail for panic() without WithAllowPanic, but it succeeded")
	}
	if !strings.Contains(err.Error(), "panic()") {
		t.Errorf("error should mention panic(), got: %v", err)
	}
}

// TestPanicBan_WithAllowPanicCompiles verifies that panic() compiles
// successfully when WithAllowPanic() is set.
func TestPanicBan_WithAllowPanicCompiles(t *testing.T) {
	source := `
package main

func Fail() int {
	defer func() { recover() }()
	panic("boom")
	return 0
}
`
	prog, err := Build(source, WithAllowPanic())
	if err != nil {
		t.Fatalf("Build with WithAllowPanic failed: %v", err)
	}

	result, err := prog.Run("Fail")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// After recover(), the function returns 0 (zero value) or nil
	if result != nil {
		got, ok := toInt64(result)
		if !ok {
			t.Fatalf("unexpected result type %T", result)
		}
		if got != 0 {
			t.Errorf("result = %v, want 0", got)
		}
	}
}

// TestPanicBan_NoPanicCodeCompiles verifies that code without panic()
// compiles successfully with the default settings.
func TestPanicBan_NoPanicCodeCompiles(t *testing.T) {
	source := `
package main

func Safe() int {
	return 42
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed for panic-free code: %v", err)
	}

	result, err := prog.Run("Safe")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, _ := toInt64(result)
	if got != 42 {
		t.Errorf("result = %v, want 42", got)
	}
}

// TestSafetyNet_RuntimePanicReturnsError verifies that a Go-level runtime panic
// is caught by the VM safety net and returned as an error instead of crashing
// the host process. We use panic() with WithAllowPanic to test the safety net
// for unrecovered panics.
func TestSafetyNet_RuntimePanicReturnsError(t *testing.T) {
	source := `
package main

func UnrecoveredPanic() int {
	panic("unrecovered!")
	return 0
}
`
	prog, err := Build(source, WithAllowPanic())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	_, err = prog.Run("UnrecoveredPanic")
	if err == nil {
		t.Fatal("expected Run to return error for unrecovered panic, but it succeeded")
	}
	t.Logf("Got expected error: %v", err)
	if !strings.Contains(err.Error(), "panic") {
		t.Errorf("error should mention panic, got: %v", err)
	}
}

// toInt64 converts an interface result to int64 for comparison.
func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int64:
		return n, true
	case int32:
		return int64(n), true
	default:
		return 0, false
	}
}

// TestStatefulGlobals_PersistAcrossRuns verifies that package-level variable
// mutations persist across multiple Run calls when WithStatefulGlobals is set.
func TestStatefulGlobals_PersistAcrossRuns(t *testing.T) {
	source := `
package main

var counter int

func init() {
	counter = 0
}

func Increment() int {
	counter++
	return counter
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	for i := int64(1); i <= 3; i++ {
		result, err := prog.Run("Increment")
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}
		got, ok := toInt64(result)
		if !ok {
			t.Fatalf("Run %d: unexpected type %T", i, result)
		}
		if got != i {
			t.Errorf("Run %d: got %d, want %d", i, got, i)
		}
	}
}

// TestStatefulGlobals_ConcurrentRuns verifies that concurrent Run calls on a
// stateful program are serialized (safe but not parallel).
func TestStatefulGlobals_ConcurrentRuns(t *testing.T) {
	source := `
package main

var counter int

func init() {
	counter = 0
}

func Increment() int {
	counter++
	return counter
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Run 100 increments concurrently
	const numGoroutines = 100
	var wg sync.WaitGroup
	results := make(chan int64, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := prog.Run("Increment")
			if err != nil {
				t.Errorf("Run failed: %v", err)
				return
			}
			got, _ := toInt64(result)
			results <- got
		}()
	}

	wg.Wait()
	close(results)

	// All results should be unique integers 1..100 (no duplicates from races)
	seen := make(map[int64]bool)
	for r := range results {
		if seen[r] {
			t.Errorf("duplicate result %d - indicates race condition", r)
		}
		seen[r] = true
	}

	// Final counter should be exactly numGoroutines
	finalResult, _ := prog.Run("Increment")
	final, _ := toInt64(finalResult)
	if final != numGoroutines+1 {
		t.Errorf("final counter = %d, want %d", final, numGoroutines+1)
	}

	// Verify all values 1..100 were seen (serialized execution)
	for i := int64(1); i <= numGoroutines; i++ {
		if !seen[i] {
			t.Errorf("missing result %d - execution was not properly serialized", i)
		}
	}
}

// TestStatefulGlobals_DefaultStatelessIsolation verifies that the default mode
// (no WithStatefulGlobals) still resets globals between calls.
func TestStatefulGlobals_DefaultStatelessIsolation(t *testing.T) {
	source := `
package main

var counter int

func init() {
	counter = 0
}

func Increment() int {
	counter++
	return counter
}
`
	prog, err := Build(source)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Each call should start from 0 and return 1
	for i := 0; i < 3; i++ {
		result, err := prog.Run("Increment")
		if err != nil {
			t.Fatalf("Run %d failed: %v", i+1, err)
		}
		got, ok := toInt64(result)
		if !ok {
			t.Fatalf("Run %d: unexpected type %T", i+1, result)
		}
		if got != 1 {
			t.Errorf("Run %d: got %d, want 1 (stateless isolation)", i+1, got)
		}
	}
}

// TestStatefulGlobals_InitSeeded verifies that init()-seeded globals are
// preserved and further mutated across calls in stateful mode.
func TestStatefulGlobals_InitSeeded(t *testing.T) {
	source := `
package main

var base int

func init() {
	base = 100
}

func AddAndGet(n int) int {
	base += n
	return base
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// First call: base starts at 100, add 5 → 105
	result, err := prog.Run("AddAndGet", 5)
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	got, _ := toInt64(result)
	if got != 105 {
		t.Errorf("Run 1: got %v, want 105", got)
	}

	// Second call: base is now 105, add 10 → 115
	result, err = prog.Run("AddAndGet", 10)
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	got, _ = toInt64(result)
	if got != 115 {
		t.Errorf("Run 2: got %v, want 115", got)
	}
}

// TestStatefulGlobals_MapCache verifies that a package-level map can serve as
// a cross-call cache in stateful mode.
func TestStatefulGlobals_MapCache(t *testing.T) {
	source := `
package main

var cache map[string]int

func init() {
	cache = make(map[string]int)
}

func GetOrSet(key string, val int) int {
	if v, ok := cache[key]; ok {
		return v
	}
	cache[key] = val
	return val
}

func CacheLen() int {
	return len(cache)
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// First call: cache miss, store "a" → 1
	result, err := prog.Run("GetOrSet", "a", 1)
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	got, _ := toInt64(result)
	if got != 1 {
		t.Errorf("Run 1: got %v, want 1", got)
	}

	// Second call: cache miss, store "b" → 2
	result, err = prog.Run("GetOrSet", "b", 2)
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	got, _ = toInt64(result)
	if got != 2 {
		t.Errorf("Run 2: got %v, want 2", got)
	}

	// Third call: cache hit for "a", should return 1
	result, err = prog.Run("GetOrSet", "a", 999)
	if err != nil {
		t.Fatalf("Run 3 failed: %v", err)
	}
	got, _ = toInt64(result)
	if got != 1 {
		t.Errorf("Run 3: got %v, want 1 (cache hit)", got)
	}

	// Verify cache length
	result, err = prog.Run("CacheLen")
	if err != nil {
		t.Fatalf("CacheLen failed: %v", err)
	}
	got, _ = toInt64(result)
	if got != 2 {
		t.Errorf("CacheLen: got %v, want 2", got)
	}
}

// TestStatefulGlobals_SeparateProgramsIsolated verifies that two separate
// Program instances with stateful globals have independent state.
func TestStatefulGlobals_SeparateProgramsIsolated(t *testing.T) {
	source := `
package main

var counter int

func init() {
	counter = 0
}

func Increment() int {
	counter++
	return counter
}
`
	prog1, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build prog1 failed: %v", err)
	}
	prog2, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build prog2 failed: %v", err)
	}

	// Increment prog1 twice
	prog1.Run("Increment")
	r1, _ := prog1.Run("Increment")
	got1, _ := toInt64(r1)
	if got1 != 2 {
		t.Errorf("prog1 Run 2: got %v, want 2", got1)
	}

	// prog2 should be independent, starting from 0
	r2, _ := prog2.Run("Increment")
	got2, _ := toInt64(r2)
	if got2 != 1 {
		t.Errorf("prog2 Run 1: got %v, want 1 (independent state)", got2)
	}
}

// TestStatefulGlobals_IncCounterSequence verifies that the IncCounter pattern
// from the initialize testdata works correctly with stateful globals.
// This test is isolated from the main correctness suite because it requires
// stateful globals mode.
func TestStatefulGlobals_IncCounterSequence(t *testing.T) {
	// Use the same source as tests/testdata/initialize/main.go IncCounter functions
	// Note: init() is required to properly initialize the counter to a valid int value
	source := `
package main

var counter int

func init() {
	counter = 0
}

func IncCounter1() int {
	counter++
	return counter
}

func IncCounter2() int {
	counter++
	return counter
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// First call: counter 0 → 1
	r1, err := prog.Run("IncCounter1")
	if err != nil {
		t.Fatalf("IncCounter1 failed: %v", err)
	}
	got1, _ := toInt64(r1)
	if got1 != 1 {
		t.Errorf("IncCounter1: got %d, want 1", got1)
	}

	// Second call: counter 1 → 2
	r2, err := prog.Run("IncCounter2")
	if err != nil {
		t.Fatalf("IncCounter2 failed: %v", err)
	}
	got2, _ := toInt64(r2)
	if got2 != 2 {
		t.Errorf("IncCounter2: got %d, want 2", got2)
	}

	// Third call: counter 2 → 3
	r3, err := prog.Run("IncCounter1")
	if err != nil {
		t.Fatalf("IncCounter1 (3rd) failed: %v", err)
	}
	got3, _ := toInt64(r3)
	if got3 != 3 {
		t.Errorf("IncCounter1 (3rd): got %d, want 3", got3)
	}
}

// TestStatefulGlobals_InitSeededCounter verifies that a counter initialized
// via init() works correctly with stateful globals.
func TestStatefulGlobals_InitSeededCounter(t *testing.T) {
	source := `
package main

var counter int

func init() {
	counter = 100
}

func IncCounter() int {
	counter++
	return counter
}

func GetCounter() int {
	return counter
}
`
	prog, err := Build(source, WithStatefulGlobals())
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// After init, counter is 100
	r0, _ := prog.Run("GetCounter")
	got0, _ := toInt64(r0)
	if got0 != 100 {
		t.Errorf("GetCounter after init: got %d, want 100", got0)
	}

	// Increment: 100 → 101
	r1, _ := prog.Run("IncCounter")
	got1, _ := toInt64(r1)
	if got1 != 101 {
		t.Errorf("IncCounter: got %d, want 101", got1)
	}

	// Another increment: 101 → 102
	r2, _ := prog.Run("IncCounter")
	got2, _ := toInt64(r2)
	if got2 != 102 {
		t.Errorf("IncCounter (2nd): got %d, want 102", got2)
	}

	// Verify via GetCounter
	r3, _ := prog.Run("GetCounter")
	got3, _ := toInt64(r3)
	if got3 != 102 {
		t.Errorf("GetCounter after increments: got %d, want 102", got3)
	}
}
