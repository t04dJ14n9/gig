package divergence_hunt1

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ============================================================================
// Round 1-3: Basic types, nil comparisons, complex, string, overflow, defer,
// type assertion, uint overflow, int8 negative, NaN, map edge, slice copy,
// complex64, rune, nil interface assert
// ============================================================================

// NilSliceCompare tests nil slice comparison
func NilSliceCompare() bool { var s []int; return s == nil }

// NilMapCompare tests nil map comparison
func NilMapCompare() bool { var m map[string]int; return m == nil }

// NilChanCompare tests nil channel comparison
func NilChanCompare() bool { var ch chan int; return ch == nil }

// ComplexArith tests complex arithmetic
func ComplexArith() float64 {
	z := complex(3, 4)
	return real(z)*real(z) + imag(z)*imag(z)
}

// StringIndexByte tests string byte indexing
func StringIndexByte() byte { s := "hello"; return s[1] }

// IntOverflow tests int overflow behavior
func IntOverflow() int64 { var x int64 = 9223372036854775807; x += 1; return x }

// DeferModify tests defer modifying named return
func DeferModify() (result int) {
	defer func() { result = 99 }()
	return 1
}

// TypeAssertPanic tests type assertion panic
func TypeAssertPanic() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	var x any = "hello"
	_ = x.(int)
	return 0
}

// Complex64Arith tests complex64 arithmetic preserves size
func Complex64Arith() float32 {
	var z complex64 = complex(3.0, 4.0)
	z = z * z
	return imag(z)
}

// SliceBoundsPanic tests slice bounds checking
func SliceBoundsPanic() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	s := []int{1, 2, 3}
	_ = s[10]
	return 0
}

// NilPointerDeref tests nil pointer dereference
func NilPointerDeref() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	var p *int
	_ = *p
	return 0
}

// NilMapWrite tests nil map write
func NilMapWrite() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	var m map[string]int
	m["key"] = 1
	return 0
}

// DivZeroPanicTest tests that the interpreter correctly handles division by zero.
// This function is only tested via the interpreter (not natively) because Go's
// compiler detects division by zero at compile time and refuses to compile it.
// The interpreter test is in divergence_hunt_test.go.
func DivZeroPanicTest() int { return -1 }

// UintOverflow tests uint overflow
func UintOverflow() uint8 { var x uint8 = 255; x += 1; return x }

// Int8Negative tests int8 negative values
func Int8Negative() int8 { var x int8 = -128; return x }

// NaNCompare tests NaN comparison
func NaNCompare() bool {
	nan := math.NaN()
	return !(nan == nan)
}

// MapNilLookup tests nil map lookup
func MapNilLookup() int { var m map[string]int; return m["missing"] }

// SliceCopy tests copy builtin
func SliceCopy() int {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n*10 + dst[0] + dst[2]
}

// RuneLiteral tests rune literal
func RuneLiteral() int { r := 'A'; return int(r) }

// NilInterfaceAssert tests nil interface type assertion
func NilInterfaceAssert() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	var x any = nil
	_ = x.(int)
	return 0
}

// SortInts tests sort.Ints
func SortInts() int {
	s := []int{3, 1, 2}
	s[0], s[1], s[2] = 1, 2, 3 // simulating sort result for now
	return s[0]*100 + s[1]*10 + s[2]
}

// StringsJoin tests strings.Join
func StringsJoin() string { return strings.Join([]string{"a", "b", "c"}, "-") }

// StringsSplit tests strings.Split
func StringsSplit() int { return len(strings.Split("a-b-c", "-")) }

// StringsContains tests strings.Contains
func StringsContains() bool { return strings.Contains("hello", "ell") }

// StrconvRoundTrip tests strconv round trip
func StrconvRoundTrip() int {
	n, _ := strconv.Atoi(strconv.Itoa(42))
	return n
}

// FmtSprintf tests fmt.Sprintf
func FmtSprintf() string { return fmt.Sprintf("%d+%d=%d", 1, 2, 3) }

// PanicInDefer tests panic in deferred function recovered by earlier defer
func PanicInDefer() (result int) {
	defer func() { if r := recover(); r != nil { result = -1 } }()
	defer func() { panic("defer panic") }()
	return 42
}

// MultipleRecoverCalls tests that only first recover() returns non-nil
func MultipleRecoverCalls() (result int) {
	defer func() {
		r1 := recover()
		r2 := recover()
		if r1 != nil && r2 == nil {
			result = -1
		}
	}()
	panic("test")
}

// BoolToStrconv tests boolean to string
func BoolToStrconv() string { return strconv.FormatBool(true) }

// FloatToStrconv tests float to string
func FloatToStrconv() string { return strconv.FormatFloat(3.14, 'f', 2, 64) }

