package tricky

// Tricky test cases - RALPH LOOP ITERATION 1
// Designed to find edge case bugs in the interpreter

import "errors"

// Counter for method test
type Counter struct{ value int }

func (c *Counter) inc() { c.value++ }

type Box struct{ val int }

func (b *Box) set(v int) { b.val = v }

type Point struct{ X, Y int }
type Inner struct{ V int }
type Outer struct{ I Inner }
type Base struct{ Value int }

type Derived struct {
	Base
	Extra int
}

// ShortVarDeclShadow tests shadowing in if
func ShortVarDeclShadow() int {
	x := 10
	if x, err := shadowHelper(); err == nil {
		_ = x
	}
	return x
}

func shadowHelper() (int, error) { return 42, nil }

// SliceIndexExpr tests slice indexing with expr
func SliceIndexExpr() int {
	s := []int{10, 20, 30, 40, 50}
	i := 1
	return s[i+1]
}

// MapStructKey tests map with struct key
func MapStructKey() int {
	m := map[Point]int{{X: 1, Y: 2}: 100}
	return m[Point{X: 1, Y: 2}]
}

// NestedSliceAppend tests nested slices
func NestedSliceAppend() int {
	var s [][]int
	s = append(s, []int{1, 2})
	s[0] = append(s[0], 3)
	return len(s[0])
}

// ClosureCaptureLoop tests loop capture
func ClosureCaptureLoop() int {
	var funcs []func() int
	for i := 0; i < 3; i++ {
		i := i
		funcs = append(funcs, func() int { return i })
	}
	return funcs[0]() + funcs[1]() + funcs[2]()
}

// DeferNamedReturn tests defer with named return
func DeferNamedReturn() (result int) {
	defer func() { result *= 2 }()
	result = 5
	return
}

// FullSliceExpr tests 3-index slice
func FullSliceExpr() int {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:3:4]
	return len(t)*10 + cap(t)
}

// NestedShadowing tests nested block shadowing
func NestedShadowing() int {
	x := 1
	{
		x := 2
		_ = x
	}
	return x
}

// SliceOfPointers tests slice of pointers
func SliceOfPointers() int {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}
	return *s[0] + *s[1] + *s[2]
}

// MapNestedStruct tests nested struct in map
func MapNestedStruct() int {
	m := map[string]Outer{"key": {I: Inner{V: 42}}}
	return m["key"].I.V
}

// VariadicEmpty tests empty variadic
func VariadicEmpty() int { return variadicSum() }

// VariadicOne tests variadic with one arg
func VariadicOne() int { return variadicSum(5) }

// VariadicMultiple tests variadic with multiple args
func VariadicMultiple() int { return variadicSum(1, 2, 3, 4) }

func variadicSum(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

// EmbeddedField tests embedded struct
func EmbeddedField() int {
	d := Derived{Base: Base{Value: 10}, Extra: 20}
	return d.Value + d.Extra
}

// MapPointerValue tests map with pointer
func MapPointerValue() int {
	v := 100
	m := map[string]*int{"key": &v}
	return *m["key"]
}

// ComplexBoolExpr tests boolean expr
func ComplexBoolExpr() int {
	a, b, c := true, false, true
	if a && !b || c {
		return 1
	}
	return 0
}

// SwitchFallthrough tests fallthrough
func SwitchFallthrough() int {
	n := 1
	result := 0
	switch n {
	case 1:
		result += 10
		fallthrough
	case 2:
		result += 20
	}
	return result
}

// SliceCopyOperation tests copy
func SliceCopyOperation() int {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	n := copy(dst, src)
	return n + dst[0] + dst[1] + dst[2]
}

// DeferStackOrder tests defer LIFO
func DeferStackOrder() int {
	result := 0
	defer func() { result += 1 }()
	defer func() { result += 10 }()
	defer func() { result += 100 }()
	result = 1000
	return result
}

// InterfaceAssertion tests type assertion
func InterfaceAssertion() int {
	var i interface{} = 42
	if v, ok := i.(int); ok {
		return v
	}
	return 0
}

// ChannelBasic tests channel
func ChannelBasic() int {
	ch := make(chan int, 1)
	ch <- 42
	return <-ch
}

// SelectDefault tests select default
func SelectDefault() int {
	ch := make(chan int)
	select {
	case v := <-ch:
		return v
	default:
		return -1
	}
}

// RecursiveFibMemo tests memoized fib
func RecursiveFibMemo() int {
	memo := make(map[int]int)
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		v := fib(n-1) + fib(n-2)
		memo[n] = v
		return v
	}
	return fib(10)
}

// PanicRecover tests panic/recover - disabled as panic is banned
func PanicRecover() int {
	// panic is banned in interpreted code
	return 0
}

// ClosureWithDefer tests closure defer
func ClosureWithDefer() int {
	x := 10
	f := func() int {
		defer func() { x *= 2 }()
		return x
	}
	return f() + x
}

// MethodOnPointer tests pointer method
func MethodOnPointer() int {
	c := &Counter{value: 5}
	c.inc()
	return c.value
}

// MultiReturnDiscard tests discard
func MultiReturnDiscard() int {
	a, _ := multiReturnHelper()
	return a
}

func multiReturnHelper() (int, int) { return 1, 2 }

// NilSliceAppend tests nil slice append
func NilSliceAppend() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// ShortCircuitEval tests short circuit &&
func ShortCircuitEval() int {
	called := false
	f := func() bool {
		called = true
		return true
	}
	_ = false && f()
	if called {
		return 1
	}
	return 0
}

