package divergence_hunt82

import "fmt"

// ============================================================================
// Round 82: Error interface edge cases
// ============================================================================

type MyErr struct {
	Code int
	Msg  string
}

func (e MyErr) Error() string {
	return fmt.Sprintf("code %d: %s", e.Code, e.Msg)
}

func ErrorAsInterface() string {
	var err error = MyErr{Code: 404, Msg: "not found"}
	return err.Error()
}

func ErrorNilCheck() bool {
	var err error
	return err == nil
}

func ErrorTypeAssertion() int {
	var err error = MyErr{Code: 500, Msg: "internal"}
	if e, ok := err.(MyErr); ok {
		return e.Code
	}
	return -1
}

func ErrorPointerAssertion() int {
	var err error = &MyErr{Code: 503, Msg: "unavailable"}
	if e, ok := err.(*MyErr); ok {
		return e.Code
	}
	return -1
}

func ErrorPointerDoesNotMatchValue() int {
	var err error = MyErr{Code: 500, Msg: "internal"}
	if _, ok := err.(*MyErr); ok {
		return 1 // pointer assertion on value type should fail
	}
	return 0
}

func ErrorValueDoesNotMatchPointer() int {
	var err error = &MyErr{Code: 503, Msg: "unavailable"}
	if _, ok := err.(MyErr); ok {
		return 1 // value assertion on pointer should succeed (pointer receiver but value assertion)
	}
	return 0
}

func FmtErrorf() string {
	err := fmt.Errorf("value %d is invalid", 42)
	return err.Error()
}

func ErrorInMultiReturn() (int, error) {
	div := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}
	result, err := div(10, 2)
	if err != nil {
		return -1, err
	}
	return result, nil
}

func ErrorInMultiReturnFail() string {
	div := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}
	_, err := div(10, 0)
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func ErrorSlice() int {
	errs := []error{
		fmt.Errorf("err1"),
		fmt.Errorf("err2"),
	}
	return len(errs)
}
