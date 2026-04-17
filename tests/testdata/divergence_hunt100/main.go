package divergence_hunt100

import "fmt"

// ============================================================================
// Round 100: Panic/recover patterns - nested, typed, nil
// ============================================================================

func BasicPanicRecover() string {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("caught: %v", r)
		}
	}()
	panic("hello")
}

func PanicInt() int {
	defer func() {
		if r := recover(); r != nil {
			_ = r.(int)
		}
	}()
	panic(42)
	return 0
}

func PanicStruct() string {
	type Err struct{ Code int; Msg string }
	defer func() {
		if r := recover(); r != nil {
			e := r.(Err)
			_ = fmt.Sprintf("%d:%s", e.Code, e.Msg)
		}
	}()
	panic(Err{Code: 404, Msg: "not found"})
}

func NestedPanicRecover() string {
	inner := func() {
		defer func() {
			recover()
		}()
		panic("inner")
	}
	outer := func() string {
		defer func() {
			recover()
		}()
		inner()
		return "ok"
	}
	return outer()
}

func PanicInDefer() string {
	result := "before"
	defer func() {
		recover()
		result = "recovered"
	}()
	defer func() {
		panic("defer panic")
	}()
	return result
}

func NoPanicReturn() int {
	defer func() {
		recover()
	}()
	return 42
}

func RecoverWithoutPanic() string {
	result := "none"
	defer func() {
		if r := recover(); r != nil {
			result = "caught"
		}
	}()
	return result
}

func PanicNilInterface() string {
	defer func() {
		if r := recover(); r != nil {
			_ = fmt.Sprintf("%v", r)
		}
	}()
	var err error
	panic(err)
}

func PanicSliceBounds() string {
	defer func() {
		recover()
	}()
	s := []int{1, 2, 3}
	_ = s[10]
	return "unreachable"
}

func PanicNilMap() string {
	defer func() {
		recover()
	}()
	var m map[string]int
	m["key"] = 1
	return "unreachable"
}

func PanicNilPointer() string {
	defer func() {
		recover()
	}()
	var p *int
	*p = 42
	return "unreachable"
}
