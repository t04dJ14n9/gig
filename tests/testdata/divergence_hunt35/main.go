package divergence_hunt35

// ============================================================================
// Round 35: Const expressions, iota, type aliases, named types
// ============================================================================

func IotaShift() int {
	const (
		KB = 1 << (10 * iota)
		MB
		GB
	)
	return int(MB / KB)
}

func IotaBitmask() int {
	const (
		Read    = 1 << iota // 1
		Write               // 2
		Execute             // 4
	)
	return Read + Write + Execute // 7
}

func ConstExpression() int {
	const (
		A = 10
		B = A * 2
		C = B + A
	)
	return C
}

func TypeAliasBasic() int {
	type MyInt int
	var x MyInt = 42
	return int(x)
}

func TypeAliasString() int {
	type MyStr string
	var s MyStr = "hello"
	return len(s)
}

func TypeAliasArith() int {
	type MyInt int
	var a MyInt = 10
	var b MyInt = 20
	return int(a + b)
}

func TypeAliasComparison() bool {
	type MyInt int
	var a MyInt = 42
	var b MyInt = 42
	return a == b
}

func NestedTypeAlias() int {
	type Inner int
	type Outer Inner
	var x Outer = 100
	return int(x)
}

func ConstBlockBlank() int {
	const (
		A = 1
		_
		C = 3
	)
	return A + C
}

func ConstWithString() string {
	const (
		Prefix = "pre"
		Suffix = "suf"
	)
	return Prefix + Suffix
}

func TypeAliasSlice() int {
	type IntSlice []int
	s := IntSlice{1, 2, 3}
	return len(s)
}

func TypeAliasMap() int {
	type StringMap map[string]int
	m := StringMap{"a": 1, "b": 2}
	return m["a"] + m["b"]
}

func ConstExpressionFloat() float64 {
	const (
		Pi  = 3.14159
		Tau = Pi * 2
	)
	return Tau
}

func TypeAliasConversion() int {
	type Celsius float64
	type Fahrenheit float64
	var c Celsius = 100
	f := Fahrenheit(c*9/5 + 32)
	return int(f)
}

func ConstArithComplex() int {
	const (
		A = 1 + 2
		B = A * 3
		C = B - A
	)
	return C
}

func IotaSkip() int {
	const (
		A = iota + 1 // 1
		B             // 2
		_             // skip 3
		D             // 4
	)
	return A + B + D
}

func ConstBitwiseOps() int {
	const (
		Mask = 0xFF00
		Val  = 0x1234
		Result = Val & Mask
	)
	return Result
}
