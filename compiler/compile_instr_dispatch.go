// compile_instr_dispatch.go keeps SSA instruction dispatch split by compiler concern.
package compiler

import (
	"golang.org/x/tools/go/ssa"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func (c *compiler) compileExpressionInstruction(instr ssa.Instruction) bool {
	switch i := instr.(type) {
	case *ssa.BinOp:
		c.compileBinOp(i)
	case *ssa.UnOp:
		c.compileUnOp(i)
	case *ssa.Call:
		c.compileCall(i)
	case *ssa.ChangeInterface:
		c.compileChangeInterface(i)
	case *ssa.ChangeType:
		c.compileChangeType(i)
	case *ssa.Convert:
		c.compileConvert(i)
	case *ssa.Extract:
		c.compileExtract(i)
	case *ssa.MakeInterface:
		c.compileMakeInterface(i)
	case *ssa.TypeAssert:
		c.compileTypeAssert(i)
	default:
		return false
	}
	return true
}

func (c *compiler) compileAddressingInstruction(instr ssa.Instruction) bool {
	switch i := instr.(type) {
	case *ssa.Alloc:
		c.compileAlloc(i)
	case *ssa.Field:
		c.compileField(i)
	case *ssa.FieldAddr:
		c.compileFieldAddr(i)
	case *ssa.Index:
		c.compileIndex(i)
	case *ssa.IndexAddr:
		c.compileIndexAddr(i)
	case *ssa.Lookup:
		c.compileLookup(i)
	case *ssa.Slice:
		c.compileSlice(i)
	default:
		return false
	}
	return true
}

func (c *compiler) compileConstructionInstruction(instr ssa.Instruction) bool {
	switch i := instr.(type) {
	case *ssa.MakeClosure:
		c.compileMakeClosure(i)
	case *ssa.MakeChan:
		c.compileMakeChan(i)
	case *ssa.MakeMap:
		c.compileMakeMap(i)
	case *ssa.MakeSlice:
		c.compileMakeSlice(i)
	default:
		return false
	}
	return true
}

func (c *compiler) compileControlInstruction(instr ssa.Instruction) bool {
	switch i := instr.(type) {
	case *ssa.Next:
		c.compileNext(i)
	case *ssa.Range:
		c.compileRange(i)
	case *ssa.Select:
		c.compileSelect(i)
	default:
		return false
	}
	return true
}

func (c *compiler) compileEffectInstruction(instr ssa.Instruction) bool {
	switch i := instr.(type) {
	case *ssa.Defer:
		c.compileDefer(i)
	case *ssa.Go:
		c.compileGo(i)
	case *ssa.MapUpdate:
		c.compileMapUpdate(i)
	case *ssa.Return:
		c.compileReturn(i)
	case *ssa.RunDefers:
		c.emit(bytecode.OpRunDefers)
	case *ssa.Send:
		c.compileSend(i)
	case *ssa.Store:
		c.compileStore(i)
	default:
		return false
	}
	return true
}

func (c *compiler) compileIgnoredInstruction(instr ssa.Instruction) {
	switch instr.(type) {
	case *ssa.Phi:
		// Phi nodes are lowered into predecessor moves before instruction dispatch.
	case *ssa.DebugRef:
		// Debug references describe source locations and do not emit bytecode.
	case *ssa.Panic, *ssa.Jump, *ssa.If:
		// Terminators are emitted while compiling the owning basic block.
	}
}
