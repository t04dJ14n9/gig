package known_issues

// ─────────────────────────────────────────────────────────────────────────────
// KNOWN ISSUES - Tests that STILL FAIL due to interpreter bugs.
//
// Each function documents the bug, expected (native) behavior, and actual
// (interpreter) behavior. The test harness in known_issue_test.go runs each
// function through both the interpreter and native Go, then reports mismatches.
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"sort"
	"strconv"
)

// ============================================================================
// BUG 1: Named-type conversion to external sort types panics
//
// sort.IntSlice, sort.Float64Slice, and sort.StringSlice are named types
// (e.g. type IntSlice []int). The VM keeps the underlying slice type
// ([]int) and does not convert it to the named type (sort.IntSlice).
// When passed to sort.Sort(sort.Interface), the interface assertion fails:
//   panic: interface conversion: []int is not sort.Interface: missing method Len
// ============================================================================

// SortIntSlice converts []int to sort.IntSlice and sorts.
// Native: returns 1 (sorted). Interpreter: PANIC.
func SortIntSlice() int {
	s := []int{3, 1, 2}
	sort.Sort(sort.IntSlice(s))
	return s[0]
}

// SortFloat64Slice converts []float64 to sort.Float64Slice and sorts.
// Native: returns 1. Interpreter: PANIC.
func SortFloat64Slice() int {
	s := []float64{3.0, 1.0, 2.0}
	sort.Sort(sort.Float64Slice(s))
	if s[0] == 1.0 {
		return 1
	}
	return 0
}

// SortStringSlice converts []string to sort.StringSlice and sorts.
// Native: returns 1. Interpreter: PANIC.
func SortStringSlice() int {
	s := []string{"c", "a", "b"}
	sort.Sort(sort.StringSlice(s))
	if s[0] == "a" {
		return 1
	}
	return 0
}

// SortReverse wraps sort.IntSlice in sort.Reverse.
// Native: returns 5. Interpreter: PANIC (same root cause).
func SortReverse() int {
	s := []int{1, 2, 3, 4, 5}
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	return s[0]
}

// SortIntsInPlace calls sort.Ints which mutates []int in-place.
// The VM stores []int as []int64; the reflect-based dispatch creates a
// converted []int copy, sorts it, but the original []int64 is unchanged.
// Native: returns 1. Interpreter: returns 0 (unsorted).
func SortIntsInPlace() int {
	s := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Ints(s)
	if s[0] == 1 && s[7] == 9 {
		return 1
	}
	return 0
}

// ============================================================================
// BUG 2: FIXED — time.Duration DirectCall wrappers now use cast
//
// The generator was producing args[i].Interface().(time.Duration) which panicked
// because the VM stores time.Duration as int64. The fix generates
// time.Duration(args[i].Int()) for cross-package named types with basic underlying.
// context.WithTimeout, time.Duration methods, etc. now work correctly.
// ============================================================================

// ============================================================================
// BUG 3: fmt.Stringer interface not honored on interpreted types
//
// When an interpreted struct implements fmt.Stringer (has a String() method),
// fmt.Sprintf("%v", val) should call String(). Instead, it prints the raw
// struct fields using the reflect-synthesized struct layout, including the
// internal _gig_id sentinel field: "{42 {}}" instead of "custom".
// ============================================================================

// FmtStringerNotCalled has a String() method but fmt ignores it.
// Native: returns "custom". Interpreter: returns "{42 {}}".
func FmtStringerNotCalled() string {
	return fmt.Sprintf("%v", stringerVal{42})
}

type stringerVal struct{ N int }

func (v stringerVal) String() string { return "custom" }

// ============================================================================
// BUG 4: fmt.Sprintf %T reports synthesized struct type, not declared name
//
// The VM uses reflect.StructOf() to create runtime types for interpreted
// structs, adding a _gig_id sentinel field. %T reports the synthesized
// anonymous struct type rather than the declared Go type name.
// ============================================================================

// FmtSprintfTypeWrong shows %T for an interpreted struct.
// Native: returns "known_issues.point". Interpreter: returns
// "struct { X int; Y int; _gig_id struct {} }".
func FmtSprintfTypeWrong() string {
	return fmt.Sprintf("%T", point{1, 2})
}

type point struct{ X, Y int }

// ============================================================================
// BUG 5: Interpreted structs have extra _gig_id field in %v output
//
// Because the VM synthesizes struct types with an extra _gig_id sentinel
// field for type identity, fmt.Sprintf("%v", val) includes it:
// "{1 2 {}}" instead of "{1 2}".
// ============================================================================

// FmtSprintfExtraField shows %v for an interpreted struct.
// Native: returns "{1 2}". Interpreter: returns "{1 2 {}}".
func FmtSprintfExtraField() string {
	return fmt.Sprintf("%v", point{1, 2})
}

// ============================================================================
// BUG 6: prog.Run() narrows int64→int and uint64→uint return types
//
// When an interpreted function declares a return type of int64 or uint64,
// the VM's value.Value stores the value correctly internally, but the
// Interface() conversion (used by prog.Run to return any) narrows int64
// to int and uint64 to uint. This means the Go type of the returned
// interface{} differs from the declared function return type.
//
// Inside the interpreter, the values work correctly (arithmetic, %T shows
// int64). The bug only manifests at the prog.Run() boundary.
// ============================================================================

// StrconvParseIntNarrowed returns int64 but prog.Run() returns int.
// Native: 42 (int64). Interpreter: 42 (int).
func StrconvParseIntNarrowed() int64 {
	n, _ := strconv.ParseInt("42", 10, 64)
	return n
}

// StrconvParseUintNarrowed returns uint64 but prog.Run() returns uint.
// Native: 42 (uint64). Interpreter: 42 (uint).
func StrconvParseUintNarrowed() uint64 {
	n, _ := strconv.ParseUint("42", 10, 64)
	return n
}
