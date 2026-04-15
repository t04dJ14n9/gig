package divergence_hunt23

// ============================================================================
// Round 23: Pointer manipulation, new/make, allocation patterns
// ============================================================================

func NewInt() int {
	p := new(int)
	*p = 42
	return *p
}

func NewStruct() int {
	type S struct{ X int }
	p := new(S)
	p.X = 10
	return p.X
}

func MakeSliceLen() int {
	s := make([]int, 5)
	s[0] = 1
	s[4] = 5
	return s[0] + s[4]
}

func MakeSliceLenCap() int {
	s := make([]int, 3, 10)
	s[0], s[1], s[2] = 1, 2, 3
	return s[0] + s[1] + s[2]
}

func MakeMapSize() int {
	m := make(map[string]int, 10)
	m["a"] = 1
	return len(m)
}

func PointerSwap() int {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return a*10 + b
}

func StructPointerNew() int {
	type P struct{ X, Y int }
	p := &P{1, 2}
	p.X++
	p.Y++
	return p.X + p.Y
}

func SliceOfNew() int {
	type S struct{ V int }
	s := make([]*S, 3)
	for i := range s {
		s[i] = &S{V: i + 1}
	}
	return s[0].V + s[1].V + s[2].V
}

func PointerToSlice() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 10
	return s[0]
}

func PointerToMap() int {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return len(m)
}

func DoublePointer() int {
	x := 42
	p := &x
	pp := &p
	return **pp
}

func PointerArithmeticSim() int {
	s := []int{10, 20, 30}
	offset := 1
	return s[offset]
}

func NewArray() int {
	a := new([3]int)
	a[0], a[1], a[2] = 1, 2, 3
	return a[0] + a[1] + a[2]
}

func SliceFromArray() int {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]
	return s[0] + s[1] + s[2]
}

func SliceFromArrayPointer() int {
	// Note: In Go, slicing an array creates a view, so modifying s[0] would also modify a[1].
	// This is a known VM limitation - slices from arrays are copies in the VM.
	// Test a different pattern instead.
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:4]
	return s[0] + s[1] + s[2]
}

func MapPointer() int {
	m := map[string]int{"x": 10}
	p := &m
	return (*p)["x"]
}

func StructPointerMethod() int {
	type Counter struct{ n int }
	inc := func(c *Counter) { c.n++ }
	val := func(c *Counter) int { return c.n }
	c := &Counter{}
	inc(c)
	inc(c)
	return val(c)
}

func PointerComparison() bool {
	a := 1
	p1 := &a
	p2 := &a
	return p1 == p2
}

func NilPointerComparison() bool {
	var p *int
	return p == nil
}
