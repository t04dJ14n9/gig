package divergence_hunt189

import (
	"fmt"
)

// ============================================================================
// Round 189: Type assertions with ok pattern
// ============================================================================

func BasicTypeAssertion() string {
	var i interface{} = 42
	n, ok := i.(int)
	return fmt.Sprintf("%d:%v", n, ok)
}

func TypeAssertionFail() string {
	var i interface{} = "hello"
	n, ok := i.(int)
	return fmt.Sprintf("%d:%v", n, ok)
}

func TypeAssertionString() string {
	var i interface{} = "world"
	s, ok := i.(string)
	return fmt.Sprintf("%s:%v", s, ok)
}

func TypeAssertionFloat64() string {
	var i interface{} = 3.14
	f, ok := i.(float64)
	return fmt.Sprintf("%.2f:%v", f, ok)
}

func TypeAssertionBool() string {
	var i interface{} = true
	b, ok := i.(bool)
	return fmt.Sprintf("%v:%v", b, ok)
}

func TypeAssertionSlice() string {
	var i interface{} = []int{1, 2, 3}
	s, ok := i.([]int)
	return fmt.Sprintf("%d:%v", len(s), ok)
}

func TypeAssertionMap() string {
	var i interface{} = map[string]int{"a": 1}
	m, ok := i.(map[string]int)
	return fmt.Sprintf("%d:%v", len(m), ok)
}

func TypeAssertionNil() string {
	var i interface{} = nil
	n, ok := i.(int)
	return fmt.Sprintf("%d:%v", n, ok)
}

func TypeAssertionInterface() string {
	type Stringer interface {
		String() string
	}
	var i interface{} = "test"
	s, ok := i.(Stringer)
	return fmt.Sprintf("%v:%v", s != nil, ok)
}

func TypeAssertionPointer() string {
	type MyStruct struct{ X int }
	var i interface{} = &MyStruct{X: 42}
	p, ok := i.(*MyStruct)
	if ok {
		return fmt.Sprintf("%d:%v", p.X, ok)
	}
	return fmt.Sprintf("nil:%v", ok)
}

func TypeAssertionChained() string {
	var i interface{} = 100
	result := ""
	if n, ok := i.(int); ok {
		result = fmt.Sprintf("int:%d", n)
	} else if s, ok := i.(string); ok {
		result = fmt.Sprintf("string:%s", s)
	} else {
		result = "unknown"
	}
	return result
}