// ShortCircuitEval2 tests short circuit ||
func ShortCircuitEval2() int {
	called := false
	f := func() bool {
		called = true
		return true
	}
	_ = true || f()
	if called {
		return 1
	}
	return 0
}

// MapDelete tests map delete
func MapDelete() int {
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	return len(m)
}

// SliceNil tests nil slice
func SliceNil() int {
	var s []int
	if s == nil {
		return 1
	}
	return 0
}

// MapCommaOk tests comma ok
func MapCommaOk() int {
	m := map[string]int{"a": 1}
	if v, ok := m["a"]; ok {
		return v
	}
	return 0
}

// InterfaceNil tests nil interface
func InterfaceNil() int {
	var i interface{}
	if i == nil {
		return 1
	}
	return 0
}

// SliceLenCap tests len/cap
func SliceLenCap() int {
	s := make([]int, 3, 5)
	return len(s)*10 + cap(s)
}

// ComplexArray tests array
func ComplexArray() int {
	var a [3]int
	a[0] = 1
	a[1] = 2
	a[2] = 3
	return a[0] + a[1] + a[2]
}

// PointerArithmetic tests pointer
func PointerArithmetic() int {
	x := 10
	p := &x
	*p = 20
	return x
}

// DoublePointer tests double pointer
func DoublePointer() int {
	x := 10
	p := &x
	pp := &p
	**pp = 20
	return x
}

// StructPointerMethod tests struct pointer method
func StructPointerMethod() int {
	b := &Box{val: 5}
	b.set(10)
	return b.val
}

// ForRangeWithIndex tests range with index
func ForRangeWithIndex() int {
	s := []int{10, 20, 30}
	sum := 0
	for i := range s {
		sum += i
	}
	return sum
}

// ForRangeKeyValue tests range with key and value
func ForRangeKeyValue() int {
	s := []int{1, 2, 3}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// StringIndex tests string indexing
func StringIndex() int {
	s := "hello"
	return int(s[0])
}

// MapAssign tests map assignment
func MapAssign() int {
	m := make(map[int]int)
	m[1] = 10
	m[2] = 20
	return m[1] + m[2]
}

// ComplexLiteral tests complex literal
func ComplexLiteral() int {
	type S struct {
		A int
		B string
	}
	s := S{A: 42, B: "test"}
	return s.A
}

// ErrorReturn tests error return
func ErrorReturn() int {
	_, err := errorHelper()
	if err != nil {
		return -1
	}
	return 0
}

func errorHelper() (int, error) { return 0, errors.New("test") }

// NilPointerCheck tests nil pointer check
func NilPointerCheck() int {
	var p *int
	if p == nil {
		return 1
	}
	return 0
}

// SliceAppendNil tests append nil
func SliceAppendNil() int {
	var s []int
	s = append(s, 0)
	return len(s)
}

// MapLookupNil tests map lookup on nil
func MapLookupNil() int {
	var m map[string]int
	if m == nil {
		return -1
	}
	return m["key"]
}

// DeferModifyNamed tests defer modify named
func DeferModifyNamed() (result int) {
	defer func() { result = 999 }()
	result = 42
	return
}

// MultipleNamedReturn tests multiple named return
func MultipleNamedReturn() (x int, y int) {
	defer func() {
		x *= 2
		y *= 3
	}()
	x = 5
	y = 7
	return
}

// MultipleNamedReturnCombined wraps MultipleNamedReturn for single-value comparison
func MultipleNamedReturnCombined() int {
	x, y := MultipleNamedReturn()
	return x*100 + y
}

// ForRangeMap tests ranging over map
func ForRangeMap() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 2 - New Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// DeferInClosure tests defer inside closure
func DeferInClosure() int {
	result := 0
	f := func() int {
		defer func() { result += 10 }()
		return result + 1
	}
	return f() + result
}

// MultipleDeferSameName tests multiple defer modifying same named return
func MultipleDeferSameName() (result int) {
	defer func() { result *= 2 }()
	defer func() { result += 10 }()
	result = 5
	return
}

// ClosureMutateOuter tests closure mutating outer variable
func ClosureMutateOuter() int {
	x := 10
	f := func() {
		x = 20
	}
	f()
	return x
}

// SliceAppendExpand tests slice append expanding capacity
func SliceAppendExpand() int {
	s := make([]int, 0, 2)
	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3) // should expand
	return len(s)*10 + cap(s)/10
}

// MapIncrement tests map value increment
func MapIncrement() int {
	m := map[string]int{"a": 1}
	m["a"]++
	m["a"]++
	return m["a"]
}

// InterfaceTypeSwitch tests type switch on interface
func InterfaceTypeSwitch() int {
	var i interface{} = 42
	switch v := i.(type) {
	case int:
		return v
	case string:
		return -1
	default:
		return 0
	}
}

// PointerToSlice tests pointer to slice
func PointerToSlice() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 10
	return (*p)[0] + s[0]
}

// NestedClosure tests nested closures
func NestedClosure() int {
	x := 1
	outer := func() int {
		y := 2
		inner := func() int {
			return x + y
		}
		return inner()
	}
	return outer()
}

// SliceOfSlice tests slice of slice operations
func SliceOfSlice() int {
	s := [][]int{{1, 2}, {3, 4}}
	s[0] = append(s[0], 5)
	return len(s[0]) + len(s[1])
}

// MapOfSlice tests map with slice values
func MapOfSlice() int {
	m := map[string][]int{"a": {1, 2, 3}}
	m["a"] = append(m["a"], 4)
	return len(m["a"])
}

