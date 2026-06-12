package parity

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

// ============================================================================
// Type assertions
// ============================================================================

// TypeAssertInt tests basic type assertion on an interface{} holding int.
func TypeAssertInt() any {
	var x any = 42
	v := x.(int)
	return v
}

// TypeAssertCommaOk tests comma-ok type assertion (success case).
func TypeAssertCommaOk() any {
	var x any = "hello"
	v, ok := x.(string)
	return fmt.Sprintf("%s:%v", v, ok)
}

// TypeAssertCommaOkFail tests comma-ok type assertion (failure case).
func TypeAssertCommaOkFail() any {
	var x any = 42
	v, ok := x.(string)
	return fmt.Sprintf("%q:%v", v, ok)
}

// TypeAssertInterfaceToInterface tests asserting interface{} to interface{}.
func TypeAssertInterfaceToInterface() any {
	var x any = 99
	v, ok := x.(any)
	return fmt.Sprintf("%v:%v", v, ok)
}

// TypeSwitch tests a type switch on an interface{} value.
func TypeSwitch() any {
	values := []any{42, "hello", 3.14, true}
	results := ""
	for _, v := range values {
		switch val := v.(type) {
		case int:
			results += fmt.Sprintf("int:%d ", val)
		case string:
			results += fmt.Sprintf("str:%s ", val)
		case float64:
			results += fmt.Sprintf("float:%.2f ", val)
		default:
			results += fmt.Sprintf("other:%T ", val)
		}
	}
	return strings.TrimSpace(results)
}

// TypeSwitchWithDefault tests type switch with default branch.
func TypeSwitchWithDefault() any {
	var x any = struct{}{}
	switch x.(type) {
	case int:
		return "int"
	case string:
		return "string"
	default:
		return "default"
	}
}

// ============================================================================
// Type conversions
// ============================================================================

// IntToFloat64Conv tests int to float64 conversion.
func IntToFloat64Conv() any {
	var x int = 42
	var y float64 = float64(x)
	return y
}

// Float64ToIntConv tests float64 to int conversion (truncation).
func Float64ToIntConv() any {
	var x float64 = 3.99
	var y int = int(x)
	return y
}

// UintToIntConv tests uint to int conversion.
func UintToIntConv() any {
	var x uint = 100
	var y int = int(x)
	return y
}

// IntToUintConv tests int to uint conversion.
func IntToUintConv() any {
	var x int = 42
	var y uint = uint(x)
	return y
}

// RuneToStringConv tests rune to string conversion.
func RuneToStringConv() any {
	var r rune = 65 // 'A'
	return string(r)
}

// ByteToStringConv tests byte to string conversion.
func ByteToStringConv() any {
	var b byte = 66 // 'B'
	return string(b)
}

// StringToRuneSliceConv tests string to []rune conversion.
func StringToRuneSliceConv() any {
	s := "Hello"
	runes := []rune(s)
	return fmt.Sprintf("%d:%d:%d", runes[0], runes[1], len(runes))
}

// StringToByteSliceConv tests string to []byte conversion.
func StringToByteSliceConv() any {
	s := "Hi"
	b := []byte(s)
	return fmt.Sprintf("%d:%d:%d", b[0], b[1], len(b))
}

// ByteSliceToStringConv tests []byte to string conversion.
func ByteSliceToStringConv() any {
	b := []byte{72, 105} // "Hi"
	return string(b)
}

// Int8Range tests int8 boundaries.
func Int8Range() any {
	var x int8 = 127
	var y int8 = -128
	return fmt.Sprintf("%d:%d", x, y)
}

// Uint8Range tests uint8 boundaries.
func Uint8Range() any {
	var x uint8 = 255
	var y uint8 = 0
	return fmt.Sprintf("%d:%d", x, y)
}

// ============================================================================
// Interface wrapping (typed nil semantics)
// ============================================================================

type errorImpl struct {
	msg string
}

func (e *errorImpl) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.msg
}