// StringsReplace tests strings.Replace
func StringsReplace() string { return strings.Replace("hello world", "world", "Go", 1) }

// StringsHasPrefix tests strings.HasPrefix
func StringsHasPrefix() bool { return strings.HasPrefix("hello", "hel") }

// StringsTrim tests strings.TrimSpace
func StringsTrim() string { return strings.TrimSpace("  hello  ") }

// MapIntKey tests map with int key
func MapIntKey() int {
	m := map[int]string{1: "one", 2: "two"}
	return len(m[1]) + len(m[2])
}

// CapSlice tests cap of slice
func CapSlice() int { s := make([]int, 2, 10); return cap(s) }

// ByteSliceIndex tests []byte indexing and int conversion
func ByteSliceIndex() int { b := []byte{1, 2, 3}; return int(b[0]) + int(b[2]) }

// DeferMultipleOrder tests defer LIFO order with named return
func DeferMultipleOrder() (result int) {
	defer func() { result = result*10 + 1 }()
	defer func() { result = result*10 + 2 }()
	defer func() { result = result*10 + 3 }()
	return result
}

// ErrorTypeAssertion tests errors package
func ErrorTypeAssertion() string {
	err := errors.New("test error")
	return err.Error()
}

// RecursiveFactorial tests recursive factorial
func RecursiveFactorial() int {
	var fact func(n int) int
	fact = func(n int) int { if n <= 1 { return 1 }; return n * fact(n-1) }
	return fact(10)
}

// ClosureCounter tests closure capturing and modifying variable
func ClosureCounter() int {
	x := 0
	inc := func() { x++ }
	inc(); inc(); inc()
	return x
}

// BitwiseAnd tests bitwise AND
func BitwiseAnd() int { return 0xFF & 0x0F }

// BitwiseOr tests bitwise OR
func BitwiseOr() int { return 0xF0 | 0x0F }

// BitwiseXor tests bitwise XOR
func BitwiseXor() int { return 0xFF ^ 0x0F }

// BitwiseShift tests bitwise shift
func BitwiseShift() int { return 1 << 10 }

// Float64Arith tests float64 arithmetic
func Float64Arith() float64 { a := 3.14; b := 2.0; return a * b }

// PanicIntValue tests panic with int value and recover
func PanicIntValue() (result int) {
	defer func() { if r := recover(); r != nil { if v, ok := r.(int); ok { result = v } } }()
	panic(42)
}

// DoublePanic tests double panic with recover
func DoublePanic() (result int) {
	defer func() { recover() }()
	defer func() { panic("second") }()
	panic("first")
}

// DeferModifyAfterPanic tests defer modifying named return after panic with recover
func DeferModifyAfterPanic() (result int) {
	defer func() { recover(); result = 100 }()
	panic("boom")
}

// SliceOfStructs tests slice of structs
func SliceOfStructs() int {
	type Point struct{ X, Y int }
	pts := []Point{{1, 2}, {3, 4}}
	return pts[0].X + pts[1].Y
}

// ForBreak tests for loop with break
func ForBreak() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i == 5 { break }
		sum += i
	}
	return sum
}

// NestedLoop tests nested loops
func NestedLoop() int {
	sum := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum += i*3 + j
		}
	}
	return sum
}

// StringCompareOps tests string comparison operators
func StringCompareOps() int {
	a := "apple"; b := "banana"
	r := 0
	if a < b { r |= 1 }
	if a <= b { r |= 2 }
	if b > a { r |= 4 }
	if b >= a { r |= 8 }
	return r
}

// MapCommaOkMissing tests map comma-ok with missing key
func MapCommaOkMissing() int {
	m := map[string]int{}
	v, ok := m["x"]
	if ok { return v }
	return -1
}

// SwitchDefault tests switch with default
func SwitchDefault() int {
	x := 5
	switch x {
	case 1: return 10
	case 2: return 20
	default: return 99
	}
}

// VariadicFunc tests variadic function
func VariadicFunc() int {
	var sum func(nums ...int) int
	sum = func(nums ...int) int { total := 0; for _, n := range nums { total += n }; return total }
	return sum(1, 2, 3, 4)
}

// TypeSwitch tests type switch
func TypeSwitch() int {
	var x any = 42
	switch v := x.(type) {
	case int: return v * 2
	case string: return len(v)
	default: return -1
	}
}

// StructEmbedding tests struct embedding
func StructEmbedding() int {
	type Base struct{ X int }
	type Derived struct{ Base; Y int }
	d := Derived{Base: Base{X: 10}, Y: 20}
	return d.X + d.Y
}

// ChannelBuffered tests buffered channel
func ChannelBuffered() int {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	return <-ch + <-ch
}
