package bytecode

import (
	"go/types"
	"reflect"
	"testing"
)

// TestReadWriteU16 verifies the shared-kernel helpers that encode/decode
// two-byte operands in the instruction stream.
func TestReadWriteU16(t *testing.T) {
	tests := []struct {
		name  string
		value uint16
	}{
		{"zero", 0},
		{"one", 1},
		{"max_byte", 255},
		{"two_bytes", 256},
		{"large", 0xABCD},
		{"max", 0xFFFF},
	}

	buf := make([]byte, 4)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteU16(buf, 0, tt.value)
			got := ReadU16(buf, 0)
			if got != tt.value {
				t.Errorf("ReadU16(WriteU16(%d)) = %d", tt.value, got)
			}
			// Also test with non-zero offset.
			WriteU16(buf, 2, tt.value)
			got2 := ReadU16(buf, 2)
			if got2 != tt.value {
				t.Errorf("ReadU16 at offset 2: got %d, want %d", got2, tt.value)
			}
		})
	}
}

// TestOperandWidth ensures OperandWidth returns correct values for known opcodes.
func TestOperandWidth(t *testing.T) {
	// Spot-check a few well-known opcodes.
	checks := []struct {
		op    OpCode
		name  string
		width int
	}{
		{OpNop, "OpNop", 0},
		{OpConst, "OpConst", 2},
		{OpLocal, "OpLocal", 2},
		{OpJump, "OpJump", 2},
		{OpFree, "OpFree", 1},
		{OpCall, "OpCall", 3},
		{OpChangeType, "OpChangeType", 4},
		{OpCallIndirect, "OpCallIndirect", 1},
		{OpAdd, "OpAdd", 0},
		{OpSub, "OpSub", 0},
		{OpMul, "OpMul", 0},
		{OpDiv, "OpDiv", 0},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			got := OperandWidth(c.op)
			if got != c.width {
				t.Errorf("OperandWidth(%s) = %d, want %d", c.name, got, c.width)
			}
		})
	}
}

// TestOpcodeConstants verifies that the opcode constants used by the compiler
// and VM agree on their numeric values. This prevents silent breakage if
// someone reorders the const block.
func TestOpcodeConstants(t *testing.T) {
	if OpNop != 0 {
		t.Errorf("OpNop = %d, want 0", OpNop)
	}
	if OpPop != 1 {
		t.Errorf("OpPop = %d, want 1", OpPop)
	}
	// OpHalt should be the last one; verify it is > 0.
	if OpHalt == 0 {
		t.Error("OpHalt should be non-zero")
	}
}

