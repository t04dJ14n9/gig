package peephole

import "git.woa.com/youngjin/gig/bytecode"

// deadJumpPattern eliminates OpJump(off) where off == next instruction.
type deadJumpPattern struct{}

func (deadJumpPattern) Match(code []byte, i int) (int, []byte, bool) {
	if bytecode.OpCode(code[i]) != bytecode.OpJump {
		return 0, nil, false
	}
	width := bytecode.OperandWidth(bytecode.OpJump)
	instrEnd := i + 1 + width
	if instrEnd > len(code) {
		return 0, nil, false
	}
	target := ReadU16(code, i+1)
	if int(target) == instrEnd {
		return instrEnd - i, nil, true // nil newBytes → delete
	}
	return 0, nil, false
}

func init() {
	Register(deadJumpPattern{})
}
