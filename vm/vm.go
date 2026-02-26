package vm

import (
	"context"
	"fmt"
	"go/types"
	"reflect"
	"sync"
	"sync/atomic"

	"gig/compiler"
	"gig/value"
)

// VM is the bytecode virtual machine.
type VM struct {
	program   *compiler.Program
	stack     []value.Value
	sp        int      // stack pointer
	frames    []*Frame // call frames
	fp        int      // frame pointer
	globals   []value.Value
	ctx       context.Context
	panicking bool
	panicVal  value.Value
}

// New creates a new VM.
func New(program *compiler.Program) *VM {
	return &VM{
		program: program,
		stack:   make([]value.Value, 1024), // initial stack size
		sp:      0,
		frames:  make([]*Frame, 64), // max call depth
		fp:      0,
		globals: make([]value.Value, len(program.Globals)),
	}
}

// Execute runs the specified function with the given arguments.
func (vm *VM) Execute(funcName string, ctx context.Context, args ...value.Value) (value.Value, error) {
	vm.ctx = ctx

	// Look up the function
	fn, ok := vm.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	// Convert args to []value.Value
	valArgs := make([]value.Value, len(args))
	copy(valArgs, args)

	// Create initial frame
	frame := newFrame(fn, 0, valArgs, nil)
	vm.frames[0] = frame
	vm.fp = 1

	// Run the VM
	result, err := vm.run()
	return result, err
}

// ExecuteWithValues runs the specified function with pre-converted arguments.
func (vm *VM) ExecuteWithValues(funcName string, ctx context.Context, args []value.Value) (value.Value, error) {
	vm.ctx = ctx

	// Look up the function
	fn, ok := vm.program.Functions[funcName]
	if !ok {
		return value.MakeNil(), fmt.Errorf("function %q not found", funcName)
	}

	// Create initial frame
	frame := newFrame(fn, 0, args, nil)
	vm.frames[0] = frame
	vm.fp = 1

	// Run the VM
	return vm.run()
}

// run executes the VM loop.
func (vm *VM) run() (value.Value, error) {
	instructionCount := 0

	for vm.fp > 0 {
		// Check context every 1024 instructions
		instructionCount++
		if instructionCount%1024 == 0 {
			select {
			case <-vm.ctx.Done():
				return value.MakeNil(), vm.ctx.Err()
			default:
			}
		}

		frame := vm.frames[vm.fp-1]

		// Check for end of function
		if frame.ip >= len(frame.Instructions()) {
			// Pop frame
			vm.fp--
			continue
		}

		// Fetch opcode
		op := compiler.OpCode(frame.Instructions()[frame.ip])
		frame.ip++

		// Execute opcode
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
			continue
		}
	}

	// Return top of stack (or nil if empty)
	if vm.sp > 0 {
		return vm.pop(), nil
	}
	return value.MakeNil(), nil
}

