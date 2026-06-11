package divergence_hunt117

import "fmt"

// ============================================================================
// Round 117: Struct comparison and equality
// ============================================================================

func StructEqual() string {
	type Point struct{ X, Y int }
	p1 := Point{1, 2}
	p2 := Point{1, 2}
	return fmt.Sprintf("%v", p1 == p2)
}

func StructNotEqual() string {
	type Point struct{ X, Y int }
	p1 := Point{1, 2}
	p2 := Point{1, 3}
	return fmt.Sprintf("%v", p1 != p2)
}

func StructCopy() string {
	type Data struct{ Val int }
	d1 := Data{Val: 42}
	d2 := d1
	d2.Val = 99
	return fmt.Sprintf("%d:%d", d1.Val, d2.Val)
}

func StructPointerEqual() string {
	type S struct{ Val int }
	s1 := &S{Val: 1}
	s2 := &S{Val: 1}
	return fmt.Sprintf("%v", s1 == s2)
}

func StructPointerSame() string {
	type S struct{ Val int }
	s1 := &S{Val: 1}
	s2 := s1
	return fmt.Sprintf("%v", s1 == s2)
}

func StructWithSlice() string {
	type Bag struct{ Items []int }
	b := Bag{Items: []int{1, 2, 3}}
	return fmt.Sprintf("%d", len(b.Items))
}

func StructWithMap() string {
	type Config struct{ Entries map[string]int }
	c := Config{Entries: map[string]int{"x": 1}}
	return fmt.Sprintf("%d", c.Entries["x"])
}

func StructNested() string {
	type Inner struct{ V int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{V: 42}}
	return fmt.Sprintf("%d", o.I.V)
}

func StructZeroValue() string {
	type S struct{ X int; Y string }
	var s S
	return fmt.Sprintf("%d:%s:%v", s.X, s.Y, s.Y == "")
}

func StructSliceOfPointers() string {
	type Item struct{ Val int }
	items := []*Item{{1}, {2}, {3}}
	total := 0
	for _, item := range items {
		total += item.Val
	}
	return fmt.Sprintf("%d", total)
}
