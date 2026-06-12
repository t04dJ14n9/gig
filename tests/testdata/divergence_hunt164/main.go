package divergence_hunt164

import "fmt"

// ============================================================================
// Round 164: Method sets and pointer vs value receivers
// ============================================================================

type Counter struct {
	count int
}

func (c Counter) Value() int {
	return c.count
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Reset() {
	c.count = 0
}

type ReadWrite interface {
	Read() string
	Write(string)
}

type Document struct {
	content string
}

func (d Document) Read() string {
	return d.content
}

func (d *Document) Write(s string) {
	d.content = s
}

type CounterInterface interface {
	Value() int
	Increment()
}

// ValueReceiverNoModify tests value receiver that does not modify
func ValueReceiverNoModify() string {
	c := Counter{count: 5}
	v1 := c.Value()
	c.count = 10
	v2 := c.Value()
	return fmt.Sprintf("v1=%d,v2=%d", v1, v2)
}

// PointerReceiverModifies tests pointer receiver that modifies
func PointerReceiverModifies() string {
	c := Counter{count: 5}
	c.Increment()
	c.Increment()
	return fmt.Sprintf("count=%d", c.count)
}

// PointerReceiverViaValue tests calling pointer method on value
func PointerReceiverViaValue() string {
	c := Counter{count: 5}
	c.Increment()
	return fmt.Sprintf("count=%d", c.count)
}

// ValueReceiverOnPointer tests calling value method on pointer
func ValueReceiverOnPointer() string {
	c := &Counter{count: 7}
	v := c.Value()
	return fmt.Sprintf("value=%d", v)
}

// MethodSetDifference tests difference between value and pointer method sets
func MethodSetDifference() string {
	type S struct{ x int }
	var v S
	var p *S
	_ = v
	_ = p
	// Value receiver methods work on both
	// Pointer receiver methods only work on pointer
	return "method set understood"
}

// InterfaceSatisfactionValue tests interface satisfaction with value receiver
func InterfaceSatisfactionValue() string {
	c := Counter{count: 10}
	var i interface{ Value() int } = c
	return fmt.Sprintf("value=%d", i.Value())
}

// InterfaceSatisfactionPointer tests interface satisfaction with pointer
func InterfaceSatisfactionPointer() string {
	c := &Counter{count: 10}
	var i CounterInterface = c
	i.Increment()
	return fmt.Sprintf("value=%d", i.Value())
}

// NilPointerReceiverWithValueMethod tests nil pointer with value method
func NilPointerReceiverWithValueMethod() string {
	var c *Counter
	v := c.Value()
	return fmt.Sprintf("value=%d", v)
}

// NilPointerReceiverWithPointerMethod tests nil pointer with pointer method
func NilPointerReceiverWithPointerMethod() string {
	var c *Counter
	c.Increment()
	return fmt.Sprintf("count=%d", c.count)
}

// AssignmentToInterface tests assignment to interface types
func AssignmentToInterface() string {
	type Stringer interface{ String() string }
	type MyStruct struct{ name string }
	fn := func(m MyStruct) string { return m.name }
	_ = fn
	// var s Stringer = MyStruct{"test"} // requires value receiver
	return "interface assignment"
}

// SliceOfStructsWithMethods tests slice of structs with methods
func SliceOfStructsWithMethods() string {
	counters := []Counter{{count: 1}, {count: 2}, {count: 3}}
	sum := 0
	for _, c := range counters {
		sum += c.Value()
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapOfStructsWithMethods tests map of structs with methods
func MapOfStructsWithMethods() string {
	counters := map[string]Counter{
		"a": {count: 10},
		"b": {count: 20},
	}
	sum := 0
	for _, c := range counters {
		sum += c.Value()
	}
	return fmt.Sprintf("sum=%d", sum)
}

// EmbeddedFieldMethodPromotion tests embedded field method promotion
func EmbeddedFieldMethodPromotion() string {
	type Inner struct {
		value int
	}
	type Outer struct {
		Inner
	}
	outer := Outer{Inner: Inner{value: 42}}
	return fmt.Sprintf("value=%d", outer.value)
}
