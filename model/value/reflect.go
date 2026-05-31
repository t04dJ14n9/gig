package value

import "reflect"

// MakeFromReflect creates a Value from reflect.Value.
func MakeFromReflect(rv reflect.Value) Value {
	if !rv.IsValid() {
		return MakeNil()
	}

	if val, ok := reflectedValueStruct(rv); ok {
		return val
	}
	if val, ok := reflectPrimitiveValue(rv); ok {
		return val
	}
	if val, ok := reflectNativeSliceValue(rv); ok {
		return val
	}
	return makeReflectValue(rv)
}

// FromInterface creates a Value from any Go value.
func FromInterface(v any) Value {
	if v == nil {
		return MakeNil()
	}

	if val, ok := primitiveInterfaceValue(v); ok {
		return val
	}
	if val, ok := specialInterfaceValue(v); ok {
		return val
	}
	return MakeFromReflect(reflect.ValueOf(v))
}

func reflectedValueStruct(rv reflect.Value) (Value, bool) {
	if rv.Kind() != reflect.Struct {
		return Value{}, false
	}
	if !isValueStructType(rv.Type()) {
		return Value{}, false
	}
	// Keep value.Value identity intact when reflected interpreter storage flows
	// back through host APIs; wrapping it as KindReflect would hide its tag.
	return rv.Interface().(Value), true
}

func isValueStructType(t reflect.Type) bool {
	return t.Name() == "Value" && t.PkgPath() == "github.com/t04dJ14n9/gig/model/value"
}

func reflectPrimitiveValue(rv reflect.Value) (Value, bool) {
	switch rv.Kind() {
	case reflect.Bool:
		return MakeBool(rv.Bool()), true
	case reflect.String:
		return MakeString(rv.String()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectSignedValue(rv), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflectUnsignedValue(rv), true
	case reflect.Float32, reflect.Float64:
		return reflectFloatValue(rv), true
	case reflect.Complex64, reflect.Complex128:
		return reflectComplexValue(rv), true
	default:
		return Value{}, false
	}
}

func reflectSignedValue(rv reflect.Value) Value {
	switch rv.Kind() {
	case reflect.Int:
		return MakeInt(rv.Int())
	case reflect.Int8:
		return MakeInt8(int8(rv.Int()))
	case reflect.Int16:
		return MakeInt16(int16(rv.Int()))
	case reflect.Int32:
		return MakeInt32(int32(rv.Int()))
	default:
		return MakeInt64(rv.Int())
	}
}

func reflectUnsignedValue(rv reflect.Value) Value {
	switch rv.Kind() {
	case reflect.Uint:
		return MakeUint(rv.Uint())
	case reflect.Uint8:
		return MakeUint8(uint8(rv.Uint()))
	case reflect.Uint16:
		return MakeUint16(uint16(rv.Uint()))
	case reflect.Uint32:
		return MakeUint32(uint32(rv.Uint()))
	default:
		// Uint64 and Uintptr share the existing uint64 Value representation.
		return MakeUint64(rv.Uint())
	}
}

func reflectFloatValue(rv reflect.Value) Value {
	if rv.Kind() == reflect.Float32 {
		return MakeFloat32(float32(rv.Float()))
	}
	return MakeFloat(rv.Float())
}

func reflectComplexValue(rv reflect.Value) Value {
	c := rv.Complex()
	if rv.Kind() == reflect.Complex64 {
		return MakeComplex64(float32(real(c)), float32(imag(c)))
	}
	return MakeComplex(real(c), imag(c))
}

func reflectNativeSliceValue(rv reflect.Value) (Value, bool) {
	if rv.Kind() != reflect.Slice {
		return Value{}, false
	}

	elemKind := rv.Type().Elem().Kind()
	if elemKind == reflect.Uint8 {
		// []byte is common in external APIs; store it natively so callers avoid
		// reflect overhead and preserve byte-slice formatting behavior.
		return MakeBytes(rv.Bytes()), true
	}
	if elemKind == reflect.Int64 {
		// []int64 is the interpreter's native integer slice representation.
		// Keeping it as KindSlice lets Interface() convert it back to []int.
		return MakeIntSlice(rv.Interface().([]int64)), true
	}
	return Value{}, false
}

func primitiveInterfaceValue(v any) (Value, bool) {
	// Fast path common scalar values to avoid reflect.ValueOf while preserving
	// the original Go width in the Value size tag.
	switch val := v.(type) {
	case bool:
		return MakeBool(val), true
	case int:
		return MakeInt(int64(val)), true
	case int8:
		return MakeInt8(val), true
	case int16:
		return MakeInt16(val), true
	case int32:
		return MakeInt32(val), true
	case int64:
		return MakeInt64(val), true
	case uint:
		return MakeUint(uint64(val)), true
	case uint8:
		return MakeUint8(val), true
	case uint16:
		return MakeUint16(val), true
	case uint32:
		return MakeUint32(val), true
	case uint64:
		return MakeUint64(val), true
	case float32:
		return MakeFloat32(val), true
	case float64:
		return MakeFloat(val), true
	case complex64:
		return MakeComplex64(real(val), imag(val)), true
	case complex128:
		return MakeComplex(real(val), imag(val)), true
	case string:
		return MakeString(val), true
	default:
		return Value{}, false
	}
}

func specialInterfaceValue(v any) (Value, bool) {
	switch val := v.(type) {
	case []byte:
		return MakeBytes(val), true
	case reflect.Value:
		// Unwrap reflect.Value directly (e.g., typed nil constants from the compiler).
		return MakeFromReflect(val), true
	default:
		return Value{}, false
	}
}

func makeReflectValue(rv reflect.Value) Value {
	return Value{kind: KindReflect, obj: rv}
}
