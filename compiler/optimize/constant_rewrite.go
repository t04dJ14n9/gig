package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

func foldConstantStackOps(code []byte, constants *[]any) ([]byte, bool) {
	targets := jumpTargetSet(code)
	var rewrites []rewrite
	var starts []int

	for i := 0; i < len(code); {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}
		if isFoldableConstBinaryOp(op) && len(starts) >= 2 {
			if r, ok := foldedBinaryRewrite(code, starts[len(starts)-2], starts[len(starts)-1], i, instrEnd, op, targets, constants); ok {
				rewrites = append(rewrites, r)
			}
		}
		starts = append(starts, i)
		i = instrEnd
	}
	if len(rewrites) == 0 {
		return code, false
	}
	return applyRewrites(code, rewrites), true
}

func foldedBinaryRewrite(
	code []byte,
	leftStart, rightStart, opStart, opEnd int,
	op bytecode.OpCode,
	targets map[int]bool,
	constants *[]any,
) (rewrite, bool) {
	if !isConsecutiveConstPair(code, leftStart, rightStart, opStart) || hasJumpTargetInside(targets, leftStart, opEnd) {
		return rewrite{}, false
	}
	leftIdx := bytecode.ReadU16(code, leftStart+1)
	rightIdx := bytecode.ReadU16(code, rightStart+1)
	newBytes, ok := foldedBinaryBytes(op, constantAt(*constants, leftIdx), constantAt(*constants, rightIdx), constants)
	if !ok {
		return rewrite{}, false
	}
	return rewrite{oldStart: leftStart, oldEnd: opEnd, newBytes: newBytes}, true
}

func isConsecutiveConstPair(code []byte, leftStart, rightStart, opStart int) bool {
	return bytecode.OpCode(code[leftStart]) == bytecode.OpConst &&
		bytecode.OpCode(code[rightStart]) == bytecode.OpConst &&
		leftStart+3 == rightStart &&
		rightStart+3 == opStart
}

func hasJumpTargetInside(targets map[int]bool, start, end int) bool {
	for pos := start + 1; pos < end; pos++ {
		if targets[pos] {
			return true
		}
	}
	return false
}

func constBytes(idx uint16) []byte {
	out := []byte{byte(bytecode.OpConst), 0, 0}
	bytecode.WriteU16(out, 1, idx)
	return out
}

func boolBytes(v bool) []byte {
	if v {
		return []byte{byte(bytecode.OpTrue)}
	}
	return []byte{byte(bytecode.OpFalse)}
}

func appendConstant(constants *[]any, value any) uint16 {
	*constants = append(*constants, value)
	return uint16(len(*constants) - 1)
}

func constantAt(constants []any, idx uint16) any {
	if int(idx) >= len(constants) {
		return nil
	}
	return constants[idx]
}
