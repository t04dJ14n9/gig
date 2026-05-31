package value

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type externTestError struct {
	msg string
}

func (e *externTestError) Error() string { return e.msg }

func TestFmtWrapDetectsGigStructAndFormatsTypeName(t *testing.T) {
	rv := newExternTestGigStruct(t, "example.Widget")
	rv.FieldByName("Name").SetString("alpha")

	wrapped := FmtWrap(MakeFromReflect(rv))
	wrapper, ok := wrapped.(*gigStructWrapper)
	if !ok {
		t.Fatalf("FmtWrap returned %T, want *gigStructWrapper", wrapped)
	}
	if wrapper.typeName != "example.Widget" {
		t.Fatalf("wrapper typeName = %q, want example.Widget", wrapper.typeName)
	}
	if got := fmt.Sprintf("%v", wrapped); got != "{alpha}" {
		t.Fatalf("formatted value = %q, want {alpha}", got)
	}
	if got := fmt.Sprintf("%#v", wrapped); got != `example.Widget{Name:"alpha"}` {
		t.Fatalf("go-syntax value = %q, want example.Widget{Name:\"alpha\"}", got)
	}
}

func TestIsGigStructDetectsPointersAndValues(t *testing.T) {
	rv := newExternTestGigStruct(t, "example.Widget")

	if got := isGigStruct(rv.Interface()); got != "example.Widget" {
		t.Fatalf("isGigStruct(value) = %q, want example.Widget", got)
	}
	if got := isGigStruct(rv.Addr().Interface()); got != "example.Widget" {
		t.Fatalf("isGigStruct(pointer) = %q, want example.Widget", got)
	}
}

func TestGigErrorsNativeCompatibility(t *testing.T) {
	target := &externTestError{msg: "target"}
	wrapped := fmt.Errorf("wrapped: %w", target)

	if !GigErrorsIs(FromInterface(wrapped), FromInterface(target)) {
		t.Fatal("GigErrorsIs did not match native wrapped error")
	}

	var asTarget *externTestError
	if !GigErrorsAs(wrapped, &asTarget) {
		t.Fatal("GigErrorsAs did not match native wrapped error")
	}
	if asTarget != target {
		t.Fatalf("GigErrorsAs target = %#v, want original target", asTarget)
	}

	unwrapped := GigErrorsUnwrap(FromInterface(wrapped))
	if got, ok := unwrapped.Interface().(error); !ok || !errors.Is(got, target) {
		t.Fatalf("GigErrorsUnwrap = %#v, want error wrapping target", unwrapped.Interface())
	}
}

func newExternTestGigStruct(t *testing.T, typeName string) reflect.Value {
	t.Helper()

	rt := reflect.StructOf([]reflect.StructField{
		{
			Name:    "gigType",
			PkgPath: "gig/internal",
			Type:    reflect.TypeOf(struct{}{}),
			Tag:     reflect.StructTag(`gig:"` + typeName + `"`),
		},
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
		},
	})
	return reflect.New(rt).Elem()
}