// TypedNilInterface tests that a typed nil pointer in an interface is NOT nil.
func TypedNilInterface() any {
	var p *errorImpl
	var err error = p
	if err == nil {
		return "nil"
	}
	return "not-nil"
}

// NonNilInterface tests that a non-nil pointer in an interface is not nil.
func NonNilInterface() any {
	p := &errorImpl{"oops"}
	var err error = p
	if err == nil {
		return "nil"
	}
	return "not-nil"
}

// ============================================================================
// Method dispatch
// ============================================================================

type Counter struct {
	Value int
}

func (c *Counter) Inc()     { c.Value++ }
func (c *Counter) Get() int { return c.Value }

// PointerReceiverMethod tests pointer receiver method calls.
func PointerReceiverMethod() any {
	c := &Counter{10}
	c.Inc()
	c.Inc()
	return c.Get()
}

// ValueToPointerMethod tests calling pointer receiver on addressable value.
func ValueToPointerMethod() any {
	c := Counter{5}
	c.Inc()
	return c.Get()
}

// ============================================================================
// errors.Is / errors.As
// ============================================================================

// ErrorsIsWrap tests errors.Is with fmt.Errorf %w wrapping.
func ErrorsIsWrap() any {
	base := errors.New("base")
	wrapped := fmt.Errorf("wrapped: %w", base)
	return errors.Is(wrapped, base)
}

// ErrorsIsChain tests errors.Is through a chain of wrappings.
func ErrorsIsChain() any {
	root := errors.New("root")
	mid := fmt.Errorf("mid: %w", root)
	top := fmt.Errorf("top: %w", mid)
	return errors.Is(top, root)
}

// ErrorsIsNoMatch tests errors.Is when errors don't match.
func ErrorsIsNoMatch() any {
	e1 := errors.New("one")
	e2 := errors.New("two")
	return errors.Is(e1, e2)
}

// ErrorsAsMatch tests errors.As with a matching concrete type.
func ErrorsAsMatch() any {
	err := &errorImpl{"test"}
	var target *errorImpl
	if errors.As(err, &target) {
		return target.Error()
	}
	return "no match"
}

// ErrorsAsNoMatch tests errors.As with non-matching type.
func ErrorsAsNoMatch() any {
	err := errors.New("plain")
	var target *errorImpl
	if errors.As(err, &target) {
		return "matched"
	}
	return "no match"
}

// ============================================================================
// sort callbacks
// ============================================================================

// SortInts tests sort.Ints.
func SortInts() any {
	nums := []int{5, 3, 1, 4, 2}
	sort.Ints(nums)
	return fmt.Sprintf("%v", nums)
}

// SortStrings tests sort.Strings.
func SortStrings() any {
	s := []string{"banana", "apple", "cherry"}
	sort.Strings(s)
	return fmt.Sprintf("%v", s)
}

// SortFloat64s tests sort.Float64s.
func SortFloat64s() any {
	f := []float64{3.3, 1.1, 2.2}
	sort.Float64s(f)
	return fmt.Sprintf("%v", f)
}

// SortSliceCustom tests sort.Slice with a custom less function.
func SortSliceCustom() any {
	words := []string{"pie", "apple", "banana"}
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j])
	})
	return fmt.Sprintf("%v", words)
}

// SortSliceStable tests sort.SliceStable stability.
func SortSliceStable() any {
	type item struct {
		name  string
		order int
	}
	items := []item{{"a", 2}, {"b", 1}, {"c", 2}, {"d", 1}}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].order < items[j].order
	})
	result := ""
	for _, it := range items {
		result += it.name
	}
	return result
}

// SortSearch tests sort.SearchInts.
func SortSearchInts() any {
	nums := []int{10, 20, 30, 40, 50}
	idx := sort.SearchInts(nums, 30)
	return idx
}

// SortReverseInts tests sort.Reverse with sort.IntSlice.
func SortReverseInts() any {
	nums := []int{1, 3, 2, 5, 4}
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	return fmt.Sprintf("%v", nums)
}

