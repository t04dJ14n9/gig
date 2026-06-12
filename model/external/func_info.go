// Package external defines cross-cutting external types shared across the Gig
// interpreter's compilation and execution pipeline.
//
// This package is a leaf: it imports only model/value/ and standard library packages.
// It does NOT import any other gig internal packages to avoid circular dependencies.
package external

import "github.com/t04dJ14n9/gig/model/value"

// ExternalFuncInfo contains pre-resolved external function info for fast calls.
// This allows the VM to bypass reflection when calling external functions.
type ExternalFuncInfo struct {
	// PkgPath is the Go import path for the package that owns Func.
	PkgPath string

	// FuncName is the exported function name used for diagnostics.
	FuncName string

	// IsStdlib records whether PkgPath belongs to the Go standard library.
	// The compiler computes this once so each VM call avoids reparsing PkgPath
	// while still enforcing third-party boundary checks for non-stdlib calls.
	IsStdlib bool

	// Func is the actual function value.
	Func any

	// DirectCall is a typed wrapper that avoids reflect.Call.
	DirectCall func(args []value.Value) value.Value

	// IsVariadic indicates whether the function takes variadic arguments.
	// Pre-computed at compile time so the VM avoids reflect.Type queries.
	IsVariadic bool

	// NumIn is the number of declared parameters (including the variadic slice type).
	// Pre-computed at compile time so the VM avoids reflect.Type queries.
	NumIn int
}
