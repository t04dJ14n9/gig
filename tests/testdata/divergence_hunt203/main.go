package divergence_hunt203

import "fmt"

// ============================================================================
// Round 203: Method promotion and overriding
// ============================================================================

type Animal203 struct {
	Name string
}

func (a Animal203) Speak() string {
	return fmt.Sprintf("%s makes a sound", a.Name)
}

func (a Animal203) GetName() string {
	return a.Name
}

type Dog203 struct {
	Animal203
	Breed string
}

func (d Dog203) Speak() string {
	return fmt.Sprintf("%s barks", d.Name)
}

func BasicPromotion() string {
	d := Dog203{Animal203{Name: "Rex"}, "German Shepherd"}
	return d.GetName()
}

func MethodOverride() string {
	d := Dog203{Animal203{Name: "Buddy"}, "Labrador"}
	return d.Speak()
}

func EmbeddedMethodCall() string {
	d := Dog203{Animal203{Name: "Max"}, "Poodle"}
	return d.Animal203.Speak()
}

type Cat203 struct {
	Animal203
}

func (c Cat203) Speak() string {
	return fmt.Sprintf("%s meows", c.Name)
}

type Pet203 interface {
	Speak() string
}

func InterfaceWithPromoted() string {
	var p Pet203 = Dog203{Animal203{Name: "Fido"}, "Beagle"}
	return p.Speak()
}

type Wrapper203 struct {
	Dog203
}

func DeepEmbedding() string {
	w := Wrapper203{Dog203{Animal203{Name: "Deep"}, "Mixed"}}
	return w.Speak()
}

type Thing203 struct {
	Value int
}

func (t Thing203) ValueMethod() int {
	return t.Value
}

type Other203 struct {
	Thing203
}

func PromotedValueMethod() string {
	o := Other203{Thing203{42}}
	return fmt.Sprintf("%d", o.ValueMethod())
}

func EmbeddedPointerMethod() string {
	d := &Dog203{Animal203{Name: "Ptr"}, "Husky"}
	return d.Speak()
}