// executeOp executes a single opcode.
func (vm *VM) executeOp(op compiler.OpCode, frame *Frame) error {
	switch op {
	// Stack operations
	case compiler.OpNop:
		// No operation

	case compiler.OpPop:
		vm.pop()

	case compiler.OpDup:
		val := vm.peek()
		vm.push(val)

	// Constants and locals
	case compiler.OpConst:
		idx := frame.readUint16()
		if int(idx) < len(vm.program.Constants) {
			vm.push(value.FromInterface(vm.program.Constants[idx]))
		}

	case compiler.OpNil:
		vm.push(value.MakeNil())

	case compiler.OpTrue:
		vm.push(value.MakeBool(true))

	case compiler.OpFalse:
		vm.push(value.MakeBool(false))

	case compiler.OpLocal:
		idx := frame.readUint16()
		if int(idx) < len(frame.locals) {
			vm.push(frame.locals[idx])
		}

	case compiler.OpSetLocal:
		idx := frame.readUint16()
		val := vm.pop()
		if int(idx) < len(frame.locals) {
			frame.locals[idx] = val
		}

	case compiler.OpGlobal:
		idx := frame.readUint16()
		if int(idx) < len(vm.globals) {
			vm.push(vm.globals[idx])
		}

	case compiler.OpSetGlobal:
		idx := frame.readUint16()
		val := vm.pop()
		if int(idx) < len(vm.globals) {
			vm.globals[idx] = val
		}

	case compiler.OpFree:
		idx := frame.readByte()
		if int(idx) < len(frame.freeVars) && frame.freeVars[idx] != nil {
			vm.push(*frame.freeVars[idx])
		}

	case compiler.OpSetFree:
		idx := frame.readByte()
		val := vm.pop()
		if int(idx) < len(frame.freeVars) && frame.freeVars[idx] != nil {
			*frame.freeVars[idx] = val
		}

	// Arithmetic
	case compiler.OpAdd:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Add(b))

	case compiler.OpSub:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Sub(b))

	case compiler.OpMul:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Mul(b))

	case compiler.OpDiv:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Div(b))

	case compiler.OpMod:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Mod(b))

	case compiler.OpNeg:
		a := vm.pop()
		vm.push(a.Neg())

	// Bitwise
	case compiler.OpAnd:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.And(b))

	case compiler.OpOr:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Or(b))

	case compiler.OpXor:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.Xor(b))

	case compiler.OpAndNot:
		b := vm.pop()
		a := vm.pop()
		vm.push(a.AndNot(b))

	case compiler.OpLsh:
		n := uint(vm.pop().Int())
		a := vm.pop()
		vm.push(a.Lsh(n))

	case compiler.OpRsh:
		n := uint(vm.pop().Int())
		a := vm.pop()
		vm.push(a.Rsh(n))

	// Comparison
	case compiler.OpEqual:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(a.Equal(b)))

	case compiler.OpNotEqual:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(!a.Equal(b)))

	case compiler.OpLess:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(a.Cmp(b) < 0))

	case compiler.OpLessEq:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(a.Cmp(b) <= 0))

	case compiler.OpGreater:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(a.Cmp(b) > 0))

	case compiler.OpGreaterEq:
		b := vm.pop()
		a := vm.pop()
		vm.push(value.MakeBool(a.Cmp(b) >= 0))

	// Logical
	case compiler.OpNot:
		a := vm.pop()
		vm.push(value.MakeBool(!a.Bool()))

	// Control flow
	case compiler.OpJump:
		offset := frame.readUint16()
		frame.ip = int(offset)

	case compiler.OpJumpTrue:
		offset := frame.readUint16()
		cond := vm.pop()
		if cond.Bool() {
			frame.ip = int(offset)
		}

	case compiler.OpJumpFalse:
		offset := frame.readUint16()
		cond := vm.pop()
		if !cond.Bool() {
			frame.ip = int(offset)
		}

	case compiler.OpCall:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		vm.callCompiledFunction(int(funcIdx), int(numArgs))

	case compiler.OpReturn:
		// Pop frame
		vm.fp--
		if vm.fp > 0 {
			// Restore stack pointer
			prevFrame := vm.frames[vm.fp-1]
			vm.sp = prevFrame.basePtr
		}

	case compiler.OpReturnVal:
		// Save return value
		retVal := vm.pop()
		// Pop frame
		vm.fp--
		if vm.fp > 0 {
			// Restore stack pointer
			prevFrame := vm.frames[vm.fp-1]
			vm.sp = prevFrame.basePtr
		}
		// Push return value
		vm.push(retVal)

	// Container operations
	case compiler.OpMakeSlice:
		capVal := vm.pop()
		lenVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		// Create slice using the type from the type pool
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				slice := reflect.MakeSlice(rt, int(lenVal.Int()), int(capVal.Int()))
				vm.push(value.MakeFromReflect(slice))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpMakeMap:
		sizeVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				m := reflect.MakeMap(rt)
				_ = sizeVal // Size hint ignored for simplicity
				vm.push(value.MakeFromReflect(m))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpMakeChan:
		sizeVal := vm.pop()
		typeIdxVal := vm.pop()
		typeIdx := uint16(typeIdxVal.Int())
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				ch := reflect.MakeChan(rt, int(sizeVal.Int()))
				vm.push(value.MakeFromReflect(ch))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	// Index operations
	case compiler.OpIndex:
		key := vm.pop()
		container := vm.pop()
		switch container.Kind() {
		case value.KindSlice, value.KindArray:
			idx := int(key.Int())
			vm.push(container.Index(idx))
		case value.KindMap:
			vm.push(container.MapIndex(key))
		case value.KindString:
			idx := int(key.Int())
			vm.push(container.Index(idx))
		case value.KindReflect:
			// Handle reflect.Value containing a slice, array, or map
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					vm.push(value.MakeFromReflect(rv.Index(idx)))
				case reflect.Map:
					k := key.ToReflectValue(rv.Type().Key())
					elem := rv.MapIndex(k)
					if !elem.IsValid() {
						vm.push(value.MakeNil())
					} else {
						vm.push(value.MakeFromReflect(elem))
					}
				default:
					vm.push(value.MakeNil())
				}
			} else {
				vm.push(value.MakeNil())
			}
		default:
			vm.push(value.MakeNil())
		}

	case compiler.OpSetIndex:
		val := vm.pop()
		key := vm.pop()
		container := vm.pop()
		switch container.Kind() {
		case value.KindSlice, value.KindArray:
			idx := int(key.Int())
			container.SetIndex(idx, val)
		case value.KindMap:
			container.SetMapIndex(key, val)
		case value.KindReflect:
			if rv, ok := container.ReflectValue(); ok {
				switch rv.Kind() {
				case reflect.Slice, reflect.Array:
					idx := int(key.Int())
					rv.Index(idx).Set(val.ToReflectValue(rv.Type().Elem()))
				case reflect.Map:
					container.SetMapIndex(key, val)
				}
			}
		}

	case compiler.OpSlice:
		// Slice operation: container[low:high:max]
		maxVal := vm.pop()
		highVal := vm.pop()
		lowVal := vm.pop()
		container := vm.pop()
		
		low := int(lowVal.Int())
		
		if rv, ok := container.ReflectValue(); ok {
			// If it's a pointer to an array, dereference it first
			if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Array {
				rv = rv.Elem()
			}
			
			high := int(highVal.Int())
			if high == 0xFFFF {
				high = rv.Len()
			}
			
			var sliced reflect.Value
			if maxVal.Kind() != value.KindNil && maxVal.Int() != 0xFFFF {
				// 3-index slice: container[low:high:max]
				max := int(maxVal.Int())
				sliced = rv.Slice3(low, high, max)
			} else {
				// 2-index slice: container[low:high]
				sliced = rv.Slice(low, high)
			}
			vm.push(value.MakeFromReflect(sliced))
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpField:
		fieldIdx := frame.readUint16()
		obj := vm.pop()
		vm.push(obj.Field(int(fieldIdx)))

	case compiler.OpSetField:
		fieldIdx := frame.readUint16()
		val := vm.pop()
		obj := vm.pop()
		obj.SetField(int(fieldIdx), val)

	// Pointer operations
	case compiler.OpAddr:
		// Get address of a local variable (for taking pointer)
		localIdx := frame.readUint16()
		if int(localIdx) < len(frame.locals) {
			// Create a pointer to the local
			ptr := &frame.locals[localIdx]
			vm.push(value.FromInterface(ptr))
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpIndexAddr:
		// Get address of slice/array element: &slice[index]
		index := vm.pop()
		container := vm.pop()
		idx := int(index.Int())

		if rv, ok := container.ReflectValue(); ok {
			// Dereference pointer if needed
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			// Get element address using reflect
			elem := rv.Index(idx)
			if elem.CanAddr() {
				elemPtr := elem.Addr()
				vm.push(value.MakeFromReflect(elemPtr))
			} else {
				// Can't address - set directly
				vm.push(value.MakeFromReflect(elem))
			}
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpDeref:
		ptr := vm.pop()
		vm.push(ptr.Elem())

	case compiler.OpSetDeref:
		val := vm.pop()
		ptr := vm.pop()
		ptr.SetElem(val)

	// Type operations
	case compiler.OpAssert:
		typeIdx := frame.readUint16()
		_ = typeIdx
		obj := vm.pop()
		// Type assertion - for now, just pass through
		vm.push(obj)

	case compiler.OpConvert:
		typeIdx := frame.readUint16()
		_ = typeIdx
		val := vm.pop()
		// Type conversion - for now, just pass through
		vm.push(val)

	// Function operations
	case compiler.OpClosure:
		funcIdx := frame.readUint16()
		numFree := frame.readByte()
		// Get free variables (popped in reverse order)
		freeVars := make([]*value.Value, numFree)
		for i := int(numFree) - 1; i >= 0; i-- {
			v := vm.pop()
			freeVars[i] = &v
		}
		// Look up the function by index
		var fn *compiler.CompiledFunction
		for _, f := range vm.program.Functions {
			if vm.program.FuncIndex[f.Source] == int(funcIdx) {
				fn = f
				break
			}
		}
		if fn != nil {
			closure := &Closure{Fn: fn, FreeVars: freeVars}
			vm.push(value.FromInterface(closure))
		} else {
			vm.push(value.MakeNil())
		}

	// Concurrency
	case compiler.OpGo:
		// Start goroutine - simplified implementation
		// Would need to capture the function and arguments

	case compiler.OpSend:
		val := vm.pop()
		ch := vm.pop()
		ch.Send(val)

	case compiler.OpRecv:
		ch := vm.pop()
		val, _ := ch.Recv()
		vm.push(val)

	case compiler.OpClose:
		ch := vm.pop()
		ch.Close()

	// Defer/recover
	case compiler.OpDefer:
		funcIdx := frame.readUint16()
		// Would capture function and add to defer list
		_ = funcIdx

	case compiler.OpRecover:
		// Recover from panic
		if vm.panicking {
			vm.push(vm.panicVal)
			vm.panicking = false
			vm.panicVal = value.MakeNil()
		} else {
			vm.push(value.MakeNil())
		}

	// Range operations
	case compiler.OpRange:
		// Create an iterator for the collection
		collection := vm.pop()
		vm.push(value.FromInterface(&iterator{collection: collection, index: -1}))

	case compiler.OpRangeNext:
		// Advance iterator and push a tuple (ok, key, value)
		iterVal := vm.pop()
		iter, ok := iterVal.Interface().(*iterator)
		if !ok {
			// Return tuple (false, nil, nil)
			tuple := []value.Value{value.MakeBool(false), value.MakeNil(), value.MakeNil()}
			vm.push(value.FromInterface(tuple))
			return nil
		}
		iter.index++
		key, val, ok := iter.next()
		// SSA Next returns (ok, key, value) as a tuple
		tuple := []value.Value{value.MakeBool(ok), key, val}
		vm.push(value.FromInterface(tuple))

	// Builtins
	case compiler.OpLen:
		obj := vm.pop()
		switch obj.Kind() {
		case value.KindString:
			vm.push(value.MakeInt(int64(len(obj.String()))))
		case value.KindSlice, value.KindArray, value.KindMap, value.KindChan:
			vm.push(value.MakeInt(int64(obj.Len())))
		case value.KindReflect:
			if rv, ok := obj.ReflectValue(); ok {
				vm.push(value.MakeInt(int64(rv.Len())))
			} else {
				vm.push(value.MakeInt(0))
			}
		default:
			vm.push(value.MakeInt(0))
		}

	case compiler.OpCap:
		obj := vm.pop()
		switch obj.Kind() {
		case value.KindSlice, value.KindArray, value.KindChan:
			vm.push(value.MakeInt(int64(obj.Cap())))
		case value.KindReflect:
			if rv, ok := obj.ReflectValue(); ok {
				vm.push(value.MakeInt(int64(rv.Cap())))
			} else {
				vm.push(value.MakeInt(0))
			}
		default:
			vm.push(value.MakeInt(0))
		}

	case compiler.OpAppend:
		elem := vm.pop()
		slice := vm.pop()
		// Append element to slice
		if rv, ok := slice.ReflectValue(); ok {
			sliceElemType := rv.Type().Elem()
			// Check if SSA packed variadic args into a slice (e.g., append(s, elems...))
			if elemRV, ok2 := elem.ReflectValue(); ok2 && elemRV.Kind() == reflect.Slice && elemRV.Type().Elem() == sliceElemType {
				// The element is a slice of the same element type — spread it
				newSlice := reflect.AppendSlice(rv, elemRV)
				vm.push(value.MakeFromReflect(newSlice))
			} else {
				newSlice := reflect.Append(rv, elem.ToReflectValue(sliceElemType))
				vm.push(value.MakeFromReflect(newSlice))
			}
		} else {
			vm.push(slice)
		}

	case compiler.OpCopy:
		src := vm.pop()
		dst := vm.pop()
		// Copy slice
		if dstRV, ok := dst.ReflectValue(); ok {
			if srcRV, ok := src.ReflectValue(); ok {
				n := reflect.Copy(dstRV, srcRV)
				vm.push(value.MakeInt(int64(n)))
			}
		} else {
			vm.push(value.MakeInt(0))
		}

	case compiler.OpDelete:
		key := vm.pop()
		m := vm.pop()
		m.SetMapIndex(key, value.MakeNil())

	case compiler.OpPanic:
		msg := vm.pop()
		vm.panicking = true
		vm.panicVal = msg

	case compiler.OpPrint:
		n := frame.readByte()
		for i := 0; i < int(n); i++ {
			val := vm.pop()
			fmt.Print(val.Interface())
		}

	case compiler.OpPrintln:
		n := frame.readByte()
		args := make([]any, n)
		for i := int(n) - 1; i >= 0; i-- {
			args[i] = vm.pop().Interface()
		}
		fmt.Println(args...)

	case compiler.OpNew:
		typeIdx := frame.readUint16()
		// Allocate new pointer to value of the given type
		if int(typeIdx) < len(vm.program.Types) {
			typ := vm.program.Types[typeIdx]
			if rt := typeToReflect(typ); rt != nil {
				// Create a new pointer to zero value of the type
				newPtr := reflect.New(rt)
				vm.push(value.MakeFromReflect(newPtr))
			} else {
				vm.push(value.MakeNil())
			}
		} else {
			vm.push(value.MakeNil())
		}

	case compiler.OpMake:
		_ = frame.readUint16() // typeIdx
		_ = frame.readUint16() // sizeIdx
		// Make operation handled by specific OpMakeSlice/Map/Chan

	// External call
	case compiler.OpCallExternal:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		vm.callExternal(int(funcIdx), int(numArgs))

	// Indirect call (closures, function values)
	case compiler.OpCallIndirect:
		numArgs := frame.readByte()
		// Pop arguments
		args := make([]value.Value, numArgs)
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		// Pop the callee
		callee := vm.pop()
		calleeIface := callee.Interface()
		switch fn := calleeIface.(type) {
		case *Closure:
			// Call closure: create new frame with free vars
			vm.callFunction(fn.Fn, args, fn.FreeVars)
		default:
			// Not a known callable — push nil
			vm.push(value.MakeNil())
		}

	case compiler.OpPack:
		count := frame.readUint16()
		// Pop 'count' values from stack and pack into a slice
		values := make([]value.Value, count)
		for i := int(count) - 1; i >= 0; i-- {
			values[i] = vm.pop()
		}
		vm.push(value.FromInterface(values))

	case compiler.OpUnpack:
		// Pop a slice and push each element onto the stack
		slice := vm.pop()
		if slice.Kind() == value.KindSlice || slice.Kind() == value.KindReflect {
			if rv, ok := slice.ReflectValue(); ok {
				for i := 0; i < rv.Len(); i++ {
					vm.push(value.MakeFromReflect(rv.Index(i)))
				}
			}
		}

	case compiler.OpHalt:
		return fmt.Errorf("halt")

	default:
		return fmt.Errorf("unknown opcode: %v", op)
	}

	return nil
}

// push pushes a value onto the stack.
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

// pop pops a value from the stack.
func (vm *VM) pop() value.Value {
	vm.sp--
	return vm.stack[vm.sp]
}

// peek returns the top of the stack without popping.
func (vm *VM) peek() value.Value {
	return vm.stack[vm.sp-1]
}

// callCompiledFunction calls a compiled function.
func (vm *VM) callCompiledFunction(funcIdx, numArgs int) {
	// Get function
	var fn *compiler.CompiledFunction
	for _, f := range vm.program.Functions {
		if vm.program.FuncIndex[f.Source] == funcIdx {
			fn = f
			break
		}
	}
	if fn == nil {
		vm.push(value.MakeNil())
		return
	}

	// Pop arguments
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	// Create new frame
	frame := newFrame(fn, vm.sp, args, nil)
	vm.frames[vm.fp] = frame
	vm.fp++
}

// callFunction calls a function with the given arguments.
func (vm *VM) callFunction(fn *compiler.CompiledFunction, args []value.Value, freeVars []*value.Value) {
	frame := newFrame(fn, vm.sp, args, freeVars)
	vm.frames[vm.fp] = frame
	vm.fp++
}

// callExternal calls an external function.
func (vm *VM) callExternal(funcIdx, numArgs int) {
	// Get external function from constant pool
	if funcIdx >= len(vm.program.Constants) {
		vm.push(value.MakeNil())
		return
	}

	extFunc := vm.program.Constants[funcIdx]

	// Pop arguments
	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = vm.pop()
	}

	// Call the function using reflect
	if extFunc == nil {
		vm.push(value.MakeNil())
		return
	}

	fnVal := reflect.ValueOf(extFunc)
	if fnVal.Kind() != reflect.Func {
		vm.push(value.MakeNil())
		return
	}
	fnType := fnVal.Type()

	// Build reflect.Value arguments
	var in []reflect.Value

	// For variadic calls where SSA passes the variadic slice as the last arg,
	// we need to unpack it for reflect.Call
	if fnType.IsVariadic() && numArgs == fnType.NumIn() {
		// The last arg might be the variadic slice packed by SSA
		lastArg := args[numArgs-1]
		if rv, ok := lastArg.ReflectValue(); ok && rv.Kind() == reflect.Slice {
			// Unpack: use first N-1 args normally, then spread the slice elements
			sliceLen := rv.Len()
			in = make([]reflect.Value, fnType.NumIn()-1+sliceLen)
			for i := 0; i < numArgs-1; i++ {
				in[i] = args[i].ToReflectValue(fnType.In(i))
			}
			elemType := fnType.In(fnType.NumIn() - 1).Elem()
			for i := 0; i < sliceLen; i++ {
				elem := rv.Index(i)
				// If elem is interface{}, unwrap it
				if elem.Kind() == reflect.Interface && !elem.IsNil() {
					elem = elem.Elem()
				}
				if elem.Type().ConvertibleTo(elemType) {
					in[fnType.NumIn()-1+i] = elem.Convert(elemType)
				} else {
					in[fnType.NumIn()-1+i] = elem
				}
			}
		} else {
			// Last arg is not a slice, treat normally
			in = make([]reflect.Value, numArgs)
			for i, arg := range args {
				if i >= fnType.NumIn()-1 {
					variadicType := fnType.In(fnType.NumIn() - 1).Elem()
					in[i] = arg.ToReflectValue(variadicType)
				} else {
					in[i] = arg.ToReflectValue(fnType.In(i))
				}
			}
		}
	} else {
		in = make([]reflect.Value, numArgs)
		for i, arg := range args {
			if i < fnType.NumIn() {
				in[i] = arg.ToReflectValue(fnType.In(i))
			}
		}
	}

	// Call the function
	out := fnVal.Call(in)

	// Convert result
	if len(out) == 0 {
		vm.push(value.MakeNil())
	} else if len(out) == 1 {
		vm.push(value.MakeFromReflect(out[0]))
	} else {
		// Multiple return values - pack as slice
		results := make([]value.Value, len(out))
		for i, v := range out {
			results[i] = value.MakeFromReflect(v)
		}
		vm.push(value.FromInterface(results))
	}
}

