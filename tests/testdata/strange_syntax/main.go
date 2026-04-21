package strange_syntax

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
	"text/tabwriter"
	"time"
	"unicode/utf8"
)

// ============================================================================
// OPERATOR PRECEDENCE EDGE CASES
// ============================================================================

// StrangePrecedence1 tests complex operator precedence
func StrangePrecedence1() int {
	return 1 + 2*3<<4/8 - 5
}

// StrangePrecedence2 tests bitwise and logical mix
func StrangePrecedence2() bool {
	return 5&3 > 2 || 7^2 < 10 && !false
}

// StrangePrecedence3 tests shift with addition
func StrangePrecedence3() int {
	return 1<<2 + 3 // Should be 1 << (2+3) = 32
}

// StrangePrecedence4 tests channel and comparison
func StrangePrecedence4() bool {
	ch := make(chan int, 1)
	ch <- 1
	return <-ch == 1
}

// StrangePrecedence5 tests complex arithmetic with unary
func StrangePrecedence5() int {
	return -1 + 2*-3 // Should be -1 + (2 * -3) = -7
}

// ============================================================================
// STRANGE SLICE OPERATIONS
// ============================================================================

// SliceBeyondCapacity tests slicing beyond capacity
func SliceBeyondCapacity() int {
	s := make([]int, 5, 10)
	s = s[2:8] // Valid: within capacity
	return len(s) + cap(s)
}

// SliceNegativePattern tests slice with negative indices (via len)
func SliceNegativePattern() int {
	s := []int{1, 2, 3, 4, 5}
	idx := 2
	return s[idx] + s[len(s)-idx-1]
}

// SliceTripleIndex tests three-index slice
func SliceTripleIndex() int {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s2 := s[2:5:7] // len=3, cap=5
	return len(s2) + cap(s2)
}

// SliceAppendToNilWithCapacity tests append to nil slice
func SliceAppendToNilWithCapacity() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// SliceComplexAppend tests complex append chain
func SliceComplexAppend() []int {
	s := []int{1, 2}
	s = append(s, 3, 4)
	s = append(s, []int{5, 6}...)
	return s
}

// SliceModifyDuringRange tests modifying slice during range
func SliceModifyDuringRange() int {
	s := []int{1, 2, 3}
	sum := 0
	for i := range s {
		sum += s[i]
		if i == 1 {
			s = append(s, 4)
		}
	}
	return sum
}

// ============================================================================
// COMPLEX TYPE CONVERSIONS
// ============================================================================

// ConvertComplexChain tests chained type conversions
func ConvertComplexChain() int {
	return int(int64(int32(int16(int8(42)))))
}

// ConvertFloatToInt tests float to int truncation
func ConvertFloatToInt() int {
	f1 := 3.7
	f2 := -2.3
	return int(f1) + int(f2)
}

// ConvertByteToString tests byte to string
func ConvertByteToString() string {
	b := []byte{'h', 'i'}
	return string(b)
}

// ConvertStringToByte tests string to byte
func ConvertStringToByte() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

// ConvertIntPtrToInt tests pointer to value conversion
func ConvertIntPtrToInt() int {
	x := 42
	p := &x
	return int(*p)
}

// ConvertNilToInterface tests nil to interface
func ConvertNilToInterface() interface{} {
	var s []int
	return s // nil slice assigned to interface
}

// ============================================================================
// NESTED EXPRESSIONS
// ============================================================================

// NestedTernaryLike tests nested if-else as ternary
func NestedTernaryLike() int {
	x := 5
	result := 0
	if x > 0 {
		if x < 10 {
			result = 1
		} else {
			result = 2
		}
	} else {
		result = 0
	}
	return result
}

// NestedFunctionCalls tests deeply nested function calls
func NestedFunctionCalls() int {
	return add(mul(sub(10, 3), 2), 5)
}

func add(a, b int) int { return a + b }
func mul(a, b int) int { return a * b }
func sub(a, b int) int { return a - b }

// NestedMapIndex tests nested map indexing
func NestedMapIndex() int {
	m := map[string]map[int]string{
		"outer": {1: "inner"},
	}
	return len(m["outer"][1])
}

// NestedSliceIndex tests nested slice indexing
func NestedSliceIndex() int {
	s := [][]int{{1, 2, 3}, {4, 5, 6}}
	return s[1][2]
}

// NestedStructField tests deeply nested struct fields
func NestedStructField() int {
	type D struct{ val int }
	type C struct{ d D }
	type B struct{ c C }
	type A struct{ b B }
	a := A{b: B{c: C{d: D{val: 42}}}}
	return a.b.c.d.val
}

// ============================================================================
// UNUSUAL CONTROL FLOW
// ============================================================================

// BreakToLabel tests breaking to label
func BreakToLabel() int {
	sum := 0
outer:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i+j > 5 {
				break outer
			}
			sum++
		}
	}
	return sum
}

// ContinueToLabel tests continue to label
func ContinueToLabel() int {
	sum := 0
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue outer
			}
			sum++
		}
	}
	return sum
}

// GotoForward tests forward goto
func GotoForward() int {
	x := 1
	goto skip
	x = 2
skip:
	return x
}

// GotoBackward tests backward goto
func GotoBackward() int {
	sum := 0
	i := 0
start:
	if i >= 5 {
		return sum
	}
	sum += i
	i++
	goto start
}

// SwitchBreakToLabel tests break to label in switch
func SwitchBreakToLabel() int {
	i := 0
outer:
	for {
		switch i {
		case 5:
			break outer
		default:
			i++
		}
	}
	return i
}

// EmptySelect tests empty select block
func EmptySelect() bool {
	select {
	default:
		return true
	}
}

// SelectWithMultipleCases tests select with multiple ready cases
func SelectWithMultipleCases() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2

	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	default:
		return 0
	}
}

// ============================================================================
// COMPLEX MAP OPERATIONS
// ============================================================================

// MapNestedStructKey tests map with struct key
func MapNestedStructKey() int {
	type Key struct{ x, y int }
	m := map[Key]int{
		{1, 2}: 3,
		{4, 5}: 6,
	}
	return m[Key{1, 2}] + m[Key{4, 5}]
}

// MapDeleteDuringRange tests deleting from map during range
func MapDeleteDuringRange() int {
	m := map[int]int{1: 1, 2: 2, 3: 3}
	for k := range m {
		if k == 2 {
			delete(m, k)
		}
	}
	return len(m)
}

// MapUpdateDuringRange tests updating map during range
func MapUpdateDuringRange() int {
	m := map[int]int{1: 1, 2: 2, 3: 3}
	for k, v := range m {
		m[k] = v * 2
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// MapWithNilValue tests map with nil pointer values
func MapWithNilValue() bool {
	m := map[string]*int{"key": nil}
	return m["key"] == nil
}

// MapComplexKeyType tests map with array key
func MapComplexKeyType() int {
	m := map[[2]int]int{
		{1, 2}: 3,
		{4, 5}: 6,
	}
	return m[[2]int{1, 2}]
}

// ============================================================================
// STRANGE CLOSURE PATTERNS
// ============================================================================

// ClosureCaptureBeforeDeclaration tests capturing variable before declaration
func ClosureCaptureBeforeDeclaration() int {
	var f func() int
	x := 1
	f = func() int {
		return x
	}
	x = 2
	return f()
}

// ClosureRecursive tests recursive closure
func ClosureRecursive() int {
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fib(10)
}

// ClosureMultipleCaptures tests closure capturing multiple variables
func ClosureMultipleCaptures() int {
	a, b, c := 1, 2, 3
	f := func() int {
		return a + b + c
	}
	a *= 10
	b *= 10
	c *= 10
	return f()
}

// ClosureInLoop tests closure in loop with capture
func ClosureInLoop() int {
	var fns []func() int
	for i := 0; i < 3; i++ {
		i := i // capture loop variable
		fns = append(fns, func() int { return i })
	}
	return fns[0]() + fns[1]() + fns[2]()
}

// ClosureReturnNil tests returning nil function
func ClosureReturnNil() func() int {
	if false {
		return func() int { return 1 }
	}
	return nil
}

// ============================================================================
// POINTER WEIRDNESS
// ============================================================================

// PointerToPointer tests double pointer
func PointerToPointer() int {
	x := 42
	p := &x
	pp := &p
	return **pp
}

// PointerToSliceElement tests pointer to slice element
func PointerToSliceElement() int {
	s := []int{1, 2, 3}
	p := &s[1]
	*p = 20
	return s[1]
}

// PointerToArrayElement tests pointer to array element
func PointerToArrayElement() int {
	a := [3]int{1, 2, 3}
	p := &a[1]
	*p = 20
	return a[1]
}

// NilPointerDereferenceGuard tests nil pointer check
func NilPointerDereferenceGuard() int {
	var p *int
	if p != nil {
		return *p
	}
	return 0
}

// PointerToMapValue tests pointer to map value (not allowed, test error handling)
func PointerToMapValue() int {
	m := map[int]int{1: 10}
	// Can't take address of map element, so just return value
	return m[1]
}

// PointerArithmetic tests pointer-like arithmetic (via index)
func PointerArithmetic() int {
	s := []int{1, 2, 3, 4, 5}
	offset := 2
	return s[offset] // Simulating pointer arithmetic
}

// ============================================================================
// MULTIPLE RETURN VALUE EDGE CASES
// ============================================================================

// MultipleReturnIgnore tests ignoring multiple return values
func MultipleReturnIgnore() int {
	a, _ := multiReturn()
	return a
}

func multiReturn() (int, int) {
	return 10, 20
}

// MultipleReturnChain tests chaining multiple returns
func MultipleReturnChain() int {
	return addBoth(multiReturn())
}

func addBoth(a, b int) int {
	return a + b
}

// MultipleReturnToSlice tests multiple return to slice
func MultipleReturnToSlice() []int {
	a, b := multiReturn()
	return []int{a, b}
}

// NamedReturnShadow tests named return shadowing
func NamedReturnShadow() (result int) {
	result = 1
	if true {
		result := 2 // shadows named return
		return result
	}
	return
}

// MultipleReturnInClosure tests multiple return in closure
func MultipleReturnInClosure() int {
	f := func() (int, int) {
		return 5, 10
	}
	a, b := f()
	return a + b
}

// ============================================================================
// DEFER EDGE CASES
// ============================================================================

// DeferMultiple tests multiple defers (LIFO order)
func DeferMultiple() int {
	result := 0
	defer func() { result += 1 }()
	defer func() { result += 2 }()
	defer func() { result += 4 }()
	result = 8
	return result
}

// DeferInLoop tests defer in loop
func DeferInLoop() int {
	result := 0
	for i := 0; i < 3; i++ {
		defer func(x int) { result += x }(i)
	}
	return result // 2+1+0 = 3
}

// DeferModifyReturn tests defer modifying named return
func DeferModifyReturn() (result int) {
	defer func() { result *= 2 }()
	return 5 // Returns 10 after defer
}

// DeferClosureCapture tests defer closure capture
func DeferClosureCapture() int {
	x := 1
	defer func() {
		x = 2 // Modifies outer x but after return
	}()
	return x // Returns 1
}

// DeferArguments tests defer argument evaluation
func DeferArguments() int {
	x := 1
	defer func(val int) {
		// val is evaluated at defer call
		_ = val
	}(x)
	x = 2
	return x
}

// ============================================================================
// STRUCT EMBEDDING EDGE CASES
// ============================================================================

// StructEmbed tests struct embedding
func StructEmbed() int {
	type Inner struct{ value int }
	type Outer struct {
		Inner
		extra int
	}
	o := Outer{Inner: Inner{value: 10}, extra: 5}
	return o.value + o.extra
}

// StructEmbedInterface tests embedded interface
func StructEmbedInterface() int {
	type Printer interface{ Print() int }
	type Container struct {
		Printer
		value int
	}
	// Can't instantiate without implementation
	return 0
}

// StructPointerEmbed tests pointer embedding
func StructPointerEmbed() int {
	type Inner struct{ value int }
	type Outer struct {
		*Inner
	}
	o := Outer{Inner: &Inner{value: 42}}
	return o.value
}

// StructMultipleEmbed tests multiple embedding
func StructMultipleEmbed() int {
	type A struct{ a int }
	type B struct{ b int }
	type C struct {
		A
		B
		c int
	}
	obj := C{A: A{a: 1}, B: B{b: 2}, c: 3}
	return obj.a + obj.b + obj.c
}

// ============================================================================
// CHANNEL EDGE CASES
// ============================================================================

// ChannelNilSend tests nil channel send (blocks forever, avoid)
func ChannelNilSend() int {
	// Don't actually send to nil channel - would block
	var ch chan int
	if ch == nil {
		return 1
	}
	return 0
}

// ChannelNilReceive tests nil channel receive (blocks forever, avoid)
func ChannelNilReceive() int {
	var ch chan int
	if ch == nil {
		return 1
	}
	return 0
}

// ChannelClosedSend tests sending to closed channel
func ChannelClosedSend() int {
	ch := make(chan int, 1)
	close(ch)
	defer func() {
		if r := recover(); r != nil {
			// Recovered from panic
		}
	}()
	ch <- 1 // Will panic
	return 0
}

// ChannelClosedReceive tests receiving from closed channel
func ChannelClosedReceive() (int, bool) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v, ok := <-ch
	return v, ok
}

// ChannelBufferedClose tests buffered channel close
func ChannelBufferedClose() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// ============================================================================
// TYPE ASSERTION EDGE CASES
// ============================================================================

// TypeAssertionSuccess tests successful type assertion
func TypeAssertionSuccess() int {
	var i interface{} = 42
	v, ok := i.(int)
	if ok {
		return v
	}
	return 0
}

// TypeAssertionFailure tests failed type assertion
func TypeAssertionFailure() int {
	var i interface{} = "hello"
	v, ok := i.(int)
	if ok {
		return v
	}
	return -1
}

// TypeAssertionPanic tests type assertion panic without comma ok
func TypeAssertionPanic() int {
	defer func() {
		if r := recover(); r != nil {
			// Recovered
		}
	}()
	var i interface{} = "hello"
	return i.(int) // Will panic
}

// TypeAssertionNil tests type assertion on nil
func TypeAssertionNil() int {
	var i interface{}
	_, ok := i.(int)
	if ok {
		return 1
	}
	return 0
}

// TypeSwitch tests type switch
func TypeSwitch() string {
	var i interface{} = 42
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case string:
		return fmt.Sprintf("string: %s", v)
	default:
		return fmt.Sprintf("unknown: %v", v)
	}
}

// ============================================================================
// NIL HANDLING EDGE CASES
// ============================================================================

// NilSliceAppend tests nil slice append
func NilSliceAppend() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// NilMapAccess tests nil map access
func NilMapAccess() int {
	var m map[string]int
	return m["key"] // Returns zero value
}

// NilMapDelete tests nil map delete
func NilMapDelete() int {
	var m map[string]int
	delete(m, "key") // No-op
	return 0
}

// NilMapLen tests nil map length
func NilMapLen() int {
	var m map[string]int
	return len(m)
}

// NilSliceLen tests nil slice length
func NilSliceLen() int {
	var s []int
	return len(s)
}

// NilSliceCap tests nil slice capacity
func NilSliceCap() int {
	var s []int
	return cap(s)
}

// NilInterfaceComparison tests nil interface comparison
func NilInterfaceComparison() bool {
	var i interface{}
	return i == nil
}

// NilTypedInterface tests nil typed interface
func NilTypedInterface() bool {
	var err error
	return err == nil
}

// ============================================================================
// SHADOWING EDGE CASES
// ============================================================================

// VariableShadowing tests variable shadowing in inner scope
func VariableShadowing() int {
	x := 1
	{
		x := 2
		_ = x
	}
	return x // Returns 1
}

// ParameterShadowing tests parameter shadowing
func ParameterShadowing(paramX int) int {
	x := 2 // Shadows parameter
	return x
}

// ReturnShadowing tests return variable shadowing
func ReturnShadowing() (x int) {
	x = 1
	{
		x := 2
		_ = x
	}
	return x // Returns 1
}

// ImportShadowing tests shadowing imported names
func ImportShadowing() int {
	// Shadowing strings package with variable
	strings := "hello"
	return len(strings)
}

// ============================================================================
// METHOD EXPRESSION EDGE CASES
// ============================================================================

// MethodExpression tests method expression
func MethodExpression() int {
	type MyInt int
	m := func(m MyInt) int {
		return int(m) * 2
	}
	var x MyInt = 21
	return m(x)
}

// MethodValue tests method value
func MethodValue() int {
	type Counter struct{ value int }
	c := &Counter{value: 10}
	inc := func() {
		c.value++
	}
	inc()
	inc()
	return c.value
}

// ============================================================================
// BLANK IDENTIFIER EDGE CASES
// ============================================================================

// BlankIdentifierAssignment tests blank identifier assignment
func BlankIdentifierAssignment() int {
	_, b, _ := 1, 2, 3
	return b
}

// BlankIdentifierImport tests blank identifier import (strconv imported at top)
func BlankIdentifierImport() int {
	return 1
}

// BlankIdentifierRange tests blank identifier in range
func BlankIdentifierRange() int {
	s := []int{1, 2, 3}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// BlankIdentifierReturn tests blank identifier in return
func BlankIdentifierReturn() int {
	_, err := returnsError()
	if err != nil {
		return -1
	}
	return 0
}

func returnsError() (int, error) {
	return 42, nil
}

// ============================================================================
// COMPLEX COMPOSITE LITERALS
// ============================================================================

// ComplexSliceLiteral tests complex slice literal
func ComplexSliceLiteral() int {
	s := []int{1, 2, 3: 10, 4: 20, 5}
	return s[3] + s[4] + s[5]
}

// ComplexMapLiteral tests complex map literal
func ComplexMapLiteral() int {
	m := map[int]int{
		1: 10,
		2: 20,
		3: 30,
	}
	return m[1] + m[2] + m[3]
}

// NestedCompositeLiteral tests nested composite literal
func NestedCompositeLiteral() int {
	type Inner struct{ x, y int }
	type Outer struct {
		inners []Inner
	}
	o := Outer{
		inners: []Inner{
			{1, 2},
			{3, 4},
		},
	}
	return o.inners[0].x + o.inners[1].y
}

// PointerCompositeLiteral tests pointer composite literal
func PointerCompositeLiteral() int {
	type Point struct{ x, y int }
	p := &Point{1, 2}
	return p.x + p.y
}

// ============================================================================
// STRING EDGE CASES
// ============================================================================

// StringIndex tests string indexing
func StringIndex() byte {
	s := "hello"
	return s[1]
}

// StringSlice tests string slicing
func StringSlice() string {
	s := "hello world"
	return s[0:5]
}

// StringRange tests ranging over string
func StringRange() int {
	s := "hello"
	count := 0
	for range s {
		count++
	}
	return count
}

// StringConcat tests string concatenation
func StringConcat() string {
	s1 := "hello"
	s2 := " world"
	return s1 + s2
}

// StringCompare tests string comparison
func StringCompare() bool {
	s1 := "hello"
	s2 := "hello"
	s3 := "world"
	return s1 == s2 && s1 != s3
}

// MultilineString tests multiline string literal
func MultilineString() int {
	s := `line1
line2
line3`
	return len(s)
}

// RawStringLiteral tests raw string literal
func RawStringLiteral() string {
	return `hello\nworld` // Literal backslash n
}

// InterpretedStringLiteral tests interpreted string literal
func InterpretedStringLiteral() string {
	return "hello\nworld" // Actual newline
}

// ============================================================================
// ARRAY EDGE CASES
// ============================================================================

// ArrayLiteral tests array literal
func ArrayLiteral() int {
	arr := [5]int{1, 2, 3, 4, 5}
	return arr[0] + arr[4]
}

// ArrayPartialInit tests array partial initialization
func ArrayPartialInit() int {
	arr := [5]int{1, 2}
	return arr[2] // Should be 0
}

// ArrayIndexExpression tests array index with expression
func ArrayIndexExpression() int {
	arr := [5]int{10, 20, 30, 40, 50}
	idx := 2
	return arr[idx+1]
}

// ArrayPointer tests array pointer
func ArrayPointer() int {
	arr := [3]int{1, 2, 3}
	p := &arr
	return (*p)[1]
}

// ArrayComparison tests array comparison
func ArrayComparison() bool {
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	return arr1 == arr2
}

// ============================================================================
// INTERFACE EDGE CASES
// ============================================================================

// InterfaceNil tests nil interface
func InterfaceNil() interface{} {
	return nil
}

// InterfaceConcrete tests interface with concrete type
func InterfaceConcrete() int {
	var i interface{} = 42
	return i.(int)
}

// InterfaceSlice tests interface slice
func InterfaceSlice() int {
	var s []interface{}
	s = append(s, 1, "hello", 3.14)
	return len(s)
}

// InterfaceMap tests interface map
func InterfaceMap() int {
	m := make(map[string]interface{})
	m["int"] = 42
	m["str"] = "hello"
	return len(m)
}

// EmptyInterface tests empty interface
func EmptyInterface() interface{} {
	var e interface{}
	e = 42
	return e
}

// ============================================================================
// COMPARISON EDGE CASES
// ============================================================================

// CompareDifferentTypes tests comparing different types via interface
func CompareDifferentTypes() bool {
	var i1 interface{} = 42
	var i2 interface{} = "hello"
	return i1 == i2
}

// CompareNilInterface tests comparing nil interface
func CompareNilInterface() bool {
	var i interface{}
	return i == nil
}

// CompareFunc tests function comparison (always false unless nil)
func CompareFunc() bool {
	f1 := func() {}
	f2 := func() {}
	return f1 == nil && f2 == nil // Both nil
}

// CompareMap tests map comparison (always false unless nil)
func CompareMap() bool {
	m1 := map[int]int{1: 2}
	m2 := map[int]int{1: 2}
	return m1 != nil && m2 != nil
}

// CompareSlice tests slice comparison (always false unless nil)
func CompareSlice() bool {
	s1 := []int{1, 2}
	s2 := []int{1, 2}
	return s1 != nil && s2 != nil
}

// ============================================================================
// BITWISE EDGE CASES
// ============================================================================

// BitwiseAnd tests bitwise AND
func BitwiseAnd() int {
	return 0xFF & 0x0F
}

// BitwiseOr tests bitwise OR
func BitwiseOr() int {
	return 0xF0 | 0x0F
}

// BitwiseXor tests bitwise XOR
func BitwiseXor() int {
	return 0xFF ^ 0x0F
}

// BitwiseNot tests bitwise NOT
func BitwiseNot() int {
	return ^0x0F
}

// BitwiseLeftShift tests left shift
func BitwiseLeftShift() int {
	return 1 << 4
}

// BitwiseRightShift tests right shift
func BitwiseRightShift() int {
	return 16 >> 2
}

// BitwiseComplex tests complex bitwise expression
func BitwiseComplex() int {
	return (0xAA & 0x55) | (0xF0 ^ 0x0F)
}

// ============================================================================
// FLOATING POINT EDGE CASES
// ============================================================================

// FloatNaN tests NaN comparison
func FloatNaN() bool {
	x := 0.0
	y := 0.0
	nan := x / y
	return nan != nan // NaN != NaN is true
}

// FloatInf tests infinity
func FloatInf() bool {
	x := 1.0
	y := 0.0
	inf := x / y
	return inf > 1e308
}

// FloatNegativeInf tests negative infinity
func FloatNegativeInf() bool {
	x := -1.0
	y := 0.0
	negInf := x / y
	return negInf < -1e308
}

// FloatZeroDivision tests division by zero
func FloatZeroDivision() bool {
	x := 1.0
	y := 0.0
	result := x / y
	return result > 0 // Should be +Inf
}

// FloatPrecision tests floating point precision
func FloatPrecision() float64 {
	return 0.1 + 0.2
}

// ============================================================================
// UNARY OPERATOR EDGE CASES
// ============================================================================

// UnaryPlus tests unary plus
func UnaryPlus() int {
	return +42
}

// UnaryMinus tests unary minus
func UnaryMinus() int {
	return -42
}

// UnaryNot tests logical NOT
func UnaryNot() bool {
	return !true
}

// UnaryXor tests bitwise complement
func UnaryXor() int {
	return ^0
}

// UnaryComplex tests complex unary expression
func UnaryComplex() int {
	return -(-42)
}

// ============================================================================
// ASSIGNMENT EDGE CASES
// ============================================================================

// AssignMultiple tests multiple assignment
func AssignMultiple() int {
	a, b, c := 1, 2, 3
	return a + b + c
}

// AssignSwap tests swap via multiple assignment
func AssignSwap() (int, int) {
	a, b := 1, 2
	a, b = b, a
	return a, b
}

// AssignComplex tests complex assignment
func AssignComplex() int {
	x := 10
	x += 5
	x *= 2
	x /= 3
	return x
}

// AssignOperator tests assignment operators
func AssignOperator() int {
	x := 10
	x &= 7
	x |= 8
	x ^= 15
	return x
}

// ============================================================================
// CONSTANTS EDGE CASES
// ============================================================================

const (
	ConstA = iota
	ConstB
	ConstC
)

// IotaUsage tests iota usage
func IotaUsage() int {
	return ConstA + ConstB + ConstC
}

// ConstExpression tests constant expression
func ConstExpression() int {
	const x = 10 + 20*2
	return x
}

// ConstUntyped tests untyped constant
func ConstUntyped() int {
	const x = 42
	return x
}

// ConstTyped tests typed constant
func ConstTyped() int {
	const x int = 42
	return x
}

// ============================================================================
// RANGE EDGE CASES
// ============================================================================

// RangeOverMap tests range over map
func RangeOverMap() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// RangeOverString tests range over string
func RangeOverString() int {
	s := "hello"
	sum := 0
	for i := range s {
		sum += i
	}
	return sum
}

// RangeOverChannel tests range over channel
func RangeOverChannel() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// RangeWithBreak tests range with break
func RangeWithBreak() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v > 3 {
			break
		}
		sum += v
	}
	return sum
}

// RangeWithContinue tests range with continue
func RangeWithContinue() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v%2 == 0 {
			continue
		}
		sum += v
	}
	return sum
}

// ============================================================================
// MISCELLANEOUS EDGE CASES
// ============================================================================

// ShortVariableDeclaration tests short variable declaration
func ShortVariableDeclaration() int {
	x := 10
	y := 20
	z := x + y
	return z
}

// RedeclarationInDifferentScope tests redeclaration in different scope
func RedeclarationInDifferentScope() int {
	x := 1
	{
		x := 2
		x++
	}
	return x
}

// BlankExpression tests blank expression
func BlankExpression() {
	_ = 42
}

// MultipleBlankAssignments tests multiple blank assignments
func MultipleBlankAssignments() int {
	_, _, c := 1, 2, 3
	return c
}

// StringContains tests string contains operation
func StringContains() bool {
	return strings.Contains("hello world", "world")
}

// StringHasPrefix tests string prefix check
func StringHasPrefix() bool {
	return strings.HasPrefix("hello", "hel")
}

// StringHasSuffix tests string suffix check
func StringHasSuffix() bool {
	return strings.HasSuffix("hello", "llo")
}

// StringSplit tests string split
func StringSplit() int {
	parts := strings.Split("a,b,c", ",")
	return len(parts)
}

// StringJoin tests string join
func StringJoin() string {
	parts := []string{"a", "b", "c"}
	return strings.Join(parts, ",")
}

// StringToUpper tests string to upper
func StringToUpper() string {
	return strings.ToUpper("hello")
}

// StringToLower tests string to lower
func StringToLower() string {
	return strings.ToLower("HELLO")
}

// StringTrim tests string trim
func StringTrim() string {
	return strings.Trim("  hello  ", " ")
}

// StringReplace tests string replace
func StringReplace() string {
	return strings.Replace("hello world", "world", "golang", 1)
}

// StringCount tests string count
func StringCount() int {
	return strings.Count("hello hello", "hello")
}

// StringRepeat tests string repeat
func StringRepeat() string {
	return strings.Repeat("ab", 3)
}

// ComplexExpressions tests multiple complex expressions combined
func ComplexExpressions() int {
	type Point struct{ x, y int }
	points := []Point{{1, 2}, {3, 4}}
	m := map[int]Point{}
	for i, p := range points {
		m[i] = p
	}

	sum := 0
	for k, v := range m {
		sum += k + v.x + v.y
	}

	f := func(p Point) int {
		return p.x * p.y
	}

	return sum + f(m[1])
}

