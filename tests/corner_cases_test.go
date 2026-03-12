package tests

import (
	"testing"

	"git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
)

// ============================================================================
// Corner Case Tests - Comprehensive Edge Cases for Robustness
// ============================================================================

// cornerCaseTest defines a corner case test
type cornerCaseTest struct {
	name     string
	src      string
	funcName string
	expected any
}

// allCornerCases contains all corner case tests
var allCornerCases = []cornerCaseTest{
	// ------------------------------------------------------------------------
	// Zero Value Tests
	// ------------------------------------------------------------------------
	{
		name: "ZeroValue_Int",
		src: `
package main

func Test() int {
	var x int
	return x
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "ZeroValue_Int64",
		src: `
package main

func Test() int64 {
	var x int64
	return x
}`,
		funcName: "Test",
		expected: int64(0),
	},
	{
		name: "ZeroValue_Float64",
		src: `
package main

func Test() float64 {
	var x float64
	return x
}`,
		funcName: "Test",
		expected: float64(0),
	},
	{
		name: "ZeroValue_String",
		src: `
package main

func Test() string {
	var s string
	return s
}`,
		funcName: "Test",
		expected: "",
	},
	{
		name: "ZeroValue_Bool",
		src: `
package main

func Test() bool {
	var b bool
	return b
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "ZeroValue_Slice",
		src: `
package main

func Test() int {
	var s []int
	if s == nil {
		return 1
	}
	return 0
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "ZeroValue_Map",
		src: `
package main

func Test() int {
	var m map[string]int
	if m == nil {
		return 1
	}
	return 0
}`,
		funcName: "Test",
		expected: int(1),
	},

	// ------------------------------------------------------------------------
	// Integer Boundary Tests
	// ------------------------------------------------------------------------
	{
		name: "IntBoundary_MaxInt32",
		src: `
package main

func Test() int32 {
	return 2147483647
}`,
		funcName: "Test",
		expected: int64(2147483647), // Gig returns int64
	},
	{
		name: "IntBoundary_MinInt32",
		src: `
package main

func Test() int32 {
	return -2147483648
}`,
		funcName: "Test",
		expected: int64(-2147483648), // Gig returns int64
	},
	{
		name: "IntBoundary_MaxInt64",
		src: `
package main

func Test() int64 {
	return 9223372036854775807
}`,
		funcName: "Test",
		expected: int64(9223372036854775807),
	},
	{
		name: "IntBoundary_MinInt64",
		src: `
package main

func Test() int64 {
	return -9223372036854775808
}`,
		funcName: "Test",
		expected: int64(-9223372036854775808),
	},
	{
		name: "IntBoundary_MaxUint32",
		src: `
package main

func Test() uint32 {
	return 4294967295
}`,
		funcName: "Test",
		expected: int64(4294967295), // Gig returns int64 for uint32
	},
	{
		name: "IntBoundary_NearMaxInt",
		src: `
package main

func Test() int {
	return 2147483646
}`,
		funcName: "Test",
		expected: int(2147483646),
	},
	{
		name: "IntBoundary_NearMinInt",
		src: `
package main

func Test() int {
	return -2147483647
}`,
		funcName: "Test",
		expected: int(-2147483647),
	},

	// ------------------------------------------------------------------------
	// Integer Overflow Tests
	// ------------------------------------------------------------------------
	{
		name: "Overflow_Int32Add",
		src: `
package main

func Test() int32 {
	var x int32 = 2147483647
	return x + 1
}`,
		funcName: "Test",
		expected: int64(-2147483648), // Gig returns int64, overflow wraps around
	},
	{
		name: "Overflow_Int32Sub",
		src: `
package main

func Test() int32 {
	var x int32 = -2147483648
	return x - 1
}`,
		funcName: "Test",
		expected: int64(2147483647), // Gig returns int64, overflow wraps around
	},
	{
		name: "Overflow_Int32Mul",
		src: `
package main

func Test() int32 {
	var x int32 = 65536
	return x * 65536
}`,
		funcName: "Test",
		expected: int64(0), // Gig returns int64, overflow
	},

	// ------------------------------------------------------------------------
	// Float Boundary Tests
	// ------------------------------------------------------------------------
	{
		name: "FloatBoundary_SmallPositive",
		src: `
package main

func Test() float64 {
	return 1e-300
}`,
		funcName: "Test",
		expected: 1e-300,
	},
	{
		name: "FloatBoundary_SmallNegative",
		src: `
package main

func Test() float64 {
	return -1e-300
}`,
		funcName: "Test",
		expected: -1e-300,
	},
	{
		name: "FloatBoundary_LargePositive",
		src: `
package main

func Test() float64 {
	return 1e300
}`,
		funcName: "Test",
		expected: 1e300,
	},
	{
		name: "FloatBoundary_LargeNegative",
		src: `
package main

func Test() float64 {
	return -1e300
}`,
		funcName: "Test",
		expected: -1e300,
	},

	// ------------------------------------------------------------------------
	// Empty Collection Tests
	// ------------------------------------------------------------------------
	{
		name: "EmptySlice_Len",
		src: `
package main

func Test() int {
	s := []int{}
	return len(s)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "EmptySlice_Cap",
		src: `
package main

func Test() int {
	s := []int{}
	return cap(s)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "EmptySlice_Make",
		src: `
package main

func Test() int {
	s := make([]int, 0)
	return len(s)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "EmptyMap_Len",
		src: `
package main

func Test() int {
	m := map[string]int{}
	return len(m)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "EmptyMap_Make",
		src: `
package main

func Test() int {
	m := make(map[string]int)
	return len(m)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "EmptyString_Len",
		src: `
package main

func Test() int {
	s := ""
	return len(s)
}`,
		funcName: "Test",
		expected: int(0),
	},

	// ------------------------------------------------------------------------
	// Slice Operations Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Slice_ZeroToZero",
		src: `
package main

func Test() int {
	s := []int{1, 2, 3}
	sub := s[0:0]
	return len(sub)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Slice_EndToEnd",
		src: `
package main

func Test() int {
	s := []int{1, 2, 3}
	sub := s[3:3]
	return len(sub)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Slice_FullSlice",
		src: `
package main

func Test() int {
	s := []int{1, 2, 3}
	sub := s[:]
	return len(sub)
}`,
		funcName: "Test",
		expected: int(3),
	},
	{
		name: "Slice_NilSlice",
		src: `
package main

func Test() int {
	var s []int
	if s == nil {
		return 1
	}
	return 0
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Slice_AppendToNil",
		src: `
package main

func Test() int {
	var s []int
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return len(s)
}`,
		funcName: "Test",
		expected: int(3),
	},
	{
		name: "Slice_AppendEmpty",
		src: `
package main

func Test() int {
	s := []int{1, 2}
	_ = append(s)
	return len(s)
}`,
		funcName: "Test",
		expected: int(2),
	},

	// ------------------------------------------------------------------------
	// Map Operations Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Map_NilMap",
		src: `
package main

func Test() int {
	var m map[string]int
	if m == nil {
		return 1
	}
	return 0
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Map_AccessMissingKey",
		src: `
package main

func Test() int {
	m := map[string]int{"a": 1}
	return m["b"]
}`,
		funcName: "Test",
		expected: int(0), // Zero value for missing key
	},
	{
		name: "Map_DeleteMissingKey",
		src: `
package main

func Test() int {
	m := map[string]int{"a": 1}
	delete(m, "b")
	return len(m)
}`,
		funcName: "Test",
		expected: int(1), // No effect
	},
	{
		name: "Map_OverwriteKey",
		src: `
package main

func Test() int {
	m := map[string]int{"a": 1}
	m["a"] = 2
	return m["a"]
}`,
		funcName: "Test",
		expected: int(2),
	},
	{
		name: "Map_NilKeyString",
		src: `
package main

func Test() int {
	m := map[string]int{"" : 42}
	return m[""]
}`,
		funcName: "Test",
		expected: int(42),
	},
	{
		name: "Map_ZeroIntKey",
		src: `
package main

func Test() int {
	m := map[int]string{0: "zero"}
	return len(m[0])
}`,
		funcName: "Test",
		expected: int(4),
	},

	// ------------------------------------------------------------------------
	// String Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "String_Empty",
		src: `
package main

func Test() int {
	s := ""
	return len(s)
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
	 name: "String_SingleChar",
		src: `
package main

func Test() int {
	s := "a"
	return len(s)
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "String_SingleChar",
		src: `
package main

func Test() int {
	s := "a"
	return len(s)
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
	 name: "String_SingleByteIndex",
		src: `
package main

func Test() uint8 {
	s := "abc"
	return s[1]
}`,
		funcName: "Test",
		expected: int64('b'), // Gig returns int64 for uint8
	},
	{
		name: "String_LastByte",
		src: `
package main

func Test() uint8 {
	s := "hello"
	return s[len(s)-1]
}`,
		funcName: "Test",
		expected: int64('o'), // Gig returns int64 for uint8
	},
	{
		name: "String_Whitespace",
		src: `
package main

func Test() int {
	s := " \t\n"
	return len(s)
}`,
		funcName: "Test",
		expected: int(3),
	},
	{
		name: "String_UnicodeMultibyte",
		src: `
package main

func Test() int {
	s := "你好"
	return len(s)
}`,
		funcName: "Test",
		expected: int(6), // 3 bytes per Chinese character
	},

	// ------------------------------------------------------------------------
	// Boolean Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Bool_True",
		src: `
package main

func Test() bool {
	return true
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Bool_False",
		src: `
package main

func Test() bool {
	return false
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "Bool_NotTrue",
		src: `
package main

func Test() bool {
	return !true
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "Bool_NotFalse",
		src: `
package main

func Test() bool {
	return !false
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Bool_DoubleNegation",
		src: `
package main

func Test() bool {
	return !!true
}`,
		funcName: "Test",
		expected: true,
	},

	// ------------------------------------------------------------------------
	// Arithmetic Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Arith_DivByOne",
		src: `
package main

func Test() int {
	return 100 / 1
}`,
		funcName: "Test",
		expected: int(100),
	},
	{
		name: "Arith_ModByOne",
		src: `
package main

func Test() int {
	return 100 % 1
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Arith_MulByZero",
		src: `
package main

func Test() int {
	return 100 * 0
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Arith_MulByOne",
		src: `
package main

func Test() int {
	return 100 * 1
}`,
		funcName: "Test",
		expected: int(100),
	},
	{
		name: "Arith_AddZero",
		src: `
package main

func Test() int {
	return 100 + 0
}`,
		funcName: "Test",
		expected: int(100),
	},
	{
		name: "Arith_SubZero",
		src: `
package main

func Test() int {
	return 100 - 0
}`,
		funcName: "Test",
		expected: int(100),
	},
	{
		name: "Arith_NegNeg",
		src: `
package main

func Test() int {
	return -(-100)
}`,
		funcName: "Test",
		expected: int(100),
	},
	{
		name: "Arith_NegAddNeg",
		src: `
package main

func Test() int {
	return -10 + (-20)
}`,
		funcName: "Test",
		expected: int(-30),
	},

	// ------------------------------------------------------------------------
	// Comparison Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Compare_IntEqual",
		src: `
package main

func Test() bool {
	return 5 == 5
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_IntNotEqual",
		src: `
package main

func Test() bool {
	return 5 != 6
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_IntLess",
		src: `
package main

func Test() bool {
	return 5 < 6
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_IntLessEqual",
		src: `
package main

func Test() bool {
	return 5 <= 5
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_IntGreater",
		src: `
package main

func Test() bool {
	return 6 > 5
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_IntGreaterEqual",
		src: `
package main

func Test() bool {
	return 5 >= 5
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_StringEqual",
		src: `
package main

func Test() bool {
	return "hello" == "hello"
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_StringNotEqual",
		src: `
package main

func Test() bool {
	return "hello" != "world"
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Compare_EmptyStringEqual",
		src: `
package main

func Test() bool {
	return "" == ""
}`,
		funcName: "Test",
		expected: true,
	},

	// ------------------------------------------------------------------------
	// Logical Operation Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Logic_TrueAndTrue",
		src: `
package main

func Test() bool {
	return true && true
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Logic_TrueAndFalse",
		src: `
package main

func Test() bool {
	return true && false
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "Logic_FalseAndTrue",
		src: `
package main

func Test() bool {
	return false && true
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "Logic_TrueOrFalse",
		src: `
package main

func Test() bool {
	return true || false
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Logic_FalseOrTrue",
		src: `
package main

func Test() bool {
	return false || true
}`,
		funcName: "Test",
		expected: true,
	},
	{
		name: "Logic_FalseOrFalse",
		src: `
package main

func Test() bool {
	return false || false
}`,
		funcName: "Test",
		expected: false,
	},

	// ------------------------------------------------------------------------
	// Short Circuit Evaluation Tests
	// ------------------------------------------------------------------------
	{
		name: "ShortCircuit_AndFalse",
		src: `
package main

func panicFunc() bool {
	panic("should not be called")
}

func Test() bool {
	return false && panicFunc()
}`,
		funcName: "Test",
		expected: false,
	},
	{
		name: "ShortCircuit_OrTrue",
		src: `
package main

func panicFunc() bool {
	panic("should not be called")
}

func Test() bool {
	return true || panicFunc()
}`,
		funcName: "Test",
		expected: true,
	},

	// ------------------------------------------------------------------------
	// Control Flow Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Control_IfNoElse",
		src: `
package main

func Test() int {
	x := 0
	if true {
		x = 1
	}
	return x
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Control_IfFalseNoElse",
		src: `
package main

func Test() int {
	x := 0
	if false {
		x = 1
	}
	return x
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Control_ForZeroIter",
		src: `
package main

func Test() int {
	count := 0
	for i := 0; i < 0; i++ {
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Control_ForOneIter",
		src: `
package main

func Test() int {
	count := 0
	for i := 0; i < 1; i++ {
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Control_ForBreakFirst",
		src: `
package main

func Test() int {
	count := 0
	for i := 0; i < 10; i++ {
		break
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Control_ForContinueAll",
		src: `
package main

func Test() int {
	count := 0
	for i := 0; i < 5; i++ {
		continue
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Control_SwitchNoMatch",
		src: `
package main

func Test() int {
	x := 100
	switch x {
	case 1:
		return 1
	case 2:
		return 2
	}
	return 0
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Control_SwitchDefault",
		src: `
package main

func Test() int {
	x := 100
	switch x {
	case 1:
		return 1
	default:
		return 99
	}
}`,
		funcName: "Test",
		expected: int(99),
	},

	// ------------------------------------------------------------------------
	// Function Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Func_NoReturn",
		src: `
package main

func noop() {}

func Test() int {
	noop()
	return 42
}`,
		funcName: "Test",
		expected: int(42),
	},
	{
		name: "Func_MultipleReturnAll",
		src: `
package main

func multi() (int, int, int) {
	return 1, 2, 3
}

func Test() int {
	a, b, c := multi()
	return a + b + c
}`,
		funcName: "Test",
		expected: int(6),
	},
	{
		name: "Func_MultipleReturnIgnore",
		src: `
package main

func multi() (int, int, int) {
	return 1, 2, 3
}

func Test() int {
	a, _, c := multi()
	return a + c
}`,
		funcName: "Test",
		expected: int(4),
	},
	{
		name: "Func_NamedReturn",
		src: `
package main

func named() (result int) {
	result = 42
	return
}

func Test() int {
	return named()
}`,
		funcName: "Test",
		expected: int(42),
	},
	{
		name: "Func_VariadicEmpty",
		src: `
package main

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Test() int {
	return sum()
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Func_VariadicOne",
		src: `
package main

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Test() int {
	return sum(42)
}`,
		funcName: "Test",
		expected: int(42),
	},
	{
		name: "Func_VariadicMultiple",
		src: `
package main

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func Test() int {
	return sum(1, 2, 3, 4, 5)
}`,
		funcName: "Test",
		expected: int(15),
	},
	{
		name: "Func_RecursionBase",
		src: `
package main

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func Test() int {
	return fib(10)
}`,
		funcName: "Test",
		expected: int(55),
	},

	// ------------------------------------------------------------------------
	// Closure Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Closure_CaptureVariable",
		src: `
package main

func Test() int {
	x := 10
	f := func() int {
		return x
	}
	return f()
}`,
		funcName: "Test",
		expected: int(10),
	},
	{
		name: "Closure_ModifyCaptured",
		src: `
package main

func Test() int {
	x := 10
	f := func() {
		x = 20
	}
	f()
	return x
}`,
		funcName: "Test",
		expected: int(20),
	},
	{
		name: "Closure_ReturnClosure",
		src: `
package main

func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func Test() int {
	c := counter()
	c()
	c()
	return c()
}`,
		funcName: "Test",
		expected: int(3),
	},
	{
		name: "Closure_LoopCapture",
		src: `
package main

func Test() int {
	var funcs []func() int
	for i := 0; i < 3; i++ {
		i := i // Capture local copy
		funcs = append(funcs, func() int {
			return i
		})
	}
	return funcs[0]() + funcs[1]() + funcs[2]()
}`,
		funcName: "Test",
		expected: int(3), // 0 + 1 + 2
	},

	// ------------------------------------------------------------------------
	// Struct Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Struct_EmptyStruct",
		src: `
package main

type Empty struct{}

func Test() int {
	var e Empty
	return 0
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Struct_ZeroValueFields",
		src: `
package main

type Data struct {
	x int
	y string
	z bool
}

func Test() int {
	var d Data
	if d.x == 0 && d.y == "" && d.z == false {
		return 1
	}
	return 0
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Struct_PointerReceiver",
		src: `
package main

type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func Test() int {
	c := &Counter{value: 0}
	c.Inc()
	c.Inc()
	return c.value
}`,
		funcName: "Test",
		expected: int(2),
	},
	{
		name: "Struct_NestedStruct",
		src: `
package main

type Inner struct {
	x int
}

type Outer struct {
	inner Inner
}

func Test() int {
	o := Outer{inner: Inner{x: 42}}
	return o.inner.x
}`,
		funcName: "Test",
		expected: int(42),
	},

	// ------------------------------------------------------------------------
	// Type Conversion Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Convert_IntToFloat",
		src: `
package main

func Test() float64 {
	var x int = 42
	return float64(x)
}`,
		funcName: "Test",
		expected: float64(42),
	},
	{
		name: "Convert_FloatToInt",
		src: `
package main

func Test() int {
	var x float64 = 42.9
	return int(x)
}`,
		funcName: "Test",
		expected: int(42), // Truncates
	},
	{
		name: "Convert_Int64ToInt32",
		src: `
package main

func Test() int32 {
	var x int64 = 100
	return int32(x)
}`,
		funcName: "Test",
		expected: int32(100),
	},
	{
		name: "Convert_Int32ToInt64",
		src: `
package main

func Test() int64 {
	var x int32 = 100
	return int64(x)
}`,
		funcName: "Test",
		expected: int64(100),
	},

	// ------------------------------------------------------------------------
	// Complex Expression Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Expr_ComplexArithmetic",
		src: `
package main

func Test() int {
	return (1 + 2) * 3 - 4 / 2
}`,
		funcName: "Test",
		expected: int(7), // (1+2)*3 - 4/2 = 9 - 2 = 7
	},
	{
		name: "Expr_NestedTernaryLike",
		src: `
package main

func Test() int {
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
}`,
		funcName: "Test",
		expected: int(2),
	},
	{
		name: "Expr_MultipleAssignment",
		src: `
package main

func Test() int {
	a, b := 1, 2
	a, b = b, a
	return a*10 + b
}`,
		funcName: "Test",
		expected: int(21), // a=2, b=1
	},
	{
		name: "Expr_ChainedComparison",
		src: `
package main

func Test() bool {
	x := 5
	return x > 0 && x < 10
}`,
		funcName: "Test",
		expected: true,
	},

	// ------------------------------------------------------------------------
	// Map with Complex Keys/Values
	// ------------------------------------------------------------------------
	{
		name: "Map_IntKey",
		src: `
package main

func Test() string {
	m := map[int]string{
		1: "one",
		2: "two",
	}
	return m[1]
}`,
		funcName: "Test",
		expected: "one",
	},
	{
		name: "Map_NegativeKey",
		src: `
package main

func Test() string {
	m := map[int]string{
		-1: "negative",
		1:  "positive",
	}
	return m[-1]
}`,
		funcName: "Test",
		expected: "negative",
	},
	{
		name: "Map_SliceNotValidKey",
		src: `
package main

func Test() int {
	// Slices cannot be map keys, but we can use arrays
	type Key [2]int
	m := map[Key]int{
		{1, 2}: 3,
	}
	return m[Key{1, 2}]
}`,
		funcName: "Test",
		expected: int(3),
	},

	// ------------------------------------------------------------------------
	// Edge Cases with Make
	// ------------------------------------------------------------------------
	{
		name: "Make_SliceWithCap",
		src: `
package main

func Test() int {
	s := make([]int, 5, 10)
	return len(s) * 100 + cap(s)
}`,
		funcName: "Test",
		expected: int(510), // len=5, cap=10 -> 5*100+10=510
	},
	{
		name: "Make_MapWithSize",
		src: `
package main

func Test() int {
	m := make(map[string]int, 10)
	m["a"] = 1
	return len(m)
}`,
		funcName: "Test",
		expected: int(1),
	},
	{
		name: "Make_ZeroLenZeroCap",
		src: `
package main

func Test() int {
	s := make([]int, 0, 0)
	return len(s) + cap(s)
}`,
		funcName: "Test",
		expected: int(0),
	},

	// ------------------------------------------------------------------------
	// Range Corner Cases
	// ------------------------------------------------------------------------
	{
		name: "Range_EmptySlice",
		src: `
package main

func Test() int {
	s := []int{}
	count := 0
	for range s {
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Range_EmptyMap",
		src: `
package main

func Test() int {
	m := map[string]int{}
	count := 0
	for range m {
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Range_EmptyString",
		src: `
package main

func Test() int {
	s := ""
	count := 0
	for range s {
		count++
	}
	return count
}`,
		funcName: "Test",
		expected: int(0),
	},
	{
		name: "Range_SingleElement",
		src: `
package main

func Test() int {
	s := []int{42}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}`,
		funcName: "Test",
		expected: int(42),
	},
}

// TestCornerCases runs all corner case tests
func TestCornerCases(t *testing.T) {
	for _, tc := range allCornerCases {
		t.Run(tc.name, func(t *testing.T) {
			prog, err := gig.Build(tc.src)
			if err != nil {
				t.Fatalf("Build error: %v", err)
			}

			result, err := prog.Run(tc.funcName)
			if err != nil {
				t.Fatalf("Run error: %v", err)
			}

			compareCornerCaseResult(t, result, tc.expected)
		})
	}
}

// compareCornerCaseResult compares result with expected value
func compareCornerCaseResult(t *testing.T, result, expected any) {
	t.Helper()

	switch exp := expected.(type) {
	case int:
		var got int64
		switch v := result.(type) {
		case int64:
			got = v
		case int:
			got = int64(v)
		default:
			t.Fatalf("expected int, got %T", result)
		}
		if got != int64(exp) {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case int32:
		got, ok := result.(int32)
		if !ok {
			t.Fatalf("expected int32, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case int64:
		got, ok := result.(int64)
		if !ok {
			t.Fatalf("expected int64, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case uint32:
		got, ok := result.(uint32)
		if !ok {
			t.Fatalf("expected uint32, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case uint8:
		got, ok := result.(uint8)
		if !ok {
			t.Fatalf("expected uint8, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %d, got %d", exp, got)
		}

	case float64:
		got, ok := result.(float64)
		if !ok {
			t.Fatalf("expected float64, got %T", result)
		}
		// Use approximate comparison for floats
		diff := got - exp
		if diff < 0 {
			diff = -diff
		}
		if diff > 1e-10 {
			t.Errorf("expected %v, got %v", exp, got)
		}

	case string:
		got, ok := result.(string)
		if !ok {
			t.Fatalf("expected string, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %q, got %q", exp, got)
		}

	case bool:
		got, ok := result.(bool)
		if !ok {
			t.Fatalf("expected bool, got %T", result)
		}
		if got != exp {
			t.Errorf("expected %v, got %v", exp, got)
		}

	default:
		t.Fatalf("unsupported expected type: %T", expected)
	}
}
