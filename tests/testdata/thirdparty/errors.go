package thirdparty

import (
	"errors"
	"fmt"
)

// ErrorsNew tests errors.New.
func ErrorsNew() int {
	err := errors.New("test error")
	if err != nil && err.Error() == "test error" {
		return 1
	}
	return 0
}

// ErrorsIs tests errors.Is.
func ErrorsIs() int {
	err1 := errors.New("error")
	err2 := fmt.Errorf("wrapped: %w", err1)
	if errors.Is(err2, err1) {
		return 1
	}
	return 0
}

// ErrorsAs tests errors.As.
func ErrorsAs() int {
	err := errors.New("test")
	var target error
	if errors.As(err, &target) {
		return 1
	}
	return 0
}

// ErrorsJoin tests errors.Join.
func ErrorsJoin() string {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	return errors.Join(err1, err2).Error()
}
