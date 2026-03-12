package cornercases_src

// ============================================================================
// Corner Case Tests - Source Functions for Native Go Reference
// ============================================================================

// ------------------------------------------------------------------------
// Zero Value Tests
// ------------------------------------------------------------------------

func ZeroValue_Int() int {
	var x int
	return x
}

func ZeroValue_Int64() int64 {
	var x int64
	return x
}

func ZeroValue_Float64() float64 {
	var x float64
	return x
}

func ZeroValue_String() string {
	var s string
	return s
}

func ZeroValue_Bool() bool {
	var b bool
	return b
}

func ZeroValue_Slice() int {
	var s []int
	if s == nil {
		return 1
	}
	return 0
}

func ZeroValue_Map() int {
	var m map[string]int
	if m == nil {
		return 1
	}
	return 0
}

// ------------------------------------------------------------------------
// Integer Boundary Tests
// ------------------------------------------------------------------------

func IntBoundary_MaxInt32() int32 {
	return 2147483647
}

func IntBoundary_MinInt32() int32 {
	return -2147483648
}

func IntBoundary_MaxInt64() int64 {
	return 9223372036854775807
}

func IntBoundary_MinInt64() int64 {
	return -9223372036854775808
}

func IntBoundary_MaxUint32() uint32 {
	return 4294967295
}

func IntBoundary_NearMaxInt() int {
	return 2147483646
}

func IntBoundary_NearMinInt() int {
	return -2147483647
}

// ------------------------------------------------------------------------
// Integer Overflow Tests
// ------------------------------------------------------------------------

func Overflow_Int32Add() int32 {
	var x int32 = 2147483647
	return x + 1
}

func Overflow_Int32Sub() int32 {
	var x int32 = -2147483648
	return x - 1
}

func Overflow_Int32Mul() int32 {
	var x int32 = 65536
	return x * 65536
}

// ------------------------------------------------------------------------
// Float Boundary Tests
// ------------------------------------------------------------------------

func FloatBoundary_SmallPositive() float64 {
	return 1e-300
}

func FloatBoundary_SmallNegative() float64 {
	return -1e-300
}

func FloatBoundary_LargePositive() float64 {
	return 1e300
}

func FloatBoundary_LargeNegative() float64 {
	return -1e300
}

// ------------------------------------------------------------------------
// Empty Collection Tests
// ------------------------------------------------------------------------

func EmptySlice_Len() int {
	s := []int{}
	return len(s)
}

func EmptySlice_Cap() int {
	s := []int{}
	return cap(s)
}

func EmptySlice_Make() int {
	s := make([]int, 0)
	return len(s)
}

func EmptyMap_Len() int {
	m := map[string]int{}
	return len(m)
}

func EmptyMap_Make() int {
	m := make(map[string]int)
	return len(m)
}

func EmptyString_Len() int {
	s := ""
	return len(s)
}

// ------------------------------------------------------------------------
// Slice Operations Corner Cases
// ------------------------------------------------------------------------

func Slice_ZeroToZero() int {
	s := []int{1, 2, 3}
	sub := s[0:0]
	return len(sub)
}

func Slice_EndToEnd() int {
	s := []int{1, 2, 3}
	sub := s[3:3]
	return len(sub)
}

func Slice_FullSlice() int {
	s := []int{1, 2, 3}
	sub := s[:]
	return len(sub)
}

func Slice_NilSlice() int {
	var s []int
	if s == nil {
		return 1
	}
	return 0
}

func Slice_AppendToNil() int {
	var s []int
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return len(s)
}

func Slice_AppendEmpty() int {
	s := []int{1, 2}
	_ = append(s)
	return len(s)
}

// ------------------------------------------------------------------------
// Map Operations Corner Cases
// ------------------------------------------------------------------------

func Map_NilMap() int {
	var m map[string]int
	if m == nil {
		return 1
	}
	return 0
}

func Map_AccessMissingKey() int {
	m := map[string]int{"a": 1}
	return m["b"]
}

func Map_DeleteMissingKey() int {
	m := map[string]int{"a": 1}
	delete(m, "b")
	return len(m)
}

func Map_OverwriteKey() int {
	m := map[string]int{"a": 1}
	m["a"] = 2
	return m["a"]
}

func Map_NilKeyString() int {
	m := map[string]int{"": 42}
	return m[""]
}

func Map_ZeroIntKey() int {
	m := map[int]string{0: "zero"}
	return len(m[0])
}

