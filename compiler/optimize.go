package compiler

import "github.com/t04dJ14n9/gig/bytecode"

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

		// ============================================================
		// Extended compare patterns: Local(A) Const/Local(B) CmpOp SetLocal(X) Local(X) Jump*
		// SSA generates `t = a cmp b; if t goto ...` which compiles to
		// LOCAL(A) CONST/LOCAL(B) CMP SETLOCAL(X) LOCAL(X) JUMP*(off)
		// We fuse these 6-instruction sequences into a single compare+jump.
		// ============================================================

		// Pattern: OpLocal(A) OpConst(B) OpLessEq OpSetLocal(X) OpLocal(X) OpJumpFalse(off)
		// -> OpLessEqLocalConstJumpFalse (which acts as "if a > b, jump")
		// Since we don't have OpLessEqJumpFalse, we note that "NOT(a<=b)" == "a>b"
		// We reuse the LessEq+JumpFalse semantics but need a dedicated approach.
		// Simplest: fuse into OpLessEqLocalConstJumpTrue by flipping the jump target.
		// Actually, OpJumpFalse means "if condition false, jump". So if (a<=b) is false, jump.
		// The simplest approach: add a new combined opcode, or reuse existing infrastructure.
		//
		// For now, fuse the entire 6 instructions to eliminate SetLocal+Local overhead.
		// The 6-byte pattern is:
		//   LOCAL(A) [3] + CONST(B) [3] + LESSEQ [1] + SETLOCAL(X) [3] + LOCAL(X) [3] + JUMPFALSE(off) [3] = 16 bytes
		// Fused: OpLessEqLocalConstJumpFalse(A, B, off) = 7 bytes

		if op == bytecode.OpLocal && i+16 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			op5 := bytecode.OpCode(code[i+10])
			op6 := bytecode.OpCode(code[i+13])
			if op4 == bytecode.OpSetLocal && op5 == bytecode.OpLocal {
				setIdx := readU16(code, i+8)
				getIdx := readU16(code, i+11)
				if setIdx == getIdx {
					a := readU16(code, i+1)
					b := readU16(code, i+4)
					// Local(A) Const(B) LessEq SetLocal(X) Local(X) JumpFalse(off)
					if op2 == bytecode.OpConst && op3 == bytecode.OpLessEq && op6 == bytecode.OpJumpFalse {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessEqLocalConstJumpFalse)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Const(B) LessEq SetLocal(X) Local(X) JumpTrue(off)
					if op2 == bytecode.OpConst && op3 == bytecode.OpLessEq && op6 == bytecode.OpJumpTrue {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessEqLocalConstJumpTrue)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Const(B) Less SetLocal(X) Local(X) JumpFalse(off)
					if op2 == bytecode.OpConst && op3 == bytecode.OpLess && op6 == bytecode.OpJumpFalse {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessLocalConstJumpFalse)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Const(B) Less SetLocal(X) Local(X) JumpTrue(off)
					if op2 == bytecode.OpConst && op3 == bytecode.OpLess && op6 == bytecode.OpJumpTrue {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessLocalConstJumpTrue)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Local(B) Less SetLocal(X) Local(X) JumpFalse(off)
					if op2 == bytecode.OpLocal && op3 == bytecode.OpLess && op6 == bytecode.OpJumpFalse {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessLocalLocalJumpFalse)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Local(B) Less SetLocal(X) Local(X) JumpTrue(off)
					if op2 == bytecode.OpLocal && op3 == bytecode.OpLess && op6 == bytecode.OpJumpTrue {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpLessLocalLocalJumpTrue)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
					// Local(A) Local(B) Greater SetLocal(X) Local(X) JumpTrue(off)
					if op2 == bytecode.OpLocal && op3 == bytecode.OpGreater && op6 == bytecode.OpJumpTrue {
						offset := readU16(code, i+14)
						newInstr := make([]byte, 7)
						newInstr[0] = byte(bytecode.OpGreaterLocalLocalJumpTrue)
						writeU16(newInstr, 1, a)
						writeU16(newInstr, 3, b)
						writeU16(newInstr, 5, offset)
						rewrites = append(rewrites, rewrite{i, i + 16, newInstr})
						i += 16
						continue
					}
				}
			}
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

		// Pattern: OpLocal(A) OpLocal(B) OpSub OpSetLocal(C) -> OpLocalLocalSubSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpSub && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalLocalSubSetLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, c)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpLocal(B) OpMul OpSetLocal(C) -> OpLocalLocalMulSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpLocal && op3 == bytecode.OpMul && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalLocalMulSetLocal)
				writeU16(newInstr, 1, a)
				writeU16(newInstr, 3, b)
				writeU16(newInstr, 5, c)
				rewrites = append(rewrites, rewrite{i, i + 10, newInstr})
				i += 10
				continue
			}
		}

		// Pattern: OpLocal(A) OpConst(B) OpMul OpSetLocal(C) -> OpLocalConstMulSetLocal(A,B,C)
		if op == bytecode.OpLocal && i+12 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			if op2 == bytecode.OpConst && op3 == bytecode.OpMul && op4 == bytecode.OpSetLocal {
				a := readU16(code, i+1)
				b := readU16(code, i+4)
				c := readU16(code, i+8)
				newInstr := make([]byte, 7)
				newInstr[0] = byte(bytecode.OpLocalConstMulSetLocal)
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

		// Pattern: OpJump that jumps to the very next instruction → remove
		if op == bytecode.OpJump {
			target := readU16(code, i+1)
			if int(target) == instrEnd {
				rewrites = append(rewrites, rewrite{i, instrEnd, nil})
				i = instrEnd
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
			bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpLessEqLocalConstJumpFalse,
			bytecode.OpGreaterLocalLocalJumpTrue,
			bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalConstJumpFalse,
			bytecode.OpIntLessLocalConstJumpFalse, bytecode.OpIntLessEqLocalConstJumpTrue,
			bytecode.OpIntLessEqLocalConstJumpFalse,
			bytecode.OpIntLessLocalLocalJumpFalse, bytecode.OpIntGreaterLocalLocalJumpTrue,
			bytecode.OpIntLessLocalConstJumpTrue, bytecode.OpIntLessLocalLocalJumpTrue:
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
// Uses O(1) array lookup instead of map lookup.
func opcodeWidth(op bytecode.OpCode) int {
	return bytecode.OperandWidth(op)
}

func readU16(code []byte, offset int) uint16 {
	return uint16(code[offset])<<8 | uint16(code[offset+1])
}

func writeU16(code []byte, offset int, val uint16) {
	code[offset] = byte(val >> 8)
	code[offset+1] = byte(val)
}

// intSpecialize performs a two-pass upgrade of Value-based superinstructions
// to OpInt* variants when all involved locals and constants are int-typed.
//
// Pass 1: Scan for superinstructions that CAN be upgraded. Collect which local
// indices participate in int-specialized ops ("intUsed" set).
//
// Pass 2: Upgrade the superinstructions to OpInt* variants AND convert
// OpSetLocal/OpLocal for intUsed locals to OpIntSetLocal/OpIntLocal bridges.
//
// Returns the modified code and whether any OpInt* opcodes were emitted.
//
//nolint:gocyclo,cyclop,funlen,maintidx,gocognit
func intSpecialize(code []byte, localIsInt, constIsInt []bool) ([]byte, bool) {
	// Pass 1: identify which local indices will participate in int ops.
	intUsed := make([]bool, len(localIsInt))
	hasInt := false

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		switch op {
		case bytecode.OpLocalConstAddSetLocal, bytecode.OpLocalConstSubSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalLocalAddSetLocal, bytecode.OpLocalLocalSubSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalLocalMulSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalConstMulSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLessLocalConstJumpFalse, bytecode.OpLessLocalConstJumpTrue,
			bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpLessEqLocalConstJumpFalse:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				intUsed[a] = true
				hasInt = true
			}
		case bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalLocalJumpTrue,
			bytecode.OpGreaterLocalLocalJumpTrue:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				hasInt = true
			}
		case bytecode.OpIntSliceGet:
			// Operands: slice(2) index(2) dest(2)
			// index and dest are int locals that need intLocals sync
			j := int(readU16(code, i+3))
			v := int(readU16(code, i+5))
			intUsed[j] = true
			intUsed[v] = true
			hasInt = true
		case bytecode.OpIntSliceSet:
			// Operands: slice(2) index(2) val(2)
			// index and val are int locals that need intLocals sync
			j := int(readU16(code, i+3))
			val := int(readU16(code, i+5))
			intUsed[j] = true
			intUsed[val] = true
			hasInt = true
		case bytecode.OpIntSliceSetConst:
			// Operands: slice(2) index(2) const(2)
			// index is an int local that needs intLocals sync
			j := int(readU16(code, i+3))
			intUsed[j] = true
			hasInt = true
		}
		i = instrEnd
	}

	if !hasInt {
		return code, false
	}

	// Pass 2: upgrade superinstructions and bridge OpSetLocal/OpLocal.
	i = 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		switch op {
		case bytecode.OpLocalConstAddSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstAddSetLocal)
			}
		case bytecode.OpLocalConstSubSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstSubSetLocal)
			}
		case bytecode.OpLocalLocalAddSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalAddSetLocal)
			}
		case bytecode.OpLocalLocalSubSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalSubSetLocal)
			}
		case bytecode.OpLocalLocalMulSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalMulSetLocal)
			}
		case bytecode.OpLocalConstMulSetLocal:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			c := int(readU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstMulSetLocal)
			}
		case bytecode.OpLessLocalConstJumpFalse:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalConstJumpFalse)
			}
		case bytecode.OpLessLocalConstJumpTrue:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalConstJumpTrue)
			}
		case bytecode.OpLessEqLocalConstJumpTrue:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessEqLocalConstJumpTrue)
			}
		case bytecode.OpLessEqLocalConstJumpFalse:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessEqLocalConstJumpFalse)
			}
		case bytecode.OpLessLocalLocalJumpFalse:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalLocalJumpFalse)
			}
		case bytecode.OpLessLocalLocalJumpTrue:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalLocalJumpTrue)
			}
		case bytecode.OpGreaterLocalLocalJumpTrue:
			a := int(readU16(code, i+1))
			b := int(readU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntGreaterLocalLocalJumpTrue)
			}
		// Bridge: upgrade OpSetLocal/OpLocal for int-participating locals
		case bytecode.OpSetLocal:
			a := int(readU16(code, i+1))
			if a < len(intUsed) && intUsed[a] {
				code[i] = byte(bytecode.OpIntSetLocal)
			}
		case bytecode.OpLocal:
			a := int(readU16(code, i+1))
			if a < len(intUsed) && intUsed[a] {
				code[i] = byte(bytecode.OpIntLocal)
			}
		}

		i = instrEnd
	}

	return code, hasInt
}

