package value

import (
	"math"
	"reflect"
	"testing"
)

// ---------------------------------------------------------------------------
// Constructors & Accessors
// ---------------------------------------------------------------------------

func TestMakeNilAndIsNil(t *testing.T) {
	v := MakeNil()
	if !v.IsNil() {
		t.Error("MakeNil().IsNil() should be true")
	}
	if v.Kind() != KindNil {
		t.Errorf("MakeNil().Kind() = %d, want KindNil", v.Kind())
	}
}

func TestMakeBool(t *testing.T) {
	tr := MakeBool(true)
	fa := MakeBool(false)
	if tr.Kind() != KindBool || !tr.Bool() {
		t.Error("MakeBool(true) failed")
	}
	if fa.Bool() {
		t.Error("MakeBool(false) should return false")
	}
}

func TestMakeInt(t *testing.T) {
	v := MakeInt(42)
	if v.Kind() != KindInt || v.Int() != 42 {
		t.Errorf("MakeInt(42): kind=%d, val=%d", v.Kind(), v.Int())
	}
}

func TestMakeUint(t *testing.T) {
	v := MakeUint(100)
	if v.Kind() != KindUint || v.Uint() != 100 {
		t.Errorf("MakeUint(100): kind=%d, val=%d", v.Kind(), v.Uint())
	}
}

func TestMakeFloat(t *testing.T) {
	v := MakeFloat(3.14)
	if v.Kind() != KindFloat || v.Float() != 3.14 {
		t.Errorf("MakeFloat(3.14): kind=%d, val=%f", v.Kind(), v.Float())
	}
}

func TestMakeString(t *testing.T) {
	v := MakeString("hello")
	if v.Kind() != KindString || v.String() != "hello" {
		t.Errorf("MakeString: kind=%d, val=%q", v.Kind(), v.String())
	}
}

func TestMakeComplex(t *testing.T) {
	v := MakeComplex(1.0, 2.0)
	if v.Kind() != KindComplex {
		t.Fatalf("MakeComplex kind = %d", v.Kind())
	}
	c := v.Complex()
	if real(c) != 1.0 || imag(c) != 2.0 {
		t.Errorf("MakeComplex(1,2) = %v", c)
	}
}

func TestMakeFromReflect(t *testing.T) {
	rv := reflect.ValueOf([]int{1, 2, 3})
	v := MakeFromReflect(rv)
	if v.Kind() != KindReflect {
		t.Fatalf("MakeFromReflect kind = %d", v.Kind())
	}
	got, ok := v.ReflectValue()
	if !ok {
		t.Fatal("ReflectValue should succeed for MakeFromReflect")
	}
	if got.Len() != 3 {
		t.Errorf("ReflectValue slice len = %d, want 3", got.Len())
	}
}

