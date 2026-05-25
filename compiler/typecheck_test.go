package compiler

import (
	"go/types"
	"testing"
)

func TestIsStdlibPath(t *testing.T) {
	tests := []struct {
		path   string
		expect bool
	}{
		{"fmt", true},
		{"encoding/json", true},
		{"sort", true},
		{"io", true},
		{"net/http", true},
		{"golang.org/x/tools", false},
		{"github.com/foo/bar", false},
		{"github.com/t04dJ14n9/gig", false},
		{"command-line-arguments", false},
		{"main", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isStdlibPath(tt.path)
			if got != tt.expect {
				t.Errorf("isStdlibPath(%q) = %v, want %v", tt.path, got, tt.expect)
			}
		})
	}
}

func TestIsUserDefinedNamedType(t *testing.T) {
	userPkg := types.NewPackage("command-line-arguments", "main")
	userType := types.NewNamed(
		types.NewTypeName(0, userPkg, "MySlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	sortPkg := types.NewPackage("sort", "sort")
	sortType := types.NewNamed(
		types.NewTypeName(0, sortPkg, "IntSlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	intType := types.Typ[types.Int]
	stringType := types.Typ[types.String]

	tests := []struct {
		name   string
		typ    types.Type
		expect bool
	}{
		{"user named type", userType, true},
		{"user pointer to named type", types.NewPointer(userType), true},
		{"stdlib named type", sortType, false},
		{"bare int", intType, false},
		{"bare string", stringType, false},
		{"slice of int", types.NewSlice(intType), false},
		{"nil type", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserDefinedNamedType(tt.typ)
			if got != tt.expect {
				t.Errorf("isUserDefinedNamedType(%v) = %v, want %v", tt.typ, got, tt.expect)
			}
		})
	}
}

func TestRejectUserTypeToThirdParty(t *testing.T) {
	userPkg := types.NewPackage("command-line-arguments", "main")
	userType := types.NewNamed(
		types.NewTypeName(0, userPkg, "MySlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	sortPkg := types.NewPackage("sort", "sort")
	sortType := types.NewNamed(
		types.NewTypeName(0, sortPkg, "IntSlice", nil),
		types.NewSlice(types.Typ[types.Int]),
		nil,
	)

	tests := []struct {
		name      string
		pkgPath   string
		argTypes  []types.Type
		expectErr bool
	}{
		{"user type to stdlib is allowed", "sort", []types.Type{userType}, false},
		{"user type to third-party is rejected", "github.com/foo/bar", []types.Type{userType}, true},
		{"stdlib type to third-party is allowed", "github.com/foo/bar", []types.Type{sortType}, false},
		{"primitive to third-party is allowed", "github.com/foo/bar", []types.Type{types.Typ[types.Int]}, false},
		{"slice of user type to third-party is rejected", "github.com/foo/bar", []types.Type{types.NewSlice(userType)}, true},
		{"pointer to user type to third-party is rejected", "github.com/foo/bar", []types.Type{types.NewPointer(userType)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExternalCallArgs(tt.pkgPath, "SomeFunc", tt.argTypes)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
