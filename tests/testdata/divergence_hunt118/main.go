package divergence_hunt118

import (
	"fmt"
	"strings"
)

// ============================================================================
// Round 118: Func type as variable/parameter/return
// ============================================================================

func FuncVariable() string {
	fn := func() string { return "hello" }
	return fn()
}

func FuncParam() string {
	apply := func(fn func(int) int, x int) int {
		return fn(x)
	}
	double := func(n int) int { return n * 2 }
	return fmt.Sprintf("%d", apply(double, 5))
}

func FuncReturn() string {
	makeAdder := func(n int) func(int) int {
		return func(x int) int { return x + n }
	}
	add5 := makeAdder(5)
	return fmt.Sprintf("%d", add5(10))
}

func FuncSlice() string {
	fns := []func(int) int{
		func(x int) int { return x + 1 },
		func(x int) int { return x * 2 },
	}
	return fmt.Sprintf("%d:%d", fns[0](5), fns[1](5))
}

func FuncMap() string {
	ops := map[string]func(int, int) int{
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	}
	return fmt.Sprintf("%d:%d", ops["add"](3, 4), ops["mul"](3, 4))
}

func FuncChaining() string {
	process := func(s string) string {
		return strings.TrimSpace(strings.ToUpper(s))
	}
	return process("  hello  ")
}

func FuncCompose() string {
	compose := func(f, g func(int) int) func(int) int {
		return func(x int) int { return f(g(x)) }
	}
	double := func(x int) int { return x * 2 }
	addOne := func(x int) int { return x + 1 }
	doubleThenAdd := compose(addOne, double)
	return fmt.Sprintf("%d", doubleThenAdd(5))
}

func FuncAsField() string {
	type Processor struct {
		Transform func(string) string
	}
	p := Processor{
		Transform: strings.ToUpper,
	}
	return p.Transform("hello")
}

func FuncComparison() string {
	// Functions can only be compared to nil
	fn := func() {}
	return fmt.Sprintf("%v", fn != nil)
}

func FuncNilCheck() string {
	var fn func()
	return fmt.Sprintf("%v", fn == nil)
}
