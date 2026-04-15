package divergence_hunt17

import "fmt"

// ============================================================================
// Round 17: Interface composition, method patterns, polymorphism-like patterns
// ============================================================================

// InterfaceComposition tests interface composition
func InterfaceComposition() int {
	type Reader interface{ Read() int }
	type Writer interface{ Write(int) }
	type ReadWriter interface {
		Reader
		Writer
	}
	_ = ReadWriter(nil)
	return 42
}

// InterfaceEmpty tests empty interface usage
func InterfaceEmpty() int {
	var x any = 42
	switch v := x.(type) {
	case int: return v
	case string: return len(v)
	default: return -1
	}
}

// InterfaceSlice tests slice of empty interface
func InterfaceSlice() int {
	s := []any{1, "hello", true, 3.14}
	return len(s)
}

// InterfaceMap tests map with interface value
func InterfaceMap() int {
	m := map[string]any{
		"int":    42,
		"string": "hello",
		"bool":   true,
	}
	return len(m)
}

// StructMethodOnPointer tests struct method via pointer
func StructMethodOnPointer() int {
	type Counter struct{ n int }
	inc := func(c *Counter) { c.n++ }
	c := &Counter{n: 10}
	inc(c)
	inc(c)
	return c.n
}

// StructMethodOnValue tests struct method via value
func StructMethodOnValue() int {
	type Rect struct{ W, H int }
	area := func(r Rect) int { return r.W * r.H }
	r := Rect{W: 3, H: 4}
	return area(r)
}

// MethodChain tests method chaining pattern
func MethodChain() int {
	type Builder struct{ val int }
	add := func(b *Builder, n int) *Builder { b.val += n; return b }
	b := &Builder{val: 0}
	add(add(add(b, 1), 2), 3)
	return b.val
}

// PolymorphismPattern tests polymorphism via interface
func PolymorphismPattern() int {
	type Shape interface{ Area() int }
	type Rect struct{ W, H int }
	type Circle struct{ R int }
	rectArea := func(r Rect) int { return r.W * r.H }
	circleArea := func(c Circle) int { return c.R * c.R * 3 }
	r := Rect{W: 3, H: 4}
	c := Circle{R: 2}
	return rectArea(r) + circleArea(c)
}

// NullableInterface tests nil interface
func NullableInterface() bool {
	var x any = nil
	return x == nil
}

// InterfaceTypeAssertion tests type assertion on interface
func InterfaceTypeAssertion() int {
	var x any = []int{1, 2, 3}
	if v, ok := x.([]int); ok {
		return v[0] + v[1] + v[2]
	}
	return -1
}

// EmbeddedStructAccess tests embedded struct access
func EmbeddedStructAccess() int {
	type Base struct{ X int }
	type Derived struct {
		Base
		Y int
	}
	d := Derived{Base: Base{X: 10}, Y: 20}
	return d.X + d.Y
}

// NestedStructAccess tests nested struct access
func NestedStructAccess() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 42}}
	return o.I.X
}

// StructSliceMethods tests struct methods in slice
func StructSliceMethods() int {
	type Item struct{ Name string; Val int }
	sum := func(items []Item) int {
		total := 0
		for _, it := range items { total += it.Val }
		return total
	}
	items := []Item{{"a", 1}, {"b", 2}, {"c", 3}}
	return sum(items)
}

// FmtInterface tests fmt with interface
func FmtInterface() string {
	var x any = 42
	return fmt.Sprintf("%v", x)
}

// FmtNilInterface tests fmt with nil interface
func FmtNilInterface() string {
	var x any
	return fmt.Sprintf("%v", x)
}

// StructComparison tests struct comparison
func StructComparison() bool {
	type P struct{ X, Y int }
	return P{1, 2} == P{1, 2} && P{1, 2} != P{1, 3}
}

// InterfaceEquality tests interface equality
func InterfaceEquality() bool {
	var a any = 42
	var b any = 42
	return a == b
}

// InterfaceInequality tests interface inequality
func InterfaceInequality() bool {
	var a any = 42
	var b any = "42"
	return a != b
}
