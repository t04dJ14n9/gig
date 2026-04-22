package divergence_hunt286

import (
	"fmt"
)

// ============================================================================
// Round 286: Interface method sets — value vs pointer receiver, method set rules

type Describer interface {
	Describe() string
}

type FullDescriber interface {
	Describe() string
	SetName(string)
}

type Animal struct {
	Name string
}

func (a Animal) Describe() string {
	return "animal:" + a.Name
}

func (a *Animal) SetName(name string) {
	a.Name = name
}

// ValueReceiverOnInterface tests value type satisfies interface with value method
func ValueReceiverOnInterface() string {
	a := Animal{Name: "cat"}
	var d Describer = a
	return d.Describe()
}

// PointerReceiverOnInterface tests pointer type satisfies interface
func PointerReceiverOnInterface() string {
	a := Animal{Name: "dog"}
	var d Describer = &a
	return d.Describe()
}

// PointerSatisfiesFullInterface tests *T satisfies interface with both value+pointer methods
func PointerSatisfiesFullInterface() string {
	a := Animal{Name: "fish"}
	var d FullDescriber = &a
	d.SetName("whale")
	return d.Describe()
}

// InterfaceSlice tests slice of interfaces
func InterfaceSlice() string {
	items := []Describer{
		Animal{Name: "cat"},
		&Animal{Name: "dog"},
	}
	return fmt.Sprintf("%s,%s", items[0].Describe(), items[1].Describe())
}

// InterfaceNil tests nil interface value
func InterfaceNil() string {
	var d Describer
	return fmt.Sprintf("nil=%t", d == nil)
}

// TypeAssertionOnConcrete tests type assertion to concrete type
func TypeAssertionOnConcrete() string {
	var d Describer = Animal{Name: "bird"}
	a, ok := d.(Animal)
	return fmt.Sprintf("name=%s,ok=%t", a.Name, ok)
}

// TypeAssertionOnPointer tests type assertion to pointer type
func TypeAssertionOnPointer() string {
	var d Describer = &Animal{Name: "bear"}
	a, ok := d.(*Animal)
	return fmt.Sprintf("name=%s,ok=%t", a.Name, ok)
}

// TypeAssertionWrongType tests type assertion to wrong type
func TypeAssertionWrongType() string {
	var d Describer = Animal{Name: "cat"}
	_, ok := d.(*Animal) // value doesn't satisfy pointer interface
	return fmt.Sprintf("ok=%t", ok)
}

// EmptyInterfaceHoldsAnything tests interface{} holds any type
func EmptyInterfaceHoldsAnything() string {
	var i interface{} = 42
	_, ok1 := i.(int)
	i = "hello"
	_, ok2 := i.(string)
	return fmt.Sprintf("int=%t,string=%t", ok1, ok2)
}

// InterfaceMethodOnNil tests calling method on nil typed value in interface
// In Go, calling any method (even value receiver) on nil *T through an
// interface panics because the runtime dereferences the nil pointer.
func InterfaceMethodOnNil() string {
	var a *Animal = nil
	var d Describer = a
	result := "no_panic"
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = "panic"
			}
		}()
		d.Describe()
	}()
	return result
}
