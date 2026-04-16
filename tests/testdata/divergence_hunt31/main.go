package divergence_hunt31

import "fmt"

// ============================================================================
// Round 31: Struct methods - value receiver vs pointer receiver
// ============================================================================

type Counter31 struct {
	n int
}

func (c Counter31) Value() int  { return c.n }
func (c *Counter31) Inc()      { c.n++ }
func (c *Counter31) Add(n int) { c.n += n }

func ValueReceiverNoMutation() int {
	c := Counter31{n: 10}
	c.Inc() // pointer receiver mutates
	return c.n
}

func PointerReceiverChain() int {
	c := &Counter31{n: 0}
	c.Inc()
	c.Inc()
	c.Add(5)
	return c.Value()
}

func ValueReceiverCopy() int {
	c := Counter31{n: 5}
	v := c.Value()
	c.Inc()
	return v + c.Value() // 5 + 6
}

func StructMethodOnLiteral() int {
	c := &Counter31{n: 100}
	return c.Value()
}

func NestedMethodCall() int {
	c := &Counter31{n: 1}
	c.Inc()
	c.Add(c.Value())
	return c.Value() // 1 + 1 + 2 = 4
}

type Rect31 struct {
	W, H int
}

func (r Rect31) Area() int      { return r.W * r.H }
func (r *Rect31) Scale(n int)   { r.W *= n; r.H *= n }
func (r Rect31) Perimeter() int { return 2*r.W + 2*r.H }

func MethodValueVsPointer() int {
	r := &Rect31{W: 3, H: 4}
	area := r.Area()
	r.Scale(2)
	return area + r.Area() // 12 + 48
}

func MethodOnValueStruct() int {
	r := Rect31{W: 5, H: 6}
	return r.Area() + r.Perimeter() // 30 + 22
}

type Stringer31 interface {
	String() string
}

type Name31 struct {
	First, Last string
}

func (n Name31) String() string { return n.First + " " + n.Last }

func InterfaceMethodCall() string {
	var s Stringer31 = Name31{First: "Alice", Last: "Smith"}
	return s.String()
}

func InterfaceMethodOnPointer() string {
	var s Stringer31 = &Name31{First: "Bob", Last: "Jones"}
	return s.String()
}

func MethodReturnsMultipleValues() int {
	divide := func(a, b int) (int, int) { return a / b, a % b }
	q, r := divide(17, 5)
	return q*10 + r
}

func StructWithBoolMethod() int {
	type Validator struct{ Min, Max int }
	isValid := func(v Validator, x int) bool { return x >= v.Min && x <= v.Max }
	v := Validator{Min: 1, Max: 10}
	result := 0
	if isValid(v, 5) { result += 1 }
	if isValid(v, 15) { result += 10 }
	return result
}

func FmtStructWithMethods() string {
	n := Name31{First: "Hello", Last: "World"}
	return fmt.Sprintf("%v", n.String())
}

func StructSliceWithMethods() int {
	items := []Name31{
		{"Alice", "A"},
		{"Bob", "B"},
	}
	sum := 0
	for _, n := range items {
		sum += len(n.First) + len(n.Last)
	}
	return sum
}

func EmbedStructMethod() int {
	type Base struct{ X int }
	type Derived struct {
		Base
		Y int
	}
	d := Derived{Base: Base{X: 10}, Y: 20}
	return d.X + d.Y
}
