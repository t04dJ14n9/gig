package vm

import (
	"fmt"

	"gig/bytecode"
	"gig/value"
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
func (vm *VM) run() (value.Value, error) {
	// Hoist hot fields into local variables for better register allocation.
	// The Go compiler can keep these in CPU registers across iterations,
	// avoiding repeated loads from vm.* on each instruction.
	stack := vm.stack
	sp := vm.sp
	prebaked := vm.program.PrebakedConstants

	// Cache current frame state to avoid re-reading from vm.frames[] each iteration.
	// These are only invalidated on call/return/executeOp.
	var frame *Frame
	var ins []byte
	var locals []value.Value
	var intLocals []int64
	intConsts := vm.program.IntConstants

	// loadFrame caches the current frame's hot fields into local variables.
	loadFrame := func() {
		frame = vm.frames[vm.fp-1]
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

	// backJumpCount throttles context checks: only check every 128 backward jumps.
	// This replaces per-instruction counting — cheaper because backward jumps only
	// occur in loops, and the counter increment is a single bitwise AND.
	backJumpCount := 0

	if vm.fp > 0 {
		loadFrame()
	}

	for vm.fp > 0 {
		// Check for end of function
		if frame.ip >= len(ins) {
			// Pop frame and return it to pool
			vm.fp--
			vm.fpool.put(frame)
			if vm.fp > 0 {
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
			} else if int(idx) < len(vm.program.Constants) {
				stack[sp] = value.FromInterface(vm.program.Constants[idx])
			}
			sp++
			continue

		case bytecode.OpAdd:
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() - b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() * b.RawInt())
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
			if int(offset) < frame.ip {
				backJumpCount++
				if backJumpCount&0x7F == 0 {
					vm.sp = sp
					select {
					case <-vm.ctx.Done():
						return value.MakeNil(), vm.ctx.Err()
					default:
					}
				}
			}
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
			vm.sp = sp
			vm.callCompiledFunction(int(funcIdx), int(numArgs))
			sp = vm.sp
			stack = vm.stack
			loadFrame()
			continue

		case bytecode.OpReturn:
			vm.fpool.put(frame)
			vm.fp--
			if vm.fp > 0 {
				loadFrame()
				sp = frame.basePtr
			}
			stack[sp] = value.MakeNil()
			sp++
			continue

		case bytecode.OpReturnVal:
			sp--
			retVal := stack[sp]
			vm.fpool.put(frame)
			vm.fp--
			if vm.fp > 0 {
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
			// Fast path: *int64 pointer (from native int slice OpIndexAddr)
			if p, ok := ptr.IntPtr(); ok {
				*p = val.RawInt()
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
				stack[sp] = value.MakeIntPtr(&s[index.RawInt()])
				sp++
				continue
			}
			// Slow path: go through executeOp
			vm.sp = sp
			vm.push(container)
			vm.push(index)
			if err := vm.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			sp = vm.sp
			stack = vm.stack
			if vm.fp > 0 {
				loadFrame()
			}
			continue

		case bytecode.OpDeref:
			sp--
			ptr := stack[sp]
			// Fast path: *int64 pointer (from native int slice OpIndexAddr)
			if p, ok := ptr.IntPtr(); ok {
				stack[sp] = value.MakeInt(*p)
				sp++
				continue
			}
			// Slow path: go through executeOp
			vm.sp = sp
			vm.push(ptr)
			if err := vm.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			sp = vm.sp
			stack = vm.stack
			if vm.fp > 0 {
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
			vm.sp = sp
			vm.push(obj)
			if err := vm.executeOp(op, frame); err != nil {
				return value.MakeNil(), err
			}
			sp = vm.sp
			stack = vm.stack
			if vm.fp > 0 {
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
				stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() - b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() * b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() + b.RawInt())
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
				stack[sp] = value.MakeInt(a.RawInt() - b.RawInt())
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
				locals[idx] = value.MakeInt(r)
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
				locals[idx] = value.MakeInt(r)
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
				locals[idxC] = value.MakeInt(a.RawInt() + b.RawInt())
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
				locals[idxC] = value.MakeInt(a.RawInt() + b.RawInt())
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
				locals[idxC] = value.MakeInt(a.RawInt() - b.RawInt())
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
				locals[idxC] = value.MakeInt(a.RawInt() - b.RawInt())
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
				locals[idxC] = value.MakeInt(a.RawInt() * b.RawInt())
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
				locals[idxC] = value.MakeInt(a.RawInt() * b.RawInt())
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
				r := s[intLocals[jIdx]]
				intLocals[vIdx] = r
				locals[vIdx] = value.MakeInt(r)
			} else {
				// Fallback: execute as IndexAddr + Deref manually
				vm.sp = sp
				vm.push(locals[sIdx])
				vm.push(value.MakeInt(intLocals[jIdx]))
				if err := vm.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				if err := vm.executeOp(bytecode.OpDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				v := vm.pop()
				intLocals[vIdx] = v.RawInt()
				locals[vIdx] = v
				sp = vm.sp
				stack = vm.stack
			}
			continue

		case bytecode.OpIntSliceSet:
			sIdx := readU16()
			jIdx := readU16()
			valIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				s[intLocals[jIdx]] = intLocals[valIdx]
			} else {
				// Fallback: execute as IndexAddr + SetDeref manually
				vm.sp = sp
				vm.push(locals[sIdx])
				vm.push(value.MakeInt(intLocals[jIdx]))
				if err := vm.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				vm.push(value.MakeInt(intLocals[valIdx]))
				if err := vm.executeOp(bytecode.OpSetDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				sp = vm.sp
				stack = vm.stack
			}
			continue

		case bytecode.OpIntSliceSetConst:
			sIdx := readU16()
			jIdx := readU16()
			cIdx := readU16()
			if s, ok := locals[sIdx].IntSlice(); ok {
				s[intLocals[jIdx]] = intConsts[cIdx]
			} else {
				// Fallback: execute as IndexAddr + SetDeref manually
				vm.sp = sp
				vm.push(locals[sIdx])
				vm.push(value.MakeInt(intLocals[jIdx]))
				if err := vm.executeOp(bytecode.OpIndexAddr, frame); err != nil {
					return value.MakeNil(), err
				}
				vm.push(prebaked[cIdx])
				if err := vm.executeOp(bytecode.OpSetDeref, frame); err != nil {
					return value.MakeNil(), err
				}
				sp = vm.sp
				stack = vm.stack
			}
			continue

		case bytecode.OpCallExternal:
			funcIdx := readU16()
			numArgs := int(frame.readByte())
			vm.sp = sp
			vm.callExternal(int(funcIdx), numArgs)
			sp = vm.sp
			stack = vm.stack
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
				vm.sp = sp
				vm.callFunction(closure.Fn, args, closure.FreeVars)
				sp = vm.sp
				stack = vm.stack
				loadFrame()
			} else {
				stack[sp] = value.MakeNil()
				sp++
			}
			continue

		default:
			// Fall through to executeOp for all other opcodes
		}

		// Non-hot-path: dispatch to the full handler.
		// Sync sp back before calling executeOp (it uses vm.push/vm.pop).
		vm.sp = sp
		if err := vm.executeOp(op, frame); err != nil {
			return value.MakeNil(), err
		}
		sp = vm.sp
		stack = vm.stack
		// Reload frame state in case executeOp changed it (call/return within executeOp)
		if vm.fp > 0 {
			loadFrame()
		}

		// Handle panic
		if vm.panicking {
			// Run deferred functions
			if len(frame.defers) > 0 {
				// Execute deferred functions in reverse order
				for i := len(frame.defers) - 1; i >= 0; i-- {
					d := frame.defers[i]
					if d.external != nil {
						// External defer - not supported for now
					} else if d.fn != nil {
						// Internal defer
						vm.sp = sp
						vm.callFunction(d.fn, d.args, nil)
						_, _ = vm.run() // Run the deferred function
						sp = vm.sp
					}
				}
				frame.defers = nil
			}

			// If this is the last frame, return the panic
			if vm.fp == 1 {
				err := fmt.Errorf("panic: %v", vm.panicVal.Interface())
				vm.panicking = false
				vm.panicVal = value.MakeNil()
				return value.MakeNil(), err
			}

			// Propagate panic to caller
			vm.fp--
			vm.fpool.put(frame)
			if vm.fp > 0 {
				loadFrame()
			}
			continue
		}
	}

	// Return top of stack (or nil if empty)
	vm.sp = sp
	if sp > 0 {
		sp--
		vm.sp = sp
		return stack[sp], nil
	}
	return value.MakeNil(), nil
}
