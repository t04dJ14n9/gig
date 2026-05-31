package vm

import (
	"context"
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
	"github.com/t04dJ14n9/gig/model/value"
)

// ---------------------------------------------------------------------------
// Bitwise opcodes
// ---------------------------------------------------------------------------

func TestVM_AndIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAnd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("and", instr, 0, int64(0xFF), int64(0x0F))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 0x0F {
		t.Errorf("result = %d, want %d", result.Int(), 0x0F)
	}
}

func TestVM_OrIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpOr),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("or", instr, 0, int64(0xF0), int64(0x0F))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 0xFF {
		t.Errorf("result = %d, want %d", result.Int(), 0xFF)
	}
}

func TestVM_XorIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpXor),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("xor", instr, 0, int64(0xFF), int64(0x0F))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 0xF0 {
		t.Errorf("result = %d, want %d", result.Int(), 0xF0)
	}
}

func TestVM_AndNotIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAndNot),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("andnot", instr, 0, int64(0xFF), int64(0x0F))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 0xF0 {
		t.Errorf("result = %d, want %d", result.Int(), 0xF0)
	}
}

func TestVM_LshIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpLsh),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("lsh", instr, 0, int64(1), int64(4))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 16 {
		t.Errorf("result = %d, want 16", result.Int())
	}
}

func TestVM_RshIntegers(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpRsh),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("rsh", instr, 0, int64(16), int64(4))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 1 {
		t.Errorf("result = %d, want 1", result.Int())
	}
}

// ---------------------------------------------------------------------------
// Comparison opcodes (additional)
// ---------------------------------------------------------------------------

func TestVM_NotEqualTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpNotEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("neq_true", instr, 0, int64(1), int64(2))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

func TestVM_NotEqualFalse(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpNotEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("neq_false", instr, 0, int64(42), int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Bool() {
		t.Errorf("result = %v, want false", result.Bool())
	}
}

func TestVM_LessEqTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpLessEq),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("lesseq_true", instr, 0, int64(1), int64(2))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

func TestVM_LessEqEqual(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpLessEq),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("lesseq_eq", instr, 0, int64(42), int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

func TestVM_GreaterEqTrue(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpGreaterEq),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("greatereq_true", instr, 0, int64(2), int64(1))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("result = %v, want true", result.Bool())
	}
}

func TestVM_GreaterEqEqual(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpGreaterEq),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("greatereq_eq", instr, 0, int64(42), int64(42))
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
// OpPop
// ---------------------------------------------------------------------------

func TestVM_Pop(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0, // push 42
		byte(bytecode.OpPop),             // discard
		byte(bytecode.OpConst), hi0, lo0, // push 42 again
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("pop", instr, 0, int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// OpNop was removed as dead code. No VM test needed.

// ---------------------------------------------------------------------------
// OpNil
// ---------------------------------------------------------------------------

func TestVM_Nil(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpNil),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("nil", instr, 0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.IsNil() {
		t.Errorf("result = %v, want nil", result)
	}
}

// ---------------------------------------------------------------------------
// Global variable opcodes
// ---------------------------------------------------------------------------

func TestVM_SetAndGetGlobal(t *testing.T) {
	hiConst, loConst := u16(0) // const index 0
	hiG, loG := u16(0)         // global index 0
	// In non-shared mode, OpGlobal pushes a *value.Value pointer.
	// To read the actual value, we need OpDeref after OpGlobal.
	instr := makeInstructions(
		byte(bytecode.OpConst), hiConst, loConst,
		byte(bytecode.OpSetGlobal), hiG, loG,
		byte(bytecode.OpGlobal), hiG, loG,
		byte(bytecode.OpDeref),
		byte(bytecode.OpReturnVal),
	)
	fn := &bytecode.CompiledFunction{
		Name:         "global_test",
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     8,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{"global_test": fn},
		Constants: []any{int64(42)},
		Globals:   map[string]int{"x": 0},
	}
	v := New(prog)
	result, err := v.Execute("global_test", context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	iface := result.Interface()
	if got, ok := iface.(int64); !ok || got != 42 {
		t.Errorf("result = %v (type %T), want int64(42)", iface, iface)
	}
}

// ---------------------------------------------------------------------------
// OpLen and OpCap
// ---------------------------------------------------------------------------

func TestVM_LenString(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpLen),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("lenstr", instr, 0, "hello")
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 5 {
		t.Errorf("len(\"hello\") = %d, want 5", result.Int())
	}
}

func TestVM_LenSlice(t *testing.T) {
	// Push a slice constant, then OpLen
	s := []int{1, 2, 3}
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpLen),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("lenslice", instr, 0, s)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 3 {
		t.Errorf("len([]int{1,2,3}) = %d, want 3", result.Int())
	}
}

// ---------------------------------------------------------------------------
// Float arithmetic
// ---------------------------------------------------------------------------

func TestVM_AddFloats(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpAdd),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("addfloat", instr, 0, 3.14, 2.86)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	f := result.Float()
	if f < 5.99 || f > 6.01 {
		t.Errorf("3.14 + 2.86 = %f, want ~6.0", f)
	}
}

