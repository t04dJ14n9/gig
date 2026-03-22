package typeconv

import "strconv"

// IntToFloat64 tests int to float64 conversion
func IntToFloat64() int {
	x := 42
	f := float64(x)
	return int(f)
}

// Float64Arithmetic tests float64 arithmetic
func Float64Arithmetic() int {
	a := 10
	b := 3
	return a / b
}

// StringToByteConversion tests string to byte conversion
func StringToByteConversion() string {
	s := "hello"
	b := string(s[0])
	return b
}

// IntStringConversion tests int to string conversion
func IntStringConversion() string {
	n := 12345
	return strconv.Itoa(n)
}

// StringIntConversion tests string to int conversion
func StringIntConversion() int {
	n, _ := strconv.Atoi("54321")
	return n
}

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// IntToString converts int to string using strconv.Itoa
func IntToString(n int) string { return strconv.Itoa(n) }

// StringToInt converts string to int using strconv.Atoi
func StringToInt(s string) (int, error) { return strconv.Atoi(s) }

// IntToFloatToInt converts int to float64 and back
func IntToFloatToInt(x int) int { return int(float64(x)) }
