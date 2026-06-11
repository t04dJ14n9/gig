package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runColdOp(frame *Frame, op bytecode.OpCode, sp int) (int, []value.Value, bool, error) {
	// Cold opcodes use executeOp, which operates on v.sp/v.stack directly.
	// Return the refreshed stack state and whether the caller should reload the
	// cached frame fields before the next fetch.
	v.sp = sp
	if err := v.executeOp(op, frame); err != nil {
		return v.sp, v.stack, false, err
	}
	return v.sp, v.stack, v.fp > 0, nil
}
