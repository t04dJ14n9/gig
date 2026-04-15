// ops_control.go implements control flow, channels, select, defer, panic/recover, print, and halt.
package vm

import (
	"fmt"
	"reflect"

	"git.woa.com/youngjin/gig/model/bytecode"
	"git.woa.com/youngjin/gig/model/external"
	"git.woa.com/youngjin/gig/model/value"
)

// executeControl handles channels, select, defer, panic/recover, print, and halt opcodes.
// Note: OpJump, OpJumpTrue, OpJumpFalse, OpReturn, OpReturnVal are inlined in run.go's
// hot path and never reach this handler.
func (v *vm) executeControl(op bytecode.OpCode, frame *Frame) error { //nolint:gocyclo,cyclop,funlen,maintidx
	switch op {
	case bytecode.OpSend:
		val := v.pop()
		ch := v.pop()
		// Use a Go-level recover to catch "send on closed channel" panic
		// and convert it to a guest-level panic (recoverable by defer/recover).
		var sendErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					sendErr = fmt.Errorf("%v", r)
				}
			}()
			sendErr = ch.SendContext(v.ctx, val)
		}()
		if sendErr != nil {
			// Trigger guest-level panic so defer/recover can handle it
			v.panicking = true
			v.panicVal = value.FromInterface(sendErr.Error())
			break
		}

	case bytecode.OpRecv:
		ch := v.pop()
		val, _, err := ch.RecvContext(v.ctx)
		if err != nil {
			return err
		}
		v.push(val)

	case bytecode.OpRecvOk:
		// Receive with comma-ok: returns (value, ok) tuple
		ch := v.pop()
		val, recvOK, err := ch.RecvContext(v.ctx)
		if err != nil {
			return err
		}
		// Push as tuple (value, ok)
		v.pushCommaOk(val, recvOK)

	case bytecode.OpClose:
		ch := v.pop()
		ch.Close()

	case bytecode.OpSelect:
		// OpSelect performs a select statement using reflect.Select.
		// Operands: [meta_idx:2]
		// Stack (bottom to top): for each state, Chan; if send, also SendVal.
		// Result pushed: tuple (index, recvOk, recv_0, ..., recv_{n-1})
		metaIdx := frame.readUint16()
		meta, ok := v.program.Constants[metaIdx].(bytecode.SelectMeta)
		if !ok {
			return fmt.Errorf("OpSelect: invalid meta at index %d", metaIdx)
		}

		// Pop channels and send values from stack (they were pushed in order,
		// so we need to pop in reverse).
		type stateData struct {
			ch      value.Value
			sendVal value.Value
			isSend  bool
		}
		states := make([]stateData, meta.NumStates)
		// Pop in reverse order
		for i := meta.NumStates - 1; i >= 0; i-- {
			if meta.Dirs[i] { // send
				states[i].sendVal = v.pop()
				states[i].ch = v.pop()
				states[i].isSend = true
			} else { // recv
				states[i].ch = v.pop()
			}
		}

		// Build reflect.SelectCase slice
		// Add 1 for default case (non-blocking) or context cancellation case (blocking)
		numCases := meta.NumStates + 1
		cases := make([]reflect.SelectCase, numCases)
		for i := 0; i < meta.NumStates; i++ {
			rv, _ := states[i].ch.ReflectValue()
			if states[i].isSend {
				sendRV := states[i].sendVal.ToReflectValue(rv.Type().Elem())
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: rv,
					Send: sendRV,
				}
			} else {
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: rv,
				}
			}
		}
		if !meta.Blocking {
			cases[meta.NumStates] = reflect.SelectCase{Dir: reflect.SelectDefault}
		} else {
			// Inject context cancellation case for blocking select
			cases[meta.NumStates] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(v.ctx.Done()),
			}
		}

		// Perform the select
		chosen, recv, recvOK := reflect.Select(cases)

		// Check if context was cancelled (chosen == meta.NumStates in blocking mode)
		if meta.Blocking && chosen == meta.NumStates {
			return v.ctx.Err()
		}

		// Adjust chosen index: if default was selected, chosen == meta.NumStates → map to -1
		if !meta.Blocking && chosen == meta.NumStates {
			chosen = -1
		}

		// Build result tuple: (index, recvOk, recv_0, ..., recv_{n-1})
		tupleLen := 2 + meta.NumRecv
		tuple := make([]value.Value, tupleLen)
		tuple[0] = value.MakeInt(int64(chosen))
		tuple[1] = value.MakeBool(recvOK)

		// Fill recv values: for each recv state (in order), if it was the chosen one, set the value
		recvIdx := 0
		for i := 0; i < meta.NumStates; i++ {
			if !meta.Dirs[i] { // recv state
				if i == chosen {
					tuple[2+recvIdx] = value.MakeFromReflect(recv)
				} else {
					tuple[2+recvIdx] = value.MakeNil()
				}
				recvIdx++
			}
		}

		v.push(value.FromInterface(tuple))

	// Defer/recover
	case bytecode.OpDefer:
		funcIdx := frame.readUint16()
		fn := v.program.FuncByIndex[funcIdx]
		numArgs := fn.NumParams

		// Pop arguments from stack
		args := make([]value.Value, numArgs)
		for i := numArgs - 1; i >= 0; i-- {
			args[i] = v.pop()
		}

		// Add to defer list (will be executed in LIFO order on return)
		frame.defers = append(frame.defers, DeferInfo{
			fn:   fn,
			args: args,
		})

	case bytecode.OpDeferExternal:
		funcIdx := frame.readUint16()
		numArgs := int(frame.readByte())

		// Pop arguments from stack
		args := make([]value.Value, numArgs)
		for i := numArgs - 1; i >= 0; i-- {
			args[i] = v.pop()
		}

		// Get the external function/method info
		externalInfo := v.program.Constants[funcIdx]

		// Store as external defer
		frame.defers = append(frame.defers, DeferInfo{
			args:         args,
			externalInfo: externalInfo,
		})

	case bytecode.OpDeferIndirect:
		numArgs := int(frame.readUint16())

		// Pop arguments from stack (pushed after closure)
		args := make([]value.Value, numArgs)
		for i := numArgs - 1; i >= 0; i-- {
			args[i] = v.pop()
		}

		// Pop closure from stack
		closureVal := v.pop()
		closure, ok := closureVal.RawObj().(*Closure)
		if ok {
			// Gig closure - add to defer list
			frame.defers = append(frame.defers, DeferInfo{
				fn:      closure.Fn,
				args:    args,
				closure: closure,
			})
		} else if rv, ok := closureVal.ReflectValue(); ok {
			if rv.Kind() == reflect.Func {
				// External function/method value - add to defer list
				frame.defers = append(frame.defers, DeferInfo{
					args:         args,
					externalFunc: rv,
				})
			} else if rv.Kind() == reflect.Interface && !rv.IsNil() {
				// Interface wrapping a function
				concrete := rv.Elem()
				if concrete.Kind() == reflect.Func {
					frame.defers = append(frame.defers, DeferInfo{
						args:         args,
						externalFunc: concrete,
					})
				} else {
					// Invalid defer value
					return nil
				}
			} else {
				// Invalid defer value
				return nil
			}
		} else {
			// Invalid defer value - this shouldn't happen in well-formed programs
			return nil
		}

	case bytecode.OpRunDefers:
		// Execute all pending deferred calls synchronously in LIFO order.
		// This is critical for named return values: the code after RunDefers
		// reads the (potentially modified) return values.
		for len(frame.defers) > 0 {
			// Pop the last defer (LIFO)
			d := frame.defers[len(frame.defers)-1]
			frame.defers = frame.defers[:len(frame.defers)-1]

			// Handle external info (OpDeferExternal for interface methods)
			if d.externalInfo != nil {
				if methodInfo, ok := d.externalInfo.(*external.ExternalMethodInfo); ok {
					if err := v.callExternalMethod(methodInfo, d.args); err != nil {
						// Propagate error
						return err
					}
					// Pop the result (deferred calls should not return values)
					if methodInfo.DirectCall == nil {
						// Reflection call may push a result
						_ = v.pop()
					}
				}
				continue
			}

			// Handle external function defers
			if d.externalFunc.IsValid() {
				// Convert arguments to reflect.Value
				argVals := make([]reflect.Value, len(d.args))
				for i, arg := range d.args {
					// Get the argument type from function signature
					funcType := d.externalFunc.Type()
					if i < funcType.NumIn() {
						argType := funcType.In(i)
						argVals[i] = arg.ToReflectValue(argType)
					} else {
						// Variadic argument
						argVals[i] = reflect.ValueOf(arg.Interface())
					}
				}

				// Call the external function
				d.externalFunc.Call(argVals)
				continue
			}

			// Get free variables from closure if present
			var freeVars []*value.Value
			if d.closure != nil {
				freeVars = d.closure.FreeVars
			}

			// Execute the deferred function using a child VM.
			// A child VM isolates the defer's frame stack from the parent,
			// so nested calls within the defer don't interfere with the parent.
			childVM := v.newDeferVM()
			deferFrame := newFrame(d.fn, d.args, freeVars)
			childVM.frames[0] = deferFrame
			childVM.fp = 1
			_, runErr := childVM.run()

			// If the child VM panicked (returned an error), we need to
			// switch to panic mode and run remaining defers using
			// runDefersDuringPanic so that recover() in earlier defers
			// can catch this panic.
			// childVM.run() clears panicking at the top frame, so we
			// detect panics via the error return.
			if runErr != nil {
				v.panicking = true
				// Use the preserved original panic value from the child VM.
				// The child VM's run() formats the error as "panic: <value>"
				// (losing type info), but also saves the original typed value
				// in lastPanicVal before clearing.
				if childVM.lastPanicVal.IsValid() {
					v.panicVal = childVM.lastPanicVal
				} else if childVM.panicVal.IsValid() && childVM.panicVal.Kind() != value.KindNil {
					v.panicVal = childVM.panicVal
				} else {
					// Fallback: parse from error message (loses type information)
					panicMsg := runErr.Error()
					if len(panicMsg) > 7 && panicMsg[:7] == "panic: " {
						v.panicVal = value.FromInterface(panicMsg[7:])
					} else {
						v.panicVal = value.FromInterface(panicMsg)
					}
				}
				v.runDefersDuringPanic(frame)
				break
			}
		}

	case bytecode.OpRecover:
		// Recover from panic. recover() only works when called from a deferred function
		// during panic unwinding. The panic state is stored on panicStack when defers
		// are being executed, or in v.panicking for direct panic context.
		var panicVal value.Value
		recovered := false
		if v.panicking {
			// Direct panic context (shouldn't normally happen inside deferred functions
			// since we save to panicStack, but handle it for safety)
			panicVal = v.panicVal
			v.panicking = false
			v.panicVal = value.MakeNil()
			recovered = true
		} else if len(v.panicStack) > 0 && v.panicStack[len(v.panicStack)-1].panicking {
			// Inside a deferred function: the panic state was saved on the stack.
			// Consume it — mark as recovered.
			panicVal = v.panicStack[len(v.panicStack)-1].panicVal
			v.panicStack[len(v.panicStack)-1].panicking = false
			v.panicStack[len(v.panicStack)-1].panicVal = value.MakeNil()
			recovered = true
		}
		if recovered {
			// Wrap the panic value as a reflect.Value containing an interface{}
			// so that subsequent type assertions (r.(int), r.(string), etc.) work correctly.
			iface := panicVal.Interface()
			if iface != nil {
				// Create a reflect.Value of type interface{} that wraps the concrete value.
				var i any = iface
				rv := reflect.ValueOf(&i).Elem() // type is interface{}, value is int(42)
				v.push(value.MakeFromReflect(rv))
			} else {
				v.push(value.MakeNil())
			}
		} else {
			v.push(value.MakeNil())
		}

	case bytecode.OpPanic:
		msg := v.pop()
		v.panicking = true
		// Go 1.21+ wraps panic(nil) in a PanicNilError so recover() returns non-nil.
		// Match this behavior by wrapping nil in an error-like value.
		if msg.IsNil() {
			v.panicVal = value.FromInterface("panic called with nil argument")
		} else {
			v.panicVal = msg
		}

	case bytecode.OpPrint:
		n := frame.readByte()
		for i := 0; i < int(n); i++ {
			val := v.pop()
			fmt.Print(val.Interface())
		}

	case bytecode.OpPrintln:
		n := frame.readByte()
		args := make([]any, n)
		for i := int(n) - 1; i >= 0; i-- {
			args[i] = v.pop().Interface()
		}
		fmt.Println(args...)

	case bytecode.OpHalt:
		return fmt.Errorf("halt")

	default:
		return fmt.Errorf("unknown opcode: %v", op)
	}

	return nil
}
