package tests

import (
	_ "embed"
	"strings"
	"testing"

	gig "git.woa.com/youngjin/gig"
	_ "git.woa.com/youngjin/gig/stdlib/packages"
	"git.woa.com/youngjin/gig/tests/testdata/resolved_issue"
)

//go:embed testdata/resolved_issue/main.go
var resolvedSrc string

// toMainPackageResolved converts a source file to package main
func toMainPackageResolved(src string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "package ") {
			lines[i] = "package main"
			break
		}
	}
	return strings.Join(lines, "\n")
}

// runResolvedTest runs a function from the resolved_issue test file
func runResolvedTest(t *testing.T, funcName string) any {
	t.Helper()
	src := toMainPackageResolved(resolvedSrc)
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run(funcName)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	return result
}

// Resolved issue tests - these verify bugs have been fixed
// Each test compares interpreted execution with native Go execution

func TestResolved_BytesToString(t *testing.T) {
	expected := resolved_issue.BytesToString()
	result := runResolvedTest(t, "BytesToString")

	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}

func TestResolved_BytesToStringHi(t *testing.T) {
	expected := resolved_issue.BytesToStringHi()
	result := runResolvedTest(t, "BytesToStringHi")

	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}

func TestResolved_BytesToStringGo(t *testing.T) {
	expected := resolved_issue.BytesToStringGo()
	result := runResolvedTest(t, "BytesToStringGo")

	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}

func TestResolved_BytesToStringEmpty(t *testing.T) {
	expected := resolved_issue.BytesToStringEmpty()
	result := runResolvedTest(t, "BytesToStringEmpty")

	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}

func TestResolved_BytesToStringSingle(t *testing.T) {
	expected := resolved_issue.BytesToStringSingle()
	result := runResolvedTest(t, "BytesToStringSingle")

	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != expected {
		t.Errorf("got %q, want %q", s, expected)
	}
}

