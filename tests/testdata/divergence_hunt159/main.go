package divergence_hunt159

import "fmt"

// ============================================================================
// Round 159: Defer and panic edge cases
// ============================================================================

// NamedReturnWithDefer tests named return modification via defer
func NamedReturnWithDefer() (result int) {
	defer func() {
		result = 42
	}()
	return 0
}

// NamedReturnWithDeferAndValue tests named return with value
func NamedReturnWithDeferAndValue() (result int) {
	defer func() {
		result *= 2
	}()
	return 21
}

// MultipleDefersExecutionOrder tests multiple defers execute in LIFO order
func MultipleDefersExecutionOrder() string {
	result := ""
	defer func() { result += "3" }()
	defer func() { result += "2" }()
	defer func() { result += "1" }()
	result = "start-"
	return result
}

// DeferWithArguments tests defer captures arguments at defer time
func DeferWithArguments() string {
	result := ""
	x := 1
	defer func(n int) {
		result = fmt.Sprintf("deferred-%d", n)
	}(x)
	x = 100
	return result
}

// DeferInLoopLastValue tests defer in loop capturing loop variable
func DeferInLoopLastValue() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered-%v", r)
		}
	}()

	// This pattern is tricky - defers in loops
	for i := 0; i < 3; i++ {
		if i == 2 {
			panic("loop-panic")
		}
	}
	return "no-panic"
}

// PanicWithString tests panic with string value
func PanicWithString() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("caught-%v", r)
		}
	}()
	panic("string-panic")
}

// PanicWithError tests panic with error value
func PanicWithError() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("caught-%v", r)
		}
	}()
	panic(fmt.Errorf("error-panic"))
}

// PanicWithInt tests panic with int value
func PanicWithInt() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("caught-%v", r)
		}
	}()
	panic(42)
}

// NestedDeferPanic tests panic in deferred function
func NestedDeferPanic() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("outer-recovered-%v", r)
		}
	}()
	defer func() {
		panic("inner-panic")
	}()
	return "no-panic"
}

// RecoverOnlyInDefer tests that recover only works in defer
func RecoverOnlyInDefer() string {
	// This recover won't catch anything since it's not in a deferred function
	if r := recover(); r != nil {
		return fmt.Sprintf("recovered-%v", r)
	}
	return "no-recover"
}

// PanicNilInterface tests panic with nil interface
func PanicNilInterface() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("caught-nil=%t", r == nil)
		} else {
			result = "recovered-nil"
		}
	}()
	var i interface{} = nil
	panic(i)
}

// DeferReturnValue tests defer modifying return value after return
func DeferReturnValue() (x int) {
	defer func() {
		x = 100
	}()
	return 10
}
