package divergence_hunt148

import (
	"errors"
	"fmt"
)

// ============================================================================
// Round 148: Error chain traversal with standard library
// ============================================================================

func ErrorsNewCheck() string {
	err := errors.New("test error")
	return err.Error()
}

func ErrorsIsMatch() string {
	err1 := errors.New("base")
	err2 := fmt.Errorf("wrapped: %w", err1)
	if errors.Is(err2, err1) {
		return "match"
	}
	return "no-match"
}

func ErrorsAsInterface() string {
	err := fmt.Errorf("error with code %d", 404)
	// Test errors.As with a concrete struct target (not interface target)
	var target *CustomErr
	if errors.As(err, &target) {
		return "matched"
	}
	return "no-match"
}

type CustomErr struct {
	Code int
}

func (e *CustomErr) Error() string {
	return fmt.Sprintf("code=%d", e.Code)
}

func ErrorsUnwrapNil() string {
	err := errors.New("simple")
	if unwrapped := errors.Unwrap(err); unwrapped == nil {
		return "nil"
	}
	return "has-unwrap"
}

func ErrorfWrapUnwrap() string {
	base := errors.New("base")
	wrapped := fmt.Errorf("context: %w", base)
	if errors.Unwrap(wrapped) == base {
		return "unwrapped"
	}
	return "failed"
}

func ErrorfMultiWrap() string {
	e1 := errors.New("e1")
	e2 := fmt.Errorf("e2: %w", e1)
	e3 := fmt.Errorf("e3: %w", e2)
	if errors.Is(e3, e1) {
		return "deep-match"
	}
	return "no-deep-match"
}

func ErrorJoin() string {
	e1 := errors.New("err1")
	e2 := errors.New("err2")
	joined := errors.Join(e1, e2)
	if joined != nil {
		return "has-errors"
	}
	return "nil"
}

func ErrorJoinIs() string {
	e1 := errors.New("err1")
	e2 := errors.New("err2")
	joined := errors.Join(e1, e2)
	if errors.Is(joined, e1) && errors.Is(joined, e2) {
		return "both-found"
	}
	return "missing"
}

func ErrorNilIs() string {
	var err error
	if errors.Is(err, nil) {
		return "nil-is-nil"
	}
	return "not-nil"
}