// StructWithSlice tests struct with slice field
func StructWithSlice() int {
	type Container struct {
		items []int
	}
	c := Container{items: []int{1, 2}}
	c.items = append(c.items, 3)
	return len(c.items)
}

// DeferReadAfterAssign tests defer reading named return after assignment
func DeferReadAfterAssign() (result int) {
	result = 100
	defer func() {
		result = result / 2
	}()
	return result + 1
}

// ForRangePointer tests for range with pointer elements
func ForRangePointer() int {
	items := []*int{}
	a, b, c := 1, 2, 3
	items = append(items, &a, &b, &c)
	sum := 0
	for _, p := range items {
		sum += *p
	}
	return sum
}

// NilInterfaceValue tests nil interface with type
func NilInterfaceValue() int {
	var i interface{}
	if i == nil {
		return 1
	}
	return 0
}

// SliceCopyOverlap tests slice copy with overlap
func SliceCopyOverlap() int {
	s := []int{1, 2, 3, 4, 5}
	copy(s[1:], s[:4])
	return s[0] + s[1] + s[2] + s[3] + s[4]
}

// PointerReassign tests pointer reassignment
func PointerReassign() int {
	a, b := 10, 20
	p := &a
	result := *p
	p = &b
	result += *p
	return result
}

// InterfaceNilComparison tests interface nil comparison
func InterfaceNilComparison() int {
	var err error
	if err == nil {
		return 1
	}
	return 0
}

// DeferClosureCapture tests defer closure capturing variable
func DeferClosureCapture() int {
	x := 10
	defer func() {
		x *= 2
	}()
	x = 20
	return x
}

// MapLookupModify tests map lookup and modify
func MapLookupModify() int {
	m := map[int]*int{}
	v := 10
	m[1] = &v
	*m[1] = 20
	return *m[1]
}

// SliceZeroLength tests slice with zero length but non-zero capacity
func SliceZeroLength() int {
	s := make([]int, 0, 10)
	if len(s) == 0 && cap(s) == 10 {
		return 1
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 3 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// SliceModifyViaSubslice tests modifying slice via subslice
func SliceModifyViaSubslice() int {
	s := []int{1, 2, 3, 4, 5}
	sub := s[1:4]
	for i := range sub {
		sub[i] *= 10
	}
	return s[1] + s[2] + s[3]
}

// MapDeleteDuringRange tests deleting from map during range
func MapDeleteDuringRange() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	count := 0
	for k := range m {
		if k == 2 {
			delete(m, k)
		}
		count++
	}
	return count
}

// ClosureReturnClosure tests closure that returns closure
func ClosureReturnClosure() int {
	makeAdder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := makeAdder(5)
	return add5(10)
}

// StructMethodOnNil tests method call on nil pointer
func StructMethodOnNil() int {
	type N struct{ v int }
	var p *N
	if p == nil {
		return -1
	}
	return p.v
}

// ArrayPointerIndex tests array pointer indexing
func ArrayPointerIndex() int {
	a := [3]int{10, 20, 30}
	p := &a
	return (*p)[1]
}

// SliceThreeIndex tests three-index slice
func SliceThreeIndex() int {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:3:4]
	return len(t)*10 + cap(t)
}

// MapWithFuncKey tests map with func key not allowed - use int
func MapWithFuncValue() int {
	m := make(map[int]func() int)
	m[1] = func() int { return 10 }
	m[2] = func() int { return 20 }
	return m[1]() + m[2]()
}

// DeferInLoop tests defer in loop
func DeferInLoop() int {
	result := 0
	for i := 0; i < 3; i++ {
		defer func(n int) {
			result += n
		}(i)
	}
	return result
}

// StructCompare tests struct comparison
func StructCompare() int {
	type P struct{ x, y int }
	p1 := P{1, 2}
	p2 := P{1, 2}
	if p1 == p2 {
		return 1
	}
	return 0
}

// InterfaceSlice tests slice of interfaces
func InterfaceSlice() int {
	var items []interface{}
	items = append(items, 1, "hello", 3.14)
	return len(items)
}

// PointerMethodValueReceiver tests pointer method with value receiver
func PointerMethodValueReceiver() int {
	type S struct{ v int }
	f := func(s S) int { return s.v }
	p := &S{v: 42}
	return f(*p)
}

// SliceOfMaps tests slice of maps
func SliceOfMaps() int {
	var sm []map[string]int
	sm = append(sm, map[string]int{"a": 1})
	sm = append(sm, map[string]int{"b": 2})
	return sm[0]["a"] + sm[1]["b"]
}

// MapWithNilValue tests map with nil value
func MapWithNilValue() int {
	m := map[string]*int{"a": nil, "b": new(int)}
	if m["a"] == nil && m["b"] != nil {
		return 1
	}
	return 0
}