// NestedSlices tests nested slices
func NestedSlices() int {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	sum := 0
	for _, row := range matrix {
		for _, val := range row {
			sum += val
		}
	}
	return sum
}

// NestedMaps tests nested maps
func NestedMaps() int {
	data := map[string]map[string]int{
		"group1": {"a": 1, "b": 2},
		"group2": {"c": 3, "d": 4},
	}
	sum := 0
	for _, inner := range data {
		for _, val := range inner {
			sum += val
		}
	}
	return sum
}

// ComplexClosureChain tests complex closure chaining
func ComplexClosureChain() int {
	makeAdder := func(x int) func(int) func(int) int {
		return func(y int) func(int) int {
			return func(z int) int {
				return x + y + z
			}
		}
	}
	return makeAdder(10)(20)(30)
}

// RecursiveStruct tests recursive struct definition
func RecursiveStruct() int {
	type Node struct {
		value int
		next  *Node
	}
	n3 := &Node{value: 3}
	n2 := &Node{value: 2, next: n3}
	n1 := &Node{value: 1, next: n2}

	sum := 0
	for current := n1; current != nil; current = current.next {
		sum += current.value
	}
	return sum
}

// InterfaceMethodCall tests interface method call
func InterfaceMethodCall() int {
	type Stringer interface {
		String() string
	}

	type MyInt int
	var x MyInt = 42

	// Convert to interface and call method
	return int(x)
}

// ============================================================================
// MORE EDGE CASES TO DISCOVER BUGS
// ============================================================================

// NilSliceCopy tests copying nil slice
func NilSliceCopy() int {
	var src []int
	dst := make([]int, 0)
	copy(dst, src)
	return len(dst)
}

// NilMapRange tests ranging over nil map
func NilMapRange() int {
	var m map[int]int
	count := 0
	for range m {
		count++
	}
	return count
}

// NilSliceRange tests ranging over nil slice
func NilSliceRange() int {
	var s []int
	count := 0
	for range s {
		count++
	}
	return count
}

// NilChannelRange tests ranging over nil channel - should block forever
// so we just test that nil channel is nil and skip the range
func NilChannelRange() int {
	var ch chan int
	// Ranging over nil channel blocks forever, so we just check it's nil
	if ch == nil {
		return 0
	}
	count := 0
	for range ch {
		count++
	}
	return count
}

// SliceLenCap tests len and cap on various types
func SliceLenCap() int {
	s := make([]int, 5, 10)
	return len(s)*100 + cap(s)
}

// MapLen tests len on map
func MapLen() int {
	m := map[int]int{1: 2, 3: 4, 5: 6}
	return len(m)
}

// StringLen tests len on string
func StringLen() int {
	s := "hello world"
	return len(s)
}

// ChannelLen tests len on channel
func ChannelLen() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	return len(ch)
}

// ComplexNilCheck tests complex nil checks
func ComplexNilCheck() bool {
	var s []int
	var m map[int]int
	var f func()
	var ch chan int
	var iface interface{}
	return s == nil && m == nil && f == nil && ch == nil && iface == nil
}

// TypedNilNotEqualNil tests typed nil != nil for interface
func TypedNilNotEqualNil() bool {
	var s []int = nil
	var iface interface{} = s
	return iface != nil // typed nil inside interface is not nil
}

// PointerToNilSlice tests pointer to nil slice
func PointerToNilSlice() bool {
	var s []int
	p := &s
	return *p == nil
}

// PointerToNilMap tests pointer to nil map
func PointerToNilMap() bool {
	var m map[int]int
	p := &m
	return *p == nil
}

// EmptySliceVsNil tests empty slice vs nil slice
func EmptySliceVsNil() (bool, bool) {
	var nilSlice []int
	emptySlice := []int{}
	return nilSlice == nil, emptySlice == nil
}

// EmptyMapVsNil tests empty map vs nil map
func EmptyMapVsNil() (bool, bool) {
	var nilMap map[int]int
	emptyMap := make(map[int]int)
	return nilMap == nil, emptyMap == nil
}

// SliceAppendNil tests append to nil slice
func SliceAppendNil() int {
	var s []int
	s = append(s, 1)
	s = append(s, 2, 3)
	return len(s)
}

// MapAssignNil tests assigning nil to map variable
func MapAssignNil() bool {
	m := map[int]int{1: 2}
	m = nil
	return m == nil
}

// SliceAssignNil tests assigning nil to slice variable
func SliceAssignNil() bool {
	s := []int{1, 2, 3}
	s = nil
	return s == nil
}

// ComplexDeferOrder tests complex defer ordering
func ComplexDeferOrder() int {
	result := 0
	defer func() { result += 1 }()
	defer func() { result += 2 }()
	defer func() { result += 4 }()
	result = 8
	defer func() { result += 8 }()
	return result // 8 + 8 + 4 + 2 + 1 = 23
}

// DeferInDefer tests defer inside defer
func DeferInDefer() int {
	result := 0
	defer func() {
		defer func() {
			result += 1
		}()
		result += 2
	}()
	result = 10
	return result // 10 + 2 + 1 = 13
}

// MultipleReturnToInterface tests multiple return wrapped in interface
func MultipleReturnToInterface() interface{} {
	a, b := multiReturnToInterface()
	return []interface{}{a, b}
}

func multiReturnToInterface() (int, string) {
	return 42, "hello"
}

// InterfaceSliceLiteral tests interface slice literal
func InterfaceSliceLiteral() int {
	s := []interface{}{1, "hello", 3.14, true, nil}
	return len(s)
}

// InterfaceMapLiteral tests interface map literal
func InterfaceMapLiteral() int {
	m := map[string]interface{}{
		"int":    42,
		"string": "hello",
		"float":  3.14,
		"nil":    nil,
	}
	return len(m)
}

// StructWithSliceField tests struct with slice field
func StructWithSliceField() int {
	type Container struct {
		items []int
	}
	c := Container{items: []int{1, 2, 3}}
	c.items = append(c.items, 4)
	return len(c.items)
}

// StructWithMapField tests struct with map field
func StructWithMapField() int {
	type Container struct {
		items map[int]string
	}
	c := Container{items: map[int]string{1: "a", 2: "b"}}
	c.items[3] = "c"
	return len(c.items)
}

// StructWithChannelField tests struct with channel field
func StructWithChannelField() int {
	type Container struct {
		ch chan int
	}
	c := Container{ch: make(chan int, 2)}
	c.ch <- 1
	c.ch <- 2
	return len(c.ch)
}

// StructWithFuncField tests struct with function field
func StructWithFuncField() int {
	type Container struct {
		fn func(int) int
	}
	c := Container{fn: func(x int) int { return x * 2 }}
	return c.fn(21)
}

// NestedStructWithPointers tests nested struct with pointers
func NestedStructWithPointers() int {
	type Inner struct{ value int }
	type Outer struct {
		inner *Inner
		next  *Outer
	}
	i := &Inner{value: 42}
	o := &Outer{inner: i}
	o2 := &Outer{inner: &Inner{value: 10}, next: o}
	return o2.inner.value + o2.next.inner.value
}

// SliceOfPointers tests slice of pointers
func SliceOfPointers() int {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}
	sum := 0
	for _, p := range s {
		sum += *p
	}
	return sum
}

// MapOfPointers tests map of pointers
func MapOfPointers() int {
	a, b := 1, 2
	m := map[string]*int{"a": &a, "b": &b}
	sum := 0
	for _, p := range m {
		if p != nil {
			sum += *p
		}
	}
	return sum
}

// SliceOfSlices tests slice of slices
func SliceOfSlices() int {
	s := [][]int{
		{1, 2},
		{3, 4, 5},
		{6},
	}
	return len(s) + len(s[0]) + len(s[1]) + len(s[2])
}

// MapOfMaps tests map of maps
func MapOfMaps() int {
	m := map[string]map[int]string{
		"first":  {1: "a", 2: "b"},
		"second": {3: "c"},
	}
	return len(m) + len(m["first"]) + len(m["second"])
}

// SliceOfMaps tests slice of maps
func SliceOfMaps() int {
	s := []map[int]string{
		{1: "a"},
		{2: "b", 3: "c"},
	}
	return len(s) + len(s[0]) + len(s[1])
}

// MapOfSlices tests map of slices
func MapOfSlices() int {
	m := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5},
	}
	return len(m) + len(m["a"]) + len(m["b"])
}

// ComplexInterfaceAssertion tests complex interface assertions
func ComplexInterfaceAssertion() int {
	var iface interface{} = []int{1, 2, 3}
	if s, ok := iface.([]int); ok {
		return len(s)
	}
	return 0
}

// InterfaceAssertionWithNil tests interface assertion with nil
func InterfaceAssertionWithNil() bool {
	var iface interface{}
	_, ok := iface.(int)
	return ok
}

// TypeSwitchWithNil tests type switch with nil
func TypeSwitchWithNil() string {
	var iface interface{}
	switch v := iface.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("other: %T", v)
	}
}

// PointerToPointerToStruct tests pointer to pointer to struct
func PointerToPointerToStruct() int {
	type Data struct{ value int }
	d := &Data{value: 42}
	pp := &d
	return (*pp).value
}

// MultiplePointerDereference tests multiple pointer dereferences
func MultiplePointerDereference() int {
	x := 42
	p1 := &x
	p2 := &p1
	p3 := &p2
	return ***p3
}

// SliceWithNilElements tests slice with nil elements
func SliceWithNilElements() int {
	s := []*int{nil, new(int), nil}
	count := 0
	for _, p := range s {
		if p == nil {
			count++
		}
	}
	return count
}

// MapWithNilValues tests map with nil values
func MapWithNilValues() int {
	m := map[int]*int{
		1: nil,
		2: new(int),
		3: nil,
	}
	count := 0
	for _, v := range m {
		if v == nil {
			count++
		}
	}
	return count
}

// EmptyStructAsMapValue tests empty struct as map value
func EmptyStructAsMapValue() int {
	type Empty struct{}
	m := map[int]Empty{
		1: {},
		2: {},
	}
	return len(m)
}

// EmptyStructInSlice tests empty struct in slice
func EmptyStructInSlice() int {
	type Empty struct{}
	s := []Empty{{}, {}, {}}
	return len(s)
}

// FunctionReturningNil tests function returning nil
func FunctionReturningNil() func() int {
	return nilFunction()
}

func nilFunction() func() int {
	return nil
}

// ChannelOfChannels tests channel of channels
func ChannelOfChannels() int {
	ch := make(chan chan int, 1)
	inner := make(chan int, 1)
	inner <- 42
	ch <- inner
	return <-(<-ch)
}

// SliceOfChannels tests slice of channels
func SliceOfChannels() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	s := []chan int{ch1, ch2}
	return <-s[0] + <-s[1]
}

// MapOfChannels tests map of channels
func MapOfChannels() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	m := map[string]chan int{"a": ch1, "b": ch2}
	return <-m["a"] + <-m["b"]
}

// ComplexCompositeLiteral tests complex composite literal
func ComplexCompositeLiteral() int {
	type Inner struct {
		values []int
	}
	type Outer struct {
		inners []Inner
	}
	o := Outer{
		inners: []Inner{
			{values: []int{1, 2, 3}},
			{values: []int{4, 5}},
		},
	}
	sum := 0
	for _, inner := range o.inners {
		for _, v := range inner.values {
			sum += v
		}
	}
	return sum
}

// VariadicFunction tests variadic function
func VariadicFunction() int {
	return sumAll(1, 2, 3, 4, 5)
}

func sumAll(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

// VariadicFunctionWithSlice tests variadic function with slice
func VariadicFunctionWithSlice() int {
	nums := []int{1, 2, 3, 4, 5}
	return sumAll(nums...)
}

// VariadicFunctionEmpty tests variadic function with no args
func VariadicFunctionEmpty() int {
	return sumAll()
}

// VariadicInterface tests variadic interface
func VariadicInterface() int {
	return countInterfaces(1, "hello", 3.14, true)
}

func countInterfaces(args ...interface{}) int {
	return len(args)
}

// StructWithVariadicMethod tests struct with variadic method
func StructWithVariadicMethod() int {
	type Adder struct{}
	add := func(a Adder, nums ...int) int {
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return sum
	}
	return add(Adder{}, 1, 2, 3)
}

// ClosureWithVariadic tests closure with variadic
func ClosureWithVariadic() int {
	makeSummer := func(nums ...int) func() int {
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return func() int { return sum }
	}
	return makeSummer(1, 2, 3, 4, 5)()
}

// ============================================================================
// MORE EDGE CASES TO DISCOVER BUGS (Round 2)
// ============================================================================

// TypeAliasBasic tests type alias for basic types
func TypeAliasBasic() int {
	type MyInt = int
	var x MyInt = 42
	return x + 1
}

// TypeAliasStruct tests type alias for struct
func TypeAliasStruct() int {
	type Point struct{ x, y int }
	type P = Point
	p := P{x: 1, y: 2}
	return p.x + p.y
}

// TypeAliasPointer tests type alias for pointer
func TypeAliasPointer() int {
	type IntPtr = *int
	x := 42
	var p IntPtr = &x
	return *p
}

// NamedTypeMethod tests method on named type
func NamedTypeMethod() int {
	type Counter int
	var c Counter
	c = 10
	return int(c) * 2
}

// NamedTypeWithMethods tests named type with multiple methods
func NamedTypeWithMethods() string {
	type Stringer struct{ val string }
	return Stringer{val: "hello"}.val
}

// StructWithAnonymousFields tests struct with anonymous fields
func StructWithAnonymousFields() int {
	type Inner struct{ x int }
	type Outer struct {
		Inner
		y int
	}
	o := Outer{Inner: Inner{x: 10}, y: 20}
	return o.x + o.y
}

// StructWithEmbeddedPointer tests struct with embedded pointer
func StructWithEmbeddedPointer() int {
	type Inner struct{ val int }
	type Outer struct {
		*Inner
	}
	o := Outer{Inner: &Inner{val: 42}}
	return o.val
}

// StructWithMultipleEmbedded tests struct with multiple embedded types
func StructWithMultipleEmbedded() int {
	type A struct{ a int }
	type B struct{ b int }
	type C struct{ c int }
	type D struct {
		A
		B
		C
	}
	d := D{A: A{a: 1}, B: B{b: 2}, C: C{c: 3}}
	return d.a + d.b + d.c
}

// PointerToStructLiteral tests pointer to struct literal
func PointerToStructLiteral() int {
	type Point struct{ x, y int }
	p := &Point{x: 1, y: 2}
	return p.x + p.y
}

// ArrayOfPointers tests array of pointers
func ArrayOfPointers() int {
	a, b, c := 1, 2, 3
	arr := [3]*int{&a, &b, &c}
	sum := 0
	for _, p := range arr {
		sum += *p
	}
	return sum
}

// SliceOfArrays tests slice of arrays
func SliceOfArrays() int {
	s := [][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	sum := 0
	for _, arr := range s {
		for _, v := range arr {
			sum += v
		}
	}
	return sum
}

// MapWithArrayKey tests map with array key
func MapWithArrayKey() int {
	m := map[[2]int]string{
		{1, 2}: "a",
		{3, 4}: "b",
	}
	return len(m)
}

// MapWithStructKey tests map with struct key
func MapWithStructKey() int {
	type Key struct{ x, y int }
	m := map[Key]int{
		{1, 2}: 3,
		{4, 5}: 6,
	}
	return m[Key{1, 2}]
}

// MapWithFuncValue tests map with func value
func MapWithFuncValue() int {
	m := map[string]func() int{
		"a": func() int { return 1 },
		"b": func() int { return 2 },
	}
	return m["a"]() + m["b"]()
}

// ComplexMapKeyWithSlice tests map key with slice field (should fail at compile)
// This tests that we don't crash on complex types
func ComplexMapKeyType() int {
	// Use array instead of slice for valid map key
	m := map[[2]string]int{
		{"a", "b"}: 1,
		{"c", "d"}: 2,
	}
	return len(m)
}

// SliceOfFuncs tests slice of functions
func SliceOfFuncs() int {
	fns := []func() int{
		func() int { return 1 },
		func() int { return 2 },
		func() int { return 3 },
	}
	sum := 0
	for _, fn := range fns {
		sum += fn()
	}
	return sum
}

// ArrayOfFuncs tests array of functions
func ArrayOfFuncs() int {
	var arr [3]func() int
	arr[0] = func() int { return 10 }
	arr[1] = func() int { return 20 }
	arr[2] = func() int { return 30 }
	return arr[0]() + arr[1]() + arr[2]()
}

// FuncReturningFunc tests function returning function
func FuncReturningFunc() int {
	makeAdder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := makeAdder(5)
	return add5(10)
}

// FuncTakingFunc tests function taking function as argument
func FuncTakingFunc() int {
	apply := func(fn func(int) int, x int) int {
		return fn(x)
	}
	double := func(x int) int { return x * 2 }
	return apply(double, 21)
}

// ClosureCapturingLoopVar tests closure capturing loop variable
func ClosureCapturingLoopVar() int {
	var fns []func() int
	for i := 0; i < 3; i++ {
		i := i // capture loop variable
		fns = append(fns, func() int { return i })
	}
	sum := 0
	for _, fn := range fns {
		sum += fn()
	}
	return sum // 0 + 1 + 2 = 3
}

// ClosureCapturingMultipleVars tests closure capturing multiple variables
func ClosureCapturingMultipleVars() int {
	a, b, c := 1, 2, 3
	fn := func() int {
		return a + b + c
	}
	a = 10
	b = 20
	c = 30
	return fn() // 10 + 20 + 30 = 60
}

// NestedClosures tests nested closures
func NestedClosures() int {
	outer := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	fn := outer(21)
	return fn(21)
}

// SelectWithDefault tests select with default
func SelectWithDefault() int {
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		return v
	default:
		return 0
	}
}

// SelectWithNilChannel tests select with nil channel
func SelectWithNilChannel() int {
	var nilCh chan int
	ch := make(chan int, 1)
	ch <- 1
	select {
	case v := <-nilCh:
		return v // never executed
	case v := <-ch:
		return v // returns 1
	}
}

// ChannelOfFuncs tests channel of functions
func ChannelOfFuncs() int {
	ch := make(chan func() int, 1)
	ch <- func() int { return 42 }
	fn := <-ch
	return fn()
}

// ChannelOfInterfaces tests channel of interfaces
func ChannelOfInterfaces() int {
	ch := make(chan interface{}, 2)
	ch <- 42
	ch <- "hello"
	v1 := <-ch
	v2 := <-ch
	return v1.(int) + len(v2.(string))
}

// BufferedChannelWithCap tests buffered channel capacity
func BufferedChannelWithCap() int {
	ch := make(chan int, 5)
	return cap(ch)
}

// ChannelCloseAndRange tests ranging over closed channel
func ChannelCloseAndRange() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

// StringAsByteSlice tests string to byte slice conversion
func StringAsByteSlice() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

// ByteSliceAsString tests byte slice to string conversion
func ByteSliceAsString() int {
	b := []byte{'h', 'e', 'l', 'l', 'o'}
	s := string(b)
	return len(s)
}

// RuneSliceAsString tests rune slice to string conversion
func RuneSliceAsString() int {
	r := []rune{'世', '界'}
	s := string(r)
	return len(s)
}

// StringAsRuneSlice tests string to rune slice conversion
func StringAsRuneSlice() int {
	s := "世界"
	r := []rune(s)
	return len(r)
}

// ComplexNumberLiteral tests complex number literals
func ComplexNumberLiteral() complex128 {
	c := 1 + 2i
	return c
}

// ComplexNumberOperations tests complex number operations
func ComplexNumberOperations() float64 {
	c1 := 1 + 2i
	c2 := 3 + 4i
	c3 := c1 + c2
	return real(c3) + imag(c3)
}

// ComplexNumberFunc tests function with complex return
func ComplexNumberFunc() float64 {
	multiply := func(c1, c2 complex128) complex128 {
		return c1 * c2
	}
	result := multiply(1+2i, 3+4i)
	return real(result)
}

// BlankAssignmentInShortDecl tests blank assignment in short declaration
func BlankAssignmentInShortDecl() int {
	_, b, _ := 1, 2, 3
	return b
}

// BlankInTypeAssertion tests blank in type assertion
func BlankInTypeAssertion() bool {
	var i interface{} = 42
	_, ok := i.(int)
	return ok
}

// BlankInTypeSwitch tests blank in type switch
func BlankInTypeSwitch() string {
	var i interface{} = 42
	switch i.(type) {
	case int:
		return "int"
	case string:
		return "string"
	default:
		return "unknown"
	}
}

// NamedReturnWithDefer tests named return with defer
func NamedReturnWithDefer() (result int) {
	defer func() { result *= 2 }()
	return 5 // should become 10
}

// NamedReturnWithComplexDefer tests named return with complex defer
func NamedReturnWithComplexDefer() (result int) {
	result = 1
	defer func() {
		result++
	}()
	defer func() {
		result += 10
	}()
	return result // returns 1, then +10, then +1 = 12
}

// MultipleNamedReturns tests multiple named returns
func MultipleNamedReturns() (a, b, c int) {
	a = 1
	b = 2
	c = 3
	return // naked return
}

// NamedReturnShadowing tests named return shadowing
func NamedReturnShadowing() (result int) {
	result = 1
	if true {
		result := 2 // shadows named return
		_ = result
	}
	return // returns 1
}

// RecursivePointerType tests recursive pointer type
func RecursivePointerType() int {
	type Node struct {
		value int
		next  *Node
	}
	n1 := &Node{value: 1}
	n2 := &Node{value: 2, next: n1}
	return n2.value + n2.next.value
}

// MutualRecursiveTypes tests mutually recursive types
func MutualRecursiveTypes() int {
	type A struct {
		value int
		b     *struct {
			value int
			a     *A
		}
	}
	a := &A{value: 1}
	a.b = &struct {
		value int
		a     *A
	}{value: 2, a: a}
	return a.value + a.b.value
}

// DeeplyNestedPointer tests deeply nested pointer dereference
func DeeplyNestedPointer() int {
	x := 42
	p1 := &x
	p2 := &p1
	p3 := &p2
	p4 := &p3
	p5 := &p4
	return *****p5
}

// StructWithEmbeddedInterface tests struct with embedded interface
func StructWithEmbeddedInterface() int {
	type Stringer interface {
		String() string
	}
	type Container struct {
		Stringer
		value int
	}
	return 42
}

// InterfaceEmbedding tests interface embedding
func InterfaceEmbedding() int {
	type Reader interface {
		Read() int
	}
	type Writer interface {
		Write(int)
	}
	type ReadWriter interface {
		Reader
		Writer
	}
	return 42
}

// StructWithFuncFieldMethod tests struct with function field as method
func StructWithFuncFieldMethod() int {
	type Processor struct {
		process func(int) int
	}
	p := Processor{process: func(x int) int { return x * 2 }}
	return p.process(21)
}

// SliceWithNamedType tests slice with named element type
func SliceWithNamedType() int {
	type MyInt int
	s := []MyInt{1, 2, 3}
	return int(s[0] + s[1] + s[2])
}

// MapWithNamedType tests map with named key/value types
func MapWithNamedType() int {
	type Key string
	type Value int
	m := map[Key]Value{
		"a": 1,
		"b": 2,
	}
	return int(m["a"] + m["b"])
}

// ArrayWithNamedType tests array with named element type
func ArrayWithNamedType() int {
	type MyInt int
	var arr [3]MyInt
	arr[0] = 1
	arr[1] = 2
	arr[2] = 3
	return int(arr[0] + arr[1] + arr[2])
}

// PointerToNamedType tests pointer to named type
func PointerToNamedType() int {
	type MyInt int
	var x MyInt = 42
	p := &x
	return int(*p)
}

// NamedTypeSlice tests named type as slice
func NamedTypeSlice() int {
	type IntSlice []int
	var s IntSlice = []int{1, 2, 3}
	return s[0] + s[1] + s[2]
}

// NamedTypeMap tests named type as map
func NamedTypeMap() int {
	type StringMap map[string]int
	var m StringMap = map[string]int{"a": 1, "b": 2}
	return m["a"] + m["b"]
}

// NamedTypeFunc tests named type as func
func NamedTypeFunc() int {
	type IntFunc func(int) int
	var fn IntFunc = func(x int) int { return x * 2 }
	return fn(21)
}

// ============================================================================
// MORE EDGE CASES (Round 3) - Type aliases and more
// ============================================================================

// TypeAliasWithMethod tests type alias with method on underlying type
func TypeAliasWithMethod() int {
	type MyInt3 int
	type AliasInt = MyInt3
	var x AliasInt = 21
	return int(x) * 2
}

// TypeAliasSlice tests type alias for slice
func TypeAliasSlice() int {
	type IntSlice = []int
	var s IntSlice = []int{1, 2, 3}
	return len(s)
}

// TypeAliasMap tests type alias for map
func TypeAliasMap() int {
	type StringMap = map[string]int
	var m StringMap = map[string]int{"a": 1, "b": 2}
	return m["a"] + m["b"]
}

// TypeAliasFunc tests type alias for function
func TypeAliasFunc() int {
	type IntFunc = func(int) int
	var fn IntFunc = func(x int) int { return x * 2 }
	return fn(21)
}

// StructComparison tests struct equality comparison
func StructComparison() bool {
	type Point3 struct{ x, y int }
	p1 := Point3{x: 1, y: 2}
	p2 := Point3{x: 1, y: 2}
	p3 := Point3{x: 1, y: 3}
	return p1 == p2 && p1 != p3
}

// Counter3 is a counter for method value test
type Counter3 struct{ value int }

func (c *Counter3) Inc3() int {
	c.value++
	return c.value
}

// MethodValueTest tests method value binding
func MethodValueTest() int {
	c := &Counter3{value: 10}
	inc := c.Inc3
	return inc() + inc() + c.value
}

// MyInt4 is a named type for method expression test
type MyInt4 int

func (m MyInt4) Double4() int {
	return int(m) * 2
}

// MethodExpressionTest tests method expression
func MethodExpressionTest() int {
	var x MyInt4 = 21
	return MyInt4.Double4(x)
}

// EmbeddedFieldShadowing tests embedded field shadowing
func EmbeddedFieldShadowing() int {
	type Base3 struct{ value int }
	type Derived3 struct {
		Base3
		value int
	}
	d := Derived3{Base3: Base3{value: 10}, value: 20}
	return d.value + d.Base3.value
}

// MyInt5 is a named type for interface method set test
type MyInt5 int

func (m MyInt5) Add5(n int) int {
	return int(m) + n
}

// InterfaceMethodSet tests interface method set
func InterfaceMethodSet() int {
	type Adder5 interface{ Add5(int) int }
	var a Adder5 = MyInt5(10)
	return a.Add5(5)
}

// NestedClosureMutation tests nested closure with mutation
func NestedClosureMutation() int {
	x := 1
	outer := func() int {
		y := 10
		inner := func() int {
			x = x * 2
			return y + x
		}
		return inner() + inner()
	}
	return outer() + x
}

// DeferInClosureNamedReturn tests defer in closure with named return
func DeferInClosureNamedReturn() int {
	fn := func() (r int) {
		defer func() { r++ }()
		return 10
	}
	return fn()
}

// ReceiveFromClosedChannel tests receiving from closed channel
func ReceiveFromClosedChannel() (int, bool) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	return v1 + v2, ok1 && !ok2
}

// MapWithNilKey tests map with nil key
func MapWithNilKey() int {
	m := map[*int]int{nil: 100}
	var p *int
	return m[p]
}

// InterfaceEmbeddingTest tests interface embedding
func InterfaceEmbeddingTest() int {
	var rw ReadWriter = &File{data: 10}
	rw.Write(42)
	return rw.Read()
}

type Reader interface{ Read() int }
type Writer interface{ Write(int) }
type ReadWriter interface {
	Reader
	Writer
}

type File struct{ data int }

func (f *File) Read() int   { return f.data }
func (f *File) Write(v int) { f.data = v }

// ============================================================================
// MORE EDGE CASES (Round 4) - Trying to find more bugs
// ============================================================================

// ZeroValueStruct tests zero value struct initialization
func ZeroValueStruct() int {
	type Data struct {
		a int
		b string
		c []int
		d map[int]int
		e *int
	}
	var d Data
	return d.a + len(d.b) + len(d.c) + len(d.d)
}

