package divergence_hunt210

import "fmt"

// ============================================================================
// Round 210: Closure variable capture edge cases
// ============================================================================

func LoopCaptureClassic() string {
	funcs := []func() int{}
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int { return i })
	}
	result := ""
	for _, f := range funcs {
		result += fmt.Sprintf("%d", f())
	}
	return result
}

func LoopCaptureFixed() string {
	funcs := []func() int{}
	for i := 0; i < 3; i++ {
		i := i
		funcs = append(funcs, func() int { return i })
	}
	result := ""
	for _, f := range funcs {
		result += fmt.Sprintf("%d", f())
	}
	return result
}

func ClosureModifiesCaptured() string {
	x := 0
	f := func() {
		x++
	}
	f()
	f()
	f()
	return fmt.Sprintf("%d", x)
}

func MultipleClosuresShareVar() string {
	x := 10
	inc := func() { x++ }
	dec := func() { x-- }
	inc()
	inc()
	dec()
	return fmt.Sprintf("%d", x)
}

func ClosureCapturesSlice() string {
	s := []int{1, 2, 3}
	f := func() int {
		s[0] = 99
		return len(s)
	}
	f()
	return fmt.Sprintf("%v", s)
}

func ClosureCapturesMap() string {
	m := map[string]int{"a": 1}
	f := func() {
		m["b"] = 2
	}
	f()
	return fmt.Sprintf("%d", len(m))
}

func NestedClosure() string {
	x := 5
	outer := func() func() int {
		y := 10
		return func() int { return x + y }
	}
	inner := outer()
	x = 100
	return fmt.Sprintf("%d", inner())
}

func ClosureAsReturn() string {
	makeAdder := func(n int) func(int) int {
		return func(x int) int { return x + n }
	}
	add5 := makeAdder(5)
	add10 := makeAdder(10)
	return fmt.Sprintf("%d:%d", add5(3), add10(3))
}

func ClosureReceivesItself() string {
	var fib func(n int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fmt.Sprintf("%d", fib(6))
}

func DeferredClosureCapture() string {
	result := ""
	for i := 0; i < 3; i++ {
		defer func() {
			result += fmt.Sprintf("%d", i)
		}()
	}
	return result
}

func RangeCaptureIssue() string {
	funcs := []func() int{}
	items := []int{10, 20, 30}
	for _, v := range items {
		funcs = append(funcs, func() int { return v })
	}
	result := ""
	for _, f := range funcs {
		result += fmt.Sprintf("%d", f())
	}
	return result
}
