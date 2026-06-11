package divergence_hunt273

import (
	"fmt"
)

// ============================================================================
// Round 273: Struct embedding, shadowing, method promotion

type Base struct {
	Name string
}

func (b Base) Greet() string {
	return "Hello, " + b.Name
}

func (b *Base) GreetPointer() string {
	return "Hi, " + b.Name
}

type Container struct {
	Base
	Age int
}

type Container2 struct {
	Base
	Name string // shadows Base.Name
}

type DeepBase1 struct {
	Val int
}

func (d DeepBase1) Method() string { return "deep1" }

type DeepBase2 struct {
	Val int
}

func (d DeepBase2) Method() string { return "deep2" }

type MultiEmbed struct {
	DeepBase1
	DeepBase2
}

// PromotedField tests accessing promoted field from embedded struct
func PromotedField() string {
	c := Container{Base: Base{Name: "Alice"}, Age: 30}
	return fmt.Sprintf("name=%s,age=%d", c.Name, c.Age)
}

// PromotedMethod tests calling promoted method from embedded struct
func PromotedMethod() string {
	c := Container{Base: Base{Name: "Bob"}}
	return c.Greet()
}

// PromotedPointerMethod tests calling promoted pointer method
func PromotedPointerMethod() string {
	c := Container{Base: Base{Name: "Carol"}}
	return c.GreetPointer()
}

// ShadowedField tests that Container2.Name shadows Base.Name
func ShadowedField() string {
	c := Container2{Base: Base{Name: "inner"}, Name: "outer"}
	return fmt.Sprintf("c.Name=%s,c.Base.Name=%s", c.Name, c.Base.Name)
}

// EmbedZeroValue tests zero value of embedded struct
func EmbedZeroValue() string {
	var c Container
	return fmt.Sprintf("name=%q,age=%d", c.Name, c.Age)
}

// StructLiteral tests struct literal with embedded field
func StructLiteral() string {
	c := Container{Base: Base{Name: "X"}, Age: 5}
	return fmt.Sprintf("greet=%s", c.Greet())
}

// MultiEmbedAmbiguous tests that ambiguous method must be accessed explicitly
func MultiEmbedAmbiguous() string {
	m := MultiEmbed{DeepBase1: DeepBase1{Val: 1}, DeepBase2: DeepBase2{Val: 2}}
	return fmt.Sprintf("d1=%s,d2=%s", m.DeepBase1.Method(), m.DeepBase2.Method())
}

// MultiEmbedFieldAmbiguous tests ambiguous field access
func MultiEmbedFieldAmbiguous() string {
	m := MultiEmbed{DeepBase1: DeepBase1{Val: 10}, DeepBase2: DeepBase2{Val: 20}}
	return fmt.Sprintf("d1=%d,d2=%d", m.DeepBase1.Val, m.DeepBase2.Val)
}

// EmbeddedSliceHeader tests embedded struct with slice
type WithSlice struct {
	Items []int
}

type ContainerSlice struct {
	WithSlice
}

func EmbeddedSliceHeader() string {
	c := ContainerSlice{WithSlice: WithSlice{Items: []int{1, 2, 3}}}
	c.Items = append(c.Items, 4)
	return fmt.Sprintf("len=%d", len(c.Items))
}

// AddressOfEmbedded tests taking address of embedded field
func AddressOfEmbedded() string {
	c := Container{Base: Base{Name: "test"}}
	p := &c.Base
	p.Name = "changed"
	return fmt.Sprintf("name=%s", c.Name)
}
