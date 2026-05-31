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
	intUsed := make([]bool, len(localIsInt))
	hasInt := false

	// Build a lookup table: generic opcode -> rule index, for O(1) dispatch.
	var ruleByOp [256]*intRule
	for i := range intRules {
		ruleByOp[intRules[i].from] = &intRules[i]
	}

	// Pass 1: discover which locals need int shadows.
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		if r := ruleByOp[op]; r != nil {
			allInt := true
			for j := byte(0); j < r.n; j++ {
				idx := int(bytecode.ReadU16(code, i+1+int(j)*2))
				if !safeIdx(idx, flagsFor(r.ops[j], localIsInt, constIsInt)) {
					allInt = false
					break
				}
			}
			if allInt {
				hasInt = true
				// Mark local-typed operands as needing int shadows.
				for j := byte(0); j < r.n; j++ {
					if r.ops[j] == opLocal {
						intUsed[bytecode.ReadU16(code, i+1+int(j)*2)] = true
					}
				}
			}
		} else {
			// Handle pre-fused slice ops from FuseSliceOps.
			switch op {
			case bytecode.OpIntSliceGet:
				intUsed[bytecode.ReadU16(code, i+3)] = true // j
				intUsed[bytecode.ReadU16(code, i+5)] = true // v
				hasInt = true
			case bytecode.OpIntSliceSet:
				intUsed[bytecode.ReadU16(code, i+3)] = true // j
				intUsed[bytecode.ReadU16(code, i+5)] = true // val
				hasInt = true
			case bytecode.OpIntSliceSetConst:
				intUsed[bytecode.ReadU16(code, i+3)] = true // j
				hasInt = true
			}
		}
		i = instrEnd
	}

	if !hasInt {
		return code, false
	}

	// Pass 2: rewrite opcodes and bridge OpLocal/OpSetLocal to int variants.
	i = 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		if r := ruleByOp[op]; r != nil {
			allInt := true
			for j := byte(0); j < r.n; j++ {
				idx := int(bytecode.ReadU16(code, i+1+int(j)*2))
				if !safeIdx(idx, flagsFor(r.ops[j], localIsInt, constIsInt)) {
					allInt = false
					break
				}
			}
			if allInt {
				code[i] = byte(r.to)
			}
		} else {
			switch op {
			case bytecode.OpSetLocal:
				a := int(bytecode.ReadU16(code, i+1))
				if a < len(intUsed) && intUsed[a] {
					code[i] = byte(bytecode.OpIntSetLocal)
				}
			case bytecode.OpLocal:
				a := int(bytecode.ReadU16(code, i+1))
				if a < len(intUsed) && intUsed[a] {
					code[i] = byte(bytecode.OpIntLocal)
				}
			}
		}
		i = instrEnd
	}

	return code, hasInt
}