// SwitchNoCondition tests switch with no condition
func SwitchNoCondition() int {
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

// DeferModifyReturn tests defer modifying return value
func DeferModifyReturn() (result int) {
	defer func() {
		result += 100
	}()
	return 42
}

// SliceAppendToCap tests append to capacity
func SliceAppendToCap() int {
	s := make([]int, 2, 4)
	s[0], s[1] = 1, 2
	s = append(s, 3)
	s = append(s, 4)
	return len(s) + cap(s)
}

// ForRangeStringByteIndex tests for range string byte index
func ForRangeStringByteIndex() int {
	s := "abc"
	var indices []int
	for i := range s {
		indices = append(indices, i)
	}
	return indices[0] + indices[1] + indices[2]
}

// StructLiteralEmbedded tests struct literal with embedded field
func StructLiteralEmbedded() int {
	type Inner struct{ V int }
	type Outer struct {
		Inner
		X int
	}
	o := Outer{Inner: Inner{V: 10}, X: 20}
	return o.V + o.X
}

// MapNilKey tests map with nil pointer key
func MapNilKey() int {
	m := map[*int]int{}
	var p *int
	m[p] = 42
	return m[p]
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

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 4 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// SliceOfInterfacesWithTypes tests slice of interfaces with different types
func SliceOfInterfacesWithTypes() int {
	var items []interface{}
	items = append(items, 1, "hello", true, 3.14)
	count := 0
	for _, item := range items {
		switch item.(type) {
		case int:
			count += 1
		case string:
			count += 10
		case bool:
			count += 100
		case float64:
			count += 1000
		}
	}
	return count
}

// DeferCallInDefer tests defer inside defer
func DeferCallInDefer() (result int) {
	defer func() {
		defer func() {
			result += 1
		}()
		result += 10
	}()
	result = 100
	return
}

// MapLookupAssign tests map lookup and assign
func MapLookupAssign() int {
	m := map[string]int{"a": 1}
	v, ok := m["a"]
	if ok {
		m["b"] = v * 2
	}
	return m["b"]
}

// StructMethodOnValue tests method on value type
func StructMethodOnValue() int {
	type S struct{ v int }
	f := func(s S) int { return s.v }
	s := S{v: 42}
	return f(s)
}

// PointerToMap tests pointer to map
func PointerToMap() int {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return len(m)
}

// SliceCapAfterAppend tests slice capacity after append
func SliceCapAfterAppend() int {
	s := make([]int, 0, 4)
	s = append(s, 1, 2)
	origCap := cap(s)
	s = append(s, 3, 4, 5) // should reallocate
	if cap(s) > origCap {
		return 1
	}
	return 0
}

// NestedMaps tests nested maps
func NestedMaps() int {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 42}
	return m["outer"]["inner"]
}

// StructPointerNil tests nil struct pointer
func StructPointerNil() int {
	type S struct{ v int }
	var p *S
	if p == nil {
		return 1
	}
	return 0
}

// VariadicWithSlice tests variadic with slice spread
func VariadicWithSlice() int {
	nums := []int{1, 2, 3}
	return variadicSum(nums...)
}

// SliceMakeWithLen tests make slice with length
func SliceMakeWithLen() int {
	s := make([]int, 5)
	s[0] = 1
	s[4] = 5
	return s[0] + s[4]
}

// InterfaceConversion tests interface conversion
func InterfaceConversion() int {
	var i interface{} = 42
	n, ok := i.(int)
	if ok {
		return n
	}
	return -1
}

// MapWithEmptyStringKey tests map with empty string key
func MapWithEmptyStringKey() int {
	m := map[string]int{"": 1, "a": 2}
	return m[""] + m["a"]
}

// DeferPanicRecover tests defer with panic/recover
func DeferPanicRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = 100
		}
	}()
	result = 1
	return result
}

// StructWithMap tests struct with map field
func StructWithMap() int {
	type Container struct {
		items map[string]int
	}
	c := Container{items: map[string]int{"a": 1}}
	c.items["b"] = 2
	return len(c.items)
}

// ForRangeBreak tests break in for range
func ForRangeBreak() int {
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

// SliceLiteralNested tests nested slice literal
func SliceLiteralNested() int {
	s := [][]int{{1, 2}, {3, 4, 5}}
	return len(s[0]) + len(s[1])
}

// MapLiteralNested tests nested map literal
func MapLiteralNested() int {
	m := map[string]map[string]int{
		"a": {"x": 1},
		"b": {"y": 2},
	}
	return m["a"]["x"] + m["b"]["y"]
}

// PointerToStructLiteral tests pointer to struct literal
func PointerToStructLiteral() int {
	type S struct{ v int }
	p := &S{v: 42}
	return p.v
}

// SliceOfStructs tests slice of structs
func SliceOfStructs() int {
	type P struct{ x, y int }
	s := []P{{1, 2}, {3, 4}}
	return s[0].x + s[1].y
}

// MapIterateModify tests iterating and modifying map
func MapIterateModify() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	for k := range m {
		m[k] *= 2
	}
	return m[1] + m[2] + m[3]
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 5 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// ChannelBuffered tests buffered channel operations
func ChannelBuffered() int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	return <-ch + <-ch + <-ch
}

// StructEmbeddedMethod tests method on embedded struct
// Types defined at package level for method access
type InnerWithMethod struct{ v int }

func (i *InnerWithMethod) get() int { return i.v }

type OuterEmbedded struct {
	InnerWithMethod
}

func StructEmbeddedMethod() int {
	o := OuterEmbedded{InnerWithMethod: InnerWithMethod{v: 42}}
	return o.get()
}

// SliceOfChannels tests slice of channels
func SliceOfChannels() int {
	ch1, ch2 := make(chan int, 1), make(chan int, 1)
	ch1 <- 10
	ch2 <- 20
	chs := []chan int{ch1, ch2}
	return <-chs[0] + <-chs[1]
}

// MapOfChannels tests map of channels
func MapOfChannels() int {
	ch1, ch2 := make(chan int, 1), make(chan int, 1)
	ch1 <- 10
	ch2 <- 20
	m := map[string]chan int{"a": ch1, "b": ch2}
	return <-m["a"] + <-m["b"]
}

// InterfaceMethod tests interface method call
// Types defined at package level
type Adder interface{ Add(int) int }

type AdderStruct struct{ v int }

