package vm

import "github.com/t04dJ14n9/gig/value"

// push pushes a value onto the operand stack.
// Grows the stack if necessary.
func (v *vm) push(val value.Value) {
	if v.sp >= len(v.stack) {
		// Grow stack
		newStack := make([]value.Value, len(v.stack)*2)
		copy(newStack, v.stack)
		v.stack = newStack
	}
	v.stack[v.sp] = val
	v.sp++
}

// pop pops a value from the operand stack.
// Does not check for underflow - caller must ensure stack is not empty.
func (v *vm) pop() value.Value {
	v.sp--
	return v.stack[v.sp]
}

// peek returns the top of the stack without popping.
func (v *vm) peek() value.Value {
	return v.stack[v.sp-1]
}
