package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runExternalCall(funcIdx, numArgs, sp int) (int, []value.Value, bool, error) {
	prevFP := v.fp
	v.sp = sp

	var callErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				v.panicking = true
				v.panicVal = value.FromInterface(r)
			}
		}()
		callErr = v.callExternal(funcIdx, numArgs)
	}()
	if callErr != nil {
		return v.sp, v.stack, false, callErr
	}
	if v.panicking {
		return v.sp, v.stack, false, nil
	}
	if err := v.checkCtx(); err != nil {
		return v.sp, v.stack, false, err
	}
	return v.sp, v.stack, v.fp != prevFP, nil
}

func (v *vm) runCallComplete(err error, frameChanged bool, reloadOnPanic bool) (bool, error) {
	if err != nil {
		return false, err
	}
	if v.panicking {
		return reloadOnPanic && v.fp > 0, nil
	}
	return frameChanged, nil
}

func (v *vm) runIndirectCall(sp, numArgs int) (int, []value.Value, bool, error) {
	stack := v.stack
	var argsBuf [8]value.Value
	var args []value.Value
	if numArgs <= len(argsBuf) {
		args = argsBuf[:numArgs]
	} else {
		args = make([]value.Value, numArgs)
	}

	spLocal := sp
	for i := numArgs - 1; i >= 0; i-- {
		spLocal--
		args[i] = stack[spLocal]
	}
	spLocal--
	callee := stack[spLocal]
	sp = spLocal

	if closure, ok := callee.RawObj().(*Closure); ok {
		v.sp = sp
		v.callFunction(closure.Fn, args, closure.FreeVars)
		return v.sp, v.stack, true, nil
	}
	if rv, ok := callee.ReflectValue(); ok && rv.Kind() == reflect.Func {
		return v.runReflectIndirectCall(rv, args, sp)
	}

	stack[sp] = value.MakeNil()
	return sp + 1, stack, false, nil
}

func (v *vm) runReflectIndirectCall(rv reflect.Value, args []value.Value, sp int) (int, []value.Value, bool, error) {
	stack := v.stack
	if rv.IsNil() {
		v.sp = sp
		v.panicking = true
		v.panicVal = value.FromInterface("invalid memory address or nil pointer dereference")
		return sp, stack, false, nil
	}

	in := make([]reflect.Value, len(args))
	fnType := rv.Type()
	for i := 0; i < len(args); i++ {
		if i < fnType.NumIn() {
			in[i] = args[i].ToReflectValue(fnType.In(i))
		}
	}

	var out []reflect.Value
	func() {
		defer func() {
			if r := recover(); r != nil {
				v.sp = sp
				v.panicking = true
				v.panicVal = value.FromInterface(r)
			}
		}()
		out = rv.Call(in)
	}()
	if v.panicking {
		return sp, v.stack, false, nil
	}
	if len(out) == 0 {
		stack[sp] = value.MakeNil()
	} else {
		stack[sp] = value.MakeFromReflect(out[0])
	}
	return sp + 1, stack, false, nil
}
