package divergence_hunt182

import (
	"fmt"
	"strings"
)

// ============================================================================
// Round 182: strings.Builder advanced usage
// ============================================================================

func BuilderBasic() string {
	var b strings.Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return fmt.Sprintf("%s", b.String())
}

func BuilderWriteByte() string {
	var b strings.Builder
	b.WriteByte('H')
	b.WriteByte('i')
	return fmt.Sprintf("%s", b.String())
}

func BuilderWriteRune() string {
	var b strings.Builder
	b.WriteRune('你')
	b.WriteRune('好')
	return fmt.Sprintf("%s", b.String())
}

func BuilderGrow() string {
	var b strings.Builder
	b.Grow(100)
	b.WriteString("preallocated")
	return fmt.Sprintf("%s:%d", b.String(), b.Len())
}

func BuilderReset() string {
	var b strings.Builder
	b.WriteString("first")
	b.Reset()
	b.WriteString("second")
	return fmt.Sprintf("%s", b.String())
}

func BuilderMultipleWrites() string {
	var b strings.Builder
	for i := 0; i < 5; i++ {
		b.WriteString("x")
	}
	return fmt.Sprintf("%s:%d", b.String(), b.Len())
}

func BuilderMixedWrites() string {
	var b strings.Builder
	b.WriteString("num:")
	b.WriteByte('1')
	b.WriteString("")
	b.WriteByte('2')
	return fmt.Sprintf("%s", b.String())
}

func BuilderEmpty() string {
	var b strings.Builder
	return fmt.Sprintf("%s:%d", b.String(), b.Len())
}

func BuilderNested() string {
	outer := func() string {
		var b strings.Builder
		b.WriteString("outer[")
		inner := func() string {
			var b strings.Builder
			b.WriteString("inner")
			return b.String()
		}
		b.WriteString(inner())
		b.WriteString("]")
		return b.String()
	}
	return fmt.Sprintf("%s", outer())
}

func BuilderLargeString() string {
	var b strings.Builder
	for i := 0; i < 10; i++ {
		b.WriteString("abc")
	}
	return fmt.Sprintf("%d:%s", b.Len(), b.String()[:9])
}
