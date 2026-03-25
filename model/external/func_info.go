// Package external defines cross-cutting external types shared across the Gig
// interpreter's compilation and execution pipeline.
//
// This package is a leaf: it imports only model/value/ and standard library packages.
// It does NOT import any other gig internal packages to avoid circular dependencies.
package external

import "git.woa.com/youngjin/gig/model/value"

// ExternalFuncInfo contains pre-resolved external function info for fast calls.
// This allows the VM to bypass reflection when calling external functions.
type ExternalFuncInfo struct {
	// Func is the actual function value.
	Func any

	// DirectCall is a typed wrapper that avoids reflect.Call.
	DirectCall func(args []value.Value) value.Value
}
