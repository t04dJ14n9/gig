package compiler

import (
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// builtinResultOps maps simple builtin names to their single-opcode implementations.
// Builtins not in this map have custom compilation logic in compileBuiltinCall.
var builtinResultOps = map[string]bytecode.OpCode{
	"len":     bytecode.OpLen,
	"cap":     bytecode.OpCap,
	"copy":    bytecode.OpCopy,
	"recover": bytecode.OpRecover,
}

// compileBuiltinCall compiles a call to a builtin function.
func (c *compiler) compileBuiltinCall(builtin *ssa.Builtin, args []ssa.Value, resultIdx int) {
	name := builtin.Name()

	if c.compileNoResultBuiltin(name, args) {
		return
	}

	c.compileResultBuiltin(name, args)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

func (c *compiler) compileNoResultBuiltin(name string, args []ssa.Value) bool {
	switch name {
	case "delete":
		c.compileValue(args[0])
		c.compileValue(args[1])
		c.emit(bytecode.OpDelete)
		return true
	case "print":
		c.compileValues(args)
		c.emit(bytecode.OpPrint, uint16(len(args)))
		return true
	case "println":
		c.compileValues(args)
		c.emit(bytecode.OpPrintln, uint16(len(args)))
		return true
	case "panic":
		c.compileValue(args[0])
		c.emit(bytecode.OpPanic)
		return true
	case "close":
		c.compileValue(args[0])
		c.emit(bytecode.OpClose)
		return true
	default:
		return false
	}
}

func (c *compiler) compileValues(args []ssa.Value) {
	for _, arg := range args {
		c.compileValue(arg)
	}
}

func (c *compiler) compileResultBuiltin(name string, args []ssa.Value) {
	if c.compileSimpleResultBuiltin(name, args) {
		return
	}
	if c.compileCustomResultBuiltin(name, args) {
		return
	}
	c.emit(bytecode.OpNil)
}

func (c *compiler) compileSimpleResultBuiltin(name string, args []ssa.Value) bool {
	op, ok := builtinResultOps[name]
	if !ok {
		return false
	}
	switch name {
	case "len", "cap":
		c.compileValue(args[0])
	case "copy":
		c.compileValue(args[0])
		c.compileValue(args[1])
	case "recover":
		// recover() takes no arguments.
	}
	c.emit(op)
	return true
}

func (c *compiler) compileCustomResultBuiltin(name string, args []ssa.Value) bool {
	switch name {
	case "append":
		c.compileAppendBuiltin(args)
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
		return false
	}
	return true
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
	alloc, length, ok := packedVarargsAlloc(v)
	if !ok {
		return nil, false
	}

	values := make([]ssa.Value, length)
	if !collectPackedVarargsStores(alloc, values) || !packedVarargsComplete(values) {
		return nil, false
	}
	return values, true
}

func packedVarargsAlloc(v ssa.Value) (*ssa.Alloc, int, bool) {
	slice, ok := v.(*ssa.Slice)
	if !ok || slice.Low != nil || slice.High != nil || slice.Max != nil {
		return nil, 0, false
	}
	alloc, ok := slice.X.(*ssa.Alloc)
	if !ok || alloc.Comment != "varargs" {
		return nil, 0, false
	}
	length, ok := packedVarargsArrayLen(alloc)
	return alloc, length, ok
}

func packedVarargsArrayLen(alloc *ssa.Alloc) (int, bool) {
	ptr, ok := alloc.Type().Underlying().(*types.Pointer)
	if !ok {
		return 0, false
	}
	arr, ok := ptr.Elem().Underlying().(*types.Array)
	if !ok || arr.Len() < 0 {
		return 0, false
	}
	return int(arr.Len()), true
}

func collectPackedVarargsStores(alloc *ssa.Alloc, values []ssa.Value) bool {
	refs := alloc.Referrers()
	if refs == nil {
		return false
	}
	for _, ref := range *refs {
		storePackedVararg(values, alloc, ref)
	}
	return true
}

func storePackedVararg(values []ssa.Value, alloc *ssa.Alloc, ref ssa.Instruction) {
	indexAddr, ok := ref.(*ssa.IndexAddr)
	if !ok || indexAddr.X != alloc {
		return
	}
	idx, ok := packedVarargIndex(indexAddr.Index, len(values))
	if !ok {
		return
	}
	if val, ok := packedVarargStoreValue(indexAddr); ok {
		values[idx] = val
	}
}

func packedVarargIndex(index ssa.Value, length int) (int, bool) {
	idxConst, ok := index.(*ssa.Const)
	if !ok || idxConst.Value == nil {
		return 0, false
	}
	idx, exact := constant.Int64Val(idxConst.Value)
	if !exact || idx < 0 || idx >= int64(length) {
		return 0, false
	}
	return int(idx), true
}

func packedVarargStoreValue(indexAddr *ssa.IndexAddr) (ssa.Value, bool) {
	indexRefs := indexAddr.Referrers()
	if indexRefs == nil {
		return nil, false
	}
	for _, indexRef := range *indexRefs {
		store, ok := indexRef.(*ssa.Store)
		if ok && store.Addr == indexAddr {
			return store.Val, true
		}
	}
	return nil, false
}

func packedVarargsComplete(values []ssa.Value) bool {
	for _, val := range values {
		if val == nil {
			return false
		}
	}
	return true
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
