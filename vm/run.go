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

	for vm.fp > 0 {
		// Check context every 1024 instructions (bitwise AND is faster than modulus)
		instructionCount++
		if instructionCount&0x3FF == 0 {
			select {
			case <-vm.ctx.Done():
				return value.MakeNil(), vm.ctx.Err()
			default:
			}
		}

		frame := vm.frames[vm.fp-1]
		ins := frame.fn.Instructions

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
			vm.push(frame.locals[idx])
			continue

		case bytecode.OpSetLocal:
			idx := frame.readUint16()
			frame.locals[idx] = vm.pop()
			continue

		case bytecode.OpConst:
			idx := frame.readUint16()
			if int(idx) < len(vm.program.PrebakedConstants) {
				vm.push(vm.program.PrebakedConstants[idx])
			} else if int(idx) < len(vm.program.Constants) {
				vm.push(value.FromInterface(vm.program.Constants[idx]))
			}
			continue

		case bytecode.OpAdd:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeInt(a.RawInt() + b.RawInt()))
			} else {
				vm.push(a.Add(b))
			}
			continue

		case bytecode.OpSub:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeInt(a.RawInt() - b.RawInt()))
			} else {
				vm.push(a.Sub(b))
			}
			continue

		case bytecode.OpMul:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeInt(a.RawInt() * b.RawInt()))
			} else {
				vm.push(a.Mul(b))
			}
			continue

		case bytecode.OpLess:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() < b.RawInt()))
			} else {
				vm.push(value.MakeBool(a.Cmp(b) < 0))
			}
			continue

		case bytecode.OpLessEq:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() <= b.RawInt()))
			} else {
				vm.push(value.MakeBool(a.Cmp(b) <= 0))
			}
			continue

		case bytecode.OpGreater:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() > b.RawInt()))
			} else {
				vm.push(value.MakeBool(a.Cmp(b) > 0))
			}
			continue

		case bytecode.OpGreaterEq:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() >= b.RawInt()))
			} else {
				vm.push(value.MakeBool(a.Cmp(b) >= 0))
			}
			continue

		case bytecode.OpEqual:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() == b.RawInt()))
			} else {
				vm.push(value.MakeBool(a.Equal(b)))
			}
			continue

		case bytecode.OpNotEqual:
			b := vm.pop()
			a := vm.pop()
			if a.Kind() == value.KindInt && b.Kind() == value.KindInt {
				vm.push(value.MakeBool(a.RawInt() != b.RawInt()))
			} else {
				vm.push(value.MakeBool(!a.Equal(b)))
			}
			continue

		case bytecode.OpJump:
			offset := frame.readUint16()
			frame.ip = int(offset)
			continue

		case bytecode.OpJumpTrue:
			offset := frame.readUint16()
			cond := vm.pop()
			if cond.RawBool() {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpJumpFalse:
			offset := frame.readUint16()
			cond := vm.pop()
			if !cond.RawBool() {
				frame.ip = int(offset)
			}
			continue

		case bytecode.OpNot:
			a := vm.pop()
			vm.push(value.MakeBool(!a.RawBool()))
			continue

		case bytecode.OpNil:
			vm.push(value.MakeNil())
			continue

		case bytecode.OpTrue:
			vm.push(value.MakeBool(true))
			continue

		case bytecode.OpFalse:
			vm.push(value.MakeBool(false))
			continue

		case bytecode.OpPop:
			vm.pop()
			continue

		case bytecode.OpDup:
			vm.push(vm.peek())
			continue

		case bytecode.OpCall:
			funcIdx := frame.readUint16()
			numArgs := frame.readByte()
			vm.callCompiledFunction(int(funcIdx), int(numArgs))
			continue

		case bytecode.OpReturn:
			vm.fpool.put(frame)
			vm.fp--
			if vm.fp > 0 {
				prevFrame := vm.frames[vm.fp-1]
				vm.sp = prevFrame.basePtr
			}
			vm.push(value.MakeNil())
			continue

		case bytecode.OpReturnVal:
			retVal := vm.pop()
			vm.fpool.put(frame)
			vm.fp--
			if vm.fp > 0 {
				prevFrame := vm.frames[vm.fp-1]
				vm.sp = prevFrame.basePtr
			}
			vm.push(retVal)
			continue

		case bytecode.OpSetDeref:
			val := vm.pop()
			ptr := vm.pop()
			ptr.SetElem(val)
			continue

		default:
			// Fall through to executeOp for all other opcodes
		}

		// Non-hot-path: dispatch to the full handler
		if err := vm.executeOp(op, frame); err != nil {
			return value.MakeNil(), err
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
						vm.callFunction(d.fn, d.args, nil)
						_, _ = vm.run() // Run the deferred function
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
	if vm.sp > 0 {
		return vm.pop(), nil
	}
	return value.MakeNil(), nil
}
