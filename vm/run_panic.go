package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

// derefAllocLocal dereferences an Alloc pointer stored in a frame's local slot.
// In SSA, named return variables and captured locals are represented as pointer
// allocations (OpNew). This function reads the value behind the pointer, mirroring
// what OpDeref does in the normal return path.
func derefAllocLocal(ptr value.Value) value.Value {
	switch ptr.Kind() {
	case value.KindPointer:
		return ptr.Elem()
	case value.KindInterface:
		if rv, ok := ptr.ReflectValue(); ok {
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				return value.MakeFromReflect(rv.Elem())
			}
		}
	}
	// If it's a *value.Value wrapper (from OpNew for function types)
	if rv, ok := ptr.ReflectValue(); ok {
		if rv.Kind() == reflect.Ptr && !rv.IsNil() {
			return value.MakeFromReflect(rv.Elem())
		}
	}
	return ptr
}

// runDefersDuringPanic runs deferred functions in LIFO order during a panic.
// It uses the shared VM + recursive run() so that OpRecover (inside the deferred
// function) can clear v.panicking on the same VM instance.
// Returns true if the panic was recovered by a deferred function.
//
// Caller must sync v.sp before calling; v.sp is updated on return.
func (v *vm) runDefersDuringPanic(frame *Frame) bool {
	recovered := false
	if len(frame.defers) == 0 {
		return recovered
	}

	for i := len(frame.defers) - 1; i >= 0; i-- {
		d := frame.defers[i]

		// Handle external method defers (OpDeferExternal, e.g. defer mu.Unlock())
		// These must run during panic recovery, just like in Go.
		if d.externalInfo != nil {
			switch info := d.externalInfo.(type) {
			case *external.ExternalMethodInfo:
				_ = v.callExternalMethod(info, d.args)
				if info.DirectCall == nil {
					_ = v.pop()
				}
			case *external.ExternalFuncInfo:
				before := v.sp
				_ = v.callResolvedExternal(bytecode.ResolveConstant(info), d.args)
				if v.sp > before {
					_ = v.pop()
				}
			}
			continue
		}

		// Handle external function value defers (OpDeferIndirect with external func)
		// e.g. defer fn() where fn is an external function variable
		if d.externalFunc.IsValid() {
			argVals := make([]reflect.Value, len(d.args))
			for j, arg := range d.args {
				funcType := d.externalFunc.Type()
				if j < funcType.NumIn() {
					argType := funcType.In(j)
					argVals[j] = arg.ToReflectValue(argType)
				} else {
					argVals[j] = reflect.ValueOf(arg.Interface())
				}
			}
			d.externalFunc.Call(argVals)
			continue
		}

		// Skip nil function entries (shouldn't happen, but defensive)
		if d.fn == nil {
			continue
		}

		// Get free variables from closure if present
		var freeVars []*value.Value
		if d.closure != nil {
			freeVars = d.closure.FreeVars
		}

		if v.panicking {
			// Panic is active: push state onto panicStack so that
			// OpRecover (inside the deferred function) can access it.
			// Clear v.panicking so the recursive run() doesn't immediately
			// re-enter the panic handler on the defer's frame.
			v.panicStack = append(v.panicStack, panicState{
				panicking: true,
				panicVal:  v.panicVal,
			})
			v.panicking = false
			v.panicVal = value.MakeNil()

			v.callFunction(d.fn, d.args, freeVars)
			v.deferDepth++
			_, _ = v.run()
			v.deferDepth--

			if v.panicking {
				// The defer itself panicked (and was not recovered).
				// Pop the saved state — this new panic replaces the old one.
				v.panicStack = v.panicStack[:len(v.panicStack)-1]
				continue
			}

			// Pop the saved state and check if recover() consumed it.
			saved := v.panicStack[len(v.panicStack)-1]
			v.panicStack = v.panicStack[:len(v.panicStack)-1]

			if saved.panicking {
				// The defer didn't call recover() — restore the panic.
				v.panicking = true
				v.panicVal = saved.panicVal
			} else {
				// recover() was called — panic is resolved.
				// Continue running remaining defers in normal mode.
				recovered = true
			}
		} else {
			// Panic already recovered: run remaining defers normally.
			// Use child VM to avoid interfering with the parent frame stack.
			childVM := v.newDeferVM()
			deferFrame := newFrame(d.fn, d.args, freeVars)
			childVM.frames[0] = deferFrame
			childVM.fp = 1
			_, _ = childVM.run()

			// If the child VM panicked, re-enter panic mode.
			if childVM.lastPanicVal.IsValid() {
				v.panicking = true
				v.panicVal = childVM.lastPanicVal
				recovered = false
			} else if childVM.panicking {
				v.panicking = true
				v.panicVal = childVM.panicVal
				recovered = false
			}
		}
	}
	frame.defers = nil
	return recovered
}
