// ops_memory.go routes cold memory-family opcodes to focused domain handlers.
//
// Hot stack, constant, local, and boolean opcodes are intentionally inlined in
// run.go. The dispatcher here keeps the less common memory domains readable:
// variable slots, fields, pointer operations, and allocation.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// executeMemory handles global/free variable, field, address, dereference, and new opcodes.
// Note: OpPop, OpDup, OpConst, OpNil, OpTrue, OpFalse, OpLocal, OpSetLocal
// are inlined in run.go's hot path and never reach this handler.
func (v *vm) executeMemory(op bytecode.OpCode, frame *Frame) error {
	switch op {
	case bytecode.OpGlobal, bytecode.OpSetGlobal, bytecode.OpFree, bytecode.OpSetFree:
		v.executeVariableMemory(op, frame)
	case bytecode.OpField, bytecode.OpSetField:
		v.executeFieldMemory(op, frame)
	case bytecode.OpAddr, bytecode.OpFieldAddr, bytecode.OpIndexAddr, bytecode.OpDeref, bytecode.OpSetDeref:
		v.executePointerMemory(op, frame)
	case bytecode.OpNew:
		v.executeNew(frame)
	}

	return nil
}
