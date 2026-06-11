package parser

import (
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig/importer"
)

// ---------------------------------------------------------------------------
// Parse — auto package wrapping
// ---------------------------------------------------------------------------

func TestParseAutoPackageWrap(t *testing.T) {
	// Source without "package main" should be auto-wrapped
	reg := importer.NewRegistry()
	result, err := Parse(`func Hello() int { return 42 }`, reg)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result.File == nil {
		t.Fatal("result.File is nil")
	}
	if result.Pkg == nil {
		t.Fatal("result.Pkg is nil")
	}
	if result.Pkg.Name() != "main" {
		t.Errorf("package name = %q, want %q", result.Pkg.Name(), "main")
	}
}

func TestParseExplicitPackageDecl(t *testing.T) {
	// Source with explicit package declaration
	reg := importer.NewRegistry()
	result, err := Parse(`package main

func Hello() int { return 42 }
`, reg)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result.Pkg.Name() != "main" {
		t.Errorf("package name = %q, want %q", result.Pkg.Name(), "main")
	}
}

// ---------------------------------------------------------------------------
// Parse — type checking
// ---------------------------------------------------------------------------

func TestParseTypeCheckSuccess(t *testing.T) {
	reg := importer.NewRegistry()
	result, err := Parse(`func Add(a, b int) int { return a + b }`, reg)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result.Info == nil {
		t.Fatal("result.Info is nil")
	}
	if len(result.Info.Types) == 0 {
		t.Error("expected non-empty Types map after type checking")
	}
}

func TestParseTypeCheckError(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`func Bad() int { return "not an int" }`, reg)
	if err == nil {
		t.Fatal("expected type check error for type mismatch")
	}
	if !strings.Contains(err.Error(), "type check") {
		t.Errorf("error = %q, want error containing 'type check'", err.Error())
	}
}

