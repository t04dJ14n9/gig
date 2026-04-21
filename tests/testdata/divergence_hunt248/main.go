package divergence_hunt248

import (
	"errors"
	"fmt"
)

// ============================================================================
// Custom error types for testing
// ============================================================================

// appError248 for testing
type appError248 struct {
	Code    int
	Message string
}

func (e *appError248) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// myError248 for testing
type myError248 string

func (e myError248) Error() string { return string(e) }

// wrappedError248 for testing
type wrappedError248 struct {
	msg string
	err error
}

func (e *wrappedError248) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}
func (e *wrappedError248) Unwrap() error {
	return e.err
}

// comparableError248 for testing
type comparableError248 struct {
	code int
}

func (e *comparableError248) Error() string {
	return fmt.Sprintf("error code %d", e.code)
}
func (e *comparableError248) Is(target error) bool {
	if t, ok := target.(*comparableError248); ok {
		return e.code == t.code
	}
	return false
}

// richError248 for testing
type richError248 struct {
	Code int
	Data string
}

func (e *richError248) Error() string {
	return fmt.Sprintf("rich error: %d", e.Code)
}

// ptrError248 for testing
type ptrError248 struct{ msg string }

func (e *ptrError248) Error() string { return e.msg }

// valError248 for testing
type valError248 struct{ msg string }

func (e valError248) Error() string { return e.msg }

// nilError248 for testing
type nilError248 struct{ msg string }

func (e *nilError248) Error() string {
	if e == nil {
		return "nil"
	}
	return e.msg
}

// multiError248 for testing
type multiError248 struct {
	errors []error
}

func (e *multiError248) Error() string {
	return fmt.Sprintf("multiple errors (%d)", len(e.errors))
}

// statefulError248 for testing
type statefulError248 struct {
	msg   string
	count int
}

func (e *statefulError248) Error() string {
	e.count++
	return fmt.Sprintf("%s (count: %d)", e.msg, e.count)
}

// ============================================================================
// Round 248: Custom error types
// ============================================================================

// CustomErrorStruct tests custom error as struct
func CustomErrorStruct() string {
	err := &appError248{Code: 404, Message: "not found"}
	return err.Error()
}

// CustomErrorString tests custom error as string type
func CustomErrorString() string {
	var err error = myError248("my custom error")
	return err.Error()
}

// CustomErrorWithUnwrap tests custom error with Unwrap method
func CustomErrorWithUnwrap() string {
	base := errors.New("base error")
	wrapped := &wrappedError248{msg: "wrapped", err: base}
	return fmt.Sprintf("%v", errors.Unwrap(wrapped) == base)
}

// CustomErrorIs tests custom error with Is method
func CustomErrorIs() string {
	err1 := &comparableError248{code: 404}
	err2 := &comparableError248{code: 404}
	err3 := &comparableError248{code: 500}
	return fmt.Sprintf("%v:%v", errors.Is(err1, err2), errors.Is(err1, err3))
}

// CustomErrorAs tests custom error with As method
func CustomErrorAs() string {
	var err error = &richError248{Code: 42, Data: "test"}
	var target *richError248
	if errors.As(err, &target) {
		return fmt.Sprintf("%d:%s", target.Code, target.Data)
	}
	return "as failed"
}

// baseError248 and specificError248 for hierarchy testing
type baseError248 struct{ msg string }

type specificError248 struct {
	baseError248
	code int
}

func (e *baseError248) Error() string { return e.msg }

// CustomErrorHierarchy tests error type hierarchy
func CustomErrorHierarchy() string {
	var err error = &specificError248{
		baseError248: baseError248{msg: "specific"},
		code:         123,
	}
	return fmt.Sprintf("%v:%v", err.Error(), errors.As(err, new(*specificError248)))
}

// CustomErrorPointerValue tests pointer vs value receiver error
func CustomErrorPointerValue() string {
	var pErr error = &ptrError248{msg: "ptr"}
	var vErr error = valError248{msg: "val"}
	return fmt.Sprintf("%v:%v", pErr.Error(), vErr.Error())
}

// CustomErrorNilPointer tests nil pointer custom error
func CustomErrorNilPointer() string {
	var err error = (*nilError248)(nil)
	return fmt.Sprintf("%v:%v", err != nil, err.Error())
}

// CustomErrorComposite tests composite error type
func CustomErrorComposite() string {
	mErr := &multiError248{
		errors: []error{
			errors.New("err1"),
			errors.New("err2"),
		},
	}
	return mErr.Error()
}

// CustomErrorState tests error with mutable state
func CustomErrorState() string {
	err := &statefulError248{msg: "stateful"}
	_ = err.Error()
	_ = err.Error()
	return err.Error()
}
