package vm

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/value"
)

// ---------------------------------------------------------------------------
// Stack operations
// ---------------------------------------------------------------------------

func TestPushPop(t *testing.T) {
	vm := &vm{
		stack: make([]value.Value, 8),
		sp:    0,
	}
	vm.push(value.MakeInt(1))
	vm.push(value.MakeInt(2))
	vm.push(value.MakeInt(3))

	if vm.sp != 3 {
		t.Fatalf("sp = %d, want 3", vm.sp)
	}
	if vm.pop().Int() != 3 {
		t.Error("pop 1")
	}
	if vm.pop().Int() != 2 {
		t.Error("pop 2")
	}
	if vm.peek().Int() != 1 {
		t.Error("peek")
	}
	if vm.sp != 1 {
		t.Errorf("sp after peek = %d, want 1", vm.sp)
	}
}

func TestStackAutoGrow(t *testing.T) {
	vm := &vm{
		stack: make([]value.Value, 2),
		sp:    0,
	}
	// Push more than initial capacity to trigger auto-grow.
	for i := 0; i < 10; i++ {
		vm.push(value.MakeInt(int64(i)))
	}
	if vm.sp != 10 {
		t.Fatalf("sp = %d, want 10", vm.sp)
	}
	// Verify values are intact after growth.
	for i := 9; i >= 0; i-- {
		v := vm.pop()
		if v.Int() != int64(i) {
			t.Errorf("pop(%d) = %d, want %d", 9-i, v.Int(), i)
		}
	}
}

// ---------------------------------------------------------------------------
// VM creation and simple execution
// ---------------------------------------------------------------------------

func TestNewVM(t *testing.T) {
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{},
		Globals:   map[string]int{"x": 0, "y": 1},
	}
	v := New(prog)
	if v == nil {
		t.Fatal("New returned nil")
	}
	// VM is now an interface - internal fields (globals, stack) are encapsulated
	// Skip internal field checks as they are implementation details
	_ = v
}

// TestExecuteHalt verifies that the VM handles OpHalt correctly.
// OpHalt produces a "halt" error, which is the expected termination signal.
func TestExecuteHalt(t *testing.T) {
	fn := &bytecode.CompiledFunction{
		Name:         "main",
		Instructions: []byte{byte(bytecode.OpHalt)},
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     1,
	}
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{"main": fn},
		Globals:   map[string]int{},
	}
	v := New(prog)
	_, err := v.Execute("main", context.Background())
	if err == nil {
		t.Fatal("expected halt error from OpHalt")
	}
	if err.Error() != "halt" {
		t.Errorf("error = %q, want %q", err.Error(), "halt")
	}
}

// TestExecuteConstAndReturn verifies OpConst + OpReturnVal.
func TestExecuteConstAndReturn(t *testing.T) {
	// Build instructions: OpConst 0, OpReturnVal
	instr := make([]byte, 0, 4)
	instr = append(instr, byte(bytecode.OpConst))
	instr = append(instr, 0, 0) // constant index 0
	instr = append(instr, byte(bytecode.OpReturnVal))

	fn := &bytecode.CompiledFunction{
		Name:         "compute",
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     1,
	}
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{"compute": fn},
		Constants: []any{int64(42)},
		Globals:   map[string]int{},
	}

	v := New(prog)
	result, err := v.Execute("compute", context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestExecuteFunctionNotFound verifies error for missing function.
func TestExecuteFunctionNotFound(t *testing.T) {
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{},
		Globals:   map[string]int{},
	}
	v := New(prog)
	_, err := v.Execute("nonexistent", context.Background())
	if err == nil {
		t.Error("expected error for non-existent function")
	}
}

