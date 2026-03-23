package complex_tests

func VariadicBasic(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func VariadicBasicCheck() int {
	return VariadicBasic(1, 2, 3, 4, 5)
}

func VariadicEmpty() int {
	return VariadicBasic()
}

func VariadicWithRegular(base int, nums ...int) int {
	sum := base
	for _, n := range nums {
		sum += n
	}
	return sum
}

func VariadicWithRegularCheck() int {
	return VariadicWithRegular(100, 1, 2, 3)
}

func VariadicSpread() int {
	nums := []int{1, 2, 3, 4, 5}
	return VariadicBasic(nums...)
}

func VariadicOneArg() int {
	return VariadicBasic(42)
}
