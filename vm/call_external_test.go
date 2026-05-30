package vm

import (
	"context"
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

func TestCallExternalMissingResolvedCallReturnsError(t *testing.T) {
	v := &vm{
		program: &bytecode.CompiledProgram{},
		stack:   make([]value.Value, initialStackSize),
		ctx:     context.Background(),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("callExternal panicked for missing call entry: %v", r)
		}
	}()

	err := v.callExternal(99, 0)
	if err == nil {
		t.Fatal("expected error for missing external call entry, got nil")
	}
	if !strings.Contains(err.Error(), "unresolved external call") {
		t.Fatalf("error = %q, want unresolved external call", err.Error())
	}
}

func TestCallExternalInvalidFunctionReturnsError(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Constants: []any{
			&external.ExternalFuncInfo{Func: "not a function"},
		},
	}
	prog.ResolveExternCalls()

	v := &vm{
		program: prog,
		stack:   make([]value.Value, initialStackSize),
		ctx:     context.Background(),
	}

	err := v.callExternal(0, 0)
	if err == nil {
		t.Fatal("expected error for invalid external function, got nil")
	}
	if !strings.Contains(err.Error(), "invalid external function") {
		t.Fatalf("error = %q, want invalid external function", err.Error())
	}
}

func TestValidateExternalMethodBoundaryAllowsMissingReceiver(t *testing.T) {
	v := &vm{
		program: &bytecode.CompiledProgram{},
		ctx:     context.Background(),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("validateExternalMethodBoundary panicked with missing receiver: %v", r)
		}
	}()

	if err := v.validateExternalMethodBoundary(&external.ExternalMethodInfo{
		PkgPath:    "example.com/host",
		MethodName: "AcceptAny",
	}, nil); err != nil {
		t.Fatalf("validateExternalMethodBoundary returned error for missing receiver: %v", err)
	}
}
