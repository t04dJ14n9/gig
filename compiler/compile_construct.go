package compiler

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"golang.org/x/tools/go/ssa"
)

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
	ifaceTypeIdx := c.addType(i.Type())      // interface type (e.g., error)
	concreteTypeIdx := c.addType(i.X.Type()) // concrete type (e.g., *MyError3)
	c.compileValue(i.X)
	// Emit: OpMakeInterface, ifaceTypeIdx_hi, ifaceTypeIdx_lo, concreteTypeIdx_hi, concreteTypeIdx_lo
	c.currentFunc.Instructions = append(c.currentFunc.Instructions,
		byte(bytecode.OpMakeInterface),
		byte(ifaceTypeIdx>>8), byte(ifaceTypeIdx),
		byte(concreteTypeIdx>>8), byte(concreteTypeIdx))
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

	c.emitClosure(fnIdx, len(i.Bindings))
	c.emit(bytecode.OpSetLocal, uint16(resultIdx))
}
