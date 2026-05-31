package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

type memoryHelperBase struct {
	msg string
}

type memoryHelperSpecific struct {
	memoryHelperBase
	code int
}

func TestFieldAddressValueHandlesUnexportedEmbeddedField(t *testing.T) {
	input := &memoryHelperSpecific{
		memoryHelperBase: memoryHelperBase{msg: "specific"},
		code:             123,
	}

	got := fieldAddressValue(value.MakeFromReflect(reflect.ValueOf(input)), 0)
	rv, ok := got.ReflectValue()
	if !ok {
		t.Fatalf("fieldAddressValue returned non-reflect value: %v", got)
	}
	if rv.Kind() != reflect.Ptr {
		t.Fatalf("fieldAddressValue kind = %v, want pointer", rv.Kind())
	}
	msg := rv.Elem().FieldByName("msg")
	if msg.Kind() != reflect.String || msg.String() != "specific" {
		t.Fatalf("fieldAddressValue msg = %v, want specific", msg)
	}
}

func TestReferenceHelpersDoNotInterfaceUnexportedFields(t *testing.T) {
	input := memoryHelperSpecific{
		memoryHelperBase: memoryHelperBase{msg: "specific"},
		code:             123,
	}
	field := reflect.ValueOf(input).Field(0)
	if field.CanInterface() {
		t.Fatal("test field unexpectedly allows Interface")
	}
	wrapped := value.MakeFromReflect(field)

	if _, ok := unwrapValueSlot(wrapped); ok {
		t.Fatal("unwrapValueSlot matched an unexported field")
	}
	if _, ok := globalRefFromValue(wrapped); ok {
		t.Fatal("globalRefFromValue matched an unexported field")
	}
}

func TestReferenceHelpersUnwrapGlobalSlotAndRef(t *testing.T) {
	slot := value.MakeInt(41)
	wrappedSlot := value.FromInterface(&slot)

	got, ok := unwrapValueSlot(wrappedSlot)
	if !ok {
		t.Fatal("unwrapValueSlot did not match *value.Value")
	}
	if got.Int() != 41 {
		t.Fatalf("unwrapValueSlot = %d, want 41", got.Int())
	}

	sg := NewSharedGlobals([]value.Value{value.MakeInt(1)}, 1)
	ref := &GlobalRef{sg: sg, idx: 0}
	gotRef, ok := globalRefFromValue(value.FromInterface(ref))
	if !ok {
		t.Fatal("globalRefFromValue did not match *GlobalRef")
	}
	gotRef.Store(value.MakeInt(42))
	if got := sg.Get(0).Int(); got != 42 {
		t.Fatalf("GlobalRef Store result = %d, want 42", got)
	}
}

func TestIndexAddressValueKeepsNativeIntSliceFastPath(t *testing.T) {
	ints := []int64{10, 20, 30}

	got := indexAddressValue(value.MakeIntSlice(ints), 1)
	ptr, ok := got.IntPtr()
	if !ok {
		t.Fatalf("indexAddressValue returned %v, want native int pointer", got)
	}
	*ptr = 99
	if ints[1] != 99 {
		t.Fatalf("slice after pointer write = %v, want index 1 updated", ints)
	}
}

func TestIndexAddressValueHandlesByteSlices(t *testing.T) {
	bytes := []byte("abc")

	got := indexAddressValue(value.MakeBytes(bytes), 1)
	rv, ok := got.ReflectValue()
	if !ok {
		t.Fatalf("indexAddressValue returned %v, want reflect pointer", got)
	}
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Uint8 {
		t.Fatalf("indexAddressValue reflect kind = %v/%v, want *uint8", rv.Kind(), rv.Elem().Kind())
	}
	rv.Elem().SetUint('z')
	if string(bytes) != "azc" {
		t.Fatalf("bytes after pointer write = %q, want azc", string(bytes))
	}
}

func TestDereferenceValueLoadsGlobalRefAndValueSlot(t *testing.T) {
	sg := NewSharedGlobals([]value.Value{value.MakeInt(7)}, 1)
	ref := &GlobalRef{sg: sg, idx: 0}

	if got := dereferenceValue(value.FromInterface(ref)); got.Int() != 7 {
		t.Fatalf("dereferenceValue(GlobalRef) = %d, want 7", got.Int())
	}

	slot := value.MakeString("slot")
	if got := dereferenceValue(value.FromInterface(&slot)); got.String() != "slot" {
		t.Fatalf("dereferenceValue(*Value) = %q, want slot", got.String())
	}
}

func TestSetDereferenceValueStoresGlobalRefAndInterfaceSlot(t *testing.T) {
	prog := &bytecode.CompiledProgram{}
	v := New(prog).(*vm)

	sg := NewSharedGlobals([]value.Value{value.MakeInt(1)}, 1)
	ref := &GlobalRef{sg: sg, idx: 0}
	v.setDereferenceValue(value.FromInterface(ref), value.MakeInt(9))
	if got := sg.Get(0).Int(); got != 9 {
		t.Fatalf("setDereferenceValue(GlobalRef) = %d, want 9", got)
	}

	var iface any
	ptr := value.MakeFromReflect(reflect.ValueOf(&iface))
	v.setDereferenceValue(ptr, value.MakeString("iface"))
	if iface != "iface" {
		t.Fatalf("interface after setDereferenceValue = %#v, want iface", iface)
	}
}
