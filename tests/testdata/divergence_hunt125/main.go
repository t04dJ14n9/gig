package divergence_hunt125

import "fmt"

// ============================================================================
// Round 125: String and rune operations
// ============================================================================

func StringLenBytes() string {
	s := "Hello"
	return fmt.Sprintf("len=%d", len(s))
}

func StringLenRunes() string {
	s := "Hello, 世界"
	return fmt.Sprintf("bytes=%d-runes=%d", len(s), len([]rune(s)))
}

func StringRuneAt() string {
	s := "Hello, 世界"
	r := []rune(s)
	return fmt.Sprintf("rune7=%c", r[7])
}

func StringRangeRunes() string {
	s := "Go语言"
	count := 0
	for range s {
		count++
	}
	return fmt.Sprintf("bytes=%d", count)
}

func StringConcat() string {
	a := "hello"
	b := " "
	c := "world"
	return a + b + c
}

func StringCompare() string {
	a := "apple"
	b := "banana"
	if a < b {
		return "less"
	}
	return "greater"
}

func StringSliceBytes() string {
	s := "Hello, World"
	sub := s[7:12]
	return sub
}

func StringByteConversion() string {
	s := "ABC"
	b := []byte(s)
	return fmt.Sprintf("%d-%d-%d", b[0], b[1], b[2])
}

func StringRuneConversion() string {
	r := '世'
	return fmt.Sprintf("codepoint=%U", r)
}

func StringEmptyCheck() string {
	s := ""
	if s == "" {
		return "empty"
	}
	return "not-empty"
}
