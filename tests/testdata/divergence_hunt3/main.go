package divergence_hunt3

import "strings"

// ============================================================================
// Round 3: Interface methods, string builder, const blocks, type aliases,
// multiple assignment, select, goroutine-free concurrency patterns,
// rune iteration, nested maps, struct methods
// ============================================================================

// String Builder pattern
func StringBuilder() string {
	var b strings.Builder
	b.WriteString("hello")
	b.WriteString(" ")
	b.WriteString("world")
	return b.String()
}

// Const block
func ConstBlock() int {
	const (
		A = 1
		B = 2
		C = 3
	)
	return A + B + C
}

// Iota enum
func IotaEnum() int {
	const (
		Sunday = iota
		Monday
		Tuesday
		Wednesday
	)
	return Wednesday
}

// Multiple assignment
func MultipleAssign() int {
	x, y := 1, 2
	x, y = y, x
	return x*10 + y
}

// NestedMap tests map of maps
func NestedMap() int {
	m := map[string]map[string]int{}
	m["outer"] = map[string]int{"inner": 42}
	return m["outer"]["inner"]
}

// RuneIteration tests ranging over string runes
func RuneIteration() int {
	s := "Hello, 世界"
	count := 0
	for _, r := range s {
		if r > 127 {
			count++
		}
	}
	return count
}

// StringIndexRune tests strings.IndexRune
func StringIndexRune() int {
	return strings.IndexRune("hello世界", '界')
}

// StringCount tests strings.Count
func StringCount() int { return strings.Count("banana", "a") }

// MapBoolKey tests map with bool key
func MapBoolKey() int {
	m := map[bool]int{true: 1, false: 0}
	return m[true] + m[false]
}

// SliceReverse tests slice reversal
func SliceReverse() int {
	s := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s[0]*10000 + s[1]*1000 + s[2]*100 + s[3]*10 + s[4]
}

// StructMethod tests struct with method
func StructMethod() int {
	type Counter struct{ n int }
	c := Counter{n: 5}
	c.n += 10
	return c.n
}

// InterfaceEmpty tests empty interface
func InterfaceEmpty() int {
	var x any = 42
	return x.(int)
}

// InterfaceNil tests nil interface
func InterfaceNil() bool {
	var x any = nil
	return x == nil
}

// SliceOfInterface tests slice of interface
func SliceOfInterface() int {
	s := []any{1, "hello", true}
	return len(s)
}

// MapWithStructValue tests map with struct value
func MapWithStructValue() int {
	type Point struct{ X, Y int }
	m := map[string]Point{"origin": {0, 0}, "p1": {1, 2}}
	return m["p1"].X + m["p1"].Y
}

// StringFields tests strings.Fields
func StringFields() int { return len(strings.Fields("  hello   world  ")) }

// StringRepeat tests strings.Repeat
func StringRepeat() string { return strings.Repeat("ab", 3) }

// StringMap tests mapping over string
func StringMap() string {
	s := "hello"
	result := strings.Map(func(r rune) rune {
		if r == 'l' { return 'L' }
		return r
	}, s)
	return result
}

// MapStructKey tests map with struct key
func MapStructKey() int {
	type Key struct{ X, Y int }
	m := map[Key]string{{1, 2}: "a", {3, 4}: "b"}
	return len(m)
}

// SliceMinMax tests finding min/max in slice
func SliceMinMax() int {
	s := []int{5, 3, 8, 1, 9, 2}
	min, max := s[0], s[0]
	for _, v := range s {
		if v < min { min = v }
		if v > max { max = v }
	}
	return min*10 + max
}

// NestedIf tests nested if
func NestedIf() int {
	x := 15
	if x > 10 {
		if x > 20 {
			return 3
		}
		return 2
	}
	return 1
}

// StringToLower tests strings.ToLower
func StringToLower() string { return strings.ToLower("HELLO World") }

// StringToUpper tests strings.ToUpper
func StringToUpper() string { return strings.ToUpper("hello world") }

// ContinueLoop tests continue in loop
func ContinueLoop() int {
	sum := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 { continue }
		sum += i
	}
	return sum
}

// LabeledBreak tests labeled break
func LabeledBreak() int {
	sum := 0
outer:
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if i+j > 5 { break outer }
			sum += i + j
		}
	}
	return sum
}

// SliceMakeZero tests make([]T, 0)
func SliceMakeZero() int {
	s := make([]int, 0)
	s = append(s, 1)
	return len(s)
}

// ArrayIteration tests array range
func ArrayIteration() int {
	a := [3]int{10, 20, 30}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

// Float32Arith tests float32 arithmetic
func Float32Arith() float32 {
	var a float32 = 3.14
	var b float32 = 2.0
	return a * b
}

// Int8Arith tests int8 arithmetic
func Int8Arith() int8 {
	var a int8 = 10
	var b int8 = 20
	return a + b
}

// Uint16Arith tests uint16 arithmetic
func Uint16Arith() uint16 {
	var a uint16 = 100
	var b uint16 = 200
	return a + b
}
