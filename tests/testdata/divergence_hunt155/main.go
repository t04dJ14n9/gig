package divergence_hunt155

import (
	"fmt"
	"io"
)

// ============================================================================
// Round 155: Interface composition and type sets
// ============================================================================

// ReaderWriter interface composition
type ReaderWriter interface {
	io.Reader
	io.Writer
}

// Stringer interface
type Stringer interface {
	String() string
}

// Combined interface
type Combined interface {
	Stringer
	fmt.Stringer
}

// InterfaceSatisfaction tests interface satisfaction with composite interface
func InterfaceSatisfaction() string {
	// A type that implements io.Writer
	type MyWriter struct{}
	// Implement Write method
	_ = func(w *MyWriter, p []byte) (n int, err error) { return len(p), nil }
	return "satisfied"
}

// EmptyInterface tests empty interface behavior
func EmptyInterface() string {
	var x interface{} = 42
	var y interface{} = "hello"
	var z interface{} = []int{1, 2, 3}
	return fmt.Sprintf("int=%T-str=%T-slice=%T", x, y, z)
}

// InterfaceNilComparison tests nil interface comparison
func InterfaceNilComparison() string {
	var i interface{}
	isNil := i == nil
	var p *int
	i = p
	isTypedNil := i != nil
	return fmt.Sprintf("nil=%t-typednil=%t", isNil, isTypedNil)
}

// InterfaceTypeSwitch tests type switch with multiple types
func InterfaceTypeSwitch() string {
	check := func(v interface{}) string {
		switch val := v.(type) {
		case int:
			return fmt.Sprintf("int:%d", val)
		case string:
			return fmt.Sprintf("str:%s", val)
		case bool:
			return fmt.Sprintf("bool:%t", val)
		case []int:
			return fmt.Sprintf("slice:%d", len(val))
		case map[string]int:
			return fmt.Sprintf("map:%d", len(val))
		default:
			return fmt.Sprintf("other:%T", val)
		}
	}
	return check(42) + "-" + check("hi") + "-" + check([]int{1, 2})
}

// InterfaceTypeAssertion tests type assertion with ok pattern
func InterfaceTypeAssertion() string {
	var i interface{} = 42
	n, ok1 := i.(int)
	s, ok2 := i.(string)
	return fmt.Sprintf("n=%d-ok1=%t-s=%q-ok2=%t", n, ok1, s, ok2)
}

// InterfaceSlice tests slice of interfaces
func InterfaceSlice() string {
	s := []interface{}{
		1,
		"two",
		3.0,
		true,
	}
	result := ""
	for _, v := range s {
		switch val := v.(type) {
		case int:
			result += fmt.Sprintf("i%d", val)
		case string:
			result += fmt.Sprintf("s%s", val)
		case float64:
			result += fmt.Sprintf("f%.0f", val)
		case bool:
			result += fmt.Sprintf("b%t", val)
		}
	}
	return result
}

// InterfaceMap tests map with interface values
func InterfaceMap() string {
	m := map[string]interface{}{
		"age":   30,
		"name":  "Alice",
		"score": 95.5,
	}
	return fmt.Sprintf("age=%v-name=%v", m["age"], m["name"])
}

// InterfaceEmbedding tests interface embedding
func InterfaceEmbedding() string {
	// Test that embedded interfaces work
	type ReadWriter interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
	}
	var _ ReadWriter = (interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
	})(nil)
	return "embedded"
}
