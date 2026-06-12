package divergence_hunt225

import "fmt"

// ============================================================================
// Round 225: Map with pointer keys
// ============================================================================

// MapPointerKeyBasic tests basic pointer key operations
func MapPointerKeyBasic() string {
	a, b, c := 1, 2, 3
	m := map[*int]string{
		&a: "a",
		&b: "b",
		&c: "c",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyLookup tests pointer key lookup
func MapPointerKeyLookup() string {
	x := 42
	m := map[*int]string{&x: "found"}
	v, ok := m[&x]
	return fmt.Sprintf("v=%s,ok=%t", v, ok)
}

// MapPointerKeyDifferentAddresses different pointers with same value
func MapPointerKeyDifferentAddresses() string {
	x, y := 10, 10
	m := map[*int]int{}
	m[&x] = 1
	m[&y] = 2
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyFromSlice uses pointers from slice elements
func MapPointerKeyFromSlice() string {
	nums := []int{10, 20, 30}
	m := map[*int]string{}
	for i := range nums {
		m[&nums[i]] = fmt.Sprintf("idx_%d", i)
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyStruct tests pointer to struct as key
func MapPointerKeyStruct() string {
	type Node struct {
		Val int
	}
	n1 := Node{Val: 1}
	n2 := Node{Val: 2}
	m := map[*Node]string{
		&n1: "node1",
		&n2: "node2",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyArray tests pointer to array as key
func MapPointerKeyArray() string {
	arr1 := [3]int{1, 2, 3}
	arr2 := [3]int{1, 2, 3}
	m := map[*[3]int]string{
		&arr1: "first",
		&arr2: "second",
	}
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyIterate tests iterating over pointer-keyed map
func MapPointerKeyIterate() string {
	a, b, c := 1, 2, 3
	m := map[*int]int{
		&a: 10,
		&b: 20,
		&c: 30,
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}

// MapPointerKeyDelete tests deleting pointer keys
func MapPointerKeyDelete() string {
	x, y, z := 1, 2, 3
	m := map[*int]string{
		&x: "x",
		&y: "y",
		&z: "z",
	}
	delete(m, &y)
	return fmt.Sprintf("len=%d", len(m))
}

// MapPointerKeyNil tests nil pointer as key
func MapPointerKeyNil() string {
	m := map[*int]int{}
	var p *int = nil
	m[p] = 100
	v, ok := m[nil]
	return fmt.Sprintf("v=%d,ok=%t", v, ok)
}

// MapPointerKeyModifyValue modifies value through pointer key
func MapPointerKeyModifyValue() string {
	x := 5
	m := map[*int]string{&x: "original"}
	x = 10
	_, ok := m[&x]
	return fmt.Sprintf("x=%d,ok=%t", x, ok)
}

// MapPointerKeyComposite tests composite type with pointer
func MapPointerKeyComposite() string {
	type Item struct {
		ID   int
		Name string
	}
	items := []Item{{1, "a"}, {2, "b"}, {3, "c"}}
	m := map[*Item]int{}
	for i := range items {
		m[&items[i]] = items[i].ID * 10
	}
	sum := 0
	for _, v := range m {
		sum += v
	}
	return fmt.Sprintf("sum=%d", sum)
}
