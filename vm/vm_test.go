package vm

import (
	"context"
	"sync/atomic"
	"testing"

	"gig/bytecode"
	"gig/value"
)

// ---------------------------------------------------------------------------
// Stack operations
// ---------------------------------------------------------------------------

func TestPushPop(t *testing.T) {
	vm := &VM{
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
	vm := &VM{
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
// Frame creation
// ---------------------------------------------------------------------------

func TestNewFrame(t *testing.T) {
	fn := &bytecode.CompiledFunction{
		Name:      "add",
		NumLocals: 4,
		NumParams: 2,
	}
	args := []value.Value{value.MakeInt(10), value.MakeInt(20)}
	f := newFrame(fn, 0, args, nil)

	if f.fn != fn {
		t.Error("frame fn mismatch")
	}
	if len(f.locals) != 4 {
		t.Fatalf("locals len = %d, want 4", len(f.locals))
	}
	// First two locals should be the args.
	if f.locals[0].Int() != 10 || f.locals[1].Int() != 20 {
		t.Errorf("locals = [%v, %v], want [10, 20]", f.locals[0], f.locals[1])
	}
}

// ---------------------------------------------------------------------------
// Closure
// ---------------------------------------------------------------------------

func TestClosureStruct(t *testing.T) {
	fn := &bytecode.CompiledFunction{Name: "closure_fn"}
	free1 := value.MakeInt(42)
	cl := &Closure{
		Fn:       fn,
		FreeVars: []*value.Value{&free1},
	}
	if cl.Fn.Name != "closure_fn" {
		t.Error("closure fn name")
	}
	if cl.FreeVars[0].Int() != 42 {
		t.Errorf("free var = %d, want 42", cl.FreeVars[0].Int())
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
	if len(v.globals) != 2 {
		t.Errorf("globals len = %d, want 2", len(v.globals))
	}
	if len(v.stack) != 1024 {
		t.Errorf("stack len = %d, want 1024", len(v.stack))
	}
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

// ---------------------------------------------------------------------------
// Goroutine tracking
// ---------------------------------------------------------------------------

func TestStartAndWaitGoroutines(t *testing.T) {
	var counter int64
	for i := 0; i < 5; i++ {
		StartGoroutine(func() {
			atomic.AddInt64(&counter, 1)
		})
	}
	WaitGoroutines()
	if atomic.LoadInt64(&counter) != 5 {
		t.Errorf("counter = %d, want 5", counter)
	}
}

// ---------------------------------------------------------------------------
// VM Registry
// ---------------------------------------------------------------------------

func TestVMRegistry(t *testing.T) {
	vm := &VM{stack: make([]value.Value, 8)}
	id := RegisterVM(vm)
	if id <= 0 {
		t.Fatalf("RegisterVM returned %d", id)
	}

	vmRegistryMutex.Lock()
	got, ok := vmRegistry[id]
	vmRegistryMutex.Unlock()
	if !ok || got != vm {
		t.Error("VM not found in registry")
	}

	UnregisterVM(id)
	vmRegistryMutex.Lock()
	_, ok = vmRegistry[id]
	vmRegistryMutex.Unlock()
	if ok {
		t.Error("VM should have been unregistered")
	}
}

// ---------------------------------------------------------------------------
// Child VM creation
// ---------------------------------------------------------------------------

func TestNewChildVM(t *testing.T) {
	prog := &bytecode.Program{
		Functions: map[string]*bytecode.CompiledFunction{},
		Globals:   map[string]int{"a": 0},
	}
	parent := New(prog)
	parent.globals[0] = value.MakeInt(99)
	parent.ctx = context.Background()

	child := parent.newChildVM()
	if child.program != parent.program {
		t.Error("child should share parent's program")
	}
	if child.globalsPtr == nil {
		t.Fatal("child globalsPtr should not be nil")
	}
	// Child should see the parent's globals through the pointer.
	globals := child.getGlobals()
	if globals[0].Int() != 99 {
		t.Errorf("child globals[0] = %d, want 99", globals[0].Int())
	}
}
