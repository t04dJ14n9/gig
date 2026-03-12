package peephole

import "github.com/t04dJ14n9/gig/bytecode"

// cmpSetJumpPattern fuses the 6-instruction sequence (16 bytes → 7 bytes):
// LOCAL(A) op2(B) cmpOp SETLOCAL(X) LOCAL(X) jumpOp(off)
type cmpSetJumpPattern struct {
	op2    bytecode.OpCode // OpConst or OpLocal
	cmpOp  bytecode.OpCode // OpLess, OpLessEq, OpGreater
	jumpOp bytecode.OpCode // OpJumpTrue or OpJumpFalse
	fused  bytecode.OpCode
}

func init() {
	Register(
		// LOCAL(A) CONST(B) LESSEQ  SETLOCAL(X) LOCAL(X) JUMPFALSE(off)
		cmpSetJumpPattern{bytecode.OpConst, bytecode.OpLessEq, bytecode.OpJumpFalse, bytecode.OpLessEqLocalConstJumpFalse},
		// LOCAL(A) CONST(B) LESSEQ  SETLOCAL(X) LOCAL(X) JUMPTRUE(off)
		cmpSetJumpPattern{bytecode.OpConst, bytecode.OpLessEq, bytecode.OpJumpTrue, bytecode.OpLessEqLocalConstJumpTrue},
		// LOCAL(A) CONST(B) LESS    SETLOCAL(X) LOCAL(X) JUMPFALSE(off)
		cmpSetJumpPattern{bytecode.OpConst, bytecode.OpLess, bytecode.OpJumpFalse, bytecode.OpLessLocalConstJumpFalse},
		// LOCAL(A) CONST(B) LESS    SETLOCAL(X) LOCAL(X) JUMPTRUE(off)
		cmpSetJumpPattern{bytecode.OpConst, bytecode.OpLess, bytecode.OpJumpTrue, bytecode.OpLessLocalConstJumpTrue},
		// LOCAL(A) LOCAL(B) LESS    SETLOCAL(X) LOCAL(X) JUMPFALSE(off)
		cmpSetJumpPattern{bytecode.OpLocal, bytecode.OpLess, bytecode.OpJumpFalse, bytecode.OpLessLocalLocalJumpFalse},
		// LOCAL(A) LOCAL(B) LESS    SETLOCAL(X) LOCAL(X) JUMPTRUE(off)
		cmpSetJumpPattern{bytecode.OpLocal, bytecode.OpLess, bytecode.OpJumpTrue, bytecode.OpLessLocalLocalJumpTrue},
		// LOCAL(A) LOCAL(B) GREATER SETLOCAL(X) LOCAL(X) JUMPTRUE(off)
		cmpSetJumpPattern{bytecode.OpLocal, bytecode.OpGreater, bytecode.OpJumpTrue, bytecode.OpGreaterLocalLocalJumpTrue},
	)
}

func (p cmpSetJumpPattern) Match(code []byte, i int) (int, []byte, bool) {
	const size = 16
	if !MatchOp(code, i, bytecode.OpLocal) || i+size > len(code) {
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+7]) != bytecode.OpSetLocal ||
		bytecode.OpCode(code[i+10]) != bytecode.OpLocal {
		return 0, nil, false
	}
	if ReadU16(code, i+8) != ReadU16(code, i+11) { // setIdx == getIdx
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+3]) != p.op2 ||
		bytecode.OpCode(code[i+6]) != p.cmpOp ||
		bytecode.OpCode(code[i+13]) != p.jumpOp {
		return 0, nil, false
	}
	a := ReadU16(code, i+1)
	b := ReadU16(code, i+4)
	off := ReadU16(code, i+14)
	return size, Make3Op(p.fused, a, b, off), true
}