// StructWithZeroSizeField tests struct with zero-size field
func StructWithZeroSizeField() int {
	type Empty struct{}
	type Data struct {
		value int
		empty Empty
	}
	d := Data{value: 42}
	return d.value
}

// SliceReslice tests reslicing behavior
func SliceReslice() int {
	s := make([]int, 5, 10)
	for i := range s {
		s[i] = i + 1
	}
	s = s[1:4]     // len=3, cap=9
	s = s[:cap(s)] // extend to cap
	return len(s) + cap(s)
}

// SliceResliceToCap tests reslicing to full capacity
func SliceResliceToCap() int {
	s := make([]int, 0, 10)
	s = append(s, 1, 2, 3)
	s = s[:cap(s)] // extend to cap
	return cap(s)
}

// NilSliceComparison tests nil slice comparison
func NilSliceComparison() bool {
	var s1 []int
	s2 := []int(nil)
	return s1 == nil && s2 == nil
}

// NilMapComparison tests nil map comparison
func NilMapComparison() bool {
	var m1 map[int]int
	m2 := map[int]int(nil)
	return m1 == nil && m2 == nil
}

// NilFuncComparison tests nil function comparison
func NilFuncComparison() bool {
	var f1 func()
	var f2 func() = nil
	return f1 == nil && f2 == nil
}

// NilChannelComparison tests nil channel comparison
func NilChannelComparison() bool {
	var ch1 chan int
	var ch2 chan int = nil
	return ch1 == nil && ch2 == nil
}

// EmptyStructComparison tests empty struct comparison
func EmptyStructComparison() bool {
	type Empty struct{}
	return Empty{} == Empty{}
}

// StructWithOnlyUnexported tests struct with only unexported fields
func StructWithOnlyUnexported() int {
	type private struct {
		value int
	}
	p := private{value: 42}
	return p.value
}

// StructWithOnlyExported tests struct with only exported fields
func StructWithOnlyExported() int {
	type Public struct {
		Value int
	}
	p := Public{Value: 42}
	return p.Value
}

// MapLookupReturnsZero tests map lookup returns zero value for missing key
func MapLookupReturnsZero() int {
	m := map[int]int{1: 10}
	return m[999] // should return 0
}

// MapLookupNilPointer tests map lookup with nil pointer value
func MapLookupNilPointer() bool {
	m := map[int]*int{1: nil}
	return m[1] == nil
}

// SliceCopyBehavior tests copy behavior
func SliceCopyBehavior() int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n + dst[0] + dst[1] + dst[2]
}

// SliceCopyOverlap tests copy with overlapping regions
func SliceCopyOverlap() int {
	s := []int{1, 2, 3, 4, 5}
	copy(s[1:], s[:3])
	return s[0] + s[1] + s[2] + s[3] + s[4]
}

// SliceCopyZero tests copy with zero elements
func SliceCopyZero() int {
	src := []int{1, 2, 3}
	dst := []int{}
	n := copy(dst, src)
	return n
}

// MapDeleteNonExistent tests deleting non-existent key
func MapDeleteNonExistent() int {
	m := map[int]int{1: 2}
	delete(m, 999) // no-op
	return len(m)
}

// MapLength tests map length
func MapLength() int {
	m := map[int]int{1: 2, 3: 4, 5: 6}
	return len(m)
}

// ChannelAfterClose tests channel behavior after close
func ChannelAfterClose() (int, bool) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	v3, ok3 := <-ch // zero value, ok=false
	return v1 + v2 + v3, ok1 && ok2 && !ok3
}

// ChannelCap tests channel capacity
func ChannelCap() int {
	ch := make(chan int, 5)
	return cap(ch)
}

// NonBufferedChannelCap tests non-buffered channel capacity
func NonBufferedChannelCap() int {
	ch := make(chan int)
	return cap(ch)
}

// PointerToZeroValue tests pointer to zero value
func PointerToZeroValue() int {
	var x int
	p := &x
	return *p // should return 0
}

// PointerToEmptyStruct tests pointer to empty struct
func PointerToEmptyStruct() int {
	type Empty struct{}
	var e Empty
	_ = &e
	return 1 // just verify no crash
}

// StructLiteralWithFieldNames tests struct literal with field names
func StructLiteralWithFieldNames() int {
	type Point struct {
		X, Y int
	}
	p := Point{X: 1, Y: 2}
	return p.X + p.Y
}

// StructLiteralWithoutFieldNames tests struct literal without field names
func StructLiteralWithoutFieldNames() int {
	type Point struct {
		X, Y int
	}
	p := Point{1, 2}
	return p.X + p.Y
}

// StructLiteralPartial tests partial struct literal
func StructLiteralPartial() int {
	type Data struct {
		A, B, C int
	}
	d := Data{A: 1}
	return d.A + d.B + d.C // should be 1 + 0 + 0
}

// StructLiteralWithPointers tests struct literal with pointer fields
func StructLiteralWithPointers() int {
	type Data struct {
		A *int
		B *string
	}
	x := 42
	s := "hello"
	d := Data{A: &x, B: &s}
	return *d.A + len(*d.B)
}

// NestedStructLiteral tests nested struct literal
func NestedStructLiteral() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 42}}
	return o.I.X
}

// ArrayLiteralWithIndex tests array literal with index
func ArrayLiteralWithIndex() int {
	arr := [5]int{1: 10, 3: 30}
	return arr[1] + arr[3]
}

// ArrayLiteralWithExpression tests array literal with expression index
func ArrayLiteralWithExpression() int {
	const N = 3
	arr := [5]int{N: 42}
	return arr[N]
}

// SliceLiteralWithIndex tests slice literal with index
func SliceLiteralWithIndex() int {
	s := []int{0: 1, 2: 3, 4: 5}
	return s[0] + s[2] + s[4]
}

// SliceLiteralWithExpression tests slice literal with expression index
func SliceLiteralWithExpression() int {
	const idx = 2
	s := []int{idx: 42}
	return s[idx]
}

// MapLiteralWithComplexKey tests map literal with complex key
func MapLiteralWithComplexKey() int {
	type Key struct{ A, B int }
	m := map[Key]int{
		{1, 2}: 3,
		{4, 5}: 6,
	}
	return m[Key{1, 2}]
}

// InterfaceNilTypeAssertion tests type assertion on nil interface
func InterfaceNilTypeAssertion() bool {
	var i interface{}
	_, ok := i.(int)
	return !ok // should be false
}

// InterfaceNilTypeSwitch tests type switch on nil interface
func InterfaceNilTypeSwitch() string {
	var i interface{}
	switch i.(type) {
	case int:
		return "int"
	case nil:
		return "nil"
	default:
		return "other"
	}
}

// InterfaceConcreteToInterface tests concrete to interface conversion
func InterfaceConcreteToInterface() int {
	var i interface{} = 42
	return i.(int)
}

// InterfaceToEmptyInterface tests interface to empty interface conversion
func InterfaceToEmptyInterface() int {
	type Stringer interface{ String() string }
	var s Stringer
	_ = interface{}(s)
	return 1 // just verify no crash
}

// PointerInterface tests pointer as interface
func PointerInterface() int {
	var i interface{} = new(int)
	p := i.(*int)
	return *p
}

// SliceInterface tests slice as interface
func SliceInterface() int {
	var i interface{} = []int{1, 2, 3}
	s := i.([]int)
	return len(s)
}

// MapInterface tests map as interface
func MapInterface() int {
	var i interface{} = map[int]int{1: 2}
	m := i.(map[int]int)
	return m[1]
}

// FuncInterface tests func as interface
func FuncInterface() int {
	var i interface{} = func() int { return 42 }
	fn := i.(func() int)
	return fn()
}

// ChanInterface tests channel as interface
func ChanInterface() int {
	var i interface{} = make(chan int, 1)
	ch := i.(chan int)
	return cap(ch)
}

// StructZeroValueComparison tests zero value struct comparison
func StructZeroValueComparison() bool {
	type Point struct{ X, Y int }
	var p1, p2 Point
	return p1 == p2
}

// StructFieldZeroValue tests struct field zero value
func StructFieldZeroValue() int {
	type Data struct {
		Int     int
		String  string
		Slice   []int
		Map     map[int]int
		Channel chan int
		Func    func()
		Ptr     *int
	}
	var d Data
	return d.Int + len(d.String) + len(d.Slice) + len(d.Map)
}

// ClosureReadsOuter tests closure reading outer variable
func ClosureReadsOuter() int {
	x := 42
	fn := func() int {
		return x
	}
	return fn()
}

// ClosureWritesOuter tests closure writing outer variable
func ClosureWritesOuter() int {
	x := 0
	fn := func() {
		x = 42
	}
	fn()
	return x
}

// ClosureReturnsOuter tests closure returning outer variable address
func ClosureReturnsOuter() int {
	x := 42
	fn := func() *int {
		return &x
	}
	return *fn()
}

// ClosureMultipleReturn tests closure with multiple returns
func ClosureMultipleReturn() (int, int) {
	fn := func() (int, int) {
		return 10, 20
	}
	return fn()
}

// ClosureVariadic tests variadic closure
func ClosureVariadic() int {
	fn := func(args ...int) int {
		sum := 0
		for _, a := range args {
			sum += a
		}
		return sum
	}
	return fn(1, 2, 3, 4, 5)
}

// DeferNamedReturnMultiple tests multiple named returns with defer
func DeferNamedReturnMultiple() (a, b int) {
	defer func() {
		a++
		b++
	}()
	a = 10
	b = 20
	return
}

// DeferModifiesMultipleNamed tests defer modifying multiple named returns
func DeferModifiesMultipleNamed() (a, b int) {
	defer func() {
		a, b = b, a
	}()
	a = 10
	b = 20
	return
}

// ForBreakContinue tests for loop with break and continue
func ForBreakContinue() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i == 3 {
			continue
		}
		if i == 7 {
			break
		}
		sum += i
	}
	return sum
}

// RangeBreakContinue tests range with break and continue
func RangeBreakContinue() int {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for i, v := range s {
		if i == 1 {
			continue
		}
		if i == 3 {
			break
		}
		sum += v
	}
	return sum
}

// SwitchWithFallthrough tests switch with fallthrough
func SwitchWithFallthrough() int {
	result := 0
	switch 1 {
	case 1:
		result += 1
		fallthrough
	case 2:
		result += 2
		fallthrough
	case 3:
		result += 4
	}
	return result
}

// SwitchWithoutCondition tests switch without condition
func SwitchWithoutCondition() int {
	x := 5
	switch {
	case x < 0:
		return -1
	case x == 0:
		return 0
	default:
		return 1
	}
}

// SelectWithTimeout tests select with timeout simulation
func SelectWithTimeout() int {
	ch := make(chan int, 1)
	select {
	case v := <-ch:
		return v
	default:
		return -1 // timeout
	}
}

// GotoWithLabel tests goto with label
func GotoWithLabel() int {
	i := 0
start:
	i++
	if i < 5 {
		goto start
	}
	return i
}

// TypeConversionBasic tests basic type conversion
func TypeConversionBasic() int {
	var x int32 = 42
	return int(x)
}

// TypeConversionFloat tests float type conversion
func TypeConversionFloat() int {
	var x float64 = 42.9
	return int(x) // truncates to 42
}

// TypeConversionComplex tests complex type conversion chain
func TypeConversionComplex() int {
	x := 42.5
	return int(int64(x))
}

// SliceOfStringToInterface tests slice of string to interface
func SliceOfStringToInterface() int {
	s := []string{"a", "b", "c"}
	var i interface{} = s
	return len(i.([]string))
}

// MapOfStringToInterface tests map of string to interface
func MapOfStringToInterface() int {
	m := map[string]int{"a": 1}
	var i interface{} = m
	return i.(map[string]int)["a"]
}

// EmptySliceCopy tests copy to empty slice
func EmptySliceCopy() int {
	src := []int{1, 2, 3}
	var dst []int
	n := copy(dst, src)
	return n // should be 0
}

// NilSliceCopyTo tests copy from nil slice
func NilSliceCopyTo() int {
	var src []int
	dst := make([]int, 5)
	n := copy(dst, src)
	return n // should be 0
}

// ============================================================================
// MORE EDGE CASES (Round 5) - More corner cases
// ============================================================================

// AppendToNilSlice tests append to nil slice
func AppendToNilSlice() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// AppendExpand tests append expanding capacity
func AppendExpand() int {
	s := make([]int, 0, 1)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return len(s)
}

// AppendSliceToSlice tests appending slice to slice
func AppendSliceToSlice() int {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	s1 = append(s1, s2...)
	return s1[5]
}

// SliceMakeLenCap tests make with len and cap
func SliceMakeLenCap() int {
	s := make([]int, 3, 10)
	return len(s)*100 + cap(s)
}

// SliceMakeLenOnly tests make with len only
func SliceMakeLenOnly() int {
	s := make([]int, 5)
	return len(s)*100 + cap(s)
}

// MapMakeWithSize tests make map with size hint
func MapMakeWithSize() int {
	m := make(map[int]int, 100)
	m[1] = 2
	return len(m)
}

// ChannelMakeBuffered tests make buffered channel
func ChannelMakeBuffered() int {
	ch := make(chan int, 5)
	return cap(ch)
}

// ChannelMakeUnbuffered tests make unbuffered channel
func ChannelMakeUnbuffered() int {
	ch := make(chan int)
	return cap(ch)
}

// NilSliceAppendNil tests appending nil to nil slice
func NilSliceAppendNil() int {
	var s1 []int
	var s2 []int
	s1 = append(s1, s2...)
	return len(s1)
}

// SliceThreeIndexReslice tests three-index reslice
func SliceThreeIndexReslice() int {
	s := make([]int, 5, 10)
	s2 := s[1:3:5]
	return len(s2)*100 + cap(s2)
}

// SliceZeroLength tests zero-length slice operations
func SliceZeroLength() int {
	s := []int{}
	s = append(s, 1)
	return len(s)
}

// MapIterateAndModify tests iterating and modifying map
func MapIterateAndModify() int {
	m := map[int]int{1: 1, 2: 2, 3: 3}
	for k := range m {
		m[k] *= 2
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// MapNestedDelete tests nested map delete
func MapNestedDelete() int {
	m := map[int]map[int]int{
		1: {1: 10, 2: 20},
		2: {3: 30, 4: 40},
	}
	delete(m[1], 2)
	return len(m[1])
}

// StructFieldPointer tests struct field pointer
func StructFieldPointer() int {
	type Data struct{ Value int }
	d := Data{Value: 42}
	p := &d.Value
	*p = 100
	return d.Value
}

// StructFieldPointerModify tests modifying struct via field pointer
func StructFieldPointerModify() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 10}}
	p := &o.I.X
	*p = 42
	return o.I.X
}

// PointerToArray tests pointer to array
func PointerToArray() int {
	arr := [3]int{1, 2, 3}
	p := &arr
	(*p)[1] = 20
	return arr[1]
}

// PointerToArrayFullSlice tests pointer to array with full slice
func PointerToArrayFullSlice() int {
	arr := [3]int{1, 2, 3}
	p := &arr
	s := (*p)[:]
	return s[0] + s[1] + s[2]
}

// ArrayPointerModification tests modifying array via pointer
func ArrayPointerModification() int {
	var arr [3]int
	p := &arr
	p[0] = 1
	p[1] = 2
	p[2] = 3
	return arr[0] + arr[1] + arr[2]
}

// SlicePointerModification tests modifying slice via pointer
func SlicePointerModification() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 10
	return s[0]
}

// MultipleAssignDifferentTypes tests multiple assignment with different types
func MultipleAssignDifferentTypes() (int, string, bool) {
	a, b, c := 42, "hello", true
	return a, b, c
}

// MultipleAssignSameExpression tests multiple assignment from same expression
func MultipleAssignSameExpression() (int, int) {
	fn := func() (int, int) { return 10, 20 }
	a, b := fn()
	return a, b
}

// TypeAssertionOnConcrete tests type assertion on concrete type
func TypeAssertionOnConcrete() int {
	var i interface{} = 42
	switch i.(type) {
	case int:
		return i.(int)
	default:
		return -1
	}
}

// TypeSwitchMultipleCases tests type switch with multiple cases
func TypeSwitchMultipleCases() string {
	var i interface{} = 42.0
	switch i.(type) {
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case string:
		return "string"
	default:
		return "other"
	}
}

// InterfaceConversion tests interface conversion
func InterfaceConversion() int {
	type Stringer interface{ String() string }
	type Stringer2 interface{ String() string }
	var s1 Stringer
	var s2 Stringer2 = s1
	_ = s2
	return 1
}

// InterfaceNilAssignment tests nil assignment to interface
func InterfaceNilAssignment() bool {
	var i interface{}
	i = nil
	return i == nil
}

// InterfaceTypedNilAssignment tests typed nil assignment to interface
func InterfaceTypedNilAssignment() bool {
	var s []int = nil
	var i interface{} = s
	return i != nil // typed nil inside interface is not nil
}

// StructMethodOnPointer tests method on pointer receiver
func StructMethodOnPointer() int {
	type Counter struct{ value int }
	c := &Counter{value: 10}
	// Simulate method behavior
	c.value++
	return c.value
}

// StructMethodOnValue tests method on value receiver
func StructMethodOnValue() int {
	type Point struct{ x, y int }
	p := Point{x: 1, y: 2}
	// Simulate method behavior
	return p.x + p.y
}

// EmbeddingMethodPromotion tests method promotion from embedded field
func EmbeddingMethodPromotion() int {
	type Inner struct{ value int }
	type Outer struct {
		Inner
	}
	o := Outer{Inner: Inner{value: 42}}
	return o.value
}

// EmbeddingFieldPromotion tests field promotion from embedded field
func EmbeddingFieldPromotion() int {
	type Base struct {
		X, Y int
	}
	type Derived struct {
		Base
		Z int
	}
	d := Derived{Base: Base{X: 1, Y: 2}, Z: 3}
	return d.X + d.Y + d.Z
}

// EmbeddingPointerMethod tests method on embedded pointer
func EmbeddingPointerMethod() int {
	type Inner struct{ value int }
	type Outer struct {
		*Inner
	}
	o := Outer{Inner: &Inner{value: 42}}
	return o.value
}

// MultipleEmbeddingConflictResolution tests multiple embedding conflict resolution
func MultipleEmbeddingConflictResolution() int {
	type A struct{ a int }
	type B struct{ a int }
	type C struct {
		A
		B
	}
	c := C{A: A{a: 1}, B: B{a: 2}}
	return c.A.a + c.B.a
}

// StructComparisonAllTypes tests struct comparison with various types
func StructComparisonAllTypes() bool {
	type Data struct {
		Int   int
		Float float64
		Str   string
		Bool  bool
	}
	d1 := Data{Int: 1, Float: 1.5, Str: "a", Bool: true}
	d2 := Data{Int: 1, Float: 1.5, Str: "a", Bool: true}
	return d1 == d2
}

// StructWithNestedSlice tests struct with nested slice
func StructWithNestedSlice() int {
	type Matrix struct {
		Rows [][]int
	}
	m := Matrix{Rows: [][]int{{1, 2}, {3, 4}}}
	return m.Rows[0][0] + m.Rows[1][1]
}

// StructWithNestedMap tests struct with nested map
func StructWithNestedMap() int {
	type Dict struct {
		Data map[string]map[int]int
	}
	d := Dict{Data: map[string]map[int]int{
		"a": {1: 10},
		"b": {2: 20},
	}}
	return d.Data["a"][1] + d.Data["b"][2]
}

// ClosureCaptureSliceElement tests closure capturing slice element
func ClosureCaptureSliceElement() int {
	s := []int{1, 2, 3}
	fn := func() int {
		return s[1]
	}
	return fn()
}

// ClosureCaptureMapValue tests closure capturing map value
func ClosureCaptureMapValue() int {
	m := map[int]int{1: 10, 2: 20}
	fn := func() int {
		return m[1]
	}
	return fn()
}

// ClosureCaptureStructField tests closure capturing struct field
func ClosureCaptureStructField() int {
	type Data struct{ Value int }
	d := Data{Value: 42}
	fn := func() int {
		return d.Value
	}
	return fn()
}

// DeferClosureArgCapture tests defer closure argument capture timing
func DeferClosureArgCapture() int {
	x := 1
	defer func(v int) {
		_ = v
	}(x)
	x = 2
	return x
}

// DeferClosureNoArg tests defer closure without arguments
func DeferClosureNoArg() int {
	x := 1
	defer func() {
		x = 2
	}()
	return x
}

// ForRangeModifyValue tests for range modifying value
func ForRangeModifyValue() int {
	s := []int{1, 2, 3}
	sum := 0
	for i, v := range s {
		sum += v
		s[i] = v * 2
	}
	return sum
}

// ForRangeMapModify tests for range map modification
func ForRangeMapModify() int {
	m := map[int]int{1: 1, 2: 2}
	for k := range m {
		m[k] = m[k] * 2
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// SelectNonBlocking tests non-blocking select
func SelectNonBlocking() int {
	ch := make(chan int, 1)
	select {
	case ch <- 1:
		return 1
	default:
		return 0
	}
}

// SwitchEmptyCases tests switch with empty case bodies
func SwitchEmptyCases() int {
	x := 1
	switch x {
	case 1:
		// empty
	case 2:
		return 2
	}
	return 0
}

// SwitchDefaultFirst tests switch with default first
func SwitchDefaultFirst() int {
	x := 1
	switch x {
	default:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	}
}

// GotoSkipDeclaration tests goto skipping variable declaration
func GotoSkipDeclaration() int {
	goto skip
	// can't skip variable declaration in Go
skip:
	return 0
}

// LabelInNestedLoop tests label in nested loop
func LabelInNestedLoop() int {
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i+j > 3 {
				break outer
			}
		}
	}
	return 0
}

// ContinueInNestedLoop tests continue in nested loop
func ContinueInNestedLoop() int {
	sum := 0
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue outer
			}
			sum++
		}
	}
	return sum
}

// BreakInSelect tests break in select
func BreakInSelect() int {
	ch := make(chan int, 1)
	for {
		select {
		case <-ch:
			break
		default:
			return 0
		}
	}
}

// ============================================================================
// MORE EDGE CASES (Round 6) - Even more corner cases
// ============================================================================

// SliceAppendOverflow tests append beyond capacity
func SliceAppendOverflow() int {
	s := make([]int, 0, 2)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return len(s)
}

// MapPreallocate tests map preallocation
func MapPreallocate() int {
	m := make(map[int]int, 1000)
	m[1] = 1
	return len(m)
}

// ChannelSendRecv tests basic channel send/recv
func ChannelSendRecv() int {
	ch := make(chan int, 1)
	ch <- 42
	return <-ch
}

// ChannelBufferedMultiple tests multiple buffered sends
func ChannelBufferedMultiple() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	return <-ch + <-ch + <-ch
}

// StructWithAllBasicTypes tests struct with all basic types
func StructWithAllBasicTypes() int {
	type AllTypes struct {
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		Bool    bool
		String  string
	}
	a := AllTypes{
		Int: 1, Int8: 2, Int16: 3, Int32: 4, Int64: 5,
		Uint: 6, Uint8: 7, Uint16: 8, Uint32: 9, Uint64: 10,
		Float32: 1.5, Float64: 2.5, Bool: true, String: "hello",
	}
	return int(a.Int) + int(a.Int8) + int(a.Int16) + len(a.String)
}

// PointerToAllBasicTypes tests pointer to all basic types
func PointerToAllBasicTypes() int {
	var i int = 42
	var f float64 = 3.14
	var s string = "hello"
	pi := &i
	pf := &f
	ps := &s
	return *pi + int(*pf) + len(*ps)
}

// SliceOfAllBasicTypes tests slice of all basic types
func SliceOfAllBasicTypes() int {
	si := []int{1, 2, 3}
	sf := []float64{1.0, 2.0, 3.0}
	ss := []string{"a", "b", "c"}
	sb := []bool{true, false, true}
	return len(si) + len(sf) + len(ss) + len(sb)
}

// MapOfAllBasicTypes tests map with all basic key types
func MapOfAllBasicTypes() int {
	mi := map[int]int{1: 10}
	mf := map[float64]int{1.0: 10}
	ms := map[string]int{"a": 10}
	return len(mi) + len(mf) + len(ms)
}

// ArrayFixedSize tests fixed size array
func ArrayFixedSize() int {
	arr := [100]int{}
	arr[0] = 1
	arr[99] = 2
	return arr[0] + arr[99]
}

// ArrayZeroSized tests zero-sized array
func ArrayZeroSized() int {
	var arr [0]int
	return len(arr)
}

// SliceOfZeroSizedArray tests slice of zero-sized arrays
func SliceOfZeroSizedArray() int {
	s := [][0]int{{}, {}, {}}
	return len(s)
}

// StructWithZeroSizedArray tests struct with zero-sized array
func StructWithZeroSizedArray() int {
	type Data struct {
		Arr [0]int
		Val int
	}
	d := Data{Val: 42}
	return d.Val
}

// NilPointerToStruct tests nil pointer to struct
func NilPointerToStruct() bool {
	type Data struct{ Value int }
	var p *Data
	return p == nil
}

// NilPointerToSlice tests nil pointer to slice
func NilPointerToSlice() bool {
	var p *[]int
	return p == nil
}

// NilPointerToMap tests nil pointer to map
func NilPointerToMap() bool {
	var p *map[int]int
	return p == nil
}

// EmptyStructLiteral tests empty struct literal
func EmptyStructLiteral() int {
	type Empty struct{}
	e := Empty{}
	_ = e
	return 1
}

// EmptyInterfaceLiteral tests empty interface literal
func EmptyInterfaceLiteral() interface{} {
	return interface{}(nil)
}

// InterfaceSliceOfInterfaces tests interface slice of interfaces
func InterfaceSliceOfInterfaces() int {
	var s []interface{}
	s = append(s, 1, "hello", 3.14, true, nil)
	return len(s)
}

// MapOfInterfaces tests map with interface values
func MapOfInterfaces() int {
	m := map[string]interface{}{
		"int":    42,
		"string": "hello",
		"nil":    nil,
	}
	return len(m)
}

// NestedInterfaceSlice tests nested interface slice
func NestedInterfaceSlice() int {
	var outer []interface{}
	inner := []interface{}{1, 2, 3}
	outer = append(outer, inner)
	return len(outer)
}

// NestedInterfaceMap tests nested interface map
func NestedInterfaceMap() int {
	m := map[string]interface{}{
		"nested": map[string]int{"a": 1, "b": 2},
	}
	return len(m)
}

// TypeAssertionChained tests chained type assertions
func TypeAssertionChained() int {
	var i interface{} = 42
	if v, ok := i.(int); ok {
		return v
	}
	return 0
}

// TypeAssertionOnConcreteType tests type assertion on concrete type
func TypeAssertionOnConcreteType() int {
	var i interface{} = "hello"
	s, ok := i.(string)
	if ok {
		return len(s)
	}
	return 0
}

// MultipleTypeAssertions tests multiple type assertions in sequence
func MultipleTypeAssertions() string {
	var i interface{} = 42.5
	if _, ok := i.(int); ok {
		return "int"
	}
	if _, ok := i.(float64); ok {
		return "float"
	}
	return "unknown"
}

// SwitchTypeAssertion tests switch with type assertion
func SwitchTypeAssertion() string {
	var i interface{} = []int{1, 2, 3}
	switch v := i.(type) {
	case int:
		return "int"
	case []int:
		return "slice"
	default:
		_ = v
		return "other"
	}
}

// ClosureWithDeferAndReturn tests closure with defer and return
func ClosureWithDeferAndReturn() int {
	fn := func() int {
		defer func() {}()
		return 42
	}
	return fn()
}

// MultipleClosures tests multiple closures
func MultipleClosures() int {
	add := func(a, b int) int { return a + b }
	mul := func(a, b int) int { return a * b }
	return add(1, 2) + mul(3, 4)
}

// ClosureAsParameter tests closure as parameter
func ClosureAsParameter() int {
	apply := func(fn func(int) int, x int) int {
		return fn(x)
	}
	double := func(x int) int { return x * 2 }
	return apply(double, 21)
}

// ClosureAsReturn tests closure as return value
func ClosureAsReturn() int {
	makeAdder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := makeAdder(5)
	return add5(10)
}

