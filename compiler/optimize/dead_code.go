package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

func removeUnreachableAfterJumps(code []byte) ([]byte, bool) {
	targets := jumpTargetSet(code)
	var rewrites []rewrite

	for i := 0; i < len(code); {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}
		if op == bytecode.OpJump {
			if r, ok := unreachableAfterJumpRewrite(instrEnd, len(code), targets); ok {
				rewrites = append(rewrites, r)
				i = r.oldEnd
				continue
			}
		}
		i = instrEnd
	}
	if len(rewrites) == 0 {
		return code, false
	}
	return applyRewrites(code, rewrites), true
}

func unreachableAfterJumpRewrite(instrEnd, codeLen int, targets map[int]bool) (rewrite, bool) {
	deadEnd := nextJumpTargetAtOrAfter(targets, instrEnd, codeLen)
	if instrEnd >= deadEnd {
		return rewrite{}, false
	}
	return rewrite{oldStart: instrEnd, oldEnd: deadEnd, newBytes: nil}, true
}

func nextJumpTargetAtOrAfter(targets map[int]bool, start, fallback int) int {
	next := fallback
	for target := range targets {
		if target >= start && target < next {
			next = target
		}
	}
	return next
}
