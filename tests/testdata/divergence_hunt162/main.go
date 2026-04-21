package divergence_hunt162

import "fmt"

// ============================================================================
// Round 162: Interface method calls with nil receivers
// ============================================================================

type Stringer interface {
	String() string
}

type Printer interface {
	Print() string
}

type ValueReceiver struct {
	Name string
}

func (v ValueReceiver) String() string {
	return v.Name
}

type PointerReceiver struct {
	Name string
}

func (p *PointerReceiver) String() string {
	if p == nil {
		return "<nil>"
	}
	return p.Name
}

type MixedReceiver struct {
	Value int
}

func (m MixedReceiver) ValueMethod() int {
	return m.Value
}

func (m *MixedReceiver) PointerMethod() string {
	if m == nil {
		return "pointer is nil"
	}
	return fmt.Sprintf("value=%d", m.Value)
}

// NilValueReceiver tests nil value receiver
func NilValueReceiver() string {
	var v *ValueReceiver
	result := v.String()
	return fmt.Sprintf("result=%s", result)
}

// NilPointerReceiverSafe tests nil pointer receiver that handles nil
func NilPointerReceiverSafe() string {
	var p *PointerReceiver
	result := p.String()
	return fmt.Sprintf("result=%s", result)
}

// NonNilValueReceiver tests non-nil value receiver
func NonNilValueReceiver() string {
	v := ValueReceiver{Name: "test"}
	return fmt.Sprintf("result=%s", v.String())
}

// NonNilPointerReceiver tests non-nil pointer receiver
func NonNilPointerReceiver() string {
	p := &PointerReceiver{Name: "test"}
	return fmt.Sprintf("result=%s", p.String())
}

// NilMixedReceiverValueMethod tests nil mixed receiver with value method
func NilMixedReceiverValueMethod() string {
	var m *MixedReceiver
	result := m.ValueMethod()
	return fmt.Sprintf("result=%d", result)
}

// NilMixedReceiverPointerMethod tests nil mixed receiver with pointer method
func NilMixedReceiverPointerMethod() string {
	var m *MixedReceiver
	result := m.PointerMethod()
	return fmt.Sprintf("result=%s", result)
}

// InterfaceWithNilConcrete tests interface holding nil concrete type
func InterfaceWithNilConcrete() string {
	var p *PointerReceiver
	var s Stringer = p
	result := ""
	if s == nil {
		result = "interface is nil"
	} else {
		result = s.String()
	}
	return fmt.Sprintf("result=%s", result)
}

// NilInterfaceCall tests calling method on nil interface
func NilInterfaceCall() string {
	var s Stringer
	defer func() {
		recover()
	}()
	_ = s.String()
	return "no panic"
}

// NestedNilReceiver tests nested struct with nil pointer
func NestedNilReceiver() string {
	type Inner struct {
		Name string
	}
	type Outer struct {
		*Inner
	}
	outer := Outer{}
	// outer.Inner is nil, but accessing outer.Name should work (returns zero value)
	result := outer.Name
	return fmt.Sprintf("result=%s", result)
}

// SliceOfInterfacesWithNil tests slice containing nil interfaces
func SliceOfInterfacesWithNil() string {
	items := []Stringer{
		&PointerReceiver{Name: "first"},
		nil,
		&PointerReceiver{Name: "third"},
	}
	result := ""
	for i, item := range items {
		if item == nil {
			result += fmt.Sprintf("[%d]=nil ", i)
		} else {
			result += fmt.Sprintf("[%d]=%s ", i, item.String())
		}
	}
	return result
}

// MapWithNilValues tests map with nil interface values
func MapWithNilValues() string {
	items := map[string]Stringer{
		"a": &PointerReceiver{Name: "alpha"},
		"b": nil,
		"c": &PointerReceiver{Name: "charlie"},
	}
	result := ""
	for k, v := range items {
		if v == nil {
			result += fmt.Sprintf("%s=nil ", k)
		} else {
			result += fmt.Sprintf("%s=%s ", k, v.String())
		}
	}
	return result
}
