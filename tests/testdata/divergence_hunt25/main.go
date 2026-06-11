package divergence_hunt25

import "fmt"

// ============================================================================
// Round 25: Complex closure, defer, and recover patterns
// ============================================================================

func DeferStack() (result int) {
	for i := 1; i <= 5; i++ {
		defer func(n int) { result += n }(i)
	}
	return 0
}

func DeferInClosure() (result int) {
	fn := func() {
		defer func() { result++ }()
		result += 10
	}
	fn()
	return result
}

func RecoverInNestedDefer() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	defer func() {
		panic(42)
	}()
	return 0
}

func MultipleRecover() (result int) {
	defer func() {
		r1 := recover()
		r2 := recover()
		if r1 != nil && r2 == nil {
			result = 1
		}
	}()
	panic("test")
}

func DeferClosureCapture() (result int) {
	x := 10
	defer func() { result = x }()
	x = 20
	return 0
}

func DeferClosureCopy() (result int) {
	x := 10
	defer func(v int) { result = v }(x)
	x = 20
	return 0
}

func PanicInDeferRecover() (result int) {
	defer func() {
		if r := recover(); r != nil { result = -1 }
	}()
	defer func() { panic("defer panic") }()
	return 42
}

func DeferModifyNamedReturn() (result int) {
	defer func() { result *= 2 }()
	return 5
}

func NestedPanicRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(int); ok { result = v }
		}
	}()
	func() {
		defer func() { recover() }()
		panic(100)
	}()
	panic(200)
}

func ClosureWithDefer() (result int) {
	fn := func() (r int) {
		defer func() { r++ }()
		return 10
	}
	return fn()
}

func RecursiveWithDefer() (result int) {
	var fib func(n int) (r int)
	fib = func(n int) (r int) {
		defer func() { r++ }()
		if n <= 1 { return n }
		return fib(n-1) + fib(n-2)
	}
	return fib(5)
}

func PanicRecoverTypeSwitch() (result string) {
	defer func() {
		switch r := recover().(type) {
		case int: result = fmt.Sprintf("int:%d", r)
		case string: result = fmt.Sprintf("string:%s", r)
		default: result = "unknown"
		}
	}()
	panic(42)
}

func DeferMultipleModifies() (result int) {
	defer func() { result += 1 }()
	defer func() { result += 10 }()
	defer func() { result += 100 }()
	return 0
}

func RecoverReturnsPanicValue() (result int) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(int); ok { result = v }
		}
	}()
	panic(99)
}

func DeferInMethod() (result int) {
	type S struct{ val int }
	s := &S{val: 10}
	func() {
		defer func() { s.val++ }()
		s.val += 5
	}()
	return s.val
}

func ClosureState() int {
	x := 0
	inc := func() int { x++; return x }
	inc()
	inc()
	return inc()
}

func ClosureSharedState() int {
	x := 0
	inc := func() { x++ }
	get := func() int { return x }
	inc()
	inc()
	inc()
	return get()
}

func FmtDefer() string {
	var result string
	defer func() { result = "done" }()
	_ = result
	return "done"
}
