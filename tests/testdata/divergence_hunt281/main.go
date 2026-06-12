package divergence_hunt281

import (
	"fmt"
)

// ============================================================================
// Round 281: String conversion edge cases — []byte↔string, rune↔string, byte↔string

// ByteSliceToStringAndBack tests []byte→string→[]byte round trip
func ByteSliceToStringAndBack() string {
	original := []byte{72, 101, 108, 108, 111}
	s := string(original)
	b := []byte(s)
	return fmt.Sprintf("s=%s,b=%v", s, b)
}

// StringToRuneSlice tests converting string to []rune
func StringToRuneSlice() string {
	s := "Hello, 世界"
	runes := []rune(s)
	return fmt.Sprintf("len=%d,r0=%d,r7=%d", len(runes), runes[0], runes[7])
}

// RuneSliceToString tests converting []rune to string
func RuneSliceToString() string {
	runes := []rune{72, 105, 0x4E16, 0x754C}
	s := string(runes)
	return fmt.Sprintf("s=%s,len=%d", s, len(s))
}

// ByteToIntConversion tests byte to int conversion
func ByteToIntConversion() string {
	var b byte = 255
	i := int(b)
	return fmt.Sprintf("i=%d", i)
}

// IntToByteTruncation tests int to byte truncation
func IntToByteTruncation() string {
	i := 300
	b := byte(i)
	return fmt.Sprintf("b=%d", b)
}

// StringIndexOutOfRangePanics tests string index out of range
func StringIndexOutOfRangePanics() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	s := "hi"
	_ = s[10]
	return "no_panic"
}

// EmptyStringOperations tests operations on empty string
func EmptyStringOperations() string {
	s := ""
	return fmt.Sprintf("len=%d,empty=%t,quote=%q", len(s), s == "", s)
}

// StringComparisonWithOperators tests string comparison
func StringComparisonWithOperators() string {
	a := "abc"
	b := "abd"
	c := "abc"
	return fmt.Sprintf("a<b=%t,a==c=%t,a<=c=%t,b>a=%t", a < b, a == c, a <= c, b > a)
}

// StringFromIntConversion tests string(int) — NOT allowed in Go, must use strconv
func StringFromIntConversion() string {
	// This tests that we use fmt.Sprintf or strconv, not string(42) which gives '*'
	i := 42
	s := fmt.Sprintf("%d", i)
	return fmt.Sprintf("s=%s", s)
}

// StringConcatInLoop tests string concatenation performance (correctness)
func StringConcatInLoop() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += "x"
	}
	return fmt.Sprintf("len=%d", len(s))
}
