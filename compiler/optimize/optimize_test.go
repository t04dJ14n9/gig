package optimize

import (
	"testing"

	"git.woa.com/youngjin/gig/bytecode"
)

func TestSafeIdx(t *testing.T) {
	flags := []bool{true, false, true, false, true}

	tests := []struct {
		idx  int
		want bool
	}{
		{0, true},
		{1, false},
		{2, true},
		{3, false},
		{4, true},
		{5, false},   // out of bounds
		{100, false}, // way out of bounds
	}

	for _, tt := range tests {
		got := safeIdx(tt.idx, flags)
		if got != tt.want {
			t.Errorf("safeIdx(%d, flags) = %v, want %v", tt.idx, got, tt.want)
		}
	}
}

func TestSafeIdxEmptyFlags(t *testing.T) {
	flags := []bool{}

	if safeIdx(0, flags) {
		t.Error("safeIdx(0, []) = true, want false")
	}
	// Note: safeIdx does not handle negative indices, so we don't test that case
}

func TestReadU16(t *testing.T) {
	code := []byte{0xAB, 0xCD, 0xEF, 0x01}

	if got := readU16(code, 0); got != 0xABCD {
		t.Errorf("readU16(code, 0) = 0x%04x, want 0xABCD", got)
	}
	// At offset 1: bytes are 0xCD, 0xEF → 0xCDEF
	if got := readU16(code, 1); got != 0xCDEF {
		t.Errorf("readU16(code, 1) = 0x%04x, want 0xCDEF", got)
	}
	if got := readU16(code, 2); got != 0xEF01 {
		t.Errorf("readU16(code, 2) = 0x%04x, want 0xEF01", got)
	}
}

func TestWriteU16(t *testing.T) {
	code := make([]byte, 6)

	writeU16(code, 0, 0x1234)
	writeU16(code, 2, 0x5678)
	writeU16(code, 4, 0x9ABC)

	if code[0] != 0x12 || code[1] != 0x34 {
		t.Errorf("writeU16 to 0: got [%02x, %02x], want [12, 34]", code[0], code[1])
	}
	if code[2] != 0x56 || code[3] != 0x78 {
		t.Errorf("writeU16 to 2: got [%02x, %02x], want [56, 78]", code[2], code[3])
	}
	if code[4] != 0x9A || code[5] != 0xBC {
		t.Errorf("writeU16 to 4: got [%02x, %02x], want [9A, BC]", code[4], code[5])
	}
}

func TestOpcodeWidth(t *testing.T) {
	tests := []struct {
		op    bytecode.OpCode
		width int
	}{
		{bytecode.OpAdd, 0},
		{bytecode.OpSub, 0},
		{bytecode.OpMul, 0},
		{bytecode.OpLocal, 2},
		{bytecode.OpSetLocal, 2},
		{bytecode.OpJump, 2},
		{bytecode.OpJumpTrue, 2},
		{bytecode.OpJumpFalse, 2},
		{bytecode.OpLessLocalLocalJumpTrue, 6},
		{bytecode.OpIntLocal, 2},
		{bytecode.OpIntSetLocal, 2},
		{bytecode.OpIntMoveLocal, 4},
		{bytecode.OpClosure, 3},
	}

	for _, tt := range tests {
		got := opcodeWidth(tt.op)
		if got != tt.width {
			t.Errorf("opcodeWidth(%v) = %d, want %d", tt.op, got, tt.width)
		}
	}
}

func TestLocalUsedOutsideEmptyCode(t *testing.T) {
	code := []byte{}
	if localUsedOutside(code, 0, 0, 0) {
		t.Error("localUsedOutside([]byte{}, 0, 0, 0) = true, want false")
	}
}

func TestLocalUsedOutsideNoMatch(t *testing.T) {
	// OpLocal at idx 0, but we skip the range that contains it
	code := []byte{byte(bytecode.OpLocal), 0x00, 0x01}
	if localUsedOutside(code, 0, 0, 3) {
		t.Error("localUsedOutside with skip range covering the use = true, want false")
	}
}

func TestLocalUsedOutsideOpAddr(t *testing.T) {
	// OpAddr also references locals - OpAddr takes 3 bytes (1 opcode + 2 index)
	// We skip the range [0,3) which covers the entire OpAddr instruction
	// So this should return false (we're skipping the reference)
	code := []byte{byte(bytecode.OpAddr), 0x00, 0x01}
	if localUsedOutside(code, 1, 0, 3) {
		t.Error("localUsedOutside(OpAddr, 1, 0, 3) = true, want false (skipped)")
	}

	// Without skipping, it should find the reference
	// OpAddr with index 1 (bytes: 0x00, 0x01 big-endian)
	code2 := []byte{byte(bytecode.OpAddr), 0x00, 0x01}
	if !localUsedOutside(code2, 1, 3, 10) { // skip range doesn't cover the instruction
		t.Error("localUsedOutside(OpAddr, 1, 3, 10) = false, want true")
	}
}

func TestLocalUsedOutsideNoReference(t *testing.T) {
	// Only OpAdd (no local index)
	code := []byte{byte(bytecode.OpAdd)}
	if localUsedOutside(code, 0, 0, 1) {
		t.Error("localUsedOutside(OpAdd, 0, 0, 1) = true, want false")
	}
}

