package divergence_hunt127

import "fmt"

// ============================================================================
// Round 127: Type switches and interface type assertions
// ============================================================================

func TypeSwitchBasic() string {
	var x interface{} = 42
	switch v := x.(type) {
	case int:
		return fmt.Sprintf("int-%d", v)
	case string:
		return fmt.Sprintf("string-%s", v)
	default:
		return "unknown"
	}
}

func TypeSwitchString() string {
	var x interface{} = "hello"
	switch v := x.(type) {
	case int:
		return fmt.Sprintf("int-%d", v)
	case string:
		return fmt.Sprintf("string-%s", v)
	default:
		return "unknown"
	}
}

func TypeSwitchNil() string {
	var x interface{}
	switch x.(type) {
	case int:
		return "int"
	case nil:
		return "nil"
	default:
		return "other"
	}
}

func TypeAssertionOk() string {
	var x interface{} = "hello"
	v, ok := x.(string)
	return fmt.Sprintf("val=%s-ok=%t", v, ok)
}

func TypeAssertionFail() string {
	var x interface{} = 42
	v, ok := x.(string)
	return fmt.Sprintf("val=%s-ok=%t", v, ok)
}

func TypeAssertionPanicFree() string {
	var x interface{} = []int{1, 2, 3}
	// Safe assertion with ok
	v, ok := x.(string)
	if !ok {
		return fmt.Sprintf("not-string-type=%T", v)
	}
	return v
}

func TypeSwitchMultiCase() string {
	var x interface{} = 3.14
	switch x.(type) {
	case int, float64:
		return "numeric"
	case string:
		return "string"
	default:
		return "other"
	}
}

func TypeSwitchStruct() string {
	type Point struct{ X, Y int }
	var x interface{} = Point{X: 1, Y: 2}
	switch v := x.(type) {
	case Point:
		return fmt.Sprintf("point-%d-%d", v.X, v.Y)
	default:
		return "other"
	}
}

func TypeAssertionChain() string {
	var x interface{} = "test"
	if s, ok := x.(string); ok {
		return fmt.Sprintf("string-%s", s)
	}
	if i, ok := x.(int); ok {
		return fmt.Sprintf("int-%d", i)
	}
	return "none"
}