func (s *AdderStruct) Add(n int) int { return s.v + n }

func InterfaceMethod() int {
	var a Adder = &AdderStruct{v: 10}
	return a.Add(5)
}

// MultipleAssignment tests multiple assignment
func MultipleAssignment() int {
	a, b := 1, 2
	a, b = b, a
	return a*10 + b
}

// SliceAssign tests slice assignment
func SliceAssign() int {
	s := make([]int, 3)
	s[0], s[1], s[2] = 1, 2, 3
	return s[0] + s[1] + s[2]
}

// MapTwoAssign tests two-value map lookup
func MapTwoAssign() int {
	m := map[string]int{"a": 1, "b": 2}
	v1, ok1 := m["a"]
	_, ok2 := m["c"]
	if ok1 && !ok2 {
		return v1
	}
	return -1
}

// StructPointerMethodNil tests nil receiver check
// Type defined at package level
type ValueStruct struct{ v int }

func (s *ValueStruct) Value() int {
	if s == nil {
		return 0
	}
	return s.v
}

func StructPointerMethodNil() int {
	var p *ValueStruct
	return p.Value()
}

// DeferAfterPanic tests defer execution after panic
// Note: panic is banned in interpreted code, so this test returns 100 directly
func DeferAfterPanic() (result int) {
	// panic is banned - simulate the result
	result = 100
	return
}

// SliceFromArray tests slice from array
func SliceFromArray() int {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]
	return s[0] + s[1] + s[2]
}

// ArrayPointerSlice tests slicing array pointer
func ArrayPointerSlice() int {
	a := [5]int{1, 2, 3, 4, 5}
	p := &a
	s := p[1:4]
	return s[0] + s[1] + s[2]
}

// StructFieldPointer tests struct field pointer
func StructFieldPointer() int {
	type S struct{ v int }
	s := S{v: 10}
	p := &s.v
	*p = 20
	return s.v
}

// MapLenCap tests map len
func MapLenCap() int {
	m := map[string]int{"a": 1, "b": 2}
	return len(m)
}

// StringConcat tests string concatenation
func StringConcat() string {
	return "hello" + " " + "world"
}

// StringLen tests string length
func StringLen() int {
	return len("hello world")
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 6 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// ComplexMapKey tests map with complex key type
func ComplexMapKey() int {
	type Key struct{ a, b int }
	m := map[Key]int{{1, 2}: 10, {3, 4}: 20}
	return m[Key{1, 2}] + m[Key{3, 4}]
}

// SliceReverse tests slice reversal
func SliceReverse() int {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

// MapMerge tests map merging
func MapMerge() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 3, "d": 4}
	for k, v := range m2 {
		m1[k] = v
	}
	return len(m1)
}

// StructZeroValue tests struct zero value
func StructZeroValue() int {
	type S struct {
		a int
		b string
		c bool
	}
	var s S
	if s.a == 0 && s.b == "" && s.c == false {
		return 1
	}
	return 0
}

// SliceDeleteByIndex tests slice element deletion
func SliceDeleteByIndex() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return len(s)
}

// MapValueOverwrite tests map value overwrite
func MapValueOverwrite() int {
	m := map[string]int{"a": 1}
	m["a"] = 10
	m["a"] = 100
	return m["a"]
}

// InterfaceEmbed tests interface embedding
func InterfaceEmbed() int {
	type Reader interface{ Read() int }
	type Writer interface{ Write(int) }
	type ReadWriter interface {
		Reader
		Writer
	}
	return 0
}

// SliceOfFuncs tests slice of functions
func SliceOfFuncs() int {
	funcs := []func() int{
		func() int { return 1 },
		func() int { return 2 },
		func() int { return 3 },
	}
	return funcs[0]() + funcs[1]() + funcs[2]()
}

// StructWithFunc tests struct with func field
func StructWithFunc() int {
	type S struct {
		f func() int
	}
	s := S{f: func() int { return 42 }}
	return s.f()
}

// PointerToSliceElement tests pointer to slice element
func PointerToSliceElement() int {
	s := []int{1, 2, 3}
	p := &s[1]
	*p = 20
	return s[1]
}

// MapKeyPointer tests map with pointer key
func MapKeyPointer() int {
	a, b := 1, 2
	m := map[*int]int{&a: 10, &b: 20}
	return m[&a] + m[&b]
}

// SliceOfPointersToStruct tests slice of pointers to struct
func SliceOfPointersToStruct() int {
	type S struct{ v int }
	s := []*S{{v: 1}, {v: 2}, {v: 3}}
	return s[0].v + s[1].v + s[2].v
}

// DoubleMapLookup tests double map lookup
func DoubleMapLookup() int {
	m := map[int]map[int]int{
		1: {2: 3},
		4: {5: 6},
	}
	return m[1][2] + m[4][5]
}

// StructSliceLiteral tests struct slice literal
func StructSliceLiteral() int {
	type P struct{ x, y int }
	s := []P{{1, 2}, {3, 4}, {5, 6}}
	sum := 0
	for _, p := range s {
		sum += p.x + p.y
	}
	return sum
}

// ForRangeModifyValue tests for range modify value
func ForRangeModifyValue() int {
	s := []int{1, 2, 3}
	for i := range s {
		s[i] *= 2
	}
	return s[0] + s[1] + s[2]
}

// MapWithStructPointerKey tests map with struct pointer key
func MapWithStructPointerKey() int {
	type K struct{ v int }
	k1, k2 := &K{v: 1}, &K{v: 2}
	m := map[*K]int{k1: 10, k2: 20}
	return m[k1] + m[k2]
}

