package divergence_hunt266

import (
	"fmt"
)

// ============================================================================
// Round 266: Defer edge cases — modify returns, multiple defers, closures
// ============================================================================

// DeferModifyNamedReturn tests defer modifying named return value
func DeferModifyNamedReturn() (result string) {
	result = "before"
	defer func() {
		result = "after"
	}()
	return result
}

// DeferOrder tests LIFO defer execution order
func DeferOrder() string {
	result := ""
	defer func() { result += "A" }()
	defer func() { result += "B" }()
	defer func() { result += "C" }()
	return result // returns before defers run, but defers modify
}

// DeferInLoop tests defer in a loop (common mistake pattern)
func DeferInLoop() string {
	result := ""
	for i := 0; i < 3; i++ {
		v := i
		defer func() {
			result += fmt.Sprintf("%d", v)
		}()
	}
	return result // defers haven't run yet
}

// DeferReadModified tests reading named return in defer after modification
func DeferReadModified() (x int) {
	x = 10
	defer func() {
		x += 5
	}()
	return x
}

// DeferMultipleModify tests multiple defers modifying return
func DeferMultipleModify() (x int) {
	x = 1
	defer func() { x *= 2 }()
	defer func() { x += 10 }()
	return x
}

// DeferClosureCapture tests defer closure capturing variable
func DeferClosureCapture() string {
	v := "original"
	defer func() {
		v = "modified"
	}()
	return v
}

// DeferArgsEvaluatedNow tests defer evaluating args at defer time
func DeferArgsEvaluatedNow() string {
	x := 10
	result := ""
	defer func(val int) {
		result = fmt.Sprintf("val=%d", val)
	}(x)
	x = 20
	return result // "val=10" because args evaluated at defer call
}

// DeferRecoverPanic tests defer recovering from panic
func DeferRecoverPanic() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered:%v", r)
		}
	}()
	panic("boom")
}

// DeferWithoutRecover tests normal defer without panic
func DeferWithoutRecover() string {
	result := "start"
	defer func() {
		result = "deferred"
	}()
	return result
}

// DeferPanicInDefer tests panic inside a defer
func DeferPanicInDefer() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("outer_recover:%v", r)
		}
	}()
	defer func() {
		panic("inner_panic")
	}()
	return "normal"
}
