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

		// Inline dispatch — avoids function call overhead per instruction
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
