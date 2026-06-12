// ops_call.go routes cold call-family opcodes to focused domain handlers.
//
// The hot direct-call instructions (OpCall, OpCallExternal, OpCallIndirect)
// are intentionally inlined in run.go; only closure construction, goroutine
// spawning, and tuple shape helpers reach this dispatcher.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// executeCall handles closure creation, goroutine spawning, and pack/unpack opcodes.
// Note: OpCall, OpCallExternal, OpCallIndirect are inlined in run.go's hot path
// and never reach this handler.
func (v *vm) executeCall(op bytecode.OpCode, frame *Frame) error {
	switch op {
	case bytecode.OpClosure:
		v.executeClosure(frame)
	case bytecode.OpGoCall:
		return v.executeGoCall(frame)
	case bytecode.OpGoCallExternal:
		return v.executeGoCallExternal(frame)
	case bytecode.OpGoCallIndirect:
		return v.executeGoCallIndirect(frame)
	case bytecode.OpPack:
		v.executePack(frame)
	case bytecode.OpUnpack:
		v.executeUnpack()
	}

	return nil
}
