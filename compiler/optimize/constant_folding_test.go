package optimize

import (
	"testing"

	"github.com/t04dJ14n9/gig/model/bytecode"
)

func TestFoldConstantsCollapsesIntegerArithmetic(t *testing.T) {
	constants := []any{int64(3), int64(4)}
	code := encodeOps(
		op(bytecode.OpConst, 0),
		op(bytecode.OpConst, 1),
		op(bytecode.OpMul),
	)

	got := FoldConstants(code, &constants)

	if len(got) != 3 {
		t.Fatalf("folded code len = %d, want 3 bytes: %v", len(got), got)
	}
	if bytecode.OpCode(got[0]) != bytecode.OpConst {
		t.Fatalf("folded op = %v, want OpConst", bytecode.OpCode(got[0]))
	}
	foldedIdx := bytecode.ReadU16(got, 1)
	if constants[foldedIdx] != int64(12) {
		t.Fatalf("folded constant = %v, want int64(12)", constants[foldedIdx])
	}
}

func TestFoldConstantsPropagatesLocalConstantsInsideBlock(t *testing.T) {
	constants := []any{int64(3), int64(4), int64(5)}
	code := encodeOps(
		op(bytecode.OpConst, 0),
		op(bytecode.OpSetLocal, 0),
		op(bytecode.OpConst, 1),
		op(bytecode.OpSetLocal, 1),
		op(bytecode.OpLocal, 0),
		op(bytecode.OpLocal, 1),
		op(bytecode.OpMul),
		op(bytecode.OpConst, 2),
		op(bytecode.OpAdd),
		op(bytecode.OpReturnVal),
	)

	got := FoldConstants(code, &constants)

	counts := countOps(got)
	if counts[bytecode.OpMul] != 0 || counts[bytecode.OpAdd] != 0 {
		t.Fatalf("folded code still has arithmetic ops: counts=%v code=%v", counts, got)
	}
	if constants[constantOperandBeforeReturn(t, got)] != int64(17) {
		t.Fatalf("constant before return = %v, want int64(17)", constants[constantOperandBeforeReturn(t, got)])
	}
}

func TestFoldConstantsSkipsDivideByZero(t *testing.T) {
	constants := []any{int64(12), int64(0)}
	code := encodeOps(
		op(bytecode.OpConst, 0),
		op(bytecode.OpConst, 1),
		op(bytecode.OpDiv),
	)

	got := FoldConstants(code, &constants)

	if counts := countOps(got); counts[bytecode.OpDiv] != 1 {
		t.Fatalf("OpDiv count = %d, want 1 to preserve runtime divide-by-zero behavior", counts[bytecode.OpDiv])
	}
}

func TestFoldConstantsDoesNotFoldAcrossJumpTarget(t *testing.T) {
	constants := []any{int64(3), int64(4)}
	code := encodeOps(
		op(bytecode.OpJump, 6),
		op(bytecode.OpConst, 0),
		op(bytecode.OpConst, 1),
		op(bytecode.OpAdd),
	)

	got := FoldConstants(code, &constants)

	if counts := countOps(got); counts[bytecode.OpAdd] != 1 {
		t.Fatalf("OpAdd count = %d, want 1 because right operand is a jump target", counts[bytecode.OpAdd])
	}
}

func TestFoldConstantsTurnsKnownTrueBranchIntoJumpAndRemovesDeadFallthrough(t *testing.T) {
	constants := []any{true, int64(1)}
	code := encodeOps(
		op(bytecode.OpConst, 0),
		op(bytecode.OpJumpTrue, 9),
		op(bytecode.OpConst, 1),
	)

	got := FoldConstants(code, &constants)

	if len(got) != 3 {
		t.Fatalf("folded code len = %d, want only jump to end: %v", len(got), got)
	}
	if bytecode.OpCode(got[0]) != bytecode.OpJump {
		t.Fatalf("folded branch op = %v, want OpJump", bytecode.OpCode(got[0]))
	}
	if target := bytecode.ReadU16(got, 1); target != uint16(len(got)) {
		t.Fatalf("folded jump target = %d, want new end offset %d", target, len(got))
	}
	if counts := countOps(got); counts[bytecode.OpJumpTrue] != 0 {
		t.Fatalf("OpJumpTrue count = %d, want 0 after constant branch folding", counts[bytecode.OpJumpTrue])
	}
}

