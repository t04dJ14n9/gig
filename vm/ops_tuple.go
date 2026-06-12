package vm

import "github.com/t04dJ14n9/gig/model/value"

func (v *vm) executePack(frame *Frame) {
	count := frame.readUint16()
	values := make([]value.Value, count)
	for i := int(count) - 1; i >= 0; i-- {
		values[i] = v.pop()
	}
	v.push(value.FromInterface(values))
}

func (v *vm) executeUnpack() {
	slice := v.pop()
	if vals, ok := slice.ValueSlice(); ok {
		v.pushValueSlice(vals)
		return
	}
	v.pushReflectSlice(slice)
}

func (v *vm) pushValueSlice(vals []value.Value) {
	// Direct-call multi-return wrappers already produce []value.Value, so this
	// fast path avoids reflect conversion for the common tuple-unpack case.
	for _, elem := range vals {
		v.push(elem)
	}
}

func (v *vm) pushReflectSlice(slice value.Value) {
	if slice.Kind() != value.KindSlice && slice.Kind() != value.KindReflect {
		return
	}
	rv, ok := slice.ReflectValue()
	if !ok {
		return
	}
	for i := 0; i < rv.Len(); i++ {
		v.push(value.MakeFromReflect(rv.Index(i)))
	}
}