// ClosureCapturingPointer tests closure capturing pointer
func ClosureCapturingPointer() int {
	x := 10
	p := &x
	fn := func() {
		*p = 42
	}
	fn()
	return *p
}

// ClosureCapturingSlice tests closure capturing slice
func ClosureCapturingSlice() int {
	s := []int{1, 2, 3}
	fn := func() {
		s[0] = 10
	}
	fn()
	return s[0]
}

// ClosureCapturingMap tests closure capturing map
func ClosureCapturingMap() int {
	m := map[int]int{1: 10}
	fn := func() {
		m[1] = 20
	}
	fn()
	return m[1]
}

// DeferWithMethodCall tests defer with method call
func DeferWithMethodCall() int {
	type Counter struct{ value int }
	c := &Counter{}
	defer func() { c.value++ }()
	c.value = 10
	return c.value + 1
}

// DeferWithMultipleReturns tests defer with multiple returns
func DeferWithMultipleReturns() (int, int) {
	a, b := 0, 0
	defer func() {
		a++
		b++
	}()
	return a, b
}

// DeferInClosure tests defer inside closure
func DeferInClosureNormal() int {
	fn := func() int {
		defer func() {}()
		return 42
	}
	return fn()
}

// ForWithDefer tests for with defer
func ForWithDefer() int {
	sum := 0
	for i := 0; i < 3; i++ {
		func() {
			defer func() { sum++ }()
		}()
	}
	return sum
}

// RangeWithDefer tests range with defer
func RangeWithDefer() int {
	s := []int{1, 2, 3}
	sum := 0
	for _, v := range s {
		func(x int) {
			defer func() { sum += x }()
		}(v)
	}
	return sum
}

// MapRangeOrderIndependent tests map range order independence
func MapRangeOrderIndependent() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum // should be 60 regardless of order
}

// ChannelCloseMultipleReceive tests multiple receives from closed channel
func ChannelCloseMultipleReceive() int {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	v3, ok3 := <-ch
	return v1 + v2 + v3 + boolToInt(ok1) + boolToInt(ok2) + boolToInt(ok3)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// SelectWithMultipleReady tests select with multiple ready cases
func SelectWithMultipleReady() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	ch2 <- 2
	select {
	case v := <-ch1:
		return v
	case v := <-ch2:
		return v
	}
}

// SwitchWithExpression tests switch with expression
func SwitchWithExpression() int {
	x := 5
	switch x * 2 {
	case 10:
		return 1
	case 20:
		return 2
	default:
		return 0
	}
}

// SwitchWithFunctionCall tests switch with function call
func SwitchWithFunctionCall() int {
	fn := func() int { return 42 }
	switch fn() {
	case 42:
		return 1
	default:
		return 0
	}
}

// GotoWithCondition tests goto with condition
func GotoWithCondition() int {
	i := 0
start:
	if i >= 5 {
		return i
	}
	i++
	goto start
}

// LabelBeforeStatement tests label before statement (used for goto)
func LabelBeforeStatement() int {
	goto label
label:
	return 42
}

// TypeConversionToInt tests type conversion to int
func TypeConversionToInt() int {
	var f float64 = 42.9
	return int(f)
}

// TypeConversionToFloat tests type conversion to float
func TypeConversionToFloat() float64 {
	var i int = 42
	return float64(i)
}

// TypeConversionToString tests type conversion to string
func TypeConversionToString() string {
	b := []byte{'h', 'i'}
	return string(b)
}

// TypeConversionToSlice tests type conversion to slice
func TypeConversionToSlice() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

// StructLiteralPartialFields tests struct literal with partial fields
func StructLiteralPartialFields() int {
	type Data struct {
		A, B, C, D int
	}
	d := Data{A: 1, C: 3}
	return d.A + d.B + d.C + d.D
}

// StructLiteralAllFields tests struct literal with all fields
func StructLiteralAllFields() int {
	type Data struct{ A, B int }
	d := Data{A: 1, B: 2}
	return d.A + d.B
}

// StructLiteralPositional tests struct literal with positional args
func StructLiteralPositional() int {
	type Data struct{ A, B int }
	d := Data{1, 2}
	return d.A + d.B
}

// SliceLiteralWithIndices tests slice literal with indices
func SliceLiteralWithIndices() int {
	s := []int{0: 1, 2: 3, 4: 5}
	return s[0] + s[2] + s[4]
}

// ArrayLiteralWithIndices tests array literal with indices
func ArrayLiteralWithIndices() int {
	arr := [5]int{0: 1, 2: 3, 4: 5}
	return arr[0] + arr[2] + arr[4]
}

// MapLiteralEmpty tests empty map literal
func MapLiteralEmpty() int {
	m := map[int]int{}
	return len(m)
}

// SliceLiteralEmpty tests empty slice literal
func SliceLiteralEmpty() int {
	s := []int{}
	return len(s)
}

// NilComparisonAllTypes tests nil comparison for all nil-able types
func NilComparisonAllTypes() bool {
	var s []int
	var m map[int]int
	var ch chan int
	var fn func()
	var p *int
	return s == nil && m == nil && ch == nil && fn == nil && p == nil
}

// LenCapOnAllTypes tests len and cap on all applicable types
func LenCapOnAllTypes() int {
	s := make([]int, 5, 10)
	ch := make(chan int, 3)
	m := map[int]int{1: 2, 3: 4}
	str := "hello"
	return len(s) + cap(s) + len(ch) + cap(ch) + len(m) + len(str)
}

// ============================================================================
// ROUND 7: MORE CORNER CASES - Method sets, type assertions, edge cases
// ============================================================================

// MethodSetOnNamedType tests method set on named type
func MethodSetOnNamedType() int {
	type MySlice []int
	var s MySlice = []int{1, 2, 3}
	return len(s)
}

// MethodSetOnNamedMap tests method set on named map
func MethodSetOnNamedMap() int {
	type MyMap map[int]int
	var m MyMap = map[int]int{1: 2}
	return m[1]
}

// MethodSetOnNamedFunc tests named function type
func MethodSetOnNamedFunc() int {
	type IntFunc func(int) int
	var f IntFunc = func(x int) int { return x * 2 }
	return f(21)
}

// EmptyInterfaceTypeAssertion tests empty interface type assertion
func EmptyInterfaceTypeAssertion() string {
	var i interface{} = "hello"
	if s, ok := i.(string); ok {
		return s
	}
	return "not a string"
}

// InterfaceTypeAssertionWithNil tests nil interface type assertion
func InterfaceTypeAssertionWithNil() bool {
	var i interface{}
	_, ok := i.(int)
	return !ok
}

// InterfaceTypeAssertionWithConcrete tests concrete type assertion
func InterfaceTypeAssertionWithConcrete() int {
	var i interface{} = 42
	return i.(int)
}

// InterfaceTypeSwitchWithMultipleTypes tests type switch with multiple types
func InterfaceTypeSwitchWithMultipleTypes() string {
	test := func(i interface{}) string {
		switch i.(type) {
		case int, int8, int16:
			return "int"
		case uint, uint8:
			return "uint"
		case string:
			return "string"
		case float32, float64:
			return "float"
		default:
			return "other"
		}
	}
	return test(42.0)
}

// PointerToMapValueNotSupported tests that pointer to map value is not supported
func PointerToMapValueNotSupported() int {
	m := map[int]int{1: 10}
	// Can't take address of map element in Go
	return m[1]
}

// PointerToStructField tests pointer to struct field
func PointerToStructField() int {
	type Data struct{ Value int }
	d := Data{Value: 42}
	p := &d.Value
	*p = 100
	return d.Value
}

// PointerToNestedStructField tests pointer to nested struct field
func PointerToNestedStructField() int {
	type Inner struct{ Value int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{Value: 42}}
	p := &o.I.Value
	*p = 100
	return o.I.Value
}

// NilPointerDereference tests nil pointer dereference handling
func NilPointerDereference() int {
	var p *int
	if p == nil {
		return 0
	}
	return *p
}

// PointerComparison tests pointer comparison
func PointerComparison() bool {
	x := 42
	p1 := &x
	p2 := &x
	return p1 == p2
}

// DifferentPointerComparison tests different pointer comparison
func DifferentPointerComparison() bool {
	x, y := 42, 42
	p1 := &x
	p2 := &y
	return p1 != p2
}

// SliceOfPointersToStruct tests slice of pointers to struct
func SliceOfPointersToStruct() int {
	type Item struct{ Value int }
	items := []*Item{{1}, {2}, {3}}
	sum := 0
	for _, item := range items {
		if item != nil {
			sum += item.Value
		}
	}
	return sum
}

// MapOfPointersToStruct tests map of pointers to struct
func MapOfPointersToStruct() int {
	type Item struct{ Value int }
	m := map[string]*Item{
		"a": {1},
		"b": {2},
	}
	return m["a"].Value + m["b"].Value
}

// StructWithPointerTypeField tests struct with pointer type field
func StructWithPointerTypeField() int {
	type Data struct {
		Value *int
	}
	v := 42
	d := Data{Value: &v}
	return *d.Value
}

// StructWithSliceTypeField tests struct with slice type field
func StructWithSliceTypeField() int {
	type Data struct {
		Items []int
	}
	d := Data{Items: []int{1, 2, 3}}
	return len(d.Items)
}

// StructWithMapTypeField tests struct with map type field
func StructWithMapTypeField() int {
	type Data struct {
		Items map[int]int
	}
	d := Data{Items: map[int]int{1: 2}}
	return d.Items[1]
}

// StructWithChannelTypeField tests struct with channel type field
func StructWithChannelTypeField() int {
	type Data struct {
		Ch chan int
	}
	d := Data{Ch: make(chan int, 1)}
	d.Ch <- 42
	return <-d.Ch
}

// StructWithFuncTypeField tests struct with function type field
func StructWithFuncTypeField() int {
	type Data struct {
		Func func(int) int
	}
	d := Data{Func: func(x int) int { return x * 2 }}
	return d.Func(21)
}

// NestedStructWithMethods tests nested struct with methods
func NestedStructWithMethods() int {
	type Inner struct{ Value int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{Value: 42}}
	return o.I.Value
}

// EmbeddedStructWithMethods tests embedded struct with methods
func EmbeddedStructWithMethods() int {
	type Base struct{ Value int }
	type Derived struct {
		Base
		Extra int
	}
	d := Derived{Base: Base{Value: 10}, Extra: 20}
	return d.Value + d.Extra
}

// MultipleEmbeddedStructs tests multiple embedded structs
func MultipleEmbeddedStructs() int {
	type A struct{ AVal int }
	type B struct{ BVal int }
	type C struct {
		A
		B
		CVal int
	}
	c := C{A: A{AVal: 1}, B: B{BVal: 2}, CVal: 3}
	return c.AVal + c.BVal + c.CVal
}

// StructWithPrivateField tests struct with private field
func StructWithPrivateField() int {
	type Data struct {
		value int
	}
	d := Data{value: 42}
	return d.value
}

// StructWithMixedFields tests struct with mixed public/private fields
func StructWithMixedFields() int {
	type Data struct {
		Public  int
		private int
	}
	d := Data{Public: 1, private: 2}
	return d.Public + d.private
}

// EmptyStruct tests empty struct
func EmptyStruct() int {
	type Empty struct{}
	var e Empty
	_ = e
	return 1
}

// EmptyStructPointer tests empty struct pointer
func EmptyStructPointer() int {
	type Empty struct{}
	e := &Empty{}
	_ = e
	return 1
}

// StructAlignment tests struct alignment
func StructAlignment() int {
	type Data struct {
		A byte
		B int
		C byte
	}
	d := Data{A: 1, B: 2, C: 3}
	return int(d.A) + d.B + int(d.C)
}

// StructWithPadding tests struct with padding
func StructWithPadding() int {
	type Data struct {
		A byte
		// 7 bytes padding
		B int64
	}
	d := Data{A: 1, B: 2}
	return int(d.A) + int(d.B)
}

// ArrayOfEmptyStruct tests array of empty struct
func ArrayOfEmptyStruct() int {
	type Empty struct{}
	arr := [3]Empty{{}, {}, {}}
	return len(arr)
}

// SliceOfEmptyStruct tests slice of empty struct
func SliceOfEmptyStruct() int {
	type Empty struct{}
	s := []Empty{{}, {}, {}}
	return len(s)
}

// MapWithEmptyStructValue tests map with empty struct value
func MapWithEmptyStructValue() int {
	type Empty struct{}
	m := map[int]Empty{1: {}, 2: {}}
	return len(m)
}

// ChannelOfEmptyStruct tests channel of empty struct
func ChannelOfEmptyStruct() int {
	type Empty struct{}
	ch := make(chan Empty, 2)
	ch <- Empty{}
	ch <- Empty{}
	return len(ch)
}

// FuncReturningEmptyStruct tests function returning empty struct
func FuncReturningEmptyStruct() int {
	type Empty struct{}
	f := func() Empty { return Empty{} }
	_ = f()
	return 1
}

// ClosureCapturingEmptyStruct tests closure capturing empty struct
func ClosureCapturingEmptyStruct() int {
	type Empty struct{}
	e := Empty{}
	f := func() Empty { return e }
	_ = f()
	return 1
}

// ZeroValueComparison tests zero value comparison
func ZeroValueComparison() bool {
	var i int
	var s string
	var b bool
	var f float64
	return i == 0 && s == "" && b == false && f == 0.0
}

// NamedTypeZeroValue tests named type zero value
func NamedTypeZeroValue() int {
	type MyInt int
	var x MyInt
	return int(x)
}

// NamedTypeZeroValueComparison tests named type zero value comparison
func NamedTypeZeroValueComparison() bool {
	type MyInt int
	var x MyInt
	return x == 0
}

// SliceZeroValue tests slice zero value
func SliceZeroValue() bool {
	var s []int
	return s == nil
}

// MapZeroValue tests map zero value
func MapZeroValue() bool {
	var m map[int]int
	return m == nil
}

// ChannelZeroValue tests channel zero value
func ChannelZeroValue() bool {
	var ch chan int
	return ch == nil
}

// FuncZeroValue tests function zero value
func FuncZeroValue() bool {
	var f func()
	return f == nil
}

// InterfaceZeroValue tests interface zero value
func InterfaceZeroValue() bool {
	var i interface{}
	return i == nil
}

// PointerZeroValue tests pointer zero value
func PointerZeroValue() bool {
	var p *int
	return p == nil
}

// CompositeLiteralWithZeroValues tests composite literal with zero values
func CompositeLiteralWithZeroValues() int {
	type Data struct {
		A int
		B string
		C bool
	}
	d := Data{}
	return d.A
}

// CompositeLiteralWithPartialValues tests composite literal with partial values
func CompositeLiteralWithPartialValues() int {
	type Data struct {
		A int
		B string
		C bool
	}
	d := Data{A: 42}
	return d.A
}

// NestedCompositeLiteralWithZeroValues tests nested composite literal with zero values
func NestedCompositeLiteralWithZeroValues() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{}
	return o.I.X
}

// SliceLiteralWithZeroElements tests slice literal with zero elements
func SliceLiteralWithZeroElements() int {
	s := []int{}
	return len(s)
}

// MapLiteralWithZeroElements tests map literal with zero elements
func MapLiteralWithZeroElements() int {
	m := map[int]int{}
	return len(m)
}

// ArrayLiteralWithZeroElements tests array literal with zero elements
func ArrayLiteralWithZeroElements() int {
	arr := [0]int{}
	return len(arr)
}

// ============================================================================
// FMT.STRINGER INTERFACE TESTS - Third-party library dependency on String()
// ============================================================================

// StringerBasic tests basic fmt.Stringer implementation
type StringerBasic struct{ Value int }

func (s StringerBasic) String() string {
	return fmt.Sprintf("StringerBasic(%d)", s.Value)
}

// FmtStringerBasic tests fmt.Stringer with value receiver
func FmtStringerBasic() string {
	s := StringerBasic{Value: 42}
	return fmt.Sprintf("%v", s)
}

// StringerPointer tests Stringer with pointer receiver
type StringerPointer struct{ Value int }

func (s *StringerPointer) String() string {
	return fmt.Sprintf("StringerPointer(%d)", s.Value)
}

// FmtStringerPointer tests fmt.Stringer with pointer receiver
func FmtStringerPointer() string {
	s := &StringerPointer{Value: 42}
	return fmt.Sprintf("%v", s)
}

// FmtStringerPointerFromValue tests Stringer pointer method on value
func FmtStringerPointerFromValue() string {
	s := StringerPointer{Value: 42}
	return fmt.Sprintf("%v", &s)
}

// StringerNested tests nested struct with Stringer
type InnerStringer struct{ Value int }

func (i InnerStringer) String() string {
	return fmt.Sprintf("Inner(%d)", i.Value)
}

type OuterStringer struct {
	Inner InnerStringer
	Name  string
}

// FmtStringerNested tests nested Stringer
func FmtStringerNested() string {
	o := OuterStringer{Inner: InnerStringer{Value: 10}, Name: "test"}
	return fmt.Sprintf("%v", o)
}

// StringerInSlice tests Stringer in slice
func StringerInSlice() string {
	items := []StringerBasic{{1}, {2}, {3}}
	return fmt.Sprintf("%v", items)
}

// StringerInMap tests Stringer as map value
func StringerInMap() string {
	m := map[string]StringerBasic{
		"a": {1},
		"b": {2},
	}
	return fmt.Sprintf("%v", m)
}

// StringerInArray tests Stringer in array
func StringerInArray() string {
	arr := [3]StringerBasic{{1}, {2}, {3}}
	return fmt.Sprintf("%v", arr)
}

// StringerAsInterface tests Stringer via interface
func StringerAsInterface() string {
	var s fmt.Stringer = StringerBasic{Value: 42}
	return s.String()
}

// StringerInInterfaceSlice tests Stringer in interface slice
func StringerInInterfaceSlice() string {
	items := []fmt.Stringer{
		StringerBasic{1},
		&StringerPointer{2},
	}
	return fmt.Sprintf("%v %v", items[0], items[1])
}

// StringerEmbedded tests embedded Stringer
type BaseStringer struct{ Value int }

func (b BaseStringer) String() string {
	return fmt.Sprintf("Base(%d)", b.Value)
}

type DerivedStringer struct {
	BaseStringer
	Extra int
}

// FmtStringerEmbedded tests embedded Stringer
func FmtStringerEmbedded() string {
	d := DerivedStringer{BaseStringer: BaseStringer{Value: 10}, Extra: 20}
	return fmt.Sprintf("%v", d)
}

// StringerWithPointerField tests Stringer with pointer field
type StringerWithPointerField struct {
	Value *int
}

func (s StringerWithPointerField) String() string {
	if s.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("Value(%d)", *s.Value)
}

// FmtStringerWithPointerField tests Stringer with pointer field
func FmtStringerWithPointerField() string {
	v := 42
	s := StringerWithPointerField{Value: &v}
	return fmt.Sprintf("%v", s)
}

// StringerWithSliceField tests Stringer with slice field
type StringerWithSliceField struct {
	Items []int
}

func (s StringerWithSliceField) String() string {
	return fmt.Sprintf("Items%v", s.Items)
}

// FmtStringerWithSliceField tests Stringer with slice field
func FmtStringerWithSliceField() string {
	s := StringerWithSliceField{Items: []int{1, 2, 3}}
	return fmt.Sprintf("%v", s)
}

// StringerWithMapField tests Stringer with map field
type StringerWithMapField struct {
	Data map[string]int
}

func (s StringerWithMapField) String() string {
	return fmt.Sprintf("Data%v", s.Data)
}

// FmtStringerWithMapField tests Stringer with map field
func FmtStringerWithMapField() string {
	s := StringerWithMapField{Data: map[string]int{"a": 1}}
	return fmt.Sprintf("%v", s)
}

// StringerRecursive tests recursive Stringer
type StringerRecursive struct {
	Value int
	Next  *StringerRecursive
}

func (s StringerRecursive) String() string {
	if s.Next == nil {
		return fmt.Sprintf("%d", s.Value)
	}
	return fmt.Sprintf("%d->%s", s.Value, s.Next.String())
}

// FmtStringerRecursive tests recursive Stringer
func FmtStringerRecursive() string {
	s3 := &StringerRecursive{Value: 3}
	s2 := &StringerRecursive{Value: 2, Next: s3}
	s1 := &StringerRecursive{Value: 1, Next: s2}
	return fmt.Sprintf("%v", s1)
}

// MultipleStringers tests multiple different Stringer types
type StringerA struct{ Value int }

func (s StringerA) String() string { return fmt.Sprintf("A(%d)", s.Value) }

type StringerB struct{ Value int }

func (s StringerB) String() string { return fmt.Sprintf("B(%d)", s.Value) }

// FmtMultipleStringers tests multiple Stringer types together
func FmtMultipleStringers() string {
	a := StringerA{1}
	b := StringerB{2}
	return fmt.Sprintf("%v %v", a, b)
}

// StringerInStruct tests Stringer as struct field
type Container struct {
	Item StringerBasic
}

// FmtStringerInStruct tests Stringer as struct field
func FmtStringerInStruct() string {
	c := Container{Item: StringerBasic{Value: 42}}
	return fmt.Sprintf("%v", c)
}

// StringerWithMultipleFields tests Stringer with multiple fields
type MultiFieldStringer struct {
	Name  string
	Value int
	Flag  bool
}

func (s MultiFieldStringer) String() string {
	return fmt.Sprintf("%s:%d:%v", s.Name, s.Value, s.Flag)
}

// FmtStringerWithMultipleFields tests Stringer with multiple fields
func FmtStringerWithMultipleFields() string {
	s := MultiFieldStringer{Name: "test", Value: 42, Flag: true}
	return fmt.Sprintf("%v", s)
}

// StringerWithEmptyStruct tests Stringer with empty struct field
type EmptyStructStringer struct {
	Value int
	Empty struct{}
}

func (s EmptyStructStringer) String() string {
	return fmt.Sprintf("Value(%d)", s.Value)
}

// FmtStringerWithEmptyStruct tests Stringer with empty struct field
func FmtStringerWithEmptyStruct() string {
	s := EmptyStructStringer{Value: 42}
	return fmt.Sprintf("%v", s)
}

// StringerZeroValue tests Stringer zero value
// FmtStringerZeroValue tests Stringer zero value
func FmtStringerZeroValue() string {
	var s StringerBasic
	return fmt.Sprintf("%v", s)
}

// StringerNilPointer tests nil Stringer pointer
// FmtStringerNilPointer tests nil Stringer pointer
func FmtStringerNilPointer() string {
	var s *StringerPointer
	return fmt.Sprintf("%v", s)
}

// StringerInVariadic tests Stringer in variadic function
// FmtStringerInVariadic tests Stringer in variadic function
func FmtStringerInVariadic() string {
	items := []fmt.Stringer{StringerBasic{1}, StringerBasic{2}}
	return fmt.Sprintf("%v %v", items[0], items[1])
}

// StringerMethodCall tests direct String() method call
func FmtStringerMethodCall() string {
	s := StringerBasic{Value: 42}
	return s.String()
}

// StringerViaInterface tests Stringer via interface type
func FmtStringerViaInterface() string {
	var i interface{} = StringerBasic{Value: 42}
	if s, ok := i.(fmt.Stringer); ok {
		return s.String()
	}
	return "not a stringer"
}

// StringerInChan tests Stringer in channel
func StringerInChan() string {
	ch := make(chan StringerBasic, 1)
	ch <- StringerBasic{Value: 42}
	s := <-ch
	return fmt.Sprintf("%v", s)
}

// StringerComparison tests Stringer comparison
func StringerComparison() bool {
	s1 := StringerBasic{Value: 42}
	s2 := StringerBasic{Value: 42}
	return s1 == s2
}

// StringerAsMapKey tests Stringer as map key
func StringerAsMapKey() int {
	m := map[StringerBasic]int{
		{1}: 10,
		{2}: 20,
	}
	return m[StringerBasic{1}]
}

// CustomError tests custom error with String()
type CustomError struct {
	Msg string
}

func (e CustomError) Error() string {
	return e.Msg
}

func (e CustomError) String() string {
	return "STRING:" + e.Msg
}

// FmtCustomError tests custom error String() method
func FmtCustomError() string {
	e := CustomError{Msg: "test error"}
	return fmt.Sprintf("%v", e)
}

// StringerWithPrivateFields tests Stringer with private fields
type stringerPrivate struct {
	value int
	name  string
}

func (s stringerPrivate) String() string {
	return fmt.Sprintf("%s:%d", s.name, s.value)
}

// FmtStringerWithPrivateFields tests Stringer with private fields
func FmtStringerWithPrivateFields() string {
	s := stringerPrivate{value: 42, name: "test"}
	return fmt.Sprintf("%v", s)
}

// ============================================================================
// ROUND 8: MORE CORNER CASES - Type conversions, interfaces, closures
// ============================================================================

// TypeConversionOverflow tests type conversion overflow
func TypeConversionOverflow() int8 {
	var x int = 300
	return int8(x) // Truncates to 44 (300 % 256 - 128)
}

// TypeConversionNegative tests negative number conversion
func TypeConversionNegative() uint8 {
	var x int8 = -1
	return uint8(x) // 255
}

// TypeConversionFloatTruncate tests float to int truncation
func TypeConversionFloatTruncate() int {
	x := 3.7
	return int(x) // 3
}

// TypeConversionIntToFloat tests int to float
func TypeConversionIntToFloat() float64 {
	return float64(42)
}

// TypeConversionBoolToInt tests bool can't convert to int (error case)
func TypeConversionBoolToInt() int {
	// Can't convert bool to int in Go, return a value instead
	return 1
}

// InterfaceConversionToInt tests interface to int
func InterfaceConversionToInt() int {
	var i interface{} = 42
	return i.(int)
}

// InterfaceConversionToSlice tests interface to slice
func InterfaceConversionToSlice() int {
	var i interface{} = []int{1, 2, 3}
	return len(i.([]int))
}

// InterfaceConversionToMap tests interface to map
func InterfaceConversionToMap() int {
	var i interface{} = map[string]int{"a": 1}
	return i.(map[string]int)["a"]
}

// InterfaceConversionToFunc tests interface to func
func InterfaceConversionToFunc() int {
	var i interface{} = func(x int) int { return x * 2 }
	return i.(func(int) int)(21)
}

// NilInterfaceTypeAssertion tests nil interface type assertion
func NilInterfaceTypeAssertion() bool {
	var i interface{}
	_, ok := i.(int)
	return !ok
}

// TypedNilInterface tests typed nil to interface
func TypedNilInterface() bool {
	var s []int
	var i interface{} = s
	return i == nil // false! typed nil != nil interface
}

// SliceOfInterfaces tests slice of interfaces
func SliceOfInterfaces() string {
	items := []interface{}{1, "hello", 3.14}
	return fmt.Sprintf("%v", items)
}

// MapOfInterfaces2 tests map of interfaces
func MapOfInterfaces2() string {
	m := map[string]interface{}{
		"int":    1,
		"string": "hello",
		"float":  3.14,
	}
	return fmt.Sprintf("%v", m)
}

// ClosureWithMultipleCaptures tests closure capturing multiple variables
func ClosureWithMultipleCaptures() int {
	a, b, c := 1, 2, 3
	f := func() int {
		return a + b + c
	}
	return f()
}

// ClosureWithNestedCapture tests nested closure capture
func ClosureWithNestedCapture() int {
	x := 10
	f1 := func() int {
		y := 20
		f2 := func() int {
			return x + y
		}
		return f2()
	}
	return f1()
}

// ClosureModifyCapture tests closure modifying captured variable
func ClosureModifyCapture() int {
	x := 10
	f := func() {
		x = 20
	}
	f()
	return x
}

// ClosureReturnClosure tests closure returning closure
func ClosureReturnClosure() int {
	f := func() func() int {
		x := 42
		return func() int {
			return x
		}
	}
	return f()()
}

// ClosureTakeClosure tests closure taking closure as argument
func ClosureTakeClosure() int {
	apply := func(f func() int) int {
		return f()
	}
	return apply(func() int { return 42 })
}

// ClosureInMap tests closure in map
func ClosureInMap() int {
	m := map[string]func() int{
		"a": func() int { return 1 },
		"b": func() int { return 2 },
	}
	return m["a"]() + m["b"]()
}

// ClosureInSlice tests closure in slice
func ClosureInSlice() int {
	items := []func() int{
		func() int { return 1 },
		func() int { return 2 },
	}
	return items[0]() + items[1]()
}

