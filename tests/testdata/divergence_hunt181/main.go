package divergence_hunt181

import (
	"fmt"
	"strconv"
)

// ============================================================================
// Round 181: strconv package conversions
// ============================================================================

func AtoiBasic() string {
	n, err := strconv.Atoi("42")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%d", n)
}

func AtoiNegative() string {
	n, err := strconv.Atoi("-99")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%d", n)
}

func AtoiInvalid() string {
	_, err := strconv.Atoi("abc")
	if err != nil {
		return fmt.Sprintf("has_error")
	}
	return fmt.Sprintf("no_error")
}

func ParseIntBase10() string {
	n, err := strconv.ParseInt("255", 10, 64)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%d", n)
}

func ParseIntBase16() string {
	n, err := strconv.ParseInt("FF", 16, 64)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%d", n)
}

func ParseIntBase2() string {
	n, err := strconv.ParseInt("1010", 2, 64)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%d", n)
}

func FormatIntBase10() string {
	s := strconv.FormatInt(255, 10)
	return fmt.Sprintf("%s", s)
}

func FormatIntBase16() string {
	s := strconv.FormatInt(255, 16)
	return fmt.Sprintf("%s", s)
}

func FormatIntBase2() string {
	s := strconv.FormatInt(10, 2)
	return fmt.Sprintf("%s", s)
}

func ParseBool() string {
	b1, _ := strconv.ParseBool("true")
	b2, _ := strconv.ParseBool("1")
	b3, _ := strconv.ParseBool("false")
	b4, _ := strconv.ParseBool("0")
	return fmt.Sprintf("%v:%v:%v:%v", b1, b2, b3, b4)
}

func ItoaBasic() string {
	s := strconv.Itoa(123)
	return fmt.Sprintf("%s", s)
}

func ItoaNegative() string {
	s := strconv.Itoa(-456)
	return fmt.Sprintf("%s", s)
}
