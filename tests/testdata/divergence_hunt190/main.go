package divergence_hunt190

import (
	"fmt"
)

// ============================================================================
// Round 190: Type switches with multiple cases
// ============================================================================

func BasicTypeSwitch() string {
	var i interface{} = 42
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case string:
		return fmt.Sprintf("string:%s", v)
	default:
		return fmt.Sprintf("unknown")
	}
}

func TypeSwitchString() string {
	var i interface{} = "hello"
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case string:
		return fmt.Sprintf("string:%s", v)
	default:
		return fmt.Sprintf("unknown")
	}
}

func TypeSwitchMultipleTypes() string {
	test := func(i interface{}) string {
		switch v := i.(type) {
		case int, int8, int16, int32, int64:
			return fmt.Sprintf("integer:%v", v)
		case uint, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("unsigned:%v", v)
		case float32, float64:
			return fmt.Sprintf("float:%v", v)
		default:
			return fmt.Sprintf("other")
		}
	}
	return fmt.Sprintf("%s:%s:%s", test(42), test(uint(10)), test(3.14))
}

func TypeSwitchBool() string {
	var i interface{} = true
	switch v := i.(type) {
	case bool:
		return fmt.Sprintf("bool:%v", v)
	default:
		return fmt.Sprintf("other")
	}
}

func TypeSwitchSlice() string {
	var i interface{} = []int{1, 2, 3}
	switch v := i.(type) {
	case []int:
		return fmt.Sprintf("[]int:%d", len(v))
	case []string:
		return fmt.Sprintf("[]string:%d", len(v))
	default:
		return fmt.Sprintf("other")
	}
}

func TypeSwitchMap() string {
	var i interface{} = map[string]int{"a": 1, "b": 2}
	switch v := i.(type) {
	case map[string]int:
		return fmt.Sprintf("map[string]int:%d", len(v))
	case map[int]string:
		return fmt.Sprintf("map[int]string:%d", len(v))
	default:
		return fmt.Sprintf("other")
	}
}

func TypeSwitchNil() string {
	var i interface{} = nil
	switch v := i.(type) {
	case nil:
		return fmt.Sprintf("nil:%v", v)
	case int:
		return fmt.Sprintf("int:%d", v)
	default:
		return fmt.Sprintf("other:%v", v)
	}
}

func TypeSwitchPointer() string {
	type Point struct{ X, Y int }
	var i interface{} = &Point{X: 10, Y: 20}
	switch v := i.(type) {
	case *Point:
		return fmt.Sprintf("*Point:%d,%d", v.X, v.Y)
	default:
		return fmt.Sprintf("other")
	}
}

func TypeSwitchInterface() string {
	type Stringer interface {
		String() string
	}
	type MyType struct{}
	var i interface{} = MyType{}
	switch v := i.(type) {
	case Stringer:
		return fmt.Sprintf("Stringer:%v", v)
	case MyType:
		return fmt.Sprintf("MyType:%v", v)
	default:
		return fmt.Sprintf("other")
	}
}

func TypeSwitchDefault() string {
	var i interface{} = struct{}{}
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int:%d", v)
	case string:
		return fmt.Sprintf("string:%s", v)
	default:
		return fmt.Sprintf("default:%v", v)
	}
}

func TypeSwitchFunction() string {
	var i interface{} = func(x int) int { return x * 2 }
	switch v := i.(type) {
	case func(int) int:
		result := v(5)
		return fmt.Sprintf("func(int)int:%d", result)
	default:
		return fmt.Sprintf("other")
	}
}
