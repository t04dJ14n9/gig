package strange_syntax_panic

// ============================================================================
// PANIC/RECOVER EDGE CASES
// ============================================================================

// MultipleDefersWithRecover tests multiple defers with recover
func MultipleDefersWithRecover() int {
	result := 0
	defer func() {
		if r := recover(); r != nil {
			result += 100
		}
	}()
	defer func() { result++ }()
	defer func() { result += 2 }()
	panic("test")
}

// PanicInDefer tests panic in defer
func PanicInDefer() int {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic in defer
		}
	}()
	defer func() {
		panic("defer panic")
	}()
	return 1
}

// NestedPanics tests nested panics
func NestedPanics() int {
	result := 0
	defer func() {
		if r := recover(); r != nil {
			result++
			defer func() {
				if r2 := recover(); r2 != nil {
					result++
				}
			}()
			panic("second panic")
		}
	}()
	panic("first panic")
}

// ClosureWithDefer tests closure with defer
func ClosureWithDefer() int {
	fn := func() int {
		result := 1
		defer func() { result *= 2 }()
		return result // returns 1, then defer makes it 2
	}
	return fn()
}

// DeferInClosure tests defer inside closure
func DeferInClosure() int {
	result := 0
	fn := func() {
		defer func() { result++ }()
		result = 10
	}
	fn()
	return result // 10 + 1 = 11
}

// DeferWithPanicAndRecover tests defer with panic and recover
func DeferWithPanicAndRecover() int {
	defer func() {
		if r := recover(); r != nil {
			// recovered
		}
	}()
	panic("test panic")
}

// MultipleDeferRecover tests multiple defer with recover
func MultipleDeferRecover() int {
	result := 0
	defer func() {
		result++ // third
	}()
	defer func() {
		if r := recover(); r != nil {
			result += 10 // second
		}
	}()
	defer func() {
		result++ // first
	}()
	panic("test")
}

// ClosureWithPanicAndRecover tests closure with panic and recover
func ClosureWithPanicAndRecover() int {
	fn := func() (result int) {
		defer func() {
			if r := recover(); r != nil {
				result = -1
			}
		}()
		panic("test")
	}
	return fn()
}
