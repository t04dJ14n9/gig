package divergence_hunt91

import "fmt"

// ============================================================================
// Round 91: Struct methods - value vs pointer receivers, method sets
// ============================================================================

type Counter struct {
	count int
}

func (c Counter) Value() int {
	return c.count
}

func (c *Counter) Inc() {
	c.count++
}

func (c *Counter) IncBy(n int) {
	c.count += n
}

func (c Counter) Double() Counter {
	return Counter{count: c.count * 2}
}

type NamedCounter struct {
	Counter
	name string
}

func (nc *NamedCounter) Reset() {
	nc.count = 0
}

func ValueReceiver() int {
	c := Counter{count: 5}
	d := c.Double()
	c.Inc()
	return c.count*100 + d.count
}

func PointerReceiver() int {
	c := &Counter{count: 5}
	c.Inc()
	c.IncBy(3)
	return c.Value()
}

func EmbeddedMethod() int {
	nc := &NamedCounter{name: "test"}
	nc.Inc()
	nc.Inc()
	nc.IncBy(7)
	return nc.Value()
}

func EmbeddedReset() int {
	nc := &NamedCounter{name: "test"}
	nc.Inc()
	nc.Inc()
	nc.Reset()
	nc.IncBy(5)
	return nc.Value()
}

func MethodOnLiteral() int {
	c := Counter{count: 10}.Double()
	return c.count
}

func MethodChain() int {
	c := &Counter{count: 1}
	c.Inc()
	d := c.Double()
	c.Inc()
	return c.count*100 + d.count
}

func ValueCopySemantics() int {
	c := Counter{count: 10}
	d := c
	d.Inc()
	return c.count*100 + d.count
}

func PointerSharedSemantics() int {
	c := &Counter{count: 10}
	d := c
	d.Inc()
	return c.count*100 + d.count
}

func ReceiverOnStructLiteral() string {
	c := &Counter{count: 5}
	return fmt.Sprintf("%d", c.Value())
}

func EmbeddedPromoteMethod() int {
	nc := NamedCounter{name: "test", Counter: Counter{count: 42}}
	return nc.Value()
}
