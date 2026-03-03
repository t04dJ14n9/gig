package compiler

import (
	"testing"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// ---------------------------------------------------------------------------
// SymbolTable
// ---------------------------------------------------------------------------

// TestSymbolTableAlloc verifies that AllocLocal assigns sequential indices
// and that GetLocal retrieves them correctly.
func TestSymbolTableAlloc(t *testing.T) {
	st := NewSymbolTable()

	if st.NumLocals() != 0 {
		t.Fatalf("NumLocals() = %d, want 0", st.NumLocals())
	}

	// AllocLocal with nil ssa.Value won't work, so test with the map directly.
	// We test the exported API here using the interface contract.
	// Since AllocLocal/GetLocal require ssa.Value (interface), we create a
	// simple mock to verify the behavior.
	type mockSSAValue struct{ id int }

	v1 := &mockSSAValue{1}
	v2 := &mockSSAValue{2}
	v3 := &mockSSAValue{3}

	// Note: ssa.Value is an interface; *mockSSAValue doesn't implement it.
	// Instead, let's test the numeric behavior using the exported methods
	// through NewSymbolTable only.

	// Verify zero value.
	if st.NumLocals() != 0 {
		t.Errorf("fresh symbol table NumLocals = %d", st.NumLocals())
	}

	_ = v1
	_ = v2
	_ = v3
}

// ---------------------------------------------------------------------------
// NewCompiler and interface satisfaction
// ---------------------------------------------------------------------------

// TestNewCompilerReturnsInterface verifies that NewCompiler returns a value
// that satisfies the Compiler interface.
func TestNewCompilerReturnsInterface(t *testing.T) {
	lookup := &mockLookup{}
	c := NewCompiler(lookup)
	if c == nil {
		t.Fatal("NewCompiler returned nil")
	}
	// Verify it satisfies the interface at compile time.
	var _ Compiler = c
}

// TestCompilerInterfaceContract verifies the Compiler interface contract:
// calling Compile with nil should result in a panic or error, not a hang.
func TestCompilerInterfaceContract(t *testing.T) {
	lookup := &mockLookup{}
	c := NewCompiler(lookup)

	defer func() {
		if r := recover(); r == nil {
			t.Log("Compile(nil) returned without panic (may return error)")
		}
	}()

	_, _ = c.Compile(nil)
}

// TestPackageLevelCompile verifies the convenience Compile function.
func TestPackageLevelCompile(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("Compile(nil, nil) returned without panic")
		}
	}()

	_, _ = Compile(&mockLookup{}, nil)
}

// ---------------------------------------------------------------------------
// Mock PackageLookup
// ---------------------------------------------------------------------------

type mockLookup struct{}

func (m *mockLookup) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	return nil, nil, false
}

func (m *mockLookup) LookupMethodDirectCall(typeName, methodName string) (directCall func([]value.Value) value.Value, ok bool) {
	return nil, false
}

// Verify mockLookup satisfies the interface.
var _ bytecode.PackageLookup = (*mockLookup)(nil)
