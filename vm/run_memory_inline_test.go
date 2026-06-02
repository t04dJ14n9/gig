package vm

import (
	"testing"

	"github.com/t04dJ14n9/gig/model/value"
)

func TestRunDerefLoadsValueSlotDirectly(t *testing.T) {
	slot := value.MakeInt(41)
	v := &vm{stack: make([]value.Value, 4)}
	v.stack[0] = value.FromInterface(&slot)

	sp, stack, err := v.runDeref(nil, 1)
	if err != nil {
		t.Fatalf("runDeref returned error: %v", err)
	}
	if sp != 1 {
		t.Fatalf("runDeref sp = %d, want 1", sp)
	}
	if got := stack[0].RawInt(); got != 41 {
		t.Fatalf("runDeref loaded %d, want 41", got)
	}
}

func TestRunSetDerefStoresValueSlotDirectly(t *testing.T) {
	slot := value.MakeInt(1)
	v := &vm{stack: make([]value.Value, 4)}
	v.stack[0] = value.FromInterface(&slot)
	v.stack[1] = value.MakeInt(42)

	if sp := v.runSetDeref(2); sp != 0 {
		t.Fatalf("runSetDeref sp = %d, want 0", sp)
	}
	if got := slot.RawInt(); got != 42 {
		t.Fatalf("runSetDeref stored %d, want 42", got)
	}
}