// safeIdx returns true if idx is within bounds and the flag is true.
func safeIdx(idx int, flags []bool) bool {
	return idx < len(flags) && flags[idx]
}

// fuseSliceOps replaces common slice access patterns with fused superinstructions.
// This must run after optimizeBytecode (which doesn't consume these patterns) and
// before intSpecialize (which upgrades LOCAL→INTLOCAL).
//
// Pattern 1 (slice read): LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) DEREF SETLOCAL(v)
//
//	→ OpIntSliceGet(s, j, v)  when s is []int and j, v are int
//
// Pattern 2 (slice write): LOCAL(s) LOCAL(j) INDEXADDR SETLOCAL(ptr) LOCAL(ptr) LOCAL(val) SETDEREF
//
//	→ OpIntSliceSet(s, j, val)  when s is []int and j, val are int
//
//nolint:gocyclo,cyclop,funlen,gocognit
func fuseSliceOps(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		// Both patterns start with LOCAL INDEXADDR SETLOCAL LOCAL = 13 bytes
		// Pattern 1 continues with DEREF SETLOCAL = +4 bytes = 17 total
		// Pattern 2 continues with LOCAL SETDEREF = +4 bytes = 17 total
		if op == bytecode.OpLocal && i+17 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			op5 := bytecode.OpCode(code[i+10])

			if op2 == bytecode.OpLocal && op3 == bytecode.OpIndexAddr &&
				op4 == bytecode.OpSetLocal && op5 == bytecode.OpLocal {
				s := readU16(code, i+1)       // slice local
				j := readU16(code, i+4)       // index local
				ptr := readU16(code, i+8)     // ptr temp local
				ptrGet := readU16(code, i+11) // must match ptr

				if ptr == ptrGet {
					op6 := bytecode.OpCode(code[i+13])

					// Pattern 1: ... DEREF SETLOCAL(v) → OpIntSliceGet(s, j, v)
					if op6 == bytecode.OpDeref && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+14])
						if op7 == bytecode.OpSetLocal {
							v := readU16(code, i+15)
							// Check types: s is []int, j is int, v is int
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) &&
								safeIdx(int(v), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceGet)
								writeU16(newInstr, 1, s)
								writeU16(newInstr, 3, j)
								writeU16(newInstr, 5, v)
								rewrites = append(rewrites, rewrite{i, i + 17, newInstr})
								i += 17
								continue
							}
						}
					}

					// Pattern 2: ... LOCAL(val) SETDEREF → OpIntSliceSet(s, j, val)
					if op6 == bytecode.OpLocal && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+16])
						if op7 == bytecode.OpSetDeref {
							val := readU16(code, i+14)
							// Check types: s is []int, j is int, val is int
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) &&
								safeIdx(int(val), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceSet)
								writeU16(newInstr, 1, s)
								writeU16(newInstr, 3, j)
								writeU16(newInstr, 5, val)
								rewrites = append(rewrites, rewrite{i, i + 17, newInstr})
								i += 17
								continue
							}
						}
					}

					// Pattern 3: ... CONST(val) SETDEREF → OpIntSliceSetConst(s, j, const_idx)
					if op6 == bytecode.OpConst && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+16])
						if op7 == bytecode.OpSetDeref {
							constIdx := readU16(code, i+14)
							// Check types: s is []int, j is int
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceSetConst)
								writeU16(newInstr, 1, s)
								writeU16(newInstr, 3, j)
								writeU16(newInstr, 5, constIdx)
								rewrites = append(rewrites, rewrite{i, i + 17, newInstr})
								i += 17
								continue
							}
						}
					}
				}
			}
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

// fuseIntMoves replaces OpIntLocal(A) OpIntSetLocal(B) pairs with OpIntMoveLocal(A, B).
// This eliminates unnecessary push+pop for phi-move patterns in int-specialized code.
func fuseIntMoves(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		// Pattern: OpIntLocal(A) OpIntSetLocal(B) → OpIntMoveLocal(A, B)
		if op == bytecode.OpIntLocal && i+6 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			if op2 == bytecode.OpIntSetLocal {
				src := readU16(code, i+1)
				dst := readU16(code, i+4)
				newInstr := make([]byte, 5)
				newInstr[0] = byte(bytecode.OpIntMoveLocal)
				writeU16(newInstr, 1, src)
				writeU16(newInstr, 3, dst)
				rewrites = append(rewrites, rewrite{i, i + 6, newInstr})
				i += 6
				continue
			}
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}
