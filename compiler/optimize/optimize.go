// Package optimize implements a 4-pass bytecode optimization pipeline.
//
// The passes run in order after SSA→bytecode compilation:
//   1. Peephole — pattern-based superinstruction fusion (17 rules)
//   2. Slice fusion — OpIntSliceGet/Set/SetConst for native []int64
//   3. Int specialization — generic ops → OpInt* variants for int-typed locals
//   4. Move fusion — OpIntLocal+OpIntSetLocal → OpIntMoveLocal
//
// See docs/optimization-report.md for detailed performance analysis.
package optimize

import (
	"git.woa.com/youngjin/gig/bytecode"
	"git.woa.com/youngjin/gig/compiler/peephole"
)

// Optimize applies all optimization passes to compiled bytecode in the correct order:
//  1. Peephole optimization (pattern matching → superinstructions)
//  2. Slice operation fusion (→ OpIntSliceGet/Set/SetConst)
//  3. Int specialization (generic ops → OpInt* variants)
//  4. Int move fusion (OpIntLocal+OpIntSetLocal → OpIntMoveLocal)
//
// localIsInt indicates which local slots hold int-typed values.
// constIsInt indicates which constant pool entries hold int-typed values.
// localIsIntSlice indicates which locals hold []int-typed values.
//
// Returns the optimized code and whether int-specialized opcodes were emitted.
func Optimize(code []byte, localIsInt, constIsInt, localIsIntSlice []bool) ([]byte, bool) {
	code = Peephole(code)
	code = FuseSliceOps(code, localIsInt, localIsIntSlice)
	code, hasInt := IntSpecialize(code, localIsInt, constIsInt)
	code = FuseIntMoves(code)
	return code, hasInt
}

// Peephole performs peephole optimization on compiled bytecode.
func Peephole(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width

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

// IntSpecialize upgrades Value-based superinstructions to OpInt* variants
// when all involved locals and constants are int-typed.
func IntSpecialize(code []byte, localIsInt, constIsInt []bool) ([]byte, bool) {
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
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalLocalAddSetLocal, bytecode.OpLocalLocalSubSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalLocalMulSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLocalConstMulSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				intUsed[a] = true
				intUsed[c] = true
				hasInt = true
			}
		case bytecode.OpLessLocalConstJumpFalse, bytecode.OpLessLocalConstJumpTrue,
			bytecode.OpLessEqLocalConstJumpTrue, bytecode.OpLessEqLocalConstJumpFalse:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				intUsed[a] = true
				hasInt = true
			}
		case bytecode.OpLessLocalLocalJumpFalse, bytecode.OpLessLocalLocalJumpTrue,
			bytecode.OpGreaterLocalLocalJumpTrue:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				intUsed[a] = true
				intUsed[b] = true
				hasInt = true
			}
		case bytecode.OpIntSliceGet:
			j := int(bytecode.ReadU16(code, i+3))
			v := int(bytecode.ReadU16(code, i+5))
			intUsed[j] = true
			intUsed[v] = true
			hasInt = true
		case bytecode.OpIntSliceSet:
			j := int(bytecode.ReadU16(code, i+3))
			val := int(bytecode.ReadU16(code, i+5))
			intUsed[j] = true
			intUsed[val] = true
			hasInt = true
		case bytecode.OpIntSliceSetConst:
			j := int(bytecode.ReadU16(code, i+3))
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
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstAddSetLocal)
			}
		case bytecode.OpLocalConstSubSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstSubSetLocal)
			}
		case bytecode.OpLocalLocalAddSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalAddSetLocal)
			}
		case bytecode.OpLocalLocalSubSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalSubSetLocal)
			}
		case bytecode.OpLocalLocalMulSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalLocalMulSetLocal)
			}
		case bytecode.OpLocalConstMulSetLocal:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			c := int(bytecode.ReadU16(code, i+5))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) && safeIdx(c, localIsInt) {
				code[i] = byte(bytecode.OpIntLocalConstMulSetLocal)
			}
		case bytecode.OpLessLocalConstJumpFalse:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalConstJumpFalse)
			}
		case bytecode.OpLessLocalConstJumpTrue:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalConstJumpTrue)
			}
		case bytecode.OpLessEqLocalConstJumpTrue:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessEqLocalConstJumpTrue)
			}
		case bytecode.OpLessEqLocalConstJumpFalse:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, constIsInt) {
				code[i] = byte(bytecode.OpIntLessEqLocalConstJumpFalse)
			}
		case bytecode.OpLessLocalLocalJumpFalse:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalLocalJumpFalse)
			}
		case bytecode.OpLessLocalLocalJumpTrue:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntLessLocalLocalJumpTrue)
			}
		case bytecode.OpGreaterLocalLocalJumpTrue:
			a := int(bytecode.ReadU16(code, i+1))
			b := int(bytecode.ReadU16(code, i+3))
			if safeIdx(a, localIsInt) && safeIdx(b, localIsInt) {
				code[i] = byte(bytecode.OpIntGreaterLocalLocalJumpTrue)
			}
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

		i = instrEnd
	}

	return code, hasInt
}

