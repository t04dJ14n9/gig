package divergence_hunt202

import "fmt"

// ============================================================================
// Round 202: Struct embedding and shadowing
// ============================================================================

type Base202 struct {
	X int
	Y int
}

func (b Base202) String() string {
	return fmt.Sprintf("Base(%d,%d)", b.X, b.Y)
}

type Inner202 struct {
	Value int
}

type Outer202 struct {
	Inner202
	Value int
}

func BasicEmbedding() string {
	o := Outer202{Inner202{10}, 20}
	return fmt.Sprintf("%d:%d", o.Inner202.Value, o.Value)
}

func EmbeddedFieldAccess() string {
	o := Outer202{Inner202{5}, 10}
	return fmt.Sprintf("%d", o.Value)
}

func ShadowedField() string {
	o := Outer202{Inner202{100}, 200}
	innerVal := o.Inner202.Value
	outerVal := o.Value
	return fmt.Sprintf("%d:%d", innerVal, outerVal)
}

type DoubleEmbed struct {
	Base202
	Z int
}

func DoubleEmbedded() string {
	d := DoubleEmbed{Base202{1, 2}, 3}
	return fmt.Sprintf("%d:%d:%d", d.X, d.Y, d.Z)
}

type A202 struct {
	Name string
}

type B202 struct {
	A202
	Name string
}

func DeepShadowing() string {
	b := B202{A202{"inner"}, "outer"}
	return fmt.Sprintf("%s:%s", b.A202.Name, b.Name)
}

type EmbeddedPtr struct {
	*Base202
	Z int
}

func EmbeddedPointer() string {
	b := &Base202{X: 10, Y: 20}
	e := EmbeddedPtr{Base202: b, Z: 30}
	return fmt.Sprintf("%d:%d:%d", e.X, e.Y, e.Z)
}

type NamedEmbed struct {
	Base202
}

func NamedEmbeddedAccess() string {
	n := NamedEmbed{Base202{X: 5, Y: 10}}
	return fmt.Sprintf("%d", n.X)
}

func EmbeddedMethodAccess() string {
	o := Outer202{Inner202{42}, 0}
	return fmt.Sprintf("Inner=%d", o.Value)
}
