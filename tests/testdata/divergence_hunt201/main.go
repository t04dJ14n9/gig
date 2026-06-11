package divergence_hunt201

import "fmt"

// ============================================================================
// Round 201: Array copy semantics and comparison
// ============================================================================

func ArrayCopyByValue() string {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 99
	return fmt.Sprintf("%v:%v", a, b)
}

func ArrayEquality() string {
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}
	c := [3]int{3, 2, 1}
	return fmt.Sprintf("%v:%v", a == b, a == c)
}

func ArrayInequality() string {
	a := [2]int{1, 2}
	b := [2]int{1, 3}
	return fmt.Sprintf("%v", a != b)
}

func ArraySliceShare() string {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]
	s[0] = 99
	return fmt.Sprintf("%v:%v", a, s)
}

func ArrayOfArrays() string {
	a := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	b := a
	b[0][0] = 99
	return fmt.Sprintf("%v", a[0][0])
}

func ArrayZeroValue() string {
	var a [5]int
	return fmt.Sprintf("%v:%d", a, len(a))
}

func ArrayLenCap() string {
	a := [10]int{1, 2, 3}
	return fmt.Sprintf("%d:%d", len(a), cap(a))
}

func ArrayRange() string {
	sum := 0
	a := [4]int{10, 20, 30, 40}
	for _, v := range a {
		sum += v
	}
	return fmt.Sprintf("%d", sum)
}

func ArrayLiteral() string {
	a := [...]int{1, 2, 3, 4, 5}
	return fmt.Sprintf("%d", len(a))
}

func ArrayIndexAccess() string {
	a := [3]string{"a", "b", "c"}
	return fmt.Sprintf("%s%s%s", a[0], a[1], a[2])
}

func ArrayPointerDeref() string {
	a := [3]int{1, 2, 3}
	p := &a
	(*p)[0] = 99
	return fmt.Sprintf("%v", a)
}
