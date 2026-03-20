package thirdparty

import "strconv"

// StrconvParseBool tests strconv.ParseBool.
func StrconvParseBool() int {
	b, _ := strconv.ParseBool("true")
	if b {
		return 1
	}
	return 0
}

// StrconvFormatBool tests strconv.FormatBool.
func StrconvFormatBool() string {
	return strconv.FormatBool(false)
}

// StrconvParseInt tests strconv.ParseInt.
func StrconvParseInt() int64 {
	n, _ := strconv.ParseInt("12345", 10, 64)
	return n
}

// StrconvParseUint tests strconv.ParseUint.
func StrconvParseUint() uint64 {
	n, _ := strconv.ParseUint("12345", 10, 64)
	return n
}

// StrconvFormatInt tests strconv.FormatInt.
func StrconvFormatInt() string {
	return strconv.FormatInt(-12345, 10)
}

// StrconvFormatUint tests strconv.FormatUint.
func StrconvFormatUint() string {
	return strconv.FormatUint(12345, 10)
}

// StrconvParseFloat tests strconv.ParseFloat.
func StrconvParseFloat() float64 {
	f, _ := strconv.ParseFloat("123.45", 64)
	return f
}

// StrconvFormatFloat tests strconv.FormatFloat.
func StrconvFormatFloat() string {
	return strconv.FormatFloat(123.45, 'f', 2, 64)
}

// StrconvQuote tests strconv.Quote.
func StrconvQuote() string {
	return strconv.Quote("hello\nworld")
}

// StrconvQuoteToASCII tests strconv.QuoteToASCII.
func StrconvQuoteToASCII() string {
	return strconv.QuoteToASCII("hello")
}

// StrconvUnquote tests strconv.Unquote.
func StrconvUnquote() string {
	s, _ := strconv.Unquote(`"hello"`)
	return s
}

// StrconvAppendInt tests strconv.AppendInt.
func StrconvAppendInt() string {
	b := make([]byte, 0, 20)
	b = strconv.AppendInt(b, 12345, 10)
	return string(b)
}

// StrconvAppendFloat tests strconv.AppendFloat.
func StrconvAppendFloat() string {
	b := make([]byte, 0, 20)
	b = strconv.AppendFloat(b, 123.45, 'f', 2, 64)
	return string(b)
}
