package main

import "fmt"

type Counter struct{ n int }

func (c Counter) Value() int {
	return c.n + 1
}

func (c *Counter) Inc(delta int) int {
	c.n += delta
	return c.n
}

// Result checks method expressions/values and loop-variable captures in closures.
func Result() string {
	c := Counter{n: 1}
	v := c.Value
	mv := Counter.Value
	p := (&c).Inc

	valViaValue := v()
	valViaMethodExpr := mv(c)
	valViaPtrMethodExpr := p(2)

	closures := make([]func() int, 0, 4)
	for i := 0; i < 4; i++ {
		closures = append(closures, func() int { return i })
	}
	capture := fmt.Sprintf("%d,%d,%d,%d", closures[0](), closures[1](), closures[2](), closures[3]())

	return fmt.Sprintf("%d:%d:%d:%s", valViaValue, valViaMethodExpr, valViaPtrMethodExpr, capture)
}