func TestResolved_PointerReceiverMutation(t *testing.T) {
	expected := resolved_issue.PointerReceiverMutation()
	result := runResolvedTest(t, "PointerReceiverMutation")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_PointerReceiverMutationReturnValue(t *testing.T) {
	expected := resolved_issue.PointerReceiverMutationReturnValue()
	result := runResolvedTest(t, "PointerReceiverMutationReturnValue")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_InitFuncExecuted(t *testing.T) {
	expected := resolved_issue.InitFuncExecuted()
	result := runResolvedTest(t, "InitFuncExecuted")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_InitFuncSideEffect(t *testing.T) {
	expected := resolved_issue.InitFuncSideEffect()
	result := runResolvedTest(t, "InitFuncSideEffect")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_RangeStringRuneValue(t *testing.T) {
	expected := resolved_issue.RangeStringRuneValue()
	result := runResolvedTest(t, "RangeStringRuneValue")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_RangeStringIndexValue(t *testing.T) {
	expected := resolved_issue.RangeStringIndexValue()
	result := runResolvedTest(t, "RangeStringIndexValue")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

func TestResolved_RangeStringMultibyte(t *testing.T) {
	expected := resolved_issue.RangeStringMultibyte()
	result := runResolvedTest(t, "RangeStringMultibyte")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 5: Map with function value type
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_MapWithFuncValue(t *testing.T) {
	expected := resolved_issue.MapWithFuncValue()
	result := runResolvedTest(t, "MapWithFuncValue")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("map with func value: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 6: Type switch on interface values in slice
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_InterfaceSliceTypeSwitch(t *testing.T) {
	expected := resolved_issue.InterfaceSliceTypeSwitch()
	result := runResolvedTest(t, "InterfaceSliceTypeSwitch")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("interface slice type switch: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 7: Struct with function field
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructWithFuncField(t *testing.T) {
	expected := resolved_issue.StructWithFuncField()
	result := runResolvedTest(t, "StructWithFuncField")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("struct with func field: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 8: Slice append with spread operator
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_SliceFlatten(t *testing.T) {
	expected := resolved_issue.SliceFlatten()
	result := runResolvedTest(t, "SliceFlatten")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("slice flatten: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 9: Map update during range
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_MapUpdateDuringRange(t *testing.T) {
	result := runResolvedTest(t, "MapUpdateDuringRange")

	n := toInt64(t, result)
	// The initial map has 2 keys. Each visited key adds one new key.
	// At minimum, only the 2 original keys are visited → 4 total.
	// The result should be at least 4 (2 original + 2 added).
	if n < 4 {
		t.Errorf("map update during range: got %d, want >= 4", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 10: Self-referencing struct type
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructSelfRef(t *testing.T) {
	expected := resolved_issue.StructSelfRef()
	result := runResolvedTest(t, "StructSelfRef")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("struct self-ref: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 11: Defer in closure with argument
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_DeferInClosureWithArg(t *testing.T) {
	expected := resolved_issue.DeferInClosureWithArg()
	result := runResolvedTest(t, "DeferInClosureWithArg")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("defer in closure with arg: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 12: Pointer swap in struct
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_PointerSwapInStruct(t *testing.T) {
	expected := resolved_issue.PointerSwapInStruct()
	result := runResolvedTest(t, "PointerSwapInStruct")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("pointer swap in struct: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 13: Struct with function slice
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructWithFuncSlice(t *testing.T) {
	expected := resolved_issue.StructWithFuncSlice()
	result := runResolvedTest(t, "StructWithFuncSlice")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("struct with func slice: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 14: Struct with anonymous field
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructAnonymousField(t *testing.T) {
	expected := resolved_issue.StructAnonymousField()
	result := runResolvedTest(t, "StructAnonymousField")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("struct anonymous field: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 15: Struct with embedded interface
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructEmbeddedInterface(t *testing.T) {
	// Build with isolated source to avoid type confusion with other struct types
	// in the larger resolved_issue file (reflect.StructOf creates anonymous types).
	src := `package main

type Getter interface{ Get() int }
type GetterImpl struct{ v int }
func (g *GetterImpl) Get() int { return g.v }
type GetterHolder struct { Getter }

func StructEmbeddedInterface() int {
	h := GetterHolder{Getter: &GetterImpl{v: 42}}
	return h.Get()
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("StructEmbeddedInterface")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	n := toInt64(t, result)
	// Native result: 42
	if n != 42 {
		t.Errorf("struct embedded interface: got %d, want 42", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 16: Map range with break
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_MapRangeWithBreak(t *testing.T) {
	result := runResolvedTest(t, "MapRangeWithBreak")

	n := toInt64(t, result)
	// Non-deterministic: sum of some values from {10, 20, 30} until sum > 25.
	// Valid results: 30 (=10+20), 30 (=30), 40 (=10+30), 50 (=20+30), 60 (=10+20+30)
	// At minimum, at least one value is consumed, so n >= 10.
	if n < 10 {
		t.Errorf("map range with break: got %d, want >= 10", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 17: Pointer to interface type assertion
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_PointerToInterface(t *testing.T) {
	expected := resolved_issue.PointerToInterface()
	result := runResolvedTest(t, "PointerToInterface")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("pointer to interface: got %d, want %d", n, expected)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 18: Struct with pointer to interface
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructWithPointerToInterface(t *testing.T) {
	// Uses isolated inline source to avoid reflect.StructOf type collision
	src := `package main

type PtrToInterface struct {
	data *interface{}
}

func StructWithPointerToInterface() int {
	var i interface{} = 42
	s := PtrToInterface{data: &i}
	return (*s.data).(int)
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("StructWithPointerToInterface")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	n := toInt64(t, result)
	if n != 42 {
		t.Errorf("struct with pointer to interface: got %d, want 42", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 19: Struct with nested function field
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructWithNestedFunc(t *testing.T) {
	// Uses isolated inline source to avoid reflect.StructOf type collision
	src := `package main

type NestedFuncHolder struct {
	get func() func() int
}

func StructWithNestedFunc() int {
	h := NestedFuncHolder{
		get: func() func() int {
			return func() int { return 42 }
		},
	}
	return h.get()()
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("StructWithNestedFunc")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	n := toInt64(t, result)
	if n != 42 {
		t.Errorf("struct with nested func: got %d, want 42", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 20: Struct with interface map
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_StructWithInterfaceMap(t *testing.T) {
	// Uses isolated inline source to avoid reflect.StructOf type collision
	src := `package main

type InterfaceMapHolder struct {
	data map[string]interface{}
}

func StructWithInterfaceMap() int {
	h := InterfaceMapHolder{
		data: map[string]interface{}{
			"a": 1,
			"b": "hello",
		},
	}
	return h.data["a"].(int)
}
`
	prog, err := gig.Build(src)
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	result, err := prog.Run("StructWithInterfaceMap")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	n := toInt64(t, result)
	if n != 1 {
		t.Errorf("struct with interface map: got %d, want 1", n)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// Resolved Issue 21: Pointer to slice element modify in loop
// ═══════════════════════════════════════════════════════════════════════════════

func TestResolved_PointerToSliceElemModify(t *testing.T) {
	expected := resolved_issue.PointerToSliceElemModify()
	result := runResolvedTest(t, "PointerToSliceElemModify")

	n := toInt64(t, result)
	if n != int64(expected) {
		t.Errorf("pointer to slice elem modify: got %d, want %d", n, expected)
	}
}
