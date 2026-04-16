package divergence_hunt48

import (
	"fmt"
	"strconv"
	"strings"
)

// ============================================================================
// Round 48: String manipulation edge cases - trim, split, join, replace
// ============================================================================

func TrimSpace() string {
	return strings.TrimSpace("  hello world  ")
}

func TrimPrefix() string {
	return strings.TrimPrefix("HelloWorld", "Hello")
}

func TrimSuffix() string {
	return strings.TrimSuffix("HelloWorld", "World")
}

func SplitN() int {
	parts := strings.SplitN("a,b,c,d", ",", 2)
	return len(parts)
}

func SplitAfter() int {
	parts := strings.SplitAfter("a,b,c", ",")
	return len(parts)
}

func ReplaceN() string {
	return strings.Replace("hello hello hello", "hello", "hi", 2)
}

func ReplaceAll() string {
	return strings.ReplaceAll("hello hello hello", "hello", "hi")
}

func Repeat() string {
	return strings.Repeat("ab", 3)
}

func Contains() bool {
	return strings.Contains("hello world", "world")
}

func ContainsAny() bool {
	return strings.ContainsAny("hello", "aeiou")
}

func HasPrefix() bool {
	return strings.HasPrefix("hello world", "hello")
}

func HasSuffix() bool {
	return strings.HasSuffix("hello world", "world")
}

func IndexFunc() int {
	return strings.IndexFunc("hello123", func(r rune) bool { return r >= '0' && r <= '9' })
}

func TitleCase() string {
	return strings.Title("hello world")
}

func ToTitle() string {
	return strings.ToTitle("hello world")
}

func MapFunc() string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' { return r - 32 }
		return r
	}, "hello")
}

func BuilderString() string {
	var b strings.Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return b.String()
}

func BuilderLen() int {
	var b strings.Builder
	b.WriteString("hello")
	return b.Len()
}

func StrconvQuote() string {
	return strconv.Quote(`hello "world"`)
}

func FmtStringOps() string {
	return fmt.Sprintf("%s|%q|%5s", "hi", "hi", "hi")
}
