// resolve.go contains the small handful of name-resolution helpers
// the v2 generator needs: a reflect.Type expression for named types,
// a constant cast helper for unsigned overflow, and an import-path
// sanitizer.
package gentool

import (
	"fmt"
	"go/constant"
	"go/types"
	"strings"
)

// needsUintCast reports whether a constant value would overflow when
// passed to AddConstant as an untyped int. Untyped/uint values larger
// than math.MaxInt64 must be wrapped in uint64() to compile.
func needsUintCast(c *types.Const) bool {
	basic, ok := c.Type().Underlying().(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint, types.Uint64, types.Uintptr, types.UntypedInt:
		val := c.Val()
		if val.Kind() == constant.Int {
			if v, ok := constant.Uint64Val(val); ok {
				return v > (1<<63 - 1) // > math.MaxInt64
			}
		}
	}
	return false
}

// typeToReflectExpr returns a Go expression that evaluates to the
// reflect.Type for a named type t. Returns "" for unnamed types.
func typeToReflectExpr(t types.Type, pkgRef string) string {
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	name := named.Obj().Name()
	switch t.Underlying().(type) {
	case *types.Interface:
		return fmt.Sprintf("reflect.TypeOf((*%s.%s)(nil)).Elem()", pkgRef, name)
	case *types.Struct:
		return fmt.Sprintf("reflect.TypeOf(%s.%s{})", pkgRef, name)
	default:
		return fmt.Sprintf("reflect.TypeOf((*%s.%s)(nil)).Elem()", pkgRef, name)
	}
}

// sanitizePkgName turns a Go import path into an identifier we can use
// as a Go variable / file name (e.g. "encoding/json" -> "encoding_json").
func sanitizePkgName(path string) string {
	return strings.NewReplacer(
		"/", "_",
		"-", "_",
		".", "_",
	).Replace(path)
}
