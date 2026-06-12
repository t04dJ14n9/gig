package gentool

import (
	"fmt"
	"go/constant"
	"go/types"
)

// needsUintCast checks whether a constant needs a uint64() cast to prevent
// overflow when passed as an untyped int to AddConstant.
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
				return v > (1<<63 - 1)
			}
		}
	}
	return false
}

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
