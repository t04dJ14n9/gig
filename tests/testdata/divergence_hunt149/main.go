package divergence_hunt149

import (
	"fmt"
	"sort"
	"strconv"
)

// ============================================================================
// Round 149: Variadic functions and interface{} patterns
// ============================================================================

func VariadicSum(nums ...int) string {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return fmt.Sprintf("sum=%d", sum)
}

func VariadicStringJoin() string {
	parts := []string{"hello", "world", "test"}
	sep := "-"
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

func VariadicEmpty() string {
	return VariadicSumHelper()
}

func VariadicSumHelper(nums ...int) string {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return fmt.Sprintf("sum=%d", sum)
}

func VariadicInterface(args ...interface{}) string {
	result := ""
	for _, a := range args {
		result += fmt.Sprintf("%v", a)
	}
	return result
}

func VariadicSpread() string {
	nums := []int{1, 2, 3, 4, 5}
	return VariadicSumHelper(nums...)
}

func VariadicFmt() string {
	return fmt.Sprintf("%s-%d-%t", "hello", 42, true)
}

func VariadicSort() string {
	nums := []int{5, 3, 1, 4, 2}
	sort.Ints(nums)
	return fmt.Sprintf("%v", nums)
}

func VariadicStrconv() string {
	i, err := strconv.Atoi("42")
	if err != nil {
		return "error"
	}
	return fmt.Sprintf("val=%d", i)
}

func VariadicPrintf() string {
	return fmt.Sprintf("name=%s age=%d", "Alice", 30)
}
