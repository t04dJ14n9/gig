package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

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
			if oldTarget <= oldLen {
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
			if oldTarget <= oldLen {
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
