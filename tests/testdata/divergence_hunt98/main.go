package divergence_hunt98

import (
	"fmt"
	"unicode/utf8"
)

// ============================================================================
// Round 98: String manipulation - runes, indices, multi-byte
// ============================================================================

func RuneCount() string {
	s := "Hello, 世界"
	return fmt.Sprintf("%d", utf8.RuneCountInString(s))
}

func ByteLen() string {
	s := "Hello, 世界"
	return fmt.Sprintf("%d", len(s))
}

func RuneAt() string {
	s := "Hello, 世界"
	r := []rune(s)
	return fmt.Sprintf("%c", r[7])
}

func StringFromRunes() string {
	runes := []rune{'H', 'e', 'l', 'l', 'o'}
	return string(runes)
}

func StringSliceByte() string {
	s := "Hello, World"
	return s[7:12]
}

func StringConcat() string {
	s := "Hello" + " " + "World"
	return s
}

func StringRangeRunes() string {
	s := "Go语言"
	count := 0
	for range s {
		count++
	}
	return fmt.Sprintf("%d", count)
}

func RuneSliceModify() string {
	s := "Hello"
	runes := []rune(s)
	runes[0] = 'J'
	return string(runes)
}

func MultiByteIndex() string {
	s := "abc你好"
	// Byte index 3 is the start of first Chinese char
	return fmt.Sprintf("%d:%d", len(s), utf8.RuneCountInString(s))
}

func StringCompare() string {
	a := "apple"
	b := "banana"
	if a < b {
		return "a < b"
	}
	return "a >= b"
}

func StringPrefixSuffix() string {
	s := "Hello, World"
	return fmt.Sprintf("%v:%v", len(s) > 5, len(s) > 20)
}

func EmptyString() string {
	var s string
	return fmt.Sprintf("%v:%d", s == "", len(s))
}
