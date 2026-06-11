package divergence_hunt43

import (
	"fmt"
	"strings"
)

// ============================================================================
// Round 43: Closure patterns - capturing, modifying, returning, currying
// ============================================================================

func ClosureCaptureValue() int {
	x := 10
	fn := func() int { return x }
	x = 20
	return fn() // captures reference, returns 20
}

func ClosureCapturePointer() int {
	x := 10
	p := &x
	fn := func() int { return *p }
	x = 20
	return fn() // 20
}

func ClosureModifyCaptured() int {
	x := 10
	fn := func() { x = 42 }
	fn()
	return x
}

func ClosureReturnFunc() int {
	makeAdder := func(n int) func(int) int {
		return func(x int) int { return x + n }
	}
	add5 := makeAdder(5)
	return add5(10)
}

func ClosureCurry() int {
	curry := func(a int) func(int) func(int) int {
		return func(b int) func(int) int {
			return func(c int) int { return a + b + c }
		}
	}
	return curry(1)(2)(3)
}

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

func ClosureAccumulator() int {
	makeAccum := func() func(int) int {
		total := 0
		return func(n int) int {
			total += n
			return total
		}
	}
	acc := makeAccum()
	acc(10)
	acc(20)
	return acc(5)
}

func ClosureOverSlice() int {
	s := []int{1, 2, 3}
	fn := func() { s[0] = 99 }
	fn()
	return s[0]
}

func ClosureOverMap() int {
	m := map[string]int{"a": 1}
	fn := func() { m["a"] = 42 }
	fn()
	return m["a"]
}

func ClosureOverLoopCopy() int {
	fns := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		v := i // copy
		fns[i] = func() int { return v }
	}
	return fns[0]() + fns[1]() + fns[2]()
}

func ClosureOverLoopNoCopy() int {
	fns := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		fns[i] = func() int { return i } // captures loop variable
	}
	// In Go 1.22+, i is per-iteration, but we test the VM behavior
	return fns[0]() + fns[1]() + fns[2]()
}

func ClosurePartialApplication() int {
	multiply := func(a, b int) int { return a * b }
	double := func(x int) int { return multiply(2, x) }
	triple := func(x int) int { return multiply(3, x) }
	return double(5) + triple(5)
}

func ClosureFilter() int {
	filter := func(s []int, pred func(int) bool) []int {
		result := []int{}
		for _, v := range s {
			if pred(v) { result = append(result, v) }
		}
		return result
	}
	data := []int{1, 2, 3, 4, 5, 6}
	evens := filter(data, func(n int) bool { return n%2 == 0 })
	return len(evens)
}

func ClosureMapFunc() int {
	mapFn := func(s []int, f func(int) int) []int {
		result := make([]int, len(s))
		for i, v := range s { result[i] = f(v) }
		return result
	}
	data := []int{1, 2, 3}
	doubled := mapFn(data, func(n int) int { return n * 2 })
	return doubled[0] + doubled[1] + doubled[2]
}

func ClosureReduce() int {
	reduce := func(s []int, init int, f func(int, int) int) int {
		acc := init
		for _, v := range s { acc = f(acc, v) }
		return acc
	}
	data := []int{1, 2, 3, 4, 5}
	return reduce(data, 0, func(a, b int) int { return a + b })
}

func ClosureInStruct() int {
	type Processor struct {
		Transform func(int) int
	}
	p := Processor{Transform: func(x int) int { return x * x }}
	return p.Transform(5)
}

func ClosureStringProcessor() string {
	process := func(s string) string {
		return strings.TrimSpace(s)
	}
	return process("  hello  ")
}

func FmtClosure() string {
	fn := func() int { return 42 }
	return fmt.Sprintf("%d", fn())
}
