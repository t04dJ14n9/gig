package gig

import (
	"strings"
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