// ------------------------------------------------------------------------
// String Corner Cases
// ------------------------------------------------------------------------

func String_Empty() int {
	s := ""
	return len(s)
}

func String_SingleChar() int {
	s := "a"
	return len(s)
}

func String_SingleByteIndex() uint8 {
	s := "abc"
	return s[1]
}

func String_LastByte() uint8 {
	s := "hello"
	return s[len(s)-1]
}

func String_Whitespace() int {
	s := " \t\n"
	return len(s)
}

func String_UnicodeMultibyte() int {
	s := "你好"
	return len(s)
}

// ------------------------------------------------------------------------
// Boolean Corner Cases
// ------------------------------------------------------------------------

func Bool_True() bool {
	return true
}

func Bool_False() bool {
	return false
}

func Bool_NotTrue() bool {
	return !true
}

func Bool_NotFalse() bool {
	return !false
}

func Bool_DoubleNegation() bool {
	return !!true
}

// ------------------------------------------------------------------------
// Arithmetic Corner Cases
// ------------------------------------------------------------------------

func Arith_DivByOne() int {
	return 100 / 1
}

func Arith_ModByOne() int {
	return 100 % 1
}

func Arith_MulByZero() int {
	return 100 * 0
}

func Arith_MulByOne() int {
	return 100 * 1
}

func Arith_AddZero() int {
	return 100 + 0
}

func Arith_SubZero() int {
	return 100 - 0
}

func Arith_NegNeg() int {
	return -(-100)
}

func Arith_NegAddNeg() int {
	return -10 + (-20)
}

// ------------------------------------------------------------------------
// Comparison Corner Cases
// ------------------------------------------------------------------------

func Compare_IntEqual() bool {
	return 5 == 5
}

func Compare_IntNotEqual() bool {
	return 5 != 6
}

func Compare_IntLess() bool {
	return 5 < 6
}

func Compare_IntLessEqual() bool {
	return 5 <= 5
}

func Compare_IntGreater() bool {
	return 6 > 5
}

func Compare_IntGreaterEqual() bool {
	return 5 >= 5
}

func Compare_StringEqual() bool {
	return "hello" == "hello"
}

func Compare_StringNotEqual() bool {
	return "hello" != "world"
}

func Compare_EmptyStringEqual() bool {
	return "" == ""
}

// ------------------------------------------------------------------------
// Logical Operation Corner Cases
// ------------------------------------------------------------------------

func Logic_TrueAndTrue() bool {
	return true && true
}

func Logic_TrueAndFalse() bool {
	return true && false
}

func Logic_FalseAndTrue() bool {
	return false && true
}

func Logic_TrueOrFalse() bool {
	return true || false
}

func Logic_FalseOrTrue() bool {
	return false || true
}

func Logic_FalseOrFalse() bool {
	return false || false
}

// ------------------------------------------------------------------------
// Control Flow Corner Cases
// ------------------------------------------------------------------------

func Control_IfNoElse() int {
	x := 0
	if true {
		x = 1
	}
	return x
}

func Control_IfFalseNoElse() int {
	x := 0
	if false {
		x = 1
	}
	return x
}

func Control_ForZeroIter() int {
	count := 0
	for i := 0; i < 0; i++ {
		count++
	}
	return count
}

func Control_ForOneIter() int {
	count := 0
	for i := 0; i < 1; i++ {
		count++
	}
	return count
}

func Control_ForBreakFirst() int {
	count := 0
	for i := 0; i < 10; i++ {
		break
		count++
	}
	return count
}

func Control_ForContinueAll() int {
	count := 0
	for i := 0; i < 5; i++ {
		continue
		count++
	}
	return count
}

func Control_SwitchNoMatch() int {
	x := 100
	switch x {
	case 1:
		return 1
	case 2:
		return 2
	}
	return 0
}

func Control_SwitchDefault() int {
	x := 100
	switch x {
	case 1:
		return 1
	default:
		return 99
	}
}

// ------------------------------------------------------------------------
// Function Corner Cases
// ------------------------------------------------------------------------

func noop() {}

func Func_NoReturn() int {
	noop()
	return 42
}

func multi() (int, int, int) {
	return 1, 2, 3
}

func Func_MultipleReturnAll() int {
	a, b, c := multi()
	return a + b + c
}

func Func_MultipleReturnIgnore() int {
	a, _, c := multi()
	return a + c
}

func named() (result int) {
	result = 42
	return
}

