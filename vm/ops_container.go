// ops_container.go dispatches VM container opcodes to focused handlers.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// executeContainer handles slice, map, channel creation, index, append,
// copy, delete, range, len, and cap opcodes.
func (v *vm) executeContainer(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx,unparam // frame: uniform dispatch signature
	switch op {
	case bytecode.OpMakeSlice:
		v.executeMakeSlice()

	case bytecode.OpMakeMap:
		v.executeMakeMap()

	case bytecode.OpMakeChan:
		v.executeMakeChan()

	// Index operations
	case bytecode.OpIndex:
		v.executeIndex()

	case bytecode.OpIndexOk:
		v.executeIndexOk()

	case bytecode.OpSetIndex:
		v.executeSetIndex()

	case bytecode.OpSlice:
		v.executeSlice()

	case bytecode.OpRange:
		v.executeRange()

	case bytecode.OpRangeNext:
		v.executeRangeNext()

	case bytecode.OpLen:
		v.executeLen()

	case bytecode.OpCap:
		v.executeCap()

	case bytecode.OpAppend:
		elem := v.pop()
		slice := v.pop()
		v.push(appendValue(slice, elem))

	case bytecode.OpCopy:
		v.executeCopy()

	case bytecode.OpDelete:
		v.executeDelete()

	}

	return nil
}
