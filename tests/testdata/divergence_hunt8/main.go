package divergence_hunt8

import "sync"

// ============================================================================
// Round 8: Sync primitives, mutex, once, wait group patterns (goroutine-free),
// complex struct hierarchies, slice-of-slice, map-of-map
// ============================================================================

// MutexBasic tests basic mutex lock/unlock
func MutexBasic() int {
	var mu sync.Mutex
	x := 0
	mu.Lock()
	x++
	mu.Unlock()
	return x
}

// OnceBasic tests sync.Once
func OnceBasic() int {
	var once sync.Once
	count := 0
	once.Do(func() { count++ })
	once.Do(func() { count++ })
	return count
}

// SliceOfSlice tests 2D slice
func SliceOfSlice() int {
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	return grid[0][1] + grid[1][2]
}

// MapOfMap tests 2D map
func MapOfMap() int {
	m := map[int]map[int]string{}
	m[1] = map[int]string{10: "a", 20: "b"}
	m[2] = map[int]string{30: "c"}
	return len(m[1]) + len(m[2])
}

// StructWithSlice tests struct with slice field
func StructWithSlice() int {
	type Team struct {
		Name    string
		Members []string
	}
	t := Team{Name: "A", Members: []string{"alice", "bob", "charlie"}}
	return len(t.Members)
}

// StructWithMap tests struct with map field
func StructWithMap() int {
	type Config struct {
		Values map[string]int
	}
	c := Config{Values: map[string]int{"x": 1, "y": 2}}
	return c.Values["x"] + c.Values["y"]
}

// NestedSliceAppend tests appending to nested slices
func NestedSliceAppend() int {
	matrix := [][]int{}
	matrix = append(matrix, []int{1, 2})
	matrix = append(matrix, []int{3, 4})
	return matrix[0][0] + matrix[1][1]
}

// DeepStruct tests deeply nested struct
func DeepStruct() int {
	type A struct{ X int }
	type B struct{ A A }
	type C struct{ B B }
	c := C{B: B{A: A{X: 42}}}
	return c.B.A.X
}

// SliceOfStructAppend tests appending structs to slice
func SliceOfStructAppend() int {
	type P struct{ X, Y int }
	pts := []P{}
	pts = append(pts, P{1, 2})
	pts = append(pts, P{3, 4})
	return pts[0].X + pts[1].Y
}

// MapWithSliceValue tests map with slice value
func MapWithSliceValue() int {
	m := map[string][]int{}
	m["a"] = []int{1, 2}
	m["a"] = append(m["a"], 3)
	return len(m["a"])
}

// MutexInDefer tests mutex with defer
func MutexInDefer() int {
	var mu sync.Mutex
	x := 0
	func() {
		mu.Lock()
		defer mu.Unlock()
		x++
	}()
	return x
}

// RWMutexBasic tests RWMutex
func RWMutexBasic() int {
	var mu sync.RWMutex
	x := 0
	mu.Lock()
	x = 42
	mu.Unlock()
	mu.RLock()
	v := x
	mu.RUnlock()
	return v
}

// StructWithFunc tests struct with function field
func StructWithFunc() int {
	type Op struct {
		Apply func(int) int
	}
	double := Op{Apply: func(x int) int { return x * 2 }}
	return double.Apply(21)
}

// StructWithPointer tests struct with pointer field
func StructWithPointer() int {
	type Node struct {
		Value int
		Next  *Node
	}
	n2 := &Node{Value: 2, Next: nil}
	n1 := &Node{Value: 1, Next: n2}
	return n1.Value + n1.Next.Value
}

// SliceGrowPattern tests slice growth pattern
func SliceGrowPattern() int {
	var s []int
	for i := 0; i < 100; i++ {
		s = append(s, i)
	}
	return s[99]
}

// MapGrowPattern tests map growth
func MapGrowPattern() int {
	m := map[int]int{}
	for i := 0; i < 100; i++ {
		m[i] = i * 2
	}
	return m[50]
}

// CompositeLiteralNested tests nested composite literals
func CompositeLiteralNested() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 10}}
	return o.I.X
}
