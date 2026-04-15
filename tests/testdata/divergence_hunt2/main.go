package divergence_hunt2

// ============================================================================
// Round 2: Map operations, slice tricks, string edge cases, integer conversion,
// multi-return, blank identifier, nil handling, pointer operations
// ============================================================================

// MapLen tests map length
func MapLen() int { m := map[string]int{"a": 1, "b": 2}; return len(m) }

// MapDelete tests map delete
func MapDelete() int {
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	return len(m)
}

// MapOverwrite tests map overwrite
func MapOverwrite() int {
	m := map[string]int{"a": 1}
	m["a"] = 2
	return m["a"]
}

// SliceNilAppend tests appending to nil slice
func SliceNilAppend() int {
	var s []int
	s = append(s, 1, 2, 3)
	return len(s)
}

// SliceGrow tests slice grow via append
func SliceGrow() int {
	s := make([]int, 2)
	s = append(s, 3)
	return len(s) + s[2]
}

// StringLen tests string length
func StringLen() int { return len("hello世界") }

// StringConcat tests string concatenation
func StringConcat() string { return "hello" + " " + "world" }

// IntConversion tests various int conversions
func IntConversion() int64 {
	var x int = 42
	var y int64 = int64(x)
	return y
}

// UintConversion tests uint conversions
func UintConversion() uint64 {
	var x uint = 42
	var y uint64 = uint64(x)
	return y
}

// MultiReturnSwap tests multiple return value swap
func MultiReturnSwap() int {
	swap := func(a, b int) (int, int) { return b, a }
	x, y := swap(1, 2)
	return x*10 + y
}

// BlankIdentifier tests blank identifier
func BlankIdentifier() int {
	_, b := 1, 2
	return b
}

// NilSliceLen tests nil slice len
func NilSliceLen() int { var s []int; return len(s) }

// NilMapLen tests nil map len
func NilMapLen() int { var m map[string]int; return len(m) }

// PointerDeref tests pointer dereference
func PointerDeref() int {
	x := 42
	p := &x
	return *p
}

// PointerAssign tests pointer assignment
func PointerAssign() int {
	x := 10
	p := &x
	*p = 20
	return x
}

// SliceOfPointers tests slice of pointers
func SliceOfPointers() int {
	a, b := 1, 2
	s := []*int{&a, &b}
	return *s[0] + *s[1]
}

// MapIteration tests map iteration (just count)
func MapIteration() int {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	count := 0
	for range m {
		count++
	}
	return count
}

// StringRange tests ranging over string
func StringRange() int {
	s := "abc"
	count := 0
	for i, ch := range s {
		count += i
		_ = ch
	}
	return count
}

// FloatConversion tests float to int conversion
func FloatConversion() int {
	x := 3.7
	return int(x)
}

// ByteSliceAppend tests byte slice append
func ByteSliceAppend() int {
	b := []byte{1, 2}
	b = append(b, 3)
	return int(b[0]) + int(b[1]) + int(b[2])
}

// ByteSliceWrite tests byte slice write via index
func ByteSliceWrite() int {
	b := make([]byte, 3)
	b[0] = 10
	b[1] = 20
	b[2] = 30
	return int(b[0]) + int(b[1]) + int(b[2])
}

// StructCompare tests struct comparison
func StructCompare() bool {
	type P struct{ X, Y int }
	return P{1, 2} == P{1, 2}
}

// ArrayLen tests array length
func ArrayLen() int {
	a := [3]int{1, 2, 3}
	return len(a)
}

// ArrayValue tests array value semantics
func ArrayValue() int {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 99
	return a[0]
}

// StringIndexOutOfRange tests string index with recover
func StringIndexOutOfRange() (result int) {
	defer func() { if recover() != nil { result = -1 } }()
	s := "hi"
	_ = s[10]
	return 0
}

// MapKeyIntFloat tests map with different key types
func MapKeyIntFloat() int {
	m1 := map[int]string{1: "a"}
	m2 := map[float64]string{1.0: "b"}
	return len(m1) + len(m2)
}

// ShortVarDecl tests short variable declaration in if
func ShortVarDecl() int {
	if x := 42; x > 0 {
		return x
	}
	return 0
}

// MultipleShortVar tests multiple short var decl
func MultipleShortVar() int {
	x, y, z := 1, 2, 3
	return x + y + z
}

// SliceThreeIndex tests three-index slice
func SliceThreeIndex() int {
	s := []int{0, 1, 2, 3, 4}
	sub := s[1:3:3]
	return len(sub) + cap(sub)
}

// NilFuncCall tests nil function pointer with recover
func NilFuncCall() (result int) {
	defer func() { if recover() != nil { result = -1 } }()
	var f func()
	f()
	return 0
}

// StringByteSliceConversion tests string ↔ []byte conversion
func StringByteSliceConversion() int {
	s := "hello"
	b := []byte(s)
	return len(b)
}
