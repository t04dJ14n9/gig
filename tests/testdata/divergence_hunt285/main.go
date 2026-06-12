package divergence_hunt285

import (
	"fmt"
)

// ============================================================================
// Round 285: Defer in tricky contexts — defer in method, defer return values, defer in closure

type DeferTracker struct {
	logs []string
}

func (d *DeferTracker) Log(msg string) {
	d.logs = append(d.logs, msg)
}

func (d *DeferTracker) String() string {
	result := ""
	for i, l := range d.logs {
		if i > 0 {
			result += ","
		}
		result += l
	}
	return result
}

// DeferInMethod tests defer inside a method
func DeferInMethod() string {
	d := &DeferTracker{}
	d.Log("start")
	defer d.Log("defer1")
	d.Log("end")
	return d.String()
}

// DeferNamedReturnOrder tests that named return is set before defers run
func DeferNamedReturnOrder() (result string) {
	result = "initial"
	defer func() {
		result = "deferred:" + result
	}()
	result = "final"
	return
}

// DeferMultipleModifyOrder tests multiple defers modify in LIFO order
func DeferMultipleModifyOrder() (result string) {
	result = "A"
	defer func() { result += "1" }()
	defer func() { result += "2" }()
	defer func() { result += "3" }()
	result = "B"
	return
}

// DeferInClosure tests defer called inside a closure
func DeferInClosure() string {
	result := ""
	func() {
		defer func() { result += "deferred" }()
		result += "inner"
	}()
	result += "outer"
	return result
}

// DeferArgEvalAtDeferTime tests defer args evaluated when defer is called
func DeferArgEvalAtDeferTime() string {
	x := 10
	result := ""
	defer func(val int) {
		result = fmt.Sprintf("val=%d", val)
	}(x)
	x = 20
	return result
}

// DeferInLoopAccumulate tests defer accumulating in loop
func DeferInLoopAccumulate() (result string) {
	for i := 0; i < 3; i++ {
		v := i
		defer func() {
			result += fmt.Sprintf("%d", v)
		}()
	}
	return
}

// DeferRecoverAndContinue tests recover returns to calling function
func DeferRecoverAndContinue() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered:%v", r)
		}
	}()
	panic("test_panic")
}

// DeferModifySlice tests defer modifying a slice (reference type)
func DeferModifySlice() string {
	s := []int{1, 2, 3}
	defer func() {
		s[0] = 99
	}()
	_ = s                                   // suppress unused warning
	return fmt.Sprintf("after_defer=%v", s) // note: return value captured before defers run
}

// DeferReturnsFirstThenDefers tests return value vs defer modification
func DeferReturnsFirstThenDefers() (result string) {
	result = "return_val"
	defer func() {
		result = "defer_modified"
	}()
	return result
}

// DeferInIfBranch tests defer in conditional branch
func DeferInIfBranch() string {
	result := "start"
	x := 5
	if x > 3 {
		defer func() { result += ":defer_if" }()
	}
	result += ":after_if"
	return result
}
