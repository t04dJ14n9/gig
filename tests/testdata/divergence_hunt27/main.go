package divergence_hunt27

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Round 27: String manipulation, sorting, data processing
// ============================================================================

func StringSort() string {
	s := []string{"banana", "apple", "cherry"}
	sort.Strings(s)
	return strings.Join(s, ",")
}

func StringUnique() int {
	s := "hello"
	seen := map[rune]bool{}
	for _, c := range s { seen[c] = true }
	return len(seen)
}

func StringIsDigit() int {
	s := "a1b2c3"
	digits := 0
	for _, c := range s {
		if c >= '0' && c <= '9' { digits++ }
	}
	return digits
}

func StringIsAlpha() int {
	s := "a1b2c3"
	alpha := 0
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') { alpha++ }
	}
	return alpha
}

func StringToUpperLower() string {
	s := "Hello World"
	return strings.ToUpper(s) + strings.ToLower(s)
}

func StringCapitalize() string {
	s := "hello world"
	return strings.ToUpper(s[:1]) + s[1:]
}

func StringCountWords() int {
	s := "the quick brown fox"
	return len(strings.Fields(s))
}

func StringReverseWords() string {
	words := strings.Fields("hello world foo")
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return strings.Join(words, " ")
}

func FmtInteger() string {
	return fmt.Sprintf("%d", 42)
}

func FmtHexInt() string {
	return fmt.Sprintf("%x", 255)
}

func FmtOctalInt() string {
	return fmt.Sprintf("%o", 8)
}

func FmtBinaryInt() string {
	return fmt.Sprintf("%b", 10)
}

func FmtCharFromInt() string {
	return fmt.Sprintf("%c", 65)
}

func FmtUnicode() string {
	return fmt.Sprintf("%U", 'A')
}

func SortIntSliceDesc() int {
	s := []int{5, 3, 1, 4, 2}
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

func SortFloatSliceDesc() float64 {
	s := []float64{3.14, 1.41, 2.71}
	sort.Float64s(s)
	return s[0] + s[1] + s[2]
}

func StringJoinWithSep() string {
	parts := []string{"a", "b", "c"}
	return strings.Join(parts, "-")
}

func StringSplitN() int {
	parts := strings.SplitN("a-b-c-d", "-", 2)
	return len(parts)
}

func StringRepeatN() string {
	return strings.Repeat("ab", 3)
}

func StringMapFunc() string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' { return r - 32 }
		return r
	}, "hello")
}
