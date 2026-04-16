package divergence_hunt45

// ============================================================================
// Round 45: Defer edge cases - named return modification, closure capture,
// stack ordering, defer in loops
// ============================================================================

func DeferNamedReturn() (result int) {
	defer func() { result++ }()
	return 10
}

func DeferMultipleNamedReturn() (result int) {
	defer func() { result += 1 }()
	defer func() { result += 10 }()
	defer func() { result += 100 }()
	return 0
}

func DeferCaptureByValue() (result int) {
	x := 10
	defer func(v int) { result = v }(x)
	x = 20
	return 0
}

func DeferCaptureByRef() (result int) {
	x := 10
	defer func() { result = x }()
	x = 20
	return 0
}

func DeferInLoop() (result int) {
	for i := 0; i < 5; i++ {
		v := i
		defer func() { result += v }()
	}
	return 0
}

func DeferOrder() (result int) {
	defer func() { result = result*10 + 1 }()
	defer func() { result = result*10 + 2 }()
	defer func() { result = result*10 + 3 }()
	return 0
}

func DeferModifyBeforeReturn() (result int) {
	result = 5
	defer func() { result *= 2 }()
	return result // defer runs after, modifies to 10
}

func DeferWithRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	panic(42)
}

func DeferAfterRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = 10
		}
	}()
	defer func() { result += 5 }() // runs before recover defer
	panic("error")
}

func DeferWithNilRecover() (result int) {
	defer func() {
		r := recover()
		if r == nil {
			result = 1
		} else {
			result = 2
		}
	}()
	panic(nil)
}

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

func DeferExternalFunc() (result int) {
	x := 0
	increment := func() { x++ }
	defer increment()
	defer increment()
	defer increment()
	return x // defers run after, x becomes 3
}

func DeferReturnOverride() (result int) {
	defer func() { result = 99 }()
	return 1
}

func DeferConditional() (result int) {
	condition := true
	if condition {
		defer func() { result += 10 }()
	}
	return 5
}

func FmtDeferCapture() string {
	x := "hello"
	defer func() { _ = x }()
	return x
}
