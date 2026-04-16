package divergence_hunt59

import (
	"fmt"
	"strings"
)

// ============================================================================
// Round 59: Multi-return values, variadic functions, blank identifiers
// ============================================================================

func MultiReturnSwap() int {
	swap := func(a, b int) (int, int) { return b, a }
	x, y := swap(1, 2)
	return x*10 + y
}

func MultiReturnDivMod() int {
	divmod := func(a, b int) (int, int) { return a / b, a % b }
	q, r := divmod(17, 5)
	return q*10 + r
}

func MultiReturnMinMax() int {
	minmax := func(a, b int) (int, int) {
		if a < b { return a, b }
		return b, a
	}
	min, max := minmax(7, 3)
	return min*10 + max
}

func BlankIdentifier() int {
	_, b := 1, 2
	return b
}

func BlankInMultiReturn() int {
	divmod := func(a, b int) (int, int) { return a / b, a % b }
	q, _ := divmod(17, 5)
	return q
}

func VariadicSum() int {
	sum := func(nums ...int) int {
		total := 0
		for _, n := range nums { total += n }
		return total
	}
	return sum(1, 2, 3, 4, 5)
}

func VariadicSpread() int {
	sum := func(nums ...int) int {
		total := 0
		for _, n := range nums { total += n }
		return total
	}
	nums := []int{1, 2, 3}
	return sum(nums...)
}

func VariadicEmpty() int {
	sum := func(nums ...int) int {
		total := 0
		for _, n := range nums { total += n }
		return total
	}
	return sum()
}

func VariadicWithRegular() string {
	greet := func(prefix string, names ...string) string {
		if len(names) == 0 { return prefix }
		return prefix + " " + strings.Join(names, ",")
	}
	return greet("Hello", "Alice", "Bob")
}

func NamedReturnBare() (result int) {
	result = 42
	return
}

func NamedReturnWithDefer() (result int) {
	defer func() { result++ }()
	return 10
}

func MultipleNamedReturn() (a int, b int) {
	a, b = 1, 2
	return
}

func MultiReturnWithInterface() string {
	process := func(x any) (any, error) {
		return fmt.Sprintf("%v", x), nil
	}
	v, _ := process(42)
	return v.(string)
}

func MultiReturnInLoop() int {
	find := func(s []int, target int) (int, bool) {
		for i, v := range s {
			if v == target { return i, true }
		}
		return -1, false
	}
	idx, found := find([]int{10, 20, 30}, 20)
	if found { return idx }
	return -1
}

func BlankInLoop() int {
	data := []int{10, 20, 30}
	sum := 0
	for _, v := range data { sum += v }
	return sum
}

func MultiReturnError() int {
	mightFail := func(ok bool) (int, error) {
		if !ok { return 0, fmt.Errorf("failed") }
		return 42, nil
	}
	v, err := mightFail(true)
	if err != nil { return -1 }
	return v
}
