package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runCallComplete(err error, frameChanged bool, reloadOnPanic bool) (bool, error) {
	if err != nil {
		return false, err
	}
	if v.panicking {
		return reloadOnPanic && v.fp > 0, nil
	}
	return frameChanged, nil
}

func (v *vm) runIndirectCallFallback(
	callee value.Value,
	args []value.Value,
	sp int,
) (int, []value.Value, bool, error) {
	stack := v.stack
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
