package vm

import (
	"context"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestAppendValueAppendsIntToByteSlice(t *testing.T) {
	got := appendValue(value.MakeBytes([]byte("a")), value.MakeInt('b'))

	bytes, ok := got.Bytes()
	if !ok {
		t.Fatalf("appendValue returned %s, want byte slice", got.Kind())
	}
	if string(bytes) != "ab" {
		t.Fatalf("appendValue bytes = %q, want %q", string(bytes), "ab")
	}
}

func TestVM_AppendOpcodeUsesNativeByteFastPath(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAppend),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("append_byte", instr, 0, []byte("a"), byte('b'))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	bytes, ok := result.Bytes()
	if !ok {
		t.Fatalf("OpAppend returned %s, want byte slice", result.Kind())
	}
	if string(bytes) != "ab" {
		t.Fatalf("OpAppend bytes = %q, want %q", string(bytes), "ab")
	}
}
