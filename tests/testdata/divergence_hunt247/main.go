package divergence_hunt247

import (
	"errors"
	"fmt"
)

// errorString for testing
type errorString247 struct {
	s string
}

func (e *errorString247) Error() string {
	return e.s
}

// myError247 for testing
type myError247 string

func (e myError247) Error() string { return string(e) }

// multiMethodError247 for testing
type multiMethodError247 struct {
	msg  string
	code int
}

func (e *multiMethodError247) Error() string { return e.msg }
func (e *multiMethodError247) Code() int     { return e.code }

// ============================================================================
// Round 247: Error interface satisfaction
// ============================================================================

// ErrorInterfaceBasic tests basic error interface satisfaction
func ErrorInterfaceBasic() string {
	var err error = errors.New("test error")
	return fmt.Sprintf("%v:%T", err != nil, err)
}

// ErrorInterfaceNil tests nil error interface
func ErrorInterfaceNil() string {
	var err error = nil
	return fmt.Sprintf("%v", err == nil)
}

// ErrorInterfacePointerNil tests typed nil in error interface
func ErrorInterfacePointerNil() string {
	var p *errorString247 = nil
	var err error = p
	return fmt.Sprintf("%v:%v", err == nil, p == nil)
}

// ErrorInterfaceAssignability tests error interface assignability
func ErrorInterfaceAssignability() string {
	var err error = myError247("custom")
	return fmt.Sprintf("%v", err.Error())
}

// ErrorInterfaceMethodSet tests error interface method set
func ErrorInterfaceMethodSet() string {
	var err error = &multiMethodError247{msg: "test", code: 500}
	_, ok := err.(*multiMethodError247)
	return fmt.Sprintf("%v", ok)
}

// ErrorInterfaceAssertion tests type assertion on error
func ErrorInterfaceAssertion() string {
	type customError struct {
		error
		code int
	}
	base := errors.New("base")
	wrapped := &customError{error: base, code: 404}
	var err error = wrapped

	if ce, ok := err.(*customError); ok {
		return fmt.Sprintf("%d", ce.code)
	}
	return "not custom"
}

// ErrorInterfaceEmptyString tests error with empty string
func ErrorInterfaceEmptyString() string {
	err := errors.New("")
	return fmt.Sprintf("%v:%v", err != nil, err.Error() == "")
}

// ErrorInterfaceComparison tests error interface comparison
func ErrorInterfaceComparison() string {
	err1 := errors.New("same")
	err2 := errors.New("same")
	var e1 error = err1
	var e2 error = err2
	return fmt.Sprintf("%v:%v", e1 == e1, e1 == e2)
}

// ErrorInterfaceSlice tests slice of error interfaces
func ErrorInterfaceSlice() string {
	errs := []error{
		errors.New("first"),
		errors.New("second"),
		nil,
	}
	return fmt.Sprintf("%d:%v", len(errs), errs[2] == nil)
}

// ErrorInterfaceMap tests map with error keys/values
func ErrorInterfaceMap() string {
	errMap := map[string]error{
		"a": errors.New("error a"),
		"b": nil,
	}
	return fmt.Sprintf("%v:%v", errMap["a"] != nil, errMap["b"] == nil)
}

// ErrorInterfaceFunction tests function returning error
func ErrorInterfaceFunction() string {
	mayError := func(fail bool) error {
		if fail {
			return errors.New("failed")
		}
		return nil
	}
	err1 := mayError(true)
	err2 := mayError(false)
	return fmt.Sprintf("%v:%v", err1 != nil, err2 == nil)
}
