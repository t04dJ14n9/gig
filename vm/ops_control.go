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
	case bytecode.OpSend, bytecode.OpRecv, bytecode.OpRecvOk, bytecode.OpClose:
		return v.executeChannelControl(op)
	case bytecode.OpSelect:
		return v.executeSelect(frame)
	case bytecode.OpDefer, bytecode.OpDeferExternal, bytecode.OpDeferIndirect, bytecode.OpRunDefers:
		return v.executeDeferControl(op, frame)
	case bytecode.OpRecover, bytecode.OpPanic:
		v.executePanicControl(op)
		return nil
	case bytecode.OpPrint, bytecode.OpPrintln:
		v.executePrintControl(op, frame)
		return nil
	case bytecode.OpHalt:
		return fmt.Errorf("halt")
	default:
		return fmt.Errorf("unknown opcode: %v", op)
	}
}

func (v *vm) executeChannelControl(op bytecode.OpCode) error {
	switch op {
	case bytecode.OpSend:
		return v.executeSend()
	case bytecode.OpRecv:
		return v.executeRecv()
	case bytecode.OpRecvOk:
		return v.executeRecvOk()
	case bytecode.OpClose:
		v.executeClose()
	}
	return nil
}

func (v *vm) executeDeferControl(op bytecode.OpCode, frame *Frame) error {
	switch op {
	case bytecode.OpDefer:
		v.executeDefer(frame)
	case bytecode.OpDeferExternal:
		v.executeDeferExternal(frame)
	case bytecode.OpDeferIndirect:
		return v.executeDeferIndirect(frame)
	case bytecode.OpRunDefers:
		return v.executeRunDefers(frame)
	}
	return nil
}

func (v *vm) executePanicControl(op bytecode.OpCode) {
	switch op {
	case bytecode.OpRecover:
		v.executeRecover()
	case bytecode.OpPanic:
		v.executePanic()
	}
}

func (v *vm) executePrintControl(op bytecode.OpCode, frame *Frame) {
	switch op {
	case bytecode.OpPrint:
		v.executePrint(frame)
	case bytecode.OpPrintln:
		v.executePrintln(frame)
	}
}
