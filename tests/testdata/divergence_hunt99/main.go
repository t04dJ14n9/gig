package divergence_hunt99

import "fmt"

// ============================================================================
// Round 99: Struct embedding with method overrides
// ============================================================================

type Base struct {
	val int
}

func (b Base) Get() int    { return b.val }
func (b *Base) Set(v int)  { b.val = v }

type Derived struct {
	Base
	extra string
}

func (d *Derived) Set(v int) {
	d.val = v * 10
}

type DeepDerived struct {
	Derived
	deep string
}

func (dd *DeepDerived) Get() int {
	return dd.val * 100
}

func OverrideMethod() int {
	d := &Derived{Base: Base{val: 5}, extra: "test"}
	d.Set(3)
	return d.Get()
}

func PromotedMethod() int {
	d := &Derived{Base: Base{val: 5}, extra: "test"}
	return d.Get()
}

func DirectBaseMethod() int {
	d := &Derived{Base: Base{val: 5}, extra: "test"}
	d.Base.Set(3)
	return d.Get()
}

func DeepEmbedding() int {
	dd := &DeepDerived{Derived: Derived{Base: Base{val: 2}, extra: "mid"}, deep: "bottom"}
	return dd.Get()
}

func DeepSetViaBase() int {
	dd := &DeepDerived{Derived: Derived{Base: Base{val: 2}, extra: "mid"}, deep: "bottom"}
	dd.Base.Set(3)
	return dd.Get()
}

type A struct{ x int }
type B struct{ A }
type C struct{ B }

func (a A) getX() int { return a.x }

func TripleEmbedding() int {
	c := C{B: B{A: A{x: 42}}}
	return c.getX()
}

func EmbeddedLiteral() int {
	d := Derived{
		Base:  Base{val: 10},
		extra: "hi",
	}
	return d.Get()
}

func OverrideVsPromote() string {
	d := &Derived{Base: Base{val: 1}, extra: "test"}
	d.Set(5)       // calls Derived.Set (override)
	baseVal := d.Get()
	d.Base.Set(5)  // calls Base.Set (direct)
	baseDirectVal := d.Get()
	return fmt.Sprintf("%d:%d", baseVal, baseDirectVal)
}
