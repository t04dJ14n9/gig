package peephole

import "git.woa.com/youngjin/gig/bytecode"

// arithSetLocalPattern fuses the 4-instruction sequence (10 bytes → 7 bytes):
// LOCAL(A) op2(B) arithOp SETLOCAL(C)
type arithSetLocalPattern struct {
	op2   bytecode.OpCode // OpLocal or OpConst
	arith bytecode.OpCode // OpAdd, OpSub, OpMul
	fused bytecode.OpCode
}

func init() {
	Register(
		arithSetLocalPattern{bytecode.OpLocal, bytecode.OpAdd, bytecode.OpLocalLocalAddSetLocal},
		arithSetLocalPattern{bytecode.OpConst, bytecode.OpAdd, bytecode.OpLocalConstAddSetLocal},
		arithSetLocalPattern{bytecode.OpConst, bytecode.OpSub, bytecode.OpLocalConstSubSetLocal},
		arithSetLocalPattern{bytecode.OpLocal, bytecode.OpSub, bytecode.OpLocalLocalSubSetLocal},
		arithSetLocalPattern{bytecode.OpLocal, bytecode.OpMul, bytecode.OpLocalLocalMulSetLocal},
		arithSetLocalPattern{bytecode.OpConst, bytecode.OpMul, bytecode.OpLocalConstMulSetLocal},
	)
}

func (p arithSetLocalPattern) Match(code []byte, i int) (int, []byte, bool) {
	const size = 10
	if !MatchOp(code, i, bytecode.OpLocal) || i+size > len(code) {
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+3]) != p.op2 ||
		bytecode.OpCode(code[i+6]) != p.arith ||
		bytecode.OpCode(code[i+7]) != bytecode.OpSetLocal {
		return 0, nil, false
	}
	a := ReadU16(code, i+1)
	b := ReadU16(code, i+4)
	c := ReadU16(code, i+8)
	return size, Make3Op(p.fused, a, b, c), true
}
