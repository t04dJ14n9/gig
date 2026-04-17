package divergence_hunt104

// ============================================================================
// Round 104: Named return values with bare return
// ============================================================================

func NamedBareReturn() (result int) {
	result = 42
	return
}

func NamedBareReturnModify() (result int) {
	result = 10
	result *= 3
	return
}

func NamedBareReturnConditional() (result string) {
	x := 5
	if x > 3 {
		result = "big"
		return
	}
	result = "small"
	return
}

func NamedMultiBareReturn() (a int, b string) {
	a = 42
	b = "hello"
	return
}

func NamedReturnWithDefer() (result int) {
	defer func() {
		result++
	}()
	result = 10
	return
}

func NamedReturnDeferChain() (result int) {
	defer func() { result += 10 }()
	defer func() { result *= 2 }()
	result = 5
	return
}

func NamedReturnZeroValue() (x int, s string) {
	return
}

func NamedReturnPartial() (a, b int) {
	a = 10
	return
}

func NamedReturnLoop() (sum int) {
	for i := 1; i <= 10; i++ {
		sum += i
	}
	return
}

func NamedReturnClosure() (result int) {
	fn := func() {
		result = 99
	}
	fn()
	return
}