// ClosureAsMapKey tests closure can't be map key (use int key)
func ClosureAsMapKey() int {
	// Functions can't be map keys, use int keys instead
	m := map[int]string{
		1: "one",
		2: "two",
	}
	return len(m)
}

// ChannelSendRecvOrder tests channel send/receive order
func ChannelSendRecvOrder() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	return <-ch + <-ch // 1 + 2 = 3
}

// ChannelCloseThenRecv tests receiving from closed channel
func ChannelCloseThenRecv() (int, bool) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	v, ok := <-ch
	return v, ok
}

// ChannelClosedRecvZero tests receiving from closed empty channel
func ChannelClosedRecvZero() int {
	ch := make(chan int)
	close(ch)
	v, ok := <-ch
	_ = ok
	return v // zero value
}

// NilChannelBlocks tests nil channel blocks (test with select)
func NilChannelBlocks() string {
	var ch chan int
	select {
	case <-ch:
		return "received"
	default:
		return "default"
	}
}

// ChannelOfChannels tests channel of channels
func ChannelOfChannelsSend() int {
	ch := make(chan chan int, 1)
	subCh := make(chan int, 1)
	subCh <- 42
	ch <- subCh
	return <-(<-ch)
}

// BufferedChannelFull tests buffered channel full behavior
func BufferedChannelFull() int {
	ch := make(chan int, 1)
	ch <- 1
	select {
	case ch <- 2:
		return 2
	default:
		return 1 // channel full, default case
	}
}

// StructWithGetMethod tests struct with Get method
type Getter struct{ Value int }

func (g Getter) Get() int { return g.Value }

// StructGetMethod tests struct get method
func StructGetMethod() int {
	g := Getter{Value: 42}
	return g.Get()
}

// StructWithSetMethod tests struct with Set method
type Setter struct{ Value int }

func (s *Setter) Set(v int) { s.Value = v }

// StructSetMethod tests struct set method
func StructSetMethod() int {
	s := Setter{Value: 0}
	s.Set(42)
	return s.Value
}

// StructWithBothMethods tests struct with both Get and Set
type Accessor struct{ Value int }

func (a Accessor) Get() int   { return a.Value }
func (a *Accessor) Set(v int) { a.Value = v }

// StructBothMethods tests both methods
func StructBothMethods() int {
	a := Accessor{Value: 10}
	a.Set(42)
	return a.Get()
}

// InterfaceWithMultipleMethods2 tests interface with multiple methods
type Reader2 interface{ Read() int }
type Writer2 interface{ Write(int) }
type ReadWriter2 interface {
	Reader2
	Writer2
}

type RWImpl struct{ Value int }

func (r RWImpl) Read() int    { return r.Value }
func (r *RWImpl) Write(v int) { r.Value = v }

// InterfaceMultipleMethods tests multiple interface methods
func InterfaceMultipleMethods() int {
	rw := &RWImpl{Value: 10}
	rw.Write(42)
	return rw.Read()
}

// EmptyInterfaceWithMethods tests empty interface
func EmptyInterfaceWithMethods() int {
	var i interface{} = 42
	return i.(int)
}

// InterfaceEmbeddingMultiple tests multiple interface embedding
type ReaderA interface{ ReadA() int }
type ReaderB interface{ ReadB() int }
type CombinedReader interface {
	ReaderA
	ReaderB
}

type CombinedImpl struct{ A, B int }

func (c CombinedImpl) ReadA() int { return c.A }
func (c CombinedImpl) ReadB() int { return c.B }

// InterfaceEmbeddingMultipleTest tests multiple interface embedding
func InterfaceEmbeddingMultipleTest() int {
	c := CombinedImpl{A: 1, B: 2}
	return c.ReadA() + c.ReadB()
}

// MethodOnPointerType tests method on pointer type
type Counter struct{ Count int }

func (c *Counter) Increment() { c.Count++ }
func (c Counter) Value() int  { return c.Count }

// MethodOnPointerTypeTest tests pointer method
func MethodOnPointerTypeTest() int {
	c := &Counter{Count: 0}
	c.Increment()
	c.Increment()
	return c.Value()
}

// MethodOnValueType tests method on value type
func MethodOnValueTypeTest() int {
	c := Counter{Count: 10}
	return c.Value()
}

// MethodPointerOnValue tests pointer method on value (auto-address)
func MethodPointerOnValueTest() int {
	c := Counter{Count: 10}
	c.Increment() // Auto-address taken
	return c.Value()
}

// SliceAppendMake tests append with make
func SliceAppendMake() int {
	s := make([]int, 0, 10)
	s = append(s, 1, 2, 3)
	return len(s)
}

// MapMakeDelete tests make and delete
func MapMakeDelete() int {
	m := make(map[int]int)
	m[1] = 10
	m[2] = 20
	delete(m, 1)
	return len(m)
}

// SliceCopyMake tests copy with make
func SliceCopyMake() int {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	return dst[0] + dst[1] + dst[2]
}

// NilSliceAppendNil tests nil slice append nil
func NilSliceAppendNil2() int {
	var s []int
	s = append(s, nil...)
	return len(s)
}

// SliceAppendFunc tests append with function result
func SliceAppendFunc() int {
	getSlice := func() []int { return []int{3, 4} }
	s := []int{1, 2}
	s = append(s, getSlice()...)
	return len(s)
}

// MapWithFuncKey tests map with func key (not allowed, use string)
func MapWithFuncKey() int {
	// Func can't be map key, use string
	m := map[string]int{"key": 42}
	return m["key"]
}

// StructWithFuncFieldMethod tests struct with func field and method
type FuncHolder struct {
	Fn func() int
}

func (f FuncHolder) Call() int { return f.Fn() }

// StructWithFuncFieldMethodTest tests struct with func field
func StructWithFuncFieldMethodTest() int {
	f := FuncHolder{Fn: func() int { return 42 }}
	return f.Call()
}

// ============================================================================
// ROUND 9: MORE CORNER CASES - Type switches, embedded fields, method values
// ============================================================================

// TypeSwitchWithFallthrough tests type switch with fallthrough (not allowed)
func TypeSwitchWithFallthrough() string {
	var i interface{} = 42
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case string:
		return fmt.Sprintf("string: %s", v)
	default:
		return "unknown"
	}
}

// TypeSwitchMultipleInOne tests multiple types in one case
func TypeSwitchMultipleInOne() string {
	var i interface{} = 42
	switch i.(type) {
	case int, int8, int16, int32, int64:
		return "int type"
	case uint, uint8, uint16, uint32, uint64:
		return "uint type"
	default:
		return "other"
	}
}

// EmbeddedInner for EmbeddedFieldAccess test
type EmbeddedInner struct{ X int }

// EmbeddedOuter for EmbeddedFieldAccess test
type EmbeddedOuter struct {
	EmbeddedInner
	Y int
}

// EmbeddedFieldAccess tests accessing embedded field
func EmbeddedFieldAccess() int {
	o := EmbeddedOuter{EmbeddedInner: EmbeddedInner{X: 10}, Y: 20}
	return o.X + o.Y // X is promoted from EmbeddedInner
}

// EmbeddedBase for EmbeddedMethodAccess test
type EmbeddedBase struct{ Value int }

// GetValue method for EmbeddedBase
func (b EmbeddedBase) GetValue() int { return b.Value }

// EmbeddedDerived for EmbeddedMethodAccess test
type EmbeddedDerived struct {
	EmbeddedBase
	Extra int
}

// EmbeddedMethodAccess tests accessing embedded method
func EmbeddedMethodAccess() int {
	d := EmbeddedDerived{EmbeddedBase: EmbeddedBase{Value: 10}, Extra: 20}
	return d.GetValue() + d.Extra // GetValue is promoted
}

// EmbeddedPtrInner for EmbeddedPointerField test
type EmbeddedPtrInner struct{ X int }

// EmbeddedPtrOuter for EmbeddedPointerField test
type EmbeddedPtrOuter struct {
	*EmbeddedPtrInner
	Y int
}

// EmbeddedPointerField tests embedded pointer field
func EmbeddedPointerField() int {
	o := EmbeddedPtrOuter{EmbeddedPtrInner: &EmbeddedPtrInner{X: 10}, Y: 20}
	return o.X + o.Y
}

// EmbeddedPtrBase for EmbeddedPointerMethod test
type EmbeddedPtrBase struct{ Value int }

// GetPtrValue method for EmbeddedPtrBase
func (b *EmbeddedPtrBase) GetPtrValue() int { return b.Value }

// EmbeddedPtrDerived for EmbeddedPointerMethod test
type EmbeddedPtrDerived struct {
	*EmbeddedPtrBase
	Extra int
}

// EmbeddedPointerMethod tests embedded pointer method
func EmbeddedPointerMethod() int {
	d := EmbeddedPtrDerived{EmbeddedPtrBase: &EmbeddedPtrBase{Value: 10}, Extra: 20}
	return d.GetPtrValue() + d.Extra
}

// Reader3 for EmbeddedInterfaceField test
type Reader3 interface{ Read() int }

// Writer3 for EmbeddedInterfaceField test
type Writer3 interface{ Write(int) }

// ReadWriter3 for EmbeddedInterfaceField test
type ReadWriter3 interface {
	Reader3
	Writer3
}

// Impl3 for EmbeddedInterfaceField test
type Impl3 struct{ Value int }

// Read method for Impl3
func (i Impl3) Read() int { return i.Value }

// Write method for Impl3
func (i *Impl3) Write(v int) { i.Value = v }

// EmbeddedInterfaceField tests embedded interface
func EmbeddedInterfaceField() int {
	var rw ReadWriter3 = &Impl3{Value: 42}
	return rw.Read()
}

// Counter2 for MethodValue test
type Counter2 struct{ Count int }

// Increment method for Counter2
func (c *Counter2) Increment() { c.Count++ }

// MethodValueTest2 tests method value
func MethodValueTest2() int {
	c := &Counter2{Count: 0}
	inc := c.Increment
	inc()
	inc()
	return c.Count
}

// Counter3MethodExpr for MethodExpressionTest test
type Counter3MethodExpr struct{ Count int }

// Inc method for Counter3MethodExpr
func (c *Counter3MethodExpr) Inc() { c.Count++ }

// MethodExpressionTest2 tests method expression
func MethodExpressionTest2() int {
	c := &Counter3MethodExpr{Count: 10}
	inc := (*Counter3MethodExpr).Inc
	inc(c)
	return c.Count
}

// Counter4 for MethodValueCapturesReceiver test
type Counter4 struct{ Count int }

// GetCount method for Counter4
func (c Counter4) GetCount() int { return c.Count }

// MethodValueCapturesReceiver tests method value captures receiver
func MethodValueCapturesReceiver() int {
	c := Counter4{Count: 42}
	v := c.GetCount
	return v()
}

// Adder for SliceOfMethodValues test
type Adder struct{ Value int }

// Add method for Adder
func (a Adder) Add(x int) int { return a.Value + x }

// SliceOfMethodValues tests slice of method values
func SliceOfMethodValues() int {
	a := Adder{Value: 10}
	adds := []func(int) int{a.Add, a.Add}
	return adds[0](1) + adds[1](2)
}

// MapOfMethodValues tests map of method values
func MapOfMethodValues() int {
	a := Adder{Value: 10}
	m := map[string]func(int) int{
		"add": a.Add,
	}
	return m["add"](5)
}

// NilSliceLen tests nil slice len
func NilSliceLen2() int {
	var s []int
	return len(s)
}

// NilSliceCap tests nil slice cap
func NilSliceCap2() int {
	var s []int
	return cap(s)
}

// NilMapLen tests nil map len
func NilMapLen2() int {
	var m map[int]int
	return len(m)
}

// NilMapDeleteTest tests nil map delete (no-op)
func NilMapDeleteTest() int {
	var m map[int]int
	delete(m, 1) // no-op on nil map
	return 0
}

// EmptySliceLen tests empty slice len
func EmptySliceLen() int {
	s := []int{}
	return len(s)
}

// EmptyMapLen tests empty map len
func EmptyMapLen() int {
	m := map[int]int{}
	return len(m)
}

// SliceMakeZeroLen tests make with zero len
func SliceMakeZeroLen() int {
	s := make([]int, 0)
	return len(s)
}

// MapMakeZeroSize tests make with zero size
func MapMakeZeroSize() int {
	m := make(map[int]int)
	return len(m)
}

// FuncsHolder for StructWithSliceOfFuncs test
type FuncsHolder struct {
	Funcs []func() int
}

// StructWithSliceOfFuncs tests struct with slice of funcs
func StructWithSliceOfFuncs() int {
	f := FuncsHolder{
		Funcs: []func() int{
			func() int { return 1 },
			func() int { return 2 },
		},
	}
	return f.Funcs[0]() + f.Funcs[1]()
}

// FuncsMapHolder for StructWithMapOfFuncs test
type FuncsMapHolder struct {
	Funcs map[string]func() int
}

// StructWithMapOfFuncs tests struct with map of funcs
func StructWithMapOfFuncs() int {
	f := FuncsMapHolder{
		Funcs: map[string]func() int{
			"a": func() int { return 1 },
			"b": func() int { return 2 },
		},
	}
	return f.Funcs["a"]() + f.Funcs["b"]()
}

// InnerWithSlice for NestedStructWithSlice test
type InnerWithSlice struct{ Items []int }

// OuterWithSlice for NestedStructWithSlice test
type OuterWithSlice struct{ I InnerWithSlice }

// NestedStructWithSlice tests nested struct with slice
func NestedStructWithSlice() int {
	o := OuterWithSlice{I: InnerWithSlice{Items: []int{1, 2, 3}}}
	return len(o.I.Items)
}

// InnerWithMap for NestedStructWithMap test
type InnerWithMap struct{ Data map[int]int }

// OuterWithMap for NestedStructWithMap test
type OuterWithMap struct{ I InnerWithMap }

// NestedStructWithMap tests nested struct with map
func NestedStructWithMap() int {
	o := OuterWithMap{I: InnerWithMap{Data: map[int]int{1: 10}}}
	return o.I.Data[1]
}

// DataForModify for PointerToStructModify test
type DataForModify struct{ Value int }

// PointerToSliceModify tests modifying slice via pointer
func PointerToSliceModify() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 10
	return s[0]
}

// PointerToMapModify tests modifying map via pointer
func PointerToMapModify() int {
	m := map[int]int{1: 10}
	p := &m
	(*p)[2] = 20
	return len(m)
}

// PointerToStructModify tests modifying struct via pointer
func PointerToStructModify() int {
	d := DataForModify{Value: 10}
	p := &d
	p.Value = 20
	return d.Value
}

// SliceOfPointersModify tests modifying via slice of pointers
func SliceOfPointersModify() int {
	items := []*int{}
	v := 42
	items = append(items, &v)
	*items[0] = 100
	return v
}

// MapOfPointersModify tests modifying via map of pointers
func MapOfPointersModify() int {
	m := map[string]*int{}
	v := 42
	m["key"] = &v
	*m["key"] = 100
	return v
}

// InterfaceSliceTypeAssertion tests type assertion on interface slice
func InterfaceSliceTypeAssertion() int {
	var i interface{} = []int{1, 2, 3}
	s := i.([]int)
	return len(s)
}

// InterfaceMapTypeAssertion tests type assertion on interface map
func InterfaceMapTypeAssertion() int {
	var i interface{} = map[int]int{1: 2}
	m := i.(map[int]int)
	return m[1]
}

// InterfaceFuncTypeAssertion tests type assertion on interface func
func InterfaceFuncTypeAssertion() int {
	var i interface{} = func(x int) int { return x * 2 }
	f := i.(func(int) int)
	return f(21)
}

// InterfaceChanTypeAssertion tests type assertion on interface chan
func InterfaceChanTypeAssertion() int {
	ch := make(chan int, 1)
	ch <- 42
	var i interface{} = ch
	c := i.(chan int)
	return <-c
}

// MultipleInterfaceTypeAssertion tests multiple interface type assertions
func MultipleInterfaceTypeAssertion() int {
	var i interface{} = 42
	if v, ok := i.(int); ok {
		return v
	}
	return 0
}

// MyIntNamed for TypeAssertionOnNamed test
type MyIntNamed int

// TypeAssertionOnNamed tests type assertion on named type
func TypeAssertionOnNamed() int {
	var i interface{} = MyIntNamed(42)
	v := i.(MyIntNamed)
	return int(v)
}

// DataForAssert for TypeAssertionOnStruct test
type DataForAssert struct{ Value int }

// TypeAssertionOnStruct tests type assertion on struct
func TypeAssertionOnStruct() int {
	var i interface{} = DataForAssert{Value: 42}
	v := i.(DataForAssert)
	return v.Value
}

// TypeAssertionOnPointer tests type assertion on pointer
func TypeAssertionOnPointer() int {
	var i interface{} = &DataForAssert{Value: 42}
	v := i.(*DataForAssert)
	return v.Value
}

// ============================================================================
// ROUND 10: MORE CORNER CASES - Concurrency, complex types, edge cases
// ============================================================================

// ChannelBidirectional tests bidirectional channel
func ChannelBidirectional() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	return <-ch + <-ch
}

// ChannelSendOnly tests send-only channel
func ChannelSendOnly() int {
	ch := make(chan int, 1)
	var sendCh chan<- int = ch
	sendCh <- 42
	return <-ch
}

// ChannelRecvOnly tests receive-only channel
func ChannelRecvOnly() int {
	ch := make(chan int, 1)
	ch <- 42
	var recvCh <-chan int = ch
	return <-recvCh
}

// SelectNonBlockingDefault tests non-blocking select with default
func SelectNonBlockingDefault2() int {
	ch := make(chan int)
	select {
	case v := <-ch:
		return v
	default:
		return 42
	}
}

// BufferedChannelLen tests buffered channel len
func BufferedChannelLen() int {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	return len(ch)
}

// StructWithNilPointerField tests struct with nil pointer field
func StructWithNilPointerField() bool {
	type Data struct{ Value *int }
	d := Data{Value: nil}
	return d.Value == nil
}

// StructWithNilSliceField tests struct with nil slice field
func StructWithNilSliceField() bool {
	type Data struct{ Items []int }
	d := Data{Items: nil}
	return d.Items == nil
}

// InterfaceNilVsTypedNil tests interface nil vs typed nil
func InterfaceNilVsTypedNil() bool {
	var i interface{} = nil
	var s []int = nil
	var j interface{} = s
	return i == nil && j != nil
}

// SliceOfNilInterfaces tests slice of nil interfaces
func SliceOfNilInterfaces() int {
	items := []interface{}{nil, nil, nil}
	count := 0
	for _, item := range items {
		if item == nil {
			count++
		}
	}
	return count
}

// NilInterfaceTypeSwitch tests type switch on nil interface
func NilInterfaceTypeSwitch() string {
	var i interface{}
	switch i.(type) {
	case int:
		return "int"
	case string:
		return "string"
	default:
		return "nil or unknown"
	}
}

// StructComparisonWithNilPointer tests struct comparison with nil pointer
func StructComparisonWithNilPointer() bool {
	type Data struct{ P *int }
	d1 := Data{P: nil}
	d2 := Data{P: nil}
	return d1 == d2
}

// StructWithSameTypeFields tests struct with multiple same type fields
func StructWithSameTypeFields() int {
	type Data struct{ A, B, C int }
	d := Data{A: 1, B: 2, C: 3}
	return d.A + d.B + d.C
}

// StructWithMixedTypeFields tests struct with mixed type fields
func StructWithMixedTypeFields() int {
	type Data struct {
		A int
		B string
		C bool
		D float64
	}
	d := Data{A: 1, B: "test", C: true, D: 3.14}
	return d.A + len(d.B)
}

// TimeDurationOperation tests time.Duration operations
func TimeDurationOperation() int {
	d := 100 * time.Millisecond
	return int(d / time.Millisecond)
}

// TimeNowOperation tests time.Now() operation
func TimeNowOperation() bool {
	t := time.Now()
	return !t.IsZero()
}

// ============================================================================
// ROUND 11: MORE CORNER CASES - Error handling, complex types, more edge cases
// ============================================================================

// ErrorTypeAssertion tests error type assertion
func ErrorTypeAssertion() string {
	var err error
	_, ok := err.(*strconv.NumError)
	_ = ok
	return "ok"
}

// ErrorNilComparison tests nil error comparison
func ErrorNilComparison() bool {
	var err error
	return err == nil
}

// ErrorWithNil tests error with nil value
func ErrorWithNil() bool {
	var err error = nil
	return err == nil
}

// ErrorFromFunc tests error returned from function
func ErrorFromFunc() bool {
	f := func() error { return nil }
	return f() == nil
}

// SliceAppendVariadic tests append with variadic
func SliceAppendVariadic() int {
	s := []int{1, 2}
	s = append(s, 3, 4, 5)
	return len(s)
}

// SliceAppendSlice tests append with slice
func SliceAppendSlice() int {
	s1 := []int{1, 2}
	s2 := []int{3, 4}
	s1 = append(s1, s2...)
	return len(s1)
}

// MapIterateAndDelete tests iterating map and deleting
func MapIterateAndDelete() int {
	m := map[int]int{1: 1, 2: 2, 3: 3}
	for k := range m {
		if k == 2 {
			delete(m, k)
		}
	}
	return len(m)
}

// StructWithEmptyInterface tests struct with empty interface field
func StructWithEmptyInterface() int {
	type Data struct {
		Value interface{}
	}
	d := Data{Value: 42}
	return d.Value.(int)
}

// simpleReader for StructWithTwoInterfaces test - defined before use
type simpleReader struct{ val int }

func (s *simpleReader) Read() int { return s.val }

// StructWithTwoInterfaces tests struct with two interface fields
func StructWithTwoInterfaces() int {
	type Data struct {
		Reader interface{ Read() int }
	}
	sr := &simpleReader{val: 42}
	d := Data{Reader: sr}
	return d.Reader.Read()
}

// ComplexRealImag tests complex real/imag
func ComplexRealImag() float64 {
	c := complex(3.0, 4.0)
	return real(c) + imag(c)
}

// ComplexFromRealImag tests complex from real/imag
func ComplexFromRealImag() complex128 {
	return complex(3.0, 4.0)
}

// ComplexOperations tests complex operations
func ComplexOperations() float64 {
	c1 := complex(1, 2)
	c2 := complex(3, 4)
	c3 := c1 + c2
	return real(c3) + imag(c3)
}

// StringCompareRound11 tests string comparison
func StringCompareRound11() bool {
	s1 := "hello"
	s2 := "hello"
	return s1 == s2
}

// StringToByteSliceRound11 tests string to byte slice
func StringToByteSliceRound11() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}

// ByteSliceToStringRound11 tests byte slice to string
func ByteSliceToStringRound11() string {
	b := []byte{'h', 'i'}
	return string(b)
}

// RuneSliceToString tests rune slice to string
func RuneSliceToString() string {
	r := []rune{'h', 'i'}
	return string(r)
}

// StringToRuneSlice tests string to rune slice
func StringToRuneSlice() int {
	s := "hello"
	r := []rune(s)
	return len(r)
}

// RangeOverStringCount tests range over string counting runes
func RangeOverStringCount() int {
	s := "hello"
	count := 0
	for range s {
		count++
	}
	return count
}

// RangeOverStringIndex tests range over string with index
func RangeOverStringIndex() int {
	s := "hello"
	sum := 0
	for i := range s {
		sum += i
	}
	return sum
}

// RangeOverStringRune tests range over string with index and rune
func RangeOverStringRune() int {
	s := "hello"
	count := 0
	for _, r := range s {
		if r == 'l' {
			count++
		}
	}
	return count
}

// ============================================================================
// ROUND 12: MORE CORNER CASES - Type aliases, unsafe, more edge cases
// ============================================================================

// TypeAliasBasic tests basic type alias
func TypeAliasBasicR12() int {
	type MyInt = int
	var x MyInt = 42
	return x
}

// TypeAliasSlice tests type alias for slice
func TypeAliasSliceR12() int {
	type IntSlice = []int
	var s IntSlice = []int{1, 2, 3}
	return len(s)
}

// TypeAliasMap tests type alias for map
func TypeAliasMapR12() int {
	type StringMap = map[string]int
	var m StringMap = map[string]int{"a": 1}
	return m["a"]
}

// TypeAliasFunc tests type alias for func
func TypeAliasFuncR12() int {
	type IntFunc = func(int) int
	var f IntFunc = func(x int) int { return x * 2 }
	return f(21)
}

// TypeAliasStruct tests type alias for struct
func TypeAliasStructR12() int {
	type Point = struct{ X, Y int }
	var p Point = struct{ X, Y int }{X: 1, Y: 2}
	return p.X + p.Y
}

// TypeAliasPointer tests type alias for pointer
func TypeAliasPointerR12() int {
	type IntPtr = *int
	x := 42
	var p IntPtr = &x
	return *p
}

// TypeAliasChan tests type alias for channel
func TypeAliasChan() int {
	type IntChan = chan int
	var ch IntChan = make(chan int, 1)
	ch <- 42
	return <-ch
}

// TypeAliasInterface tests type alias for interface
func TypeAliasInterface() int {
	type Stringer = interface{ String() string }
	return 1
}

// NestedTypeAlias tests nested type alias
func NestedTypeAlias() int {
	type Int = int
	type IntPtr = *Int
	type IntPtrSlice = []IntPtr
	x := Int(42)
	s := IntPtrSlice{&x}
	return *s[0]
}

// StructWithTag tests struct with tags
func StructWithTag() string {
	type Data struct {
		Value int `json:"value"`
	}
	d := Data{Value: 42}
	_ = d
	return "ok"
}

// MultipleTags tests multiple tags
func MultipleTags() string {
	type Data struct {
		Value int `json:"value" xml:"value"`
	}
	d := Data{Value: 42}
	_ = d
	return "ok"
}

// StructWithOmitEmpty tests struct with omitempty tag
func StructWithOmitEmpty() string {
	type Data struct {
		Value int `json:"value,omitempty"`
	}
	d := Data{Value: 0}
	_ = d
	return "ok"
}

// BlankIdentifierInVar tests blank identifier in var
func BlankIdentifierInVar() int {
	_, b, _ := 1, 2, 3
	return b
}

// BlankIdentifierInFor tests blank identifier in for
func BlankIdentifierInFor() int {
	s := []int{1, 2, 3}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// BlankIdentifierInImport tests blank identifier import
func BlankIdentifierInImport() int {
	// strconv imported as blank at top
	return 1
}

// BlankIdentifierInReturn tests blank identifier in return
func BlankIdentifierInReturn() (int, int) {
	return 1, 2
}

// BlankIdentifierInTypeSwitch tests blank identifier in type switch
func BlankIdentifierInTypeSwitch() string {
	var i interface{} = 42
	switch i.(type) {
	case int:
		return "int"
	default:
		return "other"
	}
}

// NamedReturnWithDeferModify tests named return with defer modification
func NamedReturnWithDeferModify() (result int) {
	defer func() { result++ }()
	return 10
}

// NamedReturnMultiple tests multiple named returns
func NamedReturnMultiple() (a, b int) {
	a = 1
	b = 2
	return
}

// LiteralEllipsis tests literal with ellipsis
func LiteralEllipsis() int {
	arr := [...]int{1, 2, 3}
	return len(arr)
}

// ArrayLiteralEllipsis tests array literal with ellipsis
func ArrayLiteralEllipsis() int {
	arr := [...]int{1, 2: 10}
	return arr[0] + arr[2]
}

// SliceLiteralFromArr tests slice literal from array
func SliceLiteralFromArr() int {
	arr := [3]int{1, 2, 3}
	s := arr[:]
	return s[0] + s[1] + s[2]
}

// ArrayPointerLiteral tests array pointer literal
func ArrayPointerLiteral() int {
	arr := &[3]int{1, 2, 3}
	return arr[0] + arr[1] + arr[2]
}

// StructPointerLiteral tests struct pointer literal
func StructPointerLiteral() int {
	type Data struct{ Value int }
	d := &Data{Value: 42}
	return d.Value
}

// MapLiteralWithStructKey tests map with struct key
func MapLiteralWithStructKey2() int {
	type Key struct{ X, Y int }
	m := map[Key]int{
		{1, 2}: 10,
		{3, 4}: 20,
	}
	return m[Key{1, 2}]
}

// SliceLiteralWithMaxIndex tests slice literal with max index
func SliceLiteralWithMaxIndex() int {
	s := []int{100: 42}
	return len(s)
}

// ArrayLiteralWithMaxIndex tests array literal with max index
func ArrayLiteralWithMaxIndex() int {
	arr := [101]int{100: 42}
	return arr[100]
}

// ConstExpression tests constant expression
func ConstExpressionR12() int {
	const x = 1 + 2*3
	return x
}

// ConstIota tests iota constant
func ConstIota() int {
	const (
		a = iota
		b
		c
	)
	return a + b + c
}

// ConstIotaExpression tests iota with expression
func ConstIotaExpression() int {
	const (
		a = 1 << iota
		b
		c
	)
	return a + b + c
}

// ConstIotaSkip tests iota with skip
func ConstIotaSkip() int {
	const (
		a = iota
		_ = iota
		b = iota
	)
	return a + b
}

// VarBlock tests var block
func VarBlock() int {
	var (
		a = 1
		b = 2
		c = 3
	)
	return a + b + c
}

// ConstBlock tests const block
func ConstBlock() int {
	const (
		a = 1
		b = 2
		c = 3
	)
	return a + b + c
}

// TypeBlock tests type block
func TypeBlock() int {
	type (
		Point struct{ X, Y int }
		Rect  struct{ Min, Max Point }
	)
	r := Rect{Min: Point{X: 0, Y: 0}, Max: Point{X: 10, Y: 10}}
	return r.Max.X
}

// ShortVarDeclInIf tests short var decl in if
func ShortVarDeclInIf() int {
	if x := 42; x > 0 {
		return x
	}
	return 0
}

// ShortVarDeclInSwitch tests short var decl in switch
func ShortVarDeclInSwitch() string {
	switch x := 1; x {
	case 1:
		return "one"
	default:
		return "other"
	}
}

// ShortVarDeclInFor tests short var decl in for
func ShortVarDeclInFor() int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += i
	}
	return sum
}

