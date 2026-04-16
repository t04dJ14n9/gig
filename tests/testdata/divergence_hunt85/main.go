package divergence_hunt85

import "strings"

// ============================================================================
// Round 85: String builder edge cases
// ============================================================================

func BuilderBasic() string {
	var b strings.Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return b.String()
}

func BuilderLen() int {
	var b strings.Builder
	b.WriteString("hello")
	return b.Len()
}

func BuilderGrow() int {
	var b strings.Builder
	b.Grow(100)
	return b.Len() // Len is still 0
}

func BuilderReset() string {
	var b strings.Builder
	b.WriteString("hello")
	b.Reset()
	b.WriteString("world")
	return b.String()
}

func BuilderWriteByte() string {
	var b strings.Builder
	b.WriteByte('A')
	b.WriteByte('B')
	b.WriteByte('C')
	return b.String()
}

func BuilderWriteString() string {
	var b strings.Builder
	for i := 0; i < 5; i++ {
		b.WriteString("x")
	}
	return b.String()
}

func BuilderCap() int {
	var b strings.Builder
	b.WriteString("hello")
	return b.Cap()
}

func BuilderEmpty() string {
	var b strings.Builder
	return b.String()
}

func BuilderLarge() string {
	var b strings.Builder
	for i := 0; i < 100; i++ {
		b.WriteString("a")
	}
	return b.String()[:5]
}

func StringConcatMany() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += "x"
	}
	return s
}

func StringJoin() string {
	parts := []string{"hello", "world", "foo"}
	return strings.Join(parts, ", ")
}

func StringRepeat() string {
	return strings.Repeat("ab", 5)
}

func StringReplace() string {
	return strings.Replace("hello world hello", "hello", "hi", 1)
}

func StringReplaceAll() string {
	return strings.ReplaceAll("hello world hello", "hello", "hi")
}

func StringContains() bool {
	return strings.Contains("hello world", "world")
}
