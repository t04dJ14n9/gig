package divergence_hunt223

import "fmt"

// ============================================================================
// Round 223: Map with struct keys
// ============================================================================

type Point struct {
	X, Y int
}

type Person struct {
	Name string
	Age  int
}

type Coord3D struct {
	X, Y, Z float64
}

// MapStructKeyBasic tests basic struct key operations
func MapStructKeyBasic() string {
	m := map[Point]string{
		{1, 2}: "A",
		{3, 4}: "B",
		{5, 6}: "C",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapStructKeyLookup tests struct key lookup
func MapStructKeyLookup() string {
	m := map[Point]string{
		{1, 2}: "origin",
		{10, 20}: "far",
	}
	p := Point{X: 1, Y: 2}
	v, ok := m[p]
	return fmt.Sprintf("v=%s,ok=%t", v, ok)
}

// MapStructKeyInsert tests inserting with struct keys
func MapStructKeyInsert() string {
	m := map[Point]int{}
	m[Point{1, 1}] = 1
	m[Point{2, 2}] = 4
	m[Point{3, 3}] = 9
	return fmt.Sprintf("len=%d", len(m))
}

// MapStructKeyDelete tests deleting struct keys
func MapStructKeyDelete() string {
	m := map[Point]string{
		{1, 1}: "a",
		{2, 2}: "b",
		{3, 3}: "c",
	}
	delete(m, Point{2, 2})
	return fmt.Sprintf("len=%d", len(m))
}

// MapStructKeyIterate tests iterating over struct-keyed map
func MapStructKeyIterate() string {
	m := map[Point]int{
		{1, 0}: 1,
		{0, 1}: 2,
		{1, 1}: 3,
	}
	sum := 0
	for p, v := range m {
		sum += p.X + p.Y + v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapStructKeyWithStringField tests struct with string field as key
func MapStructKeyWithStringField() string {
	m := map[Person]int{
		{"Alice", 30}: 1,
		{"Bob", 25}: 2,
		{"Alice", 25}: 3,
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapStructKeyZeroValue tests zero value struct key
func MapStructKeyZeroValue() string {
	m := map[Point]int{}
	m[Point{}] = 100
	v, ok := m[Point{0, 0}]
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// MapStructKeyFloat tests struct with float fields as key
func MapStructKeyFloat() string {
	m := map[Coord3D]string{
		{1.5, 2.5, 3.5}: "A",
		{1.5, 2.5, 3.5}: "B",
		{1.0, 2.0, 3.0}: "C",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapStructKeyComposite tests composite struct operations
func MapStructKeyComposite() string {
	type Key struct {
		A int
		B string
	}
	m := map[Key]int{
		{1, "a"}: 10,
		{1, "b"}: 20,
		{2, "a"}: 30,
		{2, "b"}: 40,
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapStructKeyUpdate updates values for struct keys
func MapStructKeyUpdate() string {
	m := map[Point]int{{0, 0}: 1}
	p := Point{X: 0, Y: 0}
	m[p] = 100
	return fmt.Sprintf("v=%d", m[Point{}])
}

// MapStructKeyCommaOk tests comma-ok with struct keys
func MapStructKeyCommaOk() string {
	m := map[Point]string{
		{1, 2}: "found",
	}
	_, ok1 := m[Point{1, 2}]
	_, ok2 := m[Point{3, 4}]
	return fmt.Sprintf("ok1=%t,ok2=%t", ok1, ok2)
}