func TestVM_MulFloats(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpMul),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("mulfloat", instr, 0, 3.0, 7.0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	f := result.Float()
	if f < 20.99 || f > 21.01 {
		t.Errorf("3.0 * 7.0 = %f, want 21.0", f)
	}
}

func TestVM_NegFloat(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpNeg),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("negfloat", instr, 0, 3.14)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	f := result.Float()
	if f > -3.13 || f < -3.15 {
		t.Errorf("neg(3.14) = %f, want ~-3.14", f)
	}
}

// ---------------------------------------------------------------------------
// Complex number opcodes
// ---------------------------------------------------------------------------

func TestVM_ComplexRealImag(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	// Push real, push imag, OpComplex, OpReal, OpReturnVal
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0, // real: 3.0
		byte(bytecode.OpConst), hi1, lo1, // imag: 4.0
		byte(bytecode.OpComplex),
		byte(bytecode.OpReal),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("complex_real", instr, 0, 3.0, 4.0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	f := result.Float()
	if f < 2.99 || f > 3.01 {
		t.Errorf("real(3+4i) = %f, want 3.0", f)
	}
}

func TestVM_ComplexImag(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpComplex),
		byte(bytecode.OpImag),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("complex_imag", instr, 0, 3.0, 4.0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	f := result.Float()
	if f < 3.99 || f > 4.01 {
		t.Errorf("imag(3+4i) = %f, want 4.0", f)
	}
}

func TestVM_ComplexFromFloat32KeepsComplex64Size(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpComplex),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("complex64_size", instr, 0, float32(3.0), 4.0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.RawSize() != value.Size32 {
		t.Fatalf("complex(float32, float64) size = %v, want Size32", result.RawSize())
	}
	got := result.Interface()
	if _, ok := got.(complex64); !ok {
		t.Fatalf("complex(float32, float64) Interface() = %T, want complex64", got)
	}
}

// ---------------------------------------------------------------------------
// OpCall — function calls
// ---------------------------------------------------------------------------

func TestVM_CallFunction(t *testing.T) {
	// Define a callee that returns 42
	callee := &bytecode.CompiledFunction{
		Name: "callee",
		Instructions: []byte{
			byte(bytecode.OpConst), 0, 0,
			byte(bytecode.OpReturnVal),
		},
		NumLocals: 0,
		NumParams: 0,
		MaxStack:  1,
		FuncIdx:   1,
	}
	// Define a caller that calls callee
	hiF, loF := u16(1) // func index 1
	instr := makeInstructions(
		byte(bytecode.OpCall), byte(hiF), byte(loF), 0, // call func 1 with 0 args
		byte(bytecode.OpReturnVal),
	)
	caller := &bytecode.CompiledFunction{
		Name:         "caller",
		Instructions: instr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     8,
		FuncIdx:      0,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{
			"caller": caller,
			"callee": callee,
		},
		FuncByIndex: []*bytecode.CompiledFunction{caller, callee},
		Constants:   []any{int64(42)},
		Globals:     map[string]int{},
	}
	v := New(prog)
	result, err := v.Execute("caller", context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

func TestVM_CallFunctionWithArgs(t *testing.T) {
	// Define an add function: local 0 + local 1
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	addFn := &bytecode.CompiledFunction{
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
		FuncIdx:   1,
	}
	// Caller: push 10 and 32, call add(10, 32)
	hiConst0, loConst0 := u16(0)
	hiConst1, loConst1 := u16(1)
	hiF, loF := u16(1)
	callerInstr := makeInstructions(
		byte(bytecode.OpConst), byte(hiConst0), byte(loConst0),
		byte(bytecode.OpConst), byte(hiConst1), byte(loConst1),
		byte(bytecode.OpCall), byte(hiF), byte(loF), 2, // call func 1 with 2 args
		byte(bytecode.OpReturnVal),
	)
	caller := &bytecode.CompiledFunction{
		Name:         "caller",
		Instructions: callerInstr,
		NumLocals:    0,
		NumParams:    0,
		MaxStack:     8,
		FuncIdx:      0,
	}
	prog := &bytecode.CompiledProgram{
		Functions: map[string]*bytecode.CompiledFunction{
			"caller": caller,
			"add":    addFn,
		},
		FuncByIndex: []*bytecode.CompiledFunction{caller, addFn},
		Constants:   []any{int64(10), int64(32)},
		Globals:     map[string]int{},
	}
	v := New(prog)
	result, err := v.Execute("caller", context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}
}

// ---------------------------------------------------------------------------
// SharedGlobals
// ---------------------------------------------------------------------------

func TestSharedGlobalsGetSet(t *testing.T) {
	sg := NewSharedGlobals([]value.Value{value.MakeInt(0)}, 1)

	// Get initial value
	got := sg.Get(0)
	if got.Int() != 0 {
		t.Errorf("Get(0) = %d, want 0", got.Int())
	}

	// Set and verify
	sg.Set(0, value.MakeInt(42))
	got = sg.Get(0)
	if got.Int() != 42 {
		t.Errorf("Get(0) after Set = %d, want 42", got.Int())
	}
}

func TestSharedGlobalsLen(t *testing.T) {
	sg := NewSharedGlobals(nil, 5)
	if sg.Len() != 5 {
		t.Errorf("Len() = %d, want 5", sg.Len())
	}
}

func TestSharedGlobalsGlobals(t *testing.T) {
	initial := []value.Value{value.MakeInt(1), value.MakeString("hello")}
	sg := NewSharedGlobals(initial, 2)
	globals := sg.Globals()
	if len(globals) != 2 {
		t.Fatalf("len(Globals()) = %d, want 2", len(globals))
	}
	if globals[0].Int() != 1 {
		t.Errorf("globals[0] = %d, want 1", globals[0].Int())
	}
	if globals[1].String() != "hello" {
		t.Errorf("globals[1] = %q, want %q", globals[1].String(), "hello")
	}
}

func TestSharedGlobalsInitExternalVars(t *testing.T) {
	sg := NewSharedGlobals(nil, 3)
	// Simulate external variable at index 1 — stored as a pointer
	extVal := 42
	sg.InitExternalVars(map[int]any{1: &extVal})
	got := sg.Get(1)
	// The value is a *int pointer wrapped via FromInterface
	iface := got.Interface()
	if ptr, ok := iface.(*int); !ok || *ptr != 42 {
		t.Errorf("Get(1) after InitExternalVars = %v (type %T), want *int(42)", iface, iface)
	}
}

// ---------------------------------------------------------------------------
// GlobalRef
// ---------------------------------------------------------------------------

func TestGlobalRefLoadStore(t *testing.T) {
	sg := NewSharedGlobals([]value.Value{value.MakeInt(0)}, 1)
	ref := &GlobalRef{sg: sg, idx: 0}

	// Load initial
	got := ref.Load()
	if got.Int() != 0 {
		t.Errorf("Load() = %d, want 0", got.Int())
	}

	// Store and load
	ref.Store(value.MakeInt(99))
	got = ref.Load()
	if got.Int() != 99 {
		t.Errorf("Load() after Store = %d, want 99", got.Int())
	}
}

// ---------------------------------------------------------------------------
// Bind/Unbind SharedGlobals
// ---------------------------------------------------------------------------

func TestVM_BindUnbindSharedGlobals(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("shared", instr, 0, int64(42))
	v := New(prog)

	sg := NewSharedGlobals(nil, 0)
	v.BindSharedGlobals(sg)

	// Should still execute fine
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Int() != 42 {
		t.Errorf("result = %d, want 42", result.Int())
	}

	v.UnbindSharedGlobals()
}

// ---------------------------------------------------------------------------
// OpReturn (no value)
// ---------------------------------------------------------------------------

func TestVM_ReturnNoValue(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpConst), 0, 0,
		byte(bytecode.OpPop),
		byte(bytecode.OpReturn),
	)
	prog, name := buildProg("ret_none", instr, 0, int64(42))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	// OpReturn with no value should still complete
	_ = result
}

// ---------------------------------------------------------------------------
// OpPanic
// ---------------------------------------------------------------------------

func TestVM_PanicOpcode(t *testing.T) {
	hi0, lo0 := u16(0)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0, // push "error message"
		byte(bytecode.OpPanic),
	)
	prog, name := buildProg("panic_test", instr, 0, "error message")
	v := New(prog)
	_, err := v.Execute(name, context.Background())
	if err == nil {
		t.Error("expected error from OpPanic")
	}
}

