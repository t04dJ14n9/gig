package divergence_hunt36

import "strconv"

// ============================================================================
// Round 36: String/rune edge cases - multi-byte, indexing, slicing, conversion
// ============================================================================

func StringByteLen() int {
	s := "Hello, 世界"
	return len(s) // 12 bytes (7 ASCII + 3 for '世' + 3 for '界'... wait, Chinese chars are 3 bytes each in UTF-8)
}

func StringRuneLen() int {
	s := "Hello, 世界"
	return len([]rune(s))
}

func StringByteIndex() byte {
	s := "Hello"
	return s[1] // 'e'
}

func StringSliceMultiByte() string {
	s := "Hello, 世界"
	return s[:7] // "Hello, "
}

func RuneFromInt() string {
	return string(rune(65)) // "A"
}

func StringFromBytes() string {
	b := []byte{72, 101, 108, 108, 111}
	return string(b)
}

func BytesFromString() int {
	s := "Hello"
	b := []byte(s)
	return len(b)
}

func RuneSliceFromString() int {
	s := "Hello, 世界"
	r := []rune(s)
	return len(r)
}

func StringFromRuneSlice() string {
	r := []rune{'H', 'i', '世', '界'}
	return string(r)
}

func StrconvAtoiNegative() int {
	n, _ := strconv.Atoi("-42")
	return n
}

func StrconvItoaNegative() string {
	return strconv.Itoa(-42)
}

func StrconvFormatUint() string {
	return strconv.FormatUint(42, 16)
}

func StrconvFormatIntBase() string {
	return strconv.FormatInt(255, 16)
}

func StringRangeRuneIndex() int {
	s := "Hello, 世界"
	result := 0
	for i, r := range s {
		if r == '世' {
			result = i
			break
		}
	}
	return result
}

func StringCompareOps() int {
	a := "apple"
	b := "banana"
	r := 0
	if a < b { r |= 1 }
	if a <= b { r |= 2 }
	if b > a { r |= 4 }
	if b >= a { r |= 8 }
	return r
}

func StringConcatMulti() string {
	a := "Hello"
	b := " "
	c := "World"
	return a + b + c
}

func StringEmptyLen() int {
	return len("")
}

func StringMultiByteIndex() int {
	s := "世界"
	b := []byte(s)
	return len(b) // 6 (3 bytes per Chinese char)
}

func RuneValue() int {
	var r rune = 'A'
	return int(r)
}
