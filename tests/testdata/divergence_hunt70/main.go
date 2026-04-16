package divergence_hunt70

// ============================================================================
// Round 70: Type alias and named type edge cases
// ============================================================================

type MyInt int

func (m MyInt) Double() MyInt {
	return m * 2
}

func NamedTypeMethod() int {
	var x MyInt = 5
	return int(x.Double())
}

func NamedTypeConversion() int {
	x := 42
	m := MyInt(x)
	return int(m)
}

func NamedTypeArith() int {
	var a MyInt = 10
	var b MyInt = 20
	return int(a + b)
}

type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func TypeAliasConversion() float64 {
	c := Celsius(100)
	f := CToF(c)
	return float64(f)
}

type StringAlias string

func NamedStringType() string {
	var s StringAlias = "hello"
	return string(s)
}

type BoolAlias bool

func NamedBoolType() bool {
	var b BoolAlias = true
	return bool(b)
}

type SliceAlias []int

func NamedSliceType() int {
	var s SliceAlias = SliceAlias{1, 2, 3}
	return len(s)
}

type MapAlias map[string]int

func NamedMapType() int {
	var m MapAlias = MapAlias{"a": 1}
	return m["a"]
}

type FuncAlias func(int) int

func NamedFuncType() int {
	var f FuncAlias = func(x int) int { return x * 2 }
	return f(5)
}

type IntPtr *int

func NamedPointerType() int {
	v := 42
	var p IntPtr = IntPtr(&v)
	return *p
}

func NamedTypeCompare() bool {
	var a MyInt = 5
	var b MyInt = 5
	return a == b
}

func NamedTypeLessThan() bool {
	var a MyInt = 3
	var b MyInt = 5
	return a < b
}
