package divergence_hunt107

import "fmt"

// ============================================================================
// Round 107: Pointer patterns - double pointers, pointer to struct, new()
// ============================================================================

func PointerBasic() int {
	x := 42
	p := &x
	return *p
}

func PointerModify() int {
	x := 10
	p := &x
	*p = 20
	return x
}

func PointerToStruct() string {
	type S struct{ Val int }
	s := &S{Val: 42}
	s.Val = 99
	return fmt.Sprintf("%d", s.Val)
}

func PointerSwap() string {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return fmt.Sprintf("%d:%d", a, b)
}

func NewKeyword() int {
	p := new(int)
	*p = 42
	return *p
}

func NewStruct() string {
	type Item struct{ Name string; Val int }
	p := new(Item)
	p.Name = "test"
	p.Val = 10
	return fmt.Sprintf("%s:%d", p.Name, p.Val)
}

func PointerSlice() string {
	x, y, z := 1, 2, 3
	ptrs := []*int{&x, &y, &z}
	total := 0
	for _, p := range ptrs {
		total += *p
	}
	return fmt.Sprintf("%d", total)
}

func NilPointerCheck() string {
	var p *int
	if p == nil {
		return "nil"
	}
	return "not nil"
}

func PointerAsParam() int {
	double := func(p *int) {
		*p *= 2
	}
	x := 5
	double(&x)
	return x
}

func PointerReturn() int {
	makePtr := func(v int) *int {
		return &v
	}
	p := makePtr(42)
	return *p
}
