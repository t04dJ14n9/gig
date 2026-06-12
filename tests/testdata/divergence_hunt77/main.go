package divergence_hunt77

// ============================================================================
// Round 77: Composite literal edge cases - nested, slice of struct, map of struct
// ============================================================================

type Item struct {
	ID    int
	Name  string
	Price float64
}

func SliceOfStructLiteral() int {
	items := []Item{
		{ID: 1, Name: "A", Price: 10.0},
		{ID: 2, Name: "B", Price: 20.0},
	}
	return items[0].ID + items[1].ID
}

func MapOfStructLiteral() int {
	m := map[string]Item{
		"first":  {ID: 1, Name: "A", Price: 10.0},
		"second": {ID: 2, Name: "B", Price: 20.0},
	}
	return m["first"].ID + m["second"].ID
}

func NestedSliceLiteral() int {
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	return grid[0][1] + grid[1][0]
}

func ArrayLiteral() int {
	arr := [3]int{10, 20, 30}
	return arr[0] + arr[1] + arr[2]
}

func ArrayAutoLen() int {
	arr := [...]int{10, 20, 30, 40}
	return len(arr)
}

func MapLiteralEmpty() int {
	m := map[string]int{}
	m["a"] = 1
	return m["a"]
}

func SliceLiteralEmpty() int {
	s := []int{}
	s = append(s, 1)
	return len(s)
}

func StructLiteralPositional() int {
	p := struct {
		X, Y int
	}{10, 20}
	return p.X + p.Y
}

func StructLiteralNamed() int {
	p := struct {
		X int
		Y int
	}{X: 10, Y: 20}
	return p.X + p.Y
}

func SliceOfMap() int {
	sm := []map[string]int{
		{"a": 1},
		{"b": 2},
	}
	return sm[0]["a"] + sm[1]["b"]
}

func MapKeyStruct() int {
	type Key struct{ K string }
	m := map[Key]int{
		{K: "x"}: 10,
		{K: "y"}: 20,
	}
	return m[Key{K: "x"}]
}

func NestedMapLiteral() int {
	m := map[string]map[string]int{
		"outer": {"inner": 42},
	}
	return m["outer"]["inner"]
}

func PointerStructLiteral() int {
	p := &Item{ID: 1, Name: "test", Price: 9.99}
	return p.ID
}

func SliceOfPointer() int {
	items := []*Item{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	}
	return items[0].ID + items[1].ID
}
