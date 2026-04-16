package divergence_hunt32

import "fmt"

// ============================================================================
// Round 32: Fixed arrays - arrays vs slices, array copy semantics, len/cap
// ============================================================================

func ArrayLenCap() int {
	a := [5]int{1, 2, 3, 4, 5}
	return len(a)*10 + cap(a)
}

func ArrayCopyValue() int {
	a := [3]int{10, 20, 30}
	b := a
	b[0] = 99
	return a[0] // value copy: a unchanged
}

func ArrayPointerModify() int {
	a := [3]int{10, 20, 30}
	p := &a
	p[0] = 99
	return a[0] // pointer: a changed
}

func ArrayAsArg() int {
	sum := func(a [3]int) int { return a[0] + a[1] + a[2] }
	a := [3]int{1, 2, 3}
	return sum(a)
}

func ArrayPointerAsArg() int {
	modify := func(a *[3]int) { a[0] = 99 }
	a := [3]int{1, 2, 3}
	modify(&a)
	return a[0]
}

func ArrayIteration() int {
	a := [4]int{10, 20, 30, 40}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

func ArrayIndexAccess() int {
	a := [5]int{0, 10, 20, 30, 40}
	return a[2] + a[4]
}

func ArrayZeroValue() int {
	var a [3]int
	return a[0] + a[1] + a[2] // all zeros
}

func ArrayOfString() int {
	a := [3]string{"hello", "world", "foo"}
	return len(a[0]) + len(a[1]) + len(a[2])
}

func ArrayOfStruct() int {
	type P struct{ X, Y int }
	a := [2]P{{1, 2}, {3, 4}}
	return a[0].X + a[1].Y
}

func SliceFromArray() int {
	a := [5]int{10, 20, 30, 40, 50}
	s := a[1:4]
	return s[0] + s[1] + s[2]
}

func ArrayComparison() bool {
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}
	c := [3]int{1, 2, 4}
	return a == b && a != c
}

func MultiDimensionalArray() int {
	a := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	return a[0][1] + a[1][2]
}

func ArrayInStruct() int {
	type Matrix struct{ Data [4]int }
	m := Matrix{Data: [4]int{1, 2, 3, 4}}
	return m.Data[0] + m.Data[3]
}

func FmtArray() string {
	a := [3]int{1, 2, 3}
	return fmt.Sprintf("%v", a)
}

func ArrayLiteralPartial() int {
	a := [5]int{1, 2} // rest are zeros
	return a[0] + a[1] + a[2] + a[3] + a[4]
}

func ArrayLiteralIndex() int {
	a := [5]int{2: 10, 4: 20} // index: value
	return a[0] + a[2] + a[4]
}
