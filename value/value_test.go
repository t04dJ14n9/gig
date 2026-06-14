package value

import (
	"go/types"
	"math"
	"reflect"
	"testing"
)

// --- Value: constructors and accessors --------------------------------------

func TestValue_Primitives_RoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   any
	}{
		{"bool true", true},
		{"bool false", false},
		{"int", int(42)},
		{"int8", int8(-7)},
		{"int16", int16(-300)},
		{"int32", int32(-50_000)},
		{"int64", int64(-1 << 40)},
		{"uint", uint(42)},
		{"uint8", uint8(255)},
		{"uint16", uint16(65535)},
		{"uint32", uint32(1 << 31)},
		{"uint64", uint64(1 << 60)},
		{"float32", float32(3.14)},
		{"float64", float64(2.718281828)},
		{"complex64", complex64(complex(1, 2))},
		{"complex128", complex(3.0, 4.0)},
		{"string", "hello"},
		{"empty string", ""},
	}

	c := DefaultConverter()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := c.FromAny(tc.in)
			if err != nil {
				t.Fatalf("FromAny: %v", err)
			}
			out, err := c.ToAny(v)
			if err != nil {
				t.Fatalf("ToAny: %v", err)
			}
			if !reflect.DeepEqual(out, tc.in) {
				t.Fatalf("round-trip mismatch: got %v (%T), want %v (%T)",
					out, out, tc.in, tc.in)
			}
		})
	}
}

func TestValue_TypedAccessors(t *testing.T) {
	if !MakeBool(true).Bool() {
		t.Fatal("MakeBool(true).Bool() should be true")
	}
	if MakeInt(7).Int() != 7 {
		t.Fatalf("MakeInt(7).Int() = %d", MakeInt(7).Int())
	}
	if MakeUint(9).Uint() != 9 {
		t.Fatalf("MakeUint(9).Uint() = %d", MakeUint(9).Uint())
	}
	if MakeFloat(1.5).Float() != 1.5 {
		t.Fatalf("MakeFloat(1.5).Float() = %v", MakeFloat(1.5).Float())
	}
	if MakeString("x").Str() != "x" {
		t.Fatalf("MakeString(\"x\").Str() = %q", MakeString("x").Str())
	}
	if MakeComplex(1, 2).Complex() != complex(1, 2) {
		t.Fatalf("MakeComplex round-trip failed")
	}
}

func TestValue_AccessorPanicsOnKindMismatch(t *testing.T) {
	mustPanic(t, "Bool on int", func() { _ = MakeInt(1).Bool() })
	mustPanic(t, "Int on bool", func() { _ = MakeBool(true).Int() })
	mustPanic(t, "Float on string", func() { _ = MakeString("x").Float() })
	mustPanic(t, "Str on int", func() { _ = MakeInt(1).Str() })
}

// --- Value: identity, validity, nil-semantics -------------------------------

func TestValue_IsValid(t *testing.T) {
	var zero Value
	if zero.IsValid() {
		t.Fatal("zero Value should not be valid")
	}
	if !MakeInt(0).IsValid() {
		t.Fatal("MakeInt(0) should be valid")
	}
	if !MakeNil().IsValid() {
		t.Fatal("MakeNil() should be valid (it has KindNil, not KindInvalid)")
	}
}

func TestValue_IsNil(t *testing.T) {
	if !MakeNil().IsNil() {
		t.Fatal("MakeNil should be nil")
	}
	if MakeInt(0).IsNil() {
		t.Fatal("MakeInt(0) is not nil")
	}
	// reflect-backed nil slice
	var s []int
	rv, _ := DefaultConverter().FromAny(s)
	if !rv.IsNil() {
		t.Fatalf("reflect-backed nil slice should report IsNil, got %#v", rv)
	}
	// reflect-backed non-nil slice
	rv2, _ := DefaultConverter().FromAny([]int{1})
	if rv2.IsNil() {
		t.Fatal("non-nil slice should not report IsNil")
	}
}

// --- Value: size preservation through Interface() ---------------------------

