package divergence_hunt105

import "fmt"

// ============================================================================
// Round 105: Composite literal edge cases
// ============================================================================

func NestedSliceLiteral() string {
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	return fmt.Sprintf("%d:%d", grid[0][1], grid[1][2])
}

func MapLiteralWithStruct() string {
	type Point struct{ X, Y int }
	m := map[string]Point{
		"origin": {0, 0},
		"a":      {1, 2},
	}
	return fmt.Sprintf("%d:%d", m["a"].X, m["a"].Y)
}

func SliceOfMap() string {
	items := []map[string]int{
		{"a": 1},
		{"b": 2, "c": 3},
	}
	return fmt.Sprintf("%d:%d", len(items), items[1]["c"])
}

func StructWithSlice() string {
	type Bag struct {
		Items []string
		Size  int
	}
	b := Bag{Items: []string{"x", "y", "z"}, Size: 3}
	return fmt.Sprintf("%d:%d", len(b.Items), b.Size)
}

func NestedMapLiteral() string {
	m := map[string]map[string]int{
		"outer": {"inner": 42},
	}
	return fmt.Sprintf("%d", m["outer"]["inner"])
}

func SliceOfFunc() string {
	type Op func(int, int) int
	ops := []Op{
		func(a, b int) int { return a + b },
		func(a, b int) int { return a * b },
	}
	return fmt.Sprintf("%d:%d", ops[0](3, 4), ops[1](3, 4))
}

func EmptyCompositeLiterals() string {
	s := []int{}
	m := map[string]int{}
	return fmt.Sprintf("%d:%d", len(s), len(m))
}

func PointerStructLiteral() string {
	type Item struct{ Val int }
	p := &Item{Val: 42}
	return fmt.Sprintf("%d", p.Val)
}

func NestedStructLiteral() string {
	type Inner struct{ V int }
	type Outer struct{ I Inner; Name string }
	o := Outer{I: Inner{V: 10}, Name: "test"}
	return fmt.Sprintf("%d:%s", o.I.V, o.Name)
}

func ArrayLiteral() string {
	arr := [5]int{10, 20, 30, 40, 50}
	return fmt.Sprintf("%d", arr[2])
}