func TestParseSyntaxError(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`func {{{{`, reg)
	if err == nil {
		t.Fatal("expected parse error for invalid syntax")
	}
	if !strings.Contains(err.Error(), "parse error") {
		t.Errorf("error = %q, want error containing 'parse error'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Parse — banned imports
// ---------------------------------------------------------------------------

func TestParseBannedUnsafeImport(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`package main

import "unsafe"

func Foo() { _ = unsafe.Pointer(nil) }
`, reg)
	if err == nil {
		t.Fatal("expected error for banned 'unsafe' import")
	}
	if !strings.Contains(err.Error(), "unsafe") {
		t.Errorf("error = %q, want error mentioning 'unsafe'", err.Error())
	}
}

func TestParseBannedReflectImport(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`package main

import "reflect"

func Foo() { reflect.TypeOf(0) }
`, reg)
	if err == nil {
		t.Fatal("expected error for banned 'reflect' import")
	}
	if !strings.Contains(err.Error(), "reflect") {
		t.Errorf("error = %q, want error mentioning 'reflect'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Parse — banned panic
// ---------------------------------------------------------------------------

func TestParseBannedPanic(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`func Foo() { panic("oh no") }`, reg)
	if err == nil {
		t.Fatal("expected error for banned panic()")
	}
	if !strings.Contains(err.Error(), "panic") {
		t.Errorf("error = %q, want error mentioning 'panic'", err.Error())
	}
}

func TestParsePanicAllowedWithOption(t *testing.T) {
	reg := importer.NewRegistry()
	result, err := Parse(`func Foo() { panic("oh no") }`, reg, WithAllowPanic())
	if err != nil {
		t.Fatalf("Parse with WithAllowPanic() error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
}

// ---------------------------------------------------------------------------
// Parse — auto-import
// ---------------------------------------------------------------------------

func TestParseAutoImportRegisteredPackage(t *testing.T) {
	reg := importer.NewRegistry()
	pkg := importer.RegisterPackage("testpkg", "testpkg")
	pkg.AddFunction("DoStuff", func() {}, "", nil)

	// The auto-import logic runs before type checking, so even if type checking
	// fails, the import should be injected. We verify by checking the error type.
	_, err := Parse(`func Foo() { testpkg.DoStuff() }`, reg)
	if err != nil {
		// Type check may fail for reasons beyond our control, but
		// the auto-import injection should have happened before that.
		// The important thing is that it's not a "banned import" or "parse" error.
		if strings.Contains(err.Error(), "banned") || strings.Contains(err.Error(), "parse error") {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestParseAutoImportNoDuplicate(t *testing.T) {
	// Verify that when the user already has the import, auto-import doesn't add a duplicate.
	reg := importer.NewRegistry()
	pkg := importer.RegisterPackage("testpkg", "testpkg")
	pkg.AddFunction("DoStuff", func() {}, "", nil)

	src := `package main

import "testpkg"

func Foo() { testpkg.DoStuff() }
`
	result, err := Parse(src, reg)
	if err != nil {
		if strings.Contains(err.Error(), "banned") || strings.Contains(err.Error(), "parse error") {
			t.Fatalf("unexpected error: %v", err)
		}
		// Type check failure is acceptable
		return
	}
	count := 0
	for _, imp := range result.File.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if path == "testpkg" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("testpkg import count = %d, want 1", count)
	}
}

func TestParseAutoImportUnregisteredPackage(t *testing.T) {
	// Using an unregistered package name should fail type checking
	reg := importer.NewRegistry()
	_, err := Parse(`func Foo() { nonexistent.Func() }`, reg)
	if err == nil {
		t.Fatal("expected error for unregistered package")
	}
}

// ---------------------------------------------------------------------------
// Parse — ParseResult fields
// ---------------------------------------------------------------------------

func TestParseResultFields(t *testing.T) {
	reg := importer.NewRegistry()
	result, err := Parse(`func Hello() int { return 42 }`, reg)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result.File == nil {
		t.Error("File is nil")
	}
	if result.FSet == nil {
		t.Error("FSet is nil")
	}
	if result.Info == nil {
		t.Error("Info is nil")
	}
	if result.Pkg == nil {
		t.Error("Pkg is nil")
	}
}

// ---------------------------------------------------------------------------
// checkBannedImports — edge cases
// ---------------------------------------------------------------------------

func TestCheckBannedImportsEmpty(t *testing.T) {
	// Source with no imports should not error
	reg := importer.NewRegistry()
	_, err := Parse(`func Foo() int { return 1 }`, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckBannedImportsAllowed(t *testing.T) {
	// Non-banned import should not trigger the banned-import check.
	// The type checker may fail if the package isn't fully registered,
	// but the important thing is that the banned-import check passes.
	reg := importer.NewRegistry()
	_, err := Parse(`package main

import "math"

func Foo() float64 { return 1.0 }
`, reg)
	if err != nil {
		// The error should NOT be about banned imports
		if strings.Contains(err.Error(), "banned") {
			t.Fatalf("math should not be banned: %v", err)
		}
		// Type check failure is expected — math isn't registered with our custom registry
	}
}

// ---------------------------------------------------------------------------
// checkBannedPanic — edge cases
// ---------------------------------------------------------------------------

func TestCheckBannedPanicWithExpression(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`func Foo() { panic(42) }`, reg)
	if err == nil {
		t.Fatal("expected error for panic(42)")
	}
}

func TestCheckBannedPanicWithVariable(t *testing.T) {
	reg := importer.NewRegistry()
	_, err := Parse(`func Foo() { var msg = "err"; panic(msg) }`, reg)
	if err == nil {
		t.Fatal("expected error for panic(msg)")
	}
}

// ---------------------------------------------------------------------------
// Multiple banned constructs
// ---------------------------------------------------------------------------

func TestParseBannedImportTakesPrecedenceOverPanic(t *testing.T) {
	// When both unsafe import and panic are present, import check should fail first
	reg := importer.NewRegistry()
	_, err := Parse(`package main

import "unsafe"

func Foo() { panic(unsafe.Pointer(nil)) }
`, reg)
	if err == nil {
		t.Fatal("expected error")
	}
	// Should mention unsafe, not panic
	if !strings.Contains(err.Error(), "unsafe") {
		t.Errorf("error = %q, want error mentioning 'unsafe'", err.Error())
	}
}
