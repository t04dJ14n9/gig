package divergence_hunt138

import "fmt"

// ============================================================================
// Round 138: Pointer patterns — new, &, dereference, nil pointer
// ============================================================================

func NewInt() string {
	p := new(int)
	*p = 42
	return fmt.Sprintf("val=%d", *p)
}

func NewStruct() string {
	type Point struct{ X, Y int }
	p := new(Point)
	p.X = 10
	p.Y = 20
	return fmt.Sprintf("%d-%d", p.X, p.Y)
}

func AddressOf() string {
	x := 100
	p := &x
	*p = 200
	return fmt.Sprintf("x=%d", x)
}

func NilPointerCheck() string {
	var p *int
	if p == nil {
		return "nil"
	}
	return "not-nil"
}

func PointerSwap() string {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return fmt.Sprintf("a=%d-b=%d", a, b)
}

func PointerToSlice() string {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 99
	return fmt.Sprintf("%v", s)
}

func PointerToMap() string {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return fmt.Sprintf("len=%d", len(m))
}

func PointerStructMethod() string {
	type Wrapper struct{ Val int }
	w := &Wrapper{Val: 42}
	w.Val = 100
	return fmt.Sprintf("val=%d", w.Val)
}

func DoublePointer() string {
	x := 5
	p := &x
	pp := &p
	**pp = 99
	return fmt.Sprintf("x=%d", x)
}

func PointerArray() string {
	arr := [3]int{1, 2, 3}
	p := &arr
	p[1] = 99
	return fmt.Sprintf("%v", arr)
}
