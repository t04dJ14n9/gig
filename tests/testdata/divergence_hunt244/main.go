package divergence_hunt244

import (
	"fmt"
)

// ============================================================================
// Round 244: Panic recovery patterns
// ============================================================================

// BasicPanicRecover tests basic panic and recover
func BasicPanicRecover() string {
	result := "no panic"
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("caught:%v", r)
			}
		}()
		panic("boom")
	}()
	return result
}

// PanicWithIntRecover tests panic with integer value
func PanicWithIntRecover() string {
	result := 0
	func() {
		defer func() {
			if r := recover(); r != nil {
				if v, ok := r.(int); ok {
					result = v
				}
			}
		}()
		panic(42)
	}()
	return fmt.Sprintf("%d", result)
}

// PanicWithStructRecover tests panic with struct value
func PanicWithStructRecover() string {
	type Err struct {
		Code int
		Msg  string
	}
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(Err); ok {
					result = fmt.Sprintf("%d:%s", e.Code, e.Msg)
				}
			}
		}()
		panic(Err{Code: 404, Msg: "not found"})
	}()
	return result
}

// PanicWithPointerRecover tests panic with pointer value
func PanicWithPointerRecover() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				if s, ok := r.(*string); ok {
					result = *s
				}
			}
		}()
		msg := "pointer panic"
		panic(&msg)
	}()
	return result
}

// RecoverReturnsNil tests that recover returns nil without panic
func RecoverReturnsNil() string {
	result := "has value"
	func() {
		defer func() {
			if r := recover(); r == nil {
				result = "nil"
			}
		}()
		// no panic
	}()
	return result
}

// MultipleRecoverCalls tests multiple recover calls in same defer
func MultipleRecoverCalls() string {
	result := ""
	func() {
		defer func() {
			r1 := recover()
			r2 := recover()
			result = fmt.Sprintf("%v:%v", r1 != nil, r2 == nil)
		}()
		panic("test")
	}()
	return result
}

// PanicInDeferRecovered tests panic in defer recovered by outer defer
func PanicInDeferRecovered() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("caught:%v", r)
			}
		}()
		defer func() {
			panic("defer panic")
		}()
		panic("original")
	}()
	return result
}

// RepanicAfterRecover tests re-panicking after recover
func RepanicAfterRecover() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("caught:%v", r)
			}
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Sprintf("re-%v", r))
				}
			}()
			panic("original")
		}()
	}()
	return result
}

// PanicWithNilInterface tests panic with nil interface
func PanicWithNilInterface() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = "not nil"
			} else {
				result = "nil"
			}
		}()
		var err error
		panic(err)
	}()
	return result
}

// PanicNilLiteral tests panic with nil literal
func PanicNilLiteral() string {
	result := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = fmt.Sprintf("%v", r)
			} else {
				result = "nil"
			}
		}()
		panic(nil)
	}()
	return result
}
