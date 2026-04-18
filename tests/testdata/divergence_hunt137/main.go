package divergence_hunt137

import "fmt"

// ============================================================================
// Round 137: Method value expressions and method calls
// ============================================================================

type Counter struct {
	n int
}

func (c *Counter) Inc() {
	c.n++
}

func (c *Counter) Get() int {
	return c.n
}

type Valuer struct {
	val int
}

func (v Valuer) Value() int {
	return v.val
}

func (v *Valuer) PtrValue() int {
	return v.val
}

func MethodValueExpr() string {
	c := &Counter{}
	f := c.Inc
	f()
	f()
	f()
	return fmt.Sprintf("n=%d", c.Get())
}

func MethodCallDirect() string {
	c := &Counter{}
	c.Inc()
	c.Inc()
	return fmt.Sprintf("n=%d", c.Get())
}

func MethodValueReceiver() string {
	v := Valuer{val: 42}
	return fmt.Sprintf("val=%d", v.Value())
}

func MethodPtrReceiver() string {
	v := &Valuer{val: 99}
	return fmt.Sprintf("val=%d", v.PtrValue())
}

func MethodOnLiteral() string {
	return fmt.Sprintf("val=%d", Valuer{val: 7}.Value())
}

type Adder struct{ X int }

func (a Adder) Add(y int) int {
	return a.X + y
}

func MethodWithArgs() string {
	a := Adder{X: 10}
	return fmt.Sprintf("10+5=%d", a.Add(5))
}

type Stringer struct {
	s string
}

func (s Stringer) String() string {
	return fmt.Sprintf("str:%s", s.s)
}

func MethodStringer() string {
	s := Stringer{s: "hello"}
	return s.String()
}

type Stack struct {
	items []int
}

func (s *Stack) Push(v int) {
	s.items = append(s.items, v)
}

func (s *Stack) Pop() int {
	n := len(s.items)
	v := s.items[n-1]
	s.items = s.items[:n-1]
	return v
}

func MethodStackPushPop() string {
	s := &Stack{}
	s.Push(10)
	s.Push(20)
	s.Push(30)
	v := s.Pop()
	return fmt.Sprintf("popped=%d-remaining=%d", v, len(s.items))
}
