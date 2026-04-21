package divergence_hunt249

import (
	"errors"
	"fmt"
)

// ============================================================================
// Custom error types for testing
// ============================================================================

// customError249 for testing
type customError249 struct {
	Code int
}

func (e *customError249) Error() string { return fmt.Sprintf("code: %d", e.Code) }

// coder249 interface for testing
type coder249 interface {
	Code() int
}

// codedError249 for testing
type codedError249 struct {
	code int
}

func (e *codedError249) Error() string { return "coded error" }
func (e *codedError249) Code() int     { return e.code }

// errorA249 for testing
type errorA249 struct{}

func (e *errorA249) Error() string { return "a" }

// errorB249 for testing
type errorB249 struct{}

func (e *errorB249) Error() string { return "b" }

// codedError249b for testing (for sentinel)
type codedError249b struct {
	code int
}

func (e *codedError249b) Error() string { return fmt.Sprintf("error %d", e.code) }

// valueError249 for testing
type valueError249 struct {
	val int
}

func (e valueError249) Error() string { return fmt.Sprintf("value error %d", e.val) }

// ============================================================================
// Round 249: Error comparison with Is/As
// ============================================================================

// ErrorsIsBasic tests basic errors.Is usage
func ErrorsIsBasic() string {
	target := errors.New("target")
	wrapped := fmt.Errorf("wrapped: %w", target)
	return fmt.Sprintf("%v", errors.Is(wrapped, target))
}

// ErrorsIsDeep tests errors.Is through multiple wraps
func ErrorsIsDeep() string {
	target := errors.New("target")
	wrap1 := fmt.Errorf("layer1: %w", target)
	wrap2 := fmt.Errorf("layer2: %w", wrap1)
	wrap3 := fmt.Errorf("layer3: %w", wrap2)
	return fmt.Sprintf("%v", errors.Is(wrap3, target))
}

// ErrorsIsNotFound tests errors.Is returns false for different errors
func ErrorsIsNotFound() string {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	return fmt.Sprintf("%v", errors.Is(err1, err2))
}

// ErrorsIsNilTarget tests errors.Is with nil target
func ErrorsIsNilTarget() string {
	err := errors.New("some error")
	return fmt.Sprintf("%v:%v", errors.Is(err, nil), errors.Is(nil, nil))
}

// ErrorsAsBasic tests basic errors.As usage
func ErrorsAsBasic() string {
	var err error = &customError249{Code: 42}
	var target *customError249
	if errors.As(err, &target) {
		return fmt.Sprintf("%d", target.Code)
	}
	return "failed"
}

// ErrorsAsWrapped tests errors.As through wrapping
func ErrorsAsWrapped() string {
	base := &customError249{Code: 100}
	wrapped := fmt.Errorf("wrapped: %w", base)
	doubleWrapped := fmt.Errorf("outer: %w", wrapped)

	var target *customError249
	if errors.As(doubleWrapped, &target) {
		return fmt.Sprintf("%d", target.Code)
	}
	return "failed"
}

// ErrorsAsInterface tests errors.As with interface target
func ErrorsAsInterface() string {
	var err error = &codedError249{code: 500}
	var target coder249
	if errors.As(err, &target) {
		return fmt.Sprintf("%d", target.Code())
	}
	return "failed"
}

// ErrorsAsNotMatching tests errors.As returns false for non-matching type
func ErrorsAsNotMatching() string {
	var err error = &errorA249{}
	var target *errorB249
	if errors.As(err, &target) {
		return "matched incorrectly"
	}
	return "no match"
}

// ErrorsIsSentinel tests errors.Is with sentinel errors
func ErrorsIsSentinel() string {
	var (
		ErrNotFound = errors.New("not found")
		ErrInvalid  = errors.New("invalid")
	)

	wrapped := fmt.Errorf("query failed: %w", ErrNotFound)
	return fmt.Sprintf("%v:%v", errors.Is(wrapped, ErrNotFound), errors.Is(wrapped, ErrInvalid))
}

// ErrorsIsAndAsTogether tests errors.Is and errors.As together
func ErrorsIsAndAsTogether() string {
	sentinel := &codedError249b{code: 404}
	wrapped := fmt.Errorf("wrapped: %w", sentinel)

	isMatch := errors.Is(wrapped, sentinel)

	var target *codedError249b
	asMatch := errors.As(wrapped, &target)

	return fmt.Sprintf("%v:%v:%d", isMatch, asMatch, target.code)
}

// ErrorsAsPointerValue tests errors.As with pointer and value types
func ErrorsAsPointerValue() string {
	var err error = valueError249{val: 42}

	// Try to match with pointer target
	var ptrTarget *valueError249
	ptrMatch := errors.As(err, &ptrTarget)

	// Try to match with value target
	var valTarget valueError249
	valMatch := errors.As(err, &valTarget)

	return fmt.Sprintf("%v:%v", ptrMatch, valMatch)
}
