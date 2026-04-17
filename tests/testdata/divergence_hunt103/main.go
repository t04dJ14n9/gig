package divergence_hunt103

import (
	"errors"
	"fmt"
)

// ============================================================================
// Round 103: Multiple return values and blank identifier
// ============================================================================

func DivMod(a, b int) (int, int) {
	return a / b, a % b
}

func Swap(a, b int) (int, int) {
	return b, a
}

func MinMax(a, b int) (min, max int) {
	if a < b {
		return a, b
	}
	return b, a
}

func MultiReturnBlank() int {
	_, mod := DivMod(17, 5)
	return mod
}

func MultiReturnAll() string {
	div, mod := DivMod(17, 5)
	return fmt.Sprintf("%d:%d", div, mod)
}

func ErrorReturn() string {
	safeDiv := func(a, b int) (int, error) {
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	if r, err := safeDiv(10, 3); err == nil {
		return fmt.Sprintf("%d", r)
	}
	return "error"
}

func ErrorReturnFail() string {
	safeDiv := func(a, b int) (int, error) {
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	}
	if _, err := safeDiv(10, 0); err != nil {
		return "error"
	}
	return "ok"
}

func NamedReturnBare() int {
	result := 42
	return result
}

func NamedReturnOverride() string {
	val := func() (x int) {
		x = 10
		return
	}()
	return fmt.Sprintf("%d", val)
}

func SwapValues() string {
	a, b := 1, 2
	a, b = b, a
	return fmt.Sprintf("%d:%d", a, b)
}

func MultiAssignExpression() string {
	a, b := 1, 2
	a, b = a+b, a*b
	return fmt.Sprintf("%d:%d", a, b)
}

func BlankInLoop() int {
	data := []string{"a", "b", "c"}
	count := 0
	for range data {
		count++
	}
	return count
}
