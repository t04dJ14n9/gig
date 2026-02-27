// Package value implements a tagged-union Value system for high-performance interpretation.
// Primitive types (bool, int, uint, float, string) use native Go operations with zero reflect overhead.
// Complex types fall back to reflect.Value.
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

// --- Accessors ---

// Bool returns the bool value. Panics if not KindBool.
func (v Value) Bool() bool {
	if v.kind != KindBool {
		panic(fmt.Sprintf("not a bool: %v", v.kind))
	}
	return v.num != 0
}

// Int returns the int value. Panics if not KindInt.
func (v Value) Int() int64 {
	if v.kind != KindInt {
		panic(fmt.Sprintf("not an int: %v", v.kind))
	}
	return v.num
}

// Uint returns the uint value. Panics if not KindUint.
func (v Value) Uint() uint64 {
	if v.kind != KindUint {
		panic(fmt.Sprintf("not a uint: %v", v.kind))
	}
	return uint64(v.num)
}

// Float returns the float value. Panics if not KindFloat.
func (v Value) Float() float64 {
	if v.kind != KindFloat {
		panic(fmt.Sprintf("not a float: %v", v.kind))
	}
	return math.Float64frombits(uint64(v.num))
}

// String returns the string value. Panics if not KindString.
func (v Value) String() string {
	if v.kind != KindString {
		panic(fmt.Sprintf("not a string: %v", v.kind))
	}
	return v.str
}

// Complex returns the complex value. Panics if not KindComplex.
func (v Value) Complex() complex128 {
	if v.kind != KindComplex {
		panic(fmt.Sprintf("not a complex: %v", v.kind))
	}
	return complex(math.Float64frombits(uint64(v.num)), math.Float64frombits(uint64(v.num2)))
}

// Interface returns the value as an interface{}.
func (v Value) Interface() any {
	switch v.kind {
	case KindNil:
		return nil
	case KindBool:
		return v.Bool()
	case KindInt:
		return v.Int()
	case KindUint:
		return v.Uint()
	case KindFloat:
		return v.Float()
	case KindString:
		return v.str
	case KindComplex:
		return v.Complex()
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Interface()
		}
		return v.obj
	default:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Interface()
		}
		return v.obj
	}
}

// ToReflectValue converts to reflect.Value.
func (v Value) ToReflectValue(typ reflect.Type) reflect.Value {
	switch v.kind {
	case KindNil:
		return reflect.Zero(typ)
	case KindBool:
		return reflect.ValueOf(v.Bool())
	case KindInt:
		return reflect.ValueOf(v.num).Convert(typ)
	case KindUint:
		return reflect.ValueOf(uint64(v.num)).Convert(typ)
	case KindFloat:
		return reflect.ValueOf(v.Float()).Convert(typ)
	case KindString:
		return reflect.ValueOf(v.str)
	case KindComplex:
		return reflect.ValueOf(v.Complex())
	case KindReflect:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	default:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv
		}
		return reflect.ValueOf(v.obj)
	}
}

// ReflectValue returns the internal reflect.Value if stored.
func (v Value) ReflectValue() (reflect.Value, bool) {
	rv, ok := v.obj.(reflect.Value)
	return rv, ok
}

// --- Arithmetic Operations ---

// Add returns v + other.
func (v Value) Add(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num + other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) + other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() + other.Float())
	case KindString:
		return MakeString(v.str + other.String())
	case KindComplex:
		r1, i1 := v.realImag()
		r2, i2 := other.realImag()
		return MakeComplex(r1+r2, i1+i2)
	default:
		panic(fmt.Sprintf("cannot add %v", v.kind))
	}
}

// Sub returns v - other.
func (v Value) Sub(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num - other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) - other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() - other.Float())
	case KindComplex:
		r1, i1 := v.realImag()
		r2, i2 := other.realImag()
		return MakeComplex(r1-r2, i1-i2)
	default:
		panic(fmt.Sprintf("cannot sub %v", v.kind))
	}
}

// Mul returns v * other.
func (v Value) Mul(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num * other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) * other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() * other.Float())
	case KindComplex:
		r1, i1 := v.realImag()
		r2, i2 := other.realImag()
		return MakeComplex(r1*r2-i1*i2, r1*i2+r2*i1)
	default:
		panic(fmt.Sprintf("cannot mul %v", v.kind))
	}
}

