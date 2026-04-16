package divergence_hunt62

import "fmt"

// ============================================================================
// Round 62: Typed nil vs nil interface edge cases
// ============================================================================

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

// Typed nil pointer satisfies error interface but is not nil
func TypedNilError() string {
	var err *MyError // nil pointer
	if err != nil {
		return "not nil"
	}
	return "nil"
}

// Typed nil as interface is NOT nil
func TypedNilInterface() string {
	var err error // interface type
	var myErr *MyError // nil pointer
	err = myErr // typed nil assigned to interface
	if err == nil {
		return "nil"
	}
	return "not nil"
}

// Nil interface check
func NilInterfaceCheck() string {
	var err error
	if err == nil {
		return "nil"
	}
	return "not nil"
}

func NilSliceVsEmptySlice() string {
	var s1 []int        // nil slice
	s2 := []int{}       // empty slice
	r1 := "nil"
	r2 := "empty"
	if s1 != nil {
		r1 = "not nil"
	}
	if s2 != nil {
		r2 = "not nil"
	}
	return r1 + ":" + r2
}

func NilMapVsEmptyMap() string {
	var m1 map[string]int // nil map
	m2 := map[string]int{} // empty map
	r1 := "nil"
	r2 := "empty"
	if m1 != nil {
		r1 = "not nil"
	}
	if m2 != nil {
		r2 = "not nil"
	}
	return r1 + ":" + r2
}

func NilChanVsMakeChan() string {
	var ch1 chan int // nil chan
	ch2 := make(chan int) // made chan
	r1 := "nil"
	r2 := "not nil"
	if ch1 != nil {
		r1 = "not nil"
	}
	if ch2 == nil {
		r2 = "nil"
	}
	return r1 + ":" + r2
}

func NilFuncCheck() string {
	var f func() int // nil func
	if f == nil {
		return "nil"
	}
	return "not nil"
}

func NilPointerCheck() string {
	var p *int
	if p == nil {
		return "nil"
	}
	return "not nil"
}

func TypeAssertNil() string {
	var x any = nil
	if x == nil {
		return "nil"
	}
	return "not nil"
}

func TypeAssertTypedNil() string {
	var p *int
	var x any = p
	if _, ok := x.(*int); ok {
		return "ok"
	}
	return "fail"
}

func FmtTypedNil() string {
	var p *MyError
	return fmt.Sprintf("%v", p)
}

func InterfaceMethodOnTypedNil() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	var p *MyError // typed nil
	var e error = p
	_ = e.Error() // should panic - nil pointer dereference
	return "ok"
}
