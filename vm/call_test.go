package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ---------------------------------------------------------------------------
// ExternalCallCancelledError tests
// ---------------------------------------------------------------------------

func TestExternalCallCancelledError(t *testing.T) {
	cause := &assertError{msg: "test cause"}
	err := &ExternalCallCancelledError{Cause: cause}

	if err.Error() != "external call cancelled: test cause" {
		t.Errorf("Error() = %q, want %q", err.Error(), "external call cancelled: test cause")
	}

	if !errors.Is(err.Unwrap(), cause) {
		t.Error("Unwrap() did not return the cause")
	}
}

type assertError struct {
	msg string
}

func (e *assertError) Error() string { return e.msg }

// ---------------------------------------------------------------------------
// GoroutineTracker tests
// ---------------------------------------------------------------------------

func TestGoroutineTrackerBasic(t *testing.T) {
	tracker := NewGoroutineTracker()

	// Start a goroutine that completes quickly
	done := make(chan bool, 1)
	tracker.Start(func() { //nolint:errcheck
		done <- true
	})

	// Wait for goroutines - should complete
	tracker.Wait()

	select {
	case <-done:
		// Success - goroutine completed
	default:
		t.Error("Goroutine did not complete")
	}
}

func TestGoroutineTrackerWithContext(t *testing.T) {
	tracker := NewGoroutineTracker()
	ctx := context.Background()

	// Start a quick goroutine
	tracker.Start(func() { //nolint:errcheck
		// Do nothing, just return
	})

	err := tracker.WaitContext(ctx)
	if err != nil {
		t.Errorf("WaitContext failed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// VM Execution tests via bytecode
// ---------------------------------------------------------------------------

func buildSimpleProg(name string, instr []byte, constants []any) (*bytecode.CompiledProgram, string) {
	fn := &bytecode.CompiledFunction{
		Name:         name,
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		NumFreeVars:  0,
		MaxStack:     8,
	}
	return &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{name: fn},
		Constants: constants,
		Globals:   map[string]int{},
	}, name
}

func TestExecuteBasic(t *testing.T) {
	// Test basic OpConst and OpReturnVal
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("test_basic", instr, []any{int64(42)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

func TestExecuteAddIntegers(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("add_test", instr, []any{int64(10), int64(20)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 30 {
		t.Errorf("result = %d, want 30", result.Int())
	}
}

func TestExecuteSubIntegers(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpSub),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("sub_test", instr, []any{int64(100), int64(30)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 70 {
		t.Errorf("result = %d, want 70", result.Int())
	}
}

func TestExecuteMulIntegers(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpMul),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("mul_test", instr, []any{int64(6), int64(7)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

func TestExecuteDivIntegers(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpDiv),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("div_test", instr, []any{int64(42), int64(6)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 7 {
		t.Errorf("result = %d, want 7", result.Int())
	}
}

func TestExecuteModIntegers(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpMod),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("mod_test", instr, []any{int64(17), int64(5)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 2 {
		t.Errorf("result = %d, want 2", result.Int())
	}
}

func TestExecuteNegInteger(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpNeg),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("neg_test", instr, []any{int64(42)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != -42 {
		t.Errorf("result = %d, want -42", result.Int())
	}
}

func TestExecuteNotBool(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpTrue),
		byte(bytecode.OpNot),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("not_test", instr, nil)

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Kind() != value.KindBool {
		t.Errorf("result.Kind() = %v, want KindBool", result.Kind())
	}
	if result.Bool() != false {
		t.Errorf("result.Bool() = %v, want false", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// Comparison operations
// ---------------------------------------------------------------------------

func TestExecuteEqualTrue(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("eq_test", instr, []any{int64(42), int64(42)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Bool() != true {
		t.Errorf("result.Bool() = %v, want true", result.Bool())
	}
}

func TestExecuteEqualFalse(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("neq_test", instr, []any{int64(10), int64(20)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Bool() != false {
		t.Errorf("result.Bool() = %v, want false", result.Bool())
	}
}

func TestExecuteLessTrue(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpLess),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("less_test", instr, []any{int64(10), int64(20)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Bool() != true {
		t.Errorf("result.Bool() = %v, want true", result.Bool())
	}
}

func TestExecuteGreaterTrue(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpGreater),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("greater_test", instr, []any{int64(20), int64(10)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Bool() != true {
		t.Errorf("result.Bool() = %v, want true", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// Bitwise operations
// ---------------------------------------------------------------------------

func TestExecuteBitwiseAnd(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpAnd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("bitwise_and", instr, []any{int64(0xFF), int64(0x0F)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 15 {
		t.Errorf("result = %d, want 15", result.Int())
	}
}

func TestExecuteBitwiseOr(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpOr),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("bitwise_or", instr, []any{int64(0xF0), int64(0x0F)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 255 {
		t.Errorf("result = %d, want 255", result.Int())
	}
}

func TestExecuteBitwiseXor(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpXor),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("bitwise_xor", instr, []any{int64(0xFF), int64(0x0F)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 240 {
		t.Errorf("result = %d, want 240", result.Int())
	}
}

func TestExecuteShiftLeft(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpLsh),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("shift_left", instr, []any{int64(1), int64(2)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 4 {
		t.Errorf("result = %d, want 4", result.Int())
	}
}

func TestExecuteShiftRight(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpRsh),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("shift_right", instr, []any{int64(8), int64(2)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 2 {
		t.Errorf("result = %d, want 2", result.Int())
	}
}

// ---------------------------------------------------------------------------
// String operations
// ---------------------------------------------------------------------------

func TestExecuteAddStrings(t *testing.T) {
	hiA, loA := u16(0)
	hiB, loB := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hiA, loA,
		byte(bytecode.OpConst), hiB, loB,
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("add_strings", instr, []any{"hello", "world"})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.String() != "helloworld" {
		t.Errorf("result = %q, want %q", result.String(), "helloworld")
	}
}

// ---------------------------------------------------------------------------
// Nil and boolean constants
// ---------------------------------------------------------------------------

func TestExecutePushNil(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpNil),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("nil_test", instr, nil)

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Kind() != value.KindNil {
		t.Errorf("result.Kind() = %v, want KindNil", result.Kind())
	}
}

func TestExecutePushTrue(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpTrue),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("true_test", instr, nil)

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Kind() != value.KindBool {
		t.Errorf("result.Kind() = %v, want KindBool", result.Kind())
	}
	if result.Bool() != true {
		t.Errorf("result.Bool() = %v, want true", result.Bool())
	}
}

func TestExecutePushFalse(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpFalse),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("false_test", instr, nil)

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Kind() != value.KindBool {
		t.Errorf("result.Kind() = %v, want KindBool", result.Kind())
	}
	if result.Bool() != false {
		t.Errorf("result.Bool() = %v, want false", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// Pop and Dup
// ---------------------------------------------------------------------------

func TestExecutePop(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpPop), // discard top value
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("pop_test", instr, []any{int64(1), int64(2)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 1 {
		t.Errorf("result = %d, want 1", result.Int())
	}
}

func TestExecuteDup(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpDup), // duplicate top value
		byte(bytecode.OpAdd), // add them
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("dup_test", instr, []any{int64(21)})

	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// ---------------------------------------------------------------------------
// VM Reset test
// ---------------------------------------------------------------------------

func TestExecuteResetAfterExecution(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("reset_test", instr, []any{int64(99)})

	v := New(prog)
	_, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	v.Reset()
	// VM is now reset - verify it can execute again
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute after Reset error: %v", err)
	}

	if result.Int() != 99 {
		t.Errorf("result = %d, want 99", result.Int())
	}
}

// ---------------------------------------------------------------------------
// VM Pool test
// ---------------------------------------------------------------------------

func TestExecuteVMPoolGetPut(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildSimpleProg("pool_test", instr, []any{int64(7)})

	pool := NewVMPool(prog, nil, NewGoroutineTracker())

	// First Get: execute to prove the VM works
	vm1 := pool.Get()
	if vm1 == nil {
		t.Fatal("Get returned nil")
	}
	result, err := vm1.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 7 {
		t.Errorf("first result = %d, want 7", result.Int())
	}

	// Put it back and Get again; the reused VM must also execute correctly
	pool.Put(vm1)
	vm2 := pool.Get()
	if vm2 == nil {
		t.Fatal("second Get returned nil")
	}
	result2, err := vm2.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute (reused) error: %v", err)
	}
	if result2.Int() != 7 {
		t.Errorf("reused result = %d, want 7", result2.Int())
	}
	pool.Put(vm2)
}
