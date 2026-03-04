package peephole

import "git.woa.com/youngjin/gig/bytecode"

// cmpJumpPattern fuses the 3-instruction sequence (10 bytes → 7 bytes):
// LOCAL(A) op2(B) cmpOp jumpOp(off)
type cmpJumpPattern struct {
	op2    bytecode.OpCode
	cmpOp  bytecode.OpCode
	jumpOp bytecode.OpCode
	fused  bytecode.OpCode
}

func (p cmpJumpPattern) Match(code []byte, i int) (int, []byte, bool) {
	const size = 10
	if !MatchOp(code, i, bytecode.OpLocal) || i+size > len(code) {
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+3]) != p.op2 ||
		bytecode.OpCode(code[i+6]) != p.cmpOp ||
		bytecode.OpCode(code[i+7]) != p.jumpOp {
		return 0, nil, false
	}
	a := ReadU16(code, i+1)
	b := ReadU16(code, i+4)
	off := ReadU16(code, i+8)
	return size, Make3Op(p.fused, a, b, off), true
}

func init() {
	Register(
		cmpJumpPattern{bytecode.OpLocal, bytecode.OpLess, bytecode.OpJumpTrue, bytecode.OpLessLocalLocalJumpTrue},
		cmpJumpPattern{bytecode.OpLocal, bytecode.OpLess, bytecode.OpJumpFalse, bytecode.OpLessLocalLocalJumpFalse},
		cmpJumpPattern{bytecode.OpConst, bytecode.OpLess, bytecode.OpJumpTrue, bytecode.OpLessLocalConstJumpTrue},
		cmpJumpPattern{bytecode.OpConst, bytecode.OpLess, bytecode.OpJumpFalse, bytecode.OpLessLocalConstJumpFalse},
		cmpJumpPattern{bytecode.OpConst, bytecode.OpLessEq, bytecode.OpJumpTrue, bytecode.OpLessEqLocalConstJumpTrue},
		cmpJumpPattern{bytecode.OpLocal, bytecode.OpGreater, bytecode.OpJumpTrue, bytecode.OpGreaterLocalLocalJumpTrue},
	)
}