// Closure represents a closure with captured free variables.
type Closure struct {
	Fn       *compiler.CompiledFunction
	FreeVars []*value.Value
}

// iterator is a helper for range iteration.
type iterator struct {
	collection value.Value
	index      int
	mapIter    *reflect.MapIter // for map iteration
}

// next advances the iterator and returns the next key, value, and whether there are more elements.
func (it *iterator) next() (key, val value.Value, ok bool) {
	switch it.collection.Kind() {
	case value.KindSlice, value.KindArray, value.KindString:
		if it.index >= it.collection.Len() {
			return value.MakeNil(), value.MakeNil(), false
		}
		key = value.MakeInt(int64(it.index))
		val = it.collection.Index(it.index)
		return key, val, true
	case value.KindMap:
		if it.mapIter == nil {
			if rv, isValid := it.collection.ReflectValue(); isValid {
				it.mapIter = rv.MapRange()
			} else {
				return value.MakeNil(), value.MakeNil(), false
			}
		}
		if !it.mapIter.Next() {
			return value.MakeNil(), value.MakeNil(), false
		}
		key = value.MakeFromReflect(it.mapIter.Key())
		val = value.MakeFromReflect(it.mapIter.Value())
		return key, val, true
	default:
		// Try to use reflect for other types
		if rv, isValid := it.collection.ReflectValue(); isValid {
			switch rv.Kind() {
			case reflect.Slice, reflect.Array, reflect.String:
				if it.index >= rv.Len() {
					return value.MakeNil(), value.MakeNil(), false
				}
				key = value.MakeInt(int64(it.index))
				val = value.MakeFromReflect(rv.Index(it.index))
				return key, val, true
			case reflect.Map:
				if it.mapIter == nil {
					it.mapIter = rv.MapRange()
				}
				if !it.mapIter.Next() {
					return value.MakeNil(), value.MakeNil(), false
				}
				key = value.MakeFromReflect(it.mapIter.Key())
				val = value.MakeFromReflect(it.mapIter.Value())
				return key, val, true
			}
		}
		return value.MakeNil(), value.MakeNil(), false
	}
}

