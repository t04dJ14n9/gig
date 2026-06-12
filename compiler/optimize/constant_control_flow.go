package optimize

import "github.com/t04dJ14n9/gig/model/bytecode"

func jumpTargetSet(code []byte) map[int]bool {
	targets := make(map[int]bool)
	for i := 0; i < len(code); {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}
		switch op {
		case bytecode.OpJump, bytecode.OpJumpTrue, bytecode.OpJumpFalse:
			targets[int(bytecode.ReadU16(code, i+1))] = true
		}
		i = instrEnd
	}
	return targets
}