func TestFoldConstantsDeletesKnownFalseJumpTrue(t *testing.T) {
	constants := []any{false, int64(1)}
	code := encodeOps(
		op(bytecode.OpConst, 0),
		op(bytecode.OpJumpTrue, 9),
		op(bytecode.OpConst, 1),
	)

	got := FoldConstants(code, &constants)

	if len(got) != 3 {
		t.Fatalf("folded code len = %d, want only trailing OpConst: %v", len(got), got)
	}
	if bytecode.OpCode(got[0]) != bytecode.OpConst || bytecode.ReadU16(got, 1) != 1 {
		t.Fatalf("folded code = %v, want OpConst 1", got)
	}
	if counts := countOps(got); counts[bytecode.OpJumpTrue] != 0 {
		t.Fatalf("OpJumpTrue count = %d, want 0 after deleting untaken branch", counts[bytecode.OpJumpTrue])
	}
}

func TestFoldConstantsKeepsBareConstantBool(t *testing.T) {
	constants := []any{}
	code := encodeOps(op(bytecode.OpTrue))

	got := FoldConstants(code, &constants)

	if len(got) != len(code) || bytecode.OpCode(got[0]) != bytecode.OpTrue {
		t.Fatalf("FoldConstants(OpTrue) = %v, want unchanged OpTrue", got)
	}
}

func TestFoldConstantsRemovesUnreachableAfterJump(t *testing.T) {
	constants := []any{int64(1), int64(2)}
	code := encodeOps(
		op(bytecode.OpJump, 6),
		op(bytecode.OpConst, 0),
		op(bytecode.OpConst, 1),
	)

	got := FoldConstants(code, &constants)

	if len(got) != 6 {
		t.Fatalf("folded code len = %d, want jump plus reachable const: %v", len(got), got)
	}
	if bytecode.OpCode(got[0]) != bytecode.OpJump || bytecode.ReadU16(got, 1) != 3 {
		t.Fatalf("folded jump = %v target %d, want OpJump target 3", bytecode.OpCode(got[0]), bytecode.ReadU16(got, 1))
	}
	if bytecode.OpCode(got[3]) != bytecode.OpConst || bytecode.ReadU16(got, 4) != 1 {
		t.Fatalf("reachable instruction = %v, want OpConst 1", got[3:])
	}
}

type testInstr struct {
	op       bytecode.OpCode
	operands []uint16
}

func op(code bytecode.OpCode, operands ...uint16) testInstr {
	return testInstr{op: code, operands: operands}
}

func encodeOps(instrs ...testInstr) []byte {
	var code []byte
	for _, instr := range instrs {
		code = append(code, byte(instr.op))
		for _, operand := range instr.operands {
			code = append(code, byte(operand>>8), byte(operand))
		}
	}
	return code
}

func countOps(code []byte) map[bytecode.OpCode]int {
	counts := make(map[bytecode.OpCode]int)
	for i := 0; i < len(code); {
		codeOp := bytecode.OpCode(code[i])
		counts[codeOp]++
		i += 1 + bytecode.OperandWidth(codeOp)
	}
	return counts
}

func constantOperandBeforeReturn(t *testing.T, code []byte) uint16 {
	t.Helper()
	for i := 0; i < len(code); {
		codeOp := bytecode.OpCode(code[i])
		next := i + 1 + bytecode.OperandWidth(codeOp)
		if codeOp == bytecode.OpConst && next < len(code) && bytecode.OpCode(code[next]) == bytecode.OpReturnVal {
			return bytecode.ReadU16(code, i+1)
		}
		i = next
	}
	t.Fatalf("no OpConst immediately before OpReturnVal in %v", code)
	return 0
}
