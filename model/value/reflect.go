package value

import "reflect"

// MakeFromReflect creates a Value from reflect.Value.
func MakeFromReflect(rv reflect.Value) Value {
	if !rv.IsValid() {
		return MakeNil()
	}

	// Check if the underlying value is already a Value - unwrap it
	if rv.Kind() == reflect.Struct {
		if t := rv.Type(); t.Name() == "Value" && t.PkgPath() == "github.com/t04dJ14n9/gig/model/value" {
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
	case reflect.Complex64:
		c := rv.Complex()
		return MakeComplex64(float32(real(c)), float32(imag(c)))
	case reflect.Complex128:
		c := rv.Complex()
		return MakeComplex(real(c), imag(c))
	case reflect.Slice:
		// Native []byte storage: avoid reflect overhead for the common []byte case
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return MakeBytes(rv.Bytes())
		}
		// Native []int64 storage: use KindSlice so Interface() converts to []int correctly
		if rv.Type().Elem().Kind() == reflect.Int64 {
			return MakeIntSlice(rv.Interface().([]int64))
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
	case complex64:
		return MakeComplex64(real(val), imag(val))
	case complex128:
		return MakeComplex(real(val), imag(val))
	case string:
		return MakeString(val)
	case []byte:
		return MakeBytes(val)
	case reflect.Value:
		// Unwrap reflect.Value directly (e.g., typed nil constants from the compiler).
		return MakeFromReflect(val)
	}
	return MakeFromReflect(reflect.ValueOf(v))
}
