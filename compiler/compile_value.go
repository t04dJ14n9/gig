package compiler

import (
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/bytecode"
)

// compileValue compiles an SSA value to push it onto the stack.
func (c *compiler) compileValue(v ssa.Value) {
	switch val := v.(type) {
	case *ssa.Const:
		c.compileConst(val)
	case *ssa.Function:
		if fnIdx, ok := c.funcIndex[val]; ok {
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpClosure),
				byte(fnIdx>>8), byte(fnIdx),
				byte(0))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.Phi:
		if slot, ok := c.phiSlots[val]; ok {
			c.emit(bytecode.OpLocal, uint16(slot))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.FreeVar:
		if idx, ok := c.symbolTable.freeVars[val]; ok {
			c.emit(bytecode.OpFree, uint16(idx))
		} else {
			c.emit(bytecode.OpNil)
		}
	case *ssa.Global:
		globalName := val.Name()
		globalIdx, ok := c.globals[globalName]
		if !ok {
			globalIdx = len(c.globals)
			c.globals[globalName] = globalIdx
		}
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpGlobal),
			byte(globalIdx>>8), byte(globalIdx))
	default:
		if idx, ok := c.symbolTable.GetLocal(v); ok {
			c.emit(bytecode.OpLocal, uint16(idx))
		} else {
			if idx, ok := c.symbolTable.freeVars[v]; ok {
				c.emit(bytecode.OpFree, uint16(idx))
			} else {
				c.emit(bytecode.OpNil)
			}
		}
	}
}

// compileConst compiles a constant value.
func (c *compiler) compileConst(cnst *ssa.Const) {
	var v any
	switch t := cnst.Type().(type) {
	case *types.Basic:
		switch t.Kind() { //nolint:exhaustive
		case types.Bool, types.UntypedBool:
			v = cnst.Value != nil && cnst.Value.Kind() == constant.Bool && constant.BoolVal(cnst.Value)
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.UntypedInt, types.UntypedRune:
			if cnst.Value != nil {
				i, exact := constant.Int64Val(cnst.Value)
				if exact {
					v = i
				} else {
					v = int64(0)
				}
			} else {
				v = int64(0)
			}
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
			if cnst.Value != nil {
				v, _ = constant.Uint64Val(cnst.Value)
			} else {
				v = uint64(0)
			}
		case types.Float32, types.Float64, types.UntypedFloat:
			if cnst.Value != nil {
				f, _ := constant.Float64Val(cnst.Value)
				v = f
			} else {
				v = 0.0
			}
		case types.String, types.UntypedString:
			if cnst.Value != nil {
				v = constant.StringVal(cnst.Value)
			} else {
				v = ""
			}
		default:
			v = nil
		}
	default:
		v = nil
	}

	idx := c.addConstant(v)
	c.emit(bytecode.OpConst, idx)
}

// compileField compiles a Field instruction.
func (c *compiler) compileField(i *ssa.Field) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpField, uint16(i.Field))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileFieldAddr compiles a FieldAddr instruction.
func (c *compiler) compileFieldAddr(i *ssa.FieldAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpFieldAddr, uint16(i.Field))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndex compiles an Index instruction.
func (c *compiler) compileIndex(i *ssa.Index) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(bytecode.OpIndex)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileIndexAddr compiles an IndexAddr instruction.
func (c *compiler) compileIndexAddr(i *ssa.IndexAddr) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)
	c.emit(bytecode.OpIndexAddr)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileLookup compiles a Lookup instruction.
func (c *compiler) compileLookup(i *ssa.Lookup) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.compileValue(i.Index)

	if i.CommaOk {
		c.emit(bytecode.OpIndexOk)
	} else {
		c.emit(bytecode.OpIndex)
	}
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileStore compiles a Store instruction.
func (c *compiler) compileStore(i *ssa.Store) {
	c.compileValue(i.Addr)
	c.compileValue(i.Val)
	c.emit(bytecode.OpSetDeref)
}

// compileMakeSlice compiles a MakeSlice instruction.
func (c *compiler) compileMakeSlice(i *ssa.MakeSlice) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)
	c.compileValue(i.Len)
	c.compileValue(i.Cap)
	c.emit(bytecode.OpMakeSlice)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeMap compiles a MakeMap instruction.
