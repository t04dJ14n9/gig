package peephole

import "git.woa.com/youngjin/gig/bytecode"

// arithLocalPattern fuses the 2-instruction sequence (7 bytes → 5 bytes):
// LOCAL(A) op2(B) arithOp
type arithLocalPattern struct {
	op2   bytecode.OpCode
	arith bytecode.OpCode
	fused bytecode.OpCode
}

func init() {
	Register(
		arithLocalPattern{bytecode.OpLocal, bytecode.OpAdd, bytecode.OpAddLocalLocal},
		arithLocalPattern{bytecode.OpLocal, bytecode.OpSub, bytecode.OpSubLocalLocal},
		arithLocalPattern{bytecode.OpLocal, bytecode.OpMul, bytecode.OpMulLocalLocal},
		arithLocalPattern{bytecode.OpConst, bytecode.OpAdd, bytecode.OpAddLocalConst},
		arithLocalPattern{bytecode.OpConst, bytecode.OpSub, bytecode.OpSubLocalConst},
	)
}

func (p arithLocalPattern) Match(code []byte, i int) (int, []byte, bool) {
	const size = 7
	if !MatchOp(code, i, bytecode.OpLocal) || i+size > len(code) {
		return 0, nil, false
	}
	if bytecode.OpCode(code[i+3]) != p.op2 ||
		bytecode.OpCode(code[i+6]) != p.arith {
		return 0, nil, false
	}
	a := ReadU16(code, i+1)
	b := ReadU16(code, i+4)
	return size, Make2Op(p.fused, a, b), true
}

func init() {
	Register(
		arithLocalPattern{bytecode.OpLocal, bytecode.OpAdd, bytecode.OpAddLocalLocal},
		arithLocalPattern{bytecode.OpLocal, bytecode.OpSub, bytecode.OpSubLocalLocal},
		arithLocalPattern{bytecode.OpLocal, bytecode.OpMul, bytecode.OpMulLocalLocal},
		arithLocalPattern{bytecode.OpConst, bytecode.OpAdd, bytecode.OpAddLocalConst},
		arithLocalPattern{bytecode.OpConst, bytecode.OpSub, bytecode.OpSubLocalConst},
	)
}
