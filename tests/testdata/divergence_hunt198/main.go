package divergence_hunt198

import (
	"fmt"
)

// ============================================================================
// Round 198: Unsafe-like operations via standard means
// Note: Since unsafe package may not be supported, we test similar concepts
// using only standard, safe Go operations
// ============================================================================

// TypeAliasConversion tests type alias behavior
func TypeAliasConversion() string {
	type MyInt int
	var a MyInt = 42
	var b int = int(a)
	return fmt.Sprintf("%d:%d", a, b)
}

// TypeAliasComparison tests type alias comparison
func TypeAliasComparison() string {
	type MyInt int
	var a MyInt = 42
	var b MyInt = 42
	return fmt.Sprintf("%v", a == b)
}

// UnderlyingType tests underlying type access
func UnderlyingType() string {
	type MyString string
	s := MyString("hello")
	underlying := string(s)
	return fmt.Sprintf("%s", underlying)
}

// ByteSliceToStringConversion tests byte slice to string conversion
func ByteSliceToStringConversion() string {
	bytes := []byte{104, 101, 108, 108, 111}
	str := string(bytes)
	return fmt.Sprintf("%s", str)
}

// StringToByteSliceConversion tests string to byte slice conversion
func StringToByteSliceConversion() string {
	str := "hello"
	bytes := []byte(str)
	return fmt.Sprintf("%d:%d:%d", len(bytes), bytes[0], bytes[4])
}

// RuneSliceConversion tests rune slice conversion
func RuneSliceConversion() string {
	runes := []rune{'h', 'e', 'l', 'l', 'o'}
	str := string(runes)
	return fmt.Sprintf("%s", str)
}

// StringToRuneSlice tests string to rune slice
func StringToRuneSlice() string {
	str := "hello"
	runes := []rune(str)
	return fmt.Sprintf("%d:%d", len(runes), runes[0])
}

// BinaryDataRepresentation tests binary data handling
func BinaryDataRepresentation() string {
	data := []byte{0x00, 0xFF, 0xAB, 0xCD}
	return fmt.Sprintf("%d:%d:%d:%d", data[0], data[1], data[2], data[3])
}

// StructLayoutExploration tests struct field access patterns
func StructLayoutExploration() string {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alice", Age: 30}
	return fmt.Sprintf("%s:%d", p.Name, p.Age)
}

// SliceHeaderConcept tests slice header concept
func SliceHeaderConcept() string {
	original := []int{1, 2, 3, 4, 5}
	slice := original[1:3]
	// Modify through slice affects original
	slice[0] = 100
	return fmt.Sprintf("%d:%d", len(slice), cap(slice))
}

// StringImmutabilityPattern tests string immutability
func StringImmutabilityPattern() string {
	s := "hello"
	original := s
	s = s + " world"
	return fmt.Sprintf("%s:%s", original, s)
}