// FuseSliceOps replaces common slice access patterns with fused superinstructions.
func FuseSliceOps(code []byte, localIsInt, localIsIntSlice []bool) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		if op == bytecode.OpLocal && i+17 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			op3 := bytecode.OpCode(code[i+6])
			op4 := bytecode.OpCode(code[i+7])
			op5 := bytecode.OpCode(code[i+10])

			if op2 == bytecode.OpLocal && op3 == bytecode.OpIndexAddr &&
				op4 == bytecode.OpSetLocal && op5 == bytecode.OpLocal {
				s := bytecode.ReadU16(code, i+1)
				j := bytecode.ReadU16(code, i+4)
				ptr := bytecode.ReadU16(code, i+8)
				ptrGet := bytecode.ReadU16(code, i+11)

				if ptr == ptrGet {
					op6 := bytecode.OpCode(code[i+13])
					ptrEscapes := localUsedOutside(code, ptr, i, i+17)

					if !ptrEscapes && op6 == bytecode.OpDeref && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+14])
						if op7 == bytecode.OpSetLocal {
							v := bytecode.ReadU16(code, i+15)
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) &&
								safeIdx(int(v), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceGet)
								bytecode.WriteU16(newInstr, 1, s)
								bytecode.WriteU16(newInstr, 3, j)
								bytecode.WriteU16(newInstr, 5, v)
								rewrites = append(rewrites, rewrite{i, i + 17, newInstr})
								i += 17
								continue
							}
						}
					}

					if !ptrEscapes && op6 == bytecode.OpLocal && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+16])
						if op7 == bytecode.OpSetDeref {
							val := bytecode.ReadU16(code, i+14)
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) &&
								safeIdx(int(val), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceSet)
								bytecode.WriteU16(newInstr, 1, s)
								bytecode.WriteU16(newInstr, 3, j)
								bytecode.WriteU16(newInstr, 5, val)
								rewrites = append(rewrites, rewrite{i, i + 17, newInstr})
								i += 17
								continue
							}
						}
					}

					if !ptrEscapes && op6 == bytecode.OpConst && i+17 <= len(code) {
						op7 := bytecode.OpCode(code[i+16])
						if op7 == bytecode.OpSetDeref {
							constIdx := bytecode.ReadU16(code, i+14)
							if safeIdx(int(s), localIsIntSlice) &&
								safeIdx(int(j), localIsInt) {
								newInstr := make([]byte, 7)
								newInstr[0] = byte(bytecode.OpIntSliceSetConst)
								bytecode.WriteU16(newInstr, 1, s)
								bytecode.WriteU16(newInstr, 3, j)
								bytecode.WriteU16(newInstr, 5, constIdx)
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

// FuseIntMoves replaces OpIntLocal(A) OpIntSetLocal(B) pairs with OpIntMoveLocal(A, B).
func FuseIntMoves(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		width := opcodeWidth(op)
		instrEnd := i + 1 + width
		if instrEnd > len(code) {
			break
		}

		if op == bytecode.OpIntLocal && i+6 <= len(code) {
			op2 := bytecode.OpCode(code[i+3])
			if op2 == bytecode.OpIntSetLocal {
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
		}

		i = instrEnd
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}

// --- Internal helpers ---

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
		case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
			oldTarget := int(bytecode.ReadU16(code, i+1))
			if oldTarget < oldLen {
				newTarget := offsetMap[oldTarget]
				bytecode.WriteU16(code, i+1, uint16(newTarget))
			}

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
				newTarget := offsetMap[oldTarget]
				bytecode.WriteU16(code, i+5, uint16(newTarget))
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
			idx := bytecode.ReadU16(code, i+1)
			if idx == localIdx {
				return true
			}
		}
		i = instrEnd
	}
	return false
}
