package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// These helpers keep the native []int fast path in run.go while isolating the
// reflective fallback used when an optimized int-slice opcode sees another
// indexable value. They intentionally preserve executeOp's panic/error order.
func (v *vm) runIntSliceGetFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, vIdx uint16,
	sp int,
) (int, []value.Value, error) {
	if err := v.runIntSliceFallbackIndex(frame, locals, intLocals, sIdx, jIdx, sp); err != nil || v.panicking {
		return v.sp, v.stack, err
	}
	if err := v.executeOp(bytecode.OpDeref, frame); err != nil || v.panicking {
		return v.sp, v.stack, err
	}

	ret := v.pop()
	intLocals[vIdx] = ret.RawInt()
	locals[vIdx] = ret
	return v.sp, v.stack, nil
}

func (v *vm) runIntSliceSetFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, valIdx uint16,
	sp int,
) (int, []value.Value, error) {
	if err := v.runIntSliceFallbackIndex(frame, locals, intLocals, sIdx, jIdx, sp); err != nil || v.panicking {
		return v.sp, v.stack, err
	}
	v.push(value.MakeInt(intLocals[valIdx]))
	if err := v.executeOp(bytecode.OpSetDeref, frame); err != nil || v.panicking {
		return v.sp, v.stack, err
	}
	return v.sp, v.stack, nil
}

func (v *vm) runIntSliceSetConstFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	prebaked []value.Value,
	sIdx, jIdx, cIdx uint16,
	sp int,
) (int, []value.Value, error) {
	if err := v.runIntSliceFallbackIndex(frame, locals, intLocals, sIdx, jIdx, sp); err != nil || v.panicking {
		return v.sp, v.stack, err
	}
	v.push(prebaked[cIdx])
	if err := v.executeOp(bytecode.OpSetDeref, frame); err != nil || v.panicking {
		return v.sp, v.stack, err
	}
	return v.sp, v.stack, nil
}

func (v *vm) runIntSliceFallbackIndex(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx uint16,
	sp int,
) error {
	v.sp = sp
	v.push(locals[sIdx])
	v.push(value.MakeInt(intLocals[jIdx]))
	return v.executeOp(bytecode.OpIndexAddr, frame)
}
