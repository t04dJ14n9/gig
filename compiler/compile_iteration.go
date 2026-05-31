package compiler

import (
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// compileMapUpdate compiles a MapUpdate instruction.
func (c *compiler) compileMapUpdate(i *ssa.MapUpdate) {
	c.compileValue(i.Map)
	c.compileValue(i.Key)
	c.compileValue(i.Value)
	c.emit(bytecode.OpSetIndex)
}

// compileRange compiles a Range instruction.
func (c *compiler) compileRange(i *ssa.Range) {
	c.compileSimpleUnaryOp(i, i.X, bytecode.OpRange)
}

// compileNext compiles a Next instruction.
func (c *compiler) compileNext(i *ssa.Next) {
	c.compileSimpleUnaryOp(i, i.Iter, bytecode.OpRangeNext)
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
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(bytecode.SliceEndSentinel))))
	}

	if i.Max != nil {
		c.compileValue(i.Max)
	} else {
		c.emit(bytecode.OpConst, uint16(c.addConstant(int64(bytecode.SliceEndSentinel))))
	}

	c.emit(bytecode.OpSlice)

	// If the result type is a named slice type (e.g., sort.IntSlice from [5]int[:]),
	// emit OpChangeType to convert the underlying []int to the named type.
	if named, ok := i.Type().(*types.Named); ok {
		if _, isSlice := named.Underlying().(*types.Slice); isSlice {
			typeIdx := c.addType(named)
			srcLocalIdx := uint16(bytecode.NoSourceLocal)
			c.currentFunc.Instructions = append(c.currentFunc.Instructions,
				byte(bytecode.OpChangeType),
				byte(typeIdx>>8), byte(typeIdx),
				byte(srcLocalIdx>>8), byte(srcLocalIdx))
		}
	}

	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
