package divergence_hunt26

// ============================================================================
// Round 26: Type system edge cases - conversions, constants, enums
// ============================================================================

func Int8Range() int8 { var x int8 = 127; return x }

func Int8MinRange() int8 { var x int8 = -128; return x }

func Uint8Max() uint8 { return 255 }

func Int16Range() int16 { var x int16 = 32767; return x }

func Uint16Max() uint16 { return 65535 }

func Float32Smallest() float32 { return 1.0 }

func Complex64Basic() complex64 {
	var z complex64 = 1 + 2i
	return z
}

func Complex128Basic() complex128 {
	z := complex(3.0, 4.0)
	return z
}

func RuneType() rune { return 'A' }

func ByteType() byte { return 'Z' }

func StringType() string { return "hello" }

func BoolType() bool { return true }

func IntType() int { return 42 }

func Int64Type() int64 { return 1234567890123 }

func UintType() uint { return 42 }

func Uint64Type() uint64 { return 1234567890123 }

func Float64Type() float64 { return 3.14 }

func Float32Type() float32 { return 2.71 }

func TypeConversionChain() int {
	var a int = 10
	var b int8 = int8(a)
	var c int16 = int16(b)
	var d int32 = int32(c)
	return int(d)
}

func UnsignedConversion() int {
	var x uint8 = 200
	var y int16 = int16(x)
	return int(y)
}

func SignedToUnsigned() uint {
	var x int = -1
	return uint(x)
}

func FloatToIntTrunc() int {
	x := 9.99
	return int(x)
}

func IntToFloatPrecise() float64 {
	x := 1000000
	return float64(x)
}

func StringToSlice() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

func SliceToString() string {
	b := []byte{'h', 'i'}
	return string(b)
}

func RuneToString() string { return string('A') }

func RunesToString() string { return string([]rune{'A', 'B', 'C'}) }

func StringToRunes() int { return len([]rune("Hello, 世界")) }
