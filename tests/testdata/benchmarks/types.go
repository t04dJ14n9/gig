package benchmarks

// ============================================================================
// Type System Operations
// ============================================================================

type Counter struct {
	value int
}

func (c *Counter) Add(x int) { c.value = c.value + x }
func (c *Counter) Get() int  { return c.value }

func StructMethod() int {
	c := &Counter{}
	for i := 0; i < 100; i++ {
		c.Add(i)
	}
	return c.Get()
}

type Adder interface{ Add(int) }

func Interface() int {
	var a Adder = &Counter{}
	for i := 0; i < 100; i++ {
		a.Add(i)
	}
	return a.(*Counter).value
}

type Any interface{}

func TypeAssertion() int {
	var x Any = 42
	sum := 0
	for i := 0; i < 100; i++ {
		if v, ok := x.(int); ok {
			sum = sum + v
		}
	}
	return sum
}

func process(x Any) int {
	switch v := x.(type) {
	case int:
		return v * 2
	case string:
		return len(v)
	default:
		return 0
	}
}

func TypeSwitch() int {
	values := []Any{1, "hello", 2.5, 3, "world", 4.0}
	sum := 0
	for i := 0; i < 100; i++ {
		for _, v := range values {
			sum = sum + process(v)
		}
	}
	return sum
}

func SliceInterface() int {
	arr := make([]*Counter, 10)
	for i := 0; i < 10; i++ {
		arr[i] = &Counter{}
	}
	sum := 0
	for i := 0; i < 100; i++ {
		for _, c := range arr {
			c.Add(i)
			sum = sum + c.value
		}
	}
	return sum
}

type Point struct{ X, Y int }

func CompositeLiteral() int {
	points := []Point{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}}
	sum := 0
	for i := 0; i < 100; i++ {
		for _, p := range points {
			sum = sum + p.X + p.Y
		}
	}
	return sum
}
