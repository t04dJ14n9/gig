package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeRecover() {
	// recover() only works when called from a deferred function during panic
	// unwinding. Panic state is saved on panicStack while defers execute, or in
	// v.panicking for direct panic context.
	var panicVal value.Value
	recovered := false
	if v.panicking {
		panicVal = v.panicVal
		v.panicking = false
		v.panicVal = value.MakeNil()
		recovered = true
	} else if len(v.panicStack) > 0 && v.panicStack[len(v.panicStack)-1].panicking {
		panicVal = v.panicStack[len(v.panicStack)-1].panicVal
		v.panicStack[len(v.panicStack)-1].panicking = false
		v.panicStack[len(v.panicStack)-1].panicVal = value.MakeNil()
		recovered = true
	}
	if !recovered {
		v.push(value.MakeNil())
		return
	}

	// Wrap the panic value as interface{} so type assertions like r.(int) work.
	iface := panicVal.Interface()
	if iface == nil {
		v.push(value.MakeNil())
		return
	}
	var i any = iface
	rv := reflect.ValueOf(&i).Elem()
	v.push(value.MakeFromReflect(rv))
}

func (v *vm) executePanic() {
	msg := v.pop()
	v.panicking = true
	// Go 1.21+ wraps panic(nil) in a PanicNilError so recover returns non-nil.
	if msg.IsNil() {
		v.panicVal = value.FromInterface("panic called with nil argument")
		return
	}
	v.panicVal = msg
}
