package divergence_hunt280

import (
	"fmt"
)

// ============================================================================
// Round 280: Methods on defined types — pointer vs value receiver, method on basic type

type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	return float64(c)*9/5 + 32
}

func (c *Celsius) Add(delta Celsius) {
	*c += delta
}

type IntSlice []int

func (s IntSlice) Sum() int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total
}

func (s *IntSlice) Append(vals ...int) {
	*s = append(*s, vals...)
}

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c Counter) Value() int {
	return c.count
}

type MyString string

func (s MyString) Shout() string {
	return string(s) + "!!!"
}

func (s *MyString) Append(suffix string) {
	*s += MyString(suffix)
}

type Named struct {
	Name string
}

func (n Named) Greet() string {
	return "Hello, " + n.Name
}

func (n *Named) Rename(newName string) {
	n.Name = newName
}

// MethodOnBasicType tests method on defined basic type
func MethodOnBasicType() string {
	var c Celsius = 100
	return fmt.Sprintf("f=%.1f", c.ToFahrenheit())
}

// PointerReceiverOnBasicType tests pointer receiver on defined basic type
func PointerReceiverOnBasicType() string {
	var c Celsius = 0
	c.Add(100)
	return fmt.Sprintf("c=%.0f", float64(c))
}

// MethodOnSlice tests method on defined slice type
func MethodOnSlice() string {
	s := IntSlice{1, 2, 3, 4, 5}
	return fmt.Sprintf("sum=%d", s.Sum())
}

// PointerReceiverOnSlice tests pointer receiver on defined slice type
func PointerReceiverOnSlice() string {
	s := IntSlice{1, 2}
	s.Append(3, 4)
	return fmt.Sprintf("s=%v,sum=%d", s, s.Sum())
}

// ValueReceiverDoesNotModify tests value receiver doesn't modify original
func ValueReceiverDoesNotModify() string {
	c := Counter{count: 5}
	v := c.Value()
	c.Increment()
	return fmt.Sprintf("v=%d,after=%d", v, c.Value())
}

// PointerReceiverModifies tests pointer receiver modifies original
func PointerReceiverModifies() string {
	c := Counter{}
	c.Increment()
	c.Increment()
	c.Increment()
	return fmt.Sprintf("count=%d", c.Value())
}

// MethodOnStringType tests method on defined string type
func MethodOnStringType() string {
	s := MyString("hello")
	return s.Shout()
}

// PointerMethodOnStringType tests pointer method on defined string type
func PointerMethodOnStringType() string {
	s := MyString("hello")
	s.Append(" world")
	return string(s)
}

// AddressableValueReceiver tests calling pointer method on addressable value
func AddressableValueReceiver() string {
	n := Named{Name: "Alice"}
	n.Rename("Bob")
	return n.Greet()
}

// ValueMethodOnPointer tests calling value method on pointer
func ValueMethodOnPointer() string {
	n := &Named{Name: "Carol"}
	return n.Greet()
}

// MethodChain tests chaining methods
func MethodChain() string {
	var c Celsius = 0
	c.Add(50)
	c.Add(50)
	return fmt.Sprintf("f=%.1f", c.ToFahrenheit())
}

// NilPointerReceiver tests calling method on nil pointer of struct type
func NilPointerReceiver() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	var c *Counter
	c.Increment() // should panic
	return "no_panic"
}
