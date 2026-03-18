package multiassign

// Swap tests multiple assignment swap
func Swap() int {
	a := 10
	b := 20
	a, b = b, a
	return a*100 + b
}

// FromFunction tests multiple assignment from function
func FromFunction() int {
	a, b := twoVals()
	return a + b
}

func twoVals() (int, int) { return 42, 58 }

// ThreeValues tests three value assignment
func ThreeValues() int {
	a, b, c := threeVals(10)
	return a + b + c
}

func threeVals(x int) (int, int, int) {
	return x, x * 2, x * 3
}

// InLoop tests multiple assignment in loop
func InLoop() int {
	a, b := 0, 1
	for i := 0; i < 10; i++ {
		a, b = b, a+b
	}
	return a
}

// DiscardWithBlank tests discarding with blank
func DiscardWithBlank() int {
	q, _ := divmod(17, 5)
	return q
}

func divmod(a, b int) (int, int) { return a / b, a % b }

// ============================================================================
// Exported wrappers for parameterized testing
// ============================================================================

// TwoVals returns (42, 58)
func TwoVals() (int, int) { return twoVals() }

// ThreeValsN returns (x, x*2, x*3)
func ThreeValsN(x int) (int, int, int) { return threeVals(x) }

// DivmodAB returns (a/b, a%b)
func DivmodAB(a, b int) (int, int) { return divmod(a, b) }
