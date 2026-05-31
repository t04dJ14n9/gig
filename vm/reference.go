package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

var (
	referenceValueType = reflect.TypeOf(value.Value{})
	valueSlotType      = reflect.TypeOf((*value.Value)(nil))
)

// unwrapValueSlot resolves a *value.Value reference produced by OpGlobal,
// OpAddr, or closure capture into the value currently stored in that slot.
func unwrapValueSlot(v value.Value) (value.Value, bool) {
	slot, ok := valueSlotFromValue(v)
	if !ok {
		return value.MakeNil(), false
	}
	return *slot, true
}

// valueSlotFromValue detects a *value.Value without calling Interface on
// non-interfaceable reflect.Values such as unexported embedded fields.
func valueSlotFromValue(v value.Value) (*value.Value, bool) {
	if !v.IsValid() || v.IsNil() {
		return nil, false
	}
	if rv, ok := v.ReflectValue(); ok && rv.IsValid() && rv.Kind() == reflect.Ptr && !rv.IsNil() && rv.Type() == valueSlotType && rv.CanInterface() {
		slot, ok := rv.Interface().(*value.Value)
		return slot, ok
	}
	if !v.CanInterface() {
		return nil, false
	}
	slot, ok := v.Interface().(*value.Value)
	return slot, ok
}

// globalRefFromValue detects a shared-global reference. Shared stateful
// execution must go through GlobalRef so load/store operations remain locked.
func globalRefFromValue(v value.Value) (*GlobalRef, bool) {
	if !v.IsValid() || v.IsNil() || !v.CanInterface() {
		return nil, false
	}
	ref, ok := v.Interface().(*GlobalRef)
	return ref, ok
}

// fieldAddressValue implements OpFieldAddr's reference semantics.
func fieldAddressValue(structPtr value.Value, fieldIdx int) value.Value {
	if slotVal, ok := unwrapValueSlot(structPtr); ok {
		structPtr = slotVal
	}

	rv, ok := structPtr.ReflectValue()
	if !ok {
		return value.MakeNil()
	}

	s, ok := reflectAddressRoot(rv)
	if !ok {
		return value.MakeNil()
	}
	s = unwrapInterfaceStructValue(s)
	if !s.IsValid() || s.Kind() != reflect.Struct {
		return value.MakeNil()
	}
	s = unwrapValueStruct(s)
	if !s.IsValid() || s.Kind() != reflect.Struct {
		return value.MakeNil()
	}
	return addressableValue(s.Field(fieldIdx))
}

// indexAddressValue implements OpIndexAddr's reference semantics.
func indexAddressValue(container value.Value, idx int) value.Value {
	if slotVal, ok := unwrapValueSlot(container); ok {
		container = slotVal
	}

	if s, ok := container.IntSlice(); ok {
		return value.MakeIntPtr(&s[idx])
	}

	if container.Kind() == value.KindBytes {
		b, ok := container.Bytes()
		if !ok {
			return value.MakeNil()
		}
		if idx < 0 || idx >= len(b) {
			return value.MakeNil()
		}
		elem := reflect.ValueOf(b).Index(idx)
		return addressableValue(elem)
	}

	rv, ok := container.ReflectValue()
	if !ok {
		return value.MakeNil()
	}
	if rv.Kind() == reflect.Ptr {
		rv, ok = reflectAddressRoot(rv)
		if !ok {
			return value.MakeNil()
		}
	}
	return addressableValue(rv.Index(idx))
}

