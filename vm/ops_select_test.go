package vm

import (
	"context"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestVM_SelectNonBlockingDefaultTuple(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpSelect), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	meta := bytecode.SelectMeta{
		NumStates: 0,
		Blocking:  false,
		NumRecv:   0,
	}
	prog, name := buildProg("select_default", instr, 0, meta)
	v := New(prog)

	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	tuple, ok := result.Interface().([]value.Value)
	if !ok {
		t.Fatalf("OpSelect result = %T, want []value.Value", result.Interface())
	}
	if len(tuple) != 2 {
		t.Fatalf("OpSelect tuple len = %d, want 2", len(tuple))
	}
	if got := tuple[0].Int(); got != -1 {
		t.Fatalf("OpSelect chosen index = %d, want -1 for default", got)
	}
	if got := tuple[1].Bool(); got {
		t.Fatalf("OpSelect recvOK = %v, want false for default", got)
	}
}
