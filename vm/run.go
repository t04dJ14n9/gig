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
	instructionCount := 0

	// Hoist hot fields into local variables for better register allocation.
	// The Go compiler can keep these in CPU registers across iterations,
	// avoiding repeated loads from vm.* on each instruction.
	stack := vm.stack
	sp := vm.sp
	prebaked := vm.program.PrebakedConstants

	for vm.fp > 0 {
		// Check context every 1024 instructions (bitwise AND is faster than modulus)
		instructionCount++
		if instructionCount&0x3FF == 0 {
			// Sync back before potential return
			vm.sp = sp
			select {
			case <-vm.ctx.Done():
				return value.MakeNil(), vm.ctx.Err()
			default:
			}
		}

		frame := vm.frames[vm.fp-1]
		ins := frame.fn.Instructions
		locals := frame.locals

		// Check for end of function
		if frame.ip >= len(ins) {
			// Pop frame and return it to pool
			vm.fp--
			vm.fpool.put(frame)
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
			idx := frame.readUint16()
			stack[sp] = locals[idx]
			sp++
			continue

		case bytecode.OpSetLocal:
			idx := frame.readUint16()
			sp--
			locals[idx] = stack[sp]
			continue

		case bytecode.OpConst:
			idx := frame.readUint16()
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
			offset := frame.readUint16()
			frame.ip = int(offset)
			continue

		case bytecode.OpJumpTrue:
			offset := frame.readUint16()
			sp--
			if stack[sp].RawBool() {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpJumpFalse:
			offset := frame.readUint16()
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
			funcIdx := frame.readUint16()
			numArgs := frame.readByte()
			vm.sp = sp
			vm.callCompiledFunction(int(funcIdx), int(numArgs))
			sp = vm.sp
			stack = vm.stack
			continue

		case bytecode.OpReturn:
			vm.fpool.put(frame)
			vm.fp--
			if vm.fp > 0 {
				prevFrame := vm.frames[vm.fp-1]
				sp = prevFrame.basePtr
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
				prevFrame := vm.frames[vm.fp-1]
				sp = prevFrame.basePtr
			}
			stack[sp] = retVal
			sp++
			continue

		case bytecode.OpSetDeref:
			sp--
			val := stack[sp]
			sp--
			ptr := stack[sp]
			ptr.SetElem(val)
			continue

		// ========================================
		// Superinstructions: fused ops for hot loops
		// ========================================

		case bytecode.OpAddLocalLocal:
			idxA := frame.readUint16()
			idxB := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			offset := frame.readUint16()
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

		case bytecode.OpAddSetLocal:
			idx := frame.readUint16()
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idx] = value.MakeInt(a.RawInt() + b.RawInt())
			} else {
				locals[idx] = a.Add(b)
			}
			continue

		case bytecode.OpSubSetLocal:
			idx := frame.readUint16()
			sp--
			b := stack[sp]
			sp--
			a := stack[sp]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idx] = value.MakeInt(a.RawInt() - b.RawInt())
			} else {
				locals[idx] = a.Sub(b)
			}
			continue

		case bytecode.OpLocalLocalAddSetLocal:
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			idxC := frame.readUint16()
			a := locals[idxA]
			b := locals[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeInt(a.RawInt() + b.RawInt())
			} else {
				locals[idxC] = a.Add(b)
			}
			continue

		case bytecode.OpLocalConstAddSetLocal:
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			idxC := frame.readUint16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeInt(a.RawInt() + b.RawInt())
			} else {
				locals[idxC] = a.Add(b)
			}
			continue

		case bytecode.OpLocalConstSubSetLocal:
			idxA := frame.readUint16()
			idxB := frame.readUint16()
			idxC := frame.readUint16()
			a := locals[idxA]
			b := prebaked[idxB]
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				locals[idxC] = value.MakeInt(a.RawInt() - b.RawInt())
			} else {
				locals[idxC] = a.Sub(b)
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
