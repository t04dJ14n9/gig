package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runSlowConst(idx uint16) value.Value {
	// Most constants are prebaked and loaded directly in run.go. This fallback
	// keeps legacy interface-backed constants correct without making the common
	// constant load carry the extra bounds branch.
	if int(idx) < len(v.program.Constants) {
		return value.FromInterface(v.program.Constants[idx])
	}
	return value.Value{}
}

func runLiteralValue(op bytecode.OpCode) value.Value {
	// Nil/boolean opcodes all push a single literal and advance sp in the same
	// way. Keeping the literal selection here lets run() group those dispatch
	// cases without hiding any stack mutation.
	switch op { //nolint:exhaustive
	case bytecode.OpTrue:
		return value.MakeBool(true)
	case bytecode.OpFalse:
		return value.MakeBool(false)
	default:
		return value.MakeNil()
	}
}
