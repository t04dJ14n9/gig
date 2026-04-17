package divergence_hunt94

import "fmt"

// ============================================================================
// Round 94: Complex type assertions and type switches
// ============================================================================

func TypeSwitchBasic() string {
	identify := func(v interface{}) string {
		switch v := v.(type) {
		case int:
			return fmt.Sprintf("int:%d", v)
		case string:
			return fmt.Sprintf("string:%s", v)
		case bool:
			return fmt.Sprintf("bool:%v", v)
		default:
			return "unknown"
		}
	}
	return identify(42) + ";" + identify("hi") + ";" + identify(true)
}

func TypeSwitchMultiple() string {
	categorize := func(v interface{}) string {
		switch v.(type) {
		case int, int8, int16:
			return "small_int"
		case int32, int64:
			return "big_int"
		case float32, float64:
			return "float"
		default:
			return "other"
		}
	}
	return categorize(int(5)) + ";" + categorize(int64(5)) + ";" + categorize(3.14)
}

func TypeAssertionCommaOk() string {
	items := []interface{}{42, "hello", 3.14}
	result := ""
	for _, item := range items {
		if s, ok := item.(string); ok {
			result += s
		} else if i, ok := item.(int); ok {
			result += fmt.Sprintf("%d", i)
		}
	}
	return result
}

func TypeAssertionPanicSafe() string {
	safe := func(v interface{}) string {
		defer func() {
			recover()
		}()
		return v.(string)
	}
	return safe(42)
}

func NestedTypeSwitch() string {
	process := func(v interface{}) string {
		switch v := v.(type) {
		case []interface{}:
			sum := ""
			for _, item := range v {
				if s, ok := item.(string); ok {
					sum += s
				}
			}
			return sum
		case map[string]interface{}:
			return fmt.Sprintf("%d keys", len(v))
		default:
			return "other"
		}
	}
	return process([]interface{}{"a", "b", "c"})
}

func TypeSwitchWithNil() string {
	check := func(v interface{}) string {
		switch v.(type) {
		case nil:
			return "nil"
		default:
			return "non-nil"
		}
	}
	return check(nil) + ";" + check(0)
}

func AssertToInterface() string {
	var x interface{} = "hello"
	var y interface{} = x
	if s, ok := y.(string); ok {
		return s
	}
	return "fail"
}

func AssertSliceTypes() string {
	items := []interface{}{
		[]int{1, 2, 3},
		[]string{"a", "b"},
	}
	result := ""
	for _, item := range items {
		if s, ok := item.([]string); ok {
			result += fmt.Sprintf("%d", len(s))
		}
	}
	return result
}

func AssertMapType() string {
	var v interface{} = map[string]int{"a": 1, "b": 2}
	if m, ok := v.(map[string]int); ok {
		return fmt.Sprintf("%d", m["a"])
	}
	return "fail"
}

func TypeSwitchFallthrough() string {
	// Note: fallthrough not allowed in type switches, using if-else
	check := func(v interface{}) string {
		if _, ok := v.(int); ok {
			return "integer"
		} else if _, ok := v.(string); ok {
			return "string"
		}
		return "other"
	}
	return check(42)
}
