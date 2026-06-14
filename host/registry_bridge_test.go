package host

import (
	"testing"

	"github.com/t04dJ14n9/gig/importer"
	"github.com/t04dJ14n9/gig/value"
)

func TestRegistryBridgeUsesFunctionDirectCall(t *testing.T) {
	reg := importer.NewRegistry()
	pkg := reg.RegisterPackage("example/direct", "direct")
	reflectCalled := false
	pkg.AddFunction("Add", func(int, int) int {
		reflectCalled = true
		return -1
	}, "", func(args []value.Value) ([]value.Value, error) {
		return []value.Value{value.MakeInt(args[0].Int() + args[1].Int())}, nil
	})

	env := FromRegistry(reg)
	fn, ok := env.LookupFunc("example/direct", "Add")
	if !ok {
		t.Fatal("LookupFunc did not find Add")
	}
	got, err := fn.Call([]value.Value{value.MakeInt(2), value.MakeInt(3)})
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if reflectCalled {
		t.Fatal("reflect function body was called; DirectCall wrapper was bypassed")
	}
	if len(got) != 1 || got[0].Int() != 5 {
		t.Fatalf("Add returned %v, want 5", got)
	}
	direct, ok := fn.(DirectFunction)
	if !ok {
		t.Fatal("LookupFunc result does not implement DirectFunction")
	}
	gotDirect, ok, err := direct.CallDirect([]value.Value{value.MakeInt(4), value.MakeInt(6)})
	if err != nil {
		t.Fatalf("CallDirect: %v", err)
	}
	if !ok || len(gotDirect) != 1 || gotDirect[0].Int() != 10 {
		t.Fatalf("CallDirect returned %v/%v, want 10/true", gotDirect, ok)
	}
}

func TestRegistryBridgeUsesMultiResultFunctionDirectCall(t *testing.T) {
	reg := importer.NewRegistry()
	pkg := reg.RegisterPackage("example/direct", "direct")
	reflectCalled := false
	pkg.AddFunction("Split", func(string) (string, string) {
		reflectCalled = true
		return "reflect", "fallback"
	}, "", func(args []value.Value) ([]value.Value, error) {
		s := args[0].Str()
		return []value.Value{value.MakeString(s[:1]), value.MakeString(s[1:])}, nil
	})

	env := FromRegistry(reg)
	fn, ok := env.LookupFunc("example/direct", "Split")
	if !ok {
		t.Fatal("LookupFunc did not find Split")
	}
	got, err := fn.Call([]value.Value{value.MakeString("go")})
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if reflectCalled {
		t.Fatal("reflect function body was called; DirectCall wrapper was bypassed")
	}
	if len(got) != 2 || got[0].Str() != "g" || got[1].Str() != "o" {
		t.Fatalf("Split returned %v, want [g o]", got)
	}
}

func TestRegistryBridgeUsesMethodDirectCall(t *testing.T) {
	reg := importer.NewRegistry()
	pkg := reg.RegisterPackage("example/direct", "direct")
	pkg.AddMethodDirectCall("Counter", "Len", func(recv value.Value, _ []value.Value) value.Value {
		return value.MakeInt(42 + recv.Int())
	})

	env := FromRegistry(reg)
	method, ok := env.LookupMethod("example/direct.Counter", "Len")
	if !ok {
		t.Fatal("LookupMethod did not find Counter.Len")
	}
	got, err := method.Call(value.MakeInt(8), nil)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if len(got) != 1 || got[0].Int() != 50 {
		t.Fatalf("Len returned %v, want 50", got)
	}
	direct, ok := method.(DirectMethod)
	if !ok {
		t.Fatal("LookupMethod result does not implement DirectMethod")
	}
	gotDirect, ok, err := direct.CallDirect(value.MakeInt(9), nil)
	if err != nil {
		t.Fatalf("CallDirect: %v", err)
	}
	if !ok || gotDirect.Int() != 51 {
		t.Fatalf("CallDirect returned %v/%v, want 51/true", gotDirect, ok)
	}
}
