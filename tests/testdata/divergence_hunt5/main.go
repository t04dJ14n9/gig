package divergence_hunt5

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Round 5: Error handling patterns, fmt formatting, string manipulation,
// slice tricks, closure patterns, interface composition
// ============================================================================

// ErrorNew tests errors.New
func ErrorNew() string {
	return errors.New("test").Error()
}

// Errorf tests fmt.Errorf
func Errorf() string {
	return fmt.Errorf("value %d", 42).Error()
}

// FmtPrintln tests fmt.Sprintln
func FmtPrintln() string {
	return strings.TrimSpace(fmt.Sprintln("hello", "world"))
}

// FmtIntWidth tests fmt integer width
func FmtIntWidth() string {
	return fmt.Sprintf("%05d", 42)
}

// FmtFloat tests fmt float formatting
func FmtFloat() string {
	return fmt.Sprintf("%.2f", 3.14159)
}

// FmtBool tests fmt boolean
func FmtBool() string {
	return fmt.Sprintf("%t", true)
}

// FmtHex tests fmt hex
func FmtHex() string {
	return fmt.Sprintf("%x", 255)
}

// FmtOctal tests fmt octal
func FmtOctal() string {
	return fmt.Sprintf("%o", 8)
}

// FmtBinary tests fmt binary
func FmtBinary() string {
	return fmt.Sprintf("%b", 10)
}

// FmtChar tests fmt char
func FmtChar() string {
	return fmt.Sprintf("%c", 65)
}

// FmtStringWidth tests fmt string width
func FmtStringWidth() string {
	return fmt.Sprintf("%10s", "hi")
}

// SliceFilter tests filtering a slice
func SliceFilter() int {
	s := []int{1, 2, 3, 4, 5, 6}
	var result []int
	for _, v := range s {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return len(result)
}

// SliceMap tests mapping a slice
func SliceMap() int {
	s := []int{1, 2, 3}
	result := make([]int, len(s))
	for i, v := range s {
		result[i] = v * 2
	}
	return result[0] + result[1] + result[2]
}

// ClosureSum tests closure for sum
func ClosureSum() int {
	s := []int{1, 2, 3, 4, 5}
	acc := 0
	add := func(n int) { acc += n }
	for _, v := range s {
		add(v)
	}
	return acc
}

// ClosureCapture tests closure capturing variable
func ClosureCapture() int {
	x := 10
	fn := func() int { return x }
	x = 20
	return fn()
}

// InterfaceSlice tests slice of different interfaces
func InterfaceSlice() int {
	type Stringer interface{ String() string }
	items := []Stringer{}
	return len(items)
}

// MultipleReturnIgnore tests ignoring some return values
func MultipleReturnIgnore() int {
	divmod := func(a, b int) (int, int) { return a / b, a % b }
	q, _ := divmod(17, 5)
	return q
}

// NamedReturn tests named return
func NamedReturn() (result int) {
	result = 42
	return
}

// NamedReturnBare tests bare return
func NamedReturnBare() (result int) {
	defer func() { result++ }()
	return 10
}

// StringJoinInts tests joining ints as strings
func StringJoinInts() string {
	parts := []string{}
	for i := 0; i < 5; i++ {
		parts = append(parts, fmt.Sprintf("%d", i))
	}
	return strings.Join(parts, ",")
}

// MapStringSlice tests map of string to slice
func MapStringSlice() int {
	m := map[string][]int{}
	m["a"] = []int{1, 2, 3}
	m["b"] = []int{4, 5}
	return len(m["a"]) + len(m["b"])
}

// NestedStruct tests nested struct
func NestedStruct() int {
	type Inner struct{ X int }
	type Outer struct{ I Inner }
	o := Outer{I: Inner{X: 42}}
	return o.I.X
}

// StructLiteral tests struct literal with field names
func StructLiteral() int {
	type Point struct{ X, Y int }
	p := Point{X: 10, Y: 20}
	return p.X + p.Y
}

// StructPointer tests struct pointer access
func StructPointer() int {
	type Point struct{ X, Y int }
	p := &Point{X: 10, Y: 20}
	p.X = 30
	return p.X + p.Y
}

// DeferReturn tests defer modifying return value
func DeferReturn() (result int) {
	defer func() { result += 10 }()
	return 5
}

// DeferClosure tests defer with closure
func DeferClosure() (result int) {
	x := 10
	defer func() { result = x }()
	x = 20
	return 0
}

// StringEqual tests string equality
func StringEqual() bool {
	return "hello" == "hello" && "hello" != "world"
}

// MapLookup tests map lookup
func MapLookup() int {
	m := map[string]int{"a": 1, "b": 2}
	if v, ok := m["c"]; ok {
		return v
	}
	return -1
}
