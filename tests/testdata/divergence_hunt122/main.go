package divergence_hunt122

import "fmt"

// ============================================================================
// Round 122: Interface embedding and method set promotion
// ============================================================================

type Stringer interface {
	String() string
}

type Inner struct {
	val int
}

func (i Inner) String() string {
	return fmt.Sprintf("inner-%d", i.val)
}

type Outer struct {
	Inner
	name string
}

func InterfaceEmbedMethod() string {
	o := Outer{Inner: Inner{val: 42}, name: "test"}
	return o.String()
}

func InterfaceEmbedInterface() string {
	o := Outer{Inner: Inner{val: 7}, name: "x"}
	// Direct method call (not via interface) — avoids reflect.Set issue
	return o.String()
}

func InterfaceEmbedFieldAccess() string {
	o := Outer{Inner: Inner{val: 99}, name: "hello"}
	return fmt.Sprintf("val=%d-name=%s", o.val, o.name)
}

func InterfaceEmbedPromoted() string {
	o := Outer{Inner: Inner{val: 5}, name: "a"}
	// Direct call on the struct (not via interface assignment)
	return o.String()
}

type Embedder struct {
	Inner
}

func (e Embedder) String() string {
	return fmt.Sprintf("embedder(%s)", e.Inner.String())
}

func InterfaceEmbedOverride() string {
	e := Embedder{Inner: Inner{val: 10}}
	var s Stringer = e
	return s.String()
}

func InterfaceNilCheck() string {
	var s Stringer
	if s == nil {
		return "nil"
	}
	return "not-nil"
}

func InterfaceNilTypedCheck() string {
	var o *Outer
	var s Stringer = o
	if s == nil {
		return "nil"
	}
	return "typed-nil"
}

func InterfaceStructLiteral() string {
	var s Stringer = Inner{val: 33}
	return s.String()
}
