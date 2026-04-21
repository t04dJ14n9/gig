package divergence_hunt232

import "fmt"

// ============================================================================
// Round 232: Rune iteration
// ============================================================================

func RangeOverStringASCII() string {
	s := "ABC"
	result := ""
	for _, r := range s {
		result += fmt.Sprintf("%c", r)
	}
	return result
}

func RangeOverStringUnicode() string {
	s := "Hello, 世界"
	count := 0
	for i, r := range s {
		count += i + int(r)
	}
	return fmt.Sprintf("%d", count)
}

func RangeStringIndexValues() string {
	s := "Go"
	result := ""
	for i, r := range s {
		result += fmt.Sprintf("%d:%c ", i, r)
	}
	return result
}

func CountRunesInString() string {
	s := "Hello, 世界"
	count := 0
	for range s {
		count++
	}
	return fmt.Sprintf("bytes=%d,runes=%d", len(s), count)
}

func RuneSliceFromString() string {
	s := "test"
	runes := []rune(s)
	return fmt.Sprintf("%d:%c:%c", len(runes), runes[0], runes[len(runes)-1])
}

func StringFromRuneSlice() string {
	runes := []rune{'H', 'e', 'l', 'l', 'o'}
	return string(runes)
}

func IterateEmptyString() string {
	count := 0
	for range "" {
		count++
	}
	return fmt.Sprintf("%d", count)
}

func UnicodeByteVsRuneCount() string {
	s := "日本語"
	return fmt.Sprintf("bytes=%d,runes=%d", len(s), len([]rune(s)))
}

func RangeWithEmoji() string {
	s := "Go: emoji"
	count := 0
	for i := range s {
		count += i
	}
	return fmt.Sprintf("%d", count)
}

func RuneComparison() string {
	r1, r2 := 'A', 'B'
	return fmt.Sprintf("%v:%v:%v", r1 < r2, r1 == r2, r1 > r2)
}

func RuneArithmetic() string {
	r := 'A'
	return fmt.Sprintf("%c:%c:%c", r, r+1, r+32)
}
