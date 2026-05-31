package value

import (
	"fmt"
	"reflect"
)

// Size records the original Go bit-width for numeric kinds.
// It occupies 1 byte in the padding gap between kind and num (zero extra memory).
type Size uint8

const (
	Size0   Size = 0  // default / unspecified (treated as widest: int64, uint64, float64)
	Size8   Size = 8  // int8, uint8
	Size16  Size = 16 // int16, uint16
	Size32  Size = 32 // int32, uint32, float32
	Size64  Size = 64 // int64, uint64, float64
	SizePtr Size = 1  // int, uint (platform-dependent, but always 64-bit for Gig)
)

// Value is a tagged-union that stores Go values with minimal overhead.
// For primitives (int, float, bool, uint, nil), operations are done directly
// on the num field without touching obj. For complex types (string, complex,
// reflect.Value, slices, maps, etc.), obj stores the value.
//
// Layout: 32 bytes (kind:1 + size:1 + pad:6 + num:8 + obj:16)
type Value struct {
	kind Kind
	size Size  // original Go bit-width (lives in padding, zero extra memory)
	num  int64 // Stores: bool (0/1), int, uint bits, float64 bits
	obj  any   // string, complex128, reflect.Value, native Go composites, or nil for primitives
}

// InterpretedInterfaceValue preserves the dynamic type of a script-defined
// named value stored in an interface{}.
type InterpretedInterfaceValue struct {
	Value     Value
	TypeName  string
	IsPointer bool
}

// MakeInterpretedInterface creates an interface value carrying interpreter
// dynamic type metadata.
func MakeInterpretedInterface(val Value, typeName string, isPointer bool) Value {
	return Value{
		kind: KindInterface,
		obj: &InterpretedInterfaceValue{
			Value:     val,
			TypeName:  typeName,
			IsPointer: isPointer,
		},
	}
}

// InterpretedInterface returns interpreter dynamic type metadata, if present.
func (v Value) InterpretedInterface() (*InterpretedInterfaceValue, bool) {
	if v.kind != KindInterface {
		return nil, false
	}
	dyn, ok := v.obj.(*InterpretedInterfaceValue)
	return dyn, ok
}

// Kind returns the kind of the value.
func (v Value) Kind() Kind { return v.kind }

// RawInt returns the raw int64 value without kind checking.
// Only valid when Kind() == KindInt. Designed to be inlined by the Go compiler.
func (v Value) RawInt() int64 { return v.num }

// RawBool returns the raw bool value without kind checking.
// Only valid when Kind() == KindBool. Designed to be inlined by the Go compiler.
func (v Value) RawBool() bool { return v.num != 0 }

// IsNil returns true if the value is nil.
func (v Value) IsNil() bool {
	if v.kind == KindNil {
		return true
	}
	if v.kind == KindReflect {
		if rv, ok := v.obj.(reflect.Value); ok {
			if !rv.IsValid() {
				return true
			}
			// Only call IsNil on types that support it
			switch rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
				return rv.IsNil()
			}
			return false
		}
	}
	return false
}

// IsValid returns true if the value is valid.
func (v Value) IsValid() bool {
	return v.kind != KindInvalid
}

// RawObj returns the raw obj field for direct type assertions in the hot path.
// This avoids the overhead of Interface() which goes through a full kind-switch.
func (v Value) RawObj() any { return v.obj }

// MakeFunc creates a Value storing a function/closure pointer directly in obj.
// This avoids the reflect.ValueOf overhead of FromInterface for callable objects.
func MakeFunc(fn any) Value {
	return Value{kind: KindFunc, obj: fn}
}

// MakeValueSlice creates a Value backed by a native []Value slice.
// Used by DirectCall wrappers for multi-return packing — zero reflection.
func MakeValueSlice(vals []Value) Value {
	return Value{kind: KindSlice, obj: vals}
}

// ValueSlice returns the underlying []Value if this is a native value slice.
// Returns nil, false if not a native value slice.
func (v Value) ValueSlice() ([]Value, bool) {
	if v.kind == KindSlice {
		s, ok := v.obj.([]Value)
		return s, ok
	}
	return nil, false
}

// GoString returns a Go-syntax representation of the value.
func (v Value) GoString() string {
	return fmt.Sprintf("value.Value{kind:%v, num:%d, obj:%v}", v.kind, v.num, v.obj)
}

// ClosureExecutor is implemented by closure objects that can be executed.
// This interface breaks the circular dependency between value/ and vm/ —
// vm.Closure implements it, and value.ToReflectValue uses it to convert
// closures into real Go functions via reflect.MakeFunc.
type ClosureExecutor interface {
	Execute(args []reflect.Value, outTypes []reflect.Type) []reflect.Value
}
