package vm

import (
	"fmt"

	"github.com/t04dJ14n9/gig/model/value"
)

type runDisposition uint8

const (
	runContinue runDisposition = iota
	runReturn
)

type runFrameReturnResult struct {
	sp        int
	done      bool
	frame     *Frame
	ins       []byte
	locals    []value.Value
	intLocals []int64
}

type runPanicStepResult struct {
	sp        int
	done      bool
	retVal    value.Value
	err       error
	frame     *Frame
	ins       []byte
	locals    []value.Value
	intLocals []int64
}

func (v *vm) runFrameReturn(frame *Frame, stack []value.Value, sp int, retVal value.Value) runFrameReturnResult {
	v.fpool.put(frame)
	v.fp--

	// Deferred functions are executed by a child run loop. Returning here hands
	// control back to the caller instead of continuing in the parent frame.
	if v.deferDepth > 0 {
		return runFrameReturnResult{sp: sp, done: true}
	}

	var next *Frame
	if v.fp > 0 {
		next = v.frames[v.fp-1]
		sp = next.basePtr
	}
	stack[sp] = retVal
	sp++

	if next == nil {
		return runFrameReturnResult{sp: sp}
	}
	return runFrameReturnResult{
		sp:        sp,
		frame:     next,
		ins:       next.fn.Instructions,
		locals:    next.locals,
		intLocals: next.intLocals,
	}
}

func (v *vm) runPanicStep(frame *Frame, sp int) runPanicStepResult {
	// Panic unwinding is cold, but it mutates the frame stack and therefore must
	// return refreshed cached fields for the main loop before execution resumes.
	disposition, nextSP, retVal, err := v.handlePendingPanic(frame, sp)
	result := runPanicStepResult{
		sp:     nextSP,
		retVal: retVal,
		err:    err,
	}
	if err != nil || disposition == runReturn {
		result.done = true
		return result
	}
	if v.fp == 0 {
		return result
	}

	next := v.frames[v.fp-1]
	result.frame = next
	result.ins = next.fn.Instructions
	result.locals = next.locals
	result.intLocals = next.intLocals
	return result
}

func (v *vm) runFinalStackValue(stack []value.Value, sp int) value.Value {
	v.sp = sp
	if sp == 0 {
		return value.MakeNil()
	}
	sp--
	v.sp = sp
	return stack[sp]
}

func (v *vm) handlePendingPanic(frame *Frame, sp int) (runDisposition, int, value.Value, error) {
	// Sync sp so runDefersDuringPanic can use v.sp for recursive run() calls.
	v.sp = sp
	recovered := v.runDefersDuringPanic(frame)
	sp = v.sp

	if recovered || !v.panicking {
		retVal := recoveredPanicReturnValue(frame)
		v.fpool.put(frame)
		v.fp--

		if v.deferDepth > 0 {
			v.sp = sp
			return runReturn, sp, retVal, nil
		}
		if v.fp > 0 {
			sp = v.frames[v.fp-1].basePtr
		}
		v.stack[sp] = retVal
		sp++
		return runContinue, sp, value.Value{}, nil
	}

	if v.fp == 1 {
		// Preserve the original typed panic value before clearing, so callers
		// such as OpRunDefers can recover it without parsing an error string.
		v.lastPanicVal = v.panicVal
		err := fmt.Errorf("panic: %v", v.panicVal.Interface())
		v.panicking = false
		v.panicVal = value.MakeNil()
		return runReturn, sp, value.MakeNil(), err
	}

	if v.deferDepth > 0 {
		v.fp--
		v.fpool.put(frame)
		// Keep v.panicking set so the caller's runDefersDuringPanic sees it.
		return runReturn, sp, value.MakeNil(), nil
	}

	v.fp--
	v.fpool.put(frame)
	return runContinue, sp, value.Value{}, nil
}

func recoveredPanicReturnValue(frame *Frame) value.Value {
	slots := frame.fn.ResultAllocSlots
	if len(slots) == 0 {
		return value.MakeNil()
	}
	if len(slots) == 1 {
		return derefAllocLocal(frame.locals[slots[0]])
	}
	results := make([]value.Value, len(slots))
	for i, slot := range slots {
		results[i] = derefAllocLocal(frame.locals[slot])
	}
	return value.FromInterface(results)
}

func (v *vm) handleFrameEnd(frame *Frame) runDisposition {
	v.fp--
	v.fpool.put(frame)
	if v.deferDepth > 0 {
		return runReturn
	}
	return runContinue
}