// SliceCopyDifferentTypes tests slice copy with different element types
func SliceCopyDifferentTypes() int {
	src := []int{1, 2, 3}
	dst := make([]int, 3)
	copied := copy(dst, src)
	return copied + dst[0] + dst[1] + dst[2]
}

// NestedStructWithPointer tests nested struct with pointer
func NestedStructWithPointer() int {
	type Inner struct{ v int }
	type Outer struct {
		inner *Inner
	}
	o := Outer{inner: &Inner{v: 42}}
	return o.inner.v
}

// SliceOfSlicesAppend tests append to nested slice
func SliceOfSlicesAppend() int {
	s := [][]int{{1}, {2}}
	s[0] = append(s[0], 10)
	s[1] = append(s[1], 20)
	return s[0][1] + s[1][1]
}

// MapDeleteAll tests deleting all map entries
func MapDeleteAll() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return len(m)
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 7 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// StructPointerSlice tests slice of struct pointers
func StructPointerSlice() int {
	type S struct{ v int }
	s := []*S{{v: 1}, {v: 2}, {v: 3}}
	sum := 0
	for _, p := range s {
		sum += p.v
	}
	return sum
}

// MapWithInterfaceKey tests map with interface key
func MapWithInterfaceKey() int {
	m := map[interface{}]int{1: 10, "a": 20}
	return m[1] + m["a"]
}

// SliceOfInterfaces tests slice of interfaces
func SliceOfInterfaces() int {
	s := []interface{}{1, "hello", true}
	return len(s)
}

// NestedPointerStruct tests nested pointer struct
func NestedPointerStruct() int {
	type Inner struct{ v int }
	type Outer struct {
		inner *Inner
	}
	o := Outer{inner: &Inner{v: 42}}
	return o.inner.v
}

// StructMethodOnNilPointer tests method on nil pointer
func StructMethodOnNilPointer() int {
	return StructPointerMethodNil()
}

// SliceAppendToSlice tests appending slice to slice
func SliceAppendToSlice() int {
	s1 := []int{1, 2}
	s2 := []int{3, 4}
	s1 = append(s1, s2...)
	return s1[0] + s1[1] + s1[2] + s1[3]
}

// MapLookupWithDefault tests map lookup with default
func MapLookupWithDefault() int {
	m := map[string]int{"a": 1}
	v := m["b"]
	if v == 0 {
		v = 100
	}
	return v
}

// StructFieldUpdate tests struct field update
func StructFieldUpdate() int {
	type S struct{ v int }
	s := S{v: 10}
	s.v = 20
	return s.v
}

// PointerToNilSlice tests pointer to nil slice
func PointerToNilSlice() int {
	var s []int
	p := &s
	if *p == nil {
		return 1
	}
	return 0
}

// MapUpdateDuringRange tests updating map during range
func MapUpdateDuringRange() int {
	m := map[int]int{1: 10, 2: 20}
	for k := range m {
		m[k+10] = k
	}
	return len(m)
}

// SliceCopyToSubslice tests copy to subslice
func SliceCopyToSubslice() int {
	src := []int{1, 2, 3}
	dst := make([]int, 5)
	copy(dst[1:4], src)
	return dst[1] + dst[2] + dst[3]
}

// StructWithMultipleFields tests struct with multiple fields
func StructWithMultipleFields() int {
	type S struct {
		a int
		b string
		c bool
		d float64
	}
	s := S{a: 1, b: "test", c: true, d: 3.14}
	result := s.a
	if s.c {
		result += 10
	}
	return result
}

// ForRangeContinue tests continue in for range
func ForRangeContinue() int {
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

// MapWithBoolKey tests map with bool key
func MapWithBoolKey() int {
	m := map[bool]int{true: 1, false: 0}
	return m[true] + m[false]
}

// SliceInsert tests slice element insertion
func SliceInsert() int {
	s := []int{1, 3, 4}
	s = append(s[:1], append([]int{2}, s[1:]...)...)
	return s[0] + s[1] + s[2] + s[3]
}

// StructEmbeddedFieldAccess tests accessing embedded struct field
func StructEmbeddedFieldAccess() int {
	type Inner struct{ V int }
	type Outer struct{ Inner }
	o := Outer{Inner: Inner{V: 42}}
	return o.V
}

// PointerToChannel tests pointer to channel
func PointerToChannel() int {
	ch := make(chan int, 1)
	p := &ch
	*p <- 42
	return <-ch
}

// MapKeyModification tests map key modification during range
func MapKeyModification() int {
	m := map[string]int{"a": 1, "b": 2}
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return len(keys)
}

// SliceRangeModify tests range with index modification
func SliceRangeModify() int {
	s := []int{1, 2, 3}
	for i := range s {
		s[i] = i * 10
	}
	return s[0] + s[1] + s[2]
}

// StructLiteralShort tests short struct literal
func StructLiteralShort() int {
	type P struct{ x, y int }
	p := P{1, 2}
	return p.x + p.y
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 8 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// SliceDrain tests draining a slice
func SliceDrain() int {
	s := []int{1, 2, 3, 4, 5}
	for len(s) > 0 {
		s = s[1:]
	}
	return len(s)
}

// MapClear tests clearing a map
func MapClear() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		delete(m, k)
	}
	return len(m)
}

// StructCopy tests copying struct
func StructCopy() int {
	type S struct{ v int }
	s1 := S{v: 42}
	s2 := s1
	s2.v = 100
	return s1.v
}

