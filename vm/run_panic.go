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
		recovered = v.runDeferDuringPanic(frame.defers[i], recovered)
	}
	frame.defers = nil
	return recovered
}

func (v *vm) runDeferDuringPanic(d DeferInfo, recovered bool) bool {
	if v.runExternalInfoDefer(d) {
		return recovered
	}
	if v.runExternalFuncDefer(d) {
		return recovered
	}
	if d.fn == nil {
		return recovered
	}

	freeVars := deferFreeVars(d)
	if v.panicking {
		return v.runInterpretedDeferWithActivePanic(d, freeVars, recovered)
	}
	return v.runInterpretedDeferAfterRecovery(d, freeVars, recovered)
}

func (v *vm) runExternalInfoDefer(d DeferInfo) bool {
	// OpDeferExternal stores resolved metadata for calls such as
	// `defer mu.Unlock()`. They must run even while unwinding a panic.
	if d.externalInfo == nil {
		return false
	}
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
	return true
}

func (v *vm) runExternalFuncDefer(d DeferInfo) bool {
	// OpDeferIndirect can capture an external function value. These calls use
	// reflect because their concrete function value is known only at runtime.
	if !d.externalFunc.IsValid() {
		return false
	}
	d.externalFunc.Call(reflectDeferArgs(d.externalFunc, d.args))
	return true
}

func reflectDeferArgs(fn reflect.Value, args []value.Value) []reflect.Value {
	argVals := make([]reflect.Value, len(args))
	funcType := fn.Type()
	for j, arg := range args {
		if j < funcType.NumIn() {
			argVals[j] = arg.ToReflectValue(funcType.In(j))
		} else {
			argVals[j] = reflect.ValueOf(arg.Interface())
		}
	}
	return argVals
}

func deferFreeVars(d DeferInfo) []*value.Value {
	if d.closure == nil {
		return nil
	}
	return d.closure.FreeVars
}

func (v *vm) runInterpretedDeferWithActivePanic(d DeferInfo, freeVars []*value.Value, recovered bool) bool {
	// OpRecover must observe the same VM panic stack as the panicking frame, so
	// interpreted defers run recursively on the current VM while the saved panic
	// state is temporarily moved onto panicStack.
	v.pushActivePanicForDefer()
	v.invokeInterpretedDefer(d, freeVars)

	if v.panicking {
		v.dropSavedPanicForReplacedPanic()
		return recovered
	}

	saved := v.popSavedPanic()
	if saved.panicking {
		v.panicking = true
		v.panicVal = saved.panicVal
		return recovered
	}
	return true
}

func (v *vm) pushActivePanicForDefer() {
	v.panicStack = append(v.panicStack, panicState{
		panicking: true,
		panicVal:  v.panicVal,
	})
	v.panicking = false
	v.panicVal = value.MakeNil()
}

func (v *vm) invokeInterpretedDefer(d DeferInfo, freeVars []*value.Value) {
	v.callFunction(d.fn, d.args, freeVars)
	v.deferDepth++
	_, _ = v.run()
	v.deferDepth--
}

func (v *vm) dropSavedPanicForReplacedPanic() {
	v.panicStack = v.panicStack[:len(v.panicStack)-1]
}

func (v *vm) popSavedPanic() panicState {
	saved := v.panicStack[len(v.panicStack)-1]
	v.panicStack = v.panicStack[:len(v.panicStack)-1]
	return saved
}

func (v *vm) runInterpretedDeferAfterRecovery(d DeferInfo, freeVars []*value.Value, recovered bool) bool {
	// After recover has resolved the active panic, remaining interpreted defers
	// run on a lightweight child VM so their stack frames do not disturb the
	// parent frame that is finishing panic unwinding.
	childVM := v.runRecoveredDeferInChildVM(d, freeVars)
	if childVM.lastPanicVal.IsValid() {
		v.panicking = true
		v.panicVal = childVM.lastPanicVal
		return false
	}
	if childVM.panicking {
		v.panicking = true
		v.panicVal = childVM.panicVal
		return false
	}
	return recovered
}

func (v *vm) runRecoveredDeferInChildVM(d DeferInfo, freeVars []*value.Value) *vm {
	childVM := v.newDeferVM()
	deferFrame := newFrame(d.fn, d.args, freeVars)
	childVM.frames[0] = deferFrame
	childVM.fp = 1
	_, _ = childVM.run()
	return childVM
}
