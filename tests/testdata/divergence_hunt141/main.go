package divergence_hunt141

import "fmt"

// ============================================================================
// Round 141: Closure variable capture — common Go gotchas
// ============================================================================

func ClosureLoopCapture() string {
	var results []string
	for i := 0; i < 3; i++ {
		func() {
			results = append(results, fmt.Sprintf("i=%d", i))
		}()
	}
	return fmt.Sprintf("%v", results)
}

func ClosureLoopDeferred() string {
	var result string
	for i := 0; i < 3; i++ {
		defer func() {
			result += fmt.Sprintf("%d", i)
		}()
	}
	_ = result
	return "deferred"
}

func ClosureShadowVar() string {
	x := "outer"
	f := func() string {
		x := "inner"
		return x
	}
	return fmt.Sprintf("f=%s-x=%s", f(), x)
}

func ClosureMutateOuter() string {
	x := 10
	f := func() {
		x = 20
	}
	f()
	return fmt.Sprintf("x=%d", x)
}

func ClosureMultipleCaptures() string {
	x := 1
	y := 2
	f := func() int {
		return x + y
	}
	x = 10
	y = 20
	return fmt.Sprintf("sum=%d", f())
}

func ClosureReturned() string {
	makeFunc := func() func() string {
		s := "captured"
		return func() string {
			return s
		}
	}
	f := makeFunc()
	return f()
}

func ClosureSliceAppend() string {
	s := []int{1, 2, 3}
	f := func() {
		s = append(s, 4)
	}
	f()
	return fmt.Sprintf("%v", s)
}

func ClosureMapModify() string {
	m := map[string]int{"a": 1}
	f := func() {
		m["b"] = 2
	}
	f()
	return fmt.Sprintf("len=%d", len(m))
}

func ClosureNested() string {
	outer := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	add3 := outer(3)
	return fmt.Sprintf("3+4=%d", add3(4))
}
