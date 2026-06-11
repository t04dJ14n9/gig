package divergence_hunt34

import "fmt"

// ============================================================================
// Round 34: Interface edge cases - typed nil, nil interface, empty interface,
// interface wrapping, type switch with multiple types
// ============================================================================

func TypedNilSlice() bool {
	var s []int = nil
	var x any = s
	return x != nil // typed nil is not nil interface
}

func TypedNilMap() bool {
	var m map[string]int = nil
	var x any = m
	return x != nil
}

func TypedNilPointer() bool {
	var p *int = nil
	var x any = p
	return x != nil
}

func TypedNilFunc() bool {
	var f func() = nil
	var x any = f
	return x != nil
}

func TypedNilChan() bool {
	var ch chan int = nil
	var x any = ch
	return x != nil
}

func InterfaceEqualSame() bool {
	var a any = 42
	var b any = 42
	return a == b
}

func InterfaceEqualDifferent() bool {
	var a any = 42
	var b any = "42"
	return a != b
}

func InterfaceEqualNil() bool {
	var a any = nil
	var b any = nil
	return a == b
}

func TypeSwitchMultiCase() int {
	describe := func(x any) int {
		switch v := x.(type) {
		case int8:
			return 1
		case int16:
			return 2
		case int32:
			return 3
		case int64:
			return 4
		case int:
			return 5
		default:
			_ = v
			return -1
		}
	}
	return describe(int8(1)) + describe(int16(2)) + describe(int(5))
}

func TypeSwitchUintFamily() int {
	describe := func(x any) int {
		switch v := x.(type) {
		case uint8:
			return 1
		case uint16:
			return 2
		case uint32:
			return 3
		case uint64:
			return 4
		default:
			_ = v
			return -1
		}
	}
	return describe(uint8(1)) + describe(uint16(2)) + describe(uint32(3))
}

func TypeSwitchFloatFamily() int {
	describe := func(x any) int {
		switch v := x.(type) {
		case float32:
			return 32
		case float64:
			return 64
		default:
			_ = v
			return -1
		}
	}
	return describe(float32(1)) + describe(float64(2))
}

func AssertToSliceType() int {
	var x any = []int{1, 2, 3}
	if v, ok := x.([]int); ok {
		return v[0] + v[1] + v[2]
	}
	return -1
}

func AssertToMapType() int {
	var x any = map[string]int{"a": 1}
	if v, ok := x.(map[string]int); ok {
		return v["a"]
	}
	return -1
}

func AssertToFuncType() int {
	fn := func(x int) int { return x * 2 }
	var x any = fn
	if v, ok := x.(func(int) int); ok {
		return v(21)
	}
	return -1
}

func FmtTypedNil() string {
	var s []int = nil
	var x any = s
	return fmt.Sprintf("%v", x)
}

func FmtNilInterface() string {
	var x any
	return fmt.Sprintf("%v", x)
}

func InterfaceSliceOfTypeSwitch() int {
	items := []any{42, "hello", true, 3.14}
	count := 0
	for _, item := range items {
		switch item.(type) {
		case int:
			count += 1
		case string:
			count += 10
		case bool:
			count += 100
		case float64:
			count += 1000
		}
	}
	return count
}

func NestedTypeSwitch() int {
	outer := func(x any) int {
		switch v := x.(type) {
		case int:
			return v * 2
		case string:
			return len(v)
		default:
			return -1
		}
	}
	return outer(21) + outer("hi")
}
