package divergence_hunt209

import "fmt"

// ============================================================================
// Round 209: Function value comparisons
// ============================================================================

func simpleFunc209() {}

func adder209(a, b int) int {
	return a + b
}

type FuncHolder209 struct {
	Fn func()
}

type MultiFunc209 struct {
	F1 func() int
	F2 func() int
}

func FunctionNilComparison() string {
	var f func()
	return fmt.Sprintf("%v", f == nil)
}

func FunctionNonNilComparison() string {
	f := simpleFunc209
	return fmt.Sprintf("%v", f != nil)
}

func FunctionSameVar() string {
	f1 := simpleFunc209
	return fmt.Sprintf("%v", f1 != nil)
}

func FunctionDifferentVars() string {
	f1 := simpleFunc209
	f2 := simpleFunc209
	_ = f1
	_ = f2
	return "ok"
}

func FunctionInStruct() string {
	h := FuncHolder209{Fn: simpleFunc209}
	return fmt.Sprintf("%v", h.Fn != nil)
}

func FunctionInMap() string {
	m := map[string]func(){
		"a": simpleFunc209,
	}
	f, ok := m["a"]
	return fmt.Sprintf("ok:%v,nil:%v", ok, f == nil)
}

func FunctionSlice() string {
	funcs := []func() int{
		func() int { return 1 },
		func() int { return 2 },
	}
	sum := 0
	for _, f := range funcs {
		sum += f()
	}
	return fmt.Sprintf("%d", sum)
}

func FunctionAsInterface() string {
	var i interface{} = simpleFunc209
	_, ok := i.(func())
	return fmt.Sprintf("ok:%v", ok)
}

func FunctionVariableCapture() string {
	x := 10
	f := func() int { return x }
	x = 20
	return fmt.Sprintf("%d", f())
}

func FunctionAssignedToVar() string {
	f := simpleFunc209
	return fmt.Sprintf("%v", f != nil)
}

func ClosureComparison() string {
	f1 := func() {}
	_ = f1
	return "closures not comparable"
}