// dereferenceValue implements OpDeref's load semantics.
func dereferenceValue(ptr value.Value) value.Value {
	if ref, ok := globalRefFromValue(ptr); ok {
		return ref.Load()
	}

	switch ptr.Kind() {
	case value.KindPointer:
		if ptr.Elem().IsValid() {
			return ptr.Elem()
		}
		panic("runtime error: invalid memory address or nil pointer dereference")
	case value.KindInterface:
		return ptr
	case value.KindReflect:
		rv, ok := ptr.ReflectValue()
		if !ok || rv.Kind() != reflect.Ptr {
			return ptr
		}
		if rv.IsNil() {
			panic("runtime error: invalid memory address or nil pointer dereference")
		}
		if slot, ok := valueSlotFromValue(ptr); ok {
			return *slot
		}
		elem := rv.Elem()
		if elem.Kind() == reflect.Ptr && elem.CanSet() && elem.CanInterface() {
			return value.MakeFromReflect(reflect.ValueOf(elem.Interface()))
		}
		if elem.Kind() == reflect.Interface && elem.CanSet() && elem.Type().NumMethod() == 0 {
			if elem.IsNil() {
				return value.MakeNil()
			}
			concrete := elem.Elem()
			return value.MakeFromReflect(reflect.ValueOf(concrete.Interface()))
		}
		if elem.CanAddr() {
			return value.MakeFromReflect(cloneReflectValue(elem))
		}
		return value.MakeFromReflect(elem)
	default:
		if ptr.IsNil() || !ptr.IsValid() {
			panic("runtime error: invalid memory address or nil pointer dereference")
		}
		return ptr
	}
}

// setDereferenceValue implements OpSetDeref's store semantics.
func (v *vm) setDereferenceValue(ptr value.Value, val value.Value) {
	if ptr.IsNil() || !ptr.IsValid() {
		panic("runtime error: invalid memory address or nil pointer dereference")
	}
	if ref, ok := globalRefFromValue(ptr); ok {
		ref.Store(val)
		return
	}
	if rv, ok := ptr.ReflectValue(); ok && rv.Kind() == reflect.Ptr && !rv.IsNil() {
		elem := rv.Elem()
		if elem.IsValid() && elem.CanSet() && elem.Kind() == reflect.Interface {
			elem.Set(v.valueForReflectSet(val, elem.Type()))
			return
		}
	}
	ptr.SetElem(val)
}

// reflectAddressRoot normalizes a reflect.Value before field/index addressing.
// It unwraps VM slot pointers while preserving unexported-field safety.
func reflectAddressRoot(rv reflect.Value) (reflect.Value, bool) {
	if rv.Kind() != reflect.Ptr {
		return rv, true
	}
	if rv.IsNil() {
		return rv.Elem(), true
	}
	if rv.Type() == valueSlotType && rv.CanInterface() {
		slot, ok := rv.Interface().(*value.Value)
		if !ok {
			return reflect.Value{}, false
		}
		inner, ok := slot.ReflectValue()
		if !ok || !inner.IsValid() {
			return reflect.Value{}, false
		}
		if inner.Kind() == reflect.Ptr {
			return inner.Elem(), true
		}
		return inner, true
	}
	return rv.Elem(), true
}

func unwrapInterfaceStructValue(rv reflect.Value) reflect.Value {
	if rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}
	return rv
}

func unwrapValueStruct(rv reflect.Value) reflect.Value {
	if rv.Kind() != reflect.Struct || rv.Type() != referenceValueType || !rv.CanInterface() {
		return rv
	}
	inner, ok := rv.Interface().(value.Value).ReflectValue()
	if !ok || !inner.IsValid() {
		return rv
	}
	if inner.Kind() == reflect.Ptr {
		return inner.Elem()
	}
	return inner
}

func addressableValue(rv reflect.Value) value.Value {
	if rv.CanAddr() {
		return value.MakeFromReflect(reflect.NewAt(rv.Type(), value.UnsafeAddrOf(rv)))
	}
	return value.MakeFromReflect(rv)
}

// cloneReflectValue creates an independent copy of a reflect.Value that
// references addressable memory (slice element, struct field, etc.).
// This breaks the alias so subsequent writes through the original pointer
// don't corrupt the copy.
func cloneReflectValue(rv reflect.Value) reflect.Value {
	copy := reflect.New(rv.Type()).Elem()
	copy.Set(rv)
	return copy
}
