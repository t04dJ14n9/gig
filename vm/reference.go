package vm

import (
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

var (
	referenceValueType = reflect.TypeOf(value.Value{})
	valueSlotType      = reflect.TypeOf((*value.Value)(nil))
	globalRefType      = reflect.TypeOf((*GlobalRef)(nil))
)

const nilPointerDereference = "runtime error: invalid memory address or nil pointer dereference"

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
	if rv, ok := v.ReflectValue(); ok {
		// Reflect-backed values are common on slice/field address paths. Return
		// false as soon as the reflect shape cannot be a VM slot; falling through
		// to Value.Interface would materialize the wrapped Go value on every index.
		if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Type() != valueSlotType || !rv.CanInterface() {
			return nil, false
		}
		return rv.Interface().(*value.Value), true
	}
	if !v.IsValid() || v.IsNil() {
		return nil, false
	}
	if slot, ok := v.RawObj().(*value.Value); ok {
		return slot, true
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
	if rv, ok := v.ReflectValue(); ok {
		// Most reflect-backed values are not shared globals. Avoid converting
		// arbitrary reflect.Values to interfaces unless the type is exactly the
		// GlobalRef pointer wrapper used by shared global storage.
		if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Type() != globalRefType || !rv.CanInterface() {
			return nil, false
		}
		return rv.Interface().(*GlobalRef), true
	}
	if !v.IsValid() || v.IsNil() {
		return nil, false
	}
	if ref, ok := v.RawObj().(*GlobalRef); ok {
		return ref, true
	}
	if !v.CanInterface() {
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
	if direct, ok := directIndexAddressValue(container, idx); ok {
		return direct
	}

	if slotVal, ok := unwrapValueSlot(container); ok {
		if direct, ok := directIndexAddressValue(slotVal, idx); ok {
			return direct
		}
		container = slotVal
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

func directIndexAddressValue(container value.Value, idx int) (value.Value, bool) {
	if s, ok := container.IntSlice(); ok {
		return value.MakeIntPtr(&s[idx]), true
	}

	if container.Kind() == value.KindBytes {
		b, ok := container.Bytes()
		if !ok {
			return value.MakeNil(), true
		}
		if idx < 0 || idx >= len(b) {
			return value.MakeNil(), true
		}
		elem := reflect.ValueOf(b).Index(idx)
		return addressableValue(elem), true
	}
	return value.Value{}, false
}

// dereferenceValue implements OpDeref's load semantics.
func dereferenceValue(ptr value.Value) value.Value {
	if ref, ok := globalRefFromValue(ptr); ok {
		return ref.Load()
	}

	switch ptr.Kind() {
	case value.KindPointer:
		return dereferencePointerValue(ptr)
	case value.KindInterface:
		return ptr
	case value.KindReflect:
		return dereferenceReflectValue(ptr)
	default:
		return dereferenceConcreteValue(ptr)
	}
}

func dereferenceReflectValue(ptr value.Value) value.Value {
	rv, ok := ptr.ReflectValue()
	if !ok || rv.Kind() != reflect.Ptr {
		return ptr
	}
	if rv.IsNil() {
		return nilDereferenceValue()
	}
	if slot, ok := valueSlotFromValue(ptr); ok {
		return *slot
	}
	return dereferenceReflectElement(rv.Elem())
}

func dereferenceReflectElement(elem reflect.Value) value.Value {
	if isSettableReflectPointerElement(elem) {
		return value.MakeFromReflect(reflect.ValueOf(elem.Interface()))
	}
	if isSettableEmptyInterfaceElement(elem) {
		return dereferenceReflectInterfaceElement(elem)
	}
	if elem.CanAddr() && !isReflectScalar(elem.Kind()) {
		return value.MakeFromReflect(cloneReflectValue(elem))
	}
	return value.MakeFromReflect(elem)
}

func isReflectScalar(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func isSettableReflectPointerElement(elem reflect.Value) bool {
	return elem.Kind() == reflect.Ptr && elem.CanSet() && elem.CanInterface()
}

func isSettableEmptyInterfaceElement(elem reflect.Value) bool {
	return elem.Kind() == reflect.Interface && elem.CanSet() && elem.Type().NumMethod() == 0
}

func dereferenceReflectInterfaceElement(elem reflect.Value) value.Value {
	if elem.IsNil() {
		return value.MakeNil()
	}
	concrete := elem.Elem()
	return value.MakeFromReflect(reflect.ValueOf(concrete.Interface()))
}

func dereferencePointerValue(ptr value.Value) value.Value {
	if ptr.Elem().IsValid() {
		return ptr.Elem()
	}
	return nilDereferenceValue()
}

func dereferenceConcreteValue(ptr value.Value) value.Value {
	if ptr.IsNil() || !ptr.IsValid() {
		return nilDereferenceValue()
	}
	return ptr
}

func panicNilDereference() {
	panic(nilPointerDereference)
}

func nilDereferenceValue() value.Value {
	panic(nilPointerDereference)
}

// setDereferenceValue implements OpSetDeref's store semantics.
func (v *vm) setDereferenceValue(ptr value.Value, val value.Value) {
	if ptr.IsNil() || !ptr.IsValid() {
		panicNilDereference()
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
		return value.MakeFromReflect(addressableReflectValue(rv))
	}
	return value.MakeFromReflect(rv)
}

func addressableReflectValue(rv reflect.Value) reflect.Value {
	if rv.CanInterface() {
		// Normal slice/array elements can use Addr directly. This is the common
		// indexed-write path and avoids reflect.NewAt's slower type/unsafe setup.
		return rv.Addr()
	}
	// Unexported fields cannot be interfaced through Addr. NewAt deliberately
	// keeps the old ability to mutate such fields through VM pointer operations.
	return reflect.NewAt(rv.Type(), value.UnsafeAddrOf(rv))
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
