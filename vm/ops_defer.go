package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeDefer(frame *Frame) {
	funcIdx := frame.readUint16()
	fn := v.program.FuncByIndex[funcIdx]
	numArgs := fn.NumParams

	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = v.pop()
	}

	frame.defers = append(frame.defers, DeferInfo{
		fn:   fn,
		args: args,
	})
}

func (v *vm) executeDeferExternal(frame *Frame) {
	funcIdx := frame.readUint16()
	numArgs := int(frame.readByte())

	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = v.pop()
	}

	externalInfo := v.program.Constants[funcIdx]
	frame.defers = append(frame.defers, DeferInfo{
		args:         args,
		externalInfo: externalInfo,
	})
}

func (v *vm) executeDeferIndirect(frame *Frame) error {
	numArgs := int(frame.readUint16())

	args := make([]value.Value, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		args[i] = v.pop()
	}

	closureVal := v.pop()
	closure, ok := closureVal.RawObj().(*Closure)
	if ok {
		frame.defers = append(frame.defers, DeferInfo{
			fn:      closure.Fn,
			args:    args,
			closure: closure,
		})
		return nil
	}

	if rv, ok := closureVal.ReflectValue(); ok {
		if rv.Kind() == reflect.Func {
			frame.defers = append(frame.defers, DeferInfo{
				args:         args,
				externalFunc: rv,
			})
			return nil
		}
		if rv.Kind() == reflect.Interface && !rv.IsNil() {
			concrete := rv.Elem()
			if concrete.Kind() == reflect.Func {
				frame.defers = append(frame.defers, DeferInfo{
					args:         args,
					externalFunc: concrete,
				})
			}
			return nil
		}
	}

	// Invalid defer value. Well-formed programs should not reach this path.
	return nil
}

func (v *vm) executeRunDefers(frame *Frame) error {
	// Execute all pending deferred calls synchronously in LIFO order.
	// This is critical for named return values: the code after RunDefers
	// reads the potentially modified return values.
	for len(frame.defers) > 0 {
		d := frame.defers[len(frame.defers)-1]
		frame.defers = frame.defers[:len(frame.defers)-1]

		if d.externalInfo != nil {
			if err := v.executeExternalDefer(d); err != nil {
				return err
			}
			continue
		}

		if d.externalFunc.IsValid() {
			v.executeReflectFuncDefer(d)
			continue
		}

		if stop := v.executeClosureDefer(frame, d); stop {
			break
		}
	}
	return nil
}

func (v *vm) executeExternalDefer(d DeferInfo) error {
	switch info := d.externalInfo.(type) {
	case *external.ExternalMethodInfo:
		if err := v.callExternalMethod(info, d.args); err != nil {
			return err
		}
		if info.DirectCall == nil {
			_ = v.pop()
		}
	case *external.ExternalFuncInfo:
		before := v.sp
		if err := v.callResolvedExternal(bytecode.ResolveConstant(info), d.args); err != nil {
			return err
		}
		if v.sp > before {
			_ = v.pop()
		}
	}
	return nil
}

func (v *vm) executeReflectFuncDefer(d DeferInfo) {
	argVals := make([]reflect.Value, len(d.args))
	funcType := d.externalFunc.Type()
	for i, arg := range d.args {
		if i < funcType.NumIn() {
			argVals[i] = arg.ToReflectValue(funcType.In(i))
		} else {
			argVals[i] = reflect.ValueOf(arg.Interface())
		}
	}
	d.externalFunc.Call(argVals)
}

func (v *vm) executeClosureDefer(frame *Frame, d DeferInfo) bool {
	var freeVars []*value.Value
	if d.closure != nil {
		freeVars = d.closure.FreeVars
	}

	childVM := v.newDeferVM()
	deferFrame := newFrame(d.fn, d.args, freeVars)
	childVM.frames[0] = deferFrame
	childVM.fp = 1
	_, runErr := childVM.run()
	if runErr == nil {
		return false
	}

	v.panicking = true
	if childVM.lastPanicVal.IsValid() {
		v.panicVal = childVM.lastPanicVal
	} else if childVM.panicVal.IsValid() && childVM.panicVal.Kind() != value.KindNil {
		v.panicVal = childVM.panicVal
	} else {
		// Fallback: parse from error message, which loses type information.
		panicMsg := runErr.Error()
		if len(panicMsg) > 7 && panicMsg[:7] == "panic: " {
			v.panicVal = value.FromInterface(panicMsg[7:])
		} else {
			v.panicVal = value.FromInterface(panicMsg)
		}
	}
	v.runDefersDuringPanic(frame)
	return true
}
