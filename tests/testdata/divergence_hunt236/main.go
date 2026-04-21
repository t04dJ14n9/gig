package divergence_hunt236

import "fmt"

// ============================================================================
// Round 236: String comparison and sorting
// ============================================================================

func StringEquality() string {
	s1 := "hello"
	s2 := "hello"
	s3 := "world"
	return fmt.Sprintf("%v:%v", s1 == s2, s1 == s3)
}

func StringInequality() string {
	s1 := "abc"
	s2 := "def"
	return fmt.Sprintf("%v:%v", s1 != s2, s1 != "abc")
}

func StringLessThan() string {
	s1 := "apple"
	s2 := "banana"
	return fmt.Sprintf("%v:%v", s1 < s2, s2 < s1)
}

func StringGreaterThan() string {
	s1 := "zebra"
	s2 := "apple"
	return fmt.Sprintf("%v:%v", s1 > s2, s2 > s1)
}

func StringLessOrEqual() string {
	s1 := "same"
	s2 := "same"
	s3 := "different"
	return fmt.Sprintf("%v:%v", s1 <= s2, s3 <= s1)
}

func StringCaseSensitivity() string {
	s1 := "Hello"
	s2 := "hello"
	return fmt.Sprintf("%v:%v", s1 == s2, s1 < s2)
}

func EmptyStringComparison() string {
	s := "test"
	return fmt.Sprintf("%v:%v:%v", s == "", "" == "", s != "")
}

func StringLengthComparison() string {
	s1 := "hi"
	s2 := "hello"
	return fmt.Sprintf("%v:%v", len(s1) < len(s2), len(s1) == len(s2))
}

func StringLexicographicOrder() string {
	words := []string{"cherry", "apple", "banana"}
	sorted := false
	for i := 0; i < len(words)-1; i++ {
		if words[i] > words[i+1] {
			sorted = false
			break
		}
	}
	return fmt.Sprintf("sorted=%v", sorted)
}

func UnicodeStringComparison() string {
	s1 := ""
	s2 := ""
	return fmt.Sprintf("%v:%v", s1 < s2, s1 == s2)
}

func StringPrefixComparison() string {
	s := "hello world"
	prefix := "hello"
	return fmt.Sprintf("%v", len(s) >= len(prefix) && s[:len(prefix)] == prefix)
}
