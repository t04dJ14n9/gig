package external

import "fmt"
import "strconv"
import "strings"

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