func TestValue_SizePreservation_Int(t *testing.T) {
	c := DefaultConverter()
	cases := []struct {
		in   any
		want reflect.Kind
	}{
		{int8(1), reflect.Int8},
		{int16(1), reflect.Int16},
		{int32(1), reflect.Int32},
		{int64(1), reflect.Int64},
		{int(1), reflect.Int},
	}
	for _, tc := range cases {
		v, _ := c.FromAny(tc.in)
		got := reflect.TypeOf(v.Interface()).Kind()
		if got != tc.want {
			t.Errorf("FromAny(%T(%v)).Interface() kind = %s, want %s",
				tc.in, tc.in, got, tc.want)
		}
	}
}

func TestValue_SizePreservation_Float(t *testing.T) {
	c := DefaultConverter()
	v32, _ := c.FromAny(float32(1.5))
	if reflect.TypeOf(v32.Interface()).Kind() != reflect.Float32 {
		t.Errorf("float32 should round-trip as float32")
	}
	v64, _ := c.FromAny(float64(1.5))
	if reflect.TypeOf(v64.Interface()).Kind() != reflect.Float64 {
		t.Errorf("float64 should round-trip as float64")
	}
}

// --- Converter: reflect handling --------------------------------------------

func TestConverter_FromReflect_Primitives(t *testing.T) {
	c := DefaultConverter()
	v, err := c.FromReflect(reflect.ValueOf(int32(7)))
	if err != nil {
		t.Fatalf("FromReflect: %v", err)
	}
	if v.Kind() != KindInt || v.Int() != 7 || v.SizeTag() != Size32 {
		t.Fatalf("FromReflect lost info: kind=%s int=%d size=%d",
			v.Kind(), v.Int(), v.SizeTag())
	}
}

func TestConverter_FromReflect_Composite(t *testing.T) {
	c := DefaultConverter()
	v, err := c.FromReflect(reflect.ValueOf([]string{"a", "b"}))
	if err != nil {
		t.Fatalf("FromReflect: %v", err)
	}
	if v.Kind() != KindReflect {
		t.Fatalf("composite should land in KindReflect, got %s", v.Kind())
	}
	out := v.Interface().([]string)
	if len(out) != 2 || out[0] != "a" {
		t.Fatalf("composite round-trip wrong: %v", out)
	}
}

func TestConverter_FromReflect_Invalid(t *testing.T) {
	c := DefaultConverter()
	v, err := c.FromReflect(reflect.Value{})
	if err != nil {
		t.Fatalf("FromReflect on invalid: %v", err)
	}
	if !v.IsNil() {
		t.Fatal("invalid reflect.Value should produce nil Value")
	}
}

func TestConverter_ToReflect(t *testing.T) {
	c := DefaultConverter()
	v := MakeInt8(5)
	rv, err := c.ToReflect(v, reflect.TypeOf(int(0)))
	if err != nil {
		t.Fatalf("ToReflect: %v", err)
	}
	if rv.Kind() != reflect.Int || rv.Int() != 5 {
		t.Fatalf("ToReflect int8->int: %v %d", rv.Kind(), rv.Int())
	}
}

func TestConverter_ToReflect_NilProducesZero(t *testing.T) {
	c := DefaultConverter()
	rv, err := c.ToReflect(MakeNil(), reflect.TypeOf(int(0)))
	if err != nil {
		t.Fatalf("ToReflect: %v", err)
	}
	if rv.Int() != 0 {
		t.Fatalf("nil should map to zero, got %d", rv.Int())
	}
}

// --- Converter: FromAny passthrough cases -----------------------------------

func TestConverter_FromAny_PassthroughValue(t *testing.T) {
	c := DefaultConverter()
	v := MakeInt(7)
	out, err := c.FromAny(v)
	if err != nil {
		t.Fatalf("FromAny: %v", err)
	}
	if out.Kind() != KindInt || out.Int() != 7 {
		t.Fatalf("Value pass-through failed: %#v", out)
	}
}

func TestConverter_FromAny_NilProducesNilKind(t *testing.T) {
	c := DefaultConverter()
	out, err := c.FromAny(nil)
	if err != nil {
		t.Fatalf("FromAny(nil): %v", err)
	}
	if out.Kind() != KindNil {
		t.Fatalf("nil should produce KindNil, got %s", out.Kind())
	}
}

// --- Converter: Zero from go/types ------------------------------------------

