package vm

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/external"
	"github.com/t04dJ14n9/gig/model/value"
)

type externalBoundaryPolicyHost struct{}

func (externalBoundaryPolicyHost) AcceptAny(any) {}

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

func TestCallExternalSmallArityDoesNotAllocateArgSlice(t *testing.T) {
	prog := &bytecode.CompiledProgram{
		Constants: []any{
			&external.ExternalFuncInfo{
				PkgPath:  "testing",
				FuncName: "Add",
				IsStdlib: true,
				DirectCall: func(args []value.Value) value.Value {
					return value.MakeInt(args[0].RawInt() + args[1].RawInt())
				},
				NumIn: 2,
			},
		},
	}
	prog.ResolveExternCalls()

	v := &vm{
		program: prog,
		stack:   make([]value.Value, initialStackSize),
		ctx:     context.Background(),
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

func TestExternalBoundaryPolicyTrustsStdlib(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	arg := value.MakeInterpretedInterface(value.MakeInt(1), "ScriptStruct", false)

	for _, rc := range []*bytecode.ResolvedCall{
		{PkgPath: "fmt", FuncName: "Sprint", IsStdlib: true, FnType: reflect.TypeOf(func(any) {})},
		{PkgPath: "strings", FuncName: "Contains", FnType: reflect.TypeOf(func(any) {})},
		{PkgPath: "main", FuncName: "F", FnType: reflect.TypeOf(func(any) {})},
	} {
		if err := v.validateExternalBoundary(rc, []value.Value{arg}); err != nil {
			t.Fatalf("trusted call %s.%s was rejected: %v", rc.PkgPath, rc.FuncName, err)
		}
	}

	methodInfo := &external.ExternalMethodInfo{PkgPath: "strings", MethodName: "Len", IsStdlib: true}
	args := []value.Value{value.FromInterface(externalBoundaryPolicyHost{}), arg}
	if err := v.validateExternalMethodBoundary(methodInfo, args); err != nil {
		t.Fatalf("trusted method was rejected: %v", err)
	}
}

func TestExternalBoundaryPolicyRequiresThirdPartyValidation(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	arg := value.MakeInterpretedInterface(value.MakeInt(1), "ScriptStruct", false)

	rc := &bytecode.ResolvedCall{
		PkgPath:  "example.com/host",
		FuncName: "AcceptAny",
		FnType:   reflect.TypeOf(func(any) {}),
	}
	if err := v.validateExternalBoundary(rc, []value.Value{arg}); err == nil {
		t.Fatal("third-party call skipped boundary validation")
	}

	methodInfo := &external.ExternalMethodInfo{PkgPath: "example.com/host", MethodName: "AcceptAny"}
	args := []value.Value{value.FromInterface(externalBoundaryPolicyHost{}), arg}
	if err := v.validateExternalMethodBoundary(methodInfo, args); err == nil {
		t.Fatal("third-party method skipped boundary validation")
	}
}

func TestExternalBoundaryPolicyUnsafeOverrideSkipsValidation(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{AllowUnsafeTypePass: true}}
	arg := value.MakeInterpretedInterface(value.MakeInt(1), "ScriptStruct", false)

	rc := &bytecode.ResolvedCall{
		PkgPath:  "example.com/host",
		FuncName: "AcceptAny",
		FnType:   reflect.TypeOf(func(any) {}),
	}
	if err := v.validateExternalBoundary(rc, []value.Value{arg}); err != nil {
		t.Fatalf("unsafe override should skip function boundary validation: %v", err)
	}

	methodInfo := &external.ExternalMethodInfo{PkgPath: "example.com/host", MethodName: "AcceptAny"}
	args := []value.Value{value.FromInterface(externalBoundaryPolicyHost{}), arg}
	if err := v.validateExternalMethodBoundary(methodInfo, args); err != nil {
		t.Fatalf("unsafe override should skip method boundary validation: %v", err)
	}
}

func TestValidateExternalBoundaryRejectsInterpretedInterfaceToThirdPartyAny(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	rc := &bytecode.ResolvedCall{
		PkgPath:  "example.com/host",
		FuncName: "AcceptAny",
		FnType:   reflect.TypeOf(func(any) {}),
	}
	arg := value.MakeInterpretedInterface(value.MakeInt(1), "ScriptStruct", false)

	err := v.validateExternalBoundary(rc, []value.Value{arg})
	if err == nil {
		t.Fatal("expected interpreted interface to be rejected")
	}
	if !strings.Contains(err.Error(), `interpreter-defined type "ScriptStruct"`) {
		t.Fatalf("error = %q, want ScriptStruct boundary diagnostic", err.Error())
	}
}

func TestValidateExternalBoundaryRejectsInterpretedFuncToThirdPartyAny(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	rc := &bytecode.ResolvedCall{
		PkgPath:  "example.com/host",
		FuncName: "AcceptAny",
		FnType:   reflect.TypeOf(func(any) {}),
	}
	arg := value.MakeFunc(&Closure{Fn: &bytecode.CompiledFunction{Name: "Callback"}})

	err := v.validateExternalBoundary(rc, []value.Value{arg})
	if err == nil {
		t.Fatal("expected interpreted function to be rejected")
	}
	if !strings.Contains(err.Error(), `interpreter-defined type "func Callback"`) {
		t.Fatalf("error = %q, want func Callback boundary diagnostic", err.Error())
	}
}

func TestValidateExternalBoundaryAllowsTypedInterpretedFuncCallback(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	rc := &bytecode.ResolvedCall{
		PkgPath:  "example.com/host",
		FuncName: "AcceptIntFunc",
		FnType:   reflect.TypeOf(func(func(int) int) {}),
	}
	arg := value.MakeFunc(&Closure{Fn: &bytecode.CompiledFunction{Name: "Callback"}})

	if err := v.validateExternalBoundary(rc, []value.Value{arg}); err != nil {
		t.Fatalf("typed callback should be allowed: %v", err)
	}
}

func TestValidateExternalMethodBoundaryRejectsInterpretedInterfaceArg(t *testing.T) {
	v := &vm{program: &bytecode.CompiledProgram{}}
	methodInfo := &external.ExternalMethodInfo{
		PkgPath:    "example.com/host",
		MethodName: "AcceptAny",
	}
	args := []value.Value{
		value.FromInterface(externalBoundaryPolicyHost{}),
		value.MakeInterpretedInterface(value.MakeInt(1), "ScriptStruct", false),
	}

	err := v.validateExternalMethodBoundary(methodInfo, args)
	if err == nil {
		t.Fatal("expected interpreted method argument to be rejected")
	}
	if !strings.Contains(err.Error(), `interpreter-defined type "ScriptStruct"`) {
		t.Fatalf("error = %q, want ScriptStruct boundary diagnostic", err.Error())
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

func TestBuildReflectArgsExactVariadicArgUsesElementType(t *testing.T) {
	fnType := reflect.TypeOf(func(prefix string, parts ...string) {})

	got := buildReflectArgs([]value.Value{
		value.MakeString("p"),
		value.MakeString("x"),
	}, fnType)

	if len(got) != 2 {
		t.Fatalf("buildReflectArgs returned %d args, want 2", len(got))
	}
	if got[1].Type() != reflect.TypeOf("") {
		t.Fatalf("exact variadic arg type = %v, want string element type", got[1].Type())
	}
	if got[1].String() != "x" {
		t.Fatalf("exact variadic arg value = %q, want x", got[1].String())
	}
}

func TestBuildReflectArgsExpandsPackedVariadicReflectSlice(t *testing.T) {
	fnType := reflect.TypeOf(func(prefix string, parts ...string) {})

	got := buildReflectArgs([]value.Value{
		value.MakeString("p"),
		value.FromInterface([]string{"x", "y"}),
	}, fnType)

	if len(got) != 3 {
		t.Fatalf("buildReflectArgs returned %d args, want 3", len(got))
	}
	for i, arg := range got {
		if arg.Type() != reflect.TypeOf("") {
			t.Fatalf("arg %d type = %v, want string", i, arg.Type())
		}
	}
	if got[1].String() != "x" || got[2].String() != "y" {
		t.Fatalf("expanded variadic values = %q, %q; want x, y", got[1].String(), got[2].String())
	}
}
