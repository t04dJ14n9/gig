package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

type findMethodValueReceiver struct{}

func (findMethodValueReceiver) ValueMethod() string { return "value" }

type findMethodPointerReceiver struct {
	N int
}

func (r *findMethodPointerReceiver) PointerMethod() int { return r.N }

type findMethodEmbedded interface {
	EmbeddedMethod() string
}

type findMethodEmbeddedConcrete struct{}

func (findMethodEmbeddedConcrete) EmbeddedMethod() string { return "embedded" }

type findMethodHolder struct {
	Inner findMethodEmbedded
}

func TestFindMethodResolvesDirectValueMethod(t *testing.T) {
	method, ok := findMethod(reflect.ValueOf(findMethodValueReceiver{}), "ValueMethod", nil)
	if !ok {
		t.Fatal("findMethod did not find direct value method")
	}

	out := method.Call(nil)
	if len(out) != 1 || out[0].String() != "value" {
		t.Fatalf("direct value method returned %#v, want value", out)
	}
}

func TestFindMethodResolvesPointerMethodFromNonAddressableStructCopy(t *testing.T) {
	method, ok := findMethod(reflect.ValueOf(findMethodPointerReceiver{N: 7}), "PointerMethod", nil)
	if !ok {
		t.Fatal("findMethod did not find pointer method through addressable copy")
	}

	out := method.Call(nil)
	if len(out) != 1 || out[0].Int() != 7 {
		t.Fatalf("pointer method returned %#v, want 7", out)
	}
}

func TestFindMethodResolvesEmbeddedInterfaceMethodAndRewritesReceiver(t *testing.T) {
	args := []value.Value{value.MakeString("placeholder")}
	method, ok := findMethod(
		reflect.ValueOf(findMethodHolder{Inner: findMethodEmbeddedConcrete{}}),
		"EmbeddedMethod",
		args,
	)
	if !ok {
		t.Fatal("findMethod did not find embedded interface method")
	}

	out := method.Call(nil)
	if len(out) != 1 || out[0].String() != "embedded" {
		t.Fatalf("embedded method returned %#v, want embedded", out)
	}

	receiver, ok := args[0].ReflectValue()
	if !ok {
		t.Fatal("embedded interface receiver was not rewritten to reflect value")
	}
	if receiver.Type() != reflect.TypeOf(findMethodEmbeddedConcrete{}) {
		t.Fatalf("receiver type = %v, want findMethodEmbeddedConcrete", receiver.Type())
	}
}