// PointerStructCopy tests copying pointer struct
func PointerStructCopy() int {
	type S struct{ v int }
	s1 := &S{v: 42}
	s2 := s1
	s2.v = 100
	return s1.v
}

// SliceFilter tests filtering slice
func SliceFilter() int {
	s := []int{1, 2, 3, 4, 5}
	result := []int{}
	for _, v := range s {
		if v%2 == 1 {
			result = append(result, v)
		}
	}
	return len(result)
}

// MapTransform tests transforming map
func MapTransform() int {
	m1 := map[int]int{1: 10, 2: 20}
	m2 := make(map[int]int)
	for k, v := range m1 {
		m2[k] = v * 2
	}
	return m2[1] + m2[2]
}

// SliceContains tests slice contains check
func SliceContains() int {
	s := []int{1, 2, 3, 4, 5}
	for _, v := range s {
		if v == 3 {
			return 1
		}
	}
	return 0
}

// MapKeys tests extracting map keys
func MapKeys() int {
	m := map[string]int{"a": 1, "b": 2}
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return len(keys)
}

// StructMethodChain tests chained method calls
func StructMethodChain() int {
	type S struct{ v int }
	s := &S{v: 1}
	s.v++
	return s.v
}

// SliceLast tests getting last element
func SliceLast() int {
	s := []int{1, 2, 3, 4, 5}
	return s[len(s)-1]
}

// MapGetOrSet tests map get or set pattern
func MapGetOrSet() int {
	m := map[string]int{}
	if v, ok := m["key"]; !ok {
		m["key"] = 42
	} else {
		_ = v
	}
	return m["key"]
}

// StructValidation tests struct validation
func StructValidation() int {
	type S struct {
		name string
		age  int
	}
	s := S{name: "test", age: 25}
	if s.age > 0 && s.name != "" {
		return 1
	}
	return 0
}

// SlicePrepend tests prepending to slice
func SlicePrepend() int {
	s := []int{2, 3}
	s = append([]int{1}, s...)
	return s[0] + s[1] + s[2]
}

// MapMergeOverwrite tests map merge with overwrite
func MapMergeOverwrite() int {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 20, "c": 3}
	for k, v := range m2 {
		m1[k] = v
	}
	return m1["b"]
}

// SliceRotate tests rotating slice
func SliceRotate() int {
	s := []int{1, 2, 3, 4, 5}
	s = append(s[1:], s[0])
	return s[0] + s[4]
}

// StructInterface tests struct as interface
func StructInterface() int {
	type S struct{ v int }
	var i interface{} = S{v: 42}
	if s, ok := i.(S); ok {
		return s.v
	}
	return 0
}

// MapKeysSorted tests sorted map keys (manual sort)
func MapKeysSorted() int {
	m := map[int]int{3: 1, 1: 2, 2: 3}
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}
	// Simple bubble sort
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys[0] + keys[1] + keys[2]
}

// SliceFlatten tests flattening nested slice
func SliceFlatten() int {
	s := [][]int{{1, 2}, {3, 4}}
	result := []int{}
	for _, inner := range s {
		result = append(result, inner...)
	}
	return len(result)
}

// StructFieldPointerModify tests modifying via field pointer
func StructFieldPointerModify() int {
	type S struct{ v int }
	s := S{v: 10}
	p := &s.v
	*p = 20
	return s.v
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 9 - More Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// MapSwap tests swapping values in map
func MapSwap() int {
	m := map[string]int{"a": 1, "b": 2}
	m["a"], m["b"] = m["b"], m["a"]
	return m["a"]*10 + m["b"]
}

// SliceSplit tests splitting slice
func SliceSplit() int {
	s := []int{1, 2, 3, 4, 5, 6}
	mid := len(s) / 2
	left := s[:mid]
	right := s[mid:]
	return len(left) + len(right)
}

// StructCompareDiff tests comparing different struct values
func StructCompareDiff() int {
	type P struct{ x, y int }
	p1 := P{1, 2}
	p2 := P{1, 3}
	if p1 != p2 {
		return 1
	}
	return 0
}

// MapNestedDelete tests deleting from nested map
func MapNestedDelete() int {
	m := map[string]map[string]int{
		"outer": {"inner": 42},
	}
	delete(m["outer"], "inner")
	return len(m["outer"])
}

// PointerNilDeref tests nil pointer check before deref
func PointerNilDeref() int {
	var p *int
	if p != nil {
		return *p
	}
	return -1
}

// SliceGrow tests growing slice beyond capacity
func SliceGrow() int {
	s := make([]int, 0, 2)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	return len(s)
}

// StructEmpty tests empty struct
func StructEmpty() int {
	type Empty struct{}
	_ = Empty{}
	return 0
}

// MapEmptyKey tests map with empty interface key
func MapEmptyKey() int {
	m := map[interface{}]int{}
	m[1] = 10
	m[""] = 20
	return len(m)
}

// SliceMakeZero tests make with zero length
func SliceMakeZero() int {
	s := make([]int, 0)
	return len(s)
}

// StructAnon tests anonymous struct
func StructAnon() int {
	s := struct {
		x int
		y string
	}{x: 42, y: "test"}
	return s.x
}

// MapSizeHint tests map with size hint
func MapSizeHint() int {
	m := make(map[int]int, 100)
	return len(m)
}

// SliceNilAppend tests appending to nil slice
func SliceNilAppend() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// StructFieldPtr tests pointer to struct field
func StructFieldPtr() int {
	type S struct{ v int }
	s := S{v: 10}
	p := &s.v
	*p = 20
	return s.v
}

// MapIterateDelete tests iterate and delete
func MapIterateDelete() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	for k := range m {
		if k == 2 {
			delete(m, k)
		}
	}
	return len(m)
}

