package divergence_hunt268

import (
	"fmt"
)

// ============================================================================
// Round 268: Interface composition and embedding
// ============================================================================

type Reader268 interface {
	Read() string
}

type Writer268 interface {
	Write(s string)
}

type ReadWriter268 interface {
	Reader268
	Writer268
}

type Document268 struct {
	content string
}

func (d *Document268) Read() string    { return d.content }
func (d *Document268) Write(s string) { d.content = s }

// InterfaceEmbedding tests interface embedding
func InterfaceEmbedding() string {
	var rw ReadWriter268 = &Document268{}
	rw.Write("hello")
	return fmt.Sprintf("read=%s", rw.Read())
}

// InterfaceAssignStruct tests assigning struct to interface
func InterfaceAssignStruct() string {
	var r Reader268 = &Document268{content: "test"}
	return fmt.Sprintf("read=%s", r.Read())
}

// InterfaceNilCheck tests nil interface check
func InterfaceNilCheck() string {
	var r Reader268
	return fmt.Sprintf("nil=%t", r == nil)
}

// InterfaceTypeAssertionFromEmbedded tests asserting from embedded interface
func InterfaceTypeAssertionFromEmbedded() string {
	var rw ReadWriter268 = &Document268{}
	d, ok := rw.(*Document268)
	return fmt.Sprintf("ok=%t,content=%s", ok, d.Read())
}

// EmptyInterface tests empty interface holding various types
func EmptyInterface() string {
	var a interface{} = 42
	var b interface{} = "hello"
	var c interface{} = []int{1, 2, 3}
	return fmt.Sprintf("a=%v,b=%v,c=%v", a, b, c)
}

// InterfaceSlice tests slice of interfaces
func InterfaceSlice() string {
	s := []interface{}{1, "two", 3.0}
	return fmt.Sprintf("len=%d,0=%v,1=%v", len(s), s[0], s[1])
}

// InterfaceMap tests map with interface values
func InterfaceMap() string {
	m := map[string]interface{}{
		"num":   42,
		"str":   "hello",
		"bool":  true,
	}
	return fmt.Sprintf("num=%v,str=%v,bool=%v", m["num"], m["str"], m["bool"])
}

// StructEmbeddingMethodPromotion tests method promotion through embedding
func StructEmbeddingMethodPromotion() string {
	type Base268 struct{ X int }
	type Derived268 struct {
		Base268
		Y int
	}
	d := Derived268{Base268: Base268{X: 10}, Y: 20}
	return fmt.Sprintf("x=%d,y=%d", d.X, d.Y)
}

// StructEmbeddingFieldAccess tests accessing promoted field
func StructEmbeddingFieldAccess() string {
	type Inner struct{ Val int }
	type Outer struct {
		Inner
		Name string
	}
	o := Outer{Inner: Inner{Val: 5}, Name: "test"}
	return fmt.Sprintf("val=%d,name=%s", o.Val, o.Name)
}

// InterfaceComparison tests interface value comparison
func InterfaceComparison() string {
	var a interface{} = 42
	var b interface{} = 42
	return fmt.Sprintf("eq=%t", a == b)
}
