// ops_call.go handles function/closure calls, goroutine spawning, and tuple pack/unpack.
package vm

import (
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/value"
)

// executeCall handles function call, closure creation, goroutine spawning,
// and pack/unpack opcodes.
func (v *vm) executeCall(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen
	switch op {
	case bytecode.OpCall:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		v.callCompiledFunction(int(funcIdx), int(numArgs))

	case bytecode.OpCallExternal:
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()
		if err := v.callExternal(int(funcIdx), int(numArgs)); err != nil {
			return err
		}

	case bytecode.OpCallIndirect:
		numArgs := frame.readByte()
		// Pop arguments using stack-allocated buffer
		var argsBuf [8]value.Value
		var args []value.Value
		if int(numArgs) <= len(argsBuf) {
			args = argsBuf[:numArgs]
		} else {
			args = make([]value.Value, numArgs)
		}
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = v.pop()
		}
		// Pop the callee
		callee := v.pop()
		switch fn := callee.RawObj().(type) {
		case *Closure:
			// Call closure: create new frame with free vars
			v.callFunction(fn.Fn, args, fn.FreeVars)
		case *bytecode.CompiledFunction:
			// Call compiled function
			v.callFunction(fn, args, nil)
		default:
			// Check if callee is a reflect-based function (e.g., from a typed container)
			if rv, ok := callee.ReflectValue(); ok && rv.Kind() == reflect.Func {
				in := make([]reflect.Value, numArgs)
				fnType := rv.Type()
				for i := 0; i < int(numArgs); i++ {
					if i < fnType.NumIn() {
						in[i] = args[i].ToReflectValue(fnType.In(i))
					}
				}
				out := rv.Call(in)
				if len(out) == 0 {
					v.push(value.MakeNil())
				} else {
					v.push(value.MakeFromReflect(out[0]))
				}
			} else {
				// Not a known callable — push nil
				v.push(value.MakeNil())
			}
		}

	case bytecode.OpClosure:
		funcIdx := frame.readUint16()
		numFree := frame.readByte()
		// Look up the function by index (O(1))
		var fn *bytecode.CompiledFunction
		if int(funcIdx) < len(v.program.FuncByIndex) {
			fn = v.program.FuncByIndex[funcIdx]
		}
		if fn != nil {
			closure := getClosure(fn, int(numFree))
			closure.Program = v.program
			closure.InitialGlobals = v.initialGlobals
			// Propagate runtime context so that closures converted to Go
			// functions (via reflect.MakeFunc for sync.Once.Do etc.) can
			// access shared globals, spawn goroutines, and use the same
			// external call cache as the parent VM.
			closure.Shared = v.shared
			closure.Goroutines = v.goroutines
			closure.ExtCallCache = v.extCallCache
			closure.Ctx = v.ctx
			// Get free variables (popped in reverse order)
			for i := int(numFree) - 1; i >= 0; i-- {
				v := v.pop()
				// Create a new *value.Value slot holding the captured value.
				// This allows the closure to read/write the slot via OpFree/OpSetFree.
				// If the captured value is a reflect pointer (e.g., *int from Alloc),
				// all closures sharing that pointer will see each other's modifications.
				slot := new(value.Value)
				*slot = v
				closure.FreeVars[i] = slot
			}
			v.push(value.MakeFunc(closure))
		} else {
			// Still need to pop free vars to keep stack balanced
			for i := 0; i < int(numFree); i++ {
				v.pop()
			}
			v.push(value.MakeNil())
		}

	case bytecode.OpGoCall:
		// OpGoCall spawns a new goroutine to execute a function call.
		// Operands: [func_idx:2, num_args:1]
		// Stack: [... args] -> [...] (arguments consumed)
		funcIdx := frame.readUint16()
		numArgs := frame.readByte()

		// Pop arguments from current goroutine's stack
		args := make([]value.Value, numArgs)
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = v.pop()
		}

		// Get the function to call (O(1))
		var goFn *bytecode.CompiledFunction
		if int(funcIdx) < len(v.program.FuncByIndex) {
			goFn = v.program.FuncByIndex[funcIdx]
		}

		if goFn != nil {
			// Create a child VM with shared globals
			childVM := v.newChildVM()

			// Capture for closure
			capturedFn := goFn
			capturedArgs := args

			// Track the goroutine
			if err := v.goroutines.Start(func() {
				// Create initial frame for the child goroutine
				childFrame := newFrame(capturedFn, capturedArgs, nil)
				childVM.frames[0] = childFrame
				childVM.fp = 1

				// Run the child VM (ignore return value - goroutine result is discarded)
				_, _ = childVM.run()
			}); err != nil {
				return err
			}
		}

	case bytecode.OpGoCallIndirect:
		// OpGoCallIndirect spawns a new goroutine to execute a closure call.
		// Operands: [num_args:1]
		// Stack: [... closure args...] -> [...] (closure and arguments consumed)
		numArgs := frame.readByte()

		// Pop arguments from current goroutine's stack
		args := make([]value.Value, numArgs)
		for i := int(numArgs) - 1; i >= 0; i-- {
			args[i] = v.pop()
		}

		// Pop the closure
		callee := v.pop()

		if closure, ok := callee.RawObj().(*Closure); ok {
			// Create a child VM with shared globals
			childVM := v.newChildVM()

			// Capture for closure
			capturedClosure := closure
			capturedArgs := args

			// Track the goroutine
			if err := v.goroutines.Start(func() {
				// Create initial frame for the child goroutine with free vars
				childFrame := newFrame(capturedClosure.Fn, capturedArgs, capturedClosure.FreeVars)
				childVM.frames[0] = childFrame
				childVM.fp = 1

				// Run the child VM (ignore return value - goroutine result is discarded)
				_, _ = childVM.run()
			}); err != nil {
				return err
			}
		}

	case bytecode.OpPack:
		count := frame.readUint16()
		// Pop 'count' values from stack and pack into a slice
		values := make([]value.Value, count)
		for i := int(count) - 1; i >= 0; i-- {
			values[i] = v.pop()
		}
		v.push(value.FromInterface(values))

	case bytecode.OpUnpack:
		// Pop a slice and push each element onto the stack
		slice := v.pop()
		// Fast path: native []value.Value (produced by DirectCall multi-return wrappers)
		if vals, ok := slice.ValueSlice(); ok {
			for _, elem := range vals {
				v.push(elem)
			}
			break
		}
		if slice.Kind() == value.KindSlice || slice.Kind() == value.KindReflect {
			if rv, ok := slice.ReflectValue(); ok {
				for i := 0; i < rv.Len(); i++ {
					v.push(value.MakeFromReflect(rv.Index(i)))
				}
			}
		}
	}

	return nil
}
