// ops_control.go dispatches non-hot control opcodes.
package vm

import (
	"fmt"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

// executeControl handles channels, select, defer, panic/recover, print, and halt opcodes.
// Note: OpJump, OpJumpTrue, OpJumpFalse, OpReturn, OpReturnVal are inlined in run.go's
// hot path and never reach this handler.
func (v *vm) executeControl(op bytecode.OpCode, frame *Frame) error {
	switch op {
	case bytecode.OpSend:
		return v.executeSend()
	case bytecode.OpRecv:
		return v.executeRecv()
	case bytecode.OpRecvOk:
		return v.executeRecvOk()
	case bytecode.OpClose:
		v.executeClose()
		return nil
	case bytecode.OpSelect:
		return v.executeSelect(frame)
	case bytecode.OpDefer:
		v.executeDefer(frame)
		return nil
	case bytecode.OpDeferExternal:
		v.executeDeferExternal(frame)
		return nil
	case bytecode.OpDeferIndirect:
		return v.executeDeferIndirect(frame)
	case bytecode.OpRunDefers:
		return v.executeRunDefers(frame)
	case bytecode.OpRecover:
		v.executeRecover()
		return nil
	case bytecode.OpPanic:
		v.executePanic()
		return nil
	case bytecode.OpPrint:
		v.executePrint(frame)
		return nil
	case bytecode.OpPrintln:
		v.executePrintln(frame)
		return nil
	case bytecode.OpHalt:
		return fmt.Errorf("halt")
	default:
		return fmt.Errorf("unknown opcode: %v", op)
	}
}
