package divergence_hunt74

// ============================================================================
// Round 74: Variadic and spread edge cases
// ============================================================================

func VariadicSum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func VariadicDirect() int {
	return VariadicSum(1, 2, 3, 4, 5)
}

func VariadicEmpty() int {
	return VariadicSum()
}

func VariadicSpread() int {
	nums := []int{10, 20, 30}
	return VariadicSum(nums...)
}

func VariadicWithPrefix() int {
	return VariadicWithPrefixImpl(99, 1, 2, 3)
}

func VariadicWithPrefixImpl(prefix int, nums ...int) int {
	result := prefix
	for _, n := range nums {
		result += n
	}
	return result
}

func VariadicString(parts ...string) string {
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

func VariadicInterface(items ...any) int {
	return len(items)
}

func VariadicNilSpread() int {
	var nums []int
	return VariadicSum(nums...)
}

func VariadicInClosure() int {
	f := func(nums ...int) int {
		return len(nums)
	}
	return f(1, 2, 3)
}

func VariadicAppend() []int {
	base := []int{1, 2, 3}
	result := append(base, 4, 5)
	return result
}

func VariadicAppendSpread() []int {
	base := []int{1, 2, 3}
	extra := []int{4, 5}
	result := append(base, extra...)
	return result
}

func VariadicFmt() string {
	return sprintf3("%d %d %d", 1, 2, 3)
}

func sprintf3(format string, args ...int) string {
	_ = format
	result := 0
	for _, a := range args {
		result += a
	}
	return "sum"
}

func VariadicReturnSlice() []int {
	return makeVariadic(1, 2, 3)
}

func makeVariadic(nums ...int) []int {
	return nums
}