// Div returns v / other.
func (v Value) Div(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num / other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) / other.Uint())
	case KindFloat:
		return MakeFloat(v.Float() / other.Float())
	case KindComplex:
		r1, i1 := v.realImag()
		r2, i2 := other.realImag()
		denom := r2*r2 + i2*i2
		return MakeComplex((r1*r2+i1*i2)/denom, (i1*r2-r1*i2)/denom)
	default:
		panic(fmt.Sprintf("cannot div %v", v.kind))
	}
}

// Mod returns v % other.
func (v Value) Mod(other Value) Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(v.num % other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) % other.Uint())
	case KindFloat:
		return MakeFloat(math.Mod(v.Float(), other.Float()))
	default:
		panic(fmt.Sprintf("cannot mod %v", v.kind))
	}
}

// Neg returns -v.
func (v Value) Neg() Value {
	switch v.kind { //nolint:exhaustive
	case KindInt:
		return MakeInt(-v.num)
	case KindFloat:
		return MakeFloat(-v.Float())
	case KindComplex:
		r, i := v.realImag()
		return MakeComplex(-r, -i)
	default:
		panic(fmt.Sprintf("cannot neg %v", v.kind))
	}
}

// --- Comparison Operations ---

// Cmp compares v with other. Returns -1, 0, or 1.
func (v Value) Cmp(other Value) int {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		a, b := v.Bool(), other.Bool()
		if a == b {
			return 0
		}
		if !a {
			return -1
		}
		return 1
	case KindInt:
		a, b := v.num, other.Int()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindUint:
		a, b := uint64(v.num), other.Uint()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindFloat:
		a, b := v.Float(), other.Float()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case KindString:
		a, b := v.str, other.String()
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("cannot compare %v", v.kind))
	}
}

// Equal returns v == other.
func (v Value) Equal(other Value) bool {
	if v.kind != other.kind {
		// Handle nil comparison
		if v.kind == KindNil || other.kind == KindNil {
			return v.IsNil() && other.IsNil()
		}
		return false
	}
	switch v.kind {
	case KindNil:
		return true
	case KindBool:
		return v.num == other.num
	case KindInt:
		return v.num == other.num
	case KindUint:
		return v.num == other.num
	case KindFloat:
		return v.Float() == other.Float()
	case KindString:
		return v.str == other.str
	case KindComplex:
		return v.num == other.num && v.num2 == other.num2
	default:
		// For complex types, use reflect.DeepEqual
		return reflect.DeepEqual(v.Interface(), other.Interface())
	}
}

// --- Bitwise Operations ---

// And returns v & other.
func (v Value) And(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num & other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) & other.Uint())
	default:
		panic(fmt.Sprintf("cannot and %v", v.kind))
	}
}

// Or returns v | other.
func (v Value) Or(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num | other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) | other.Uint())
	default:
		panic(fmt.Sprintf("cannot or %v", v.kind))
	}
}

// Xor returns v ^ other.
func (v Value) Xor(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num ^ other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) ^ other.Uint())
	default:
		panic(fmt.Sprintf("cannot xor %v", v.kind))
	}
}

// AndNot returns v &^ other.
func (v Value) AndNot(other Value) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num &^ other.Int())
	case KindUint:
		return MakeUint(uint64(v.num) &^ other.Uint())
	default:
		panic(fmt.Sprintf("cannot andnot %v", v.kind))
	}
}

// Lsh returns v << n.
func (v Value) Lsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num << n)
	case KindUint:
		return MakeUint(uint64(v.num) << n)
	default:
		panic(fmt.Sprintf("cannot lsh %v", v.kind))
	}
}

// Rsh returns v >> n.
func (v Value) Rsh(n uint) Value {
	switch v.kind {
	case KindInt:
		return MakeInt(v.num >> n)
	case KindUint:
		return MakeUint(uint64(v.num) >> n)
	default:
		panic(fmt.Sprintf("cannot rsh %v", v.kind))
	}
}

// --- Type Conversions ---

// ToInt converts to int.
func (v Value) ToInt() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeInt(1)
		}
		return MakeInt(0)
	case KindInt:
		return v
	case KindUint:
		return MakeInt(int64(v.Uint()))
	case KindFloat:
		return MakeInt(int64(v.Float()))
	default:
		panic(fmt.Sprintf("cannot convert %v to int", v.kind))
	}
}

// ToUint converts to uint.
func (v Value) ToUint() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeUint(1)
		}
		return MakeUint(0)
	case KindInt:
		return MakeUint(uint64(v.num))
	case KindUint:
		return v
	case KindFloat:
		return MakeUint(uint64(v.Float()))
	default:
		panic(fmt.Sprintf("cannot convert %v to uint", v.kind))
	}
}

