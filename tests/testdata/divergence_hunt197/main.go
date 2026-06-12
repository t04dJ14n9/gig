package divergence_hunt197

import (
	"fmt"
)

// ============================================================================
// Round 197: Pointer operations and semantics
// ============================================================================

// BasicPointer tests basic pointer operations
func BasicPointer() string {
	x := 42
	p := &x
	return fmt.Sprintf("%d", *p)
}

// PointerAssignment tests pointer assignment
func PointerAssignment() string {
	x := 10
	y := 20
	p := &x
	p = &y
	*p = 30
	return fmt.Sprintf("%d:%d", x, y)
}

// PointerToPointer tests pointer to pointer
func PointerToPointer() string {
	x := 5
	p := &x
	pp := &p
	**pp = 10
	return fmt.Sprintf("%d", x)
}

// PointerComparison tests pointer comparison
func PointerComparison() string {
	x := 42
	p1 := &x
	p2 := &x
	p3 := p1
	return fmt.Sprintf("%v:%v", p1 == p2, p1 == p3)
}

// PointerZeroValue tests pointer zero value
func PointerZeroValue() string {
	var p *int
	return fmt.Sprintf("%v", p == nil)
}

// PointerToArray tests pointer to array
func PointerToArray() string {
	a := [3]int{1, 2, 3}
	p := &a
	(*p)[0] = 100
	return fmt.Sprintf("%d", a[0])
}

// PointerToStruct tests pointer to struct
func PointerToStruct() string {
	type Point struct{ X, Y int }
	p := &Point{X: 1, Y: 2}
	p.X = 10
	return fmt.Sprintf("%d:%d", p.X, p.Y)
}

// PointerArithmeticSimulated tests simulated pointer arithmetic via slice
func PointerArithmeticSimulated() string {
	arr := []int{10, 20, 30, 40, 50}
	base := &arr[0]
	// Access via index instead of arithmetic
	val0 := *base
	val1 := arr[1]
	return fmt.Sprintf("%d:%d", val0, val1)
}

// PointerInStruct tests pointer field in struct
func PointerInStruct() string {
	type Node struct {
		Value int
		Next  *Node
	}
	n2 := &Node{Value: 20}
	n1 := &Node{Value: 10, Next: n2}
	return fmt.Sprintf("%d:%d", n1.Value, n1.Next.Value)
}

// PointerSwap tests swapping via pointers
func PointerSwap() string {
	swap := func(a, b *int) {
		*a, *b = *b, *a
	}
	x, y := 1, 2
	swap(&x, &y)
	return fmt.Sprintf("%d:%d", x, y)
}

// PointerToInterface tests pointer to interface
func PointerToInterface() string {
	var i interface{} = 42
	p := &i
	*p = "hello"
	return fmt.Sprintf("%v", *p)
}
