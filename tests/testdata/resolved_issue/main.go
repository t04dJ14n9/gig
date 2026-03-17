package resolved_issue

// Resolved Issues - Tests that previously failed but now pass
// These verify that bugs have been fixed and no regressions occur.

// ── Resolved Issue 1: string([]byte{...}) conversion ───────────────────────────

// BytesToString tests basic byte to string conversion
func BytesToString() string {
	return string([]byte{104, 105})
}

// BytesToStringHi tests "hi" conversion
func BytesToStringHi() string {
	return string([]byte{104, 105})
}

// BytesToStringGo tests "Go" conversion
func BytesToStringGo() string {
	return string([]byte{71, 111})
}

// BytesToStringEmpty tests empty byte slice conversion
func BytesToStringEmpty() string {
	return string([]byte{})
}

// BytesToStringSingle tests single byte conversion
func BytesToStringSingle() string {
	return string([]byte{65})
}

// ── Resolved Issue 2: Pointer-receiver method mutation ─────────────────────────

type Counter struct{ n int }

func (c *Counter) Inc() { c.n++ }

// PointerReceiverMutation tests pointer receiver mutation
func PointerReceiverMutation() int {
	c := &Counter{}
	c.Inc()
	c.Inc()
	return c.n
}

type Box struct{ val int }

func (b *Box) Set(v int) { b.val = v }
func (b Box) Get() int   { return b.val }

// PointerReceiverMutationReturnValue tests Set and Get consistency
func PointerReceiverMutationReturnValue() int {
	b := &Box{}
	b.Set(99)
	return b.Get()
}

// ── Resolved Issue 3: init() execution ───────────────────────────────────────

var initVal int

func init() { initVal = 42 }

// InitFuncExecuted tests init() execution
func InitFuncExecuted() int { return initVal }

var registry []string

func init() {
	registry = append(registry, "alpha")
	registry = append(registry, "beta")
}

// InitFuncSideEffect tests init() side effects
func InitFuncSideEffect() int { return len(registry) }

// ── Resolved Issue 4: range-over-string rune values ───────────────────────────

// RangeStringRuneValue tests rune values from range over string
func RangeStringRuneValue() int {
	sum := 0
	for _, r := range "abc" {
		sum += int(r)
	}
	return sum
}

// RangeStringIndexValue tests index values from range over string
func RangeStringIndexValue() int {
	sum := 0
	for i := range "xyz" {
		sum += i
	}
	return sum
}

// RangeStringMultibyte tests multibyte rune values
func RangeStringMultibyte() int {
	sum := 0
	for _, r := range "中文" {
		sum += int(r)
	}
	return sum
}

// ── Resolved Issue 5: Map with function value type ────────────────────────────

// MapWithFuncValue tests storing closures in a map.
// Previously panicked: reflect.Value.SetMapIndex: value of type *vm.Closure
// is not assignable to type func() int.
// Expected: 30
func MapWithFuncValue() int {
	m := make(map[int]func() int)
	m[1] = func() int { return 10 }
	m[2] = func() int { return 20 }
	return m[1]() + m[2]()
}

// ── Resolved Issue 6: Type switch on interface values in slice ────────────────

// InterfaceSliceTypeSwitch tests type switch on interface slice elements.
// Previously returned 0 because type switch always fell through to default.
// Expected: 1111 (1 + 10 + 100 + 1000)
func InterfaceSliceTypeSwitch() int {
	var items []interface{}
	items = append(items, 1, "hello", true, 3.14)
	count := 0
	for _, item := range items {
		switch item.(type) {
		case int:
			count += 1
		case string:
			count += 10
		case bool:
			count += 100
		case float64:
			count += 1000
		}
	}
	return count
}

// ── Resolved Issue 7: Struct with function field ──────────────────────────────

type structWithFunc struct {
	f func() int
}

// StructWithFuncField tests struct with function field.
// Previously panicked: reflect.Set: value of type value.Value is not
// assignable to type func() int.
// Expected: 42
func StructWithFuncField() int {
	s := structWithFunc{f: func() int { return 42 }}
	return s.f()
}

// ── Resolved Issue 8: Slice append with spread operator ──────────────────────

// SliceFlatten tests appending slice with spread operator in a loop.
// Previously only appended the first element, resulting in wrong length.
// Expected: 4
func SliceFlatten() int {
	s := [][]int{{1, 2}, {3, 4}}
	result := []int{}
	for _, inner := range s {
		result = append(result, inner...)
	}
	return len(result)
}

// ── Resolved Issue 9: Map update during range ────────────────────────────────

// MapUpdateDuringRange tests modifying map during range iteration.
// Previously the map iteration count was incorrect.
// Go spec: adding keys during range iteration is allowed; range may or may not
// visit newly-added keys (non-deterministic). The result is at least 4.
func MapUpdateDuringRange() int {
	m := map[int]int{1: 10, 2: 20}
	for k := range m {
		m[k+10] = k
	}
	return len(m)
}

