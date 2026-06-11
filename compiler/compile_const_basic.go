package compiler

import (
	"go/constant"
	"go/types"
)

type basicConstConverter func(constant.Value) any

// basicConstValue extracts a Go value from a constant.Value based on the basic type kind.
// Returns nil for unsupported kinds.
func basicConstValue(kind types.BasicKind, val constant.Value) any {
	if val == nil {
		return basicZeroValue(kind)
	}

	conv := basicConstConverters[kind]
	if conv == nil {
		return nil
	}
	return conv(val)
}

func complex64ConstValue(val constant.Value) any {
	re := constant.Real(val)
	im := constant.Imag(val)
	reVal, _ := constant.Float64Val(re)
	imVal, _ := constant.Float64Val(im)
	return complex(float32(reVal), float32(imVal))
}

func complex128ConstValue(val constant.Value) any {
	re := constant.Real(val)
	im := constant.Imag(val)
	reVal, _ := constant.Float64Val(re)
	imVal, _ := constant.Float64Val(im)
	return complex(reVal, imVal)
}

// basicConstConverters makes the supported-kind matrix explicit without a large switch.
// The Uintptr entry intentionally maps to uint64 to preserve the old constant-pool shape.
var basicConstConverters = map[types.BasicKind]basicConstConverter{
	types.Bool: boolConstValue, types.UntypedBool: boolConstValue,
	types.Int: exactSignedConstValue[int], types.UntypedInt: exactSignedConstValue[int], types.UntypedRune: exactSignedConstValue[int],
	types.Int8: signedConstValue[int8], types.Int16: signedConstValue[int16], types.Int32: signedConstValue[int32], types.Int64: exactSignedConstValue[int64],
	types.Uint: unsignedConstValue[uint], types.Uint8: unsignedConstValue[uint8], types.Uint16: unsignedConstValue[uint16],
	types.Uint32: unsignedConstValue[uint32], types.Uint64: unsignedConstValue[uint64], types.Uintptr: unsignedConstValue[uint64],
	types.Float32: floatConstValue[float32], types.Float64: floatConstValue[float64], types.UntypedFloat: floatConstValue[float64],
	types.String: stringConstValue, types.UntypedString: stringConstValue,
	types.Complex64: complex64ConstValue, types.Complex128: complex128ConstValue, types.UntypedComplex: complex128ConstValue,
}

func boolConstValue(val constant.Value) any {
	return val.Kind() == constant.Bool && constant.BoolVal(val)
}

func exactSignedConstValue[T ~int | ~int64](val constant.Value) any {
	i, exact := constant.Int64Val(val)
	if exact {
		return T(i)
	}
	var zero T
	return zero
}

func signedConstValue[T ~int8 | ~int16 | ~int32](val constant.Value) any {
	i, _ := constant.Int64Val(val)
	return T(i)
}

func unsignedConstValue[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](val constant.Value) any {
	u, _ := constant.Uint64Val(val)
	return T(u)
}

func floatConstValue[T ~float32 | ~float64](val constant.Value) any {
	f, _ := constant.Float64Val(val)
	return T(f)
}

func stringConstValue(val constant.Value) any {
	return constant.StringVal(val)
}

var basicZeroValues = map[types.BasicKind]any{
	types.Bool: false, types.UntypedBool: false,
	types.Int: int(0), types.UntypedInt: int(0), types.UntypedRune: int(0),
	types.Int8: int8(0), types.Int16: int16(0), types.Int32: int32(0), types.Int64: int64(0),
	types.Uint: uint(0), types.Uint8: uint8(0), types.Uint16: uint16(0),
	types.Uint32: uint32(0), types.Uint64: uint64(0), types.Uintptr: uint64(0),
	types.Float32: float32(0), types.Float64: 0.0, types.UntypedFloat: 0.0,
	types.String: "", types.UntypedString: "",
	types.Complex64: complex64(0), types.Complex128: complex128(0), types.UntypedComplex: complex128(0),
}

func basicZeroValue(kind types.BasicKind) any {
	return basicZeroValues[kind] // nil for unsupported kinds
}