func TestConverter_Zero_BasicTypes(t *testing.T) {
	c := DefaultConverter()
	cases := []struct {
		t        types.Type
		wantKind Kind
	}{
		{types.Typ[types.Bool], KindBool},
		{types.Typ[types.Int], KindInt},
		{types.Typ[types.Int8], KindInt},
		{types.Typ[types.Uint64], KindUint},
		{types.Typ[types.Float32], KindFloat},
		{types.Typ[types.Float64], KindFloat},
		{types.Typ[types.String], KindString},
		{types.Typ[types.Complex128], KindComplex},
	}
	for _, tc := range cases {
		t.Run(tc.t.String(), func(t *testing.T) {
			v, err := c.Zero(tc.t, nil)
			if err != nil {
				t.Fatalf("Zero(%s): %v", tc.t, err)
			}
			if v.Kind() != tc.wantKind {
				t.Fatalf("Zero(%s) kind = %s, want %s",
					tc.t, v.Kind(), tc.wantKind)
			}
		})
	}
}

func TestConverter_Zero_BasicSizesAreCorrect(t *testing.T) {
	c := DefaultConverter()
	v, err := c.Zero(types.Typ[types.Int8], nil)
	if err != nil {
		t.Fatalf("Zero: %v", err)
	}
	got, ok := v.Interface().(int8)
	if !ok {
		t.Fatalf("Zero(int8).Interface() = %T, want int8", v.Interface())
	}
	if got != 0 {
		t.Fatalf("Zero(int8) = %d, want 0", got)
	}
}

func TestConverter_Zero_RequiresResolverForCompositeTypes(t *testing.T) {
	c := DefaultConverter()
	composite := types.NewSlice(types.Typ[types.Int])
	if _, err := c.Zero(composite, nil); err == nil {
		t.Fatal("Zero of slice without TypeResolver should error")
	}
}

func TestConverter_Zero_UsesResolverForCompositeTypes(t *testing.T) {
	c := DefaultConverter()
	composite := types.NewSlice(types.Typ[types.Int])
	r := stubResolver{m: map[string]reflect.Type{
		composite.String(): reflect.TypeOf([]int(nil)),
	}}
	v, err := c.Zero(composite, r)
	if err != nil {
		t.Fatalf("Zero: %v", err)
	}
	if !v.IsValid() {
		t.Fatalf("zero composite should be valid")
	}
}

// --- Converter: Convert -----------------------------------------------------

func TestConverter_Convert_NarrowingInt(t *testing.T) {
	c := DefaultConverter()
	r := stubResolver{m: map[string]reflect.Type{
		types.Typ[types.Int8].String(): reflect.TypeOf(int8(0)),
	}}
	v := MakeInt(127)
	out, err := c.Convert(v, types.Typ[types.Int8], r)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if out.Kind() != KindInt || out.SizeTag() != Size8 || out.Int() != 127 {
		t.Fatalf("Convert int->int8 wrong: %#v", out)
	}
}

func TestConverter_Convert_FloatToInt(t *testing.T) {
	c := DefaultConverter()
	r := stubResolver{m: map[string]reflect.Type{
		types.Typ[types.Int].String(): reflect.TypeOf(int(0)),
	}}
	v := MakeFloat(3.7)
	out, err := c.Convert(v, types.Typ[types.Int], r)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if out.Int() != 3 { // Go truncates float -> int
		t.Fatalf("float->int should truncate to 3, got %d", out.Int())
	}
}

// --- Float bit-pattern preservation ----------------------------------------

func TestValue_Float_PreservesBitPattern(t *testing.T) {
	bits := math.Float64bits(math.Pi)
	v := MakeFloat(math.Pi)
	if math.Float64bits(v.Float()) != bits {
		t.Fatal("MakeFloat lost bits")
	}
}

func TestValue_Float32_PreservesAfterRoundTrip(t *testing.T) {
	c := DefaultConverter()
	v, _ := c.FromAny(float32(0.1))
	out := v.Interface().(float32)
	if out != float32(0.1) {
		t.Fatalf("float32 round-trip lost precision: got %v", out)
	}
}

// --- helpers ----------------------------------------------------------------

type stubResolver struct{ m map[string]reflect.Type }

func (s stubResolver) ResolveType(t types.Type) (reflect.Type, error) {
	if rt, ok := s.m[t.String()]; ok {
		return rt, nil
	}
	return nil, &resolverError{t: t}
}

type resolverError struct{ t types.Type }

func (e *resolverError) Error() string { return "stub: no type for " + e.t.String() }

func mustPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: expected panic, got none", name)
		}
	}()
	f()
}
