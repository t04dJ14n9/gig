// frame.go defines Frame (call stack entry) and DeferInfo (deferred call metadata).
package vm

import (
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
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

	// intLocals is an integer-specialized local variable array.
	// Used by OpInt* superinstructions for 8-byte operations (vs 32 bytes for Value).
	// Allocated alongside locals when the function uses int-specialized opcodes.
	intLocals []int64

	// freeVars are free variables for closures.
	// These are pointers to allow shared state with the enclosing scope.
	freeVars []*value.Value

	// defers is the list of deferred function calls.
	defers []DeferInfo

	// addrTaken is set to true when OpAddr creates a pointer into this frame's locals.
	// Frames with addrTaken must NOT be returned to the pool because closures
	// may hold live references to the locals slice.
	addrTaken bool
}

// DeferInfo represents a deferred function call.
// Deferred calls are executed in LIFO order when the function returns.
type DeferInfo struct {
	// fn is the compiled function to call.
	fn *bytecode.CompiledFunction

	// args are the arguments to pass.
	args []value.Value

	// closure is the closure to call (for OpDeferIndirect).
	closure *Closure

	// externalFunc is an external function or method value to call via reflection.
	// This is used for defer statements that capture external type methods, e.g.:
	//   encoder := base64.NewEncoder(...)
	//   defer encoder.Close()  // externalFunc will hold the Close method value
	externalFunc reflect.Value
	
	// externalInfo holds external function/method metadata (for OpDeferExternal).
	// This is used for interface method invocations in defer.
	externalInfo interface{}
}

// newFrame creates a new call frame for a function with a zero base pointer.
// Used for goroutine and defer frames that start with a fresh operand stack.
// It initializes the local variable array and copies arguments into the first slots.
func newFrame(fn *bytecode.CompiledFunction, args []value.Value, freeVars []*value.Value) *Frame {
	locals := make([]value.Value, fn.NumLocals)

	// Copy arguments to local slots
	for i, arg := range args {
		if i < fn.NumLocals {
			locals[i] = arg
		}
	}

	f := &Frame{
		fn:       fn,
		ip:       0,
		basePtr:  0,
		locals:   locals,
		freeVars: freeVars,
		defers:   nil,
	}

	// Allocate intLocals and mirror int parameters for OpInt* opcodes
	if fn.HasIntLocals {
		f.intLocals = make([]int64, fn.NumLocals)
		for i, arg := range args {
			if i < fn.NumLocals {
				f.intLocals[i] = arg.RawInt()
			}
		}
	}

	return f
}

// framePool is a VM-local pool for reusing Frame objects and their locals slices.
// This eliminates heap allocations in the hot call path (e.g., recursive calls).
type framePool struct {
	frames []*Frame
}

// get returns a recycled Frame reset for the given function, or allocates a new one.
func (p *framePool) get(fn *bytecode.CompiledFunction, basePtr int, freeVars []*value.Value) *Frame {
	var f *Frame
	n := len(p.frames)
	if n > 0 {
		f = p.frames[n-1]
		p.frames = p.frames[:n-1]
		// Reuse the locals slice if it has enough capacity
		if cap(f.locals) >= fn.NumLocals {
			f.locals = f.locals[:fn.NumLocals]
			// Zero out the locals (important for correctness)
			for i := range f.locals {
				f.locals[i] = value.Value{}
			}
		} else {
			f.locals = make([]value.Value, fn.NumLocals)
		}
		// Reuse or allocate intLocals
		if fn.HasIntLocals {
			if cap(f.intLocals) >= fn.NumLocals {
				f.intLocals = f.intLocals[:fn.NumLocals]
				for i := range f.intLocals {
					f.intLocals[i] = 0
				}
			} else {
				f.intLocals = make([]int64, fn.NumLocals)
			}
		} else {
			f.intLocals = nil
		}
	} else {
		f = &Frame{
			locals: make([]value.Value, fn.NumLocals),
		}
		if fn.HasIntLocals {
			f.intLocals = make([]int64, fn.NumLocals)
		}
	}
	f.fn = fn
	f.ip = 0
	f.basePtr = basePtr
	f.freeVars = freeVars
	f.addrTaken = false
	if f.defers != nil {
		f.defers = f.defers[:0]
	}
	return f
}

// put returns a Frame to the pool for reuse.
// Frames with addrTaken are not pooled because closures may hold
// live references to the locals slice.
func (p *framePool) put(f *Frame) {
	if f.addrTaken {
		return
	}
	// Clear references to allow GC of captured values
	f.fn = nil
	f.freeVars = nil
	if f.defers != nil {
		f.defers = f.defers[:0]
	}
	p.frames = append(p.frames, f)
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