// ShortVarDeclInSelect tests short var decl in select
func ShortVarDeclInSelect() int {
	ch := make(chan int, 1)
	ch <- 42
	select {
	case v := <-ch:
		return v
	default:
		return 0
	}
}

// ExpressionStatement tests expression statement
func ExpressionStatement() int {
	x := 1
	x++
	return x
}

// IncDecStatement tests inc/dec statement
func IncDecStatement() int {
	x := 10
	x++
	x--
	return x
}

// AssignmentStatement tests assignment statement
func AssignmentStatement() int {
	x, y := 1, 2
	x, y = y, x
	return x + y
}

// AssignmentWithOp tests assignment with operation
func AssignmentWithOp() int {
	x := 10
	x += 5
	x -= 3
	x *= 2
	return x
}

// SendStatement tests send statement
func SendStatement() int {
	ch := make(chan int, 1)
	ch <- 42
	return <-ch
}

// RangeStatement tests range statement
func RangeStatement() int {
	s := []int{1, 2, 3}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// DeferStatement tests defer statement
func DeferStatement() int {
	x := 1
	defer func() { x++ }()
	return x
}

// MultipleDefer tests multiple defer
func MultipleDefer() int {
	x := 0
	defer func() { x += 1 }()
	defer func() { x += 2 }()
	defer func() { x += 4 }()
	return x
}

// ============================================================================
// MORE EDGE CASES (Round 7) - Trying to find bugs
// ============================================================================

// DeepNestedClosureChain tests deeply nested closure returning closure
func DeepNestedClosureChain() int {
	f1 := func(x int) func(int) func(int) int {
		return func(y int) func(int) int {
			return func(z int) int {
				return x + y + z
			}
		}
	}
	return f1(1)(2)(3)
}

// ClosureMutatingCapturedVars tests closure mutating multiple captured vars
func ClosureMutatingCapturedVars() int {
	x, y, z := 1, 2, 3
	f := func() {
		x *= 10
		y *= 20
		z *= 30
	}
	g := func() {
		x++
		y++
		z++
	}
	f()
	g()
	return x + y + z // 10+1 + 40+1 + 90+1 = 143
}

// RecursiveClosureWithDefer tests recursive closure with defer
func RecursiveClosureWithDefer() int {
	result := 0
	var f func(int)
	f = func(n int) {
		if n <= 0 {
			return
		}
		defer func() {
			result += n
		}()
		f(n - 1)
	}
	f(5)
	return result // 1+2+3+4+5 = 15
}

// DeferInRecursiveCall tests defer in recursive function
func DeferInRecursiveCall() int {
	var helper func(int) int
	helper = func(n int) int {
		if n <= 0 {
			return 0
		}
		defer func() {
			// This should execute in reverse order
		}()
		return n + helper(n-1)
	}
	return helper(10) // 55
}

// DeferModifyingNamedReturns tests defer modifying named returns
func DeferModifyingNamedReturns() (result int) {
	defer func() {
		result *= 2
	}()
	defer func() {
		result += 10
	}()
	result = 5
	return // Returns (5 + 10) * 2 = 30
}

// DeferInClosureReturning tests defer in closure that returns
func DeferInClosureReturning() int {
	f := func() (result int) {
		defer func() {
			result = 999
		}()
		return 0
	}
	return f()
}

// LinkedListCycle tests linked list with cycle
func LinkedListCycle() int {
	type Node struct {
		Value int
		Next  *Node
	}

	n1 := &Node{Value: 1}
	n2 := &Node{Value: 2}
	n3 := &Node{Value: 3}
	n1.Next = n2
	n2.Next = n3
	n3.Next = n1 // Cycle

	// Traverse up to 10 nodes to avoid infinite loop
	sum := 0
	current := n1
	for i := 0; i < 10; i++ {
		sum += current.Value
		current = current.Next
	}
	return sum // 1+2+3+1+2+3+1+2+3+1 = 19
}

// TreeWithParentRef tests tree with parent references
// ChainedTypeAssertions tests chained type assertions
func ChainedTypeAssertions() int {
	var i interface{} = "hello"

	// Chain of assertions
	s, ok := i.(string)
	if !ok {
		return -1
	}

	var j interface{} = s
	s2, ok2 := j.(string)
	if !ok2 {
		return -2
	}

	return len(s2)
}

// NestedInterfaceAssertions tests nested interface assertions
func NestedInterfaceAssertions() int {
	var outer interface{}
	var inner interface{} = 42
	outer = inner

	// Assert twice
	i1, ok1 := outer.(interface{})
	if !ok1 {
		return -1
	}
	i2, ok2 := i1.(int)
	if !ok2 {
		return -2
	}
	return i2
}

// TypeSwitchWithMultipleTypes tests type switch with many types
func TypeSwitchWithMultipleTypes() int {
	values := []interface{}{42, "hello", 3.14, true, []int{1, 2, 3}}
	result := 0
	for _, v := range values {
		switch v.(type) {
		case int:
			result += 1
		case string:
			result += 2
		case float64:
			result += 3
		case bool:
			result += 4
		default:
			result += 5
		}
	}
	return result // 1 + 2 + 3 + 4 + 5 = 15
}

// MapWithComplexKeys2 tests map with struct keys
func MapWithComplexKeys2() int {
	type Key struct {
		X, Y int
	}
	m := make(map[Key]int)
	m[Key{1, 2}] = 10
	m[Key{3, 4}] = 20

	return m[Key{1, 2}] + m[Key{3, 4}]
}

// NestedMapOperations tests nested maps
func NestedMapOperations() int {
	m := make(map[string]map[string]int)
	m["outer"] = make(map[string]int)
	m["outer"]["inner"] = 42

	return m["outer"]["inner"]
}

// MapIterationWithDelete tests map iteration with deletion
func MapIterationWithDelete() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for k, v := range m {
		if v > 1 {
			delete(m, k)
		}
		count++
	}
	return len(m) + count // Remaining items + iterations
}

// NestedSliceAppend tests nested slice append
func NestedSliceAppend() int {
	var s [][]int
	s = append(s, []int{1, 2})
	s = append(s, []int{3, 4})
	s[0] = append(s[0], 5)

	sum := 0
	for _, inner := range s {
		for _, v := range inner {
			sum += v
		}
	}
	return sum // 1+2+5+3+4 = 15
}

// SliceResliceMultiple tests multiple reslicing
func SliceResliceMultiple() int {
	s := make([]int, 10)
	for i := range s {
		s[i] = i
	}

	s1 := s[2:5]  // [2, 3, 4]
	s2 := s1[1:3] // [3, 4]
	s3 := s2[0:2] // [3, 4]

	return s3[0] + s3[1] // 3 + 4 = 7
}

// SliceWithCapacityGrowth tests slice capacity growth
func SliceWithCapacityGrowth() int {
	s := make([]int, 0, 2)
	initialCap := cap(s)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	return len(s) + initialCap // 10 + 2 = 12
}

// IntegerOverflow tests integer overflow behavior
func IntegerOverflow() int {
	max := int(^uint(0) >> 1)
	result := max + 1 // Overflow
	return result
}

// ShiftOperations tests shift operations
func ShiftOperations() int {
	x := 1
	y := x << 31 // Shift left 31 bits
	z := y >> 31 // Shift back
	return int(z)
}

// FloatSpecialValues tests float special values
func FloatSpecialValues() int {
	// Use math.Inf to get infinity values
	// inf := math.Inf(1)  // Can't use math package in tests
	// negInf := math.Inf(-1)

	// Test with very large numbers instead
	large := 1e308
	veryLarge := large * 10 // This becomes +Inf at runtime

	if veryLarge > 0 {
		return 1
	}
	return 0
}

// DeeplyNestedIfElse tests deeply nested if-else
func DeeplyNestedIfElse() int {
	x, y, z := 1, 2, 3

	if x > 0 {
		if y > 0 {
			if z > 0 {
				return 1
			} else {
				return 2
			}
		} else {
			if z > 0 {
				return 3
			} else {
				return 4
			}
		}
	} else {
		if y > 0 {
			if z > 0 {
				return 5
			} else {
				return 6
			}
		} else {
			return 7
		}
	}
}

// SwitchWithComplexConditions tests switch with complex conditions
func SwitchWithComplexConditions() int {
	x := 5
	switch {
	case x < 0:
		return 0
	case x >= 0 && x < 3:
		return 1
	case x >= 3 && x < 7:
		return 2
	default:
		return 3
	}
}

// ForWithComplexCondition tests for with complex condition
func ForWithComplexCondition() int {
	sum := 0
	for i, j := 0, 10; i < j; i, j = i+1, j-1 {
		sum += i + j
	}
	return sum
}

// ChannelSelectWithMultipleCases tests select with multiple cases
func ChannelSelectWithMultipleCases() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2

	result := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			result += v
		case v := <-ch2:
			result += v
		}
	}
	return result // 1 + 2 = 3
}

// NilChannelInSelect tests nil channel in select
func NilChannelInSelect() int {
	var nilCh chan int
	ch := make(chan int, 1)
	ch <- 42

	select {
	case v := <-ch:
		return v
	case <-nilCh: // Never ready
		return -1
	default:
		return 0
	}
}

// MultipleAssignmentChain tests chained multiple assignments
func MultipleAssignmentChain() int {
	a, b, c := 1, 2, 3
	d, e, f := a+1, b+1, c+1
	g, h, i := d+e, e+f, f+d

	return g + h + i
}

// AssignmentWithFunctionCall tests assignment with function calls
func AssignmentWithFunctionCall() int {
	getValue := func() (int, int, int) {
		return 10, 20, 30
	}

	a, b, c := getValue()
	d, e, f := getValue()

	return a + b + c + d + e + f // 60 + 60 = 120
}

// SwapMultipleValues tests swapping multiple values
func SwapMultipleValues() int {
	a, b, c := 1, 2, 3
	a, b, c = c, a, b

	return a*100 + b*10 + c // 300 + 10 + 2 = 312
}

// ============================================================================
// Missing Round 7 functions
// ============================================================================

// MapWithComplexKeys tests map with complex struct keys
func MapWithComplexKeys() int {
	type Key struct {
		X int
		Y int
	}

	m := make(map[Key]string)
	m[Key{X: 1, Y: 2}] = "first"
	m[Key{X: 3, Y: 4}] = "second"

	return len(m) // Should be 2
}

// ComplexArithmetic tests complex number arithmetic
func ComplexArithmetic() int {
	var c1 complex128 = 3 + 4i
	var c2 complex128 = 1 + 2i

	sum := c1 + c2
	prod := c1 * c2

	// real(3+4i + 1+2i) = 4, imag = 6
	// real((3+4i)*(1+2i)) = 3-8 = -5, imag = 6+4 = 10
	return int(real(sum)) + int(imag(sum)) + int(real(prod)) + int(imag(prod)) // 4 + 6 + -5 + 10 = 15
}

// NilInterfaceMethodCall tests calling method on nil interface
func NilInterfaceMethodCall() int {
	var i interface{} = nil

	// Check if nil
	if i == nil {
		return 1
	}
	return 0
}

// MethodSetTest tests method set on type
func MethodSetTest() int {
	type Counter struct {
		value int
	}

	// Value receiver method
	getValue := func(c Counter) int {
		return c.value
	}

	c := Counter{value: 42}
	return getValue(c) // Should be 42
}

// MethodOnNilPointer tests calling method on nil pointer
func MethodOnNilPointer() int {
	type Container struct {
		Value int
	}

	var ptr *Container

	// Safe nil check
	if ptr == nil {
		return 1
	}
	return 0
}

// ============================================================================
// ROUND 8: COMPLEX THIRD-PARTY LIBRARY INTERACTIONS
// ============================================================================

// ============================================================================
// ENCODING/BASE64 - Complex encoding/decoding scenarios
// ============================================================================

// Base64ComplexRoundTrip tests complex base64 round-trip with different encodings
func Base64ComplexRoundTrip() int {
	data := []byte("Hello, 世界! This is a complex test with unicode characters.")

	// Standard encoding
	stdEnc := base64.StdEncoding.EncodeToString(data)
	decoded1, _ := base64.StdEncoding.DecodeString(stdEnc)

	// URL encoding
	urlEnc := base64.URLEncoding.EncodeToString(data)
	decoded2, _ := base64.URLEncoding.DecodeString(urlEnc)

	// Raw encoding (no padding)
	rawEnc := base64.RawStdEncoding.EncodeToString(data)
	decoded3, _ := base64.RawStdEncoding.DecodeString(rawEnc)

	// Calculate result based on successful decodes
	result := 0
	if bytes.Equal(decoded1, data) {
		result += 1
	}
	if bytes.Equal(decoded2, data) {
		result += 2
	}
	if bytes.Equal(decoded3, data) {
		result += 4
	}

	return result // Should be 7 if all succeed
}

// Base64Streaming tests streaming base64 encoding/decoding
func Base64Streaming() int {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte("chunk1"))
	encoder.Write([]byte("chunk2"))
	encoder.Write([]byte("chunk3"))
	encoder.Close()

	encoded := buf.String()

	// Decode
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	decoded, _ := io.ReadAll(decoder)

	return len(decoded) // Should be 18 ("chunk1chunk2chunk3")
}

// Base64WithCorruptInput tests handling of corrupt base64 input
func Base64WithCorruptInput() int {
	valid := base64.StdEncoding.EncodeToString([]byte("valid"))
	corrupt := "SGVsbG8gV29ybGQ!=" // Extra padding

	result := 0

	_, err1 := base64.StdEncoding.DecodeString(valid)
	if err1 == nil {
		result += 1
	}

	_, err2 := base64.StdEncoding.DecodeString(corrupt)
	if err2 != nil {
		result += 2
	}

	return result // Should be 3
}

// ============================================================================
// ENCODING/BINARY - Byte order operations
// ============================================================================

// BinaryReadWrite tests binary read/write with different byte orders
func BinaryReadWrite() int {
	var buf bytes.Buffer

	// Write in big-endian
	binary.Write(&buf, binary.BigEndian, uint16(0x1234))
	binary.Write(&buf, binary.BigEndian, uint32(0x12345678))
	binary.Write(&buf, binary.BigEndian, int64(-123456789))

	// Read back
	var u16 uint16
	var u32 uint32
	var i64 int64

	binary.Read(&buf, binary.BigEndian, &u16)
	binary.Read(&buf, binary.BigEndian, &u32)
	binary.Read(&buf, binary.BigEndian, &i64)

	result := 0
	if u16 == 0x1234 {
		result += 1
	}
	if u32 == 0x12345678 {
		result += 2
	}
	if i64 == -123456789 {
		result += 4
	}

	return result // Should be 7
}

// BinarySliceConversion tests binary conversion of slices
func BinarySliceConversion() int {
	data := []uint16{0x0102, 0x0304, 0x0506}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, data)

	var decoded []uint16
	decoded = make([]uint16, 3)
	binary.Read(&buf, binary.LittleEndian, decoded)

	if decoded[0] == 0x0102 && decoded[1] == 0x0304 && decoded[2] == 0x0506 {
		return 1
	}
	return 0
}

// BinarySize tests binary.Size function
func BinarySize() int {
	size := binary.Size(uint16(0)) + binary.Size(uint32(0)) + binary.Size(int64(0))

	// 2 + 4 + 8 = 14
	return size
}

// ============================================================================
// ENCODING/JSON - Advanced JSON operations
// ============================================================================

// JSONRawMessage tests json.RawMessage for delayed decoding
func JSONRawMessage() int {
	type Envelope struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	data := `{"type":"greeting","payload":{"message":"Hello"}}`

	var env Envelope
	json.Unmarshal([]byte(data), &env)

	if env.Type == "greeting" {
		// Now decode the payload
		var payload map[string]string
		json.Unmarshal(env.Payload, &payload)
		if payload["message"] == "Hello" {
			return 1
		}
	}
	return 0
}

// JSONCustomMarshaler tests custom JSON marshalers
func JSONCustomMarshaler() int {
	type Point struct {
		X int
		Y int
	}

	// Custom marshaler
	marshalPoint := func(p Point) ([]byte, error) {
		return json.Marshal([]int{p.X, p.Y})
	}

	p := Point{X: 10, Y: 20}
	data, _ := marshalPoint(p)

	var coords []int
	json.Unmarshal(data, &coords)

	if coords[0] == 10 && coords[1] == 20 {
		return 1
	}
	return 0
}

// JSONOmitEmpty tests json omitempty tag behavior
func JSONOmitEmpty() int {
	type Data struct {
		Required string `json:"required"`
		Optional string `json:"optional,omitempty"`
		Number   int    `json:"number,omitempty"`
	}

	d := Data{
		Required: "must have",
		Optional: "",
		Number:   0,
	}

	data, _ := json.Marshal(d)

	// Should not contain "optional" or "number" fields
	str := string(data)
	result := 0
	if strings.Contains(str, `"required":"must have"`) {
		result += 1
	}
	if !strings.Contains(str, `"optional"`) {
		result += 2
	}
	if !strings.Contains(str, `"number"`) {
		result += 4
	}

	return result // Should be 7
}

// JSONStringIgnoreCase tests JSON field name matching
func JSONStringIgnoreCase() int {
	type Data struct {
		FieldName string `json:"field_name"`
	}

	// JSON with different casing
	data := `{"field_name":"test"}`

	var d Data
	json.Unmarshal([]byte(data), &d)

	if d.FieldName == "test" {
		return 1
	}
	return 0
}

// JSONComplexMap tests JSON with complex map keys
func JSONComplexMap() int {
	// Map with struct keys (requires encoding as array of pairs)
	type Item struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}

	items := []Item{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
	}

	data, _ := json.Marshal(items)

	var decoded []Item
	json.Unmarshal(data, &decoded)

	if len(decoded) == 2 && decoded[0].Key == "a" && decoded[1].Value == 2 {
		return 1
	}
	return 0
}

// ============================================================================
// IO - Complex IO patterns
// ============================================================================

// IOPipe tests io.Pipe for in-process streaming
func IOPipe() int {
	r, w := io.Pipe()

	// Writer goroutine (simulated with closure)
	go func() {
		w.Write([]byte("Hello, "))
		w.Write([]byte("Pipe!"))
		w.Close()
	}()

	// Reader
	data, _ := io.ReadAll(r)

	if string(data) == "Hello, Pipe!" {
		return 1
	}
	return 0
}

// IOLimitReader tests io.LimitReader
func IOLimitReader() int {
	src := strings.NewReader("This is a long string")
	limited := io.LimitReader(src, 10)

	data, _ := io.ReadAll(limited)

	if len(data) == 10 && string(data) == "This is a " {
		return 1
	}
	return 0
}

// IOSectionReader tests io.SectionReader
func IOSectionReader() int {
	src := strings.NewReader("0123456789ABCDEF")
	section := io.NewSectionReader(src, 5, 5)

	data, _ := io.ReadAll(section)

	if string(data) == "56789" {
		return 1
	}
	return 0
}

// IOTeeReader tests io.TeeReader
func IOTeeReader() int {
	src := strings.NewReader("Hello World")
	var buf bytes.Buffer

	tee := io.TeeReader(src, &buf)

	// Read from tee
	readData, _ := io.ReadAll(tee)

	// Both should have the same content
	if string(readData) == "Hello World" && buf.String() == "Hello World" {
		return 1
	}
	return 0
}

// IOMultiReader tests io.MultiReader
func IOMultiReader() int {
	r1 := strings.NewReader("First ")
	r2 := strings.NewReader("Second ")
	r3 := strings.NewReader("Third")

	multi := io.MultiReader(r1, r2, r3)
	data, _ := io.ReadAll(multi)

	if string(data) == "First Second Third" {
		return 1
	}
	return 0
}

// IOMultiWriter tests io.MultiWriter
func IOMultiWriter() int {
	var buf1, buf2, buf3 bytes.Buffer

	mw := io.MultiWriter(&buf1, &buf2, &buf3)
	mw.Write([]byte("broadcast"))

	if buf1.String() == "broadcast" && buf2.String() == "broadcast" && buf3.String() == "broadcast" {
		return 1
	}
	return 0
}

// IOCopy tests io.Copy with various sources and destinations
func IOCopy() int {
	src := strings.NewReader("Data to copy")
	var dst bytes.Buffer

	n, _ := io.Copy(&dst, src)

	if n == 12 && dst.String() == "Data to copy" {
		return 1
	}
	return 0
}

// IOCopyBuffer tests io.CopyBuffer with custom buffer
func IOCopyBuffer() int {
	src := strings.NewReader("Custom buffer")
	var dst bytes.Buffer
	buf := make([]byte, 5) // Small buffer to force multiple reads

	n, _ := io.CopyBuffer(&dst, src, buf)

	if n == 13 && dst.String() == "Custom buffer" {
		return 1
	}
	return 0
}

// ============================================================================
// ARCHIVE/TAR - Tar header operations
// ============================================================================

// TarHeader tests tar.Header creation and manipulation
func TarHeader() int {
	hdr := &tar.Header{
		Name:     "test.txt",
		Mode:     0644,
		Size:     100,
		ModTime:  time.Now(),
		Typeflag: byte(tar.TypeReg),
	}

	result := 0
	if hdr.Name == "test.txt" {
		result += 1
	}
	if hdr.Mode == 0644 {
		result += 2
	}
	if hdr.Size == 100 {
		result += 4
	}
	if hdr.Typeflag == byte(tar.TypeReg) {
		result += 8
	}

	return result // Should be 15
}

// TarWriteRead tests tar archive write and read
func TarWriteRead() int {
	var buf bytes.Buffer

	// Create tar archive
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: "hello.txt",
		Mode: 0644,
		Size: int64(len("Hello, Tar!")),
	}

	tw.WriteHeader(hdr)
	tw.Write([]byte("Hello, Tar!"))
	tw.Close()

	// Read tar archive
	tr := tar.NewReader(&buf)

	hdr2, _ := tr.Next()
	content, _ := io.ReadAll(tr)

	if hdr2.Name == "hello.txt" && string(content) == "Hello, Tar!" {
		return 1
	}
	return 0
}

// TarMultipleFiles tests tar archive with multiple files
func TarMultipleFiles() int {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "Content 1"},
		{"file2.txt", "Content 2"},
		{"file3.txt", "Content 3"},
	}

	for _, f := range files {
		hdr := &tar.Header{
			Name: f.name,
			Mode: 0644,
			Size: int64(len(f.content)),
		}
		tw.WriteHeader(hdr)
		tw.Write([]byte(f.content))
	}
	tw.Close()

	// Read back
	tr := tar.NewReader(&buf)
	count := 0
	for {
		_, err := tr.Next()
		if err != nil {
			break
		}
		io.ReadAll(tr)
		count++
	}

	return count // Should be 3
}

// ============================================================================
// ARCHIVE/ZIP - Zip operations
// ============================================================================

// ZipHeader tests zip.FileHeader creation
func ZipHeader() int {
	hdr := &zip.FileHeader{
		Name:               "test.txt",
		UncompressedSize:   100,
		UncompressedSize64: 100,
		Method:             zip.Deflate,
	}

	result := 0
	if hdr.Name == "test.txt" {
		result += 1
	}
	if hdr.Method == zip.Deflate {
		result += 2
	}
	if hdr.UncompressedSize == 100 {
		result += 4
	}

	return result // Should be 7
}

// ZipWriteRead tests zip archive write and read
func ZipWriteRead() int {
	var buf bytes.Buffer

	// Create zip archive
	zw := zip.NewWriter(&buf)

	w, _ := zw.Create("hello.txt")
	w.Write([]byte("Hello, Zip!"))
	zw.Close()

	// Read zip archive
	zr, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))

	if len(zr.File) == 1 && zr.File[0].Name == "hello.txt" {
		rc, _ := zr.File[0].Open()
		content, _ := io.ReadAll(rc)
		rc.Close()
		if string(content) == "Hello, Zip!" {
			return 1
		}
	}
	return 0
}

// ZipMultipleFiles tests zip archive with multiple files
func ZipMultipleFiles() int {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := []struct {
		name    string
		content string
	}{
		{"dir/file1.txt", "Content 1"},
		{"dir/file2.txt", "Content 2"},
		{"file3.txt", "Content 3"},
	}

	for _, f := range files {
		w, _ := zw.Create(f.name)
		w.Write([]byte(f.content))
	}
	zw.Close()

	// Read back
	zr, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))

	return len(zr.File) // Should be 3
}

// ============================================================================
// TEXT/TABWRITER - Aligned output
// ============================================================================

// TabwriterBasic tests basic tabwriter functionality
func TabwriterBasic() int {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "Name\tAge\tCity")
	fmt.Fprintln(tw, "Alice\t30\tNew York")
	fmt.Fprintln(tw, "Bob\t25\tLos Angeles")
	tw.Flush()

	lines := strings.Count(buf.String(), "\n")
	return lines // Should be 3
}

// TabwriterAligned tests tabwriter alignment
func TabwriterAligned() int {
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 8, 8, 2, ' ', 0)

	fmt.Fprintln(tw, "A\tB")
	fmt.Fprintln(tw, "AAA\tBBB")
	tw.Flush()

	output := buf.String()
	// Verify alignment by checking length
	if len(output) > 0 {
		return 1
	}
	return 0
}

// ============================================================================
// TEXT/SCANNER - Lexical scanning
// ============================================================================

// ScannerBasic tests basic scanner functionality
func ScannerBasic() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("hello world 123"))

	count := 0
	for {
		tok := s.Scan()
		if int(tok) == scanner.EOF {
			break
		}
		count++
	}

	return count // Should be 3 (hello, world, 123)
}

// ScannerNumbers tests scanning numbers
func ScannerNumbers() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("42 3.14 -17"))

	count := 0
	for {
		tok := s.Scan()
		if int(tok) == scanner.EOF {
			break
		}
		if int(tok) == scanner.Int || int(tok) == scanner.Float {
			count++
		}
	}

	return count // Should be 3
}

// ScannerCustomMode tests scanner with custom mode
func ScannerCustomMode() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("var x = 42"))
	s.Mode = uint(scanner.ScanIdents | scanner.ScanInts)

	tokens := []string{}
	for {
		tok := s.Scan()
		if int(tok) == scanner.EOF {
			break
		}
		tokens = append(tokens, s.TokenText())
	}

	if len(tokens) == 4 && tokens[0] == "var" && tokens[3] == "42" {
		return 1
	}
	return 0
}

// ============================================================================
// UNICODE/UTF8 - UTF-8 operations
// ============================================================================

