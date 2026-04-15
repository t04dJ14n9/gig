package cornercases

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
// Slice Operation Tests
// ------------------------------------------------------------------------

func Slice_ZeroToZero() int {
	s := []int{1, 2, 3}
	return len(s[0:0])
}

func Slice_EndToEnd() int {
	s := []int{1, 2, 3}
	return len(s[0:len(s)])
}

func Slice_FullSlice() []int {
	s := []int{1, 2, 3}
	return s[:]
}

func Slice_NilSlice() int {
	var s []int
	return len(s)
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
// Map Operation Tests
// ------------------------------------------------------------------------

func Map_NilMap() int {
	var m map[string]int
	return len(m)
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
	m := map[string]int{"": 1}
	return m[""]
}

func Map_ZeroIntKey() int {
	m := map[int]int{0: 1}
	return m[0]
}

// ------------------------------------------------------------------------
// String Boundary Tests
// ------------------------------------------------------------------------

func String_Empty() int {
	s := ""
	return len(s)
}

func String_SingleChar() int {
	s := "a"
	return len(s)
}

func String_UnicodeMultibyte() int {
	s := "你好"
	return len(s)
}

func String_Whitespace() int {
	s := " \t\n"
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

// ------------------------------------------------------------------------
// Boolean Boundary Tests
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
	b := false
	return !!b
}

// ------------------------------------------------------------------------
// Arithmetic Boundary Tests
// ------------------------------------------------------------------------

func Arith_AddZero() int {
	x := 42
	return x + 0
}

func Arith_SubZero() int {
	x := 42
	return x - 0
}

func Arith_MulByOne() int {
	x := 42
	return x * 1
}

func Arith_DivByOne() int {
	x := 42
	return x / 1
}

func Arith_ModByOne() int {
	x := 42
	return x % 1
}

func Arith_MulByZero() int {
	x := 42
	return x * 0
}

func Arith_NegNeg() int {
	x := -42
	return -x
}

func Arith_NegAddNeg() int {
	x := -10
	y := -20
	return x + y
}

// ------------------------------------------------------------------------
// Comparison Operation Tests
// ------------------------------------------------------------------------

func Compare_IntEqual() bool {
	return 5 == 5
}

func Compare_IntNotEqual() bool {
	return 5 != 6
}

func Compare_IntGreater() bool {
	return 6 > 5
}

func Compare_IntGreaterEqual() bool {
	return 5 >= 5
}

func Compare_IntLess() bool {
	return 4 < 5
}

func Compare_IntLessEqual() bool {
	return 5 <= 5
}

func Compare_StringEqual() bool {
	return "abc" == "abc"
}

func Compare_StringNotEqual() bool {
	return "abc" != "def"
}

func Compare_EmptyStringEqual() bool {
	return "" == ""
}

// ------------------------------------------------------------------------
// Logical Operation Tests
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
// Control Flow Tests
// ------------------------------------------------------------------------

func Control_IfNoElse() int {
	if true {
		return 1
	}
	return 0
}

func Control_IfFalseNoElse() int {
	if false {
		return 1
	}
	return 0
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
	for i := 0; i < 10; i++ {
		return i
	}
	return -1
}

func Control_ForContinueAll() int {
	sum := 0
	for i := 0; i < 3; i++ {
		if i < 3 {
			continue
		}
		sum += i
	}
	return sum
}

func Control_SwitchNoMatch() int {
	x := 5
	switch x {
	case 1:
		return 1
	case 2:
		return 2
	}
	return 0
}

func Control_SwitchDefault() int {
	x := 5
	switch x {
	case 1:
		return 1
	default:
		return -1
	}
}

// ------------------------------------------------------------------------
// Function Tests
// ------------------------------------------------------------------------

func Func_NoReturn() int {
	return 0
}

func Func_MultipleReturnAll() (int, int) {
	return 1, 2
}

func Func_MultipleReturnIgnore() int {
	a, _ := Func_MultipleReturnAll()
	return a
}

func Func_NamedReturn() (result int) {
	result = 42
	return
}

func Func_VariadicEmpty() int {
	return len([]int{})
}

func Func_VariadicOne() int {
	return len([]int{1})
}

func Func_VariadicMultiple() int {
	return len([]int{1, 2, 3})
}

func Func_RecursionBase() int {
	if 1 == 1 {
		return 1
	}
	return 0
}

// ------------------------------------------------------------------------
// Closure Tests
// ------------------------------------------------------------------------

func Closure_ReturnClosure() int {
	f := func() int {
		return 42
	}
	return f()
}

func Closure_CaptureVariable() int {
	x := 10
	f := func() int {
		return x
	}
	return f()
}

func Closure_ModifyCaptured() int {
	x := 10
	f := func() int {
		x = 20
		return x
	}
	return f()
}

// ------------------------------------------------------------------------
// Struct Tests
// ------------------------------------------------------------------------

func Struct_ZeroValueFields() int {
	type Point struct {
		X int
		Y int
	}
	var p Point
	return p.X + p.Y
}

func Struct_PointerReceiver() int {
	type Counter struct {
		value int
	}
	c := &Counter{value: 10}
	c.value++
	return c.value
}

func Struct_NestedStruct() int {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Inner Inner
	}
	o := Outer{Inner: Inner{Value: 42}}
	return o.Inner.Value
}

