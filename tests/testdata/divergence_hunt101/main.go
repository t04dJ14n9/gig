package divergence_hunt101

import "fmt"

// ============================================================================
// Round 101: Variadic functions and interface{} args
// ============================================================================

func variadicSum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func variadicConcat(sep string, parts ...string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}

func VariadicSumDirect() int {
	return variadicSum(1, 2, 3, 4, 5)
}

func VariadicConcatDirect() string {
	return variadicConcat("-", "a", "b", "c")
}

func VariadicEmpty() int {
	return variadicSum()
}

func VariadicFromSlice() int {
	nums := []int{10, 20, 30}
	return variadicSum(nums...)
}

func VariadicInterface() string {
	print := func(args ...interface{}) string {
		result := ""
		for _, a := range args {
			result += fmt.Sprintf("%v", a)
		}
		return result
	}
	return print(1, "two", 3.0)
}

func VariadicNil() int {
	return variadicSum(nil...)
}

func VariadicStrings() string {
	return variadicConcat(" ", "hello", "world")
}

func VariadicIntfType() string {
	check := func(args ...interface{}) string {
		if len(args) == 0 {
			return "empty"
		}
		return fmt.Sprintf("%T", args[0])
	}
	return check(42)
}

func VariadicAppend() string {
	base := []int{1, 2}
	result := append(base, 3, 4, 5)
	return fmt.Sprintf("%v", result)
}

func VariadicSpread() string {
	words := []string{"hello", "world"}
	return variadicConcat(" ", words...)
}