func (c *compiler) compileMakeMap(i *ssa.MakeMap) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)

	if i.Reserve != nil {
		c.compileValue(i.Reserve)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(bytecode.OpMakeMap)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeChan compiles a MakeChan instruction.
func (c *compiler) compileMakeChan(i *ssa.MakeChan) {
	typeIdx := c.addType(i.Type())
	resultIdx := c.symbolTable.AllocLocal(i)

	typeIdxConst := c.addConstant(int64(typeIdx))
	c.emit(bytecode.OpConst, typeIdxConst)

	if i.Size != nil {
		c.compileValue(i.Size)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	c.emit(bytecode.OpMakeChan)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeInterface compiles a MakeInterface instruction.
func (c *compiler) compileMakeInterface(i *ssa.MakeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMakeClosure compiles a MakeClosure instruction.
func (c *compiler) compileMakeClosure(i *ssa.MakeClosure) {
	fnIdx := c.funcIndex[i.Fn.(*ssa.Function)]
	resultIdx := c.symbolTable.AllocLocal(i)

	for _, binding := range i.Bindings {
		if alloc, ok := binding.(*ssa.Alloc); ok {
			if slotIdx, ok := c.symbolTable.GetLocal(alloc); ok {
				// For *ssa.Alloc (which is already a pointer type), we need to
				// get the pointer value itself (OpLocal), not the address of the slot (OpAddr).
				// Each Alloc creates a new pointer in heap/stack, and closures should
				// capture this pointer value, not the slot address which gets overwritten
				// in loop iterations.
				c.emit(bytecode.OpLocal, uint16(slotIdx))
				continue
			}
		}
		c.compileValue(binding)
	}

	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpClosure),
		byte(fnIdx>>8), byte(fnIdx),
		byte(len(i.Bindings)))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileMapUpdate compiles a MapUpdate instruction.
func (c *compiler) compileMapUpdate(i *ssa.MapUpdate) {
	c.compileValue(i.Map)
	c.compileValue(i.Key)
	c.compileValue(i.Value)
	c.emit(bytecode.OpSetIndex)
}

// compileRange compiles a Range instruction.
func (c *compiler) compileRange(i *ssa.Range) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpRange)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileNext compiles a Next instruction.
func (c *compiler) compileNext(i *ssa.Next) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.Iter)
	c.emit(bytecode.OpRangeNext)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSelect compiles a Select instruction.
func (c *compiler) compileSelect(i *ssa.Select) {
	numRecv := 0
	for _, st := range i.States {
		if st.Dir == types.RecvOnly {
			numRecv++
		}
	}

	dirs := make([]bool, len(i.States))
	for idx, st := range i.States {
		dirs[idx] = (st.Dir == types.SendOnly)
	}

	meta := bytecode.SelectMeta{
		NumStates: len(i.States),
		Blocking:  i.Blocking,
		Dirs:      dirs,
		NumRecv:   numRecv,
	}

	for _, st := range i.States {
		c.compileValue(st.Chan)
		if st.Dir == types.SendOnly {
			c.compileValue(st.Send)
		}
	}

	metaIdx := c.addConstant(meta)
	c.emit(bytecode.OpSelect, metaIdx)

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSlice compiles a Slice instruction.
func (c *compiler) compileSlice(i *ssa.Slice) {
	resultIdx := c.symbolTable.AllocLocal(i)

	c.compileValue(i.X)

	if i.Low != nil {
		c.compileValue(i.Low)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0))))
	}

	if i.High != nil {
		c.compileValue(i.High)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	if i.Max != nil {
		c.compileValue(i.Max)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(0xFFFF))))
	}

	c.emit(bytecode.OpSlice)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileTypeAssert compiles a TypeAssert instruction.
func (c *compiler) compileTypeAssert(i *ssa.TypeAssert) {
	typeIdx := c.addType(i.AssertedType)
	c.compileValue(i.X)
	c.emit(bytecode.OpAssert, uint16(typeIdx))

	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeInterface compiles a ChangeInterface instruction.
func (c *compiler) compileChangeInterface(i *ssa.ChangeInterface) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileChangeType compiles a ChangeType instruction.
func (c *compiler) compileChangeType(i *ssa.ChangeType) {
	resultIdx := c.symbolTable.AllocLocal(i)
	c.compileValue(i.X)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileConvert compiles a Convert instruction.
func (c *compiler) compileConvert(i *ssa.Convert) {
	resultIdx := c.symbolTable.AllocLocal(i)
	typeIdx := c.addType(i.Type())
	c.compileValue(i.X)
	c.emit(bytecode.OpConvert, uint16(typeIdx))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileExtract compiles an Extract instruction.
func (c *compiler) compileExtract(i *ssa.Extract) {
	c.compileValue(i.Tuple)
	c.emit(bytecode.OpConst, uint16(c.addConstant(i.Index)))
	c.emit(bytecode.OpIndex)
	resultIdx := c.symbolTable.AllocLocal(i)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}

// compileSend compiles a Send instruction.
func (c *compiler) compileSend(i *ssa.Send) {
	c.compileValue(i.Chan)
	c.compileValue(i.X)
	c.emit(bytecode.OpSend)
}

// compileDefer compiles a Defer instruction.
func (c *compiler) compileDefer(i *ssa.Defer) {
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		fnIdx := c.funcIndex[fn]
		c.emit(bytecode.OpDefer, uint16(fnIdx))
	}
}

// compileGo compiles a Go instruction.
func (c *compiler) compileGo(i *ssa.Go) {
	for _, arg := range i.Call.Args {
		c.compileValue(arg)
	}

	if fn, ok := i.Call.Value.(*ssa.Function); ok {
		funcIdx := c.funcIndex[fn]
		numArgs := len(i.Call.Args)
		c.currentFunc.Instructions = append(c.currentFunc.Instructions,
			byte(bytecode.OpGoCall),
			byte(funcIdx>>8), byte(funcIdx),
			byte(numArgs))
		return
	}

	c.compileValue(i.Call.Value)

	numArgs := len(i.Call.Args)
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpGoCallIndirect),
		byte(numArgs))
}
