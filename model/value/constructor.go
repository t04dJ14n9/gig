package value

import "math"

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

// MakeComplex creates a complex128 value.
func MakeComplex(real, imag float64) Value {
	return Value{
		kind: KindComplex,
		size: Size64, // complex128 (2x float64)
		obj:  complex(real, imag),
	}
}

// MakeComplexSized creates a complex value with the given size (Size32=complex64, Size64=complex128).
func MakeComplexSized(real, imag float64, sz Size) Value {
	return Value{
		kind: KindComplex,
		size: sz,
		obj:  complex(real, imag),
	}
}

// MakeComplex64 creates a complex64 value.
func MakeComplex64(real, imag float32) Value {
	return Value{
		kind: KindComplex,
		size: Size32, // complex64 (2x float32)
		obj:  complex(float64(real), float64(imag)),
	}
}

// MakeIntPtr creates a Value wrapping a *int64 pointer (KindPointer).
// Used by OpIndexAddr on native int slices to avoid reflect overhead.
func MakeIntPtr(p *int64) Value {
	return Value{kind: KindPointer, obj: p}
}

// MakeExternal stores a host Go object directly for external-call fast paths.
func MakeExternal(v any) Value {
	if v == nil {
		return MakeNil()
	}
	return Value{kind: KindExternal, obj: v}
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
