package main

import "fmt"

type receiver struct{ value int }

func (r receiver) ValueMethod() string {
	return fmt.Sprintf("value:%d", r.value)
}

func (r *receiver) PointerMethod() string {
	r.value++
	return fmt.Sprintf("pointer:%d", r.value)
}

type HasValue interface{ ValueMethod() string }
type HasPointer interface{ PointerMethod() string }

// Result checks value/pointer receiver dispatch through interfaces.
func Result() string {
	v := receiver{value: 10}
	p := &v

	s := []HasValue{v}
	pMethods := []HasPointer{p}

	return fmt.Sprintf("%s:%s:%s:%s", s[0].ValueMethod(), pMethods[0].PointerMethod(), v.ValueMethod(), p.PointerMethod())
}
