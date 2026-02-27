package structs

// BasicStruct tests struct creation and field access.
func BasicStruct() int {
	type Point struct {
		X int
		Y int
	}
	p := Point{X: 3, Y: 4}
	return p.X + p.Y
}

// StructPointer tests struct pointer creation and field mutation.
func StructPointer() int {
	type Rect struct {
		W int
		H int
	}
	r := &Rect{W: 10, H: 20}
	r.W = 15
	return r.W + r.H
}

// NestedStruct tests structs containing other structs.
func NestedStruct() int {
	type Inner struct {
		Val int
	}
	type Outer struct {
		A Inner
		B Inner
	}
	o := Outer{
		A: Inner{Val: 10},
		B: Inner{Val: 20},
	}
	return o.A.Val + o.B.Val
}

// EmbeddedField tests anonymous (embedded) struct fields.
func EmbeddedField() int {
	type Base struct {
		ID   int
		Name int
	}
	type Extended struct {
		Base
		Extra int
	}
	e := Extended{
		Base:  Base{ID: 42, Name: 7},
		Extra: 100,
	}
	// Access promoted fields from Base directly
	return e.ID + e.Name + e.Extra
}

// StructInSlice tests using structs in a slice.
func StructInSlice() int {
	type Item struct {
		Value int
	}
	items := make([]Item, 3)
	items[0] = Item{Value: 10}
	items[1] = Item{Value: 20}
	items[2] = Item{Value: 30}
	sum := 0
	for _, item := range items {
		sum += item.Value
	}
	return sum
}

// StructAsParam tests passing structs to and returning from functions.
func StructAsParam() int {
	type Pair struct {
		A int
		B int
	}
	add := func(p Pair) int {
		return p.A + p.B
	}
	p := Pair{A: 13, B: 29}
	return add(p)
}

// StructZeroValue tests zero values of struct fields.
func StructZeroValue() int {
	type Config struct {
		Width  int
		Height int
		Depth  int
	}
	var c Config
	// All fields should be zero
	c.Width = 5
	c.Height = 10
	return c.Width + c.Height + c.Depth
}

// MultipleEmbedded tests multiple embedded structs.
func MultipleEmbedded() int {
	type Position struct {
		X int
		Y int
	}
	type Velocity struct {
		DX int
		DY int
	}
	type Entity struct {
		Position
		Velocity
		HP int
	}
	e := Entity{
		Position: Position{X: 1, Y: 2},
		Velocity: Velocity{DX: 3, DY: 4},
		HP:       100,
	}
	return e.X + e.Y + e.DX + e.DY + e.HP
}

// DeepNesting tests deeply nested struct access.
func DeepNesting() int {
	type A struct {
		Val int
	}
	type B struct {
		Inner A
	}
	type C struct {
		Inner B
	}
	c := C{Inner: B{Inner: A{Val: 42}}}
	return c.Inner.Inner.Val
}

// StructFieldMutation tests mutating struct fields through a pointer.
func StructFieldMutation() int {
	type Counter struct {
		Count int
	}
	c := &Counter{Count: 0}
	for i := 0; i < 10; i++ {
		c.Count += i
	}
	return c.Count
}

// StructWithBool tests struct with different field types.
func StructWithBool() int {
	type Flags struct {
		Active  bool
		Count   int
		Score   int
	}
	f := Flags{Active: true, Count: 5, Score: 10}
	result := f.Count + f.Score
	if f.Active {
		result += 100
	}
	return result
}
