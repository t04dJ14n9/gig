package divergence_hunt204

import "fmt"

// ============================================================================
// Round 204: Pointer receiver vs value receiver
// ============================================================================

type Counter204 struct {
	Value int
}

func (c Counter204) GetValue() int {
	return c.Value
}

func (c *Counter204) Increment() {
	c.Value++
}

type Container204 struct {
	Items []int
}

func (c Container204) Length() int {
	return len(c.Items)
}

func (c *Container204) Add(item int) {
	c.Items = append(c.Items, item)
}

func ValueReceiverNoMutate() string {
	c := Counter204{Value: 10}
	v := c.GetValue()
	c.Value = 20
	return fmt.Sprintf("%d:%d", v, c.Value)
}

func PointerReceiverMutates() string {
	c := Counter204{Value: 5}
	c.Increment()
	c.Increment()
	return fmt.Sprintf("%d", c.Value)
}

func ValueReceiverOnPointer() string {
	c := &Counter204{Value: 100}
	v := c.GetValue()
	return fmt.Sprintf("%d", v)
}

func PointerReceiverOnValue() string {
	c := Counter204{Value: 0}
	c.Increment()
	return fmt.Sprintf("%d", c.Value)
}

func MethodSetOnValue() string {
	c := Counter204{Value: 50}
	l1 := c.GetValue()
	c.Increment()
	l2 := c.GetValue()
	return fmt.Sprintf("%d:%d", l1, l2)
}

func ContainerValueReceiver() string {
	c := Container204{Items: []int{1, 2, 3}}
	len1 := c.Length()
	c.Add(4)
	len2 := c.Length()
	return fmt.Sprintf("%d:%d", len1, len2)
}

type MyInt204 int

func (m MyInt204) String() string { return fmt.Sprintf("%d", m) }

func InterfaceWithPointer() string {
	type Stringer interface {
		String() string
	}
	var _ Stringer = MyInt204(0)
	return "ok"
}

type Modifiable204 struct {
	X int
}

func (m Modifiable204) GetX() int {
	return m.X
}

func (m *Modifiable204) SetX(x int) {
	m.X = x
}

func MixedReceivers() string {
	m := Modifiable204{X: 10}
	get1 := m.GetX()
	m.SetX(20)
	get2 := m.GetX()
	return fmt.Sprintf("%d:%d", get1, get2)
}

func CopyInValueReceiver() string {
	c := Counter204{Value: 99}
	temp := c
	temp.Increment()
	return fmt.Sprintf("%d:%d", c.Value, temp.Value)
}
