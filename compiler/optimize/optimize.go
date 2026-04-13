// Package optimize implements a 4-pass bytecode optimization pipeline.
//
// The passes run in order after SSA→bytecode compilation:
//  1. Peephole — pattern-based superinstruction fusion (17 rules)
//  2. Slice fusion — OpIntSliceGet/Set/SetConst for native []int64
//  3. Int specialization — generic ops → OpInt* variants for int-typed locals
//  4. Move fusion — OpIntLocal+OpIntSetLocal → OpIntMoveLocal
//
// See docs/optimization-report.md for detailed performance analysis.
package optimize

import (
	"github.com/t04dJ14n9/gig/compiler/peephole"
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// Optimize applies all optimization passes to compiled bytecode in the correct order.
// localIsInt/constIsInt/localIsIntSlice flag which slots hold int-typed values.
// Returns the optimized code and whether int-specialized opcodes were emitted.
func Optimize(code []byte, localIsInt, constIsInt, localIsIntSlice []bool) ([]byte, bool) {
	code = Peephole(code)
	code = FuseSliceOps(code, localIsInt, localIsIntSlice)
	code, hasInt := IntSpecialize(code, localIsInt, constIsInt)
	code = FuseIntMoves(code)
	return code, hasInt
}

// --- Peephole ---

// Peephole performs peephole optimization on compiled bytecode.
func Peephole(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		matched := false
		for _, p := range peephole.Patterns() {
			if consumed, newBytes, ok := p.Match(code, i); ok {
				rewrites = append(rewrites, rewrite{i, i + consumed, newBytes})
				i += consumed
				matched = true
				break
			}
		}
		if !matched {
			i = instrEnd
		}
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

// --- Int specialization (table-driven) ---

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
	// 3-operand: local OP const → setlocal
	{bytecode.OpLocalConstAddSetLocal, bytecode.OpIntLocalConstAddSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	{bytecode.OpLocalConstSubSetLocal, bytecode.OpIntLocalConstSubSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	{bytecode.OpLocalConstMulSetLocal, bytecode.OpIntLocalConstMulSetLocal, [3]operandKind{opLocal, opConst, opLocal}, 3},
	// 3-operand: local OP local → setlocal
	{bytecode.OpLocalLocalAddSetLocal, bytecode.OpIntLocalLocalAddSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	{bytecode.OpLocalLocalSubSetLocal, bytecode.OpIntLocalLocalSubSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	{bytecode.OpLocalLocalMulSetLocal, bytecode.OpIntLocalLocalMulSetLocal, [3]operandKind{opLocal, opLocal, opLocal}, 3},
	// 2-operand: local CMP const → jump
	{bytecode.OpLessLocalConstJumpFalse, bytecode.OpIntLessLocalConstJumpFalse, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessLocalConstJumpTrue, bytecode.OpIntLessLocalConstJumpTrue, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpIntLessEqLocalConstJumpTrue, [3]operandKind{opLocal, opConst}, 2},
	{bytecode.OpLessEqLocalConstJumpFalse, bytecode.OpIntLessEqLocalConstJumpFalse, [3]operandKind{opLocal, opConst}, 2},
	// 2-operand: local CMP local → jump
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

	// Build a lookup table: generic opcode → rule index, for O(1) dispatch.
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
				// Mark local-typed operands as needing int shadows
				for j := byte(0); j < r.n; j++ {
					if r.ops[j] == opLocal {
						intUsed[bytecode.ReadU16(code, i+1+int(j)*2)] = true
					}
				}
			}
		} else {
			// Handle pre-fused slice ops (from FuseSliceOps)
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

// --- Slice operation fusion ---

// FuseSliceOps replaces common slice access patterns with fused superinstructions.
// It matches 17-byte sequences of [LOCAL, LOCAL, INDEXADDR, SETLOCAL, LOCAL, ...]
// and replaces them with 7-byte OpIntSlice{Get,Set,SetConst}.
func FuseSliceOps(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	var rewrites []rewrite

	// The pattern we match is always 17 bytes:
	//   LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) <op6> <op7>
	// where <op6>/<op7> determines the fused instruction:
	//   DEREF + SETLOCAL(v)      → OpIntSliceGet(s, j, v)
	//   LOCAL(val) + SETDEREF    → OpIntSliceSet(s, j, val)
	//   CONST(c) + SETDEREF      → OpIntSliceSetConst(s, j, c)
	const patternLen = 17

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		if op != bytecode.OpLocal || i+patternLen > len(code) {
			i = instrEnd
			continue
		}

		// Check the fixed prefix: LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr)
		if bytecode.OpCode(code[i+3]) != bytecode.OpLocal ||
			bytecode.OpCode(code[i+6]) != bytecode.OpIndexAddr ||
			bytecode.OpCode(code[i+7]) != bytecode.OpSetLocal ||
			bytecode.OpCode(code[i+10]) != bytecode.OpLocal {
			i = instrEnd
			continue
		}

		s := bytecode.ReadU16(code, i+1)
		j := bytecode.ReadU16(code, i+4)
		ptr := bytecode.ReadU16(code, i+8)
		ptrGet := bytecode.ReadU16(code, i+11)
		if ptr != ptrGet {
			i = instrEnd
			continue
		}

		// ptr must not be used elsewhere (it's a temporary)
		if localUsedOutside(code, ptr, i, i+patternLen) {
			i = instrEnd
			continue
		}

		op6 := bytecode.OpCode(code[i+13])
		op7 := bytecode.OpCode(code[i+14]) // for DEREF case, this is SETLOCAL opcode; for others, it's at i+16

		var fused []byte
		switch {
		case op6 == bytecode.OpDeref && op7 == bytecode.OpSetLocal:
			// s[j] read → OpIntSliceGet(s, j, v)
			v := bytecode.ReadU16(code, i+15)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) && safeIdx(int(v), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceGet, s, j, v)
			}
		case op6 == bytecode.OpLocal && bytecode.OpCode(code[i+16]) == bytecode.OpSetDeref:
			// s[j] = val → OpIntSliceSet(s, j, val)
			val := bytecode.ReadU16(code, i+14)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) && safeIdx(int(val), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceSet, s, j, val)
			}
		case op6 == bytecode.OpConst && bytecode.OpCode(code[i+16]) == bytecode.OpSetDeref:
			// s[j] = const → OpIntSliceSetConst(s, j, constIdx)
			constIdx := bytecode.ReadU16(code, i+14)
			if safeIdx(int(s), localIsIntSlice) && safeIdx(int(j), localIsInt) {
				fused = makeSliceOp(bytecode.OpIntSliceSetConst, s, j, constIdx)
			}
		}

		if fused != nil {
			rewrites = append(rewrites, rewrite{i, i + patternLen, fused})
			i += patternLen
		} else {
			i = instrEnd
		}
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

func makeSliceOp(op bytecode.OpCode, a, b, c uint16) []byte {
	out := make([]byte, 7)
	out[0] = byte(op)
	bytecode.WriteU16(out, 1, a)
	bytecode.WriteU16(out, 3, b)
	bytecode.WriteU16(out, 5, c)
	return out
}

// --- Int move fusion ---

// FuseIntMoves replaces OpIntLocal(A) OpIntSetLocal(B) pairs with OpIntMoveLocal(A, B).
func FuseIntMoves(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		if op == bytecode.OpIntLocal && i+6 <= len(code) &&
			bytecode.OpCode(code[i+3]) == bytecode.OpIntSetLocal {
			src := bytecode.ReadU16(code, i+1)
			dst := bytecode.ReadU16(code, i+4)
			newInstr := make([]byte, 5)
			newInstr[0] = byte(bytecode.OpIntMoveLocal)
			bytecode.WriteU16(newInstr, 1, src)
			bytecode.WriteU16(newInstr, 3, dst)
			rewrites = append(rewrites, rewrite{i, i + 6, newInstr})
			i += 6
			continue
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

// --- Rewrite infrastructure ---

type rewrite struct {
	oldStart int
	oldEnd   int
	newBytes []byte
}

func applyRewrites(code []byte, rewrites []rewrite) []byte {
	offsetMap := make([]int, len(code)+1)
	rIdx := 0
	shift := 0
	for pos := 0; pos <= len(code); pos++ {
		if rIdx < len(rewrites) && pos == rewrites[rIdx].oldStart {
			r := rewrites[rIdx]
			offsetMap[pos] = pos - shift
			shrink := (r.oldEnd - r.oldStart) - len(r.newBytes)
			for p := pos + 1; p < r.oldEnd; p++ {
				offsetMap[p] = pos - shift
			}
			shift += shrink
			pos = r.oldEnd - 1
			offsetMap[r.oldEnd-1] = r.oldEnd - 1 - shift
			rIdx++
			continue
		}
		offsetMap[pos] = pos - shift
	}

	newCode := make([]byte, 0, len(code))
	rIdx = 0
	for pos := 0; pos < len(code); {
		if rIdx < len(rewrites) && pos == rewrites[rIdx].oldStart {
			r := rewrites[rIdx]
			newCode = append(newCode, r.newBytes...)
			pos = r.oldEnd
			rIdx++
		} else {
			newCode = append(newCode, code[pos])
			pos++
		}
	}

	fixJumpTargets(newCode, offsetMap, len(code))
	return newCode
}

func fixJumpTargets(code []byte, offsetMap []int, oldLen int) {
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)

		switch op {
		// Simple jumps: target at offset +1
		case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
			oldTarget := int(bytecode.ReadU16(code, i+1))
			if oldTarget < oldLen {
				bytecode.WriteU16(code, i+1, uint16(offsetMap[oldTarget]))
			}
		// Fused compare-jump: target at offset +5
		case bytecode.OpLessLocalLocalJumpTrue, bytecode.OpLessLocalConstJumpTrue,
			bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpLessEqLocalConstJumpFalse,
			bytecode.OpGreaterLocalLocalJumpTrue,
			bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalConstJumpFalse,
			bytecode.OpIntLessLocalConstJumpFalse, bytecode.OpIntLessEqLocalConstJumpTrue,
			bytecode.OpIntLessEqLocalConstJumpFalse,
			bytecode.OpIntLessLocalLocalJumpFalse, bytecode.OpIntGreaterLocalLocalJumpTrue,
			bytecode.OpIntLessLocalConstJumpTrue, bytecode.OpIntLessLocalLocalJumpTrue:
			oldTarget := int(bytecode.ReadU16(code, i+5))
			if oldTarget < oldLen {
				bytecode.WriteU16(code, i+5, uint16(offsetMap[oldTarget]))
			}
		}

		i += 1 + width
	}
}

func opcodeWidth(op bytecode.OpCode) int {
	return bytecode.OperandWidth(op)
}

func safeIdx(idx int, flags []bool) bool {
	return idx < len(flags) && flags[idx]
}

func localUsedOutside(code []byte, localIdx uint16, skipStart, skipEnd int) bool {
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}
		if i >= skipStart && i < skipEnd {
			i = instrEnd
			continue
		}
		if (op == bytecode.OpLocal || op == bytecode.OpAddr) && width >= 2 {
			if bytecode.ReadU16(code, i+1) == localIdx {
				return true
			}
		}
		i = instrEnd
	}
	return false
}
