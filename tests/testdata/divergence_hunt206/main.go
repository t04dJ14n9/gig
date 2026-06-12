package divergence_hunt206

import "fmt"

// ============================================================================
// Round 206: Type assertions with interfaces
// ============================================================================

type Shape206 interface {
	Area() int
}

type Rectangle206 struct {
	Width, Height int
}

func (r Rectangle206) Area() int {
	return r.Width * r.Height
}

type Circle206 struct {
	Radius int
}

func (c Circle206) Area() int {
	return 3 * c.Radius * c.Radius
}

func BasicTypeAssertion() string {
	var s Shape206 = Rectangle206{3, 4}
	r, ok := s.(Rectangle206)
	return fmt.Sprintf("ok:%v,area:%d", ok, r.Area())
}

func TypeAssertionFail() string {
	var s Shape206 = Rectangle206{1, 2}
	_, ok := s.(Circle206)
	return fmt.Sprintf("ok:%v", ok)
}

func TypeAssertionPanic() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panicked"
		}
	}()
	var s Shape206 = Rectangle206{1, 1}
	_ = s.(Circle206)
	return "no panic"
}

func TypeSwitchWithInterface() string {
	shapes := []Shape206{
		Rectangle206{2, 3},
		Circle206{5},
	}
	result := ""
	for _, s := range shapes {
		switch v := s.(type) {
		case Rectangle206:
			result += fmt.Sprintf("R%d", v.Area())
		case Circle206:
			result += fmt.Sprintf("C%d", v.Area())
		}
	}
	return result
}

type Stringer206 interface {
	String() string
}

type MyString206 string

func (m MyString206) String() string { return string(m) }

func NestedInterfaceAssertion() string {
	var i interface{} = MyString206("hello")
	if s, ok := i.(Stringer206); ok {
		return s.String()
	}
	return "fail"
}

func InterfaceToConcrete() string {
	var s Shape206 = Rectangle206{5, 5}
	var i interface{} = s
	r, ok := i.(Rectangle206)
	return fmt.Sprintf("ok:%v,area:%d", ok, r.Area())
}

func PointerTypeAssertion() string {
	var s Shape206 = &Rectangle206{2, 4}
	r, ok := s.(*Rectangle206)
	return fmt.Sprintf("ok:%v", ok && r != nil)
}

func NilInterfaceAssertion() string {
	var s Shape206
	_, ok := s.(Rectangle206)
	return fmt.Sprintf("ok:%v", ok)
}

func MultipleTypeAssertions() string {
	var items []interface{} = []interface{}{
		Rectangle206{2, 3},
		Circle206{4},
		"string",
		42,
	}
	count := 0
	for _, item := range items {
		if _, ok := item.(Shape206); ok {
			count++
		}
	}
	return fmt.Sprintf("shapes:%d", count)
}

func AssertToInterface() string {
	var s Shape206 = Rectangle206{1, 1}
	var i interface{} = s
	_, ok := i.(Shape206)
	return fmt.Sprintf("ok:%v", ok)
}
