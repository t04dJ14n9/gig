package main

import "fmt"

type named interface {
	Name() string
}

type pointerNamed interface {
	PointerName() string
}

type node struct {
	name string
}

func (n node) Name() string {
	return "node:" + n.name
}

func (n *node) PointerName() string {
	if n == nil {
		return "nil-pointer"
	}
	return "ptr:" + n.name
}

// Result checks nil pointer method dispatch through interfaces.
func Result() string {
	var p *node
	var ptrIface pointerNamed = p
	ptrResult := ptrIface.PointerName()

	panicResult := "no-panic"
	func() {
		defer func() {
			if recover() != nil {
				panicResult = "panic"
			}
		}()
		var valueIface named = p
		_ = valueIface.Name()
	}()

	var valueIface named = node{name: "value"}
	var ptrIface2 pointerNamed = &node{name: "ptr"}
	return fmt.Sprintf("%s:%s:%s:%s", ptrResult, panicResult, valueIface.Name(), ptrIface2.PointerName())
}
