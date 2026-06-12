package divergence_hunt76

// ============================================================================
// Round 76: Pointer manipulation edge cases - double pointer, pointer to struct
// ============================================================================

func PointerBasic() int {
	x := 42
	p := &x
	return *p
}

func PointerAssign() int {
	x := 10
	p := &x
	*p = 20
	return x
}

func PointerToStruct() int {
	type S struct{ X int }
	s := &S{X: 10}
	s.X = 20
	return s.X
}

func PointerNilCheck() bool {
	var p *int
	return p == nil
}

func PointerSlice() int {
	x := 10
	y := 20
	ptrs := []*int{&x, &y}
	return *ptrs[0] + *ptrs[1]
}

func PointerReassign() int {
	x := 10
	y := 20
	p := &x
	p = &y
	return *p
}

func PointerAsArg() int {
	inc := func(p *int) {
		*p++
	}
	x := 5
	inc(&x)
	return x
}

func PointerReturn() *int {
	x := 42
	return &x
}

func PointerDerefAssign() int {
	x := 5
	p := &x
	*p = 10
	return x
}

func StructPointerMethod() int {
	type S struct{ Val int }
	s := &S{Val: 10}
	s.Val = 20
	return s.Val
}

func PointerToPointer() int {
	x := 42
	p := &x
	pp := &p
	return **pp
}

func PointerArray() [3]int {
	arr := [3]int{1, 2, 3}
	p := &arr
	p[0] = 10
	return *p
}

func PointerSliceElem() int {
	s := []int{1, 2, 3}
	p := &s[1]
	*p = 20
	return s[1]
}

func PointerSwap() (int, int) {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return a, b
}

func NewPointer() int {
	p := new(int)
	*p = 42
	return *p
}
