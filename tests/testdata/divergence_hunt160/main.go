package divergence_hunt160

import "fmt"

// ============================================================================
// Round 160: Type conversion edge cases
// ============================================================================

// MyInt is a defined type based on int
type MyInt int

// MyString is a defined type based on string
type MyString string

// TypeAliasConversion tests conversion between defined types and underlying types
func TypeAliasConversion() string {
	var i int = 42
	mi := MyInt(i)
	back := int(mi)
	return fmt.Sprintf("myint=%d-back=%d", mi, back)
}

// TypeAliasString tests string type alias conversion
func TypeAliasString() string {
	s := "hello"
	ms := MyString(s)
	back := string(ms)
	return fmt.Sprintf("len=%d-back=%s", len(ms), back)
}

// ByteSliceToString tests []byte to string conversion
func ByteSliceToString() string {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	s := string(b)
	return s
}

// StringToByteSlice tests string to []byte conversion
func StringToByteSlice() string {
	s := "hello"
	b := []byte(s)
	return fmt.Sprintf("len=%d-first=%c", len(b), b[0])
}

// RuneSliceToString tests []rune to string conversion
func RuneSliceToString() string {
	r := []rune{'h', 'e', 'l', 'l', 'o'}
	s := string(r)
	return s
}

// IntToFloatConversion tests int to float conversion
func IntToFloatConversion() string {
	i := 42
	f := float64(i)
	return fmt.Sprintf("float=%.1f", f)
}

// FloatToIntConversion tests float to int conversion (truncates)
func FloatToIntConversion() string {
	f := 3.99
	i := int(f)
	return fmt.Sprintf("int=%d", i)
}

// UintToIntConversion tests uint to int conversion
func UintToIntConversion() string {
	var u uint = 42
	i := int(u)
	return fmt.Sprintf("int=%d", i)
}

// IntToUintConversion tests int to uint conversion
func IntToUintConversion() string {
	i := 42
	u := uint(i)
	return fmt.Sprintf("uint=%d", u)
}

// LargeIntToUint8 tests large int to small uint conversion (overflow)
func LargeIntToUint8() string {
	i := 300
	u := uint8(i)
	return fmt.Sprintf("uint8=%d", u) // Will wrap: 300 % 256 = 44
}

// NegativeIntToUint tests negative int to uint conversion
func NegativeIntToUint() string {
	i := -1
	u := uint(i)
	return fmt.Sprintf("uint=%d", u) // Will be max uint
}

// FloatSpecialToInt tests special float to int conversion
func FloatSpecialToInt() string {
	// Float too large for int results in implementation-defined value
	// Use a regular float instead
	f := 1e20
	i := int(f)
	return fmt.Sprintf("float=%.0e-int=%d", f, i)
}

// ByteToIntConversion tests byte to int conversion
func ByteToIntConversion() string {
	var b byte = 255
	i := int(b)
	return fmt.Sprintf("int=%d", i)
}

// Int8ToInt16Conversion tests signed int extension
func Int8ToInt16Conversion() string {
	var i8 int8 = -5
	i16 := int16(i8)
	return fmt.Sprintf("i16=%d", i16)
}

// ComplexConversion tests complex number conversions
func ComplexConversion() string {
	f := 3.0
	c := complex(f, 4.0)
	return fmt.Sprintf("real=%.0f-imag=%.0f", real(c), imag(c))
}

// Complex64ToComplex128 tests complex64 to complex128
func Complex64ToComplex128() string {
	var c64 complex64 = complex(3, 4)
	c128 := complex128(c64)
	return fmt.Sprintf("c128=%v", c128)
}
