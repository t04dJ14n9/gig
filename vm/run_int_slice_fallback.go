package vm

import (
	"github.com/t04dJ14n9/gig/model/value"
)

// These helpers keep the native []int fast path in run.go while isolating the
// reflective fallback used when an optimized int-slice opcode sees another
// indexable value. The fallback calls the same reference-boundary functions as
// OpIndexAddr/OpDeref/OpSetDeref, but avoids routing each access back through
// executeOp. That dispatcher path is readable for cold opcodes but too costly
// for fused slice operations in tight loops.
func (v *vm) runIntSliceGetFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, vIdx uint16,
	sp int,
) (int, []value.Value, error) {
	v.sp = sp
	if v.runIntSliceGetFallbackRecovered(locals, intLocals, sIdx, jIdx, vIdx) {
		return v.sp, v.stack, nil
	}
	return sp, v.stack, nil
}

func (v *vm) runIntSliceSetFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, valIdx uint16,
	sp int,
) (int, []value.Value, error) {
	v.sp = sp
	if v.runIntSliceSetFallbackRecovered(locals, intLocals, sIdx, jIdx, valIdx) {
		return v.sp, v.stack, nil
	}
	return sp, v.stack, nil
}

func (v *vm) runIntSliceSetConstFallback(
	frame *Frame,
	locals []value.Value,
	intLocals []int64,
	prebaked []value.Value,
	sIdx, jIdx, cIdx uint16,
	sp int,
) (int, []value.Value, error) {
	v.sp = sp
	if v.runIntSliceSetConstFallbackRecovered(locals, intLocals, prebaked, sIdx, jIdx, cIdx) {
		return v.sp, v.stack, nil
	}
	return sp, v.stack, nil
}

func (v *vm) runIntSliceGetFallbackRecovered(
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, vIdx uint16,
) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
			panicked = true
		}
	}()
	ptr := indexAddressValue(locals[sIdx], int(intLocals[jIdx]))
	ret := dereferenceValue(ptr)
	intLocals[vIdx] = ret.RawInt()
	locals[vIdx] = ret
	return false
}

func (v *vm) runIntSliceSetFallbackRecovered(
	locals []value.Value,
	intLocals []int64,
	sIdx, jIdx, valIdx uint16,
) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
			panicked = true
		}
	}()
	ptr := indexAddressValue(locals[sIdx], int(intLocals[jIdx]))
	v.setDereferenceValue(ptr, value.MakeInt(intLocals[valIdx]))
	return false
}

func (v *vm) runIntSliceSetConstFallbackRecovered(
	locals []value.Value,
	intLocals []int64,
	prebaked []value.Value,
	sIdx, jIdx, cIdx uint16,
) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
			panicked = true
		}
	}()
	ptr := indexAddressValue(locals[sIdx], int(intLocals[jIdx]))
	v.setDereferenceValue(ptr, prebaked[cIdx])
	return false
}
