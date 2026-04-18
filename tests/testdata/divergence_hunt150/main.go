package divergence_hunt150

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Round 150: Comprehensive integration stress test — final round
// ============================================================================

type Result struct {
	Val  int
	Ok   bool
	Err  error
}

func NewResult(val int, ok bool) *Result {
	return &Result{Val: val, Ok: ok}
}

func (r *Result) String() string {
	if r.Ok {
		return fmt.Sprintf("ok(%d)", r.Val)
	}
	return fmt.Sprintf("err(%d)", r.Val)
}

func (r *Result) SetErr(msg string) {
	r.Err = errors.New(msg)
	r.Ok = false
}

func IntegrationStructMethod() string {
	r := NewResult(42, true)
	return r.String()
}

func IntegrationStructMutation() string {
	r := NewResult(10, true)
	r.SetErr("failed")
	return r.String()
}

func IntegrationSliceMapFilter() string {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	m := make(map[int]bool)
	for _, n := range nums {
		if n%2 == 0 {
			m[n] = true
		}
	}
	return fmt.Sprintf("even-count=%d", len(m))
}

func IntegrationErrorChain() string {
	base := errors.New("root cause")
	mid := fmt.Errorf("middleware: %w", base)
	top := fmt.Errorf("handler: %w", mid)
	if errors.Is(top, base) {
		return "root-found"
	}
	return "root-missing"
}

func IntegrationStringProcess() string {
	s := "  Hello, World!  "
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "!", "")
	return s
}

func IntegrationClosureCounter() string {
	counter := func() func() int {
		n := 0
		return func() int {
			n++
			return n
		}
	}()
	counter()
	counter()
	return fmt.Sprintf("count=%d", counter())
}

func IntegrationDeferRecover() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("recovered: %v", r)
		}
	}()
	panic("integration-panic")
}

func IntegrationPointerChain() string {
	type Inner struct{ Val int }
	type Outer struct{ Ptr *Inner }
	o := &Outer{Ptr: &Inner{Val: 99}}
	o.Ptr.Val = 100
	return fmt.Sprintf("val=%d", o.Ptr.Val)
}

func IntegrationTypeSwitch() string {
	check := func(v interface{}) string {
		switch v := v.(type) {
		case int:
			return fmt.Sprintf("int:%d", v)
		case string:
			return fmt.Sprintf("str:%s", v)
		case bool:
			return fmt.Sprintf("bool:%t", v)
		default:
			return fmt.Sprintf("other:%T", v)
		}
	}
	return check(42) + "-" + check("hi") + "-" + check(true)
}

func IntegrationNamedReturn() (result string) {
	defer func() {
		result = strings.ToUpper(result)
	}()
	result = "hello world"
	return
}
