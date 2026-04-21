package divergence_hunt243

import (
	"fmt"
)

// ============================================================================
// Round 243: Defer in loops
// ============================================================================

// DeferInForLoop tests defer inside a for loop
func DeferInForLoop() string {
	result := ""
	for i := 0; i < 3; i++ {
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(i)
	}
	return result
}

// DeferInReverseLoop tests defer in reverse for loop
func DeferInReverseLoop() string {
	result := ""
	for i := 3; i > 0; i-- {
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(i)
	}
	return result
}

// DeferAccumulateInLoop tests defer accumulating values in loop
func DeferAccumulateInLoop() string {
	sum := 0
	for i := 1; i <= 5; i++ {
		defer func(n int) {
			sum += n
		}(i)
	}
	return fmt.Sprintf("%d", sum)
}

// DeferInNestedLoop tests defer in nested loops
func DeferInNestedLoop() string {
	result := ""
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			defer func(a, b int) {
				result = fmt.Sprintf("%s%d%d", result, a, b)
			}(i, j)
		}
	}
	return result
}

// DeferWithLoopVariableCapture tests defer capturing loop variable
func DeferWithLoopVariableCapture() string {
	result := ""
	for i := 0; i < 3; i++ {
		v := i
		defer func() {
			result = fmt.Sprintf("%s%d", result, v)
		}()
	}
	return result
}

// DeferInRangeLoop tests defer in range loop over slice
func DeferInRangeLoop() string {
	result := ""
	for _, v := range []int{10, 20, 30} {
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(v)
	}
	return result
}

// DeferInRangeLoopIndex tests defer in range loop capturing index
func DeferInRangeLoopIndex() string {
	result := ""
	for i := range []int{100, 200, 300} {
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(i)
	}
	return result
}

// DeferInRangeMap tests defer in range loop over map
func DeferInRangeMap() string {
	count := 0
	for k, v := range map[string]int{"a": 1, "b": 2, "c": 3} {
		defer func(key string, val int) {
			count += val
			_ = key
		}(k, v)
	}
	return fmt.Sprintf("%d", count)
}

// DeferWithBreakLoop tests defer with break in loop
func DeferWithBreakLoop() string {
	result := ""
	for i := 0; i < 10; i++ {
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(i)
		if i == 2 {
			break
		}
	}
	return result
}

// DeferWithContinueLoop tests defer with continue in loop
func DeferWithContinueLoop() string {
	result := ""
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			continue
		}
		defer func(n int) {
			result = fmt.Sprintf("%s%d", result, n)
		}(i)
	}
	return result
}

// DeferConditionalInLoop tests conditional defer in loop
func DeferConditionalInLoop() string {
	result := ""
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			defer func(n int) {
				result = fmt.Sprintf("%s%d", result, n)
			}(i)
		}
	}
	return result
}
