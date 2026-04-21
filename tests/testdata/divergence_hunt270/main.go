package divergence_hunt270

import (
	"fmt"
)

// ============================================================================
// Round 270: Complex control flow — labeled breaks, continues, goto-like patterns
// ============================================================================

// LabeledBreakOuter tests breaking out of outer loop
func LabeledBreakOuter() string {
	result := ""
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer
			}
			result += fmt.Sprintf("%d%d ", i, j)
		}
	}
	return result[:len(result)-1]
}

// LabeledContinueOuter tests continuing outer loop
func LabeledContinueOuter() string {
	result := ""
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue outer
			}
			result += fmt.Sprintf("%d%d ", i, j)
		}
	}
	return result[:len(result)-1]
}

// NestedLoopWithBreak tests breaking inner loop only
func NestedLoopWithBreak() string {
	count := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 2 {
				break
			}
			count++
		}
	}
	return fmt.Sprintf("count=%d", count)
}

// ForRangeBreak tests break in for-range
func ForRangeBreak() string {
	s := []int{10, 20, 30, 40, 50}
	sum := 0
	for _, v := range s {
		if v > 25 {
			break
		}
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// SwitchBreakInLoop tests break in switch inside loop
func SwitchBreakInLoop() string {
	result := ""
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			result += "one "
		default:
			result += fmt.Sprintf("%d ", i)
		}
	}
	return result[:len(result)-1]
}

// LabeledSwitchBreak tests labeled break in switch inside loop
func LabeledSwitchBreak() string {
	result := ""
loop:
	for i := 0; i < 5; i++ {
		switch {
		case i == 3:
			break loop
		default:
			result += fmt.Sprintf("%d ", i)
		}
	}
	return result[:len(result)-1]
}

// NestedSelectBreak tests break in select inside loop
func NestedSelectBreak() string {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	result := ""
	for v := range ch {
		if v == 2 {
			break
		}
		result += fmt.Sprintf("%d", v)
	}
	return result
}

// InfiniteLoopWithBreak tests infinite loop with break condition
func InfiniteLoopWithBreak() string {
	i := 0
	sum := 0
	for {
		if i >= 5 {
			break
		}
		sum += i
		i++
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ContinueInForRange tests continue in for-range
func ContinueInForRange() string {
	s := []int{1, 2, 3, 4, 5}
	sum := 0
	for _, v := range s {
		if v%2 == 0 {
			continue
		}
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// NestedDeferInLoop tests defers accumulated in loop
func NestedDeferInLoop() string {
	var result string
	func() {
		for i := 0; i < 3; i++ {
			v := i
			defer func() {
				result += fmt.Sprintf("%d", v)
			}()
		}
	}()
	return result // "210" - LIFO order
}