// ============================================================================
// String operations
// ============================================================================

// StringContains tests strings.Contains.
func StringContains() any {
	return strings.Contains("hello world", "world")
}

// StringReplace tests strings.Replace.
func StringReplace() any {
	return strings.Replace("hello world", "world", "golang", 1)
}

// StringSplitJoin tests strings.Split + strings.Join roundtrip.
func StringSplitJoin() any {
	parts := strings.Split("a,b,c", ",")
	return strings.Join(parts, "-")
}

// StringHasPrefix tests strings.HasPrefix.
func StringHasPrefix() any {
	return strings.HasPrefix("hello.go", "hello")
}

// StringTrimSpace tests strings.TrimSpace.
func StringTrimSpace() any {
	return strings.TrimSpace("  hello  ")
}

// ============================================================================
// Slicing edge cases
// ============================================================================

// ThreeIndexSlice tests three-index slicing [low:high:max].
func ThreeIndexSlice() any {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s2 := s[2:5:8]
	return fmt.Sprintf("len=%d,cap=%d,vals=%v", len(s2), cap(s2), s2)
}

// SliceFromSlice tests creating a slice from another slice.
func SliceFromSlice() any {
	s := []int{10, 20, 30, 40, 50}
	s2 := s[1:4]
	return fmt.Sprintf("len=%d,cap=%d,vals=%v", len(s2), cap(s2), s2)
}

// SliceAppendShared tests that append to a sub-slice may share backing array.
func SliceAppendShared() any {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[0:3]
	s2 = append(s2, 99)
	return fmt.Sprintf("s=%v,s2=%v", s, s2)
}

// NilSliceLenCap tests len/cap on nil slice.
func NilSliceLenCap() any {
	var s []int
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// ============================================================================
// Map operations
// ============================================================================

// MapDeleteKey tests deleting a key from a map.
func MapDeleteKey() any {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return fmt.Sprintf("len=%d,a=%d,c=%d", len(m), m["a"], m["c"])
}

// MapMissingKey tests accessing a missing key returns zero value.
func MapMissingKey() any {
	m := map[string]int{"a": 1}
	return m["missing"]
}

// MapCommaOk tests comma-ok map access.
func MapCommaOk() any {
	m := map[string]int{"a": 1}
	_, ok1 := m["a"]
	_, ok2 := m["b"]
	return fmt.Sprintf("%v:%v", ok1, ok2)
}

// MapRange tests ranging over a map (order-independent).
func MapRange() any {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// ============================================================================
// Closure and capture
// ============================================================================

// ClosureCaptureByRef tests that closures capture variables by reference.
func ClosureCaptureByRef() any {
	x := 0
	inc := func() { x++ }
	inc()
	inc()
	inc()
	return x
}

// ClosureReturnFunc tests returning a closure.
func ClosureReturnFunc() any {
	makeAdder := func(n int) func(int) int {
		return func(m int) int { return n + m }
	}
	add5 := makeAdder(5)
	return add5(10)
}

// ============================================================================
// Defer / Panic / Recover
// ============================================================================

// DeferOrder tests that defers execute in LIFO order.
func DeferOrder() any {
	result := ""
	defer func() { result += "c" }()
	defer func() { result += "b" }()
	defer func() { result += "a" }()
	return result
}

// PanicRecover tests panic/recover.
func PanicRecover() any {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("recovered: %v", r)
			}
		}()
		panic("boom")
	}()
	return result
}

// RecoverNil tests that recover() returns nil outside a panic.
func RecoverNil() any {
	r := recover()
	return fmt.Sprintf("%v", r == nil)
}

// ============================================================================
// Float edge cases (IEEE 754)
// ============================================================================

// NaNLess tests NaN < NaN (should be false).
func NaNLess() any {
	nan := math.NaN()
	return nan < nan
}

// NaNLessEq tests NaN <= NaN (should be false).
func NaNLessEq() any {
	nan := math.NaN()
	return nan <= nan
}

