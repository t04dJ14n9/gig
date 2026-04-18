package known_issues

import (
	"errors"
	"fmt"
)

// ============================================================================
// Known issues test data.
//
// This file contains test cases for known interpreter bugs/limitations.
// When a bug is fixed, remove its function from here and promote to a
// passing test (e.g. in divergence_hunt_test.go or correctness_test.go).
//
// As of 2026-04-17: ALL previously known issues have been fixed!
// Including errors.As with struct pointer target (via gigStructWrapper
// implementing error interface + GigErrorsAs type name matching).
// ============================================================================

// CustomError is a custom error type used to test errors.As.
type CustomError struct {
	Code int
	Msg  string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// ErrorAsStructPointer tests errors.As with a struct pointer target.
// The interpreter wraps arguments in interface{} when calling native functions,
// which causes the runtime type introspection inside errors.As to fail — it
// cannot match the concrete *CustomError type hidden behind the interpreter's
// value wrapper to the **CustomError target.
func ErrorAsStructPointer() any {
	err := &CustomError{Code: 404, Msg: "not found"}
	var ce *CustomError
	if errors.As(err, &ce) {
		return ce.Code
	}
	return -1
}
