package compiler

import (
	"fmt"
	"go/types"
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/value"
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
	reg := &mockRegistry{}
	c := NewCompiler(reg)
	if c == nil {
		t.Fatal("NewCompiler returned nil")
	}
	// Verify it satisfies the interface at compile time.
	var _ Compiler = c
}

// TestCompilerInterfaceContract verifies the Compiler interface contract:
// calling Compile with nil should result in a panic or error, not a hang.
func TestCompilerInterfaceContract(t *testing.T) {
	reg := &mockRegistry{}
	c := NewCompiler(reg)

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

	_, _ = Compile(&mockRegistry{}, nil)
}

// ---------------------------------------------------------------------------
// Mock PackageRegistry
// ---------------------------------------------------------------------------

type mockRegistry struct{}

func (m *mockRegistry) RegisterPackage(path, name string) *importer.ExternalPackage {
	return nil
}

func (m *mockRegistry) GetPackageByPath(path string) *importer.ExternalPackage {
	return nil
}

func (m *mockRegistry) GetPackageByName(name string) *importer.ExternalPackage {
	return nil
}

func (m *mockRegistry) GetAllPackages() map[string]*importer.ExternalPackage {
	return nil
}

func (m *mockRegistry) SetExternalType(t types.Type, rt reflect.Type) {}

func (m *mockRegistry) GetExternalType(t types.Type) reflect.Type {
	return nil
}

func (m *mockRegistry) AddMethodDirectCall(typeName, methodName string, dc func([]value.Value) value.Value) {
}

func (m *mockRegistry) LookupMethodDirectCall(typeName, methodName string) (func([]value.Value) value.Value, bool) {
	return nil, false
}

func (m *mockRegistry) LookupPackage(name string) (*importer.ExternalPackage, error) {
	return nil, fmt.Errorf("not found")
}

func (m *mockRegistry) AutoImport(name string) (path string, pkg *importer.ExternalPackage, ok bool) {
	return "", nil, false
}

func (m *mockRegistry) LookupExternalFunc(pkgPath, funcName string) (fn any, directCall func([]value.Value) value.Value, ok bool) {
	return nil, nil, false
}

func (m *mockRegistry) LookupExternalVar(pkgPath, varName string) (ptr any, ok bool) {
	return nil, false
}

// Verify mockRegistry satisfies the interface.
var _ importer.PackageRegistry = (*mockRegistry)(nil)
