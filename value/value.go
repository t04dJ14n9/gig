// Package value implements a tagged-union Value system for high-performance interpretation.
//
// The Value type is the fundamental data unit in the Gig interpreter. It uses a tagged-union
// design that stores primitive types (bool, int, uint, float, string, complex) directly in
// the struct fields, avoiding allocation and reflection overhead for common operations.
//
// # Design Philosophy
//
// The Value type is designed for:
//   - Zero allocation for primitive values
//   - Fast arithmetic and comparison operations without reflection
//   - Seamless interop with Go's reflect package for complex types
//   - Type safety through explicit kind checking
//
// # Memory Layout
//
// The Value struct is 48 bytes on 64-bit systems:
//   - kind: 1 byte (type tag)
//   - num: 8 bytes (stores bool, int, uint bits, float bits, complex real)
//   - num2: 8 bytes (stores complex imaginary)
//   - str: 16 bytes (string pointer + length)
//   - obj: 16 bytes (interface for reflect.Value or composite types)
//
// # Kind Types
//
//   - KindNil: null value
//   - KindBool: boolean (stored in num)
//   - KindInt: signed integers (stored in num)
//   - KindUint: unsigned integers (stored in num as bits)
//   - KindFloat: floating point (stored in num as float64 bits)
//   - KindString: string (stored in str)
//   - KindComplex: complex number (real in num, imag in num2)
//   - KindPointer, KindSlice, KindArray, KindMap, KindChan, KindFunc, KindStruct, KindInterface:
//     stored in obj as reflect.Value or native Go value
//   - KindReflect: fallback for types not directly supported
//
// # Example Usage
//
//	// Create values
//	i := value.MakeInt(42)
//	s := value.MakeString("hello")
//	f := value.MakeFloat(3.14)
//
//	// Arithmetic
//	sum := i.Add(value.MakeInt(8)) // sum.Int() == 50
//
//	// Comparison
//	if i.Cmp(value.MakeInt(40)) > 0 {
//	    fmt.Println("42 > 40")
//	}
//
//	// Convert to/from interface{}
//	v := value.FromInterface(myStruct)
//	obj := v.Interface()
package value

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// Kind represents the type of a Value.
type Kind uint8

const (
	KindInvalid Kind = iota
	KindNil
	KindBool
	KindInt     // int, int8, int16, int32, int64
	KindUint    // uint, uint8, uint16, uint32, uint64
	KindFloat   // float32, float64
	KindString  // string
	KindComplex // complex64, complex128
	KindPointer // *T
	KindSlice   // []T
	KindArray   // [N]T
	KindMap     // map[K]V
	KindChan    // chan T
	KindFunc    // func
	KindStruct  // struct{}
	KindInterface
	KindReflect // fallback to reflect.Value
)

// String returns the name of the kind.
func (k Kind) String() string {
	switch k {
	case KindInvalid:
		return "invalid"
	case KindNil:
		return "nil"
	case KindBool:
		return "bool"
	case KindInt:
		return "int"
	case KindUint:
		return "uint"
	case KindFloat:
		return "float"
	case KindString:
		return "string"
	case KindComplex:
		return "complex"
	case KindPointer:
		return "pointer"
	case KindSlice:
		return "slice"
	case KindArray:
		return "array"
	case KindMap:
		return "map"
	case KindChan:
		return "chan"
	case KindFunc:
		return "func"
	case KindStruct:
		return "struct"
	case KindInterface:
		return "interface"
	case KindReflect:
		return "reflect"
	default:
		return "unknown"
	}
}

// Value is a tagged-union that stores Go values with minimal overhead.
// For primitives, operations are done directly on the fields without reflection.
// For complex types, obj stores reflect.Value or native Go values.
type Value struct {
	kind Kind
	num  int64 // Stores: bool (0/1), int, uint bits, float64 bits, complex real part
	num2 int64 // Stores: complex imag part (for KindComplex only)
	str  string
	obj  any // reflect.Value or native Go composite (fallback)
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

// --- Constructors ---

// MakeNil creates a nil value.
func MakeNil() Value {
	return Value{kind: KindNil}
}

// MakeBool creates a bool value.
func MakeBool(b bool) Value {
	var n int64
	if b {
		n = 1
	}
	return Value{kind: KindBool, num: n}
}

// MakeInt creates an int value.
func MakeInt(i int64) Value {
	return Value{kind: KindInt, num: i}
}

// MakeUint creates a uint value.
func MakeUint(u uint64) Value {
	return Value{kind: KindUint, num: int64(u)}
}

// MakeFloat creates a float value.
func MakeFloat(f float64) Value {
	return Value{kind: KindFloat, num: int64(math.Float64bits(f))}
}

// MakeString creates a string value.
func MakeString(s string) Value {
	return Value{kind: KindString, str: s}
}

// MakeComplex creates a complex value.
func MakeComplex(real, imag float64) Value {
	return Value{
		kind: KindComplex,
		num:  int64(math.Float64bits(real)),
		num2: int64(math.Float64bits(imag)),
	}
}

// MakePointer creates a pointer value.
func MakePointer(ptr unsafe.Pointer, elemType reflect.Type) Value {
	return Value{
		kind: KindPointer,
		obj:  reflect.NewAt(elemType, ptr).Elem(),
	}
}

// MakeIntPtr creates a Value wrapping a *int64 pointer (KindPointer).
// Used by OpIndexAddr on native int slices to avoid reflect overhead.
func MakeIntPtr(p *int64) Value {
	return Value{kind: KindPointer, obj: p}
}

// MakeIntSlice creates a Value backed by a native []int64 (KindSlice).
// This avoids reflect overhead for the common []int case.
func MakeIntSlice(s []int64) Value {
	return Value{kind: KindSlice, obj: s}
}

// IntSlice returns the underlying []int64 if this is a native int slice.
// Returns nil, false if not a native int slice.
func (v Value) IntSlice() ([]int64, bool) {
	if v.kind == KindSlice {
		s, ok := v.obj.([]int64)
		return s, ok
	}
	return nil, false
}

// MakeFromReflect creates a Value from reflect.Value.
func MakeFromReflect(rv reflect.Value) Value {
	if !rv.IsValid() {
		return MakeNil()
	}

	// Check if the underlying value is already a Value - unwrap it
	if rv.Kind() == reflect.Struct {
		if t := rv.Type(); t.Name() == "Value" && t.PkgPath() == "gig/value" {
			// This is a value.Value, extract it directly
			return rv.Interface().(Value)
		}
	}

	kind := rv.Kind()
	switch kind {
	case reflect.Bool:
		return MakeBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return MakeInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return MakeUint(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return MakeFloat(rv.Float())
	case reflect.String:
		return MakeString(rv.String())
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		return MakeComplex(real(c), imag(c))
	default:
		return Value{kind: KindReflect, obj: rv}
	}
}

// FromInterface creates a Value from any Go value.
func FromInterface(v any) Value {
	if v == nil {
		return MakeNil()
	}
	return MakeFromReflect(reflect.ValueOf(v))
}

// GoString returns a Go-syntax representation of the value.
func (v Value) GoString() string {
	return fmt.Sprintf("value.Value{kind:%v, num:%d, str:%q, obj:%v}", v.kind, v.num, v.str, v.obj)
}
