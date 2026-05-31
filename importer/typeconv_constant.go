package importer

import (
	"go/constant"
	"go/token"
	"reflect"
)

// convertToConstantValue converts a Go value to a constant.Value for use in types.Const.
// Uses reflection to handle all basic types uniformly.
func convertToConstantValue(val any) constant.Value {
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Bool:
		return constant.MakeBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return constant.MakeInt64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return constant.MakeUint64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return constant.MakeFloat64(rv.Float())
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		re := constant.MakeFloat64(real(c))
		im := constant.MakeFloat64(imag(c))
		return constant.BinaryOp(re, token.ADD, constant.BinaryOp(im, token.MUL, constant.MakeImag(constant.MakeInt64(1))))
	case reflect.String:
		return constant.MakeString(rv.String())
	default:
		return constant.MakeUnknown()
	}
}