// ToFloat converts to float.
func (v Value) ToFloat() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		if v.Bool() {
			return MakeFloat(1.0)
		}
		return MakeFloat(0.0)
	case KindInt:
		return MakeFloat(float64(v.num))
	case KindUint:
		return MakeFloat(float64(v.Uint()))
	case KindFloat:
		return v
	default:
		panic(fmt.Sprintf("cannot convert %v to float", v.kind))
	}
}

// ToBool converts to bool.
func (v Value) ToBool() Value {
	switch v.kind { //nolint:exhaustive
	case KindBool:
		return v
	case KindInt:
		return MakeBool(v.num != 0)
	case KindUint:
		return MakeBool(v.num != 0)
	case KindFloat:
		return MakeBool(v.Float() != 0)
	case KindString:
		return MakeBool(v.str != "")
	default:
		return MakeBool(!v.IsNil())
	}
}

// ToString converts to string representation.
func (v Value) ToString() Value {
	return MakeString(fmt.Sprintf("%v", v.Interface()))
}

// --- Helper methods ---

func (v Value) realImag() (float64, float64) {
	return math.Float64frombits(uint64(v.num)), math.Float64frombits(uint64(v.num2))
}

// Len returns the length of string, slice, array, map, or chan.
func (v Value) Len() int {
	switch v.kind {
	case KindString:
		return len(v.str)
	case KindSlice, KindArray, KindMap, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Len()
		}
		panic("invalid reflect.Value in Len()")
	default:
		panic(fmt.Sprintf("cannot take len of %v", v.kind))
	}
}

// Cap returns the capacity of slice, array, or chan.
func (v Value) Cap() int {
	switch v.kind {
	case KindSlice, KindArray, KindChan:
		if rv, ok := v.obj.(reflect.Value); ok {
			return rv.Cap()
		}
		panic("invalid reflect.Value in Cap()")
	default:
		panic(fmt.Sprintf("cannot take cap of %v", v.kind))
	}
}

// Index returns element at index i for slice, array, or string.
func (v Value) Index(i int) Value {
	switch v.kind {
	case KindString:
		// s[i] returns a byte (uint8), not a string
		return MakeUint(uint64(v.str[i]))
	case KindSlice, KindArray:
		if rv, ok := v.obj.(reflect.Value); ok {
			elem := rv.Index(i)
			// For function element types, unwrap the stored Value
			if rv.Type().Elem().Kind() == reflect.Func {
				if val, ok := elem.Interface().(Value); ok {
					return val
				}
			}
			// For []value.Value slices (used for function slices)
			if rv.Type().Elem() == reflect.TypeOf(Value{}) {
				return elem.Interface().(Value)
			}
			return MakeFromReflect(elem)
		}
		// Handle native []value.Value slice
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	case KindReflect:
		// Handle reflect.Value containing a slice
		if rv, ok := v.obj.(reflect.Value); ok {
			elem := rv.Index(i)
			// For function element types, unwrap the stored Value
			if rv.Type().Elem().Kind() == reflect.Func {
				if val, ok := elem.Interface().(Value); ok {
					return val
				}
			}
			// For []value.Value slices (used for function slices)
			if rv.Type().Elem() == reflect.TypeOf(Value{}) {
				return elem.Interface().(Value)
			}
			return MakeFromReflect(elem)
		}
		// Handle native []value.Value slice
		if slice, ok := v.obj.([]Value); ok {
			return slice[i]
		}
		panic("invalid reflect.Value in Index()")
	default:
		panic(fmt.Sprintf("cannot index %v", v.kind))
	}
}

// SetIndex sets element at index i for slice or array.
func (v Value) SetIndex(i int, val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		elemType := rv.Type().Elem()
		// For function element types, store the Value directly (closures are *Closure)
		if elemType.Kind() == reflect.Func {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		// For []value.Value slices (used for function slices)
		if elemType == reflect.TypeOf(Value{}) {
			rv.Index(i).Set(reflect.ValueOf(val))
			return
		}
		rv.Index(i).Set(val.ToReflectValue(elemType))
		return
	}
	// Handle native []value.Value slice
	if slice, ok := v.obj.([]Value); ok {
		slice[i] = val
		return
	}
	panic("invalid reflect.Value in SetIndex()")
}

// MapIndex returns value at key k for map.
func (v Value) MapIndex(k Value) Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		key := k.ToReflectValue(rv.Type().Key())
		elem := rv.MapIndex(key)
		if !elem.IsValid() {
			// Return zero value of element type, not nil (Go semantics)
			return MakeFromReflect(reflect.Zero(rv.Type().Elem()))
		}
		return MakeFromReflect(elem)
	}
	panic("invalid reflect.Value in MapIndex()")
}

