package divergence_hunt274

import (
	"fmt"
)

// ============================================================================
// Round 274: Closure capture edge cases — loop variables, mutation, goroutines

// ClosureCaptureByRef tests closure capturing variable by reference
func ClosureCaptureByRef() string {
	x := 10
	f := func() int {
		x++
		return x
	}
	a := f()
	b := f()
	return fmt.Sprintf("a=%d,b=%d,x=%d", a, b, x)
}

// ClosureOverLoopVar tests closure capturing loop variable (pre-Go 1.22 semantics with explicit copy)
func ClosureOverLoopVar() string {
	results := []int{}
	for i := 0; i < 3; i++ {
		v := i // explicit copy
		results = append(results, v)
	}
	return fmt.Sprintf("results=%v", results)
}

// ClosureInMap tests closures stored in map
func ClosureInMap() string {
	m := map[string]func() int{}
	for i := 0; i < 3; i++ {
		v := i
		m[fmt.Sprintf("f%d", v)] = func() int { return v * 10 }
	}
	return fmt.Sprintf("f0=%d,f1=%d,f2=%d", m["f0"](), m["f1"](), m["f2"]())
}

// ClosureReturnsClosure tests closure that returns another closure
func ClosureReturnsClosure() string {
	adder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add5 := adder(5)
	return fmt.Sprintf("result=%d", add5(3))
}

// ClosureMultiCapture tests closure capturing multiple variables
func ClosureMultiCapture() string {
	a := 1
	b := 2
	c := 3
	f := func() int {
		return a + b + c
	}
	return fmt.Sprintf("sum=%d", f())
}

// ClosureReassignAndCapture tests closure captures latest value
func ClosureReassignAndCapture() string {
	x := 1
	f := func() int { return x }
	x = 2
	result := f()
	return fmt.Sprintf("result=%d", result)
}

// ClosureSliceOfFuncs tests slice of closures
func ClosureSliceOfFuncs() string {
	funcs := []func() int{}
	for i := 0; i < 4; i++ {
		v := i
		funcs = append(funcs, func() int { return v * v })
	}
	result := ""
	for _, f := range funcs {
		result += fmt.Sprintf("%d,", f())
	}
	return result[:len(result)-1]
}

// ClosureModifyOuter tests closure modifying outer variable
func ClosureModifyOuter() string {
	sum := 0
	nums := []int{1, 2, 3, 4, 5}
	for _, n := range nums {
		func() {
			sum += n
		}()
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ClosureAsArg tests passing closure as argument
func ClosureAsArg() string {
	apply := func(f func(int) int, x int) int {
		return f(x)
	}
	double := func(x int) int { return x * 2 }
	return fmt.Sprintf("result=%d", apply(double, 5))
}

// ClosureRecursive tests closure that calls itself (via variable)
func ClosureRecursive() string {
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	return fmt.Sprintf("fib10=%d", fib(10))
}

// ClosureOverShadowed tests closure over shadowed variable
func ClosureOverShadowed() string {
	x := 10
	f := func() int {
		x := 20 // shadows outer x
		return x
	}
	return fmt.Sprintf("f=%d,x=%d", f(), x)
}

// ClosureImmediateInvoke tests immediately invoked closure
func ClosureImmediateInvoke() string {
	result := func(x int) int {
		return x * x
	}(7)
	return fmt.Sprintf("result=%d", result)
}
