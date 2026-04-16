package divergence_hunt78

// ============================================================================
// Round 78: Complex type conversions - truncation, sign extension, etc.
// ============================================================================

func IntToInt8() int8 {
	var x int = 300
	return int8(x) // truncation
}

func IntToUint() uint {
	var x int = -1
	return uint(x)
}

func UintToInt() int {
	var x uint = 2147483648 // 2^31
	return int(x)
}

func Float64ToInt() int {
	x := 3.7
	return int(x)
}

func Float64ToUint() uint {
	x := 3.7
	return uint(x)
}

func IntToFloat64() float64 {
	return float64(42)
}

func Int8ToInt16() int16 {
	var x int8 = -1
	return int16(x) // sign extension
}

func Uint8ToUint16() uint16 {
	var x uint8 = 255
	return uint16(x) // zero extension
}

func Int32ToInt8() int8 {
	var x int32 = 128
	return int8(x) // overflow
}

func Float32ToFloat64() float64 {
	var x float32 = 3.14
	return float64(x)
}

func Float64ToFloat32() float32 {
	return float32(3.14159265358979323846)
}

func ByteToInt() int {
	var b byte = 200
	return int(b)
}

func RuneToInt() int {
	var r rune = '世'
	return int(r)
}

func SliceToInterface() int {
	s := []int{1, 2, 3}
	var x any = s
	return len(x.([]int))
}

func StringToByteSlice() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

func ByteSliceToString() string {
	b := []byte{72, 101, 108, 108, 111}
	return string(b)
}

func IntSliceToFloatSlice() int {
	input := []int{1, 2, 3}
	output := make([]float64, len(input))
	for i, v := range input {
		output[i] = float64(v)
	}
	return len(output)
}
