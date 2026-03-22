package external

import (
	"fmt"
	"strconv"
	"strings"
)

// FmtSprintf tests fmt.Sprintf
func FmtSprintf() string {
	return fmt.Sprintf("hello %d", 42)
}

// FmtSprintfMulti tests fmt.Sprintf with multiple args
func FmtSprintfMulti() string {
	return fmt.Sprintf("%s is %d years old", "Alice", 30)
}

// StringsToUpper tests strings.ToUpper
func StringsToUpper() string {
	return strings.ToUpper("hello world")
}

// StringsToLower tests strings.ToLower
func StringsToLower() string {
	return strings.ToLower("HELLO")
}

// StringsContains tests strings.Contains
func StringsContains() int {
	if strings.Contains("hello world", "world") {
		return 1
	}
	return 0
}

// StringsReplace tests strings.ReplaceAll
func StringsReplace() string {
	return strings.ReplaceAll("foo bar foo", "foo", "baz")
}

// StringsHasPrefix tests strings.HasPrefix
func StringsHasPrefix() int {
	if strings.HasPrefix("hello world", "hello") {
		return 1
	}
	return 0
}

// StrconvItoa tests strconv.Itoa
func StrconvItoa() string {
	return strconv.Itoa(42)
}

// StrconvAtoi tests strconv.Atoi
func StrconvAtoi() int {
	n, _ := strconv.Atoi("123")
	return n
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// FmtSprintfInt returns fmt.Sprintf("value: %d", n)
func FmtSprintfInt(n int) string { return fmt.Sprintf("value: %d", n) }

// StringsToUpperStr returns strings.ToUpper(s)
func StringsToUpperStr(s string) string { return strings.ToUpper(s) }

// StringsToLowerStr returns strings.ToLower(s)
func StringsToLowerStr(s string) string { return strings.ToLower(s) }

// StringsContainsStr returns true if s contains substr
func StringsContainsStr(s, substr string) bool { return strings.Contains(s, substr) }

// StrconvItoaN returns strconv.Itoa(n)
func StrconvItoaN(n int) string { return strconv.Itoa(n) }

// StrconvAtoiStr returns the integer parsed from s
func StrconvAtoiStr(s string) (int, error) { return strconv.Atoi(s) }
