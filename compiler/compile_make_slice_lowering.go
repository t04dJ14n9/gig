package compiler

import (
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func (c *compiler) compileSyntheticMakeSlice(slice *ssa.Slice) bool {
	_, arr, ok := syntheticMakeSliceArrayAllocForSlice(slice)
	if !ok {
		return false
	}

	resultIdx := c.symbolTable.AllocLocal(slice)
	typeIdxConst := c.addConstant(int64(c.addType(slice.Type())))

	c.emit(bytecode.OpConst, typeIdxConst)
	if slice.High != nil {
		c.compileValue(slice.High)
	} else {
		c.emitIntConstant(arr.Len())
	}
	c.emitIntConstant(arr.Len())
	c.emit(bytecode.OpMakeSlice)
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
	return true
}

func (c *compiler) emitIntConstant(n int64) {
	c.emit(bytecode.OpConst, uint16(c.addConstant(n)))
}

func syntheticMakeSliceArrayAllocForSlice(slice *ssa.Slice) (*ssa.Alloc, *types.Array, bool) {
	alloc, ok := slice.X.(*ssa.Alloc)
	if !ok {
		return nil, nil, false
	}
	arr, ok := syntheticMakeSliceArrayAlloc(alloc)
	if !ok || !canLowerSyntheticMakeSlice(slice, arr) {
		return nil, nil, false
	}
	return alloc, arr, true
}

func syntheticMakeSliceArrayAlloc(alloc *ssa.Alloc) (*types.Array, bool) {
	if alloc.Comment != "makeslice" {
		return nil, false
	}
	arr, ok := intArrayAllocType(alloc)
	if !ok {
		return nil, false
	}
	if !hasSingleSyntheticMakeSliceRef(alloc, arr) {
		return nil, false
	}
	return arr, true
}

func hasSingleSyntheticMakeSliceRef(alloc *ssa.Alloc, arr *types.Array) bool {
	refs := alloc.Referrers()
	if refs == nil {
		return false
	}
	count, ok := countSyntheticMakeSliceRefs(*refs, alloc, arr)
	return ok && count == 1
}

func countSyntheticMakeSliceRefs(refs []ssa.Instruction, alloc *ssa.Alloc, arr *types.Array) (int, bool) {
	sliceRefs := 0
	for _, ref := range refs {
		if isDebugRef(ref) {
			continue
		}
		if !isSyntheticMakeSliceRef(ref, alloc, arr) {
			return 0, false
		}
		sliceRefs++
	}
	return sliceRefs, true
}

func isDebugRef(ref ssa.Instruction) bool {
	_, ok := ref.(*ssa.DebugRef)
	return ok
}

func isSyntheticMakeSliceRef(ref ssa.Instruction, alloc *ssa.Alloc, arr *types.Array) bool {
	slice, ok := ref.(*ssa.Slice)
	return ok && slice.X == alloc && canLowerSyntheticMakeSlice(slice, arr)
}

func intArrayAllocType(alloc *ssa.Alloc) (*types.Array, bool) {
	ptr, ok := alloc.Type().Underlying().(*types.Pointer)
	if !ok {
		return nil, false
	}
	arr, ok := ptr.Elem().Underlying().(*types.Array)
	if !ok || arr.Len() < 0 {
		return nil, false
	}
	elem, ok := arr.Elem().Underlying().(*types.Basic)
	if !ok {
		return nil, false
	}
	switch elem.Kind() {
	case types.Int, types.Int64:
		return arr, true
	default:
		return nil, false
	}
}

func canLowerSyntheticMakeSlice(slice *ssa.Slice, _ *types.Array) bool {
	return isIntSliceType(slice.Type()) &&
		isNilOrZeroConst(slice.Low) &&
		slice.Max == nil
}

func isNilOrZeroConst(v ssa.Value) bool {
	if v == nil {
		return true
	}
	n, ok := constInt64(v)
	return ok && n == 0
}

func constInt64(v ssa.Value) (int64, bool) {
	c, ok := v.(*ssa.Const)
	if !ok || c.Value == nil {
		return 0, false
	}
	n, exact := constant.Int64Val(c.Value)
	return n, exact
}
