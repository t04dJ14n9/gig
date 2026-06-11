package divergence_hunt250

import (
	"errors"
	"fmt"
)

// codedError for testing
type codedError250 struct {
	Code int
}

func (e *codedError250) Error() string { return fmt.Sprintf("code: %d", e.Code) }

// ============================================================================
// Round 250: Multiple error handling
// ============================================================================

// JoinTwoErrors tests joining two errors
func JoinTwoErrors() string {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	joined := errors.Join(err1, err2)
	return fmt.Sprintf("%v", joined != nil)
}

// JoinMultipleErrors tests joining multiple errors
func JoinMultipleErrors() string {
	joined := errors.Join(
		errors.New("first"),
		errors.New("second"),
		errors.New("third"),
	)
	return fmt.Sprintf("%v", joined != nil)
}

// JoinWithNilErrors tests joining with some nil errors
func JoinWithNilErrors() string {
	err1 := errors.New("error 1")
	joined := errors.Join(err1, nil, nil)
	return fmt.Sprintf("%v", joined != nil)
}

// JoinAllNilErrors tests joining only nil errors
func JoinAllNilErrors() string {
	joined := errors.Join(nil, nil, nil)
	return fmt.Sprintf("%v", joined == nil)
}

// JoinErrorsIs tests errors.Is on joined errors
func JoinErrorsIs() string {
	target := errors.New("target")
	err1 := errors.New("other")
	err2 := fmt.Errorf("wrapped: %w", target)
	joined := errors.Join(err1, err2)
	return fmt.Sprintf("%v", errors.Is(joined, target))
}

// JoinErrorsAs tests errors.As on joined errors
func JoinErrorsAs() string {
	err1 := errors.New("plain")
	err2 := &codedError250{Code: 42}
	joined := errors.Join(err1, err2)

	var target *codedError250
	return fmt.Sprintf("%v", errors.As(joined, &target))
}

// CollectErrorsAccumulates tests collecting multiple errors
func CollectErrorsAccumulates() string {
	var errs []error

	for i := 0; i < 3; i++ {
		if err := doWork(i); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Sprintf("errors:%d", len(errs))
	}
	return "success"
}

func doWork(n int) error {
	if n%2 == 0 {
		return fmt.Errorf("work %d failed", n)
	}
	return nil
}

// FirstNonNilError tests getting first non-nil error
func FirstNonNilError() string {
	errs := []error{nil, nil, errors.New("first actual"), errors.New("second")}
	var first error
	for _, err := range errs {
		if err != nil {
			first = err
			break
		}
	}
	return fmt.Sprintf("%v", first != nil)
}

// CombineErrorsInLoop tests combining errors from loop
func CombineErrorsInLoop() string {
	var result error
	for i := 0; i < 3; i++ {
		err := fmt.Errorf("step %d failed", i)
		result = errors.Join(result, err)
	}
	return fmt.Sprintf("%v", result != nil)
}

// ErrorSliceToJoined tests converting error slice to joined error
func ErrorSliceToJoined() string {
	errs := make([]error, 0, 3)
	for i := 1; i <= 3; i++ {
		errs = append(errs, fmt.Errorf("error %d", i))
	}
	joined := errors.Join(errs...)
	return fmt.Sprintf("%v", joined != nil)
}

// ProcessMultipleErrors tests processing multiple errors
func ProcessMultipleErrors() string {
	errs := []error{
		errors.New("err1"),
		nil,
		errors.New("err3"),
	}

	count := 0
	for _, err := range errs {
		if err != nil {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}
