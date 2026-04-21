package divergence_hunt196

import (
	"fmt"
)

// ============================================================================
// Round 196: Array value semantics
// ============================================================================

// ArrayCopyByValue tests array copy behavior
func ArrayCopyByValue() string {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 100
	return fmt.Sprintf("%d:%d", a[0], b[0])
}

// ArrayPassByValue tests array passed to function
func ArrayPassByValue() string {
	modify := func(arr [3]int) {
		arr[0] = 100
	}
	a := [3]int{1, 2, 3}
	modify(a)
	return fmt.Sprintf("%d", a[0])
}

// ArrayReturnByValue tests array returned from function
func ArrayReturnByValue() string {
	create := func() [3]int {
		return [3]int{4, 5, 6}
	}
	a := create()
	return fmt.Sprintf("%d:%d:%d", a[0], a[1], a[2])
}

// ArrayEquality tests array equality comparison
func ArrayEquality() string {
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}
	c := [3]int{1, 2, 4}
	return fmt.Sprintf("%v:%v", a == b, a == c)
}

// ArrayOfStructs tests array of structs value semantics
func ArrayOfStructs() string {
	type Point struct{ X, Y int }
	a := [2]Point{{1, 2}, {3, 4}}
	b := a
	b[0].X = 100
	return fmt.Sprintf("%d:%d", a[0].X, b[0].X)
}

// ArraySliceConversion tests array to slice conversion
func ArraySliceConversion() string {
	a := [3]int{1, 2, 3}
	s := a[:]
	s[0] = 100
	return fmt.Sprintf("%d:%d", a[0], s[0])
}

// ArrayIndexing tests array indexing
func ArrayIndexing() string {
	a := [5]int{10, 20, 30, 40, 50}
	return fmt.Sprintf("%d:%d:%d", a[0], a[2], a[4])
}

// ArrayLength tests array length
func ArrayLength() string {
	a := [5]int{1, 2, 3, 4, 5}
	b := [0]int{}
	return fmt.Sprintf("%d:%d", len(a), len(b))
}

// ArrayLiteral tests array literal syntax
func ArrayLiteral() string {
	a := [...]int{1, 2, 3, 4, 5}
	b := [5]int{1, 2}
	return fmt.Sprintf("%d:%d:%d", len(a), b[0], b[4])
}

// ArrayIteration tests array iteration
func ArrayIteration() string {
	a := [3]int{2, 4, 6}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return fmt.Sprintf("%d", sum)
}

// ArrayOfArrays tests multidimensional arrays
func ArrayOfArrays() string {
	a := [2][3]int{{1, 2, 3}, {4, 5, 6}}
	b := a
	b[0][0] = 100
	return fmt.Sprintf("%d:%d", a[0][0], b[0][0])
}
