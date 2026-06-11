package divergence_hunt119

import "fmt"

// ============================================================================
// Round 119: Array operations and fixed-size types
// ============================================================================

func ArrayBasic() string {
	var arr [5]int
	arr[2] = 42
	return fmt.Sprintf("%d", arr[2])
}

func ArrayLiteral() string {
	arr := [3]string{"a", "b", "c"}
	return fmt.Sprintf("%s", arr[1])
}

func ArrayAutoLen() string {
	arr := [...]int{10, 20, 30}
	return fmt.Sprintf("%d", len(arr))
}

func ArrayCopy() string {
	a := [3]int{1, 2, 3}
	b := a
	b[0] = 99
	return fmt.Sprintf("%d:%d", a[0], b[0])
}

func ArrayRange() string {
	arr := [4]int{10, 20, 30, 40}
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return fmt.Sprintf("%d", sum)
}

func ArrayPointer() string {
	arr := [3]int{1, 2, 3}
	p := &arr
	p[1] = 99
	return fmt.Sprintf("%d", arr[1])
}

func ArrayCompare() string {
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}
	return fmt.Sprintf("%v", a == b)
}

func ArrayNotEqual() string {
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 4}
	return fmt.Sprintf("%v", a != b)
}

func ArrayZeroValue() string {
	var arr [5]int
	return fmt.Sprintf("%d", arr[0])
}

func ArrayOfStruct() string {
	type Point struct{ X, Y int }
	arr := [3]Point{{1, 2}, {3, 4}, {5, 6}}
	return fmt.Sprintf("%d:%d", arr[1].X, arr[1].Y)
}

func ArrayMultiDim() string {
	var grid [2][3]int
	grid[0][1] = 7
	grid[1][2] = 9
	return fmt.Sprintf("%d:%d", grid[0][1], grid[1][2])
}

func ArrayLenCap() string {
	arr := [5]int{}
	return fmt.Sprintf("%d:%d", len(arr), cap(arr))
}
