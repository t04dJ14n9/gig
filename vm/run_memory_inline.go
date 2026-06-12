package vm

import (
	"fmt"
	"reflect"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) runSetDeref(sp int) int {
	stack := v.stack
	sp--
	val := stack[sp]
	sp--
	ptr := stack[sp]
	if ptr.IsNil() || !ptr.IsValid() {
		v.panicking = true
		v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
		return sp
	}
	if p, ok := ptr.IntPtr(); ok {
		*p = val.RawInt()
		return sp
	}
	if slot, ok := valueSlotFromValue(ptr); ok {
		*slot = val
		return sp
	}
	if iface := ptr.Interface(); iface != nil {
		if ref, ok := iface.(*GlobalRef); ok {
			ref.Store(val)
			return sp
		}
	}
	if setReflectPointerElemFast(ptr, val) {
		return sp
	}
	v.setElemForSetDeref(ptr, val)
	return sp
}

func setReflectPointerElemFast(ptr value.Value, val value.Value) bool {
	rv, ok := ptr.ReflectValue()
	if !ok || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}
	elem := rv.Elem()
	if !elem.CanSet() {
		return false
	}
	return setReflectElemFast(elem, val)
}

func setReflectElemFast(elem reflect.Value, val value.Value) bool {
	switch elem.Kind() {
	case reflect.Bool:
		if val.Kind() != value.KindBool {
			return false
		}
		elem.SetBool(val.RawBool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Kind() != value.KindInt {
			return false
		}
		elem.SetInt(val.RawInt())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if val.Kind() != value.KindUint {
			return false
		}
		elem.SetUint(uint64(val.RawInt()))
	case reflect.Float32, reflect.Float64:
		if val.Kind() != value.KindFloat {
			return false
		}
		elem.SetFloat(val.Float())
	case reflect.String:
		if val.Kind() != value.KindString {
			return false
		}
		elem.SetString(val.String())
	case reflect.Struct:
		if elem.Type() != reflect.TypeOf(value.Value{}) {
			return false
		}
		elem.Set(reflect.ValueOf(val))
	default:
		return false
	}
	return true
}

func (v *vm) setElemForSetDeref(ptr value.Value, val value.Value) {
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
		}
	}()
	ptr.SetElem(val)
}

func (v *vm) runIndexAddr(frame *Frame, sp int) (int, []value.Value, error) {
	stack := v.stack
	sp--
	index := stack[sp]
	sp--
	container := stack[sp]

	if s, ok := container.IntSlice(); ok {
		idx := index.RawInt()
		if idx < 0 || idx >= int64(len(s)) {
			v.panicking = true
			v.panicVal = value.FromInterface(fmt.Sprintf("runtime error: index out of range [%d] with length %d", idx, len(s)))
			return sp, stack, nil
		}
		stack[sp] = value.MakeIntPtr(&s[idx])
		return sp + 1, stack, nil
	}
	return v.runSlowStackOp(frame, bytecode.OpIndexAddr, sp, container, index)
}

func (v *vm) runDeref(frame *Frame, sp int) (int, []value.Value, error) {
	stack := v.stack
	sp--
	ptr := stack[sp]
	if ptr.Kind() == value.KindReflect || ptr.Kind() == value.KindPointer {
		if ptr.IsNil() {
			v.panicking = true
			v.panicVal = value.FromInterface("runtime error: invalid memory address or nil pointer dereference")
			return sp, stack, nil
		}
	}
	if p, ok := ptr.IntPtr(); ok {
		stack[sp] = value.MakeInt(*p)
		return sp + 1, stack, nil
	}
	if v.runDerefFallback(stack, sp, ptr) {
		return sp, stack, nil
	}
	return sp + 1, stack, nil
}

func (v *vm) runLen(frame *Frame, sp int) (int, []value.Value, error) {
	stack := v.stack
	sp--
	obj := stack[sp]
	switch obj.Kind() {
	case value.KindSlice:
		stack[sp] = value.MakeInt(int64(obj.Len()))
		return sp + 1, stack, nil
	case value.KindString:
		stack[sp] = value.MakeInt(int64(len(obj.String())))
		return sp + 1, stack, nil
	default:
		return v.runSlowStackOp(frame, bytecode.OpLen, sp, obj)
	}
}

func (v *vm) runInlineStackOpComplete(err error) (bool, error) {
	if err != nil {
		return false, err
	}
	if v.panicking {
		return false, nil
	}
	return v.fp > 0, nil
}

func (v *vm) runSlowStackOp(frame *Frame, op bytecode.OpCode, sp int, operands ...value.Value) (int, []value.Value, error) {
	v.sp = sp
	for _, operand := range operands {
		v.push(operand)
	}
	if err := v.executeOp(op, frame); err != nil {
		return v.sp, v.stack, err
	}
	return v.sp, v.stack, nil
}

func (v *vm) runDerefFallback(stack []value.Value, sp int, ptr value.Value) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			v.panicking = true
			v.panicVal = value.FromInterface(r)
			panicked = true
		}
	}()
	stack[sp] = dereferenceValue(ptr)
	return false
}
