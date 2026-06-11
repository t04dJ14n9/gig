package divergence_hunt269

import (
	"fmt"
)

// ============================================================================
// Round 269: Pointer edge cases — double pointers, pointer arithmetic patterns, nil deref
// ============================================================================

// DoublePointer tests pointer to pointer
func DoublePointer() string {
	x := 42
	p := &x
	pp := &p
	**pp = 100
	return fmt.Sprintf("x=%d", x)
}

// PointerSwap tests swapping values through pointers
func PointerSwap() string {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return fmt.Sprintf("a=%d,b=%d", a, b)
}

// NilPointerGuard tests guarding against nil pointer
func NilPointerGuard() string {
	var p *int
	if p != nil {
		return fmt.Sprintf("%d", *p)
	}
	return "nil"
}

// PointerToSlice tests pointer to slice
func PointerToSlice() string {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 99
	return fmt.Sprintf("s=%v", s)
}

// PointerToMap tests pointer to map
func PointerToMap() string {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return fmt.Sprintf("len=%d,a=%d,b=%d", len(m), m["a"], m["b"])
}

// PointerToStruct tests pointer to struct
func PointerToStruct() string {
	type Point struct{ X, Y int }
	p := &Point{X: 1, Y: 2}
	p.X = 10
	p.Y = 20
	return fmt.Sprintf("x=%d,y=%d", p.X, p.Y)
}

// PointerToArray tests pointer to array
func PointerToArray() string {
	a := [3]int{1, 2, 3}
	p := &a
	p[0] = 99
	return fmt.Sprintf("a=%v", a)
}

// PointerReassignment tests reassigning a pointer variable
func PointerReassignment() string {
	x := 10
	y := 20
	p := &x
	v1 := *p
	p = &y
	v2 := *p
	return fmt.Sprintf("v1=%d,v2=%d", v1, v2)
}

// StructPointerMethod modifies through pointer receiver
func StructPointerMethod() string {
	type Counter struct{ n int }
	c := &Counter{n: 0}
	c.n++
	c.n++
	c.n++
	return fmt.Sprintf("n=%d", c.n)
}

// PointerInSlice tests slice of pointers
func PointerInSlice() string {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}
	sum := 0
	for _, p := range s {
		sum += *p
	}
	return fmt.Sprintf("sum=%d", sum)
}
