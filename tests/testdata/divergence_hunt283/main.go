package divergence_hunt283

import (
	"fmt"
)

// ============================================================================
// Round 283: Struct constructor patterns — factory functions, method chaining, self-reference

type Builder struct {
	name  string
	age   int
	items []string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithName(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) WithAge(age int) *Builder {
	b.age = age
	return b
}

func (b *Builder) WithItem(item string) *Builder {
	b.items = append(b.items, item)
	return b
}

func (b *Builder) Build() string {
	return fmt.Sprintf("name=%s,age=%d,items=%v", b.name, b.age, b.items)
}

type ImmutablePoint struct {
	X, Y float64
}

func NewPoint(x, y float64) ImmutablePoint {
	return ImmutablePoint{X: x, Y: y}
}

func (p ImmutablePoint) Add(other ImmutablePoint) ImmutablePoint {
	return ImmutablePoint{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p ImmutablePoint) String() string {
	return fmt.Sprintf("(%.0f,%.0f)", p.X, p.Y)
}

// BuilderPattern tests fluent builder pattern
func BuilderPattern() string {
	b := NewBuilder().WithName("Alice").WithAge(30).WithItem("book").WithItem("pen")
	return b.Build()
}

// BuilderPartial tests partial builder
func BuilderPartial() string {
	b := NewBuilder().WithName("Bob")
	return b.Build()
}

// ImmutableValueMethodChaining tests immutable value type chaining
func ImmutableValueMethodChaining() string {
	p1 := NewPoint(1, 2)
	p2 := NewPoint(3, 4)
	p3 := p1.Add(p2)
	return p3.String()
}

// StructCopyByValue tests that struct assignment copies
func StructCopyByValue() string {
	type Data struct{ Value int }
	a := Data{Value: 10}
	b := a
	b.Value = 20
	return fmt.Sprintf("a=%d,b=%d", a.Value, b.Value)
}

// StructCopyPointer tests that pointer assignment shares
func StructCopyPointer() string {
	type Data struct{ Value int }
	a := &Data{Value: 10}
	b := a
	b.Value = 20
	return fmt.Sprintf("a=%d,b=%d", a.Value, b.Value)
}

// StructWithSliceCopy tests struct with slice field — shallow copy
func StructWithSliceCopy() string {
	type Data struct{ Items []int }
	a := Data{Items: []int{1, 2, 3}}
	b := a
	b.Items[0] = 99
	return fmt.Sprintf("a=%v,b=%v", a.Items, b.Items)
}

// NestedStructInit tests nested struct initialization
func NestedStructInit() string {
	type Inner struct{ Val int }
	type Outer struct {
		Inner Inner
		Name  string
	}
	o := Outer{Inner: Inner{Val: 42}, Name: "test"}
	return fmt.Sprintf("val=%d,name=%s", o.Inner.Val, o.Name)
}

// StructZeroValueFields tests that all zero values are correct
func StructZeroValueFields() string {
	type All struct {
		I  int
		F  float64
		S  string
		B  bool
		P  *int
		M  map[string]int
		S_ []int
	}
	var a All
	return fmt.Sprintf("i=%d,f=%f,s=%q,b=%t,p=%t,m=%t,s=%t",
		a.I, a.F, a.S, a.B, a.P == nil, a.M == nil, a.S_ == nil)
}

// StructComparison tests struct equality (comparable fields only)
func StructComparison() string {
	type Pair struct{ A, B int }
	p1 := Pair{1, 2}
	p2 := Pair{1, 2}
	p3 := Pair{1, 3}
	return fmt.Sprintf("eq12=%t,eq13=%t", p1 == p2, p1 == p3)
}
