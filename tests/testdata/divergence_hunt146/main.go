package divergence_hunt146

import (
	"errors"
	"fmt"
)

// ============================================================================
// Round 146: Complex type assertions and error interface patterns
// ============================================================================

type CustomErr struct {
	Code int
	Msg  string
}

func (e *CustomErr) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

func ErrorAsStructPointer() string {
	err := &CustomErr{Code: 404, Msg: "not found"}
	var ce *CustomErr
	if errors.As(err, &ce) {
		return fmt.Sprintf("code=%d", ce.Code)
	}
	return "no match"
}

type WrapErr struct {
	inner error
	msg   string
}

func (e *WrapErr) Error() string {
	return e.msg
}

func (e *WrapErr) Unwrap() error {
	return e.inner
}

func ErrorChainUnwrap() string {
	e1 := &CustomErr{Code: 500, Msg: "internal"}
	e2 := &WrapErr{inner: e1, msg: "wrapped"}
	return fmt.Sprintf("inner=%s-outer=%s", e2.Unwrap().Error(), e2.Error())
}

func ErrorNilInterface() string {
	var err error
	if err == nil {
		return "nil-error"
	}
	return "has-error"
}

func ErrorTypeAssertion() string {
	var err error = &CustomErr{Code: 403, Msg: "forbidden"}
	if ce, ok := err.(*CustomErr); ok {
		return fmt.Sprintf("code=%d", ce.Code)
	}
	return "not-custom"
}

func ErrorInterfaceAssertion() string {
	var err error = &CustomErr{Code: 401, Msg: "unauthorized"}
	_ = err.Error()
	return fmt.Sprintf("is-error=%t", err != nil)
}

type TimeoutErr struct {
	Duration int
}

func (e *TimeoutErr) Error() string {
	return fmt.Sprintf("timeout after %dms", e.Duration)
}

func (e *TimeoutErr) Timeout() bool {
	return true
}

func ErrorSpecificMethod() string {
	var err error = &TimeoutErr{Duration: 5000}
	if te, ok := err.(interface{ Timeout() bool }); ok {
		return fmt.Sprintf("timeout=%t-dur=%d", te.Timeout(), err.(*TimeoutErr).Duration)
	}
	return "no-timeout"
}

func ErrorMultiWrap() string {
	e1 := &CustomErr{Code: 100, Msg: "base"}
	e2 := &WrapErr{inner: e1, msg: "level2"}
	e3 := &WrapErr{inner: e2, msg: "level3"}
	inner := e3.Unwrap()
	inner2 := inner.(*WrapErr).Unwrap()
	return fmt.Sprintf("base=%s", inner2.Error())
}
