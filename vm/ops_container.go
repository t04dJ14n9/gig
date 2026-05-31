// ops_container.go dispatches VM container opcodes to focused handlers.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// executeContainer handles slice, map, channel creation, index, append,
// copy, delete, range, len, and cap opcodes.
func (v *vm) executeContainer(op bytecode.OpCode, frame *Frame) error { //nolint:unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpMakeSlice, bytecode.OpMakeMap, bytecode.OpMakeChan:
		v.executeMakeContainer(op)
	case bytecode.OpIndex, bytecode.OpIndexOk, bytecode.OpSetIndex, bytecode.OpSlice:
		v.executeIndexContainer(op)
	case bytecode.OpRange, bytecode.OpRangeNext:
		v.executeRangeContainer(op)
	case bytecode.OpLen, bytecode.OpCap:
		v.executeSizeContainer(op)
	case bytecode.OpAppend:
		v.executeAppend()
	case bytecode.OpCopy:
		v.executeCopy()
	case bytecode.OpDelete:
		v.executeDelete()
	}

	return nil
}

func (v *vm) executeMakeContainer(op bytecode.OpCode) {
	switch op {
	case bytecode.OpMakeSlice:
		v.executeMakeSlice()
	case bytecode.OpMakeMap:
		v.executeMakeMap()
	case bytecode.OpMakeChan:
		v.executeMakeChan()
	}
}

func (v *vm) executeIndexContainer(op bytecode.OpCode) {
	switch op {
	case bytecode.OpIndex:
		v.executeIndex()
	case bytecode.OpIndexOk:
		v.executeIndexOk()
	case bytecode.OpSetIndex:
		v.executeSetIndex()
	case bytecode.OpSlice:
		v.executeSlice()
	}
}

func (v *vm) executeRangeContainer(op bytecode.OpCode) {
	switch op {
	case bytecode.OpRange:
		v.executeRange()
	case bytecode.OpRangeNext:
		v.executeRangeNext()
	}
}

func (v *vm) executeSizeContainer(op bytecode.OpCode) {
	switch op {
	case bytecode.OpLen:
		v.executeLen()
	case bytecode.OpCap:
		v.executeCap()
	}
}

func (v *vm) executeAppend() {
	elem := v.pop()
	slice := v.pop()
	v.push(appendValue(slice, elem))
}
