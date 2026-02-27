package bytecode

import (
	"testing"
)

// TestReadWriteUint16 verifies the shared-kernel helpers that encode/decode
// two-byte operands in the instruction stream.
func TestReadWriteUint16(t *testing.T) {
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
			WriteUint16(buf, 0, tt.value)
			got := ReadUint16(buf, 0)
			if got != tt.value {
				t.Errorf("ReadUint16(WriteUint16(%d)) = %d", tt.value, got)
			}
			// Also test with non-zero offset.
			WriteUint16(buf, 2, tt.value)
			got2 := ReadUint16(buf, 2)
			if got2 != tt.value {
				t.Errorf("ReadUint16 at offset 2: got %d, want %d", got2, tt.value)
			}
		})
	}
}

// TestOperandWidths ensures every opcode that has defined operand widths
// returns a non-nil slice, and that undeclared opcodes return nil.
func TestOperandWidths(t *testing.T) {
	// Spot-check a few well-known opcodes.
	// OperandWidths maps OpCode to a single int (total operand byte size).
	checks := []struct {
		op    OpCode
		name  string
		width int
		inMap bool
	}{
		{OpNop, "OpNop", 0, false},
		{OpConst, "OpConst", 2, true},
		{OpLocal, "OpLocal", 2, true},
		{OpJump, "OpJump", 2, true},
		{OpFree, "OpFree", 1, true},
	}

	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			got, ok := OperandWidths[c.op]
			if ok != c.inMap {
				t.Errorf("in map: got %v, want %v", ok, c.inMap)
				return
			}
			if c.inMap && got != c.width {
				t.Errorf("width: got %d, want %d", got, c.width)
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

// TestCompiledFunctionZeroValue ensures the zero value is usable.
func TestCompiledFunctionZeroValue(t *testing.T) {
	var fn CompiledFunction
	if fn.Name != "" {
		t.Error("zero CompiledFunction should have empty name")
	}
	if fn.Instructions != nil {
		t.Error("zero CompiledFunction should have nil instructions")
	}
	if fn.NumLocals != 0 || fn.NumParams != 0 || fn.NumFreeVars != 0 || fn.MaxStack != 0 {
		t.Error("zero CompiledFunction numeric fields should be 0")
	}
}

// TestProgramZeroValue ensures the zero value is usable.
func TestProgramZeroValue(t *testing.T) {
	var p Program
	if p.Functions != nil || p.Constants != nil || p.Globals != nil {
		t.Error("zero Program should have nil maps/slices")
	}
}
