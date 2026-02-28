package compiler

import "gig/bytecode"

// rewrite represents a bytecode rewrite from an old range to new bytes.
type rewrite struct {
	oldStart int
	oldEnd   int // exclusive
	newBytes []byte
}

// optimizeBytecode performs a peephole optimization pass on the compiled bytecode.
// It scans for common instruction sequences and replaces them with fused superinstructions.
// This must be called after patchJumps so all offsets are resolved.
//
// The optimizer rewrites instructions in-place. Because superinstructions are always
// shorter than or equal to the sequences they replace, the bytecode may shrink. We
// rebuild the instruction slice and fix jump targets accordingly.
//
//nolint:gocyclo,cyclop,funlen,maintidx,gocognit
func optimizeBytecode(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width

		if instrEnd > len(code) {
			break
		}

		// Pattern: OpLocal(A) OpLocal(B) OpAdd OpSetLocal(C) -> OpLocalLocalAddSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpAdd && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalLocalAddSetLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, c)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpAdd OpSetLocal(C) -> OpLocalConstAddSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpAdd && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalConstAddSetLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, c)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpSub OpSetLocal(C) -> OpLocalConstSubSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpSub && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalConstSubSetLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, c)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpLess OpJumpTrue(offset) -> OpLessLocalLocalJumpTrue(A,B,offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpLess && op4 == bytecode.OpJumpTrue {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLessLocalLocalJumpTrue)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpLess OpJumpFalse(offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpLess && op4 == bytecode.OpJumpFalse {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLessLocalLocalJumpFalse)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpLess OpJumpTrue(offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpLess && op4 == bytecode.OpJumpTrue {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLessLocalConstJumpTrue)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpLess OpJumpFalse(offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpLess && op4 == bytecode.OpJumpFalse {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLessLocalConstJumpFalse)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpLessEq OpJumpTrue(offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpLessEq && op4 == bytecode.OpJumpTrue {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLessEqLocalConstJumpTrue)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpGreater OpJumpTrue(offset)
		if op == bytecode.OpLocal && i+10 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpGreater && op4 == bytecode.OpJumpTrue {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				offset := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpGreaterLocalLocalJumpTrue)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, offset)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpAdd -> OpAddLocalLocal(A,B)
		if op == bytecode.OpLocal && i+7 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpAdd {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpAddLocalLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				rewrites = append(rewrites, rewrite{i, i + 7, newInstr})
				i += 7
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpSub -> OpSubLocalLocal(A,B)
		if op == bytecode.OpLocal && i+7 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpSub {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpSubLocalLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				rewrites = append(rewrites, rewrite{i, i + 7, newInstr})
				i += 7
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpMul -> OpMulLocalLocal(A,B)
		if op == bytecode.OpLocal && i+7 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpMul {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpMulLocalLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				rewrites = append(rewrites, rewrite{i, i + 7, newInstr})
				i += 7
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpAdd -> OpAddLocalConst(A,B)
		if op == bytecode.OpLocal && i+7 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			if op2 == bytecode.OpConst && op3 == bytecode.OpAdd {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpAddLocalConst)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				rewrites = append(rewrites, rewrite{i, i + 7, newInstr})
				i += 7
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpSub -> OpSubLocalConst(A,B)
		if op == bytecode.OpLocal && i+7 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			if op2 == bytecode.OpConst && op3 == bytecode.OpSub {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpSubLocalConst)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				rewrites = append(rewrites, rewrite{i, i + 7, newInstr})
				i += 7
				continue
			}
		}

		// Pattern: OpAdd OpSetLocal(A) -> OpAddSetLocal(A)
		if op == bytecode.OpAdd && i+4 <= len(code) {
			op2 := bytecode.OpCode(code[i+1])
			if op2 == bytecode.OpSetLocal {
				a := readU16(code, i+2)
				newInstr := make([]byte, 3)
				newInstr[0] = byte(bytecode.OpAddSetLocal)
				writeU16(newInstr, 1, a)
				rewrites = append(rewrites, rewrite{i, i + 4, newInstr})
				i += 4
				continue
			}
		}

		// Pattern: OpSub OpSetLocal(A) -> OpSubSetLocal(A)
		if op == bytecode.OpSub && i+4 <= len(code) {
			op2 := bytecode.OpCode(code[i+1])
			if op2 == bytecode.OpSetLocal {
				a := readU16(code, i+2)
				newInstr := make([]byte, 3)
				newInstr[0] = byte(bytecode.OpSubSetLocal)
				writeU16(newInstr, 1, a)
				rewrites = append(rewrites, rewrite{i, i + 4, newInstr})
				i += 4
				continue
			}
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}

	// Apply rewrites: build offset mapping and new code
	return applyRewrites(code, rewrites)
}

// applyRewrites applies a list of non-overlapping rewrites to the bytecode and fixes jump targets.
func applyRewrites(code []byte, rewrites []rewrite) []byte {
	// Build cumulative shift map: for every old offset, how much did it shrink?
	// We only need to track offsets at instruction boundaries.
	offsetMap := make([]int, len(code)+1) // maps old byte offset to new byte offset
	rIdx := 0
	shift := 0
	for pos := 0; pos <= len(code); pos++ {
		if rIdx < len(rewrites) && pos == rewrites[rIdx].oldStart {
			r := rewrites[rIdx]
			// Map the start
			offsetMap[pos] = pos - shift
			// Map intermediate positions to the same new position
			shrink := (r.oldEnd - r.oldStart) - len(r.newBytes)
			for p := pos + 1; p < r.oldEnd; p++ {
				offsetMap[p] = pos - shift
			}
			shift += shrink
			pos = r.oldEnd - 1                           // loop will pos++
			offsetMap[r.oldEnd-1] = r.oldEnd - 1 - shift // map end-1
			rIdx++
			continue
		}
		offsetMap[pos] = pos - shift
	}

	// Build new code with rewrites applied
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

	// Fix jump targets in the new code
	fixJumpTargets(newCode, offsetMap, len(code))

	return newCode
}

// fixJumpTargets scans the new bytecode for all jump instructions and remaps their targets.
func fixJumpTargets(code []byte, offsetMap []int, oldLen int) {
	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)

		switch op {
		case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
			oldTarget := int(readU16(code, i+1))
			if oldTarget < oldLen {
				newTarget := offsetMap[oldTarget]
				writeU16(code, i+1, uint16(newTarget))
			}

		case bytecode.OpLessLocalLocalJumpTrue, bytecode.OpLessLocalConstJumpTrue,
			bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpGreaterLocalLocalJumpTrue,
			bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalConstJumpFalse:
			oldTarget := int(readU16(code, i+5))
			if oldTarget < oldLen {
				newTarget := offsetMap[oldTarget]
				writeU16(code, i+5, uint16(newTarget))
			}
		}

		i += 1 + width
	}
}

// opcodeWidth returns the total operand byte width for an opcode.
func opcodeWidth(op bytecode.OpCode) int {
	if w, ok := bytecode.OperandWidths[op]; ok {
		return w
	}
	return 0
}

func readU16(code []byte, offset int) uint16 {
	return uint16(code[offset])<<8 | uint16(code[offset+1])
}

func writeU16(code []byte, offset int, val uint16) {
	code[offset] = byte(val >> 8)
	code[offset+1] = byte(val)
}
