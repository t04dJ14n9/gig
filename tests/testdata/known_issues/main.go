package known_issues

import (
	"sort"
	"strings"
)

// ============================================================================
// Known Issue 1: Passing stdlib functions as callback arguments
//
// Assigning a stdlib function to a variable and calling it directly works:
//
//	f := strings.TrimSpace
//	f("  hello  ")  // works, returns "hello"
//
// But passing it as a callback argument fails — the callee receives nil:
//
//	process("hello", strings.TrimSpace)  // f is nil inside process
// ============================================================================

func StdlibFuncAsCallback() string {
	process := func(s string, f func(string) string) string {
		return f(s)
	}
	return process("  hello  ", strings.TrimSpace)
}

// ============================================================================
// Known Issue 2: sort.IntSlice composite literal passed to sort.Sort
//
// Converting an existing []int slice via sort.IntSlice(s) works fine, but
// creating a sort.IntSlice via composite literal fails:
//
//	s := sort.IntSlice{5, 3, 1, 4, 2}
//	sort.Sort(s)  // panic: interface conversion: []int is not sort.Interface
//
// The composite literal creates a plain []int without sort.Interface methods.
// ============================================================================

func SortIntSliceCompositeLiteral() int {
	s := sort.IntSlice{5, 3, 1, 4, 2}
	sort.Sort(s)
	return int(s[0])*10000 + int(s[1])*1000 + int(s[2])*100 + int(s[3])*10 + int(s[4])
}

// ============================================================================
// Known Issue 3: Typed nil assigned to interface is incorrectly treated as nil
//
// In Go, assigning a typed nil pointer to an interface variable makes the
// interface NOT nil (it has a type but the value is nil). Gig incorrectly
// treats the interface as nil.
//
//	var err *MyError  // nil pointer
//	var e error = err // typed nil assigned to interface
//	e == nil          // should be false in Go, but Gig returns true
// ============================================================================

type MyError3 struct {
	msg string
}

func (e *MyError3) Error() string {
	return e.msg
}

func TypedNilInterface() string {
	var err *MyError3 // nil pointer
	var e error = err // typed nil assigned to interface
	if e == nil {
		return "nil"
	}
	return "not nil"
}

// ============================================================================
// Known Issue 4: Method call on typed nil interface doesn't panic
//
// Calling a method on a typed nil pointer through an interface should
// panic with nil pointer dereference. Gig incorrectly executes the
// method without panicking.
// ============================================================================

func InterfaceMethodOnTypedNil() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	var p *MyError3 // typed nil
	var e error = p
	_ = e.Error() // should panic - nil pointer dereference
	return "ok"
}

// ============================================================================
// Known Issue 5: Named type arithmetic fails
//
// Arithmetic operations on named types (type Celsius float64) fail
// with "cannot mul invalid" runtime panic.
//
//	var c Celsius = 100
//	c * 9  // should work, but Gig panics
// ============================================================================

type Celsius float64
type Fahrenheit float64

func CToF() Fahrenheit {
	c := Celsius(100)
	return Fahrenheit(c*9/5 + 32)
}

// ============================================================================
// Known Issue 6: Type assertion comma-ok returns nil instead of zero value
//
// When a type assertion with comma-ok fails, the value should be the
// zero value of the target type, not nil.
//
//	var x any = "hello"
//	v, ok := x.(int)  // v should be 0, but Gig returns nil
// ============================================================================

func AssignTypeAssertion() []any {
	var x any = "hello"
	v, ok := x.(int)
	return []any{v, ok}
}

// ============================================================================
// Known Issue 7: nil slice subslice [0:0] returns non-nil slice
//
// Sub-slicing a nil slice with [0:0] should return a nil slice,
// but Gig returns an empty non-nil slice.
// ============================================================================

func SliceNilSubslice() []int {
	var s []int
	return s[0:0]
}

// ============================================================================
// Known Issue 8: Linked list reverse with pointer reassignment fails
//
// Reversing a linked list by reassigning .Next pointers on struct
// nodes doesn't work correctly in Gig. The interpreter appears to
// not properly propagate pointer-based struct field modifications
// through the linked list traversal.
// ============================================================================

type ListNode8 struct {
	Val  int
	Next *ListNode8
}

func LinkedListReverse() int {
	head := &ListNode8{Val: 1, Next: &ListNode8{Val: 2, Next: &ListNode8{Val: 3}}}
	var prev *ListNode8
	for n := head; n != nil; {
		next := n.Next
		n.Next = prev
		prev = n
		n = next
	}
	return prev.Val // should be 3
}

// ============================================================================
// Known Issue 9: Package-level variable initialization
//
// Package-level var declarations with non-zero initializers are not
// properly initialized in Gig. Variables like `var s = "hello"` or
// `var m = map[string]int{"key": 42}` appear as zero values.
// ============================================================================

var GlobalSlice = []int{1, 2, 3}
var GlobalMap = map[string]int{"key": 42}
var GlobalString = "hello"
var GlobalBool = true
var GlobalFloat = 3.14
var GlobalPointer *int

func GlobalSliceAccess() int {
	return len(GlobalSlice)
}

func GlobalMapAccess() int {
	return GlobalMap["key"]
}

func GlobalStringAccess() string {
	return GlobalString
}

func GlobalBoolAccess() bool {
	return GlobalBool
}

func GlobalFloatAccess() float64 {
	return GlobalFloat
}

func GlobalPointerNil() bool {
	return GlobalPointer == nil
}
