package divergence_hunt205

import "fmt"

// ============================================================================
// Round 205: Nil interface vs nil value in interface
// ============================================================================

type Printer205 interface {
	Print() string
}

type MyPrinter205 struct {
	Name string
}

func (m *MyPrinter205) Print() string {
	return m.Name
}

func NilInterface() string {
	var p Printer205
	return fmt.Sprintf("%v", p == nil)
}

func NilPointerInInterface() string {
	var m *MyPrinter205
	var p Printer205 = m
	return fmt.Sprintf("p==nil:%v,m==nil:%v", p == nil, m == nil)
}

func TypedNilInterface() string {
	var p Printer205
	if p == nil {
		return "nil"
	}
	return "not nil"
}

func NonNilInterface() string {
	m := &MyPrinter205{Name: "test"}
	var p Printer205 = m
	return fmt.Sprintf("%v", p != nil)
}

type Error205 interface {
	Error() string
}

type MyError205 struct {
	Msg string
}

func (e *MyError205) Error() string {
	return e.Msg
}

func NilErrorInterface() string {
	var err Error205
	return fmt.Sprintf("%v", err == nil)
}

func NilPointerAsError() string {
	var e *MyError205
	var err Error205 = e
	return fmt.Sprintf("err==nil:%v", err == nil)
}

func InterfaceTypeWithNil() string {
	var p Printer205
	typeOf := fmt.Sprintf("%T", p)
	return typeOf
}

func CompareNilInterfaces() string {
	var a, b Printer205
	return fmt.Sprintf("%v", a == b)
}

func ReturnNilInterface() Printer205 {
	return nil
}

func ReturnNilPointerInterface() Printer205 {
	var m *MyPrinter205
	return m
}

func NilReturnComparison() string {
	a := ReturnNilInterface()
	b := ReturnNilPointerInterface()
	return fmt.Sprintf("a==nil:%v,b==nil:%v", a == nil, b == nil)
}
