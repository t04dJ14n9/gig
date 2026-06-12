package runner

import (
	"context"
	"testing"
	"time"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ---------------------------------------------------------------------------
// Helper: build a minimal compiled program for testing
// ---------------------------------------------------------------------------

// buildTestProgram creates a minimal CompiledProgram with a single function
// that returns a constant integer value.
func buildTestProgram(funcName string, constVal int64) *bytecode.CompiledProgram {
	fn := &bytecode.CompiledFunction{
		Name: funcName,
		Instructions: []byte{
			byte(bytecode.OpConst), 0, 0, // push constant 0
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	return &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{funcName: fn},
		Constants: []any{constVal},
		Globals:   map[string]int{},
	}
}

// buildTestProgramWithGlobals creates a CompiledProgram with a function
// that reads a global and returns it.
func buildTestProgramWithGlobals(funcName string, globalCount int, globalIdx int) *bytecode.CompiledProgram {
	hi, lo := uint16(globalIdx)>>8, uint16(globalIdx)&0xFF
	fn := &bytecode.CompiledFunction{
		Name: funcName,
		Instructions: []byte{
			byte(bytecode.OpGlobal), byte(hi), byte(lo),
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	return &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{funcName: fn},
		Constants: []any{},
		Globals:   map[string]int{"g": 0},
	}
}

// ---------------------------------------------------------------------------
// New Runner
// ---------------------------------------------------------------------------

func TestNewRunner(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil)
	if r == nil {
		t.Fatal("New returned nil")
	}
	r.Close()
}

func TestNewRunnerStateful(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil, WithStatefulGlobals())
	if r == nil {
		t.Fatal("New returned nil")
	}
	if !r.stateful {
		t.Error("expected stateful=true")
	}
	if r.shared == nil {
		t.Error("expected shared globals to be initialized")
	}
	r.Close()
}

func TestNewRunnerWithMaxGoroutines(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil, WithMaxGoroutines(100))
	if r == nil {
		t.Fatal("New returned nil")
	}
	r.Close()
}

// ---------------------------------------------------------------------------
// Runner.Run
// ---------------------------------------------------------------------------

func TestRunnerRunBasic(t *testing.T) {
	prog := buildTestProgram("compute", 42)
	r := New(prog, nil)
	defer r.Close()

	result, err := r.Run("compute")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if result.(int64) != 42 {
		t.Errorf("result = %v, want 42", result)
	}
}

func TestRunnerRunFunctionNotFound(t *testing.T) {
	prog := buildTestProgram("compute", 42)
	r := New(prog, nil)
	defer r.Close()

	_, err := r.Run("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent function")
	}
}