// ------------------------------------------------------------------------
// Type Conversion Tests
// ------------------------------------------------------------------------

func Convert_IntToFloat() float64 {
	return float64(42)
}

func Convert_FloatToInt() int {
	return int(42.0)
}

func Convert_Int32ToInt64() int64 {
	var x int32 = 100
	return int64(x)
}

func Convert_Int64ToInt32() int32 {
	var x int64 = 100
	return int32(x)
}

// ------------------------------------------------------------------------
// Complex Expression Tests
// ------------------------------------------------------------------------

func Expr_ComplexArithmetic() int {
	return (10+5)*2 - 3
}

func Expr_ChainedComparison() bool {
	x := 5
	return x > 0 && x < 10
}

func Expr_MultipleAssignment() int {
	a, b, c := 1, 2, 3
	return a + b + c
}

func Expr_NestedTernaryLike() int {
	x := 5
	if x > 0 {
		if x < 10 {
			return 1
		}
		return 2
	}
	return 0
}

// ------------------------------------------------------------------------
// Make Tests
// ------------------------------------------------------------------------

func Make_ZeroLenZeroCap() int {
	s := make([]int, 0, 0)
	return len(s) + cap(s)
}

func Make_SliceWithCap() int {
	s := make([]int, 0, 10)
	return cap(s)
}

func Make_MapWithSize() int {
	m := make(map[string]int, 10)
	return len(m)
}

// ------------------------------------------------------------------------
// Range Tests
// ------------------------------------------------------------------------

func Range_EmptySlice() int {
	count := 0
	for range []int{} {
		count++
	}
	return count
}

func Range_EmptyMap() int {
	count := 0
	for range map[string]int{} {
		count++
	}
	return count
}

func Range_EmptyString() int {
	count := 0
	for range "" {
		count++
	}
	return count
}

func Range_SingleElement() int {
	count := 0
	for range []int{1} {
		count++
	}
	return count
}

// ------------------------------------------------------------------------
// Additional Corner Case Tests
// ------------------------------------------------------------------------

// Slice_ThreeIndexSlice tests 3-index slice expression s[low:high:max] with cap control
func Slice_ThreeIndexSlice() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3] // len=2, cap=2
	return len(sub)*10 + cap(sub)
}

// Slice_ThreeIndexSliceCapIsolation tests 3-index slice doesn't share cap with original
func Slice_ThreeIndexSliceCapIsolation() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3] // cap limited to 2
	_ = append(sub, 10)
	// s[3] should still be 4, not modified by append to sub
	return s[3]
}

// Slice_AppendToNilString tests append to nil slice with non-int type
func Slice_AppendToNilString() int {
	var s []string
	s = append(s, "hello")
	s = append(s, "world")
	return len(s)
}

// Slice_AppendToNilFloat tests append to nil slice with float type
func Slice_AppendToNilFloat() float64 {
	var s []float64
	s = append(s, 1.5)
	s = append(s, 2.5)
	return s[0] + s[1]
}

// Map_NilMapReadOk tests reading from nil map with comma-ok
func Map_NilMapReadOk() int {
	var m map[string]int
	_, ok := m["key"]
	if ok {
		return 1
	}
	return 0
}

// Map_DeleteNilMap tests delete on nil map (should be no-op)
func Map_DeleteNilMap() int {
	var m map[string]int
	delete(m, "key")
	return 0 // no panic
}

// ComplexTypeAssertion tests type assertion with complex64/128
func ComplexTypeAssertion() int {
	var i interface{} = complex(3.0, 4.0)
	switch v := i.(type) {
	case complex64:
		_ = v
		return 1
	case complex128:
		_ = v
		return 2
	default:
		_ = v
		return 0
	}
}

// CrossKindTypeAssertion tests type assertion across numeric kinds
func CrossKindTypeAssertion() int {
	var i interface{} = int64(42)
	switch i.(type) {
	case int:
		return 1
	case int64:
		return 2
	case int32:
		return 3
	default:
		return 0
	}
}
