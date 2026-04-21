package divergence_hunt246

import (
	"errors"
	"fmt"
)

// ============================================================================
// Round 246: Error wrapping and unwrapping
// ============================================================================

// BasicErrorWrap tests basic error wrapping with fmt.Errorf
func BasicErrorWrap() string {
	base := errors.New("base error")
	wrapped := fmt.Errorf("wrapped: %w", base)
	return fmt.Sprintf("%v", errors.Unwrap(wrapped) == base)
}

// DoubleErrorWrap tests double wrapping an error
func DoubleErrorWrap() string {
	base := errors.New("root")
	wrap1 := fmt.Errorf("layer1: %w", base)
	wrap2 := fmt.Errorf("layer2: %w", wrap1)
	unwrapped := errors.Unwrap(wrap2)
	return fmt.Sprintf("%v:%v", unwrapped == wrap1, errors.Unwrap(unwrapped) == base)
}

// TripleErrorWrap tests triple wrapping an error
func TripleErrorWrap() string {
	base := errors.New("root")
	wrap1 := fmt.Errorf("layer1: %w", base)
	wrap2 := fmt.Errorf("layer2: %w", wrap1)
	wrap3 := fmt.Errorf("layer3: %w", wrap2)

	u1 := errors.Unwrap(wrap3)
	u2 := errors.Unwrap(u1)
	u3 := errors.Unwrap(u2)
	return fmt.Sprintf("%v:%v:%v", u1 == wrap2, u2 == wrap1, u3 == base)
}

// ErrorUnwrapChain tests unwrapping through a chain
func ErrorUnwrapChain() string {
	base := errors.New("base")
	current := base
	for i := 0; i < 3; i++ {
		current = fmt.Errorf("layer%d: %w", i, current)
	}

	count := 0
	for current != nil {
		current = errors.Unwrap(current)
		count++
	}
	return fmt.Sprintf("%d", count)
}

// ErrorWrapNil tests wrapping a nil error
func ErrorWrapNil() string {
	var err error = nil
	wrapped := fmt.Errorf("wrapped: %w", err)
	return fmt.Sprintf("%v:%v", wrapped != nil, errors.Unwrap(wrapped) == nil)
}

// ErrorUnwrapNonWrapped tests unwrapping a non-wrapped error
func ErrorUnwrapNonWrapped() string {
	err := errors.New("plain")
	unwrapped := errors.Unwrap(err)
	return fmt.Sprintf("%v", unwrapped == nil)
}

// ErrorWrapWithFormat tests error wrapping with formatting
func ErrorWrapWithFormat() string {
	base := errors.New("base")
	wrapped := fmt.Errorf("code=%d: %w", 404, base)
	return fmt.Sprintf("%v", errors.Unwrap(wrapped) == base)
}

// ErrorMultipleWrapSame tests wrapping same error multiple times
func ErrorMultipleWrapSame() string {
	base := errors.New("base")
	wrap1 := fmt.Errorf("wrap1: %w", base)
	wrap2 := fmt.Errorf("wrap2: %w", base)
	return fmt.Sprintf("%v:%v", errors.Unwrap(wrap1) == base, errors.Unwrap(wrap2) == base)
}

// ErrorWrapWithContext tests wrapping with context information
func ErrorWrapWithContext() string {
	base := errors.New("file not found")
	wrapped := fmt.Errorf("processing config: %w", base)
	doubleWrapped := fmt.Errorf("initializing service: %w", wrapped)
	return fmt.Sprintf("%v", errors.Unwrap(errors.Unwrap(doubleWrapped)) == base)
}

// ErrorUnwrapNonWrapped tests unwrapping a non-wrapped error
func ErrorUnwrapNonWrapped2() string {
	wrapped := fmt.Errorf("wrapped: %w", errors.New("base"))
	first := errors.Unwrap(wrapped)
	second := errors.Unwrap(first)
	return fmt.Sprintf("%v", second == nil)
}

// ErrorWrapAndError tests wrapped error Error() output
func ErrorWrapAndError() string {
	base := errors.New("base")
	wrapped := fmt.Errorf("outer: %w", base)
	return fmt.Sprintf("%v:%v", wrapped.Error(), errors.Unwrap(wrapped).Error())
}
