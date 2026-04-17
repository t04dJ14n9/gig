package divergence_hunt95

import "fmt"

// ============================================================================
// Round 95: Defer edge cases - argument evaluation, stacked defers, defer+return
// ============================================================================

func DeferArgEval() int {
	x := 10
	defer func() {
		x = 999
	}()
	return x
}

func DeferArgCapture() string {
	x := "hello"
	defer func() {
		x = "deferred"
	}()
	return x
}

func DeferModifyReturn() int {
	// Named return - defer can modify it
	result := 0
	defer func() {
		result++
	}()
	result = 42
	return result
}

func StackedDefers() string {
	order := ""
	defer func() { order += "1" }()
	defer func() { order += "2" }()
	defer func() { order += "3" }()
	return order
}

func DeferInLoop() int {
	sum := 0
	for i := 0; i < 5; i++ {
		defer func(v int) {
			sum += v
		}(i)
	}
	return sum
}

func DeferClosureCapture() string {
	result := "start"
	defer func() {
		result += ":defer"
	}()
	result += ":middle"
	return result
}

func DeferWithRecover() string {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("recovered: %v", r)
		}
	}()
	panic("test")
}

func MultipleDefersOrder() string {
	var stack []int
	defer func() { stack = append(stack, 1) }()
	defer func() { stack = append(stack, 2) }()
	defer func() { stack = append(stack, 3) }()
	_ = stack
	return "3,2,1"
}

func DeferReturnOrder() string {
	result := "init"
	defer func() {
		result += ":defer"
	}()
	return result + ":return"
}

func DeferWithMethod() string {
	type S struct{ val string }
	s := S{val: "hello"}
	defer func() {
		s.val += ":deferred"
	}()
	return s.val
}

func DeferClosureArgVsCapture() string {
	x := "original"
	defer func(v string) {
		_ = v // captured at defer time
	}(x)
	x = "modified"
	return x
}
