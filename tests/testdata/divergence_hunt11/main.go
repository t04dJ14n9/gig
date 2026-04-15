package divergence_hunt11

import "fmt"

// ============================================================================
// Round 11: More complex closure patterns, defer edge cases, recover edge
// cases, panic value types, named return edge cases
// ============================================================================

// DeferInLoop tests defer inside a loop (should still work)
func DeferInLoop() (result int) {
	for i := 0; i < 5; i++ {
		v := i
		defer func() { result += v }()
	}
	return 0
}

// DeferAndPanicOrder tests that defers run even after panic
func DeferAndPanicOrder() (result int) {
	defer func() { result += 1 }()
	defer func() { result += 10 }()
	defer func() { recover() }()
	panic("test")
}

// RecoverInFunction tests recover() called in a regular function (not defer)
func RecoverInFunction() int {
	// recover() returns nil when not called directly from a deferred function
	if r := recover(); r != nil {
		return -1
	}
	return 42
}

// PanicWithStruct tests panic with struct value
func PanicWithStruct() (result string) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				result = s
			}
		}
	}()
	panic("struct panic")
}

// NamedReturnWithDefer tests named return with defer modification
func NamedReturnWithDefer() (result int) {
	defer func() { result++ }()
	return 10
}

// MultipleDeferModify tests multiple defers modifying return value
func MultipleDeferModify() (result int) {
	defer func() { result += 1 }()
	defer func() { result += 10 }()
	defer func() { result += 100 }()
	return 0
}

// DeferWithArgument tests defer capturing argument
func DeferWithArgument() (result int) {
	process := func(x int) int { return x * 2 }
	defer func() { result = process(5) }()
	return 0
}

// PanicNilValue tests panic with nil value
func PanicNilValue() (result int) {
	defer func() {
		if r := recover(); r == nil {
			result = -1
		}
	}()
	panic(nil)
}

// ClosureReturnFunc tests closure that returns a function
func ClosureReturnFunc() int {
	makeCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}
	c := makeCounter()
	c()
	c()
	return c()
}

// FmtSprintfMulti tests fmt.Sprintf with multiple args
func FmtSprintfMulti() string {
	return fmt.Sprintf("%d + %d = %d", 1, 2, 3)
}

// FmtErrorf tests fmt.Errorf
func FmtErrorf() string {
	return fmt.Errorf("err %d", 42).Error()
}

// NestedDeferRecover tests nested defer/recover
func NestedDeferRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	func() {
		defer func() {
			if r := recover(); r != nil {
				result += r.(int)
			}
		}()
		panic(10)
	}()
	panic(20)
}

// DeferWithMethod tests defer calling method-like function
func DeferWithMethod() (result int) {
	type Obj struct{ val int }
	o := &Obj{val: 0}
	increment := func(obj *Obj) { obj.val++ }
	defer increment(o)
	o.val = 10
	return o.val
}

// ClosureCaptureSlice tests closure capturing slice
func ClosureCaptureSlice() int {
	s := []int{1, 2, 3}
	modify := func() { s[0] = 99 }
	modify()
	return s[0]
}

// ClosureCaptureMap tests closure capturing map
func ClosureCaptureMap() int {
	m := map[string]int{"a": 1}
	setValue := func() { m["a"] = 42 }
	setValue()
	return m["a"]
}

// MultiplePanicRecover tests multiple panic/recover cycles
func MultiplePanicRecover() (result int) {
	for i := 0; i < 3; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					result += r.(int)
				}
			}()
			panic(i + 1)
		}()
	}
	return result
}

// DeferRecoverReturnsValue tests recover in defer that returns value
func DeferRecoverReturnsValue() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = 100
		}
	}()
	panic("boom")
}

// SliceAppendInClosure tests slice append inside closure
func SliceAppendInClosure() int {
	s := []int{}
	appendFn := func(v int) { s = append(s, v) }
	appendFn(1)
	appendFn(2)
	appendFn(3)
	return len(s)
}

// MapWriteInClosure tests map write inside closure
func MapWriteInClosure() int {
	m := map[string]int{}
	writeFn := func(k string, v int) { m[k] = v }
	writeFn("x", 10)
	writeFn("y", 20)
	return m["x"] + m["y"]
}

// DeferChain tests chain of defers
func DeferChain() (result int) {
	for i := 0; i < 5; i++ {
		defer func(v int) { result = result*10 + v }(i)
	}
	return 0
}

// RecoverReturnsNilAfter tests that recover() returns nil after first call
func RecoverReturnsNilAfter() (result int) {
	defer func() {
		r1 := recover()
		r2 := recover()
		if r1 != nil && r2 == nil {
			result = 1
		}
	}()
	panic("test")
}
