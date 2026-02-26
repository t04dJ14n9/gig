// Package vm provides the bytecode virtual machine.
package vm

import (
	"gig/compiler"
	"gig/value"
)

// Frame represents a call frame.
type Frame struct {
	fn       *compiler.CompiledFunction // compiled function
	ip       int                        // instruction pointer
	basePtr  int                        // operand stack base pointer
	locals   []value.Value              // local variables
	freeVars []*value.Value             // free variables (for closures)
	defers   []DeferInfo                // deferred calls
}

// DeferInfo represents a deferred function call.
type DeferInfo struct {
	fn       *compiler.CompiledFunction
	args     []value.Value
	external any // external function if not nil
}

// newFrame creates a new call frame.
func newFrame(fn *compiler.CompiledFunction, basePtr int, args []value.Value, freeVars []*value.Value) *Frame {
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

// Instructions returns the function's bytecode.
func (f *Frame) Instructions() []byte {
	return f.fn.Instructions
}

// readUint16 reads a 2-byte operand at the current IP.
func (f *Frame) readUint16() uint16 {
	val := uint16(f.fn.Instructions[f.ip])<<8 | uint16(f.fn.Instructions[f.ip+1])
	f.ip += 2
	return val
}

// readByte reads a 1-byte operand at the current IP.
func (f *Frame) readByte() byte {
	val := f.fn.Instructions[f.ip]
	f.ip++
	return val
}
