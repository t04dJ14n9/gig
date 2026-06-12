package divergence_hunt68

// ============================================================================
// Round 68: Struct/method receiver edge cases - pointer vs value receiver
// ============================================================================

type Counter struct {
	count int
}

func (c *Counter) Inc() {
	c.count++
}

func (c *Counter) Get() int {
	return c.count
}

func (c Counter) ValueGet() int {
	return c.count
}

func PointerReceiverModify() int {
	c := Counter{}
	c.Inc()
	c.Inc()
	return c.Get()
}

func ValueReceiverNoModify() int {
	c := Counter{count: 5}
	c.ValueGet() // value receiver doesn't modify
	return c.count
}

type Pair struct {
	X, Y int
}

func (p Pair) Sum() int {
	return p.X + p.Y
}

func (p *Pair) Scale(f int) {
	p.X *= f
	p.Y *= f
}

func MixedReceiver() int {
	p := Pair{X: 2, Y: 3}
	p.Scale(10) // pointer receiver modifies
	return p.Sum()
}

func StructLiteral() int {
	p := Pair{X: 10, Y: 20}
	return p.X + p.Y
}

func StructZeroValue() int {
	var p Pair
	return p.X + p.Y
}

func StructPointerLiteral() int {
	p := &Pair{X: 5, Y: 15}
	return p.X + p.Y
}

func StructFieldAssign() int {
	p := Pair{}
	p.X = 10
	p.Y = 20
	return p.X + p.Y
}

func StructPointerFieldAssign() int {
	p := &Pair{}
	p.X = 10
	p.Y = 20
	return p.X + p.Y
}

type Nested struct {
	Inner Pair
	Name  string
}

func StructNested() int {
	n := Nested{Inner: Pair{X: 1, Y: 2}, Name: "test"}
	return n.Inner.X + n.Inner.Y
}

func StructNestedFieldAssign() int {
	n := Nested{}
	n.Inner.X = 10
	n.Inner.Y = 20
	return n.Inner.X + n.Inner.Y
}

type Rect struct {
	W, H int
}

func (r Rect) Area() int {
	return r.W * r.H
}

func (r *Rect) Double() {
	r.W *= 2
	r.H *= 2
}

func StructMethodChain() int {
	r := Rect{W: 3, H: 4}
	area1 := r.Area()
	r.Double()
	area2 := r.Area()
	return area1 + area2
}

func StructCopySemantics() int {
	a := Pair{X: 10, Y: 20}
	b := a // value copy
	b.X = 99
	return a.X // should still be 10
}

func StructPointerCopy() int {
	a := &Pair{X: 10, Y: 20}
	b := a // pointer copy
	b.X = 99
	return a.X // should be 99 (shared)
}

func MethodOnLiteral() int {
	return Pair{X: 3, Y: 4}.Sum()
}