// TestOpCodeString tests opcode String() methods.
func TestOpCodeString(t *testing.T) {
	tests := []struct {
		op       OpCode
		expected string
	}{
		{OpNop, "NOP"},
		{OpPop, "POP"},
		{OpDup, "DUP"},
		{OpConst, "CONST"},
		{OpNil, "NIL"},
		{OpTrue, "TRUE"},
		{OpFalse, "FALSE"},
		{OpLocal, "LOCAL"},
		{OpSetLocal, "SETLOCAL"},
		{OpGlobal, "GLOBAL"},
		{OpSetGlobal, "SETGLOBAL"},
		{OpFree, "FREE"},
		{OpSetFree, "SETFREE"},
		{OpAdd, "ADD"},
		{OpSub, "SUB"},
		{OpMul, "MUL"},
		{OpDiv, "DIV"},
		{OpMod, "MOD"},
		{OpNeg, "NEG"},
		{OpAnd, "AND"},
		{OpOr, "OR"},
		{OpXor, "XOR"},
		{OpAndNot, "ANDNOT"},
		{OpLsh, "LSH"},
		{OpRsh, "RSH"},
		{OpEqual, "EQUAL"},
		{OpNotEqual, "NOTEQUAL"},
		{OpLess, "LESS"},
		{OpLessEq, "LESSEQ"},
		{OpGreater, "GREATER"},
		{OpGreaterEq, "GREATEREQ"},
		{OpNot, "NOT"},
		{OpJump, "JUMP"},
		{OpJumpTrue, "JUMPTRUE"},
		{OpJumpFalse, "JUMPFALSE"},
		{OpCall, "CALL"},
		{OpReturn, "RETURN"},
		{OpReturnVal, "RETURNVAL"},
		{OpMakeSlice, "MAKESLICE"},
		{OpMakeMap, "MAKEMAP"},
		{OpMakeChan, "MAKECHAN"},
		{OpMakeArray, "MAKEARRAY"},
		{OpMakeStruct, "MAKESTRUCT"},
		{OpIndex, "INDEX"},
		{OpIndexOk, "INDEXOK"},
		{OpSetIndex, "SETINDEX"},
		{OpSlice, "SLICE"},
		{OpMapIter, "MAPITER"},
		{OpMapIterNext, "MAPITERNEXT"},
		{OpField, "FIELD"},
		{OpSetField, "SETFIELD"},
		{OpAddr, "ADDR"},
		{OpFieldAddr, "FIELDADDR"},
		{OpIndexAddr, "INDEXADDR"},
		{OpDeref, "DEREF"},
		{OpSetDeref, "SETDEREF"},
		{OpAssert, "ASSERT"},
		{OpConvert, "CONVERT"},
		{OpChangeType, "CHANGETYPE"},
		{OpClosure, "CLOSURE"},
		{OpMethod, "METHOD"},
		{OpMethodCall, "METHODCALL"},
		{OpGoCall, "GOCALL"},
		{OpGoCallIndirect, "GOCALLINDIRECT"},
		{OpSend, "SEND"},
		{OpRecv, "RECV"},
		{OpRecvOk, "RECVOK"},
		{OpTrySend, "TRYSEND"},
		{OpTryRecv, "TRYRECV"},
		{OpClose, "CLOSE"},
		{OpSelect, "SELECT"},
		{OpDefer, "DEFER"},
		{OpDeferIndirect, "DEFERINDIRECT"},
		{OpRunDefers, "RUNDEFERS"},
		{OpRecover, "RECOVER"},
		{OpRange, "RANGE"},
		{OpRangeNext, "RANGENEXT"},
		{OpLen, "LEN"},
		{OpCap, "CAP"},
		{OpAppend, "APPEND"},
		{OpCopy, "COPY"},
		{OpDelete, "DELETE"},
		{OpPanic, "PANIC"},
		{OpPrint, "PRINT"},
		{OpPrintln, "PRINTLN"},
		{OpNew, "NEW"},
		{OpMake, "MAKE"},
		{OpCallExternal, "CALLEXTERNAL"},
		{OpCallIndirect, "CALLINDIRECT"},
		{OpPack, "PACK"},
		{OpUnpack, "UNPACK"},
		{OpHalt, "HALT"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.op.String()
			if got != tt.expected {
				t.Errorf("OpCode(%d).String() = %q, want %q", tt.op, got, tt.expected)
			}
		})
	}
}

// TestOpCodeStringUnknown tests that unknown opcodes return "UNKNOWN".
func TestOpCodeStringUnknown(t *testing.T) {
	unknown := OpCode(255)
	if got := unknown.String(); got != "UNKNOWN" {
		t.Errorf("Unknown opcode String() = %q, want %q", got, "UNKNOWN")
	}
}

// TestProgramReflectTypeCache tests the ReflectTypeCache methods.
func TestProgramReflectTypeCache(t *testing.T) {
	p := &Program{}

	// Test caching with actual types.Type
	key := types.Typ[types.Int] // int type
	rt := reflect.TypeOf(42)

	// Cache it
	result := p.CacheReflectType(key, rt)
	if result != rt {
		t.Error("CacheReflectType should return the stored type")
	}

	// Retrieve it
	got, ok := p.CachedReflectType(key)
	if !ok {
		t.Error("CachedReflectType should find the cached type")
	}
	if got != rt {
		t.Error("CachedReflectType should return the same type")
	}

	// Test non-existent key (use a different type)
	_, ok = p.CachedReflectType(types.Typ[types.String])
	if ok {
		t.Error("CachedReflectType should return false for non-existent key")
	}

	// Test LoadOrStore behavior - same key should return existing type
	rt2 := reflect.TypeOf(int64(0))
	result2 := p.CacheReflectType(key, rt2)
	if result2 != rt {
		t.Error("CacheReflectType should return existing type on duplicate key")
	}
}

// TestCompiledFunction tests CompiledFunction with actual values.
func TestCompiledFunction(t *testing.T) {
	fn := &CompiledFunction{
		Name:         "testFunc",
		Instructions: []byte{byte(OpNop), byte(OpConst), 0, 1, byte(OpReturnVal)},
		NumLocals:    5,
		NumParams:    2,
		NumFreeVars:  1,
		MaxStack:     20,
		HasIntLocals: true,
	}

	if fn.Name != "testFunc" {
		t.Errorf("Expected name 'testFunc', got %q", fn.Name)
	}
	if len(fn.Instructions) != 5 {
		t.Errorf("Expected 5 instruction bytes, got %d", len(fn.Instructions))
	}
	if fn.NumLocals != 5 {
		t.Errorf("Expected NumLocals 5, got %d", fn.NumLocals)
	}
	if fn.NumParams != 2 {
		t.Errorf("Expected NumParams 2, got %d", fn.NumParams)
	}
	if fn.NumFreeVars != 1 {
		t.Errorf("Expected NumFreeVars 1, got %d", fn.NumFreeVars)
	}
	if fn.MaxStack != 20 {
		t.Errorf("Expected MaxStack 20, got %d", fn.MaxStack)
	}
	if !fn.HasIntLocals {
		t.Error("Expected HasIntLocals to be true")
	}
}

