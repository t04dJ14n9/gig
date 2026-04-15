package divergence_hunt18

import (
	"fmt"
	"strconv"
	"strings"
)

// ============================================================================
// Round 18: String processing, parsing, formatting, conversion
// ============================================================================

// StringToIntConversion tests string to int conversion
func StringToIntConversion() int {
	n, _ := strconv.Atoi("12345")
	return n
}

// IntToStringConversion tests int to string conversion
func IntToStringConversion() string {
	return strconv.Itoa(42)
}

// FloatToStringConversion tests float to string conversion
func FloatToStringConversion() string {
	return strconv.FormatFloat(3.14, 'f', 2, 64)
}

// StringToFloatConversion tests string to float conversion
func StringToFloatConversion() float64 {
	f, _ := strconv.ParseFloat("3.14", 64)
	return f
}

// BoolToStringConversion tests bool to string conversion
func BoolToStringConversion() string {
	return strconv.FormatBool(true)
}

// StringToBoolConversion tests string to bool conversion
func StringToBoolConversion() bool {
	b, _ := strconv.ParseBool("true")
	return b
}

// StringSplitJoin tests split and join round trip
func StringSplitJoin() string {
	parts := strings.Split("a,b,c", ",")
	return strings.Join(parts, "-")
}

// StringTrimSpace tests trimming whitespace
func StringTrimSpace() string {
	return strings.TrimSpace("  hello  world  ")
}

// StringTrimPrefix tests trimming prefix
func StringTrimPrefix() string {
	return strings.TrimPrefix("hello world", "hello ")
}

// StringTrimSuffix tests trimming suffix
func StringTrimSuffix() string {
	return strings.TrimSuffix("hello world", " world")
}

// StringReplaceAll tests replace all occurrences
func StringReplaceAll() string {
	return strings.ReplaceAll("hello hello", "hello", "hi")
}

// StringBuilderPattern tests string builder
func StringBuilderPattern() string {
	var b strings.Builder
	for i := 0; i < 5; i++ {
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}

// StringRuneConversion tests string to rune conversion
func StringRuneConversion() int {
	s := "Hello, 世界"
	runes := []rune(s)
	return len(runes)
}

// RuneToStringConversion tests rune to string conversion
func RuneToStringConversion() string {
	return string([]rune{'H', 'i'})
}

// StringByteConversion tests string to byte conversion
func StringByteConversion() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

// ByteToStringConversion tests byte to string conversion
func ByteToStringConversion() string {
	b := []byte("world")
	return string(b)
}

// FmtSprintfComplex tests complex sprintf
func FmtSprintfComplex() string {
	return fmt.Sprintf("Name: %s, Age: %d, Score: %.1f", "Alice", 30, 95.5)
}

// FmtSprintfPadding tests sprintf with padding
func FmtSprintfPadding() string {
	return fmt.Sprintf("[%5d][%-5d][%05d]", 42, 42, 42)
}

// StringPadLeft tests left padding
func StringPadLeft() string {
	s := "42"
	for len(s) < 5 { s = "0" + s }
	return s
}

// StringPadRight tests right padding
func StringPadRight() string {
	s := "42"
	for len(s) < 5 { s = s + " " }
	return s
}

// CamelCaseSplit tests splitting camelCase
func CamelCaseSplit() int {
	s := "helloWorldFoo"
	count := 1
	for _, c := range s {
		if c >= 'A' && c <= 'Z' { count++ }
	}
	return count
}

// StringReverseWords tests reversing words in string
func StringReverseWords() string {
	words := strings.Split("hello world foo", " ")
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}
