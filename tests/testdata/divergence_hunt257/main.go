package divergence_hunt257

import (
	"fmt"
)

// ============================================================================
// Round 257: Multiple assignments
// ============================================================================

// MultiAssignBasic tests basic multiple assignment
func MultiAssignBasic() string {
	x, y := 1, 2
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// MultiAssignSwap tests variable swapping
func MultiAssignSwap() string {
	a, b := 10, 20
	a, b = b, a
	return fmt.Sprintf("a=%d,b=%d", a, b)
}

// MultiAssignTripleSwap tests triple swap
func MultiAssignTripleSwap() string {
	a, b, c := 1, 2, 3
	a, b, c = c, a, b
	return fmt.Sprintf("a=%d,b=%d,c=%d", a, b, c)
}

// MultiAssignFunctionReturn tests assignment from function
func MultiAssignFunctionReturn() string {
	x, y := getTwoValues()
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

func getTwoValues() (int, int) {
	return 100, 200
}

// MultiAssignMapLookup tests assignment from map lookup
func MultiAssignMapLookup() string {
	m := map[string]int{"key": 42}
	v, ok := m["key"]
	return fmt.Sprintf("v=%d,ok=%v", v, ok)
}

// MultiAssignTypeAssertion tests assignment from type assertion
func MultiAssignTypeAssertion() string {
	var i interface{} = 42
	n, ok := i.(int)
	return fmt.Sprintf("n=%d,ok=%v", n, ok)
}

// MultiAssignChannelRecv tests assignment from channel receive
func MultiAssignChannelRecv() string {
	ch := make(chan int, 1)
	ch <- 42
	v, ok := <-ch
	return fmt.Sprintf("v=%d,ok=%v", v, ok)
}

// MultiAssignSameValue tests assigning same value to multiple vars
func MultiAssignSameValue() string {
	x := 5
	y := 5
	z := 5
	return fmt.Sprintf("x=%d,y=%d,z=%d", x, y, z)
}

// MultiAssignExpression tests assignment with expressions
func MultiAssignExpression() string {
	x, y := 1+2, 3*4
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// MultiAssignArrayIndex tests assignment with array indexing
func MultiAssignArrayIndex() string {
	arr := []int{10, 20, 30}
	x, y := arr[0], arr[2]
	return fmt.Sprintf("x=%d,y=%d", x, y)
}

// MultiAssignMixedTypes tests assignment with different types
func MultiAssignMixedTypes() string {
	i, s, b := 42, "hello", true
	return fmt.Sprintf("i=%d,s=%s,b=%v", i, s, b)
}

// MultiAssignReassign tests reassignment with multiple values
func MultiAssignReassign() string {
	x, y := 1, 2
	x, y = x+10, y+20
	return fmt.Sprintf("x=%d,y=%d", x, y)
}