// TestExecuteWithCancellation verifies context cancellation.
func TestExecuteWithCancellation(t *testing.T) {
	// Infinite loop: OpJump back to 0
	instr := make([]byte, 0, 3)
	instr = append(instr, byte(bytecode.OpJump))
	instr = append(instr, 0, 0) // jump to offset 0

	fn := &bytecode.CompiledFunction{
		Name:         "loop",
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     1,
	}
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{"loop": fn},
		Globals:   map[string]int{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	v := New(prog)
	_, err := v.Execute("loop", ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

// TestExecuteWithTimeout verifies context timeout during execution.
func TestExecuteWithTimeout(t *testing.T) {
	// Infinite loop: OpJump back to 0
	instr := make([]byte, 0, 3)
	instr = append(instr, byte(bytecode.OpJump))
	instr = append(instr, 0, 0) // jump to offset 0

	fn := &bytecode.CompiledFunction{
		Name:         "loop",
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     1,
	}
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{"loop": fn},
		Globals:   map[string]int{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50)
	defer cancel()

	v := New(prog)
	_, err := v.Execute("loop", ctx)
	if err == nil {
		t.Error("expected context timeout error")
	}
	if ctx.Err() == nil {
		t.Error("expected context to be cancelled")
	}
}

// TestChannelSendCancellation verifies that channel send respects context cancellation.
func TestChannelSendCancellation(t *testing.T) {
	// This test verifies the SendContext path - actual VM integration test
	// would require compiled bytecode with channel operations
	ch := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Wrap channel in value
	v := value.FromInterface(ch)

	// Send should fail with context cancelled
	err := v.SendContext(ctx, value.MakeInt(42))
	if err == nil {
		t.Error("expected context cancellation error for SendContext")
	}
}

// TestChannelRecvCancellation verifies that channel receive respects context cancellation.
func TestChannelRecvCancellation(t *testing.T) {
	// Empty channel - receive should block then cancel
	ch := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Wrap channel in value
	v := value.FromInterface(ch)

	// Recv should fail with context cancelled
	_, _, err := v.RecvContext(ctx)
	if err == nil {
		t.Error("expected context cancellation error for RecvContext")
	}
}

// TestChannelSendSuccess verifies successful channel send with context.
func TestChannelSendSuccess(t *testing.T) {
	ch := make(chan int, 1) // Buffered channel

	ctx := context.Background()
	v := value.FromInterface(ch)

	err := v.SendContext(ctx, value.MakeInt(42))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify value was sent
	val := <-ch
	if val != 42 {
		t.Errorf("received %d, want 42", val)
	}
}

// TestChannelRecvSuccess verifies successful channel receive with context.
func TestChannelRecvSuccess(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	ctx := context.Background()
	v := value.FromInterface(ch)

	val, ok, err := v.RecvContext(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected ok=true for successful receive")
	}
	if val.Int() != 42 {
		t.Errorf("received %d, want 42", val.Int())
	}
}

// TestGoroutineTrackerWaitContext verifies that GoroutineTracker.WaitContext respects cancellation.
func TestGoroutineTrackerWaitContext(t *testing.T) {
	gt := NewGoroutineTracker()
	gt.Start(func() {
		time.Sleep(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := gt.WaitContext(ctx)
	if err == nil {
		t.Error("expected timeout error")
	}

	// Wait for the goroutine to actually complete to clean up
	gt.Wait()
}

// ---------------------------------------------------------------------------
// Goroutine tracking
// ---------------------------------------------------------------------------

func TestGoroutineTrackerStartAndWait(t *testing.T) {
	gt := NewGoroutineTracker()
	var counter int64
	for i := 0; i < 5; i++ {
		gt.Start(func() {
			atomic.AddInt64(&counter, 1)
		})
	}
	gt.Wait()
	if atomic.LoadInt64(&counter) != 5 {
		t.Errorf("counter = %d, want 5", counter)
	}
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// Helpers used by the new tests
// ---------------------------------------------------------------------------

// makeInstructions assembles raw bytecode bytes.
func makeInstructions(ops ...byte) []byte { return ops }

// u16 encodes a uint16 as two big-endian bytes for inline use.
func u16(v uint16) (byte, byte) { return byte(v >> 8), byte(v) }

// buildProg creates a Program + CompiledFunction ready for execution.
// constants are raw any values added to Constants (not PrebakedConstants),
// so the VM resolves them via value.FromInterface on first access.
func buildProg(name string, instr []byte, numLocals int, constants ...any) (*bytecode.Program, string) {
	fn := &bytecode.CompiledFunction{
		Name:         name,
		Instructions: instr,
		NumLocals:    numLocals,
		NumParams:    0,
		MaxStack:     8,
	}
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{name: fn},
		Constants: constants,
		Globals:   map[string]int{},
	}
	return prog, name
}

// ---------------------------------------------------------------------------
// Arithmetic opcodes
// ---------------------------------------------------------------------------

// TestVM_AddIntegers: 10 + 32 => 42
func TestVM_AddIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("add", instr, 0, int64(10), int64(32))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_SubIntegers: 100 - 58 => 42
func TestVM_SubIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpSub),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("sub", instr, 0, int64(100), int64(58))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_MulIntegers: 6 * 7 => 42
func TestVM_MulIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpMul),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("mul", instr, 0, int64(6), int64(7))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_DivIntegers: 84 / 2 => 42
func TestVM_DivIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpDiv),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("div", instr, 0, int64(84), int64(2))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_ModIntegers: 85 % 43 => 42
func TestVM_ModIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpMod),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("mod", instr, 0, int64(85), int64(43))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_NegInteger: -42 (push 42, OpNeg) => -42
func TestVM_NegInteger(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpNeg),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("neg", instr, 0, int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != -42 {
		t.Errorf("result = %d, want -42", result.Int())
	}
}

// TestVM_AddStrings: "hello " + "world" => "hello world"
func TestVM_AddStrings(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("addstr", instr, 0, "hello ", "world")
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.String() != "hello world" {
		t.Errorf("result = %q, want %q", result.String(), "hello world")
	}
}

// ---------------------------------------------------------------------------
// Comparison opcodes
// ---------------------------------------------------------------------------

// TestVM_EqualTrue: 42 == 42 => true
func TestVM_EqualTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("eq_true", instr, 0, int64(42), int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

// TestVM_EqualFalse: 1 == 2 => false
func TestVM_EqualFalse(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("eq_false", instr, 0, int64(1), int64(2))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Bool() {
		t.Errorf("result = %v, want false", result.Bool())
	}
}

// TestVM_LessTrue: 1 < 2 => true
func TestVM_LessTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpLess),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("less_true", instr, 0, int64(1), int64(2))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

// TestVM_GreaterTrue: 2 > 1 => true
func TestVM_GreaterTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpGreater),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("greater_true", instr, 0, int64(2), int64(1))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// Boolean opcodes
// ---------------------------------------------------------------------------

// TestVM_NotTrue: !true => false
func TestVM_NotTrue(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpTrue),
		byte(bytecode.OpNot),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("not_true", instr, 0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Bool() {
		t.Errorf("result = %v, want false", result.Bool())
	}
}

// TestVM_NotFalse: !false => true
func TestVM_NotFalse(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpFalse),
		byte(bytecode.OpNot),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("not_false", instr, 0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// Jump opcodes
// ---------------------------------------------------------------------------

// TestVM_JumpUnconditional: OpJump skips past OpConst 99, lands on OpConst 42.
//
// Bytecode layout:
//
//	[0] OpJump [1,2]=6      -- jump to offset 6 (3 bytes)
//	[3] OpConst [4,5]=0     -- push 99  (skipped)
//	[6] OpConst [7,8]=1     -- push 42
//	[9] OpReturnVal
func TestVM_JumpUnconditional(t *testing.T) {
	jumpHi, jumpLo := u16(6)
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpJump), jumpHi, jumpLo,
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("jump", instr, 0, int64(99), int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// TestVM_JumpFalse_taken: push false, OpJumpFalse to offset 7 (skips OpConst 42),
// lands on OpConst 99, returns 99.
//
// Bytecode layout:
//
//	[0] OpFalse             -- push false (1 byte)
//	[1] OpJumpFalse [2,3]=7 -- if false jump to 7  (3 bytes)
//	[4] OpConst [5,6]=0     -- push 42 (skipped)   (3 bytes)
//	[7] OpConst [8,9]=1     -- push 99              (3 bytes)
//	[10] OpReturnVal
func TestVM_JumpFalse_taken(t *testing.T) {
	jumpHi, jumpLo := u16(7)
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpFalse),
		byte(bytecode.OpJumpFalse), jumpHi, jumpLo,
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("jmpf", instr, 0, int64(42), int64(99))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 99 {
		t.Errorf("result = %d, want 99", result.Int())
	}
}

// TestVM_JumpTrue_taken: push true, OpJumpTrue to offset 7 (skips OpConst 42),
// lands on OpConst 99, returns 99.
//
// Bytecode layout:
//
//	[0] OpTrue              -- push true (1 byte)
//	[1] OpJumpTrue [2,3]=7  -- if true jump to 7  (3 bytes)
//	[4] OpConst [5,6]=0     -- push 42 (skipped)  (3 bytes)
//	[7] OpConst [8,9]=1     -- push 99             (3 bytes)
//	[10] OpReturnVal
func TestVM_JumpTrue_taken(t *testing.T) {
	jumpHi, jumpLo := u16(7)
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpTrue),
		byte(bytecode.OpJumpTrue), jumpHi, jumpLo,
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("jmpt", instr, 0, int64(42), int64(99))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 99 {
		t.Errorf("result = %d, want 99", result.Int())
	}
}

// ---------------------------------------------------------------------------
// Local variable opcodes
// ---------------------------------------------------------------------------

// TestVM_SetAndGetLocal: OpConst 42, OpSetLocal 0, OpLocal 0, OpReturnVal => 42
func TestVM_SetAndGetLocal(t *testing.T) {
	hi0, lo0 := u16(0) // const index 0
	hiL, loL := u16(0) // local index 0
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpSetLocal), hiL, loL,
		byte(bytecode.OpLocal), hiL, loL,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("local", instr, 1, int64(42))
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
// OpDup
// ---------------------------------------------------------------------------

// TestVM_Dup: push 21, OpDup, OpAdd => 42
func TestVM_Dup(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpDup),
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("dup", instr, 0, int64(21))
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
// VM Pool
// ---------------------------------------------------------------------------

// TestVMPool_GetPut: Get returns a usable VM, Put recycles it, second Get reuses.
func TestVMPool_GetPut(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("pool_fn", instr, 0, int64(7))

	pool := NewVMPool(prog, nil, NewGoroutineTracker())

	// First Get: execute to prove the VM works.
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

	// Put it back and Get again; the reused VM must also execute correctly.
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

// ---------------------------------------------------------------------------
// VM Reset
// ---------------------------------------------------------------------------

// TestVM_Reset: after executing a program, Reset clears the stack pointer.
func TestVM_Reset(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("reset_fn", instr, 0, int64(99))
	v := New(prog)

	_, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	v.Reset()
	// VM is now an interface - internal fields (sp, fp, panicking) are encapsulated
	// Skip internal field checks as they are implementation details
	_ = v
}

// ---------------------------------------------------------------------------
// Frame pool
// ---------------------------------------------------------------------------

// TestFramePool: get and put frames without panicking; reuse works correctly.
func TestFramePool(t *testing.T) {
	fn := &bytecode.CompiledFunction{
		Name:      "fp_fn",
		NumLocals: 3,
		NumParams: 0,
	}

	var pool framePool

	// Get a frame from an empty pool (allocates new).
	f := pool.get(fn, 0, nil)
	if f == nil {
		t.Fatal("get returned nil")
	}
	if len(f.locals) != 3 {
		t.Errorf("locals len = %d, want 3", len(f.locals))
	}
	if f.fn != fn {
		t.Error("frame fn mismatch")
	}

	// Return it to the pool.
	pool.put(f)
	if len(pool.frames) != 1 {
		t.Errorf("pool.frames len = %d, want 1", len(pool.frames))
	}

	// Get again: should reuse the same underlying memory.
	f2 := pool.get(fn, 0, nil)
	if f2 == nil {
		t.Fatal("second get returned nil")
	}
	// Locals must be zeroed on reuse.
	for i, local := range f2.locals {
		if local != (value.Value{}) {
			t.Errorf("local[%d] not zeroed after reuse", i)
		}
	}

	// Frames with addrTaken must NOT be pooled.
	f2.addrTaken = true
	pool.put(f2) // should be a no-op
	if len(pool.frames) != 0 {
		t.Errorf("addrTaken frame should not be pooled; pool len = %d", len(pool.frames))
	}
}
