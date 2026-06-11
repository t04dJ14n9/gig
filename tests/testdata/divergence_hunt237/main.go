package divergence_hunt237

import "fmt"

// ============================================================================
// Round 237: Unicode string handling
// ============================================================================

func UnicodeStringLength() string {
	s := "Hello, 世界"
	return fmt.Sprintf("bytes=%d,runes=%d", len(s), len([]rune(s)))
}

func UnicodeRuneAccess() string {
	s := "日本語"
	runes := []rune(s)
	return fmt.Sprintf("%c:%c:%c", runes[0], runes[1], runes[2])
}

func UnicodeStringIndexing() string {
	s := "Hello, 世界"
	return fmt.Sprintf("%d:%d", s[0], s[7])
}

func UnicodeRangeLoop() string {
	s := "Go"
	result := ""
	for _, r := range s {
		result += fmt.Sprintf("%c", r)
	}
	return result
}

func UnicodeConcatenation() string {
	s1 := "Hello"
	s2 := "世界"
	s3 := ""
	return s1 + " " + s2 + " " + s3
}

func MixedASCIIUnicode() string {
	s := "Price: $100 (USD)"
	return fmt.Sprintf("len=%d", len(s))
}

func UnicodeComparison() string {
	s1 := ""
	s2 := ""
	return fmt.Sprintf("%v:%v", s1 == s2, s1 < s2)
}

func EmojiStringHandling() string {
	s := "Go rocks!"
	return fmt.Sprintf("bytes=%d,runes=%d", len(s), len([]rune(s)))
}

func UnicodeCaseConversion() string {
	s := "Hello"
	result := ""
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			r = r - 'a' + 'A'
		}
		result += string(r)
	}
	return result
}

func UnicodeByteAccess() string {
	s := "中"
	return fmt.Sprintf("bytes=%v", []byte(s))
}

func StringWithCombiningChars() string {
	s := "caf"
	return fmt.Sprintf("len=%d,runes=%d", len(s), len([]rune(s)))
}
