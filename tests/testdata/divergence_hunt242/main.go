package divergence_hunt242

import (
	"fmt"
)

// ============================================================================
// Round 242: Defer with multiple returns
// ============================================================================

// MultipleNamedReturnDefer tests defer modifying multiple named returns
func MultipleNamedReturnDefer() string {
	a, b := func() (x, y int) {
		defer func() { x = 100; y = 200 }()
		return 1, 2
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferOne tests defer modifying only one named return
func MultipleNamedReturnDeferOne() string {
	a, b := func() (x, y int) {
		defer func() { x = 99 }()
		return 1, 2
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferSwap tests defer swapping named returns
func MultipleNamedReturnDeferSwap() string {
	a, b := func() (x, y int) {
		defer func() { x, y = y, x }()
		return 10, 20
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferIncrement tests defer incrementing both returns
func MultipleNamedReturnDeferIncrement() string {
	a, b := func() (x, y int) {
		defer func() { x++; y++ }()
		return 5, 7
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferMultiply tests defer multiplying both returns
func MultipleNamedReturnDeferMultiply() string {
	a, b := func() (x, y int) {
		defer func() { x *= 2; y *= 3 }()
		return 3, 4
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferChain tests chained defers with multiple returns
func MultipleNamedReturnDeferChain() string {
	a, b := func() (x, y int) {
		defer func() { x = x + 1; y = y + 1 }()
		defer func() { x = x * 10; y = y * 10 }()
		return 1, 2
	}()
	return fmt.Sprintf("%d:%d", a, b)
}

// MultipleNamedReturnDeferStringInt tests defer with string and int returns
func MultipleNamedReturnDeferStringInt() string {
	s, n := func() (str string, num int) {
		defer func() { str = str + "!"; num = num * 10 }()
		return "hello", 5
	}()
	return fmt.Sprintf("%s:%d", s, n)
}

// MultipleNamedReturnDeferBoolInt tests defer with bool and int returns
func MultipleNamedReturnDeferBoolInt() string {
	b, n := func() (ok bool, count int) {
		defer func() { ok = true; count = 42 }()
		return false, 0
	}()
	return fmt.Sprintf("%v:%d", b, n)
}

// MultipleNamedReturnDeferSliceMap tests defer with slice and map returns
func MultipleNamedReturnDeferSliceMap() string {
	s, m := func() (sl []int, mp map[string]int) {
		defer func() {
			sl = append(sl, 4)
			mp["b"] = 2
		}()
		return []int{1, 2, 3}, map[string]int{"a": 1}
	}()
	return fmt.Sprintf("%v:%v", s, m)
}

// MultipleNamedReturnDeferPanicRecover tests defer with panic/recover and multiple returns
func MultipleNamedReturnDeferPanicRecover() string {
	a, b := func() (x, y int) {
		defer func() {
			if r := recover(); r != nil {
				x = -1
				y = -1
			}
		}()
		panic("error")
	}()
	return fmt.Sprintf("%d:%d", a, b)
}
