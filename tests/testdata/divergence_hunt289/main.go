package divergence_hunt289

import (
	"fmt"
)

// ============================================================================
// Round 289: Complex control flow — for patterns, break/continue, labeled loops

// ForAsWhile tests for used as while loop
func ForAsWhile() string {
	x := 1
	for x < 100 {
		x *= 2
	}
	return fmt.Sprintf("x=%d", x)
}

// ForInfiniteWithBreak tests infinite for with break
func ForInfiniteWithBreak() string {
	i := 0
	for {
		i++
		if i >= 5 {
			break
		}
	}
	return fmt.Sprintf("i=%d", i)
}

// NestedForBreakOuter tests breaking out of outer loop with label
func NestedForBreakOuter() string {
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

// NestedForContinueOuter tests continuing outer loop with label
func NestedForContinueOuter() string {
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

// ForRangeWithIndexOnly tests for range with only index
func ForRangeWithIndexOnly() string {
	s := []string{"a", "b", "c"}
	count := 0
	for i := range s {
		count += i
	}
	return fmt.Sprintf("count=%d", count)
}

// ForRangeWithValueOnly tests for range with only value
func ForRangeWithValueOnly() string {
	s := []int{10, 20, 30}
	sum := 0
	for _, v := range s {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ForRangeMap tests ranging over map
func ForRangeMap() string {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// ForRangeString tests ranging over string (produces runes)
func ForRangeString() string {
	s := "Go语"
	result := ""
	for i, r := range s {
		result += fmt.Sprintf("%d:%d ", i, r)
	}
	return result[:len(result)-1]
}

// ContinueInForRange tests continue in for-range
func ContinueInForRange() string {
	s := []int{1, 2, 3, 4, 5, 6}
	sum := 0
	for _, v := range s {
		if v%2 == 0 {
			continue
		}
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// NestedForRange tests nested for-range
func NestedForRange() string {
	matrix := [][]int{{1, 2}, {3, 4}}
	total := 0
	for _, row := range matrix {
		for _, v := range row {
			total += v
		}
	}
	return fmt.Sprintf("total=%d", total)
}
