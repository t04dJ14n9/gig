package divergence_hunt51

import "fmt"

// ============================================================================
// Round 51: Type conversion edge cases - numeric, string, rune, bool
// ============================================================================

func IntToInt8() int8 {
	var x int = 200
	return int8(x) // wraps to -56
}

func IntToInt16() int16 {
	var x int = 40000
	return int16(x) // wraps
}

func IntToUint() uint {
	var x int = -1
	return uint(x) // max uint
}

func UintToInt() int {
	var x uint = 42
	return int(x)
}

func Float32ToInt() int {
	var f float32 = 3.7
	return int(f)
}

func Float64ToInt() int {
	var f float64 = -2.9
	return int(f) // truncates toward zero
}

func IntToFloat32() float32 {
	var x int = 42
	return float32(x)
}

func IntToFloat64() float64 {
	var x int = 42
	return float64(x)
}

func RuneToString() string {
	return string('A')
}

func IntRuneToString() string {
	return string(rune(65))
}

func BoolToInt() int {
	b := true
	if b { return 1 }
	return 0
}

func ByteToString() string {
	return string(byte(72))
}

func BytesToString() string {
	b := []byte{72, 101, 108, 108, 111}
	return string(b)
}

func StringToBytes() int {
	s := "Hello"
	b := []byte(s)
	return len(b)
}

func RunesToString() string {
	r := []rune{72, 101, 108, 108, 111}
	return string(r)
}

func StringToRunes() int {
	s := "Hello"
	r := []rune(s)
	return len(r)
}

func Int64ToInt32() int32 {
	var x int64 = 2147483648 // overflow int32
	return int32(x)
}

func Uint32ToInt32() int32 {
	var x uint32 = 2147483648
	return int32(x)
}

func FmtConversion() string {
	return fmt.Sprintf("%d %f %s %t", 42, 3.14, "hi", true)
}
