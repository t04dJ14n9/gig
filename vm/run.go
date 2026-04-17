// run.go contains the main fetch-decode-execute loop with hot-path inlined instructions.
package vm

import (
	"fmt"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
)

// derefAllocLocal dereferences an Alloc pointer stored in a frame's local slot.
// In SSA, named return variables and captured locals are represented as pointer
// allocations (OpNew). This function reads the value behind the pointer, mirroring
// what OpDeref does in the normal return path.
func derefAllocLocal(ptr value.Value) value.Value {
	switch ptr.Kind() {
	case value.KindPointer:
		return ptr.Elem()
	case value.KindInterface:
		if rv, ok := ptr.ReflectValue(); ok {
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				return value.MakeFromReflect(rv.Elem())
			}
		}
	}
	// If it's a *value.Value wrapper (from OpNew for function types)
	if rv, ok := ptr.ReflectValue(); ok {
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			return value.MakeFromReflect(rv.Elem())
		}
	}
	return ptr
}

// runDefersDuringPanic runs deferred functions in LIFO order during a panic.
// It uses the shared VM + recursive run() so that OpRecover (inside the deferred
// function) can clear v.panicking on the same VM instance.
// Returns true if the panic was recovered by a deferred function.
//
// Caller must sync v.sp before calling; v.sp is updated on return.
func (v *vm) runDefersDuringPanic(frame *Frame) bool {
	recovered := false
	if len(frame.defers) == 0 {
		return recovered
	}

	for i := len(frame.defers) - 1; i >= 0; i-- {
		d := frame.defers[i]

		// Handle external method defers (OpDeferExternal, e.g. defer mu.Unlock())
		// These must run during panic recovery, just like in Go.
		if d.externalInfo != nil {
			if methodInfo, ok := d.externalInfo.(*external.ExternalMethodInfo); ok {
				_ = v.callExternalMethod(methodInfo, d.args)
				if methodInfo.DirectCall == nil {
					_ = v.pop()
				}
			}
			continue
		}

		// Handle external function value defers (OpDeferIndirect with external func)
		// e.g. defer fn() where fn is an external function variable
		if d.externalFunc.IsValid() {
			argVals := make([]reflect.Value, len(d.args))
			for j, arg := range d.args {
				funcType := d.externalFunc.Type()
				if j < funcType.NumIn() {
					argType := funcType.In(j)
					argVals[j] = arg.ToReflectValue(argType)
				} else {
					argVals[j] = reflect.ValueOf(arg.Interface())
				}
			}
			d.externalFunc.Call(argVals)
			continue
		}

		// Skip nil function entries (shouldn't happen, but defensive)
		if d.fn == nil {
			continue
		}

		// Get free variables from closure if present
		var freeVars []*value.Value
		if d.closure != nil {
			freeVars = d.closure.FreeVars
		}

		if v.panicking {
			// Panic is active: push state onto panicStack so that
			// OpRecover (inside the deferred function) can access it.
			// Clear v.panicking so the recursive run() doesn't immediately
			// re-enter the panic handler on the defer's frame.
			v.panicStack = append(v.panicStack, panicState{
				panicking: true,
				panicVal:  v.panicVal,
			})
			v.panicking = false
			v.panicVal = value.MakeNil()

			v.callFunction(d.fn, d.args, freeVars)
			v.deferDepth++
			_, _ = v.run()
			v.deferDepth--

			if v.panicking {
				// The defer itself panicked (and was not recovered).
				// Pop the saved state — this new panic replaces the old one.
				v.panicStack = v.panicStack[:len(v.panicStack)-1]
				continue
			}

			// Pop the saved state and check if recover() consumed it.
			saved := v.panicStack[len(v.panicStack)-1]
			v.panicStack = v.panicStack[:len(v.panicStack)-1]

			if saved.panicking {
				// The defer didn't call recover() — restore the panic.
				v.panicking = true
				v.panicVal = saved.panicVal
			} else {
				// recover() was called — panic is resolved.
				// Continue running remaining defers in normal mode.
				recovered = true
			}
		} else {
			// Panic already recovered: run remaining defers normally.
			// Use child VM to avoid interfering with the parent frame stack.
			childVM := v.newDeferVM()
			deferFrame := newFrame(d.fn, d.args, freeVars)
			childVM.frames[0] = deferFrame
			childVM.fp = 1
			_, _ = childVM.run()

			// If the child VM panicked, re-enter panic mode.
			if childVM.lastPanicVal.IsValid() {
				v.panicking = true
				v.panicVal = childVM.lastPanicVal
				recovered = false
			} else if childVM.panicking {
				v.panicking = true
				v.panicVal = childVM.panicVal
				recovered = false
			}
		}
	}
	frame.defers = nil
	return recovered
}

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
			// Sync sp so runDefersDuringPanic can use v.sp for recursive run() calls.
			v.sp = sp
			recovered := v.runDefersDuringPanic(frame)
			sp = v.sp

			// Check if panic was recovered during deferred execution
			if recovered || !v.panicking {
				// Panic was recovered — return value from this frame.
				// If the function has ResultAllocSlots (named returns),
				// deref those Alloc pointers to get the value that deferred closures
				// may have written. Otherwise fall back to nil (zero value).
				retVal := value.MakeNil()
				if slots := frame.fn.ResultAllocSlots; len(slots) > 0 {
					if len(slots) == 1 {
						// Single result: deref the Alloc pointer in the local slot
						ptr := frame.locals[slots[0]]
						retVal = derefAllocLocal(ptr)
					} else {
						// Multiple results: pack them
						results := make([]value.Value, len(slots))
						for i, slot := range slots {
							results[i] = derefAllocLocal(frame.locals[slot])
						}
						retVal = value.FromInterface(results)
					}
				}
				v.fpool.put(frame)
				v.fp--
				// If running inside a deferred function (deferDepth > 0),
				// return immediately. Don't continue executing the outer
				// function's frame — that's handled by the caller's run().
				if v.deferDepth > 0 {
					v.sp = sp
					return retVal, nil
				}
				if v.fp > 0 {
					loadFrame()
					sp = frame.basePtr
				}
				stack[sp] = retVal
				sp++
				continue
			}

			// If this is the last frame, return the panic as an error
			if v.fp == 1 {
				// Preserve the original typed panic value before clearing,
				// so callers (e.g. OpRunDefers) can recover it instead of
				// parsing the error string (which loses type information).
				v.lastPanicVal = v.panicVal
				err := fmt.Errorf("panic: %v", v.panicVal.Interface())
				v.panicking = false
				v.panicVal = value.MakeNil()
				return value.MakeNil(), err
			}

			// If running inside a deferred function (deferDepth > 0) and panic wasn't recovered,
			// return immediately to let the outer runDefersDuringPanic handle it.
			// This prevents re-running defers on the same frame after a nested panic.
			// Must pop the defer's frame before returning.
			if v.deferDepth > 0 {
				v.fp--
				v.fpool.put(frame)
				// Don't clear v.panicking - let the caller's runDefersDuringPanic see it
				return value.MakeNil(), nil
			}

			// Propagate panic to caller
			v.fp--
			v.fpool.put(frame)
			if v.fp > 0 {
				loadFrame()
			}
			continue
		}

		// Check for end of function
		if frame.ip >= len(ins) {
			// If running a deferred function, return immediately after the function ends.
			// Don't continue with the caller's frame - that's handled by the outer run().
			if v.deferDepth > 0 {
				v.fp--
				v.fpool.put(frame)
				return value.MakeNil(), nil
			}
			// Pop frame and return it to pool
			v.fp--
			v.fpool.put(frame)
			if v.fp > 0 {
				loadFrame()
			}
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
			} else if int(idx) < len(v.program.Constants) {
				stack[sp] = value.FromInterface(v.program.Constants[idx])
			}
			sp++
			continue

		case bytecode.OpAdd:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeBool(a.RawInt() <= b.RawInt())
			} else {
				stack[sp] = value.MakeBool(a.Cmp(b) <= 0)
			}
			sp++
			continue

		case bytecode.OpGreater:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeBool(a.RawInt() >= b.RawInt())
			} else {
				stack[sp] = value.MakeBool(a.Cmp(b) >= 0)
			}
			sp++
			continue

		case bytecode.OpEqual:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
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
			if stack[sp].RawBool() {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpJumpFalse:
			offset := readU16()
			sp--
			if !stack[sp].RawBool() {
				frame.ip = int(offset)
			}
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

		case bytecode.OpReturn:
			v.fpool.put(frame)
			v.fp--
			// If running a deferred function, return immediately.
			// Don't continue with the caller's frame - that's handled by the outer run().
			if v.deferDepth > 0 {
				return value.MakeNil(), nil
			}
			if v.fp > 0 {
				loadFrame()
				sp = frame.basePtr
			}
			stack[sp] = value.MakeNil()
			sp++
			continue

		case bytecode.OpReturnVal:
			sp--
			retVal := stack[sp]
			v.fpool.put(frame)
			v.fp--
			// If running a deferred function, return immediately.
			if v.deferDepth > 0 {
				return retVal, nil
			}
			if v.fp > 0 {
				loadFrame()
				sp = frame.basePtr
			}
			stack[sp] = retVal
			sp++
			continue

		case bytecode.OpSetDeref:
			sp--
			val := stack[sp]
			sp--
			ptr := stack[sp]
			// Nil pointer dereference check
			if ptr.IsNil() || !ptr.IsValid() {
				v.panicking = true
				v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
				continue
			}
			// Fast path: *int64 pointer (from native int slice OpIndexAddr)
			if p, ok := ptr.IntPtr(); ok {
				*p = val.RawInt()
			} else if iface := ptr.Interface(); iface != nil {
				// GlobalRef from shared-mode OpGlobal — use locked write
				if ref, ok := iface.(*GlobalRef); ok {
					ref.Store(val)
				} else {
					ptr.SetElem(val)
				}
			} else {
				ptr.SetElem(val)
			}
			continue

		case bytecode.OpIndexAddr:
			sp--
			index := stack[sp]
			sp--
			container := stack[sp]
			// Fast path: native []int64 slice (covers make([]int, N) in interpreted code)
			if s, ok := container.IntSlice(); ok {
				idx := index.RawInt()
				if idx < 0 || idx >= int64(len(s)) {
					// Bounds check failed — convert to VM panic so guest recover() can catch it.
					v.panicking = true
					v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, len(s)))
					continue
				}
				stack[sp] = value.MakeIntPtr(&s[idx])
				sp++
				continue
			}
			// Slow path: go through executeOp
			v.sp = sp
			v.push(container)
			v.push(index)
			if err := v.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			if v.panicking {
				continue
			}
			sp = v.sp
			stack = v.stack
			if v.fp > 0 {
				loadFrame()
			}
			continue

		case bytecode.OpDeref:
			sp--
			ptr := stack[sp]
			// Nil pointer dereference check for reflect-based pointers.
			// IntPtr fast path handles native *int64 pointers.
			if ptr.Kind() == value.KindReflect || ptr.Kind() == value.KindPointer {
				if ptr.IsNil() {
					v.panicking = true
					v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
					continue
				}
			}
			// Fast path: *int64 pointer (from native int slice OpIndexAddr)
			if p, ok := ptr.IntPtr(); ok {
				stack[sp] = value.MakeInt(*p)
				sp++
				continue
			}
			// Slow path: go through executeOp
			v.sp = sp
			v.push(ptr)
			if err := v.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			if v.panicking {
				continue
			}
			sp = v.sp
			stack = v.stack
			if v.fp > 0 {
				loadFrame()
			}
			continue

		case bytecode.OpLen:
			sp--
			obj := stack[sp]
			switch obj.Kind() {
			case value.KindSlice:
				stack[sp] = value.MakeInt(int64(obj.Len()))
				sp++
				continue
			case value.KindString:
				stack[sp] = value.MakeInt(int64(len(obj.String())))
				sp++
				continue
			}
			// Slow path
			v.sp = sp
			v.push(obj)
			if err := v.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			if v.panicking {
				continue
			}
			sp = v.sp
			stack = v.stack
			if v.fp > 0 {
				loadFrame()
			}
			continue

			// ========================================
			// Superinstructions: fused ops for hot loops
			// ========================================

		case bytecode.OpAddLocalLocal:
			idxA := readU16()
			idxB := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Add(b)
			}
			sp++
			continue

		case bytecode.OpSubLocalLocal:
			idxA := readU16()
			idxB := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Sub(b)
			}
			sp++
			continue

		case bytecode.OpMulLocalLocal:
			idxA := readU16()
			idxB := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Mul(b)
			}
			sp++
			continue

		case bytecode.OpAddLocalConst:
			idxA := readU16()
			idxB := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Add(b)
			}
			sp++
			continue

		case bytecode.OpSubLocalConst:
			idxA := readU16()
			idxB := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
			} else {
				stack[sp] = a.Sub(b)
			}
			sp++
			continue

		case bytecode.OpLessLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() < b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) < 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpLessLocalConstJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() < b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) < 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpLessEqLocalConstJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() <= b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) <= 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpGreaterLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() > b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) > 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpLessLocalLocalJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() >= b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) >= 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpLessLocalConstJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() >= b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) >= 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpLessEqLocalConstJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				if a.RawInt() > b.RawInt() {
					frame.ip = int(offset)
				}
			} else {
				if a.Cmp(b) > 0 {
					frame.ip = int(offset)
				}
			}
			continue

		case bytecode.OpAddSetLocal:
			idx := readU16()
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				r := a.RawInt() + b.RawInt()
				locals[idx] = value.MakeIntSized(r, a.RawSize())
				if intLocals != nil {
					intLocals[idx] = r
				}
			} else {
				locals[idx] = a.Add(b)
			}
			continue

		case bytecode.OpSubSetLocal:
			idx := readU16()
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				r := a.RawInt() - b.RawInt()
				locals[idx] = value.MakeIntSized(r, a.RawSize())
				if intLocals != nil {
					intLocals[idx] = r
				}
			} else {
				locals[idx] = a.Sub(b)
			}
			continue

		case bytecode.OpLocalLocalAddSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Add(b)
			}
			continue

		case bytecode.OpLocalConstAddSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()+b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Add(b)
			}
			continue

		case bytecode.OpLocalConstSubSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Sub(b)
			}
			continue

		case bytecode.OpLocalLocalSubSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()-b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Sub(b)
			}
			continue

		case bytecode.OpLocalLocalMulSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Mul(b)
			}
			continue

		case bytecode.OpLocalConstMulSetLocal:
			idxA := readU16()
			idxB := readU16()
			idxC := readU16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeIntSized(a.RawInt()*b.RawInt(), a.RawSize())
			} else {
				locals[idxC] = a.Mul(b)
			}
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
			if intLocals[idxA] >= intConsts[idxB] {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpIntLessEqLocalConstJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			if intLocals[idxA] <= intConsts[idxB] {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpIntLessEqLocalConstJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			if intLocals[idxA] > intConsts[idxB] {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpIntLessLocalLocalJumpFalse:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			if intLocals[idxA] >= intLocals[idxB] {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpIntGreaterLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			if intLocals[idxA] > intLocals[idxB] {
				frame.ip = int(offset)
			}
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
			if intLocals[idxA] < intConsts[idxB] {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpIntLessLocalLocalJumpTrue:
			idxA := readU16()
			idxB := readU16()
			offset := readU16()
			if intLocals[idxA] < intLocals[idxB] {
				frame.ip = int(offset)
			}
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
					v.panicking = true
					v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, len(s)))
					continue
				}
				r := s[idx]
				intLocals[vIdx] = r
				locals[vIdx] = value.MakeInt(r)
			} else {
				// Fallback: execute as IndexAddr + Deref manually
				v.sp = sp
				v.push(locals[sIdx])
				v.push(value.MakeInt(intLocals[jIdx]))
				if err := v.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				if err := v.executeOp(bytecode.OpDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				ret := v.pop()
				intLocals[vIdx] = ret.RawInt()
				locals[vIdx] = ret
				sp = v.sp
				stack = v.stack
			}
			continue

		case bytecode.OpIntSliceSet:
			sIdx := readU16()
			jIdx := readU16()
			valIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				idx := intLocals[jIdx]
				if idx < 0 || idx >= int64(len(s)) {
					v.panicking = true
					v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, len(s)))
					continue
				}
				s[idx] = intLocals[valIdx]
			} else {
				// Fallback: execute as IndexAddr + SetDeref manually
				v.sp = sp
				v.push(locals[sIdx])
				v.push(value.MakeInt(intLocals[jIdx]))
				if err := v.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				v.push(value.MakeInt(intLocals[valIdx]))
				if err := v.executeOp(bytecode.OpSetDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				sp = v.sp
				stack = v.stack
			}
			continue

		case bytecode.OpIntSliceSetConst:
			sIdx := readU16()
			jIdx := readU16()
			cIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				idx := intLocals[jIdx]
				if idx < 0 || idx >= int64(len(s)) {
					v.panicking = true
					v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, len(s)))
					continue
				}
				s[idx] = intConsts[cIdx]
			} else {
				// Fallback: execute as IndexAddr + SetDeref manually
				v.sp = sp
				v.push(locals[sIdx])
				v.push(value.MakeInt(intLocals[jIdx]))
				if err := v.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				v.push(prebaked[cIdx])
				if err := v.executeOp(bytecode.OpSetDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				if v.panicking {
					continue
				}
				sp = v.sp
				stack = v.stack
			}
			continue

		case bytecode.OpCallExternal:
			funcIdx := readU16()
			numArgs := int(frame.readByte())
			prevFP := v.fp
			v.sp = sp
			if err := v.callExternal(int(funcIdx), numArgs); err != nil {
				return value.MakeNil(), err
			}
			sp = v.sp
			stack = v.stack
			// If callExternal pushed a new compiled frame (e.g., compiled method
			// dispatch for invoke calls), reload frame state so the main loop
			// executes the new frame before continuing.
			if v.fp != prevFP {
				loadFrame()
			}
			continue

		case bytecode.OpCallIndirect:
			numArgs := int(frame.readByte())
			// Pop arguments using stack-allocated buffer to avoid heap allocation
			var argsBuf [8]value.Value
			var args []value.Value
			if numArgs <= len(argsBuf) {
				args = argsBuf[:numArgs]
			} else {
				args = make([]value.Value, numArgs)
			}
			spLocal := sp
			for i := numArgs - 1; i >= 0; i-- {
				spLocal--
				args[i] = stack[spLocal]
			}
			// Pop the callee
			spLocal--
			callee := stack[spLocal]
			sp = spLocal
			// Fast path: direct obj type assertion for *Closure avoids Interface() overhead
			if closure, ok := callee.RawObj().(*Closure); ok {
				v.sp = sp
				v.callFunction(closure.Fn, args, closure.FreeVars)
				sp = v.sp
				stack = v.stack
				loadFrame()
		} else if rv, ok := callee.ReflectValue(); ok && rv.Kind() == reflect.Func {
			// Nil function call: trigger VM panic so guest recover() can catch it
			if rv.IsNil() {
				v.sp = sp
				v.panicking = true
				v.panicVal = value.FromInterface("invalid memory address or nil pointer dereference")
				continue
			}
			// Reflect-based function call with panic safety:
			// catch Go-level panics from external code and convert to VM panics
			// so guest recover() can catch them.
			in := make([]reflect.Value, numArgs)
			fnType := rv.Type()
			for i := 0; i < numArgs; i++ {
				if i < fnType.NumIn() {
					in[i] = args[i].ToReflectValue(fnType.In(i))
				}
			}
			var out []reflect.Value
			func() {
				defer func() {
					if r := recover(); r != nil {
						v.sp = sp
						v.panicking = true
						v.panicVal = value.FromInterface(r)
					}
				}()
				out = rv.Call(in)
			}()
			if v.panicking {
				continue
			}
			if len(out) == 0 {
				stack[sp] = value.MakeNil()
			} else {
				stack[sp] = value.MakeFromReflect(out[0])
			}
			sp++
		} else {
			stack[sp] = value.MakeNil()
			sp++
		}
		continue

		default:
			// Fall through to executeOp for all other opcodes
		}

		// Non-hot-path: dispatch to the full handler.
		// Sync sp back before calling executeOp (it uses v.push/v.pop).
		v.sp = sp
		if err := v.executeOp(op, frame); err != nil {
			return value.MakeNil(), err
		}
		if v.panicking {
			sp = v.sp
			stack = v.stack
			if v.fp > 0 {
				loadFrame()
			}
			continue
		}
		sp = v.sp
		stack = v.stack
		// Reload frame state in case executeOp changed it (call/return within executeOp)
		if v.fp > 0 {
			loadFrame()
		}
	}

	// Return top of stack (or nil if empty)
	v.sp = sp
	if sp > 0 {
		sp--
		v.sp = sp
		return stack[sp], nil
	}
	return value.MakeNil(), nil
}
