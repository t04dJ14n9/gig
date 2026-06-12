package divergence_hunt58

import (
	"fmt"
	"strconv"
	"strings"
)

// ============================================================================
// Round 58: String formatting and parsing - fmt, strconv, builders
// ============================================================================

func FmtSprintfInt() string {
	return fmt.Sprintf("%d", 42)
}

func FmtSprintfFloat() string {
	return fmt.Sprintf("%.2f", 3.14159)
}

func FmtSprintfString() string {
	return fmt.Sprintf("%s", "hello")
}

func FmtSprintfBool() string {
	return fmt.Sprintf("%t", true)
}

func FmtSprintfWidth() string {
	return fmt.Sprintf("|%5d|%-5d|", 42, 42)
}

func FmtSprintfHex() string {
	return fmt.Sprintf("%x %X %#x", 255, 255, 255)
}

func FmtSprintfOctal() string {
	return fmt.Sprintf("%o %#o", 8, 8)
}

func FmtSprintfBinary() string {
	return fmt.Sprintf("%b", 10)
}

func FmtSprintfChar() string {
	return fmt.Sprintf("%c", 65)
}

func FmtSprintfPadZero() string {
	return fmt.Sprintf("%05d", 42)
}

func FmtSprintfQuoted() string {
	return fmt.Sprintf("%q", "hello")
}

func FmtSprintfDefault() string {
	return fmt.Sprintf("%v", []int{1, 2, 3})
}

func FmtErrorf() string {
	return fmt.Errorf("error %d", 42).Error()
}

func StrconvAtoiPositive() int {
	n, _ := strconv.Atoi("12345")
	return n
}

func StrconvAtoiNegative() int {
	n, _ := strconv.Atoi("-42")
	return n
}

func StrconvItoaPositive() string {
	return strconv.Itoa(42)
}

func StrconvItoaNegative() string {
	return strconv.Itoa(-42)
}

func StrconvFormatBool() string {
	return strconv.FormatBool(true)
}

func StrconvParseBool() bool {
	b, _ := strconv.ParseBool("true")
	return b
}

func StrconvFormatFloat() string {
	return strconv.FormatFloat(3.14, 'f', 2, 64)
}

func StrconvParseFloat() float64 {
	f, _ := strconv.ParseFloat("3.14", 64)
	return f
}

func StringBuilderConcat() string {
	var b strings.Builder
	for i := 0; i < 5; i++ {
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}
