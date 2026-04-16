package divergence_hunt49

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// Round 49: Pointer operations - new, address-of, dereference, double pointer
// ============================================================================

func NewInt() int {
	p := new(int)
	*p = 42
	return *p
}

func NewStruct() int {
	type P struct{ X, Y int }
	p := new(P)
	p.X = 10
	p.Y = 20
	return p.X + p.Y
}

func AddressOf() int {
	x := 42
	p := &x
	return *p
}

func AddressOfModify() int {
	x := 10
	p := &x
	*p = 20
	return x
}

func PointerToSlice() int {
	s := []int{1, 2, 3}
	p := &s
	(*p)[0] = 99
	return s[0]
}

func PointerToMap() int {
	m := map[string]int{"a": 1}
	p := &m
	(*p)["b"] = 2
	return len(m)
}

func PointerToStruct() int {
	type P struct{ X int }
	p := &P{X: 42}
	return p.X
}

func DoublePointer() int {
	x := 10
	p := &x
	pp := &p
	**pp = 20
	return x
}

func NilPointerComparison() bool {
	var p *int
	return p == nil
}

func PointerComparison() bool {
	x := 42
	p1 := &x
	p2 := &x
	return p1 == p2 // same variable, same address
}

func PointerSlice() int {
	a, b, c := 1, 2, 3
	s := []*int{&a, &b, &c}
	return *s[0] + *s[1] + *s[2]
}

func PointerArray() int {
	a := [3]int{10, 20, 30}
	p := &a
	return (*p)[0] + (*p)[2]
}

func StructPointerMethod() int {
	type Counter struct{ n int }
	c := &Counter{n: 0}
	c.n++
	c.n++
	c.n++
	return c.n
}

func JSONPointerRoundTrip() int {
	type Item struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	items := []*Item{
		{Name: "a", Value: 1},
		{Name: "b", Value: 2},
	}
	data, _ := json.Marshal(items)
	var decoded []*Item
	json.Unmarshal(data, &decoded)
	return decoded[0].Value + decoded[1].Value
}

func FmtPointer() string {
	x := 42
	p := &x
	return fmt.Sprintf("%d", *p)
}

func PointerSwap() int {
	a, b := 1, 2
	pa, pb := &a, &b
	*pa, *pb = *pb, *pa
	return a*10 + b
}
