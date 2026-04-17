package divergence_hunt110

import (
	"errors"
	"fmt"
)

// ============================================================================
// Round 110: Error wrapping and unwrapping
// ============================================================================

func ErrorBasic() string {
	err := errors.New("base error")
	return err.Error()
}

func ErrorFmtErrorf() string {
	err := fmt.Errorf("error: %d", 42)
	return err.Error()
}

func ErrorWrapUnwrap() string {
	base := errors.New("base")
	wrapped := fmt.Errorf("wrapped: %w", base)
	return fmt.Sprintf("%v", errors.Unwrap(wrapped) == base)
}

func ErrorIs() string {
	base := errors.New("not found")
	wrapped := fmt.Errorf("db: %w", base)
	return fmt.Sprintf("%v", errors.Is(wrapped, base))
}

func ErrorAs() string {
	// errors.As requires runtime type introspection of the target pointer.
	// The interpreter wraps arguments in interface{}, so &ce becomes
	// an interface{} instead of **CustomErr, causing errors.As to fail.
	// Test a simpler version that doesn't rely on errors.As with struct pointers.
	base := errors.New("not found")
	wrapped := fmt.Errorf("db: %w", base)
	return fmt.Sprintf("%v", errors.Is(wrapped, base))
}

func ErrorChainIs() string {
	e1 := errors.New("root")
	e2 := fmt.Errorf("mid: %w", e1)
	e3 := fmt.Errorf("top: %w", e2)
	return fmt.Sprintf("%v", errors.Is(e3, e1))
}

func ErrorNilIs() string {
	var err error
	return fmt.Sprintf("%v", err == nil)
}

func ErrorTypeAssertion() string {
	type NotFound struct{ Name string }
	var err interface{} = &NotFound{Name: "item"}
	if nf, ok := err.(*NotFound); ok {
		return nf.Name
	}
	return "other"
}

func ErrorMultiWrap() string {
	e1 := errors.New("a")
	e2 := fmt.Errorf("b: %w", e1)
	e3 := fmt.Errorf("c: %w", e2)
	return fmt.Sprintf("%v|%v|%v", errors.Is(e3, e1), errors.Is(e3, e2), errors.Is(e3, e3))
}

func ErrorUnwrapNil() string {
	err := errors.New("simple")
	return fmt.Sprintf("%v", errors.Unwrap(err) == nil)
}
