package divergence_hunt282

import (
	"fmt"
)

// ============================================================================
// Round 282: Pointer edge cases — new, &, double pointer, pointer arithmetic simulation

// NewInt tests new(int) creates zero-valued int pointer
func NewInt() string {
	p := new(int)
	return fmt.Sprintf("val=%d,nil=%t", *p, p == nil)
}

// NewStruct tests new(struct) creates zero-valued struct pointer
func NewStruct() string {
	type S struct{ X int }
	p := new(S)
	return fmt.Sprintf("x=%d,nil=%t", p.X, p == nil)
}

// DoublePointer tests **int (pointer to pointer)
func DoublePointer() string {
	x := 42
	p := &x
	pp := &p
	**pp = 99
	return fmt.Sprintf("x=%d,*p=%d,**pp=%d", x, *p, **pp)
}

// PointerSwap tests swapping pointers
func PointerSwap() string {
	a, b := 1, 2
	pa, pb := &a, &b
	pa, pb = pb, pa
	return fmt.Sprintf("*pa=%d,*pb=%d", *pa, *pb)
}

// NilPointerDereferencePanics tests nil pointer dereference
func NilPointerDereferencePanics() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	var p *int
	_ = *p
	return "no_panic"
}

// PointerToArray tests pointer to array
func PointerToArray() string {
	arr := [3]int{1, 2, 3}
	p := &arr
	p[1] = 99
	return fmt.Sprintf("arr=%v", arr)
}

// PointerToSliceHeader tests taking address of slice variable
func PointerToSliceHeader() string {
	s := []int{1, 2, 3}
	p := &s
	*p = append(*p, 4)
	return fmt.Sprintf("s=%v,len=%d", s, len(s))
}

// PointerComparison tests pointer comparison
func PointerComparison() string {
	x := 42
	p1 := &x
	p2 := &x
	return fmt.Sprintf("same=%t", p1 == p2)
}

// NewMap tests new(map) creates nil map
func NewMap() string {
	p := new(map[string]int)
	return fmt.Sprintf("nil=%t", *p == nil)
}

// AddressOfLiteral tests taking address of composite literal
func AddressOfLiteral() string {
	type Point struct{ X, Y int }
	p := &Point{1, 2}
	p.X = 10
	return fmt.Sprintf("x=%d,y=%d", p.X, p.Y)
}