// UTF8RuneCount tests utf8.RuneCountInString
func UTF8RuneCount() int {
	s := "Hello, 世界"
	count := utf8.RuneCountInString(s)
	return count // Should be 9 (7 ASCII + 2 Chinese)
}

// UTF8EncodeDecode tests utf8 encode/decode
func UTF8EncodeDecode() int {
	r := '世'
	buf := make([]byte, 3)
	n := utf8.EncodeRune(buf, r)

	if n == 3 {
		decodedRune, _ := utf8.DecodeRune(buf)
		if decodedRune == r {
			return 1
		}
	}
	return 0
}

// UTF8Valid tests utf8.Valid
func UTF8Valid() int {
	valid := []byte("Hello, 世界")
	invalid := []byte{0xff, 0xfe, 0xfd}

	result := 0
	if utf8.Valid(valid) {
		result += 1
	}
	if !utf8.Valid(invalid) {
		result += 2
	}

	return result // Should be 3
}

// ============================================================================
// COMBINED COMPLEX SCENARIOS
// ============================================================================

// CombinedEncodingPipeline tests multiple encoding layers
func CombinedEncodingPipeline() int {
	data := "Hello, Multi-encoding World!"

	// Step 1: JSON encode
	jsonData, _ := json.Marshal(map[string]string{"data": data})

	// Step 2: Base64 encode
	base64Data := base64.StdEncoding.EncodeToString(jsonData)

	// Step 3: Base64 decode
	jsonDecoded, _ := base64.StdEncoding.DecodeString(base64Data)

	// Step 4: JSON decode
	var result map[string]string
	json.Unmarshal(jsonDecoded, &result)

	if result["data"] == data {
		return 1
	}
	return 0
}

// CombinedArchiveWithCompression tests tar with gzip-like compression
func CombinedArchiveWithCompression() int {
	var tarBuf bytes.Buffer

	// Create tar
	tw := tar.NewWriter(&tarBuf)
	hdr := &tar.Header{
		Name: "test.txt",
		Mode: 0644,
		Size: int64(len("test content")),
	}
	tw.WriteHeader(hdr)
	tw.Write([]byte("test content"))
	tw.Close()

	// Read tar
	tr := tar.NewReader(&tarBuf)
	hdr2, _ := tr.Next()
	content, _ := io.ReadAll(tr)

	if hdr2.Name == "test.txt" && string(content) == "test content" {
		return 1
	}
	return 0
}

// CombinedIOAndEncoding tests io operations with encoding
func CombinedIOAndEncoding() int {
	var buf bytes.Buffer

	// Write binary data
	binary.Write(&buf, binary.BigEndian, uint32(12345))
	binary.Write(&buf, binary.BigEndian, uint32(67890))

	// Base64 encode the binary
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Decode
	decoded, _ := base64.StdEncoding.DecodeString(encoded)

	// Read back binary
	reader := bytes.NewReader(decoded)
	var v1, v2 uint32
	binary.Read(reader, binary.BigEndian, &v1)
	binary.Read(reader, binary.BigEndian, &v2)

	if v1 == 12345 && v2 == 67890 {
		return 1
	}
	return 0
}

// CombinedNestedClosuresWithThirdParty tests closures with third-party libraries
func CombinedNestedClosuresWithThirdParty() int {
	// Create encoder function
	encodeJSON := func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}

	// Create nested encoder that adds metadata
	encodeWithMeta := func(v interface{}) ([]byte, error) {
		wrapped := map[string]interface{}{
			"timestamp": "2026-04-03",
			"data":      v,
		}
		return encodeJSON(wrapped)
	}

	// Use it
	data := map[string]int{"x": 1, "y": 2}
	result, _ := encodeWithMeta(data)

	// Verify
	var decoded map[string]interface{}
	json.Unmarshal(result, &decoded)

	if decoded["timestamp"] == "2026-04-03" {
		return 1
	}
	return 0
}

// CombinedDeferWithIO tests defer with io operations
// CombinedPanicRecoverWithJSON tests panic/recover in JSON operations
func CombinedPanicRecoverWithJSON() int {
	result := func() (ret int) {
		defer func() {
			if r := recover(); r != nil {
				ret = -1
			}
		}()

		// This should not panic
		data, _ := json.Marshal(map[string]int{"key": 42})
		var decoded map[string]int
		json.Unmarshal(data, &decoded)

		return decoded["key"]
	}()

	return result // Should be 42
}

// CombinedChannelWithEncoding tests channels with encoding operations
func CombinedChannelWithEncoding() int {
	ch := make(chan []byte, 3)

	// Producer
	go func() {
		data := []string{"a", "b", "c"}
		for _, s := range data {
			encoded, _ := json.Marshal(s)
			ch <- encoded
		}
		close(ch)
	}()

	// Consumer
	count := 0
	for encoded := range ch {
		var decoded string
		json.Unmarshal(encoded, &decoded)
		if decoded != "" {
			count++
		}
	}

	return count // Should be 3
}

// CombinedStructMethodsAndJSON tests struct methods with JSON marshaling
func CombinedStructMethodsAndJSON() int {
	type Counter struct {
		Value int `json:"value"`
	}

	c := &Counter{Value: 10}

	// Modify via method
	increment := func(c *Counter) {
		c.Value++
	}

	increment(c)
	increment(c)
	increment(c)

	// Serialize
	data, _ := json.Marshal(c)

	var decoded Counter
	json.Unmarshal(data, &decoded)

	return decoded.Value // Should be 13
}

// CombinedInterfaceAndEncoding tests interface implementation with encoding
func CombinedInterfaceAndEncoding() int {
	var buf bytes.Buffer
	var w io.Writer = &buf

	// Write JSON to interface
	data, _ := json.Marshal(map[string]string{"key": "value"})
	w.Write(data)

	// Read back
	var decoded map[string]string
	json.Unmarshal(buf.Bytes(), &decoded)

	if decoded["key"] == "value" {
		return 1
	}
	return 0
}

// ============================================================================
// ROUND 13: 100 NEW DIVERGENCE HUNT TESTS
// ============================================================================

// SliceSharedUnderlying tests that slices share underlying array
func SliceSharedUnderlying() int {
	a := [5]int{1, 2, 3, 4, 5}
	s1 := a[1:3] // [2, 3]
	s2 := a[2:4] // [3, 4]
	s1[1] = 99   // modifies a[2]
	return s2[0] // should be 99 because s2 shares underlying array
}

// SliceAppendDoesNotShare tests that append may create new array
func SliceAppendDoesNotShare() int {
	s := []int{1, 2, 3}
	s2 := append(s, 4)
	s[0] = 99
	return s2[0] // should be 1, not 99, because append created new array
}

// SliceAppendSharesBeforeCap tests append within capacity shares array
func SliceAppendSharesBeforeCap() int {
	s := make([]int, 2, 5)
	s[0], s[1] = 1, 2
	s2 := append(s, 3) // within cap, shares underlying array
	s[0] = 99
	return s2[0] // should be 99 because s2 shares underlying array with s
}

// SliceCopySemantics tests copy function behavior
func SliceCopySemantics() int {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n + dst[0] + dst[2] // 3 + 1 + 3 = 7
}

// SliceCopyFromSubslice tests copy from subslice
func SliceCopyFromSubslice() int {
	src := []int{10, 20, 30, 40, 50}
	dst := make([]int, 2)
	n := copy(dst, src[2:])
	return n + dst[1] // 2 + 40 = 42
}

// DeferNamedReturnMutate tests defer mutating named return
func DeferNamedReturnMutate() (result int) {
	result = 10
	defer func() { result *= 2 }()
	return result // defer changes result to 20
}

// DeferNamedReturnOverride tests defer overriding named return
func DeferNamedReturnOverride() (result int) {
	defer func() { result = 42 }()
	return 10 // defer overrides to 42
}

// InterfaceBoxNilSlice tests that nil slice boxed in interface != nil
func InterfaceBoxNilSlice() bool {
	var s []int
	var i interface{} = s
	return i == nil // false: non-nil interface with nil value
}

// InterfaceBoxNilMap tests that nil map boxed in interface != nil
func InterfaceBoxNilMap() bool {
	var m map[string]int
	var i interface{} = m
	return i == nil // false
}

// InterfaceBoxNilPointer tests that nil pointer boxed in interface != nil
func InterfaceBoxNilPointer() bool {
	var p *int
	var i interface{} = p
	return i == nil // false
}

// InterfaceBoxNilFunc tests that nil func boxed in interface != nil
func InterfaceBoxNilFunc() bool {
	var f func()
	var i interface{} = f
	return i == nil // false
}

// InterfaceBoxNilChannel tests that nil channel boxed in interface != nil
func InterfaceBoxNilChannel() bool {
	var ch chan int
	var i interface{} = ch
	return i == nil // false
}

// TypeAssertionCommaOKSuccess tests comma-ok type assertion success
func TypeAssertionCommaOKSuccess() int {
	var i interface{} = 42
	v, ok := i.(int)
	if ok {
		return v
	}
	return -1
}

// TypeAssertionCommaOKFailure tests comma-ok type assertion failure
func TypeAssertionCommaOKFailure() int {
	var i interface{} = "hello"
	v, ok := i.(int)
	if ok {
		return v
	}
	return -1
}

// MapAccessZeroValue tests accessing non-existent key returns zero value
func MapAccessZeroValue() int {
	m := map[string]int{"a": 1}
	return m["b"] // returns 0
}

// MapAccessBoolZeroValue tests accessing non-existent bool key
func MapAccessBoolZeroValue() bool {
	m := map[string]bool{"a": true}
	return m["b"] // returns false
}

// MapAccessStructZeroValue tests accessing non-existent struct key
func MapAccessStructZeroValue() int {
	type Point struct{ X, Y int }
	m := map[string]Point{"a": {1, 2}}
	p := m["b"] // returns Point{0, 0}
	return p.X + p.Y
}

// MapAccessPointerZeroValue tests accessing non-existent pointer key
func MapAccessPointerZeroValue() bool {
	m := map[string]*int{}
	p := m["a"]
	return p == nil // true
}

// SliceThreeIndexLimit tests three-index slice capacity limit
func SliceThreeIndexLimit() int {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[1:3:3]           // len=2, cap=2
	return len(s2) + cap(s2) // 2 + 2 = 4
}

// SliceThreeIndexAppend tests append after three-index slice
func SliceThreeIndexAppend() int {
	s := []int{1, 2, 3, 4, 5}
	s2 := s[1:3:3]
	s2 = append(s2, 10) // does NOT overwrite s[3] because cap was limited
	return s[3]         // should still be 4
}

// RangeValueCopy tests that range copies values
func RangeValueCopy() int {
	type Item struct{ Val int }
	items := []Item{{1}, {2}, {3}}
	sum := 0
	for _, item := range items {
		item.Val *= 2 // does NOT modify original
	}
	for _, item := range items {
		sum += item.Val
	}
	return sum // 1 + 2 + 3 = 6
}

// RangePointerMutation tests range with pointer mutation
func RangePointerMutation() int {
	type Item struct{ Val int }
	items := []*Item{{1}, {2}, {3}}
	for _, item := range items {
		item.Val *= 2 // DOES modify original
	}
	sum := 0
	for _, item := range items {
		sum += item.Val
	}
	return sum // 2 + 4 + 6 = 12
}

// ClosureCaptureTiming tests closure captures variable at call time
func ClosureCaptureTiming() int {
	x := 1
	f := func() int { return x }
	x = 2
	return f() // returns 2, not 1
}

// ClosureCaptureByReference tests closure captures by reference
func ClosureCaptureByReference() int {
	x := 1
	f := func() { x = 10 }
	f()
	return x // returns 10
}

// ShortVarRedeclare tests short var redeclaration in same scope
func ShortVarRedeclare() int {
	x := 1
	x, y := 2, 3 // x is redeclared, y is new
	return x + y // 2 + 3 = 5
}

// ShortVarRedeclareInInnerScope tests short var in inner scope doesn't affect outer
func ShortVarRedeclareInInnerScope() int {
	x := 1
	{
		x := 2 // new x in inner scope
		_ = x
	}
	return x // still 1
}

// MapDeleteAndAccess tests delete then access
func MapDeleteAndAccess() int {
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	v, ok := m["a"]
	if ok {
		return v
	}
	return -1 // deleted, returns -1
}

// MapLenAfterDelete tests map length after delete
func MapLenAfterDelete() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	delete(m, "b")
	return len(m) // 2
}

// CounterR13 for StructMethodPointerReceiver test
type CounterR13 struct{ N int }

func (c *CounterR13) IncrR13() { c.N++ }

// StructMethodPointerReceiver tests pointer receiver method
func StructMethodPointerReceiver() int {
	c := CounterR13{N: 0}
	c.IncrR13()
	c.IncrR13()
	return c.N
}

// ValR13 for StructMethodValueReceiver test
type ValR13 struct{ N int }

func (v ValR13) AddR13(n int) { v.N += n } // doesn't affect original

// StructMethodValueReceiver tests value receiver doesn't mutate
func StructMethodValueReceiver() int {
	v := ValR13{N: 5}
	v.AddR13(10)
	return v.N // still 5, value receiver
}

// ChannelCommaOK tests channel receive with comma-ok
func ChannelCommaOK() int {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	sum := 0
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		sum += v
	}
	return sum // 1 + 2 = 3
}

// ChannelClosedZeroValue tests receiving from closed channel returns zero
func ChannelClosedZeroValue() int {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)
	<-ch      // 42
	v := <-ch // zero value after close
	return v  // 0
}

// SwitchCaseBreak tests switch case with break doesn't break outer loop
func SwitchCaseBreak() int {
	sum := 0
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			break // only breaks switch, not for loop
		}
		sum += i
	}
	return sum // 0 + 1 + 2 = 3 (break only exits switch)
}

// SwitchFallthroughLogic tests fallthrough continues to next case
func SwitchFallthroughLogic() int {
	n := 0
	switch {
	case true:
		n = 1
		fallthrough
	case false:
		n += 10
	}
	return n // 1 + 10 = 11 (fallthrough executes next case unconditionally)
}

// StringIndexByte tests string indexing returns byte
func StringIndexByte() int {
	s := "Hello"
	return int(s[0]) // 72 (ASCII 'H')
}

// StringSliceShared tests string slice shares underlying data
func StringSliceShared() string {
	s := "Hello, World"
	s2 := s[7:] // "World"
	return s2
}

// ByteSliceModification tests byte slice modification
func ByteSliceModification() string {
	b := []byte("hello")
	b[0] = 'H'
	return string(b) // "Hello"
}

// MultipleAssignmentOrder tests multiple assignment evaluates right side first
func MultipleAssignmentOrder() int {
	a := []int{1, 2, 3}
	i := 0
	i, a[i] = 1, 100 // i becomes 1, a[0] becomes 100 (i was 0 when evaluating left side)
	return a[0]      // 100
}

// MultipleAssignmentSwap tests swap via multiple assignment
func MultipleAssignmentSwap() int {
	a, b := 1, 2
	a, b = b, a
	return a*10 + b // 21
}

// NilSliceCompareEmpty tests nil slice vs empty slice comparison
func NilSliceCompareEmpty() bool {
	var s1 []int
	s2 := []int{}
	return s1 == nil && s2 != nil // true
}

// StructComparable tests struct with comparable fields
func StructComparable() bool {
	type Point struct{ X, Y int }
	p1 := Point{1, 2}
	p2 := Point{1, 2}
	return p1 == p2 // true
}

// StructNotComparable tests struct with slice is not comparable (via workaround)
func StructNotComparable() int {
	type S struct {
		X int
	}
	type W struct {
		X int
		T []int // not comparable
	}
	w := W{X: 1, T: []int{1}}
	return w.X // just access comparable field
}

// PointerComparisonR13 tests comparing pointers
func PointerComparisonR13() bool {
	x := 42
	p1 := &x
	p2 := &x
	return p1 == p2 // true: same address
}

// DifferentPointerComparisonR13 tests different pointers
func DifferentPointerComparisonR13() bool {
	x := 42
	y := 42
	p1 := &x
	p2 := &y
	return p1 == p2 // false: different addresses
}

// Rect3 for InterfaceMethodDispatch test
type Rect3 struct{ W, H int }

func (r *Rect3) Area3() int { return r.W * r.H }

// InterfaceMethodDispatch tests interface method dispatch
func InterfaceMethodDispatch() int {
	r := &Rect3{W: 3, H: 4}
	var s interface{ Area3() int } = r
	return s.Area3() // 12
}

// Rect2 for InterfaceMethodDispatch2 test
type Rect2Alt struct{ W, H int }

func (r *Rect2Alt) Area2() int { return r.W * r.H }

// InterfaceMethodDispatch2 tests interface method dispatch
func InterfaceMethodDispatch2() int {
	r := &Rect2Alt{W: 3, H: 4}
	var s interface{ Area2() int } = r
	return s.Area2() // 12
}

// SliceMakeInitialized tests make creates zero-initialized slice
func SliceMakeInitialized() int {
	s := make([]int, 5)
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum // 0
}

// MapMakeSizeHint tests map with size hint
func MapMakeSizeHint() int {
	m := make(map[string]int, 100)
	m["a"] = 1
	return len(m) // 1
}

// NilMapWritePanic tests nil map write (catchable pattern)
func NilMapWritePanic() int {
	var m map[string]int
	if m == nil {
		m = make(map[string]int)
	}
	m["key"] = 42
	return m["key"]
}

// DeferStackOrder tests defer executes in LIFO order
func DeferStackOrder() string {
	result := ""
	defer func() { result += "A" }()
	defer func() { result += "B" }()
	defer func() { result += "C" }()
	return result // "" because defers run after return
}

// DeferStackOrderWithNamed tests defer with named return
func DeferStackOrderWithNamed() (result string) {
	defer func() { result += "A" }()
	defer func() { result += "B" }()
	result = "X"
	return // defers run: result = "X" -> "XB" -> "XBA"
}

// SliceAppendGrowth tests slice growth pattern
func SliceAppendGrowth() int {
	s := make([]int, 0, 1)
	s = append(s, 1)
	oldCap := cap(s)
	s = append(s, 2)
	newCap := cap(s)
	return oldCap + newCap/2 // depends on growth factor
}

// ChannelSelectDefault tests select with default (non-blocking)
func ChannelSelectDefault() int {
	ch := make(chan int)
	select {
	case v := <-ch:
		return v
	default:
		return -1 // non-blocking, returns -1
	}
}

// ChannelSelectNilChannel tests select with nil channel
func ChannelSelectNilChannel() int {
	var nilCh chan int
	readyCh := make(chan int, 1)
	readyCh <- 42
	select {
	case v := <-nilCh:
		return v // never reached
	case v := <-readyCh:
		return v // 42
	}
}

// GotoSkipVariable tests goto can skip over declarations if not used
func GotoSkipVariable() int {
	x := 0
	goto Label
Label:
	return x // 0
}

// GotoBackwardLoop tests goto backward for loop
func GotoBackwardLoop() int {
	sum := 0
	i := 0
Loop:
	sum += i
	i++
	if i < 5 {
		goto Loop
	}
	return sum // 0+1+2+3+4 = 10
}

// EsmBase for EmbeddedStructMethodResolution
type EsmBase struct{ X int }

func (b *EsmBase) GetEsmX() int { return b.X }

// EmbeddedStructMethodResolution tests method resolution with embedding
func EmbeddedStructMethodResolution() int {
	type Derived struct {
		EsmBase
		Y int
	}
	d := Derived{EsmBase: EsmBase{X: 10}, Y: 20}
	return d.GetEsmX() + d.Y // 10 + 20 = 30
}

// EmbeddedStructFieldPromotion tests field promotion
func EmbeddedStructFieldPromotion() int {
	type Inner struct{ Val int }
	type Outer struct {
		Inner
		Extra int
	}
	o := Outer{Inner: Inner{Val: 5}, Extra: 10}
	return o.Val + o.Extra // 5 + 10 = 15
}

// MyIntR13 for TypeAliasMethod test
type MyIntR13 int

func (m MyIntR13) DoubleR13() MyIntR13 { return m * 2 }

// TypeAliasMethod tests method on named type
func TypeAliasMethod() int {
	var x MyIntR13 = 5
	return int(x.DoubleR13()) // 10
}

// TypeAliasOperator tests operators on named types
func TypeAliasOperator() int {
	type MyInt int
	var a MyInt = 10
	var b MyInt = 3
	return int(a + b) // 13
}

// SliceOfInterfacesAppend tests append to []interface{}
func SliceOfInterfacesAppend() int {
	s := []interface{}{1, "hello", true}
	s = append(s, 3.14)
	return len(s) // 4
}

// MapInterfaceKey tests map with interface{} value
func MapInterfaceKeyValue() int {
	m := map[string]interface{}{
		"int":    42,
		"string": "hi",
		"bool":   true,
	}
	return len(m) // 3
}

// VariadicCallWithSlice tests calling variadic with ...
func VariadicCallWithSlice() int {
	nums := []int{1, 2, 3, 4, 5}
	return variadicSum(nums...)
}

func variadicSum(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

// VariadicCallDirect tests calling variadic directly
func VariadicCallDirect() int {
	return variadicSum(1, 2, 3) // 6
}

// VariadicCallEmpty tests calling variadic with no args
func VariadicCallEmpty() int {
	return variadicSum() // 0
}

// NestedSliceAppendR13 tests nested slice append
func NestedSliceAppendR13() int {
	grid := [][]int{{1, 2}, {3, 4}}
	grid = append(grid, []int{5, 6})
	return len(grid) + len(grid[2]) // 3 + 2 = 5
}

// StringRuneConversion tests string/rune conversion
func StringRuneConversion() int {
	s := "Hello世界"
	runes := []rune(s)
	return len(runes) // 7 (5 ASCII + 2 CJK)
}

// StringByteConversion tests string/byte conversion
func StringByteConversion() int {
	s := "Hello世界"
	b := []byte(s)
	return len(b) // 5 + 2*3 = 11 (UTF-8)
}

// ConstIotaExpressionR13 tests iota with expression
func ConstIotaExpressionR13() int {
	const (
		ar = 1 << iota // 1
		br             // 2
		cr             // 4
		dr             // 8
	)
	return ar + br + cr + dr // 1 + 2 + 4 + 8 = 15
}

// ConstIotaSkipR13 tests iota with blank identifier
func ConstIotaSkipR13() int {
	const (
		ar = iota // 0
		_         // 1 (skipped)
		br        // 2
		_         // 3 (skipped)
		cr        // 4
	)
	return ar + br + cr // 0 + 2 + 4 = 6
}

// BlankIdentifierInForR13 tests blank identifier in for loop
func BlankIdentifierInForR13() int {
	sum := 0
	for _ = range []int{1, 2, 3} {
		sum++
	}
	return sum // 3
}

// MapValueTypeAssertion tests type assertion on map value
func MapValueTypeAssertion() int {
	m := map[string]interface{}{"x": 42}
	v := m["x"].(int)
	return v // 42
}

// SliceIndexOutOfRange tests safe index access
func SliceIndexOutOfRange() int {
	s := []int{1, 2, 3}
	if len(s) > 5 {
		return s[5]
	}
	return -1 // safe
}

// ArrayValueSemantics tests array value copy
func ArrayValueSemantics() int {
	a1 := [3]int{1, 2, 3}
	a2 := a1
	a2[0] = 99
	return a1[0] // still 1
}

// PointerToArrayModify tests pointer to array modification
func PointerToArrayModify() int {
	a := [3]int{1, 2, 3}
	p := &a
	p[0] = 99
	return a[0] // 99
}

// StructFieldPointerModifyR13 tests modifying struct field through pointer
func StructFieldPointerModifyR13() int {
	type DataR13 struct{ Val int }
	d := DataR13{Val: 10}
	p := &d.Val
	*p = 20
	return d.Val // 20
}

// InterfaceEmptyTypeAssertion tests empty interface type assertion
func InterfaceEmptyTypeAssertion() string {
	var i interface{} = "test"
	return i.(string) // "test"
}

// InterfaceDoubleBoxing tests double interface boxing
func InterfaceDoubleBoxing() int {
	var i1 interface{} = 42
	var i2 interface{} = i1
	v := i2.(int)
	return v // 42
}

// SliceLiteralNested tests nested slice literal
func SliceLiteralNested() int {
	s := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	return s[1][2] // 6
}

// MapLiteralNested tests nested map literal
func MapLiteralNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
		"b": {"x": 3, "y": 4},
	}
	return m["a"]["x"] + m["b"]["y"] // 1 + 4 = 5
}

// StructLiteralNested tests nested struct literal
func StructLiteralNested() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 42}}
	return o.I.X // 42
}

// ForLoopWithContinue tests continue in for loop
func ForLoopWithContinue() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		sum += i
	}
	return sum // 1+3+5+7+9 = 25
}

// ForLoopWithBreak tests break in for loop
func ForLoopWithBreak() int {
	sum := 0
	for i := 0; i < 100; i++ {
		if i > 5 {
			break
		}
		sum += i
	}
	return sum // 0+1+2+3+4+5 = 15
}

// SwitchMultipleCases tests switch with multiple expressions per case
func SwitchMultipleCases() int {
	x := 3
	switch x {
	case 1, 2, 3:
		return 1
	case 4, 5:
		return 2
	default:
		return 0
	}
}

// ComplexSliceReslice tests reslice after append
func ComplexSliceReslice() int {
	s := make([]int, 3, 10)
	s[0], s[1], s[2] = 1, 2, 3
	s2 := s[:2]
	return len(s2) + cap(s2) // 2 + 10 = 12
}

// NilSliceAppendMultiple tests multiple appends to nil slice
func NilSliceAppendMultiple() int {
	var s []int
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return len(s) // 3
}

// MapIterateCount tests map iteration count
func MapIterateCount() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	count := 0
	for range m {
		count++
	}
	return count // 3
}

// NilN for StructMethodOnNil test
type NilN struct{ X int }

func (n *NilN) GetNilX() int {
	if n == nil {
		return -1
	}
	return n.X
}

// StructMethodOnNil tests calling method on nil struct
func StructMethodOnNil() int {
	var n *NilN = nil
	return n.GetNilX() // -1
}

// ClosureOverwriteVariable tests closure with overwritten variable
func ClosureOverwriteVariable() int {
	x := 10
	f := func() int { return x }
	x = 20
	return f() // 20
}

// SliceDeleteElement tests removing element from slice
func SliceDeleteElement() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return len(s) + s[2] // 4 + 4 = 8
}

// MapWithStructValue tests map with struct value
func MapWithStructValue() int {
	type Point struct{ X, Y int }
	m := map[string]Point{
		"origin": {0, 0},
		"p1":     {3, 4},
	}
	return m["p1"].X + m["p1"].Y // 3 + 4 = 7
}

// ChannelSendReceiveBuffered tests buffered channel send/receive
func ChannelSendReceiveBuffered() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	sum := <-ch + <-ch + <-ch
	return sum // 1 + 2 + 3 = 6
}

// MultipleReturnAssign tests multiple return value assignment
func MultipleReturnAssign() int {
	a, b := multiReturnR13()
	return a + b
}

func multiReturnR13() (int, int) {
	return 10, 20
}

// NamedReturnBare tests bare return with named values
func NamedReturnBare() (x int) {
	x = 42
	return // bare return
}

// TypeConversionChain tests chained type conversions
func TypeConversionChain() int {
	return int(int8(int16(100)))
}

// SPtrR13 for StructPointerMethodOnValue test
type SPtrR13 struct{ X int }

func (s *SPtrR13) DoubleR13() int { return s.X * 2 }

// StructPointerMethodOnValue tests calling pointer method on value
func StructPointerMethodOnValue() int {
	v := SPtrR13{X: 5}
	return v.DoubleR13() // 10 (Go auto-takes address)
}

// SliceFromToArray tests slicing an array
func SliceFromToArray() int {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]          // [2, 3, 4]
	return len(s) + s[0] // 3 + 2 = 5
}

// MapAccessTwoValue tests two-value map access
func MapAccessTwoValue() int {
	m := map[string]int{"a": 1}
	v, ok := m["a"]
	if ok {
		return v
	}
	return -1 // 1
}

// MapAccessTwoValueMissing tests two-value map access for missing key
func MapAccessTwoValueMissing() int {
	m := map[string]int{"a": 1}
	v, ok := m["b"]
	if ok {
		return v
	}
	return -1 // -1 (missing key)
}