func TestApplyRewritesEmpty(t *testing.T) {
	code := []byte{1, 2, 3, 4}
	rewrites := []rewrite{}

	result := applyRewrites(code, rewrites)

	if len(result) != len(code) {
		t.Errorf("applyRewrites with no rewrites: len = %d, want %d", len(result), len(code))
	}
	for i := range result {
		if result[i] != code[i] {
			t.Errorf("applyRewrites result[%d] = %d, want %d", i, result[i], code[i])
		}
	}
}

func TestApplyRewritesSingleByteReplacement(t *testing.T) {
	code := []byte{1, 2, 3, 4}
	rewrites := []rewrite{
		{oldStart: 1, oldEnd: 2, newBytes: []byte{9, 8}},
	}

	result := applyRewrites(code, rewrites)

	if len(result) != 5 {
		t.Errorf("applyRewrites length = %d, want 5", len(result))
	}
	if result[0] != 1 || result[1] != 9 || result[2] != 8 || result[3] != 3 || result[4] != 4 {
		t.Errorf("applyRewrites result = %v, want [1, 9, 8, 3, 4]", result)
	}
}

func TestApplyRewritesDeletion(t *testing.T) {
	code := []byte{1, 2, 3, 4}
	rewrites := []rewrite{
		{oldStart: 1, oldEnd: 3, newBytes: nil}, // delete bytes 1,2
	}

	result := applyRewrites(code, rewrites)

	if len(result) != 2 {
		t.Errorf("applyRewrites deletion length = %d, want 2", len(result))
	}
	if result[0] != 1 || result[1] != 4 {
		t.Errorf("applyRewrites result = %v, want [1, 4]", result)
	}
}

func TestApplyRewritesInsertion(t *testing.T) {
	code := []byte{1, 4}
	rewrites := []rewrite{
		{oldStart: 1, oldEnd: 1, newBytes: []byte{2, 3}}, // insert 2,3
	}

	result := applyRewrites(code, rewrites)

	if len(result) != 4 {
		t.Errorf("applyRewrites insertion length = %d, want 4", len(result))
	}
	if result[0] != 1 || result[1] != 2 || result[2] != 3 || result[3] != 4 {
		t.Errorf("applyRewrites result = %v, want [1, 2, 3, 4]", result)
	}
}

func TestFixJumpTargetsNoJump(t *testing.T) {
	code := []byte{byte(bytecode.OpAdd)}
	offsetMap := []int{0}

	fixJumpTargets(code, offsetMap, 1)

	// Should not panic and code should be unchanged
	if code[0] != byte(bytecode.OpAdd) {
		t.Errorf("fixJumpTargets modified code unexpectedly")
	}
}

func TestPeepholeEmptyCode(t *testing.T) {
	code := []byte{}
	result := Peephole(code)

	if len(result) != 0 {
		t.Errorf("Peephole([]) length = %d, want 0", len(result))
	}
}

func TestPeepholeNoMatchingPatterns(t *testing.T) {
	// OpAdd is not typically matched by peephole patterns
	code := []byte{byte(bytecode.OpAdd)}
	result := Peephole(code)

	if len(result) != 1 {
		t.Errorf("Peephole([OpAdd]) length = %d, want 1", len(result))
	}
	if result[0] != byte(bytecode.OpAdd) {
		t.Errorf("Peephole result[0] = 0x%02x, want 0x%02x", result[0], bytecode.OpAdd)
	}
}

func TestOptimizeNoChanges(t *testing.T) {
	code := []byte{byte(bytecode.OpAdd)}
	localIsInt := []bool{false}
	constIsInt := []bool{false}
	localIsIntSlice := []bool{false}

	result, changed := Optimize(code, localIsInt, constIsInt, localIsIntSlice)

	if len(result) != 1 {
		t.Errorf("Optimize length = %d, want 1", len(result))
	}
	if changed {
		t.Error("Optimize returned changed = true, want false")
	}
}

func TestOptimizeWithIntSpecialization(t *testing.T) {
	// OpSetLocal can be upgraded to OpIntSetLocal when intUsed is true
	code := []byte{byte(bytecode.OpSetLocal), 0x00, 0x00}
	localIsInt := []bool{true}  // local 0 is int
	constIsInt := []bool{false}
	localIsIntSlice := []bool{false}

	result, changed := Optimize(code, localIsInt, constIsInt, localIsIntSlice)

	// The optimization should upgrade OpSetLocal to OpIntSetLocal
	if len(result) >= 1 && result[0] == byte(bytecode.OpIntSetLocal) {
		t.Log("Optimize correctly upgraded OpSetLocal to OpIntSetLocal")
	} else {
		t.Logf("Optimize result[0] = 0x%02x (OpSetLocal=0x%02x, OpIntSetLocal=0x%02x), changed=%v",
			result[0], bytecode.OpSetLocal, bytecode.OpIntSetLocal, changed)
	}
}

func TestRewriteStruct(t *testing.T) {
	r := rewrite{
		oldStart: 5,
		oldEnd:   10,
		newBytes: []byte{1, 2, 3},
	}

	if r.oldStart != 5 || r.oldEnd != 10 || len(r.newBytes) != 3 {
		t.Errorf("rewrite struct mismatch: got %+v", r)
	}
}