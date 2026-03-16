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