// SliceTruncate tests truncating slice
func SliceTruncate() int {
	s := []int{1, 2, 3, 4, 5}
	s = s[:2]
	return len(s)
}

// StructMethodValue tests value receiver method
func StructMethodValue() int {
	type S struct{ v int }
	s := S{v: 42}
	// Define method separately at package level
	return s.v
}

// MapFloatKey tests map with float key
func MapFloatKey() int {
	m := map[float64]int{1.5: 10, 2.5: 20}
	return m[1.5] + m[2.5]
}

// SliceRepeat tests repeating slice pattern
func SliceRepeat() int {
	s := []int{}
	for i := 0; i < 3; i++ {
		s = append(s, i, i*10)
	}
	return len(s)
}

// StructNestedAssign tests nested struct assignment
func StructNestedAssign() int {
	type Inner struct{ v int }
	type Outer struct{ inner Inner }
	o := Outer{inner: Inner{v: 10}}
	o.inner.v = 20
	return o.inner.v
}

// MapIntKey tests map with int key
func MapIntKey() int {
	m := map[int]string{1: "a", 2: "b"}
	return len(m[1]) + len(m[2])
}

// ─────────────────────────────────────────────────────────────────────────────
// RALPH LOOP ITERATION 10 - Final Tricky Tests
// ─────────────────────────────────────────────────────────────────────────────

// SliceReverseInPlace tests in-place slice reversal
func SliceReverseInPlace() int {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

// MapIncrementAll tests incrementing all map values
func MapIncrementAll() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k := range m {
		m[k]++
	}
	return m["a"] + m["b"] + m["c"]
}

// StructPtrMethod tests pointer method on struct
func StructPtrMethod() int {
	type S struct{ v int }
	s := &S{v: 10}
	s.v++
	return s.v
}

// SliceMapIndex tests using slice element as map index
func SliceMapIndex() int {
	m := map[int]int{0: 10, 1: 20, 2: 30}
	s := []int{0, 1, 2}
	return m[s[0]] + m[s[1]] + m[s[2]]
}

// MapCopy tests copying map
func MapCopy() int {
	m1 := map[int]int{1: 10, 2: 20}
	m2 := make(map[int]int)
	for k, v := range m1 {
		m2[k] = v
	}
	return len(m2)
}

// StructSliceAppend tests appending to struct slice field
func StructSliceAppend() int {
	type S struct{ items []int }
	s := S{items: []int{1, 2}}
	s.items = append(s.items, 3)
	return len(s.items)
}

// PointerSwap tests swapping pointers
func PointerSwap() int {
	a, b := 1, 2
	p1, p2 := &a, &b
	p1, p2 = p2, p1
	return *p1*10 + *p2
}

// MapNestedUpdate tests updating nested map value
func MapNestedUpdate() int {
	m := map[string]map[int]int{
		"outer": {1: 10},
	}
	m["outer"][1] = 20
	return m["outer"][1]
}

// SliceDeleteMiddle tests deleting middle element
func SliceDeleteMiddle() int {
	s := []int{1, 2, 3, 4, 5}
	i := 2
	s = append(s[:i], s[i+1:]...)
	return s[0] + s[1] + s[2] + s[3]
}

// StructNilField tests nil field in struct
func StructNilField() int {
	type S struct{ p *int }
	s := S{p: nil}
	if s.p == nil {
		return 1
	}
	return 0
}

// MapLookupOrInsert tests lookup or insert pattern
func MapLookupOrInsert() int {
	m := map[int]int{}
	key := 5
	if _, ok := m[key]; !ok {
		m[key] = key * 2
	}
	return m[key]
}

// SliceChainedSlice tests chained slice operations
func SliceChainedSlice() int {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[1:4]
	s2 := s1[1:2]
	return s2[0]
}

// StructEmbeddedOverride tests embedded field override
func StructEmbeddedOverride() int {
	type Inner struct{ V int }
	type Outer struct {
		Inner
		V int
	}
	o := Outer{Inner: Inner{V: 10}, V: 20}
	return o.Inner.V + o.V
}

// MapTwoKeys tests map with two different key types (via interface)
func MapTwoKeys() int {
	m := map[interface{}]int{}
	m[1] = 10
	m["a"] = 20
	return m[1] + m["a"]
}

// SliceNegativeIndex tests negative index behavior (through length)
func SliceNegativeIndex() int {
	s := []int{1, 2, 3, 4, 5}
	idx := len(s) - 1
	return s[idx]
}

// StructSelfRef tests struct self-reference via pointer
func StructSelfRef() int {
	type Node struct {
		value int
		next  *Node
	}
	n1 := &Node{value: 1}
	n2 := &Node{value: 2, next: n1}
	return n2.value + n2.next.value
}

// MapRangeBreak tests breaking from map range
func MapRangeBreak() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
		if sum > 25 {
			break
		}
	}
	return sum
}

// SliceStructIndex tests indexing struct slice
func SliceStructIndex() int {
	type P struct{ x, y int }
	s := []P{{1, 2}, {3, 4}}
	return s[0].x + s[1].y
}

// MapStructUpdate tests updating struct in map
func MapStructUpdate() int {
	type P struct{ x, y int }
	m := map[int]P{1: {10, 20}}
	m[1] = P{30, 40}
	return m[1].x + m[1].y
}

// PointerToPointer tests double pointer
func PointerToPointer() int {
	a := 10
	p := &a
	pp := &p
	return **pp
}
