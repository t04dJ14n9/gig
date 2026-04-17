package divergence_hunt93

import "fmt"

// ============================================================================
// Round 93: Interface satisfaction and dynamic dispatch
// ============================================================================

type Speaker interface {
	Speak() string
}

type Dog struct{ Name string }

func (d Dog) Speak() string { return d.Name + " says woof" }

type Cat struct{ Name string }

func (c Cat) Speak() string { return c.Name + " says meow" }

type Robot struct{ ID int }

func (r Robot) Speak() string { return fmt.Sprintf("robot %d beep", r.ID) }

func InterfaceSlice() string {
	speakers := []Speaker{Dog{"Rex"}, Cat{"Whiskers"}, Robot{7}}
	result := ""
	for _, s := range speakers {
		result += s.Speak() + ";"
	}
	return result
}

func InterfaceMap() string {
	m := map[string]Speaker{
		"dog":  Dog{"Buddy"},
		"cat":  Cat{"Mittens"},
		"bot":  Robot{42},
	}
	return m["dog"].Speak()
}

func InterfaceParam() string {
	describe := func(s Speaker) string {
		return "speaker: " + s.Speak()
	}
	return describe(Dog{"Fido"})
}

func InterfaceReturn() string {
	getSpeaker := func(kind string) Speaker {
		if kind == "cat" {
			return Cat{"Luna"}
		}
		return Dog{"Max"}
	}
	return getSpeaker("cat").Speak()
}

func InterfaceNil() string {
	var s Speaker
	if s == nil {
		return "nil"
	}
	return "not nil"
}

func InterfaceTypedNil() string {
	var d *Dog
	var s Speaker = d
	if s == nil {
		return "nil"
	}
	return "not nil"
}

func InterfaceSliceOfInterface() int {
	items := []interface{}{42, "hello", 3.14, true}
	return len(items)
}

func InterfaceSliceTypeAssert() string {
	items := []interface{}{42, "hello", true}
	result := ""
	for _, item := range items {
		if s, ok := item.(string); ok {
			result += s
		}
	}
	return result
}

func DoubleInterfaceEmbedding() string {
	type Speaker2 interface {
		Speak() string
	}
	typeDescriber := func(v Speaker2) string {
		return v.Speak()
	}
	return typeDescriber(Dog{"Rex"})
}
