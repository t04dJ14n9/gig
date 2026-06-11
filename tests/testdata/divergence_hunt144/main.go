package divergence_hunt144

import "fmt"

// ============================================================================
// Round 144: Struct comparison and equality
// ============================================================================

func StructEqual() string {
	type Point struct{ X, Y int }
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 1, Y: 2}
	if p1 == p2 {
		return "equal"
	}
	return "not-equal"
}

func StructNotEqual() string {
	type Point struct{ X, Y int }
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 1, Y: 3}
	if p1 != p2 {
		return "not-equal"
	}
	return "equal"
}

func StructZeroValue() string {
	type Point struct{ X, Y int }
	p := Point{}
	if p.X == 0 && p.Y == 0 {
		return "zero"
	}
	return "non-zero"
}

func StructCopy() string {
	type Data struct{ Val int }
	d1 := Data{Val: 42}
	d2 := d1
	d2.Val = 99
	return fmt.Sprintf("d1=%d-d2=%d", d1.Val, d2.Val)
}

func StructPointerCopy() string {
	type Data struct{ Val int }
	d1 := &Data{Val: 42}
	d2 := d1
	d2.Val = 99
	return fmt.Sprintf("d1=%d-d2=%d", d1.Val, d2.Val)
}

func StructStringField() string {
	type Info struct{ Name string; Age int }
	i := Info{Name: "test", Age: 10}
	return fmt.Sprintf("%s-%d", i.Name, i.Age)
}

func StructBoolField() string {
	type Flags struct{ Active, Visible bool }
	f := Flags{Active: true, Visible: false}
	if f.Active && !f.Visible {
		return "active-hidden"
	}
	return "other"
}

func StructSliceField() string {
	type Container struct{ Items []int }
	c := Container{Items: []int{1, 2, 3}}
	c.Items = append(c.Items, 4)
	return fmt.Sprintf("len=%d", len(c.Items))
}

func StructMapField() string {
	type Cache struct{ Data map[string]int }
	c := Cache{Data: map[string]int{"a": 1}}
	c.Data["b"] = 2
	return fmt.Sprintf("len=%d", len(c.Data))
}

func StructEmbeddedCompare() string {
	type Base struct{ Val int }
	type Extended struct{ Base; Label string }
	e1 := Extended{Base: Base{Val: 1}, Label: "a"}
	e2 := Extended{Base: Base{Val: 1}, Label: "a"}
	if e1 == e2 {
		return "equal"
	}
	return "not-equal"
}
