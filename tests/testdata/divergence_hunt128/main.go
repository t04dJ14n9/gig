package divergence_hunt128

import "fmt"

// ============================================================================
// Round 128: Multiple assignment and tuple returns
// ============================================================================

func SwapVariables() string {
	a, b := 1, 2
	a, b = b, a
	return fmt.Sprintf("a=%d-b=%d", a, b)
}

func MultiReturnAssign() string {
	fn := func() (int, string) { return 42, "hello" }
	x, y := fn()
	return fmt.Sprintf("x=%d-y=%s", x, y)
}

func BlankAssign() string {
	_, y := func() (int, string) { return 1, "two" }()
	return y
}

func MultiAssignExpression() string {
	x := 0
	y := 0
	x, y = x+1, y+2
	return fmt.Sprintf("x=%d-y=%d", x, y)
}

func MultiAssignSwap() string {
	arr := []int{10, 20, 30}
	arr[0], arr[2] = arr[2], arr[0]
	return fmt.Sprintf("%v", arr)
}

func MultiAssignMap() string {
	m := map[string]int{"a": 1, "b": 2}
	v1, ok1 := m["a"]
	v2, ok2 := m["c"]
	return fmt.Sprintf("a=%d-%t-c=%d-%t", v1, ok1, v2, ok2)
}

func NestedMultiReturn() string {
	inner := func() (int, int) { return 3, 4 }
	outer := func() (int, int, int) {
		a, b := inner()
		return a, b, a + b
	}
	x, y, z := outer()
	return fmt.Sprintf("%d-%d-%d", x, y, z)
}

func AssignDifferentTypes() string {
	var i int
	var s string
	var f float64
	i, s, f = 1, "two", 3.0
	return fmt.Sprintf("i=%d-s=%s-f=%.1f", i, s, f)
}

func MultiAssignStruct() string {
	type Pair struct{ A, B int }
	p := Pair{A: 1, B: 2}
	p.A, p.B = p.B, p.A
	return fmt.Sprintf("%d-%d", p.A, p.B)
}
