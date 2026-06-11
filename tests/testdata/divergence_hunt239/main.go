package divergence_hunt239

import "fmt"

// ============================================================================
// Round 239: fmt.Stringer interface
// ============================================================================

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("Person(%s, %d)", p.Name, p.Age)
}

type Point struct {
	X, Y int
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

type Counter struct {
	Value int
}

func (c Counter) String() string {
	return fmt.Sprintf("Counter: %d", c.Value)
}

func StringerWithStruct() string {
	p := Person{Name: "Alice", Age: 30}
	return p.String()
}

func StringerWithPointerReceiver() string {
	type Item struct {
		Name  string
		Price float64
	}
	i := Item{Name: "Book", Price: 19.99}
	return fmt.Sprintf("%v", i)
}

func StringerInSlice() string {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	result := ""
	for _, p := range people {
		result += p.String() + ";"
	}
	return result
}

func StringerWithFmtSprintf() string {
	p := Point{X: 10, Y: 20}
	return fmt.Sprintf("%s", p)
}

func StringerWithFmtV() string {
	c := Counter{Value: 42}
	return fmt.Sprintf("%v", c)
}

func StringerInterfaceAssertion() string {
	var s fmt.Stringer = Person{Name: "Bob", Age: 25}
	return s.String()
}

func StringerNilReceiver() string {
	type Maybe struct {
		Value *int
	}
	m := Maybe{Value: nil}
	return fmt.Sprintf("%v", m)
}

func StringerNestedCall() string {
	type Address struct {
		City    string
		Country string
	}
	a := Address{City: "Tokyo", Country: "Japan"}
	return fmt.Sprintf("Address: %v", a)
}

func StringerWithSpecialChars() string {
	type Message struct {
		Text string
	}
	m := Message{Text: "Hello\nWorld\t!"}
	return fmt.Sprintf("%v", m)
}
