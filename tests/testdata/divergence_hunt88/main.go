package divergence_hunt88

// ============================================================================
// Round 88: Method value edge cases - bound methods, method expressions
// ============================================================================

type Adder struct{ X int }

func (a Adder) Add(y int) int {
	return a.X + y
}

func (a *Adder) AddPtr(y int) int {
	a.X += y
	return a.X
}

func MethodValue() int {
	a := Adder{X: 10}
	f := a.Add // bound method value
	return f(5)
}

func MethodValuePointer() int {
	a := &Adder{X: 10}
	f := a.AddPtr // bound method value on pointer
	return f(5)
}

func MethodCall() int {
	a := Adder{X: 10}
	return a.Add(5)
}

func MethodCallPointer() int {
	a := &Adder{X: 10}
	return a.Add(5) // value method on pointer
}

type Greeter struct {
	Greeting string
}

func (g Greeter) Greet(name string) string {
	return g.Greeting + " " + name
}

func MethodValueString() string {
	g := Greeter{Greeting: "Hello"}
	f := g.Greet
	return f("World")
}

type Accumulator struct {
	total int
}

func (a *Accumulator) Add(n int) {
	a.total += n
}

func (a *Accumulator) Total() int {
	return a.total
}

func MethodValueModify() int {
	a := &Accumulator{}
	f := a.Add
	f(10)
	f(20)
	return a.Total()
}

func MethodOnStructLiteral() int {
	return Adder{X: 100}.Add(50)
}

type Scale struct {
	Factor int
}

func (s Scale) Apply(x int) int {
	return x * s.Factor
}

func MethodValueInLoop() int {
	fns := []func(int) int{}
	for i := 0; i < 3; i++ {
		s := Scale{Factor: i + 1}
		fns = append(fns, s.Apply)
	}
	return fns[0](10) + fns[1](10) + fns[2](10)
}

type Stringer struct {
	val string
}

func (s Stringer) String() string {
	return s.val
}

func MethodValueReturn() string {
	getStringer := func() Stringer {
		return Stringer{val: "hello"}
	}
	return getStringer().String()
}

type Transformer struct {
	prefix string
}

func (t Transformer) Transform(s string) string {
	return t.prefix + s
}

func MethodValueChain() int {
	s := "abc"
	t := Transformer{prefix: "x"}
	result := t.Transform(s)
	return len(result)
}
