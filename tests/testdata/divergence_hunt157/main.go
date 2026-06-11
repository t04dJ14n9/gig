package divergence_hunt157

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ============================================================================
// Round 157: String advanced operations and utf8
// ============================================================================

// Utf8RuneCount tests utf8.RuneCountInString
func Utf8RuneCount() string {
	s := "Hello, 世界"
	return fmt.Sprintf("bytes=%d-runecount=%d", len(s), utf8.RuneCountInString(s))
}

// Utf8DecodeRune tests utf8.DecodeRuneInString
func Utf8DecodeRune() string {
	s := "世"
	r, size := utf8.DecodeRuneInString(s)
	return fmt.Sprintf("rune=%c-size=%d", r, size)
}

// Utf8Valid tests utf8.ValidString
func Utf8Valid() string {
	valid := "Hello"
	return fmt.Sprintf("valid=%t", utf8.ValidString(valid))
}

// Utf8RuneLen tests utf8.RuneLen
func Utf8RuneLen() string {
	r1 := 'a'   // ASCII
	r2 := '世' // 3 bytes in UTF-8
	return fmt.Sprintf("ascii=%d-cjk=%d", utf8.RuneLen(r1), utf8.RuneLen(r2))
}

// StringBuilderBasic tests strings.Builder basic usage
func StringBuilderBasic() string {
	var b strings.Builder
	b.WriteString("Hello")
	b.WriteString(" ")
	b.WriteString("World")
	return b.String()
}

// StringBuilderGrow tests strings.Builder with Grow
func StringBuilderGrow() string {
	var b strings.Builder
	b.Grow(100)
	b.WriteString("test")
	return fmt.Sprintf("len=%d", b.Len())
}

// StringBuilderByte tests strings.Builder WriteByte
func StringBuilderByte() string {
	var b strings.Builder
	b.WriteByte('H')
	b.WriteByte('i')
	return b.String()
}

// StringBuilderRune tests strings.Builder WriteRune
func StringBuilderRune() string {
	var b strings.Builder
	b.WriteRune('世')
	b.WriteRune('界')
	return b.String()
}

// StringCompare tests strings.Compare
func StringCompare() string {
	return fmt.Sprintf("ab-cd=%d-cd-ab=%d-ab-ab=%d",
		strings.Compare("ab", "cd"),
		strings.Compare("cd", "ab"),
		strings.Compare("ab", "ab"))
}

// StringEqualFold tests strings.EqualFold (case-insensitive)
func StringEqualFold() string {
	return fmt.Sprintf("Hello-hello=%t-Go-go=%t",
		strings.EqualFold("Hello", "hello"),
		strings.EqualFold("Go", "go"))
}

// StringIndexAny tests strings.IndexAny
func StringIndexAny() string {
	return fmt.Sprintf("idx=%d", strings.IndexAny("hello", "aeiou"))
}

// StringLastIndex tests strings.LastIndex
func StringLastIndex() string {
	return fmt.Sprintf("idx=%d", strings.LastIndex("hello hello", "hello"))
}

// StringCutPrefix tests strings.CutPrefix (Go 1.20+)
func StringCutPrefix() string {
	s := "hello world"
	after, found := strings.CutPrefix(s, "hello ")
	return fmt.Sprintf("found=%t-after=%s", found, after)
}

// StringCutSuffix tests strings.CutSuffix
func StringCutSuffix() string {
	s := "hello world"
	before, found := strings.CutSuffix(s, " world")
	return fmt.Sprintf("found=%t-before=%s", found, before)
}

// StringClone tests strings.Clone
func StringClone() string {
	original := "hello"
	cloned := strings.Clone(original)
	return fmt.Sprintf("equal=%t", original == cloned)
}

// StringCount tests strings.Count
func StringCount() string {
	return fmt.Sprintf("count=%d", strings.Count("banana", "a"))
}
