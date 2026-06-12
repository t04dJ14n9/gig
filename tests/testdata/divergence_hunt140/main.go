package divergence_hunt140

import "fmt"

// ============================================================================
// Round 140: Embedded field access and method promotion
// ============================================================================

type Engine struct {
	Power int
}

func (e Engine) Start() string {
	return fmt.Sprintf("engine-%d", e.Power)
}

type Car struct {
	Engine
	Brand string
}

func EmbeddedFieldAccess() string {
	c := Car{Engine: Engine{Power: 200}, Brand: "Toyota"}
	return fmt.Sprintf("power=%d", c.Power)
}

func EmbeddedMethodCall() string {
	c := Car{Engine: Engine{Power: 150}, Brand: "Honda"}
	return c.Start()
}

func EmbeddedFieldExplicit() string {
	c := Car{Engine: Engine{Power: 100}, Brand: "Ford"}
	return fmt.Sprintf("power=%d", c.Engine.Power)
}

type Base struct {
	Name string
}

func (b Base) Greet() string {
	return fmt.Sprintf("hello %s", b.Name)
}

type Derived struct {
	Base
	Age int
}

func EmbeddedChain() string {
	d := Derived{Base: Base{Name: "Alice"}, Age: 30}
	return fmt.Sprintf("name=%s-age=%d", d.Name, d.Age)
}

func EmbeddedChainMethod() string {
	d := Derived{Base: Base{Name: "Bob"}, Age: 25}
	return d.Greet()
}

type Animal struct {
	Sound string
}

func (a Animal) Speak() string {
	return a.Sound
}

type Dog struct {
	Animal
	Name string
}

func EmbeddedOverride() string {
	d := Dog{Animal: Animal{Sound: "woof"}, Name: "Rex"}
	return fmt.Sprintf("%s-says-%s", d.Name, d.Speak())
}

type Wrapper struct {
	Val int
}

type Container struct {
	Wrapper
	Label string
}

func EmbeddedPointer() string {
	c := &Container{Wrapper: Wrapper{Val: 77}, Label: "test"}
	return fmt.Sprintf("val=%d", c.Val)
}

type Config struct {
	Debug bool
}

type App struct {
	Config
	Name string
}

func EmbeddedBoolField() string {
	a := App{Config: Config{Debug: true}, Name: "myapp"}
	if a.Debug {
		return fmt.Sprintf("%s-debug", a.Name)
	}
	return a.Name
}