func TestFromInterface(t *testing.T) {
	tests := []struct {
		name string
		in   any
		kind Kind
	}{
		{"nil", nil, KindNil},
		{"bool", true, KindBool},
		{"int", 42, KindInt},
		{"string", "hi", KindString},
		{"float64", 1.5, KindFloat},
		{"slice", []int{1}, KindReflect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := FromInterface(tt.in)
			if v.Kind() != tt.kind {
				t.Errorf("FromInterface(%v).Kind() = %d, want %d", tt.in, v.Kind(), tt.kind)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	var zero Value
	if zero.IsValid() {
		t.Error("zero Value should not be valid")
	}
	if !MakeInt(0).IsValid() {
		t.Error("MakeInt(0) should be valid")
	}
}

// ---------------------------------------------------------------------------
// Arithmetic
// ---------------------------------------------------------------------------

func TestAddInts(t *testing.T) {
	a, b := MakeInt(10), MakeInt(32)
	r := a.Add(b)
	if r.Int() != 42 {
		t.Errorf("10 + 32 = %d", r.Int())
	}
}

func TestAddStrings(t *testing.T) {
	a, b := MakeString("foo"), MakeString("bar")
	r := a.Add(b)
	if r.String() != "foobar" {
		t.Errorf(`"foo" + "bar" = %q`, r.String())
	}
}

func TestSubMulDivMod(t *testing.T) {
	a, b := MakeInt(10), MakeInt(3)
	if a.Sub(b).Int() != 7 {
		t.Error("Sub failed")
	}
	if a.Mul(b).Int() != 30 {
		t.Error("Mul failed")
	}
	if a.Div(b).Int() != 3 {
		t.Error("Div failed")
	}
	if a.Mod(b).Int() != 1 {
		t.Error("Mod failed")
	}
}

func TestNeg(t *testing.T) {
	v := MakeInt(5)
	if v.Neg().Int() != -5 {
		t.Errorf("Neg(5) = %d", v.Neg().Int())
	}
}

func TestFloatArithmetic(t *testing.T) {
	a, b := MakeFloat(2.5), MakeFloat(1.5)
	if a.Add(b).Float() != 4.0 {
		t.Error("float add")
	}
	if a.Sub(b).Float() != 1.0 {
		t.Error("float sub")
	}
	if a.Mul(b).Float() != 3.75 {
		t.Error("float mul")
	}
	if a.Div(b).Float() != 2.5/1.5 {
		t.Error("float div")
	}
}

func TestUintArithmetic(t *testing.T) {
	a, b := MakeUint(10), MakeUint(3)
	if a.Add(b).Uint() != 13 {
		t.Error("uint add")
	}
	if a.Mul(b).Uint() != 30 {
		t.Error("uint mul")
	}
}

// ---------------------------------------------------------------------------
// Comparison & Equality
// ---------------------------------------------------------------------------

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b  Value
		equal bool
	}{
		{MakeNil(), MakeNil(), true},
		{MakeInt(1), MakeInt(1), true},
		{MakeInt(1), MakeInt(2), false},
		{MakeString("x"), MakeString("x"), true},
		{MakeString("x"), MakeString("y"), false},
		{MakeBool(true), MakeBool(true), true},
		{MakeBool(true), MakeBool(false), false},
		{MakeFloat(1.0), MakeFloat(1.0), true},
	}
	for i, tt := range tests {
		if got := tt.a.Equal(tt.b); got != tt.equal {
			t.Errorf("case %d: %v.Equal(%v) = %v, want %v", i, tt.a, tt.b, got, tt.equal)
		}
	}
}

func TestCmp(t *testing.T) {
	if MakeInt(1).Cmp(MakeInt(2)) >= 0 {
		t.Error("1 should be < 2")
	}
	if MakeInt(2).Cmp(MakeInt(1)) <= 0 {
		t.Error("2 should be > 1")
	}
	if MakeInt(5).Cmp(MakeInt(5)) != 0 {
		t.Error("5 should == 5")
	}
	if MakeString("a").Cmp(MakeString("b")) >= 0 {
		t.Error(`"a" should be < "b"`)
	}
}

// ---------------------------------------------------------------------------
// Bitwise
// ---------------------------------------------------------------------------

func TestBitwiseOps(t *testing.T) {
	a, b := MakeInt(0xFF), MakeInt(0x0F)
	if a.And(b).Int() != 0x0F {
		t.Error("And")
	}
	if a.Or(MakeInt(0x100)).Int() != 0x1FF {
		t.Error("Or")
	}
	if a.Xor(b).Int() != 0xF0 {
		t.Error("Xor")
	}
	if MakeInt(1).Lsh(4).Int() != 16 {
		t.Error("Lsh")
	}
	if MakeInt(16).Rsh(4).Int() != 1 {
		t.Error("Rsh")
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

func TestConversions(t *testing.T) {
	// Int -> other types
	i := MakeInt(42)
	if i.ToFloat().Float() != 42.0 {
		t.Error("ToFloat")
	}
	if i.ToUint().Uint() != 42 {
		t.Error("ToUint")
	}
	if i.ToBool().Bool() != true {
		t.Error("ToBool non-zero")
	}
	if MakeInt(0).ToBool().Bool() {
		t.Error("ToBool zero")
	}
	if i.ToString().String() != "42" {
		t.Errorf("ToString = %q", i.ToString().String())
	}

	// Float -> Int
	f := MakeFloat(3.9)
	if f.ToInt().Int() != 3 {
		t.Errorf("Float.ToInt = %d", f.ToInt().Int())
	}

	// String -> ToString (identity)
	s := MakeString("hello")
	if s.ToString().String() != "hello" {
		t.Error("String.ToString")
	}
}

// ---------------------------------------------------------------------------
// Interface round-trip
// ---------------------------------------------------------------------------

func TestInterface(t *testing.T) {
	tests := []struct {
		v    Value
		want any
	}{
		{MakeNil(), nil},
		{MakeBool(true), true},
		{MakeInt(7), int(7)},
		{MakeInt64(7), int64(7)},
		{MakeInt8(7), int8(7)},
		{MakeInt16(7), int16(7)},
		{MakeInt32(7), int32(7)},
		{MakeUint(8), uint(8)},
		{MakeUint8(8), uint8(8)},
		{MakeUint16(8), uint16(8)},
		{MakeUint32(8), uint32(8)},
		{MakeUint64(8), uint64(8)},
		{MakeFloat(1.5), 1.5},
		{MakeFloat32(1.5), float32(1.5)},
		{MakeString("s"), "s"},
	}
	for _, tt := range tests {
		got := tt.v.Interface()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Interface() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestComplexArithmetic(t *testing.T) {
	a := MakeComplex(1, 2)
	b := MakeComplex(3, 4)
	sum := a.Add(b)
	c := sum.Complex()
	if real(c) != 4 || imag(c) != 6 {
		t.Errorf("complex add: %v", c)
	}
}

func TestFloatNaN(t *testing.T) {
	nan := MakeFloat(math.NaN())
	if nan.Equal(nan) {
		t.Error("NaN should not equal itself")
	}
}

func TestSizePreservation(t *testing.T) {
	// Arithmetic should preserve the size from the left operand
	a := MakeInt8(10)
	b := MakeInt8(3)
	sum := a.Add(b)
	if got := sum.Interface(); got != int8(13) {
		t.Errorf("int8(10) + int8(3) = %v (%T), want int8(13)", got, got)
	}

	a32 := MakeInt32(100)
	b32 := MakeInt32(50)
	diff := a32.Sub(b32)
	if got := diff.Interface(); got != int32(50) {
		t.Errorf("int32(100) - int32(50) = %v (%T), want int32(50)", got, got)
	}

	f32 := MakeFloat32(2.5)
	g32 := MakeFloat32(1.0)
	prod := f32.Mul(g32)
	if got := prod.Interface(); got != float32(2.5) {
		t.Errorf("float32(2.5) * float32(1.0) = %v (%T), want float32(2.5)", got, got)
	}

	u16 := MakeUint16(100)
	v16 := MakeUint16(10)
	quotient := u16.Div(v16)
	if got := quotient.Interface(); got != uint16(10) {
		t.Errorf("uint16(100) / uint16(10) = %v (%T), want uint16(10)", got, got)
	}
}

func TestFromInterfaceRoundTrip(t *testing.T) {
	// FromInterface should preserve the exact type through Interface()
	tests := []struct {
		in   any
		want any
	}{
		{int(42), int(42)},
		{int8(42), int8(42)},
		{int16(42), int16(42)},
		{int32(42), int32(42)},
		{int64(42), int64(42)},
		{uint(42), uint(42)},
		{uint8(42), uint8(42)},
		{uint16(42), uint16(42)},
		{uint32(42), uint32(42)},
		{uint64(42), uint64(42)},
		{float32(3.14), float32(3.14)},
		{float64(3.14), float64(3.14)},
	}
	for _, tt := range tests {
		v := FromInterface(tt.in)
		got := v.Interface()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("FromInterface(%v (%T)).Interface() = %v (%T), want %v (%T)",
				tt.in, tt.in, got, got, tt.want, tt.want)
		}
	}
}
