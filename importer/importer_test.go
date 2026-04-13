package importer

import (
	"fmt"
	"reflect"
	"testing"

	"git.woa.com/youngjin/gig/model/external"
)

// ---------------------------------------------------------------------------
// Package Registration
// ---------------------------------------------------------------------------

// TestRegisterAndLookupPackage verifies the complete lifecycle of registering
// a package and looking it up by path and by name.
func TestRegisterAndLookupPackage(t *testing.T) {
	path := "test/mypkg"
	name := "mypkg"
	pkg := RegisterPackage(path, name)

	if pkg.Path != path {
		t.Errorf("Path = %q, want %q", pkg.Path, path)
	}
	if pkg.Name != name {
		t.Errorf("Name = %q, want %q", pkg.Name, name)
	}

	// Lookup by path.
	byPath := GetPackageByPath(path)
	if byPath != pkg {
		t.Error("GetPackageByPath did not return the registered package")
	}

	// Lookup by name.
	byName := GetPackageByName(name)
	if byName != pkg {
		t.Error("GetPackageByName did not return the registered package")
	}
}

// TestGetPackageNotFound verifies nil is returned for unregistered packages.
func TestGetPackageNotFound(t *testing.T) {
	if got := GetPackageByPath("no/such/pkg"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
	if got := GetPackageByName("nosuchpkg"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// TestGetAllPackages verifies the returned map is a copy.
func TestGetAllPackages(t *testing.T) {
	// Register a test package to ensure at least one exists.
	RegisterPackage("test/allpkgs", "allpkgs")

	all := GetAllPackages()
	if len(all) == 0 {
		t.Fatal("GetAllPackages returned empty map")
	}

	// Mutating the returned map should not affect the registry.
	sizeBefore := len(all)
	all["fake"] = &ExternalPackage{}
	allAgain := GetAllPackages()
	if len(allAgain) != sizeBefore {
		t.Error("mutation of returned map affected the registry")
	}
}

// ---------------------------------------------------------------------------
// AddFunction / AddConstant / AddVariable / AddType
// ---------------------------------------------------------------------------

// TestAddFunction verifies adding a function to a package.
func TestAddFunction(t *testing.T) {
	pkg := RegisterPackage("test/funcpkg", "funcpkg")

	fn := func(a, b int) int { return a + b }
	pkg.AddFunction("Add", fn, "adds two ints", nil)

	obj, ok := pkg.Objects["Add"]
	if !ok {
		t.Fatal("Add not found in Objects")
	}
	if obj.Kind != external.ObjectKindFunction {
		t.Errorf("Kind = %d, want ObjectKindFunction", obj.Kind)
	}
	if obj.Name != "Add" {
		t.Errorf("Name = %q", obj.Name)
	}
	if obj.Doc != "adds two ints" {
		t.Errorf("Doc = %q", obj.Doc)
	}
	if obj.Type == nil {
		t.Error("Type should not be nil")
	}

	// The underlying function should be callable.
	result := obj.Value.(func(int, int) int)(3, 4)
	if result != 7 {
		t.Errorf("Add(3,4) = %d, want 7", result)
	}
}

// TestAddConstant verifies adding a constant.
func TestAddConstant(t *testing.T) {
	pkg := RegisterPackage("test/constpkg", "constpkg")
	pkg.AddConstant("Pi", 3.14, "approximate pi")

	obj := pkg.Objects["Pi"]
	if obj == nil {
		t.Fatal("Pi not found")
	}
	if obj.Kind != external.ObjectKindConstant {
		t.Errorf("Kind = %d", obj.Kind)
	}
	if obj.Value.(float64) != 3.14 {
		t.Errorf("Value = %v", obj.Value)
	}
}

// TestAddVariable verifies adding a variable.
func TestAddVariable(t *testing.T) {
	pkg := RegisterPackage("test/varpkg", "varpkg")
	x := 42
	pkg.AddVariable("X", &x, "a variable")

	obj := pkg.Objects["X"]
	if obj == nil {
		t.Fatal("X not found")
	}
	if obj.Kind != external.ObjectKindVariable {
		t.Errorf("Kind = %d", obj.Kind)
	}
}

// TestAddType verifies adding a named type.
func TestAddType(t *testing.T) {
	pkg := RegisterPackage("test/typepkg", "typepkg")
	type MyStruct struct{ X int }
	pkg.AddType("MyStruct", reflect.TypeOf(MyStruct{}), "a struct type")

	obj := pkg.Objects["MyStruct"]
	if obj == nil {
		t.Fatal("MyStruct not found")
	}
	if obj.Kind != external.ObjectKindType {
		t.Errorf("Kind = %d", obj.Kind)
	}
	rt, ok := pkg.Types["MyStruct"]
	if !ok {
		t.Fatal("type not in Types map")
	}
	if rt.Name() != "MyStruct" {
		t.Errorf("type name = %q", rt.Name())
	}
}

// TestAddTypeNil verifies that nil types are silently skipped.
func TestAddTypeNil(t *testing.T) {
	pkg := RegisterPackage("test/niltypepkg", "niltypepkg")
	pkg.AddType("Nil", nil, "")
	if _, ok := pkg.Objects["Nil"]; ok {
		t.Error("nil type should not create an object")
	}
}

// ---------------------------------------------------------------------------
// ExternalType registry
// ---------------------------------------------------------------------------

// TestSetGetExternalType verifies the types.Type <-> reflect.Type mapping.
func TestSetGetExternalType(t *testing.T) {
	reg := globalRegistry
	pkg := RegisterPackage("test/exttype", "exttype")
	pkg.AddFunction("Sprintf", fmt.Sprintf, "", nil)
	obj := pkg.Objects["Sprintf"]
	if obj.Type == nil {
		t.Skip("no types.Type available for test")
	}

	rt := reflect.TypeOf(fmt.Sprintf)
	reg.SetExternalType(obj.Type, rt)

	got, ok := reg.LookupExternalType(obj.Type)
	if !ok || got != rt {
		t.Errorf("LookupExternalType returned %v, want %v", got, rt)
	}
}

// TestGetExternalTypeNotFound verifies nil is returned for unmapped types.
func TestGetExternalTypeNotFound(t *testing.T) {
	got, ok := globalRegistry.LookupExternalType(nil)
	if ok || got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// ---------------------------------------------------------------------------
// Importer (types.Importer interface)
// ---------------------------------------------------------------------------

// TestNewImporter verifies that NewImporter creates a valid Importer.
func TestNewImporter(t *testing.T) {
	imp := NewImporter(GlobalRegistry())
	if imp == nil {
		t.Fatal("NewImporter returned nil")
	}
}

// TestImportRegisteredPackage verifies that Import can resolve a registered package.
func TestImportRegisteredPackage(t *testing.T) {
	path := "test/importable"
	pkg := RegisterPackage(path, "importable")
	pkg.AddFunction("Hello", func() string { return "world" }, "", nil)

	imp := NewImporter(GlobalRegistry())
	typesPkg, err := imp.Import(path)
	if err != nil {
		t.Fatalf("Import(%q) error: %v", path, err)
	}
	if typesPkg == nil {
		t.Fatal("Import returned nil package")
	}
	if typesPkg.Path() != path {
		t.Errorf("Path() = %q, want %q", typesPkg.Path(), path)
	}
	if typesPkg.Name() != "importable" {
		t.Errorf("Name() = %q, want %q", typesPkg.Name(), "importable")
	}
}

// TestImportCaching verifies that repeated imports return the same package.
func TestImportCaching(t *testing.T) {
	path := "test/cached"
	RegisterPackage(path, "cached")

	imp := NewImporter(GlobalRegistry())
	p1, err1 := imp.Import(path)
	p2, err2 := imp.Import(path)
	if err1 != nil || err2 != nil {
		t.Fatalf("errors: %v, %v", err1, err2)
	}
	if p1 != p2 {
		t.Error("repeated Import should return the same *types.Package")
	}
}

// TestImportUnregistered verifies error for unregistered package.
func TestImportUnregistered(t *testing.T) {
	imp := NewImporter(GlobalRegistry())
	_, err := imp.Import("no/such/package/ever")
	if err == nil {
		t.Error("expected error for unregistered package")
	}
}

// ---------------------------------------------------------------------------
// LookupPackage (Registry method)
// ---------------------------------------------------------------------------

// TestLookupPackage verifies lookup by path and name.
func TestLookupPackage(t *testing.T) {
	path := "test/lookuppkg"
	name := "lookuppkg"
	RegisterPackage(path, name)

	reg := GlobalRegistry()

	// Lookup by path.
	pkg, err := reg.LookupPackage(path)
	if err != nil || pkg == nil {
		t.Fatalf("LookupPackage(%q) failed: %v", path, err)
	}

	// Lookup by name.
	pkg2, err := reg.LookupPackage(name)
	if err != nil || pkg2 == nil {
		t.Fatalf("LookupPackage(%q) failed: %v", name, err)
	}
}

// TestLookupPackageNotFound verifies error for missing package.
func TestLookupPackageNotFound(t *testing.T) {
	_, err := GlobalRegistry().LookupPackage("nonexistent_pkg_xyz")
	if err == nil {
		t.Error("expected error for nonexistent package")
	}
}

// ---------------------------------------------------------------------------
// AutoImport (Registry method)
// ---------------------------------------------------------------------------

// TestAutoImport verifies automatic import resolution by name.
func TestAutoImport(t *testing.T) {
	path := "test/autoimp"
	name := "autoimp"
	RegisterPackage(path, name)

	gotPath, pkg, ok := GlobalRegistry().AutoImport(name)
	if !ok {
		t.Fatal("AutoImport should succeed for registered package")
	}
	if gotPath != path {
		t.Errorf("path = %q, want %q", gotPath, path)
	}
	if pkg == nil {
		t.Error("pkg should not be nil")
	}
}

// TestAutoImportNotFound verifies false for missing package.
func TestAutoImportNotFound(t *testing.T) {
	_, _, ok := GlobalRegistry().AutoImport("no_such_auto_import_xyz")
	if ok {
		t.Error("AutoImport should return false for unregistered package")
	}
}
