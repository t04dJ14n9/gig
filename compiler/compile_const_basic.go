package compiler

import (
	"go/constant"
	"go/types"
)

// basicConstValue extracts a Go value from a constant.Value based on the basic type kind.
// Returns nil for unsupported kinds.
func basicConstValue(kind types.BasicKind, val constant.Value) any { //nolint:gocyclo,cyclop
	if val == nil {
		return basicZeroValue(kind)
	}

	switch kind { //nolint:exhaustive
	case types.Bool, types.UntypedBool:
		return val.Kind() == constant.Bool && constant.BoolVal(val)
	case types.Int, types.UntypedInt, types.UntypedRune:
		i, exact := constant.Int64Val(val)
		if exact {
			return int(i)
		}
		return int(0)
	case types.Int8:
		i, _ := constant.Int64Val(val)
		return int8(i)
	case types.Int16:
		i, _ := constant.Int64Val(val)
		return int16(i)
	case types.Int32:
		i, _ := constant.Int64Val(val)
		return int32(i)
	case types.Int64:
		i, exact := constant.Int64Val(val)
		if exact {
			return i
		}
		return int64(0)
	case types.Uint:
		u, _ := constant.Uint64Val(val)
		return uint(u)
	case types.Uint8:
		u, _ := constant.Uint64Val(val)
		return uint8(u)
	case types.Uint16:
		u, _ := constant.Uint64Val(val)
		return uint16(u)
	case types.Uint32:
		u, _ := constant.Uint64Val(val)
		return uint32(u)
	case types.Uint64:
		u, _ := constant.Uint64Val(val)
		return u
	case types.Uintptr:
		u, _ := constant.Uint64Val(val)
		return uint64(u)
	case types.Float32:
		f, _ := constant.Float64Val(val)
		return float32(f)
	case types.Float64, types.UntypedFloat:
		f, _ := constant.Float64Val(val)
		return f
	case types.String, types.UntypedString:
		return constant.StringVal(val)
	case types.Complex64:
		return complex64ConstValue(val)
	case types.Complex128, types.UntypedComplex:
		return complex128ConstValue(val)
	default:
		return nil
	}
}

func complex64ConstValue(val constant.Value) complex64 {
	re := constant.Real(val)
	im := constant.Imag(val)
	reVal, _ := constant.Float64Val(re)
	imVal, _ := constant.Float64Val(im)
	return complex(float32(reVal), float32(imVal))
}

func complex128ConstValue(val constant.Value) complex128 {
	re := constant.Real(val)
	im := constant.Imag(val)
	reVal, _ := constant.Float64Val(re)
	imVal, _ := constant.Float64Val(im)
	return complex(reVal, imVal)
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
