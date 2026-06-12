package vm

import (
	"context"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestCallExternalSmallArityDoesNotAllocateArgSlice(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Constants: []any{
			&external.ExternalFuncInfo{
				DirectCall: func(args []value.Value) value.Value {
					return value.MakeInt(args[0].RawInt() + args[1].RawInt())
				},
			},
		},
	}
	v := &vm{
		program: prog,
		stack:   make([]value.Value, initialStackSize),
		ctx:     context.Background(),
		extCallCache: &externalCallCache{
			cache: make([]*extCallCacheEntry, len(prog.Constants)),
		},
	}

	allocs := testing.AllocsPerRun(1000, func() {
		v.sp = 2
		v.stack[0] = value.MakeInt(20)
		v.stack[1] = value.MakeInt(22)
		if err := v.callExternal(0, 2); err != nil {
			panic(err)
		}
		got := v.pop()
		if got.RawInt() != 42 {
			panic("unexpected direct external result")
		}
	})

	if allocs != 0 {
		t.Fatalf("callExternal allocations per small-arity call = %v, want 0", allocs)
	}
}
