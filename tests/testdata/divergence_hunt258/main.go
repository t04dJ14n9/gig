package divergence_hunt258

import (
	"fmt"
)

// ============================================================================
// Round 258: Blank identifier uses
// ============================================================================

// BlankInShortDecl tests blank in short declaration
func BlankInShortDecl() string {
	x, _ := 1, 2
	_, y := 3, 4
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// BlankInMultiAssign tests blank in multiple assignment
func BlankInMultiAssign() string {
	x, _ := getTwoValues()
	return fmt.Sprintf("x=%d", x)
}

func getTwoValues() (int, int) {
	return 10, 20
}

// BlankInRange tests blank in range loop
func BlankInRange() string {
	nums := []int{10, 20, 30}
	result := ""
	for _, v := range nums {
		result += fmt.Sprintf("%d,", v)
	}
	return result
}

// BlankInRangeIndex tests blank for value in range
func BlankInRangeIndex() string {
	nums := []int{10, 20, 30}
	result := ""
	for i, _ := range nums {
		result += fmt.Sprintf("%d,", i)
	}
	return result
}

// BlankInMapRange tests blank in map range
func BlankInMapRange() string {
	// Map iteration order is non-deterministic; verify only set membership.
	m := map[string]int{"a": 1, "b": 2}
	count := 0
	hasA := false
	hasB := false
	for k, _ := range m {
		count++
		if k == "a" {
			hasA = true
		}
		if k == "b" {
			hasB = true
		}
	}
	return fmt.Sprintf("n=%d,a=%v,b=%v", count, hasA, hasB)
}

// BlankInTypeAssertion tests blank in type assertion
func BlankInTypeAssertion() string {
	var i interface{} = 42
	n, _ := i.(int)
	return fmt.Sprintf("n=%d", n)
}

// BlankInChannelRecv tests blank in channel receive
func BlankInChannelRecv() string {
	ch := make(chan int, 1)
	ch <- 100
	v, _ := <-ch
	return fmt.Sprintf("v=%d", v)
}

// BlankInFuncParams tests blank identifier in function parameters
func BlankInFuncParams() string {
	result := funcWithUnusedParam(10, 20)
	return fmt.Sprintf("result=%d", result)
}

func funcWithUnusedParam(x, _ int) int {
	return x * 2
}

// BlankMultiple tests multiple blanks
func BlankMultiple() string {
	_, _, x := 1, 2, 3
	return fmt.Sprintf("x=%d", x)
}

// BlankInImport tests blank import (no-op for test)
func BlankInImport() string {
	// import _ "fmt" would go at top
	return fmt.Sprintf("blank import works")
}

// BlankWithSideEffects tests blank with side-effect function
func BlankWithSideEffects() string {
	counter := 0
	increment := func() int {
		counter++
		return counter
	}
	_, _ = increment(), increment()
	return fmt.Sprintf("counter=%d", counter)
}
