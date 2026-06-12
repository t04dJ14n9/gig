package divergence_hunt277

import (
	"fmt"
)

// ============================================================================
// Round 277: Variadic function edge cases

// VariadicSum tests basic variadic function
func VariadicSum(nums ...int) string {
	total := 0
	for _, n := range nums {
		total += n
	}
	return fmt.Sprintf("sum=%d,count=%d", total, len(nums))
}

// VariadicEmpty tests calling variadic with no args
func VariadicEmpty() string {
	return VariadicSum()
}

// VariadicOneArg tests calling variadic with one arg
func VariadicOneArg() string {
	return VariadicSum(42)
}

// VariadicMultipleArgs tests calling variadic with multiple args
func VariadicMultipleArgs() string {
	return VariadicSum(1, 2, 3, 4, 5)
}

// VariadicSpread tests spreading a slice with ...
func VariadicSpread() string {
	nums := []int{10, 20, 30}
	return VariadicSum(nums...)
}

// VariadicNilSlice tests spreading nil slice
func VariadicNilSlice() string {
	var nums []int
	return VariadicSum(nums...)
}

// VariadicEmptySlice tests spreading empty slice
func VariadicEmptySlice() string {
	nums := []int{}
	return VariadicSum(nums...)
}

// VariadicInterface tests variadic with interface{}
func VariadicInterface() string {
	print := func(args ...interface{}) string {
		result := ""
		for i, a := range args {
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("%v", a)
		}
		return result
	}
	return print(1, "two", 3.0, true)
}

// VariadicWithRegularParams tests variadic with regular params before it
func VariadicWithRegularParams() string {
	join := func(sep string, parts ...string) string {
		result := ""
		for i, p := range parts {
			if i > 0 {
				result += sep
			}
			result += p
		}
		return result
	}
	return join("-", "a", "b", "c")
}

// VariadicModifySlice tests that modifying variadic arg inside doesn't affect original
func VariadicModifySlice() string {
	nums := []int{1, 2, 3}
	modify := func(args ...int) string {
		if len(args) > 0 {
			args[0] = 99
		}
		return fmt.Sprintf("inner=%v", args)
	}
	inner := modify(nums...)
	return fmt.Sprintf("inner=%s,outer=%v", inner, nums)
}

// VariadicLenCap tests len and cap of variadic parameter
func VariadicLenCap() string {
	check := func(args ...int) string {
		return fmt.Sprintf("len=%d,cap=%d", len(args), cap(args))
	}
	return check(1, 2, 3)
}

// VariadicSpreadFromSubslice tests spreading a subslice
func VariadicSpreadFromSubslice() string {
	nums := []int{1, 2, 3, 4, 5}
	return VariadicSum(nums[1:4]...)
}
