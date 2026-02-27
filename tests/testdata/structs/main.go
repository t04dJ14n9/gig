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
		Active bool
		Count  int
		Score  int
	}
	f := Flags{Active: true, Count: 5, Score: 10}
	result := f.Count + f.Score
	if f.Active {
		result += 100
	}
	return result
}

// StructCopySemantics tests that struct assignment copies values (not references).
func StructCopySemantics() int {
	type Point struct {
		X int
		Y int
	}
	a := Point{X: 10, Y: 20}
	b := a // copy
	b.X = 99
	b.Y = 88
	// a should be unchanged
	return a.X + a.Y + b.X + b.Y
}

// StructPointerSharing tests that pointer assignment shares the struct.
func StructPointerSharing() int {
	type Box struct {
		Size int
	}
	a := &Box{Size: 10}
	b := a // both point to the same struct
	b.Size = 50
	return a.Size + b.Size // both should be 50
}

// StructReturnFromFunc tests returning a struct from a function.
func StructReturnFromFunc() int {
	type Pair struct {
		A int
		B int
	}
	makePair := func(a, b int) Pair {
		return Pair{A: a, B: b}
	}
	p := makePair(17, 25)
	return p.A + p.B
}

// StructPointerReturnFromFunc tests returning a struct pointer from a function.
func StructPointerReturnFromFunc() int {
	type Node struct {
		Val  int
		Next int // simplified: using int instead of *Node
	}
	makeNode := func(v int) *Node {
		return &Node{Val: v, Next: 0}
	}
	n := makeNode(42)
	n.Next = 100
	return n.Val + n.Next
}

// StructSliceAppend tests appending structs to a slice.
func StructSliceAppend() int {
	type Entry struct {
		Key   int
		Value int
	}
	var entries []Entry
	for i := 0; i < 5; i++ {
		entries = append(entries, Entry{Key: i, Value: i * 10})
	}
	sum := 0
	for _, e := range entries {
		sum += e.Key + e.Value
	}
	return sum
}

// StructPointerSlice tests a slice of struct pointers.
func StructPointerSlice() int {
	type Item struct {
		Score int
	}
	items := make([]*Item, 5)
	for i := 0; i < 5; i++ {
		items[i] = &Item{Score: (i + 1) * 10}
	}
	// Mutate through pointers
	for _, item := range items {
		item.Score += 5
	}
	sum := 0
	for _, item := range items {
		sum += item.Score
	}
	return sum
}

// StructInMap tests using structs as map values.
func StructInMap() int {
	type Score struct {
		Points int
		Bonus  int
	}
	scores := make(map[string]Score)
	scores["alice"] = Score{Points: 100, Bonus: 20}
	scores["bob"] = Score{Points: 85, Bonus: 15}
	a := scores["alice"]
	b := scores["bob"]
	return a.Points + a.Bonus + b.Points + b.Bonus
}

// StructConditionalInit tests conditional struct initialization.
func StructConditionalInit() int {
	type Config struct {
		Mode  int
		Level int
	}
	var c Config
	flag := true
	if flag {
		c = Config{Mode: 1, Level: 10}
	} else {
		c = Config{Mode: 2, Level: 20}
	}
	return c.Mode + c.Level
}

// StructFieldLoop tests iterating and accumulating struct field values.
func StructFieldLoop() int {
	type Triple struct {
		A int
		B int
		C int
	}
	triples := []Triple{
		{A: 1, B: 2, C: 3},
		{A: 4, B: 5, C: 6},
		{A: 7, B: 8, C: 9},
	}
	sum := 0
	for _, t := range triples {
		sum += t.A + t.B + t.C
	}
	return sum
}

// StructNestedMutation tests mutating deeply nested struct fields via pointer.
func StructNestedMutation() int {
	type Inner struct {
		Val int
	}
	type Outer struct {
		Data Inner
	}
	o := &Outer{Data: Inner{Val: 10}}
	o.Data.Val = 42
	return o.Data.Val
}

// StructEmbeddedOverride tests that a field in the outer struct shadows an embedded field.
func StructEmbeddedOverride() int {
	type Base struct {
		X int
	}
	type Derived struct {
		Base
		X int // shadows Base.X
	}
	d := Derived{
		Base: Base{X: 10},
		X:    99,
	}
	// d.X should be the Derived.X (99), not Base.X (10)
	return d.X + d.Base.X
}

// StructWithClosure tests struct fields captured by closures.
func StructWithClosure() int {
	type Accum struct {
		Total int
	}
	a := &Accum{Total: 0}
	add := func(x int) {
		a.Total += x
	}
	for i := 1; i <= 10; i++ {
		add(i)
	}
	return a.Total
}

// StructReassign tests reassigning struct fields after creation.
func StructReassign() int {
	type Pair struct {
		X int
		Y int
	}
	p := Pair{X: 1, Y: 2}
	p.X = 10
	p.Y = 20
	return p.X + p.Y
}

// StructSliceOfNested tests slice of nested structs with accumulation.
func StructSliceOfNested() int {
	type Coord struct {
		X int
		Y int
	}
	type Line struct {
		Start Coord
		End   Coord
	}
	lines := []Line{
		{Start: Coord{X: 0, Y: 0}, End: Coord{X: 3, Y: 4}},
		{Start: Coord{X: 1, Y: 1}, End: Coord{X: 5, Y: 6}},
	}
	sum := 0
	for _, l := range lines {
		sum += l.Start.X + l.Start.Y + l.End.X + l.End.Y
	}
	return sum
}

// StructMultiReturn tests returning multiple values including a struct.
func StructMultiReturn() int {
	type Result struct {
		Val int
		Ok  bool
	}
	compute := func(x int) (Result, int) {
		return Result{Val: x * 2, Ok: true}, x + 1
	}
	r, extra := compute(20)
	total := r.Val + extra
	if r.Ok {
		total += 100
	}
	return total
}

// StructBuilderPattern tests building up a struct incrementally.
func StructBuilderPattern() int {
	type Builder struct {
		Width  int
		Height int
		Depth  int
	}
	setWidth := func(b *Builder, w int) {
		b.Width = w
	}
	setHeight := func(b *Builder, h int) {
		b.Height = h
	}
	setDepth := func(b *Builder, d int) {
		b.Depth = d
	}
	b := &Builder{}
	setWidth(b, 10)
	setHeight(b, 20)
	setDepth(b, 30)
	return b.Width + b.Height + b.Depth
}

// StructArrayField tests struct containing an array-like field (simulated with int fields).
func StructArrayField() int {
	type Matrix struct {
		R0C0 int
		R0C1 int
		R1C0 int
		R1C1 int
	}
	// 2x2 matrix multiply identity
	m := Matrix{R0C0: 1, R0C1: 0, R1C0: 0, R1C1: 1}
	// Trace of identity matrix = sum of diagonal
	return m.R0C0 + m.R1C1
}

// StructEmbeddedChain tests chained embedded struct access (3 levels).
func StructEmbeddedChain() int {
	type Level1 struct {
		A int
	}
	type Level2 struct {
		Level1
		B int
	}
	type Level3 struct {
		Level2
		C int
	}
	v := Level3{
		Level2: Level2{
			Level1: Level1{A: 10},
			B:      20,
		},
		C: 30,
	}
	// Promoted field access through two levels
	return v.A + v.B + v.C
}
