package divergence_hunt261

import (
	"fmt"
)

// ============================================================================
// Round 261: Struct method sets and pointer receivers
// ============================================================================

type Counter261 struct {
	val int
}

func (c *Counter261) Inc()   { c.val++ }
func (c *Counter261) Get() int { return c.val }
func (c Counter261) Val() int  { return c.val }

// PointerReceiverOnValue tests calling pointer receiver method on value
func PointerReceiverOnValue() string {
	c := Counter261{val: 5}
	c.Inc() // Go auto-takes address
	return fmt.Sprintf("val=%d", c.Get())
}

// MethodOnPointer tests method on explicitly pointer-typed variable
func MethodOnPointer() string {
	c := &Counter261{val: 10}
	c.Inc()
	c.Inc()
	return fmt.Sprintf("get=%d,val=%d", c.Get(), c.Val())
}

// ValueReceiverNoMutation tests value receiver doesn't mutate
func ValueReceiverNoMutation() string {
	c := Counter261{val: 42}
	v := c.Val()
	c.Inc()
	return fmt.Sprintf("val_before=%d,val_after=%d", v, c.Val())
}

// SliceOfStructPointers tests methods on slice of pointers
func SliceOfStructPointers() string {
	s := []*Counter261{{val: 1}, {val: 2}, {val: 3}}
	for _, c := range s {
		c.Inc()
	}
	result := ""
	for i, c := range s {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", c.Get())
	}
	return result
}

// NilPointerMethodCall tests calling method on nil pointer with nil check
func NilPointerMethodCall() string {
	var c *Counter261
	if c != nil {
		return fmt.Sprintf("%d", c.Get())
	}
	return "nil"
}

// StructCopyVsPointer tests struct copy doesn't affect original
func StructCopyVsPointer() string {
	a := Counter261{val: 100}
	b := a // copy
	b.Inc()
	return fmt.Sprintf("a=%d,b=%d", a.Val(), b.Get())
}

// PointerCopySameUnderlying tests pointer copy shares data
func PointerCopySameUnderlying() string {
	a := &Counter261{val: 100}
	b := a // same pointer
	b.Inc()
	return fmt.Sprintf("a=%d,b=%d", a.Get(), b.Get())
}

// MethodChaining tests chained method calls on pointer receiver
func MethodChaining() string {
	c := &Counter261{val: 0}
	c.Inc()
	c.Inc()
	c.Inc()
	return fmt.Sprintf("result=%d", c.Get())
}

// StructLiteralMethodCall tests calling method on struct literal address
func StructLiteralMethodCall() string {
	c := &Counter261{val: 7}
	c.Inc()
	return fmt.Sprintf("v=%d", c.Get())
}

// ReassignPointer tests reassigning pointer changes target
func ReassignPointer() string {
	a := &Counter261{val: 1}
	b := &Counter261{val: 2}
	a = b
	a.Inc()
	return fmt.Sprintf("a=%d,b=%d", a.Get(), b.Get())
}
