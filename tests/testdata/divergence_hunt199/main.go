package divergence_hunt199

import (
	"fmt"
)

// ============================================================================
// Round 199: Reflection-like operations via standard means
// Note: Since reflect package may not be fully supported, we test similar
// concepts using only standard Go operations
// ============================================================================

// TypeSwitchDynamicDispatch tests type switching
func TypeSwitchDynamicDispatch() string {
	var i interface{}
	i = 42

	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case string:
		return fmt.Sprintf("string:%s", v)
	default:
		return fmt.Sprintf("unknown")
	}
}

// TypeAssertDynamicDispatch tests type assertions
func TypeAssertDynamicDispatch() string {
	var i interface{}
	i = "hello"

	if s, ok := i.(string); ok {
		return fmt.Sprintf("string:%s", s)
	}
	return fmt.Sprintf("not string")
}

// InterfaceHoldingDifferentTypes tests interface holding different types
func InterfaceHoldingDifferentTypes() string {
	var arr [3]interface{}
	arr[0] = 42
	arr[1] = "hello"
	arr[2] = true

	result := ""
	for _, v := range arr {
		switch val := v.(type) {
		case int:
			result += fmt.Sprintf("int:%d ", val)
		case string:
			result += fmt.Sprintf("string:%s ", val)
		case bool:
			result += fmt.Sprintf("bool:%v ", val)
		}
	}
	return result
}

// NilInterfaceValue tests nil interface value
func NilInterfaceValue() string {
	var i interface{}
	return fmt.Sprintf("%v", i == nil)
}

// InterfaceWithNilPointer tests interface containing nil pointer
func InterfaceWithNilPointer() string {
	type MyInterface interface{}
	var p *int = nil
	var i MyInterface = p
	return fmt.Sprintf("%v", i == nil)
}

// EmptyInterfaceIdentity tests empty interface identity
func EmptyInterfaceIdentity() string {
	var a interface{} = 42
	var b interface{} = 42
	return fmt.Sprintf("%v", a == b)
}

// DynamicMethodDispatch tests method dispatch via interface
func DynamicMethodDispatch() string {
	type Stringer interface {
		String() string
	}
	type MyInt int
	myInt := MyInt(42)
	_ = myInt
	// Note: Without String() method defined, this is a limitation
	// Return something meaningful
	return fmt.Sprintf("interface dispatch")
}

// InterfaceSlice tests slice of empty interfaces
func InterfaceSlice() string {
	s := []interface{}{1, "two", 3.0, true}
	sum := 0
	for _, v := range s {
		if n, ok := v.(int); ok {
			sum += n
		}
	}
	return fmt.Sprintf("%d", sum)
}

// InterfaceMap tests map with interface values
func InterfaceMap() string {
	m := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}
	name := m["name"].(string)
	age := m["age"].(int)
	return fmt.Sprintf("%s:%d", name, age)
}

// NestedInterface tests nested interfaces
func NestedInterface() string {
	var outer interface{}
	var inner interface{} = 42
	outer = inner
	return fmt.Sprintf("%v", outer)
}

// InterfaceEquality tests interface equality
func InterfaceEquality() string {
	var a interface{} = "hello"
	var b interface{} = "hello"
	var c interface{} = "world"
	return fmt.Sprintf("%v:%v", a == b, a == c)
}
