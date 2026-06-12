package divergence_hunt71

import "fmt"

// ============================================================================
// Round 71: Interface satisfaction and method sets
// ============================================================================

type Describer interface {
	Describe() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Describe() string {
	return fmt.Sprintf("%s(%d)", p.Name, p.Age)
}

func InterfaceSatisfaction() string {
	var d Describer = Person{Name: "Alice", Age: 30}
	return d.Describe()
}

type Animal struct {
	Species string
}

func (a *Animal) Describe() string {
	return a.Species
}

func PointerReceiverInterface() string {
	a := Animal{Species: "Cat"}
	var d Describer = &a
	return d.Describe()
}

func InterfaceNilCheck() bool {
	var d Describer
	return d == nil
}

func InterfaceTypeSwitch() int {
	var x any = 42
	switch v := x.(type) {
	case int:
		return v
	case string:
		return len(v)
	default:
		return -1
	}
}

func InterfaceAssertionOk() int {
	var x any = "hello"
	if v, ok := x.(string); ok {
		return len(v)
	}
	return -1
}

func InterfaceAssertionFail() int {
	var x any = "hello"
	if _, ok := x.(int); ok {
		return 1
	}
	return 0
}

func EmptyInterface() int {
	var x any = 42
	return x.(int)
}

func InterfaceSlice() int {
	items := []any{1, "hello", true}
	count := 0
	for _, item := range items {
		switch item.(type) {
		case int:
			count += 1
		case string:
			count += 2
		case bool:
			count += 4
		}
	}
	return count
}

func InterfaceMap() int {
	m := map[string]any{
		"name":  "Alice",
		"age":   30,
		"admin": true,
	}
	return len(m)
}

func InterfaceMethodCall() string {
	var d Describer = Person{Name: "Bob", Age: 25}
	return d.Describe()
}

type Wrapper struct {
	val any
}

func (w Wrapper) Unwrap() any {
	return w.val
}

func InterfaceAsField() int {
	w := Wrapper{val: 42}
	v := w.Unwrap()
	return v.(int)
}
