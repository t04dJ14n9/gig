package vm

import "github.com/t04dJ14n9/gig/value"

// push pushes a value onto the operand stack.
// Grows the stack if necessary.
func (vm *VM) push(val value.Value) {
	if vm.sp >= len(vm.stack) {
		// Grow stack
		newStack := make([]value.Value, len(vm.stack)*2)
		copy(newStack, vm.stack)
		vm.stack = newStack
	}
	vm.stack[vm.sp] = val
	vm.sp++
}

// pop pops a value from the operand stack.
// Does not check for underflow - caller must ensure stack is not empty.
func (vm *VM) pop() value.Value {
	vm.sp--
	return vm.stack[vm.sp]
}

// peek returns the top of the stack without popping.
func (vm *VM) peek() value.Value {
	return vm.stack[vm.sp-1]
}
