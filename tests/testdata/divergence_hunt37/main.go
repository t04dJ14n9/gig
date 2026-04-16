package divergence_hunt37

import "fmt"

// ============================================================================
// Round 37: Panic/recover with different value types
// ============================================================================

func PanicIntRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(int); ok {
				result = v
			}
		}
	}()
	panic(42)
}

func PanicStringRecover() (result string) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(string); ok {
				result = v
			}
		}
	}()
	panic("error!")
}

func PanicFloatRecover() (result float64) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(float64); ok {
				result = v
			}
		}
	}()
	panic(3.14)
}

func PanicBoolRecover() (result bool) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(bool); ok {
				result = v
			}
		}
	}()
	panic(true)
}

func PanicInt32Recover() (result int32) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(int32); ok {
				result = v
			}
		}
	}()
	panic(int32(100))
}

func PanicUint8Recover() (result uint8) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(uint8); ok {
				result = v
			}
		}
	}()
	panic(uint8(255))
}

func PanicSliceRecover() (result int) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.([]int); ok {
				result = v[0]
			}
		}
	}()
	panic([]int{99, 100})
}

func PanicMapRecover() (result string) {
	defer func() {
		if r := recover(); r != nil {
			if v, ok := r.(map[string]int); ok {
				result = fmt.Sprintf("%d", v["key"])
			}
		}
	}()
	panic(map[string]int{"key": 42})
}

func RecoverInMultipleDefers() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result += r.(int)
		}
	}()
	defer func() {
		// This runs first (LIFO), but recover() returns nil because no panic yet
		if r := recover(); r != nil {
			result += r.(int) * 10
		}
	}()
	panic(5)
}

func RecoverTypeSwitch() (result string) {
	defer func() {
		switch r := recover().(type) {
		case int:
			result = fmt.Sprintf("int:%d", r)
		case string:
			result = fmt.Sprintf("string:%s", r)
		case float64:
			result = fmt.Sprintf("float:%v", r)
		case bool:
			result = fmt.Sprintf("bool:%v", r)
		default:
			result = "unknown"
		}
	}()
	panic(true)
}

func PanicInNestedFunc() (result int) {
	defer func() {
		if r := recover(); r != nil {
			result = r.(int)
		}
	}()
	func() {
		func() {
			panic(99)
		}()
	}()
	return 0
}

func PanicWithNilInterface() (result int) {
	defer func() {
		r := recover()
		if r == nil {
			result = 1 // Go 1.21+: panic(nil) returns PanicNilError
		} else {
			result = 2
		}
	}()
	panic(nil)
}
