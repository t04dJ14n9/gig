package divergence_hunt275

import (
	"fmt"
)

// ============================================================================
// Round 275: Typed nil vs untyped nil interface

// TypedNilInterface tests that a typed nil is NOT nil interface
func TypedNilInterface() string {
	var s *[]int = nil
	var i interface{} = s
	return fmt.Sprintf("nil=%t", i == nil)
}

// UntypedNilInterface tests that untyped nil IS nil interface
func UntypedNilInterface() string {
	var i interface{} = nil
	return fmt.Sprintf("nil=%t", i == nil)
}

// TypedNilMap tests nil map in interface
func TypedNilMap() string {
	var m map[string]int = nil
	var i interface{} = m
	return fmt.Sprintf("nil=%t", i == nil)
}

// TypedNilSlice tests nil slice in interface
func TypedNilSlice() string {
	var s []int = nil
	var i interface{} = s
	return fmt.Sprintf("nil=%t", i == nil)
}

// TypedNilFunc tests nil func in interface
func TypedNilFunc() string {
	var f func() = nil
	var i interface{} = f
	return fmt.Sprintf("nil=%t", i == nil)
}

// TypedNilChan tests nil chan in interface
func TypedNilChan() string {
	var ch chan int = nil
	var i interface{} = ch
	return fmt.Sprintf("nil=%t", i == nil)
}

// InterfaceWithConcreteValue tests interface with concrete value
func InterfaceWithConcreteValue() string {
	x := 42
	var i interface{} = &x
	return fmt.Sprintf("nil=%t", i == nil)
}

// NilCheckWithReflection tests checking nil via type switch
func NilCheckWithReflection() string {
	var s *[]int = nil
	var i interface{} = s
	result := "not_nil"
	switch v := i.(type) {
	case *[]int:
		if v == nil {
			result = "typed_nil"
		}
	case nil:
		result = "untyped_nil"
	}
	return result
}

// ReassignNilInterface tests reassigning nil to interface
func ReassignNilInterface() string {
	x := 42
	var i interface{} = x
	i = nil
	return fmt.Sprintf("nil=%t", i == nil)
}

// InterfaceSliceOfNil tests slice of nil interfaces
func InterfaceSliceOfNil() string {
	s := make([]interface{}, 3)
	return fmt.Sprintf("len=%d,nil0=%t,nil1=%t,nil2=%t", len(s), s[0] == nil, s[1] == nil, s[2] == nil)
}

// InterfaceMapNilValue tests map with nil interface values
func InterfaceMapNilValue() string {
	m := map[string]interface{}{"a": nil, "b": 42}
	return fmt.Sprintf("a_nil=%t,b_nil=%t,b=%d", m["a"] == nil, m["b"] == nil, m["b"])
}

// NilInterfaceMethodCallPanics tests calling method on nil interface panics
func NilInterfaceMethodCallPanics() (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = "panic"
		}
	}()
	var i interface{} = nil
	_ = i.(int) // will panic
	return "no_panic"
}