// Goroutine tracking for concurrent execution
var activeGoroutines int64

// StartGoroutine starts a new goroutine and tracks it.
func StartGoroutine(fn func()) {
	atomic.AddInt64(&activeGoroutines, 1)
	go func() {
		defer atomic.AddInt64(&activeGoroutines, -1)
		fn()
	}()
}

// WaitGoroutines waits for all goroutines to complete.
func WaitGoroutines() {
	for atomic.LoadInt64(&activeGoroutines) > 0 {
		// Busy wait - could use a WaitGroup instead
	}
}

// Global VM registry for concurrent execution
var (
	vmRegistryMutex sync.Mutex
	vmRegistry      = make(map[int64]*VM)
	vmIDCounter     int64
)

// RegisterVM registers a VM for later use.
func RegisterVM(vm *VM) int64 {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	vmIDCounter++
	vmRegistry[vmIDCounter] = vm
	return vmIDCounter
}

// UnregisterVM unregisters a VM.
func UnregisterVM(id int64) {
	vmRegistryMutex.Lock()
	defer vmRegistryMutex.Unlock()
	delete(vmRegistry, id)
}

// typeToReflect converts a go/types.Type to reflect.Type.
// This is a simplified implementation that handles common cases.
func typeToReflect(t types.Type) reflect.Type {
	if t == nil {
		return nil
	}

	switch tt := t.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return reflect.TypeOf(false)
		case types.Int:
			return reflect.TypeOf(int(0))
		case types.Int8:
			return reflect.TypeOf(int8(0))
		case types.Int16:
			return reflect.TypeOf(int16(0))
		case types.Int32:
			return reflect.TypeOf(int32(0))
		case types.Int64:
			return reflect.TypeOf(int64(0))
		case types.Uint:
			return reflect.TypeOf(uint(0))
		case types.Uint8:
			return reflect.TypeOf(uint8(0))
		case types.Uint16:
			return reflect.TypeOf(uint16(0))
		case types.Uint32:
			return reflect.TypeOf(uint32(0))
		case types.Uint64:
			return reflect.TypeOf(uint64(0))
		case types.Uintptr:
			return reflect.TypeOf(uintptr(0))
		case types.Float32:
			return reflect.TypeOf(float32(0))
		case types.Float64:
			return reflect.TypeOf(float64(0))
		case types.Complex64:
			return reflect.TypeOf(complex64(0))
		case types.Complex128:
			return reflect.TypeOf(complex128(0))
		case types.String:
			return reflect.TypeOf("")
		default:
			return nil
		}
	case *types.Slice:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.SliceOf(elem)
		}
		return nil
	case *types.Array:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.ArrayOf(int(tt.Len()), elem)
		}
		return nil
	case *types.Map:
		key := typeToReflect(tt.Key())
		val := typeToReflect(tt.Elem())
		if key != nil && val != nil {
			return reflect.MapOf(key, val)
		}
		return nil
	case *types.Chan:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.ChanOf(reflect.BothDir, elem)
		}
		return nil
	case *types.Pointer:
		elem := typeToReflect(tt.Elem())
		if elem != nil {
			return reflect.PtrTo(elem)
		}
		return nil
	case *types.Interface:
		// Interface type — use the empty interface (any) type
		// For the VM, all interfaces are represented as interface{}
		var emptyIface interface{}
		return reflect.TypeOf(&emptyIface).Elem()
	case *types.Named:
		// For named types, try to get the underlying type
		return typeToReflect(tt.Underlying())
	default:
		return nil
	}
}
