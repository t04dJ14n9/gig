package compiler

import (
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileConst compiles a constant value.
func (c *compiler) compileConst(cnst *ssa.Const) {
	// Constants are stored once in the bytecode pool and loaded by OpConst.
	// The value placed in the pool must preserve Go type information for typed
	// nils and named zero-size structs, because the VM cannot recover that from
	// an untyped nil later.
	v := ssaConstValue(cnst)
	idx := c.addConstant(v)
	c.emit(bytecode.OpConst, idx)
}

func ssaConstValue(cnst *ssa.Const) any {
	switch t := cnst.Type().(type) {
	case *types.Basic:
		return basicConstValue(t.Kind(), cnst.Value)
	case *types.Named, *types.Alias:
		return namedConstValue(t, cnst.Value)
	case *types.Struct:
		return typedNilConstValue(t, cnst.Value)
	default:
		return typedNilConstValue(cnst.Type(), cnst.Value)
	}
}

func namedConstValue(t types.Type, val constant.Value) any {
	// Non-nil named constants compile as their underlying basic value. Nil named
	// constants need reflect.Zero so map/slice/func/pointer/interface nils keep
	// the declared type through the VM constant pool.
	if val != nil {
		if underlying, ok := t.Underlying().(*types.Basic); ok {
			return basicConstValue(underlying.Kind(), val)
		}
		return nil
	}
	return typedNilConstValue(t, val)
}