// ---------------------------------------------------------------------------
// OpPack / OpUnpack
// ---------------------------------------------------------------------------

func TestVM_PackUnpack(t *testing.T) {
	// Push 2 values, pack them into a tuple, unpack, and return first
	hi0, lo0 := u16(0)       // const 0 = 10
	hi1, lo1 := u16(1)       // const 1 = 32
	hiPack, loPack := u16(2) // count = 2
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpPack), byte(hiPack), byte(loPack),
		byte(bytecode.OpUnpack),
		byte(bytecode.OpReturnVal), // top of stack should be 32 (last unpacked)
	)
	prog, name := buildProg("pack_unpack", instr, 0, int64(10), int64(32))
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	// After unpack, the values are pushed back in order, so top = last value
	if result.Int() != 32 {
		t.Errorf("result = %d, want 32", result.Int())
	}
}

// ---------------------------------------------------------------------------
// OpBool comparisons
// ---------------------------------------------------------------------------

func TestVM_EqualStrings(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("eq_str", instr, 0, "hello", "hello")
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("equal strings result = %v, want true", result.Bool())
	}
}

func TestVM_EqualBools(t *testing.T) {
	instr := makeInstructions(
		byte(bytecode.OpTrue),
		byte(bytecode.OpTrue),
		byte(bytecode.OpEqual),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("eq_bool", instr, 0)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("true == true = %v, want true", result.Bool())
	}
}

// ---------------------------------------------------------------------------
// OpLess / OpGreater with floats
// ---------------------------------------------------------------------------

func TestVM_LessFloats(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpLess),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("less_float", instr, 0, 1.5, 2.5)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("1.5 < 2.5 = %v, want true", result.Bool())
	}
}

func TestVM_GreaterFloats(t *testing.T) {
	hi0, lo0 := u16(0)
	hi1, lo1 := u16(1)
	instr := makeInstructions(
		byte(bytecode.OpConst), hi0, lo0,
		byte(bytecode.OpConst), hi1, lo1,
		byte(bytecode.OpGreater),
		byte(bytecode.OpReturnVal),
	)
	prog, name := buildProg("greater_float", instr, 0, 2.5, 1.5)
	v := New(prog)
	result, err := v.Execute(name, context.Background())
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if !result.Bool() {
		t.Errorf("2.5 > 1.5 = %v, want true", result.Bool())
	}
}
