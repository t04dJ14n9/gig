package divergence_hunt67

import "strings"

// ============================================================================
// Round 67: String/rune edge cases - multi-byte, emoji, byte conversions
// ============================================================================

func StringLenBytes() int {
	s := "Hello"
	return len(s)
}

func StringLenMultiByte() int {
	s := "你好"
	return len(s) // 6 bytes (3 per Chinese char in UTF-8)
}

func RuneCount() int {
	s := "你好世界"
	count := 0
	for _ = range s {
		count++
	}
	return count
}

func StringIndexByte() byte {
	s := "Hello"
	return s[1] // 'e'
}

func StringSlice() string {
	s := "Hello, World"
	return s[7:12]
}

func StringConcatEmpty() string {
	s := "Hello" + ""
	return s
}

func StringConcatMulti() string {
	s := "Hello" + " " + "World"
	return s
}

func StringCompare() int {
	a := "apple"
	b := "banana"
	if a < b {
		return -1
	}
	return 1
}

func StringEqual() bool {
	return "hello" == "hello"
}

func StringEmptyCompare() bool {
	return "" == ""
}

func RuneValue() int {
	r := 'A'
	return int(r)
}

func RuneChineseValue() int {
	r := '中'
	return int(r)
}

func StringFromBytes() string {
	b := []byte{72, 101, 108, 108, 111}
	return string(b)
}

func StringToBytes() int {
	s := "Hello"
	b := []byte(s)
	return len(b)
}

func StringFromRunes() string {
	r := []rune{'H', 'e', 'l', 'l', 'o'}
	return string(r)
}

func StringRangeIndex() string {
	s := "abc"
	result := ""
	for i := range s {
		result += string(rune('0' + i))
	}
	return result
}

func StringsRepeat() string {
	return strings.Repeat("ab", 3)
}

func StringsTrimCutset() string {
	return strings.Trim("!!hello!!", "!")
}

func StringContainsEmpty() bool {
	return strings.Contains("hello", "")
}