// TestProgramWithValues tests Program with actual values.
func TestProgramWithValues(t *testing.T) {
	p := &Program{
		Functions: map[string]*CompiledFunction{
			"main": {
				Name:         "main",
				Instructions: []byte{byte(OpNop), byte(OpReturn)},
				NumLocals:    2,
				NumParams:    0,
				MaxStack:     10,
			},
			"add": {
				Name:         "add",
				Instructions: []byte{byte(OpLocal), 0, 0, byte(OpLocal), 0, 1, byte(OpAdd), byte(OpReturnVal)},
				NumLocals:    2,
				NumParams:    2,
				MaxStack:     5,
			},
		},
		FuncByIndex: []*CompiledFunction{
			{Name: "main"},
			{Name: "add"},
		},
		Constants: []any{int64(1), int64(2), "hello"},
		Globals:   map[string]int{"x": 0, "y": 1},
	}

	// Test Functions map
	if len(p.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(p.Functions))
	}

	if fn, ok := p.Functions["main"]; !ok {
		t.Error("Expected 'main' function to exist")
	} else if fn.Name != "main" {
		t.Errorf("Expected function name 'main', got %q", fn.Name)
	}

	// Test FuncByIndex
	if len(p.FuncByIndex) != 2 {
		t.Errorf("Expected 2 FuncByIndex entries, got %d", len(p.FuncByIndex))
	}

	// Test Constants
	if len(p.Constants) != 3 {
		t.Errorf("Expected 3 constants, got %d", len(p.Constants))
	}

	// Test Globals
	if len(p.Globals) != 2 {
		t.Errorf("Expected 2 globals, got %d", len(p.Globals))
	}
}

// TestReadWriteU16ByteOrder verifies big-endian byte order.
func TestReadWriteU16ByteOrder(t *testing.T) {
	buf := make([]byte, 2)
	WriteU16(buf, 0, 0x1234)

	// Big-endian: high byte first
	if buf[0] != 0x12 {
		t.Errorf("Expected high byte 0x12, got 0x%02X", buf[0])
	}
	if buf[1] != 0x34 {
		t.Errorf("Expected low byte 0x34, got 0x%02X", buf[1])
	}

	// Verify reading
	buf2 := []byte{0xAB, 0xCD}
	got := ReadU16(buf2, 0)
	if got != 0xABCD {
		t.Errorf("Expected 0xABCD, got 0x%04X", got)
	}
}

// TestOperandWidthSuperinstructions tests superinstruction operand widths.
func TestOperandWidthSuperinstructions(t *testing.T) {
	tests := []struct {
		op    OpCode
		name  string
		width int
	}{
		{OpAddLocalLocal, "OpAddLocalLocal", 4},
		{OpSubLocalLocal, "OpSubLocalLocal", 4},
		{OpMulLocalLocal, "OpMulLocalLocal", 4},
		{OpAddLocalConst, "OpAddLocalConst", 4},
		{OpSubLocalConst, "OpSubLocalConst", 4},
		{OpAddSetLocal, "OpAddSetLocal", 2},
		{OpSubSetLocal, "OpSubSetLocal", 2},
		{OpLessLocalLocalJumpTrue, "OpLessLocalLocalJumpTrue", 6},
		{OpLessLocalConstJumpTrue, "OpLessLocalConstJumpTrue", 6},
		{OpLocalLocalAddSetLocal, "OpLocalLocalAddSetLocal", 6},
		{OpLocalConstAddSetLocal, "OpLocalConstAddSetLocal", 6},
		{OpIntSetLocal, "OpIntSetLocal", 2},
		{OpIntLocal, "OpIntLocal", 2},
		{OpIntMoveLocal, "OpIntMoveLocal", 4},
		{OpIntSliceGet, "OpIntSliceGet", 6},
		{OpIntSliceSet, "OpIntSliceSet", 6},
		{OpIntSliceSetConst, "OpIntSliceSetConst", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OperandWidth(tt.op)
			if got != tt.width {
				t.Errorf("OperandWidth(%s) = %d, want %d", tt.name, got, tt.width)
			}
		})
	}
}
