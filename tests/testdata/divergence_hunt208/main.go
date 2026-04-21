package divergence_hunt208

import (
	"fmt"
)

// ============================================================================
// Round 208: Struct field ordering and access patterns
// ============================================================================

type SimpleStruct208 struct {
	A int8
	B int32
}

type PackedStruct208 struct {
	A int32
	B int8
	C int32
}

type AlignedStruct208 struct {
	X int64
	Y int32
	Z int16
}

type BoolStruct208 struct {
	A bool
	B int64
	C bool
}

type PointerStruct208 struct {
	A *int
	B int32
}

type SliceStruct208 struct {
	A []int
	B int32
}

type StringStruct208 struct {
	A string
	B int32
}

type NestedAlign208 struct {
	Inner SimpleStruct208
	Value int64
}

// SimpleStructAccess tests simple struct field access
func SimpleStructAccess() string {
	s := SimpleStruct208{A: 1, B: 2}
	return fmt.Sprintf("A=%d,B=%d", s.A, s.B)
}

// PackedStructAccess tests packed struct field access
func PackedStructAccess() string {
	s := PackedStruct208{A: 1, B: 2, C: 3}
	return fmt.Sprintf("A=%d,B=%d,C=%d", s.A, s.B, s.C)
}

// AlignedStructAccess tests aligned struct field access
func AlignedStructAccess() string {
	s := AlignedStruct208{X: 1, Y: 2, Z: 3}
	return fmt.Sprintf("X=%d,Y=%d,Z=%d", s.X, s.Y, s.Z)
}

// BoolStructAccess tests struct with bool fields
func BoolStructAccess() string {
	s := BoolStruct208{A: true, B: 42, C: false}
	return fmt.Sprintf("A=%v,B=%d,C=%v", s.A, s.B, s.C)
}

// PointerStructAccess tests struct with pointer field
func PointerStructAccess() string {
	val := 42
	s := PointerStruct208{A: &val, B: 10}
	return fmt.Sprintf("A=%d,B=%d", *s.A, s.B)
}

// SliceStructAccess tests struct with slice field
func SliceStructAccess() string {
	s := SliceStruct208{A: []int{1, 2, 3}, B: 10}
	return fmt.Sprintf("len(A)=%d,B=%d", len(s.A), s.B)
}

// StringStructAccess tests struct with string field
func StringStructAccess() string {
	s := StringStruct208{A: "hello", B: 10}
	return fmt.Sprintf("A=%s,B=%d", s.A, s.B)
}

// NestedStructAccess tests nested struct access
func NestedStructAccess() string {
	s := NestedAlign208{
		Inner: SimpleStruct208{A: 1, B: 2},
		Value: 3,
	}
	return fmt.Sprintf("Inner.A=%d,Inner.B=%d,Value=%d", s.Inner.A, s.Inner.B, s.Value)
}

// StructArrayAccess tests struct array access
func StructArrayAccess() string {
	arr := [3]SimpleStruct208{
		{A: 1, B: 2},
		{A: 3, B: 4},
		{A: 5, B: 6},
	}
	return fmt.Sprintf("arr[0]=(%d,%d),arr[1]=(%d,%d),arr[2]=(%d,%d)",
		arr[0].A, arr[0].B, arr[1].A, arr[1].B, arr[2].A, arr[2].B)
}

// StructSliceAccess tests struct slice access
func StructSliceAccess() string {
	slice := []SimpleStruct208{
		{A: 1, B: 2},
		{A: 3, B: 4},
	}
	return fmt.Sprintf("len=%d,slice[0]=(%d,%d)", len(slice), slice[0].A, slice[0].B)
}

// StructEquality tests struct equality
func StructEquality() string {
	a := SimpleStruct208{A: 1, B: 2}
	b := SimpleStruct208{A: 1, B: 2}
	c := SimpleStruct208{A: 2, B: 1}
	return fmt.Sprintf("a==b:%v,a==c:%v", a == b, a == c)
}
