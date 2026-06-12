package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

// operandKind describes how to type-check an operand.
type operandKind byte

const (
	opLocal operandKind = iota // check against localIsInt
	opConst                    // check against constIsInt
)

// intRule describes one int-specialization rule:
//   - from: the generic superinstruction opcode to match
//   - to: the int-specialized opcode to rewrite to
//   - ops: operand kinds for type checking (2 or 3 elements)
//
// For 3-operand rules (arithmetic), operands are at code offsets +1, +3, +5.
// For 2-operand rules (compare-jump), operands are at code offsets +1, +3.
type intRule struct {
	from bytecode.OpCode
	to   bytecode.OpCode
	ops  [3]operandKind // ops[i] is only used if i < numOps()
	n    byte           // 2 or 3 operands
}

// intRuleTable gives each bytecode opcode an O(1) path to its int rewrite rule.
// It is rebuilt per optimizer call so the table cannot drift from intRules.
type intRuleTable [256]*intRule

// intInstruction carries the decoded instruction bounds shared by both passes.
// Keeping bounds in one small value avoids duplicating the truncate-at-end rule.
type intInstruction struct {
	op    bytecode.OpCode
	start int
	end   int
}

// intRules is the table of all int-specialization rules.
// Grouped by operand count for clarity.
var intRules = [...]intRule{
	// 3-operand: local OP const -> setlocal
	{bytecode.OpLocalConstAddSetLocal, bytecode.OpIntLocalConstAddSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	{bytecode.OpLocalConstSubSetLocal, bytecode.OpIntLocalConstSubSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	{bytecode.OpLocalConstMulSetLocal, bytecode.OpIntLocalConstMulSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	// 3-operand: local OP local -> setlocal
	{bytecode.OpLocalLocalAddSetLocal, bytecode.OpIntLocalLocalAddSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	{bytecode.OpLocalLocalSubSetLocal, bytecode.OpIntLocalLocalSubSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	{bytecode.OpLocalLocalMulSetLocal, bytecode.OpIntLocalLocalMulSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	// 2-operand: local CMP const -> jump
	{bytecode.OpLessLocalConstJumpFalse, bytecode.OpIntLessLocalConstJumpFalse, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessLocalConstJumpTrue, bytecode.OpIntLessLocalConstJumpTrue, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpIntLessEqLocalConstJumpTrue, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessEqLocalConstJumpFalse, bytecode.OpIntLessEqLocalConstJumpFalse, [3]operandKind{opLocal, opConst}, 2},
	// 2-operand: local CMP local -> jump
	{bytecode.OpLessLocalLocalJumpFalse, bytecode.OpIntLessLocalLocalJumpFalse, [3]operandKind{opLocal, opLocal}, 2},
	{bytecode.OpLessLocalLocalJumpTrue, bytecode.OpIntLessLocalLocalJumpTrue, [3]operandKind{opLocal, opLocal}, 2},
	{bytecode.OpGreaterLocalLocalJumpTrue, bytecode.OpIntGreaterLocalLocalJumpTrue, [3]operandKind{opLocal, opLocal}, 2},
}

// flagsFor returns the type-check flags for an operand kind.
func flagsFor(k operandKind, localIsInt, constIsInt []bool) []bool {
	if k == opConst {
		return constIsInt
	}
	return localIsInt
}

// IntSpecialize upgrades Value-based superinstructions to OpInt* variants
// when all involved locals and constants are int-typed. Two-pass: first pass
// discovers which locals need int shadows, second pass rewrites opcodes.
func IntSpecialize(code []byte, localIsInt, constIsInt []bool) ([]byte, bool) {
	ruleByOp := newIntRuleTable()
	intUsed, hasInt := discoverIntShadowUses(code, localIsInt, constIsInt, ruleByOp)

	if !hasInt {
		return code, false
	}

	rewriteIntSpecializedOps(code, localIsInt, constIsInt, intUsed, ruleByOp)
	return code, true
}

// newIntRuleTable centralizes opcode-to-rule indexing for discovery and rewrite.
func newIntRuleTable() intRuleTable {
	var ruleByOp intRuleTable
	for i := range intRules {
		ruleByOp[intRules[i].from] = &intRules[i]
	}
	return ruleByOp
}

// discoverIntShadowUses is the first pass: find locals that need int shadow slots.
// Rule-backed instructions mark all local operands; fused slice ops mark their int index/value locals.
func discoverIntShadowUses(code []byte, localIsInt, constIsInt []bool, ruleByOp intRuleTable) ([]bool, bool) {
	intUsed := make([]bool, len(localIsInt))
	hasInt := false

	for i := 0; i < len(code); {
		instr, ok := nextIntInstruction(code, i)
		if !ok {
			break
		}

		if markIntUsesForInstruction(code, instr, intUsed, localIsInt, constIsInt, ruleByOp) {
			hasInt = true
		}
		i = instr.end
	}

	return intUsed, hasInt
}

// nextIntInstruction decodes only enough shape for this optimizer and preserves the old behavior of stopping at truncated bytecode.
func nextIntInstruction(code []byte, start int) (intInstruction, bool) {
	op := bytecode.OpCode(code[start])
	end := start + 1 + opcodeWidth(op)
	if end > len(code) {
		return intInstruction{}, false
	}
	return intInstruction{op: op, start: start, end: end}, true
}

// markIntUsesForInstruction handles one first-pass instruction and reports whether any int specialization is present.
func markIntUsesForInstruction(code []byte, instr intInstruction, intUsed []bool, localIsInt, constIsInt []bool, ruleByOp intRuleTable) bool {
	if r := ruleByOp[instr.op]; r != nil {
		if !ruleOperandsAreInt(code, instr.start, r, localIsInt, constIsInt) {
			return false
		}
		markRuleLocalOperands(code, instr.start, r, intUsed)
		return true
	}

	return markFusedSliceIntUses(code, instr.start, instr.op, intUsed)
}

// ruleOperandsAreInt is shared by both passes so the rewrite pass uses the same type predicate discovered in pass one.
func ruleOperandsAreInt(code []byte, start int, r *intRule, localIsInt, constIsInt []bool) bool {
	for j := byte(0); j < r.n; j++ {
		idx := int(bytecode.ReadU16(code, intOperandOffset(start, j)))
		if !safeIdx(idx, flagsFor(r.ops[j], localIsInt, constIsInt)) {
			return false
		}
	}
	return true
}

func markRuleLocalOperands(code []byte, start int, r *intRule, intUsed []bool) {
	for j := byte(0); j < r.n; j++ {
		if r.ops[j] == opLocal {
			markIntLocalUse(code, intOperandOffset(start, j), intUsed)
		}
	}
}

// markFusedSliceIntUses covers slice superinstructions that were already converted by FuseSliceOps before this pass.
func markFusedSliceIntUses(code []byte, start int, op bytecode.OpCode, intUsed []bool) bool {
	switch op {
	case bytecode.OpIntSliceGet:
		markIntLocalUse(code, start+3, intUsed) // j
		markIntLocalUse(code, start+5, intUsed) // v
	case bytecode.OpIntSliceSet:
		markIntLocalUse(code, start+3, intUsed) // j
		markIntLocalUse(code, start+5, intUsed) // val
	case bytecode.OpIntSliceSetConst:
		markIntLocalUse(code, start+3, intUsed) // j
	default:
		return false
	}
	return true
}

// rewriteIntSpecializedOps is the second pass: rewrite eligible superinstructions and bridge OpLocal/OpSetLocal users.
func rewriteIntSpecializedOps(code []byte, localIsInt, constIsInt, intUsed []bool, ruleByOp intRuleTable) {
	for i := 0; i < len(code); {
		instr, ok := nextIntInstruction(code, i)
		if !ok {
			break
		}

		rewriteIntInstruction(code, instr, localIsInt, constIsInt, intUsed, ruleByOp)
		i = instr.end
	}
}

func rewriteIntInstruction(code []byte, instr intInstruction, localIsInt, constIsInt, intUsed []bool, ruleByOp intRuleTable) {
	if r := ruleByOp[instr.op]; r != nil {
		if ruleOperandsAreInt(code, instr.start, r, localIsInt, constIsInt) {
			code[instr.start] = byte(r.to)
		}
		return
	}

	rewriteIntLocalBridge(code, instr.start, instr.op, intUsed)
}

func rewriteIntLocalBridge(code []byte, start int, op bytecode.OpCode, intUsed []bool) {
	switch op {
	case bytecode.OpSetLocal:
		if isIntShadowedLocal(code, start+1, intUsed) {
			code[start] = byte(bytecode.OpIntSetLocal)
		}
	case bytecode.OpLocal:
		if isIntShadowedLocal(code, start+1, intUsed) {
			code[start] = byte(bytecode.OpIntLocal)
		}
	}
}

func intOperandOffset(start int, operand byte) int {
	return start + 1 + int(operand)*2
}

// markIntLocalUse intentionally preserves the previous direct indexing semantics for compiler-emitted bytecode.
// The rule path validates local operands before calling this; fused slice bytecode is generated internally by FuseSliceOps.
func markIntLocalUse(code []byte, offset int, intUsed []bool) {
	intUsed[bytecode.ReadU16(code, offset)] = true
}

func isIntShadowedLocal(code []byte, offset int, intUsed []bool) bool {
	idx := int(bytecode.ReadU16(code, offset))
	return idx < len(intUsed) && intUsed[idx]
}
