package compiler

import (
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// builtinOps maps simple builtin names to their single-opcode implementations.
// Builtins not in this map have custom compilation logic in compileBuiltinCall.
var builtinOps = map[string]bytecode.OpCode{
	"len":     bytecode.OpLen,
	"cap":     bytecode.OpCap,
	"copy":    bytecode.OpCopy,
	"panic":   bytecode.OpPanic,
	"recover": bytecode.OpRecover,
	"close":   bytecode.OpClose,
}

// compileBuiltinCall compiles a call to a builtin function.
func (c *compiler) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	// Fast path: single-opcode builtins with 1 arg that push a value
	if op, ok := builtinOps[name]; ok {
		switch name {
		case "len", "cap":
			c.compileValue(args[0])
			c.emit(op)
		case "copy":
			c.compileValue(args[0])
			c.compileValue(args[1])
			c.emit(op)
		case "panic":
			c.compileValue(args[0])
			c.emit(op)
			return
		case "recover":
			// recover() takes no arguments
			c.emit(op)
		case "close":
			c.compileValue(args[0])
			c.emit(op)
			return
		}
		c.emit(bytecode.OpSetLocal, uint16(resultIdx))
		return
	}

	switch name {
	case "append":
		c.compileAppendBuiltin(args)
	case "delete":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(bytecode.OpDelete)
		return
	case "print":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(bytecode.OpPrint, uint16(len(args)))
		return
	case "println":
		for _, arg := range args {
			c.compileValue(arg)
		}
		c.emit(bytecode.OpPrintln, uint16(len(args)))
		return
	case "new":
		typeIdx := c.addType(args[0].Type())
		c.emit(bytecode.OpNew, uint16(typeIdx))
	case "make":
		c.compileMakeBuiltin(args)
	case "ssa:wrapnilchk":
		c.compileValue(args[0])
	case "real":
		// real(complex) -> float
		c.compileValue(args[0])
		c.emit(bytecode.OpReal)
	case "imag":
		// imag(complex) -> float
		c.compileValue(args[0])
		c.emit(bytecode.OpImag)
	case "complex":
		// complex(real, imag) -> complex
		c.compileValue(args[0]) // real part
		c.compileValue(args[1]) // imag part
		c.emit(bytecode.OpComplex)
	default:
		c.emit(bytecode.OpNil)
	}

	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

func (c *compiler) compileAppendBuiltin(args []ssa.Value) {
	if len(args) == 0 {
		c.emit(bytecode.OpNil)
		return
	}
	c.compileValue(args[0])

	if len(args) == 2 && !isNilConst(args[0]) {
		if packed, ok := packedVarargsValues(args[1]); ok {
			for _, arg := range packed {
				c.compileValue(arg)
				c.emit(bytecode.OpAppend)
			}
			return
		}
	}

	for _, arg := range args[1:] {
		c.compileValue(arg)
		c.emit(bytecode.OpAppend)
	}
}

func isNilConst(v ssa.Value) bool {
	c, ok := v.(*ssa.Const)
	return ok && c.Value == nil
}

func packedVarargsValues(v ssa.Value) ([]ssa.Value, bool) {
	slice, ok := v.(*ssa.Slice)
	if !ok || slice.Low != nil || slice.High != nil || slice.Max != nil {
		return nil, false
	}
	alloc, ok := slice.X.(*ssa.Alloc)
	if !ok || alloc.Comment != "varargs" {
		return nil, false
	}
	ptr, ok := alloc.Type().Underlying().(*types.Pointer)
	if !ok {
		return nil, false
	}
	arr, ok := ptr.Elem().Underlying().(*types.Array)
	if !ok || arr.Len() < 0 {
		return nil, false
	}

	values := make([]ssa.Value, int(arr.Len()))
	refs := alloc.Referrers()
	if refs == nil {
		return nil, false
	}
	for _, ref := range *refs {
		indexAddr, ok := ref.(*ssa.IndexAddr)
		if !ok || indexAddr.X != alloc {
			continue
		}
		idxConst, ok := indexAddr.Index.(*ssa.Const)
		if !ok || idxConst.Value == nil {
			continue
		}
		idx, exact := constant.Int64Val(idxConst.Value)
		if !exact || idx < 0 || idx >= int64(len(values)) {
			continue
		}
		indexRefs := indexAddr.Referrers()
		if indexRefs == nil {
			continue
		}
		for _, indexRef := range *indexRefs {
			store, ok := indexRef.(*ssa.Store)
			if ok && store.Addr == indexAddr {
				values[idx] = store.Val
				break
			}
		}
	}
	for _, val := range values {
		if val == nil {
			return nil, false
		}
	}
	return values, true
}

// compileMakeBuiltin compiles the make builtin.
func (c *compiler) compileMakeBuiltin(args []ssa.Value) {
	t := args[0].Type()
	typeIdx := c.addType(t)
	typeIdxConst := c.addConstant(int64(typeIdx))
	zeroIdx := c.addConstant(int64(0))

	switch t.(type) {
	case *types.Slice:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) >= 2 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		if len(args) >= 3 {
			c.compileValue(args[2])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeSlice)
	case *types.Map:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeMap)
	case *types.Chan:
		c.emit(bytecode.OpConst, typeIdxConst)
		if len(args) > 1 {
			c.compileValue(args[1])
		} else {
			c.emit(bytecode.OpConst, zeroIdx)
		}
		c.emit(bytecode.OpMakeChan)
	}
}
