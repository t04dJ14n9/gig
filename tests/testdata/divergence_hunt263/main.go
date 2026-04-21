package divergence_hunt263

import (
	"fmt"
)

// ============================================================================
// Round 263: Type conversions and assertions — edge cases
// ============================================================================

type MyInt int
type MyFloat float64
type MyString string

// TypeConversionNumeric tests converting between numeric types
func TypeConversionNumeric() string {
	var x MyInt = 42
	y := int(x)
	z := float64(y)
	w := MyFloat(z)
	return fmt.Sprintf("x=%d,y=%d,z=%.0f,w=%.0f", x, y, z, w)
}

// TypeConversionString tests converting between string types
func TypeConversionString() string {
	var s MyString = "hello"
	t := string(s)
	return fmt.Sprintf("s=%s,t=%s", s, t)
}

// SliceTypeConversion tests converting between slice types (same underlying)
func SliceTypeConversion() string {
	a := []int{1, 2, 3}
	b := []MyInt{}
	for _, v := range a {
		b = append(b, MyInt(v))
	}
	return fmt.Sprintf("b=%v", b)
}

// InterfaceTypeAssertion tests type assertion on interface
func InterfaceTypeAssertion() string {
	var i interface{} = 42
	v, ok := i.(int)
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// InterfaceTypeAssertionFail tests failed type assertion with comma-ok
func InterfaceTypeAssertionFail() string {
	var i interface{} = "hello"
	v, ok := i.(int)
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// InterfaceTypeSwitch tests type switch
func InterfaceTypeSwitch() string {
	var val interface{} = 3.14
	switch v := val.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case float64:
		return fmt.Sprintf("float64:%v", v)
	case string:
		return fmt.Sprintf("string:%s", v)
	default:
		return "unknown"
	}
}

// NestedTypeAssertion tests assertion through multiple levels
func NestedTypeAssertion() string {
	var i interface{} = MyInt(100)
	v, ok := i.(MyInt)
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// ConversionBetweenSameUnderlying tests converting between types with same underlying
func ConversionBetweenSameUnderlying() string {
	a := MyInt(5)
	b := int(a)
	c := MyInt(b)
	return fmt.Sprintf("a=%d,b=%d,c=%d", a, b, c)
}

// FloatToIntTruncation tests float to int truncation
func FloatToIntTruncation() string {
	f := 3.9
	i := int(f)
	return fmt.Sprintf("f=%.1f,i=%d", f, i)
}

// UintToIntConversion tests unsigned to signed conversion
func UintToIntConversion() string {
	var u uint = 42
	i := int(u)
	return fmt.Sprintf("u=%d,i=%d", u, i)
}

// ByteSliceToString tests []byte to string conversion
func ByteSliceToString() string {
	b := []byte{72, 101, 108, 108, 111}
	s := string(b)
	return fmt.Sprintf("s=%s", s)
}