// NaNGreater tests NaN > NaN (should be false).
func NaNGreater() any {
	nan := math.NaN()
	return nan > nan
}

// NaNGreaterEq tests NaN >= NaN (should be false).
func NaNGreaterEq() any {
	nan := math.NaN()
	return nan >= nan
}

// NaNEq tests NaN == NaN (should be false).
func NaNEq() any {
	nan := math.NaN()
	return nan == nan
}

// NaNNeq tests NaN != NaN (should be true).
func NaNNeq() any {
	nan := math.NaN()
	return nan != nan
}

// NaNLessThanNumber tests NaN < 1.0 (should be false).
func NaNLessThanNumber() any {
	nan := math.NaN()
	return nan < 1.0
}

// NaNLessEqNumber tests NaN <= 1.0 (should be false).
func NaNLessEqNumber() any {
	nan := math.NaN()
	return nan <= 1.0
}

// NumberLessEqNaN tests 1.0 <= NaN (should be false).
func NumberLessEqNaN() any {
	nan := math.NaN()
	return 1.0 <= nan
}

// InfArithmetic tests Inf behavior.
func InfArithmetic() any {
	inf := math.Inf(1)
	negInf := math.Inf(-1)
	return fmt.Sprintf("inf>0:%v,neginf<0:%v,inf+inf:%v,inf==inf:%v",
		inf > 0, negInf < 0, inf+inf == inf, inf == inf)
}

// InfTimesZero tests Inf * 0 = NaN.
func InfTimesZero() any {
	inf := math.Inf(1)
	result := inf * 0
	return math.IsNaN(result)
}

// NegativeZero tests -0.0 == 0.0.
func NegativeZero() any {
	negZero := math.Copysign(0, -1)
	posZero := 0.0
	return fmt.Sprintf("eq:%v,signbit:%v", negZero == posZero, math.Signbit(negZero))
}

// FloatCompareNaN tests comparing float variables that might be NaN.
func FloatCompareNaN() any {
	x := math.NaN()
	y := 1.0
	// All comparisons with NaN are false except !=
	results := ""
	if x < y {
		results += "lt "
	}
	if x <= y {
		results += "le "
	}
	if x > y {
		results += "gt "
	}
	if x >= y {
		results += "ge "
	}
	if x == y {
		results += "eq "
	}
	if x != y {
		results += "ne "
	}
	return strings.TrimSpace(results)
}

// ============================================================================
// Integer overflow
// ============================================================================

// Uint8Overflow tests uint8 wrapping.
func Uint8Overflow() any {
	var x uint8 = 255
	x = x + 1
	return x
}

// Int8Overflow tests int8 wrapping.
func Int8Overflow() any {
	var x int8 = 127
	x = x + 1
	return x
}

// Int8Underflow tests int8 underflow.
func Int8Underflow() any {
	var x int8 = -128
	x = x - 1
	return x
}

// Uint16Overflow tests uint16 wrapping.
func Uint16Overflow() any {
	var x uint16 = 65535
	x = x + 1
	return x
}

// ============================================================================
// Slice capacity
// ============================================================================

// SliceAppendCap tests that append doubles capacity when full.
func SliceAppendCap() any {
	s := make([]int, 3, 3)
	s[0], s[1], s[2] = 1, 2, 3
	s = append(s, 4)
	return fmt.Sprintf("len=%d,cap=%d", len(s), cap(s))
}

// SliceSubAppend tests append to sub-slice with remaining capacity.
func SliceSubAppend() any {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3]
	sub = append(sub, 99)
	return fmt.Sprintf("s=%v,sub=%v,len=%d,cap=%d", s, sub, len(sub), cap(sub))
}

// SliceSubAppendNoCap tests append to sub-slice when at capacity.
func SliceSubAppendNoCap() any {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:3:3] // cap limited to 3
	sub = append(sub, 99)
	return fmt.Sprintf("s=%v,sub=%v", s, sub)
}
