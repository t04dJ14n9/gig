package strings_pkg

// Concat tests string concatenation
func Concat() string {
	s := "hello"
	return s + " world"
}

// ConcatLoop tests string concatenation in loop
func ConcatLoop() string {
	s := ""
	for i := 0; i < 3; i++ {
		s = s + "ab"
	}
	return s
}

// Len tests string length
func Len() int {
	return len("hello")
}

// Index tests string indexing
func Index() string {
	s := "abcde"
	return string(s[0]) + string(s[4])
}

// Comparison tests string comparison
func Comparison() int {
	a := "abc"
	b := "abd"
	if a < b {
		return 1
	}
	return 0
}

// Equality tests string equality
func Equality() int {
	a := "hello"
	b := "hello"
	c := "world"
	result := 0
	if a == b {
		result = result + 1
	}
	if a != c {
		result = result + 10
	}
	return result
}

// EmptyCheck tests empty string check
func EmptyCheck() int {
	s := ""
	if len(s) == 0 {
		return 1
	}
	return 0
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// StrConcat concatenates two strings
func StrConcat(a, b string) string { return a + b }

// StrLen returns the length of s
func StrLen(s string) int { return len(s) }

// StrCompare returns -1 if a<b, 0 if a==b, 1 if a>b
func StrCompare(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// StrEqual returns true if a == b
func StrEqual(a, b string) bool { return a == b }
