// run.go contains the main fetch-decode-execute loop with hot-path inlined instructions.
package vm

import (
	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// run is the main execution loop for the VM.
// It fetches, decodes, and executes bytecode instructions until:
//   - All call frames return (normal termination)
//   - Context is cancelled (timeout/cancellation)
//   - A panic propagates to the top frame (error return)
//
// Hot-path instructions (OpLocal, OpSetLocal, OpConst, arithmetic, comparisons,
// jumps) are inlined directly in the loop to avoid per-instruction function call
// overhead. Less frequent opcodes fall through to executeOp.
//

//nolint:gocyclo,cyclop,funlen,maintidx,gocognit
func (v *vm) run() (value.Value, error) {
	// Hoist hot fields into local variables for better register allocation.
	// The Go compiler can keep these in CPU registers across iterations,
	// avoiding repeated loads from v.* on each instruction.
	stack := v.stack
	sp := v.sp
	prebaked := v.program.PrebakedConstants

	// Cache current frame state to avoid re-reading from v.frames[] each iteration.
	// These are only invalidated on call/return/executeOp.
	var frame *Frame
	var ins []byte
	var locals []value.Value
	var intLocals []int64
	intConsts := v.program.IntConstants

	// loadFrame caches the current frame's hot fields into local variables.
	loadFrame := func() {
		frame = v.frames[v.fp-1]
		ins = frame.fn.Instructions
		locals = frame.locals
		intLocals = frame.intLocals
	}

	// readU16 reads a 2-byte big-endian operand from the cached ins slice.
	// This is faster than frame.readUint16() which dereferences frame.fn.Instructions.
	readU16 := func() uint16 {
		v := uint16(ins[frame.ip])<<8 | uint16(ins[frame.ip+1])
		frame.ip += 2
		return v
	}

	// instructionCount tracks total instructions executed for periodic context checks.
	instructionCount := uint64(0)

	if v.fp > 0 {
		loadFrame()
	}

	for v.fp > 0 {
		// Periodic context check counter
		instructionCount++
		if instructionCount&contextCheckMask == 0 {
			select {
			case <-v.ctx.Done():
				v.sp = sp
				return value.MakeNil(), v.ctx.Err()
			default:
			}
		}

		// Handle panic FIRST (before end-of-function check)
		// This is critical: when a function panics as its last instruction,
		// we need to run deferred functions before the frame is popped.
		// Allow panic handling at any defer depth — this enables nested panics
		// (panic inside a deferred function) to be properly recovered.
		if v.panicking {
			ret := v.runPanicStep(frame, sp)
			sp = ret.sp
			stack = v.stack
			if ret.done {
				return ret.retVal, ret.err
			}
			frame, ins, locals, intLocals = ret.frame, ret.ins, ret.locals, ret.intLocals
			continue
		}

		// Check for end of function
		if frame.ip >= len(ins) {
			ret := v.runFrameEndStep(frame)
			if ret.done {
				return value.MakeNil(), nil
			}
			frame, ins, locals, intLocals = ret.frame, ret.ins, ret.locals, ret.intLocals
			continue
		}

		// Fetch opcode
		op := bytecode.OpCode(ins[frame.ip])
		frame.ip++

		// Inline hot-path instructions to eliminate per-instruction function call overhead.
		// These opcodes cover >90% of instructions in typical numeric programs.
		// Instructions handled here use 'continue' to skip the executeOp call below.
		switch op { //nolint:exhaustive
		case bytecode.OpLocal:
			idx := readU16()
			stack[sp] = locals[idx]
			sp++
			continue

		case bytecode.OpSetLocal:
			idx := readU16()
			sp--
			locals[idx] = stack[sp]
			continue

		case bytecode.OpConst:
			idx := readU16()
			if int(idx) < len(prebaked) {
				stack[sp] = prebaked[idx]
			} else {
				stack[sp] = v.runSlowConst(idx)
			}
			sp++
			continue

		case bytecode.OpAdd:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Add(b)
			}
			sp++
			continue

		case bytecode.OpSub:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Sub(b)
			}
			sp++
			continue

		case bytecode.OpMul:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Mul(b)
			}
			sp++
			continue

		case bytecode.OpLess:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() < b.RawInt())
			} else {
				stack[sp] = value.MakeBool(a.Cmp(b) < 0)
			}
			sp++
			continue

		case bytecode.OpLessEq:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() <= b.RawInt())
			} else {
				stack[sp] = value.MakeBool(lessEqCmp(a, b))
			}
			sp++
			continue

		case bytecode.OpGreater:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() > b.RawInt())
			} else {
				stack[sp] = value.MakeBool(a.Cmp(b) > 0)
			}
			sp++
			continue

		case bytecode.OpGreaterEq:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runBothInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() >= b.RawInt())
			} else {
				stack[sp] = value.MakeBool(greaterEqCmp(a, b))
			}
			sp++
			continue

		case bytecode.OpEqual:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runSameSizedInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() == b.RawInt())
			} else {
				stack[sp] = value.MakeBool(a.Equal(b))
			}
			sp++
			continue

		case bytecode.OpNotEqual:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if runSameSizedInts(a, b) {
				stack[sp] = value.MakeBool(a.RawInt() != b.RawInt())
			} else {
				stack[sp] = value.MakeBool(!a.Equal(b))
			}
			sp++
			continue

		case bytecode.OpJump:
			offset := readU16()
			frame.ip = int(offset)
			continue

		case bytecode.OpJumpTrue:
			offset := readU16()
			sp--
			runJumpIf(frame, offset, stack[sp].RawBool())
			continue

		case bytecode.OpJumpFalse:
			offset := readU16()
			sp--
			runJumpIf(frame, offset, !stack[sp].RawBool())
			continue

		case bytecode.OpNot:
			sp--
			stack[sp] = value.MakeBool(!stack[sp].RawBool())
			sp++
			continue

		case bytecode.OpNil:
			stack[sp] = value.MakeNil()
			sp++
			continue

		case bytecode.OpTrue:
			stack[sp] = value.MakeBool(true)
			sp++
			continue

		case bytecode.OpFalse:
			stack[sp] = value.MakeBool(false)
			sp++
			continue

		case bytecode.OpPop:
			sp--
			continue

		case bytecode.OpDup:
			stack[sp] = stack[sp-1]
			sp++
			continue

		case bytecode.OpCall:
			funcIdx := readU16()
			numArgs := frame.readByte()
			v.sp = sp
			v.callCompiledFunction(int(funcIdx), int(numArgs))
			sp = v.sp
			stack = v.stack
			loadFrame()
			continue

		case bytecode.OpReturn, bytecode.OpReturnVal:
			retVal := value.MakeNil()
			if op == bytecode.OpReturnVal {
				sp--
				retVal = stack[sp]
			}
			ret := v.runFrameReturn(frame, stack, sp, retVal)
			if ret.done {
				return retVal, nil
			}
			sp, frame, ins, locals, intLocals = ret.sp, ret.frame, ret.ins, ret.locals, ret.intLocals
			continue

		case bytecode.OpSetDeref:
			sp = v.runSetDeref(sp)
			continue

		case bytecode.OpIndexAddr:
			var err error
			sp, stack, err = v.runIndexAddr(frame, sp)
			reloadFrame, err := v.runInlineStackOpComplete(err)
			if err != nil {
				return value.MakeNil(), err
			}
			if reloadFrame {
				loadFrame()
			}
			continue

		case bytecode.OpDeref:
			var err error
			sp, stack, err = v.runDeref(frame, sp)
			reloadFrame, err := v.runInlineStackOpComplete(err)
			if err != nil {
				return value.MakeNil(), err
			}
			if reloadFrame {
				loadFrame()
			}
			continue

		case bytecode.OpLen:
			var err error
			sp, stack, err = v.runLen(frame, sp)
			reloadFrame, err := v.runInlineStackOpComplete(err)
			if err != nil {
				return value.MakeNil(), err
			}
			if reloadFrame {
				loadFrame()
			}
			continue

			// ========================================
			// Superinstructions: fused ops for hot loops
			// ========================================

		case bytecode.OpAddLocalLocal,
			bytecode.OpSubLocalLocal,
			bytecode.OpMulLocalLocal,
			bytecode.OpAddLocalConst,
			bytecode.OpSubLocalConst,
			bytecode.OpLessLocalLocalJumpTrue,
			bytecode.OpLessLocalConstJumpTrue,
			bytecode.OpLessEqLocalConstJumpTrue,
			bytecode.OpGreaterLocalLocalJumpTrue,
			bytecode.OpLessLocalLocalJumpFalse,
			bytecode.OpLessLocalConstJumpFalse,
			bytecode.OpLessEqLocalConstJumpFalse,
			bytecode.OpAddSetLocal,
			bytecode.OpSubSetLocal,
			bytecode.OpLocalLocalAddSetLocal,
			bytecode.OpLocalConstAddSetLocal,
			bytecode.OpLocalConstSubSetLocal,
			bytecode.OpLocalLocalSubSetLocal,
			bytecode.OpLocalLocalMulSetLocal,
			bytecode.OpLocalConstMulSetLocal:
			sp = v.runGenericSuperinstruction(op, frame, sp, locals, intLocals, prebaked)
			continue

		// ========================================
		// Integer-specialized superinstructions
		// Operate on intLocals []int64 directly (8 bytes vs 32 bytes per op)
		// ========================================

		case bytecode.OpIntLocalConstAddSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] + intConsts[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLocalConstSubSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] - intConsts[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLocalLocalAddSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] + intLocals[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLocalLocalSubSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] - intLocals[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLocalLocalMulSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] * intLocals[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLocalConstMulSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			r := intLocals[idxA] * intConsts[idxB]
			intLocals[idxC] = r
			locals[idxC] = value.MakeInt(r)
			continue

		case bytecode.OpIntLessLocalConstJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] >= intConsts[idxB])
			continue

		case bytecode.OpIntLessEqLocalConstJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] <= intConsts[idxB])
			continue

		case bytecode.OpIntLessEqLocalConstJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] > intConsts[idxB])
			continue

		case bytecode.OpIntLessLocalLocalJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] >= intLocals[idxB])
			continue

		case bytecode.OpIntGreaterLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] > intLocals[idxB])
			continue

		case bytecode.OpIntSetLocal:
			idx := readU16()
			sp--
			v := stack[sp]
			intLocals[idx] = v.RawInt()
			locals[idx] = v
			continue

		case bytecode.OpIntLocal:
			idx := readU16()
			stack[sp] = value.MakeInt(intLocals[idx])
			sp++
			continue

		case bytecode.OpIntLessLocalConstJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] < intConsts[idxB])
			continue

		case bytecode.OpIntLessLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			runJumpIf(frame, offset, intLocals[idxA] < intLocals[idxB])
			continue

		case bytecode.OpIntMoveLocal:
			src := readU16()
			dst := readU16()
			intLocals[dst] = intLocals[src]
			locals[dst] = locals[src]
			continue

		case bytecode.OpIntSliceGet:
			sIdx := readU16()
			jIdx := readU16()
			vIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				idx := intLocals[jIdx]
				if idx < 0 || idx >= int64(len(s)) {
					v.setIntSliceIndexPanic(idx, len(s))
					continue
				}
				r := s[idx]
				intLocals[vIdx] = r
				locals[vIdx] = value.MakeInt(r)
			} else {
				var err error
				sp, stack, err = v.runIntSliceGetFallback(frame, locals, intLocals, sIdx, jIdx, vIdx, sp)
				if err != nil {
					return value.MakeNil(), err
				}
			}
			continue

		case bytecode.OpIntSliceSet:
			sIdx := readU16()
			jIdx := readU16()
			valIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				idx := intLocals[jIdx]
				if idx < 0 || idx >= int64(len(s)) {
					v.setIntSliceIndexPanic(idx, len(s))
					continue
				}
				s[idx] = intLocals[valIdx]
			} else {
				var err error
				sp, stack, err = v.runIntSliceSetFallback(frame, locals, intLocals, sIdx, jIdx, valIdx, sp)
				if err != nil {
					return value.MakeNil(), err
				}
			}
			continue

		case bytecode.OpIntSliceSetConst:
			sIdx := readU16()
			jIdx := readU16()
			cIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				idx := intLocals[jIdx]
				if idx < 0 || idx >= int64(len(s)) {
					v.setIntSliceIndexPanic(idx, len(s))
					continue
				}
				s[idx] = intConsts[cIdx]
			} else {
				var err error
				sp, stack, err = v.runIntSliceSetConstFallback(frame, locals, intLocals, prebaked, sIdx, jIdx, cIdx, sp)
				if err != nil {
					return value.MakeNil(), err
				}
			}
			continue

		case bytecode.OpCallExternal:
			funcIdx := readU16()
			numArgs := int(frame.readByte())
			frameChanged := false
			var err error
			sp, stack, frameChanged, err = v.runExternalCall(int(funcIdx), numArgs, sp)
			reloadFrame, err := v.runCallComplete(err, frameChanged, true)
			if err != nil {
				return value.MakeNil(), err
			}
			if reloadFrame {
				loadFrame()
			}
			continue

		case bytecode.OpCallIndirect:
			numArgs := int(frame.readByte())
			frameChanged := false
			var err error
			sp, stack, frameChanged, err = v.runIndirectCall(sp, numArgs)
			reloadFrame, err := v.runCallComplete(err, frameChanged, false)
			if err != nil {
				return value.MakeNil(), err
			}
			if reloadFrame {
				loadFrame()
			}
			continue

		default:
			// Fall through to executeOp for all other opcodes
		}

		// Non-hot-path: dispatch to the full handler.
		var reloadFrame bool
		var err error
		sp, stack, reloadFrame, err = v.runColdOp(frame, op, sp)
		if err != nil {
			return value.MakeNil(), err
		}
		if reloadFrame {
			loadFrame()
		}
		continue
	}

	// Return top of stack (or nil if empty)
	return v.runFinalStackValue(stack, sp), nil
}
