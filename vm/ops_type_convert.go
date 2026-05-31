package vm

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/t04dJ14n9/gig/model/value"
)

func (v *vm) executeTypeConvert(frame *Frame) {
	typeIdx := frame.readUint16()
	targetType := v.program.Types[typeIdx]
	val := v.pop()

	v.push(v.convertValueToType(val, targetType))
}

func (v *vm) convertValueToType(val value.Value, targetType types.Type) value.Value {
	switch t := targetType.(type) {
	case *types.Basic:
		return convertBasicValue(val, t.Kind())
	case *types.Slice:
		return convertSliceValue(val, t)
	case *types.Named:
		return v.convertNamedValue(val, t)
	default:
		return val
	}
}

func convertBasicValue(val value.Value, kind types.BasicKind) value.Value {
	switch kind {
	case types.String:
		return convertValueToString(val)
	case types.Int:
		return value.MakeInt(toInt64(val))
	case types.Int8:
		return value.MakeInt8(int8(toInt64(val)))
	case types.Int16:
		return value.MakeInt16(int16(toInt64(val)))
	case types.Int32:
		return value.MakeInt32(int32(toInt64(val)))
	case types.Int64:
		return value.MakeInt64(toInt64(val))
	case types.Uint:
		return value.MakeUint(toUint64(val))
	case types.Uint8:
		return value.MakeUint8(uint8(toUint64(val)))
	case types.Uint16:
		return value.MakeUint16(uint16(toUint64(val)))
	case types.Uint32:
		return value.MakeUint32(uint32(toUint64(val)))
	case types.Uint64, types.Uintptr:
		return value.MakeUint64(toUint64(val))
	case types.Float32:
		return value.MakeFloat32(float32(toFloat64(val)))
	case types.Float64:
		return value.MakeFloat(toFloat64(val))
	default:
		return val
	}
}

func convertValueToString(val value.Value) value.Value {
	switch val.Kind() {
	case value.KindInt:
		return value.MakeString(string(rune(val.Int())))
	case value.KindUint:
		return value.MakeString(string(byte(val.Uint())))
	case value.KindString:
		return val
	case value.KindBytes:
		return convertBytesToString(val)
	case value.KindReflect:
		return convertReflectValueToString(val)
	default:
		return value.MakeString(fmt.Sprintf("%v", val.Interface()))
	}
}

func convertBytesToString(val value.Value) value.Value {
	if b, ok := val.Bytes(); ok {
		return value.MakeString(string(b))
	}
	return value.MakeString("")
}

func convertReflectValueToString(val value.Value) value.Value {
	rv, ok := val.ReflectValue()
	if !ok || rv.Kind() != reflect.Slice {
		return value.MakeString(fmt.Sprintf("%v", val.Interface()))
	}
	if rv.Type().Elem().Kind() != reflect.Int32 {
		return value.MakeString(fmt.Sprintf("%v", val.Interface()))
	}
	runes := make([]rune, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		runes[i] = rune(rv.Index(i).Int())
	}
	return value.MakeString(string(runes))
}

func convertSliceValue(val value.Value, target *types.Slice) value.Value {
	basic, ok := target.Elem().(*types.Basic)
	if !ok || val.Kind() != value.KindString {
		return val
	}
	switch basic.Kind() {
	case types.Int32:
		return convertStringToRuneSlice(val.String())
	case types.Uint8:
		return convertStringToByteSlice(val.String())
	default:
		return val
	}
}

func convertStringToRuneSlice(s string) value.Value {
	runes := []rune(s)
	rs := reflect.MakeSlice(reflect.TypeOf([]int32{}), len(runes), len(runes))
	for i, r := range runes {
		rs.Index(i).SetInt(int64(r))
	}
	return value.MakeFromReflect(rs)
}

func convertStringToByteSlice(s string) value.Value {
	// Use make+copy to ensure cap==len, matching native Go compiler behavior.
	// Direct []byte(s) uses runtime's stringtoslicebyte which rounds capacity
	// up to allocator size classes, while compiler-optimized conversions do not.
	b := make([]byte, len(s))
	copy(b, s)
	return value.MakeBytes(b)
}

func (v *vm) convertNamedValue(val value.Value, target *types.Named) value.Value {
	// Named-type conversion (for example []int -> sort.IntSlice) must resolve
	// the real host type before reflecting, otherwise the underlying slice type
	// is correct but the named external type is lost.
	targetRT := typeToReflect(target, v.program)
	if targetRT == nil {
		return val
	}
	rv := val.ToReflectValue(targetRT)
	if !rv.IsValid() {
		return val
	}
	if rv.Type() != targetRT && rv.Type().ConvertibleTo(targetRT) {
		rv = rv.Convert(targetRT)
	}
	return value.MakeFromReflect(rv)
}
