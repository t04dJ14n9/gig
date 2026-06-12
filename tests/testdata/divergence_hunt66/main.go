package divergence_hunt66

// ============================================================================
// Round 66: Panic/recover nested - recover in nested defer, panic after recover
// ============================================================================

func RecoverInNestedDefer() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	defer func() {
		// This runs first (LIFO), before the recover defer
	}()
	panic(42)
}

func PanicAfterRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	defer func() {
		defer func() {
			recover() // recover the second panic
		}()
		panic(100) // second panic
	}()
	panic(42) // first panic
}

func RecoverReturnsNilAfterCall() (result int) {
	defer func() {
		r := recover()
		if r == nil {
			result = 1
		} else {
			result = 2
		}
	}()
	// no panic, recover returns nil
	return 0
}

func MultipleRecoverSameDefer() (result int) {
	defer func() {
		r1 := recover()
		r2 := recover() // should also return the panic value
		if r1 != nil && r2 != nil {
			result = r1.(int) + r2.(int)
		} else if r1 != nil {
			result = r1.(int)
		} else {
			result = 0
		}
	}()
	panic(21)
}

func RecoverOnlyInDefer() (result int) {
	defer func() {
		if recover() != nil {
			result = 1
		}
	}()
	panic("boom")
}

func NestedPanicRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int) * 10
		}
	}()
	inner := func() {
		defer func() {
			recover()
		}()
		panic(1)
	}
	inner()
	panic(2)
}

func PanicString() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(string)
		}
	}()
	panic("hello")
}

func PanicNilInterface() (result string) {
	defer func() {
		if r := recover(); r != nil {
			if r == nil {
				result = "nil"
			} else {
				result = "not nil"
			}
		}
	}()
	var err error
	panic(err)
}

func DeferPanicOrder() (result int) {
	defer func() { result += 1 }()
	defer func() {
		result += 10
		recover()
	}()
	defer func() { result += 100 }()
	panic("x")
}

func RecoverTypeAssertion() (result int) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case int:
				result = v
			case string:
				result = -1
			default:
				result = -2
			}
		}
	}()
	panic(42)
}
