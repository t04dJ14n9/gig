package divergence_hunt265

import (
	"fmt"
)

// ============================================================================
// Round 265: String manipulation — runes, indexing, slicing, concatenation
// ============================================================================

// StringRuneCount tests len vs rune count for unicode
func StringRuneCount() string {
	s := "Hello, 世界"
	return fmt.Sprintf("byte_len=%d", len(s))
}

// StringIndexByte tests indexing a byte from string
func StringIndexByte() string {
	s := "ABC"
	return fmt.Sprintf("byte0=%d", s[0])
}

// StringSliceRange tests slicing a string
func StringSliceRange() string {
	s := "Hello, World"
	return fmt.Sprintf("sub=%s", s[7:12])
}

// StringConcat tests string concatenation
func StringConcat() string {
	a := "Hello"
	b := " "
	c := "World"
	return a + b + c
}

// StringConcatWithInt tests fmt.Sprint for concatenation
func StringConcatWithInt() string {
	s := "value="
	return fmt.Sprintf("%s%d", s, 42)
}

// StringEmptyVsNil tests empty string literal
func StringEmptyVsNil() string {
	var s string
	return fmt.Sprintf("empty=%t,len=%d", s == "", len(s))
}

// StringRangeRunes tests ranging over string produces runes
func StringRangeRunes() string {
	s := "Go"
	result := ""
	for i, r := range s {
		result += fmt.Sprintf("%d:%d ", i, r)
	}
	return result[:len(result)-1]
}

// StringMultiByte tests string with multi-byte characters
func StringMultiByte() string {
	s := "日本語"
	return fmt.Sprintf("byte_len=%d", len(s))
}

// StringCompare tests string comparison operators
func StringCompare() string {
	a := "apple"
	b := "banana"
	return fmt.Sprintf("lt=%t,eq=%t,gt=%t", a < b, a == b, a > b)
}

// StringRepeatConcat tests repeated string concatenation in loop
func StringRepeatConcat() string {
	s := ""
	for i := 0; i < 5; i++ {
		s += "x"
	}
	return fmt.Sprintf("s=%s,len=%d", s, len(s))
}
