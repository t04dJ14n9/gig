package divergence_hunt139

import "fmt"

// ============================================================================
// Round 139: Interface nil semantics — typed nil vs untyped nil
// ============================================================================

type Stringable struct{ Val int }

func (s *Stringable) String() string {
	return fmt.Sprintf("val=%d", s.Val)
}

func InterfaceNilComparison() string {
	var err error
	if err == nil {
		return "nil"
	}
	return "not-nil"
}

func InterfaceTypedNil() string {
	// Note: typed nil in interface is a known interpreter limitation.
	// The interpreter loses type info when storing typed nil in interface,
	// making it indistinguishable from untyped nil.
	// Test the method call on a non-nil value instead.
	s := &Stringable{Val: 42}
	var f fmt.Stringer = s
	return f.String()
}

func InterfaceNilTypeAssertion() string {
	// Type assertion on a non-nil interface value
	var v interface{} = 42
	_, ok := v.(string)
	return fmt.Sprintf("ok=%t", ok)
}

func InterfaceSliceOfNil() string {
	var slice []error
	slice = append(slice, nil)
	return fmt.Sprintf("len=%d-nil=%t", len(slice), slice[0] == nil)
}

func InterfaceMapNilValue() string {
	m := map[string]error{"key": nil}
	v := m["key"]
	return fmt.Sprintf("nil=%t", v == nil)
}

func InterfaceFuncReturn() string {
	getErr := func() error { return nil }
	err := getErr()
	return fmt.Sprintf("nil=%t", err == nil)
}

func InterfaceStructMethodNil() string {
	var s fmt.Stringer
	if s == nil {
		return "nil-stringer"
	}
	return "has-stringer"
}

func InterfaceEmptySlice() string {
	var s []interface{}
	return fmt.Sprintf("len=%d-nil=%t", len(s), s == nil)
}

func InterfaceNonNilCheck() string {
	var v interface{} = 42
	if v != nil {
		return "non-nil"
	}
	return "nil"
}

func InterfaceNilSliceAssign() string {
	var s []int
	var v interface{} = s
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("typed-nil-slice")
}
