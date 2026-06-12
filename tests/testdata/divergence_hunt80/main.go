package divergence_hunt80

// ============================================================================
// Round 80: Multiple assignment edge cases - swap, multi-return, blank
// ============================================================================

func SwapInt() (int, int) {
	a, b := 1, 2
	a, b = b, a
	return a, b
}

func MultiReturnAssign() (int, int) {
	divmod := func(a, b int) (int, int) {
		return a / b, a % b
	}
	q, r := divmod(17, 5)
	return q, r
}

func BlankAssign() int {
	_, b := 10, 20
	return b
}

func MultiAssignSameVar() int {
	x := 0
	x, x = 1, 2
	return x
}

func MultiAssignSwap() int {
	a, b := 1, 2
	a, b = b, a+b
	return a*10 + b
}

func AssignMapAccess() (int, bool) {
	m := map[string]int{"key": 42}
	v, ok := m["key"]
	return v, ok
}

func AssignTypeAssertion() (int, bool) {
	var x any = "hello"
	v, ok := x.(int)
	return v, ok
}

func AssignTypeAssertionString() (string, bool) {
	var x any = "hello"
	v, ok := x.(string)
	return v, ok
}

func MultiReturnBlank() int {
	fn := func() (int, int) { return 10, 20 }
	_, b := fn()
	return b
}

func SwapSliceElements() []int {
	s := []int{1, 2, 3}
	s[0], s[2] = s[2], s[0]
	return s
}

func AssignStructFields() (int, int) {
	type S struct{ X, Y int }
	s := S{X: 10, Y: 20}
	s.X, s.Y = s.Y, s.X
	return s.X, s.Y
}

func MultiAssignExpression() (int, int) {
	x := 5
	y := 10
	x, y = x+y, x*y
	return x, y
}

func AssignPointerDeref() (int, int) {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = 100, 200
	return a, b
}

func NestedMultiReturn() int {
	fn := func() (int, int) {
		inner := func() (int, int) { return 3, 4 }
		a, b := inner()
		return a + b, a * b
	}
	x, y := fn()
	return x + y
}
