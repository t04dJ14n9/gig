package divergence_hunt287

import (
	"fmt"
)

// ============================================================================
// Round 287: Multiple return values — named/unnamed, blank identifier, complex patterns

// MultipleReturnBasic tests basic multiple return
func MultipleReturnBasic() string {
	divide := func(a, b int) (int, int) {
		return a / b, a % b
	}
	q, r := divide(17, 5)
	return fmt.Sprintf("q=%d,r=%d", q, r)
}

// MultipleReturnNamed tests named return values
func MultipleReturnNamed() string {
	divide := func(a, b int) (quotient, remainder int) {
		quotient = a / b
		remainder = a % b
		return
	}
	q, r := divide(17, 5)
	return fmt.Sprintf("q=%d,r=%d", q, r)
}

// MultipleReturnBlank tests blank identifier for unused return
func MultipleReturnBlank() string {
	divide := func(a, b int) (int, int) {
		return a / b, a % b
	}
	q, _ := divide(17, 5)
	return fmt.Sprintf("q=%d", q)
}

// MultipleReturnSwap tests using multiple returns for swap
func MultipleReturnSwap() string {
	swap := func(a, b int) (int, int) {
		return b, a
	}
	x, y := swap(1, 2)
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// MultipleReturnWithError tests error return pattern
func MultipleReturnWithError() string {
	safeDivide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, fmt.Errorf("divide by zero")
		}
		return a / b, nil
	}
	r1, e1 := safeDivide(10, 2)
	r2, e2 := safeDivide(10, 0)
	return fmt.Sprintf("r1=%d,e1=%t,r2=%d,e2=%t", r1, e1 == nil, r2, e2 == nil)
}

// MultipleReturnAllNamedBareReturn tests all named returns with bare return
func MultipleReturnAllNamedBareReturn() string {
	compute := func() (sum, product int) {
		sum = 3 + 7
		product = 3 * 7
		return
	}
	s, p := compute()
	return fmt.Sprintf("sum=%d,product=%d", s, p)
}

// NestedFunctionReturns tests function returning function
func NestedFunctionReturns() string {
	adder := func(x int) func(int) int {
		return func(y int) int {
			return x + y
		}
	}
	f := adder(10)
	return fmt.Sprintf("result=%d", f(5))
}

// MultipleReturnInDefer tests multiple return values modified by defer
func MultipleReturnInDefer() (a, b int) {
	a = 1
	b = 2
	defer func() {
		a = 10
		b = 20
	}()
	return
}

// MultipleReturnToInterface tests multiple returns assigned to variables of different types
func MultipleReturnToInterface() string {
	getBoth := func() (int, string) {
		return 42, "hello"
	}
	var i interface{} 
	var s interface{}
	i, s = getBoth()
	return fmt.Sprintf("i=%v,s=%v", i, s)
}

// MultipleReturnOnlyOneUsed tests using only one of multiple return values
func MultipleReturnOnlyOneUsed() string {
	get := func() (int, string) {
		return 42, "unused"
	}
	_, s := get()
	return s
}