func Func_NamedReturn() int {
	return named()
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Func_VariadicEmpty() int {
	return sum()
}

func Func_VariadicOne() int {
	return sum(42)
}

func Func_VariadicMultiple() int {
	return sum(1, 2, 3, 4, 5)
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func Func_RecursionBase() int {
	return fib(10)
}

// ------------------------------------------------------------------------
// Closure Corner Cases
// ------------------------------------------------------------------------

func Closure_CaptureVariable() int {
	x := 10
	f := func() int {
		return x
	}
	return f()
}

func Closure_ModifyCaptured() int {
	x := 10
	f := func() {
		x = 20
	}
	f()
	return x
}

func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func Closure_ReturnClosure() int {
	c := counter()
	c()
	c()
	return c()
}

func Closure_LoopCapture() int {
	var funcs []func() int
	for i := 0; i < 3; i++ {
		i := i // Capture local copy
		funcs = append(funcs, func() int {
			return i
		})
	}
	return funcs[0]() + funcs[1]() + funcs[2]()
}

// ------------------------------------------------------------------------
// Struct Corner Cases
// ------------------------------------------------------------------------

type Empty struct{}

func Struct_EmptyStruct() int {
	var e Empty
	_ = e
	return 0
}

type Data struct {
	x int
	y string
	z bool
}

func Struct_ZeroValueFields() int {
	var d Data
	if d.x == 0 && d.y == "" && d.z == false {
		return 1
	}
	return 0
}

type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func Struct_PointerReceiver() int {
	c := &Counter{value: 0}
	c.Inc()
	c.Inc()
	return c.value
}

type Inner struct {
	x int
}

type Outer struct {
	inner Inner
}

func Struct_NestedStruct() int {
	o := Outer{inner: Inner{x: 42}}
	return o.inner.x
}

// ------------------------------------------------------------------------
// Type Conversion Corner Cases
// ------------------------------------------------------------------------

func Convert_IntToFloat() float64 {
	var x int = 42
	return float64(x)
}

func Convert_FloatToInt() int {
	var x float64 = 42.9
	return int(x)
}

func Convert_Int64ToInt32() int32 {
	var x int64 = 100
	return int32(x)
}

func Convert_Int32ToInt64() int64 {
	var x int32 = 100
	return int64(x)
}

// ------------------------------------------------------------------------
// Complex Expression Corner Cases
// ------------------------------------------------------------------------

func Expr_ComplexArithmetic() int {
	return (1 + 2)*3 - 4/2
}

func Expr_NestedTernaryLike() int {
	x := 10
	result := 0
	if x > 5 {
		if x > 15 {
			result = 3
		} else {
			result = 2
		}
	} else {
		result = 1
	}
	return result
}

func Expr_MultipleAssignment() int {
	a, b := 1, 2
	a, b = b, a
	return a*10 + b
}

func Expr_ChainedComparison() bool {
	x := 5
	return x > 0 && x < 10
}

// ------------------------------------------------------------------------
// Map with Complex Keys/Values
// ------------------------------------------------------------------------

func Map_IntKey() string {
	m := map[int]string{
		1: "one",
		2: "two",
	}
	return m[1]
}

func Map_NegativeKey() string {
	m := map[int]string{
		-1: "negative",
		1:  "positive",
	}
	return m[-1]
}

type Key [2]int

func Map_SliceNotValidKey() int {
	m := map[Key]int{
		{1, 2}: 3,
	}
	return m[Key{1, 2}]
}

// ------------------------------------------------------------------------
// Edge Cases with Make
// ------------------------------------------------------------------------

func Make_SliceWithCap() int {
	s := make([]int, 5, 10)
	return len(s)*100 + cap(s)
}

func Make_MapWithSize() int {
	m := make(map[string]int, 10)
	m["a"] = 1
	return len(m)
}

func Make_ZeroLenZeroCap() int {
	s := make([]int, 0, 0)
	return len(s) + cap(s)
}

// ------------------------------------------------------------------------
// Range Corner Cases
// ------------------------------------------------------------------------

func Range_EmptySlice() int {
	s := []int{}
	count := 0
	for range s {
		count++
	}
	return count
}

func Range_EmptyMap() int {
	m := map[string]int{}
	count := 0
	for range m {
		count++
	}
	return count
}

func Range_EmptyString() int {
	s := ""
	count := 0
	for range s {
		count++
	}
	return count
}

func Range_SingleElement() int {
	s := []int{42}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}
