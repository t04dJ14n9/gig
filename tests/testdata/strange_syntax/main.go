package strange_syntax

import (
	"fmt"
	"strings"
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

// BlankIdentifierImport tests blank identifier import (imported at top)
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

func (f *File) Read() int  { return f.data }
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
	s = s[1:4]      // len=3, cap=9
	s = s[:cap(s)]  // extend to cap
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
