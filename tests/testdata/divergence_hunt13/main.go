package divergence_hunt13

import "fmt"

// ============================================================================
// Round 13: Struct edge cases, interface edge cases, type system edge cases,
// conversion edge cases, pointer edge cases
// ============================================================================

// StructZeroValue tests struct zero value
func StructZeroValue() int {
	type P struct{ X, Y int }
	var p P
	return p.X + p.Y
}

// StructPointerNil tests nil struct pointer
func StructPointerNil() bool {
	type P struct{ X int }
	var p *P
	return p == nil
}

// StructCopyOnAssign tests struct copy semantics
func StructCopyOnAssign() int {
	type P struct{ X int }
	a := P{X: 10}
	b := a
	b.X = 20
	return a.X
}

// StructFieldAccess tests struct field access patterns
func StructFieldAccess() int {
	type P struct{ X, Y, Z int }
	p := P{X: 1, Y: 2, Z: 3}
	return p.X + p.Y + p.Z
}

// InterfaceNilComparison tests nil interface comparison
func InterfaceNilComparison() bool {
	var x any
	return x == nil
}

// InterfaceTypedNil tests typed nil (not nil interface)
func InterfaceTypedNil() bool {
	var s []int
	var x any = s
	// x is not nil even though s is nil
	return x != nil
}

// TypeAssertionWithBool tests type assertion returning bool
func TypeAssertionWithBool() int {
	var x any = 42
	if v, ok := x.(int); ok {
		return v
	}
	return -1
}

// MultipleTypeAssertions tests multiple type assertions
func MultipleTypeAssertions() int {
	var x any = "hello"
	if _, ok := x.(int); ok { return 1 }
	if v, ok := x.(string); ok { return len(v) }
	return -1
}

// PointerToStruct tests pointer to struct
func PointerToStruct() int {
	type P struct{ X int }
	p := &P{X: 42}
	return p.X
}

// PointerToStructModify tests pointer to struct modification
func PointerToStructModify() int {
	type P struct{ X int }
	p := &P{X: 10}
	p.X = 20
	return p.X
}

// StructAsMapValue tests struct as map value
func StructAsMapValue() int {
	type P struct{ X, Y int }
	m := map[string]P{"a": {1, 2}, "b": {3, 4}}
	return m["a"].X + m["b"].Y
}

// StructInSlice tests struct in slice
func StructInSlice() int {
	type P struct{ X int }
	s := []P{{1}, {2}, {3}}
	return s[0].X + s[1].X + s[2].X
}

// IntTypeAlias tests type alias
func IntTypeAlias() int {
	type MyInt int
	var x MyInt = 42
	return int(x)
}

// StringTypeAlias tests string type alias
func StringTypeAlias() int {
	type MyStr string
	var s MyStr = "hello"
	return len(s)
}

// SliceOfAlias tests slice of alias type
func SliceOfAlias() int {
	type MyInt int
	s := []MyInt{1, 2, 3}
	return int(s[0] + s[1] + s[2])
}

// NestedTypeDefinitions tests nested type definitions
func NestedTypeDefinitions() int {
	type Inner int
	type Outer struct{ V Inner }
	o := Outer{V: Inner(42)}
	return int(o.V)
}

// FmtStruct tests fmt formatting of struct
func FmtStruct() string {
	type P struct{ X, Y int }
	p := P{1, 2}
	return fmt.Sprintf("%v", p)
}

// FmtPointer tests fmt formatting of pointer
func FmtPointer() string {
	x := 42
	p := &x
	return fmt.Sprintf("%d", *p)
}

// ConversionBetweenNumericTypes tests conversion between numeric types
func ConversionBetweenNumericTypes() int64 {
	var a int8 = 10
	var b int16 = int16(a)
	var c int32 = int32(b)
	var d int64 = int64(c)
	return d
}

// UnsignedToSigned tests unsigned to signed conversion
func UnsignedToSigned() int {
	var u uint = 42
	return int(u)
}