func TestRunnerRunRejectsWrongArity(t *testing.T) {
	prog := buildTestProgram("compute", 42)
	prog.Functions["two"] = &bytecode.CompiledFunction{
		Name: "two",
		Instructions: []byte{
			byte(bytecode.OpConst), 0, 0,
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 2,
		NumParams: 2,
		MaxStack:  1,
	}
	prog.Functions["variadic"] = &bytecode.CompiledFunction{
		Name: "variadic",
		Instructions: []byte{
			byte(bytecode.OpConst), 0, 0,
			byte(bytecode.OpReturnVal),
		},
		NumLocals:  2,
		NumParams:  2,
		IsVariadic: true,
		MaxStack:   1,
	}

	r := New(prog, nil)
	defer r.Close()

	if _, err := r.Run("compute", 1); err == nil {
		t.Fatal("expected error for extra argument to non-variadic function")
	}
	if _, err := r.Run("two", 1); err == nil {
		t.Fatal("expected error for missing argument to non-variadic function")
	}
	if _, err := r.Run("variadic"); err == nil {
		t.Fatal("expected error for missing fixed argument to variadic function")
	}
}

// ---------------------------------------------------------------------------
// Runner.RunWithContext
// ---------------------------------------------------------------------------

func TestRunnerRunWithContext(t *testing.T) {
	prog := buildTestProgram("compute", 42)
	r := New(prog, nil)
	defer r.Close()

	ctx := context.Background()
	result, err := r.RunWithContext(ctx, "compute")
	if err != nil {
		t.Fatalf("RunWithContext error: %v", err)
	}
	if result.(int64) != 42 {
		t.Errorf("result = %v, want 42", result)
	}
}

func TestRunnerRunWithContextCancelled(t *testing.T) {
	// Build a program with an infinite loop
	fn := &bytecode.CompiledFunction{
		Name: "loop",
		Instructions: []byte{
			byte(bytecode.OpJump), 0, 0, // jump to offset 0 (infinite loop)
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{"loop": fn},
		Constants: []any{},
		Globals:   map[string]int{},
	}

	r := New(prog, nil)
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := r.RunWithContext(ctx, "loop")
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestRunnerRunWithContextTimeout(t *testing.T) {
	fn := &bytecode.CompiledFunction{
		Name: "loop",
		Instructions: []byte{
			byte(bytecode.OpJump), 0, 0,
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{"loop": fn},
		Constants: []any{},
		Globals:   map[string]int{},
	}

	r := New(prog, nil)
	defer r.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := r.RunWithContext(ctx, "loop")
	if err == nil {
		t.Error("expected timeout error")
	}
}

// ---------------------------------------------------------------------------
// Runner.RunWithValues
// ---------------------------------------------------------------------------

func TestRunnerRunWithValues(t *testing.T) {
	// Build a program that takes 2 params and returns their sum
	// param0 = locals[0], param1 = locals[1]
	// OpLocal 0, OpLocal 1, OpAdd, OpReturnVal
	hi0, lo0 := uint16(0)>>8, uint16(0)&0xFF
	hi1, lo1 := uint16(1)>>8, uint16(1)&0xFF
	fn := &bytecode.CompiledFunction{
		Name: "add",
		Instructions: []byte{
			byte(bytecode.OpLocal), byte(hi0), byte(lo0),
			byte(bytecode.OpLocal), byte(hi1), byte(lo1),
			byte(bytecode.OpAdd),
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 2,
		NumParams: 2,
		MaxStack:  4,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{"add": fn},
		Constants: []any{},
		Globals:   map[string]int{},
	}

	r := New(prog, nil)
	defer r.Close()

	args := []value.Value{value.MakeInt(10), value.MakeInt(32)}
	result, err := r.RunWithValues(context.Background(), "add", args)
	if err != nil {
		t.Fatalf("RunWithValues error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// ---------------------------------------------------------------------------
// Runner.Close
// ---------------------------------------------------------------------------

func TestRunnerClose(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil)
	// Should not panic
	r.Close()
	// Double close should also not panic
	r.Close()
}

// ---------------------------------------------------------------------------
// Runner.InternalProgram
// ---------------------------------------------------------------------------

func TestRunnerInternalProgram(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil)
	defer r.Close()

	ip := r.InternalProgram()
	if ip == nil {
		t.Fatal("InternalProgram returned nil")
	}
	if ip != prog {
		t.Error("InternalProgram returned wrong program")
	}
}

// ---------------------------------------------------------------------------
// Runner.Wait / WaitContext
// ---------------------------------------------------------------------------

func TestRunnerWait(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil)
	defer r.Close()

	// Wait with no running goroutines should return immediately
	r.Wait()
}

func TestRunnerWaitContext(t *testing.T) {
	prog := buildTestProgram("test", 42)
	r := New(prog, nil)
	defer r.Close()

	ctx := context.Background()
	err := r.WaitContext(ctx)
	if err != nil {
		t.Errorf("WaitContext error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ExecuteInit
// ---------------------------------------------------------------------------

func TestExecuteInitNoInit(t *testing.T) {
	prog := buildTestProgram("main", 42)
	globals, err := ExecuteInit(prog)
	if err != nil {
		t.Fatalf("ExecuteInit error: %v", err)
	}
	if globals != nil {
		t.Errorf("expected nil globals for program without init(), got %v", globals)
	}
}

func TestExecuteInitWithInit(t *testing.T) {
	// Build a program with an init function that sets a global
	// init: OpConst 0, OpSetGlobal 0, OpReturn
	// main: OpGlobal 0, OpReturnVal
	hi0, lo0 := uint16(0)>>8, uint16(0)&0xFF
	hiG, loG := uint16(0)>>8, uint16(0)&0xFF

	initFn := &bytecode.CompiledFunction{
		Name: "init",
		Instructions: []byte{
			byte(bytecode.OpConst), byte(hi0), byte(lo0),
			byte(bytecode.OpSetGlobal), byte(hiG), byte(loG),
			byte(bytecode.OpReturn),
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	mainFn := &bytecode.CompiledFunction{
		Name: "main",
		Instructions: []byte{
			byte(bytecode.OpGlobal), byte(hiG), byte(loG),
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{
			"init": initFn,
			"main": mainFn,
		},
		Constants: []any{int64(99)},
		Globals:   map[string]int{"x": 0},
	}

	globals, err := ExecuteInit(prog)
	if err != nil {
		t.Fatalf("ExecuteInit error: %v", err)
	}
	if globals == nil {
		t.Fatal("expected non-nil globals after init()")
	}
	if len(globals) != 1 {
		t.Fatalf("globals len = %d, want 1", len(globals))
	}
	if globals[0].Int() != 99 {
		t.Errorf("globals[0] = %d, want 99", globals[0].Int())
	}
}

// ---------------------------------------------------------------------------
// UnwrapResult
// ---------------------------------------------------------------------------

func TestUnwrapResultSingleValue(t *testing.T) {
	v := value.MakeInt(42)
	result := UnwrapResult(v)
	// value.MakeInt stores int, so Interface() returns int
	if result.(int) != 42 {
		t.Errorf("result = %v, want 42", result)
	}
}

func TestUnwrapResultString(t *testing.T) {
	v := value.MakeString("hello")
	result := UnwrapResult(v)
	if result.(string) != "hello" {
		t.Errorf("result = %v, want 'hello'", result)
	}
}

func TestUnwrapResultMultiValue(t *testing.T) {
	vals := []value.Value{value.MakeInt(1), value.MakeString("two")}
	v := value.FromInterface(vals)
	result := UnwrapResult(v)
	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("result type = %T, want []any", result)
	}
	if len(arr) != 2 {
		t.Fatalf("result len = %d, want 2", len(arr))
	}
	if arr[0].(int) != 1 {
		t.Errorf("arr[0] = %v, want 1", arr[0])
	}
	if arr[1].(string) != "two" {
		t.Errorf("arr[1] = %v, want 'two'", arr[1])
	}
}

func TestUnwrapResultNil(t *testing.T) {
	v := value.MakeNil()
	result := UnwrapResult(v)
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
}

// ---------------------------------------------------------------------------
// Runner stateful mode
// ---------------------------------------------------------------------------

func TestRunnerStatefulGlobals(t *testing.T) {
	// Verify stateful globals mode is created without error
	// and that the shared globals object is initialized.
	prog := buildTestProgram("test", 42)
	r := New(prog, nil, WithStatefulGlobals())
	defer r.Close()

	if !r.stateful {
		t.Error("expected stateful=true")
	}
	if r.shared == nil {
		t.Error("expected shared globals to be initialized")
	}
	// Verify we can execute a function
	_, err := r.Run("test")
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// DefaultTimeout
// ---------------------------------------------------------------------------

func TestDefaultTimeout(t *testing.T) {
	if DefaultTimeout != 10*time.Second {
		t.Errorf("DefaultTimeout = %v, want 10s", DefaultTimeout)
	}
}
