// Package value implements a tagged-union Value system for high-performance interpretation.
//
// The Value type is the fundamental data unit in the Gig interpreter. It uses a tagged-union
// design that stores primitive types (bool, int, uint, float) directly in the num field,
// avoiding allocation and reflection overhead for common operations.
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
// The Value struct is 32 bytes on 64-bit systems:
//   - kind: 1 byte (type tag)
//   - size: 1 byte (original Go type bit-width for numeric kinds) + 6 bytes padding
//   - num: 8 bytes (stores bool, int, uint bits, float bits)
//   - obj: 16 bytes (interface for string, complex, reflect.Value, composite types)
//
// The size field records the original Go type (e.g. int8 vs int32 vs int64) so that
// Interface() can return the exact Go type declared in the user's source code.
// It lives in the padding gap between kind and num, so it adds zero extra memory.
//
// Primitives (int, float, bool, uint, nil) are stored entirely in kind+size+num with obj=nil,
// so they never cause GC pressure.
//
// # Kind Types
//
//   - KindNil: null value
//   - KindBool: boolean (stored in num)
//   - KindInt: signed integers (stored in num)
//   - KindUint: unsigned integers (stored in num as bits)
//   - KindFloat: floating point (stored in num as float64 bits)
//   - KindString: string (stored in obj)
//   - KindComplex: complex number (stored in obj as complex128)
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
	KindBytes   // []byte stored natively (zero reflection)
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
	case KindBytes:
		return "bytes"
	default:
		return "unknown"
	}
}

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

// MakeInt creates an int value (platform-dependent, mapped to int64 internally).
func MakeInt(i int64) Value {
	return Value{kind: KindInt, size: SizePtr, num: i}
}

// MakeInt8 creates an int8 value.
func MakeInt8(i int8) Value {
	return Value{kind: KindInt, size: Size8, num: int64(i)}
}

// MakeInt16 creates an int16 value.
func MakeInt16(i int16) Value {
	return Value{kind: KindInt, size: Size16, num: int64(i)}
}

// MakeInt32 creates an int32 value.
func MakeInt32(i int32) Value {
	return Value{kind: KindInt, size: Size32, num: int64(i)}
}

// MakeInt64 creates an int64 value.
func MakeInt64(i int64) Value {
	return Value{kind: KindInt, size: Size64, num: i}
}

// MakeUint creates a uint value (platform-dependent, mapped to uint64 internally).
func MakeUint(u uint64) Value {
	return Value{kind: KindUint, size: SizePtr, num: int64(u)}
}

// MakeUint8 creates a uint8 value.
func MakeUint8(u uint8) Value {
	return Value{kind: KindUint, size: Size8, num: int64(u)}
}

// MakeUint16 creates a uint16 value.
func MakeUint16(u uint16) Value {
	return Value{kind: KindUint, size: Size16, num: int64(u)}
}

// MakeUint32 creates a uint32 value.
func MakeUint32(u uint32) Value {
	return Value{kind: KindUint, size: Size32, num: int64(u)}
}

// MakeUint64 creates a uint64 value.
func MakeUint64(u uint64) Value {
	return Value{kind: KindUint, size: Size64, num: int64(u)}
}

// MakeFloat creates a float64 value.
func MakeFloat(f float64) Value {
	return Value{kind: KindFloat, size: Size64, num: int64(math.Float64bits(f))}
}

// MakeFloat32 creates a float32 value.
func MakeFloat32(f float32) Value {
	return Value{kind: KindFloat, size: Size32, num: int64(math.Float64bits(float64(f)))}
}

// RawSize returns the size tag of a Value. Used by arithmetic ops to propagate
// the original type width through computations.
func (v Value) RawSize() Size { return v.size }

// MakeString creates a string value.
func MakeString(s string) Value {
	return Value{kind: KindString, obj: s}
}

// MakeComplex creates a complex value.
func MakeComplex(real, imag float64) Value {
	return Value{
		kind: KindComplex,
		obj:  complex(real, imag),
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

// MakeBytes creates a Value backed by a native []byte (KindBytes).
// This avoids reflect overhead for []byte arguments and return values.
func MakeBytes(b []byte) Value {
	return Value{kind: KindBytes, obj: b}
}

// Bytes returns the underlying []byte if this is a KindBytes value.
// Returns nil, false if not a KindBytes value.
func (v Value) Bytes() ([]byte, bool) {
	if v.kind == KindBytes {
		b, ok := v.obj.([]byte)
		return b, ok
	}
	return nil, false
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

// IntPtr returns the underlying *int64 if this is a native int pointer (from IndexAddr on []int64).
// Returns nil, false if not a *int64.
func (v Value) IntPtr() (*int64, bool) {
	if v.kind == KindPointer {
		p, ok := v.obj.(*int64)
		return p, ok
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
		if t := rv.Type(); t.Name() == "Value" && t.PkgPath() == "git.woa.com/youngjin/gig/model/value" {
			// This is a value.Value, extract it directly
			return rv.Interface().(Value)
		}
	}

	kind := rv.Kind()
	switch kind {
	case reflect.Bool:
		return MakeBool(rv.Bool())
	case reflect.Int:
		return MakeInt(rv.Int())
	case reflect.Int8:
		return MakeInt8(int8(rv.Int()))
	case reflect.Int16:
		return MakeInt16(int16(rv.Int()))
	case reflect.Int32:
		return MakeInt32(int32(rv.Int()))
	case reflect.Int64:
		return MakeInt64(rv.Int())
	case reflect.Uint:
		return MakeUint(rv.Uint())
	case reflect.Uint8:
		return MakeUint8(uint8(rv.Uint()))
	case reflect.Uint16:
		return MakeUint16(uint16(rv.Uint()))
	case reflect.Uint32:
		return MakeUint32(uint32(rv.Uint()))
	case reflect.Uint64, reflect.Uintptr:
		return MakeUint64(rv.Uint())
	case reflect.Float32:
		return MakeFloat32(float32(rv.Float()))
	case reflect.Float64:
		return MakeFloat(rv.Float())
	case reflect.String:
		return MakeString(rv.String())
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		return MakeComplex(real(c), imag(c))
	case reflect.Slice:
		// Native []byte storage: avoid reflect overhead for the common []byte case
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return MakeBytes(rv.Bytes())
		}
		return Value{kind: KindReflect, obj: rv}
	default:
		return Value{kind: KindReflect, obj: rv}
	}
}

// FromInterface creates a Value from any Go value.
func FromInterface(v any) Value {
	if v == nil {
		return MakeNil()
	}
	// Fast path: detect common types via type switch to avoid reflect.ValueOf.
	// Each case uses the typed constructor to preserve the original Go type.
	switch val := v.(type) {
	case bool:
		return MakeBool(val)
	case int:
		return MakeInt(int64(val))
	case int8:
		return MakeInt8(val)
	case int16:
		return MakeInt16(val)
	case int32:
		return MakeInt32(val)
	case int64:
		return MakeInt64(val)
	case uint:
		return MakeUint(uint64(val))
	case uint8:
		return MakeUint8(val)
	case uint16:
		return MakeUint16(val)
	case uint32:
		return MakeUint32(val)
	case uint64:
		return MakeUint64(val)
	case float32:
		return MakeFloat32(val)
	case float64:
		return MakeFloat(val)
	case string:
		return MakeString(val)
	case []byte:
		return MakeBytes(val)
	}
	return MakeFromReflect(reflect.ValueOf(v))
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
