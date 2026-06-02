package vm

import (
	"reflect"
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

func TestRunIntSliceFallbackHandlesReflectIntSlice(t *testing.T) {
	ints := []int{3, 1, 2}
	v := &vm{stack: make([]value.Value, 8)}
	locals := []value.Value{
		value.MakeFromReflect(reflect.ValueOf(ints)),
		value.MakeInt(1),
		value.MakeNil(),
	}
	intLocals := []int64{0, 1, 9}

	sp, _, err := v.runIntSliceGetFallback(nil, locals, intLocals, 0, 1, 2, 0)
	if err != nil {
		t.Fatalf("runIntSliceGetFallback returned error: %v", err)
	}
	if sp != 0 {
		t.Fatalf("runIntSliceGetFallback sp = %d, want 0", sp)
	}
	if intLocals[2] != 1 || locals[2].RawInt() != 1 {
		t.Fatalf("fallback get = intLocals %d / locals %d, want 1", intLocals[2], locals[2].RawInt())
	}

	intLocals[2] = 7
	sp, _, err = v.runIntSliceSetFallback(nil, locals, intLocals, 0, 1, 2, 0)
	if err != nil {
		t.Fatalf("runIntSliceSetFallback returned error: %v", err)
	}
	if sp != 0 {
		t.Fatalf("runIntSliceSetFallback sp = %d, want 0", sp)
	}
	if ints[1] != 7 {
		t.Fatalf("fallback set left slice as %v, want index 1 updated to 7", ints)
	}
}

func TestRunIntSliceFallbackConvertsIndexPanic(t *testing.T) {
	ints := []int{3, 1, 2}
	v := &vm{stack: make([]value.Value, 8)}
	locals := []value.Value{
		value.MakeFromReflect(reflect.ValueOf(ints)),
		value.MakeInt(99),
		value.MakeNil(),
	}
	intLocals := []int64{0, 99, 0}

	sp, _, err := v.runIntSliceGetFallback(nil, locals, intLocals, 0, 1, 2, 0)
	if err != nil {
		t.Fatalf("runIntSliceGetFallback returned error: %v", err)
	}
	if sp != 0 {
		t.Fatalf("runIntSliceGetFallback sp = %d, want 0", sp)
	}
	if !v.panicking {
		t.Fatal("runIntSliceGetFallback did not convert index panic into VM panic")
	}
	if v.panicVal.String() == "" {
		t.Fatal("runIntSliceGetFallback stored empty panic value")
	}
}