// ── Resolved Issue 10: Self-referencing struct type ──────────────────────────

type node struct {
	value int
	next  *node
}

// StructSelfRef tests self-referencing struct types.
// Previously caused stack overflow in typeToReflect due to infinite recursion.
// Expected: 3
func StructSelfRef() int {
	n1 := &node{value: 1}
	n2 := &node{value: 2, next: n1}
	return n2.value + n2.next.value
}

// ── Resolved Issue 11: Defer in closure with argument ────────────────────────

// DeferInClosureWithArg tests defer in closure with argument.
// Previously returned 1 because compileDefer pushed args before the closure,
// causing OpDeferIndirect to pop them in the wrong order.
// Expected: 11 (1 + 10 from defer)
func DeferInClosureWithArg() int {
	result := 0
	f := func() {
		defer func(v int) {
			result += v
		}(10)
		result = 1
	}
	f()
	return result
}

// ── Resolved Issue 12: Pointer swap in struct ────────────────────────────────

// PtrPair is a pair of int pointers for swap testing.
type PtrPair struct {
	a, b *int
}

// PointerSwapInStruct tests swapping pointer fields in struct.
// Previously returned 22 because OpDeref returned a reference to the struct
// field instead of an independent copy, causing the swap to alias.
// Expected: 21 (2*10 + 1)
func PointerSwapInStruct() int {
	x, y := 1, 2
	p := PtrPair{a: &x, b: &y}
	p.a, p.b = p.b, p.a
	return *p.a*10 + *p.b
}

// ── Resolved Issue 13: Struct with function slice ────────────────────────────

// StructWithFuncSlice tests struct with slice of functions.
// Previously panicked: []value.Value is not assignable to []func() int.
// Fixed by using typeToReflect for proper typed arrays/slices in OpNew.
// Expected: 3
func StructWithFuncSlice() int {
	type FuncSliceHolder struct {
		funcs []func() int
	}
	h := FuncSliceHolder{
		funcs: []func() int{
			func() int { return 1 },
			func() int { return 2 },
		},
	}
	return h.funcs[0]() + h.funcs[1]()
}

// ── Resolved Issue 14: Struct with anonymous field ───────────────────────────

// StructAnonymousField tests struct with anonymous embedded field.
// Previously panicked: reflect.StructOf field "int" is unexported but missing PkgPath.
// Fixed by demoting anonymous unexported fields in typeToReflect (reflect.StructOf limitation).
// Expected: 42
func StructAnonymousField() int {
	type AnonField struct {
		int
		name string
	}
	s := AnonField{int: 42, name: "test"}
	return s.int
}

// ── Resolved Issue 15: Struct with embedded interface ────────────────────────
// NOTE: StructEmbeddedInterface types and test function are defined inline
// in resolved_issue_test.go to avoid type collision with other struct types
// in this file (reflect.StructOf creates anonymous types that can conflict).

// ── Resolved Issue 16: Map range with break ──────────────────────────────────

// MapRangeWithBreak tests breaking from map range.
// Map iteration order is non-deterministic by Go spec.
// The result varies but should be a valid partial sum.
// Expected: varies (sum of some values from {10, 20, 30} until sum > 25)
func MapRangeWithBreak() int {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	sum := 0
	for _, v := range m {
		sum += v
		if sum > 25 {
			break
		}
	}
	return sum
}

// ── Resolved Issue 17: Pointer to interface type assertion ───────────────────

// PointerToInterface tests dereferencing a pointer to interface and type-asserting.
// Previously returned []value.Value (the raw [result, ok] tuple) instead of int
// because compileTypeAssert didn't handle non-comma-ok assertions.
// Expected: 42
func PointerToInterface() int {
	var i interface{} = 42
	p := &i
	return (*p).(int)
}

// ── Resolved Issue 18-20: see inline tests in resolved_issue_test.go ────────
// StructWithPointerToInterface, StructWithNestedFunc, StructWithInterfaceMap
// use package-level types that can cause reflect.StructOf type collision,
// so their tests are defined inline in resolved_issue_test.go.

// ── Resolved Issue 21: Pointer to slice element modify in loop ──────────────

// PointerToSliceElemModify tests modifying slice elements via pointer in a loop.
// Previously panicked with "invalid reflect.Value in SetElem()" because the
// fuseSliceOps peephole optimizer fused the read pattern (IndexAddr+Deref→IntSliceGet)
// but eliminated the SETLOCAL for the pointer temporary, leaving it uninitialized
// when the subsequent store (*p = expr) tried to use it.
// Fix: fuseSliceOps now checks localUsedOutside() to skip fusion when the pointer
// temporary is referenced outside the fused region.
// Expected: 120 (= 20 + 40 + 60)
func PointerToSliceElemModify() int {
	s := []int{10, 20, 30}
	for i := range s {
		p := &s[i]
		*p = *p * 2
	}
	return s[0] + s[1] + s[2]
}
