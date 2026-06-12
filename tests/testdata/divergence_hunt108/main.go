package divergence_hunt108

import "fmt"

// ============================================================================
// Round 108: Closure capturing and mutation
// ============================================================================

func ClosureCounter() int {
	makeCounter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}
	c := makeCounter()
	c()
	c()
	return c()
}

func ClosureCapture() string {
	x := "hello"
	fn := func() string {
		return x
	}
	x = "world"
	return fn()
}

func ClosureMultiCapture() string {
	a, b := 10, 20
	fn := func() string {
		return fmt.Sprintf("%d:%d", a, b)
	}
	a = 30
	return fn()
}

func ClosureInLoop() string {
	var fns []func() string
	for i := 0; i < 3; i++ {
		i := i // capture loop var
		fns = append(fns, func() string {
			return fmt.Sprintf("%d", i)
		})
	}
	result := ""
	for _, fn := range fns {
		result += fn()
	}
	return result
}

func ClosureModifyOuter() int {
	x := 10
	fn := func() {
		x = 20
	}
	fn()
	return x
}

func ClosureReturnClosure() string {
	adder := func(base int) func(int) int {
		return func(n int) int {
			return base + n
		}
	}
	add5 := adder(5)
	return fmt.Sprintf("%d", add5(10))
}

func ClosureSlice() string {
	var fns []func() int
	for i := 0; i < 3; i++ {
		v := i * 10
		fns = append(fns, func() int { return v })
	}
	result := 0
	for _, fn := range fns {
		result += fn()
	}
	return fmt.Sprintf("%d", result)
}

func ClosureAsParam() string {
	apply := func(fn func(string) string, s string) string {
		return fn(s)
	}
	upper := func(s string) string {
		return fmt.Sprintf(">%s<", s)
	}
	return apply(upper, "hi")
}

func ClosureCaptureSlice() string {
	data := []int{1, 2, 3}
	fn := func() int {
		sum := 0
		for _, v := range data {
			sum += v
		}
		return sum
	}
	data = append(data, 4)
	return fmt.Sprintf("%d", fn())
}

func ClosureStacked() int {
	outer := func() func() int {
		x := 1
		middle := func() func() int {
			x *= 2
			return func() int {
				x *= 3
				return x
			}
		}
		inner := middle()
		return inner
	}
	fn := outer()
	return fn()
}
