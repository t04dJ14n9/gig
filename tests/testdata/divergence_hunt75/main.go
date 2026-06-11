package divergence_hunt75

// ============================================================================
// Round 75: Struct embedding edge cases - promotion, shadowing, method override
// ============================================================================

type Base1 struct {
	X int
}

func (b Base1) GetX() int {
	return b.X
}

type Derived1 struct {
	Base1
	Y int
}

func EmbeddingPromotion() int {
	d := Derived1{Base1: Base1{X: 10}, Y: 20}
	return d.GetX() + d.Y
}

func EmbeddingFieldAccess() int {
	d := Derived1{Base1: Base1{X: 10}, Y: 20}
	return d.X + d.Y
}

func EmbeddingExplicitBase() int {
	d := Derived1{Base1: Base1{X: 10}, Y: 20}
	return d.Base1.X + d.Y
}

type Inner struct {
	Value int
}

type Outer struct {
	Inner
	Label string
}

func NestedEmbedding() int {
	o := Outer{Inner: Inner{Value: 42}, Label: "test"}
	return o.Value
}

type A struct {
	Val int
}

type B struct {
	A
	Val string // shadows A.Val
}

func ShadowingEmbed() string {
	b := B{A: A{Val: 42}, Val: "hello"}
	return b.Val // should be "hello" (shadows)
}

func ShadowingExplicit() int {
	b := B{A: A{Val: 42}, Val: "hello"}
	return b.A.Val // should be 42 (explicit)
}

type Mover interface {
	Move() string
}

type Dog struct {
	Name string
}

func (d Dog) Move() string {
	return d.Name + " walks"
}

func EmbeddingInterface() string {
	var m Mover = Dog{Name: "Rex"}
	return m.Move()
}

type Point3D struct {
	X, Y, Z int
}

type ColoredPoint struct {
	Point3D
	Color string
}

func EmbeddingLiteral() int {
	cp := ColoredPoint{
		Point3D: Point3D{X: 1, Y: 2, Z: 3},
		Color:   "red",
	}
	return cp.X + cp.Y + cp.Z
}

func EmbeddingFieldAssign() int {
	cp := ColoredPoint{}
	cp.X = 10
	cp.Y = 20
	cp.Z = 30
	return cp.X + cp.Y + cp.Z
}

type DoubleEmbed1 struct {
	Val int
}

type DoubleEmbed2 struct {
	DoubleEmbed1
	Extra string
}

type DoubleEmbed3 struct {
	DoubleEmbed2
	Flag bool
}

func DoubleEmbedding() int {
	d := DoubleEmbed3{
		DoubleEmbed2: DoubleEmbed2{
			DoubleEmbed1: DoubleEmbed1{Val: 99},
			Extra:        "x",
		},
		Flag: true,
	}
	return d.Val
}

type WithMethod struct {
	Data int
}

func (w WithMethod) Process() int {
	return w.Data * 2
}

type Wrapper struct {
	WithMethod
	Name string
}

func EmbeddingMethodPromotion() int {
	w := Wrapper{WithMethod: WithMethod{Data: 21}, Name: "test"}
	return w.Process()
}
