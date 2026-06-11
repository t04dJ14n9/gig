package divergence_hunt116

import "fmt"

// ============================================================================
// Round 116: Interface nil semantics deep test
// ============================================================================

type Niler interface{ NilMethod() }

type NilerImpl struct{ Val int }

func (n *NilerImpl) NilMethod() { _ = n.Val }

func InterfaceNilCompare() string {
	var i interface{}
	if i == nil {
		return "nil"
	}
	return "not nil"
}

func TypedNilInterface() string {
	var p *NilerImpl
	var i Niler = p
	if i == nil {
		return "nil"
	}
	return "not nil"
}

func NilInterfaceTypeAssert() string {
	var i interface{}
	_, ok := i.(string)
	return fmt.Sprintf("%v", ok)
}

func NilInterfaceTypeSwitch() string {
	var i interface{}
	switch i.(type) {
	case nil:
		return "nil"
	default:
		return "non-nil"
	}
}

func EmptyInterfaceVsNil() string {
	var i interface{}
	return fmt.Sprintf("%v:%v", i == nil, i != nil)
}

func NilSliceVsNilInterface() string {
	var s []int
	var i interface{} = s
	return fmt.Sprintf("%v:%v", s == nil, i == nil)
}

func NilMapVsNilInterface() string {
	var m map[string]int
	var i interface{} = m
	return fmt.Sprintf("%v:%v", m == nil, i == nil)
}

func NilFuncVsNilInterface() string {
	var f func()
	var i interface{} = f
	return fmt.Sprintf("%v:%v", f == nil, i == nil)
}

func NilChanVsNilInterface() string {
	var ch chan int
	var i interface{} = ch
	return fmt.Sprintf("%v:%v", ch == nil, i == nil)
}

func InterfaceReturnNil() string {
	getNil := func() interface{} { return nil }
	v := getNil()
	return fmt.Sprintf("%v", v == nil)
}