// SetMapIndex sets value at key k for map.
func (v Value) SetMapIndex(k, val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		key := k.ToReflectValue(rv.Type().Key())
		if val.IsNil() {
			rv.SetMapIndex(key, reflect.Value{})
		} else {
			rv.SetMapIndex(key, val.ToReflectValue(rv.Type().Elem()))
		}
		return
	}
	panic("invalid reflect.Value in SetMapIndex()")
}

// MapIter iterates over a map.
func (v Value) MapIter(f func(key, val Value) bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		iter := rv.MapRange()
		for iter.Next() {
			key := MakeFromReflect(iter.Key())
			val := MakeFromReflect(iter.Value())
			if !f(key, val) {
				break
			}
		}
		return
	}
	panic("invalid reflect.Value in MapIter()")
}

// Field returns struct field at index i.
func (v Value) Field(i int) Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		return MakeFromReflect(rv.Field(i))
	}
	panic("invalid reflect.Value in Field()")
}

// SetField sets struct field at index i.
func (v Value) SetField(i int, val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Field(i).Set(val.ToReflectValue(rv.Type().Field(i).Type))
		return
	}
	panic("invalid reflect.Value in SetField()")
}

// Elem dereferences a pointer or returns the underlying value of interface.
func (v Value) Elem() Value {
	if rv, ok := v.obj.(reflect.Value); ok {
		return MakeFromReflect(rv.Elem())
	}
	panic("invalid reflect.Value in Elem()")
}

// SetElem sets the value pointed to by a pointer.
func (v Value) SetElem(val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		// Handle different reflect.Value kinds
		kind := rv.Kind()
		if kind == reflect.Ptr {
			// Handle pointer case
			elemType := rv.Type().Elem()
			if elemType.Kind() == reflect.Func {
				rv.Elem().Set(reflect.ValueOf(val))
				return
			}
			if elemType.Name() == "Value" && elemType.PkgPath() == "gig/value" {
				ptr := rv.Interface().(*Value)
				*ptr = val
				return
			}
			targetRV := rv.Elem()
			if targetRV.CanSet() {
				targetRV.Set(val.ToReflectValue(elemType))
			}
			return
		}
		if kind == reflect.Interface {
			// For interface, just set the underlying value
			rv.Set(val.ToReflectValue(rv.Type()))
			return
		}
		if kind == reflect.Struct {
			// For struct values, we can't set elements - this shouldn't happen
			// but handle gracefully
			return
		}
	}
	panic("invalid reflect.Value in SetElem()")
}

// Pointer returns the underlying pointer value.
func (v Value) Pointer() uintptr {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.Pointer()
	}
	return 0
}

// Send sends a value on a channel.
func (v Value) Send(val Value) {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Send(val.ToReflectValue(rv.Type().Elem()))
		return
	}
	panic("invalid reflect.Value in Send()")
}

// TrySend tries to send a value on a channel (non-blocking).
func (v Value) TrySend(val Value) bool {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.TrySend(val.ToReflectValue(rv.Type().Elem()))
	}
	panic("invalid reflect.Value in TrySend()")
}

// Recv receives a value from a channel.
func (v Value) Recv() (Value, bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		val, ok := rv.Recv()
		return MakeFromReflect(val), ok
	}
	panic("invalid reflect.Value in Recv()")
}

// TryRecv tries to receive a value from a channel (non-blocking).
func (v Value) TryRecv() (Value, bool) {
	if rv, ok := v.obj.(reflect.Value); ok {
		val, ok := rv.TryRecv()
		return MakeFromReflect(val), ok
	}
	panic("invalid reflect.Value in TryRecv()")
}

// Close closes a channel.
func (v Value) Close() {
	if rv, ok := v.obj.(reflect.Value); ok {
		rv.Close()
		return
	}
	panic("invalid reflect.Value in Close()")
}

// CanInterface reports whether Interface can be used without panicking.
func (v Value) CanInterface() bool {
	if rv, ok := v.obj.(reflect.Value); ok {
		return rv.CanInterface()
	}
	return true
}

// Package packs multiple values into a slice.
func Package(vals ...Value) []Value {
	return vals
}

// Unpackage unpacks a slice of values.
func Unpackage(vals []Value) []Value {
	return vals
}

// GoString returns a Go-syntax representation of the value.
func (v Value) GoString() string {
	return fmt.Sprintf("value.Value{kind:%v, num:%d, str:%q, obj:%v}", v.kind, v.num, v.str, v.obj)
}
