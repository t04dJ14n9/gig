// Package vm provides the bytecode virtual machine.
package vm

import (
	"gig/bytecode"
	"gig/value"
)

// Frame represents a call frame on the VM's call stack.
// Each function call creates a new frame with its own local variables
// and instruction pointer.
type Frame struct {
	// fn is the compiled function being executed.
	fn *bytecode.CompiledFunction

	// ip is the instruction pointer (current bytecode offset).
	ip int

	// basePtr is the operand stack base pointer for this frame.
	basePtr int

	// locals is the local variable array.
	// Parameters are stored in the first slots.
	locals []value.Value

	// freeVars are free variables for closures.
	// These are pointers to allow shared state with the enclosing scope.
	freeVars []*value.Value

	// defers is the list of deferred function calls.
	defers []DeferInfo
}

// DeferInfo represents a deferred function call.
// Deferred calls are executed in LIFO order when the function returns.
type DeferInfo struct {
	// fn is the compiled function to call.
	fn *bytecode.CompiledFunction

	// args are the arguments to pass.
	args []value.Value

	// external is the external function to call (if not nil).
	external any
}

// newFrame creates a new call frame for a function.
// It initializes the local variable array and copies arguments into the first slots.
func newFrame(fn *bytecode.CompiledFunction, basePtr int, args []value.Value, freeVars []*value.Value) *Frame {
	locals := make([]value.Value, fn.NumLocals)

	// Copy arguments to local slots
	for i, arg := range args {
		if i < fn.NumLocals {
			locals[i] = arg
		}
	}

	return &Frame{
		fn:       fn,
		ip:       0,
		basePtr:  basePtr,
		locals:   locals,
		freeVars: freeVars,
		defers:   nil,
	}
}

// Instructions returns the function's bytecode instructions.
func (f *Frame) Instructions() []byte {
	return f.fn.Instructions
}

// readUint16 reads a 2-byte operand at the current instruction pointer.
// Advances the instruction pointer by 2 bytes.
func (f *Frame) readUint16() uint16 {
	val := uint16(f.fn.Instructions[f.ip])<<8 | uint16(f.fn.Instructions[f.ip+1])
	f.ip += 2
	return val
}

// readByte reads a 1-byte operand at the current instruction pointer.
// Advances the instruction pointer by 1 byte.
func (f *Frame) readByte() byte {
	val := f.fn.Instructions[f.ip]
	f.ip++
	return val
}