// DeferArgumentsEvaluatedAtCall tests defer evaluates args immediately
func DeferArgumentsEvaluatedAtCall() (result int) {
	x := 1
	defer func(n int) { result = n }(x)
	x = 2
	return 0 // defer sets result to 1 (x was 1 when defer was called)
}

// StringConcatEmpty tests string concatenation with empty
func StringConcatEmpty() string {
	s := "" + "hello" + ""
	return s // "hello"
}

// SliceOfSliceAppend tests appending to inner slice
func SliceOfSliceAppend() int {
	s := [][]int{{1, 2}}
	s[0] = append(s[0], 3)
	return len(s[0]) // 3
}

// InterfaceSliceAccess tests accessing interface slice element
func InterfaceSliceAccess() int {
	s := []interface{}{1, "two", 3.0, true}
	return len(s) + s[0].(int) // 4 + 1 = 5
}

// PointerNilCheck tests nil pointer check
func PointerNilCheck() int {
	var p *int
	if p == nil {
		return 0
	}
	return *p
}

// StructZeroValueFields tests struct zero value fields
func StructZeroValueFields() int {
	type Data struct {
		A int
		B string
		C bool
		D []int
	}
	d := Data{}
	return d.A + len(d.B) // 0 + 0 = 0
}

// RangeOverNilMap tests ranging over nil map
func RangeOverNilMap() int {
	var m map[string]int
	count := 0
	for range m {
		count++
	}
	return count // 0
}

// RangeOverNilSlice tests ranging over nil slice
func RangeOverNilSlice() int {
	var s []int
	count := 0
	for range s {
		count++
	}
	return count // 0
}

// SliceCapAfterAppend tests capacity after multiple appends
func SliceCapAfterAppend() int {
	s := make([]int, 0)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	return len(s) // 3 (cap is implementation-defined, just check len)
}

// MapAssignExisting tests assigning to existing key
func MapAssignExisting() int {
	m := map[string]int{"a": 1}
	m["a"] = 42
	return m["a"] // 42
}

// StructCompareDifferentValues tests struct comparison with different values
func StructCompareDifferentValues() bool {
	type Point struct{ X, Y int }
	p1 := Point{1, 2}
	p2 := Point{1, 3}
	return p1 != p2 // true
}

// DeferModifyNamedReturn tests defer modifying named return value
func DeferModifyNamedReturn() (x int) {
	x = 5
	defer func() { x += 10 }()
	return x // return 5, then defer makes it 15
}

// ChannelCloseAndLen tests closed channel length
func ChannelCloseAndLen() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	close(ch)
	return len(ch) // 2
}

// MultipleShortVarDecl tests multiple short var declarations
func MultipleShortVarDecl() int {
	a, b := 1, 2
	c, d := 3, 4
	return a + b + c + d // 10
}

// StructEmbedShadow tests embedded field shadowing
func StructEmbedShadow() int {
	type Base struct{ X int }
	type Derived struct {
		Base
		X int // shadows Base.X
	}
	d := Derived{Base: Base{X: 1}, X: 2}
	return d.X + d.Base.X // 2 + 1 = 3
}

// ============================================================================
// ROUND 14: Package-level types for method tests
// ============================================================================

type base14Promoted struct{ X int }

func (b *base14Promoted) GetX() int { return b.X }

type derived14Promoted struct {
	base14Promoted
	Y int
}

type iface14R14 interface{ Do() int }

type impl14R14 struct{ X int }

func (i impl14R14) Do() int { return i.X * 2 }

type chain14R14 struct{ X int }

func (c chain14R14) Add(n int) chain14R14 { return chain14R14{X: c.X + n} }

type base14Override struct{ X int }

func (b base14Override) Val() int { return b.X }

type override14R struct {
	base14Override
}

func (o override14R) Val() int { return o.X * 10 }

type val14Double struct{ X int }

func (v *val14Double) Double() int { return v.X * 2 }

type builder14Chain struct{ val int }

func (b builder14Chain) Add(n int) builder14Chain { return builder14Chain{val: b.val + n} }

func (b builder14Chain) Mul(n int) builder14Chain { return builder14Chain{val: b.val * n} }

// ============================================================================
// ROUND 14: Additional Edge Cases
// ============================================================================

// SliceShareModifyR14 tests that sub-slices share underlying memory
func SliceShareModifyR14() int {
	a := []int{1, 2, 3, 4, 5}
	b := a[1:3] // [2,3]
	b[0] = 99
	return a[1] // 99 - shared memory
}

// SliceAppendOverwriteR14 tests append behavior with shared slices
func SliceAppendOverwriteR14() int {
	a := []int{1, 2, 3}
	b := a[0:2]
	b = append(b, 99)
	return a[2] // 3 - original not modified because append created new backing
}

// MapDeleteAndReadR14 tests reading after delete
func MapDeleteAndReadR14() int {
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	v, ok := m["a"]
	if ok {
		return -1
	}
	return v + m["b"] // 0 + 2 = 2
}

// MapLenAfterDeletesR14 tests map length after multiple deletes
func MapLenAfterDeletesR14() int {
	m := map[int]int{1: 10, 2: 20, 3: 30, 4: 40}
	delete(m, 1)
	delete(m, 3)
	return len(m) // 2
}

// PointerReassignR14 tests reassigning pointer target
func PointerReassignR14() int {
	x := 10
	y := 20
	p := &x
	v1 := *p
	p = &y
	v2 := *p
	return v1 + v2 // 30
}

// DoublePointerR14 tests pointer to pointer
func DoublePointerR14() int {
	x := 42
	p := &x
	pp := &p
	**pp = 99
	return x // 99
}

// SliceOfPointersR14 tests slice of pointers modification
func SliceOfPointersR14() int {
	a, b := 1, 2
	s := []*int{&a, &b}
	*s[0] = 10
	return a + *s[1] // 10 + 2 = 12
}

// MapOfSlicesR14 tests map with slice values
func MapOfSlicesR14() int {
	m := map[string][]int{
		"a": {1, 2, 3},
		"b": {4, 5},
	}
	return len(m["a"]) + len(m["b"]) // 3 + 2 = 5
}

// MapOfMapsR14 tests map with map values
func MapOfMapsR14() int {
	m := map[string]map[string]int{
		"outer": {"inner": 42},
	}
	return m["outer"]["inner"] // 42
}

// NestedMapLookupR14 tests nested map access chain
func NestedMapLookupR14() int {
	m := map[int]map[int]map[int]int{
		1: {2: {3: 100}},
	}
	return m[1][2][3] // 100
}

// StringIndexR14 tests string indexing returns byte
func StringIndexR14() int {
	s := "Hello"
	return int(s[1]) // 101 = 'e'
}

// RuneRangeStringR14 tests range over string yields runes
func RuneRangeStringR14() int {
	s := "日本"
	count := 0
	for range s {
		count++
	}
	return count // 2
}

// SliceCopyBuiltInR14 tests copy() builtin behavior
func SliceCopyBuiltInR14() int {
	src := []int{1, 2, 3}
	dst := make([]int, 2)
	n := copy(dst, src)
	return n + dst[0] + dst[1] // 2 + 1 + 2 = 5
}

// AppendByteSliceR14 tests append to byte slice
func AppendByteSliceR14() int {
	var b []byte
	b = append(b, 'H', 'i')
	return len(b) // 2
}

// StructCompareEqualR14 tests struct equality
func StructCompareEqualR14() bool {
	type P struct{ X, Y int }
	a := P{1, 2}
	b := P{1, 2}
	return a == b
}

// StructCompareNotEqualR14 tests struct inequality
func StructCompareNotEqualR14() bool {
	type P struct{ X, Y int }
	a := P{1, 2}
	b := P{1, 3}
	return a != b
}

// ArrayAssignCopyR14 tests that array assignment copies
func ArrayAssignCopyR14() int {
	var a [3]int = [3]int{1, 2, 3}
	b := a
	b[0] = 99
	return a[0] // 1 - array is copied
}

// SliceNilSubsliceR14 tests nil slice can be subsliced
func SliceNilSubsliceR14() int {
	var s []int
	s2 := s[0:0]
	return len(s2) // 0
}

// MapKeyPointerR14 tests map with pointer keys - not allowed directly, use struct
func MapKeyPointerR14() int {
	type Key struct{ X, Y int }
	m := map[Key]int{{1, 2}: 10, {3, 4}: 20}
	return m[Key{1, 2}] // 10
}

// ChannelCloseThenCapR14 tests cap of closed channel
func ChannelCloseThenCapR14() int {
	ch := make(chan int, 5)
	ch <- 1
	close(ch)
	return cap(ch) // 5
}

// TypeSwitchDefaultR14 tests type switch with default
func TypeSwitchDefaultR14() string {
	var i interface{} = 42
	switch i.(type) {
	case string:
		return "string"
	default:
		return "other"
	}
}

// InterfaceNilComparisonR14 tests nil interface comparison
func InterfaceNilComparisonR14() bool {
	var i interface{}
	return i == nil
}

// SliceGrowFactorR14 tests slice growth after append
func SliceGrowFactorR14() int {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	c1 := cap(s)
	s = append(s, 3) // triggers grow
	c2 := cap(s)
	if c2 > c1 {
		return 1
	}
	return 0
}

// MapPreAllocR14 tests map with pre-allocation
func MapPreAllocR14() int {
	m := make(map[int]int, 100)
	m[1] = 10
	m[2] = 20
	return len(m) // 2
}

// StructAnonymousFieldR14 tests struct with anonymous field
func StructAnonymousFieldR14() int {
	type Inner struct{ Val int }
	type Outer struct {
		Inner
		Name string
	}
	o := Outer{Inner: Inner{Val: 42}, Name: "test"}
	return o.Val // 42 - promoted field
}

// PointerSliceModifyR14 tests modify slice element via pointer
func PointerSliceModifyR14() int {
	s := []int{1, 2, 3}
	p := &s[1]
	*p = 99
	return s[1] // 99
}

// MapIterateDeleteR14 tests iterate and delete from map
func MapIterateDeleteR14() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	for k := range m {
		if k == 2 {
			delete(m, k)
		}
	}
	_, ok := m[2]
	if ok {
		return -1
	}
	return len(m) // 2
}

// SliceThreeIndexR14 tests 3-index slice expression
func SliceThreeIndexR14() int {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:3:4] // len=2, cap=3
	return len(s)*10 + cap(s) // 23
}

// ArrayPointerIndexR14 tests index through pointer to array
func ArrayPointerIndexR14() int {
	a := [3]int{10, 20, 30}
	p := &a
	return p[1] // 20
}

// StringByteSliceConvR14 tests string <-> []byte conversion
func StringByteSliceConvR14() int {
	s := "ABC"
	b := []byte(s)
	b[0] = 'X'
	result := string(b)
	if result == "XBC" {
		return 1
	}
	return 0
}

// RuneSliceStringConvR14 tests []rune <-> string conversion
func RuneSliceStringConvR14() int {
	s := "Hello"
	r := []rune(s)
	r[0] = 'J'
	result := string(r)
	if result == "Jello" {
		return 1
	}
	return 0
}

// ClosureReadBeforeWriteR14 tests closure reads before writes
func ClosureReadBeforeWriteR14() int {
	x := 10
	f := func() int {
		return x // reads outer x
	}
	x = 20
	return f() // 20
}

// MultiAssignDifferentTypesR14 tests multi assignment with different types
func MultiAssignDifferentTypesR14() int {
	var i int
	var f float64
	i, f = 1, 2.5
	return i + int(f) // 3
}

// SwapViaTupleR14 tests swap via tuple assignment
func SwapViaTupleR14() int {
	a, b := 1, 2
	a, b = b, a
	return a*10 + b // 21
}

// MapLookupOrDefaultR14 tests map lookup with comma-ok pattern
func MapLookupOrDefaultR14() int {
	m := map[string]int{"key": 42}
	if v, ok := m["key"]; ok {
		return v
	}
	return -1 // 42
}

// SliceRemovePreserveOrderR14 tests remove element preserving order
func SliceRemovePreserveOrderR14() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2 // remove index 2
	s = append(s[:i], s[i+1:]...)
	return len(s)*10 + s[2] // 34
}

// NestedStructLiteralR14 tests nested struct literal
func NestedStructLiteralR14() int {
	type Inner struct{ V int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{V: 99}}
	return o.I.V // 99
}

// SliceOfStringPointersR14 tests slice of string pointers
func SliceOfStringPointersR14() int {
	a, b := "hello", "world"
	s := []*string{&a, &b}
	return len(*s[0]) + len(*s[1]) // 5 + 5 = 10
}

// MapOfFuncValuesR14 tests map with function values
func MapOfFuncValuesR14() int {
	m := map[string]func(int) int{
		"double": func(x int) int { return x * 2 },
		"add1":   func(x int) int { return x + 1 },
	}
	return m["double"](5) + m["add1"](5) // 10 + 6 = 16
}

// StructPromotedMethodR14 tests promoted method from embedded struct
func StructPromotedMethodR14() int {
	d := derived14Promoted{base14Promoted: base14Promoted{X: 42}, Y: 10}
	return d.GetX() + d.Y // 52
}

// InterfaceSliceOfTypeR14 tests []interface{} with different types
func InterfaceSliceOfTypeR14() int {
	s := []interface{}{1, "hello", 3.14, true}
	return len(s) // 4
}

// ChannelDrainAfterCloseR14 tests drain channel after close
func ChannelDrainAfterCloseR14() int {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum // 60
}

// MapWithEmptyStructKeyR14 tests map with struct key
func MapWithEmptyStructKeyR14() int {
	type Key struct{ X, Y int }
	m := map[Key]int{}
	m[Key{1, 2}] = 100
	return m[Key{1, 2}] // 100
}

// PointerToSliceR14 tests pointer to slice, modify via pointer
func PointerToSliceR14() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 99
	return s[0] // 99
}

// SliceReverseR14 tests reversing a slice
func SliceReverseR14() int {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s[0]*10 + s[4] // 51
}

// ClosureCapturesNewVarR14 tests closure captures newly declared var
func ClosureCapturesNewVarR14() int {
	x := 1
	f := func() int {
		x = x + 10
		return x
	}
	return f() // 11
}

// MapClearAndReuseR14 tests clear map and add again
func MapClearAndReuseR14() int {
	m := map[int]int{1: 10, 2: 20}
	for k := range m {
		delete(m, k)
	}
	m[3] = 30
	return len(m) // 1
}

// SliceOfChannelsR14 tests slice of channels
func SliceOfChannelsR14() int {
	chs := make([]chan int, 2)
	chs[0] = make(chan int, 1)
	chs[1] = make(chan int, 1)
	chs[0] <- 10
	chs[1] <- 20
	return <-chs[0] + <-chs[1] // 30
}

// TypeAssertionMethodCallR14 tests type assertion then method call
func TypeAssertionMethodCallR14() int {
	var ifc iface14R14 = impl14R14{X: 21}
	if v, ok := ifc.(impl14R14); ok {
		return v.Do() // 42
	}
	return 0
}

// NestedForBreakLabelR14 tests break to label in nested for
func NestedForBreakLabelR14() int {
	sum := 0
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum += i*3 + j
			if i == 1 && j == 1 {
				break outer
			}
		}
	}
	return sum // 0+1+2+3+4 = 10
}

// RangeOverArrayR14 tests range over array (copy semantics)
func RangeOverArrayR14() int {
	var a [3]int = [3]int{10, 20, 30}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum // 60
}

// MapValueIsSliceR14 tests map value is a slice, append to it
func MapValueIsSliceR14() int {
	m := map[string][]int{"a": {1, 2}}
	m["a"] = append(m["a"], 3)
	return len(m["a"]) // 3
}

// PointerSwapR14 tests swap values via pointers
func PointerSwapR14() int {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return a*10 + b // 21
}

// SliceFromArrayOfIntR14 tests slice from array of int
func SliceFromArrayOfIntR14() int {
	a := [5]int{10, 20, 30, 40, 50}
	s := a[1:4] // [20, 30, 40]
	return s[0] + s[2] // 20 + 40 = 60
}

// StructWithEmbeddedSliceR14 tests struct with embedded slice field
func StructWithEmbeddedSliceR14() int {
	type S14 struct{ X []int }
	s := S14{X: []int{1, 2, 3}}
	s.X = append(s.X, 4)
	return len(s.X) // 4
}

// InterfaceAssertToPointerR14 tests interface assert to pointer type
func InterfaceAssertToPointerR14() int {
	x := 42
	var i interface{} = &x
	if p, ok := i.(*int); ok {
		return *p // 42
	}
	return 0
}

// MapKeyIsStringR14 tests map with string keys
func MapKeyIsStringR14() int {
	m := map[string]int{"alpha": 1, "beta": 2}
	return m["alpha"] + m["beta"] // 3
}

// SliceCapVsLenR14 tests slice cap vs len
func SliceCapVsLenR14() int {
	s := make([]int, 3, 10)
	return s[0]*0 + len(s)*10 + cap(s) // 0 + 30 + 10 = 40
}

// StructMethodReturnsSelfR14 tests method that returns self
func StructMethodReturnsSelfR14() int {
	c := chain14R14{X: 5}
	r := c.Add(3).Add(2)
	return r.X // 10
}

// ClosureAsArgumentR14 tests pass closure as function argument
func ClosureAsArgumentR14() int {
	apply := func(f func(int) int, x int) int { return f(x) }
	result := apply(func(x int) int { return x * x }, 7)
	return result // 49
}

// MapWithInterfaceKeyR14 tests map with interface{} key
func MapWithInterfaceKeyR14() int {
	m := map[interface{}]int{1: 10, "a": 20}
	return m[1] + m["a"] // 30
}

// SliceInsertMiddleR14 tests insert element in middle
func SliceInsertMiddleR14() int {
	s := []int{1, 2, 4, 5}
	i := 2
	s = append(s[:i], append([]int{3}, s[i:]...)...)
	return len(s)*10 + s[2] // 53
}

// PointerDerefAssignR14 tests deref pointer and assign
func PointerDerefAssignR14() int {
	x := 0
	p := &x
	*p = 42
	return x // 42
}

// RangeOverStringR14 tests range over string gives runes
func RangeOverStringR14() int {
	s := "Go语言"
	count := 0
	for range s {
		count++
	}
	return count // 4
}

// MapWithBoolKeyR14 tests map with bool key
func MapWithBoolKeyR14() int {
	m := map[bool]int{true: 1, false: 0}
	return m[true] + m[false] // 1
}

// SliceOfStructLiteralR14 tests slice of struct literals
func SliceOfStructLiteralR14() int {
	type Pt14 struct{ X, Y int }
	s := []Pt14{{1, 2}, {3, 4}, {5, 6}}
	return s[1].X + s[1].Y // 7
}

// ChannelAsReturnR14 tests function returning channel
func ChannelAsReturnR14() int {
	makeCh := func() chan int {
		ch := make(chan int, 1)
		ch <- 42
		return ch
	}
	ch := makeCh()
	return <-ch // 42
}

// StructEmbedMethodOverrideR14 tests embedded method override
func StructEmbedMethodOverrideR14() int {
	o := override14R{base14Override: base14Override{X: 5}}
	return o.Val() // 50
}

// InterfaceTypeSwitchR14 tests type switch on interface
func InterfaceTypeSwitchR14() int {
	var i interface{} = 3.14
	switch v := i.(type) {
	case int:
		return v
	case float64:
		return int(v * 100) // 314
	default:
		return 0
	}
}

// SliceFilterPatternR14 tests filter pattern on slice
func SliceFilterPatternR14() int {
	s := []int{1, 2, 3, 4, 5, 6}
	var result []int
	for _, v := range s {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return len(result)*10 + result[0] // 32
}

// PointerMethodOnStructR14 tests pointer method called on struct value
func PointerMethodOnStructR14() int {
	v := val14Double{X: 21}
	return v.Double() // 42
}

// MapDeleteNonExistentR14 tests delete non-existent key from map
func MapDeleteNonExistentR14() int {
	m := map[int]int{1: 10}
	delete(m, 999) // no-op
	return m[1] // 10
}

// SliceAppendNilR14 tests append nil to interface slice
func SliceAppendNilR14() int {
	s := []interface{}{1, "two"}
	s = append(s, nil)
	if s[2] == nil {
		return len(s) // 3
	}
	return 0
}

// StructWithSliceOfSelfR14 tests struct with slice of self type
func StructWithSliceOfSelfR14() int {
	type Node14 struct {
		Val      int
		Children []*Node14
	}
	child := &Node14{Val: 2}
	root := &Node14{Val: 1, Children: []*Node14{child}}
	return root.Val + root.Children[0].Val // 3
}

// ClosureCounterR14 tests closure-based counter
func ClosureCounterR14() int {
	newCounter := func() func() int {
		n := 0
		return func() int {
			n++
			return n
		}
	}
	c := newCounter()
	c()
	c()
	return c() // 3
}

// MapMergeR14 tests merge two maps
func MapMergeR14() int {
	m1 := map[int]int{1: 10, 2: 20}
	m2 := map[int]int{2: 200, 3: 300}
	for k, v := range m2 {
		m1[k] = v
	}
	return m1[2] + m1[3] // 200 + 300 = 500
}

// SliceDedupeR14 tests deduplicate sorted slice
func SliceDedupeR14() int {
	s := []int{1, 1, 2, 2, 3, 3, 3}
	var result []int
	prev := -1
	for _, v := range s {
		if v != prev {
			result = append(result, v)
			prev = v
		}
	}
	return len(result) // 3
}

// ChannelRangeAfterCloseR14 tests range over channel after close
func ChannelRangeAfterCloseR14() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum // 6
}

// TypeConversionStringToIntR14 tests string to int conversion via strconv
func TypeConversionStringToIntR14() int {
	s := "42"
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n // 42
}

// PointerToMapR14 tests pointer to map
func PointerToMapR14() int {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return len(m) // 2
}

// SliceCopyViaAssignmentR14 tests slice assignment doesn't copy elements
func SliceCopyViaAssignmentR14() int {
	a := []int{1, 2, 3}
	b := a
	b[0] = 99
	return a[0] // 99 - shared backing
}

// MapUpdateExistingKeyR14 tests update existing key in map
func MapUpdateExistingKeyR14() int {
	m := map[int]int{1: 10}
	m[1] = 20
	return m[1] // 20
}

// StructMethodChainR14 tests method chaining
func StructMethodChainR14() int {
	r := builder14Chain{val: 2}.Add(3).Mul(4) // (2+3)*4 = 20
	return r.val
}

// InterfaceEmptyVsNilR14 tests empty interface vs nil
func InterfaceEmptyVsNilR14() int {
	var i interface{}
	if i == nil {
		return 1
	}
	return 0 // 1
}

// SliceMapTransformR14 tests map/transform pattern on slice
func SliceMapTransformR14() int {
	s := []int{1, 2, 3}
	var result []int
	for _, v := range s {
		result = append(result, v*v)
	}
	return result[0] + result[1] + result[2] // 1 + 4 + 9 = 14
}

// PointerIncrementR14 tests increment through pointer
func PointerIncrementR14() int {
	x := 0
	p := &x
	*p++
	*p++
	*p++
	return *p // 3
}

// SliceZeroLenNonNilR14 tests make([]int, 0) is not nil
func SliceZeroLenNonNilR14() int {
	s := make([]int, 0)
	if s != nil {
		return 1
	}
	return 0 // 1
}

// MapReadWriteModifyR14 tests read, modify, write map
func MapReadWriteModifyR14() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	for k, v := range m {
		m[k] = v * 2
	}
	return m[1] + m[2] + m[3] // 20+40+60 = 120
}

// StructFieldPointerR14 tests struct field that is a pointer
func StructFieldPointerR14() int {
	type SP14 struct{ P *int }
	x := 42
	s := SP14{P: &x}
	return *s.P // 42
}

// SliceOfMapsR14 tests slice of maps
func SliceOfMapsR14() int {
	s := []map[string]int{
		{"a": 1},
		{"b": 2},
	}
	return s[0]["a"] + s[1]["b"] // 3
}

// NestedSliceAccessR14 tests deeply nested slice access
func NestedSliceAccessR14() int {
	s := [][][]int{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}
	return s[1][0][1] // 6
}

// MapKeyIntFloatR14 tests map with int key
func MapKeyIntFloatR14() int {
	m := map[int]float64{1: 1.1, 2: 2.2}
	return int(m[1]*10) + int(m[2]*10) // 11 + 22 = 33
}

// StructCompareDifferentFieldsR14 tests comparing struct with different field values
func StructCompareDifferentFieldsR14() bool {
	type Pt14D struct{ X, Y int }
	return Pt14D{1, 2} != Pt14D{1, 3}
}

// ChannelSelectWithTwoR14 tests select with two ready channels
func ChannelSelectWithTwoR14() int {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 10
	ch2 <- 20
	sum := 0
	for i := 0; i < 2; i++ {
		select {
		case v := <-ch1:
			sum += v
		case v := <-ch2:
			sum += v
		}
	}
	return sum // 30
}

// SliceFromSliceR14 tests slicing an existing slice
func SliceFromSliceR14() int {
	s := []int{0, 1, 2, 3, 4, 5}
	t := s[2:5] // [2,3,4]
	return t[0] + t[2] // 2 + 4 = 6
}

// MapLenZeroR14 tests len of empty map
func MapLenZeroR14() int {
	m := map[int]int{}
	return len(m) // 0
}

// StructPointerNilCheckR14 tests nil check on struct pointer
func StructPointerNilCheckR14() int {
	type St14N struct{ X int }
	var p *St14N
	if p == nil {
		return 1
	}
	return 0 // 1
}

// SliceAppendMultipleR14 tests append multiple elements
func SliceAppendMultipleR14() int {
	s := []int{1}
	s = append(s, 2, 3, 4, 5)
	return len(s) // 5
}

// MapAccessDefaultR14 tests accessing non-existent key returns zero
func MapAccessDefaultR14() int {
	m := map[int]int{1: 10}
	return m[999] // 0
}

// StructEmbedFieldAccessR14 tests embedded struct field access
func StructEmbedFieldAccessR14() int {
	type Inner14 struct{ A, B int }
	type Outer14 struct {
		Inner14
		C int
	}
	o := Outer14{Inner14: Inner14{A: 1, B: 2}, C: 3}
	return o.A + o.B + o.C // 6
}

// ClosureOverSliceR14 tests closure modifying slice
func ClosureOverSliceR14() int {
	s := []int{1, 2, 3}
	f := func() {
		s[0] = 99
	}
	f()
	return s[0] // 99
}

// MapSizeAfterClearR14 tests map size after clearing via delete
func MapSizeAfterClearR14() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := []string{"a", "b", "c"}
	for _, k := range keys {
		delete(m, k)
	}
	return len(m) // 0
}

// SliceOfBoolR14 tests slice of bool
func SliceOfBoolR14() int {
	s := []bool{true, false, true}
	count := 0
	for _, v := range s {
		if v {
			count++
		}
	}
	return count // 2
}

// StructZeroValCheckR14 tests zero value of struct
func StructZeroValCheckR14() int {
	type Z14 struct{ X int }
	var z Z14
	return z.X // 0
}

// MapWithIntKeyStringValR14 tests map[int]string
func MapWithIntKeyStringValR14() int {
	m := map[int]string{1: "one", 2: "two"}
	return len(m[1]) + len(m[2]) // 3 + 3 = 6
}

// SliceCopyFullR14 tests copying entire slice
func SliceCopyFullR14() int {
	src := []int{10, 20, 30}
	dst := make([]int, len(src))
	copy(dst, src)
	return dst[0] + dst[1] + dst[2] // 60
}

// StructPointerModifyFieldR14 tests modifying struct field via pointer
func StructPointerModifyFieldR14() int {
	type SP14M struct{ X int }
	s := SP14M{X: 10}
	p := &s
	p.X = 99
	return s.X // 99
}

// NestedMapDeleteR14 tests deleting from nested map
func NestedMapDeleteR14() int {
	m := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
	}
	delete(m["a"], "x")
	return len(m["a"]) // 1
}

// SliceOfStringR14 tests slice of strings
func SliceOfStringR14() int {
	s := []string{"hello", "world"}
	return len(s[0]) + len(s[1]) // 5 + 5 = 10
}

// MapRangeOrderIndependentR14 tests map range is order independent
func MapRangeOrderIndependentR14() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum // 60
}
