package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

func foldConstantBranches(code []byte, constants []any) ([]byte, bool) {
	targets := jumpTargetSet(code)
	var rewrites []rewrite

	for i := 0; i < len(code); {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}
		if r, ok := constantBranchRewriteAt(code, i, constants, targets); ok {
			rewrites = append(rewrites, r)
			i = r.oldEnd
			continue
		}
		i = instrEnd
	}
	if len(rewrites) == 0 {
		return code, false
	}
	return applyRewrites(code, rewrites), true
}

func constantBranchRewriteAt(code []byte, i int, constants []any, targets map[int]bool) (rewrite, bool) {
	cond, jumpStart, ok := constantBoolAt(code, i, constants)
	if !ok || jumpStart >= len(code) || targets[jumpStart] {
		return rewrite{}, false
	}
	jumpOp := bytecode.OpCode(code[jumpStart])
	jumpEnd := jumpStart + 1 + opcodeWidth(jumpOp)
	if jumpEnd > len(code) || !isConditionalJump(jumpOp) {
		return rewrite{}, false
	}
	if branchTaken(cond, jumpOp) {
		target := bytecode.ReadU16(code, jumpStart+1)
		return rewrite{i, jumpEnd, jumpBytes(target)}, true
	}
	return rewrite{i, jumpEnd, nil}, true
}

func constantBoolAt(code []byte, i int, constants []any) (value bool, next int, ok bool) {
	switch bytecode.OpCode(code[i]) {
	case bytecode.OpTrue:
		return true, i + 1, true
	case bytecode.OpFalse:
		return false, i + 1, true
	case bytecode.OpConst:
		idx := bytecode.ReadU16(code, i+1)
		value, ok := constantAt(constants, idx).(bool)
		return value, i + 3, ok
	default:
		return false, 0, false
	}
}

func isConditionalJump(op bytecode.OpCode) bool {
	return op == bytecode.OpJumpTrue || op == bytecode.OpJumpFalse
}

func branchTaken(cond bool, jumpOp bytecode.OpCode) bool {
	return (jumpOp == bytecode.OpJumpTrue && cond) || (jumpOp == bytecode.OpJumpFalse && !cond)
}

func jumpBytes(target uint16) []byte {
	out := []byte{byte(bytecode.OpJump), 0, 0}
	bytecode.WriteU16(out, 1, target)
	return out
}
